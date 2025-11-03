package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
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
			name:           "Project Manager user cannot access dashboard",
			userRole:       models.RoleProjectManager,
			expectedStatus: http.StatusForbidden,
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

