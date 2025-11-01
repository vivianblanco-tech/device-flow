package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/yourusername/laptop-tracking-system/internal/middleware"
	"github.com/yourusername/laptop-tracking-system/internal/models"
)

// ShipmentsHandler handles shipment-related requests
type ShipmentsHandler struct {
	DB        *sql.DB
	Templates *template.Template
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
		       s.courier_name, s.tracking_number, s.pickup_scheduled_date,
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
		baseQuery += fmt.Sprintf(" AND s.client_company_id = $%d", argCount)
		// Note: In a real app, we'd link user to company via a relationship
		// For now, we'll skip this filter if user doesn't have company_id
		argCount++
	case models.RoleWarehouse:
		// Warehouse users see shipments in transit or at warehouse
		baseQuery += " AND s.status IN ('in_transit_to_warehouse', 'at_warehouse', 'released_from_warehouse')"
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

		err := rows.Scan(
			&s.ID, &s.ClientCompanyID, &s.SoftwareEngineerID, &s.Status,
			&s.CourierName, &s.TrackingNumber, &s.PickupScheduledDate,
			&s.PickedUpAt, &s.ArrivedWarehouseAt, &s.ReleasedWarehouseAt,
			&s.DeliveredAt, &s.Notes, &s.CreatedAt, &s.UpdatedAt,
			&companyName, &engineerName,
		)
		if err != nil {
			continue
		}

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

	// Get shipment ID from URL path (you'll need to parse this based on your routing)
	shipmentIDStr := r.URL.Query().Get("id")
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
		        s.courier_name, s.tracking_number, s.pickup_scheduled_date,
		        s.picked_up_at, s.arrived_warehouse_at, s.released_warehouse_at, 
		        s.delivered_at, s.notes, s.created_at, s.updated_at,
		        c.name, se.name, se.email
		FROM shipments s
		JOIN client_companies c ON c.id = s.client_company_id
		LEFT JOIN software_engineers se ON se.id = s.software_engineer_id
		WHERE s.id = $1`,
		shipmentID,
	).Scan(
		&s.ID, &s.ClientCompanyID, &s.SoftwareEngineerID, &s.Status,
		&s.CourierName, &s.TrackingNumber, &s.PickupScheduledDate,
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
	err = h.DB.QueryRowContext(r.Context(),
		`SELECT id, shipment_id, submitted_by_user_id, submitted_at, form_data
		FROM pickup_forms WHERE shipment_id = $1`,
		shipmentID,
	).Scan(&pickupForm.ID, &pickupForm.ShipmentID, &pickupForm.SubmittedByUserID,
		&pickupForm.SubmittedAt, &pickupFormData)
	if err != nil && err != sql.ErrNoRows {
		// Non-critical error
		pickupForm = nil
	}

	// Get reception report if exists
	var receptionReport *models.ReceptionReport
	err = h.DB.QueryRowContext(r.Context(),
		`SELECT id, shipment_id, warehouse_user_id, received_at, notes, photo_urls
		FROM reception_reports WHERE shipment_id = $1`,
		shipmentID,
	).Scan(&receptionReport.ID, &receptionReport.ShipmentID, &receptionReport.WarehouseUserID,
		&receptionReport.ReceivedAt, &receptionReport.Notes, &receptionReport.PhotoURLs)
	if err != nil && err != sql.ErrNoRows {
		// Non-critical error
		receptionReport = nil
	}

	// Get delivery form if exists
	var deliveryForm *models.DeliveryForm
	err = h.DB.QueryRowContext(r.Context(),
		`SELECT id, shipment_id, engineer_id, delivered_at, notes, photo_urls
		FROM delivery_forms WHERE shipment_id = $1`,
		shipmentID,
	).Scan(&deliveryForm.ID, &deliveryForm.ShipmentID, &deliveryForm.EngineerID,
		&deliveryForm.DeliveredAt, &deliveryForm.Notes, &deliveryForm.PhotoURLs)
	if err != nil && err != sql.ErrNoRows {
		// Non-critical error
		deliveryForm = nil
	}

	// Get error and success messages
	errorMsg := r.URL.Query().Get("error")
	successMsg := r.URL.Query().Get("success")

	data := map[string]interface{}{
		"Error":           errorMsg,
		"Success":         successMsg,
		"User":            user,
		"Shipment":        s,
		"CompanyName":     companyName,
		"EngineerName":    engineerName.String,
		"EngineerEmail":   engineerEmail.String,
		"Laptops":         laptops,
		"PickupForm":      pickupForm,
		"ReceptionReport": receptionReport,
		"DeliveryForm":    deliveryForm,
	}

	if h.Templates != nil {
		err := h.Templates.ExecuteTemplate(w, "shipment-detail.html", data)
		if err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
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

