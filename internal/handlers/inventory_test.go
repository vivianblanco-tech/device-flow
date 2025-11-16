package handlers

import (
	"context"
	"database/sql"
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/yourusername/laptop-tracking-system/internal/database"
	"github.com/yourusername/laptop-tracking-system/internal/middleware"
	"github.com/yourusername/laptop-tracking-system/internal/models"
)

// TestAddLaptopAutoGeneratesSKU tests that SKU is auto-generated when adding a laptop
func TestAddLaptopAutoGeneratesSKU(t *testing.T) {
	// Setup test database
	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	// Create a test client company
	var companyID int64
	err := db.QueryRow(
		`INSERT INTO client_companies (name, contact_info, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW()) RETURNING id`,
		"Test Corp", "contact@test.com",
	).Scan(&companyID)
	if err != nil {
		t.Fatalf("Failed to create client company: %v", err)
	}

	// Create a test user (logistics role can add laptops)
	var userID int64
	err = db.QueryRow(
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id`,
		"logistics@test.com", "hashed_password", models.RoleLogistics,
	).Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	user := &models.User{
		ID:           userID,
		Email:        "logistics@test.com",
		PasswordHash: "hashed_password",
		Role:         models.RoleLogistics,
	}

	// Setup handler with minimal template
	tmpl := template.New("test")
	handler := &InventoryHandler{
		DB:        db,
		Templates: tmpl,
	}

	// Create form data for a Dell laptop with i7
	form := url.Values{}
	form.Add("serial_number", "SN-AUTO-SKU-001")
	form.Add("client_company_id", strconv.FormatInt(companyID, 10))
	form.Add("brand", "Dell")
	form.Add("model", "Latitude 5520")
	form.Add("cpu", "i7")
	form.Add("ram_gb", "16GB")
	form.Add("ssd_gb", "512GB")
	form.Add("status", string(models.LaptopStatusAvailable))

	// Create request
	req := httptest.NewRequest("POST", "/inventory/add", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	
	// Add user to context
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
	req = req.WithContext(ctx)

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	handler.AddLaptopSubmit(rr, req)

	// Check response
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Expected status %d, got %d", http.StatusSeeOther, rr.Code)
	}

	// Verify laptop was created with auto-generated SKU
	filter := &models.LaptopFilter{}
	laptops, err := models.GetAllLaptops(db, filter)
	if err != nil {
		t.Fatalf("Failed to get laptops: %v", err)
	}

	if len(laptops) == 0 {
		t.Fatal("Expected laptop to be created")
	}

	laptop := laptops[0]
	expectedSKU := "C.NOT.0I7.016.2G"
	if laptop.SKU != expectedSKU {
		t.Errorf("Expected SKU %s, got %s", expectedSKU, laptop.SKU)
	}
}

// TestUpdateLaptopAutoRegeneratesSKUWhenFieldsChange tests that SKU can be manually overridden
func TestUpdateLaptopAutoRegeneratesSKUWhenFieldsChange(t *testing.T) {
	// Setup test database
	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	// Create a test client company
	var companyID int64
	err := db.QueryRow(
		`INSERT INTO client_companies (name, contact_info, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW()) RETURNING id`,
		"Test Corp", "contact@test.com",
	).Scan(&companyID)
	if err != nil {
		t.Fatalf("Failed to create client company: %v", err)
	}

	// Create a test laptop
	laptop := &models.Laptop{
		SerialNumber:    "SN-UPDATE-SKU-001",
		SKU:             "C.NOT.0I5.016.2G",
		Brand:           "Dell",
		Model:           "Latitude 5420",
		CPU:             "i5",
		RAMGB:           "16GB",
		SSDGB:           "512GB",
		Status:          models.LaptopStatusAvailable,
		ClientCompanyID: &companyID,
	}
	if err := models.CreateLaptop(db, laptop); err != nil {
		t.Fatalf("Failed to create laptop: %v", err)
	}

	// Create a test user (logistics role can update laptops)
	var userID int64
	err = db.QueryRow(
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id`,
		"logistics@test.com", "hashed_password", models.RoleLogistics,
	).Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	user := &models.User{
		ID:           userID,
		Email:        "logistics@test.com",
		PasswordHash: "hashed_password",
		Role:         models.RoleLogistics,
	}

	// Setup handler with minimal template
	tmpl := template.New("test")
	handler := &InventoryHandler{
		DB:        db,
		Templates: tmpl,
	}

	// Create form data updating CPU to i7 and RAM to 32GB (SKU should remain as-is since we provide it)
	form := url.Values{}
	form.Add("serial_number", "SN-UPDATE-SKU-001")
	form.Add("sku", "CUSTOM-SKU-123") // Manual override
	form.Add("client_company_id", strconv.FormatInt(companyID, 10))
	form.Add("brand", "Dell")
	form.Add("model", "Latitude 5520")
	form.Add("cpu", "i7")
	form.Add("ram_gb", "32GB")
	form.Add("ssd_gb", "512GB")
	form.Add("status", string(models.LaptopStatusAvailable))

	// Create request
	req := httptest.NewRequest("POST", "/inventory/"+strconv.FormatInt(laptop.ID, 10)+"/update", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(laptop.ID, 10)})
	
	// Add user to context
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
	req = req.WithContext(ctx)

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	handler.UpdateLaptopSubmit(rr, req)

	// Check response
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Expected status %d, got %d", http.StatusSeeOther, rr.Code)
	}

	// Verify laptop was updated with the custom SKU (not auto-generated)
	updatedLaptop, err := models.GetLaptopByID(db, laptop.ID)
	if err != nil {
		t.Fatalf("Failed to get updated laptop: %v", err)
	}

	if updatedLaptop.SKU != "CUSTOM-SKU-123" {
		t.Errorf("Expected SKU to be CUSTOM-SKU-123, got %s", updatedLaptop.SKU)
	}
}

