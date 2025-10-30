package models

import (
	"testing"
	"time"
)

// MagicLink tests

func TestMagicLink_Validate(t *testing.T) {
	tests := []struct {
		name    string
		link    MagicLink
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid magic link",
			link: MagicLink{
				UserID:     1,
				Token:      "secure_random_token_123",
				ExpiresAt:  time.Now().Add(24 * time.Hour),
				ShipmentID: int64Ptr(10),
			},
			wantErr: false,
		},
		{
			name: "valid - no shipment ID",
			link: MagicLink{
				UserID:    1,
				Token:     "secure_random_token_456",
				ExpiresAt: time.Now().Add(1 * time.Hour),
			},
			wantErr: false,
		},
		{
			name: "invalid - missing user ID",
			link: MagicLink{
				Token:     "secure_random_token_789",
				ExpiresAt: time.Now().Add(24 * time.Hour),
			},
			wantErr: true,
			errMsg:  "user ID is required",
		},
		{
			name: "invalid - missing token",
			link: MagicLink{
				UserID:    1,
				ExpiresAt: time.Now().Add(24 * time.Hour),
			},
			wantErr: true,
			errMsg:  "token is required",
		},
		{
			name: "invalid - missing expiration",
			link: MagicLink{
				UserID: 1,
				Token:  "secure_random_token_xyz",
			},
			wantErr: true,
			errMsg:  "expiration time is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.link.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("MagicLink.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" && err.Error() != tt.errMsg {
				t.Errorf("MagicLink.Validate() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestMagicLink_TableName(t *testing.T) {
	link := MagicLink{}
	expected := "magic_links"
	if got := link.TableName(); got != expected {
		t.Errorf("MagicLink.TableName() = %v, want %v", got, expected)
	}
}

func TestMagicLink_IsExpired(t *testing.T) {
	tests := []struct {
		name     string
		link     MagicLink
		expected bool
	}{
		{
			name: "not expired",
			link: MagicLink{
				ExpiresAt: time.Now().Add(1 * time.Hour),
			},
			expected: false,
		},
		{
			name: "expired",
			link: MagicLink{
				ExpiresAt: time.Now().Add(-1 * time.Hour),
			},
			expected: true,
		},
		{
			name: "just now (edge case)",
			link: MagicLink{
				ExpiresAt: time.Now(),
			},
			expected: false, // Could be either, depends on exact timing
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.link.IsExpired()
			// For "just now" case, we allow both results
			if tt.name != "just now (edge case)" && got != tt.expected {
				t.Errorf("MagicLink.IsExpired() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestMagicLink_IsUsed(t *testing.T) {
	tests := []struct {
		name     string
		link     MagicLink
		expected bool
	}{
		{
			name: "used",
			link: MagicLink{
				UsedAt: timePtr(time.Now()),
			},
			expected: true,
		},
		{
			name: "not used",
			link: MagicLink{
				UsedAt: nil,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.link.IsUsed(); got != tt.expected {
				t.Errorf("MagicLink.IsUsed() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestMagicLink_MarkAsUsed(t *testing.T) {
	link := &MagicLink{
		UserID:    1,
		Token:     "test_token",
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	if link.UsedAt != nil {
		t.Error("Expected UsedAt to be nil initially")
	}

	link.MarkAsUsed()

	if link.UsedAt == nil {
		t.Error("MarkAsUsed() did not set UsedAt")
	}
	if link.UsedAt.IsZero() {
		t.Error("MarkAsUsed() set UsedAt to zero time")
	}
}

// Session tests

func TestSession_Validate(t *testing.T) {
	tests := []struct {
		name    string
		session Session
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid session",
			session: Session{
				UserID:    1,
				Token:     "session_token_abc123",
				ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
			},
			wantErr: false,
		},
		{
			name: "invalid - missing user ID",
			session: Session{
				Token:     "session_token_def456",
				ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
			},
			wantErr: true,
			errMsg:  "user ID is required",
		},
		{
			name: "invalid - missing token",
			session: Session{
				UserID:    1,
				ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
			},
			wantErr: true,
			errMsg:  "token is required",
		},
		{
			name: "invalid - missing expiration",
			session: Session{
				UserID: 1,
				Token:  "session_token_ghi789",
			},
			wantErr: true,
			errMsg:  "expiration time is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.session.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Session.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" && err.Error() != tt.errMsg {
				t.Errorf("Session.Validate() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestSession_TableName(t *testing.T) {
	session := Session{}
	expected := "sessions"
	if got := session.TableName(); got != expected {
		t.Errorf("Session.TableName() = %v, want %v", got, expected)
	}
}

func TestSession_IsExpired(t *testing.T) {
	tests := []struct {
		name     string
		session  Session
		expected bool
	}{
		{
			name: "not expired",
			session: Session{
				ExpiresAt: time.Now().Add(1 * time.Hour),
			},
			expected: false,
		},
		{
			name: "expired",
			session: Session{
				ExpiresAt: time.Now().Add(-1 * time.Hour),
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.session.IsExpired(); got != tt.expected {
				t.Errorf("Session.IsExpired() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSession_BeforeCreate(t *testing.T) {
	session := &Session{
		UserID:    1,
		Token:     "test_token",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	session.BeforeCreate()

	if session.CreatedAt.IsZero() {
		t.Error("Session.BeforeCreate() did not set CreatedAt")
	}
}

