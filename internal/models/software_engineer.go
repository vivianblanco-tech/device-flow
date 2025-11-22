package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

// SoftwareEngineer represents a software engineer who will receive a laptop
type SoftwareEngineer struct {
	ID                    int64      `json:"id" db:"id"`
	Name                  string     `json:"name" db:"name"`
	Email                 string     `json:"email" db:"email"`
	Address               string     `json:"address,omitempty" db:"address"` // Legacy field, kept for backward compatibility
	AddressStreet         string     `json:"address_street,omitempty" db:"address_street"`
	AddressCity           string     `json:"address_city,omitempty" db:"address_city"`
	AddressCountry        string     `json:"address_country,omitempty" db:"address_country"`
	AddressState          string     `json:"address_state,omitempty" db:"address_state"`
	AddressPostalCode     string     `json:"address_postal_code,omitempty" db:"address_postal_code"`
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

// SoftwareEngineerFilter represents filtering options for software engineer queries
type SoftwareEngineerFilter struct {
	Search    string // Search by name, email, or employee number
	SortBy    string // Column to sort by (e.g., "name", "email", "created_at")
	SortOrder string // Sort order: "asc" or "desc"
}

// GetAllSoftwareEngineers retrieves all software engineers from the database with optional filtering
func GetAllSoftwareEngineers(db interface{ Query(query string, args ...interface{}) (*sql.Rows, error) }, filter *SoftwareEngineerFilter) ([]SoftwareEngineer, error) {
	query := `
		SELECT id, name, email, phone, address, address_street, address_city, address_country, address_state, address_postal_code, employee_number, address_confirmed, address_confirmation_at, created_at, updated_at
		FROM software_engineers
	`
	
	var conditions []string
	var args []interface{}
	argCount := 0

	// Apply search filter if provided
	if filter != nil && filter.Search != "" {
		argCount++
		searchPattern := "%" + strings.ToLower(filter.Search) + "%"
		conditions = append(conditions, fmt.Sprintf("(LOWER(name) LIKE $%d OR LOWER(email) LIKE $%d OR LOWER(COALESCE(employee_number, '')) LIKE $%d)", argCount, argCount, argCount))
		args = append(args, searchPattern)
	}

	// Add WHERE clause if there are conditions
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Add ordering
	orderBy := buildSoftwareEngineerOrderByClause(filter)
	query += " " + orderBy

	var rows *sql.Rows
	var err error
	if len(args) > 0 {
		rows, err = db.Query(query, args...)
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query software engineers: %w", err)
	}
	defer rows.Close()

	var engineers []SoftwareEngineer
	for rows.Next() {
		var engineer SoftwareEngineer
		var phone sql.NullString
		var address sql.NullString
		var addressStreet sql.NullString
		var addressCity sql.NullString
		var addressCountry sql.NullString
		var addressState sql.NullString
		var addressPostalCode sql.NullString
		var employeeNumber sql.NullString
		var addressConfirmationAt sql.NullTime

		err := rows.Scan(
			&engineer.ID,
			&engineer.Name,
			&engineer.Email,
			&phone,
			&address,
			&addressStreet,
			&addressCity,
			&addressCountry,
			&addressState,
			&addressPostalCode,
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
		if addressStreet.Valid {
			engineer.AddressStreet = addressStreet.String
		}
		if addressCity.Valid {
			engineer.AddressCity = addressCity.String
		}
		if addressCountry.Valid {
			engineer.AddressCountry = addressCountry.String
		}
		if addressState.Valid {
			engineer.AddressState = addressState.String
		}
		if addressPostalCode.Valid {
			engineer.AddressPostalCode = addressPostalCode.String
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
		SELECT id, name, email, phone, address, address_street, address_city, address_country, address_state, address_postal_code, employee_number, address_confirmed, address_confirmation_at, created_at, updated_at
		FROM software_engineers
		WHERE id = $1
	`

	var engineer SoftwareEngineer
	var phone sql.NullString
	var address sql.NullString
	var addressStreet sql.NullString
	var addressCity sql.NullString
	var addressCountry sql.NullString
	var addressState sql.NullString
	var addressPostalCode sql.NullString
	var employeeNumber sql.NullString
	var addressConfirmationAt sql.NullTime

	err := db.QueryRow(query, id).Scan(
		&engineer.ID,
		&engineer.Name,
		&engineer.Email,
		&phone,
		&address,
		&addressStreet,
		&addressCity,
		&addressCountry,
		&addressState,
		&addressPostalCode,
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
	if addressStreet.Valid {
		engineer.AddressStreet = addressStreet.String
	}
	if addressCity.Valid {
		engineer.AddressCity = addressCity.String
	}
	if addressCountry.Valid {
		engineer.AddressCountry = addressCountry.String
	}
	if addressState.Valid {
		engineer.AddressState = addressState.String
	}
	if addressPostalCode.Valid {
		engineer.AddressPostalCode = addressPostalCode.String
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
		INSERT INTO software_engineers (name, email, address, address_street, address_city, address_country, address_state, address_postal_code, phone, employee_number, address_confirmed, address_confirmation_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id
	`

	err := db.QueryRow(
		query,
		engineer.Name,
		engineer.Email,
		engineer.Address,
		engineer.AddressStreet,
		engineer.AddressCity,
		engineer.AddressCountry,
		engineer.AddressState,
		engineer.AddressPostalCode,
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
		SET name = $1, email = $2, address = $3, address_street = $4, address_city = $5, address_country = $6, address_state = $7, address_postal_code = $8, phone = $9, employee_number = $10, address_confirmed = $11, address_confirmation_at = $12, updated_at = $13
		WHERE id = $14
	`

	result, err := db.Exec(
		query,
		engineer.Name,
		engineer.Email,
		engineer.Address,
		engineer.AddressStreet,
		engineer.AddressCity,
		engineer.AddressCountry,
		engineer.AddressState,
		engineer.AddressPostalCode,
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

// buildSoftwareEngineerOrderByClause builds the ORDER BY clause based on the filter
func buildSoftwareEngineerOrderByClause(filter *SoftwareEngineerFilter) string {
	// Default sort: by name ASC
	if filter == nil {
		return "ORDER BY name ASC"
	}

	// Map of allowed sort columns to their SQL equivalents
	sortColumns := map[string]string{
		"name":       "name",
		"email":      "email",
		"phone":      "phone",
		"created_at": "created_at",
		"updated_at": "updated_at",
	}

	// Validate sort order
	sortOrder := "ASC"
	if filter.SortOrder == "desc" {
		sortOrder = "DESC"
	}

	// If no sort column specified, use default
	if filter.SortBy == "" {
		return "ORDER BY name ASC"
	}

	// Get the SQL column name
	sqlColumn, exists := sortColumns[filter.SortBy]
	if !exists {
		// If invalid column, use default
		return "ORDER BY name ASC"
	}

	// Text columns that should use COLLATE
	textColumns := map[string]bool{
		"name":  true,
		"email": true,
		"phone": true,
	}

	// Return the ORDER BY clause with the specified column and order
	if textColumns[filter.SortBy] {
		return fmt.Sprintf("ORDER BY %s COLLATE \"C\" %s", sqlColumn, sortOrder)
	}

	return fmt.Sprintf("ORDER BY %s %s", sqlColumn, sortOrder)
}

