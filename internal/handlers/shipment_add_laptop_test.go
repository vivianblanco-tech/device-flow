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

	"github.com/gorilla/mux"
	"github.com/yourusername/laptop-tracking-system/internal/database"
	"github.com/yourusername/laptop-tracking-system/internal/middleware"
	"github.com/yourusername/laptop-tracking-system/internal/models"
)

func TestAddLaptopToBulkShipment(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	templates := loadTestTemplates(t)
	handler := NewShipmentsHandler(db, templates, nil)

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

	// Create bulk shipment
	var shipmentID int64
	err = db.QueryRowContext(ctx,
		`INSERT INTO shipments (shipment_type, client_company_id, status, laptop_count, jira_ticket_number, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		models.ShipmentTypeBulkToWarehouse, companyID, models.ShipmentStatusPendingPickup, 5, "TEST-001", time.Now(), time.Now(),
	).Scan(&shipmentID)
	if err != nil {
		t.Fatalf("Failed to create bulk shipment: %v", err)
	}

	// Create logistics user
	var logisticsUserID int64
	err = db.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"logistics@test.com", "hash", models.RoleLogistics, time.Now(), time.Now(),
	).Scan(&logisticsUserID)
	if err != nil {
		t.Fatalf("Failed to create logistics user: %v", err)
	}

	// Create warehouse user (should not be able to add laptops)
	var warehouseUserID int64
	err = db.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"warehouse@test.com", "hash", models.RoleWarehouse, time.Now(), time.Now(),
	).Scan(&warehouseUserID)
	if err != nil {
		t.Fatalf("Failed to create warehouse user: %v", err)
	}

	t.Run("logistics user can add eligible laptop to bulk shipment", func(t *testing.T) {
		// Create eligible laptop (In Transit to Warehouse, not in any shipment)
		var laptopID int64
		err = db.QueryRowContext(ctx,
			`INSERT INTO laptops (serial_number, brand, model, cpu, ram_gb, ssd_gb, status, client_company_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`,
			"ELIGIBLE-001", "Dell", "Latitude", "i7", "16", "512", models.LaptopStatusInTransitToWarehouse, companyID, time.Now(), time.Now(),
		).Scan(&laptopID)
		if err != nil {
			t.Fatalf("Failed to create eligible laptop: %v", err)
		}

		formData := url.Values{}
		formData.Set("laptop_id", strconv.FormatInt(laptopID, 10))

		req := httptest.NewRequest(http.MethodPost, "/shipments/"+strconv.FormatInt(shipmentID, 10)+"/laptops/add", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		user := &models.User{ID: logisticsUserID, Email: "logistics@test.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		// Set up mux vars
		req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(shipmentID, 10)})

		w := httptest.NewRecorder()
		handler.AddLaptopToBulkShipment(w, req)

		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected status 303, got %d. Body: %s", w.Code, w.Body.String())
		}

		// Verify laptop was linked to shipment
		var count int
		err = db.QueryRowContext(ctx,
			`SELECT COUNT(*) FROM shipment_laptops WHERE shipment_id = $1 AND laptop_id = $2`,
			shipmentID, laptopID,
		).Scan(&count)
		if err != nil {
			t.Fatalf("Failed to check laptop link: %v", err)
		}
		if count != 1 {
			t.Errorf("Expected laptop to be linked to shipment, but count is %d", count)
		}
	})

	t.Run("warehouse user cannot add laptop to bulk shipment", func(t *testing.T) {
		// Create eligible laptop
		var laptopID int64
		err = db.QueryRowContext(ctx,
			`INSERT INTO laptops (serial_number, brand, model, cpu, ram_gb, ssd_gb, status, client_company_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`,
			"ELIGIBLE-002", "Dell", "Latitude", "i7", "16", "512", models.LaptopStatusInTransitToWarehouse, companyID, time.Now(), time.Now(),
		).Scan(&laptopID)
		if err != nil {
			t.Fatalf("Failed to create eligible laptop: %v", err)
		}

		formData := url.Values{}
		formData.Set("laptop_id", strconv.FormatInt(laptopID, 10))

		req := httptest.NewRequest(http.MethodPost, "/shipments/"+strconv.FormatInt(shipmentID, 10)+"/laptops/add", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		user := &models.User{ID: warehouseUserID, Email: "warehouse@test.com", Role: models.RoleWarehouse}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(shipmentID, 10)})

		w := httptest.NewRecorder()
		handler.AddLaptopToBulkShipment(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("Expected status 403, got %d. Body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("cannot add laptop with wrong status to bulk shipment", func(t *testing.T) {
		// Create laptop with wrong status (Available instead of In Transit to Warehouse)
		var laptopID int64
		err = db.QueryRowContext(ctx,
			`INSERT INTO laptops (serial_number, brand, model, cpu, ram_gb, ssd_gb, status, client_company_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`,
			"WRONG-STATUS-001", "Dell", "Latitude", "i7", "16", "512", models.LaptopStatusAvailable, companyID, time.Now(), time.Now(),
		).Scan(&laptopID)
		if err != nil {
			t.Fatalf("Failed to create laptop: %v", err)
		}

		formData := url.Values{}
		formData.Set("laptop_id", strconv.FormatInt(laptopID, 10))

		req := httptest.NewRequest(http.MethodPost, "/shipments/"+strconv.FormatInt(shipmentID, 10)+"/laptops/add", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		user := &models.User{ID: logisticsUserID, Email: "logistics@test.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(shipmentID, 10)})

		w := httptest.NewRecorder()
		handler.AddLaptopToBulkShipment(w, req)

		// Handler redirects on error with error message in query string
		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected status 303 (redirect), got %d. Body: %s", w.Code, w.Body.String())
		}
		location := w.Header().Get("Location")
		if location == "" || !strings.Contains(location, "error=") {
			t.Errorf("Expected redirect with error message, got location: %s", location)
		}
	})

	t.Run("cannot add laptop already in active shipment", func(t *testing.T) {
		// Create another bulk shipment
		var otherShipmentID int64
		err = db.QueryRowContext(ctx,
			`INSERT INTO shipments (shipment_type, client_company_id, status, laptop_count, jira_ticket_number, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
			models.ShipmentTypeBulkToWarehouse, companyID, models.ShipmentStatusPendingPickup, 3, "TEST-002", time.Now(), time.Now(),
		).Scan(&otherShipmentID)
		if err != nil {
			t.Fatalf("Failed to create other shipment: %v", err)
		}

		// Create laptop and link it to other shipment
		var laptopID int64
		err = db.QueryRowContext(ctx,
			`INSERT INTO laptops (serial_number, brand, model, cpu, ram_gb, ssd_gb, status, client_company_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`,
			"IN-SHIPMENT-001", "Dell", "Latitude", "i7", "16", "512", models.LaptopStatusInTransitToWarehouse, companyID, time.Now(), time.Now(),
		).Scan(&laptopID)
		if err != nil {
			t.Fatalf("Failed to create laptop: %v", err)
		}

		// Link laptop to other shipment
		_, err = db.ExecContext(ctx,
			`INSERT INTO shipment_laptops (shipment_id, laptop_id, created_at) VALUES ($1, $2, $3)`,
			otherShipmentID, laptopID, time.Now(),
		)
		if err != nil {
			t.Fatalf("Failed to link laptop to other shipment: %v", err)
		}

		formData := url.Values{}
		formData.Set("laptop_id", strconv.FormatInt(laptopID, 10))

		req := httptest.NewRequest(http.MethodPost, "/shipments/"+strconv.FormatInt(shipmentID, 10)+"/laptops/add", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		user := &models.User{ID: logisticsUserID, Email: "logistics@test.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(shipmentID, 10)})

		w := httptest.NewRecorder()
		handler.AddLaptopToBulkShipment(w, req)

		// Handler redirects on error with error message in query string
		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected status 303 (redirect), got %d. Body: %s", w.Code, w.Body.String())
		}
		location := w.Header().Get("Location")
		if location == "" || !strings.Contains(location, "error=") {
			t.Errorf("Expected redirect with error message, got location: %s", location)
		}
	})

	t.Run("cannot add laptop to non-bulk shipment", func(t *testing.T) {
		// Create single shipment
		var singleShipmentID int64
		err = db.QueryRowContext(ctx,
			`INSERT INTO shipments (shipment_type, client_company_id, status, laptop_count, jira_ticket_number, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
			models.ShipmentTypeSingleFullJourney, companyID, models.ShipmentStatusPendingPickup, 1, "TEST-003", time.Now(), time.Now(),
		).Scan(&singleShipmentID)
		if err != nil {
			t.Fatalf("Failed to create single shipment: %v", err)
		}

		// Create eligible laptop
		var laptopID int64
		err = db.QueryRowContext(ctx,
			`INSERT INTO laptops (serial_number, brand, model, cpu, ram_gb, ssd_gb, status, client_company_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`,
			"ELIGIBLE-003", "Dell", "Latitude", "i7", "16", "512", models.LaptopStatusInTransitToWarehouse, companyID, time.Now(), time.Now(),
		).Scan(&laptopID)
		if err != nil {
			t.Fatalf("Failed to create laptop: %v", err)
		}

		formData := url.Values{}
		formData.Set("laptop_id", strconv.FormatInt(laptopID, 10))

		req := httptest.NewRequest(http.MethodPost, "/shipments/"+strconv.FormatInt(singleShipmentID, 10)+"/laptops/add", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		user := &models.User{ID: logisticsUserID, Email: "logistics@test.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(singleShipmentID, 10)})

		w := httptest.NewRecorder()
		handler.AddLaptopToBulkShipment(w, req)

		// Handler redirects on error with error message in query string
		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected status 303 (redirect), got %d. Body: %s", w.Code, w.Body.String())
		}
		location := w.Header().Get("Location")
		if location == "" || !strings.Contains(location, "error=") {
			t.Errorf("Expected redirect with error message, got location: %s", location)
		}
	})

	t.Run("cannot add laptop with different client company to bulk shipment", func(t *testing.T) {
		// Create another company
		var otherCompanyID int64
		err = db.QueryRowContext(ctx,
			`INSERT INTO client_companies (name, contact_info, created_at)
			VALUES ($1, $2, $3) RETURNING id`,
			"Other Company", json.RawMessage(`{"email":"other@company.com"}`), time.Now(),
		).Scan(&otherCompanyID)
		if err != nil {
			t.Fatalf("Failed to create other company: %v", err)
		}

		// Create laptop with different company (but correct status)
		var laptopID int64
		err = db.QueryRowContext(ctx,
			`INSERT INTO laptops (serial_number, brand, model, cpu, ram_gb, ssd_gb, status, client_company_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`,
			"DIFF-COMPANY-001", "Dell", "Latitude", "i7", "16", "512", models.LaptopStatusInTransitToWarehouse, otherCompanyID, time.Now(), time.Now(),
		).Scan(&laptopID)
		if err != nil {
			t.Fatalf("Failed to create laptop: %v", err)
		}

		formData := url.Values{}
		formData.Set("laptop_id", strconv.FormatInt(laptopID, 10))

		req := httptest.NewRequest(http.MethodPost, "/shipments/"+strconv.FormatInt(shipmentID, 10)+"/laptops/add", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		user := &models.User{ID: logisticsUserID, Email: "logistics@test.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(shipmentID, 10)})

		w := httptest.NewRecorder()
		handler.AddLaptopToBulkShipment(w, req)

		// Handler redirects on error with error message in query string
		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected status 303 (redirect), got %d. Body: %s", w.Code, w.Body.String())
		}
		location := w.Header().Get("Location")
		if location == "" || !strings.Contains(location, "error=") {
			t.Errorf("Expected redirect with error message, got location: %s", location)
		}
		if !strings.Contains(location, "company") {
			t.Errorf("Expected error message about company mismatch, got location: %s", location)
		}
	})
}
