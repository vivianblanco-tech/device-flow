package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/lib/pq"

	"github.com/yourusername/laptop-tracking-system/internal/email"
	"github.com/yourusername/laptop-tracking-system/internal/middleware"
	"github.com/yourusername/laptop-tracking-system/internal/models"
	"github.com/yourusername/laptop-tracking-system/internal/validator"
	"github.com/yourusername/laptop-tracking-system/internal/views"
)

const (
	// MaxUploadSize is the maximum file size for uploads (10MB)
	MaxUploadSize = 10 * 1024 * 1024
	// UploadDir is the directory for uploaded files
	UploadDir = "./uploads/reception"
)

// ReceptionReportHandler handles warehouse reception report requests
type ReceptionReportHandler struct {
	DB        *sql.DB
	Templates *template.Template
	Notifier  *email.Notifier
}

// NewReceptionReportHandler creates a new ReceptionReportHandler
func NewReceptionReportHandler(db *sql.DB, templates *template.Template, notifier *email.Notifier) *ReceptionReportHandler {
	// Ensure upload directory exists
	os.MkdirAll(UploadDir, 0755)
	
	return &ReceptionReportHandler{
		DB:        db,
		Templates: templates,
		Notifier:  notifier,
	}
}

// ReceptionReportPage displays the reception report form
func (h *ReceptionReportHandler) ReceptionReportPage(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Only warehouse users can submit reception reports
	if user.Role != models.RoleWarehouse {
		http.Error(w, "Forbidden: Only warehouse users can access this page", http.StatusForbidden)
		return
	}

	// Get shipment ID from URL path
	shipmentIDStr := r.URL.Query().Get("shipment_id")
	if shipmentIDStr == "" {
		http.Error(w, "Shipment ID is required", http.StatusBadRequest)
		return
	}

	shipmentID, err := strconv.ParseInt(shipmentIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid shipment ID", http.StatusBadRequest)
		return
	}

	// Get shipment details
	var shipment models.Shipment
	var companyName string
	err = h.DB.QueryRowContext(r.Context(),
		`SELECT s.id, s.client_company_id, s.status, s.created_at, s.updated_at, c.name
		FROM shipments s
		JOIN client_companies c ON c.id = s.client_company_id
		WHERE s.id = $1`,
		shipmentID,
	).Scan(&shipment.ID, &shipment.ClientCompanyID, &shipment.Status, 
		&shipment.CreatedAt, &shipment.UpdatedAt, &companyName)

	if err == sql.ErrNoRows {
		http.Error(w, "Shipment not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed to load shipment", http.StatusInternalServerError)
		return
	}

	// Get error and success messages from query parameters
	errorMsg := r.URL.Query().Get("error")
	successMsg := r.URL.Query().Get("success")

	data := map[string]interface{}{
		"Error":       errorMsg,
		"Success":     successMsg,
		"User":        user,
		"Shipment":    shipment,
		"CompanyName": companyName,
	}

	if h.Templates != nil {
		err := h.Templates.ExecuteTemplate(w, "reception-report.html", data)
		if err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			return
		}
	} else {
		// For testing without templates
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Reception Report Page")
	}
}

