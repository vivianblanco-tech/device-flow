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

func TestPickupFormPage(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	// Create test client company
	ctx := context.Background()
	var companyID int64
	err := db.QueryRowContext(ctx,
		`INSERT INTO client_companies (name, contact_info, created_at)
		VALUES ($1, $2, $3) RETURNING id`,
		"Test Company", json.RawMessage(`{"email":"test@company.com"}`), time.Now(),
	).Scan(&companyID)
	if err != nil {
		t.Fatalf("Failed to create test company: %v", err)
	}

	// Create test user
	var userID int64
	err = db.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"client@company.com", "$2a$12$test.hash.for.testing.purposes", models.RoleClient, time.Now(), time.Now(),
	).Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	templates := loadTestTemplates(t)
	handler := NewPickupFormHandler(db, templates, nil)

	t.Run("GET request displays form", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/pickup-form?company_id="+strconv.FormatInt(companyID, 10), nil)

		// Add user to context
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, &models.User{ID: userID, Email: "client@company.com", Role: models.RoleClient})
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		handler.PickupFormPage(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})
}

func TestPickupFormSubmit(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create test client company
	var companyID int64
	err := db.QueryRowContext(ctx,
		`INSERT INTO client_companies (name, contact_info, created_at)
		VALUES ($1, $2, $3) RETURNING id`,
		"Test Company", json.RawMessage(`{"email":"test@company.com"}`), time.Now(),
	).Scan(&companyID)
	if err != nil {
		t.Fatalf("Failed to create test company: %v", err)
	}

	// Create test user
	var userID int64
	err = db.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"client@company.com", "$2a$12$test.hash.for.testing.purposes", models.RoleClient, time.Now(), time.Now(),
	).Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	templates := loadTestTemplates(t)
	handler := NewPickupFormHandler(db, templates, nil)

	t.Run("valid form submission creates shipment", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("client_company_id", strconv.FormatInt(companyID, 10))
		formData.Set("contact_name", "John Doe")
		formData.Set("contact_email", "john@company.com")
		formData.Set("contact_phone", "+1-555-0123")
		formData.Set("pickup_address", "123 Main St")
		formData.Set("pickup_city", "New York")
		formData.Set("pickup_state", "NY")
		formData.Set("pickup_zip", "10001")
		formData.Set("pickup_date", time.Now().Add(24*time.Hour).Format("2006-01-02"))
		formData.Set("pickup_time_slot", "morning")
		formData.Set("number_of_laptops", "3")
		formData.Set("jira_ticket_number", "TEST-500")
		formData.Set("special_instructions", "Please call before arrival")
		formData.Set("number_of_boxes", "2")
		formData.Set("assignment_type", "single")
		formData.Set("include_accessories", "false")

		req := httptest.NewRequest(http.MethodPost, "/pickup-form", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		// Add user to context (simulating authenticated request)
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, &models.User{ID: userID, Email: "client@company.com", Role: models.RoleClient})
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		handler.PickupFormSubmit(w, req)

		// Check redirect
		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected status 303, got %d", w.Code)
		}

		// Verify shipment was created
		var count int
		err := db.QueryRowContext(context.Background(),
			`SELECT COUNT(*) FROM shipments WHERE client_company_id = $1`,
			companyID,
		).Scan(&count)
		if err != nil {
			t.Fatalf("Failed to query shipments: %v", err)
		}
		if count != 1 {
			t.Errorf("Expected 1 shipment, got %d", count)
		}
	})

	t.Run("invalid form submission returns error", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("client_company_id", strconv.FormatInt(companyID, 10))
		// Missing required fields

		req := httptest.NewRequest(http.MethodPost, "/pickup-form", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		ctx := context.WithValue(req.Context(), middleware.UserContextKey, &models.User{ID: userID, Email: "client@company.com", Role: models.RoleClient})
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		handler.PickupFormSubmit(w, req)

		// Should get redirect with error
		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected status 303, got %d", w.Code)
		}

		location := w.Header().Get("Location")
		if !strings.Contains(location, "error=") {
			t.Errorf("Expected error in redirect URL, got: %s", location)
		}
	})
}

