package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/yourusername/laptop-tracking-system/internal/email"
	"github.com/yourusername/laptop-tracking-system/internal/middleware"
	"github.com/yourusername/laptop-tracking-system/internal/models"
	"github.com/yourusername/laptop-tracking-system/internal/validator"
	"github.com/yourusername/laptop-tracking-system/internal/views"
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

// PickupFormsLandingPage displays the pickup forms landing page with options for different form types
func (h *PickupFormHandler) PickupFormsLandingPage(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Determine which form options to show based on user role
	showWarehouseToEngineer := user.Role == models.RoleLogistics

	data := map[string]interface{}{
		"User":                    user,
		"Nav":                     views.GetNavigationLinks(user.Role),
		"CurrentPage":             "pickup-forms",
		"ShowWarehouseToEngineer": showWarehouseToEngineer,
	}

	if h.Templates != nil {
		err := h.Templates.ExecuteTemplate(w, "pickup-forms-landing.html", data)
		if err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			return
		}
	} else {
		// For testing without templates
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Pickup Forms Landing Page")
		fmt.Fprintf(w, "/shipments/create/single")
		fmt.Fprintf(w, "/shipments/create/bulk")
		if showWarehouseToEngineer {
			fmt.Fprintf(w, "/shipments/create/warehouse-to-engineer")
		}
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
		"Error":       errorMsg,
		"Success":     successMsg,
		"User":        user,
		"Nav":         views.GetNavigationLinks(user.Role),
		"CurrentPage": "pickup-forms",
		"CompanyID":   companyID,
		"Companies":   companies,
		"TimeSlots":   []string{"morning", "afternoon", "evening"},
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

// SingleShipmentFormPage displays the single full journey shipment form
func (h *PickupFormHandler) SingleShipmentFormPage(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
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
		"Error":        errorMsg,
		"Success":      successMsg,
		"User":         user,
		"Nav":          views.GetNavigationLinks(user.Role),
		"CurrentPage":  "pickup-forms",
		"Companies":    companies,
		"TimeSlots":    []string{"morning", "afternoon", "evening"},
		"ShipmentType": models.ShipmentTypeSingleFullJourney,
	}

	if h.Templates != nil {
		err := h.Templates.ExecuteTemplate(w, "single-shipment-form.html", data)
		if err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			return
		}
	} else {
		// For testing without templates
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Single Shipment Form Page")
	}
}

// BulkShipmentFormPage displays the bulk to warehouse shipment form
func (h *PickupFormHandler) BulkShipmentFormPage(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
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
		"Error":        errorMsg,
		"Success":      successMsg,
		"User":         user,
		"Nav":          views.GetNavigationLinks(user.Role),
		"CurrentPage":  "pickup-forms",
		"Companies":    companies,
		"TimeSlots":    []string{"morning", "afternoon", "evening"},
		"ShipmentType": models.ShipmentTypeBulkToWarehouse,
	}

	if h.Templates != nil {
		err := h.Templates.ExecuteTemplate(w, "bulk-shipment-form.html", data)
		if err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			return
		}
	} else {
		// For testing without templates
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Bulk Shipment Form Page")
	}
}

