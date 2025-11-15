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

