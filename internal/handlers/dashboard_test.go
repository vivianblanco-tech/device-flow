package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/yourusername/laptop-tracking-system/internal/database"
	"github.com/yourusername/laptop-tracking-system/internal/middleware"
	"github.com/yourusername/laptop-tracking-system/internal/models"
)

// TestDashboardAccessControl tests that only logistics users can access the dashboard
func TestDashboardAccessControl(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup test database
	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	// Load templates
	templates := loadTestTemplates(t)

	// Create dashboard handler
	handler := NewDashboardHandler(db, templates)

	tests := []struct {
		name           string
		userRole       models.UserRole
		expectedStatus int
		expectRedirect bool
	}{
		{
			name:           "Logistics user can access dashboard",
			userRole:       models.RoleLogistics,
			expectedStatus: http.StatusOK,
			expectRedirect: false,
		},
		{
			name:           "Client user cannot access dashboard",
			userRole:       models.RoleClient,
			expectedStatus: http.StatusForbidden,
			expectRedirect: false,
		},
		{
			name:           "Warehouse user cannot access dashboard",
			userRole:       models.RoleWarehouse,
			expectedStatus: http.StatusForbidden,
			expectRedirect: false,
		},
		{
			name:           "Project Manager user can access dashboard",
			userRole:       models.RoleProjectManager,
			expectedStatus: http.StatusOK,
			expectRedirect: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test user with specific role
			user := &models.User{
				ID:    1,
				Email: "test@bairesdev.com",
				Role:  tt.userRole,
			}

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)

			// Add user to context
			ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
			req = req.WithContext(ctx)

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call handler
			handler.Dashboard(rr, req)

			// Check status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

// TestDashboardUnauthenticated tests that unauthenticated users are redirected to login
func TestDashboardUnauthenticated(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup test database
	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	// Load templates
	templates := loadTestTemplates(t)

	// Create dashboard handler
	handler := NewDashboardHandler(db, templates)

	// Create request without user in context
	req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	rr := httptest.NewRecorder()

	// Call handler
	handler.Dashboard(rr, req)

	// Check status code (should redirect)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("expected status %d (redirect), got %d", http.StatusSeeOther, rr.Code)
	}

	// Check redirect location
	location := rr.Header().Get("Location")
	if location != "/login" {
		t.Errorf("expected redirect to /login, got %s", location)
	}
}

// TestDashboardMenuVisibility tests that dashboard menu item is only visible to authorized roles
func TestDashboardMenuVisibility(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup test database
	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	// Load templates
	templates := loadTestTemplates(t)

	// Create handlers
	dashboardHandler := NewDashboardHandler(db, templates)
	shipmentsHandler := NewShipmentsHandler(db, templates)

	tests := []struct {
		name                string
		userRole            models.UserRole
		shouldShowDashboard bool
		testHandler         string // "dashboard" or "shipments"
	}{
		{
			name:                "Logistics user sees dashboard menu on dashboard page",
			userRole:            models.RoleLogistics,
			shouldShowDashboard: true,
			testHandler:         "dashboard",
		},
		{
			name:                "Client user does not see dashboard menu on shipments page",
			userRole:            models.RoleClient,
			shouldShowDashboard: false,
			testHandler:         "shipments",
		},
		{
			name:                "Warehouse user does not see dashboard menu on shipments page",
			userRole:            models.RoleWarehouse,
			shouldShowDashboard: false,
			testHandler:         "shipments",
		},
		{
			name:                "Project Manager sees dashboard menu on dashboard page",
			userRole:            models.RoleProjectManager,
			shouldShowDashboard: true,
			testHandler:         "dashboard",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test user with specific role
			user := &models.User{
				ID:    1,
				Email: "test@bairesdev.com",
				Role:  tt.userRole,
			}

			// Create request
			var req *http.Request
			if tt.testHandler == "dashboard" {
				req = httptest.NewRequest(http.MethodGet, "/dashboard", nil)
			} else {
				req = httptest.NewRequest(http.MethodGet, "/shipments", nil)
			}

			// Add user to context
			ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
			req = req.WithContext(ctx)

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call appropriate handler
			if tt.testHandler == "dashboard" {
				dashboardHandler.Dashboard(rr, req)
			} else {
				shipmentsHandler.ShipmentsList(rr, req)
			}

			// Check status code is OK
			if rr.Code != http.StatusOK {
				t.Fatalf("expected status 200, got %d", rr.Code)
			}

			// Check dashboard link presence/absence
			body := rr.Body.String()
			dashboardLinkPresent := strings.Contains(body, `href="/dashboard"`)

			if tt.shouldShowDashboard && !dashboardLinkPresent {
				t.Errorf("expected dashboard link in HTML for %s role, but it was not found", tt.userRole)
			}

			if !tt.shouldShowDashboard && dashboardLinkPresent {
				t.Errorf("did not expect dashboard link in HTML for %s role, but it was found", tt.userRole)
			}
		})
	}
}
