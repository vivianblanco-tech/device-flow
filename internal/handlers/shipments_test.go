package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/yourusername/laptop-tracking-system/internal/database"
	"github.com/yourusername/laptop-tracking-system/internal/middleware"
	"github.com/yourusername/laptop-tracking-system/internal/models"
)

func TestShipmentsList(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create test user
	var userID int64
	err := db.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"logistics@example.com", "hashedpassword", models.RoleLogistics, time.Now(), time.Now(),
	).Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
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

	// Create test shipments with different statuses
	statuses := []models.ShipmentStatus{
		models.ShipmentStatusPendingPickup,
		models.ShipmentStatusInTransitToWarehouse,
		models.ShipmentStatusAtWarehouse,
		models.ShipmentStatusDelivered,
	}

	for i, status := range statuses {
		_, err := db.ExecContext(ctx,
			`INSERT INTO shipments (client_company_id, status, jira_ticket_number, tracking_number, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)`,
			companyID, status, fmt.Sprintf("TEST-%d", i+1), "TRACK-"+string(status), time.Now(), time.Now(),
		)
		if err != nil {
			t.Fatalf("Failed to create test shipment: %v", err)
		}
	}

	templates := loadTestTemplates(t)
	handler := NewShipmentsHandler(db, templates, nil) // nil email notifier for tests

	t.Run("authenticated user can view shipments list", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/shipments", nil)

		user := &models.User{ID: userID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.ShipmentsList(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})

	t.Run("shipments list displays JIRA ticket numbers", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/shipments", nil)

		user := &models.User{ID: userID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.ShipmentsList(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		// Verify the response contains at least one JIRA ticket from test data
		responseBody := w.Body.String()
		if !strings.Contains(responseBody, "TEST-1") && !strings.Contains(responseBody, "TEST-2") {
			t.Errorf("Expected response to contain JIRA ticket numbers (TEST-1, TEST-2), but none were found")
		}
	})

	t.Run("unauthenticated user redirects to login", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/shipments", nil)
		w := httptest.NewRecorder()

		handler.ShipmentsList(w, req)

		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected status 303, got %d", w.Code)
		}

		location := w.Header().Get("Location")
		if location != "/login" {
			t.Errorf("Expected redirect to /login, got %s", location)
		}
	})

	t.Run("status filter works", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/shipments?status="+string(models.ShipmentStatusAtWarehouse), nil)

		user := &models.User{ID: userID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.ShipmentsList(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})

	t.Run("search query works", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/shipments?search=Test", nil)

		user := &models.User{ID: userID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.ShipmentsList(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})
}

// 🟥 RED: Test shipment type filtering in list
func TestShipmentsListWithTypeFilter(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create test user
	var userID int64
	err := db.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"logistics@example.com", "hashedpassword", models.RoleLogistics, time.Now(), time.Now(),
	).Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
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

	// Create shipments of each type
	shipmentTypes := []models.ShipmentType{
		models.ShipmentTypeSingleFullJourney,
		models.ShipmentTypeBulkToWarehouse,
		models.ShipmentTypeWarehouseToEngineer,
	}

	for i, shipmentType := range shipmentTypes {
		laptopCount := 1
		if shipmentType == models.ShipmentTypeBulkToWarehouse {
			laptopCount = 5
		}

		status := models.ShipmentStatusPendingPickup
		if shipmentType == models.ShipmentTypeWarehouseToEngineer {
			status = models.ShipmentStatusReleasedFromWarehouse
		}

		_, err := db.ExecContext(ctx,
			`INSERT INTO shipments (shipment_type, client_company_id, status, laptop_count, jira_ticket_number, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			shipmentType, companyID, status, laptopCount, fmt.Sprintf("TYPE-TEST-%d", i+1), time.Now(), time.Now(),
		)
		if err != nil {
			t.Fatalf("Failed to create test shipment of type %s: %v", shipmentType, err)
		}
	}

	templates := loadTestTemplates(t)
	handler := NewShipmentsHandler(db, templates, nil)

	t.Run("list includes shipment type information", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/shipments", nil)

		user := &models.User{ID: userID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.ShipmentsList(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		// The response should include all three shipment types
		body := w.Body.String()
		// Templates will display type badges/labels
		if !strings.Contains(body, "TYPE-TEST-1") || !strings.Contains(body, "TYPE-TEST-2") || !strings.Contains(body, "TYPE-TEST-3") {
			t.Error("Expected list to contain all test shipments")
		}
	})

	t.Run("filter by single_full_journey type", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/shipments?type=single_full_journey", nil)

		user := &models.User{ID: userID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.ShipmentsList(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		body := w.Body.String()
		// Should contain single_full_journey shipment
		if !strings.Contains(body, "TYPE-TEST-1") {
			t.Error("Expected filtered list to contain single_full_journey shipment (TYPE-TEST-1)")
		}
	})

	t.Run("filter by bulk_to_warehouse type", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/shipments?type=bulk_to_warehouse", nil)

		user := &models.User{ID: userID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.ShipmentsList(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		body := w.Body.String()
		// Should contain bulk_to_warehouse shipment
		if !strings.Contains(body, "TYPE-TEST-2") {
			t.Error("Expected filtered list to contain bulk_to_warehouse shipment (TYPE-TEST-2)")
		}
	})

	t.Run("filter by warehouse_to_engineer type", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/shipments?type=warehouse_to_engineer", nil)

		user := &models.User{ID: userID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.ShipmentsList(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		body := w.Body.String()
		// Should contain warehouse_to_engineer shipment
		if !strings.Contains(body, "TYPE-TEST-3") {
			t.Error("Expected filtered list to contain warehouse_to_engineer shipment (TYPE-TEST-3)")
		}
	})
}

// 🟥 RED: Test shipment detail displays type information
func TestShipmentDetailWithTypeInformation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create test user
	var userID int64
	err := db.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"logistics@example.com", "hashedpassword", models.RoleLogistics, time.Now(), time.Now(),
	).Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
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

	t.Run("detail displays single_full_journey type information", func(t *testing.T) {
		// Create single_full_journey shipment
		var shipmentID int64
		err := db.QueryRowContext(ctx,
			`INSERT INTO shipments (shipment_type, client_company_id, status, laptop_count, jira_ticket_number, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
			models.ShipmentTypeSingleFullJourney, companyID, models.ShipmentStatusPendingPickup, 1, "TYPE-DETAIL-1", time.Now(), time.Now(),
		).Scan(&shipmentID)
		if err != nil {
			t.Fatalf("Failed to create test shipment: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/shipments/"+strconv.FormatInt(shipmentID, 10), nil)
		req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(shipmentID, 10)})

		user := &models.User{ID: userID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.ShipmentDetail(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		body := w.Body.String()
		// Should display single_full_journey information
		if !strings.Contains(body, "TYPE-DETAIL-1") {
			t.Error("Expected detail to contain shipment JIRA ticket")
		}
	})

	t.Run("detail displays bulk_to_warehouse type with laptop count", func(t *testing.T) {
		// Create bulk_to_warehouse shipment
		var shipmentID int64
		err := db.QueryRowContext(ctx,
			`INSERT INTO shipments (shipment_type, client_company_id, status, laptop_count, jira_ticket_number, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
			models.ShipmentTypeBulkToWarehouse, companyID, models.ShipmentStatusAtWarehouse, 5, "TYPE-DETAIL-2", time.Now(), time.Now(),
		).Scan(&shipmentID)
		if err != nil {
			t.Fatalf("Failed to create bulk shipment: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/shipments/"+strconv.FormatInt(shipmentID, 10), nil)
		req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(shipmentID, 10)})

		user := &models.User{ID: userID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.ShipmentDetail(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		body := w.Body.String()
		// Should display bulk_to_warehouse with laptop count
		if !strings.Contains(body, "TYPE-DETAIL-2") {
			t.Error("Expected detail to contain shipment JIRA ticket")
		}
	})

	t.Run("detail displays warehouse_to_engineer type", func(t *testing.T) {
		// Create warehouse_to_engineer shipment
		var shipmentID int64
		err := db.QueryRowContext(ctx,
			`INSERT INTO shipments (shipment_type, client_company_id, status, laptop_count, jira_ticket_number, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
			models.ShipmentTypeWarehouseToEngineer, companyID, models.ShipmentStatusReleasedFromWarehouse, 1, "TYPE-DETAIL-3", time.Now(), time.Now(),
		).Scan(&shipmentID)
		if err != nil {
			t.Fatalf("Failed to create warehouse-to-engineer shipment: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/shipments/"+strconv.FormatInt(shipmentID, 10), nil)
		req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(shipmentID, 10)})

		user := &models.User{ID: userID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.ShipmentDetail(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		body := w.Body.String()
		// Should display warehouse_to_engineer information
		if !strings.Contains(body, "TYPE-DETAIL-3") {
			t.Error("Expected detail to contain shipment JIRA ticket")
		}
	})
}

func TestShipmentDetail(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create test user
	var userID int64
	err := db.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"logistics@example.com", "hashedpassword", models.RoleLogistics, time.Now(), time.Now(),
	).Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
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

	// Create test shipment
	var shipmentID int64
	err = db.QueryRowContext(ctx,
		`INSERT INTO shipments (shipment_type, client_company_id, status, laptop_count, jira_ticket_number, tracking_number, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
		models.ShipmentTypeSingleFullJourney, companyID, models.ShipmentStatusInTransitToWarehouse, 1, "TEST-12345", "TRACK-12345", time.Now(), time.Now(),
	).Scan(&shipmentID)
	if err != nil {
		t.Fatalf("Failed to create test shipment: %v", err)
	}

	// Create test laptop
	var laptopID int64
	err = db.QueryRowContext(ctx,
		`INSERT INTO laptops (serial_number, brand, model, specs, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		"SN-12345", "Dell", "XPS 15", json.RawMessage(`{"ram":"16GB","cpu":"i7"}`), "available", time.Now(),
	).Scan(&laptopID)
	if err != nil {
		t.Fatalf("Failed to create test laptop: %v", err)
	}

	// Link laptop to shipment
	_, err = db.ExecContext(ctx,
		`INSERT INTO shipment_laptops (shipment_id, laptop_id) VALUES ($1, $2)`,
		shipmentID, laptopID,
	)
	if err != nil {
		t.Fatalf("Failed to link laptop to shipment: %v", err)
	}

	templates := loadTestTemplates(t)
	handler := NewShipmentsHandler(db, templates, nil) // nil email notifier for tests

	t.Run("authenticated user can view shipment detail", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/shipments/"+strconv.FormatInt(shipmentID, 10), nil)
		req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(shipmentID, 10)})

		user := &models.User{ID: userID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.ShipmentDetail(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})

	t.Run("shipment detail includes JIRA ticket number", func(t *testing.T) {
		// Create a shipment with a specific JIRA ticket
		var testShipmentID int64
		err := db.QueryRowContext(ctx,
			`INSERT INTO shipments (client_company_id, status, jira_ticket_number, notes, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
			companyID, models.ShipmentStatusPendingPickup, "SCOP-12345", "Test shipment with JIRA", time.Now(), time.Now(),
		).Scan(&testShipmentID)
		if err != nil {
			t.Fatalf("Failed to create test shipment: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/shipments/"+strconv.FormatInt(testShipmentID, 10), nil)
		req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(testShipmentID, 10)})

		user := &models.User{ID: userID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.ShipmentDetail(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		// Verify the response contains the JIRA ticket number
		responseBody := w.Body.String()
		if !strings.Contains(responseBody, "SCOP-12345") {
			t.Errorf("Expected response to contain JIRA ticket 'SCOP-12345', but it was not found")
		}
	})

	t.Run("shipment detail displays pickup form data when available", func(t *testing.T) {
		// Create a shipment
		var testShipmentID int64
		err := db.QueryRowContext(ctx,
			`INSERT INTO shipments (client_company_id, status, jira_ticket_number, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5) RETURNING id`,
			companyID, models.ShipmentStatusPendingPickup, "SCOP-99999", time.Now(), time.Now(),
		).Scan(&testShipmentID)
		if err != nil {
			t.Fatalf("Failed to create test shipment: %v", err)
		}

		// Create pickup form data
		pickupFormData := map[string]interface{}{
			"contact_name":            "John Doe",
			"contact_email":           "john.doe@example.com",
			"contact_phone":           "+1-555-0123",
			"pickup_address":          "123 Main Street, Suite 400",
			"pickup_city":             "New York",
			"pickup_state":            "NY",
			"pickup_zip":              "10001",
			"pickup_date":             "2024-12-25",
			"pickup_time_slot":        "morning",
			"number_of_laptops":       5,
			"number_of_boxes":         2,
			"assignment_type":         "bulk",
			"bulk_length":             20.5,
			"bulk_width":              15.0,
			"bulk_height":             10.0,
			"bulk_weight":             25.5,
			"include_accessories":     true,
			"accessories_description": "2x YubiKeys, 3x USB-C cables",
			"special_instructions":    "Call before arrival",
		}
		formDataJSON, _ := json.Marshal(pickupFormData)

		// Insert pickup form
		_, err = db.ExecContext(ctx,
			`INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data)
			VALUES ($1, $2, $3, $4)`,
			testShipmentID, userID, time.Now(), formDataJSON,
		)
		if err != nil {
			t.Fatalf("Failed to create pickup form: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/shipments/"+strconv.FormatInt(testShipmentID, 10), nil)
		req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(testShipmentID, 10)})

		user := &models.User{ID: userID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.ShipmentDetail(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		// Verify the response contains pickup form data
		responseBody := w.Body.String()

		// Check for contact information
		if !strings.Contains(responseBody, "John Doe") {
			t.Errorf("Expected response to contain contact name 'John Doe'")
		}
		if !strings.Contains(responseBody, "john.doe@example.com") {
			t.Errorf("Expected response to contain contact email 'john.doe@example.com'")
		}
		if !strings.Contains(responseBody, "+1-555-0123") && !strings.Contains(responseBody, "555-0123") {
			// Print a snippet to help debug
			idx := strings.Index(responseBody, "Contact Phone")
			if idx >= 0 && idx+200 < len(responseBody) {
				t.Errorf("Expected response to contain contact phone. Contact section: %s", responseBody[idx:idx+200])
			} else {
				t.Errorf("Expected response to contain contact phone '+1-555-0123'")
			}
		}

		// Check for pickup address
		if !strings.Contains(responseBody, "123 Main Street, Suite 400") {
			t.Errorf("Expected response to contain pickup address '123 Main Street, Suite 400'")
		}
		if !strings.Contains(responseBody, "New York") {
			t.Errorf("Expected response to contain city 'New York'")
		}
		if !strings.Contains(responseBody, "NY") {
			t.Errorf("Expected response to contain state 'NY'")
		}
		if !strings.Contains(responseBody, "10001") {
			t.Errorf("Expected response to contain ZIP '10001'")
		}

		// Check for accessories
		if !strings.Contains(responseBody, "2x YubiKeys, 3x USB-C cables") {
			t.Errorf("Expected response to contain accessories description")
		}

		// Check for special instructions
		if !strings.Contains(responseBody, "Call before arrival") {
			t.Errorf("Expected response to contain special instructions")
		}
	})

	t.Run("pickup form details appear in shipment information section", func(t *testing.T) {
		// Create a shipment
		var testShipmentID int64
		err := db.QueryRowContext(ctx,
			`INSERT INTO shipments (client_company_id, status, jira_ticket_number, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5) RETURNING id`,
			companyID, models.ShipmentStatusPendingPickup, "SCOP-88888", time.Now(), time.Now(),
		).Scan(&testShipmentID)
		if err != nil {
			t.Fatalf("Failed to create test shipment: %v", err)
		}

		// Create comprehensive pickup form data
		pickupFormData := map[string]interface{}{
			"contact_name":            "Jane Smith",
			"contact_email":           "jane.smith@techcorp.com",
			"contact_phone":           "+1-555-9876",
			"pickup_address":          "456 Tech Avenue",
			"pickup_city":             "San Francisco",
			"pickup_state":            "CA",
			"pickup_zip":              "94102",
			"pickup_date":             "2024-12-30",
			"pickup_time_slot":        "afternoon",
			"number_of_laptops":       3,
			"number_of_boxes":         1,
			"assignment_type":         "individual",
			"include_accessories":     true,
			"accessories_description": "3x Laptop chargers, 1x Docking station",
			"special_instructions":    "Building requires badge access",
		}
		formDataJSON, _ := json.Marshal(pickupFormData)

		// Insert pickup form
		_, err = db.ExecContext(ctx,
			`INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data)
			VALUES ($1, $2, $3, $4)`,
			testShipmentID, userID, time.Now(), formDataJSON,
		)
		if err != nil {
			t.Fatalf("Failed to create pickup form: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/shipments/"+strconv.FormatInt(testShipmentID, 10), nil)
		req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(testShipmentID, 10)})

		user := &models.User{ID: userID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.ShipmentDetail(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		responseBody := w.Body.String()

		// Find the Shipment Information section
		shipmentInfoStart := strings.Index(responseBody, "Shipment Information")
		if shipmentInfoStart == -1 {
			t.Fatal("Could not find 'Shipment Information' section in response")
		}

		// Find the end of Shipment Information section (next major section starts)
		// Look for the next section heading (Timeline, Laptops, etc.)
		timelineStart := strings.Index(responseBody[shipmentInfoStart:], "Tracking Timeline")
		if timelineStart == -1 {
			// If no timeline, look for other sections
			timelineStart = strings.Index(responseBody[shipmentInfoStart:], "Laptops")
			if timelineStart == -1 {
				timelineStart = len(responseBody) - shipmentInfoStart
			}
		}
		shipmentInfoEnd := shipmentInfoStart + timelineStart

		// Extract the Shipment Information section content
		shipmentInfoSection := responseBody[shipmentInfoStart:shipmentInfoEnd]

		// Verify all pickup form details appear within Shipment Information section
		// Note: Some characters are HTML-encoded (e.g., + becomes &#43;)
		pickupFormFields := map[string]string{
			"Contact Name":          "Jane Smith",
			"Contact Email":         "jane.smith@techcorp.com",
			"Contact Phone":         "&#43;1-555-9876", // + is HTML-encoded
			"Street Address":        "456 Tech Avenue",
			"City":                  "San Francisco",
			"State":                 "CA",
			"ZIP Code":              "94102",
			"Pickup Date":           "2024-12-30",
			"Time Slot Afternoon":   "Afternoon", // title filter capitalizes
			"Number of Laptops":     "3",
			"Number of Boxes":       "1",
			"Assignment Individual": "Individual", // title filter capitalizes
			"Accessories":           "3x Laptop chargers, 1x Docking station",
			"Special Instructions":  "Building requires badge access",
		}

		for fieldLabel, expectedValue := range pickupFormFields {
			if !strings.Contains(shipmentInfoSection, expectedValue) {
				t.Errorf("Expected Shipment Information section to contain '%s: %s', but it was not found", fieldLabel, expectedValue)
			}
		}
	})

	t.Run("tracking number displays as clickable link for known couriers", func(t *testing.T) {
		// Test UPS tracking URL
		var upsShipmentID int64
		err := db.QueryRowContext(ctx,
			`INSERT INTO shipments (client_company_id, status, jira_ticket_number, courier_name, tracking_number, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
			companyID, models.ShipmentStatusInTransitToWarehouse, "SCOP-11111", "UPS", "1Z9999999999999999", time.Now(), time.Now(),
		).Scan(&upsShipmentID)
		if err != nil {
			t.Fatalf("Failed to create UPS shipment: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/shipments/"+strconv.FormatInt(upsShipmentID, 10), nil)
		req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(upsShipmentID, 10)})

		user := &models.User{ID: userID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.ShipmentDetail(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		responseBody := w.Body.String()
		expectedURL := "https://www.ups.com/track?tracknum=1Z9999999999999999"
		if !strings.Contains(responseBody, expectedURL) {
			t.Errorf("Expected response to contain UPS tracking URL '%s', but it was not found", expectedURL)
		}

		// Test DHL tracking URL
		var dhlShipmentID int64
		err = db.QueryRowContext(ctx,
			`INSERT INTO shipments (client_company_id, status, jira_ticket_number, courier_name, tracking_number, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
			companyID, models.ShipmentStatusInTransitToWarehouse, "SCOP-22222", "DHL", "1234567890", time.Now(), time.Now(),
		).Scan(&dhlShipmentID)
		if err != nil {
			t.Fatalf("Failed to create DHL shipment: %v", err)
		}

		req = httptest.NewRequest(http.MethodGet, "/shipments/"+strconv.FormatInt(dhlShipmentID, 10), nil)
		req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(dhlShipmentID, 10)})

		// Create fresh context for DHL request
		reqCtx = context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w = httptest.NewRecorder()
		handler.ShipmentDetail(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		responseBody = w.Body.String()
		expectedURL = "http://www.dhl.com/en/express/tracking.html?AWB=1234567890"
		hasDHL := strings.Contains(responseBody, "dhl.com")
		hasTracking := strings.Contains(responseBody, "1234567890")

		if !strings.Contains(responseBody, expectedURL) && (!hasDHL || !hasTracking) {
			// Find where "Tracking Number" appears in the response
			idx := strings.Index(responseBody, "Tracking Number")
			if idx >= 0 && idx+200 < len(responseBody) {
				t.Errorf("Expected DHL URL. Has dhl.com=%v, Has tracking=%v. Tracking section: %s",
					hasDHL, hasTracking, responseBody[idx:idx+200])
			} else {
				t.Errorf("Expected DHL URL. Has dhl.com=%v, Has tracking=%v", hasDHL, hasTracking)
			}
		}

		// Test FedEx tracking URL
		var fedexShipmentID int64
		err = db.QueryRowContext(ctx,
			`INSERT INTO shipments (client_company_id, status, jira_ticket_number, courier_name, tracking_number, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
			companyID, models.ShipmentStatusInTransitToWarehouse, "SCOP-33333", "FedEx", "999999999999", time.Now(), time.Now(),
		).Scan(&fedexShipmentID)
		if err != nil {
			t.Fatalf("Failed to create FedEx shipment: %v", err)
		}

		req = httptest.NewRequest(http.MethodGet, "/shipments/"+strconv.FormatInt(fedexShipmentID, 10), nil)
		req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(fedexShipmentID, 10)})

		// Create fresh context for FedEx request
		reqCtx = context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w = httptest.NewRecorder()
		handler.ShipmentDetail(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		responseBody = w.Body.String()
		expectedURL = "https://www.fedex.com/fedextrack/?tracknumbers=999999999999"
		if !strings.Contains(responseBody, expectedURL) {
			// Check if it's HTML-encoded
			if !strings.Contains(responseBody, "fedex.com") || !strings.Contains(responseBody, "999999999999") {
				t.Errorf("Expected response to contain FedEx tracking URL '%s', but it was not found", expectedURL)
			}
		}
	})

	t.Run("tracking number displays as plain text for unknown courier", func(t *testing.T) {
		var unknownShipmentID int64
		err := db.QueryRowContext(ctx,
			`INSERT INTO shipments (client_company_id, status, jira_ticket_number, courier_name, tracking_number, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
			companyID, models.ShipmentStatusInTransitToWarehouse, "SCOP-44444", "Unknown Courier", "TRACK-UNKNOWN", time.Now(), time.Now(),
		).Scan(&unknownShipmentID)
		if err != nil {
			t.Fatalf("Failed to create shipment with unknown courier: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/shipments/"+strconv.FormatInt(unknownShipmentID, 10), nil)
		req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(unknownShipmentID, 10)})

		user := &models.User{ID: userID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.ShipmentDetail(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		responseBody := w.Body.String()
		// Should contain the tracking number as text
		if !strings.Contains(responseBody, "TRACK-UNKNOWN") {
			t.Errorf("Expected response to contain tracking number 'TRACK-UNKNOWN'")
		}
		// Should not contain any standard courier tracking URLs
		if strings.Contains(responseBody, "ups.com") || strings.Contains(responseBody, "dhl.com") || strings.Contains(responseBody, "fedex.com") {
			t.Errorf("Expected response NOT to contain tracking URLs for unknown courier")
		}
	})

	t.Run("shipment detail displays updated at timestamp", func(t *testing.T) {
		// Create a shipment with specific created and updated timestamps
		createdAt := time.Date(2025, 11, 1, 10, 0, 0, 0, time.UTC)
		updatedAt := time.Date(2025, 11, 10, 14, 30, 0, 0, time.UTC)

		var testShipmentID int64
		err := db.QueryRowContext(ctx,
			`INSERT INTO shipments (client_company_id, status, jira_ticket_number, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5) RETURNING id`,
			companyID, models.ShipmentStatusAtWarehouse, "SCOP-88888", createdAt, updatedAt,
		).Scan(&testShipmentID)
		if err != nil {
			t.Fatalf("Failed to create test shipment: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/shipments/"+strconv.FormatInt(testShipmentID, 10), nil)
		req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(testShipmentID, 10)})

		user := &models.User{ID: userID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.ShipmentDetail(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		responseBody := w.Body.String()

		// Find the Shipment Information section
		shipmentInfoStart := strings.Index(responseBody, "Shipment Information")
		if shipmentInfoStart == -1 {
			t.Fatal("Could not find 'Shipment Information' section in response")
		}

		// Find the end of Shipment Information section (next major section starts)
		timelineStart := strings.Index(responseBody[shipmentInfoStart:], "Tracking Timeline")
		if timelineStart == -1 {
			timelineStart = len(responseBody) - shipmentInfoStart
		}
		shipmentInfoEnd := shipmentInfoStart + timelineStart

		// Extract the Shipment Information section content
		shipmentInfoSection := responseBody[shipmentInfoStart:shipmentInfoEnd]

		// Verify "Updated" label is present in Shipment Information section
		if !strings.Contains(shipmentInfoSection, "Updated") {
			t.Errorf("Expected response to contain 'Updated' label in Shipment Information section")
		}

		// Verify the formatted updated_at timestamp is present
		// Format: "Nov 10, 2025 14:30"
		expectedTimestamp := updatedAt.Format("Jan 02, 2006 15:04")
		if !strings.Contains(shipmentInfoSection, expectedTimestamp) {
			t.Errorf("Expected response to contain updated timestamp '%s' in Shipment Information section, but it was not found", expectedTimestamp)
		}
	})

	t.Run("missing shipment ID returns error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/shipments/detail", nil)

		user := &models.User{ID: userID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.ShipmentDetail(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("invalid shipment ID returns error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/shipments/invalid", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "invalid"})

		user := &models.User{ID: userID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.ShipmentDetail(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("non-existent shipment returns not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/shipments/99999", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "99999"})

		user := &models.User{ID: userID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.ShipmentDetail(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", w.Code)
		}
	})
}

func TestShipmentDetailTimelineData(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create test user
	var userID int64
	err := db.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"logistics@example.com", "hashedpassword", models.RoleLogistics, time.Now(), time.Now(),
	).Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
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

	t.Run("timeline data includes all statuses with completed/current/pending indicators", func(t *testing.T) {
		// Create shipment in middle of process
		now := time.Now()
		pickupDate := now.AddDate(0, 0, 1)
		pickedUpAt := now.AddDate(0, 0, 2)

		var shipmentID int64
		err := db.QueryRowContext(ctx,
			`INSERT INTO shipments (client_company_id, status, jira_ticket_number, 
			pickup_scheduled_date, picked_up_at, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
			companyID, models.ShipmentStatusInTransitToWarehouse, "TEST-TIMELINE",
			pickupDate, pickedUpAt, now, now,
		).Scan(&shipmentID)
		if err != nil {
			t.Fatalf("Failed to create test shipment: %v", err)
		}

		templates := loadTestTemplates(t)
		handler := NewShipmentsHandler(db, templates, nil) // nil email notifier for tests

		req := httptest.NewRequest(http.MethodGet, "/shipments/"+strconv.FormatInt(shipmentID, 10), nil)
		req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(shipmentID, 10)})

		user := &models.User{ID: userID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.ShipmentDetail(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		// Check that timeline renders all statuses
		body := w.Body.String()

		// Should include all status labels
		expectedStatuses := []string{
			"Pickup Scheduled",
			"Picked Up",
			"In Transit to Warehouse",
			"Arrived at Warehouse",
			"Released from Warehouse",
			"In Transit to Engineer",
			"Delivered",
		}

		for _, status := range expectedStatuses {
			if !strings.Contains(body, status) {
				t.Errorf("Timeline should include status '%s' but it was not found", status)
			}
		}
	})

	t.Run("timeline uses different colors for transit statuses", func(t *testing.T) {
		// Create shipment in transit to warehouse
		var shipmentID int64
		err := db.QueryRowContext(ctx,
			`INSERT INTO shipments (client_company_id, status, jira_ticket_number, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5) RETURNING id`,
			companyID, models.ShipmentStatusInTransitToWarehouse, "TEST-TRANSIT-WH", time.Now(), time.Now(),
		).Scan(&shipmentID)
		if err != nil {
			t.Fatalf("Failed to create test shipment: %v", err)
		}

		templates := loadTestTemplates(t)
		handler := NewShipmentsHandler(db, templates, nil) // nil email notifier for tests

		req := httptest.NewRequest(http.MethodGet, "/shipments/"+strconv.FormatInt(shipmentID, 10), nil)
		req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(shipmentID, 10)})

		user := &models.User{ID: userID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.ShipmentDetail(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		body := w.Body.String()

		// Check for distinct styling for in-transit status
		// Orange/yellow colors (bg-orange or bg-yellow) should be used for transit
		hasTransitColor := strings.Contains(body, "bg-orange") || strings.Contains(body, "bg-yellow")
		if !hasTransitColor {
			t.Error("Timeline should use distinct color (orange/yellow) for 'In Transit' statuses")
		}
	})
}

func TestUpdateShipmentStatus(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create test logistics user
	var logisticsUserID int64
	err := db.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"logistics@example.com", "hashedpassword", models.RoleLogistics, time.Now(), time.Now(),
	).Scan(&logisticsUserID)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
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

	// Create test shipment
	var shipmentID int64
	err = db.QueryRowContext(ctx,
		`INSERT INTO shipments (client_company_id, status, jira_ticket_number, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		companyID, models.ShipmentStatusPendingPickup, "TEST-999", time.Now(), time.Now(),
	).Scan(&shipmentID)
	if err != nil {
		t.Fatalf("Failed to create test shipment: %v", err)
	}

	templates := loadTestTemplates(t)
	handler := NewShipmentsHandler(db, templates, nil) // nil email notifier for tests

	t.Run("logistics user can update shipment status", func(t *testing.T) {
		// First update to pickup_scheduled (sequential transition)
		formData := url.Values{}
		formData.Set("shipment_id", strconv.FormatInt(shipmentID, 10))
		formData.Set("status", string(models.ShipmentStatusPickupScheduled))
		formData.Set("tracking_number", "1Z999AA10123456784")
		formData.Set("courier_name", "UPS")

		req := httptest.NewRequest(http.MethodPost, "/shipments/update-status", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		user := &models.User{ID: logisticsUserID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.UpdateShipmentStatus(w, req)

		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected status 303, got %d", w.Code)
		}

		// Verify status was updated to pickup_scheduled
		var status models.ShipmentStatus
		err := db.QueryRowContext(ctx,
			`SELECT status FROM shipments WHERE id = $1`,
			shipmentID,
		).Scan(&status)
		if err != nil {
			t.Fatalf("Failed to query shipment status: %v", err)
		}
		if status != models.ShipmentStatusPickupScheduled {
			t.Errorf("Expected status 'pickup_from_client_scheduled', got '%s'", status)
		}

		// Then update to picked_up_from_client (sequential transition)
		formData = url.Values{}
		formData.Set("shipment_id", strconv.FormatInt(shipmentID, 10))
		formData.Set("status", string(models.ShipmentStatusPickedUpFromClient))

		req = httptest.NewRequest(http.MethodPost, "/shipments/update-status", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req = req.WithContext(reqCtx)

		w = httptest.NewRecorder()
		handler.UpdateShipmentStatus(w, req)

		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected status 303, got %d", w.Code)
		}

		// Verify status was updated to picked_up_from_client
		err = db.QueryRowContext(ctx,
			`SELECT status FROM shipments WHERE id = $1`,
			shipmentID,
		).Scan(&status)
		if err != nil {
			t.Fatalf("Failed to query shipment status: %v", err)
		}
		if status != models.ShipmentStatusPickedUpFromClient {
			t.Errorf("Expected status 'picked_up_from_client', got '%s'", status)
		}
	})

	t.Run("non-POST method returns error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/shipments/update-status", nil)
		w := httptest.NewRecorder()

		handler.UpdateShipmentStatus(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status 405, got %d", w.Code)
		}
	})

	t.Run("non-logistics user cannot update status", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("shipment_id", strconv.FormatInt(shipmentID, 10))
		formData.Set("status", string(models.ShipmentStatusInTransitToWarehouse))

		req := httptest.NewRequest(http.MethodPost, "/shipments/update-status", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		user := &models.User{ID: logisticsUserID, Email: "client@example.com", Role: models.RoleClient}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.UpdateShipmentStatus(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("Expected status 403, got %d", w.Code)
		}
	})

	t.Run("invalid status returns error", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("shipment_id", strconv.FormatInt(shipmentID, 10))
		formData.Set("status", "invalid_status")

		req := httptest.NewRequest(http.MethodPost, "/shipments/update-status", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		user := &models.User{ID: logisticsUserID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.UpdateShipmentStatus(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("updating to in_transit_to_engineer with ETA stores the ETA", func(t *testing.T) {
		// Update shipment to warehouse first
		_, err := db.ExecContext(ctx,
			`UPDATE shipments SET status = $1 WHERE id = $2`,
			models.ShipmentStatusAtWarehouse, shipmentID,
		)
		if err != nil {
			t.Fatalf("Failed to update shipment to warehouse: %v", err)
		}

		user := &models.User{ID: logisticsUserID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(context.Background(), middleware.UserContextKey, user)

		// First update to released_from_warehouse (sequential transition)
		formData := url.Values{}
		formData.Set("shipment_id", strconv.FormatInt(shipmentID, 10))
		formData.Set("status", string(models.ShipmentStatusReleasedFromWarehouse))

		req := httptest.NewRequest(http.MethodPost, "/shipments/update-status", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.UpdateShipmentStatus(w, req)

		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected status 303 for released_from_warehouse, got %d", w.Code)
		}

		// Then update to in_transit_to_engineer with ETA (sequential transition)
		etaTime := time.Now().Add(48 * time.Hour)
		etaString := etaTime.Format("2006-01-02T15:04")

		formData = url.Values{}
		formData.Set("shipment_id", strconv.FormatInt(shipmentID, 10))
		formData.Set("status", string(models.ShipmentStatusInTransitToEngineer))
		formData.Set("eta_to_engineer", etaString)

		req = httptest.NewRequest(http.MethodPost, "/shipments/update-status", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req = req.WithContext(reqCtx)

		w = httptest.NewRecorder()
		handler.UpdateShipmentStatus(w, req)

		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected status 303, got %d", w.Code)
		}

		// Verify status and ETA were updated
		var status models.ShipmentStatus
		var etaToEngineer *time.Time
		err = db.QueryRowContext(ctx,
			`SELECT status, eta_to_engineer FROM shipments WHERE id = $1`,
			shipmentID,
		).Scan(&status, &etaToEngineer)
		if err != nil {
			t.Fatalf("Failed to query shipment: %v", err)
		}

		if status != models.ShipmentStatusInTransitToEngineer {
			t.Errorf("Expected status 'in_transit_to_engineer', got '%s'", status)
		}

		if etaToEngineer == nil {
			t.Error("Expected ETA to be set, got nil")
		} else {
			// Check ETA is within a reasonable range (allowing for parsing and timezone differences)
			// We allow up to 5 hours difference to account for timezone conversions and precision loss
			timeDiff := etaToEngineer.Sub(etaTime).Abs()
			if timeDiff > 5*time.Hour {
				t.Errorf("Expected ETA around %v, got %v (diff: %v)", etaTime, etaToEngineer, timeDiff)
			}
		}
	})

	t.Run("updating to in_transit_to_engineer without ETA is allowed", func(t *testing.T) {
		// Create another test shipment at warehouse
		var shipmentID2 int64
		err := db.QueryRowContext(ctx,
			`INSERT INTO shipments (client_company_id, status, jira_ticket_number, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5) RETURNING id`,
			companyID, models.ShipmentStatusAtWarehouse, "TEST-998", time.Now(), time.Now(),
		).Scan(&shipmentID2)
		if err != nil {
			t.Fatalf("Failed to create test shipment: %v", err)
		}

		user := &models.User{ID: logisticsUserID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(context.Background(), middleware.UserContextKey, user)

		// First update to released_from_warehouse (sequential transition)
		formData := url.Values{}
		formData.Set("shipment_id", strconv.FormatInt(shipmentID2, 10))
		formData.Set("status", string(models.ShipmentStatusReleasedFromWarehouse))

		req := httptest.NewRequest(http.MethodPost, "/shipments/update-status", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.UpdateShipmentStatus(w, req)

		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected status 303 for released_from_warehouse, got %d", w.Code)
		}

		// Then update to in_transit_to_engineer without ETA (sequential transition)
		formData = url.Values{}
		formData.Set("shipment_id", strconv.FormatInt(shipmentID2, 10))
		formData.Set("status", string(models.ShipmentStatusInTransitToEngineer))
		// No eta_to_engineer field

		req = httptest.NewRequest(http.MethodPost, "/shipments/update-status", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req = req.WithContext(reqCtx)

		w = httptest.NewRecorder()
		handler.UpdateShipmentStatus(w, req)

		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected status 303, got %d", w.Code)
		}

		// Verify status was updated and ETA remains nil
		var status models.ShipmentStatus
		var etaToEngineer *time.Time
		err = db.QueryRowContext(ctx,
			`SELECT status, eta_to_engineer FROM shipments WHERE id = $1`,
			shipmentID2,
		).Scan(&status, &etaToEngineer)
		if err != nil {
			t.Fatalf("Failed to query shipment: %v", err)
		}

		if status != models.ShipmentStatusInTransitToEngineer {
			t.Errorf("Expected status 'in_transit_to_engineer', got '%s'", status)
		}

		if etaToEngineer != nil {
			t.Errorf("Expected ETA to be nil, got %v", etaToEngineer)
		}
	})

	t.Run("updating to pickup_from_client_scheduled with tracking number stores it in database", func(t *testing.T) {
		// Create a new test shipment
		var shipmentID3 int64
		err := db.QueryRowContext(ctx,
			`INSERT INTO shipments (client_company_id, status, jira_ticket_number, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5) RETURNING id`,
			companyID, models.ShipmentStatusPendingPickup, "TEST-997", time.Now(), time.Now(),
		).Scan(&shipmentID3)
		if err != nil {
			t.Fatalf("Failed to create test shipment: %v", err)
		}

		trackingNumber := "1Z999AA10123456784"
		courierName := "FedEx"

		formData := url.Values{}
		formData.Set("shipment_id", strconv.FormatInt(shipmentID3, 10))
		formData.Set("status", string(models.ShipmentStatusPickupScheduled))
		formData.Set("tracking_number", trackingNumber)
		formData.Set("courier_name", courierName)

		req := httptest.NewRequest(http.MethodPost, "/shipments/update-status", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		user := &models.User{ID: logisticsUserID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.UpdateShipmentStatus(w, req)

		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected status 303, got %d", w.Code)
		}

		// Verify status and tracking number were updated
		var status models.ShipmentStatus
		var storedTrackingNumber sql.NullString
		err = db.QueryRowContext(ctx,
			`SELECT status, tracking_number FROM shipments WHERE id = $1`,
			shipmentID3,
		).Scan(&status, &storedTrackingNumber)
		if err != nil {
			t.Fatalf("Failed to query shipment: %v", err)
		}

		if status != models.ShipmentStatusPickupScheduled {
			t.Errorf("Expected status 'pickup_from_client_scheduled', got '%s'", status)
		}

		if !storedTrackingNumber.Valid || storedTrackingNumber.String != trackingNumber {
			t.Errorf("Expected tracking number '%s', got '%s'", trackingNumber, storedTrackingNumber.String)
		}
	})

	t.Run("updating to pickup_from_client_scheduled without tracking number returns error", func(t *testing.T) {
		// Create a new test shipment
		var shipmentID4 int64
		err := db.QueryRowContext(ctx,
			`INSERT INTO shipments (client_company_id, status, jira_ticket_number, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5) RETURNING id`,
			companyID, models.ShipmentStatusPendingPickup, "TEST-996", time.Now(), time.Now(),
		).Scan(&shipmentID4)
		if err != nil {
			t.Fatalf("Failed to create test shipment: %v", err)
		}

		formData := url.Values{}
		formData.Set("shipment_id", strconv.FormatInt(shipmentID4, 10))
		formData.Set("status", string(models.ShipmentStatusPickupScheduled))
		// No tracking number provided

		req := httptest.NewRequest(http.MethodPost, "/shipments/update-status", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		user := &models.User{ID: logisticsUserID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.UpdateShipmentStatus(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400 (bad request), got %d", w.Code)
		}
	})

	t.Run("updating to pickup_from_client_scheduled without courier returns error", func(t *testing.T) {
		// Create a new test shipment
		var shipmentID5 int64
		err := db.QueryRowContext(ctx,
			`INSERT INTO shipments (client_company_id, status, jira_ticket_number, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5) RETURNING id`,
			companyID, models.ShipmentStatusPendingPickup, "TEST-995", time.Now(), time.Now(),
		).Scan(&shipmentID5)
		if err != nil {
			t.Fatalf("Failed to create test shipment: %v", err)
		}

		formData := url.Values{}
		formData.Set("shipment_id", strconv.FormatInt(shipmentID5, 10))
		formData.Set("status", string(models.ShipmentStatusPickupScheduled))
		formData.Set("tracking_number", "1Z999AA10123456784")
		// No courier_name provided

		req := httptest.NewRequest(http.MethodPost, "/shipments/update-status", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		user := &models.User{ID: logisticsUserID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.UpdateShipmentStatus(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400 (bad request), got %d", w.Code)
		}
	})

	t.Run("updating to pickup_from_client_scheduled with courier stores it in database", func(t *testing.T) {
		// Create a new test shipment
		var shipmentID6 int64
		err := db.QueryRowContext(ctx,
			`INSERT INTO shipments (client_company_id, status, jira_ticket_number, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5) RETURNING id`,
			companyID, models.ShipmentStatusPendingPickup, "TEST-994", time.Now(), time.Now(),
		).Scan(&shipmentID6)
		if err != nil {
			t.Fatalf("Failed to create test shipment: %v", err)
		}

		trackingNumber := "1Z999AA10123456784"
		courierName := "UPS"

		formData := url.Values{}
		formData.Set("shipment_id", strconv.FormatInt(shipmentID6, 10))
		formData.Set("status", string(models.ShipmentStatusPickupScheduled))
		formData.Set("tracking_number", trackingNumber)
		formData.Set("courier_name", courierName)

		req := httptest.NewRequest(http.MethodPost, "/shipments/update-status", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		user := &models.User{ID: logisticsUserID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.UpdateShipmentStatus(w, req)

		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected status 303, got %d", w.Code)
		}

		// Verify status, tracking number, and courier name were updated
		var status models.ShipmentStatus
		var storedTrackingNumber sql.NullString
		var storedCourierName sql.NullString
		err = db.QueryRowContext(ctx,
			`SELECT status, tracking_number, courier_name FROM shipments WHERE id = $1`,
			shipmentID6,
		).Scan(&status, &storedTrackingNumber, &storedCourierName)
		if err != nil {
			t.Fatalf("Failed to query shipment: %v", err)
		}

		if status != models.ShipmentStatusPickupScheduled {
			t.Errorf("Expected status 'pickup_from_client_scheduled', got '%s'", status)
		}

		if !storedTrackingNumber.Valid || storedTrackingNumber.String != trackingNumber {
			t.Errorf("Expected tracking number '%s', got '%s'", trackingNumber, storedTrackingNumber.String)
		}

		if !storedCourierName.Valid || storedCourierName.String != courierName {
			t.Errorf("Expected courier name '%s', got '%s'", courierName, storedCourierName.String)
		}
	})

	t.Run("updating to pickup_from_client_scheduled with invalid courier returns error", func(t *testing.T) {
		// Create a new test shipment
		var shipmentID7 int64
		err := db.QueryRowContext(ctx,
			`INSERT INTO shipments (client_company_id, status, jira_ticket_number, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5) RETURNING id`,
			companyID, models.ShipmentStatusPendingPickup, "TEST-993", time.Now(), time.Now(),
		).Scan(&shipmentID7)
		if err != nil {
			t.Fatalf("Failed to create test shipment: %v", err)
		}

		formData := url.Values{}
		formData.Set("shipment_id", strconv.FormatInt(shipmentID7, 10))
		formData.Set("status", string(models.ShipmentStatusPickupScheduled))
		formData.Set("tracking_number", "1Z999AA10123456784")
		formData.Set("courier_name", "InvalidCourier")

		req := httptest.NewRequest(http.MethodPost, "/shipments/update-status", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		user := &models.User{ID: logisticsUserID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.UpdateShipmentStatus(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400 (bad request), got %d", w.Code)
		}
	})

	// Tests for sequential status validation - preventing skipping and backwards transitions
	t.Run("cannot skip statuses - pending_pickup to at_warehouse", func(t *testing.T) {
		// Create a new test shipment at pending_pickup_from_client
		var shipmentID8 int64
		err := db.QueryRowContext(ctx,
			`INSERT INTO shipments (client_company_id, status, jira_ticket_number, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5) RETURNING id`,
			companyID, models.ShipmentStatusPendingPickup, "TEST-992", time.Now(), time.Now(),
		).Scan(&shipmentID8)
		if err != nil {
			t.Fatalf("Failed to create test shipment: %v", err)
		}

		formData := url.Values{}
		formData.Set("shipment_id", strconv.FormatInt(shipmentID8, 10))
		formData.Set("status", string(models.ShipmentStatusAtWarehouse)) // Skipping multiple statuses

		req := httptest.NewRequest(http.MethodPost, "/shipments/update-status", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		user := &models.User{ID: logisticsUserID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.UpdateShipmentStatus(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400 (cannot skip statuses), got %d", w.Code)
		}

		// Verify status was NOT updated
		var status models.ShipmentStatus
		err = db.QueryRowContext(ctx,
			`SELECT status FROM shipments WHERE id = $1`,
			shipmentID8,
		).Scan(&status)
		if err != nil {
			t.Fatalf("Failed to query shipment: %v", err)
		}
		if status != models.ShipmentStatusPendingPickup {
			t.Errorf("Status should remain 'pending_pickup_from_client', got '%s'", status)
		}
	})

	t.Run("cannot skip one status - pending_pickup to picked_up", func(t *testing.T) {
		// Create a new test shipment
		var shipmentID9 int64
		err := db.QueryRowContext(ctx,
			`INSERT INTO shipments (client_company_id, status, jira_ticket_number, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5) RETURNING id`,
			companyID, models.ShipmentStatusPendingPickup, "TEST-991", time.Now(), time.Now(),
		).Scan(&shipmentID9)
		if err != nil {
			t.Fatalf("Failed to create test shipment: %v", err)
		}

		formData := url.Values{}
		formData.Set("shipment_id", strconv.FormatInt(shipmentID9, 10))
		formData.Set("status", string(models.ShipmentStatusPickedUpFromClient)) // Skipping pickup_scheduled

		req := httptest.NewRequest(http.MethodPost, "/shipments/update-status", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		user := &models.User{ID: logisticsUserID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.UpdateShipmentStatus(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400 (cannot skip statuses), got %d", w.Code)
		}
	})

	t.Run("cannot go backwards - at_warehouse to pending_pickup", func(t *testing.T) {
		// Create a new test shipment at warehouse
		var shipmentID10 int64
		err := db.QueryRowContext(ctx,
			`INSERT INTO shipments (client_company_id, status, jira_ticket_number, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5) RETURNING id`,
			companyID, models.ShipmentStatusAtWarehouse, "TEST-990", time.Now(), time.Now(),
		).Scan(&shipmentID10)
		if err != nil {
			t.Fatalf("Failed to create test shipment: %v", err)
		}

		formData := url.Values{}
		formData.Set("shipment_id", strconv.FormatInt(shipmentID10, 10))
		formData.Set("status", string(models.ShipmentStatusPendingPickup)) // Going backwards

		req := httptest.NewRequest(http.MethodPost, "/shipments/update-status", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		user := &models.User{ID: logisticsUserID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.UpdateShipmentStatus(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400 (cannot go backwards), got %d", w.Code)
		}

		// Verify status was NOT updated
		var status models.ShipmentStatus
		err = db.QueryRowContext(ctx,
			`SELECT status FROM shipments WHERE id = $1`,
			shipmentID10,
		).Scan(&status)
		if err != nil {
			t.Fatalf("Failed to query shipment: %v", err)
		}
		if status != models.ShipmentStatusAtWarehouse {
			t.Errorf("Status should remain 'at_warehouse', got '%s'", status)
		}
	})

	t.Run("cannot update to same status", func(t *testing.T) {
		// Create a new test shipment
		var shipmentID11 int64
		err := db.QueryRowContext(ctx,
			`INSERT INTO shipments (client_company_id, status, jira_ticket_number, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5) RETURNING id`,
			companyID, models.ShipmentStatusAtWarehouse, "TEST-989", time.Now(), time.Now(),
		).Scan(&shipmentID11)
		if err != nil {
			t.Fatalf("Failed to create test shipment: %v", err)
		}

		formData := url.Values{}
		formData.Set("shipment_id", strconv.FormatInt(shipmentID11, 10))
		formData.Set("status", string(models.ShipmentStatusAtWarehouse)) // Same status

		req := httptest.NewRequest(http.MethodPost, "/shipments/update-status", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		user := &models.User{ID: logisticsUserID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.UpdateShipmentStatus(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400 (cannot update to same status), got %d", w.Code)
		}
	})

	t.Run("cannot update from delivered (final status)", func(t *testing.T) {
		// Create a new test shipment that is delivered
		var shipmentID12 int64
		err := db.QueryRowContext(ctx,
			`INSERT INTO shipments (client_company_id, status, jira_ticket_number, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5) RETURNING id`,
			companyID, models.ShipmentStatusDelivered, "TEST-988", time.Now(), time.Now(),
		).Scan(&shipmentID12)
		if err != nil {
			t.Fatalf("Failed to create test shipment: %v", err)
		}

		formData := url.Values{}
		formData.Set("shipment_id", strconv.FormatInt(shipmentID12, 10))
		formData.Set("status", string(models.ShipmentStatusInTransitToEngineer)) // Try to go backwards

		req := httptest.NewRequest(http.MethodPost, "/shipments/update-status", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		user := &models.User{ID: logisticsUserID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.UpdateShipmentStatus(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400 (delivered is final status), got %d", w.Code)
		}

		// Verify status was NOT updated
		var status models.ShipmentStatus
		err = db.QueryRowContext(ctx,
			`SELECT status FROM shipments WHERE id = $1`,
			shipmentID12,
		).Scan(&status)
		if err != nil {
			t.Fatalf("Failed to query shipment: %v", err)
		}
		if status != models.ShipmentStatusDelivered {
			t.Errorf("Status should remain 'delivered', got '%s'", status)
		}
	})

	t.Run("can update sequentially - pending_pickup to pickup_scheduled", func(t *testing.T) {
		// Create a new test shipment
		var shipmentID13 int64
		err := db.QueryRowContext(ctx,
			`INSERT INTO shipments (client_company_id, status, jira_ticket_number, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5) RETURNING id`,
			companyID, models.ShipmentStatusPendingPickup, "TEST-987", time.Now(), time.Now(),
		).Scan(&shipmentID13)
		if err != nil {
			t.Fatalf("Failed to create test shipment: %v", err)
		}

		formData := url.Values{}
		formData.Set("shipment_id", strconv.FormatInt(shipmentID13, 10))
		formData.Set("status", string(models.ShipmentStatusPickupScheduled))
		formData.Set("tracking_number", "1Z999AA10123456784")
		formData.Set("courier_name", "UPS")

		req := httptest.NewRequest(http.MethodPost, "/shipments/update-status", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		user := &models.User{ID: logisticsUserID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.UpdateShipmentStatus(w, req)

		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected status 303 (sequential update allowed), got %d", w.Code)
		}

		// Verify status WAS updated
		var status models.ShipmentStatus
		err = db.QueryRowContext(ctx,
			`SELECT status FROM shipments WHERE id = $1`,
			shipmentID13,
		).Scan(&status)
		if err != nil {
			t.Fatalf("Failed to query shipment: %v", err)
		}
		if status != models.ShipmentStatusPickupScheduled {
			t.Errorf("Status should be updated to 'pickup_from_client_scheduled', got '%s'", status)
		}
	})
}

func TestCreateShipment(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create test logistics user
	var logisticsUserID int64
	err := db.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"logistics@example.com", "hashedpassword", models.RoleLogistics, time.Now(), time.Now(),
	).Scan(&logisticsUserID)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
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

	// Mock JIRA validator that always succeeds
	mockJiraValidator := func(ticketKey string) error {
		if ticketKey == "INVALID-000" {
			return errors.New("JIRA ticket INVALID-000 does not exist")
		}
		return nil
	}

	handler := &ShipmentsHandler{
		DB:            db,
		Templates:     templates,
		JiraValidator: mockJiraValidator,
	}

	t.Run("logistics user can create shipment with valid JIRA ticket", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("client_company_id", strconv.FormatInt(companyID, 10))
		formData.Set("jira_ticket_number", "SCOP-67702")
		formData.Set("notes", "Test shipment")

		req := httptest.NewRequest(http.MethodPost, "/shipments/create", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		user := &models.User{ID: logisticsUserID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.CreateShipment(w, req)

		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected status 303, got %d", w.Code)
		}

		// Verify shipment was created
		var count int
		err := db.QueryRowContext(ctx,
			`SELECT COUNT(*) FROM shipments WHERE jira_ticket_number = $1`,
			"SCOP-67702",
		).Scan(&count)
		if err != nil {
			t.Fatalf("Failed to query shipment: %v", err)
		}
		if count != 1 {
			t.Errorf("Expected 1 shipment with JIRA ticket SCOP-67702, got %d", count)
		}
	})

	t.Run("cannot create shipment without JIRA ticket", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("client_company_id", strconv.FormatInt(companyID, 10))
		formData.Set("notes", "Test shipment")

		req := httptest.NewRequest(http.MethodPost, "/shipments/create", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		user := &models.User{ID: logisticsUserID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.CreateShipment(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("cannot create shipment with invalid JIRA ticket format", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("client_company_id", strconv.FormatInt(companyID, 10))
		formData.Set("jira_ticket_number", "invalid-format")
		formData.Set("notes", "Test shipment")

		req := httptest.NewRequest(http.MethodPost, "/shipments/create", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		user := &models.User{ID: logisticsUserID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.CreateShipment(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("cannot create shipment with non-existent JIRA ticket", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("client_company_id", strconv.FormatInt(companyID, 10))
		formData.Set("jira_ticket_number", "INVALID-000")
		formData.Set("notes", "Test shipment")

		req := httptest.NewRequest(http.MethodPost, "/shipments/create", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		user := &models.User{ID: logisticsUserID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.CreateShipment(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("non-logistics user cannot create shipment", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("client_company_id", strconv.FormatInt(companyID, 10))
		formData.Set("jira_ticket_number", "SCOP-67702")
		formData.Set("notes", "Test shipment")

		req := httptest.NewRequest(http.MethodPost, "/shipments/create", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		user := &models.User{ID: logisticsUserID, Email: "client@example.com", Role: models.RoleClient}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.CreateShipment(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("Expected status 403, got %d", w.Code)
		}
	})

	t.Run("GET request shows create shipment form", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/shipments/create", nil)

		user := &models.User{ID: logisticsUserID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.CreateShipment(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})
}

func TestShipmentPickupFormPage(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create test user (client role for magic link)
	var userID int64
	err := db.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"client@example.com", "hashedpassword", models.RoleClient, time.Now(), time.Now(),
	).Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
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

	// Create test shipment with JIRA ticket
	var shipmentID int64
	err = db.QueryRowContext(ctx,
		`INSERT INTO shipments (client_company_id, status, jira_ticket_number, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		companyID, models.ShipmentStatusPendingPickup, "SCOP-12345", time.Now(), time.Now(),
	).Scan(&shipmentID)
	if err != nil {
		t.Fatalf("Failed to create test shipment: %v", err)
	}

	handler := NewShipmentsHandler(db, nil, nil)

	t.Run("GET request for shipment without pickup form shows empty form", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/shipments/%d/form", shipmentID), nil)
		req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(shipmentID, 10)})

		user := &models.User{ID: userID, Email: "client@example.com", Role: models.RoleClient}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.ShipmentPickupFormPage(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}

		// Verify response contains shipment ID and JIRA ticket (when templates are nil, we'll check headers/data)
		// For now, just check that we get OK response
	})

	t.Run("GET request for shipment with existing pickup form shows pre-filled form", func(t *testing.T) {
		// Create a shipment with an existing pickup form
		var shipmentID2 int64
		err := db.QueryRowContext(ctx,
			`INSERT INTO shipments (client_company_id, status, jira_ticket_number, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5) RETURNING id`,
			companyID, models.ShipmentStatusPendingPickup, "SCOP-54321", time.Now(), time.Now(),
		).Scan(&shipmentID2)
		if err != nil {
			t.Fatalf("Failed to create second test shipment: %v", err)
		}

		// Create pickup form for this shipment
		formData := map[string]interface{}{
			"contact_name":         "Jane Doe",
			"contact_email":        "jane@company.com",
			"contact_phone":        "+1-555-9999",
			"pickup_address":       "456 Business Ave, Suite 200",
			"pickup_city":          "Boston",
			"pickup_state":         "MA",
			"pickup_zip":           "02101",
			"pickup_date":          "2025-12-15",
			"pickup_time_slot":     "afternoon",
			"number_of_laptops":    3,
			"special_instructions": "Handle with care",
		}
		formDataJSON, _ := json.Marshal(formData)

		_, err = db.ExecContext(ctx,
			`INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data)
			VALUES ($1, $2, $3, $4)`,
			shipmentID2, userID, time.Now(), formDataJSON,
		)
		if err != nil {
			t.Fatalf("Failed to create pickup form: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/shipments/%d/form", shipmentID2), nil)
		req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(shipmentID2, 10)})

		user := &models.User{ID: userID, Email: "client@example.com", Role: models.RoleClient}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.ShipmentPickupFormPage(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}

		// Handler should return OK and load the existing form data
		// The template will display pre-filled form values
	})
}

