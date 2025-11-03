package models

import (
	"database/sql"
	"fmt"
	"time"
)

// ChartDataPoint represents a single point in a time series chart
type ChartDataPoint struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

// StatusDistribution represents shipment status distribution for pie charts
type StatusDistribution struct {
	Status      ShipmentStatus `json:"status"`
	StatusLabel string         `json:"status_label"`
	Count       int            `json:"count"`
	Percentage  float64        `json:"percentage"`
}

// DeliveryTimeTrend represents delivery time trends over weeks
type DeliveryTimeTrend struct {
	WeekStart           string  `json:"week_start"`
	AverageDeliveryDays float64 `json:"average_delivery_days"`
	ShipmentCount       int     `json:"shipment_count"`
}

// GetShipmentsOverTime returns shipment counts grouped by date for timeline charts
// days parameter specifies how many days back to retrieve
func GetShipmentsOverTime(db *sql.DB, days int) ([]ChartDataPoint, error) {
	query := `
		SELECT DATE(created_at) as date, COUNT(*) as count
		FROM shipments
		WHERE DATE(created_at) >= CURRENT_DATE - INTERVAL '%d days'
		GROUP BY DATE(created_at)
		ORDER BY date ASC
	`

	rows, err := db.Query(fmt.Sprintf(query, days))
	if err != nil {
		return nil, fmt.Errorf("failed to query shipments over time: %w", err)
	}
	defer rows.Close()

	var dataPoints []ChartDataPoint
	for rows.Next() {
		var point ChartDataPoint
		var date time.Time
		if err := rows.Scan(&date, &point.Count); err != nil {
			return nil, fmt.Errorf("failed to scan data point: %w", err)
		}
		point.Date = date.Format("2006-01-02")
		dataPoints = append(dataPoints, point)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating data points: %w", err)
	}

	return dataPoints, nil
}

// GetShipmentStatusDistribution returns shipment count and percentage by status for pie charts
func GetShipmentStatusDistribution(db *sql.DB) ([]StatusDistribution, error) {
	// First get total count
	var totalCount int
	err := db.QueryRow("SELECT COUNT(*) FROM shipments").Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get total shipment count: %w", err)
	}

	if totalCount == 0 {
		return []StatusDistribution{}, nil
	}

	// Get counts by status
	query := `
		SELECT status, COUNT(*) as count
		FROM shipments
		GROUP BY status
		ORDER BY count DESC
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query status distribution: %w", err)
	}
	defer rows.Close()

	var distribution []StatusDistribution
	for rows.Next() {
		var item StatusDistribution
		if err := rows.Scan(&item.Status, &item.Count); err != nil {
			return nil, fmt.Errorf("failed to scan distribution item: %w", err)
		}
		
		// Calculate percentage
		item.Percentage = (float64(item.Count) / float64(totalCount)) * 100
		
		// Create readable label
		item.StatusLabel = formatStatusLabel(item.Status)
		
		distribution = append(distribution, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating distribution: %w", err)
	}

	return distribution, nil
}

// GetDeliveryTimeTrends returns average delivery time grouped by week
// weeks parameter specifies how many weeks back to retrieve
func GetDeliveryTimeTrends(db *sql.DB, weeks int) ([]DeliveryTimeTrend, error) {
	query := `
		SELECT 
			DATE_TRUNC('week', delivered_at) as week_start,
			AVG(EXTRACT(EPOCH FROM (delivered_at - picked_up_at)) / 86400) as avg_days,
			COUNT(*) as shipment_count
		FROM shipments
		WHERE status = $1 
		  AND picked_up_at IS NOT NULL 
		  AND delivered_at IS NOT NULL
		  AND delivered_at >= NOW() - INTERVAL '%d weeks'
		GROUP BY DATE_TRUNC('week', delivered_at)
		ORDER BY week_start ASC
	`

	rows, err := db.Query(fmt.Sprintf(query, weeks), ShipmentStatusDelivered)
	if err != nil {
		return nil, fmt.Errorf("failed to query delivery time trends: %w", err)
	}
	defer rows.Close()

	var trends []DeliveryTimeTrend
	for rows.Next() {
		var trend DeliveryTimeTrend
		var weekStart time.Time
		var avgDays sql.NullFloat64
		
		if err := rows.Scan(&weekStart, &avgDays, &trend.ShipmentCount); err != nil {
			return nil, fmt.Errorf("failed to scan trend: %w", err)
		}
		
		trend.WeekStart = weekStart.Format("2006-01-02")
		if avgDays.Valid {
			trend.AverageDeliveryDays = avgDays.Float64
		}
		
		trends = append(trends, trend)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating trends: %w", err)
	}

	return trends, nil
}

// formatStatusLabel converts a status enum to a human-readable label
func formatStatusLabel(status ShipmentStatus) string {
	labels := map[ShipmentStatus]string{
		ShipmentStatusPendingPickup:         "Pending Pickup",
		ShipmentStatusPickedUpFromClient:    "Picked Up",
		ShipmentStatusInTransitToWarehouse:  "In Transit to Warehouse",
		ShipmentStatusAtWarehouse:           "At Warehouse",
		ShipmentStatusReleasedFromWarehouse: "Released from Warehouse",
		ShipmentStatusInTransitToEngineer:   "In Transit to Engineer",
		ShipmentStatusDelivered:             "Delivered",
	}

	if label, ok := labels[status]; ok {
		return label
	}
	return string(status)
}

