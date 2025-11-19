package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// Courier represents a courier company in the system
type Courier struct {
	ID          int64     `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	ContactInfo string    `json:"contact_info,omitempty" db:"contact_info"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Validate validates the Courier model
func (c *Courier) Validate() error {
	// Name validation
	if c.Name == "" {
		return errors.New("courier name is required")
	}
	if len(c.Name) < 2 {
		return errors.New("courier name must be at least 2 characters")
	}

	return nil
}

// TableName returns the table name for the Courier model
func (c *Courier) TableName() string {
	return "couriers"
}

// BeforeCreate sets the timestamps before creating a courier
func (c *Courier) BeforeCreate() {
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now
}

// BeforeUpdate sets the updated_at timestamp before updating a courier
func (c *Courier) BeforeUpdate() {
	c.UpdatedAt = time.Now()
}

// GetAllCouriers retrieves all couriers from the database
func GetAllCouriers(db interface{ Query(query string, args ...interface{}) (*sql.Rows, error) }) ([]Courier, error) {
	query := `
		SELECT id, name, contact_info, created_at, updated_at
		FROM couriers
		ORDER BY name ASC
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query couriers: %w", err)
	}
	defer rows.Close()

	var couriers []Courier
	for rows.Next() {
		var courier Courier
		err := rows.Scan(
			&courier.ID,
			&courier.Name,
			&courier.ContactInfo,
			&courier.CreatedAt,
			&courier.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan courier: %w", err)
		}
		couriers = append(couriers, courier)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating couriers: %w", err)
	}

	return couriers, nil
}

// GetCourierByID retrieves a courier by its ID
func GetCourierByID(db *sql.DB, id int64) (*Courier, error) {
	query := `
		SELECT id, name, contact_info, created_at, updated_at
		FROM couriers
		WHERE id = $1
	`

	var courier Courier
	err := db.QueryRow(query, id).Scan(
		&courier.ID,
		&courier.Name,
		&courier.ContactInfo,
		&courier.CreatedAt,
		&courier.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("courier not found with id %d", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get courier: %w", err)
	}

	return &courier, nil
}

// CreateCourier creates a new courier in the database
func CreateCourier(db *sql.DB, courier *Courier) error {
	// Validate courier
	if err := courier.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Set timestamps
	courier.BeforeCreate()

	query := `
		INSERT INTO couriers (name, contact_info, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err := db.QueryRow(
		query,
		courier.Name,
		courier.ContactInfo,
		courier.CreatedAt,
		courier.UpdatedAt,
	).Scan(&courier.ID)

	if err != nil {
		return fmt.Errorf("failed to create courier: %w", err)
	}

	return nil
}

// UpdateCourier updates an existing courier in the database
func UpdateCourier(db *sql.DB, courier *Courier) error {
	// Validate courier
	if err := courier.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Update timestamp
	courier.BeforeUpdate()

	query := `
		UPDATE couriers
		SET name = $1, contact_info = $2, updated_at = $3
		WHERE id = $4
	`

	result, err := db.Exec(
		query,
		courier.Name,
		courier.ContactInfo,
		courier.UpdatedAt,
		courier.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update courier: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("courier not found with id %d", courier.ID)
	}

	return nil
}

// DeleteCourier deletes a courier from the database
func DeleteCourier(db *sql.DB, id int64) error {
	query := `DELETE FROM couriers WHERE id = $1`

	result, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete courier: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("courier not found with id %d", id)
	}

	return nil
}

// CourierExistsByName checks if a courier with the given name exists in the database
func CourierExistsByName(db interface{ QueryRow(query string, args ...interface{}) *sql.Row }, name string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM couriers WHERE name = $1)`
	
	var exists bool
	err := db.QueryRow(query, name).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if courier exists: %w", err)
	}
	
	return exists, nil
}

// IsValidCourierName checks if a courier name is valid by checking the database first,
// then falling back to hardcoded values (UPS, FedEx, DHL) for backward compatibility
func IsValidCourierName(db interface{ QueryRow(query string, args ...interface{}) *sql.Row }, name string) (bool, error) {
	// First check database
	exists, err := CourierExistsByName(db, name)
	if err != nil {
		return false, err
	}
	if exists {
		return true, nil
	}
	
	// Fall back to hardcoded values for backward compatibility
	switch name {
	case "UPS", "FedEx", "DHL":
		return true, nil
	}
	
	return false, nil
}

