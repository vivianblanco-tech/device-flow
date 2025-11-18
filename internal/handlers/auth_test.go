package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/yourusername/laptop-tracking-system/internal/auth"
	"github.com/yourusername/laptop-tracking-system/internal/database"
	"github.com/yourusername/laptop-tracking-system/internal/middleware"
	"github.com/yourusername/laptop-tracking-system/internal/models"
)

// loadTestTemplates loads templates for testing
func loadTestTemplates(t *testing.T) *template.Template {
	funcMap := template.FuncMap{
		"replace": func(old, new string, v interface{}) string {
			// Convert interface{} to string first
			var s string
			switch val := v.(type) {
			case string:
				s = val
			case models.UserRole:
				s = string(val)
			case models.LaptopStatus:
				s = string(val)
			default:
				s = fmt.Sprintf("%v", val)
			}
			return strings.ReplaceAll(s, old, new)
		},
		"title": func(v interface{}) string {
			// Convert interface{} to string
			var s string
			switch val := v.(type) {
			case string:
				s = val
			case models.UserRole:
				s = string(val)
			case models.LaptopStatus:
				s = string(val)
			default:
				s = fmt.Sprintf("%v", val)
			}
			return strings.Title(s)
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
		// Calendar template functions
		"formatDate": func(t time.Time) string {
			return t.Format("Jan 2, 2006")
		},
		"formatTime": func(t time.Time) string {
			return t.Format("3:04 PM")
		},
		"formatDateShort": func(t time.Time) string {
			return t.Format("Jan 2")
		},
		"daysInMonth": func(year int, month time.Month) int {
			return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
		},
		"firstWeekday": func(year int, month time.Month) time.Weekday {
			return time.Date(year, month, 1, 0, 0, 0, 0, time.UTC).Weekday()
		},
		// Dashboard template functions
		"statusColor": func(status models.ShipmentStatus) string {
			switch status {
			case models.ShipmentStatusPendingPickup:
				return "bg-yellow-400"
			case models.ShipmentStatusPickedUpFromClient:
				return "bg-orange-400"
			case models.ShipmentStatusInTransitToWarehouse:
				return "bg-purple-400"
			case models.ShipmentStatusAtWarehouse:
				return "bg-indigo-400"
			case models.ShipmentStatusReleasedFromWarehouse:
				return "bg-blue-400"
			case models.ShipmentStatusInTransitToEngineer:
				return "bg-cyan-400"
			case models.ShipmentStatusDelivered:
				return "bg-green-400"
			default:
				return "bg-gray-400"
			}
		},
		"laptopStatusColor": func(status models.LaptopStatus) string {
			switch status {
			case models.LaptopStatusAvailable:
				return "bg-green-400"
			case models.LaptopStatusInTransitToWarehouse:
				return "bg-purple-400"
			case models.LaptopStatusAtWarehouse:
				return "bg-indigo-400"
			case models.LaptopStatusInTransitToEngineer:
				return "bg-cyan-400"
			case models.LaptopStatusDelivered:
				return "bg-blue-400"
			case models.LaptopStatusRetired:
				return "bg-gray-400"
			default:
				return "bg-gray-400"
			}
		},
		"inventoryStatusColor": func(status models.LaptopStatus) string {
			switch status {
			case models.LaptopStatusAvailable:
				return "bg-green-100 text-green-800"
			case models.LaptopStatusInTransitToWarehouse:
				return "bg-purple-100 text-purple-800"
			case models.LaptopStatusAtWarehouse:
				return "bg-indigo-100 text-indigo-800"
			case models.LaptopStatusInTransitToEngineer:
				return "bg-cyan-100 text-cyan-800"
			case models.LaptopStatusDelivered:
				return "bg-blue-100 text-blue-800"
			case models.LaptopStatusRetired:
				return "bg-gray-100 text-gray-800"
			default:
				return "bg-gray-100 text-gray-800"
			}
		},
		"laptopStatusDisplayName": func(status models.LaptopStatus) string {
			return models.GetLaptopStatusDisplayName(status)
		},
		"receptionReportStatusColor": func(status string) string {
			switch models.ReceptionReportStatus(status) {
			case models.ReceptionReportStatusPendingApproval:
				return "bg-yellow-100 text-yellow-800"
			case models.ReceptionReportStatusApproved:
				return "bg-green-100 text-green-800"
			default:
				return "bg-gray-100 text-gray-800"
			}
		},
		"receptionReportStatusDisplayName": func(status string) string {
			switch models.ReceptionReportStatus(status) {
			case models.ReceptionReportStatusPendingApproval:
				return "Pending Approval"
			case models.ReceptionReportStatusApproved:
				return "Approved"
			default:
				return "Unknown"
			}
		},
		"printf": fmt.Sprintf,
	}

	templates, err := template.New("").Funcs(funcMap).ParseGlob("../../templates/pages/*.html")
	if err != nil {
		t.Fatalf("Failed to parse page templates: %v", err)
	}

	// Also load component templates
	templates, err = templates.ParseGlob("../../templates/components/*.html")
	if err != nil {
		t.Fatalf("Failed to parse component templates: %v", err)
	}

	return templates
}

func TestLoginRedirectByRole(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	templates := loadTestTemplates(t)
	handler := NewAuthHandler(db, templates)

	// Create test users with different roles
	testUsers := []struct {
		email            string
		role             models.UserRole
		expectedRedirect string
	}{
		{
			email:            "logistics.test@bairesdev.com",
			role:             models.RoleLogistics,
			expectedRedirect: "/dashboard",
		},
		{
			email:            "client.test@bairesdev.com",
			role:             models.RoleClient,
			expectedRedirect: "/shipments",
		},
		{
			email:            "warehouse.test@bairesdev.com",
			role:             models.RoleWarehouse,
			expectedRedirect: "/inventory",
		},
		{
			email:            "pm.test@bairesdev.com",
			role:             models.RoleProjectManager,
			expectedRedirect: "/dashboard",
		},
	}

	ctx := context.Background()
	password := "Test123!"
	passwordHash, _ := auth.HashPassword(password)

	for _, tt := range testUsers {
		t.Run(fmt.Sprintf("%s user redirects to %s", tt.role, tt.expectedRedirect), func(t *testing.T) {
			// Create user
			var userID int64
			err := db.QueryRowContext(ctx,
				`INSERT INTO users (email, password_hash, role, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5) RETURNING id`,
				tt.email, passwordHash, tt.role, time.Now(), time.Now(),
			).Scan(&userID)
			if err != nil {
				t.Fatalf("Failed to create test user: %v", err)
			}

			// Create login request
			form := url.Values{}
			form.Add("email", tt.email)
			form.Add("password", password)

			req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()

			// Call login handler
			handler.Login(w, req)

			// Check status code
			if w.Code != http.StatusSeeOther {
				t.Errorf("expected status %d, got %d", http.StatusSeeOther, w.Code)
			}

			// Check redirect location
			location := w.Header().Get("Location")
			if location != tt.expectedRedirect {
				t.Errorf("expected redirect to %s, got %s", tt.expectedRedirect, location)
			}

			// Cleanup: delete user
			_, _ = db.ExecContext(ctx, "DELETE FROM sessions WHERE user_id = $1", userID)
			_, _ = db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", userID)
		})
	}
}

func TestLoginPage(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	templates := loadTestTemplates(t)
	handler := NewAuthHandler(db, templates)

	t.Run("GET request displays login page", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/login", nil)
		w := httptest.NewRecorder()

		handler.LoginPage(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})

	t.Run("authenticated user redirects to dashboard", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/login", nil)

		// Create a test user in context
		user := &models.User{ID: 1, Email: "test@example.com", Role: models.RoleLogistics}
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		handler.LoginPage(w, req)

		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected status 303, got %d", w.Code)
		}

		location := w.Header().Get("Location")
		if location != "/dashboard" {
			t.Errorf("Expected redirect to /dashboard, got %s", location)
		}
	})

	t.Run("error message from query parameter", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/login?error=Invalid+credentials", nil)
		w := httptest.NewRecorder()

		handler.LoginPage(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})
}

