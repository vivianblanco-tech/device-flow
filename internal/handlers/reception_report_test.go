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

func TestReceptionReportPage(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create test warehouse user
	var warehouseUserID int64
	err := db.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"warehouse@example.com", "hashedpassword", models.RoleWarehouse, time.Now(), time.Now(),
	).Scan(&warehouseUserID)
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
		companyID, models.ShipmentStatusInTransitToWarehouse, "TEST-200", time.Now(), time.Now(),
	).Scan(&shipmentID)
	if err != nil {
		t.Fatalf("Failed to create test shipment: %v", err)
	}

	templates := loadTestTemplates(t)
	handler := NewReceptionReportHandler(db, templates, nil)

	t.Run("warehouse user can view reception report page", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/reception-report?shipment_id="+strconv.FormatInt(shipmentID, 10), nil)

		// Add warehouse user to context
		user := &models.User{ID: warehouseUserID, Email: "warehouse@example.com", Role: models.RoleWarehouse}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.ReceptionReportPage(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})

	t.Run("non-warehouse user cannot view page", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/reception-report?shipment_id="+strconv.FormatInt(shipmentID, 10), nil)

		// Add non-warehouse user to context
		user := &models.User{ID: warehouseUserID, Email: "client@example.com", Role: models.RoleClient}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.ReceptionReportPage(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("Expected status 403, got %d", w.Code)
		}
	})

	t.Run("unauthenticated user redirects to login", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/reception-report?shipment_id="+strconv.FormatInt(shipmentID, 10), nil)
		w := httptest.NewRecorder()

		handler.ReceptionReportPage(w, req)

		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected status 303, got %d", w.Code)
		}

		location := w.Header().Get("Location")
		if location != "/login" {
			t.Errorf("Expected redirect to /login, got %s", location)
		}
	})

	t.Run("missing shipment ID returns error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/reception-report", nil)

		user := &models.User{ID: warehouseUserID, Email: "warehouse@example.com", Role: models.RoleWarehouse}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.ReceptionReportPage(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})
}

func TestReceptionReportSubmit(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create test warehouse user
	var warehouseUserID int64
	err := db.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"warehouse@example.com", "hashedpassword", models.RoleWarehouse, time.Now(), time.Now(),
	).Scan(&warehouseUserID)
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
		companyID, models.ShipmentStatusInTransitToWarehouse, "TEST-201", time.Now(), time.Now(),
	).Scan(&shipmentID)
	if err != nil {
		t.Fatalf("Failed to create test shipment: %v", err)
	}

	templates := loadTestTemplates(t)
	handler := NewReceptionReportHandler(db, templates, nil)

	t.Run("valid submission creates reception report", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("shipment_id", strconv.FormatInt(shipmentID, 10))
		formData.Set("notes", "Received in good condition")

		req := httptest.NewRequest(http.MethodPost, "/reception-report", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		user := &models.User{ID: warehouseUserID, Email: "warehouse@example.com", Role: models.RoleWarehouse}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.ReceptionReportSubmit(w, req)

		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected status 303, got %d", w.Code)
		}

		// Verify reception report was created
		var count int
		err := db.QueryRowContext(ctx,
			`SELECT COUNT(*) FROM reception_reports WHERE shipment_id = $1`,
			shipmentID,
		).Scan(&count)
		if err != nil {
			t.Fatalf("Failed to query reception reports: %v", err)
		}
		if count != 1 {
			t.Errorf("Expected 1 reception report, got %d", count)
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
		if status != models.ShipmentStatusAtWarehouse {
			t.Errorf("Expected shipment status 'at_warehouse', got '%s'", status)
		}
	})

	t.Run("non-POST method returns error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/reception-report", nil)
		w := httptest.NewRecorder()

		handler.ReceptionReportSubmit(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status 405, got %d", w.Code)
		}
	})

	t.Run("non-warehouse user cannot submit", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("shipment_id", strconv.FormatInt(shipmentID, 10))
		formData.Set("notes", "Received in good condition")

		req := httptest.NewRequest(http.MethodPost, "/reception-report", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		user := &models.User{ID: warehouseUserID, Email: "client@example.com", Role: models.RoleClient}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.ReceptionReportSubmit(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("Expected status 403, got %d", w.Code)
		}
	})

	t.Run("unauthenticated user redirects to login", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("shipment_id", strconv.FormatInt(shipmentID, 10))
		formData.Set("notes", "Received in good condition")

		req := httptest.NewRequest(http.MethodPost, "/reception-report", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		w := httptest.NewRecorder()
		handler.ReceptionReportSubmit(w, req)

		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected status 303, got %d", w.Code)
		}
	})
}

