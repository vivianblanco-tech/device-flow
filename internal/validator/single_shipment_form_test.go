package validator

import (
	"strings"
	"testing"
)

func TestValidateSingleFullJourneyForm(t *testing.T) {
	tests := []struct {
		name          string
		input         SingleFullJourneyFormInput
		shouldBeValid bool
		errorContains string
	}{
		{
			name: "valid single full journey form",
			input: SingleFullJourneyFormInput{
				ClientCompanyID:     1,
				ContactName:         "John Doe",
				ContactEmail:        "john@company.com",
				ContactPhone:        "+1-555-0123",
				PickupAddress:       "123 Main St",
				PickupCity:          "New York",
				PickupState:         "NY",
				PickupZip:           "10001",
				PickupDate:          "2025-11-15",
				PickupTimeSlot:      "morning",
				JiraTicketNumber:    "SCOP-12345",
				LaptopSerialNumber:  "ABC123456",
				LaptopSpecs:         "Dell XPS 15, 16GB RAM, 512GB SSD",
				EngineerName:        "Jane Smith",
			},
			shouldBeValid: true,
		},
		{
			name: "missing laptop serial number",
			input: SingleFullJourneyFormInput{
				ClientCompanyID:  1,
				ContactName:      "John Doe",
				ContactEmail:     "john@company.com",
				ContactPhone:     "+1-555-0123",
				PickupAddress:    "123 Main St",
				PickupCity:       "New York",
				PickupState:      "NY",
				PickupZip:        "10001",
				PickupDate:       "2025-11-15",
				PickupTimeSlot:   "morning",
				JiraTicketNumber: "SCOP-12345",
				LaptopSpecs:      "Dell XPS 15",
				EngineerName:     "Jane Smith",
			},
			shouldBeValid: false,
			errorContains: "serial number is required",
		},
		{
			name: "engineer name optional",
			input: SingleFullJourneyFormInput{
				ClientCompanyID:    1,
				ContactName:        "John Doe",
				ContactEmail:       "john@company.com",
				ContactPhone:       "+1-555-0123",
				PickupAddress:      "123 Main St",
				PickupCity:         "New York",
				PickupState:        "NY",
				PickupZip:          "10001",
				PickupDate:         "2025-11-15",
				PickupTimeSlot:     "morning",
				JiraTicketNumber:   "SCOP-12345",
				LaptopSerialNumber: "ABC123456",
				LaptopSpecs:        "Dell XPS 15",
				EngineerName:       "", // Optional
			},
			shouldBeValid: true,
		},
		{
			name: "laptop specs optional",
			input: SingleFullJourneyFormInput{
				ClientCompanyID:    1,
				ContactName:        "John Doe",
				ContactEmail:       "john@company.com",
				ContactPhone:       "+1-555-0123",
				PickupAddress:      "123 Main St",
				PickupCity:         "New York",
				PickupState:        "NY",
				PickupZip:          "10001",
				PickupDate:         "2025-11-15",
				PickupTimeSlot:     "morning",
				JiraTicketNumber:   "SCOP-12345",
				LaptopSerialNumber: "ABC123456",
				LaptopSpecs:        "", // Optional
				EngineerName:       "Jane Smith",
			},
			shouldBeValid: true,
		},
		{
			name: "invalid email format",
			input: SingleFullJourneyFormInput{
				ClientCompanyID:    1,
				ContactName:        "John Doe",
				ContactEmail:       "invalid-email",
				ContactPhone:       "+1-555-0123",
				PickupAddress:      "123 Main St",
				PickupCity:         "New York",
				PickupState:        "NY",
				PickupZip:          "10001",
				PickupDate:         "2025-11-15",
				PickupTimeSlot:     "morning",
				JiraTicketNumber:   "SCOP-12345",
				LaptopSerialNumber: "ABC123456",
			},
			shouldBeValid: false,
			errorContains: "email",
		},
		{
			name: "missing client company",
			input: SingleFullJourneyFormInput{
				ClientCompanyID:    0,
				ContactName:        "John Doe",
				ContactEmail:       "john@company.com",
				ContactPhone:       "+1-555-0123",
				PickupAddress:      "123 Main St",
				PickupCity:         "New York",
				PickupState:        "NY",
				PickupZip:          "10001",
				PickupDate:         "2025-11-15",
				PickupTimeSlot:     "morning",
				JiraTicketNumber:   "SCOP-12345",
				LaptopSerialNumber: "ABC123456",
			},
			shouldBeValid: false,
			errorContains: "client company",
		},
		{
			name: "accessories description required when including accessories",
			input: SingleFullJourneyFormInput{
				ClientCompanyID:        1,
				ContactName:            "John Doe",
				ContactEmail:           "john@company.com",
				ContactPhone:           "+1-555-0123",
				PickupAddress:          "123 Main St",
				PickupCity:             "New York",
				PickupState:            "NY",
				PickupZip:              "10001",
				PickupDate:             "2025-11-15",
				PickupTimeSlot:         "morning",
				JiraTicketNumber:       "SCOP-12345",
				LaptopSerialNumber:     "ABC123456",
				IncludeAccessories:     true,
				AccessoriesDescription: "", // Missing when accessories included
			},
			shouldBeValid: false,
			errorContains: "accessories description is required",
		},
		{
			name: "valid with accessories",
			input: SingleFullJourneyFormInput{
				ClientCompanyID:        1,
				ContactName:            "John Doe",
				ContactEmail:           "john@company.com",
				ContactPhone:           "+1-555-0123",
				PickupAddress:          "123 Main St",
				PickupCity:             "New York",
				PickupState:            "NY",
				PickupZip:              "10001",
				PickupDate:             "2025-11-15",
				PickupTimeSlot:         "morning",
				JiraTicketNumber:       "SCOP-12345",
				LaptopSerialNumber:     "ABC123456",
				IncludeAccessories:     true,
				AccessoriesDescription: "Charger and mouse",
			},
			shouldBeValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSingleFullJourneyForm(tt.input)
			if tt.shouldBeValid && err != nil {
				t.Errorf("Expected valid, got error: %v", err)
			}
			if !tt.shouldBeValid {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain '%s', got: %v", tt.errorContains, err)
				}
			}
		})
	}
}

