package models

import (
	"testing"

	"github.com/yourusername/laptop-tracking-system/internal/database"
)

// TestGetAllLaptops_WarehouseRoleFilter tests that warehouse users only see relevant laptops
func TestGetAllLaptops_WarehouseRoleFilter(t *testing.T) {
	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	// Create test laptops with different statuses
	testLaptops := []Laptop{
		{SerialNumber: "SN-TRANSIT", Model: "Model A", RAMGB: "8GB", SSDGB: "256GB", Brand: "Dell", Status: LaptopStatusInTransitToWarehouse},
		{SerialNumber: "SN-AT-WH", Model: "Model B", RAMGB: "16GB", SSDGB: "512GB", Brand: "HP", Status: LaptopStatusAtWarehouse},
		{SerialNumber: "SN-AVAILABLE", Model: "Model C", RAMGB: "32GB", SSDGB: "1TB", Brand: "Lenovo", Status: LaptopStatusAvailable},
		{SerialNumber: "SN-TO-ENG", Model: "Model D", RAMGB: "16GB", SSDGB: "256GB", Brand: "Dell", Status: LaptopStatusInTransitToEngineer},
		{SerialNumber: "SN-DELIVERED", Model: "Model E", RAMGB: "8GB", SSDGB: "512GB", Brand: "HP", Status: LaptopStatusDelivered},
		{SerialNumber: "SN-RETIRED", Model: "Model F", RAMGB: "8GB", SSDGB: "256GB", Brand: "Lenovo", Status: LaptopStatusRetired},
	}

	// Insert test laptops
	for i := range testLaptops {
		err := createLaptop(db, &testLaptops[i])
		if err != nil {
			t.Fatalf("Failed to create test laptop: %v", err)
		}
	}

	// Test: Warehouse user filter
	filter := &LaptopFilter{
		UserRole: RoleWarehouse,
	}
	result, err := GetAllLaptops(db, filter)
	if err != nil {
		t.Fatalf("GetAllLaptops with warehouse role filter failed: %v", err)
	}

	// Warehouse users should only see: in_transit_to_warehouse, at_warehouse, available
	expectedCount := 3
	if len(result) != expectedCount {
		t.Errorf("Expected %d laptops for warehouse user, got %d", expectedCount, len(result))
	}

	// Verify each laptop has an allowed status
	allowedStatuses := map[LaptopStatus]bool{
		LaptopStatusInTransitToWarehouse: true,
		LaptopStatusAtWarehouse:          true,
		LaptopStatusAvailable:            true,
	}

	for _, laptop := range result {
		if !allowedStatuses[laptop.Status] {
			t.Errorf("Warehouse user should not see laptop with status %s (serial: %s)", laptop.Status, laptop.SerialNumber)
		}
	}

	// Verify each allowed status is present
	foundStatuses := make(map[LaptopStatus]bool)
	for _, laptop := range result {
		foundStatuses[laptop.Status] = true
	}

	if !foundStatuses[LaptopStatusInTransitToWarehouse] {
		t.Error("Expected to find at least one laptop with status 'in_transit_to_warehouse'")
	}
	if !foundStatuses[LaptopStatusAtWarehouse] {
		t.Error("Expected to find at least one laptop with status 'at_warehouse'")
	}
	if !foundStatuses[LaptopStatusAvailable] {
		t.Error("Expected to find at least one laptop with status 'available'")
	}
}

