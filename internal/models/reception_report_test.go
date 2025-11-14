package models

import (
	"testing"
)

// TestReceptionReport_LaptopBasedStructure verifies the reception report model has laptop-based fields
func TestReceptionReport_LaptopBasedStructure(t *testing.T) {
	shipmentID := int64(100)
	clientCompanyID := int64(5)
	
	report := &ReceptionReport{
		LaptopID:               1,
		ShipmentID:             &shipmentID, // Reference only
		ClientCompanyID:        &clientCompanyID,
		TrackingNumber:         "TRACK123",
		WarehouseUserID:        2,
		Notes:                  "Test note",
		PhotoSerialNumber:      "/uploads/serial.jpg",
		PhotoExternalCondition: "/uploads/external.jpg",
		PhotoWorkingCondition:  "/uploads/working.jpg",
		Status:                 ReceptionReportStatusPendingApproval,
	}

	if report.LaptopID == 0 {
		t.Error("Expected LaptopID to be set")
	}
	if report.ShipmentID == nil || *report.ShipmentID == 0 {
		t.Error("Expected ShipmentID reference to be set")
	}
	if report.ClientCompanyID == nil || *report.ClientCompanyID == 0 {
		t.Error("Expected ClientCompanyID to be set")
	}
	if report.TrackingNumber == "" {
		t.Error("Expected TrackingNumber to be set")
	}
	if report.PhotoSerialNumber == "" {
		t.Error("Expected PhotoSerialNumber to be set")
	}
	if report.PhotoExternalCondition == "" {
		t.Error("Expected PhotoExternalCondition to be set")
	}
	if report.PhotoWorkingCondition == "" {
		t.Error("Expected PhotoWorkingCondition to be set")
	}
	if report.Status == "" {
		t.Error("Expected Status to be set")
	}
}

// TestReceptionReport_ValidationRequiresLaptopID tests that validation requires laptop_id
func TestReceptionReport_ValidationRequiresLaptopID(t *testing.T) {
	report := &ReceptionReport{
		WarehouseUserID: 1,
	}

	err := report.Validate()
	if err == nil {
		t.Error("Expected validation error for missing laptop ID")
	}
}

// TestReceptionReport_ValidationRequiresPhotos tests that validation requires all 3 photos
func TestReceptionReport_ValidationRequiresPhotos(t *testing.T) {
	tests := []struct {
		name   string
		report *ReceptionReport
	}{
		{
			name: "Missing serial number photo",
			report: &ReceptionReport{
				LaptopID:               1,
				WarehouseUserID:        1,
				PhotoExternalCondition: "/uploads/external.jpg",
				PhotoWorkingCondition:  "/uploads/working.jpg",
			},
		},
		{
			name: "Missing external condition photo",
			report: &ReceptionReport{
				LaptopID:              1,
				WarehouseUserID:       1,
				PhotoSerialNumber:     "/uploads/serial.jpg",
				PhotoWorkingCondition: "/uploads/working.jpg",
			},
		},
		{
			name: "Missing working condition photo",
			report: &ReceptionReport{
				LaptopID:               1,
				WarehouseUserID:        1,
				PhotoSerialNumber:      "/uploads/serial.jpg",
				PhotoExternalCondition: "/uploads/external.jpg",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.report.Validate()
			if err == nil {
				t.Errorf("Expected validation error for %s", tt.name)
			}
		})
	}
}

// TestReceptionReport_StatusConstants verifies status constants exist
func TestReceptionReport_StatusConstants(t *testing.T) {
	if ReceptionReportStatusPendingApproval == "" {
		t.Error("Expected ReceptionReportStatusPendingApproval constant to exist")
	}
	if ReceptionReportStatusApproved == "" {
		t.Error("Expected ReceptionReportStatusApproved constant to exist")
	}
}

// TestReceptionReport_IsPendingApproval tests pending approval check
func TestReceptionReport_IsPendingApproval(t *testing.T) {
	report := &ReceptionReport{
		Status: ReceptionReportStatusPendingApproval,
	}

	if !report.IsPendingApproval() {
		t.Error("Expected report to be pending approval")
	}
}

// TestReceptionReport_IsApproved tests approved status check
func TestReceptionReport_IsApproved(t *testing.T) {
	report := &ReceptionReport{
		Status: ReceptionReportStatusApproved,
	}

	if !report.IsApproved() {
		t.Error("Expected report to be approved")
	}
}

// TestReceptionReport_ApproveMethod tests the Approve method
func TestReceptionReport_ApproveMethod(t *testing.T) {
	report := &ReceptionReport{
		LaptopID:        1,
		WarehouseUserID: 2,
		Status:          ReceptionReportStatusPendingApproval,
	}

	logisticsUserID := int64(10)
	report.Approve(logisticsUserID)

	if !report.IsApproved() {
		t.Error("Expected report to be approved after calling Approve()")
	}
	if report.ApprovedBy == nil || *report.ApprovedBy != logisticsUserID {
		t.Error("Expected ApprovedBy to be set to logistics user ID")
	}
	if report.ApprovedAt == nil {
		t.Error("Expected ApprovedAt timestamp to be set")
	}
}