func TestLogin(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create test user
	password := "TestPass123!"
	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	var userID int64
	err = db.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"test@example.com", passwordHash, models.RoleLogistics, time.Now(), time.Now(),
	).Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	templates := loadTestTemplates(t)
	handler := NewAuthHandler(db, templates)

	t.Run("successful login with valid credentials", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("email", "test@example.com")
		formData.Set("password", password)

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		w := httptest.NewRecorder()
		handler.Login(w, req)

		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected status 303, got %d", w.Code)
		}

		// Check session cookie was set
		cookies := w.Result().Cookies()
		found := false
		for _, cookie := range cookies {
			if cookie.Name == middleware.SessionCookieName {
				found = true
				if cookie.Value == "" {
					t.Error("Session cookie value is empty")
				}
			}
		}
		if !found {
			t.Error("Session cookie not set")
		}
	})

	t.Run("login with non-POST method", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/login", nil)
		w := httptest.NewRecorder()

		handler.Login(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status 405, got %d", w.Code)
		}
	})

	t.Run("login with missing credentials", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("email", "")
		formData.Set("password", "")

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		w := httptest.NewRecorder()
		handler.Login(w, req)

		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected status 303 (redirect with error), got %d", w.Code)
		}

		location := w.Header().Get("Location")
		if !strings.Contains(location, "error=") {
			t.Errorf("Expected error in redirect URL, got: %s", location)
		}
	})

	t.Run("login with invalid email", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("email", "nonexistent@example.com")
		formData.Set("password", password)

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		w := httptest.NewRecorder()
		handler.Login(w, req)

		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected status 303, got %d", w.Code)
		}

		location := w.Header().Get("Location")
		if !strings.Contains(location, "error=") {
			t.Errorf("Expected error in redirect URL, got: %s", location)
		}
	})

	t.Run("login with incorrect password", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("email", "test@example.com")
		formData.Set("password", "WrongPassword123!")

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		w := httptest.NewRecorder()
		handler.Login(w, req)

		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected status 303, got %d", w.Code)
		}

		location := w.Header().Get("Location")
		if !strings.Contains(location, "error=") {
			t.Errorf("Expected error in redirect URL, got: %s", location)
		}
	})
}

