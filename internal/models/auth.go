package models

import (
	"errors"
	"time"
)

// MagicLink represents a one-time login link sent via email
type MagicLink struct {
	ID         int64      `json:"id" db:"id"`
	UserID     int64      `json:"user_id" db:"user_id"`
	Token      string     `json:"token" db:"token"`
	ExpiresAt  time.Time  `json:"expires_at" db:"expires_at"`
	UsedAt     *time.Time `json:"used_at,omitempty" db:"used_at"`
	ShipmentID *int64     `json:"shipment_id,omitempty" db:"shipment_id"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`

	// Relations
	User     *User     `json:"user,omitempty" db:"-"`
	Shipment *Shipment `json:"shipment,omitempty" db:"-"`
}

// Validate validates the MagicLink model
func (m *MagicLink) Validate() error {
	if m.UserID == 0 {
		return errors.New("user ID is required")
	}
	if m.Token == "" {
		return errors.New("token is required")
	}
	if m.ExpiresAt.IsZero() {
		return errors.New("expiration time is required")
	}
	return nil
}

// TableName returns the table name for the MagicLink model
func (m *MagicLink) TableName() string {
	return "magic_links"
}

// BeforeCreate sets the timestamp before creating a magic link
func (m *MagicLink) BeforeCreate() {
	m.CreatedAt = time.Now()
}

// IsExpired returns true if the magic link has expired
func (m *MagicLink) IsExpired() bool {
	return time.Now().After(m.ExpiresAt)
}

// IsUsed returns true if the magic link has been used
func (m *MagicLink) IsUsed() bool {
	return m.UsedAt != nil
}

// MarkAsUsed marks the magic link as used
func (m *MagicLink) MarkAsUsed() {
	now := time.Now()
	m.UsedAt = &now
}

// Session represents a user session
type Session struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	Token     string    `json:"token" db:"token"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`

	// Relations
	User *User `json:"user,omitempty" db:"-"`
}

// Validate validates the Session model
func (s *Session) Validate() error {
	if s.UserID == 0 {
		return errors.New("user ID is required")
	}
	if s.Token == "" {
		return errors.New("token is required")
	}
	if s.ExpiresAt.IsZero() {
		return errors.New("expiration time is required")
	}
	return nil
}

// TableName returns the table name for the Session model
func (s *Session) TableName() string {
	return "sessions"
}

// BeforeCreate sets the timestamp before creating a session
func (s *Session) BeforeCreate() {
	s.CreatedAt = time.Now()
}

// IsExpired returns true if the session has expired
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