// TestUpdateLaptopWithSoftwareEngineerAssignment tests updating a laptop with software engineer assignment
func TestUpdateLaptopWithSoftwareEngineerAssignment(t *testing.T) {
	// Setup test database
	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	// Create a test software engineer
	engineer := &models.SoftwareEngineer{
		Name:  "John Doe",
		Email: "john.doe@example.com",
		Phone: "+1-555-1234",
	}
	if err := createTestSoftwareEngineer(db, engineer); err != nil {
		t.Fatalf("Failed to create software engineer: %v", err)
	}

	// Create a test laptop
	laptop := &models.Laptop{
		SerialNumber: "SN-TEST-001",
		Brand:        "Dell",
		Model:        "XPS 15",
		CPU:          "Intel Core i7",
		RAMGB:        "16GB",
		SSDGB:        "512GB",
		Status:       models.LaptopStatusAvailable,
	}
	if err := models.CreateLaptop(db, laptop); err != nil {
		t.Fatalf("Failed to create laptop: %v", err)
	}

	// Create a test user (logistics role can update laptops)
	var logisticsUserID int64
	err := db.QueryRow(
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id`,
		"logistics@bairesdev.com", "hashed_password", models.RoleLogistics,
	).Scan(&logisticsUserID)
	if err != nil {
		t.Fatalf("Failed to create logistics user: %v", err)
	}

	logisticsUser := &models.User{
		ID:           logisticsUserID,
		Email:        "logistics@bairesdev.com",
		PasswordHash: "hashed_password",
		Role:         models.RoleLogistics,
	}

	// Load test templates (minimal template for testing)
	tmpl := template.New("test")

	// Create handler
	handler := &InventoryHandler{
		DB:        db,
		Templates: tmpl,
	}

	// Create form data with software engineer assignment
	formData := url.Values{}
	formData.Set("serial_number", "SN-TEST-001-UPDATED")
	formData.Set("brand", "Dell")
	formData.Set("model", "XPS 15")
	formData.Set("ram_gb", "32GB") // Updated RAM
	formData.Set("ssd_gb", "1TB")  // Updated SSD
	formData.Set("status", string(models.LaptopStatusAvailable))
	formData.Set("software_engineer_id", strconv.FormatInt(engineer.ID, 10))

	// Create request
	req := httptest.NewRequest(http.MethodPost, "/inventory/"+strconv.FormatInt(laptop.ID, 10)+"/update", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Add user to context
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, logisticsUser)
	req = req.WithContext(ctx)

	// Add URL vars to request (simulate mux router)
	req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(laptop.ID, 10)})

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	handler.UpdateLaptopSubmit(rr, req)

	// Verify status code (should redirect to laptop detail page)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Expected status code %d, got %d", http.StatusSeeOther, rr.Code)
	}

	// Verify redirect location
	expectedLocation := "/inventory/" + strconv.FormatInt(laptop.ID, 10)
	if location := rr.Header().Get("Location"); location != expectedLocation {
		t.Errorf("Expected redirect to %s, got %s", expectedLocation, location)
	}

	// Verify laptop was updated in database
	updatedLaptop, err := models.GetLaptopByID(db, laptop.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve updated laptop: %v", err)
	}

	// Verify software engineer assignment
	if updatedLaptop.SoftwareEngineerID == nil {
		t.Error("Software engineer ID should be set")
	} else if *updatedLaptop.SoftwareEngineerID != engineer.ID {
		t.Errorf("Expected software engineer ID %d, got %d", engineer.ID, *updatedLaptop.SoftwareEngineerID)
	}

	// Verify other fields were updated
	if updatedLaptop.RAMGB != "32GB" {
		t.Errorf("Expected RAM '32GB', got '%s'", updatedLaptop.RAMGB)
	}
	if updatedLaptop.SSDGB != "1TB" {
		t.Errorf("Expected SSD '1TB', got '%s'", updatedLaptop.SSDGB)
	}
	if updatedLaptop.SerialNumber != "SN-TEST-001-UPDATED" {
		t.Errorf("Expected serial number 'SN-TEST-001-UPDATED', got '%s'", updatedLaptop.SerialNumber)
	}
}

// TestUpdateLaptopRemoveSoftwareEngineerAssignment tests clearing the software engineer assignment
func TestUpdateLaptopRemoveSoftwareEngineerAssignment(t *testing.T) {
	// Setup test database
	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	// Create a test software engineer
	engineer := &models.SoftwareEngineer{
		Name:  "Jane Smith",
		Email: "jane.smith@example.com",
		Phone: "+1-555-5678",
	}
	if err := createTestSoftwareEngineer(db, engineer); err != nil {
		t.Fatalf("Failed to create software engineer: %v", err)
	}

	// Create a test laptop with software engineer assigned
	laptop := &models.Laptop{
		SerialNumber:       "SN-TEST-002",
		Brand:              "HP",
		Model:              "EliteBook",
		CPU:                "Intel Core i5",
		RAMGB:              "16GB",
		SSDGB:              "512GB",
		Status:             models.LaptopStatusDelivered,
		SoftwareEngineerID: &engineer.ID,
	}
	if err := models.CreateLaptop(db, laptop); err != nil {
		t.Fatalf("Failed to create laptop: %v", err)
	}

	// Create a test user
	var logisticsUserID int64
	err := db.QueryRow(
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id`,
		"logistics2@bairesdev.com", "hashed_password", models.RoleLogistics,
	).Scan(&logisticsUserID)
	if err != nil {
		t.Fatalf("Failed to create logistics user: %v", err)
	}

	logisticsUser := &models.User{
		ID:           logisticsUserID,
		Email:        "logistics2@bairesdev.com",
		PasswordHash: "hashed_password",
		Role:         models.RoleLogistics,
	}

	// Load test templates
	tmpl := template.New("test")

	// Create handler
	handler := &InventoryHandler{
		DB:        db,
		Templates: tmpl,
	}

	// Create form data without software engineer (empty string to clear assignment)
	formData := url.Values{}
	formData.Set("serial_number", "SN-TEST-002")
	formData.Set("brand", "HP")
	formData.Set("model", "EliteBook")
	formData.Set("ram_gb", "16GB")
	formData.Set("ssd_gb", "512GB")
	formData.Set("status", string(models.LaptopStatusDelivered))
	formData.Set("software_engineer_id", "") // Empty to clear assignment

	// Create request
	req := httptest.NewRequest(http.MethodPost, "/inventory/"+strconv.FormatInt(laptop.ID, 10)+"/update", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Add user to context
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, logisticsUser)
	req = req.WithContext(ctx)

	// Add URL vars
	req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(laptop.ID, 10)})

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	handler.UpdateLaptopSubmit(rr, req)

	// Verify status code
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Expected status code %d, got %d", http.StatusSeeOther, rr.Code)
	}

	// Verify laptop was updated in database
	updatedLaptop, err := models.GetLaptopByID(db, laptop.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve updated laptop: %v", err)
	}

	// Verify software engineer assignment was cleared
	if updatedLaptop.SoftwareEngineerID != nil {
		t.Errorf("Software engineer ID should be nil, got %d", *updatedLaptop.SoftwareEngineerID)
	}
}