func TestLogout(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create test user
	var userID int64
	err := db.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"test@example.com", "hashedpassword", models.RoleLogistics, time.Now(), time.Now(),
	).Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Create session
	session, err := auth.CreateSession(ctx, db, userID, auth.DefaultSessionDuration)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	templates := loadTestTemplates(t)
	handler := NewAuthHandler(db, templates)

	t.Run("logout deletes session", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/logout", nil)

		// Add session to context
		reqCtx := context.WithValue(req.Context(), middleware.SessionContextKey, session)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.Logout(w, req)

		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected status 303, got %d", w.Code)
		}

		// Check session cookie was cleared
		cookies := w.Result().Cookies()
		found := false
		for _, cookie := range cookies {
			if cookie.Name == middleware.SessionCookieName {
				found = true
				if cookie.MaxAge != -1 {
					t.Error("Session cookie MaxAge should be -1 to clear it")
				}
			}
		}
		if !found {
			t.Error("Session cookie not cleared")
		}

		// Verify session was deleted from database
		validatedSession, _ := auth.ValidateSession(ctx, db, session.Token)
		if validatedSession != nil {
			t.Error("Session should have been deleted from database")
		}
	})
}

func TestChangePassword(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create test user
	currentPassword := "CurrentPass123!"
	passwordHash, err := auth.HashPassword(currentPassword)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	var userID int64
	err = db.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"test@example.com", passwordHash, models.RoleLogistics, time.Now(), time.Now(),
	).Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	templates := loadTestTemplates(t)
	handler := NewAuthHandler(db, templates)

	t.Run("successful password change", func(t *testing.T) {
		newPassword := "NewPass123!"
		formData := url.Values{}
		formData.Set("current_password", currentPassword)
		formData.Set("new_password", newPassword)
		formData.Set("confirm_password", newPassword)

		req := httptest.NewRequest(http.MethodPost, "/change-password", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		// Add user to context
		user := &models.User{ID: userID, Email: "test@example.com", PasswordHash: passwordHash, Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.ChangePassword(w, req)

		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected status 303, got %d", w.Code)
		}

		// Verify password was updated in database
		var newHash string
		err := db.QueryRowContext(ctx, "SELECT password_hash FROM users WHERE id = $1", userID).Scan(&newHash)
		if err != nil {
			t.Fatalf("Failed to query updated password: %v", err)
		}

		if !auth.CheckPasswordHash(newPassword, newHash) {
			t.Error("Password was not updated correctly")
		}
	})

	t.Run("change password with non-POST method", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/change-password", nil)
		w := httptest.NewRecorder()

		handler.ChangePassword(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status 405, got %d", w.Code)
		}
	})

	t.Run("change password without authentication", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("current_password", currentPassword)
		formData.Set("new_password", "NewPass123!")
		formData.Set("confirm_password", "NewPass123!")

		req := httptest.NewRequest(http.MethodPost, "/change-password", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		w := httptest.NewRecorder()
		handler.ChangePassword(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", w.Code)
		}
	})

	t.Run("change password with mismatched new passwords", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("current_password", currentPassword)
		formData.Set("new_password", "NewPass123!")
		formData.Set("confirm_password", "DifferentPass123!")

		req := httptest.NewRequest(http.MethodPost, "/change-password", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		user := &models.User{ID: userID, Email: "test@example.com", PasswordHash: passwordHash, Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.ChangePassword(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("change password with weak new password", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("current_password", currentPassword)
		formData.Set("new_password", "weak")
		formData.Set("confirm_password", "weak")

		req := httptest.NewRequest(http.MethodPost, "/change-password", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		user := &models.User{ID: userID, Email: "test@example.com", PasswordHash: passwordHash, Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.ChangePassword(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})
}

func TestMagicLinkLogin(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create test user
	var userID int64
	err := db.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"test@example.com", "hashedpassword", models.RoleClient, time.Now(), time.Now(),
	).Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Create magic link
	magicLink, err := auth.CreateMagicLink(ctx, db, userID, nil, auth.DefaultMagicLinkDuration)
	if err != nil {
		t.Fatalf("Failed to create magic link: %v", err)
	}

	templates := loadTestTemplates(t)
	handler := NewAuthHandler(db, templates)

	t.Run("successful magic link login", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/auth/magic-link?token="+magicLink.Token, nil)
		w := httptest.NewRecorder()

		handler.MagicLinkLogin(w, req)

		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected status 303, got %d", w.Code)
		}

		// Check session cookie was set
		cookies := w.Result().Cookies()
		found := false
		for _, cookie := range cookies {
			if cookie.Name == middleware.SessionCookieName && cookie.Value != "" {
				found = true
			}
		}
		if !found {
			t.Error("Session cookie not set")
		}
	})

	t.Run("magic link login without token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/auth/magic-link", nil)
		w := httptest.NewRecorder()

		handler.MagicLinkLogin(w, req)

		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected status 303, got %d", w.Code)
		}

		location := w.Header().Get("Location")
		if !strings.Contains(location, "error=") {
			t.Errorf("Expected error in redirect URL, got: %s", location)
		}
	})

	t.Run("magic link login with invalid token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/auth/magic-link?token=invalid_token", nil)
		w := httptest.NewRecorder()

		handler.MagicLinkLogin(w, req)

		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected status 303, got %d", w.Code)
		}

		location := w.Header().Get("Location")
		if !strings.Contains(location, "error=") {
			t.Errorf("Expected error in redirect URL, got: %s", location)
		}
	})

	t.Run("magic link should NOT be marked as used when clicked", func(t *testing.T) {
		// Create a new magic link for this test
		magicLink, err := auth.CreateMagicLink(ctx, db, userID, nil, auth.DefaultMagicLinkDuration)
		if err != nil {
			t.Fatalf("Failed to create magic link: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/auth/magic-link?token="+magicLink.Token, nil)
		w := httptest.NewRecorder()

		handler.MagicLinkLogin(w, req)

		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected status 303, got %d", w.Code)
		}

		// Verify magic link is still valid (not marked as used)
		validatedLink, err := auth.ValidateMagicLink(ctx, db, magicLink.Token)
		if err != nil {
			t.Fatalf("Failed to validate magic link: %v", err)
		}
		if validatedLink == nil {
			t.Error("Magic link should still be valid after clicking, but it was marked as used or expired")
		}
		if validatedLink.IsUsed() {
			t.Error("Magic link should NOT be marked as used when clicked")
		}
	})

	t.Run("magic link expiration duration should be 72 hours", func(t *testing.T) {
		// Create a magic link with default duration
		magicLink, err := auth.CreateMagicLink(ctx, db, userID, nil, auth.DefaultMagicLinkDuration)
		if err != nil {
			t.Fatalf("Failed to create magic link: %v", err)
		}

		// Verify DefaultMagicLinkDuration is 72 hours
		if auth.DefaultMagicLinkDuration != 72 {
			t.Errorf("Expected DefaultMagicLinkDuration to be 72 hours, got %d", auth.DefaultMagicLinkDuration)
		}

		// Verify expiration is approximately 72 hours from now
		expectedExpiration := time.Now().Add(72 * time.Hour)
		actualExpiration := magicLink.ExpiresAt
		diff := expectedExpiration.Sub(actualExpiration)
		if diff < 0 {
			diff = -diff
		}
		// Allow 1 minute tolerance for test execution time
		if diff > time.Minute {
			t.Errorf("Expected expiration to be approximately 72 hours from now, got %v (diff: %v)", actualExpiration, diff)
		}
	})
}

