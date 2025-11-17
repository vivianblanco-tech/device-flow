package handlers

import (
	"database/sql"
	// "encoding/json" // Not used in laptop-based reception reports
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
		"Nav":         views.GetNavigationLinks(user.Role),
		"CurrentPage": "reception-reports",
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

	// DEPRECATED: Old shipment-based reception report creation
	// This handler is deprecated and kept for backward compatibility only
	// New code should use laptop-based handlers in laptop_reception_report.go
	
	// For now, return an error directing to new system
	http.Error(w, "This endpoint is deprecated. Please use the laptop-based reception report system at /laptops/{id}/reception-report", http.StatusGone)
	return
	
	/*
	// Create reception report
	report := models.ReceptionReport{
		ShipmentID:      &shipmentID,
		WarehouseUserID: user.ID,
		Notes:           notes,
	}
	report.BeforeCreate()

	_, err = tx.ExecContext(r.Context(),
		`INSERT INTO reception_reports (shipment_id, warehouse_user_id, received_at, notes)
		VALUES ($1, $2, $3, $4)`,
		report.ShipmentID, report.WarehouseUserID, report.ReceivedAt,
		report.Notes,
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
	*/
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

	// Get sort parameters
	sortBy := r.URL.Query().Get("sort")
	sortOrder := r.URL.Query().Get("order")

	// Build query to fetch reception reports with related data
	baseQuery := `
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
	`

	// Build ORDER BY clause
	orderBy := buildReceptionReportsOrderByClause(sortBy, sortOrder)
	query := baseQuery + " " + orderBy

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
			"SortBy":           sortBy,
			"SortOrder":        sortOrder,
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

// ReceptionReportDetail displays the details of a specific reception report (laptop-based)
func (h *ReceptionReportHandler) ReceptionReportDetail(w http.ResponseWriter, r *http.Request, reportID int64) {
	// Get user from context
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Only warehouse and logistics users can view reception reports
	if user.Role != models.RoleWarehouse && user.Role != models.RoleLogistics {
		http.Error(w, "Forbidden: Only warehouse and logistics users can access this page", http.StatusForbidden)
		return
	}

	// Fetch reception report with related data using new laptop-based schema
	query := `
		SELECT 
			rr.id,
			rr.laptop_id,
			rr.shipment_id,
			rr.client_company_id,
			rr.tracking_number,
			rr.warehouse_user_id,
			rr.received_at,
			rr.notes,
			rr.photo_serial_number,
			rr.photo_external_condition,
			rr.photo_working_condition,
			rr.status,
			rr.approved_by,
			rr.approved_at,
			rr.created_at,
			rr.updated_at,
			l.serial_number as laptop_serial,
			l.brand as laptop_brand,
			l.model as laptop_model,
			l.status as laptop_status,
			cc.name as company_name,
			u.email as warehouse_user_email,
			approver.email as approver_email
		FROM reception_reports rr
		JOIN laptops l ON l.id = rr.laptop_id
		LEFT JOIN client_companies cc ON cc.id = rr.client_company_id
		JOIN users u ON u.id = rr.warehouse_user_id
		LEFT JOIN users approver ON approver.id = rr.approved_by
		WHERE rr.id = $1
	`

	type ReceptionReportDetail struct {
		ID                     int64
		LaptopID               int64
		ShipmentID             sql.NullInt64
		ClientCompanyID        sql.NullInt64
		TrackingNumber         sql.NullString
		WarehouseUserID        int64
		ReceivedAt             time.Time
		Notes                  string
		PhotoSerialNumber      string
		PhotoExternalCondition string
		PhotoWorkingCondition  string
		Status                 string
		ApprovedBy             sql.NullInt64
		ApprovedAt             sql.NullTime
		CreatedAt              time.Time
		UpdatedAt              time.Time
		LaptopSerial           string
		LaptopBrand            string
		LaptopModel            string
		LaptopStatus           string
		CompanyName            sql.NullString
		WarehouseUserEmail     string
		ApproverEmail          sql.NullString
	}

	var report ReceptionReportDetail
	err := h.DB.QueryRowContext(r.Context(), query, reportID).Scan(
		&report.ID,
		&report.LaptopID,
		&report.ShipmentID,
		&report.ClientCompanyID,
		&report.TrackingNumber,
		&report.WarehouseUserID,
		&report.ReceivedAt,
		&report.Notes,
		&report.PhotoSerialNumber,
		&report.PhotoExternalCondition,
		&report.PhotoWorkingCondition,
		&report.Status,
		&report.ApprovedBy,
		&report.ApprovedAt,
		&report.CreatedAt,
		&report.UpdatedAt,
		&report.LaptopSerial,
		&report.LaptopBrand,
		&report.LaptopModel,
		&report.LaptopStatus,
		&report.CompanyName,
		&report.WarehouseUserEmail,
		&report.ApproverEmail,
	)

	if err == sql.ErrNoRows {
		http.Error(w, "Reception report not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed to load reception report", http.StatusInternalServerError)
		return
	}

	// If templates are available, render the template
	if h.Templates != nil {
		data := map[string]interface{}{
			"User":        user,
			"Nav":         views.GetNavigationLinks(user.Role),
			"CurrentPage": "reception-reports",
			"Report":      report,
		}

		err := h.Templates.ExecuteTemplate(w, "laptop-reception-report-detail.html", data)
		if err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			return
		}
	} else {
		// For testing without templates - output plain text with the data
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Reception Report Detail\n")
		if report.CompanyName.Valid {
			fmt.Fprintf(w, "Company: %s\n", report.CompanyName.String)
		}
		fmt.Fprintf(w, "Notes: %s\n", report.Notes)
		fmt.Fprintf(w, "Warehouse User: %s\n", report.WarehouseUserEmail)
	}
}

// buildReceptionReportsOrderByClause builds the ORDER BY clause for reception reports based on sort parameters
func buildReceptionReportsOrderByClause(sortBy, sortOrder string) string {
	// Map of allowed sort columns to their SQL equivalents
	sortColumns := map[string]string{
		"id":             "rr.id",
		"shipment":       "s.jira_ticket_number",
		"company":        "c.name",
		"type":           "s.shipment_type::text",
		"received_at":    "rr.received_at",
		"warehouse_user": "u.email",
	}

	// Columns that should use COLLATE (text columns only)
	textColumns := map[string]bool{
		"shipment":       true,
		"company":        true,
		"type":           true,
		"warehouse_user": true,
	}

	// Validate sort order
	order := "ASC"
	if sortOrder == "desc" {
		order = "DESC"
	}

	// Default sort: received_at DESC
	if sortBy == "" {
		return "ORDER BY rr.received_at DESC"
	}

	// Get the SQL column name
	sqlColumn, exists := sortColumns[sortBy]
	if !exists {
		// If invalid column, use default
		return "ORDER BY rr.received_at DESC"
	}

	// Only apply COLLATE to text columns
	if textColumns[sortBy] {
		return fmt.Sprintf("ORDER BY %s COLLATE \"C\" %s", sqlColumn, order)
	}

	// For numeric and timestamp columns, don't use COLLATE
	return fmt.Sprintf("ORDER BY %s %s", sqlColumn, order)
}

