package models

import (
	"database/sql"
	"fmt"
	"strings"
)

// LaptopFilter represents filtering options for laptop queries
type LaptopFilter struct {
	Status LaptopStatus
	Brand  string
	Search string
	Limit  int
	Offset int
}

// GetAllLaptops retrieves all laptops with optional filtering
func GetAllLaptops(db *sql.DB, filter *LaptopFilter) ([]Laptop, error) {
	query := `
		SELECT 
			l.id, l.serial_number, l.sku, l.brand, l.model, l.specs, l.status, 
			l.client_company_id, l.software_engineer_id, l.created_at, l.updated_at,
			cc.name as client_company_name,
			se.name as software_engineer_name
		FROM laptops l
		LEFT JOIN client_companies cc ON cc.id = l.client_company_id
		LEFT JOIN software_engineers se ON se.id = l.software_engineer_id
	`

	var conditions []string
	var args []interface{}
	argCount := 0

	// Apply filters if provided
	if filter != nil {
		if filter.Status != "" {
			argCount++
			conditions = append(conditions, fmt.Sprintf("l.status = $%d", argCount))
			args = append(args, filter.Status)
		}

		if filter.Brand != "" {
			argCount++
			conditions = append(conditions, fmt.Sprintf("LOWER(l.brand) = LOWER($%d)", argCount))
			args = append(args, filter.Brand)
		}

		if filter.Search != "" {
			argCount++
			searchPattern := "%" + strings.ToLower(filter.Search) + "%"
			conditions = append(conditions, fmt.Sprintf("(LOWER(l.serial_number) LIKE $%d OR LOWER(l.brand) LIKE $%d OR LOWER(l.model) LIKE $%d OR LOWER(l.sku) LIKE $%d)", argCount, argCount, argCount, argCount))
			args = append(args, searchPattern)
		}
	}

	// Add WHERE clause if there are conditions
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Add ordering
	query += " ORDER BY l.created_at DESC"

	// Add pagination if specified
	if filter != nil {
		if filter.Limit > 0 {
			argCount++
			query += fmt.Sprintf(" LIMIT $%d", argCount)
			args = append(args, filter.Limit)
		}
		if filter.Offset > 0 {
			argCount++
			query += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, filter.Offset)
		}
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query laptops: %w", err)
	}
	defer rows.Close()

	var laptops []Laptop
	for rows.Next() {
		var laptop Laptop
		var sku sql.NullString
		var brand sql.NullString
		var model sql.NullString
		var specs sql.NullString
		var clientCompanyName sql.NullString
		var softwareEngineerName sql.NullString

		err := rows.Scan(
			&laptop.ID,
			&laptop.SerialNumber,
			&sku,
			&brand,
			&model,
			&specs,
			&laptop.Status,
			&laptop.ClientCompanyID,
			&laptop.SoftwareEngineerID,
			&laptop.CreatedAt,
			&laptop.UpdatedAt,
			&clientCompanyName,
			&softwareEngineerName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan laptop: %w", err)
		}

		// Set nullable fields if available
		if sku.Valid {
			laptop.SKU = sku.String
		}
		if brand.Valid {
			laptop.Brand = brand.String
		}
		if model.Valid {
			laptop.Model = model.String
		}
		if specs.Valid {
			laptop.Specs = specs.String
		}
		if clientCompanyName.Valid {
			laptop.ClientCompanyName = clientCompanyName.String
		}
		if softwareEngineerName.Valid {
			laptop.SoftwareEngineerName = softwareEngineerName.String
		}

		laptops = append(laptops, laptop)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating laptops: %w", err)
	}

	return laptops, nil
}

// SearchLaptops searches for laptops by serial number, brand, or model
func SearchLaptops(db *sql.DB, searchTerm string) ([]Laptop, error) {
	query := `
		SELECT id, serial_number, brand, model, specs, status, created_at, updated_at
		FROM laptops
		WHERE LOWER(serial_number) LIKE $1 
		   OR LOWER(brand) LIKE $1 
		   OR LOWER(model) LIKE $1
		ORDER BY created_at DESC
	`

	searchPattern := "%" + strings.ToLower(searchTerm) + "%"
	rows, err := db.Query(query, searchPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to search laptops: %w", err)
	}
	defer rows.Close()

	var laptops []Laptop
	for rows.Next() {
		var laptop Laptop
		err := rows.Scan(
			&laptop.ID,
			&laptop.SerialNumber,
			&laptop.Brand,
			&laptop.Model,
			&laptop.Specs,
			&laptop.Status,
			&laptop.CreatedAt,
			&laptop.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan laptop: %w", err)
		}
		laptops = append(laptops, laptop)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating search results: %w", err)
	}

	return laptops, nil
}

