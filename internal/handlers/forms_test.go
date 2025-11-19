package handlers

import (
	"context"
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

func TestFormsHandler_FormsPage_RequiresLogisticsRole(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	templates := loadTestTemplates(t)
	handler := NewFormsHandler(db, templates)

	tests := []struct {
		name           string
		userRole       models.UserRole
		expectedStatus int
	}{
		{
			name:           "logistics user can access forms page",
			userRole:       models.RoleLogistics,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "client user cannot access forms page",
			userRole:       models.RoleClient,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "warehouse user cannot access forms page",
			userRole:       models.RoleWarehouse,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "project manager cannot access forms page",
			userRole:       models.RoleProjectManager,
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test user
			user := &models.User{
				ID:    1,
				Email: "test@example.com",
				Role:  tt.userRole,
			}

			// Create request
			req := httptest.NewRequest("GET", "/forms", nil)
			ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
			req = req.WithContext(ctx)

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call handler
			handler.FormsPage(rr, req)

			// Check status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

// TestClientCompanyEditPage_FormatsJSONContactInfo tests that JSON contact info is converted to plain text in the form
func TestClientCompanyEditPage_FormatsJSONContactInfo(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create test user
	user := &models.User{
		ID:    1,
		Email: "logistics@bairesdev.com",
		Role:  models.RoleLogistics,
	}
	var userID int64
	err := db.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		user.Email, "$2a$12$test.hash", user.Role, time.Now(), time.Now(),
	).Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	user.ID = userID

	// Create test company with JSON contact info
	var companyID int64
	jsonContactInfo := `{"email":"contact@example.com","phone":"+1-555-0100","address":"123 Main St"}`
	err = db.QueryRowContext(ctx,
		`INSERT INTO client_companies (name, contact_info, created_at, updated_at)
		VALUES ($1, $2, $3, $4) RETURNING id`,
		"Test Company", jsonContactInfo, time.Now(), time.Now(),
	).Scan(&companyID)
	if err != nil {
		t.Fatalf("Failed to create test company: %v", err)
	}

	templates := loadTestTemplates(t)
	handler := NewFormsHandler(db, templates)

	req := httptest.NewRequest("GET", "/forms/client-companies/"+strconv.FormatInt(companyID, 10)+"/edit", nil)
	req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(companyID, 10)})
	reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
	req = req.WithContext(reqCtx)

	rr := httptest.NewRecorder()
	handler.ClientCompanyEditPage(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	// Check that the response contains formatted contact info (not raw JSON)
	body := rr.Body.String()
	if strings.Contains(body, jsonContactInfo) {
		t.Error("Form should not display raw JSON contact info")
	}
	// The formatted version should have "Email:" and "Phone:" labels
	if !strings.Contains(body, "Email:") || !strings.Contains(body, "Phone:") {
		t.Error("Form should display formatted contact info with labels")
	}
}

// TestClientCompanyAddSubmit_AcceptsPlainTextContactInfo tests that plain text contact info is stored correctly
func TestClientCompanyAddSubmit_AcceptsPlainTextContactInfo(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	user := &models.User{
		ID:    1,
		Email: "logistics@bairesdev.com",
		Role:  models.RoleLogistics,
	}
	var userID int64
	err := db.QueryRowContext(context.Background(),
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		user.Email, "$2a$12$test.hash", user.Role, time.Now(), time.Now(),
	).Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	user.ID = userID

	templates := loadTestTemplates(t)
	handler := NewFormsHandler(db, templates)

	// Prepare form data with plain text contact info
	formData := url.Values{}
	formData.Set("name", "New Company")
	formData.Set("contact_info", "Email: contact@example.com\nPhone: +1-555-0100")

	req := httptest.NewRequest("POST", "/forms/client-companies/add", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
	req = req.WithContext(reqCtx)

	rr := httptest.NewRecorder()
	handler.ClientCompanyAddSubmit(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("expected status %d, got %d", http.StatusSeeOther, rr.Code)
	}

	// Verify the company was created with plain text contact info
	var contactInfo string
	err = db.QueryRowContext(context.Background(),
		`SELECT contact_info FROM client_companies WHERE name = $1`,
		"New Company",
	).Scan(&contactInfo)
	if err != nil {
		t.Fatalf("Failed to retrieve created company: %v", err)
	}

	expected := "Email: contact@example.com\nPhone: +1-555-0100"
	if contactInfo != expected {
		t.Errorf("expected contact_info %q, got %q", expected, contactInfo)
	}
}

// TestClientCompanyEditSubmit_AcceptsPlainTextContactInfo tests that plain text contact info updates are stored correctly
func TestClientCompanyEditSubmit_AcceptsPlainTextContactInfo(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	user := &models.User{
		ID:    1,
		Email: "logistics@bairesdev.com",
		Role:  models.RoleLogistics,
	}
	var userID int64
	err := db.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		user.Email, "$2a$12$test.hash", user.Role, time.Now(), time.Now(),
	).Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	user.ID = userID

	// Create test company with JSON contact info (simulating old data)
	var companyID int64
	jsonContactInfo := `{"email":"old@example.com","phone":"+1-555-0000"}`
	err = db.QueryRowContext(ctx,
		`INSERT INTO client_companies (name, contact_info, created_at, updated_at)
		VALUES ($1, $2, $3, $4) RETURNING id`,
		"Test Company", jsonContactInfo, time.Now(), time.Now(),
	).Scan(&companyID)
	if err != nil {
		t.Fatalf("Failed to create test company: %v", err)
	}

	templates := loadTestTemplates(t)
	handler := NewFormsHandler(db, templates)

	// Prepare form data with plain text contact info
	formData := url.Values{}
	formData.Set("name", "Test Company")
	formData.Set("contact_info", "Email: new@example.com\nPhone: +1-555-9999\nAddress: 456 New St")

	req := httptest.NewRequest("POST", "/forms/client-companies/"+strconv.FormatInt(companyID, 10)+"/edit", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(companyID, 10)})
	reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
	req = req.WithContext(reqCtx)

	rr := httptest.NewRecorder()
	handler.ClientCompanyEditSubmit(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("expected status %d, got %d", http.StatusSeeOther, rr.Code)
	}

	// Verify the company was updated with plain text contact info
	var contactInfo string
	err = db.QueryRowContext(ctx,
		`SELECT contact_info FROM client_companies WHERE id = $1`,
		companyID,
	).Scan(&contactInfo)
	if err != nil {
		t.Fatalf("Failed to retrieve updated company: %v", err)
	}

	expected := "Email: new@example.com\nPhone: +1-555-9999\nAddress: 456 New St"
	if contactInfo != expected {
		t.Errorf("expected contact_info %q, got %q", expected, contactInfo)
	}
}
