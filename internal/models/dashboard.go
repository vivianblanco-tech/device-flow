package models

import (
	"database/sql"
	"fmt"
)

// DashboardStats represents dashboard statistics
type DashboardStats struct {
	TotalShipments         int                      `json:"total_shipments"`
	PendingPickups         int                      `json:"pending_pickups"`
	InTransit              int                      `json:"in_transit"`
	Delivered              int                      `json:"delivered"`
	AvgDeliveryDays        float64                  `json:"avg_delivery_days"`
	ShipmentsByStatus      map[ShipmentStatus]int   `json:"shipments_by_status"`
	LaptopsByStatus        map[LaptopStatus]int     `json:"laptops_by_status"`
	AvailableLaptops       int                      `json:"available_laptops"`
}

// GetShipmentCountsByStatus returns the count of shipments grouped by status
func GetShipmentCountsByStatus(db *sql.DB) (map[ShipmentStatus]int, error) {
	query := `
		SELECT status, COUNT(*) as count
		FROM shipments
		GROUP BY status
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query shipment counts: %w", err)
	}
	defer rows.Close()

	counts := make(map[ShipmentStatus]int)
	for rows.Next() {
		var status ShipmentStatus
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, fmt.Errorf("failed to scan shipment count: %w", err)
		}
		counts[status] = count
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating shipment counts: %w", err)
	}

	return counts, nil
}

// GetTotalShipmentCount returns the total count of all shipments
func GetTotalShipmentCount(db *sql.DB) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM shipments`
	
	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get total shipment count: %w", err)
	}

	return count, nil
}

// GetAverageDeliveryTime calculates the average delivery time in days
// for shipments that have been delivered (from pickup to delivery)
func GetAverageDeliveryTime(db *sql.DB) (float64, error) {
	query := `
		SELECT AVG(EXTRACT(EPOCH FROM (delivered_at - picked_up_at)) / 86400) as avg_days
		FROM shipments
		WHERE status = $1 
		  AND picked_up_at IS NOT NULL 
		  AND delivered_at IS NOT NULL
	`

	var avgDays sql.NullFloat64
	err := db.QueryRow(query, ShipmentStatusDelivered).Scan(&avgDays)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate average delivery time: %w", err)
	}

	if !avgDays.Valid {
		return 0, nil
	}

	return avgDays.Float64, nil
}

// GetInTransitShipmentCount returns the count of shipments currently in transit
func GetInTransitShipmentCount(db *sql.DB) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM shipments 
		WHERE status IN ($1, $2)
	`

	var count int
	err := db.QueryRow(
		query,
		ShipmentStatusInTransitToWarehouse,
		ShipmentStatusInTransitToEngineer,
	).Scan(&count)
	
	if err != nil {
		return 0, fmt.Errorf("failed to get in-transit shipment count: %w", err)
	}

	return count, nil
}

// GetPendingPickupCount returns the count of shipments pending pickup
func GetPendingPickupCount(db *sql.DB) (int, error) {
	query := `SELECT COUNT(*) FROM shipments WHERE status = $1`

	var count int
	err := db.QueryRow(query, ShipmentStatusPendingPickup).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get pending pickup count: %w", err)
	}

	return count, nil
}

// GetLaptopCountsByStatus returns the count of laptops grouped by status
func GetLaptopCountsByStatus(db *sql.DB) (map[LaptopStatus]int, error) {
	query := `
		SELECT status, COUNT(*) as count
		FROM laptops
		GROUP BY status
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query laptop counts: %w", err)
	}
	defer rows.Close()

	counts := make(map[LaptopStatus]int)
	for rows.Next() {
		var status LaptopStatus
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, fmt.Errorf("failed to scan laptop count: %w", err)
		}
		counts[status] = count
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating laptop counts: %w", err)
	}

	return counts, nil
}

// GetAvailableLaptopCount returns the count of laptops available for assignment
func GetAvailableLaptopCount(db *sql.DB) (int, error) {
	query := `SELECT COUNT(*) FROM laptops WHERE status = $1`

	var count int
	err := db.QueryRow(query, LaptopStatusAvailable).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get available laptop count: %w", err)
	}

	return count, nil
}

// GetDashboardStats retrieves all dashboard statistics in one call
func GetDashboardStats(db *sql.DB) (*DashboardStats, error) {
	stats := &DashboardStats{}

	// Get total shipments
	totalShipments, err := GetTotalShipmentCount(db)
	if err != nil {
		return nil, fmt.Errorf("failed to get total shipments: %w", err)
	}
	stats.TotalShipments = totalShipments

	// Get pending pickups
	pendingPickups, err := GetPendingPickupCount(db)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending pickups: %w", err)
	}
	stats.PendingPickups = pendingPickups

	// Get in-transit count
	inTransit, err := GetInTransitShipmentCount(db)
	if err != nil {
		return nil, fmt.Errorf("failed to get in-transit count: %w", err)
	}
	stats.InTransit = inTransit

	// Get shipments by status
	shipmentsByStatus, err := GetShipmentCountsByStatus(db)
	if err != nil {
		return nil, fmt.Errorf("failed to get shipments by status: %w", err)
	}
	stats.ShipmentsByStatus = shipmentsByStatus
	stats.Delivered = shipmentsByStatus[ShipmentStatusDelivered]

	// Get average delivery time
	avgDeliveryDays, err := GetAverageDeliveryTime(db)
	if err != nil {
		return nil, fmt.Errorf("failed to get average delivery time: %w", err)
	}
	stats.AvgDeliveryDays = avgDeliveryDays

	// Get laptops by status
	laptopsByStatus, err := GetLaptopCountsByStatus(db)
	if err != nil {
		return nil, fmt.Errorf("failed to get laptops by status: %w", err)
	}
	stats.LaptopsByStatus = laptopsByStatus

	// Get available laptops
	availableLaptops, err := GetAvailableLaptopCount(db)
	if err != nil {
		return nil, fmt.Errorf("failed to get available laptops: %w", err)
	}
	stats.AvailableLaptops = availableLaptops

	return stats, nil
}

