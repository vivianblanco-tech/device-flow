package handlers

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/yourusername/laptop-tracking-system/internal/middleware"
	"github.com/yourusername/laptop-tracking-system/internal/models"
	"github.com/yourusername/laptop-tracking-system/internal/views"
)

// LaptopReceptionReportPage displays the reception report form for a specific laptop
func (h *ReceptionReportHandler) LaptopReceptionReportPage(w http.ResponseWriter, r *http.Request, laptopID int64) {
	// Get user from context
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Only warehouse users can submit reception reports
	if user.Role != models.RoleWarehouse {
		http.Error(w, "Forbidden: Only warehouse users can create reception reports", http.StatusForbidden)
		return
	}

	// Get laptop details
	var laptop models.Laptop
	var clientCompanyName sql.NullString
	var trackingNumber sql.NullString
	var shipmentID sql.NullInt64
	
	err := h.DB.QueryRowContext(r.Context(),
		`SELECT 
			l.id, l.serial_number, l.brand, l.model, l.ram_gb, l.ssd_gb, l.status,
			l.client_company_id, l.software_engineer_id, l.created_at, l.updated_at,
			cc.name as client_company_name,
			s.id as shipment_id,
			s.tracking_number
		FROM laptops l
		LEFT JOIN client_companies cc ON cc.id = l.client_company_id
		LEFT JOIN shipment_laptops sl ON sl.laptop_id = l.id
		LEFT JOIN shipments s ON s.id = sl.shipment_id AND s.status != 'delivered'
		WHERE l.id = $1
		LIMIT 1`,
		laptopID,
	).Scan(
		&laptop.ID, &laptop.SerialNumber, &laptop.Brand, &laptop.Model,
		&laptop.RAMGB, &laptop.SSDGB, &laptop.Status,
		&laptop.ClientCompanyID, &laptop.SoftwareEngineerID,
		&laptop.CreatedAt, &laptop.UpdatedAt,
		&clientCompanyName, &shipmentID, &trackingNumber,
	)

	if err == sql.ErrNoRows {
		http.Error(w, "Laptop not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed to load laptop details", http.StatusInternalServerError)
		return
	}

	// Check if laptop status allows reception report
	if laptop.Status != models.LaptopStatusAtWarehouse {
		http.Error(w, "Reception report can only be created for laptops with 'Received at Warehouse' status", http.StatusBadRequest)
		return
	}

	// Check if reception report already exists
	existingReport, err := models.GetLaptopReceptionReport(r.Context(), h.DB, laptopID)
	if err != nil {
		http.Error(w, "Failed to check existing reception report", http.StatusInternalServerError)
		return
	}
	if existingReport != nil {
		// Redirect to view the existing report
		http.Redirect(w, r, fmt.Sprintf("/reception-reports/%d", existingReport.ID), http.StatusSeeOther)
		return
	}

	// Get error and success messages from query parameters
	errorMsg := r.URL.Query().Get("error")
	successMsg := r.URL.Query().Get("success")

	data := map[string]interface{}{
		"Error":             errorMsg,
		"Success":           successMsg,
		"User":              user,
		"Nav":               views.GetNavigationLinks(user.Role),
		"CurrentPage":       "reception-reports",
		"Laptop":            laptop,
		"ClientCompanyName": clientCompanyName.String,
		"ShipmentID":        shipmentID,
		"TrackingNumber":    trackingNumber.String,
	}

	if h.Templates != nil {
		err := h.Templates.ExecuteTemplate(w, "laptop-reception-report-form.html", data)
		if err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			return
		}
	} else {
		// For testing
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Laptop Reception Report Form")
	}
}

