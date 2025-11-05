package models

import (
	"database/sql"
	"strconv"
	"testing"
	"time"

	"github.com/yourusername/laptop-tracking-system/internal/database"
)

// TestGetShipmentsOverTime tests retrieving shipment data for timeline charts
func TestGetShipmentsOverTime(t *testing.T) {
	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	// Create test client company
	company := &ClientCompany{
		Name:        "Test Company",
		ContactInfo: "test@example.com",
	}
	err := createClientCompany(db, company)
	if err != nil {
		t.Fatalf("Failed to create client company: %v", err)
	}

	// Create shipments on different dates
	now := time.Now()
	dates := []time.Time{
		now.AddDate(0, 0, -30), // 30 days ago
		now.AddDate(0, 0, -25),
		now.AddDate(0, 0, -20),
		now.AddDate(0, 0, -15),
		now.AddDate(0, 0, -10),
		now.AddDate(0, 0, -5),
		now.AddDate(0, 0, -2),
		now, // today
	}

	for i, date := range dates {
		shipment := &Shipment{
			ClientCompanyID:  company.ID,
			Status:           ShipmentStatusPendingPickup,
			JiraTicketNumber: "TEST-" + strconv.Itoa(i+1),
			CreatedAt:        date,
			UpdatedAt:        date,
		}
		err := createShipmentWithDate(db, shipment)
		if err != nil {
			t.Fatalf("Failed to create shipment: %v", err)
		}
	}

	// Test getting shipments over last 30 days
	chartData, err := GetShipmentsOverTime(db, 30)
	if err != nil {
		t.Fatalf("GetShipmentsOverTime failed: %v", err)
	}

	// Verify we got data points
	if len(chartData) == 0 {
		t.Error("Expected chart data, got empty result")
	}

	// Verify total count matches
	totalCount := 0
	for _, point := range chartData {
		totalCount += point.Count
	}
	if totalCount != len(dates) {
		t.Errorf("Expected total count of %d, got %d", len(dates), totalCount)
	}
}

// TestGetShipmentStatusDistribution tests retrieving status breakdown for pie charts
func TestGetShipmentStatusDistribution(t *testing.T) {
	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	// Create test client company
	company := &ClientCompany{
		Name:        "Test Company",
		ContactInfo: "test@example.com",
	}
	err := createClientCompany(db, company)
	if err != nil {
		t.Fatalf("Failed to create client company: %v", err)
	}

	// Create shipments with different statuses
	statuses := []ShipmentStatus{
		ShipmentStatusPendingPickup,
		ShipmentStatusPendingPickup,
		ShipmentStatusPendingPickup,
		ShipmentStatusInTransitToWarehouse,
		ShipmentStatusInTransitToWarehouse,
		ShipmentStatusAtWarehouse,
		ShipmentStatusDelivered,
		ShipmentStatusDelivered,
		ShipmentStatusDelivered,
		ShipmentStatusDelivered,
	}

	for i, status := range statuses {
		shipment := &Shipment{
			ClientCompanyID:  company.ID,
			Status:           status,
			JiraTicketNumber: "TEST-" + strconv.Itoa(i+1),
		}
		err := createShipment(db, shipment)
		if err != nil {
			t.Fatalf("Failed to create shipment: %v", err)
		}
	}

	// Test getting status distribution
	distribution, err := GetShipmentStatusDistribution(db)
	if err != nil {
		t.Fatalf("GetShipmentStatusDistribution failed: %v", err)
	}

	// Verify we got data
	if len(distribution) == 0 {
		t.Error("Expected distribution data, got empty result")
	}

	// Verify specific counts
	foundPending := false
	foundDelivered := false
	for _, item := range distribution {
		if item.Status == ShipmentStatusPendingPickup && item.Count == 3 {
			foundPending = true
		}
		if item.Status == ShipmentStatusDelivered && item.Count == 4 {
			foundDelivered = true
		}
	}

	if !foundPending {
		t.Error("Expected to find 3 pending pickups in distribution")
	}
	if !foundDelivered {
		t.Error("Expected to find 4 delivered shipments in distribution")
	}
}

// TestGetDeliveryTimeTrends tests retrieving delivery time trends for bar charts
func TestGetDeliveryTimeTrends(t *testing.T) {
	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	// Create test client company
	company := &ClientCompany{
		Name:        "Test Company",
		ContactInfo: "test@example.com",
	}
	err := createClientCompany(db, company)
	if err != nil {
		t.Fatalf("Failed to create client company: %v", err)
	}

	// Create delivered shipments with different delivery times
	baseTime := time.Now().AddDate(0, 0, -60)
	
	deliveryTimes := []struct {
		week int
		days int
	}{
		{1, 5},  // Week 1: 5 days
		{1, 7},  // Week 1: 7 days
		{2, 10}, // Week 2: 10 days
		{2, 12}, // Week 2: 12 days
		{3, 8},  // Week 3: 8 days
		{3, 6},  // Week 3: 6 days
		{4, 9},  // Week 4: 9 days
	}

	for i, dt := range deliveryTimes {
		pickupTime := baseTime.AddDate(0, 0, (dt.week-1)*7)
		deliveryTime := pickupTime.AddDate(0, 0, dt.days)
		
		shipment := &Shipment{
			ClientCompanyID:  company.ID,
			Status:           ShipmentStatusDelivered,
			JiraTicketNumber: "TEST-" + strconv.Itoa(i+1),
			PickedUpAt:       &pickupTime,
			DeliveredAt:      &deliveryTime,
		}
		err := createShipmentWithDate(db, shipment)
		if err != nil {
			t.Fatalf("Failed to create shipment: %v", err)
		}
	}

	// Test getting delivery time trends
	trends, err := GetDeliveryTimeTrends(db, 8) // Last 8 weeks
	if err != nil {
		t.Fatalf("GetDeliveryTimeTrends failed: %v", err)
	}

	// Verify we got data
	if len(trends) == 0 {
		t.Error("Expected trends data, got empty result")
	}

	// Verify at least some weeks have data
	hasData := false
	for _, trend := range trends {
		if trend.AverageDeliveryDays > 0 {
			hasData = true
			break
		}
	}
	if !hasData {
		t.Error("Expected at least some weeks with delivery data")
	}
}

// Helper function to create a shipment with custom date
func createShipmentWithDate(db *sql.DB, s *Shipment) error {
	// Note: CreatedAt should already be set by caller
	if s.CreatedAt.IsZero() {
		s.BeforeCreate()
	}
	if s.UpdatedAt.IsZero() {
		s.UpdatedAt = s.CreatedAt
	}
	
	query := `
		INSERT INTO shipments (
			client_company_id, software_engineer_id, status, jira_ticket_number, courier_name, 
			tracking_number, pickup_scheduled_date, picked_up_at, 
			arrived_warehouse_at, released_warehouse_at, delivered_at, 
			notes, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id
	`
	
	return db.QueryRow(
		query,
		s.ClientCompanyID, s.SoftwareEngineerID, s.Status, s.JiraTicketNumber, s.CourierName,
		s.TrackingNumber, s.PickupScheduledDate, s.PickedUpAt,
		s.ArrivedWarehouseAt, s.ReleasedWarehouseAt, s.DeliveredAt,
		s.Notes, s.CreatedAt, s.UpdatedAt,
	).Scan(&s.ID)
}

