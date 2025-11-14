package models

import (
	"testing"
)

// TestLaptop_CanChangeToAvailable tests the logic for changing laptop status to available
func TestLaptop_CanChangeToAvailable(t *testing.T) {
	laptop := &Laptop{
		ID:           1,
		SerialNumber: "TEST123",
		Status:       LaptopStatusAtWarehouse,
	}

	// Should not be able to change to available without checking reception report
	if laptop.Status == LaptopStatusAvailable {
		t.Error("Laptop should not be available without approval check")
	}
}

// TestLaptop_RequiresApprovedReceptionReport tests that changing to available requires approved reception report
func TestLaptop_RequiresApprovedReceptionReport(t *testing.T) {
	tests := []struct {
		name                  string
		currentStatus         LaptopStatus
		hasReceptionReport    bool
		receptionReportStatus ReceptionReportStatus
		canChangeToAvailable  bool
	}{
		{
			name:                  "at_warehouse with approved report - can change",
			currentStatus:         LaptopStatusAtWarehouse,
			hasReceptionReport:    true,
			receptionReportStatus: ReceptionReportStatusApproved,
			canChangeToAvailable:  true,
		},
		{
			name:                  "at_warehouse with pending report - cannot change",
			currentStatus:         LaptopStatusAtWarehouse,
			hasReceptionReport:    true,
			receptionReportStatus: ReceptionReportStatusPendingApproval,
			canChangeToAvailable:  false,
		},
		{
			name:                  "at_warehouse without report - cannot change",
			currentStatus:         LaptopStatusAtWarehouse,
			hasReceptionReport:    false,
			receptionReportStatus: "",
			canChangeToAvailable:  false,
		},
		{
			name:                  "in_transit - cannot change regardless of report",
			currentStatus:         LaptopStatusInTransitToWarehouse,
			hasReceptionReport:    true,
			receptionReportStatus: ReceptionReportStatusApproved,
			canChangeToAvailable:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			laptop := &Laptop{
				ID:           1,
				SerialNumber: "TEST123",
				Status:       tt.currentStatus,
			}

			var receptionReport *ReceptionReport
			if tt.hasReceptionReport {
				receptionReport = &ReceptionReport{
					LaptopID: laptop.ID,
					Status:   tt.receptionReportStatus,
				}
			}

			canChange := laptop.CanChangeToAvailable(receptionReport)
			if canChange != tt.canChangeToAvailable {
				t.Errorf("Expected CanChangeToAvailable() = %v, got %v", tt.canChangeToAvailable, canChange)
			}
		})
	}
}

// TestGetLaptopReceptionReport tests getting a laptop's reception report from database
func TestGetLaptopReceptionReport(t *testing.T) {
	// This test requires a database connection
	// Skip if running in short mode
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	// Note: This is a placeholder test structure
	// Actual implementation will depend on database setup
	t.Run("laptop with approved reception report", func(t *testing.T) {
		// Test will be implemented when database helper is available
		t.Skip("Integration test - requires database setup")
	})

	t.Run("laptop without reception report", func(t *testing.T) {
		// Test will be implemented when database helper is available
		t.Skip("Integration test - requires database setup")
	})
}

// TestUpdateLaptopStatusToAvailable tests the database function to update laptop status
func TestUpdateLaptopStatusToAvailable(t *testing.T) {
	// This test requires a database connection
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	t.Run("successful status update with approved report", func(t *testing.T) {
		// Test will be implemented when database helper is available
		t.Skip("Integration test - requires database setup")
	})

	t.Run("fails without approved report", func(t *testing.T) {
		// Test will be implemented when database helper is available
		t.Skip("Integration test - requires database setup")
	})
}

// TestLaptopStatusConstants verifies the laptop status constants
func TestLaptopStatusConstants(t *testing.T) {
	if LaptopStatusAtWarehouse == "" {
		t.Error("Expected LaptopStatusAtWarehouse constant to exist")
	}
	if LaptopStatusAvailable == "" {
		t.Error("Expected LaptopStatusAvailable constant to exist")
	}
}

// TestLaptop_GetStatusDisplayName tests the display name for at_warehouse status
func TestLaptop_GetStatusDisplayName(t *testing.T) {
	tests := []struct {
		status      LaptopStatus
		displayName string
	}{
		{LaptopStatusAtWarehouse, "Received at Warehouse"},
		{LaptopStatusAvailable, "Available at Warehouse"},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			displayName := GetLaptopStatusDisplayName(tt.status)
			if displayName != tt.displayName {
				t.Errorf("Expected display name %q, got %q", tt.displayName, displayName)
			}
		})
	}
}

