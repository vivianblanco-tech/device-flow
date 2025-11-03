package models

import (
	"database/sql"
	"testing"
	"time"

	"github.com/yourusername/laptop-tracking-system/internal/database"
)

// TestGetShipmentCountsByStatus tests counting shipments grouped by status
func TestGetShipmentCountsByStatus(t *testing.T) {
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

	// Create test shipments with different statuses
	shipments := []Shipment{
		{ClientCompanyID: company.ID, Status: ShipmentStatusPendingPickup},
		{ClientCompanyID: company.ID, Status: ShipmentStatusPendingPickup},
		{ClientCompanyID: company.ID, Status: ShipmentStatusAtWarehouse},
		{ClientCompanyID: company.ID, Status: ShipmentStatusInTransitToEngineer},
		{ClientCompanyID: company.ID, Status: ShipmentStatusDelivered},
		{ClientCompanyID: company.ID, Status: ShipmentStatusDelivered},
		{ClientCompanyID: company.ID, Status: ShipmentStatusDelivered},
	}

	for i := range shipments {
		err := createShipment(db, &shipments[i])
		if err != nil {
			t.Fatalf("Failed to create shipment: %v", err)
		}
	}

	// Test counting shipments by status
	counts, err := GetShipmentCountsByStatus(db)
	if err != nil {
		t.Fatalf("GetShipmentCountsByStatus failed: %v", err)
	}

	// Verify counts
	expectedCounts := map[ShipmentStatus]int{
		ShipmentStatusPendingPickup:       2,
		ShipmentStatusAtWarehouse:         1,
		ShipmentStatusInTransitToEngineer: 1,
		ShipmentStatusDelivered:           3,
	}

	for status, expectedCount := range expectedCounts {
		if counts[status] != expectedCount {
			t.Errorf("Expected %d shipments with status %s, got %d",
				expectedCount, status, counts[status])
		}
	}
}

// TestGetTotalShipmentCount tests counting all shipments
func TestGetTotalShipmentCount(t *testing.T) {
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

	// Create 5 test shipments
	for i := 0; i < 5; i++ {
		shipment := &Shipment{
			ClientCompanyID: company.ID,
			Status:          ShipmentStatusPendingPickup,
		}
		err := createShipment(db, shipment)
		if err != nil {
			t.Fatalf("Failed to create shipment: %v", err)
		}
	}

	// Test total count
	count, err := GetTotalShipmentCount(db)
	if err != nil {
		t.Fatalf("GetTotalShipmentCount failed: %v", err)
	}

	if count != 5 {
		t.Errorf("Expected 5 shipments, got %d", count)
	}
}

// TestGetAverageDeliveryTime tests calculating average delivery time
func TestGetAverageDeliveryTime(t *testing.T) {
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

	// Create test shipments with delivery times
	baseTime := time.Now().Add(-10 * 24 * time.Hour)

	shipments := []struct {
		pickupDays   int
		deliveryDays int
	}{
		{0, 5},  // 5 days
		{0, 10}, // 10 days
		{0, 15}, // 15 days
	}

	for _, s := range shipments {
		pickupTime := baseTime.Add(time.Duration(s.pickupDays) * 24 * time.Hour)
		deliveryTime := baseTime.Add(time.Duration(s.deliveryDays) * 24 * time.Hour)

		shipment := &Shipment{
			ClientCompanyID: company.ID,
			Status:          ShipmentStatusDelivered,
			PickedUpAt:      &pickupTime,
			DeliveredAt:     &deliveryTime,
		}
		err := createShipment(db, shipment)
		if err != nil {
			t.Fatalf("Failed to create shipment: %v", err)
		}
	}

	// Test average delivery time
	avgDays, err := GetAverageDeliveryTime(db)
	if err != nil {
		t.Fatalf("GetAverageDeliveryTime failed: %v", err)
	}

	// Average should be 10 days
	expectedAvg := 10.0
	if avgDays < expectedAvg-0.1 || avgDays > expectedAvg+0.1 {
		t.Errorf("Expected average delivery time of %.1f days, got %.1f", expectedAvg, avgDays)
	}
}

