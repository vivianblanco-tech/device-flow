package auth

import (
	"context"
	"testing"
	"time"

	"github.com/yourusername/laptop-tracking-system/internal/database"
	"github.com/yourusername/laptop-tracking-system/internal/models"
)

func TestGenerateSessionToken(t *testing.T) {
	token1, err := GenerateSessionToken()
	if err != nil {
		t.Fatalf("GenerateSessionToken() failed: %v", err)
	}

	if token1 == "" {
		t.Error("GenerateSessionToken() returned empty token")
	}

	if len(token1) < 32 {
		t.Errorf("GenerateSessionToken() token length = %d, expected at least 32", len(token1))
	}

	// Generate another token to ensure uniqueness
	token2, err := GenerateSessionToken()
	if err != nil {
		t.Fatalf("GenerateSessionToken() failed on second call: %v", err)
	}

	if token1 == token2 {
		t.Error("GenerateSessionToken() generated duplicate tokens")
	}
}

func TestCreateSession(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup test database
	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	// Create a test user
	user := &models.User{
		Email:        "session@test.com",
		PasswordHash: "hashedpassword123",
		Role:         models.RoleLogistics,
	}
	user.BeforeCreate()

	err := db.QueryRowContext(
		context.Background(),
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		user.Email, user.PasswordHash, user.Role, user.CreatedAt, user.UpdatedAt,
	).Scan(&user.ID)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	tests := []struct {
		name           string
		userID         int64
		durationHours  int
		wantErr        bool
		checkExpiration bool
	}{
		{
			name:          "create session with standard duration",
			userID:        user.ID,
			durationHours: 24,
			wantErr:       false,
			checkExpiration: true,
		},
		{
			name:          "create session with short duration",
			userID:        user.ID,
			durationHours: 1,
			wantErr:       false,
			checkExpiration: true,
		},
		{
			name:          "create session with long duration",
			userID:        user.ID,
			durationHours: 168, // 7 days
			wantErr:       false,
			checkExpiration: true,
		},
		{
			name:          "create session with invalid user ID",
			userID:        99999,
			durationHours: 24,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session, err := CreateSession(context.Background(), db, tt.userID, tt.durationHours)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if session.ID == 0 {
				t.Error("CreateSession() returned session with zero ID")
			}

			if session.UserID != tt.userID {
				t.Errorf("CreateSession() UserID = %d, want %d", session.UserID, tt.userID)
			}

			if session.Token == "" {
				t.Error("CreateSession() returned session with empty token")
			}

			if tt.checkExpiration {
				expectedExpiry := time.Now().Add(time.Duration(tt.durationHours) * time.Hour)
				timeDiff := session.ExpiresAt.Sub(expectedExpiry).Abs()
				// Allow 1 second tolerance for test execution time
				if timeDiff > time.Second {
					t.Errorf("CreateSession() ExpiresAt = %v, expected around %v", session.ExpiresAt, expectedExpiry)
				}
			}

			if session.CreatedAt.IsZero() {
				t.Error("CreateSession() returned session with zero CreatedAt")
			}
		})
	}
}

func TestValidateSession(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup test database
	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	// Create a test user
	user := &models.User{
		Email:        "validate@test.com",
		PasswordHash: "hashedpassword123",
		Role:         models.RoleLogistics,
	}
	user.BeforeCreate()

	err := db.QueryRowContext(
		context.Background(),
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		user.Email, user.PasswordHash, user.Role, user.CreatedAt, user.UpdatedAt,
	).Scan(&user.ID)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Create a valid session
	validSession, err := CreateSession(context.Background(), db, user.ID, 24)
	if err != nil {
		t.Fatalf("Failed to create test session: %v", err)
	}

	// Create an expired session
	expiredSession := &models.Session{
		UserID:    user.ID,
		Token:     "expired-token-12345",
		ExpiresAt: time.Now().Add(-1 * time.Hour), // Already expired
	}
	expiredSession.BeforeCreate()

	err = db.QueryRowContext(
		context.Background(),
		`INSERT INTO sessions (user_id, token, expires_at, created_at)
		VALUES ($1, $2, $3, $4) RETURNING id`,
		expiredSession.UserID, expiredSession.Token, expiredSession.ExpiresAt, expiredSession.CreatedAt,
	).Scan(&expiredSession.ID)
	if err != nil {
		t.Fatalf("Failed to create expired session: %v", err)
	}

	tests := []struct {
		name    string
		token   string
		wantErr bool
		wantNil bool
	}{
		{
			name:    "valid session token",
			token:   validSession.Token,
			wantErr: false,
			wantNil: false,
		},
		{
			name:    "expired session token",
			token:   expiredSession.Token,
			wantErr: false,
			wantNil: true, // Should return nil for expired session
		},
		{
			name:    "non-existent token",
			token:   "non-existent-token",
			wantErr: false,
			wantNil: true,
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: false,
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session, err := ValidateSession(context.Background(), db, tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantNil {
				if session != nil {
					t.Errorf("ValidateSession() returned session, expected nil")
				}
			} else {
				if session == nil {
					t.Error("ValidateSession() returned nil, expected session")
					return
				}

				if session.Token != tt.token {
					t.Errorf("ValidateSession() Token = %s, want %s", session.Token, tt.token)
				}

				if session.IsExpired() {
					t.Error("ValidateSession() returned expired session")
				}

				// Verify user is loaded
				if session.User == nil {
					t.Error("ValidateSession() did not load user relation")
				}
			}
		})
	}
}

