package models

import (
	"testing"

	"github.com/yourusername/laptop-tracking-system/internal/database"
)

// TestGetAllLaptops tests retrieving all laptops with optional filtering
func TestGetAllLaptops(t *testing.T) {
	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	// Create test laptops
	laptops := []Laptop{
		{SerialNumber: "SN001", Brand: "Dell", Model: "Latitude", Status: LaptopStatusAvailable},
		{SerialNumber: "SN002", Brand: "HP", Model: "EliteBook", Status: LaptopStatusAvailable},
		{SerialNumber: "SN003", Brand: "Dell", Model: "XPS", Status: LaptopStatusDelivered},
		{SerialNumber: "SN004", Brand: "Lenovo", Model: "ThinkPad", Status: LaptopStatusAtWarehouse},
	}

	for i := range laptops {
		err := createLaptop(db, &laptops[i])
		if err != nil {
			t.Fatalf("Failed to create laptop: %v", err)
		}
	}

	// Test: Get all laptops
	result, err := GetAllLaptops(db, nil)
	if err != nil {
		t.Fatalf("GetAllLaptops failed: %v", err)
	}

	if len(result) != 4 {
		t.Errorf("Expected 4 laptops, got %d", len(result))
	}
}

// TestGetAllLaptopsWithFilter tests filtering laptops by status
func TestGetAllLaptopsWithFilter(t *testing.T) {
	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	// Create test laptops
	laptops := []Laptop{
		{SerialNumber: "SN001", Brand: "Dell", Status: LaptopStatusAvailable},
		{SerialNumber: "SN002", Brand: "HP", Status: LaptopStatusAvailable},
		{SerialNumber: "SN003", Brand: "Dell", Status: LaptopStatusDelivered},
	}

	for i := range laptops {
		err := createLaptop(db, &laptops[i])
		if err != nil {
			t.Fatalf("Failed to create laptop: %v", err)
		}
	}

	// Test: Filter by available status
	filter := &LaptopFilter{Status: LaptopStatusAvailable}
	result, err := GetAllLaptops(db, filter)
	if err != nil {
		t.Fatalf("GetAllLaptops with filter failed: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 available laptops, got %d", len(result))
	}

	// Verify all returned laptops have the correct status
	for _, laptop := range result {
		if laptop.Status != LaptopStatusAvailable {
			t.Errorf("Expected laptop with status available, got %s", laptop.Status)
		}
	}
}

// TestSearchLaptops tests searching laptops by serial number or brand
func TestSearchLaptops(t *testing.T) {
	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	// Create test laptops
	laptops := []Laptop{
		{SerialNumber: "ABC123", Brand: "Dell", Model: "Latitude", Status: LaptopStatusAvailable},
		{SerialNumber: "XYZ789", Brand: "HP", Model: "EliteBook", Status: LaptopStatusAvailable},
		{SerialNumber: "DEF456", Brand: "Dell", Model: "XPS", Status: LaptopStatusDelivered},
	}

	for i := range laptops {
		err := createLaptop(db, &laptops[i])
		if err != nil {
			t.Fatalf("Failed to create laptop: %v", err)
		}
	}

	// Test: Search by serial number
	result, err := SearchLaptops(db, "ABC")
	if err != nil {
		t.Fatalf("SearchLaptops failed: %v", err)
	}

	if len(result) != 1 {
		t.Errorf("Expected 1 laptop with 'ABC' in serial, got %d", len(result))
	}

	// Test: Search by brand
	result, err = SearchLaptops(db, "Dell")
	if err != nil {
		t.Fatalf("SearchLaptops failed: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 Dell laptops, got %d", len(result))
	}
}

// TestGetLaptopByID tests retrieving a laptop by its ID
func TestGetLaptopByID(t *testing.T) {
	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	// Create test laptop
	laptop := &Laptop{
		SerialNumber: "SN001",
		Brand:        "Dell",
		Model:        "Latitude",
		Status:       LaptopStatusAvailable,
	}
	err := createLaptop(db, laptop)
	if err != nil {
		t.Fatalf("Failed to create laptop: %v", err)
	}

	// Test: Get laptop by ID
	result, err := GetLaptopByID(db, laptop.ID)
	if err != nil {
		t.Fatalf("GetLaptopByID failed: %v", err)
	}

	if result.SerialNumber != laptop.SerialNumber {
		t.Errorf("Expected serial number %s, got %s", laptop.SerialNumber, result.SerialNumber)
	}

	if result.Brand != laptop.Brand {
		t.Errorf("Expected brand %s, got %s", laptop.Brand, result.Brand)
	}
}

