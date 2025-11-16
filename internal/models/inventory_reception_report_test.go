package models

import (
	"context"
	"testing"

	"github.com/yourusername/laptop-tracking-system/internal/database"
)

// TestGetAllLaptopsWithReceptionReportInfo tests that GetAllLaptops includes reception report information
func TestGetAllLaptopsWithReceptionReportInfo(t *testing.T) {
	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	// Create a warehouse user
	var warehouseUserID int64
	err := db.QueryRow(
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id`,
		"warehouse@bairesdev.com", "hashed_password", RoleWarehouse,
	).Scan(&warehouseUserID)
	if err != nil {
		t.Fatalf("Failed to create warehouse user: %v", err)
	}

	// Create a client company first
	var clientID int64
	err = db.QueryRow(
		`INSERT INTO client_companies (name, created_at, updated_at)
		VALUES ($1, NOW(), NOW()) RETURNING id`,
		"Test Company",
	).Scan(&clientID)
	if err != nil {
		t.Fatalf("Failed to create client company: %v", err)
	}

	// Create test laptops with different reception report scenarios
	// Laptop 1: Has reception report (pending approval)
	laptop1 := &Laptop{
		SerialNumber:    "SN001",
		Brand:           "Dell",
		Model:           "Latitude 5520",
		CPU:             "Intel Core i7",
		RAMGB:           "16GB",
		SSDGB:           "512GB",
		Status:          LaptopStatusAtWarehouse,
		ClientCompanyID: &clientID,
	}
	err = CreateLaptop(db, laptop1)
	if err != nil {
		t.Fatalf("Failed to create laptop1: %v", err)
	}

	// Create reception report for laptop1
	receptionReport1 := &ReceptionReport{
		LaptopID:               laptop1.ID,
		WarehouseUserID:        warehouseUserID,
		PhotoSerialNumber:      "/uploads/photo1.jpg",
		PhotoExternalCondition: "/uploads/photo2.jpg",
		PhotoWorkingCondition:  "/uploads/photo3.jpg",
		Status:                 ReceptionReportStatusPendingApproval,
	}
	receptionReport1.BeforeCreate()
	err = CreateReceptionReport(context.Background(), db, receptionReport1)
	if err != nil {
		t.Fatalf("Failed to create reception report for laptop1: %v", err)
	}

	// Laptop 2: No reception report
	laptop2 := &Laptop{
		SerialNumber:    "SN002",
		Brand:           "HP",
		Model:           "EliteBook 840",
		CPU:             "Intel Core i9",
		RAMGB:           "32GB",
		SSDGB:           "1024GB",
		Status:          LaptopStatusAtWarehouse,
		ClientCompanyID: &clientID,
	}
	err = CreateLaptop(db, laptop2)
	if err != nil {
		t.Fatalf("Failed to create laptop2: %v", err)
	}

	// Laptop 3: Has reception report (approved)
	laptop3 := &Laptop{
		SerialNumber:    "SN003",
		Brand:           "Lenovo",
		Model:           "ThinkPad X1",
		CPU:             "Intel Core i7",
		RAMGB:           "16GB",
		SSDGB:           "512GB",
		Status:          LaptopStatusAvailable,
		ClientCompanyID: &clientID,
	}
	err = CreateLaptop(db, laptop3)
	if err != nil {
		t.Fatalf("Failed to create laptop3: %v", err)
	}

	// Create approved reception report for laptop3
	receptionReport3 := &ReceptionReport{
		LaptopID:               laptop3.ID,
		WarehouseUserID:        warehouseUserID,
		PhotoSerialNumber:      "/uploads/photo4.jpg",
		PhotoExternalCondition: "/uploads/photo5.jpg",
		PhotoWorkingCondition:  "/uploads/photo6.jpg",
		Status:                 ReceptionReportStatusApproved,
	}
	receptionReport3.BeforeCreate()
	err = CreateReceptionReport(context.Background(), db, receptionReport3)
	if err != nil {
		t.Fatalf("Failed to create reception report for laptop3: %v", err)
	}

	// Test: Get all laptops for warehouse user
	filter := &LaptopFilter{
		UserRole: RoleWarehouse,
	}
	laptops, err := GetAllLaptops(db, filter)
	if err != nil {
		t.Fatalf("GetAllLaptops failed: %v", err)
	}

	// We should have 3 laptops
	if len(laptops) != 3 {
		t.Fatalf("Expected 3 laptops, got %d", len(laptops))
	}

	// Verify laptop1 has reception report info
	for _, laptop := range laptops {
		if laptop.ID == laptop1.ID {
			if !laptop.HasReceptionReport {
				t.Error("Expected laptop1 to have HasReceptionReport = true")
			}
			if laptop.ReceptionReportStatus != string(ReceptionReportStatusPendingApproval) {
				t.Errorf("Expected laptop1 to have ReceptionReportStatus = 'pending_approval', got '%s'", laptop.ReceptionReportStatus)
			}
			if laptop.ReceptionReportID == nil || *laptop.ReceptionReportID != receptionReport1.ID {
				t.Errorf("Expected laptop1 to have ReceptionReportID = %d, got %v", receptionReport1.ID, laptop.ReceptionReportID)
			}
		}

		if laptop.ID == laptop2.ID {
			if laptop.HasReceptionReport {
				t.Error("Expected laptop2 to have HasReceptionReport = false")
			}
			if laptop.ReceptionReportStatus != "" {
				t.Errorf("Expected laptop2 to have empty ReceptionReportStatus, got '%s'", laptop.ReceptionReportStatus)
			}
			if laptop.ReceptionReportID != nil {
				t.Errorf("Expected laptop2 to have ReceptionReportID = nil, got %v", laptop.ReceptionReportID)
			}
		}

		if laptop.ID == laptop3.ID {
			if !laptop.HasReceptionReport {
				t.Error("Expected laptop3 to have HasReceptionReport = true")
			}
			if laptop.ReceptionReportStatus != string(ReceptionReportStatusApproved) {
				t.Errorf("Expected laptop3 to have ReceptionReportStatus = 'approved', got '%s'", laptop.ReceptionReportStatus)
			}
			if laptop.ReceptionReportID == nil || *laptop.ReceptionReportID != receptionReport3.ID {
				t.Errorf("Expected laptop3 to have ReceptionReportID = %d, got %v", receptionReport3.ID, laptop.ReceptionReportID)
			}
		}
	}
}

// TestGetAllLaptopsForLogisticsUsersIncludesReceptionReports tests that logistics users also get reception report info
func TestGetAllLaptopsForLogisticsUsersIncludesReceptionReports(t *testing.T) {
	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	// Create a warehouse user to create reception reports
	var warehouseUserID int64
	err := db.QueryRow(
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id`,
		"warehouse@bairesdev.com", "hashed_password", RoleWarehouse,
	).Scan(&warehouseUserID)
	if err != nil {
		t.Fatalf("Failed to create warehouse user: %v", err)
	}

	// Create a client company first
	var clientID int64
	err = db.QueryRow(
		`INSERT INTO client_companies (name, created_at, updated_at)
		VALUES ($1, NOW(), NOW()) RETURNING id`,
		"Test Company",
	).Scan(&clientID)
	if err != nil {
		t.Fatalf("Failed to create client company: %v", err)
	}

	// Create test laptops
	// Laptop 1: At warehouse with pending approval reception report
	laptop1 := &Laptop{
		SerialNumber:    "SN100",
		Brand:           "Dell",
		Model:           "Latitude 7420",
		CPU:             "Intel Core i7",
		RAMGB:           "16GB",
		SSDGB:           "512GB",
		Status:          LaptopStatusAtWarehouse,
		ClientCompanyID: &clientID,
	}
	err = CreateLaptop(db, laptop1)
	if err != nil {
		t.Fatalf("Failed to create laptop1: %v", err)
	}

	// Create pending approval reception report for laptop1
	receptionReport1 := &ReceptionReport{
		LaptopID:               laptop1.ID,
		WarehouseUserID:        warehouseUserID,
		PhotoSerialNumber:      "/uploads/serial1.jpg",
		PhotoExternalCondition: "/uploads/external1.jpg",
		PhotoWorkingCondition:  "/uploads/working1.jpg",
		Status:                 ReceptionReportStatusPendingApproval,
	}
	receptionReport1.BeforeCreate()
	err = CreateReceptionReport(context.Background(), db, receptionReport1)
	if err != nil {
		t.Fatalf("Failed to create reception report for laptop1: %v", err)
	}

	// Laptop 2: At warehouse without reception report
	laptop2 := &Laptop{
		SerialNumber:    "SN200",
		Brand:           "HP",
		Model:           "EliteBook 850",
		CPU:             "Intel Core i9",
		RAMGB:           "32GB",
		SSDGB:           "1024GB",
		Status:          LaptopStatusAtWarehouse,
		ClientCompanyID: &clientID,
	}
	err = CreateLaptop(db, laptop2)
	if err != nil {
		t.Fatalf("Failed to create laptop2: %v", err)
	}

	// Laptop 3: Available with approved reception report
	laptop3 := &Laptop{
		SerialNumber:    "SN300",
		Brand:           "Lenovo",
		Model:           "ThinkPad X1 Carbon",
		CPU:             "Intel Core i7",
		RAMGB:           "16GB",
		SSDGB:           "512GB",
		Status:          LaptopStatusAvailable,
		ClientCompanyID: &clientID,
	}
	err = CreateLaptop(db, laptop3)
	if err != nil {
		t.Fatalf("Failed to create laptop3: %v", err)
	}

	// Create approved reception report for laptop3
	receptionReport3 := &ReceptionReport{
		LaptopID:               laptop3.ID,
		WarehouseUserID:        warehouseUserID,
		PhotoSerialNumber:      "/uploads/serial3.jpg",
		PhotoExternalCondition: "/uploads/external3.jpg",
		PhotoWorkingCondition:  "/uploads/working3.jpg",
		Status:                 ReceptionReportStatusApproved,
	}
	receptionReport3.BeforeCreate()
	err = CreateReceptionReport(context.Background(), db, receptionReport3)
	if err != nil {
		t.Fatalf("Failed to create reception report for laptop3: %v", err)
	}

	// Test: Get all laptops for logistics user
	filter := &LaptopFilter{
		UserRole: RoleLogistics,
	}
	laptops, err := GetAllLaptops(db, filter)
	if err != nil {
		t.Fatalf("GetAllLaptops failed: %v", err)
	}

	// We should have 3 laptops
	if len(laptops) != 3 {
		t.Fatalf("Expected 3 laptops, got %d", len(laptops))
	}

	// Verify laptop1 has reception report info with pending approval status
	foundLaptop1 := false
	foundLaptop2 := false
	foundLaptop3 := false

	for _, laptop := range laptops {
		if laptop.ID == laptop1.ID {
			foundLaptop1 = true
			if !laptop.HasReceptionReport {
				t.Error("Logistics user: Expected laptop1 to have HasReceptionReport = true")
			}
			if laptop.ReceptionReportStatus != string(ReceptionReportStatusPendingApproval) {
				t.Errorf("Logistics user: Expected laptop1 to have ReceptionReportStatus = 'pending_approval', got '%s'", laptop.ReceptionReportStatus)
			}
			if laptop.ReceptionReportID == nil || *laptop.ReceptionReportID != receptionReport1.ID {
				t.Errorf("Logistics user: Expected laptop1 to have ReceptionReportID = %d, got %v", receptionReport1.ID, laptop.ReceptionReportID)
			}
			if laptop.Status != LaptopStatusAtWarehouse {
				t.Errorf("Logistics user: Expected laptop1 status = 'at_warehouse', got '%s'", laptop.Status)
			}
		}

		if laptop.ID == laptop2.ID {
			foundLaptop2 = true
			if laptop.HasReceptionReport {
				t.Error("Logistics user: Expected laptop2 to have HasReceptionReport = false")
			}
			if laptop.ReceptionReportStatus != "" {
				t.Errorf("Logistics user: Expected laptop2 to have empty ReceptionReportStatus, got '%s'", laptop.ReceptionReportStatus)
			}
			if laptop.ReceptionReportID != nil {
				t.Errorf("Logistics user: Expected laptop2 to have ReceptionReportID = nil, got %v", laptop.ReceptionReportID)
			}
		}

		if laptop.ID == laptop3.ID {
			foundLaptop3 = true
			if !laptop.HasReceptionReport {
				t.Error("Logistics user: Expected laptop3 to have HasReceptionReport = true")
			}
			if laptop.ReceptionReportStatus != string(ReceptionReportStatusApproved) {
				t.Errorf("Logistics user: Expected laptop3 to have ReceptionReportStatus = 'approved', got '%s'", laptop.ReceptionReportStatus)
			}
			if laptop.ReceptionReportID == nil || *laptop.ReceptionReportID != receptionReport3.ID {
				t.Errorf("Logistics user: Expected laptop3 to have ReceptionReportID = %d, got %v", receptionReport3.ID, laptop.ReceptionReportID)
			}
		}
	}

	if !foundLaptop1 || !foundLaptop2 || !foundLaptop3 {
		t.Error("Not all test laptops were returned for logistics user")
	}
}