// WarehouseToEngineerFormPage displays the warehouse to engineer shipment form
func (h *PickupFormHandler) WarehouseToEngineerFormPage(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Get error and success messages from query parameters
	errorMsg := r.URL.Query().Get("error")
	successMsg := r.URL.Query().Get("success")

	// Get list of available laptops from inventory
	laptops := []models.Laptop{}
	rows, err := h.DB.QueryContext(r.Context(), `
		SELECT DISTINCT l.id, l.serial_number, l.sku, l.brand, l.model, l.ram_gb, l.ssd_gb,
			   l.status, l.client_company_id, l.software_engineer_id,
			   l.created_at, l.updated_at,
			   cc.name as client_company_name
		FROM laptops l
		LEFT JOIN client_companies cc ON cc.id = l.client_company_id
		WHERE l.status IN ('available', 'at_warehouse')
		  -- Must have a reception report
		  AND EXISTS (
			  SELECT 1 FROM reception_reports rr
			  JOIN shipments s ON s.id = rr.shipment_id
			  JOIN shipment_laptops sl ON sl.shipment_id = s.id
			  WHERE sl.laptop_id = l.id
		  )
		  -- Must not be in any active shipment (except bulk shipments at warehouse)
		  AND NOT EXISTS (
			  SELECT 1 FROM shipment_laptops sl
			  JOIN shipments s ON s.id = sl.shipment_id
			  WHERE sl.laptop_id = l.id
				AND s.status NOT IN ('delivered', 'at_warehouse')
		  )
		ORDER BY l.created_at DESC
	`)
	if err != nil {
		http.Error(w, "Failed to load laptops", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var laptop models.Laptop
		var clientCompanyName sql.NullString
		err := rows.Scan(
			&laptop.ID, &laptop.SerialNumber, &laptop.SKU, &laptop.Brand,
			&laptop.Model, &laptop.RAMGB, &laptop.SSDGB, &laptop.Status, &laptop.ClientCompanyID,
			&laptop.SoftwareEngineerID, &laptop.CreatedAt, &laptop.UpdatedAt,
			&clientCompanyName,
		)
		if err != nil {
			continue
		}
		if clientCompanyName.Valid {
			laptop.ClientCompanyName = clientCompanyName.String
		}
		laptops = append(laptops, laptop)
	}

	data := map[string]interface{}{
		"Error":        errorMsg,
		"Success":      successMsg,
		"User":         user,
		"Nav":          views.GetNavigationLinks(user.Role),
		"CurrentPage":  "pickup-forms",
		"Laptops":      laptops,
		"ShipmentType": models.ShipmentTypeWarehouseToEngineer,
	}

	if h.Templates != nil {
		err := h.Templates.ExecuteTemplate(w, "warehouse-to-engineer-form.html", data)
		if err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			return
		}
	} else {
		// For testing without templates
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Warehouse to Engineer Form Page")
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

	// Get shipment type first (needed to determine if we need client_company_id)
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

	// Extract company ID
	var companyID int64
	companyIDStr := r.FormValue("client_company_id")

	// For warehouse-to-engineer shipments, extract company ID from laptop if not provided
	if shipmentType == models.ShipmentTypeWarehouseToEngineer && companyIDStr == "" {
		laptopIDStr := r.FormValue("laptop_id")
		if laptopIDStr != "" {
			laptopID, err := strconv.ParseInt(laptopIDStr, 10, 64)
			if err == nil {
				// Try to get company ID from laptop first
				var nullableCompanyID sql.NullInt64
				err = h.DB.QueryRowContext(r.Context(),
					`SELECT client_company_id FROM laptops WHERE id = $1`,
					laptopID,
				).Scan(&nullableCompanyID)

				if err != nil {
					http.Redirect(w, r, "/shipments/create/warehouse-to-engineer?error=Laptop+not+found", http.StatusSeeOther)
					return
				}

				if nullableCompanyID.Valid {
					// Laptop has a company ID
					companyID = nullableCompanyID.Int64
				} else {
					// Laptop has NULL company_id, get it from the shipment it came from
					err = h.DB.QueryRowContext(r.Context(),
						`SELECT s.client_company_id
						 FROM shipments s
						 JOIN shipment_laptops sl ON sl.shipment_id = s.id
						 WHERE sl.laptop_id = $1
						 ORDER BY s.created_at DESC
						 LIMIT 1`,
						laptopID,
					).Scan(&companyID)

					if err != nil {
						http.Redirect(w, r, "/shipments/create/warehouse-to-engineer?error=Unable+to+find+laptop+company", http.StatusSeeOther)
						return
					}
				}
			}
		}
	} else if companyIDStr != "" {
		// Parse company ID from form for other shipment types
		var err error
		companyID, err = strconv.ParseInt(companyIDStr, 10, 64)
		if err != nil {
			http.Redirect(w, r, "/pickup-form?error=Invalid+company+ID", http.StatusSeeOther)
			return
		}
	} else {
		// Company ID is required for non-warehouse-to-engineer shipments
		http.Redirect(w, r, "/pickup-form?error=Company+ID+is+required", http.StatusSeeOther)
		return
	}

	// Parse pickup date (not required for warehouse-to-engineer)
	var pickupDate time.Time
	var hasPickupDate bool
	pickupDateStr := r.FormValue("pickup_date")
	if pickupDateStr != "" {
		var err error
		pickupDate, err = time.Parse("2006-01-02", pickupDateStr)
		if err != nil {
			http.Redirect(w, r, "/pickup-form?error=Invalid+date+format", http.StatusSeeOther)
			return
		}
		hasPickupDate = true
	}

	// Parse include accessories checkbox
	includeAccessories := r.FormValue("include_accessories") == "on" || r.FormValue("include_accessories") == "true"

	var shipmentID int64

	// Branch logic based on shipment type
	switch shipmentType {
	case models.ShipmentTypeSingleFullJourney:
		if !hasPickupDate {
			http.Redirect(w, r, "/pickup-form?error=Pickup+date+is+required", http.StatusSeeOther)
			return
		}
		shipmentID, err = h.handleSingleFullJourneyForm(r, user, companyID, pickupDate, includeAccessories)
		if err != nil {
			redirectURL := fmt.Sprintf("/pickup-form?error=%s&company_id=%d",
				err.Error(), companyID)
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return
		}

	case models.ShipmentTypeBulkToWarehouse:
		// Pickup date is optional for minimal bulk shipment creation (logistics only)
		// It's required for full form submission (client completing details)
		isMinimalCreation := r.FormValue("contact_name") == "" && r.FormValue("contact_email") == ""
		if !hasPickupDate && !isMinimalCreation {
			http.Redirect(w, r, "/pickup-form?error=Pickup+date+is+required", http.StatusSeeOther)
			return
		}
		shipmentID, err = h.handleBulkToWarehouseForm(r, user, companyID, pickupDate, includeAccessories)
		if err != nil {
			redirectURL := fmt.Sprintf("/pickup-form?error=%s&company_id=%d",
				err.Error(), companyID)
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return
		}

	case models.ShipmentTypeWarehouseToEngineer:
		// Warehouse-to-engineer does not require pickup date
		shipmentID, err = h.handleWarehouseToEngineerForm(r, user, companyID, includeAccessories)
		if err != nil {
			redirectURL := fmt.Sprintf("/shipments/create/warehouse-to-engineer?error=%s",
				err.Error())
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return
		}

	case "legacy":
		if !hasPickupDate {
			http.Redirect(w, r, "/pickup-form?error=Pickup+date+is+required", http.StatusSeeOther)
			return
		}
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
		LaptopBrand:            r.FormValue("laptop_brand"),
		LaptopModel:            r.FormValue("laptop_model"),
		LaptopRAMGB:            r.FormValue("laptop_ram_gb"),
		LaptopSSDGB:            r.FormValue("laptop_ssd_gb"),
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
		Brand:           formInput.LaptopBrand,
		Model:           formInput.LaptopModel,
		RAMGB:           formInput.LaptopRAMGB,
		SSDGB:           formInput.LaptopSSDGB,
		Status:          models.LaptopStatusInTransitToWarehouse,
		ClientCompanyID: &companyID,
	}
	laptop.BeforeCreate()

	var laptopID int64
	err = tx.QueryRowContext(r.Context(),
		`INSERT INTO laptops (serial_number, brand, model, ram_gb, ssd_gb, status, client_company_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`,
		laptop.SerialNumber, laptop.Brand, laptop.Model, laptop.RAMGB, laptop.SSDGB, laptop.Status, laptop.ClientCompanyID,
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
		"laptop_brand":            formInput.LaptopBrand,
		"laptop_model":            formInput.LaptopModel,
		"laptop_ram_gb":           formInput.LaptopRAMGB,
		"laptop_ssd_gb":           formInput.LaptopSSDGB,
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

// handleBulkToWarehouseForm handles bulk to warehouse shipment form submission
func (h *PickupFormHandler) handleBulkToWarehouseForm(r *http.Request, user *models.User, companyID int64, pickupDate time.Time, includeAccessories bool) (int64, error) {
	// Check if this is a minimal creation (logistics user creating with only JIRA + Company)
	// Minimal creation means no contact information is provided
	isMinimalCreation := r.FormValue("contact_name") == "" && r.FormValue("contact_email") == ""

	// Parse JIRA ticket (always required)
	jiraTicketNumber := r.FormValue("jira_ticket_number")

	// For minimal creation, validate only JIRA ticket and company
	if isMinimalCreation {
		// Only logistics users can create minimal shipments
		if user.Role != models.RoleLogistics {
			return 0, fmt.Errorf("only logistics users can create minimal bulk shipments")
		}

		// Validate JIRA ticket format
		if jiraTicketNumber == "" {
			return 0, fmt.Errorf("JIRA ticket number is required")
		}

		// Start transaction
		tx, err := h.DB.BeginTx(r.Context(), nil)
		if err != nil {
			return 0, fmt.Errorf("failed to start transaction: %w", err)
		}
		defer tx.Rollback()

		// Create minimal shipment (laptop_count defaults to 1, will be updated later)
		shipment := models.Shipment{
			ShipmentType:     models.ShipmentTypeBulkToWarehouse,
			ClientCompanyID:  companyID,
			Status:           models.ShipmentStatusPendingPickup,
			LaptopCount:      1, // Default, will be updated when client completes form
			JiraTicketNumber: jiraTicketNumber,
		}
		shipment.BeforeCreate()

		var shipmentID int64
		err = tx.QueryRowContext(r.Context(),
			`INSERT INTO shipments (shipment_type, client_company_id, status, laptop_count, jira_ticket_number, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id`,
			shipment.ShipmentType, shipment.ClientCompanyID, shipment.Status, shipment.LaptopCount,
			shipment.JiraTicketNumber, shipment.CreatedAt, shipment.UpdatedAt,
		).Scan(&shipmentID)
		if err != nil {
			return 0, fmt.Errorf("failed to create shipment: %w", err)
		}

		// Create audit log entry
		auditDetails, _ := json.Marshal(map[string]interface{}{
			"action":             "minimal_bulk_shipment_created",
			"shipment_id":        shipmentID,
			"jira_ticket_number": jiraTicketNumber,
		})

		_, err = tx.ExecContext(r.Context(),
			`INSERT INTO audit_logs (user_id, action, entity_type, entity_id, timestamp, details)
			VALUES ($1, $2, $3, $4, $5, $6)`,
			user.ID, "minimal_bulk_shipment_created", "shipment", shipmentID, time.Now(), auditDetails,
		)
		if err != nil {
			// Non-critical, log and continue
			fmt.Printf("Failed to create audit log: %v\n", err)
		}

		// Commit transaction
		if err = tx.Commit(); err != nil {
			return 0, fmt.Errorf("failed to commit transaction: %w", err)
		}

		return shipmentID, nil
	}

	// Full form submission (not minimal) - validate all fields
	// Parse bulk-specific fields
	numberOfLaptopsStr := r.FormValue("number_of_laptops")
	numberOfLaptops, err := strconv.Atoi(numberOfLaptopsStr)
	if err != nil {
		return 0, fmt.Errorf("invalid laptop count")
	}

	bulkLength, _ := strconv.ParseFloat(r.FormValue("bulk_length"), 64)
	bulkWidth, _ := strconv.ParseFloat(r.FormValue("bulk_width"), 64)
	bulkHeight, _ := strconv.ParseFloat(r.FormValue("bulk_height"), 64)
	bulkWeight, _ := strconv.ParseFloat(r.FormValue("bulk_weight"), 64)

	// Build validation input
	formInput := validator.BulkToWarehouseFormInput{
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
		JiraTicketNumber:       jiraTicketNumber,
		SpecialInstructions:    r.FormValue("special_instructions"),
		NumberOfLaptops:        numberOfLaptops,
		BulkLength:             bulkLength,
		BulkWidth:              bulkWidth,
		BulkHeight:             bulkHeight,
		BulkWeight:             bulkWeight,
		IncludeAccessories:     includeAccessories,
		AccessoriesDescription: r.FormValue("accessories_description"),
	}

	// Validate form
	if err := validator.ValidateBulkToWarehouseForm(formInput); err != nil {
		return 0, err
	}

	// Start transaction
	tx, err := h.DB.BeginTx(r.Context(), nil)
	if err != nil {
		return 0, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Create shipment with bulk_to_warehouse type (NO laptops created)
	shipment := models.Shipment{
		ShipmentType:        models.ShipmentTypeBulkToWarehouse,
		ClientCompanyID:     companyID,
		Status:              models.ShipmentStatusPendingPickup,
		LaptopCount:         numberOfLaptops,
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

	// Create pickup form with form data as JSONB (including bulk dimensions)
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
		"number_of_laptops":       formInput.NumberOfLaptops,
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
		"action":        "pickup_form_submitted",
		"shipment_id":   shipmentID,
		"shipment_type": models.ShipmentTypeBulkToWarehouse,
		"company_id":    companyID,
		"laptop_count":  numberOfLaptops,
		"bulk_length":   bulkLength,
		"bulk_width":    bulkWidth,
		"bulk_height":   bulkHeight,
		"bulk_weight":   bulkWeight,
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

// handleWarehouseToEngineerForm handles warehouse-to-engineer shipment form submission
func (h *PickupFormHandler) handleWarehouseToEngineerForm(r *http.Request, user *models.User, companyID int64, includeAccessories bool) (int64, error) {
	// Parse laptop selection
	laptopIDStr := r.FormValue("laptop_id")
	if laptopIDStr == "" {
		return 0, fmt.Errorf("laptop selection is required")
	}
	laptopID, err := strconv.ParseInt(laptopIDStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid laptop ID")
	}

	// Parse software engineer (can be ID or name)
	var softwareEngineerID *int64
	engineerIDStr := r.FormValue("software_engineer_id")
	if engineerIDStr != "" {
		engineerID, err := strconv.ParseInt(engineerIDStr, 10, 64)
		if err == nil {
			softwareEngineerID = &engineerID
		}
	}

	// Build validation input
	formInput := validator.WarehouseToEngineerFormInput{
		LaptopID:            laptopID,
		SoftwareEngineerID:  0,
		EngineerName:        r.FormValue("engineer_name"),
		EngineerEmail:       r.FormValue("engineer_email"),
		EngineerAddress:     r.FormValue("engineer_address"),
		EngineerCity:        r.FormValue("engineer_city"),
		EngineerState:       r.FormValue("engineer_state"),
		EngineerZip:         r.FormValue("engineer_zip"),
		CourierName:         r.FormValue("courier_name"),
		TrackingNumber:      r.FormValue("tracking_number"),
		JiraTicketNumber:    r.FormValue("jira_ticket_number"),
		SpecialInstructions: r.FormValue("special_instructions"),
	}
	if softwareEngineerID != nil {
		formInput.SoftwareEngineerID = *softwareEngineerID
	}

	// Validate form
	if err := validator.ValidateWarehouseToEngineerForm(formInput); err != nil {
		return 0, err
	}

	// Start transaction
	tx, err := h.DB.BeginTx(r.Context(), nil)
	if err != nil {
		return 0, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Verify laptop exists and is available
	var currentLaptopStatus models.LaptopStatus
	err = tx.QueryRowContext(r.Context(),
		`SELECT status FROM laptops WHERE id = $1`,
		laptopID,
	).Scan(&currentLaptopStatus)
	if err == sql.ErrNoRows {
		return 0, fmt.Errorf("laptop not found")
	} else if err != nil {
		return 0, fmt.Errorf("failed to query laptop: %w", err)
	}

	// Verify laptop is available (not in active shipment)
	if currentLaptopStatus != models.LaptopStatusAvailable && currentLaptopStatus != models.LaptopStatusAtWarehouse {
		return 0, fmt.Errorf("laptop is not available for shipment (current status: %s)", currentLaptopStatus)
	}

	// Verify laptop has a reception report (came through warehouse)
	var hasReceptionReport bool
	err = tx.QueryRowContext(r.Context(),
		`SELECT EXISTS(
			SELECT 1 FROM reception_reports rr
			JOIN shipment_laptops sl ON sl.shipment_id = rr.shipment_id
			WHERE sl.laptop_id = $1
		)`,
		laptopID,
	).Scan(&hasReceptionReport)
	if err != nil {
		return 0, fmt.Errorf("failed to check reception report: %w", err)
	}
	if !hasReceptionReport {
		return 0, fmt.Errorf("laptop must have a completed reception report before shipping to engineer")
	}

	// Create shipment with warehouse_to_engineer type
	shipment := models.Shipment{
		ShipmentType:       models.ShipmentTypeWarehouseToEngineer,
		ClientCompanyID:    companyID,
		Status:             models.ShipmentStatusReleasedFromWarehouse, // Start at released status
		LaptopCount:        1,
		SoftwareEngineerID: softwareEngineerID,
		JiraTicketNumber:   formInput.JiraTicketNumber,
		Notes:              formInput.SpecialInstructions,
	}
	shipment.BeforeCreate()

	var shipmentID int64
	err = tx.QueryRowContext(r.Context(),
		`INSERT INTO shipments (shipment_type, client_company_id, status, laptop_count, software_engineer_id, jira_ticket_number, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`,
		shipment.ShipmentType, shipment.ClientCompanyID, shipment.Status, shipment.LaptopCount,
		shipment.SoftwareEngineerID, shipment.JiraTicketNumber, shipment.Notes,
		shipment.CreatedAt, shipment.UpdatedAt,
	).Scan(&shipmentID)
	if err != nil {
		return 0, fmt.Errorf("failed to create shipment: %w", err)
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

	// Update laptop status to in_transit_to_engineer
	_, err = tx.ExecContext(r.Context(),
		`UPDATE laptops SET status = $1, updated_at = $2 WHERE id = $3`,
		models.LaptopStatusInTransitToEngineer, time.Now(), laptopID,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to update laptop status: %w", err)
	}

	// Create pickup form record with all data as JSON
	formData := map[string]interface{}{
		"contact_name":         formInput.EngineerName,
		"contact_email":        formInput.EngineerEmail,
		"delivery_address":     formInput.EngineerAddress,
		"delivery_city":        formInput.EngineerCity,
		"delivery_state":       formInput.EngineerState,
		"delivery_zip":         formInput.EngineerZip,
		"courier_name":         formInput.CourierName,
		"tracking_number":      formInput.TrackingNumber,
		"include_accessories":  includeAccessories,
		"special_instructions": formInput.SpecialInstructions,
		"laptop_id":            laptopID,
	}
	formDataJSON, _ := json.Marshal(formData)

	pickupForm := models.PickupForm{
		ShipmentID:        shipmentID,
		SubmittedByUserID: user.ID,
		SubmittedAt:       time.Now(),
		FormData:          formDataJSON,
	}

	_, err = tx.ExecContext(r.Context(),
		`INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data)
		VALUES ($1, $2, $3, $4)`,
		pickupForm.ShipmentID, pickupForm.SubmittedByUserID, pickupForm.SubmittedAt, pickupForm.FormData,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to create pickup form: %w", err)
	}

	// Create audit log entry
	auditDetails, _ := json.Marshal(map[string]interface{}{
		"action":               "warehouse_to_engineer_form_submitted",
		"shipment_id":          shipmentID,
		"shipment_type":        models.ShipmentTypeWarehouseToEngineer,
		"company_id":           companyID,
		"laptop_id":            laptopID,
		"software_engineer_id": softwareEngineerID,
	})

	_, err = tx.ExecContext(r.Context(),
		`INSERT INTO audit_logs (user_id, action, entity_type, entity_id, timestamp, details)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		user.ID, "warehouse_to_engineer_form_submitted", "shipment", shipmentID, time.Now(), auditDetails,
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

// CreateMinimalSingleShipment creates a minimal single shipment with only JIRA ticket and company ID
// This is used by logistics users to initiate a shipment before sending magic link to client
func (h *PickupFormHandler) CreateMinimalSingleShipment(w http.ResponseWriter, r *http.Request) {
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

	// Only logistics users can create minimal shipments
	if user.Role != models.RoleLogistics {
		http.Error(w, "Forbidden: Only logistics users can create minimal shipments", http.StatusForbidden)
		return
	}

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Redirect(w, r, "/dashboard?error=Invalid+form+data", http.StatusSeeOther)
		return
	}

	// Extract and validate company ID
	companyIDStr := r.FormValue("client_company_id")
	if companyIDStr == "" {
		http.Redirect(w, r, "/dashboard?error=Company+ID+is+required", http.StatusSeeOther)
		return
	}
	companyID, err := strconv.ParseInt(companyIDStr, 10, 64)
	if err != nil {
		http.Redirect(w, r, "/dashboard?error=Invalid+company+ID", http.StatusSeeOther)
		return
	}

	// Extract and validate JIRA ticket number
	jiraTicketNumber := r.FormValue("jira_ticket_number")
	if jiraTicketNumber == "" {
		http.Redirect(w, r, "/dashboard?error=JIRA+ticket+number+is+required", http.StatusSeeOther)
		return
	}

	// Start transaction
	tx, err := h.DB.BeginTx(r.Context(), nil)
	if err != nil {
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Create minimal shipment (no pickup form, no laptop record yet)
	shipment := models.Shipment{
		ShipmentType:     models.ShipmentTypeSingleFullJourney,
		ClientCompanyID:  companyID,
		Status:           models.ShipmentStatusPendingPickup,
		LaptopCount:      1,
		JiraTicketNumber: jiraTicketNumber,
	}
	shipment.BeforeCreate()

	var shipmentID int64
	err = tx.QueryRowContext(r.Context(),
		`INSERT INTO shipments (shipment_type, client_company_id, status, laptop_count, jira_ticket_number, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`,
		shipment.ShipmentType, shipment.ClientCompanyID, shipment.Status, shipment.LaptopCount,
		shipment.JiraTicketNumber, shipment.CreatedAt, shipment.UpdatedAt,
	).Scan(&shipmentID)
	if err != nil {
		http.Error(w, "Failed to create shipment", http.StatusInternalServerError)
		return
	}

	// Create audit log entry
	auditDetails, _ := json.Marshal(map[string]interface{}{
		"action":             "minimal_shipment_created",
		"shipment_id":        shipmentID,
		"shipment_type":      models.ShipmentTypeSingleFullJourney,
		"company_id":         companyID,
		"jira_ticket_number": jiraTicketNumber,
	})

	_, err = tx.ExecContext(r.Context(),
		`INSERT INTO audit_logs (user_id, action, entity_type, entity_id, timestamp, details)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		user.ID, "minimal_shipment_created", "shipment", shipmentID, time.Now(), auditDetails,
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

	// Redirect to shipment detail page
	redirectURL := fmt.Sprintf("/shipments/%d?success=Shipment+created+successfully.+Send+magic+link+to+client+to+complete+details", shipmentID)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

// CompleteShipmentDetails handles client completing shipment details via magic link
func (h *PickupFormHandler) CompleteShipmentDetails(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Extract and validate shipment ID
	shipmentIDStr := r.FormValue("shipment_id")
	if shipmentIDStr == "" {
		http.Error(w, "Shipment ID is required", http.StatusBadRequest)
		return
	}
	shipmentID, err := strconv.ParseInt(shipmentIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid shipment ID", http.StatusBadRequest)
		return
	}

	// Check if shipment already has a pickup form (prevent duplicate submissions)
	var existingFormCount int
	err = h.DB.QueryRowContext(r.Context(),
		`SELECT COUNT(*) FROM pickup_forms WHERE shipment_id = $1`,
		shipmentID,
	).Scan(&existingFormCount)
	if err != nil {
		http.Error(w, "Failed to check shipment status", http.StatusInternalServerError)
		return
	}
	if existingFormCount > 0 {
		redirectURL := fmt.Sprintf("/shipments/%d?error=Shipment+details+already+completed", shipmentID)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}

	// Parse pickup date
	pickupDateStr := r.FormValue("pickup_date")
	if pickupDateStr == "" {
		redirectURL := fmt.Sprintf("/shipments/%d?error=Pickup+date+is+required", shipmentID)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}
	pickupDate, err := time.Parse("2006-01-02", pickupDateStr)
	if err != nil {
		redirectURL := fmt.Sprintf("/shipments/%d?error=Invalid+date+format", shipmentID)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}

	// Parse include accessories checkbox
	includeAccessories := r.FormValue("include_accessories") == "on" || r.FormValue("include_accessories") == "true"

	// Build validation input (using CompleteShipmentDetails validator - doesn't require JIRA/Company)
	formInput := validator.CompleteShipmentDetailsInput{
		ShipmentID:             shipmentID,
		ContactName:            r.FormValue("contact_name"),
		ContactEmail:           r.FormValue("contact_email"),
		ContactPhone:           r.FormValue("contact_phone"),
		PickupAddress:          r.FormValue("pickup_address"),
		PickupCity:             r.FormValue("pickup_city"),
		PickupState:            r.FormValue("pickup_state"),
		PickupZip:              r.FormValue("pickup_zip"),
		PickupDate:             pickupDateStr,
		PickupTimeSlot:         r.FormValue("pickup_time_slot"),
		SpecialInstructions:    r.FormValue("special_instructions"),
		LaptopSerialNumber:     r.FormValue("laptop_serial_number"),
		LaptopBrand:            r.FormValue("laptop_brand"),
		LaptopModel:            r.FormValue("laptop_model"),
		LaptopRAMGB:            r.FormValue("laptop_ram_gb"),
		LaptopSSDGB:            r.FormValue("laptop_ssd_gb"),
		EngineerName:           r.FormValue("engineer_name"),
		IncludeAccessories:     includeAccessories,
		AccessoriesDescription: r.FormValue("accessories_description"),
	}

	// Validate form
	if err := validator.ValidateCompleteShipmentDetails(formInput); err != nil {
		redirectURL := fmt.Sprintf("/shipments/%d?error=%s", shipmentID, err.Error())
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

	// Get shipment and company ID
	var companyID int64
	var shipmentType models.ShipmentType
	err = tx.QueryRowContext(r.Context(),
		`SELECT client_company_id, shipment_type FROM shipments WHERE id = $1`,
		shipmentID,
	).Scan(&companyID, &shipmentType)
	if err == sql.ErrNoRows {
		http.Error(w, "Shipment not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Failed to query shipment", http.StatusInternalServerError)
		return
	}

	// Create laptop record
	laptop := models.Laptop{
		SerialNumber:    formInput.LaptopSerialNumber,
		Brand:           formInput.LaptopBrand,
		Model:           formInput.LaptopModel,
		RAMGB:           formInput.LaptopRAMGB,
		SSDGB:           formInput.LaptopSSDGB,
		Status:          models.LaptopStatusInTransitToWarehouse,
		ClientCompanyID: &companyID,
	}
	laptop.BeforeCreate()

	var laptopID int64
	err = tx.QueryRowContext(r.Context(),
		`INSERT INTO laptops (serial_number, brand, model, ram_gb, ssd_gb, status, client_company_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`,
		laptop.SerialNumber, laptop.Brand, laptop.Model, laptop.RAMGB, laptop.SSDGB, laptop.Status, laptop.ClientCompanyID,
		laptop.CreatedAt, laptop.UpdatedAt,
	).Scan(&laptopID)
	if err != nil {
		// Log the actual database error for debugging
		fmt.Printf("Failed to create laptop: %v (serial: %s, model: %s, ram: %s, ssd: %s, status: %s, company: %v)\n",
			err, laptop.SerialNumber, laptop.Model, laptop.RAMGB, laptop.SSDGB, laptop.Status, laptop.ClientCompanyID)

		// Check if it's a duplicate serial number error
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			redirectURL := fmt.Sprintf("/shipments/%d?error=Laptop+serial+number+already+exists", shipmentID)
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return
		}

		http.Error(w, "Failed to create laptop", http.StatusInternalServerError)
		return
	}

	// Link laptop to shipment
	_, err = tx.ExecContext(r.Context(),
		`INSERT INTO shipment_laptops (shipment_id, laptop_id, created_at)
		VALUES ($1, $2, $3)`,
		shipmentID, laptopID, time.Now(),
	)
	if err != nil {
		http.Error(w, "Failed to link laptop to shipment", http.StatusInternalServerError)
		return
	}

	// Update shipment with pickup_scheduled_date
	_, err = tx.ExecContext(r.Context(),
		`UPDATE shipments SET pickup_scheduled_date = $1, updated_at = $2 WHERE id = $3`,
		pickupDate, time.Now(), shipmentID,
	)
	if err != nil {
		http.Error(w, "Failed to update shipment", http.StatusInternalServerError)
		return
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
		"special_instructions":    formInput.SpecialInstructions,
		"laptop_serial_number":    formInput.LaptopSerialNumber,
		"laptop_brand":            formInput.LaptopBrand,
		"laptop_model":            formInput.LaptopModel,
		"laptop_ram_gb":           formInput.LaptopRAMGB,
		"laptop_ssd_gb":           formInput.LaptopSSDGB,
		"engineer_name":           formInput.EngineerName,
		"include_accessories":     formInput.IncludeAccessories,
		"accessories_description": formInput.AccessoriesDescription,
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
		"action":       "shipment_details_completed",
		"shipment_id":  shipmentID,
		"company_id":   companyID,
		"laptop_id":    laptopID,
		"completed_by": user.Email,
	})

	_, err = tx.ExecContext(r.Context(),
		`INSERT INTO audit_logs (user_id, action, entity_type, entity_id, timestamp, details)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		user.ID, "shipment_details_completed", "shipment", shipmentID, time.Now(), auditDetails,
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

	// Redirect to shipment detail page with success message
	redirectURL := fmt.Sprintf("/shipments/%d?success=Shipment+details+completed+successfully", shipmentID)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

// EditShipmentDetails handles logistics users editing shipment details (except JIRA ticket and company)
func (h *PickupFormHandler) EditShipmentDetails(w http.ResponseWriter, r *http.Request) {
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

	// Only logistics users can edit shipment details
	if user.Role != models.RoleLogistics {
		http.Error(w, "Forbidden: Only logistics users can edit shipment details", http.StatusForbidden)
		return
	}

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Extract and validate shipment ID
	shipmentIDStr := r.FormValue("shipment_id")
	if shipmentIDStr == "" {
		http.Error(w, "Shipment ID is required", http.StatusBadRequest)
		return
	}
	shipmentID, err := strconv.ParseInt(shipmentIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid shipment ID", http.StatusBadRequest)
		return
	}

	// Parse pickup date
	pickupDateStr := r.FormValue("pickup_date")
	if pickupDateStr == "" {
		redirectURL := fmt.Sprintf("/shipments/%d?error=Pickup+date+is+required", shipmentID)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}
	pickupDate, err := time.Parse("2006-01-02", pickupDateStr)
	if err != nil {
		redirectURL := fmt.Sprintf("/shipments/%d?error=Invalid+date+format", shipmentID)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}

	// Parse include accessories checkbox
	includeAccessories := r.FormValue("include_accessories") == "on" || r.FormValue("include_accessories") == "true"

	// Build validation input for editing (laptop serial number is NOT included, will be preserved)
	formInput := validator.EditShipmentDetailsInput{
		ShipmentID:             shipmentID,
		ContactName:            r.FormValue("contact_name"),
		ContactEmail:           r.FormValue("contact_email"),
		ContactPhone:           r.FormValue("contact_phone"),
		PickupAddress:          r.FormValue("pickup_address"),
		PickupCity:             r.FormValue("pickup_city"),
		PickupState:            r.FormValue("pickup_state"),
		PickupZip:              r.FormValue("pickup_zip"),
		PickupDate:             pickupDateStr,
		PickupTimeSlot:         r.FormValue("pickup_time_slot"),
		SpecialInstructions:    r.FormValue("special_instructions"),
		LaptopModel:            r.FormValue("laptop_model"),
		LaptopRAMGB:            r.FormValue("laptop_ram_gb"),
		LaptopSSDGB:            r.FormValue("laptop_ssd_gb"),
		EngineerName:           r.FormValue("engineer_name"),
		IncludeAccessories:     includeAccessories,
		AccessoriesDescription: r.FormValue("accessories_description"),
	}

	// Validate form using edit validator (doesn't require laptop_serial_number)
	if err := validator.ValidateEditShipmentDetails(formInput); err != nil {
		redirectURL := fmt.Sprintf("/shipments/%d?error=%s", shipmentID, err.Error())
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

	// Verify shipment exists and get related IDs
	var companyID int64
	var jiraTicket string
	err = tx.QueryRowContext(r.Context(),
		`SELECT client_company_id, jira_ticket_number FROM shipments WHERE id = $1`,
		shipmentID,
	).Scan(&companyID, &jiraTicket)
	if err == sql.ErrNoRows {
		http.Error(w, "Shipment not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Failed to query shipment", http.StatusInternalServerError)
		return
	}

	// Get existing form data to preserve fields not being updated
	var existingFormDataJSON json.RawMessage
	err = tx.QueryRowContext(r.Context(),
		`SELECT form_data FROM pickup_forms WHERE shipment_id = $1`,
		shipmentID,
	).Scan(&existingFormDataJSON)
	if err != nil {
		http.Error(w, "Failed to query existing form data", http.StatusInternalServerError)
		return
	}

	// Parse existing form data
	var existingFormData map[string]interface{}
	err = json.Unmarshal(existingFormDataJSON, &existingFormData)
	if err != nil {
		http.Error(w, "Failed to parse existing form data", http.StatusInternalServerError)
		return
	}

	// Get laptop ID from shipment
	var laptopID int64
	err = tx.QueryRowContext(r.Context(),
		`SELECT laptop_id FROM shipment_laptops WHERE shipment_id = $1 LIMIT 1`,
		shipmentID,
	).Scan(&laptopID)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, "Failed to query laptop", http.StatusInternalServerError)
		return
	}

	// Update laptop details if laptop exists and new data provided
	if laptopID > 0 {
		updateQuery := `UPDATE laptops SET updated_at = $1`
		args := []interface{}{time.Now()}
		argIdx := 2

		if formInput.LaptopModel != "" {
			updateQuery += fmt.Sprintf(", model = $%d", argIdx)
			args = append(args, formInput.LaptopModel)
			argIdx++
		}
		if formInput.LaptopRAMGB != "" {
			updateQuery += fmt.Sprintf(", ram_gb = $%d", argIdx)
			args = append(args, formInput.LaptopRAMGB)
			argIdx++
		}
		if formInput.LaptopSSDGB != "" {
			updateQuery += fmt.Sprintf(", ssd_gb = $%d", argIdx)
			args = append(args, formInput.LaptopSSDGB)
			argIdx++
		}

		if argIdx > 2 { // Only execute if we have fields to update
			updateQuery += fmt.Sprintf(" WHERE id = $%d", argIdx)
			args = append(args, laptopID)

			_, err = tx.ExecContext(r.Context(), updateQuery, args...)
			if err != nil {
				http.Error(w, "Failed to update laptop", http.StatusInternalServerError)
				return
			}
		}
	}

	// Update shipment pickup_scheduled_date
	_, err = tx.ExecContext(r.Context(),
		`UPDATE shipments SET pickup_scheduled_date = $1, updated_at = $2 WHERE id = $3`,
		pickupDate, time.Now(), shipmentID,
	)
	if err != nil {
		http.Error(w, "Failed to update shipment", http.StatusInternalServerError)
		return
	}

	// Merge new form data with existing (preserving laptop_serial_number and other fields)
	updatedFormData := map[string]interface{}{
		"contact_name":            formInput.ContactName,
		"contact_email":           formInput.ContactEmail,
		"contact_phone":           formInput.ContactPhone,
		"pickup_address":          formInput.PickupAddress,
		"pickup_city":             formInput.PickupCity,
		"pickup_state":            formInput.PickupState,
		"pickup_zip":              formInput.PickupZip,
		"pickup_date":             formInput.PickupDate,
		"pickup_time_slot":        formInput.PickupTimeSlot,
		"special_instructions":    formInput.SpecialInstructions,
		"laptop_model":            formInput.LaptopModel,
		"laptop_ram_gb":           formInput.LaptopRAMGB,
		"laptop_ssd_gb":           formInput.LaptopSSDGB,
		"engineer_name":           formInput.EngineerName,
		"include_accessories":     formInput.IncludeAccessories,
		"accessories_description": formInput.AccessoriesDescription,
		"jira_ticket_number":      jiraTicket,                               // Preserved from shipment
		"laptop_serial_number":    existingFormData["laptop_serial_number"], // Preserved from existing form
	}

	// Update pickup form with merged data
	formDataJSON, err := json.Marshal(updatedFormData)
	if err != nil {
		http.Error(w, "Failed to encode form data", http.StatusInternalServerError)
		return
	}

	_, err = tx.ExecContext(r.Context(),
		`UPDATE pickup_forms SET form_data = $1 WHERE shipment_id = $2`,
		formDataJSON, shipmentID,
	)
	if err != nil {
		http.Error(w, "Failed to update pickup form", http.StatusInternalServerError)
		return
	}

	// Create audit log entry
	auditDetails, _ := json.Marshal(map[string]interface{}{
		"action":      "shipment_details_edited",
		"shipment_id": shipmentID,
		"edited_by":   user.Email,
	})

	_, err = tx.ExecContext(r.Context(),
		`INSERT INTO audit_logs (user_id, action, entity_type, entity_id, timestamp, details)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		user.ID, "shipment_details_edited", "shipment", shipmentID, time.Now(), auditDetails,
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

	// Redirect to shipment detail page with success message
	redirectURL := fmt.Sprintf("/shipments/%d?success=Shipment+details+updated+successfully", shipmentID)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}
