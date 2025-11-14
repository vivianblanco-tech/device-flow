package integration

import (
	"context"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/yourusername/laptop-tracking-system/internal/database"
	"github.com/yourusername/laptop-tracking-system/internal/email"
	"github.com/yourusername/laptop-tracking-system/internal/handlers"
	"github.com/yourusername/laptop-tracking-system/internal/middleware"
	"github.com/yourusername/laptop-tracking-system/internal/models"
	"github.com/yourusername/laptop-tracking-system/internal/views"
)

// TestNavbarConsistencyAcrossPages verifies that all pages use the consistent navbar component
func TestNavbarConsistencyAcrossPages(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup test database
	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	// Load templates with navigation helper function
	templates := loadTemplatesWithNavigation(t)

	// Setup email client and notifier for ShipmentsHandler
	emailClient, err := email.NewClient(email.Config{
		Host: "localhost",
		Port: 1025,
		From: "test@example.com",
	})
	if err != nil {
		t.Fatalf("Failed to create email client: %v", err)
	}
	emailNotifier := email.NewNotifier(emailClient, db)

	tests := []struct {
		name              string
		handler           http.HandlerFunc
		path              string
		userRole          models.UserRole
		expectedStatus    int
		shouldHaveSticky  bool
		expectedLinks     []string
		unexpectedLinks   []string
	}{
		{
			name:             "Dashboard page for logistics user",
			handler:          handlers.NewDashboardHandler(db, templates).Dashboard,
			path:             "/dashboard",
			userRole:         models.RoleLogistics,
			expectedStatus:   http.StatusOK,
			shouldHaveSticky: true,
			expectedLinks: []string{
				`href="/dashboard"`,
				`href="/shipments"`,
				`href="/inventory"`,
				`href="/calendar"`,
				`href="/pickup-forms"`,
				`href="/reception-reports"`,
			},
			unexpectedLinks: []string{},
		},
		{
			name:             "Dashboard page for project manager",
			handler:          handlers.NewDashboardHandler(db, templates).Dashboard,
			path:             "/dashboard",
			userRole:         models.RoleProjectManager,
			expectedStatus:   http.StatusOK,
			shouldHaveSticky: true,
			expectedLinks: []string{
				`href="/dashboard"`,
				`href="/shipments"`,
				`href="/inventory"`,
				`href="/calendar"`,
			},
			unexpectedLinks: []string{
				`href="/pickup-forms"`,
				`href="/reception-reports"`,
			},
		},
		{
			name:             "Shipments page for warehouse user",
			handler:          handlers.NewShipmentsHandler(db, templates, emailNotifier).ShipmentsList,
			path:             "/shipments",
			userRole:         models.RoleWarehouse,
			expectedStatus:   http.StatusOK,
			shouldHaveSticky: true,
			expectedLinks: []string{
				`href="/shipments"`,
				`href="/inventory"`,
				`href="/calendar"`,
				`href="/reception-reports"`,
			},
			unexpectedLinks: []string{
				`href="/dashboard"`,
				`href="/pickup-forms"`,
			},
		},
		{
			name:             "Shipments page for client user",
			handler:          handlers.NewShipmentsHandler(db, templates, emailNotifier).ShipmentsList,
			path:             "/shipments",
			userRole:         models.RoleClient,
			expectedStatus:   http.StatusOK,
			shouldHaveSticky: true,
			expectedLinks: []string{
				`href="/shipments"`,
				`href="/calendar"`,
				`href="/pickup-forms"`,
			},
			unexpectedLinks: []string{
				`href="/dashboard"`,
				`href="/inventory"`,
				`href="/reception-reports"`,
			},
		},
		{
			name:             "Inventory page for logistics user",
			handler:          handlers.NewInventoryHandler(db, templates).InventoryList,
			path:             "/inventory",
			userRole:         models.RoleLogistics,
			expectedStatus:   http.StatusOK,
			shouldHaveSticky: true,
			expectedLinks: []string{
				`href="/dashboard"`,
				`href="/shipments"`,
				`href="/inventory"`,
				`href="/calendar"`,
			},
			unexpectedLinks: []string{},
		},
		{
			name:             "Calendar page for warehouse user",
			handler:          handlers.NewCalendarHandler(db, templates).Calendar,
			path:             "/calendar",
			userRole:         models.RoleWarehouse,
			expectedStatus:   http.StatusOK,
			shouldHaveSticky: true,
			expectedLinks: []string{
				`href="/shipments"`,
				`href="/inventory"`,
				`href="/calendar"`,
				`href="/reception-reports"`,
			},
			unexpectedLinks: []string{
				`href="/dashboard"`,
			},
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
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)

			// Add user to context
			ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
			req = req.WithContext(ctx)

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call handler
			tt.handler(rr, req)

			// Check status code
			if rr.Code != tt.expectedStatus {
				t.Fatalf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			html := rr.Body.String()

			// Check for consistent sticky navbar
			if tt.shouldHaveSticky {
				if !strings.Contains(html, "sticky") {
					t.Error("expected sticky positioning in navbar, but not found")
				}
			}

			// Check that navbar structure is present
			if !strings.Contains(html, "Laptop Tracking System") {
				t.Error("expected application title in navbar, but not found")
			}

			// Check for logout link
			if !strings.Contains(html, `href="/logout"`) {
				t.Error("expected logout link in navbar, but not found")
			}

			// Check user email is displayed
			if !strings.Contains(html, user.Email) {
				t.Errorf("expected user email %s in navbar, but not found", user.Email)
			}

			// Check expected links are present
			for _, link := range tt.expectedLinks {
				if !strings.Contains(html, link) {
					t.Errorf("expected link %s in navbar for %s role, but it was not found", link, tt.userRole)
				}
			}

			// Check unexpected links are NOT present
			for _, link := range tt.unexpectedLinks {
				if strings.Contains(html, link) {
					t.Errorf("did not expect link %s in navbar for %s role, but it was found", link, tt.userRole)
				}
			}
		})
	}
}

