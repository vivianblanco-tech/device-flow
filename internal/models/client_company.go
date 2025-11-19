package models

import (
	"database/sql"
	"errors"
	"fmt"
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

// GetAllClientCompanies retrieves all client companies from the database
func GetAllClientCompanies(db interface{ Query(query string, args ...interface{}) (*sql.Rows, error) }) ([]ClientCompany, error) {
	query := `
		SELECT id, name, contact_info, created_at, updated_at
		FROM client_companies
		ORDER BY name ASC
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query client companies: %w", err)
	}
	defer rows.Close()

	var companies []ClientCompany
	for rows.Next() {
		var company ClientCompany
		err := rows.Scan(
			&company.ID,
			&company.Name,
			&company.ContactInfo,
			&company.CreatedAt,
			&company.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan client company: %w", err)
		}
		companies = append(companies, company)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating client companies: %w", err)
	}

	return companies, nil
}

// GetClientCompanyByID retrieves a client company by its ID
func GetClientCompanyByID(db *sql.DB, id int64) (*ClientCompany, error) {
	query := `
		SELECT id, name, contact_info, created_at, updated_at
		FROM client_companies
		WHERE id = $1
	`

	var company ClientCompany
	err := db.QueryRow(query, id).Scan(
		&company.ID,
		&company.Name,
		&company.ContactInfo,
		&company.CreatedAt,
		&company.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("client company not found with id %d", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get client company: %w", err)
	}

	return &company, nil
}

// CreateClientCompany creates a new client company in the database
func CreateClientCompany(db *sql.DB, company *ClientCompany) error {
	// Validate company
	if err := company.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Set timestamps
	company.BeforeCreate()

	query := `
		INSERT INTO client_companies (name, contact_info, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err := db.QueryRow(
		query,
		company.Name,
		company.ContactInfo,
		company.CreatedAt,
		company.UpdatedAt,
	).Scan(&company.ID)

	if err != nil {
		return fmt.Errorf("failed to create client company: %w", err)
	}

	return nil
}

// UpdateClientCompany updates an existing client company in the database
func UpdateClientCompany(db *sql.DB, company *ClientCompany) error {
	// Validate company
	if err := company.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Update timestamp
	company.BeforeUpdate()

	query := `
		UPDATE client_companies
		SET name = $1, contact_info = $2, updated_at = $3
		WHERE id = $4
	`

	result, err := db.Exec(
		query,
		company.Name,
		company.ContactInfo,
		company.UpdatedAt,
		company.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update client company: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("client company not found with id %d", company.ID)
	}

	return nil
}

// DeleteClientCompany deletes a client company from the database
func DeleteClientCompany(db *sql.DB, id int64) error {
	query := `DELETE FROM client_companies WHERE id = $1`

	result, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete client company: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("client company not found with id %d", id)
	}

	return nil
}