// LaptopReceptionReportSubmit handles the reception report submission for a specific laptop
func (h *ReceptionReportHandler) LaptopReceptionReportSubmit(w http.ResponseWriter, r *http.Request, laptopID int64) {
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

	// Parse multipart form
	err := r.ParseMultipartForm(MaxUploadSize)
	if err != nil {
		redirectURL := fmt.Sprintf("/laptops/%d/reception-report?error=Invalid+form+data", laptopID)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}

	// Extract form values
	notes := r.FormValue("notes")

	// Handle required photo uploads
	photoSerialNumber, err := h.handleSinglePhotoUpload(r, "photo_serial_number", laptopID, "serial")
	if err != nil {
		redirectURL := fmt.Sprintf("/laptops/%d/reception-report?error=Serial+number+photo+required", laptopID)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}

	photoExternalCondition, err := h.handleSinglePhotoUpload(r, "photo_external_condition", laptopID, "external")
	if err != nil {
		// Clean up previously uploaded photo
		os.Remove("." + photoSerialNumber)
		redirectURL := fmt.Sprintf("/laptops/%d/reception-report?error=External+condition+photo+required", laptopID)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}

	photoWorkingCondition, err := h.handleSinglePhotoUpload(r, "photo_working_condition", laptopID, "working")
	if err != nil {
		// Clean up previously uploaded photos
		os.Remove("." + photoSerialNumber)
		os.Remove("." + photoExternalCondition)
		redirectURL := fmt.Sprintf("/laptops/%d/reception-report?error=Working+condition+photo+required", laptopID)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}

	// Get laptop and shipment details
	var shipmentID sql.NullInt64
	var clientCompanyID sql.NullInt64
	var trackingNumber sql.NullString
	
	err = h.DB.QueryRowContext(r.Context(),
		`SELECT 
			s.id as shipment_id,
			s.client_company_id,
			s.tracking_number
		FROM laptops l
		LEFT JOIN shipment_laptops sl ON sl.laptop_id = l.id
		LEFT JOIN shipments s ON s.id = sl.shipment_id
		WHERE l.id = $1 AND s.status != 'delivered'
		LIMIT 1`,
		laptopID,
	).Scan(&shipmentID, &clientCompanyID, &trackingNumber)

	// It's okay if there's no shipment (could be ErrNoRows)
	if err != nil && err != sql.ErrNoRows {
		// Clean up uploaded photos on error
		os.Remove("." + photoSerialNumber)
		os.Remove("." + photoExternalCondition)
		os.Remove("." + photoWorkingCondition)
		http.Error(w, "Failed to retrieve laptop details", http.StatusInternalServerError)
		return
	}

	// Create reception report
	var shipmentIDPtr *int64
	var clientCompanyIDPtr *int64
	if shipmentID.Valid {
		shipmentIDPtr = &shipmentID.Int64
	}
	if clientCompanyID.Valid {
		clientCompanyIDPtr = &clientCompanyID.Int64
	}

	report := &models.ReceptionReport{
		LaptopID:               laptopID,
		ShipmentID:             shipmentIDPtr,
		ClientCompanyID:        clientCompanyIDPtr,
		TrackingNumber:         trackingNumber.String,
		WarehouseUserID:        user.ID,
		Notes:                  notes,
		PhotoSerialNumber:      photoSerialNumber,
		PhotoExternalCondition: photoExternalCondition,
		PhotoWorkingCondition:  photoWorkingCondition,
	}

	// Validate report
	if err := report.Validate(); err != nil {
		// Clean up uploaded photos on validation error
		os.Remove("." + photoSerialNumber)
		os.Remove("." + photoExternalCondition)
		os.Remove("." + photoWorkingCondition)
		redirectURL := fmt.Sprintf("/laptops/%d/reception-report?error=%s", laptopID, err.Error())
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}

	// Save to database
	err = models.CreateReceptionReport(r.Context(), h.DB, report)
	if err != nil {
		// Clean up uploaded photos on database error
		os.Remove("." + photoSerialNumber)
		os.Remove("." + photoExternalCondition)
		os.Remove("." + photoWorkingCondition)
		http.Error(w, "Failed to save reception report", http.StatusInternalServerError)
		return
	}

	// Send email notification to international.logistics@bairesdev.com
	go h.sendReceptionReportNotification(report, user)

	// Redirect to reception report detail page
	redirectURL := fmt.Sprintf("/reception-reports/%d?success=Reception+report+created+successfully", report.ID)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

// handleSinglePhotoUpload handles uploading a single photo field
func (h *ReceptionReportHandler) handleSinglePhotoUpload(r *http.Request, fieldName string, laptopID int64, photoType string) (string, error) {
	file, fileHeader, err := r.FormFile(fieldName)
	if err != nil {
		return "", fmt.Errorf("photo required")
	}
	defer file.Close()

	// Validate file size
	if fileHeader.Size > MaxUploadSize {
		return "", fmt.Errorf("file too large")
	}

	// Generate unique filename
	ext := filepath.Ext(fileHeader.Filename)
	filename := fmt.Sprintf("%d_%s_%d%s", laptopID, photoType, time.Now().UnixNano(), ext)
	filePath := filepath.Join(UploadDir, filename)

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to save file")
	}
	defer dst.Close()

	// Copy file content
	_, err = io.Copy(dst, file)
	if err != nil {
		return "", fmt.Errorf("failed to save file")
	}

	// Return relative URL
	return fmt.Sprintf("/uploads/reception/%s", filename), nil
}