func TestReceptionReportsList(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create test warehouse user
	var warehouseUserID int64
	err := db.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"warehouse@example.com", "hashedpassword", models.RoleWarehouse, time.Now(), time.Now(),
	).Scan(&warehouseUserID)
	if err != nil {
		t.Fatalf("Failed to create warehouse user: %v", err)
	}

	// Create test logistics user
	var logisticsUserID int64
	err = db.QueryRowContext(ctx,
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

	// Create test shipment
	var shipmentID int64
	err = db.QueryRowContext(ctx,
		`INSERT INTO shipments (client_company_id, status, shipment_type, laptop_count, jira_ticket_number, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		companyID, "at_warehouse", "single_full_journey", 1, "TEST-RR-001", time.Now(), time.Now(),
	).Scan(&shipmentID)
	if err != nil {
		t.Fatalf("Failed to create test shipment: %v", err)
	}

	// Create test reception report
	var receptionReportID int64
	receivedAt := time.Now().Add(-24 * time.Hour)
	err = db.QueryRowContext(ctx,
		`INSERT INTO reception_reports (shipment_id, warehouse_user_id, received_at, notes)
		VALUES ($1, $2, $3, $4) RETURNING id`,
		shipmentID, warehouseUserID, receivedAt, "Test notes",
	).Scan(&receptionReportID)
	if err != nil {
		t.Fatalf("Failed to create test reception report: %v", err)
	}

	handler := NewReceptionReportHandler(db, nil, nil)

	t.Run("successfully displays reception reports list for warehouse user", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/reception-reports", nil)
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserContextKey, &models.User{
			ID:    warehouseUserID,
			Email: "warehouse@example.com",
			Role:  models.RoleWarehouse,
		}))

		w := httptest.NewRecorder()
		handler.ReceptionReportsList(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})

	t.Run("successfully displays reception reports list for logistics user", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/reception-reports", nil)
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserContextKey, &models.User{
			ID:    logisticsUserID,
			Email: "logistics@example.com",
			Role:  models.RoleLogistics,
		}))

		w := httptest.NewRecorder()
		handler.ReceptionReportsList(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})

	t.Run("redirects to login when user not authenticated", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/reception-reports", nil)

		w := httptest.NewRecorder()
		handler.ReceptionReportsList(w, req)

		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected status 303, got %d", w.Code)
		}

		location := w.Header().Get("Location")
		if location != "/login" {
			t.Errorf("Expected redirect to /login, got %s", location)
		}
	})

	t.Run("displays reception reports data in response", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/reception-reports", nil)
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserContextKey, &models.User{
			ID:    warehouseUserID,
			Email: "warehouse@example.com",
			Role:  models.RoleWarehouse,
		}))

		w := httptest.NewRecorder()
		handler.ReceptionReportsList(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		body := w.Body.String()

		// Verify the response contains reception report data
		if !strings.Contains(body, "Test Company") {
			t.Errorf("Expected response to contain company name 'Test Company', got: %s", body)
		}

		if !strings.Contains(body, "warehouse@example.com") {
			t.Errorf("Expected response to contain warehouse user email, got: %s", body)
		}

		if !strings.Contains(body, "Test notes") {
			t.Errorf("Expected response to contain notes 'Test notes', got: %s", body)
		}
	})
}

func TestReceptionReportDetail(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create test warehouse user
	var warehouseUserID int64
	err := db.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"warehouse@example.com", "hashedpassword", models.RoleWarehouse, time.Now(), time.Now(),
	).Scan(&warehouseUserID)
	if err != nil {
		t.Fatalf("Failed to create warehouse user: %v", err)
	}

	// Create test logistics user
	var logisticsUserID int64
	err = db.QueryRowContext(ctx,
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

	// Create test shipment
	var shipmentID int64
	err = db.QueryRowContext(ctx,
		`INSERT INTO shipments (client_company_id, status, shipment_type, laptop_count, jira_ticket_number, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		companyID, "at_warehouse", "single_full_journey", 1, "TEST-DETAIL-001", time.Now(), time.Now(),
	).Scan(&shipmentID)
	if err != nil {
		t.Fatalf("Failed to create test shipment: %v", err)
	}

	// Create test reception report
	var receptionReportID int64
	receivedAt := time.Now().Add(-24 * time.Hour)
	err = db.QueryRowContext(ctx,
		`INSERT INTO reception_reports (shipment_id, warehouse_user_id, received_at, notes)
		VALUES ($1, $2, $3, $4) RETURNING id`,
		shipmentID, warehouseUserID, receivedAt, "Detailed reception notes",
	).Scan(&receptionReportID)
	if err != nil {
		t.Fatalf("Failed to create test reception report: %v", err)
	}

	handler := NewReceptionReportHandler(db, nil, nil)

	t.Run("warehouse user can view reception report detail", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/reception-reports/"+strconv.FormatInt(receptionReportID, 10), nil)
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserContextKey, &models.User{
			ID:    warehouseUserID,
			Email: "warehouse@example.com",
			Role:  models.RoleWarehouse,
		}))

		w := httptest.NewRecorder()
		handler.ReceptionReportDetail(w, req, receptionReportID)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		body := w.Body.String()
		if !strings.Contains(body, "Detailed reception notes") {
			t.Errorf("Expected response to contain notes, got: %s", body)
		}

		if !strings.Contains(body, "Test Company") {
			t.Errorf("Expected response to contain company name, got: %s", body)
		}
	})

	t.Run("logistics user can view reception report detail", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/reception-reports/"+strconv.FormatInt(receptionReportID, 10), nil)
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserContextKey, &models.User{
			ID:    logisticsUserID,
			Email: "logistics@example.com",
			Role:  models.RoleLogistics,
		}))

		w := httptest.NewRecorder()
		handler.ReceptionReportDetail(w, req, receptionReportID)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})

	t.Run("redirects to login when user not authenticated", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/reception-reports/"+strconv.FormatInt(receptionReportID, 10), nil)

		w := httptest.NewRecorder()
		handler.ReceptionReportDetail(w, req, receptionReportID)

		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected status 303, got %d", w.Code)
		}

		location := w.Header().Get("Location")
		if location != "/login" {
			t.Errorf("Expected redirect to /login, got %s", location)
		}
	})

	t.Run("returns 404 for non-existent reception report", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/reception-reports/99999", nil)
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserContextKey, &models.User{
			ID:    warehouseUserID,
			Email: "warehouse@example.com",
			Role:  models.RoleWarehouse,
		}))

		w := httptest.NewRecorder()
		handler.ReceptionReportDetail(w, req, 99999)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", w.Code)
		}
	})
}
