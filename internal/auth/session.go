package auth

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/yourusername/laptop-tracking-system/internal/models"
)

const (
	// DefaultSessionDuration is the default session duration in hours
	DefaultSessionDuration = 24
	// SessionTokenLength is the length of the session token in bytes
	SessionTokenLength = 32
)

// GenerateSessionToken generates a cryptographically secure random token
func GenerateSessionToken() (string, error) {
	bytes := make([]byte, SessionTokenLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate session token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// CreateSession creates a new session for the user
func CreateSession(ctx context.Context, db *sql.DB, userID int64, durationHours int) (*models.Session, error) {
	// Verify user exists
	var exists bool
	err := db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", userID).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("user with ID %d does not exist", userID)
	}

	// Generate session token
	token, err := GenerateSessionToken()
	if err != nil {
		return nil, err
	}

	// Create session
	session := &models.Session{
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Duration(durationHours) * time.Hour),
	}
	session.BeforeCreate()

	// Insert session into database
	err = db.QueryRowContext(
		ctx,
		`INSERT INTO sessions (user_id, token, expires_at, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id`,
		session.UserID, session.Token, session.ExpiresAt, session.CreatedAt,
	).Scan(&session.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return session, nil
}

// ValidateSession validates a session token and returns the session with user info
// Returns nil if the session is invalid or expired
func ValidateSession(ctx context.Context, db *sql.DB, token string) (*models.Session, error) {
	if token == "" {
		return nil, nil
	}

	session := &models.Session{}
	user := &models.User{}

	// Query session with user join and company name (for client users)
	var companyName sql.NullString
	err := db.QueryRowContext(
		ctx,
		`SELECT 
			s.id, s.user_id, s.token, s.expires_at, s.created_at,
			u.id, u.email, u.password_hash, u.role, u.client_company_id, u.google_id, u.created_at, u.updated_at,
			cc.name as client_company_name
		FROM sessions s
		INNER JOIN users u ON s.user_id = u.id
		LEFT JOIN client_companies cc ON cc.id = u.client_company_id
		WHERE s.token = $1`,
		token,
	).Scan(
		&session.ID, &session.UserID, &session.Token, &session.ExpiresAt, &session.CreatedAt,
		&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.ClientCompanyID, &user.GoogleID, &user.CreatedAt, &user.UpdatedAt,
		&companyName,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to validate session: %w", err)
	}

	// Check if session is expired
	if session.IsExpired() {
		// Delete expired session
		_ = DeleteSession(ctx, db, token)
		return nil, nil
	}

	// Set company name if available
	if companyName.Valid {
		user.ClientCompanyName = companyName.String
	}

	session.User = user
	return session, nil
}

// DeleteSession deletes a session by token
func DeleteSession(ctx context.Context, db *sql.DB, token string) error {
	_, err := db.ExecContext(ctx, "DELETE FROM sessions WHERE token = $1", token)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}

// DeleteUserSessions deletes all sessions for a specific user
func DeleteUserSessions(ctx context.Context, db *sql.DB, userID int64) error {
	_, err := db.ExecContext(ctx, "DELETE FROM sessions WHERE user_id = $1", userID)
	if err != nil {
		return fmt.Errorf("failed to delete user sessions: %w", err)
	}
	return nil
}

// CleanupExpiredSessions removes all expired sessions from the database
// Returns the number of sessions deleted
func CleanupExpiredSessions(ctx context.Context, db *sql.DB) (int, error) {
	result, err := db.ExecContext(
		ctx,
		"DELETE FROM sessions WHERE expires_at < $1",
		time.Now(),
	)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup expired sessions: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return int(rowsAffected), nil
}

// ExtendSession extends the expiration time of a session
func ExtendSession(ctx context.Context, db *sql.DB, token string, durationHours int) error {
	newExpiresAt := time.Now().Add(time.Duration(durationHours) * time.Hour)
	
	result, err := db.ExecContext(
		ctx,
		"UPDATE sessions SET expires_at = $1 WHERE token = $2",
		newExpiresAt, token,
	)
	if err != nil {
		return fmt.Errorf("failed to extend session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("session not found")
	}

	return nil
}

