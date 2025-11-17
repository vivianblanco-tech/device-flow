package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/yourusername/laptop-tracking-system/internal/database"
	"github.com/yourusername/laptop-tracking-system/internal/middleware"
	"github.com/yourusername/laptop-tracking-system/internal/models"
)

func TestEditShipmentGET(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create logistics user
	var logisticsUserID int64
	err := db.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"logistics@example.com", "hashedpassword", models.RoleLogistics, time.Now(), time.Now(),
	).Scan(&logisticsUserID)
	if err != nil {
		t.Fatalf("Failed to create logistics user: %v", err)
	}

	// Create client user
	var clientUserID int64
	err = db.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"client@example.com", "hashedpassword", models.RoleClient, time.Now(), time.Now(),
	).Scan(&clientUserID)
	if err != nil {
		t.Fatalf("Failed to create client user: %v", err)
	}

	// Create test company
	var companyID int64
	err = db.QueryRowContext(ctx,
		`INSERT INTO client_companies (name, contact_info, created_at)
		VALUES ($1, $2, $3) RETURNING id`,
		"Test Company", json.RawMessage(`{"email":"test@company.com"}`), time.Now(),
	).Scan(&companyID)
	if err != nil {
		t.Fatalf("Failed to create test company: %v", err)
	}

	// Create software engineer
	var engineerID int64
	err = db.QueryRowContext(ctx,
		`INSERT INTO software_engineers (name, email, created_at)
		VALUES ($1, $2, $3) RETURNING id`,
		"John Doe", "john@example.com", time.Now(),
	).Scan(&engineerID)
	if err != nil {
		t.Fatalf("Failed to create software engineer: %v", err)
	}

	// Create shipment with pickup form (single_full_journey)
	var shipmentID int64
	err = db.QueryRowContext(ctx,
		`INSERT INTO shipments (shipment_type, client_company_id, software_engineer_id, status, laptop_count, jira_ticket_number, courier_name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
		models.ShipmentTypeSingleFullJourney, companyID, engineerID, models.ShipmentStatusPickupScheduled, 1, "TEST-123", "UPS", time.Now(), time.Now(),
	).Scan(&shipmentID)
	if err != nil {
		t.Fatalf("Failed to create shipment: %v", err)
	}

	// Create pickup form
	formData := map[string]interface{}{
		"contact_name":       "Jane Smith",
		"contact_email":      "jane@example.com",
		"contact_phone":      "555-1234",
		"pickup_address":     "123 Main St",
		"pickup_city":        "New York",
		"pickup_state":       "NY",
		"pickup_zip":         "10001",
		"pickup_date":        "2024-12-01",
		"pickup_time_slot":   "morning",
		"number_of_laptops":  1,
		"laptop_serial_number": "SN12345",
		"laptop_model":       "Dell XPS 15",
	}
	formDataJSON, _ := json.Marshal(formData)
	
	_, err = db.ExecContext(ctx,
		`INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data)
		VALUES ($1, $2, $3, $4)`,
		shipmentID, clientUserID, time.Now(), formDataJSON,
	)
	if err != nil {
		t.Fatalf("Failed to create pickup form: %v", err)
	}

	templates := loadTestTemplates(t)
	handler := NewShipmentsHandler(db, templates, nil)

	t.Run("logistics user can access edit shipment page", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/shipments/%d/edit", shipmentID), nil)
		req = mux.SetURLVars(req, map[string]string{"id": fmt.Sprintf("%d", shipmentID)})

		user := &models.User{ID: logisticsUserID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.EditShipmentGET(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})

	t.Run("non-logistics users cannot access edit shipment page", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/shipments/%d/edit", shipmentID), nil)
		req = mux.SetURLVars(req, map[string]string{"id": fmt.Sprintf("%d", shipmentID)})

		user := &models.User{ID: clientUserID, Email: "client@example.com", Role: models.RoleClient}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.EditShipmentGET(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("Expected status 403, got %d", w.Code)
		}
	})
}

func TestEditShipmentPOST(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create logistics user
	var logisticsUserID int64
	err := db.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"logistics@example.com", "hashedpassword", models.RoleLogistics, time.Now(), time.Now(),
	).Scan(&logisticsUserID)
	if err != nil {
		t.Fatalf("Failed to create logistics user: %v", err)
	}

	// Create test company
	var companyID int64
	err = db.QueryRowContext(ctx,
		`INSERT INTO client_companies (name, contact_info, created_at)
		VALUES ($1, $2, $3) RETURNING id`,
		"Test Company", json.RawMessage(`{"email":"test@company.com"}`), time.Now(),
	).Scan(&companyID)
	if err != nil {
		t.Fatalf("Failed to create test company: %v", err)
	}

	// Create software engineers
	var engineer1ID, engineer2ID int64
	err = db.QueryRowContext(ctx,
		`INSERT INTO software_engineers (name, email, created_at)
		VALUES ($1, $2, $3) RETURNING id`,
		"John Doe", "john@example.com", time.Now(),
	).Scan(&engineer1ID)
	if err != nil {
		t.Fatalf("Failed to create software engineer 1: %v", err)
	}

	err = db.QueryRowContext(ctx,
		`INSERT INTO software_engineers (name, email, created_at)
		VALUES ($1, $2, $3) RETURNING id`,
		"Jane Smith", "jane@example.com", time.Now(),
	).Scan(&engineer2ID)
	if err != nil {
		t.Fatalf("Failed to create software engineer 2: %v", err)
	}

	// Create shipment with pickup form
	var shipmentID int64
	err = db.QueryRowContext(ctx,
		`INSERT INTO shipments (shipment_type, client_company_id, software_engineer_id, status, laptop_count, jira_ticket_number, courier_name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
		models.ShipmentTypeSingleFullJourney, companyID, engineer1ID, models.ShipmentStatusPickupScheduled, 1, "TEST-123", "UPS", time.Now(), time.Now(),
	).Scan(&shipmentID)
	if err != nil {
		t.Fatalf("Failed to create shipment: %v", err)
	}

	// Create pickup form
	formData := map[string]interface{}{
		"contact_name":       "Jane Smith",
		"contact_email":      "jane@example.com",
		"contact_phone":      "555-1234",
		"pickup_address":     "123 Main St",
		"pickup_city":        "New York",
		"pickup_state":       "NY",
		"pickup_zip":         "10001",
		"pickup_date":        "2024-12-01",
		"pickup_time_slot":   "morning",
		"number_of_laptops":  1,
	}
	formDataJSON, _ := json.Marshal(formData)
	
	_, err = db.ExecContext(ctx,
		`INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data)
		VALUES ($1, $2, $3, $4)`,
		shipmentID, logisticsUserID, time.Now(), formDataJSON,
	)
	if err != nil {
		t.Fatalf("Failed to create pickup form: %v", err)
	}

	templates := loadTestTemplates(t)
	handler := NewShipmentsHandler(db, templates, nil)

	t.Run("logistics user can update shipment software engineer", func(t *testing.T) {
		form := url.Values{}
		form.Add("software_engineer_id", fmt.Sprintf("%d", engineer2ID))
		form.Add("courier_name", "UPS")

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/shipments/%d/edit", shipmentID), strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req = mux.SetURLVars(req, map[string]string{"id": fmt.Sprintf("%d", shipmentID)})

		user := &models.User{ID: logisticsUserID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.EditShipmentPOST(w, req)

		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected status 303, got %d", w.Code)
		}

		// Verify engineer was updated
		var updatedEngineerID int64
		err = db.QueryRowContext(ctx,
			`SELECT software_engineer_id FROM shipments WHERE id = $1`,
			shipmentID,
		).Scan(&updatedEngineerID)
		if err != nil {
			t.Fatalf("Failed to query updated shipment: %v", err)
		}

		if updatedEngineerID != engineer2ID {
			t.Errorf("Expected engineer ID %d, got %d", engineer2ID, updatedEngineerID)
		}
	})

	t.Run("logistics user can update courier", func(t *testing.T) {
		form := url.Values{}
		form.Add("software_engineer_id", fmt.Sprintf("%d", engineer1ID))
		form.Add("courier_name", "FedEx")

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/shipments/%d/edit", shipmentID), strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req = mux.SetURLVars(req, map[string]string{"id": fmt.Sprintf("%d", shipmentID)})

		user := &models.User{ID: logisticsUserID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.EditShipmentPOST(w, req)

		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected status 303, got %d", w.Code)
		}

		// Verify courier was updated
		var updatedCourier string
		err = db.QueryRowContext(ctx,
			`SELECT courier_name FROM shipments WHERE id = $1`,
			shipmentID,
		).Scan(&updatedCourier)
		if err != nil {
			t.Fatalf("Failed to query updated shipment: %v", err)
		}

		if updatedCourier != "FedEx" {
			t.Errorf("Expected courier 'FedEx', got '%s'", updatedCourier)
		}
	})

	t.Run("logistics user can update pickup form fields", func(t *testing.T) {
		form := url.Values{}
		form.Add("software_engineer_id", fmt.Sprintf("%d", engineer1ID))
		form.Add("courier_name", "UPS")
		form.Add("contact_name", "Updated Contact Name")
		form.Add("contact_email", "updated@example.com")
		form.Add("contact_phone", "555-9999")
		form.Add("pickup_address", "456 New Address")
		form.Add("pickup_city", "Boston")
		form.Add("pickup_state", "MA")
		form.Add("pickup_zip", "02101")
		form.Add("pickup_date", "2024-12-15")
		form.Add("pickup_time_slot", "afternoon")
		form.Add("number_of_laptops", "2")
		form.Add("special_instructions", "Handle with care")

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/shipments/%d/edit", shipmentID), strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req = mux.SetURLVars(req, map[string]string{"id": fmt.Sprintf("%d", shipmentID)})

		user := &models.User{ID: logisticsUserID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.EditShipmentPOST(w, req)

		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected status 303, got %d", w.Code)
		}

		// Verify pickup form was updated
		var updatedFormData json.RawMessage
		err = db.QueryRowContext(ctx,
			`SELECT form_data FROM pickup_forms WHERE shipment_id = $1`,
			shipmentID,
		).Scan(&updatedFormData)
		if err != nil {
			t.Fatalf("Failed to query updated pickup form: %v", err)
		}

		var updatedForm map[string]interface{}
		err = json.Unmarshal(updatedFormData, &updatedForm)
		if err != nil {
			t.Fatalf("Failed to parse updated form data: %v", err)
		}

		if updatedForm["contact_name"] != "Updated Contact Name" {
			t.Errorf("Expected contact_name 'Updated Contact Name', got '%v'", updatedForm["contact_name"])
		}
		if updatedForm["pickup_city"] != "Boston" {
			t.Errorf("Expected pickup_city 'Boston', got '%v'", updatedForm["pickup_city"])
		}
	})
}