// loadTemplatesWithNavigation loads templates with navigation helper functions
func loadTemplatesWithNavigation(t *testing.T) *template.Template {
	t.Helper()

	funcMap := template.FuncMap{
		"title": func(v interface{}) string {
			var s string
			switch val := v.(type) {
			case string:
				s = val
			case models.UserRole:
				s = string(val)
			case models.LaptopStatus:
				s = string(val)
			case models.ShipmentStatus:
				s = string(val)
			default:
				s = ""
			}
			return strings.Title(s)
		},
		"replace": func(old, new string, v interface{}) string {
			var s string
			switch val := v.(type) {
			case string:
				s = val
			case models.UserRole:
				s = string(val)
			case models.LaptopStatus:
				s = string(val)
			default:
				s = ""
			}
			return strings.ReplaceAll(s, old, new)
		},
		"add": func(a, b int) int {
			return a + b
		},
		"len": func(v interface{}) int {
			switch val := v.(type) {
			case []models.TimelineItem:
				return len(val)
			case []interface{}:
				return len(val)
			default:
				return 0
			}
		},
		"getNav": func(role models.UserRole) views.NavigationLinks {
			return views.GetNavigationLinks(role)
		},
		// Calendar template functions
		"formatDate": func(t interface{}) string {
			return ""
		},
		"formatTime": func(t interface{}) string {
			return ""
		},
		"formatDateShort": func(t interface{}) string {
			return ""
		},
		"daysInMonth": func(year int, month interface{}) int {
			return 30
		},
		"firstWeekday": func(year int, month interface{}) int {
			return 0
		},
		// Dashboard template functions
		"statusColor": func(status models.ShipmentStatus) string {
			return "bg-gray-400"
		},
		"laptopStatusColor": func(status models.LaptopStatus) string {
			return "bg-gray-400"
		},
		"inventoryStatusColor": func(status models.LaptopStatus) string {
			return "bg-gray-100 text-gray-800"
		},
		"laptopStatusDisplayName": func(status models.LaptopStatus) string {
			return models.GetLaptopStatusDisplayName(status)
		},
	}

	// Parse all templates including the navbar component
	templates, err := template.New("").Funcs(funcMap).ParseGlob("../../templates/pages/*.html")
	if err != nil {
		t.Fatalf("failed to parse page templates: %v", err)
	}

	// Parse navbar component
	templates, err = templates.ParseGlob("../../templates/components/*.html")
	if err != nil {
		t.Fatalf("failed to parse component templates: %v", err)
	}

	return templates
}

