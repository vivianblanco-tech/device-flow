package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"time"

	"golang.org/x/oauth2"

	"github.com/yourusername/laptop-tracking-system/internal/auth"
	"github.com/yourusername/laptop-tracking-system/internal/middleware"
	"github.com/yourusername/laptop-tracking-system/internal/models"
	"github.com/yourusername/laptop-tracking-system/internal/views"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	DB          *sql.DB
	Templates   *template.Template
	OAuthConfig *oauth2.Config
	OAuthDomain string // Allowed domain for Google OAuth
}

// isProduction checks if the application is running in production
func isProduction() bool {
	env := os.Getenv("APP_ENV")
	return env == "production"
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(db *sql.DB, templates *template.Template) *AuthHandler {
	return &AuthHandler{
		DB:        db,
		Templates: templates,
	}
}

// roleRedirects maps user roles to their default landing pages
var roleRedirects = map[models.UserRole]string{
	models.RoleClient:         "/shipments",
	models.RoleWarehouse:      "/inventory",
	models.RoleLogistics:      "/dashboard",
	models.RoleProjectManager: "/dashboard",
}

// getRedirectURLForRole returns the appropriate redirect URL based on user role
func getRedirectURLForRole(role models.UserRole) string {
	if url, ok := roleRedirects[role]; ok {
		return url
	}
	// Default fallback for any unmapped roles
	return "/dashboard"
}

// LoginPage displays the login form
func (h *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	// If already authenticated, redirect to appropriate page for their role
	if middleware.IsAuthenticated(r) {
		user := middleware.GetUserFromContext(r.Context())
		if user != nil {
			redirectURL := getRedirectURLForRole(user.Role)
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return
		}
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
		Secure:   isProduction(), // Only require HTTPS in production
		SameSite: http.SameSiteStrictMode,
	})

	// Redirect based on user role
	redirectURL := getRedirectURLForRole(user.Role)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
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
		Secure:   isProduction(),
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
		Secure:   isProduction(),
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
		Secure:   isProduction(),
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

	// Shipment ID is required
	if shipmentIDStr == "" {
		http.Error(w, "Shipment ID is required", http.StatusBadRequest)
		return
	}

	// Parse shipment ID
	var sid int64
	_, err = fmt.Sscanf(shipmentIDStr, "%d", &sid)
	if err != nil {
		http.Error(w, "Invalid shipment ID", http.StatusBadRequest)
		return
	}

	// Verify shipment exists and has a JIRA ticket
	var jiraTicket string
	err = h.DB.QueryRowContext(
		r.Context(),
		`SELECT jira_ticket_number FROM shipments WHERE id = $1`,
		sid,
	).Scan(&jiraTicket)
	if err == sql.ErrNoRows {
		http.Error(w, "Shipment not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Failed to verify shipment", http.StatusInternalServerError)
		return
	}

	if jiraTicket == "" {
		http.Error(w, "Shipment must have a JIRA ticket", http.StatusBadRequest)
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
		// User doesn't exist, create a client user with a placeholder password
		// (they'll authenticate via magic link)
		placeholderPassword := "MAGIC_LINK_USER" // Placeholder to satisfy the auth_method constraint
		err = h.DB.QueryRowContext(
			r.Context(),
			`INSERT INTO users (email, password_hash, role, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id`,
			email, placeholderPassword, models.RoleClient, time.Now(), time.Now(),
		).Scan(&userID)
		if err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}
	} else if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set shipment ID (now required)
	shipmentID := &sid

	// Create magic link
	magicLink, err := auth.CreateMagicLink(r.Context(), h.DB, userID, shipmentID, auth.DefaultMagicLinkDuration)
	if err != nil {
		http.Error(w, "Failed to create magic link", http.StatusInternalServerError)
		return
	}

	// TODO: Send email with magic link
	// For now, show the magic link URL (in production, this would be sent via email)
	baseURL := getBaseURL(r)
	magicLinkURL := fmt.Sprintf("%s/auth/magic-link?token=%s", baseURL, magicLink.Token)

	// Redirect back to shipment detail with success message including the URL
	redirectURL := fmt.Sprintf("/shipments/%d?success=Magic+link+sent+to+%s.+URL:+%s", 
		sid, email, url.QueryEscape(magicLinkURL))
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

// getBaseURL returns the base URL for the application
func getBaseURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	host := r.Host
	if host == "" {
		host = "localhost:8080"
	}
	return fmt.Sprintf("%s://%s", scheme, host)
}