// TestGetLaptopByIDNotFound tests retrieving a non-existent laptop
func TestGetLaptopByIDNotFound(t *testing.T) {
	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	// Test: Get non-existent laptop
	_, err := GetLaptopByID(db, 99999)
	if err == nil {
		t.Error("Expected error for non-existent laptop, got nil")
	}
}

// TestCreateLaptop tests creating a new laptop
func TestCreateLaptop(t *testing.T) {
	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	laptop := &Laptop{
		SerialNumber: "SN001",
		Brand:        "Dell",
		Model:        "Latitude 5520",
		Specs:        "i7, 16GB RAM, 512GB SSD",
		Status:       LaptopStatusAvailable,
	}

	// Test: Create laptop
	err := CreateLaptop(db, laptop)
	if err != nil {
		t.Fatalf("CreateLaptop failed: %v", err)
	}

	// Verify laptop was created with an ID
	if laptop.ID == 0 {
		t.Error("Expected laptop to have ID after creation")
	}

	// Verify laptop can be retrieved
	retrieved, err := GetLaptopByID(db, laptop.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve created laptop: %v", err)
	}

	if retrieved.SerialNumber != laptop.SerialNumber {
		t.Errorf("Expected serial number %s, got %s", laptop.SerialNumber, retrieved.SerialNumber)
	}
}

// TestCreateLaptopDuplicateSerial tests creating a laptop with duplicate serial number
func TestCreateLaptopDuplicateSerial(t *testing.T) {
	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	laptop1 := &Laptop{
		SerialNumber: "DUPLICATE001",
		Status:       LaptopStatusAvailable,
	}

	err := CreateLaptop(db, laptop1)
	if err != nil {
		t.Fatalf("Failed to create first laptop: %v", err)
	}

	// Test: Try to create laptop with same serial number
	laptop2 := &Laptop{
		SerialNumber: "DUPLICATE001",
		Status:       LaptopStatusAvailable,
	}

	err = CreateLaptop(db, laptop2)
	if err == nil {
		t.Error("Expected error for duplicate serial number, got nil")
	}
}

// TestUpdateLaptop tests updating an existing laptop
func TestUpdateLaptop(t *testing.T) {
	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	// Create test laptop
	laptop := &Laptop{
		SerialNumber: "SN001",
		Brand:        "Dell",
		Model:        "Latitude",
		Status:       LaptopStatusAvailable,
	}
	err := createLaptop(db, laptop)
	if err != nil {
		t.Fatalf("Failed to create laptop: %v", err)
	}

	// Test: Update laptop
	laptop.Model = "Latitude 5520"
	laptop.Specs = "i7, 16GB RAM"
	laptop.Status = LaptopStatusDelivered

	err = UpdateLaptop(db, laptop)
	if err != nil {
		t.Fatalf("UpdateLaptop failed: %v", err)
	}

	// Verify updates
	updated, err := GetLaptopByID(db, laptop.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve updated laptop: %v", err)
	}

	if updated.Model != "Latitude 5520" {
		t.Errorf("Expected model 'Latitude 5520', got %s", updated.Model)
	}

	if updated.Status != LaptopStatusDelivered {
		t.Errorf("Expected status delivered, got %s", updated.Status)
	}
}

// TestDeleteLaptop tests deleting a laptop
func TestDeleteLaptop(t *testing.T) {
	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	// Create test laptop
	laptop := &Laptop{
		SerialNumber: "SN001",
		Status:       LaptopStatusAvailable,
	}
	err := createLaptop(db, laptop)
	if err != nil {
		t.Fatalf("Failed to create laptop: %v", err)
	}

	// Test: Delete laptop
	err = DeleteLaptop(db, laptop.ID)
	if err != nil {
		t.Fatalf("DeleteLaptop failed: %v", err)
	}

	// Verify laptop is deleted
	_, err = GetLaptopByID(db, laptop.ID)
	if err == nil {
		t.Error("Expected error when retrieving deleted laptop, got nil")
	}
}

// TestGetLaptopsByStatus tests retrieving laptops by status
func TestGetLaptopsByStatus(t *testing.T) {
	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

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

	// Test: Get available laptops
	available, err := GetLaptopsByStatus(db, LaptopStatusAvailable)
	if err != nil {
		t.Fatalf("GetLaptopsByStatus failed: %v", err)
	}

	if len(available) != 2 {
		t.Errorf("Expected 2 available laptops, got %d", len(available))
	}

	// Test: Get delivered laptops
	delivered, err := GetLaptopsByStatus(db, LaptopStatusDelivered)
	if err != nil {
		t.Fatalf("GetLaptopsByStatus failed: %v", err)
	}

	if len(delivered) != 1 {
		t.Errorf("Expected 1 delivered laptop, got %d", len(delivered))
	}
}