func TestEditShipmentAvailability(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create logistics user
	var logisticsUserID int64
	err := db.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"logistics@example.com", "hashedpassword", models.RoleLogistics, time.Now(), time.Now(),
	).Scan(&logisticsUserID)
	if err != nil {
		t.Fatalf("Failed to create logistics user: %v", err)
	}

	// Create test company
	var companyID int64
	err = db.QueryRowContext(ctx,
		`INSERT INTO client_companies (name, contact_info, created_at)
		VALUES ($1, $2, $3) RETURNING id`,
		"Test Company", json.RawMessage(`{"email":"test@company.com"}`), time.Now(),
	).Scan(&companyID)
	if err != nil {
		t.Fatalf("Failed to create test company: %v", err)
	}

	templates := loadTestTemplates(t)
	handler := NewShipmentsHandler(db, templates, nil)

	t.Run("edit not available for single shipment without pickup form", func(t *testing.T) {
		// Create shipment WITHOUT pickup form
		var shipmentID int64
		err = db.QueryRowContext(ctx,
			`INSERT INTO shipments (shipment_type, client_company_id, status, laptop_count, jira_ticket_number, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
			models.ShipmentTypeSingleFullJourney, companyID, models.ShipmentStatusPendingPickup, 1, "TEST-123", time.Now(), time.Now(),
		).Scan(&shipmentID)
		if err != nil {
			t.Fatalf("Failed to create shipment: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/shipments/%d/edit", shipmentID), nil)
		req = mux.SetURLVars(req, map[string]string{"id": fmt.Sprintf("%d", shipmentID)})

		user := &models.User{ID: logisticsUserID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.EditShipmentGET(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("edit not available for delivered shipment", func(t *testing.T) {
		// Create delivered shipment WITH pickup form
		var shipmentID int64
		err = db.QueryRowContext(ctx,
			`INSERT INTO shipments (shipment_type, client_company_id, status, laptop_count, jira_ticket_number, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
			models.ShipmentTypeSingleFullJourney, companyID, models.ShipmentStatusDelivered, 1, "TEST-124", time.Now(), time.Now(),
		).Scan(&shipmentID)
		if err != nil {
			t.Fatalf("Failed to create shipment: %v", err)
		}

		// Create pickup form
		formData := map[string]interface{}{"contact_name": "Test"}
		formDataJSON, _ := json.Marshal(formData)
		
		_, err = db.ExecContext(ctx,
			`INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data)
			VALUES ($1, $2, $3, $4)`,
			shipmentID, logisticsUserID, time.Now(), formDataJSON,
		)
		if err != nil {
			t.Fatalf("Failed to create pickup form: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/shipments/%d/edit", shipmentID), nil)
		req = mux.SetURLVars(req, map[string]string{"id": fmt.Sprintf("%d", shipmentID)})

		user := &models.User{ID: logisticsUserID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.EditShipmentGET(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("edit available for warehouse to engineer shipment regardless of pickup form", func(t *testing.T) {
		// Create warehouse_to_engineer shipment WITHOUT pickup form
		var shipmentID int64
		err = db.QueryRowContext(ctx,
			`INSERT INTO shipments (shipment_type, client_company_id, status, laptop_count, jira_ticket_number, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
			models.ShipmentTypeWarehouseToEngineer, companyID, models.ShipmentStatusReleasedFromWarehouse, 1, "TEST-125", time.Now(), time.Now(),
		).Scan(&shipmentID)
		if err != nil {
			t.Fatalf("Failed to create shipment: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/shipments/%d/edit", shipmentID), nil)
		req = mux.SetURLVars(req, map[string]string{"id": fmt.Sprintf("%d", shipmentID)})

		user := &models.User{ID: logisticsUserID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.EditShipmentGET(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})

	t.Run("edit not available for warehouse to engineer shipment with delivered status", func(t *testing.T) {
		// Create warehouse_to_engineer shipment with delivered status
		var shipmentID int64
		err = db.QueryRowContext(ctx,
			`INSERT INTO shipments (shipment_type, client_company_id, status, laptop_count, jira_ticket_number, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
			models.ShipmentTypeWarehouseToEngineer, companyID, models.ShipmentStatusDelivered, 1, "TEST-126", time.Now(), time.Now(),
		).Scan(&shipmentID)
		if err != nil {
			t.Fatalf("Failed to create shipment: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/shipments/%d/edit", shipmentID), nil)
		req = mux.SetURLVars(req, map[string]string{"id": fmt.Sprintf("%d", shipmentID)})

		user := &models.User{ID: logisticsUserID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.EditShipmentGET(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})
}

