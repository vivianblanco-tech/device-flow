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