// MagicLinksList displays all magic links for logistics users
func (h *AuthHandler) MagicLinksList(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Only logistics users can view magic links
	if user.Role != models.RoleLogistics {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Query magic links with associated user and shipment info
	query := `
		SELECT 
			ml.id, ml.token, ml.user_id, ml.shipment_id, ml.expires_at, 
			ml.used_at, ml.created_at,
			u.email as user_email,
			COALESCE(s.id, 0) as shipment_id_coalesce,
			COALESCE(cc.name, '') as company_name
		FROM magic_links ml
		JOIN users u ON u.id = ml.user_id
		LEFT JOIN shipments s ON s.id = ml.shipment_id
		LEFT JOIN client_companies cc ON cc.id = s.client_company_id
		ORDER BY ml.created_at DESC
		LIMIT 100
	`

	rows, err := h.DB.QueryContext(r.Context(), query)
	if err != nil {
		http.Error(w, "Failed to load magic links", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type MagicLinkDisplay struct {
		ID          int64
		Token       string
		UserEmail   string
		ShipmentID  *int64
		CompanyName string
		ExpiresAt   time.Time
		UsedAt      *time.Time
		CreatedAt   time.Time
		IsExpired   bool
		IsUsed      bool
	}

	var links []MagicLinkDisplay
	for rows.Next() {
		var link MagicLinkDisplay
		var shipmentIDCoalesce int64
		var companyName string
		
		err := rows.Scan(
			&link.ID, &link.Token, new(int64), &link.ShipmentID, &link.ExpiresAt,
			&link.UsedAt, &link.CreatedAt, &link.UserEmail,
			&shipmentIDCoalesce, &companyName,
		)
		if err != nil {
			continue
		}

		if shipmentIDCoalesce > 0 {
			link.ShipmentID = &shipmentIDCoalesce
			link.CompanyName = companyName
		}

		link.IsExpired = time.Now().After(link.ExpiresAt)
		link.IsUsed = link.UsedAt != nil

		links = append(links, link)
	}

	data := map[string]interface{}{
		"User":        user,
		"Nav":         views.GetNavigationLinks(user.Role),
		"CurrentPage": "magic-links",
		"MagicLinks":  links,
	}

	err = h.Templates.ExecuteTemplate(w, "magic-links-list.html", data)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

// GoogleLogin initiates the Google OAuth flow
func (h *AuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	if h.OAuthConfig == nil {
		http.Error(w, "OAuth not configured", http.StatusInternalServerError)
		return
	}

	// Generate state token for CSRF protection
	state, err := auth.GenerateOAuthState()
	if err != nil {
		http.Error(w, "Failed to generate state token", http.StatusInternalServerError)
		return
	}

	// Store state in session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		MaxAge:   600, // 10 minutes
		HttpOnly: true,
		Secure:   isProduction(),
		SameSite: http.SameSiteStrictMode,
	})

	// Redirect to Google OAuth consent page
	url := h.OAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// GoogleCallback handles the OAuth callback from Google
func (h *AuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	if h.OAuthConfig == nil {
		http.Error(w, "OAuth not configured", http.StatusInternalServerError)
		return
	}

	// Verify state token to prevent CSRF
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil {
		http.Redirect(w, r, "/login?error=Invalid+OAuth+state", http.StatusSeeOther)
		return
	}

	state := r.URL.Query().Get("state")
	if state != stateCookie.Value {
		http.Redirect(w, r, "/login?error=Invalid+OAuth+state", http.StatusSeeOther)
		return
	}

	// Clear state cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   isProduction(),
		SameSite: http.SameSiteStrictMode,
	})

	// Exchange authorization code for token
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Redirect(w, r, "/login?error=No+authorization+code", http.StatusSeeOther)
		return
	}

	token, err := h.OAuthConfig.Exchange(r.Context(), code)
	if err != nil {
		http.Redirect(w, r, "/login?error=Failed+to+exchange+token", http.StatusSeeOther)
		return
	}

	// Get user info from Google
	userInfo, err := auth.GetGoogleUserInfo(r.Context(), token)
	if err != nil {
		http.Redirect(w, r, "/login?error=Failed+to+get+user+info", http.StatusSeeOther)
		return
	}

	// Validate email is verified
	if !userInfo.VerifiedEmail {
		http.Redirect(w, r, "/login?error=Email+not+verified", http.StatusSeeOther)
		return
	}

	// Validate domain if configured
	if h.OAuthDomain != "" && !auth.ValidateDomain(userInfo.Email, h.OAuthDomain) {
		http.Redirect(w, r, "/login?error=Email+domain+not+allowed", http.StatusSeeOther)
		return
	}

	// Find or create user
	user, err := auth.FindOrCreateGoogleUser(r.Context(), h.DB, userInfo, models.RoleLogistics)
	if err != nil {
		http.Error(w, "Failed to create/find user", http.StatusInternalServerError)
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
		Secure:   isProduction(),
		SameSite: http.SameSiteStrictMode,
	})

	// Redirect based on user role
	redirectURL := getRedirectURLForRole(user.Role)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}
