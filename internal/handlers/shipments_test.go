package handlers

import (
	"context"
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
	handler := NewShipmentsHandler(db, templates)

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
		`INSERT INTO shipments (client_company_id, status, jira_ticket_number, tracking_number, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		companyID, models.ShipmentStatusInTransitToWarehouse, "TEST-12345", "TRACK-12345", time.Now(), time.Now(),
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
	handler := NewShipmentsHandler(db, templates)

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
	handler := NewShipmentsHandler(db, templates)

	t.Run("logistics user can update shipment status", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("shipment_id", strconv.FormatInt(shipmentID, 10))
		formData.Set("status", string(models.ShipmentStatusPickedUpFromClient))

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

		// Verify status was updated
		var status models.ShipmentStatus
		err := db.QueryRowContext(ctx,
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

