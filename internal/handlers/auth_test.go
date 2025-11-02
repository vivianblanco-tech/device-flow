package handlers

import (
	"context"
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
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
		"replace": func(old, new, s string) string {
			return strings.ReplaceAll(s, old, new)
		},
		"title": func(s string) string {
			return strings.Title(s)
		},
	}

	templates, err := template.New("").Funcs(funcMap).ParseGlob("../../templates/pages/*.html")
	if err != nil {
		t.Fatalf("Failed to parse templates: %v", err)
	}
	return templates
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

	templates := loadTestTemplates(t)
	handler := NewAuthHandler(db, templates)

	t.Run("logistics user can send magic link", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("email", "newclient@example.com")

		req := httptest.NewRequest(http.MethodPost, "/auth/send-magic-link", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		user := &models.User{ID: logisticsUserID, Email: "logistics@example.com", Role: models.RoleLogistics}
		reqCtx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(reqCtx)

		w := httptest.NewRecorder()
		handler.SendMagicLink(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
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
