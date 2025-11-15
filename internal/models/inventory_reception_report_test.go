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

	// Create test laptops with different reception report scenarios
	// Laptop 1: Has reception report (pending approval)
	laptop1 := &Laptop{
		SerialNumber: "SN001",
		Brand:        "Dell",
		Model:        "Latitude 5520",
		RAMGB:        "16",
		SSDGB:        "512",
		Status:       LaptopStatusAtWarehouse,
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
		SerialNumber: "SN002",
		Brand:        "HP",
		Model:        "EliteBook 840",
		RAMGB:        "32",
		SSDGB:        "1024",
		Status:       LaptopStatusAtWarehouse,
	}
	err = CreateLaptop(db, laptop2)
	if err != nil {
		t.Fatalf("Failed to create laptop2: %v", err)
	}

	// Laptop 3: Has reception report (approved)
	laptop3 := &Laptop{
		SerialNumber: "SN003",
		Brand:        "Lenovo",
		Model:        "ThinkPad X1",
		RAMGB:        "16",
		SSDGB:        "512",
		Status:       LaptopStatusAvailable,
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