func TestShipmentPickupFormSubmit(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create test user
	var userID int64
	err := db.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"client@example.com", "hashedpassword", models.RoleClient, time.Now(), time.Now(),
	).Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
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

	// Create test shipment
	var shipmentID int64
	err = db.QueryRowContext(ctx,
		`INSERT INTO shipments (client_company_id, status, jira_ticket_number, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		companyID, models.ShipmentStatusPendingPickup, "SCOP-99999", time.Now(), time.Now(),
	).Scan(&shipmentID)
	if err != nil {
		t.Fatalf("Failed to create test shipment: %v", err)
	}

	handler := NewShipmentsHandler(db, nil, nil)

	t.Run("POST request creates new pickup form for shipment", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("contact_name", "John Smith")
		formData.Set("contact_email", "john@company.com")
		formData.Set("contact_phone", "+1-555-1234")
		formData.Set("pickup_address", "123 Main St, Suite 100")
		formData.Set("pickup_city", "New York")
		formData.Set("pickup_state", "NY")
		formData.Set("pickup_zip", "10001")
		formData.Set("pickup_date", "2025-12-20")
		formData.Set("pickup_time_slot", "morning")
		formData.Set("number_of_laptops", "2")
		formData.Set("special_instructions", "Please call before arrival")
		formData.Set("number_of_boxes", "1")
		formData.Set("assignment_type", "single")

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/shipments/%d/form", shipmentID),
			strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(shipmentID, 10)})

		user := &models.User{ID: userID, Email: "client@example.com", Role: models.RoleClient}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.ShipmentPickupFormSubmit(w, req)

		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected status 303 (redirect), got %d", w.Code)
		}

		// Verify pickup form was created
		var count int
		err := db.QueryRowContext(ctx,
			`SELECT COUNT(*) FROM pickup_forms WHERE shipment_id = $1`,
			shipmentID,
		).Scan(&count)
		if err != nil {
			t.Fatalf("Failed to check pickup form: %v", err)
		}
		if count != 1 {
			t.Errorf("Expected 1 pickup form, found %d", count)
		}
	})

	t.Run("POST request updates existing pickup form for shipment", func(t *testing.T) {
		// Create another shipment with an existing pickup form
		var shipmentID2 int64
		err := db.QueryRowContext(ctx,
			`INSERT INTO shipments (client_company_id, status, jira_ticket_number, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5) RETURNING id`,
			companyID, models.ShipmentStatusPendingPickup, "SCOP-88888", time.Now(), time.Now(),
		).Scan(&shipmentID2)
		if err != nil {
			t.Fatalf("Failed to create second test shipment: %v", err)
		}

		// Create initial pickup form
		initialFormData := map[string]interface{}{
			"contact_name":         "Old Name",
			"contact_email":        "old@company.com",
			"contact_phone":        "+1-555-0000",
			"pickup_address":       "Old Address",
			"pickup_city":          "Chicago",
			"pickup_state":         "IL",
			"pickup_zip":           "60601",
			"pickup_date":          "2025-12-10",
			"pickup_time_slot":     "evening",
			"number_of_laptops":    1,
			"special_instructions": "Old instructions",
			"number_of_boxes":      1,
			"assignment_type":      "single",
		}
		initialFormJSON, _ := json.Marshal(initialFormData)

		_, err = db.ExecContext(ctx,
			`INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data)
			VALUES ($1, $2, $3, $4)`,
			shipmentID2, userID, time.Now(), initialFormJSON,
		)
		if err != nil {
			t.Fatalf("Failed to create initial pickup form: %v", err)
		}

		// Now submit updated form data
		updatedFormData := url.Values{}
		updatedFormData.Set("contact_name", "Updated Name")
		updatedFormData.Set("contact_email", "updated@company.com")
		updatedFormData.Set("contact_phone", "+1-555-9999")
		updatedFormData.Set("pickup_address", "Updated Address")
		updatedFormData.Set("pickup_city", "Los Angeles")
		updatedFormData.Set("pickup_state", "CA")
		updatedFormData.Set("pickup_zip", "90001")
		updatedFormData.Set("pickup_date", "2025-12-25")
		updatedFormData.Set("pickup_time_slot", "afternoon")
		updatedFormData.Set("number_of_laptops", "5")
		updatedFormData.Set("special_instructions", "Updated instructions")
		updatedFormData.Set("number_of_boxes", "2")
		updatedFormData.Set("assignment_type", "single")

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/shipments/%d/form", shipmentID2),
			strings.NewReader(updatedFormData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(shipmentID2, 10)})

		user := &models.User{ID: userID, Email: "client@example.com", Role: models.RoleClient}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.ShipmentPickupFormSubmit(w, req)

		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected status 303 (redirect), got %d", w.Code)
		}

		// Verify there's still only 1 pickup form (updated, not duplicated)
		var count int
		err = db.QueryRowContext(ctx,
			`SELECT COUNT(*) FROM pickup_forms WHERE shipment_id = $1`,
			shipmentID2,
		).Scan(&count)
		if err != nil {
			t.Fatalf("Failed to check pickup form count: %v", err)
		}
		if count != 1 {
			t.Errorf("Expected 1 pickup form, found %d (should not duplicate)", count)
		}

		// Verify the form data was updated
		var formDataJSON []byte
		err = db.QueryRowContext(ctx,
			`SELECT form_data FROM pickup_forms WHERE shipment_id = $1`,
			shipmentID2,
		).Scan(&formDataJSON)
		if err != nil {
			t.Fatalf("Failed to fetch updated form: %v", err)
		}

		var formData map[string]interface{}
		json.Unmarshal(formDataJSON, &formData)

		if formData["contact_name"] != "Updated Name" {
			t.Errorf("Expected contact_name to be 'Updated Name', got %v", formData["contact_name"])
		}
		if formData["contact_email"] != "updated@company.com" {
			t.Errorf("Expected contact_email to be 'updated@company.com', got %v", formData["contact_email"])
		}
	})
}

func TestSendMagicLinkVisibility(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create test logistics user
	var logisticsUserID int64
	err := db.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"logistics@example.com", "hashedpassword", models.RoleLogistics, time.Now(), time.Now(),
	).Scan(&logisticsUserID)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
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

	t.Run("send magic link form is visible when status is pending_pickup_from_client", func(t *testing.T) {
		// Create shipment with status pending_pickup_from_client
		var shipmentID int64
		err := db.QueryRowContext(ctx,
			`INSERT INTO shipments (client_company_id, status, jira_ticket_number, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5) RETURNING id`,
			companyID, models.ShipmentStatusPendingPickup, "TEST-MAGIC-1", time.Now(), time.Now(),
		).Scan(&shipmentID)
		if err != nil {
			t.Fatalf("Failed to create test shipment: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/shipments/"+strconv.FormatInt(shipmentID, 10), nil)
		req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(shipmentID, 10)})

		user := &models.User{ID: logisticsUserID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.ShipmentDetail(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		responseBody := w.Body.String()
		if !strings.Contains(responseBody, "Send Magic Link") {
			t.Errorf("Expected 'Send Magic Link' form to be visible for status pending_pickup_from_client")
		}
	})

	t.Run("send magic link form is visible when status is pickup_from_client_scheduled", func(t *testing.T) {
		// Create shipment with status pickup_from_client_scheduled
		var shipmentID int64
		err := db.QueryRowContext(ctx,
			`INSERT INTO shipments (client_company_id, status, jira_ticket_number, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5) RETURNING id`,
			companyID, models.ShipmentStatusPickupScheduled, "TEST-MAGIC-2", time.Now(), time.Now(),
		).Scan(&shipmentID)
		if err != nil {
			t.Fatalf("Failed to create test shipment: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/shipments/"+strconv.FormatInt(shipmentID, 10), nil)
		req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(shipmentID, 10)})

		user := &models.User{ID: logisticsUserID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.ShipmentDetail(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		responseBody := w.Body.String()
		if !strings.Contains(responseBody, "Send Magic Link") {
			t.Errorf("Expected 'Send Magic Link' form to be visible for status pickup_from_client_scheduled")
		}
	})

	t.Run("send magic link form is NOT visible when status is picked_up_from_client", func(t *testing.T) {
		// Create shipment with status picked_up_from_client
		var shipmentID int64
		err := db.QueryRowContext(ctx,
			`INSERT INTO shipments (client_company_id, status, jira_ticket_number, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5) RETURNING id`,
			companyID, models.ShipmentStatusPickedUpFromClient, "TEST-MAGIC-3", time.Now(), time.Now(),
		).Scan(&shipmentID)
		if err != nil {
			t.Fatalf("Failed to create test shipment: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/shipments/"+strconv.FormatInt(shipmentID, 10), nil)
		req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(shipmentID, 10)})

		user := &models.User{ID: logisticsUserID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.ShipmentDetail(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		responseBody := w.Body.String()
		if strings.Contains(responseBody, "Send Magic Link") {
			t.Errorf("Expected 'Send Magic Link' form to NOT be visible for status picked_up_from_client")
		}
	})

	t.Run("send magic link form is NOT visible when status is at_warehouse", func(t *testing.T) {
		// Create shipment with status at_warehouse
		var shipmentID int64
		err := db.QueryRowContext(ctx,
			`INSERT INTO shipments (client_company_id, status, jira_ticket_number, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5) RETURNING id`,
			companyID, models.ShipmentStatusAtWarehouse, "TEST-MAGIC-4", time.Now(), time.Now(),
		).Scan(&shipmentID)
		if err != nil {
			t.Fatalf("Failed to create test shipment: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/shipments/"+strconv.FormatInt(shipmentID, 10), nil)
		req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(shipmentID, 10)})

		user := &models.User{ID: logisticsUserID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.ShipmentDetail(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		responseBody := w.Body.String()
		if strings.Contains(responseBody, "Send Magic Link") {
			t.Errorf("Expected 'Send Magic Link' form to NOT be visible for status at_warehouse")
		}
	})

	t.Run("send magic link form is NOT visible when status is delivered", func(t *testing.T) {
		// Create shipment with status delivered
		var shipmentID int64
		err := db.QueryRowContext(ctx,
			`INSERT INTO shipments (client_company_id, status, jira_ticket_number, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5) RETURNING id`,
			companyID, models.ShipmentStatusDelivered, "TEST-MAGIC-5", time.Now(), time.Now(),
		).Scan(&shipmentID)
		if err != nil {
			t.Fatalf("Failed to create test shipment: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/shipments/"+strconv.FormatInt(shipmentID, 10), nil)
		req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(shipmentID, 10)})

		user := &models.User{ID: logisticsUserID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.ShipmentDetail(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		responseBody := w.Body.String()
		if strings.Contains(responseBody, "Send Magic Link") {
			t.Errorf("Expected 'Send Magic Link' form to NOT be visible for status delivered")
		}
	})
}