// 🟥 RED: Test single full journey form submission
func TestPickupFormHandler_SubmitSingleFullJourney(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create test client company
	var companyID int64
	err := db.QueryRowContext(ctx,
		`INSERT INTO client_companies (name, contact_info, created_at)
		VALUES ($1, $2, $3) RETURNING id`,
		"Test Company", json.RawMessage(`{"email":"test@company.com"}`), time.Now(),
	).Scan(&companyID)
	if err != nil {
		t.Fatalf("Failed to create test company: %v", err)
	}

	// Create test user
	var userID int64
	err = db.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"client@company.com", "$2a$12$test.hash.for.testing.purposes", models.RoleClient, time.Now(), time.Now(),
	).Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	handler := NewPickupFormHandler(db, nil, nil)

	t.Run("single full journey form creates shipment with correct type", func(t *testing.T) {
		formData := url.Values{
			"shipment_type":        {string(models.ShipmentTypeSingleFullJourney)},
			"client_company_id":    {strconv.FormatInt(companyID, 10)},
			"contact_name":         {"John Doe"},
			"contact_email":        {"john@test.com"},
			"contact_phone":        {"+1-555-0123"},
			"pickup_address":       {"123 Main St"},
			"pickup_city":          {"New York"},
			"pickup_state":         {"NY"},
			"pickup_zip":           {"10001"},
			"pickup_date":          {time.Now().Add(24 * time.Hour).Format("2006-01-02")},
			"pickup_time_slot":     {"morning"},
			"jira_ticket_number":   {"SCOP-12345"},
			"laptop_serial_number": {"ABC123456"},
			"laptop_specs":         {"Dell XPS 15, 16GB RAM"},
			"engineer_name":        {"Jane Smith"},
			"include_accessories":  {"false"},
		}

		req := httptest.NewRequest(http.MethodPost, "/pickup-form", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserContextKey, &models.User{ID: userID, Role: models.RoleClient}))

		w := httptest.NewRecorder()

		handler.PickupFormSubmit(w, req)

		// Check response
		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected status 303, got %d", w.Code)
		}

		// Verify shipment was created with correct type
		var shipmentID int64
		var shipmentType models.ShipmentType
		var laptopCount int
		err := db.QueryRowContext(ctx,
			`SELECT id, shipment_type, laptop_count 
			FROM shipments 
			WHERE client_company_id = $1 AND jira_ticket_number = $2`,
			companyID, "SCOP-12345",
		).Scan(&shipmentID, &shipmentType, &laptopCount)

		if err != nil {
			t.Fatalf("Shipment not created: %v", err)
		}

		if shipmentType != models.ShipmentTypeSingleFullJourney {
			t.Errorf("Expected shipment type %s, got %s", models.ShipmentTypeSingleFullJourney, shipmentType)
		}

		if laptopCount != 1 {
			t.Errorf("Expected laptop count 1, got %d", laptopCount)
		}

		// Verify laptop was auto-created
		var laptopID int64
		var laptopSerialNumber string
		var laptopStatus models.LaptopStatus
		err = db.QueryRowContext(ctx,
			`SELECT l.id, l.serial_number, l.status
			FROM laptops l
			JOIN shipment_laptops sl ON sl.laptop_id = l.id
			WHERE sl.shipment_id = $1`,
			shipmentID,
		).Scan(&laptopID, &laptopSerialNumber, &laptopStatus)

		if err != nil {
			t.Fatalf("Laptop not created: %v", err)
		}

		if laptopSerialNumber != "ABC123456" {
			t.Errorf("Expected serial number ABC123456, got %s", laptopSerialNumber)
		}

		if laptopStatus != models.LaptopStatusInTransitToWarehouse {
			t.Errorf("Expected laptop status %s, got %s", models.LaptopStatusInTransitToWarehouse, laptopStatus)
		}
	})

	t.Run("single full journey without engineer name succeeds", func(t *testing.T) {
		formData := url.Values{
			"shipment_type":        {string(models.ShipmentTypeSingleFullJourney)},
			"client_company_id":    {strconv.FormatInt(companyID, 10)},
			"contact_name":         {"John Doe"},
			"contact_email":        {"john@test.com"},
			"contact_phone":        {"+1-555-0123"},
			"pickup_address":       {"123 Main St"},
			"pickup_city":          {"New York"},
			"pickup_state":         {"NY"},
			"pickup_zip":           {"10001"},
			"pickup_date":          {time.Now().Add(24 * time.Hour).Format("2006-01-02")},
			"pickup_time_slot":     {"morning"},
			"jira_ticket_number":   {"SCOP-12346"},
			"laptop_serial_number": {"DEF789012"},
			"laptop_specs":         {"Lenovo ThinkPad"},
			// engineer_name is optional
			"include_accessories": {"false"},
		}

		req := httptest.NewRequest(http.MethodPost, "/pickup-form", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserContextKey, &models.User{ID: userID, Role: models.RoleClient}))

		w := httptest.NewRecorder()

		handler.PickupFormSubmit(w, req)

		// Check response
		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected status 303, got %d", w.Code)
		}

		// Verify shipment was created
		var shipmentID int64
		var softwareEngineerID *int64
		err := db.QueryRowContext(ctx,
			`SELECT id, software_engineer_id
			FROM shipments 
			WHERE client_company_id = $1 AND jira_ticket_number = $2`,
			companyID, "SCOP-12346",
		).Scan(&shipmentID, &softwareEngineerID)

		if err != nil {
			t.Fatalf("Shipment not created: %v", err)
		}

		// Engineer should be nil (not assigned yet)
		if softwareEngineerID != nil {
			t.Errorf("Expected software_engineer_id to be nil, got %v", *softwareEngineerID)
		}
	})

	t.Run("single full journey without serial number fails validation", func(t *testing.T) {
		formData := url.Values{
			"shipment_type":      {string(models.ShipmentTypeSingleFullJourney)},
			"client_company_id":  {strconv.FormatInt(companyID, 10)},
			"contact_name":       {"John Doe"},
			"contact_email":      {"john@test.com"},
			"contact_phone":      {"+1-555-0123"},
			"pickup_address":     {"123 Main St"},
			"pickup_city":        {"New York"},
			"pickup_state":       {"NY"},
			"pickup_zip":         {"10001"},
			"pickup_date":        {time.Now().Add(24 * time.Hour).Format("2006-01-02")},
			"pickup_time_slot":   {"morning"},
			"jira_ticket_number": {"SCOP-12347"},
			// laptop_serial_number is MISSING (required)
			"laptop_specs":        {"Dell XPS 15"},
			"include_accessories": {"false"},
		}

		req := httptest.NewRequest(http.MethodPost, "/pickup-form", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserContextKey, &models.User{ID: userID, Role: models.RoleClient}))

		w := httptest.NewRecorder()

		handler.PickupFormSubmit(w, req)

		// Should redirect with error
		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected status 303, got %d", w.Code)
		}

		location := w.Header().Get("Location")
		if !strings.Contains(location, "error=") {
			t.Errorf("Expected error in redirect URL, got: %s", location)
		}

		// Verify shipment was NOT created
		var count int
		err := db.QueryRowContext(ctx,
			`SELECT COUNT(*) FROM shipments WHERE jira_ticket_number = $1`,
			"SCOP-12347",
		).Scan(&count)
		if err != nil {
			t.Fatalf("Failed to query shipments: %v", err)
		}
		if count != 0 {
			t.Errorf("Expected 0 shipments, got %d", count)
		}
	})
}
