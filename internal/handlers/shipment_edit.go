package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/yourusername/laptop-tracking-system/internal/middleware"
	"github.com/yourusername/laptop-tracking-system/internal/models"
	"github.com/yourusername/laptop-tracking-system/internal/views"
)

// EditShipmentGET displays the edit shipment form (logistics only)
func (h *ShipmentsHandler) EditShipmentGET(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Only logistics users can edit shipments
	if user.Role != models.RoleLogistics {
		http.Error(w, "Forbidden", http.StatusForbidden)
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
	var engineerID sql.NullInt64
	var engineerName sql.NullString

	err = h.DB.QueryRowContext(r.Context(),
		`SELECT s.id, s.shipment_type, s.laptop_count, s.client_company_id, s.software_engineer_id, s.status, 
		        COALESCE(s.jira_ticket_number, '') as jira_ticket_number,
		        COALESCE(s.courier_name, '') as courier_name, 
		        COALESCE(s.tracking_number, '') as tracking_number,
		        COALESCE(s.second_tracking_number, '') as second_tracking_number,
		        s.pickup_scheduled_date,
		        s.picked_up_at, s.arrived_warehouse_at, s.released_warehouse_at, 
		        s.eta_to_engineer, s.delivered_at, COALESCE(s.notes, '') as notes, 
		        s.created_at, s.updated_at,
		        c.name, se.id, se.name
		FROM shipments s
		JOIN client_companies c ON c.id = s.client_company_id
		LEFT JOIN software_engineers se ON se.id = s.software_engineer_id
		WHERE s.id = $1`,
		shipmentID,
	).Scan(
		&s.ID, &s.ShipmentType, &s.LaptopCount, &s.ClientCompanyID, &s.SoftwareEngineerID, &s.Status,
		&s.JiraTicketNumber, &s.CourierName, &s.TrackingNumber, &s.SecondTrackingNumber, &s.PickupScheduledDate,
		&s.PickedUpAt, &s.ArrivedWarehouseAt, &s.ReleasedWarehouseAt,
		&s.ETAToEngineer, &s.DeliveredAt, &s.Notes, &s.CreatedAt, &s.UpdatedAt,
		&companyName, &engineerID, &engineerName,
	)

	if err == sql.ErrNoRows {
		http.Error(w, "Shipment not found", http.StatusNotFound)
		return
	}
	if err != nil {
		fmt.Printf("Error loading shipment: %v\n", err)
		http.Error(w, "Failed to load shipment", http.StatusInternalServerError)
		return
	}

	// Check edit availability based on shipment type and status
	canEdit, errorMsg := canEditShipment(&s, h.DB, r.Context(), shipmentID)
	if !canEdit {
		http.Error(w, errorMsg, http.StatusBadRequest)
		return
	}

	// Get pickup form if exists
	var pickupFormData map[string]interface{}
	var formDataJSON json.RawMessage
	err = h.DB.QueryRowContext(r.Context(),
		`SELECT form_data FROM pickup_forms WHERE shipment_id = $1`,
		shipmentID,
	).Scan(&formDataJSON)
	
	if err == nil {
		// Parse the JSONB form_data into a map for template use
		if err := json.Unmarshal(formDataJSON, &pickupFormData); err != nil {
			fmt.Printf("Error parsing pickup form data: %v\n", err)
			pickupFormData = nil
		}
	} else if err != sql.ErrNoRows {
		// Non-critical error, log it but continue
		fmt.Printf("Error fetching pickup form: %v\n", err)
	}

	// Get list of software engineers
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

	data := map[string]interface{}{
		"User":           user,
		"Nav":            views.GetNavigationLinks(user.Role),
		"CurrentPage":    "shipments",
		"Shipment":       s,
		"CompanyName":    companyName,
		"EngineerID":     engineerID,
		"EngineerName":   engineerName.String,
		"Engineers":      engineers,
		"PickupFormData": pickupFormData,
		"TimeSlots":      []string{"morning", "afternoon", "evening"},
		"Couriers":       []string{"UPS", "FedEx", "DHL"},
	}

	if h.Templates != nil {
		err := h.Templates.ExecuteTemplate(w, "edit-shipment.html", data)
		if err != nil {
			fmt.Printf("Template execution error: %v\n", err)
			http.Error(w, fmt.Sprintf("Failed to render template: %v", err), http.StatusInternalServerError)
			return
		}
	} else {
		// For testing without templates
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Edit Shipment Page")
	}
}

// EditShipmentPOST processes the edit shipment form submission (logistics only)
func (h *ShipmentsHandler) EditShipmentPOST(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Only logistics users can edit shipments
	if user.Role != models.RoleLogistics {
		http.Error(w, "Forbidden", http.StatusForbidden)
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

	// Parse form data
	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Get current shipment to check edit availability
	var currentShipment models.Shipment
	err = h.DB.QueryRowContext(r.Context(),
		`SELECT id, shipment_type, status FROM shipments WHERE id = $1`,
		shipmentID,
	).Scan(&currentShipment.ID, &currentShipment.ShipmentType, &currentShipment.Status)
	if err != nil {
		http.Error(w, "Failed to fetch shipment", http.StatusInternalServerError)
		return
	}

	// Check edit availability
	canEdit, errorMsg := canEditShipment(&currentShipment, h.DB, r.Context(), shipmentID)
	if !canEdit {
		http.Error(w, errorMsg, http.StatusBadRequest)
		return
	}

	// Update software engineer if provided (and not bulk shipment)
	engineerIDStr := r.FormValue("software_engineer_id")
	if engineerIDStr != "" && currentShipment.ShipmentType != models.ShipmentTypeBulkToWarehouse {
		engineerID, err := strconv.ParseInt(engineerIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid engineer ID", http.StatusBadRequest)
			return
		}

		_, err = h.DB.ExecContext(r.Context(),
			`UPDATE shipments SET software_engineer_id = $1, updated_at = $2 WHERE id = $3`,
			engineerID, time.Now(), shipmentID,
		)
		if err != nil {
			fmt.Printf("Error updating engineer: %v\n", err)
			http.Error(w, "Failed to update engineer", http.StatusInternalServerError)
			return
		}
	}

	// Update courier if provided
	courierName := r.FormValue("courier_name")
	if courierName != "" {
		if !models.IsValidCourier(courierName) {
			http.Error(w, "Invalid courier name", http.StatusBadRequest)
			return
		}

		_, err = h.DB.ExecContext(r.Context(),
			`UPDATE shipments SET courier_name = $1, updated_at = $2 WHERE id = $3`,
			courierName, time.Now(), shipmentID,
		)
		if err != nil {
			fmt.Printf("Error updating courier: %v\n", err)
			http.Error(w, "Failed to update courier", http.StatusInternalServerError)
			return
		}
	}

	// Update second tracking number if provided
	secondTrackingNumber := r.FormValue("second_tracking_number")
	_, err = h.DB.ExecContext(r.Context(),
		`UPDATE shipments SET second_tracking_number = $1, updated_at = $2 WHERE id = $3`,
		secondTrackingNumber, time.Now(), shipmentID,
	)
	if err != nil {
		fmt.Printf("Error updating second tracking number: %v\n", err)
		http.Error(w, "Failed to update second tracking number", http.StatusInternalServerError)
		return
	}

	// Update pickup form data if pickup form exists
	var pickupFormExists bool
	err = h.DB.QueryRowContext(r.Context(),
		`SELECT EXISTS(SELECT 1 FROM pickup_forms WHERE shipment_id = $1)`,
		shipmentID,
	).Scan(&pickupFormExists)
	if err != nil {
		fmt.Printf("Error checking pickup form: %v\n", err)
		http.Error(w, "Failed to check pickup form", http.StatusInternalServerError)
		return
	}

	if pickupFormExists {
		// Build updated form data from form fields
		numberOfLaptops, _ := strconv.Atoi(r.FormValue("number_of_laptops"))
		numberOfBoxes, _ := strconv.Atoi(r.FormValue("number_of_boxes"))
		bulkLength, _ := strconv.ParseFloat(r.FormValue("bulk_length"), 64)
		bulkWidth, _ := strconv.ParseFloat(r.FormValue("bulk_width"), 64)
		bulkHeight, _ := strconv.ParseFloat(r.FormValue("bulk_height"), 64)
		bulkWeight, _ := strconv.ParseFloat(r.FormValue("bulk_weight"), 64)
		includeAccessories := r.FormValue("include_accessories") == "on" || r.FormValue("include_accessories") == "true"

		formData := map[string]interface{}{
			"contact_name":            r.FormValue("contact_name"),
			"contact_email":           r.FormValue("contact_email"),
			"contact_phone":           r.FormValue("contact_phone"),
			"pickup_address":          r.FormValue("pickup_address"),
			"pickup_city":             r.FormValue("pickup_city"),
			"pickup_state":            r.FormValue("pickup_state"),
			"pickup_zip":              r.FormValue("pickup_zip"),
			"pickup_date":             r.FormValue("pickup_date"),
			"pickup_time_slot":        r.FormValue("pickup_time_slot"),
			"number_of_laptops":       numberOfLaptops,
			"number_of_boxes":         numberOfBoxes,
			"assignment_type":         r.FormValue("assignment_type"),
			"bulk_length":             bulkLength,
			"bulk_width":              bulkWidth,
			"bulk_height":             bulkHeight,
			"bulk_weight":             bulkWeight,
			"include_accessories":     includeAccessories,
			"accessories_description": r.FormValue("accessories_description"),
			"special_instructions":    r.FormValue("special_instructions"),
			"laptop_serial_number":    r.FormValue("laptop_serial_number"),
			"laptop_model":            r.FormValue("laptop_model"),
			"laptop_ram_gb":           r.FormValue("laptop_ram_gb"),
			"laptop_ssd_gb":           r.FormValue("laptop_ssd_gb"),
			"engineer_name":           r.FormValue("engineer_name"),
		}

		formDataJSON, err := json.Marshal(formData)
		if err != nil {
			http.Error(w, "Failed to encode form data", http.StatusInternalServerError)
			return
		}

		_, err = h.DB.ExecContext(r.Context(),
			`UPDATE pickup_forms SET form_data = $1, submitted_at = $2, submitted_by_user_id = $3
			WHERE shipment_id = $4`,
			formDataJSON, time.Now(), user.ID, shipmentID,
		)
		if err != nil {
			fmt.Printf("Error updating pickup form: %v\n", err)
			http.Error(w, "Failed to update pickup form", http.StatusInternalServerError)
			return
		}
	}

	// Create audit log
	auditDetails, _ := json.Marshal(map[string]interface{}{
		"action": "shipment_edited",
	})

	_, err = h.DB.ExecContext(r.Context(),
		`INSERT INTO audit_logs (user_id, action, entity_type, entity_id, timestamp, details)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		user.ID, "shipment_edited", "shipment", shipmentID, time.Now(), auditDetails,
	)
	if err != nil {
		// Non-critical error
		fmt.Printf("Failed to create audit log: %v\n", err)
	}

	// Redirect back to shipment detail
	redirectURL := fmt.Sprintf("/shipments/%d?success=Shipment+details+updated+successfully", shipmentID)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

// canEditShipment checks if a shipment can be edited based on type and status
func canEditShipment(shipment *models.Shipment, db *sql.DB, ctx context.Context, shipmentID int64) (bool, string) {
	// Warehouse to engineer shipments: can edit if not delivered
	if shipment.ShipmentType == models.ShipmentTypeWarehouseToEngineer {
		if shipment.Status == models.ShipmentStatusDelivered {
			return false, "Cannot edit delivered shipment"
		}
		return true, ""
	}

	// Single and bulk shipments: need pickup form and must not be delivered
	if shipment.Status == models.ShipmentStatusDelivered {
		return false, "Cannot edit delivered shipment"
	}

	// Check if pickup form exists
	var pickupFormExists bool
	err := db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM pickup_forms WHERE shipment_id = $1)`,
		shipmentID,
	).Scan(&pickupFormExists)
	if err != nil {
		return false, "Failed to check pickup form existence"
	}

	if !pickupFormExists {
		return false, "Cannot edit shipment without pickup form"
	}

	return true, ""
}

