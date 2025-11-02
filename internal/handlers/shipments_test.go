package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

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

	for _, status := range statuses {
		_, err := db.ExecContext(ctx,
			`INSERT INTO shipments (client_company_id, status, tracking_number, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5)`,
			companyID, status, "TRACK-"+string(status), time.Now(), time.Now(),
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
		`INSERT INTO shipments (client_company_id, status, tracking_number, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		companyID, models.ShipmentStatusInTransitToWarehouse, "TRACK-12345", time.Now(), time.Now(),
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
		req := httptest.NewRequest(http.MethodGet, "/shipments/detail?id="+strconv.FormatInt(shipmentID, 10), nil)
		
		user := &models.User{ID: userID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.ShipmentDetail(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
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
		req := httptest.NewRequest(http.MethodGet, "/shipments/detail?id=invalid", nil)
		
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
		req := httptest.NewRequest(http.MethodGet, "/shipments/detail?id=99999", nil)
		
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
		`INSERT INTO shipments (client_company_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4) RETURNING id`,
		companyID, models.ShipmentStatusPendingPickup, time.Now(), time.Now(),
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

