package validator

import (
	"testing"
)

// TestValidateCompleteShipmentDetails_LaptopModelRequired tests that laptop model is required
func TestValidateCompleteShipmentDetails_LaptopModelRequired(t *testing.T) {
	input := CompleteShipmentDetailsInput{
		ShipmentID:         1,
		ContactName:        "John Doe",
		ContactEmail:       "john@example.com",
		ContactPhone:       "+1-555-0123",
		PickupAddress:      "123 Main St",
		PickupCity:         "New York",
		PickupState:        "NY",
		PickupZip:          "10001",
		PickupDate:         "2025-12-15",
		PickupTimeSlot:     "morning",
		LaptopSerialNumber: "ABC123456789",
		LaptopBrand:        "Dell",
		// Missing LaptopModel
		LaptopCPU:          "Intel Core i7",
		LaptopRAMGB:        "16GB",
		LaptopSSDGB:        "512GB",
	}

	err := ValidateCompleteShipmentDetails(input)
	if err == nil {
		t.Error("Expected error for missing laptop model, got nil")
	}
	if err != nil && err.Error() != "laptop model is required" {
		t.Errorf("Expected 'laptop model is required', got '%s'", err.Error())
	}
}

// TestValidateCompleteShipmentDetails_LaptopRAMRequired tests that laptop RAM is required
func TestValidateCompleteShipmentDetails_LaptopRAMRequired(t *testing.T) {
	input := CompleteShipmentDetailsInput{
		ShipmentID:         1,
		ContactName:        "John Doe",
		ContactEmail:       "john@example.com",
		ContactPhone:       "+1-555-0123",
		PickupAddress:      "123 Main St",
		PickupCity:         "New York",
		PickupState:        "NY",
		PickupZip:          "10001",
		PickupDate:         "2025-12-15",
		PickupTimeSlot:     "morning",
		LaptopSerialNumber: "ABC123456789",
		LaptopBrand:        "Dell",
		LaptopModel:        "Dell XPS 15",
		LaptopCPU:          "Intel Core i7",
		// Missing LaptopRAMGB
		LaptopSSDGB: "512GB",
	}

	err := ValidateCompleteShipmentDetails(input)
	if err == nil {
		t.Error("Expected error for missing laptop RAM, got nil")
	}
	if err != nil && err.Error() != "laptop RAM is required" {
		t.Errorf("Expected 'laptop RAM is required', got '%s'", err.Error())
	}
}

// TestValidateCompleteShipmentDetails_LaptopSSDRequired tests that laptop SSD is required
func TestValidateCompleteShipmentDetails_LaptopSSDRequired(t *testing.T) {
	input := CompleteShipmentDetailsInput{
		ShipmentID:         1,
		ContactName:        "John Doe",
		ContactEmail:       "john@example.com",
		ContactPhone:       "+1-555-0123",
		PickupAddress:      "123 Main St",
		PickupCity:         "New York",
		PickupState:        "NY",
		PickupZip:          "10001",
		PickupDate:         "2025-12-15",
		PickupTimeSlot:     "morning",
		LaptopSerialNumber: "ABC123456789",
		LaptopBrand:        "Dell",
		LaptopModel:        "Dell XPS 15",
		LaptopCPU:          "Intel Core i7",
		LaptopRAMGB:        "16GB",
		// Missing LaptopSSDGB
	}

	err := ValidateCompleteShipmentDetails(input)
	if err == nil {
		t.Error("Expected error for missing laptop SSD, got nil")
	}
	if err != nil && err.Error() != "laptop SSD is required" {
		t.Errorf("Expected 'laptop SSD is required', got '%s'", err.Error())
	}
}

// TestValidateCompleteShipmentDetails_AllLaptopFieldsProvided tests validation passes with all fields
func TestValidateCompleteShipmentDetails_AllLaptopFieldsProvided(t *testing.T) {
	input := CompleteShipmentDetailsInput{
		ShipmentID:         1,
		ContactName:        "John Doe",
		ContactEmail:       "john@example.com",
		ContactPhone:       "+1-555-0123",
		PickupAddress:      "123 Main St",
		PickupCity:         "New York",
		PickupState:        "NY",
		PickupZip:          "10001",
		PickupDate:         "2025-12-15",
		PickupTimeSlot:     "morning",
		LaptopSerialNumber: "ABC123456789",
		LaptopBrand:        "Dell",
		LaptopModel:        "Dell XPS 15",
		LaptopCPU:          "Intel Core i7",
		LaptopRAMGB:        "16GB",
		LaptopSSDGB:        "512GB",
	}

	err := ValidateCompleteShipmentDetails(input)
	if err != nil {
		t.Errorf("Expected no error with all required fields, got: %s", err.Error())
	}
}

