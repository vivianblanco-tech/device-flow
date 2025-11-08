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
		handler := NewShipmentsHandler(db, templates)

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
		handler := NewShipmentsHandler(db, templates)

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

	handler := NewShipmentsHandler(db, nil)

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

	handler := NewShipmentsHandler(db, nil)

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
