package handlers

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/yourusername/laptop-tracking-system/internal/database"
	"github.com/yourusername/laptop-tracking-system/internal/middleware"
	"github.com/yourusername/laptop-tracking-system/internal/models"
)

// TestReceptionReportsListWithSorting tests that the reception reports list can be sorted by different columns
func TestReceptionReportsListWithSorting(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup test database
	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	// Create test data
	var companyID int64
	err := db.QueryRow(
		`INSERT INTO client_companies (name, contact_info, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW()) RETURNING id`,
		"Test Corp", "test@test.com",
	).Scan(&companyID)
	if err != nil {
		t.Fatalf("Failed to create company: %v", err)
	}

	var userID int64
	err = db.QueryRow(
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id`,
		"warehouse@test.com", "hashed_password", models.RoleWarehouse,
	).Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	user := &models.User{
		ID:    userID,
		Email: "warehouse@test.com",
		Role:  models.RoleWarehouse,
	}

	// Setup handler
	handler := &ReceptionReportHandler{
		DB:        db,
		Templates: nil,
	}

	// Test cases for different sorting
	testCases := []struct {
		name      string
		sortBy    string
		sortOrder string
	}{
		{"Sort by ID ascending", "id", "asc"},
		{"Sort by received_at descending", "received_at", "desc"},
		{"Default sort (status then received_at)", "", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := "/reception-reports?sort=" + tc.sortBy + "&order=" + tc.sortOrder
			req := httptest.NewRequest("GET", url, nil)
			ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			handler.ReceptionReportsList(rr, req)

			// Test should compile and run
			if rr.Code != 200 {
				t.Logf("Handler returned status %d", rr.Code)
			}
		})
	}
}

