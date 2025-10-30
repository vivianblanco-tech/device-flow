package models

import (
	"errors"
	"time"
)

// SoftwareEngineer represents a software engineer who will receive a laptop
type SoftwareEngineer struct {
	ID                    int64      `json:"id" db:"id"`
	Name                  string     `json:"name" db:"name"`
	Email                 string     `json:"email" db:"email"`
	Address               string     `json:"address,omitempty" db:"address"`
	Phone                 string     `json:"phone,omitempty" db:"phone"`
	AddressConfirmed      bool       `json:"address_confirmed" db:"address_confirmed"`
	AddressConfirmationAt *time.Time `json:"address_confirmation_at,omitempty" db:"address_confirmation_at"`
	CreatedAt             time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at" db:"updated_at"`
}

// Validate validates the SoftwareEngineer model
func (s *SoftwareEngineer) Validate() error {
	// Name validation
	if s.Name == "" {
		return errors.New("engineer name is required")
	}

	// Email validation
	if s.Email == "" {
		return errors.New("engineer email is required")
	}
	if !emailRegex.MatchString(s.Email) {
		return errors.New("invalid email format")
	}

	return nil
}

// TableName returns the table name for the SoftwareEngineer model
func (s *SoftwareEngineer) TableName() string {
	return "software_engineers"
}

// BeforeCreate sets the timestamps before creating a software engineer
func (s *SoftwareEngineer) BeforeCreate() {
	now := time.Now()
	s.CreatedAt = now
	s.UpdatedAt = now
}

// BeforeUpdate sets the updated_at timestamp before updating a software engineer
func (s *SoftwareEngineer) BeforeUpdate() {
	s.UpdatedAt = time.Now()
}

// HasConfirmedAddress returns true if the engineer has confirmed their address
func (s *SoftwareEngineer) HasConfirmedAddress() bool {
	return s.AddressConfirmed
}

// ConfirmAddress marks the engineer's address as confirmed
func (s *SoftwareEngineer) ConfirmAddress() {
	now := time.Now()
	s.AddressConfirmed = true
	s.AddressConfirmationAt = &now
}

