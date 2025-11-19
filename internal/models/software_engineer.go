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
	EmployeeNumber        string     `json:"employee_number,omitempty" db:"employee_number"`
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
		SELECT id, name, email, phone, address, employee_number, address_confirmed, address_confirmation_at, created_at, updated_at
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
		var employeeNumber sql.NullString
		var addressConfirmationAt sql.NullTime

		err := rows.Scan(
			&engineer.ID,
			&engineer.Name,
			&engineer.Email,
			&phone,
			&address,
			&employeeNumber,
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
		if employeeNumber.Valid {
			engineer.EmployeeNumber = employeeNumber.String
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

// GetSoftwareEngineerByID retrieves a software engineer by its ID
func GetSoftwareEngineerByID(db *sql.DB, id int64) (*SoftwareEngineer, error) {
	query := `
		SELECT id, name, email, phone, address, employee_number, address_confirmed, address_confirmation_at, created_at, updated_at
		FROM software_engineers
		WHERE id = $1
	`

	var engineer SoftwareEngineer
	var phone sql.NullString
	var address sql.NullString
	var employeeNumber sql.NullString
	var addressConfirmationAt sql.NullTime

	err := db.QueryRow(query, id).Scan(
		&engineer.ID,
		&engineer.Name,
		&engineer.Email,
		&phone,
		&address,
		&employeeNumber,
		&engineer.AddressConfirmed,
		&addressConfirmationAt,
		&engineer.CreatedAt,
		&engineer.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("software engineer not found with id %d", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get software engineer: %w", err)
	}

	// Set nullable fields if available
	if phone.Valid {
		engineer.Phone = phone.String
	}
	if address.Valid {
		engineer.Address = address.String
	}
	if employeeNumber.Valid {
		engineer.EmployeeNumber = employeeNumber.String
	}
	if addressConfirmationAt.Valid {
		engineer.AddressConfirmationAt = &addressConfirmationAt.Time
	}

	return &engineer, nil
}

// CreateSoftwareEngineer creates a new software engineer in the database
func CreateSoftwareEngineer(db *sql.DB, engineer *SoftwareEngineer) error {
	// Validate engineer
	if err := engineer.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Set timestamps
	engineer.BeforeCreate()

	query := `
		INSERT INTO software_engineers (name, email, address, phone, employee_number, address_confirmed, address_confirmation_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`

	err := db.QueryRow(
		query,
		engineer.Name,
		engineer.Email,
		engineer.Address,
		engineer.Phone,
		engineer.EmployeeNumber,
		engineer.AddressConfirmed,
		engineer.AddressConfirmationAt,
		engineer.CreatedAt,
		engineer.UpdatedAt,
	).Scan(&engineer.ID)

	if err != nil {
		return fmt.Errorf("failed to create software engineer: %w", err)
	}

	return nil
}

// UpdateSoftwareEngineer updates an existing software engineer in the database
func UpdateSoftwareEngineer(db *sql.DB, engineer *SoftwareEngineer) error {
	// Validate engineer
	if err := engineer.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Update timestamp
	engineer.BeforeUpdate()

	query := `
		UPDATE software_engineers
		SET name = $1, email = $2, address = $3, phone = $4, employee_number = $5, address_confirmed = $6, address_confirmation_at = $7, updated_at = $8
		WHERE id = $9
	`

	result, err := db.Exec(
		query,
		engineer.Name,
		engineer.Email,
		engineer.Address,
		engineer.Phone,
		engineer.EmployeeNumber,
		engineer.AddressConfirmed,
		engineer.AddressConfirmationAt,
		engineer.UpdatedAt,
		engineer.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update software engineer: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("software engineer not found with id %d", engineer.ID)
	}

	return nil
}

// DeleteSoftwareEngineer deletes a software engineer from the database
func DeleteSoftwareEngineer(db *sql.DB, id int64) error {
	query := `DELETE FROM software_engineers WHERE id = $1`

	result, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete software engineer: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("software engineer not found with id %d", id)
	}

	return nil
}

