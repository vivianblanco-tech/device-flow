package unit

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/yourusername/laptop-tracking-system/internal/database"
	"github.com/yourusername/laptop-tracking-system/internal/handlers"
	"github.com/yourusername/laptop-tracking-system/internal/middleware"
	"github.com/yourusername/laptop-tracking-system/internal/models"
)

// TestLogisticsCreateMinimalBulkShipment tests that logistics users can create
// a bulk shipment with only JIRA ticket and Client Company
func TestLogisticsCreateMinimalBulkShipment(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create test client company
	var companyID int64
	err := db.QueryRowContext(ctx,
		`INSERT INTO client_companies (name, contact_info, created_at)
		VALUES ($1, $2, $3) RETURNING id`,
		"Test Company", json.RawMessage(`{"email":"test@company.com"}`), time.Now(),
	).Scan(&companyID)
	if err != nil {
		t.Fatalf("Failed to create test company: %v", err)
	}

	// Create logistics user
	var userID int64
	err = db.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"logistics@bairesdev.com", "$2a$12$test.hash.for.testing.purposes", models.RoleLogistics, time.Now(), time.Now(),
	).Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to create logistics user: %v", err)
	}

	// Create handler
	handler := handlers.NewPickupFormHandler(db, nil, nil)

	// Create form data with only JIRA ticket and company
	formData := url.Values{}
	formData.Set("shipment_type", "bulk_to_warehouse")
	formData.Set("client_company_id", strconv.FormatInt(companyID, 10))
	formData.Set("jira_ticket_number", "SCOP-12345")

	req := httptest.NewRequest(http.MethodPost, "/pickup-form", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Add logistics user to context
	ctx = context.WithValue(req.Context(), middleware.UserContextKey, &models.User{
		ID:    userID,
		Email: "logistics@bairesdev.com",
		Role:  models.RoleLogistics,
	})
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	// Submit the form
	handler.PickupFormSubmit(w, req)

	// Should redirect to shipment detail page (302)
	if w.Code != http.StatusSeeOther {
		t.Errorf("Expected status 302 (redirect), got %d. Body: %s", w.Code, w.Body.String())
	}

	// Verify shipment was created in database
	var shipmentID int64
	var shipmentType string
	var status string
	var jiraTicket string
	err = db.QueryRowContext(ctx,
		`SELECT id, shipment_type, status, jira_ticket_number 
		FROM shipments 
		WHERE client_company_id = $1 AND jira_ticket_number = $2`,
		companyID, "SCOP-12345",
	).Scan(&shipmentID, &shipmentType, &status, &jiraTicket)
	
	if err != nil {
		t.Fatalf("Shipment was not created: %v", err)
	}

	// Verify shipment details
	if shipmentType != string(models.ShipmentTypeBulkToWarehouse) {
		t.Errorf("Expected shipment type 'bulk_to_warehouse', got '%s'", shipmentType)
	}

	if status != string(models.ShipmentStatusPendingPickup) {
		t.Errorf("Expected status 'pending_pickup_from_client', got '%s'", status)
	}

	if jiraTicket != "SCOP-12345" {
		t.Errorf("Expected JIRA ticket 'SCOP-12345', got '%s'", jiraTicket)
	}

	// Verify NO pickup form was created yet
	var pickupFormCount int
	err = db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM pickup_forms WHERE shipment_id = $1`,
		shipmentID,
	).Scan(&pickupFormCount)
	
	if err != nil {
		t.Fatalf("Failed to check pickup forms: %v", err)
	}

	if pickupFormCount != 0 {
		t.Errorf("Expected 0 pickup forms, found %d (form should not be created yet)", pickupFormCount)
	}
}

// TestClientCompletesBulkShipmentDetails tests that a client can complete
// the bulk shipment details via magic link
func TestClientCompletesBulkShipmentDetails(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create test client company
	var companyID int64
	err := db.QueryRowContext(ctx,
		`INSERT INTO client_companies (name, contact_info, created_at)
		VALUES ($1, $2, $3) RETURNING id`,
		"Test Company", json.RawMessage(`{"email":"test@company.com"}`), time.Now(),
	).Scan(&companyID)
	if err != nil {
		t.Fatalf("Failed to create test company: %v", err)
	}

	// Create client user
	var userID int64
	err = db.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"client@company.com", "$2a$12$test.hash.for.testing.purposes", models.RoleClient, time.Now(), time.Now(),
	).Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to create client user: %v", err)
	}

	// Create a minimal shipment (as logistics would)
	// Use laptop_count of 1 initially (will be updated when client completes form)
	var shipmentID int64
	err = db.QueryRowContext(ctx,
		`INSERT INTO shipments (shipment_type, client_company_id, status, laptop_count, jira_ticket_number, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		models.ShipmentTypeBulkToWarehouse, companyID, models.ShipmentStatusPendingPickup, 1, "SCOP-12345", time.Now(), time.Now(),
	).Scan(&shipmentID)
	if err != nil {
		t.Fatalf("Failed to create shipment: %v", err)
	}

	// Create handler
	shipmentsHandler := handlers.NewShipmentsHandler(db, nil, nil)

	// Client completes the shipment details
	futureDate := time.Now().AddDate(0, 0, 7).Format("2006-01-02") // 7 days from now
	formData := url.Values{}
	formData.Set("assignment_type", "bulk")
	formData.Set("number_of_laptops", "5")
	formData.Set("number_of_boxes", "2")
	formData.Set("bulk_length", "50.5")
	formData.Set("bulk_width", "40.2")
	formData.Set("bulk_height", "30.8")
	formData.Set("bulk_weight", "25.5")
	formData.Set("contact_name", "John Doe")
	formData.Set("contact_email", "john@company.com")
	formData.Set("contact_phone", "+1-555-0123")
	formData.Set("pickup_address", "123 Main St")
	formData.Set("pickup_city", "New York")
	formData.Set("pickup_state", "NY")
	formData.Set("pickup_zip", "10001")
	formData.Set("pickup_date", futureDate)
	formData.Set("pickup_time_slot", "morning")
	formData.Set("special_instructions", "Please call before arriving")
	formData.Set("include_accessories", "true")
	formData.Set("accessories_description", "Chargers and mice")

	req := httptest.NewRequest(http.MethodPost, "/shipments/"+strconv.FormatInt(shipmentID, 10)+"/form", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Set URL vars for mux router
	req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(shipmentID, 10)})

	// Add client user to context
	ctx = context.WithValue(req.Context(), middleware.UserContextKey, &models.User{
		ID:    userID,
		Email: "client@company.com",
		Role:  models.RoleClient,
	})
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	// Submit the form
	shipmentsHandler.ShipmentPickupFormSubmit(w, req)

	// Should redirect to shipment detail page (302)
	if w.Code != http.StatusSeeOther {
		t.Errorf("Expected status 302 (redirect), got %d. Body: %s", w.Code, w.Body.String())
	}

	// Verify pickup form was created
	var pickupFormData map[string]interface{}
	var formDataJSON []byte
	err = db.QueryRowContext(ctx,
		`SELECT form_data FROM pickup_forms WHERE shipment_id = $1`,
		shipmentID,
	).Scan(&formDataJSON)
	
	if err == sql.ErrNoRows {
		t.Fatal("Pickup form was not created")
	}
	if err != nil {
		t.Fatalf("Failed to get pickup form: %v", err)
	}

	err = json.Unmarshal(formDataJSON, &pickupFormData)
	if err != nil {
		t.Fatalf("Failed to unmarshal form data: %v", err)
	}

	// Verify form data
	if pickupFormData["contact_name"] != "John Doe" {
		t.Errorf("Expected contact_name 'John Doe', got '%v'", pickupFormData["contact_name"])
	}

	if pickupFormData["number_of_laptops"] != float64(5) {
		t.Errorf("Expected number_of_laptops 5, got '%v'", pickupFormData["number_of_laptops"])
	}

	if pickupFormData["bulk_length"] != 50.5 {
		t.Errorf("Expected bulk_length 50.5, got '%v'", pickupFormData["bulk_length"])
	}

	// Verify shipment was updated with laptop count
	var laptopCount int
	err = db.QueryRowContext(ctx,
		`SELECT laptop_count FROM shipments WHERE id = $1`,
		shipmentID,
	).Scan(&laptopCount)
	
	if err != nil {
		t.Fatalf("Failed to get shipment: %v", err)
	}

	if laptopCount != 5 {
		t.Errorf("Expected laptop_count 5, got %d", laptopCount)
	}
}