// ReceptionReportSubmit handles the reception report submission
func (h *ReceptionReportHandler) ReceptionReportSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user from context
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Only warehouse users can submit reception reports
	if user.Role != models.RoleWarehouse {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Parse form - try multipart first (for file uploads), fallback to regular form
	err := r.ParseMultipartForm(MaxUploadSize)
	if err != nil {
		// If multipart parsing fails, try regular form parsing
		err = r.ParseForm()
		if err != nil {
			http.Redirect(w, r, "/reception-report?error=Invalid+form+data", http.StatusSeeOther)
			return
		}
	}

	// Extract form values
	shipmentIDStr := r.FormValue("shipment_id")
	shipmentID, err := strconv.ParseInt(shipmentIDStr, 10, 64)
	if err != nil {
		http.Redirect(w, r, "/reception-report?error=Invalid+shipment+ID", http.StatusSeeOther)
		return
	}

	notes := r.FormValue("notes")

	// Handle photo uploads
	photoURLs := []string{}
	var files []*multipart.FileHeader
	if r.MultipartForm != nil {
		files = r.MultipartForm.File["photos"]
	}
	
	for _, fileHeader := range files {
		// Validate file size
		if fileHeader.Size > MaxUploadSize {
			redirectURL := fmt.Sprintf("/reception-report?shipment_id=%d&error=File+too+large", shipmentID)
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return
		}

		// Open uploaded file
		file, err := fileHeader.Open()
		if err != nil {
			redirectURL := fmt.Sprintf("/reception-report?shipment_id=%d&error=Failed+to+read+file", shipmentID)
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return
		}
		defer file.Close()

		// Generate unique filename
		ext := filepath.Ext(fileHeader.Filename)
		filename := fmt.Sprintf("%d_%d%s", shipmentID, time.Now().UnixNano(), ext)
		filepath := filepath.Join(UploadDir, filename)

		// Create destination file
		dst, err := os.Create(filepath)
		if err != nil {
			redirectURL := fmt.Sprintf("/reception-report?shipment_id=%d&error=Failed+to+save+file", shipmentID)
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return
		}
		defer dst.Close()

		// Copy file content
		_, err = io.Copy(dst, file)
		if err != nil {
			redirectURL := fmt.Sprintf("/reception-report?shipment_id=%d&error=Failed+to+save+file", shipmentID)
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return
		}

		// Store relative URL
		photoURL := fmt.Sprintf("/uploads/reception/%s", filename)
		photoURLs = append(photoURLs, photoURL)
	}

	// Build validation input
	reportInput := validator.ReceptionReportInput{
		ShipmentID:      shipmentID,
		WarehouseUserID: user.ID,
		Notes:           notes,
		PhotoURLs:       photoURLs,
	}

	// Validate report
	if err := validator.ValidateReceptionReport(reportInput); err != nil {
		// Clean up uploaded files on validation error
		for _, photoURL := range photoURLs {
			filepath := "." + photoURL
			os.Remove(filepath)
		}

		redirectURL := fmt.Sprintf("/reception-report?shipment_id=%d&error=%s", 
			shipmentID, err.Error())
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}

	// Start transaction
	tx, err := h.DB.BeginTx(r.Context(), nil)
	if err != nil {
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Create reception report
	report := models.ReceptionReport{
		ShipmentID:      shipmentID,
		WarehouseUserID: user.ID,
		Notes:           notes,
		PhotoURLs:       photoURLs,
	}
	report.BeforeCreate()

	_, err = tx.ExecContext(r.Context(),
		`INSERT INTO reception_reports (shipment_id, warehouse_user_id, received_at, notes, photo_urls)
		VALUES ($1, $2, $3, $4, $5)`,
		report.ShipmentID, report.WarehouseUserID, report.ReceivedAt,
		report.Notes, pq.Array(report.PhotoURLs),
	)
	if err != nil {
		http.Error(w, "Failed to save reception report", http.StatusInternalServerError)
		return
	}

	// Update shipment status to "at_warehouse"
	now := time.Now()
	_, err = tx.ExecContext(r.Context(),
		`UPDATE shipments 
		SET status = $1, arrived_warehouse_at = $2, updated_at = $3
		WHERE id = $4`,
		models.ShipmentStatusAtWarehouse, now, now, shipmentID,
	)
	if err != nil {
		http.Error(w, "Failed to update shipment status", http.StatusInternalServerError)
		return
	}

	// Create audit log entry
	auditDetails, _ := json.Marshal(map[string]interface{}{
		"action":      "reception_report_submitted",
		"shipment_id": shipmentID,
		"photo_count": len(photoURLs),
	})

	_, err = tx.ExecContext(r.Context(),
		`INSERT INTO audit_logs (user_id, action, entity_type, entity_id, timestamp, details)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		user.ID, "reception_report_submitted", "shipment", shipmentID, time.Now(), auditDetails,
	)
	if err != nil {
		// Non-critical error, just log it
		fmt.Printf("Failed to create audit log: %v\n", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	// TODO: Send notification to logistics about warehouse reception (Step 8 in process flow)
	// This would require adding a new method to the Notifier like SendWarehouseReceptionNotification
	// For now, logging the event via audit_logs is sufficient

	// Redirect to success page or shipment detail
	redirectURL := fmt.Sprintf("/shipments/%d?success=Reception+report+submitted+successfully", shipmentID)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

// ReceptionReportsList displays a list of all reception reports
func (h *ReceptionReportHandler) ReceptionReportsList(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Only warehouse and logistics users can view reception reports list
	if user.Role != models.RoleWarehouse && user.Role != models.RoleLogistics {
		http.Error(w, "Forbidden: Only warehouse and logistics users can access this page", http.StatusForbidden)
		return
	}

	// Build query to fetch reception reports with related data
	query := `
		SELECT 
			rr.id,
			rr.shipment_id,
			rr.warehouse_user_id,
			rr.received_at,
			rr.notes,
			rr.photo_urls,
			rr.expected_serial_number,
			rr.actual_serial_number,
			rr.serial_number_corrected,
			rr.correction_note,
			rr.correction_approved_by,
			s.jira_ticket_number,
			s.shipment_type,
			s.status as shipment_status,
			c.name as company_name,
			u.email as warehouse_user_email
		FROM reception_reports rr
		JOIN shipments s ON s.id = rr.shipment_id
		JOIN client_companies c ON c.id = s.client_company_id
		JOIN users u ON u.id = rr.warehouse_user_id
		ORDER BY rr.received_at DESC
	`

	rows, err := h.DB.QueryContext(r.Context(), query)
	if err != nil {
		http.Error(w, "Failed to load reception reports", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type ReceptionReportRow struct {
		ID                     int64
		ShipmentID             int64
		WarehouseUserID        int64
		ReceivedAt             time.Time
		Notes                  string
		PhotoURLs              []string
		ExpectedSerialNumber   sql.NullString
		ActualSerialNumber     sql.NullString
		SerialNumberCorrected  bool
		CorrectionNote         sql.NullString
		CorrectionApprovedBy   sql.NullInt64
		JiraTicketNumber       string
		ShipmentType           string
		ShipmentStatus         string
		CompanyName            string
		WarehouseUserEmail     string
	}

	var receptionReports []ReceptionReportRow
	for rows.Next() {
		var row ReceptionReportRow
		err := rows.Scan(
			&row.ID,
			&row.ShipmentID,
			&row.WarehouseUserID,
			&row.ReceivedAt,
			&row.Notes,
			(*pq.StringArray)(&row.PhotoURLs),
			&row.ExpectedSerialNumber,
			&row.ActualSerialNumber,
			&row.SerialNumberCorrected,
			&row.CorrectionNote,
			&row.CorrectionApprovedBy,
			&row.JiraTicketNumber,
			&row.ShipmentType,
			&row.ShipmentStatus,
			&row.CompanyName,
			&row.WarehouseUserEmail,
		)
		if err != nil {
			http.Error(w, "Failed to parse reception reports", http.StatusInternalServerError)
			return
		}
		receptionReports = append(receptionReports, row)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, "Failed to load reception reports", http.StatusInternalServerError)
		return
	}

	// If templates are available, render the template
	if h.Templates != nil {
		data := map[string]interface{}{
			"User":             user,
			"Nav":              views.GetNavigationLinks(user.Role),
			"CurrentPage":      "reception-reports",
			"ReceptionReports": receptionReports,
		}
		
		err := h.Templates.ExecuteTemplate(w, "reception-reports-list.html", data)
		if err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			return
		}
	} else {
		// For testing without templates - output plain text with the data
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Reception Reports List\n")
		for _, rr := range receptionReports {
			fmt.Fprintf(w, "Company: %s, User: %s, Notes: %s\n", 
				rr.CompanyName, rr.WarehouseUserEmail, rr.Notes)
		}
	}
}

