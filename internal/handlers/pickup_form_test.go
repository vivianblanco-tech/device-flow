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

	handler := NewPickupFormHandler(db, nil)

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

	handler := NewPickupFormHandler(db, nil)

	t.Run("valid form submission creates shipment", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("client_company_id", strconv.FormatInt(companyID, 10))
		formData.Set("contact_name", "John Doe")
		formData.Set("contact_email", "john@company.com")
		formData.Set("contact_phone", "+1-555-0123")
		formData.Set("pickup_address", "123 Main St, City, State 12345")
		formData.Set("pickup_date", time.Now().Add(24*time.Hour).Format("2006-01-02"))
		formData.Set("pickup_time_slot", "morning")
		formData.Set("number_of_laptops", "3")
		formData.Set("special_instructions", "Please call before arrival")

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
