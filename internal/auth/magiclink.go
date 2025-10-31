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
	// MagicLinkTokenLength is the length of the magic link token in bytes
	MagicLinkTokenLength = 32
	// DefaultMagicLinkDuration is the default magic link expiration duration in hours
	DefaultMagicLinkDuration = 24
)

// GenerateMagicLinkToken generates a cryptographically secure random token
func GenerateMagicLinkToken() (string, error) {
	bytes := make([]byte, MagicLinkTokenLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate magic link token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// CreateMagicLink creates a new magic link for the user
func CreateMagicLink(ctx context.Context, db *sql.DB, userID int64, shipmentID *int64, durationHours int) (*models.MagicLink, error) {
	// Verify user exists
	var exists bool
	err := db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", userID).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("user with ID %d does not exist", userID)
	}

	// If shipmentID is provided, verify it exists
	if shipmentID != nil {
		err := db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM shipments WHERE id = $1)", *shipmentID).Scan(&exists)
		if err != nil {
			return nil, fmt.Errorf("failed to check shipment existence: %w", err)
		}
		if !exists {
			return nil, fmt.Errorf("shipment with ID %d does not exist", *shipmentID)
		}
	}

	// Generate token
	token, err := GenerateMagicLinkToken()
	if err != nil {
		return nil, err
	}

	// Create magic link
	magicLink := &models.MagicLink{
		UserID:     userID,
		Token:      token,
		ExpiresAt:  time.Now().Add(time.Duration(durationHours) * time.Hour),
		ShipmentID: shipmentID,
	}
	magicLink.BeforeCreate()

	// Insert magic link into database
	err = db.QueryRowContext(
		ctx,
		`INSERT INTO magic_links (user_id, token, expires_at, shipment_id, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`,
		magicLink.UserID, magicLink.Token, magicLink.ExpiresAt, magicLink.ShipmentID, magicLink.CreatedAt,
	).Scan(&magicLink.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create magic link: %w", err)
	}

	return magicLink, nil
}

// ValidateMagicLink validates a magic link token and returns the magic link with user info
// Returns nil if the magic link is invalid, expired, or already used
func ValidateMagicLink(ctx context.Context, db *sql.DB, token string) (*models.MagicLink, error) {
	if token == "" {
		return nil, nil
	}

	magicLink := &models.MagicLink{}
	user := &models.User{}

	var usedAt sql.NullTime
	var shipmentID sql.NullInt64
	var googleID sql.NullString

	// Query magic link with user join
	err := db.QueryRowContext(
		ctx,
		`SELECT 
			ml.id, ml.user_id, ml.token, ml.expires_at, ml.used_at, ml.shipment_id, ml.created_at,
			u.id, u.email, u.password_hash, u.role, u.google_id, u.created_at, u.updated_at
		FROM magic_links ml
		INNER JOIN users u ON ml.user_id = u.id
		WHERE ml.token = $1`,
		token,
	).Scan(
		&magicLink.ID, &magicLink.UserID, &magicLink.Token, &magicLink.ExpiresAt,
		&usedAt, &shipmentID, &magicLink.CreatedAt,
		&user.ID, &user.Email, &user.PasswordHash, &user.Role, &googleID,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to validate magic link: %w", err)
	}

	// Handle nullable fields
	if usedAt.Valid {
		magicLink.UsedAt = &usedAt.Time
	}
	if shipmentID.Valid {
		magicLink.ShipmentID = &shipmentID.Int64
	}
	if googleID.Valid {
		user.GoogleID = &googleID.String
	}

	// Check if magic link is expired
	if magicLink.IsExpired() {
		return nil, nil
	}

	// Check if magic link has been used
	if magicLink.IsUsed() {
		return nil, nil
	}

	magicLink.User = user
	return magicLink, nil
}

// MarkMagicLinkAsUsed marks a magic link as used
func MarkMagicLinkAsUsed(ctx context.Context, db *sql.DB, token string) error {
	now := time.Now()
	result, err := db.ExecContext(
		ctx,
		"UPDATE magic_links SET used_at = $1 WHERE token = $2 AND used_at IS NULL",
		now, token,
	)
	if err != nil {
		return fmt.Errorf("failed to mark magic link as used: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("magic link not found or already used")
	}

	return nil
}

// DeleteMagicLink deletes a magic link by token
func DeleteMagicLink(ctx context.Context, db *sql.DB, token string) error {
	_, err := db.ExecContext(ctx, "DELETE FROM magic_links WHERE token = $1", token)
	if err != nil {
		return fmt.Errorf("failed to delete magic link: %w", err)
	}
	return nil
}

// CleanupExpiredMagicLinks removes all expired and used magic links from the database
// Returns the number of magic links deleted
func CleanupExpiredMagicLinks(ctx context.Context, db *sql.DB) (int, error) {
	result, err := db.ExecContext(
		ctx,
		"DELETE FROM magic_links WHERE expires_at < $1 OR used_at IS NOT NULL",
		time.Now(),
	)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup expired magic links: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return int(rowsAffected), nil
}

// GetMagicLinksByUser retrieves all valid magic links for a user
func GetMagicLinksByUser(ctx context.Context, db *sql.DB, userID int64) ([]*models.MagicLink, error) {
	rows, err := db.QueryContext(
		ctx,
		`SELECT id, user_id, token, expires_at, used_at, shipment_id, created_at
		FROM magic_links
		WHERE user_id = $1 AND expires_at > $2
		ORDER BY created_at DESC`,
		userID, time.Now(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get magic links: %w", err)
	}
	defer rows.Close()

	var magicLinks []*models.MagicLink
	for rows.Next() {
		ml := &models.MagicLink{}
		var usedAt sql.NullTime
		var shipmentID sql.NullInt64

		err := rows.Scan(
			&ml.ID, &ml.UserID, &ml.Token, &ml.ExpiresAt,
			&usedAt, &shipmentID, &ml.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan magic link: %w", err)
		}

		if usedAt.Valid {
			ml.UsedAt = &usedAt.Time
		}
		if shipmentID.Valid {
			ml.ShipmentID = &shipmentID.Int64
		}

		magicLinks = append(magicLinks, ml)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating magic links: %w", err)
	}

	return magicLinks, nil
}

