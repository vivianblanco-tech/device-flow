package models

import (
	"database/sql"
	"errors"
	"fmt"
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

// GetAllSoftwareEngineers retrieves all software engineers from the database sorted by name
func GetAllSoftwareEngineers(db interface{ Query(query string, args ...interface{}) (*sql.Rows, error) }) ([]SoftwareEngineer, error) {
	query := `
		SELECT id, name, email, phone, address, address_confirmed, address_confirmation_at, created_at, updated_at
		FROM software_engineers
		ORDER BY name ASC
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query software engineers: %w", err)
	}
	defer rows.Close()

	var engineers []SoftwareEngineer
	for rows.Next() {
		var engineer SoftwareEngineer
		var phone sql.NullString
		var address sql.NullString
		var addressConfirmationAt sql.NullTime

		err := rows.Scan(
			&engineer.ID,
			&engineer.Name,
			&engineer.Email,
			&phone,
			&address,
			&engineer.AddressConfirmed,
			&addressConfirmationAt,
			&engineer.CreatedAt,
			&engineer.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan software engineer: %w", err)
		}

		// Set nullable fields if available
		if phone.Valid {
			engineer.Phone = phone.String
		}
		if address.Valid {
			engineer.Address = address.String
		}
		if addressConfirmationAt.Valid {
			engineer.AddressConfirmationAt = &addressConfirmationAt.Time
		}

		engineers = append(engineers, engineer)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating software engineers: %w", err)
	}

	return engineers, nil
}