func TestCleanupExpiredSessions(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup test database
	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	// Create a test user
	user := &models.User{
		Email:        "cleanup@test.com",
		PasswordHash: "hashedpassword123",
		Role:         models.RoleLogistics,
	}
	user.BeforeCreate()

	err := db.QueryRowContext(
		context.Background(),
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		user.Email, user.PasswordHash, user.Role, user.CreatedAt, user.UpdatedAt,
	).Scan(&user.ID)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Use a fixed reference time to avoid race conditions between session creation and cleanup
	now := time.Now()

	// Create multiple sessions with different expiration times
	sessions := []struct {
		token     string
		expiresAt time.Time
		shouldDelete bool
	}{
		{
			token:     "valid-session-1",
			expiresAt: now.Add(24 * time.Hour),
			shouldDelete: false,
		},
		{
			token:     "expired-session-1",
			expiresAt: now.Add(-1 * time.Hour),
			shouldDelete: true,
		},
		{
			token:     "expired-session-2",
			expiresAt: now.Add(-24 * time.Hour),
			shouldDelete: true,
		},
		{
			token:     "valid-session-2",
			expiresAt: now.Add(48 * time.Hour),
			shouldDelete: false,
		},
	}

	for _, s := range sessions {
		session := &models.Session{
			UserID:    user.ID,
			Token:     s.token,
			ExpiresAt: s.expiresAt,
		}
		session.BeforeCreate()

		err = db.QueryRowContext(
			context.Background(),
			`INSERT INTO sessions (user_id, token, expires_at, created_at)
			VALUES ($1, $2, $3, $4) RETURNING id`,
			session.UserID, session.Token, session.ExpiresAt, session.CreatedAt,
		).Scan(&session.ID)
		if err != nil {
			t.Fatalf("Failed to create test session: %v", err)
		}
	}

	// Clean up expired sessions
	deletedCount, err := CleanupExpiredSessions(context.Background(), db)
	if err != nil {
		t.Fatalf("CleanupExpiredSessions() failed: %v", err)
	}

	// Count expected deletions
	expectedDeletions := 0
	for _, s := range sessions {
		if s.shouldDelete {
			expectedDeletions++
		}
	}

	if deletedCount != expectedDeletions {
		t.Errorf("CleanupExpiredSessions() deleted %d sessions, expected %d", deletedCount, expectedDeletions)
	}

	// Verify that valid sessions still exist
	for _, s := range sessions {
		var count int
		err := db.QueryRowContext(
			context.Background(),
			`SELECT COUNT(*) FROM sessions WHERE token = $1`,
			s.token,
		).Scan(&count)
		if err != nil {
			t.Fatalf("Failed to check session existence: %v", err)
		}

		if s.shouldDelete && count > 0 {
			t.Errorf("Session %s should have been deleted but still exists", s.token)
		}

		if !s.shouldDelete && count == 0 {
			t.Errorf("Session %s should not have been deleted but is missing", s.token)
		}
	}
}

func TestDeleteSession(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup test database
	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	// Create a test user
	user := &models.User{
		Email:        "delete@test.com",
		PasswordHash: "hashedpassword123",
		Role:         models.RoleLogistics,
	}
	user.BeforeCreate()

	err := db.QueryRowContext(
		context.Background(),
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		user.Email, user.PasswordHash, user.Role, user.CreatedAt, user.UpdatedAt,
	).Scan(&user.ID)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Create a session to delete
	session, err := CreateSession(context.Background(), db, user.ID, 24)
	if err != nil {
		t.Fatalf("Failed to create test session: %v", err)
	}

	// Delete the session
	err = DeleteSession(context.Background(), db, session.Token)
	if err != nil {
		t.Errorf("DeleteSession() failed: %v", err)
	}

	// Verify session is deleted
	var count int
	err = db.QueryRowContext(
		context.Background(),
		`SELECT COUNT(*) FROM sessions WHERE token = $1`,
		session.Token,
	).Scan(&count)
	if err != nil {
		t.Fatalf("Failed to check session deletion: %v", err)
	}

	if count > 0 {
		t.Error("DeleteSession() did not delete the session")
	}

	// Try deleting non-existent session (should not error)
	err = DeleteSession(context.Background(), db, "non-existent-token")
	if err != nil {
		t.Errorf("DeleteSession() with non-existent token should not error: %v", err)
	}
}

