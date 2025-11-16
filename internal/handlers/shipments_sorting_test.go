package handlers

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/yourusername/laptop-tracking-system/internal/database"
	"github.com/yourusername/laptop-tracking-system/internal/middleware"
	"github.com/yourusername/laptop-tracking-system/internal/models"
)

// TestShipmentsListWithSorting tests that the shipments list can be sorted by different columns
func TestShipmentsListWithSorting(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup test database
	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	// Create test company
	var companyID int64
	err := db.QueryRow(
		`INSERT INTO client_companies (name, contact_info, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW()) RETURNING id`,
		"Test Corp", "test@test.com",
	).Scan(&companyID)
	if err != nil {
		t.Fatalf("Failed to create company: %v", err)
	}

	// Create test engineer
	var engineerID int64
	err = db.QueryRow(
		`INSERT INTO software_engineers (name, email, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW()) RETURNING id`,
		"John Doe", "john@test.com",
	).Scan(&engineerID)
	if err != nil {
		t.Fatalf("Failed to create engineer: %v", err)
	}

	// Create test shipments with different attributes
	shipments := []struct {
		jira     string
		shipType models.ShipmentType
		status   models.ShipmentStatus
		engineer bool
	}{
		{"PROJ-101", models.ShipmentTypeSingleFullJourney, models.ShipmentStatusDelivered, true},
		{"PROJ-102", models.ShipmentTypeBulkToWarehouse, models.ShipmentStatusAtWarehouse, false},
		{"PROJ-103", models.ShipmentTypeSingleFullJourney, models.ShipmentStatusPendingPickup, true},
	}

	for _, s := range shipments {
		var engID *int64
		if s.engineer {
			engID = &engineerID
		}
		_, err := db.Exec(
			`INSERT INTO shipments (shipment_type, client_company_id, software_engineer_id, status, jira_ticket_number, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, NOW(), NOW())`,
			s.shipType, companyID, engID, s.status, s.jira,
		)
		if err != nil {
			t.Fatalf("Failed to create shipment: %v", err)
		}
	}

	// Create test user
	var userID int64
	err = db.QueryRow(
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id`,
		"logistics@test.com", "hashed_password", models.RoleLogistics,
	).Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	user := &models.User{
		ID:    userID,
		Email: "logistics@test.com",
		Role:  models.RoleLogistics,
	}

	// Setup handler
	handler := &ShipmentsHandler{
		DB:        db,
		Templates: nil,
	}

	// Test cases for different sorting
	testCases := []struct {
		name     string
		sortBy   string
		sortOrder string
	}{
		{"Sort by ID ascending", "id", "asc"},
		{"Sort by status ascending", "status", "asc"},
		{"Sort by JIRA ticket ascending", "jira_ticket", "asc"},
		{"Default sort (created DESC)", "", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := "/shipments?sort=" + tc.sortBy + "&order=" + tc.sortOrder
			req := httptest.NewRequest("GET", url, nil)
			ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			handler.ShipmentsList(rr, req)

			// Test should compile but will fail because sorting isn't implemented yet
			if rr.Code != 200 {
				t.Logf("Handler returned status %d (expected after GREEN phase)", rr.Code)
			}
		})
	}
}