// TestGetInTransitShipmentCount tests counting shipments in transit
func TestGetInTransitShipmentCount(t *testing.T) {
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
		ShipmentStatusInTransitToWarehouse,
		ShipmentStatusInTransitToEngineer,
		ShipmentStatusInTransitToWarehouse,
		ShipmentStatusAtWarehouse, // Not in transit
		ShipmentStatusDelivered,   // Not in transit
	}

	for _, status := range statuses {
		shipment := &Shipment{
			ClientCompanyID: company.ID,
			Status:          status,
		}
		err := createShipment(db, shipment)
		if err != nil {
			t.Fatalf("Failed to create shipment: %v", err)
		}
	}

	// Test in-transit count
	count, err := GetInTransitShipmentCount(db)
	if err != nil {
		t.Fatalf("GetInTransitShipmentCount failed: %v", err)
	}

	// Should be 3 (2 to warehouse + 1 to engineer)
	if count != 3 {
		t.Errorf("Expected 3 in-transit shipments, got %d", count)
	}
}

// TestGetPendingPickupCount tests counting shipments pending pickup
func TestGetPendingPickupCount(t *testing.T) {
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
		ShipmentStatusPickedUpFromClient, // Not pending
		ShipmentStatusAtWarehouse,        // Not pending
	}

	for _, status := range statuses {
		shipment := &Shipment{
			ClientCompanyID: company.ID,
			Status:          status,
		}
		err := createShipment(db, shipment)
		if err != nil {
			t.Fatalf("Failed to create shipment: %v", err)
		}
	}

	// Test pending pickup count
	count, err := GetPendingPickupCount(db)
	if err != nil {
		t.Fatalf("GetPendingPickupCount failed: %v", err)
	}

	if count != 3 {
		t.Errorf("Expected 3 pending pickup shipments, got %d", count)
	}
}

// TestGetLaptopCountsByStatus tests counting laptops grouped by status
func TestGetLaptopCountsByStatus(t *testing.T) {
	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	// Create test laptops with different statuses
	laptops := []Laptop{
		{SerialNumber: "SN001", Status: LaptopStatusAvailable},
		{SerialNumber: "SN002", Status: LaptopStatusAvailable},
		{SerialNumber: "SN003", Status: LaptopStatusAvailable},
		{SerialNumber: "SN004", Status: LaptopStatusAtWarehouse},
		{SerialNumber: "SN005", Status: LaptopStatusInTransitToEngineer},
		{SerialNumber: "SN006", Status: LaptopStatusDelivered},
		{SerialNumber: "SN007", Status: LaptopStatusDelivered},
	}

	for i := range laptops {
		err := createLaptop(db, &laptops[i])
		if err != nil {
			t.Fatalf("Failed to create laptop: %v", err)
		}
	}

	// Test counting laptops by status
	counts, err := GetLaptopCountsByStatus(db)
	if err != nil {
		t.Fatalf("GetLaptopCountsByStatus failed: %v", err)
	}

	// Verify counts
	expectedCounts := map[LaptopStatus]int{
		LaptopStatusAvailable:           3,
		LaptopStatusAtWarehouse:         1,
		LaptopStatusInTransitToEngineer: 1,
		LaptopStatusDelivered:           2,
	}

	for status, expectedCount := range expectedCounts {
		if counts[status] != expectedCount {
			t.Errorf("Expected %d laptops with status %s, got %d",
				expectedCount, status, counts[status])
		}
	}
}

