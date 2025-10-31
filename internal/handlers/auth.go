package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/yourusername/laptop-tracking-system/internal/auth"
	"github.com/yourusername/laptop-tracking-system/internal/middleware"
	"github.com/yourusername/laptop-tracking-system/internal/models"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	DB        *sql.DB
	Templates *template.Template
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(db *sql.DB, templates *template.Template) *AuthHandler {
	return &AuthHandler{
		DB:        db,
		Templates: templates,
	}
}

// LoginPage displays the login form
func (h *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	// If already authenticated, redirect to dashboard
	if middleware.IsAuthenticated(r) {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	// Get error message from query parameter if any
	errorMsg := r.URL.Query().Get("error")

	data := map[string]interface{}{
		"Error": errorMsg,
	}

	err := h.Templates.ExecuteTemplate(w, "login.html", data)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// Login handles login form submission
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Redirect(w, r, "/login?error=Invalid+form+data", http.StatusSeeOther)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	// Validate input
	if email == "" || password == "" {
		http.Redirect(w, r, "/login?error=Email+and+password+are+required", http.StatusSeeOther)
		return
	}

	// Find user by email
	var user models.User
	err = h.DB.QueryRowContext(
		r.Context(),
		`SELECT id, email, password_hash, role, google_id, created_at, updated_at
		FROM users
		WHERE email = $1`,
		email,
	).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Role,
		&user.GoogleID, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		http.Redirect(w, r, "/login?error=Invalid+email+or+password", http.StatusSeeOther)
		return
	}
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Check if user authenticated via Google OAuth (no password)
	if user.IsGoogleUser() {
		http.Redirect(w, r, "/login?error=Please+sign+in+with+Google", http.StatusSeeOther)
		return
	}

	// Verify password
	if !auth.CheckPasswordHash(password, user.PasswordHash) {
		http.Redirect(w, r, "/login?error=Invalid+email+or+password", http.StatusSeeOther)
		return
	}

	// Create session
	session, err := auth.CreateSession(r.Context(), h.DB, user.ID, auth.DefaultSessionDuration)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     middleware.SessionCookieName,
		Value:    session.Token,
		Path:     "/",
		Expires:  session.ExpiresAt,
		HttpOnly: true,
		Secure:   true, // Set to true in production with HTTPS
		SameSite: http.SameSiteStrictMode,
	})

	// Redirect to dashboard
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// Logout handles user logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Get session from context
	session := middleware.GetSessionFromContext(r.Context())
	if session != nil {
		// Delete session from database
		_ = auth.DeleteSession(r.Context(), h.DB, session.Token)
	}

	// Clear session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     middleware.SessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	// Redirect to login page
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// ChangePassword handles password change requests
func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	currentPassword := r.FormValue("current_password")
	newPassword := r.FormValue("new_password")
	confirmPassword := r.FormValue("confirm_password")

	// Validate input
	if currentPassword == "" || newPassword == "" || confirmPassword == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	if newPassword != confirmPassword {
		http.Error(w, "New passwords do not match", http.StatusBadRequest)
		return
	}

	// Validate new password strength
	if err := auth.ValidatePassword(newPassword); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Verify current password
	if !auth.CheckPasswordHash(currentPassword, user.PasswordHash) {
		http.Error(w, "Current password is incorrect", http.StatusBadRequest)
		return
	}

	// Hash new password
	newPasswordHash, err := auth.HashPassword(newPassword)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Update password in database
	_, err = h.DB.ExecContext(
		r.Context(),
		`UPDATE users SET password_hash = $1, updated_at = $2 WHERE id = $3`,
		newPasswordHash, time.Now(), user.ID,
	)
	if err != nil {
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		return
	}

	// Delete all user sessions (force re-login)
	_ = auth.DeleteUserSessions(r.Context(), h.DB, user.ID)

	// Clear current session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     middleware.SessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	// Redirect to login with success message
	http.Redirect(w, r, "/login?message=Password+changed+successfully", http.StatusSeeOther)
}

// MagicLinkLogin handles magic link authentication
func (h *AuthHandler) MagicLinkLogin(w http.ResponseWriter, r *http.Request) {
	// Get token from query parameter
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Redirect(w, r, "/login?error=Invalid+magic+link", http.StatusSeeOther)
		return
	}

	// Validate magic link
	magicLink, err := auth.ValidateMagicLink(r.Context(), h.DB, token)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if magicLink == nil {
		http.Redirect(w, r, "/login?error=Magic+link+is+invalid+or+has+expired", http.StatusSeeOther)
		return
	}

	// Mark magic link as used
	err = auth.MarkMagicLinkAsUsed(r.Context(), h.DB, token)
	if err != nil {
		http.Error(w, "Failed to use magic link", http.StatusInternalServerError)
		return
	}

	// Create session
	session, err := auth.CreateSession(r.Context(), h.DB, magicLink.UserID, auth.DefaultSessionDuration)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     middleware.SessionCookieName,
		Value:    session.Token,
		Path:     "/",
		Expires:  session.ExpiresAt,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	// Redirect based on context
	redirectURL := "/dashboard"
	if magicLink.ShipmentID != nil {
		// If magic link is associated with a shipment, redirect to shipment form
		redirectURL = fmt.Sprintf("/shipments/%d/form", *magicLink.ShipmentID)
	}

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

// SendMagicLink generates and sends a magic link to the user
// This would typically be called when sending the pickup form email
func (h *AuthHandler) SendMagicLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Only logistics users can send magic links
	user := middleware.GetUserFromContext(r.Context())
	if user == nil || user.Role != models.RoleLogistics {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	shipmentIDStr := r.FormValue("shipment_id")

	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	// Find or create user by email
	var userID int64
	err = h.DB.QueryRowContext(
		r.Context(),
		`SELECT id FROM users WHERE email = $1`,
		email,
	).Scan(&userID)

	if err == sql.ErrNoRows {
		// User doesn't exist, create a client user
		err = h.DB.QueryRowContext(
			r.Context(),
			`INSERT INTO users (email, role, created_at, updated_at)
			VALUES ($1, $2, $3, $4)
			RETURNING id`,
			email, models.RoleClient, time.Now(), time.Now(),
		).Scan(&userID)
		if err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}
	} else if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Parse shipment ID if provided
	var shipmentID *int64
	if shipmentIDStr != "" {
		var sid int64
		_, err := fmt.Sscanf(shipmentIDStr, "%d", &sid)
		if err == nil {
			shipmentID = &sid
		}
	}

	// Create magic link
	magicLink, err := auth.CreateMagicLink(r.Context(), h.DB, userID, shipmentID, auth.DefaultMagicLinkDuration)
	if err != nil {
		http.Error(w, "Failed to create magic link", http.StatusInternalServerError)
		return
	}

	// TODO: Send email with magic link
	// For now, just return the magic link URL (in production, this would be sent via email)
	magicLinkURL := fmt.Sprintf("https://yourdomain.com/auth/magic-link?token=%s", magicLink.Token)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"magic_link_url": "%s", "message": "Magic link created successfully"}`, magicLinkURL)
}