// sendReceptionReportNotification sends email notification when a reception report is created
func (h *ReceptionReportHandler) sendReceptionReportNotification(report *models.ReceptionReport, submitter *models.User) {
	// Email sending is handled asynchronously, errors are logged but don't fail the request
	// For now, just log that we would send an email
	// In production, this would integrate with the email service
	
	var laptop models.Laptop
	err := h.DB.QueryRow(
		`SELECT serial_number, brand, model FROM laptops WHERE id = $1`,
		report.LaptopID,
	).Scan(&laptop.SerialNumber, &laptop.Brand, &laptop.Model)
	
	if err != nil {
		fmt.Printf("Error getting laptop details for notification: %v\n", err)
		return
	}

	fmt.Printf("📧 Email notification would be sent to international.logistics@bairesdev.com\n")
	fmt.Printf("   Subject: New Reception Report - Laptop %s\n", laptop.SerialNumber)
	fmt.Printf("   Submitted by: %s\n", submitter.Email)
	fmt.Printf("   Report ID: %d\n", report.ID)
	
	// TODO: Integrate with actual email service when ready
	// This is a placeholder for the email notification functionality
}

// ApproveReceptionReport approves a reception report (logistics only)
func (h *ReceptionReportHandler) ApproveReceptionReport(w http.ResponseWriter, r *http.Request, reportID int64) {
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

	// Only logistics users can approve reception reports
	if user.Role != models.RoleLogistics {
		http.Error(w, "Forbidden: Only logistics users can approve reception reports", http.StatusForbidden)
		return
	}

	// Approve the report (this also updates laptop status)
	err := models.ApproveReceptionReport(r.Context(), h.DB, reportID, user.ID)
	if err != nil {
		redirectURL := fmt.Sprintf("/reception-reports/%d?error=%s", reportID, err.Error())
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}

	// Redirect back to report detail with success message
	redirectURL := fmt.Sprintf("/reception-reports/%d?success=Reception+report+approved+successfully", reportID)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

// LaptopBasedReceptionReportsList displays a list of all laptop-based reception reports
func (h *ReceptionReportHandler) LaptopBasedReceptionReportsList(w http.ResponseWriter, r *http.Request) {
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
			rr.laptop_id,
			rr.shipment_id,
			rr.client_company_id,
			rr.tracking_number,
			rr.warehouse_user_id,
			rr.received_at,
			rr.notes,
			rr.status,
			rr.approved_by,
			rr.approved_at,
			l.serial_number,
			l.brand,
			l.model,
			l.status as laptop_status,
			cc.name as company_name,
			u.email as warehouse_user_email,
			approver.email as approver_email
		FROM reception_reports rr
		JOIN laptops l ON l.id = rr.laptop_id
		LEFT JOIN client_companies cc ON cc.id = rr.client_company_id
		JOIN users u ON u.id = rr.warehouse_user_id
		LEFT JOIN users approver ON approver.id = rr.approved_by
		ORDER BY rr.received_at DESC
	`

	rows, err := h.DB.QueryContext(r.Context(), query)
	if err != nil {
		http.Error(w, "Failed to load reception reports", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type ReceptionReportRow struct {
		ID                  int64
		LaptopID            int64
		ShipmentID          sql.NullInt64
		ClientCompanyID     sql.NullInt64
		TrackingNumber      sql.NullString
		WarehouseUserID     int64
		ReceivedAt          time.Time
		Notes               string
		Status              string
		ApprovedBy          sql.NullInt64
		ApprovedAt          sql.NullTime
		LaptopSerialNumber  string
		LaptopBrand         string
		LaptopModel         string
		LaptopStatus        string
		CompanyName         sql.NullString
		WarehouseUserEmail  string
		ApproverEmail       sql.NullString
	}

	var receptionReports []ReceptionReportRow
	for rows.Next() {
		var row ReceptionReportRow
		err := rows.Scan(
			&row.ID,
			&row.LaptopID,
			&row.ShipmentID,
			&row.ClientCompanyID,
			&row.TrackingNumber,
			&row.WarehouseUserID,
			&row.ReceivedAt,
			&row.Notes,
			&row.Status,
			&row.ApprovedBy,
			&row.ApprovedAt,
			&row.LaptopSerialNumber,
			&row.LaptopBrand,
			&row.LaptopModel,
			&row.LaptopStatus,
			&row.CompanyName,
			&row.WarehouseUserEmail,
			&row.ApproverEmail,
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
		
		err := h.Templates.ExecuteTemplate(w, "laptop-reception-reports-list.html", data)
		if err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			return
		}
	} else {
		// For testing without templates
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Laptop Reception Reports List\n")
	}
}

