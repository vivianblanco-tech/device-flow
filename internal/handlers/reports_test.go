package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/laptop-tracking-system/internal/database"
	"github.com/yourusername/laptop-tracking-system/internal/middleware"
	"github.com/yourusername/laptop-tracking-system/internal/models"
)

// TestReportsIndex tests the reports index page
func TestReportsIndex(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	templates := loadTestTemplates(t)
	handler := NewReportsHandler(db, templates)

	t.Run("Client user can access reports index", func(t *testing.T) {
		// Create test client company
		var companyID int64
		err := db.QueryRow(
			`INSERT INTO client_companies (name, contact_info, created_at, updated_at)
			VALUES ($1, $2, NOW(), NOW()) RETURNING id`,
			"Test Company", "test@company.com",
		).Scan(&companyID)
		if err != nil {
			t.Fatalf("Failed to create test company: %v", err)
		}

		// Create test client user
		var userID int64
		err = db.QueryRow(
			`INSERT INTO users (email, password_hash, role, client_company_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING id`,
			"client@test.com", "hashed_password", models.RoleClient, companyID,
		).Scan(&userID)
		if err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}

		user := &models.User{
			ID:              userID,
			Email:           "client@test.com",
			Role:            models.RoleClient,
			ClientCompanyID: &companyID,
		}

		req := httptest.NewRequest(http.MethodGet, "/reports", nil)
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		handler.ReportsIndex(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
		}
	})

	t.Run("Non-client user cannot access reports", func(t *testing.T) {
		user := &models.User{
			ID:    1,
			Email: "logistics@test.com",
			Role:  models.RoleLogistics,
		}

		req := httptest.NewRequest(http.MethodGet, "/reports", nil)
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		handler.ReportsIndex(rr, req)

		if rr.Code != http.StatusForbidden {
			t.Errorf("expected status %d, got %d", http.StatusForbidden, rr.Code)
		}
	})
}

// TestShipmentStatusDashboard tests the shipment status dashboard report
func TestShipmentStatusDashboard(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	templates := loadTestTemplates(t)
	handler := NewReportsHandler(db, templates)

	t.Run("Client user can view shipment status dashboard", func(t *testing.T) {
		// Create test client company
		var companyID int64
		err := db.QueryRow(
			`INSERT INTO client_companies (name, contact_info, created_at, updated_at)
			VALUES ($1, $2, NOW(), NOW()) RETURNING id`,
			"Test Company", "test@company.com",
		).Scan(&companyID)
		if err != nil {
			t.Fatalf("Failed to create test company: %v", err)
		}

		// Create test client user
		var userID int64
		err = db.QueryRow(
			`INSERT INTO users (email, password_hash, role, client_company_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING id`,
			"client@test.com", "hashed_password", models.RoleClient, companyID,
		).Scan(&userID)
		if err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}

		// Create test shipments
		now := time.Now()
		_, err = db.Exec(
			`INSERT INTO shipments (client_company_id, status, shipment_type, laptop_count, jira_ticket_number, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			companyID, models.ShipmentStatusPendingPickup, models.ShipmentTypeSingleFullJourney, 1, "TEST-100", now, now,
		)
		if err != nil {
			t.Fatalf("Failed to create test shipment: %v", err)
		}

		_, err = db.Exec(
			`INSERT INTO shipments (client_company_id, status, shipment_type, laptop_count, jira_ticket_number, created_at, updated_at, delivered_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			companyID, models.ShipmentStatusDelivered, models.ShipmentTypeSingleFullJourney, 1, "TEST-101", now.AddDate(0, 0, -30), now.AddDate(0, 0, -25), now.AddDate(0, 0, -20),
		)
		if err != nil {
			t.Fatalf("Failed to create test shipment: %v", err)
		}

		user := &models.User{
			ID:              userID,
			Email:           "client@test.com",
			Role:            models.RoleClient,
			ClientCompanyID: &companyID,
		}

		req := httptest.NewRequest(http.MethodGet, "/reports/shipment-status", nil)
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		handler.ShipmentStatusDashboard(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
		}
	})

	t.Run("Client user only sees their company's shipments", func(t *testing.T) {
		// Create two companies
		var company1ID int64
		err := db.QueryRow(
			`INSERT INTO client_companies (name, contact_info, created_at, updated_at)
			VALUES ($1, $2, NOW(), NOW()) RETURNING id`,
			"Company 1", "company1@test.com",
		).Scan(&company1ID)
		if err != nil {
			t.Fatalf("Failed to create test company: %v", err)
		}

		var company2ID int64
		err = db.QueryRow(
			`INSERT INTO client_companies (name, contact_info, created_at, updated_at)
			VALUES ($1, $2, NOW(), NOW()) RETURNING id`,
			"Company 2", "company2@test.com",
		).Scan(&company2ID)
		if err != nil {
			t.Fatalf("Failed to create test company: %v", err)
		}

		// Create client user for company 1
		var userID int64
		err = db.QueryRow(
			`INSERT INTO users (email, password_hash, role, client_company_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING id`,
			"client1@test.com", "hashed_password", models.RoleClient, company1ID,
		).Scan(&userID)
		if err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}

		// Create shipments for both companies
		now := time.Now()
		_, err = db.Exec(
			`INSERT INTO shipments (client_company_id, status, shipment_type, laptop_count, jira_ticket_number, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			company1ID, models.ShipmentStatusPendingPickup, models.ShipmentTypeSingleFullJourney, 1, "TEST-200", now, now,
		)
		if err != nil {
			t.Fatalf("Failed to create test shipment: %v", err)
		}

		_, err = db.Exec(
			`INSERT INTO shipments (client_company_id, status, shipment_type, laptop_count, jira_ticket_number, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			company2ID, models.ShipmentStatusPendingPickup, models.ShipmentTypeSingleFullJourney, 1, "TEST-201", now, now,
		)
		if err != nil {
			t.Fatalf("Failed to create test shipment: %v", err)
		}

		user := &models.User{
			ID:              userID,
			Email:           "client1@test.com",
			Role:            models.RoleClient,
			ClientCompanyID: &company1ID,
		}

		req := httptest.NewRequest(http.MethodGet, "/reports/shipment-status", nil)
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		handler.ShipmentStatusDashboard(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
		}
	})
}

