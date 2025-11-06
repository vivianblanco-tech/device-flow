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

	"github.com/gorilla/mux"
	"github.com/lib/pq"

	"github.com/yourusername/laptop-tracking-system/internal/middleware"
	"github.com/yourusername/laptop-tracking-system/internal/models"
	"github.com/yourusername/laptop-tracking-system/internal/views"
)

// ShipmentsHandler handles shipment-related requests
type ShipmentsHandler struct {
	DB            *sql.DB
	Templates     *template.Template
	JiraValidator models.JiraTicketValidator
}

// NewShipmentsHandler creates a new ShipmentsHandler
func NewShipmentsHandler(db *sql.DB, templates *template.Template) *ShipmentsHandler {
	return &ShipmentsHandler{
		DB:        db,
		Templates: templates,
	}
}

// ShipmentsList displays a list of shipments
func (h *ShipmentsHandler) ShipmentsList(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Get filter parameters
	statusFilter := r.URL.Query().Get("status")
	searchQuery := r.URL.Query().Get("search")

	// Build query based on user role
	var query string
	var args []interface{}
	argCount := 1

	baseQuery := `
		SELECT s.id, s.client_company_id, s.software_engineer_id, s.status, 
		       s.jira_ticket_number, s.courier_name, s.tracking_number, s.pickup_scheduled_date,
		       s.picked_up_at, s.arrived_warehouse_at, s.released_warehouse_at, 
		       s.delivered_at, s.notes, s.created_at, s.updated_at,
		       c.name as company_name,
		       se.name as engineer_name
		FROM shipments s
		JOIN client_companies c ON c.id = s.client_company_id
		LEFT JOIN software_engineers se ON se.id = s.software_engineer_id
		WHERE 1=1
	`

	// Role-based filtering
	switch user.Role {
	case models.RoleClient:
		// Clients can only see their own company's shipments
		// Note: In a real app, we'd link user to company via a relationship
		// For now, we skip this filter and let clients see all shipments
		// TODO: Add company_id to users table and filter properly
	case models.RoleWarehouse:
		// Warehouse users see shipments in transit or at warehouse
		baseQuery += " AND s.status IN ('in_transit_to_warehouse', 'at_warehouse', 'released_from_warehouse')"
	case models.RoleLogistics, models.RoleProjectManager:
		// Logistics and PM users can see all shipments - no additional filter needed
	}

	// Status filter
	if statusFilter != "" {
		baseQuery += fmt.Sprintf(" AND s.status = $%d", argCount)
		args = append(args, statusFilter)
		argCount++
	}

	// Search filter (by tracking number or company name)
	if searchQuery != "" {
		baseQuery += fmt.Sprintf(" AND (s.tracking_number ILIKE $%d OR c.name ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+searchQuery+"%")
		argCount++
	}

	// Order by most recent first
	query = baseQuery + " ORDER BY s.created_at DESC LIMIT 100"

	// Execute query
	rows, err := h.DB.QueryContext(r.Context(), query, args...)
	if err != nil {
		http.Error(w, "Failed to load shipments", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	shipments := []map[string]interface{}{}
	for rows.Next() {
		var s models.Shipment
		var companyName string
		var engineerName sql.NullString
		var jiraTicket sql.NullString
		var courierName sql.NullString
		var trackingNumber sql.NullString
		var notes sql.NullString

		err := rows.Scan(
			&s.ID, &s.ClientCompanyID, &s.SoftwareEngineerID, &s.Status,
			&jiraTicket, &courierName, &trackingNumber, &s.PickupScheduledDate,
			&s.PickedUpAt, &s.ArrivedWarehouseAt, &s.ReleasedWarehouseAt,
			&s.DeliveredAt, &notes, &s.CreatedAt, &s.UpdatedAt,
			&companyName, &engineerName,
		)
		if err != nil {
			continue
		}

		// Convert nullable strings
		s.JiraTicketNumber = jiraTicket.String
		s.CourierName = courierName.String
		s.TrackingNumber = trackingNumber.String
		s.Notes = notes.String

		shipment := map[string]interface{}{
			"Shipment":     s,
			"CompanyName":  companyName,
			"EngineerName": engineerName.String,
		}
		shipments = append(shipments, shipment)
	}

	// Get error and success messages
	errorMsg := r.URL.Query().Get("error")
	successMsg := r.URL.Query().Get("success")

	data := map[string]interface{}{
		"Error":        errorMsg,
		"Success":      successMsg,
		"User":         user,
		"Nav":          views.GetNavigationLinks(user.Role),
		"CurrentPage":  "shipments",
		"Shipments":    shipments,
		"StatusFilter": statusFilter,
		"SearchQuery":  searchQuery,
		"AllStatuses": []models.ShipmentStatus{
			models.ShipmentStatusPendingPickup,
			models.ShipmentStatusPickedUpFromClient,
			models.ShipmentStatusInTransitToWarehouse,
			models.ShipmentStatusAtWarehouse,
			models.ShipmentStatusReleasedFromWarehouse,
			models.ShipmentStatusInTransitToEngineer,
			models.ShipmentStatusDelivered,
		},
	}

	if h.Templates != nil {
		err := h.Templates.ExecuteTemplate(w, "shipments-list.html", data)
		if err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			return
		}
	} else {
		// For testing without templates
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Shipments List Page")
	}
}

// ShipmentDetail displays detailed information about a shipment
func (h *ShipmentsHandler) ShipmentDetail(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Get shipment ID from URL path variable
	vars := mux.Vars(r)
	shipmentIDStr := vars["id"]
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
	var s models.Shipment
	var companyName string
	var engineerName sql.NullString
	var engineerEmail sql.NullString

	err = h.DB.QueryRowContext(r.Context(),
		`SELECT s.id, s.client_company_id, s.software_engineer_id, s.status, 
		        COALESCE(s.jira_ticket_number, '') as jira_ticket_number,
		        COALESCE(s.courier_name, '') as courier_name, 
		        COALESCE(s.tracking_number, '') as tracking_number, 
		        s.pickup_scheduled_date,
		        s.picked_up_at, s.arrived_warehouse_at, s.released_warehouse_at, 
		        s.delivered_at, COALESCE(s.notes, '') as notes, 
		        s.created_at, s.updated_at,
		        c.name, se.name, se.email
		FROM shipments s
		JOIN client_companies c ON c.id = s.client_company_id
		LEFT JOIN software_engineers se ON se.id = s.software_engineer_id
		WHERE s.id = $1`,
		shipmentID,
	).Scan(
		&s.ID, &s.ClientCompanyID, &s.SoftwareEngineerID, &s.Status,
		&s.JiraTicketNumber, &s.CourierName, &s.TrackingNumber, &s.PickupScheduledDate,
		&s.PickedUpAt, &s.ArrivedWarehouseAt, &s.ReleasedWarehouseAt,
		&s.DeliveredAt, &s.Notes, &s.CreatedAt, &s.UpdatedAt,
		&companyName, &engineerName, &engineerEmail,
	)

	if err == sql.ErrNoRows {
		http.Error(w, "Shipment not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed to load shipment", http.StatusInternalServerError)
		return
	}

	// Get associated laptops
	laptopRows, err := h.DB.QueryContext(r.Context(),
		`SELECT l.id, l.serial_number, l.brand, l.model, l.specs, l.status, l.created_at
		FROM laptops l
		JOIN shipment_laptops sl ON sl.laptop_id = l.id
		WHERE sl.shipment_id = $1`,
		shipmentID,
	)
	if err != nil {
		http.Error(w, "Failed to load laptops", http.StatusInternalServerError)
		return
	}
	defer laptopRows.Close()

	laptops := []models.Laptop{}
	for laptopRows.Next() {
		var laptop models.Laptop
		err := laptopRows.Scan(
			&laptop.ID, &laptop.SerialNumber, &laptop.Brand, &laptop.Model,
			&laptop.Specs, &laptop.Status, &laptop.CreatedAt,
		)
		if err != nil {
			continue
		}
		laptops = append(laptops, laptop)
	}

	// Get pickup form if exists
	var pickupForm *models.PickupForm
	var pickupFormData json.RawMessage
	pickupFormTemp := models.PickupForm{}
	err = h.DB.QueryRowContext(r.Context(),
		`SELECT id, shipment_id, submitted_by_user_id, submitted_at, form_data
		FROM pickup_forms WHERE shipment_id = $1`,
		shipmentID,
	).Scan(&pickupFormTemp.ID, &pickupFormTemp.ShipmentID, &pickupFormTemp.SubmittedByUserID,
		&pickupFormTemp.SubmittedAt, &pickupFormData)
	if err == nil {
		pickupForm = &pickupFormTemp
	} else if err != sql.ErrNoRows {
		// Non-critical error, log it but continue
		fmt.Printf("Error fetching pickup form: %v\n", err)
	}

	// Get reception report if exists
	var receptionReport *models.ReceptionReport
	receptionReportTemp := models.ReceptionReport{}
	err = h.DB.QueryRowContext(r.Context(),
		`SELECT id, shipment_id, warehouse_user_id, received_at, notes, photo_urls
		FROM reception_reports WHERE shipment_id = $1`,
		shipmentID,
	).Scan(&receptionReportTemp.ID, &receptionReportTemp.ShipmentID, &receptionReportTemp.WarehouseUserID,
		&receptionReportTemp.ReceivedAt, &receptionReportTemp.Notes, (*pq.StringArray)(&receptionReportTemp.PhotoURLs))
	if err == nil {
		receptionReport = &receptionReportTemp
	} else if err != sql.ErrNoRows {
		// Non-critical error, log it but continue
		fmt.Printf("Error fetching reception report: %v\n", err)
	}

	// Get delivery form if exists
	var deliveryForm *models.DeliveryForm
	deliveryFormTemp := models.DeliveryForm{}
	err = h.DB.QueryRowContext(r.Context(),
		`SELECT id, shipment_id, engineer_id, delivered_at, notes, photo_urls
		FROM delivery_forms WHERE shipment_id = $1`,
		shipmentID,
	).Scan(&deliveryFormTemp.ID, &deliveryFormTemp.ShipmentID, &deliveryFormTemp.EngineerID,
		&deliveryFormTemp.DeliveredAt, &deliveryFormTemp.Notes, (*pq.StringArray)(&deliveryFormTemp.PhotoURLs))
	if err == nil {
		deliveryForm = &deliveryFormTemp
	} else if err != sql.ErrNoRows {
		// Non-critical error, log it but continue
		fmt.Printf("Error fetching delivery form: %v\n", err)
	}

	// Get list of software engineers (for assignment)
	engineerRows, err := h.DB.QueryContext(r.Context(),
		`SELECT id, name, email FROM software_engineers ORDER BY name`,
	)
	engineers := []models.SoftwareEngineer{}
	if err == nil {
		defer engineerRows.Close()
		for engineerRows.Next() {
			var engineer models.SoftwareEngineer
			if err := engineerRows.Scan(&engineer.ID, &engineer.Name, &engineer.Email); err == nil {
				engineers = append(engineers, engineer)
			}
		}
	}

	// Get error and success messages
	errorMsg := r.URL.Query().Get("error")
	successMsg := r.URL.Query().Get("success")

	// Generate tracking URL if courier and tracking number are present
	trackingURL := s.GetTrackingURL()

	data := map[string]interface{}{
		"Error":           errorMsg,
		"Success":         successMsg,
		"User":            user,
		"Nav":             views.GetNavigationLinks(user.Role),
		"CurrentPage":     "shipments",
		"Shipment":        s,
		"TrackingURL":     trackingURL,
		"CompanyName":     companyName,
		"EngineerName":    engineerName.String,
		"EngineerEmail":   engineerEmail.String,
		"Laptops":         laptops,
		"PickupForm":      pickupForm,
		"ReceptionReport": receptionReport,
		"DeliveryForm":    deliveryForm,
		"Engineers":       engineers,
	}

	if h.Templates != nil {
		err := h.Templates.ExecuteTemplate(w, "shipment-detail.html", data)
		if err != nil {
			fmt.Printf("Template execution error: %v\n", err)
			http.Error(w, fmt.Sprintf("Failed to render template: %v", err), http.StatusInternalServerError)
			return
		}
	} else {
		// For testing without templates
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Shipment Detail Page")
	}
}

// UpdateShipmentStatus updates the shipment status (logistics only)
func (h *ShipmentsHandler) UpdateShipmentStatus(w http.ResponseWriter, r *http.Request) {
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

	// Only logistics users can manually update status
	if user.Role != models.RoleLogistics {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	shipmentIDStr := r.FormValue("shipment_id")
	shipmentID, err := strconv.ParseInt(shipmentIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid shipment ID", http.StatusBadRequest)
		return
	}

	newStatus := models.ShipmentStatus(r.FormValue("status"))
	if !models.IsValidShipmentStatus(newStatus) {
		http.Error(w, "Invalid status", http.StatusBadRequest)
		return
	}

	// Update shipment status
	var shipment models.Shipment
	shipment.Status = newStatus
	shipment.UpdateStatus(newStatus)

	_, err = h.DB.ExecContext(r.Context(),
		`UPDATE shipments 
		SET status = $1, updated_at = $2,
		    picked_up_at = COALESCE($3, picked_up_at),
		    arrived_warehouse_at = COALESCE($4, arrived_warehouse_at),
		    released_warehouse_at = COALESCE($5, released_warehouse_at),
		    delivered_at = COALESCE($6, delivered_at)
		WHERE id = $7`,
		shipment.Status, shipment.UpdatedAt,
		shipment.PickedUpAt, shipment.ArrivedWarehouseAt,
		shipment.ReleasedWarehouseAt, shipment.DeliveredAt,
		shipmentID,
	)
	if err != nil {
		http.Error(w, "Failed to update shipment status", http.StatusInternalServerError)
		return
	}

	// Create audit log
	auditDetails, _ := json.Marshal(map[string]interface{}{
		"action":     "status_updated",
		"new_status": newStatus,
	})

	_, err = h.DB.ExecContext(r.Context(),
		`INSERT INTO audit_logs (user_id, action, entity_type, entity_id, timestamp, details)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		user.ID, "status_updated", "shipment", shipmentID, time.Now(), auditDetails,
	)
	if err != nil {
		// Non-critical error
		fmt.Printf("Failed to create audit log: %v\n", err)
	}

	// Redirect back to shipment detail
	redirectURL := fmt.Sprintf("/shipments/%d?success=Status+updated+successfully", shipmentID)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

// AssignEngineer assigns a software engineer to a shipment (logistics only)
func (h *ShipmentsHandler) AssignEngineer(w http.ResponseWriter, r *http.Request) {
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

	// Only logistics users can assign engineers
	if user.Role != models.RoleLogistics {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Get shipment ID from URL path variable
	vars := mux.Vars(r)
	shipmentIDStr := vars["id"]
	shipmentID, err := strconv.ParseInt(shipmentIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid shipment ID", http.StatusBadRequest)
		return
	}

	// Parse form data
	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	engineerIDStr := r.FormValue("engineer_id")
	if engineerIDStr == "" {
		redirectURL := fmt.Sprintf("/shipments/%d?error=Please+select+an+engineer", shipmentID)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}

	engineerID, err := strconv.ParseInt(engineerIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid engineer ID", http.StatusBadRequest)
		return
	}

	// Update shipment with engineer assignment
	_, err = h.DB.ExecContext(r.Context(),
		`UPDATE shipments 
		SET software_engineer_id = $1, updated_at = $2
		WHERE id = $3`,
		engineerID, time.Now(), shipmentID,
	)
	if err != nil {
		fmt.Printf("Error assigning engineer: %v\n", err)
		http.Error(w, "Failed to assign engineer", http.StatusInternalServerError)
		return
	}

	// Create audit log
	auditDetails, _ := json.Marshal(map[string]interface{}{
		"action":      "engineer_assigned",
		"engineer_id": engineerID,
	})

	_, err = h.DB.ExecContext(r.Context(),
		`INSERT INTO audit_logs (user_id, action, entity_type, entity_id, timestamp, details)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		user.ID, "engineer_assigned", "shipment", shipmentID, time.Now(), auditDetails,
	)
	if err != nil {
		// Non-critical error
		fmt.Printf("Failed to create audit log: %v\n", err)
	}

	// Redirect back to shipment detail
	redirectURL := fmt.Sprintf("/shipments/%d?success=Engineer+assigned+successfully", shipmentID)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

// CreateShipment handles creating a new shipment (GET shows form, POST creates shipment)
func (h *ShipmentsHandler) CreateShipment(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Only logistics users can create shipments
	if user.Role != models.RoleLogistics {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// GET request - show create shipment form
	if r.Method == http.MethodGet {
		// Get list of client companies for the dropdown
		rows, err := h.DB.QueryContext(r.Context(),
			`SELECT id, name FROM client_companies ORDER BY name`,
		)
		if err != nil {
			http.Error(w, "Failed to load companies", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var companies []struct {
			ID   int64
			Name string
		}
		for rows.Next() {
			var company struct {
				ID   int64
				Name string
			}
			if err := rows.Scan(&company.ID, &company.Name); err != nil {
				http.Error(w, "Failed to read companies", http.StatusInternalServerError)
				return
			}
			companies = append(companies, company)
		}

		data := map[string]interface{}{
			"User":        user,
			"Nav":         views.GetNavigationLinks(user.Role),
			"CurrentPage": "shipments",
			"Companies":   companies,
		}

		err = h.Templates.ExecuteTemplate(w, "create-shipment.html", data)
		if err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			return
		}
		return
	}

	// POST request - create shipment
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Get form fields
	clientCompanyIDStr := r.FormValue("client_company_id")
	jiraTicketNumber := strings.TrimSpace(r.FormValue("jira_ticket_number"))
	notes := r.FormValue("notes")

	// Validate client company ID
	if clientCompanyIDStr == "" {
		http.Error(w, "Client company is required", http.StatusBadRequest)
		return
	}
	clientCompanyID, err := strconv.ParseInt(clientCompanyIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid client company ID", http.StatusBadRequest)
		return
	}

	// Create shipment model for validation
	shipment := models.Shipment{
		ClientCompanyID:  clientCompanyID,
		Status:           models.ShipmentStatusPendingPickup,
		JiraTicketNumber: jiraTicketNumber,
		Notes:            notes,
	}

	// Validate shipment (includes JIRA ticket format validation)
	if err := shipment.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate JIRA ticket exists (if validator is configured)
	if err := models.ValidateJiraTicketExists(jiraTicketNumber, h.JiraValidator); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set timestamps
	shipment.BeforeCreate()

	// Insert shipment into database
	var shipmentID int64
	err = h.DB.QueryRowContext(r.Context(),
		`INSERT INTO shipments (client_company_id, status, jira_ticket_number, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`,
		shipment.ClientCompanyID, shipment.Status, shipment.JiraTicketNumber,
		shipment.Notes, shipment.CreatedAt, shipment.UpdatedAt,
	).Scan(&shipmentID)
	if err != nil {
		fmt.Printf("Error creating shipment: %v\n", err)
		http.Error(w, "Failed to create shipment", http.StatusInternalServerError)
		return
	}

	// Create audit log
	auditDetails, _ := json.Marshal(map[string]interface{}{
		"action":             "shipment_created",
		"jira_ticket_number": jiraTicketNumber,
		"client_company_id":  clientCompanyID,
	})

	_, err = h.DB.ExecContext(r.Context(),
		`INSERT INTO audit_logs (user_id, action, entity_type, entity_id, timestamp, details)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		user.ID, "shipment_created", "shipment", shipmentID, time.Now(), auditDetails,
	)
	if err != nil {
		// Non-critical error
		fmt.Printf("Failed to create audit log: %v\n", err)
	}

	// Redirect to shipment detail page
	redirectURL := fmt.Sprintf("/shipments/%d?success=Shipment+created+successfully", shipmentID)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

// ShipmentPickupFormPage displays the pickup form for a specific shipment
func (h *ShipmentsHandler) ShipmentPickupFormPage(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Get shipment ID from URL
	vars := mux.Vars(r)
	shipmentIDStr := vars["id"]
	shipmentID, err := strconv.ParseInt(shipmentIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid shipment ID", http.StatusBadRequest)
		return
	}

	// Get shipment with JIRA ticket and company information
	var shipment models.Shipment
	var companyName string
	var pickupScheduledDate sql.NullTime
	var notes sql.NullString
	err = h.DB.QueryRowContext(r.Context(),
		`SELECT s.id, s.client_company_id, s.status, s.jira_ticket_number, 
		        s.pickup_scheduled_date, s.notes, s.created_at, s.updated_at, c.name
		FROM shipments s
		JOIN client_companies c ON c.id = s.client_company_id
		WHERE s.id = $1`,
		shipmentID,
	).Scan(&shipment.ID, &shipment.ClientCompanyID, &shipment.Status, &shipment.JiraTicketNumber,
		&pickupScheduledDate, &notes, &shipment.CreatedAt, &shipment.UpdatedAt, &companyName)

	if err == sql.ErrNoRows {
		http.Error(w, "Shipment not found", http.StatusNotFound)
		return
	}
	if err != nil {
		fmt.Printf("Error loading shipment: %v\n", err)
		http.Error(w, "Failed to load shipment", http.StatusInternalServerError)
		return
	}

	// Handle nullable fields
	if pickupScheduledDate.Valid {
		shipment.PickupScheduledDate = &pickupScheduledDate.Time
	}
	if notes.Valid {
		shipment.Notes = notes.String
	}

	// Check if pickup form already exists
	var pickupFormData map[string]interface{}
	var formDataJSON []byte
	var pickupFormID int64
	err = h.DB.QueryRowContext(r.Context(),
		`SELECT id, form_data FROM pickup_forms WHERE shipment_id = $1`,
		shipmentID,
	).Scan(&pickupFormID, &formDataJSON)

	if err == nil {
		// Pickup form exists, parse the JSON
		if err := json.Unmarshal(formDataJSON, &pickupFormData); err != nil {
			// Log error but continue
			fmt.Printf("Error parsing pickup form data: %v\n", err)
		}
	} else if err != sql.ErrNoRows {
		// Real error (not just missing form)
		http.Error(w, "Failed to load pickup form", http.StatusInternalServerError)
		return
	}

	// Prepare template data
	data := map[string]interface{}{
		"User":           user,
		"Shipment":       shipment,
		"CompanyName":    companyName,
		"PickupFormData": pickupFormData,
		"IsEdit":         pickupFormData != nil,
		"TimeSlots":      []string{"morning", "afternoon", "evening"},
	}

	if h.Templates != nil {
		err := h.Templates.ExecuteTemplate(w, "shipment-pickup-form.html", data)
		if err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			return
		}
	} else {
		// For testing without templates
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Shipment Pickup Form Page")
	}
}

// ShipmentPickupFormSubmit handles pickup form submission for a specific shipment
func (h *ShipmentsHandler) ShipmentPickupFormSubmit(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Get shipment ID from URL
	vars := mux.Vars(r)
	shipmentIDStr := vars["id"]
	shipmentID, err := strconv.ParseInt(shipmentIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid shipment ID", http.StatusBadRequest)
		return
	}

	// Parse form data
	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Get form fields
	contactName := r.FormValue("contact_name")
	contactEmail := r.FormValue("contact_email")
	contactPhone := r.FormValue("contact_phone")
	pickupAddress := r.FormValue("pickup_address")
	pickupDate := r.FormValue("pickup_date")
	pickupTimeSlot := r.FormValue("pickup_time_slot")
	numberOfLaptops := r.FormValue("number_of_laptops")
	specialInstructions := r.FormValue("special_instructions")

	// Validate required fields
	if contactName == "" || contactEmail == "" || contactPhone == "" || 
	   pickupAddress == "" || pickupDate == "" || pickupTimeSlot == "" || 
	   numberOfLaptops == "" {
		http.Error(w, "All required fields must be filled", http.StatusBadRequest)
		return
	}

	// Verify shipment exists
	var existingShipmentID int64
	err = h.DB.QueryRowContext(r.Context(),
		`SELECT id FROM shipments WHERE id = $1`,
		shipmentID,
	).Scan(&existingShipmentID)
	if err == sql.ErrNoRows {
		http.Error(w, "Shipment not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed to verify shipment", http.StatusInternalServerError)
		return
	}

	// Parse pickup date
	pickupDateTime, err := time.Parse("2006-01-02", pickupDate)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	// Create form data JSON
	formData := map[string]interface{}{
		"contact_name":         contactName,
		"contact_email":        contactEmail,
		"contact_phone":        contactPhone,
		"pickup_address":       pickupAddress,
		"pickup_date":          pickupDate,
		"pickup_time_slot":     pickupTimeSlot,
		"number_of_laptops":    numberOfLaptops,
		"special_instructions": specialInstructions,
	}
	formDataJSON, err := json.Marshal(formData)
	if err != nil {
		http.Error(w, "Failed to encode form data", http.StatusInternalServerError)
		return
	}

	// Check if pickup form already exists
	var existingFormID int64
	err = h.DB.QueryRowContext(r.Context(),
		`SELECT id FROM pickup_forms WHERE shipment_id = $1`,
		shipmentID,
	).Scan(&existingFormID)

	if err == sql.ErrNoRows {
		// Create new pickup form
		_, err = h.DB.ExecContext(r.Context(),
			`INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data)
			VALUES ($1, $2, $3, $4)`,
			shipmentID, user.ID, time.Now(), formDataJSON,
		)
		if err != nil {
			http.Error(w, "Failed to create pickup form", http.StatusInternalServerError)
			return
		}

		// Update shipment pickup_scheduled_date
		_, err = h.DB.ExecContext(r.Context(),
			`UPDATE shipments SET pickup_scheduled_date = $1, updated_at = $2 WHERE id = $3`,
			pickupDateTime, time.Now(), shipmentID,
		)
		if err != nil {
			// Non-critical, log and continue
			fmt.Printf("Warning: Failed to update shipment pickup date: %v\n", err)
		}
	} else if err == nil {
		// Update existing pickup form
		_, err = h.DB.ExecContext(r.Context(),
			`UPDATE pickup_forms SET form_data = $1, submitted_at = $2, submitted_by_user_id = $3
			WHERE id = $4`,
			formDataJSON, time.Now(), user.ID, existingFormID,
		)
		if err != nil {
			http.Error(w, "Failed to update pickup form", http.StatusInternalServerError)
			return
		}

		// Update shipment pickup_scheduled_date
		_, err = h.DB.ExecContext(r.Context(),
			`UPDATE shipments SET pickup_scheduled_date = $1, updated_at = $2 WHERE id = $3`,
			pickupDateTime, time.Now(), shipmentID,
		)
		if err != nil {
			// Non-critical, log and continue
			fmt.Printf("Warning: Failed to update shipment pickup date: %v\n", err)
		}
	} else {
		http.Error(w, "Failed to check existing form", http.StatusInternalServerError)
		return
	}

	// Redirect to shipment detail page with success message
	redirectURL := fmt.Sprintf("/shipments/%d?success=Pickup+form+submitted+successfully", shipmentID)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}