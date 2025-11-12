package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
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
	shipmentsHandler := NewShipmentsHandler(db, templates, nil) // nil email notifier for tests

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

// Phase 5 Test: Test that dashboard displays three shipment type creation buttons for logistics users
func TestDashboardThreeShipmentTypeButtons(t *testing.T) {
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

	// Create test user with logistics role
	user := &models.User{
		ID:    1,
		Email: "logistics@bairesdev.com",
		Role:  models.RoleLogistics,
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
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	// Check that all three shipment type buttons are present
	body := rr.Body.String()

	// Debug: Save body to file for inspection
	err := os.WriteFile("dashboard_test_output.html", []byte(body), 0644)
	if err != nil {
		t.Logf("Failed to write body to file: %v", err)
	}

	// Debug: Output first 1000 chars if test fails
	t.Logf("Response body length: %d", len(body))
	if len(body) > 0 && len(body) < 2000 {
		t.Logf("Full body: %s", body)
	}

	// Check if Quick Actions section exists
	if !strings.Contains(body, "Quick Actions") {
		t.Error("Quick Actions section not found in dashboard")
	}

	if !strings.Contains(body, `/shipments/create/single`) {
		t.Error("Expected dashboard to contain link to single shipment form")
		t.Logf("Body snippet: %s", body[:min(len(body), 1000)])
	}
	if !strings.Contains(body, `Single Shipment`) && !strings.Contains(body, `Single Full Journey`) {
		t.Error("Expected dashboard to contain 'Single Shipment' or 'Single Full Journey' button text")
	}

	if !strings.Contains(body, `/shipments/create/bulk`) {
		t.Error("Expected dashboard to contain link to bulk shipment form")
	}
	if !strings.Contains(body, `Bulk`) {
		t.Error("Expected dashboard to contain 'Bulk' button text")
	}

	if !strings.Contains(body, `/shipments/create/warehouse-to-engineer`) {
		t.Error("Expected dashboard to contain link to warehouse-to-engineer form")
	}
	if !strings.Contains(body, `Warehouse`) && !strings.Contains(body, `Engineer`) {
		t.Error("Expected dashboard to contain 'Warehouse' or 'Engineer' button text")
	}
}

// Test that old Quick Actions buttons (New Pickup Request and View All Shipments) are NOT present
func TestDashboardQuickActionsRemovedButtons(t *testing.T) {
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

	// Test for logistics user
	logisticsUser := &models.User{
		ID:    1,
		Email: "logistics@bairesdev.com",
		Role:  models.RoleLogistics,
	}

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)

	// Add user to context
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, logisticsUser)
	req = req.WithContext(ctx)

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	handler.Dashboard(rr, req)

	// Check status code
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	body := rr.Body.String()

	// Verify "New Pickup Request" button is NOT present
	if strings.Contains(body, "New Pickup Request") || strings.Contains(body, "+ New Pickup Request") {
		t.Error("Expected 'New Pickup Request' button to be removed from Quick Actions")
	}

	// Verify direct /pickup-form link in Quick Actions is NOT present
	// (but it may exist elsewhere in navbar, so be specific about Quick Actions context)
	if strings.Contains(body, `href="/pickup-form"`) && strings.Contains(body, "Quick Actions") {
		// Check if pickup-form link appears in the Quick Actions section
		quickActionsStart := strings.Index(body, "Quick Actions")
		if quickActionsStart != -1 {
			// Find the end of Quick Actions section (next major section or end of div)
			quickActionsSection := body[quickActionsStart:]
			endOfSection := strings.Index(quickActionsSection, "</div>\n        {{end}}\n    </div>")
			if endOfSection == -1 {
				endOfSection = len(quickActionsSection)
			}
			quickActionsContent := quickActionsSection[:endOfSection]

			if strings.Contains(quickActionsContent, `href="/pickup-form"`) {
				t.Error("Expected /pickup-form link to be removed from Quick Actions section")
			}
		}
	}

	// Verify "View All Shipments" button is NOT present in Quick Actions
	// (It's OK if "Shipments" appears in the three shipment type buttons)
	if strings.Contains(body, "View All Shipments") {
		t.Error("Expected 'View All Shipments' button to be removed from Quick Actions")
	}

	// Verify direct /shipments link in Quick Actions (without create path) is NOT present
	quickActionsStart := strings.Index(body, "Quick Actions")
	if quickActionsStart != -1 {
		quickActionsSection := body[quickActionsStart:]
		endOfSection := strings.Index(quickActionsSection, "</div>\n        {{end}}\n    </div>")
		if endOfSection == -1 {
			endOfSection = len(quickActionsSection)
		}
		quickActionsContent := quickActionsSection[:endOfSection]

		// Check for standalone /shipments link (not /shipments/create/...)
		if strings.Contains(quickActionsContent, `href="/shipments"`) && !strings.Contains(quickActionsContent, `href="/shipments/create/`) {
			t.Error("Expected standalone /shipments link to be removed from Quick Actions section")
		}
	}
}
