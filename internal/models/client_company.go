package models

import (
	"errors"
	"time"
)

// ClientCompany represents a client company in the system
type ClientCompany struct {
	ID          int64     `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	ContactInfo string    `json:"contact_info,omitempty" db:"contact_info"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`

	// Relations (not stored in DB directly)
	Users []User `json:"users,omitempty" db:"-"`
}

// Validate validates the ClientCompany model
func (c *ClientCompany) Validate() error {
	// Name validation
	if c.Name == "" {
		return errors.New("company name is required")
	}
	if len(c.Name) < 3 {
		return errors.New("company name must be at least 3 characters")
	}

	return nil
}

// TableName returns the table name for the ClientCompany model
func (c *ClientCompany) TableName() string {
	return "client_companies"
}

// BeforeCreate sets the timestamps before creating a client company
func (c *ClientCompany) BeforeCreate() {
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now
}

// BeforeUpdate sets the updated_at timestamp before updating a client company
func (c *ClientCompany) BeforeUpdate() {
	c.UpdatedAt = time.Now()
}

// GetActiveUsersCount returns the count of active users for this company
// This is a placeholder that will be properly implemented with database queries
func (c *ClientCompany) GetActiveUsersCount() int {
	return len(c.Users)
}