// TestGetAvailableLaptopCount tests counting available laptops
func TestGetAvailableLaptopCount(t *testing.T) {
	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	// Create test laptops
	statuses := []LaptopStatus{
		LaptopStatusAvailable,
		LaptopStatusAvailable,
		LaptopStatusAvailable,
		LaptopStatusDelivered,
		LaptopStatusRetired,
	}

	for i, status := range statuses {
		laptop := &Laptop{
			SerialNumber: "SN" + string(rune('0'+i)),
			Status:       status,
		}
		err := createLaptop(db, laptop)
		if err != nil {
			t.Fatalf("Failed to create laptop: %v", err)
		}
	}

	// Test available count
	count, err := GetAvailableLaptopCount(db)
	if err != nil {
		t.Fatalf("GetAvailableLaptopCount failed: %v", err)
	}

	if count != 3 {
		t.Errorf("Expected 3 available laptops, got %d", count)
	}
}

// TestGetDashboardStats tests retrieving all dashboard statistics in one call
func TestGetDashboardStats(t *testing.T) {
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

	// Create test shipments with different statuses and delivery times
	baseTime := time.Now().Add(-20 * 24 * time.Hour)

	// Shipment 1: Delivered (took 5 days)
	pickup1 := baseTime
	delivery1 := baseTime.Add(5 * 24 * time.Hour)
	shipment1 := &Shipment{
		ClientCompanyID: company.ID,
		Status:          ShipmentStatusDelivered,
		PickedUpAt:      &pickup1,
		DeliveredAt:     &delivery1,
	}
	err = createShipment(db, shipment1)
	if err != nil {
		t.Fatalf("Failed to create shipment 1: %v", err)
	}

	// Shipment 2: Delivered (took 15 days)
	pickup2 := baseTime.Add(-10 * 24 * time.Hour)
	delivery2 := baseTime.Add(5 * 24 * time.Hour)
	shipment2 := &Shipment{
		ClientCompanyID: company.ID,
		Status:          ShipmentStatusDelivered,
		PickedUpAt:      &pickup2,
		DeliveredAt:     &delivery2,
	}
	err = createShipment(db, shipment2)
	if err != nil {
		t.Fatalf("Failed to create shipment 2: %v", err)
	}

	// Shipment 3: Pending pickup
	shipment3 := &Shipment{
		ClientCompanyID: company.ID,
		Status:          ShipmentStatusPendingPickup,
	}
	err = createShipment(db, shipment3)
	if err != nil {
		t.Fatalf("Failed to create shipment 3: %v", err)
	}

	// Shipment 4: In transit to warehouse
	shipment4 := &Shipment{
		ClientCompanyID: company.ID,
		Status:          ShipmentStatusInTransitToWarehouse,
	}
	err = createShipment(db, shipment4)
	if err != nil {
		t.Fatalf("Failed to create shipment 4: %v", err)
	}

	// Shipment 5: In transit to engineer
	shipment5 := &Shipment{
		ClientCompanyID: company.ID,
		Status:          ShipmentStatusInTransitToEngineer,
	}
	err = createShipment(db, shipment5)
	if err != nil {
		t.Fatalf("Failed to create shipment 5: %v", err)
	}

	// Create test laptops
	laptops := []Laptop{
		{SerialNumber: "SN001", Status: LaptopStatusAvailable},
		{SerialNumber: "SN002", Status: LaptopStatusAvailable},
		{SerialNumber: "SN003", Status: LaptopStatusDelivered},
		{SerialNumber: "SN004", Status: LaptopStatusAtWarehouse},
	}
	for i := range laptops {
		err := createLaptop(db, &laptops[i])
		if err != nil {
			t.Fatalf("Failed to create laptop: %v", err)
		}
	}

	// Get dashboard stats
	stats, err := GetDashboardStats(db)
	if err != nil {
		t.Fatalf("GetDashboardStats failed: %v", err)
	}

	// Verify total shipments
	if stats.TotalShipments != 5 {
		t.Errorf("Expected 5 total shipments, got %d", stats.TotalShipments)
	}

	// Verify pending pickups
	if stats.PendingPickups != 1 {
		t.Errorf("Expected 1 pending pickup, got %d", stats.PendingPickups)
	}

	// Verify in transit
	if stats.InTransit != 2 {
		t.Errorf("Expected 2 in-transit shipments, got %d", stats.InTransit)
	}

	// Verify delivered
	if stats.Delivered != 2 {
		t.Errorf("Expected 2 delivered shipments, got %d", stats.Delivered)
	}

	// Verify average delivery time (should be 10 days: (5+15)/2)
	expectedAvg := 10.0
	if stats.AvgDeliveryDays < expectedAvg-0.1 || stats.AvgDeliveryDays > expectedAvg+0.1 {
		t.Errorf("Expected average delivery time of %.1f days, got %.1f",
			expectedAvg, stats.AvgDeliveryDays)
	}

	// Verify shipments by status
	if stats.ShipmentsByStatus[ShipmentStatusDelivered] != 2 {
		t.Errorf("Expected 2 delivered in status map, got %d",
			stats.ShipmentsByStatus[ShipmentStatusDelivered])
	}
	if stats.ShipmentsByStatus[ShipmentStatusPendingPickup] != 1 {
		t.Errorf("Expected 1 pending pickup in status map, got %d",
			stats.ShipmentsByStatus[ShipmentStatusPendingPickup])
	}

	// Verify available laptops
	if stats.AvailableLaptops != 2 {
		t.Errorf("Expected 2 available laptops, got %d", stats.AvailableLaptops)
	}

	// Verify laptops by status
	if stats.LaptopsByStatus[LaptopStatusAvailable] != 2 {
		t.Errorf("Expected 2 available laptops in status map, got %d",
			stats.LaptopsByStatus[LaptopStatusAvailable])
	}
	if stats.LaptopsByStatus[LaptopStatusDelivered] != 1 {
		t.Errorf("Expected 1 delivered laptop in status map, got %d",
			stats.LaptopsByStatus[LaptopStatusDelivered])
	}
	if stats.LaptopsByStatus[LaptopStatusAtWarehouse] != 1 {
		t.Errorf("Expected 1 warehouse laptop in status map, got %d",
			stats.LaptopsByStatus[LaptopStatusAtWarehouse])
	}
}