// TestInventorySummaryReport tests the inventory summary report
func TestInventorySummaryReport(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	templates := loadTestTemplates(t)
	handler := NewReportsHandler(db, templates)

	t.Run("Client user can view inventory summary", func(t *testing.T) {
		// Create test client company
		var companyID int64
		err := db.QueryRow(
			`INSERT INTO client_companies (name, contact_info, created_at, updated_at)
			VALUES ($1, $2, NOW(), NOW()) RETURNING id`,
			"Test Company", "test@company.com",
		).Scan(&companyID)
		if err != nil {
			t.Fatalf("Failed to create test company: %v", err)
		}

		// Create test client user
		var userID int64
		err = db.QueryRow(
			`INSERT INTO users (email, password_hash, role, client_company_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING id`,
			"client@test.com", "hashed_password", models.RoleClient, companyID,
		).Scan(&userID)
		if err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}

		// Create test laptops
		_, err = db.Exec(
			`INSERT INTO laptops (serial_number, brand, model, cpu, ram_gb, ssd_gb, status, client_company_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())`,
			"SN-001", "Dell", "XPS 13", "Intel i7", "16GB", "512GB", models.LaptopStatusAvailable, companyID,
		)
		if err != nil {
			t.Fatalf("Failed to create test laptop: %v", err)
		}

		_, err = db.Exec(
			`INSERT INTO laptops (serial_number, brand, model, cpu, ram_gb, ssd_gb, status, client_company_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())`,
			"SN-002", "HP", "EliteBook", "Intel i5", "8GB", "256GB", models.LaptopStatusDelivered, companyID,
		)
		if err != nil {
			t.Fatalf("Failed to create test laptop: %v", err)
		}

		user := &models.User{
			ID:              userID,
			Email:           "client@test.com",
			Role:            models.RoleClient,
			ClientCompanyID: &companyID,
		}

		req := httptest.NewRequest(http.MethodGet, "/reports/inventory-summary", nil)
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		handler.InventorySummaryReport(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
		}
	})
}

