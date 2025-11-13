package auth

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/yourusername/laptop-tracking-system/internal/models"
)

// OAuthConfig holds the OAuth configuration
type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	AllowedDomain string // e.g., "bairesdev.com"
}

// GoogleUserInfo represents the user info from Google OAuth
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
	HD            string `json:"hd"` // Hosted domain (for G Suite)
}

// NewGoogleOAuthConfig creates a new Google OAuth2 configuration
func NewGoogleOAuthConfig(config OAuthConfig) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  config.RedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

// GenerateOAuthState generates a random state token for OAuth CSRF protection
func GenerateOAuthState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate OAuth state: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// GetGoogleUserInfo fetches user info from Google using the access token
func GetGoogleUserInfo(ctx context.Context, token *oauth2.Token) (*GoogleUserInfo, error) {
	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))
	
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read user info response: %w", err)
	}

	var userInfo GoogleUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user info: %w", err)
	}

	return &userInfo, nil
}

// FindOrCreateGoogleUser finds a user by Google ID or creates a new one
func FindOrCreateGoogleUser(ctx context.Context, db *sql.DB, userInfo *GoogleUserInfo, defaultRole models.UserRole) (*models.User, error) {
	// Try to find existing user by Google ID
	var user models.User
	err := db.QueryRowContext(
		ctx,
		`SELECT id, email, password_hash, role, client_company_id, google_id, created_at, updated_at
		FROM users
		WHERE google_id = $1`,
		userInfo.ID,
	).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.ClientCompanyID,
		&user.GoogleID, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == nil {
		// User found, update email if changed
		if user.Email != userInfo.Email {
			_, err := db.ExecContext(
				ctx,
				`UPDATE users SET email = $1, updated_at = $2 WHERE id = $3`,
				userInfo.Email, time.Now(), user.ID,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to update user email: %w", err)
			}
			user.Email = userInfo.Email
		}
		return &user, nil
	}

	if err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	// User not found, try to find by email
	err = db.QueryRowContext(
		ctx,
		`SELECT id, email, password_hash, role, client_company_id, google_id, created_at, updated_at
		FROM users
		WHERE email = $1`,
		userInfo.Email,
	).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.ClientCompanyID,
		&user.GoogleID, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == nil {
		// User found by email, link Google account
		googleID := userInfo.ID
		_, err := db.ExecContext(
			ctx,
			`UPDATE users SET google_id = $1, updated_at = $2 WHERE id = $3`,
			googleID, time.Now(), user.ID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to link Google account: %w", err)
		}
		user.GoogleID = &googleID
		return &user, nil
	}

	if err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to query user by email: %w", err)
	}

	// User doesn't exist, create new user
	user = models.User{
		Email:    userInfo.Email,
		Role:     defaultRole,
		GoogleID: &userInfo.ID,
	}
	user.BeforeCreate()

	err = db.QueryRowContext(
		ctx,
		`INSERT INTO users (email, role, google_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`,
		user.Email, user.Role, user.GoogleID, user.CreatedAt, user.UpdatedAt,
	).Scan(&user.ID)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

// ValidateDomain checks if the user's email domain matches the allowed domain
func ValidateDomain(email, allowedDomain string) bool {
	if allowedDomain == "" {
		return true // No domain restriction
	}

	// Extract domain from email
	at := -1
	for i := len(email) - 1; i >= 0; i-- {
		if email[i] == '@' {
			at = i
			break
		}
	}

	if at == -1 {
		return false
	}

	domain := email[at+1:]
	return domain == allowedDomain
}

