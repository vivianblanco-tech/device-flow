package models

import (
	"errors"
	"regexp"
	"time"
)

// UserRole represents the role of a user in the system
type UserRole string

// User role constants
const (
	RoleLogistics      UserRole = "logistics"
	RoleClient         UserRole = "client"
	RoleWarehouse      UserRole = "warehouse"
	RoleProjectManager UserRole = "project_manager"
)

// User represents a user in the system
type User struct {
	ID              int64     `json:"id" db:"id"`
	Email           string    `json:"email" db:"email"`
	PasswordHash    string    `json:"-" db:"password_hash"`
	Role            UserRole  `json:"role" db:"role"`
	ClientCompanyID *int64    `json:"client_company_id,omitempty" db:"client_company_id"`
	GoogleID        *string   `json:"google_id,omitempty" db:"google_id"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// Email validation regex
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// Validate validates the User model
func (u *User) Validate() error {
	// Email validation
	if u.Email == "" {
		return errors.New("email is required")
	}
	if !emailRegex.MatchString(u.Email) {
		return errors.New("invalid email format")
	}

	// Role validation
	if u.Role == "" {
		return errors.New("role is required")
	}
	if !IsValidRole(u.Role) {
		return errors.New("invalid role")
	}

	// Either password_hash or google_id must be provided
	if u.PasswordHash == "" && u.GoogleID == nil {
		return errors.New("either password_hash or google_id must be provided")
	}

	return nil
}

// IsValidRole checks if a given role is valid
func IsValidRole(role UserRole) bool {
	switch role {
	case RoleLogistics, RoleClient, RoleWarehouse, RoleProjectManager:
		return true
	}
	return false
}

// HasRole checks if the user has the specified role
func (u *User) HasRole(role UserRole) bool {
	return u.Role == role
}

// IsGoogleUser checks if the user authenticated via Google OAuth
func (u *User) IsGoogleUser() bool {
	return u.GoogleID != nil && *u.GoogleID != ""
}

// TableName returns the table name for the User model
func (u *User) TableName() string {
	return "users"
}

// BeforeCreate sets the timestamps before creating a user
func (u *User) BeforeCreate() {
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now
}

// BeforeUpdate sets the updated_at timestamp before updating a user
func (u *User) BeforeUpdate() {
	u.UpdatedAt = time.Now()
}