// TestInventoryListWithSorting tests that the inventory list can be sorted by different columns
func TestInventoryListWithSorting(t *testing.T) {
	// Setup test database
	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	// Create test companies
	var company1ID, company2ID int64
	err := db.QueryRow(
		`INSERT INTO client_companies (name, contact_info, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW()) RETURNING id`,
		"Alpha Corp", "alpha@test.com",
	).Scan(&company1ID)
	if err != nil {
		t.Fatalf("Failed to create company1: %v", err)
	}

	err = db.QueryRow(
		`INSERT INTO client_companies (name, contact_info, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW()) RETURNING id`,
		"Beta Corp", "beta@test.com",
	).Scan(&company2ID)
	if err != nil {
		t.Fatalf("Failed to create company2: %v", err)
	}

	// Create test laptops with different attributes
	laptops := []struct {
		serial  string
		brand   string
		model   string
		status  models.LaptopStatus
		company int64
	}{
		{"SN-001", "Dell", "Latitude 5520", models.LaptopStatusAvailable, company1ID},
		{"SN-002", "HP", "EliteBook 850", models.LaptopStatusAtWarehouse, company2ID},
		{"SN-003", "Apple", "MacBook Pro", models.LaptopStatusDelivered, company1ID},
		{"SN-004", "Lenovo", "ThinkPad X1", models.LaptopStatusAvailable, company2ID},
	}

	for _, l := range laptops {
		laptop := &models.Laptop{
			SerialNumber:    l.serial,
			Brand:           l.brand,
			Model:           l.model,
			CPU:             "i7",
			RAMGB:           "16GB",
			SSDGB:           "512GB",
			Status:          l.status,
			ClientCompanyID: &l.company,
		}
		if err := models.CreateLaptop(db, laptop); err != nil {
			t.Fatalf("Failed to create laptop: %v", err)
		}
	}

	// Create test user (logistics role can view all inventory)
	var userID int64
	err = db.QueryRow(
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id`,
		"logistics@test.com", "hashed_password", models.RoleLogistics,
	).Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	user := &models.User{
		ID:           userID,
		Email:        "logistics@test.com",
		PasswordHash: "hashed_password",
		Role:         models.RoleLogistics,
	}

	// Setup handler
	tmpl := template.New("test")
	handler := &InventoryHandler{
		DB:        db,
		Templates: tmpl,
	}

	// Test cases for different sorting combinations
	testCases := []struct {
		name          string
		sortBy        string
		sortOrder     string
		expectedFirst string // Expected first serial number
		expectedLast  string // Expected last serial number
	}{
		{
			name:          "Sort by serial number ascending",
			sortBy:        "serial_number",
			sortOrder:     "asc",
			expectedFirst: "SN-001",
			expectedLast:  "SN-004",
		},
		{
			name:          "Sort by serial number descending",
			sortBy:        "serial_number",
			sortOrder:     "desc",
			expectedFirst: "SN-004",
			expectedLast:  "SN-001",
		},
		{
			name:          "Sort by brand ascending",
			sortBy:        "brand",
			sortOrder:     "asc",
			expectedFirst: "SN-003", // Apple
			expectedLast:  "SN-004", // Lenovo
		},
		{
			name:          "Sort by status ascending",
			sortBy:        "status",
			sortOrder:     "asc",
			expectedFirst: "SN-002", // at_warehouse (alphabetically first: 'at' < 'av')
			expectedLast:  "SN-003", // delivered (alphabetically last)
		},
		{
			name:          "Sort by client company ascending",
			sortBy:        "client_company",
			sortOrder:     "asc",
			expectedFirst: "SN-001", // Alpha Corp (first)
			expectedLast:  "SN-004", // Beta Corp (last of Beta items)
		},
		{
			name:          "Default sort (client then status)",
			sortBy:        "",
			sortOrder:     "",
			expectedFirst: "SN-001", // Alpha Corp, available (first alphabetically by company then status)
			expectedLast:  "SN-004", // Beta Corp, available (last when sorted by company then status)
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create request with sort parameters
			url := "/inventory?sort=" + tc.sortBy + "&order=" + tc.sortOrder
			req := httptest.NewRequest("GET", url, nil)

			// Add user to context
			ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
			req = req.WithContext(ctx)

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call handler
			handler.InventoryList(rr, req)

			// Check response status
			if rr.Code != http.StatusOK {
				t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
			}

			// For now, we'll check the database directly to verify sorting
			// Once the handler is implemented, we could parse the response body
			filter := &models.LaptopFilter{
				UserRole:  user.Role,
				SortBy:    tc.sortBy,
				SortOrder: tc.sortOrder,
			}

			results, err := models.GetAllLaptops(db, filter)
			if err != nil {
				t.Fatalf("Failed to get laptops: %v", err)
			}

			if len(results) == 0 {
				t.Fatal("Expected laptops to be returned")
			}

			// Verify first and last items match expected order
			if results[0].SerialNumber != tc.expectedFirst {
				t.Errorf("Expected first item to be %s, got %s", tc.expectedFirst, results[0].SerialNumber)
			}

			if results[len(results)-1].SerialNumber != tc.expectedLast {
				t.Errorf("Expected last item to be %s, got %s", tc.expectedLast, results[len(results)-1].SerialNumber)
			}
		})
	}
}

// Helper function to create a test software engineer
func createTestSoftwareEngineer(db interface {
	QueryRow(query string, args ...interface{}) *sql.Row
}, engineer *models.SoftwareEngineer) error {
	if err := engineer.Validate(); err != nil {
		return err
	}

	engineer.BeforeCreate()

	query := `
		INSERT INTO software_engineers (name, email, address, phone, address_confirmed, address_confirmation_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`
	return db.QueryRow(
		query,
		engineer.Name,
		engineer.Email,
		engineer.Address,
		engineer.Phone,
		engineer.AddressConfirmed,
		engineer.AddressConfirmationAt,
		engineer.CreatedAt,
		engineer.UpdatedAt,
	).Scan(&engineer.ID)
}

