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
		{SerialNumber: "SN001", Brand: "Dell", Model: "Latitude", CPU: "Intel Core i5", RAMGB: "8GB", SSDGB: "256GB", Status: LaptopStatusAvailable},
		{SerialNumber: "SN002", Brand: "HP", Model: "EliteBook", CPU: "Intel Core i7", RAMGB: "16GB", SSDGB: "512GB", Status: LaptopStatusAvailable},
		{SerialNumber: "SN003", Brand: "Dell", Model: "XPS", CPU: "Intel Core i9", RAMGB: "32GB", SSDGB: "1TB", Status: LaptopStatusDelivered},
		{SerialNumber: "SN004", Brand: "Lenovo", Model: "ThinkPad", CPU: "AMD Ryzen 7", RAMGB: "16GB", SSDGB: "512GB", Status: LaptopStatusAtWarehouse},
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

	// Verify CPU field is populated for all laptops
	for i, laptop := range result {
		if laptop.CPU == "" {
			t.Errorf("Laptop %d (SN: %s) has empty CPU field", i, laptop.SerialNumber)
		}
	}
}

// TestGetAllLaptopsWithFilter tests filtering laptops by status
func TestGetAllLaptopsWithFilter(t *testing.T) {
	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	// Create test laptops
	laptops := []Laptop{
		{SerialNumber: "SN001", Brand: "Dell", Model: "Latitude", CPU: "Intel Core i5", RAMGB: "8GB", SSDGB: "256GB", Status: LaptopStatusAvailable},
		{SerialNumber: "SN002", Brand: "HP", Model: "EliteBook", CPU: "Intel Core i7", RAMGB: "16GB", SSDGB: "512GB", Status: LaptopStatusAvailable},
		{SerialNumber: "SN003", Brand: "Dell", Model: "XPS", CPU: "Intel Core i9", RAMGB: "32GB", SSDGB: "1TB", Status: LaptopStatusDelivered},
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
		{SerialNumber: "ABC123", Brand: "Dell", Model: "Latitude", CPU: "Intel Core i5", RAMGB: "8GB", SSDGB: "256GB", Status: LaptopStatusAvailable},
		{SerialNumber: "XYZ789", Brand: "HP", Model: "EliteBook", CPU: "Intel Core i7", RAMGB: "16GB", SSDGB: "512GB", Status: LaptopStatusAvailable},
		{SerialNumber: "DEF456", Brand: "Dell", Model: "XPS", CPU: "Intel Core i9", RAMGB: "32GB", SSDGB: "1TB", Status: LaptopStatusDelivered},
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
		CPU:          "Intel Core i5",
		RAMGB:        "8GB",
		SSDGB:        "256GB",
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

	// Create a client company first
	var clientID int64
	err := db.QueryRow(
		`INSERT INTO client_companies (name, created_at, updated_at)
		VALUES ($1, NOW(), NOW()) RETURNING id`,
		"Test Company",
	).Scan(&clientID)
	if err != nil {
		t.Fatalf("Failed to create client company: %v", err)
	}

	laptop := &Laptop{
		SerialNumber:    "SN001",
		Brand:           "Dell",
		Model:           "Latitude 5520",
		CPU:             "Intel Core i7",
		RAMGB:           "16GB",
		SSDGB:           "512GB",
		Status:          LaptopStatusAvailable,
		ClientCompanyID: &clientID,
	}

	// Test: Create laptop
	err = CreateLaptop(db, laptop)
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

	// Create a client company first
	var clientID int64
	err := db.QueryRow(
		`INSERT INTO client_companies (name, created_at, updated_at)
		VALUES ($1, NOW(), NOW()) RETURNING id`,
		"Test Company",
	).Scan(&clientID)
	if err != nil {
		t.Fatalf("Failed to create client company: %v", err)
	}

	laptop1 := &Laptop{
		SerialNumber:    "DUPLICATE001",
		Brand:           "Dell",
		Model:           "Dell Latitude 5520",
		CPU:             "Intel Core i5",
		RAMGB:           "16GB",
		SSDGB:           "512GB",
		Status:          LaptopStatusAvailable,
		ClientCompanyID: &clientID,
	}

	err = CreateLaptop(db, laptop1)
	if err != nil {
		t.Fatalf("Failed to create first laptop: %v", err)
	}

	// Test: Try to create laptop with same serial number
	laptop2 := &Laptop{
		SerialNumber:    "DUPLICATE001",
		Brand:           "HP",
		Model:           "Dell Latitude 5520",
		CPU:             "Intel Core i5",
		RAMGB:           "16GB",
		SSDGB:           "512GB",
		Status:          LaptopStatusAvailable,
		ClientCompanyID: &clientID,
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

	// Create a client company first
	var clientID int64
	err := db.QueryRow(
		`INSERT INTO client_companies (name, created_at, updated_at)
		VALUES ($1, NOW(), NOW()) RETURNING id`,
		"Test Company",
	).Scan(&clientID)
	if err != nil {
		t.Fatalf("Failed to create client company: %v", err)
	}

	// Create test laptop
	laptop := &Laptop{
		SerialNumber:    "SN001",
		Brand:           "Dell",
		Model:           "Latitude",
		CPU:             "Intel Core i5",
		RAMGB:           "8GB",
		SSDGB:           "256GB",
		Status:          LaptopStatusAvailable,
		ClientCompanyID: &clientID,
	}
	err = createLaptop(db, laptop)
	if err != nil {
		t.Fatalf("Failed to create laptop: %v", err)
	}

	// Test: Update laptop
	laptop.Model = "Latitude 5520"
	laptop.CPU = "Intel Core i7"
	laptop.RAMGB = "16GB"
	laptop.SSDGB = "256GB"
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
		Model:        "Test Model",
		CPU:          "Intel Core i5",
		RAMGB:        "8GB",
		SSDGB:        "256GB",
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
		{SerialNumber: "SN001", Model: "Model1", CPU: "Intel Core i5", RAMGB: "8GB", SSDGB: "256GB", Status: LaptopStatusAvailable},
		{SerialNumber: "SN002", Model: "Model2", CPU: "Intel Core i7", RAMGB: "16GB", SSDGB: "512GB", Status: LaptopStatusAvailable},
		{SerialNumber: "SN003", Model: "Model3", CPU: "Intel Core i9", RAMGB: "32GB", SSDGB: "1TB", Status: LaptopStatusDelivered},
		{SerialNumber: "SN004", Model: "Model4", CPU: "AMD Ryzen 7", RAMGB: "16GB", SSDGB: "512GB", Status: LaptopStatusAtWarehouse},
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

// TestGetAllLaptopsWithJoins tests retrieving laptops with client and engineer data
func TestGetAllLaptopsWithJoins(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	// Create test client company
	clientCompany := &ClientCompany{
		Name:        "TechCorp Inc",
		ContactInfo: "contact@techcorp.com",
	}
	err := CreateClientCompany(db, clientCompany)
	if err != nil {
		t.Fatalf("Failed to create client company: %v", err)
	}

	// Create test software engineer
	engineer := &SoftwareEngineer{
		Name:  "John Doe",
		Email: "john.doe@bairesdev.com",
	}
	err = CreateSoftwareEngineer(db, engineer)
	if err != nil {
		t.Fatalf("Failed to create software engineer: %v", err)
	}

	// Create laptop with SKU, client, and engineer
	laptop := &Laptop{
		SerialNumber:       "SN001",
		SKU:                "SKU-DELL-LAT-001",
		Brand:              "Dell",
		Model:              "Latitude 5520",
		CPU:                "Intel Core i7",
		RAMGB:              "16",
		SSDGB:              "512",
		Status:             LaptopStatusDelivered,
		ClientCompanyID:    &clientCompany.ID,
		SoftwareEngineerID: &engineer.ID,
	}
	err = CreateLaptop(db, laptop)
	if err != nil {
		t.Fatalf("Failed to create laptop: %v", err)
	}

	// Test: Get all laptops with joins
	result, err := GetAllLaptops(db, nil)
	if err != nil {
		t.Fatalf("GetAllLaptops failed: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("Expected 1 laptop, got %d", len(result))
	}

	// Verify new fields are populated
	retrieved := result[0]
	if retrieved.SKU != "SKU-DELL-LAT-001" {
		t.Errorf("Expected SKU 'SKU-DELL-LAT-001', got %s", retrieved.SKU)
	}

	if retrieved.ClientCompanyID == nil || *retrieved.ClientCompanyID != clientCompany.ID {
		t.Errorf("Expected ClientCompanyID %d, got %v", clientCompany.ID, retrieved.ClientCompanyID)
	}

	if retrieved.SoftwareEngineerID == nil || *retrieved.SoftwareEngineerID != engineer.ID {
		t.Errorf("Expected SoftwareEngineerID %d, got %v", engineer.ID, retrieved.SoftwareEngineerID)
	}

	// Verify joined data is populated
	if retrieved.ClientCompanyName != "TechCorp Inc" {
		t.Errorf("Expected ClientCompanyName 'TechCorp Inc', got %s", retrieved.ClientCompanyName)
	}

	if retrieved.SoftwareEngineerName != "John Doe" {
		t.Errorf("Expected SoftwareEngineerName 'John Doe', got %s", retrieved.SoftwareEngineerName)
	}
}

// TestGetLaptopByIDWithJoins tests retrieving a single laptop with client and engineer data
func TestGetLaptopByIDWithJoins(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	// Create test client company
	clientCompany := &ClientCompany{
		Name:        "Innovation Labs",
		ContactInfo: "info@innovationlabs.com",
	}
	err := CreateClientCompany(db, clientCompany)
	if err != nil {
		t.Fatalf("Failed to create client company: %v", err)
	}

	// Create test software engineer
	engineer := &SoftwareEngineer{
		Name:  "Jane Smith",
		Email: "jane.smith@bairesdev.com",
	}
	err = CreateSoftwareEngineer(db, engineer)
	if err != nil {
		t.Fatalf("Failed to create software engineer: %v", err)
	}

	// Create laptop with assignments
	laptop := &Laptop{
		SerialNumber:       "SN002",
		SKU:                "SKU-HP-ELITE-002",
		Brand:              "HP",
		Model:              "EliteBook 840",
		CPU:                "Intel Core i9",
		RAMGB:              "32",
		SSDGB:              "1024",
		Status:             LaptopStatusDelivered,
		ClientCompanyID:    &clientCompany.ID,
		SoftwareEngineerID: &engineer.ID,
	}
	err = CreateLaptop(db, laptop)
	if err != nil {
		t.Fatalf("Failed to create laptop: %v", err)
	}

	// Test: Get laptop by ID with joins
	retrieved, err := GetLaptopByID(db, laptop.ID)
	if err != nil {
		t.Fatalf("GetLaptopByID failed: %v", err)
	}

	// Verify new fields
	if retrieved.SKU != "SKU-HP-ELITE-002" {
		t.Errorf("Expected SKU 'SKU-HP-ELITE-002', got %s", retrieved.SKU)
	}

	// Verify joined data is populated
	if retrieved.ClientCompanyName != "Innovation Labs" {
		t.Errorf("Expected ClientCompanyName 'Innovation Labs', got %s", retrieved.ClientCompanyName)
	}

	if retrieved.SoftwareEngineerName != "Jane Smith" {
		t.Errorf("Expected SoftwareEngineerName 'Jane Smith', got %s", retrieved.SoftwareEngineerName)
	}
}

// 🟥 RED: Test that GetLaptopByID populates employee number when software engineer is assigned
func TestGetLaptopByIDPopulatesEmployeeNumber(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	// Create test client company
	clientCompany := &ClientCompany{
		Name:        "Test Company",
		ContactInfo: "test@company.com",
	}
	err := CreateClientCompany(db, clientCompany)
	if err != nil {
		t.Fatalf("Failed to create client company: %v", err)
	}

	// Create test software engineer with employee number
	engineer := &SoftwareEngineer{
		Name:          "John Doe",
		Email:         "john.doe@bairesdev.com",
		EmployeeNumber: "EMP12345",
	}
	err = CreateSoftwareEngineer(db, engineer)
	if err != nil {
		t.Fatalf("Failed to create software engineer: %v", err)
	}

	// Create laptop assigned to engineer
	laptop := &Laptop{
		SerialNumber:       "SN003",
		Brand:              "Dell",
		Model:              "Latitude 5520",
		CPU:                "Intel Core i7",
		RAMGB:              "16GB",
		SSDGB:              "512GB",
		Status:             LaptopStatusAvailable,
		ClientCompanyID:    &clientCompany.ID,
		SoftwareEngineerID: &engineer.ID,
	}
	err = CreateLaptop(db, laptop)
	if err != nil {
		t.Fatalf("Failed to create laptop: %v", err)
	}

	// Test: Get laptop by ID
	retrieved, err := GetLaptopByID(db, laptop.ID)
	if err != nil {
		t.Fatalf("GetLaptopByID failed: %v", err)
	}

	// Verify employee number is populated
	if retrieved.EmployeeID != "EMP12345" {
		t.Errorf("Expected EmployeeID 'EMP12345', got %s", retrieved.EmployeeID)
	}
}

// Helper functions removed - using actual model functions from client_company.go and software_engineer.go

// TestGetAllLaptopsHandlesNullFields tests that GetAllLaptops can handle NULL values in optional fields like brand and sku
func TestGetAllLaptopsHandlesNullFields(t *testing.T) {
	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	// Create a laptop with NULL brand and sku by inserting directly
	// (bypassing validation to simulate existing data in the database)
	// Note: model, ram_gb, ssd_gb are NOT NULL fields in the current schema
	query := `
		INSERT INTO laptops (serial_number, brand, model, ram_gb, ssd_gb, status, created_at, updated_at)
		VALUES ($1, NULL, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id
	`

	var laptopID int64
	err := db.QueryRow(query, "NULL_TEST_001", "Generic Model", "16", "512", LaptopStatusAvailable).Scan(&laptopID)
	if err != nil {
		t.Fatalf("Failed to create laptop with NULL optional fields: %v", err)
	}

	// Test: Get all laptops should not fail when encountering NULL values
	laptops, err := GetAllLaptops(db, nil)
	if err != nil {
		t.Fatalf("GetAllLaptops failed with NULL fields: %v", err)
	}

	// Find our test laptop
	found := false
	for _, laptop := range laptops {
		if laptop.ID == laptopID {
			found = true
			// Verify NULL fields are handled as empty strings
			if laptop.Brand != "" {
				t.Errorf("Expected empty brand for NULL value, got %s", laptop.Brand)
			}
			// Note: Model, RAMGB, SSDGB are now required fields, so they should have values
			// This test case should be updated if we're actually inserting NULL values
		}
	}

	if !found {
		t.Error("Test laptop with NULL fields not found in results")
	}
}