// TestShipmentTimelineReport tests the shipment timeline report
func TestShipmentTimelineReport(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	templates := loadTestTemplates(t)
	handler := NewReportsHandler(db, templates)

	t.Run("Client user can view shipment timeline", func(t *testing.T) {
		// Create test client company
		var companyID int64
		err := db.QueryRow(
			`INSERT INTO client_companies (name, contact_info, created_at, updated_at)
			VALUES ($1, $2, NOW(), NOW()) RETURNING id`,
			"Test Company", "test@company.com",
		).Scan(&companyID)
		if err != nil {
			t.Fatalf("Failed to create test company: %v", err)
		}

		// Create test client user
		var userID int64
		err = db.QueryRow(
			`INSERT INTO users (email, password_hash, role, client_company_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING id`,
			"client@test.com", "hashed_password", models.RoleClient, companyID,
		).Scan(&userID)
		if err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}

		// Create test shipment with timeline data
		now := time.Now()
		var shipmentID int64
		err = db.QueryRow(
			`INSERT INTO shipments (client_company_id, status, shipment_type, laptop_count, jira_ticket_number, 
			 pickup_scheduled_date, picked_up_at, arrived_warehouse_at, released_warehouse_at, delivered_at,
			 courier_name, tracking_number, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING id`,
			companyID, models.ShipmentStatusDelivered, models.ShipmentTypeSingleFullJourney, 1, "TEST-300",
			now.AddDate(0, 0, -10), now.AddDate(0, 0, -9), now.AddDate(0, 0, -7), now.AddDate(0, 0, -5), now.AddDate(0, 0, -2),
			"UPS", "TRACK123", now.AddDate(0, 0, -10), now.AddDate(0, 0, -2),
		).Scan(&shipmentID)
		if err != nil {
			t.Fatalf("Failed to create test shipment: %v", err)
		}

		user := &models.User{
			ID:              userID,
			Email:           "client@test.com",
			Role:            models.RoleClient,
			ClientCompanyID: &companyID,
		}

		req := httptest.NewRequest(http.MethodGet, "/reports/shipment-timeline", nil)
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		handler.ShipmentTimelineReport(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
		}
	})
}

// TestReportsExport tests export functionality
func TestReportsExport(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	templates := loadTestTemplates(t)
	handler := NewReportsHandler(db, templates)

	t.Run("Client user can export shipment status as CSV", func(t *testing.T) {
		// Create test client company
		var companyID int64
		err := db.QueryRow(
			`INSERT INTO client_companies (name, contact_info, created_at, updated_at)
			VALUES ($1, $2, NOW(), NOW()) RETURNING id`,
			"Test Company", "test@company.com",
		).Scan(&companyID)
		if err != nil {
			t.Fatalf("Failed to create test company: %v", err)
		}

		// Create test client user
		var userID int64
		err = db.QueryRow(
			`INSERT INTO users (email, password_hash, role, client_company_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING id`,
			"client@test.com", "hashed_password", models.RoleClient, companyID,
		).Scan(&userID)
		if err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}

		user := &models.User{
			ID:              userID,
			Email:           "client@test.com",
			Role:            models.RoleClient,
			ClientCompanyID: &companyID,
		}

		req := httptest.NewRequest(http.MethodGet, "/reports/shipment-status?format=csv", nil)
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		handler.ShipmentStatusDashboard(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
		}

		contentType := rr.Header().Get("Content-Type")
		if contentType != "text/csv" {
			t.Errorf("expected Content-Type text/csv, got %s", contentType)
		}
	})

	t.Run("Client user can export inventory summary as Excel", func(t *testing.T) {
		// Create test client company with unique name
		companyName := fmt.Sprintf("Test Company Export %d", time.Now().UnixNano())
		var companyID int64
		err := db.QueryRow(
			`INSERT INTO client_companies (name, contact_info, created_at, updated_at)
			VALUES ($1, $2, NOW(), NOW()) RETURNING id`,
			companyName, "test@company.com",
		).Scan(&companyID)
		if err != nil {
			t.Fatalf("Failed to create test company: %v", err)
		}

		// Create test client user with unique email
		email := fmt.Sprintf("client-export-%d@test.com", time.Now().UnixNano())
		var userID int64
		err = db.QueryRow(
			`INSERT INTO users (email, password_hash, role, client_company_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING id`,
			email, "hashed_password", models.RoleClient, companyID,
		).Scan(&userID)
		if err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}

		user := &models.User{
			ID:              userID,
			Email:           email,
			Role:            models.RoleClient,
			ClientCompanyID: &companyID,
		}

		req := httptest.NewRequest(http.MethodGet, "/reports/inventory-summary?format=xlsx", nil)
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		handler.InventorySummaryReport(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
		}

		contentType := rr.Header().Get("Content-Type")
		if contentType != "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
			t.Errorf("expected Excel Content-Type, got %s", contentType)
		}
	})
}