// TestLogisticsEditsBulkShipmentDetails tests that logistics can edit
// shipment details but NOT JIRA ticket or Company
func TestLogisticsEditsBulkShipmentDetails(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create test client company
	var companyID int64
	err := db.QueryRowContext(ctx,
		`INSERT INTO client_companies (name, contact_info, created_at)
		VALUES ($1, $2, $3) RETURNING id`,
		"Test Company", json.RawMessage(`{"email":"test@company.com"}`), time.Now(),
	).Scan(&companyID)
	if err != nil {
		t.Fatalf("Failed to create test company: %v", err)
	}

	// Create logistics user
	var userID int64
	err = db.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"logistics@bairesdev.com", "$2a$12$test.hash.for.testing.purposes", models.RoleLogistics, time.Now(), time.Now(),
	).Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to create logistics user: %v", err)
	}

	// Create a shipment with complete details
	var shipmentID int64
	err = db.QueryRowContext(ctx,
		`INSERT INTO shipments (shipment_type, client_company_id, status, laptop_count, jira_ticket_number, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		models.ShipmentTypeBulkToWarehouse, companyID, models.ShipmentStatusPendingPickup, 5, "SCOP-12345", time.Now(), time.Now(),
	).Scan(&shipmentID)
	if err != nil {
		t.Fatalf("Failed to create shipment: %v", err)
	}

	// Create existing pickup form
	existingPickupDate := time.Now().AddDate(0, 0, 7).Format("2006-01-02")
	existingFormData := map[string]interface{}{
		"contact_name":            "John Doe",
		"contact_email":           "john@company.com",
		"contact_phone":           "+1-555-0123",
		"assignment_type":         "bulk",
		"number_of_laptops":       5,
		"number_of_boxes":         2,
		"bulk_length":             50.5,
		"bulk_width":              40.2,
		"bulk_height":             30.8,
		"bulk_weight":             25.5,
		"pickup_address":          "123 Main St",
		"pickup_city":             "New York",
		"pickup_state":            "NY",
		"pickup_zip":              "10001",
		"pickup_date":             existingPickupDate,
		"pickup_time_slot":        "morning",
		"special_instructions":    "Please call before arriving",
		"include_accessories":     true,
		"accessories_description": "Chargers and mice",
	}
	existingFormDataJSON, _ := json.Marshal(existingFormData)

	_, err = db.ExecContext(ctx,
		`INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data)
		VALUES ($1, $2, $3, $4)`,
		shipmentID, userID, time.Now(), existingFormDataJSON,
	)
	if err != nil {
		t.Fatalf("Failed to create pickup form: %v", err)
	}

	// Create handler
	shipmentsHandler := handlers.NewShipmentsHandler(db, nil, nil)

	// Logistics edits the shipment details (changing contact info and laptop count)
	futureDate := time.Now().AddDate(0, 0, 7).Format("2006-01-02") // 7 days from now
	formData := url.Values{}
	formData.Set("assignment_type", "bulk")
	formData.Set("number_of_laptops", "7") // Changed from 5 to 7
	formData.Set("number_of_boxes", "3") // Changed from 2 to 3
	formData.Set("bulk_length", "50.5")
	formData.Set("bulk_width", "40.2")
	formData.Set("bulk_height", "30.8")
	formData.Set("bulk_weight", "25.5")
	formData.Set("contact_name", "Jane Smith") // Changed
	formData.Set("contact_email", "jane@company.com") // Changed
	formData.Set("contact_phone", "+1-555-9999") // Changed
	formData.Set("pickup_address", "123 Main St")
	formData.Set("pickup_city", "New York")
	formData.Set("pickup_state", "NY")
	formData.Set("pickup_zip", "10001")
	formData.Set("pickup_date", futureDate)
	formData.Set("pickup_time_slot", "afternoon") // Changed
	formData.Set("special_instructions", "Updated instructions")
	formData.Set("include_accessories", "true")
	formData.Set("accessories_description", "Updated accessories")

	req := httptest.NewRequest(http.MethodPost, "/shipments/"+strconv.FormatInt(shipmentID, 10)+"/form", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Set URL vars for mux router
	req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(shipmentID, 10)})

	// Add logistics user to context
	ctx = context.WithValue(req.Context(), middleware.UserContextKey, &models.User{
		ID:    userID,
		Email: "logistics@bairesdev.com",
		Role:  models.RoleLogistics,
	})
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	// Submit the form
	shipmentsHandler.ShipmentPickupFormSubmit(w, req)

	// Should redirect to shipment detail page (302)
	if w.Code != http.StatusSeeOther {
		t.Errorf("Expected status 302 (redirect), got %d. Body: %s", w.Code, w.Body.String())
	}

	// Verify pickup form was updated
	var updatedFormDataMap map[string]interface{}
	var updatedFormDataJSON []byte
	err = db.QueryRowContext(ctx,
		`SELECT form_data FROM pickup_forms WHERE shipment_id = $1`,
		shipmentID,
	).Scan(&updatedFormDataJSON)
	
	if err != nil {
		t.Fatalf("Failed to get updated pickup form: %v", err)
	}

	err = json.Unmarshal(updatedFormDataJSON, &updatedFormDataMap)
	if err != nil {
		t.Fatalf("Failed to unmarshal form data: %v", err)
	}

	// Verify updated fields
	if updatedFormDataMap["contact_name"] != "Jane Smith" {
		t.Errorf("Expected contact_name 'Jane Smith', got '%v'", updatedFormDataMap["contact_name"])
	}

	if updatedFormDataMap["number_of_laptops"] != float64(7) {
		t.Errorf("Expected number_of_laptops 7, got '%v'", updatedFormDataMap["number_of_laptops"])
	}

	if updatedFormDataMap["pickup_time_slot"] != "afternoon" {
		t.Errorf("Expected pickup_time_slot 'afternoon', got '%v'", updatedFormDataMap["pickup_time_slot"])
	}

	// Verify JIRA ticket and Company remain unchanged
	var jiraTicket string
	var clientCompID int64
	err = db.QueryRowContext(ctx,
		`SELECT jira_ticket_number, client_company_id FROM shipments WHERE id = $1`,
		shipmentID,
	).Scan(&jiraTicket, &clientCompID)
	
	if err != nil {
		t.Fatalf("Failed to get shipment: %v", err)
	}

	if jiraTicket != "SCOP-12345" {
		t.Errorf("JIRA ticket should not change. Expected 'SCOP-12345', got '%s'", jiraTicket)
	}

	if clientCompID != companyID {
		t.Errorf("Client company should not change. Expected %d, got %d", companyID, clientCompID)
	}
}

