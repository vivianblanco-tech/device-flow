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

	// Get shipment type
	shipmentTypeStr := r.FormValue("shipment_type")
	
	// Detect legacy form format (no shipment_type and has old fields)
	isLegacyForm := shipmentTypeStr == "" && r.FormValue("number_of_laptops") != ""
	
	// Set shipment type
	if isLegacyForm {
		// Legacy forms - will be handled by legacy handler
		shipmentTypeStr = "legacy"
	} else if shipmentTypeStr == "" {
		// New forms without explicit type default to single_full_journey
		shipmentTypeStr = string(models.ShipmentTypeSingleFullJourney)
	}
	
	shipmentType := models.ShipmentType(shipmentTypeStr)

	// Validate shipment type (skip for legacy)
	if shipmentType != "legacy" && !models.IsValidShipmentType(shipmentType) {
		http.Redirect(w, r, "/pickup-form?error=Invalid+shipment+type", http.StatusSeeOther)
		return
	}

	// Parse pickup date
	pickupDateStr := r.FormValue("pickup_date")
	pickupDate, err := time.Parse("2006-01-02", pickupDateStr)
	if err != nil {
		http.Redirect(w, r, "/pickup-form?error=Invalid+date+format", http.StatusSeeOther)
		return
	}

	// Parse include accessories checkbox
	includeAccessories := r.FormValue("include_accessories") == "on" || r.FormValue("include_accessories") == "true"

	var shipmentID int64

	// Branch logic based on shipment type
	switch shipmentType {
	case models.ShipmentTypeSingleFullJourney:
		shipmentID, err = h.handleSingleFullJourneyForm(r, user, companyID, pickupDate, includeAccessories)
		if err != nil {
			redirectURL := fmt.Sprintf("/pickup-form?error=%s&company_id=%d",
				err.Error(), companyID)
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return
		}

	case "legacy":
		// Legacy form handling for backward compatibility
		shipmentID, err = h.handleLegacyPickupForm(r, user, companyID, pickupDate, includeAccessories)
		if err != nil {
			redirectURL := fmt.Sprintf("/pickup-form?error=%s&company_id=%d",
				err.Error(), companyID)
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return
		}

	default:
		http.Redirect(w, r, "/pickup-form?error=Unsupported+shipment+type", http.StatusSeeOther)
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

// handleSingleFullJourneyForm handles single full journey shipment form submission
func (h *PickupFormHandler) handleSingleFullJourneyForm(r *http.Request, user *models.User, companyID int64, pickupDate time.Time, includeAccessories bool) (int64, error) {
	// Build validation input
	formInput := validator.SingleFullJourneyFormInput{
		ClientCompanyID:        companyID,
		ContactName:            r.FormValue("contact_name"),
		ContactEmail:           r.FormValue("contact_email"),
		ContactPhone:           r.FormValue("contact_phone"),
		PickupAddress:          r.FormValue("pickup_address"),
		PickupCity:             r.FormValue("pickup_city"),
		PickupState:            r.FormValue("pickup_state"),
		PickupZip:              r.FormValue("pickup_zip"),
		PickupDate:             r.FormValue("pickup_date"),
		PickupTimeSlot:         r.FormValue("pickup_time_slot"),
		JiraTicketNumber:       r.FormValue("jira_ticket_number"),
		SpecialInstructions:    r.FormValue("special_instructions"),
		LaptopSerialNumber:     r.FormValue("laptop_serial_number"),
		LaptopSpecs:            r.FormValue("laptop_specs"),
		EngineerName:           r.FormValue("engineer_name"),
		IncludeAccessories:     includeAccessories,
		AccessoriesDescription: r.FormValue("accessories_description"),
	}

	// Validate form
	if err := validator.ValidateSingleFullJourneyForm(formInput); err != nil {
		return 0, err
	}

	// Start transaction
	tx, err := h.DB.BeginTx(r.Context(), nil)
	if err != nil {
		return 0, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Create shipment with single_full_journey type
	shipment := models.Shipment{
		ShipmentType:        models.ShipmentTypeSingleFullJourney,
		ClientCompanyID:     companyID,
		Status:              models.ShipmentStatusPendingPickup,
		LaptopCount:         1,
		JiraTicketNumber:    formInput.JiraTicketNumber,
		PickupScheduledDate: &pickupDate,
		Notes:               formInput.SpecialInstructions,
	}
	shipment.BeforeCreate()

	var shipmentID int64
	err = tx.QueryRowContext(r.Context(),
		`INSERT INTO shipments (shipment_type, client_company_id, status, laptop_count, jira_ticket_number, pickup_scheduled_date, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`,
		shipment.ShipmentType, shipment.ClientCompanyID, shipment.Status, shipment.LaptopCount,
		shipment.JiraTicketNumber, shipment.PickupScheduledDate, shipment.Notes,
		shipment.CreatedAt, shipment.UpdatedAt,
	).Scan(&shipmentID)
	if err != nil {
		return 0, fmt.Errorf("failed to create shipment: %w", err)
	}

	// Auto-create laptop record
	laptop := models.Laptop{
		SerialNumber:    formInput.LaptopSerialNumber,
		Specs:           formInput.LaptopSpecs,
		Status:          models.LaptopStatusInTransitToWarehouse,
		ClientCompanyID: &companyID,
	}
	laptop.BeforeCreate()

	var laptopID int64
	err = tx.QueryRowContext(r.Context(),
		`INSERT INTO laptops (serial_number, specs, status, client_company_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`,
		laptop.SerialNumber, laptop.Specs, laptop.Status, laptop.ClientCompanyID,
		laptop.CreatedAt, laptop.UpdatedAt,
	).Scan(&laptopID)
	if err != nil {
		return 0, fmt.Errorf("failed to create laptop: %w", err)
	}

	// Link laptop to shipment
	_, err = tx.ExecContext(r.Context(),
		`INSERT INTO shipment_laptops (shipment_id, laptop_id, created_at)
		VALUES ($1, $2, $3)`,
		shipmentID, laptopID, time.Now(),
	)
	if err != nil {
		return 0, fmt.Errorf("failed to link laptop to shipment: %w", err)
	}

	// Create pickup form with form data as JSONB
	formDataJSON, err := json.Marshal(map[string]interface{}{
		"contact_name":            formInput.ContactName,
		"contact_email":           formInput.ContactEmail,
		"contact_phone":           formInput.ContactPhone,
		"pickup_address":          formInput.PickupAddress,
		"pickup_city":             formInput.PickupCity,
		"pickup_state":            formInput.PickupState,
		"pickup_zip":              formInput.PickupZip,
		"pickup_date":             formInput.PickupDate,
		"pickup_time_slot":        formInput.PickupTimeSlot,
		"jira_ticket_number":      formInput.JiraTicketNumber,
		"special_instructions":    formInput.SpecialInstructions,
		"laptop_serial_number":    formInput.LaptopSerialNumber,
		"laptop_specs":            formInput.LaptopSpecs,
		"engineer_name":           formInput.EngineerName,
		"include_accessories":     formInput.IncludeAccessories,
		"accessories_description": formInput.AccessoriesDescription,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to encode form data: %w", err)
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
		return 0, fmt.Errorf("failed to save pickup form: %w", err)
	}

	// Create audit log entry
	auditDetails, _ := json.Marshal(map[string]interface{}{
		"action":        "pickup_form_submitted",
		"shipment_id":   shipmentID,
		"shipment_type": models.ShipmentTypeSingleFullJourney,
		"company_id":    companyID,
		"laptop_id":     laptopID,
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
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return shipmentID, nil
}

// handleLegacyPickupForm handles legacy pickup form submission for backward compatibility
func (h *PickupFormHandler) handleLegacyPickupForm(r *http.Request, user *models.User, companyID int64, pickupDate time.Time, includeAccessories bool) (int64, error) {
	// Parse legacy form fields
	numberOfLaptopsStr := r.FormValue("number_of_laptops")
	numberOfLaptops, err := strconv.Atoi(numberOfLaptopsStr)
	if err != nil {
		numberOfLaptops = 0
	}

	numberOfBoxesStr := r.FormValue("number_of_boxes")
	numberOfBoxes, err := strconv.Atoi(numberOfBoxesStr)
	if err != nil {
		numberOfBoxes = 0
	}

	// Parse bulk dimensions and weight
	bulkLength, _ := strconv.ParseFloat(r.FormValue("bulk_length"), 64)
	bulkWidth, _ := strconv.ParseFloat(r.FormValue("bulk_width"), 64)
	bulkHeight, _ := strconv.ParseFloat(r.FormValue("bulk_height"), 64)
	bulkWeight, _ := strconv.ParseFloat(r.FormValue("bulk_weight"), 64)

	// Build validation input
	formInput := validator.PickupFormInput{
		ClientCompanyID:        companyID,
		ContactName:            r.FormValue("contact_name"),
		ContactEmail:           r.FormValue("contact_email"),
		ContactPhone:           r.FormValue("contact_phone"),
		PickupAddress:          r.FormValue("pickup_address"),
		PickupCity:             r.FormValue("pickup_city"),
		PickupState:            r.FormValue("pickup_state"),
		PickupZip:              r.FormValue("pickup_zip"),
		PickupDate:             r.FormValue("pickup_date"),
		PickupTimeSlot:         r.FormValue("pickup_time_slot"),
		NumberOfLaptops:        numberOfLaptops,
		JiraTicketNumber:       r.FormValue("jira_ticket_number"),
		SpecialInstructions:    r.FormValue("special_instructions"),
		NumberOfBoxes:          numberOfBoxes,
		AssignmentType:         r.FormValue("assignment_type"),
		BulkLength:             bulkLength,
		BulkWidth:              bulkWidth,
		BulkHeight:             bulkHeight,
		BulkWeight:             bulkWeight,
		IncludeAccessories:     includeAccessories,
		AccessoriesDescription: r.FormValue("accessories_description"),
	}

	// Validate form
	if err := validator.ValidatePickupForm(formInput); err != nil {
		return 0, err
	}

	// Start transaction
	tx, err := h.DB.BeginTx(r.Context(), nil)
	if err != nil {
		return 0, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Create shipment (legacy - defaults to single_full_journey)
	shipment := models.Shipment{
		ShipmentType:        models.ShipmentTypeSingleFullJourney,
		ClientCompanyID:     companyID,
		Status:              models.ShipmentStatusPendingPickup,
		LaptopCount:         1, // Default to 1 for legacy forms
		JiraTicketNumber:    formInput.JiraTicketNumber,
		PickupScheduledDate: &pickupDate,
		Notes:               formInput.SpecialInstructions,
	}
	shipment.BeforeCreate()

	var shipmentID int64
	err = tx.QueryRowContext(r.Context(),
		`INSERT INTO shipments (shipment_type, client_company_id, status, laptop_count, jira_ticket_number, pickup_scheduled_date, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`,
		shipment.ShipmentType, shipment.ClientCompanyID, shipment.Status, shipment.LaptopCount,
		shipment.JiraTicketNumber, shipment.PickupScheduledDate, shipment.Notes,
		shipment.CreatedAt, shipment.UpdatedAt,
	).Scan(&shipmentID)
	if err != nil {
		return 0, fmt.Errorf("failed to create shipment: %w", err)
	}

	// Create pickup form with form data as JSONB
	formDataJSON, err := json.Marshal(map[string]interface{}{
		"contact_name":            formInput.ContactName,
		"contact_email":           formInput.ContactEmail,
		"contact_phone":           formInput.ContactPhone,
		"pickup_address":          formInput.PickupAddress,
		"pickup_city":             formInput.PickupCity,
		"pickup_state":            formInput.PickupState,
		"pickup_zip":              formInput.PickupZip,
		"pickup_date":             formInput.PickupDate,
		"pickup_time_slot":        formInput.PickupTimeSlot,
		"number_of_laptops":       formInput.NumberOfLaptops,
		"jira_ticket_number":      formInput.JiraTicketNumber,
		"special_instructions":    formInput.SpecialInstructions,
		"number_of_boxes":         formInput.NumberOfBoxes,
		"assignment_type":         formInput.AssignmentType,
		"bulk_length":             formInput.BulkLength,
		"bulk_width":              formInput.BulkWidth,
		"bulk_height":             formInput.BulkHeight,
		"bulk_weight":             formInput.BulkWeight,
		"include_accessories":     formInput.IncludeAccessories,
		"accessories_description": formInput.AccessoriesDescription,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to encode form data: %w", err)
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
		return 0, fmt.Errorf("failed to save pickup form: %w", err)
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
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return shipmentID, nil
}