// GetLaptopByID retrieves a laptop by its ID
func GetLaptopByID(db *sql.DB, id int64) (*Laptop, error) {
	query := `
		SELECT 
			l.id, l.serial_number, l.sku, l.brand, l.model, l.specs, l.status, 
			l.client_company_id, l.software_engineer_id, l.created_at, l.updated_at,
			cc.name as client_company_name,
			se.name as software_engineer_name
		FROM laptops l
		LEFT JOIN client_companies cc ON cc.id = l.client_company_id
		LEFT JOIN software_engineers se ON se.id = l.software_engineer_id
		WHERE l.id = $1
	`

	var laptop Laptop
	var sku sql.NullString
	var clientCompanyName sql.NullString
	var softwareEngineerName sql.NullString

	err := db.QueryRow(query, id).Scan(
		&laptop.ID,
		&laptop.SerialNumber,
		&sku,
		&laptop.Brand,
		&laptop.Model,
		&laptop.Specs,
		&laptop.Status,
		&laptop.ClientCompanyID,
		&laptop.SoftwareEngineerID,
		&laptop.CreatedAt,
		&laptop.UpdatedAt,
		&clientCompanyName,
		&softwareEngineerName,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("laptop not found with id %d", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get laptop: %w", err)
	}

	// Set nullable fields if available
	if sku.Valid {
		laptop.SKU = sku.String
	}
	if clientCompanyName.Valid {
		laptop.ClientCompanyName = clientCompanyName.String
	}
	if softwareEngineerName.Valid {
		laptop.SoftwareEngineerName = softwareEngineerName.String
	}

	return &laptop, nil
}

// CreateLaptop creates a new laptop in the database
func CreateLaptop(db *sql.DB, laptop *Laptop) error {
	// Validate laptop
	if err := laptop.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Set timestamps
	laptop.BeforeCreate()

	query := `
		INSERT INTO laptops (serial_number, sku, brand, model, specs, status, client_company_id, software_engineer_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`

	err := db.QueryRow(
		query,
		laptop.SerialNumber,
		laptop.SKU,
		laptop.Brand,
		laptop.Model,
		laptop.Specs,
		laptop.Status,
		laptop.ClientCompanyID,
		laptop.SoftwareEngineerID,
		laptop.CreatedAt,
		laptop.UpdatedAt,
	).Scan(&laptop.ID)

	if err != nil {
		// Check for unique constraint violation (duplicate serial number)
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			return fmt.Errorf("laptop with serial number %s already exists", laptop.SerialNumber)
		}
		return fmt.Errorf("failed to create laptop: %w", err)
	}

	return nil
}

// UpdateLaptop updates an existing laptop in the database
func UpdateLaptop(db *sql.DB, laptop *Laptop) error {
	// Validate laptop
	if err := laptop.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Update timestamp
	laptop.BeforeUpdate()

	query := `
		UPDATE laptops
		SET serial_number = $1, sku = $2, brand = $3, model = $4, specs = $5, status = $6, 
		    client_company_id = $7, software_engineer_id = $8, updated_at = $9
		WHERE id = $10
	`

	result, err := db.Exec(
		query,
		laptop.SerialNumber,
		laptop.SKU,
		laptop.Brand,
		laptop.Model,
		laptop.Specs,
		laptop.Status,
		laptop.ClientCompanyID,
		laptop.SoftwareEngineerID,
		laptop.UpdatedAt,
		laptop.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update laptop: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("laptop not found with id %d", laptop.ID)
	}

	return nil
}

// DeleteLaptop deletes a laptop from the database
func DeleteLaptop(db *sql.DB, id int64) error {
	query := `DELETE FROM laptops WHERE id = $1`

	result, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete laptop: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("laptop not found with id %d", id)
	}

	return nil
}

// GetLaptopsByStatus retrieves all laptops with a specific status
func GetLaptopsByStatus(db *sql.DB, status LaptopStatus) ([]Laptop, error) {
	query := `
		SELECT id, serial_number, brand, model, specs, status, created_at, updated_at
		FROM laptops
		WHERE status = $1
		ORDER BY created_at DESC
	`

	rows, err := db.Query(query, status)
	if err != nil {
		return nil, fmt.Errorf("failed to query laptops by status: %w", err)
	}
	defer rows.Close()

	var laptops []Laptop
	for rows.Next() {
		var laptop Laptop
		err := rows.Scan(
			&laptop.ID,
			&laptop.SerialNumber,
			&laptop.Brand,
			&laptop.Model,
			&laptop.Specs,
			&laptop.Status,
			&laptop.CreatedAt,
			&laptop.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan laptop: %w", err)
		}
		laptops = append(laptops, laptop)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating laptops: %w", err)
	}

	return laptops, nil
}

