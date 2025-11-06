package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/yourusername/laptop-tracking-system/internal/email"
	"github.com/yourusername/laptop-tracking-system/internal/middleware"
	"github.com/yourusername/laptop-tracking-system/internal/models"
	"github.com/yourusername/laptop-tracking-system/internal/validator"
)

// PickupFormHandler handles pickup form requests
type PickupFormHandler struct {
	DB        *sql.DB
	Templates *template.Template
	Notifier  *email.Notifier
}

// NewPickupFormHandler creates a new PickupFormHandler
func NewPickupFormHandler(db *sql.DB, templates *template.Template, notifier *email.Notifier) *PickupFormHandler {
	return &PickupFormHandler{
		DB:        db,
		Templates: templates,
		Notifier:  notifier,
	}
}

// PickupFormPage displays the pickup form
func (h *PickupFormHandler) PickupFormPage(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Get company ID from query parameter (for client users)
	companyIDStr := r.URL.Query().Get("company_id")
	var companyID int64
	if companyIDStr != "" {
		var err error
		companyID, err = strconv.ParseInt(companyIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid company ID", http.StatusBadRequest)
			return
		}
	}

	// Get error and success messages from query parameters
	errorMsg := r.URL.Query().Get("error")
	successMsg := r.URL.Query().Get("success")

	// Get list of client companies (for logistics users)
	companies := []models.ClientCompany{}
	if user.Role == models.RoleLogistics {
		rows, err := h.DB.QueryContext(r.Context(),
			`SELECT id, name, contact_info, created_at FROM client_companies ORDER BY name`,
		)
		if err != nil {
			http.Error(w, "Failed to load companies", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var company models.ClientCompany
			err := rows.Scan(&company.ID, &company.Name, &company.ContactInfo, &company.CreatedAt)
			if err != nil {
				continue
			}
			companies = append(companies, company)
		}
	}

	data := map[string]interface{}{
		"Error":     errorMsg,
		"Success":   successMsg,
		"User":      user,
		"CompanyID": companyID,
		"Companies": companies,
		"TimeSlots": []string{"morning", "afternoon", "evening"},
	}

	if h.Templates != nil {
		err := h.Templates.ExecuteTemplate(w, "pickup-form.html", data)
		if err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			return
		}
	} else {
		// For testing without templates
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Pickup Form Page")
	}
}

// PickupFormSubmit handles the pickup form submission
func (h *PickupFormHandler) PickupFormSubmit(w http.ResponseWriter, r *http.Request) {
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

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Redirect(w, r, "/pickup-form?error=Invalid+form+data", http.StatusSeeOther)
		return
	}

	// Extract form values
	companyIDStr := r.FormValue("client_company_id")
	companyID, err := strconv.ParseInt(companyIDStr, 10, 64)
	if err != nil {
		http.Redirect(w, r, "/pickup-form?error=Invalid+company+ID", http.StatusSeeOther)
		return
	}

	numberOfLaptopsStr := r.FormValue("number_of_laptops")
	numberOfLaptops, err := strconv.Atoi(numberOfLaptopsStr)
	if err != nil {
		numberOfLaptops = 0
	}

	// Build validation input
	formInput := validator.PickupFormInput{
		ClientCompanyID:     companyID,
		ContactName:         r.FormValue("contact_name"),
		ContactEmail:        r.FormValue("contact_email"),
		ContactPhone:        r.FormValue("contact_phone"),
		PickupAddress:       r.FormValue("pickup_address"),
		PickupDate:          r.FormValue("pickup_date"),
		PickupTimeSlot:      r.FormValue("pickup_time_slot"),
		NumberOfLaptops:     numberOfLaptops,
		JiraTicketNumber:    r.FormValue("jira_ticket_number"),
		SpecialInstructions: r.FormValue("special_instructions"),
	}

	// Validate form
	if err := validator.ValidatePickupForm(formInput); err != nil {
		redirectURL := fmt.Sprintf("/pickup-form?error=%s&company_id=%d",
			err.Error(), companyID)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}

	// Parse pickup date
	pickupDate, err := time.Parse("2006-01-02", formInput.PickupDate)
	if err != nil {
		http.Redirect(w, r, "/pickup-form?error=Invalid+date+format", http.StatusSeeOther)
		return
	}

	// Start transaction
	tx, err := h.DB.BeginTx(r.Context(), nil)
	if err != nil {
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Create shipment
	shipment := models.Shipment{
		ClientCompanyID:     companyID,
		Status:              models.ShipmentStatusPendingPickup,
		JiraTicketNumber:    formInput.JiraTicketNumber,
		PickupScheduledDate: &pickupDate,
		Notes:               formInput.SpecialInstructions,
	}
	shipment.BeforeCreate()

	var shipmentID int64
	err = tx.QueryRowContext(r.Context(),
		`INSERT INTO shipments (client_company_id, status, jira_ticket_number, pickup_scheduled_date, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`,
		shipment.ClientCompanyID, shipment.Status, shipment.JiraTicketNumber, shipment.PickupScheduledDate,
		shipment.Notes, shipment.CreatedAt, shipment.UpdatedAt,
	).Scan(&shipmentID)
	if err != nil {
		http.Error(w, "Failed to create shipment", http.StatusInternalServerError)
		return
	}

	// Create pickup form with form data as JSONB
	formDataJSON, err := json.Marshal(map[string]interface{}{
		"contact_name":         formInput.ContactName,
		"contact_email":        formInput.ContactEmail,
		"contact_phone":        formInput.ContactPhone,
		"pickup_address":       formInput.PickupAddress,
		"pickup_date":          formInput.PickupDate,
		"pickup_time_slot":     formInput.PickupTimeSlot,
		"number_of_laptops":    formInput.NumberOfLaptops,
		"jira_ticket_number":   formInput.JiraTicketNumber,
		"special_instructions": formInput.SpecialInstructions,
	})
	if err != nil {
		http.Error(w, "Failed to encode form data", http.StatusInternalServerError)
		return
	}

	pickupForm := models.PickupForm{
		ShipmentID:        shipmentID,
		SubmittedByUserID: user.ID,
		FormData:          formDataJSON,
	}
	pickupForm.BeforeCreate()

	_, err = tx.ExecContext(r.Context(),
		`INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data)
		VALUES ($1, $2, $3, $4)`,
		pickupForm.ShipmentID, pickupForm.SubmittedByUserID,
		pickupForm.SubmittedAt, pickupForm.FormData,
	)
	if err != nil {
		http.Error(w, "Failed to save pickup form", http.StatusInternalServerError)
		return
	}

	// Create audit log entry
	auditDetails, _ := json.Marshal(map[string]interface{}{
		"action":      "pickup_form_submitted",
		"shipment_id": shipmentID,
		"company_id":  companyID,
	})

	_, err = tx.ExecContext(r.Context(),
		`INSERT INTO audit_logs (user_id, action, entity_type, entity_id, timestamp, details)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		user.ID, "pickup_form_submitted", "shipment", shipmentID, time.Now(), auditDetails,
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

	// Send pickup confirmation email (Step 4 in process flow)
	if h.Notifier != nil {
		if err := h.Notifier.SendPickupConfirmation(r.Context(), shipmentID); err != nil {
			// Log error but don't fail the request
			fmt.Printf("Warning: Failed to send pickup confirmation email: %v\n", err)
		}
	}

	// Redirect to success page or shipment detail
	redirectURL := fmt.Sprintf("/shipments/%d?success=Pickup+form+submitted+successfully", shipmentID)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}