// TestGetAllLaptops_NonWarehouseRoleFilter tests that non-warehouse users see all laptops
func TestGetAllLaptops_NonWarehouseRoleFilter(t *testing.T) {
	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	// Create test laptops with different statuses
	testLaptops := []Laptop{
		{SerialNumber: "SN-1", Model: "Model A", RAMGB: "8GB", SSDGB: "256GB", Brand: "Dell", Status: LaptopStatusInTransitToWarehouse},
		{SerialNumber: "SN-2", Model: "Model B", RAMGB: "16GB", SSDGB: "512GB", Brand: "HP", Status: LaptopStatusAtWarehouse},
		{SerialNumber: "SN-3", Model: "Model C", RAMGB: "32GB", SSDGB: "1TB", Brand: "Lenovo", Status: LaptopStatusInTransitToEngineer},
		{SerialNumber: "SN-4", Model: "Model D", RAMGB: "16GB", SSDGB: "256GB", Brand: "Dell", Status: LaptopStatusDelivered},
	}

	for i := range testLaptops {
		err := createLaptop(db, &testLaptops[i])
		if err != nil {
			t.Fatalf("Failed to create test laptop: %v", err)
		}
	}

	// Test: Logistics user filter (should see all laptops)
	filter := &LaptopFilter{
		UserRole: RoleLogistics,
	}
	result, err := GetAllLaptops(db, filter)
	if err != nil {
		t.Fatalf("GetAllLaptops with logistics role filter failed: %v", err)
	}

	// Logistics users should see all laptops
	expectedCount := 4
	if len(result) != expectedCount {
		t.Errorf("Expected %d laptops for logistics user, got %d", expectedCount, len(result))
	}
}

// TestGetAllLaptops_WarehouseRoleWithStatusFilter tests combining role and status filters
func TestGetAllLaptops_WarehouseRoleWithStatusFilter(t *testing.T) {
	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	// Create test laptops
	testLaptops := []Laptop{
		{SerialNumber: "SN-WH-1", Model: "Model A", RAMGB: "8GB", SSDGB: "256GB", Brand: "Dell", Status: LaptopStatusAtWarehouse},
		{SerialNumber: "SN-WH-2", Model: "Model B", RAMGB: "16GB", SSDGB: "512GB", Brand: "HP", Status: LaptopStatusAtWarehouse},
		{SerialNumber: "SN-AVAIL", Model: "Model C", RAMGB: "32GB", SSDGB: "1TB", Brand: "Lenovo", Status: LaptopStatusAvailable},
		{SerialNumber: "SN-DELIVERED", Model: "Model D", RAMGB: "16GB", SSDGB: "256GB", Brand: "Dell", Status: LaptopStatusDelivered},
	}

	for i := range testLaptops {
		err := createLaptop(db, &testLaptops[i])
		if err != nil {
			t.Fatalf("Failed to create test laptop: %v", err)
		}
	}

	// Test: Warehouse user with status filter
	filter := &LaptopFilter{
		UserRole: RoleWarehouse,
		Status:   LaptopStatusAtWarehouse,
	}
	result, err := GetAllLaptops(db, filter)
	if err != nil {
		t.Fatalf("GetAllLaptops with combined filter failed: %v", err)
	}

	// Should only return at_warehouse laptops
	expectedCount := 2
	if len(result) != expectedCount {
		t.Errorf("Expected %d laptops with status 'at_warehouse', got %d", expectedCount, len(result))
	}

	for _, laptop := range result {
		if laptop.Status != LaptopStatusAtWarehouse {
			t.Errorf("Expected all laptops to have status 'at_warehouse', got %s", laptop.Status)
		}
	}
}

// TestGetAllLaptops_NoRoleFilter tests backward compatibility (no role specified)
func TestGetAllLaptops_NoRoleFilter(t *testing.T) {
	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	// Create test laptops
	testLaptops := []Laptop{
		{SerialNumber: "SN-1", Model: "Model A", RAMGB: "8GB", SSDGB: "256GB", Brand: "Dell", Status: LaptopStatusAtWarehouse},
		{SerialNumber: "SN-2", Model: "Model B", RAMGB: "16GB", SSDGB: "512GB", Brand: "HP", Status: LaptopStatusDelivered},
	}

	for i := range testLaptops {
		err := createLaptop(db, &testLaptops[i])
		if err != nil {
			t.Fatalf("Failed to create test laptop: %v", err)
		}
	}

	// Test: No role specified (backward compatibility - should see all)
	filter := &LaptopFilter{}
	result, err := GetAllLaptops(db, filter)
	if err != nil {
		t.Fatalf("GetAllLaptops with no role filter failed: %v", err)
	}

	// Should see all laptops when no role is specified
	expectedCount := 2
	if len(result) != expectedCount {
		t.Errorf("Expected %d laptops when no role specified, got %d", expectedCount, len(result))
	}
}

