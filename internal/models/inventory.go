package models

import (
	"database/sql"
	"fmt"
	"strings"
)

// LaptopFilter represents filtering options for laptop queries
type LaptopFilter struct {
	Status   LaptopStatus
	Brand    string
	Search   string
	Limit    int
	Offset   int
	UserRole UserRole // Filter laptops based on user role permissions
}

// GetAllLaptops retrieves all laptops with optional filtering
func GetAllLaptops(db *sql.DB, filter *LaptopFilter) ([]Laptop, error) {
	query := `
		SELECT 
			l.id, l.serial_number, l.sku, l.brand, l.model, l.cpu, l.ram_gb, l.ssd_gb, l.status, 
			l.client_company_id, l.software_engineer_id, l.created_at, l.updated_at,
			cc.name as client_company_name,
			se.name as software_engineer_name,
			COALESCE(rr.id IS NOT NULL, false) as has_reception_report,
			rr.id as reception_report_id,
			rr.status as reception_report_status
		FROM laptops l
		LEFT JOIN client_companies cc ON cc.id = l.client_company_id
		LEFT JOIN software_engineers se ON se.id = l.software_engineer_id
		LEFT JOIN reception_reports rr ON rr.laptop_id = l.id
	`

	var conditions []string
	var args []interface{}
	argCount := 0

	// Apply filters if provided
	if filter != nil {
		// Role-based filtering: Warehouse users only see specific statuses
		if filter.UserRole == RoleWarehouse {
			// Warehouse users can only see: in_transit_to_warehouse, at_warehouse, available
			conditions = append(conditions, "l.status IN ('in_transit_to_warehouse', 'at_warehouse', 'available')")
		}

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
		var clientCompanyName sql.NullString
		var softwareEngineerName sql.NullString
		var receptionReportID sql.NullInt64
		var receptionReportStatus sql.NullString

		err := rows.Scan(
			&laptop.ID,
			&laptop.SerialNumber,
			&sku,
			&brand,
			&laptop.Model,
			&laptop.CPU,
			&laptop.RAMGB,
			&laptop.SSDGB,
			&laptop.Status,
			&laptop.ClientCompanyID,
			&laptop.SoftwareEngineerID,
			&laptop.CreatedAt,
			&laptop.UpdatedAt,
			&clientCompanyName,
			&softwareEngineerName,
			&laptop.HasReceptionReport,
			&receptionReportID,
			&receptionReportStatus,
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
		if clientCompanyName.Valid {
			laptop.ClientCompanyName = clientCompanyName.String
		}
		if softwareEngineerName.Valid {
			laptop.SoftwareEngineerName = softwareEngineerName.String
		}
		if receptionReportID.Valid {
			laptop.ReceptionReportID = &receptionReportID.Int64
		}
		if receptionReportStatus.Valid {
			laptop.ReceptionReportStatus = receptionReportStatus.String
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
		SELECT id, serial_number, brand, model, cpu, ram_gb, ssd_gb, status, created_at, updated_at
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
			&laptop.CPU,
			&laptop.RAMGB,
			&laptop.SSDGB,
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
			l.id, l.serial_number, l.sku, l.brand, l.model, l.cpu, l.ram_gb, l.ssd_gb, l.status, 
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
		&laptop.CPU,
		&laptop.RAMGB,
		&laptop.SSDGB,
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
	// Auto-generate SKU if not provided
	laptop.GenerateAndSetSKU()

	// Validate laptop
	if err := laptop.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Set timestamps
	laptop.BeforeCreate()

	query := `
		INSERT INTO laptops (serial_number, sku, brand, model, cpu, ram_gb, ssd_gb, status, client_company_id, software_engineer_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id
	`

	err := db.QueryRow(
		query,
		laptop.SerialNumber,
		laptop.SKU,
		laptop.Brand,
		laptop.Model,
		laptop.CPU,
		laptop.RAMGB,
		laptop.SSDGB,
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
		SET serial_number = $1, sku = $2, brand = $3, model = $4, cpu = $5, ram_gb = $6, ssd_gb = $7, status = $8, 
		    client_company_id = $9, software_engineer_id = $10, updated_at = $11
		WHERE id = $12
	`

	result, err := db.Exec(
		query,
		laptop.SerialNumber,
		laptop.SKU,
		laptop.Brand,
		laptop.Model,
		laptop.CPU,
		laptop.RAMGB,
		laptop.SSDGB,
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
		SELECT id, serial_number, brand, model, cpu, ram_gb, ssd_gb, status, created_at, updated_at
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
			&laptop.CPU,
			&laptop.RAMGB,
			&laptop.SSDGB,
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

