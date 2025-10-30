package models

import (
	"encoding/json"
	"testing"
)

// PickupForm tests

func TestPickupForm_Validate(t *testing.T) {
	tests := []struct {
		name    string
		form    PickupForm
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid pickup form",
			form: PickupForm{
				ShipmentID:     1,
				SubmittedByUserID: 5,
				FormData:       json.RawMessage(`{"field": "value"}`),
			},
			wantErr: false,
		},
		{
			name: "invalid - missing shipment ID",
			form: PickupForm{
				SubmittedByUserID: 5,
				FormData:       json.RawMessage(`{"field": "value"}`),
			},
			wantErr: true,
			errMsg:  "shipment ID is required",
		},
		{
			name: "invalid - missing submitted by user ID",
			form: PickupForm{
				ShipmentID: 1,
				FormData:   json.RawMessage(`{"field": "value"}`),
			},
			wantErr: true,
			errMsg:  "submitted by user ID is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.form.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("PickupForm.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" && err.Error() != tt.errMsg {
				t.Errorf("PickupForm.Validate() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestPickupForm_TableName(t *testing.T) {
	form := PickupForm{}
	expected := "pickup_forms"
	if got := form.TableName(); got != expected {
		t.Errorf("PickupForm.TableName() = %v, want %v", got, expected)
	}
}

// ReceptionReport tests

func TestReceptionReport_Validate(t *testing.T) {
	tests := []struct {
		name    string
		report  ReceptionReport
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid reception report",
			report: ReceptionReport{
				ShipmentID:    1,
				WarehouseUserID: 10,
				Notes:         "All items received in good condition",
				PhotoURLs:     []string{"photo1.jpg", "photo2.jpg"},
			},
			wantErr: false,
		},
		{
			name: "valid - minimal fields",
			report: ReceptionReport{
				ShipmentID:    1,
				WarehouseUserID: 10,
			},
			wantErr: false,
		},
		{
			name: "invalid - missing shipment ID",
			report: ReceptionReport{
				WarehouseUserID: 10,
			},
			wantErr: true,
			errMsg:  "shipment ID is required",
		},
		{
			name: "invalid - missing warehouse user ID",
			report: ReceptionReport{
				ShipmentID: 1,
			},
			wantErr: true,
			errMsg:  "warehouse user ID is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.report.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReceptionReport.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" && err.Error() != tt.errMsg {
				t.Errorf("ReceptionReport.Validate() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestReceptionReport_TableName(t *testing.T) {
	report := ReceptionReport{}
	expected := "reception_reports"
	if got := report.TableName(); got != expected {
		t.Errorf("ReceptionReport.TableName() = %v, want %v", got, expected)
	}
}

func TestReceptionReport_HasPhotos(t *testing.T) {
	tests := []struct {
		name     string
		report   ReceptionReport
		expected bool
	}{
		{
			name: "report with photos",
			report: ReceptionReport{
				PhotoURLs: []string{"photo1.jpg", "photo2.jpg"},
			},
			expected: true,
		},
		{
			name: "report without photos",
			report: ReceptionReport{
				PhotoURLs: []string{},
			},
			expected: false,
		},
		{
			name: "report with nil photos",
			report: ReceptionReport{
				PhotoURLs: nil,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.report.HasPhotos(); got != tt.expected {
				t.Errorf("ReceptionReport.HasPhotos() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// DeliveryForm tests

func TestDeliveryForm_Validate(t *testing.T) {
	tests := []struct {
		name    string
		form    DeliveryForm
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid delivery form",
			form: DeliveryForm{
				ShipmentID: 1,
				EngineerID: 20,
				Notes:      "Delivered successfully",
				PhotoURLs:  []string{"delivery1.jpg"},
			},
			wantErr: false,
		},
		{
			name: "valid - minimal fields",
			form: DeliveryForm{
				ShipmentID: 1,
				EngineerID: 20,
			},
			wantErr: false,
		},
		{
			name: "invalid - missing shipment ID",
			form: DeliveryForm{
				EngineerID: 20,
			},
			wantErr: true,
			errMsg:  "shipment ID is required",
		},
		{
			name: "invalid - missing engineer ID",
			form: DeliveryForm{
				ShipmentID: 1,
			},
			wantErr: true,
			errMsg:  "engineer ID is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.form.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("DeliveryForm.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" && err.Error() != tt.errMsg {
				t.Errorf("DeliveryForm.Validate() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestDeliveryForm_TableName(t *testing.T) {
	form := DeliveryForm{}
	expected := "delivery_forms"
	if got := form.TableName(); got != expected {
		t.Errorf("DeliveryForm.TableName() = %v, want %v", got, expected)
	}
}

func TestDeliveryForm_HasPhotos(t *testing.T) {
	tests := []struct {
		name     string
		form     DeliveryForm
		expected bool
	}{
		{
			name: "form with photos",
			form: DeliveryForm{
				PhotoURLs: []string{"photo1.jpg"},
			},
			expected: true,
		},
		{
			name: "form without photos",
			form: DeliveryForm{
				PhotoURLs: []string{},
			},
			expected: false,
		},
		{
			name: "form with nil photos",
			form: DeliveryForm{
				PhotoURLs: nil,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.form.HasPhotos(); got != tt.expected {
				t.Errorf("DeliveryForm.HasPhotos() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// Test BeforeCreate and BeforeUpdate for all forms

func TestForms_BeforeCreate(t *testing.T) {
	t.Run("PickupForm", func(t *testing.T) {
		form := &PickupForm{ShipmentID: 1, SubmittedByUserID: 5}
		form.BeforeCreate()
		if form.SubmittedAt.IsZero() {
			t.Error("PickupForm.BeforeCreate() did not set SubmittedAt")
		}
	})

	t.Run("ReceptionReport", func(t *testing.T) {
		report := &ReceptionReport{ShipmentID: 1, WarehouseUserID: 10}
		report.BeforeCreate()
		if report.ReceivedAt.IsZero() {
			t.Error("ReceptionReport.BeforeCreate() did not set ReceivedAt")
		}
	})

	t.Run("DeliveryForm", func(t *testing.T) {
		form := &DeliveryForm{ShipmentID: 1, EngineerID: 20}
		form.BeforeCreate()
		if form.DeliveredAt.IsZero() {
			t.Error("DeliveryForm.BeforeCreate() did not set DeliveredAt")
		}
	})
}