func TestSendMagicLink(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create test logistics user
	var logisticsUserID int64
	err := db.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"logistics@example.com", "hashedpassword", models.RoleLogistics, time.Now(), time.Now(),
	).Scan(&logisticsUserID)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Create test company
	var companyID int64
	err = db.QueryRowContext(ctx,
		`INSERT INTO client_companies (name, contact_info, created_at)
		VALUES ($1, $2, $3) RETURNING id`,
		"Test Company", json.RawMessage(`{"email":"test@company.com"}`), time.Now(),
	).Scan(&companyID)
	if err != nil {
		t.Fatalf("Failed to create test company: %v", err)
	}

	// Create test shipment with JIRA ticket
	var shipmentID int64
	err = db.QueryRowContext(ctx,
		`INSERT INTO shipments (client_company_id, status, jira_ticket_number, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		companyID, models.ShipmentStatusPendingPickup, "TEST-1234", time.Now(), time.Now(),
	).Scan(&shipmentID)
	if err != nil {
		t.Fatalf("Failed to create test shipment: %v", err)
	}

	templates := loadTestTemplates(t)
	handler := NewAuthHandler(db, templates)

	t.Run("logistics user can send magic link with valid shipment", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("email", "newclient@example.com")
		formData.Set("shipment_id", strconv.FormatInt(shipmentID, 10))

		req := httptest.NewRequest(http.MethodPost, "/auth/send-magic-link", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		user := &models.User{ID: logisticsUserID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.SendMagicLink(w, req)

		// SendMagicLink redirects to the shipment detail page with success message
		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected status 303 (SeeOther), got %d", w.Code)
		}

		// Verify redirect location contains the shipment ID
		location := w.Header().Get("Location")
		expectedPath := fmt.Sprintf("/shipments/%d", shipmentID)
		if !strings.Contains(location, expectedPath) {
			t.Errorf("Expected redirect to contain %s, got %s", expectedPath, location)
		}
	})

	t.Run("cannot send magic link without shipment ID", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("email", "newclient@example.com")

		req := httptest.NewRequest(http.MethodPost, "/auth/send-magic-link", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		user := &models.User{ID: logisticsUserID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.SendMagicLink(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("cannot send magic link with non-existent shipment", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("email", "newclient@example.com")
		formData.Set("shipment_id", "99999")

		req := httptest.NewRequest(http.MethodPost, "/auth/send-magic-link", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		user := &models.User{ID: logisticsUserID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.SendMagicLink(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", w.Code)
		}
	})

	t.Run("non-logistics user cannot send magic link", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("email", "newclient@example.com")

		req := httptest.NewRequest(http.MethodPost, "/auth/send-magic-link", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		user := &models.User{ID: logisticsUserID, Email: "client@example.com", Role: models.RoleClient}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.SendMagicLink(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("Expected status 403, got %d", w.Code)
		}
	})
}
