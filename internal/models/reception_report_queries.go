package models

import (
	"context"
	"database/sql"
	"errors"
)

// GetLaptopReceptionReport retrieves the reception report for a specific laptop
func GetLaptopReceptionReport(ctx context.Context, db *sql.DB, laptopID int64) (*ReceptionReport, error) {
	query := `
		SELECT 
			id, laptop_id, shipment_id, client_company_id, tracking_number,
			warehouse_user_id, received_at, notes,
			photo_serial_number, photo_external_condition, photo_working_condition,
			status, approved_by, approved_at, created_at, updated_at
		FROM reception_reports
		WHERE laptop_id = $1
	`

	report := &ReceptionReport{}
	var shipmentID sql.NullInt64
	var clientCompanyID sql.NullInt64
	var trackingNumber sql.NullString
	var notes sql.NullString
	var approvedBy sql.NullInt64
	var approvedAt sql.NullTime
	
	err := db.QueryRowContext(ctx, query, laptopID).Scan(
		&report.ID,
		&report.LaptopID,
		&shipmentID,
		&clientCompanyID,
		&trackingNumber,
		&report.WarehouseUserID,
		&report.ReceivedAt,
		&notes,
		&report.PhotoSerialNumber,
		&report.PhotoExternalCondition,
		&report.PhotoWorkingCondition,
		&report.Status,
		&approvedBy,
		&approvedAt,
		&report.CreatedAt,
		&report.UpdatedAt,
	)
	
	// Handle nullable fields
	if shipmentID.Valid {
		report.ShipmentID = &shipmentID.Int64
	}
	if clientCompanyID.Valid {
		report.ClientCompanyID = &clientCompanyID.Int64
	}
	if trackingNumber.Valid {
		report.TrackingNumber = trackingNumber.String
	}
	if notes.Valid {
		report.Notes = notes.String
	}
	if approvedBy.Valid {
		report.ApprovedBy = &approvedBy.Int64
	}
	if approvedAt.Valid {
		report.ApprovedAt = &approvedAt.Time
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No reception report found (not an error)
		}
		return nil, err
	}

	return report, nil
}

// CreateReceptionReport creates a new reception report in the database
func CreateReceptionReport(ctx context.Context, db *sql.DB, report *ReceptionReport) error {
	// Set timestamps
	report.BeforeCreate()

	query := `
		INSERT INTO reception_reports (
			laptop_id, shipment_id, client_company_id, tracking_number,
			warehouse_user_id, received_at, notes,
			photo_serial_number, photo_external_condition, photo_working_condition,
			status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id
	`

	err := db.QueryRowContext(ctx, query,
		report.LaptopID,
		report.ShipmentID,
		report.ClientCompanyID,
		report.TrackingNumber,
		report.WarehouseUserID,
		report.ReceivedAt,
		report.Notes,
		report.PhotoSerialNumber,
		report.PhotoExternalCondition,
		report.PhotoWorkingCondition,
		report.Status,
		report.CreatedAt,
		report.UpdatedAt,
	).Scan(&report.ID)

	return err
}

// ApproveReceptionReport approves a reception report and updates the laptop status
func ApproveReceptionReport(ctx context.Context, db *sql.DB, reportID int64, logisticsUserID int64) error {
	// Start a transaction to ensure both operations succeed or fail together
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// First, get the reception report to get the laptop ID
	var laptopID int64
	var currentStatus ReceptionReportStatus
	err = tx.QueryRowContext(ctx,
		`SELECT laptop_id, status FROM reception_reports WHERE id = $1`,
		reportID,
	).Scan(&laptopID, &currentStatus)

	if err == sql.ErrNoRows {
		return errors.New("reception report not found")
	}
	if err != nil {
		return err
	}

	// Check if already approved
	if currentStatus == ReceptionReportStatusApproved {
		return errors.New("reception report is already approved")
	}

	// Update the reception report status
	_, err = tx.ExecContext(ctx,
		`UPDATE reception_reports 
		SET status = $1, approved_by = $2, approved_at = NOW(), updated_at = NOW()
		WHERE id = $3`,
		ReceptionReportStatusApproved,
		logisticsUserID,
		reportID,
	)
	if err != nil {
		return err
	}

	// Update the laptop status to available
	_, err = tx.ExecContext(ctx,
		`UPDATE laptops 
		SET status = $1, updated_at = NOW()
		WHERE id = $2 AND status = $3`,
		LaptopStatusAvailable,
		laptopID,
		LaptopStatusAtWarehouse,
	)
	if err != nil {
		return err
	}

	// Commit the transaction
	return tx.Commit()
}

// GetReceptionReportByID retrieves a reception report by its ID
func GetReceptionReportByID(ctx context.Context, db *sql.DB, reportID int64) (*ReceptionReport, error) {
	query := `
		SELECT 
			id, laptop_id, shipment_id, client_company_id, tracking_number,
			warehouse_user_id, received_at, notes,
			photo_serial_number, photo_external_condition, photo_working_condition,
			status, approved_by, approved_at, created_at, updated_at
		FROM reception_reports
		WHERE id = $1
	`

	report := &ReceptionReport{}
	err := db.QueryRowContext(ctx, query, reportID).Scan(
		&report.ID,
		&report.LaptopID,
		&report.ShipmentID,
		&report.ClientCompanyID,
		&report.TrackingNumber,
		&report.WarehouseUserID,
		&report.ReceivedAt,
		&report.Notes,
		&report.PhotoSerialNumber,
		&report.PhotoExternalCondition,
		&report.PhotoWorkingCondition,
		&report.Status,
		&report.ApprovedBy,
		&report.ApprovedAt,
		&report.CreatedAt,
		&report.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return report, nil
}

