package handlers

import (
	"database/sql"
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

	"github.com/yourusername/laptop-tracking-system/internal/auth"
	"github.com/yourusername/laptop-tracking-system/internal/email"
	"github.com/yourusername/laptop-tracking-system/internal/middleware"
	"github.com/yourusername/laptop-tracking-system/internal/models"
	"github.com/yourusername/laptop-tracking-system/internal/validator"
	"github.com/yourusername/laptop-tracking-system/internal/views"
)

const (
	// DeliveryUploadDir is the directory for delivery photo uploads
	DeliveryUploadDir = "./uploads/delivery"
)

// DeliveryFormHandler handles delivery form requests
type DeliveryFormHandler struct {
	DB        *sql.DB
	Templates *template.Template
	Notifier  *email.Notifier
}

// NewDeliveryFormHandler creates a new DeliveryFormHandler
func NewDeliveryFormHandler(db *sql.DB, templates *template.Template, notifier *email.Notifier) *DeliveryFormHandler {
	// Ensure upload directory exists
	os.MkdirAll(DeliveryUploadDir, 0755)
	
	return &DeliveryFormHandler{
		DB:        db,
		Templates: templates,
		Notifier:  notifier,
	}
}

// DeliveryFormPage displays the delivery form
func (h *DeliveryFormHandler) DeliveryFormPage(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Get shipment ID from URL query parameter
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

	// Get shipment details with engineer information
	var shipment models.Shipment
	var companyName string
	var engineerName sql.NullString
	var engineerID sql.NullInt64
	
	err = h.DB.QueryRowContext(r.Context(),
		`SELECT s.id, s.client_company_id, s.software_engineer_id, s.status, 
		        s.created_at, s.updated_at, c.name, se.name
		FROM shipments s
		JOIN client_companies c ON c.id = s.client_company_id
		LEFT JOIN software_engineers se ON se.id = s.software_engineer_id
		WHERE s.id = $1`,
		shipmentID,
	).Scan(&shipment.ID, &shipment.ClientCompanyID, &shipment.SoftwareEngineerID,
		&shipment.Status, &shipment.CreatedAt, &shipment.UpdatedAt, 
		&companyName, &engineerName)

	if err == sql.ErrNoRows {
		http.Error(w, "Shipment not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed to load shipment", http.StatusInternalServerError)
		return
	}

	// Get list of engineers (if engineer not assigned)
	engineers := []models.SoftwareEngineer{}
	if shipment.SoftwareEngineerID == nil {
		rows, err := h.DB.QueryContext(r.Context(),
			`SELECT id, name, email, address, phone, address_confirmed, created_at 
			FROM software_engineers 
			ORDER BY name`,
		)
		if err != nil {
			http.Error(w, "Failed to load engineers", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var engineer models.SoftwareEngineer
			err := rows.Scan(&engineer.ID, &engineer.Name, &engineer.Email, 
				&engineer.Address, &engineer.Phone, &engineer.AddressConfirmed, 
				&engineer.CreatedAt)
			if err != nil {
				continue
			}
			engineers = append(engineers, engineer)
		}
	} else {
		engineerID.Valid = true
		engineerID.Int64 = *shipment.SoftwareEngineerID
	}

	// Get error and success messages from query parameters
	errorMsg := r.URL.Query().Get("error")
	successMsg := r.URL.Query().Get("success")

	data := map[string]interface{}{
		"Error":        errorMsg,
		"Success":      successMsg,
		"User":         user,
		"Nav":          views.GetNavigationLinks(user.Role),
		"CurrentPage":  "shipments",
		"Shipment":     shipment,
		"CompanyName":  companyName,
		"EngineerName": engineerName.String,
		"EngineerID":   engineerID.Int64,
		"Engineers":    engineers,
	}

	if h.Templates != nil {
		err := h.Templates.ExecuteTemplate(w, "delivery-form.html", data)
		if err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			return
		}
	} else {
		// For testing without templates
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Delivery Form Page")
	}
}

