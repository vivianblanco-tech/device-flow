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
	"github.com/yourusername/laptop-tracking-system/internal/models"
)

func TestDeliveryFormPage(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create test company
	var companyID int64
	err := db.QueryRowContext(ctx,
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
		companyID, models.ShipmentStatusInTransitToEngineer, time.Now(), time.Now(),
	).Scan(&shipmentID)
	if err != nil {
		t.Fatalf("Failed to create test shipment: %v", err)
	}

	templates := loadTestTemplates(t)
	handler := NewDeliveryFormHandler(db, templates, nil)

	t.Run("GET request displays form", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/delivery-form?shipment_id="+strconv.FormatInt(shipmentID, 10), nil)
		w := httptest.NewRecorder()

		handler.DeliveryFormPage(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})

	t.Run("missing shipment ID returns error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/delivery-form", nil)
		w := httptest.NewRecorder()

		handler.DeliveryFormPage(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("invalid shipment ID returns error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/delivery-form?shipment_id=invalid", nil)
		w := httptest.NewRecorder()

		handler.DeliveryFormPage(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("non-existent shipment returns not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/delivery-form?shipment_id=99999", nil)
		w := httptest.NewRecorder()

		handler.DeliveryFormPage(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", w.Code)
		}
	})
}

func TestDeliveryFormSubmit(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create test company
	var companyID int64
	err := db.QueryRowContext(ctx,
		`INSERT INTO client_companies (name, contact_info, created_at)
		VALUES ($1, $2, $3) RETURNING id`,
		"Test Company", json.RawMessage(`{"email":"test@company.com"}`), time.Now(),
	).Scan(&companyID)
	if err != nil {
		t.Fatalf("Failed to create test company: %v", err)
	}

	// Create test engineer
	var engineerID int64
	err = db.QueryRowContext(ctx,
		`INSERT INTO software_engineers (name, email, address, phone, address_confirmed, created_at)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		"John Engineer", "john@example.com", "123 Main St", "+1-555-0123", true, time.Now(),
	).Scan(&engineerID)
	if err != nil {
		t.Fatalf("Failed to create test engineer: %v", err)
	}

	// Create test shipment
	var shipmentID int64
	err = db.QueryRowContext(ctx,
		`INSERT INTO shipments (client_company_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4) RETURNING id`,
		companyID, models.ShipmentStatusInTransitToEngineer, time.Now(), time.Now(),
	).Scan(&shipmentID)
	if err != nil {
		t.Fatalf("Failed to create test shipment: %v", err)
	}

	templates := loadTestTemplates(t)
	handler := NewDeliveryFormHandler(db, templates, nil)

	t.Run("valid form submission creates delivery record", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("shipment_id", strconv.FormatInt(shipmentID, 10))
		formData.Set("engineer_id", strconv.FormatInt(engineerID, 10))
		formData.Set("notes", "Delivered successfully")

		req := httptest.NewRequest(http.MethodPost, "/delivery-form", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		w := httptest.NewRecorder()
		handler.DeliveryFormSubmit(w, req)

		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected status 303, got %d", w.Code)
		}

		// Verify delivery form was created
		var count int
		err := db.QueryRowContext(ctx,
			`SELECT COUNT(*) FROM delivery_forms WHERE shipment_id = $1`,
			shipmentID,
		).Scan(&count)
		if err != nil {
			t.Fatalf("Failed to query delivery forms: %v", err)
		}
		if count != 1 {
			t.Errorf("Expected 1 delivery form, got %d", count)
		}

		// Verify shipment status was updated
		var status models.ShipmentStatus
		err = db.QueryRowContext(ctx,
			`SELECT status FROM shipments WHERE id = $1`,
			shipmentID,
		).Scan(&status)
		if err != nil {
			t.Fatalf("Failed to query shipment status: %v", err)
		}
		if status != models.ShipmentStatusDelivered {
			t.Errorf("Expected shipment status 'delivered', got '%s'", status)
		}
	})

	t.Run("non-POST method returns error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/delivery-form", nil)
		w := httptest.NewRecorder()

		handler.DeliveryFormSubmit(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status 405, got %d", w.Code)
		}
	})

	t.Run("invalid shipment ID returns error", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("shipment_id", "invalid")
		formData.Set("engineer_id", strconv.FormatInt(engineerID, 10))

		req := httptest.NewRequest(http.MethodPost, "/delivery-form", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		w := httptest.NewRecorder()
		handler.DeliveryFormSubmit(w, req)

		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected status 303 (redirect with error), got %d", w.Code)
		}
	})
}