// Helper function to create a client company in the test database
func createClientCompany(db *sql.DB, c *ClientCompany) error {
	c.BeforeCreate()

	query := `
		INSERT INTO client_companies (name, contact_info, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	return db.QueryRow(query, c.Name, c.ContactInfo, c.CreatedAt, c.UpdatedAt).Scan(&c.ID)
}

// Helper function to create a shipment in the test database
func createShipment(db *sql.DB, s *Shipment) error {
	s.BeforeCreate()

	query := `
		INSERT INTO shipments (
			client_company_id, software_engineer_id, status, courier_name, 
			tracking_number, pickup_scheduled_date, picked_up_at, 
			arrived_warehouse_at, released_warehouse_at, delivered_at, 
			notes, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id
	`

	return db.QueryRow(
		query,
		s.ClientCompanyID, s.SoftwareEngineerID, s.Status, s.CourierName,
		s.TrackingNumber, s.PickupScheduledDate, s.PickedUpAt,
		s.ArrivedWarehouseAt, s.ReleasedWarehouseAt, s.DeliveredAt,
		s.Notes, s.CreatedAt, s.UpdatedAt,
	).Scan(&s.ID)
}

// Helper function to create a laptop in the test database
func createLaptop(db *sql.DB, l *Laptop) error {
	l.BeforeCreate()

	query := `
		INSERT INTO laptops (serial_number, brand, model, specs, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	return db.QueryRow(
		query,
		l.SerialNumber, l.Brand, l.Model, l.Specs, l.Status, l.CreatedAt, l.UpdatedAt,
	).Scan(&l.ID)
}
