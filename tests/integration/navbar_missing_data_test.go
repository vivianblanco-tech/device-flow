package integration

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/yourusername/laptop-tracking-system/internal/database"
	"github.com/yourusername/laptop-tracking-system/internal/email"
	"github.com/yourusername/laptop-tracking-system/internal/handlers"
	"github.com/yourusername/laptop-tracking-system/internal/middleware"
	"github.com/yourusername/laptop-tracking-system/internal/models"
	"github.com/yourusername/laptop-tracking-system/internal/views"
)

// TestNavbarDataOnAllPages verifies that all pages requiring authentication pass Nav and CurrentPage data
func TestNavbarDataOnAllPages(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup test database
	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	// Load templates with all necessary template functions
	templates := loadTemplatesWithNavigation(t)

	// Setup email client for handlers that need it
	emailClient, err := email.NewClient(email.Config{
		Host: "localhost",
		Port: 1025,
		From: "test@example.com",
	})
	if err != nil {
		t.Fatalf("Failed to create email client: %v", err)
	}
	emailNotifier := email.NewNotifier(emailClient, db)

	// Create test data
	setupTestData(t, db)

	tests := []struct {
		name             string
		setupHandler     func() http.HandlerFunc
		path             string
		method           string
		userRole         models.UserRole
		expectedStatus   int
		shouldHaveNav    bool
		expectedPage     string
	}{
		{
			name: "Laptop Detail Page",
			setupHandler: func() http.HandlerFunc {
				h := handlers.NewInventoryHandler(db, templates)
				return h.LaptopDetail
			},
			path:           "/inventory/1",
			method:         http.MethodGet,
			userRole:       models.RoleLogistics,
			expectedStatus: http.StatusOK,
			shouldHaveNav:  true,
			expectedPage:   "inventory",
		},
		{
			name: "Add Laptop Page",
			setupHandler: func() http.HandlerFunc {
				h := handlers.NewInventoryHandler(db, templates)
				return h.AddLaptopPage
			},
			path:           "/inventory/add",
			method:         http.MethodGet,
			userRole:       models.RoleLogistics,
			expectedStatus: http.StatusOK,
			shouldHaveNav:  true,
			expectedPage:   "inventory",
		},
		{
			name: "Edit Laptop Page",
			setupHandler: func() http.HandlerFunc {
				h := handlers.NewInventoryHandler(db, templates)
				return h.EditLaptopPage
			},
			path:           "/inventory/1/edit",
			method:         http.MethodGet,
			userRole:       models.RoleLogistics,
			expectedStatus: http.StatusOK,
			shouldHaveNav:  true,
			expectedPage:   "inventory",
		},
		{
			name: "Reception Report Page",
			setupHandler: func() http.HandlerFunc {
				h := handlers.NewReceptionReportHandler(db, templates, emailNotifier)
				return h.ReceptionReportPage
			},
			path:           "/reception-report?shipment_id=1",
			method:         http.MethodGet,
			userRole:       models.RoleWarehouse,
			expectedStatus: http.StatusOK,
			shouldHaveNav:  true,
			expectedPage:   "reception-reports",
		},
		{
			name: "Shipment Detail Page",
			setupHandler: func() http.HandlerFunc {
				h := handlers.NewShipmentsHandler(db, templates, emailNotifier)
				return h.ShipmentDetail
			},
			path:           "/shipments/1",
			method:         http.MethodGet,
			userRole:       models.RoleLogistics,
			expectedStatus: http.StatusOK,
			shouldHaveNav:  true,
			expectedPage:   "shipments",
		},
		{
			name: "Delivery Form Page",
			setupHandler: func() http.HandlerFunc {
				h := handlers.NewDeliveryFormHandler(db, templates, emailNotifier)
				return h.DeliveryFormPage
			},
			path:           "/delivery-form?shipment_id=1",
			method:         http.MethodGet,
			userRole:       models.RoleLogistics,
			expectedStatus: http.StatusOK,
			shouldHaveNav:  true,
			expectedPage:   "shipments",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test user
			user := &models.User{
				ID:    1,
				Email: "test@bairesdev.com",
				Role:  tt.userRole,
			}

			// Create request
			req := httptest.NewRequest(tt.method, tt.path, nil)

			// Add mux vars if path contains ID
			if strings.Contains(tt.path, "/1") {
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
			}

			// Add user to context
			ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
			req = req.WithContext(ctx)

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call handler
			handler := tt.setupHandler()
			handler(rr, req)

			// Check status code
			if rr.Code != tt.expectedStatus {
				t.Fatalf("expected status %d, got %d. Response body: %s", tt.expectedStatus, rr.Code, rr.Body.String())
			}

			html := rr.Body.String()

			// Check for navbar presence
			if tt.shouldHaveNav {
				if !strings.Contains(html, "Laptop Tracking System") {
					t.Error("expected navbar title 'Laptop Tracking System', but not found")
				}

				if !strings.Contains(html, `href="/logout"`) {
					t.Error("expected logout link in navbar, but not found")
				}

				if !strings.Contains(html, user.Email) {
					t.Errorf("expected user email %s in navbar, but not found", user.Email)
				}

				// Check for sticky positioning
				if !strings.Contains(html, "sticky") {
					t.Error("expected sticky positioning in navbar, but not found")
				}

				// Check navigation links are present based on role
				navLinks := views.GetNavigationLinks(tt.userRole)
				
				if navLinks.Dashboard && !strings.Contains(html, `href="/dashboard"`) {
					t.Errorf("expected dashboard link for %s role, but not found", tt.userRole)
				}

				if navLinks.Shipments && !strings.Contains(html, `href="/shipments"`) {
					t.Errorf("expected shipments link for %s role, but not found", tt.userRole)
				}

				if navLinks.Inventory && !strings.Contains(html, `href="/inventory"`) {
					t.Errorf("expected inventory link for %s role, but not found", tt.userRole)
				}
			}
		})
	}
}

// setupTestData creates test data in the database
func setupTestData(t *testing.T, db *sql.DB) {
	t.Helper()

	// Create test client company
	_, err := db.Exec(`
		INSERT INTO client_companies (id, name, contact_info) 
		VALUES (1, 'Test Company', 'test@company.com')
		ON CONFLICT (id) DO NOTHING
	`)
	if err != nil {
		t.Fatalf("failed to create test company: %v", err)
	}

	// Create test software engineer
	_, err = db.Exec(`
		INSERT INTO software_engineers (id, name, email) 
		VALUES (1, 'Test Engineer', 'engineer@test.com')
		ON CONFLICT (id) DO NOTHING
	`)
	if err != nil {
		t.Fatalf("failed to create test engineer: %v", err)
	}

	// Create test laptop
	_, err = db.Exec(`
		INSERT INTO laptops (id, serial_number, sku, brand, model, cpu, ram_gb, ssd_gb, status, client_company_id) 
		VALUES (1, 'TEST123', 'TST-001', 'Dell', 'Latitude', 'i7', '16GB', '512GB', 'available', 1)
		ON CONFLICT (id) DO NOTHING
	`)
	if err != nil {
		t.Fatalf("failed to create test laptop: %v", err)
	}

	// Create test shipment
	_, err = db.Exec(`
		INSERT INTO shipments (id, jira_ticket_number, client_company_id, status, shipment_type) 
		VALUES (1, 'TEST-123', 1, 'pending_pickup_from_client', 'single_full_journey')
		ON CONFLICT (id) DO NOTHING
	`)
	if err != nil {
		t.Fatalf("failed to create test shipment: %v", err)
	}
}

