package models

import (
	"testing"
	"time"
)

func TestLaptop_Validate(t *testing.T) {
	clientID := int64(1)
	engineerID := int64(10)

	tests := []struct {
		name    string
		laptop  Laptop
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid laptop with all fields including SKU, client, and engineer",
			laptop: Laptop{
				SerialNumber:       "SN123456789",
				SKU:                "SKU-DELL-LAT-5520",
				Brand:              "Dell",
				Model:              "Latitude 5520",
				Specs:              "i7, 16GB RAM, 512GB SSD",
				Status:             LaptopStatusAvailable,
				ClientCompanyID:    &clientID,
				SoftwareEngineerID: &engineerID,
			},
			wantErr: false,
		},
		{
			name: "valid laptop with all fields",
			laptop: Laptop{
				SerialNumber: "SN123456789",
				Brand:        "Dell",
				Model:        "Latitude 5520",
				Specs:        "i7, 16GB RAM, 512GB SSD",
				Status:       LaptopStatusAvailable,
			},
			wantErr: false,
		},
		{
			name: "valid laptop with minimal fields",
			laptop: Laptop{
				SerialNumber: "SN987654321",
				Status:       LaptopStatusAtWarehouse,
			},
			wantErr: false,
		},
		{
			name: "invalid - missing serial number",
			laptop: Laptop{
				Brand:  "HP",
				Model:  "EliteBook",
				Status: LaptopStatusAvailable,
			},
			wantErr: true,
			errMsg:  "serial number is required",
		},
		{
			name: "invalid - empty serial number",
			laptop: Laptop{
				SerialNumber: "",
				Status:       LaptopStatusAvailable,
			},
			wantErr: true,
			errMsg:  "serial number is required",
		},
		{
			name: "invalid - missing status",
			laptop: Laptop{
				SerialNumber: "SN123456789",
			},
			wantErr: true,
			errMsg:  "status is required",
		},
		{
			name: "invalid - invalid status",
			laptop: Laptop{
				SerialNumber: "SN123456789",
				Status:       "invalid_status",
			},
			wantErr: true,
			errMsg:  "invalid status",
		},
		{
			name: "valid - serial number exactly 3 characters",
			laptop: Laptop{
				SerialNumber: "ABC",
				Status:       LaptopStatusAvailable,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.laptop.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Laptop.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" && err.Error() != tt.errMsg {
				t.Errorf("Laptop.Validate() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestLaptop_IsValidStatus(t *testing.T) {
	tests := []struct {
		name   string
		status LaptopStatus
		want   bool
	}{
		{"available", LaptopStatusAvailable, true},
		{"in_transit_to_warehouse", LaptopStatusInTransitToWarehouse, true},
		{"at_warehouse", LaptopStatusAtWarehouse, true},
		{"in_transit_to_engineer", LaptopStatusInTransitToEngineer, true},
		{"delivered", LaptopStatusDelivered, true},
		{"retired", LaptopStatusRetired, true},
		{"invalid status", "unknown", false},
		{"empty status", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidLaptopStatus(tt.status); got != tt.want {
				t.Errorf("IsValidLaptopStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLaptop_TableName(t *testing.T) {
	laptop := Laptop{}
	expected := "laptops"
	if got := laptop.TableName(); got != expected {
		t.Errorf("Laptop.TableName() = %v, want %v", got, expected)
	}
}

func TestLaptop_BeforeCreate(t *testing.T) {
	laptop := &Laptop{
		SerialNumber: "SN123456789",
		Status:       LaptopStatusAvailable,
	}

	laptop.BeforeCreate()

	// Check that timestamps are set
	if laptop.CreatedAt.IsZero() {
		t.Error("Laptop.BeforeCreate() did not set CreatedAt")
	}
	if laptop.UpdatedAt.IsZero() {
		t.Error("Laptop.BeforeCreate() did not set UpdatedAt")
	}

	// Check that CreatedAt and UpdatedAt are approximately equal (within 1 second)
	diff := laptop.UpdatedAt.Sub(laptop.CreatedAt)
	if diff < 0 || diff > time.Second {
		t.Errorf("Laptop.BeforeCreate() CreatedAt and UpdatedAt differ by %v, expected them to be nearly equal", diff)
	}
}

func TestLaptop_BeforeUpdate(t *testing.T) {
	laptop := &Laptop{
		SerialNumber: "SN123456789",
		Status:       LaptopStatusAvailable,
		CreatedAt:    time.Now().Add(-24 * time.Hour),
		UpdatedAt:    time.Now().Add(-24 * time.Hour),
	}

	oldUpdatedAt := laptop.UpdatedAt
	time.Sleep(10 * time.Millisecond)

	laptop.BeforeUpdate()

	// Check that UpdatedAt was updated
	if !laptop.UpdatedAt.After(oldUpdatedAt) {
		t.Error("Laptop.BeforeUpdate() did not update UpdatedAt")
	}
}

func TestLaptop_IsAvailable(t *testing.T) {
	tests := []struct {
		name     string
		laptop   Laptop
		expected bool
	}{
		{
			name: "available laptop",
			laptop: Laptop{
				SerialNumber: "SN123",
				Status:       LaptopStatusAvailable,
			},
			expected: true,
		},
		{
			name: "delivered laptop",
			laptop: Laptop{
				SerialNumber: "SN456",
				Status:       LaptopStatusDelivered,
			},
			expected: false,
		},
		{
			name: "in transit laptop",
			laptop: Laptop{
				SerialNumber: "SN789",
				Status:       LaptopStatusInTransitToWarehouse,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.laptop.IsAvailable(); got != tt.expected {
				t.Errorf("Laptop.IsAvailable() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestLaptop_UpdateStatus(t *testing.T) {
	laptop := &Laptop{
		SerialNumber: "SN123456789",
		Status:       LaptopStatusAvailable,
	}

	if laptop.Status != LaptopStatusAvailable {
		t.Error("Expected initial status to be available")
	}

	laptop.UpdateStatus(LaptopStatusInTransitToWarehouse)

	if laptop.Status != LaptopStatusInTransitToWarehouse {
		t.Errorf("UpdateStatus() did not update status, got %v, want %v", laptop.Status, LaptopStatusInTransitToWarehouse)
	}
}

func TestLaptop_GetFullDescription(t *testing.T) {
	tests := []struct {
		name     string
		laptop   Laptop
		expected string
	}{
		{
			name: "laptop with all details",
			laptop: Laptop{
				Brand: "Dell",
				Model: "Latitude 5520",
				Specs: "i7, 16GB RAM",
			},
			expected: "Dell Latitude 5520 (i7, 16GB RAM)",
		},
		{
			name: "laptop with brand and model only",
			laptop: Laptop{
				Brand: "HP",
				Model: "EliteBook",
			},
			expected: "HP EliteBook",
		},
		{
			name: "laptop with brand only",
			laptop: Laptop{
				Brand: "Lenovo",
			},
			expected: "Lenovo",
		},
		{
			name: "laptop with no details",
			laptop: Laptop{
				SerialNumber: "SN123",
			},
			expected: "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.laptop.GetFullDescription(); got != tt.expected {
				t.Errorf("Laptop.GetFullDescription() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestLaptop_WithNewFields(t *testing.T) {
	clientID := int64(1)
	engineerID := int64(10)

	laptop := Laptop{
		SerialNumber:       "SN123456789",
		SKU:                "SKU-DELL-LAT-5520",
		Brand:              "Dell",
		Model:              "Latitude 5520",
		Status:             LaptopStatusAvailable,
		ClientCompanyID:    &clientID,
		SoftwareEngineerID: &engineerID,
	}

	// Test that fields are properly set
	if laptop.SKU != "SKU-DELL-LAT-5520" {
		t.Errorf("Expected SKU to be 'SKU-DELL-LAT-5520', got %s", laptop.SKU)
	}

	if laptop.ClientCompanyID == nil || *laptop.ClientCompanyID != 1 {
		t.Errorf("Expected ClientCompanyID to be 1, got %v", laptop.ClientCompanyID)
	}

	if laptop.SoftwareEngineerID == nil || *laptop.SoftwareEngineerID != 10 {
		t.Errorf("Expected SoftwareEngineerID to be 10, got %v", laptop.SoftwareEngineerID)
	}
}