// DeliveryFormSubmit handles the delivery form submission
func (h *DeliveryFormHandler) DeliveryFormSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form - try multipart first (for file uploads), fallback to regular form
	err := r.ParseMultipartForm(MaxUploadSize)
	if err != nil {
		// If multipart parsing fails, try regular form parsing
		err = r.ParseForm()
		if err != nil {
			http.Redirect(w, r, "/delivery-form?error=Invalid+form+data", http.StatusSeeOther)
			return
		}
	}

	// Extract form values
	shipmentIDStr := r.FormValue("shipment_id")
	shipmentID, err := strconv.ParseInt(shipmentIDStr, 10, 64)
	if err != nil {
		http.Redirect(w, r, "/delivery-form?error=Invalid+shipment+ID", http.StatusSeeOther)
		return
	}

	engineerIDStr := r.FormValue("engineer_id")
	engineerID, err := strconv.ParseInt(engineerIDStr, 10, 64)
	if err != nil {
		redirectURL := fmt.Sprintf("/delivery-form?shipment_id=%d&error=Invalid+engineer+ID", shipmentID)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
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
			redirectURL := fmt.Sprintf("/delivery-form?shipment_id=%d&error=File+too+large", shipmentID)
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return
		}

		// Open uploaded file
		file, err := fileHeader.Open()
		if err != nil {
			redirectURL := fmt.Sprintf("/delivery-form?shipment_id=%d&error=Failed+to+read+file", shipmentID)
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return
		}
		defer file.Close()

		// Generate unique filename
		ext := filepath.Ext(fileHeader.Filename)
		filename := fmt.Sprintf("%d_%d%s", shipmentID, time.Now().UnixNano(), ext)
		destPath := filepath.Join(DeliveryUploadDir, filename)

		// Create destination file
		dst, err := os.Create(destPath)
		if err != nil {
			redirectURL := fmt.Sprintf("/delivery-form?shipment_id=%d&error=Failed+to+save+file", shipmentID)
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return
		}
		defer dst.Close()

		// Copy file content
		_, err = io.Copy(dst, file)
		if err != nil {
			redirectURL := fmt.Sprintf("/delivery-form?shipment_id=%d&error=Failed+to+save+file", shipmentID)
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return
		}

		// Store relative URL
		photoURL := fmt.Sprintf("/uploads/delivery/%s", filename)
		photoURLs = append(photoURLs, photoURL)
	}

	// Build validation input
	formInput := validator.DeliveryFormInput{
		ShipmentID: shipmentID,
		EngineerID: engineerID,
		Notes:      notes,
		PhotoURLs:  photoURLs,
	}

	// Validate form
	if err := validator.ValidateDeliveryForm(formInput); err != nil {
		// Clean up uploaded files on validation error
		for _, photoURL := range photoURLs {
			filepath := "." + photoURL
			os.Remove(filepath)
		}

		redirectURL := fmt.Sprintf("/delivery-form?shipment_id=%d&error=%s", 
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

	// Create delivery form
	deliveryForm := models.DeliveryForm{
		ShipmentID: shipmentID,
		EngineerID: engineerID,
		Notes:      notes,
		PhotoURLs:  photoURLs,
	}
	deliveryForm.BeforeCreate()

	_, err = tx.ExecContext(r.Context(),
		`INSERT INTO delivery_forms (shipment_id, engineer_id, delivered_at, notes, photo_urls)
		VALUES ($1, $2, $3, $4, $5)`,
		deliveryForm.ShipmentID, deliveryForm.EngineerID, deliveryForm.DeliveredAt,
		deliveryForm.Notes, pq.Array(deliveryForm.PhotoURLs),
	)
	if err != nil {
		http.Error(w, "Failed to save delivery form", http.StatusInternalServerError)
		return
	}

	// Update shipment status to "delivered" and set engineer if not already set
	now := time.Now()
	_, err = tx.ExecContext(r.Context(),
		`UPDATE shipments 
		SET status = $1, software_engineer_id = $2, delivered_at = $3, updated_at = $4
		WHERE id = $5`,
		models.ShipmentStatusDelivered, engineerID, now, now, shipmentID,
	)
	if err != nil {
		http.Error(w, "Failed to update shipment status", http.StatusInternalServerError)
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	// Mark magic link as used (if accessed via magic link)
	user := middleware.GetUserFromContext(r.Context())
	if user != nil {
		magicLinkToken, err := auth.GetMagicLinkByShipmentAndUser(r.Context(), h.DB, shipmentID, user.ID)
		if err == nil && magicLinkToken != "" {
			_ = auth.MarkMagicLinkAsUsed(r.Context(), h.DB, magicLinkToken)
			// Log error but don't fail the request if marking as used fails
		}
	}

	// Send delivery confirmation email (Step 11-12 in process flow)
	if h.Notifier != nil {
		if err := h.Notifier.SendDeliveryConfirmation(r.Context(), shipmentID); err != nil {
			// Log error but don't fail the request
			fmt.Printf("Warning: Failed to send delivery confirmation email: %v\n", err)
		}
	}

	// Create audit log entry outside transaction (non-critical, user_id would need to be set)
	// Skipping for now as we don't have a valid user_id for delivery forms
	// In production, you might want to add a system user or make user_id nullable

	// Redirect to success page or shipment detail
	redirectURL := fmt.Sprintf("/shipments/%d?success=Delivery+confirmed+successfully", shipmentID)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

