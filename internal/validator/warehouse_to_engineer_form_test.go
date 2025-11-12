package validator

import (
	"strings"
	"testing"
)

func TestValidateWarehouseToEngineerForm(t *testing.T) {
	tests := []struct {
		name          string
		input         WarehouseToEngineerFormInput
		shouldBeValid bool
		errorContains string
	}{
		{
			name: "valid warehouse to engineer form",
			input: WarehouseToEngineerFormInput{
				LaptopID:            1,
				SoftwareEngineerID:  5,
				EngineerName:        "Jane Smith",
				EngineerEmail:       "jane@bairesdev.com",
				EngineerAddress:     "456 Tech Ave",
				EngineerCity:        "San Francisco",
				EngineerState:       "CA",
				EngineerZip:         "94102",
				CourierName:         "FedEx",
				TrackingNumber:      "123456789",
				JiraTicketNumber:    "SCOP-12345",
			},
			shouldBeValid: true,
		},
		{
			name: "invalid - missing laptop ID",
			input: WarehouseToEngineerFormInput{
				LaptopID:         0, // Missing
				SoftwareEngineerID: 5,
				EngineerName:     "Jane Smith",
				EngineerEmail:    "jane@bairesdev.com",
				EngineerAddress:  "456 Tech Ave",
				EngineerCity:     "San Francisco",
				EngineerState:    "CA",
				EngineerZip:      "94102",
				CourierName:      "FedEx",
				TrackingNumber:   "123456789",
				JiraTicketNumber: "SCOP-12345",
			},
			shouldBeValid: false,
			errorContains: "laptop selection is required",
		},
		{
			name: "invalid - missing engineer ID and name",
			input: WarehouseToEngineerFormInput{
				LaptopID:         1,
				SoftwareEngineerID: 0,  // Missing
				EngineerName:     "",   // Missing
				EngineerEmail:    "jane@bairesdev.com",
				EngineerAddress:  "456 Tech Ave",
				EngineerCity:     "San Francisco",
				EngineerState:    "CA",
				EngineerZip:      "94102",
				CourierName:      "FedEx",
				TrackingNumber:   "123456789",
				JiraTicketNumber: "SCOP-12345",
			},
			shouldBeValid: false,
			errorContains: "software engineer is required",
		},
		{
			name: "valid - engineer ID without name",
			input: WarehouseToEngineerFormInput{
				LaptopID:           1,
				SoftwareEngineerID: 5,
				EngineerName:       "", // Optional when ID is provided
				EngineerEmail:      "jane@bairesdev.com",
				EngineerAddress:    "456 Tech Ave",
				EngineerCity:       "San Francisco",
				EngineerState:      "CA",
				EngineerZip:        "94102",
				CourierName:        "FedEx",
				TrackingNumber:     "123456789",
				JiraTicketNumber:   "SCOP-12345",
			},
			shouldBeValid: true,
		},
		{
			name: "valid - engineer name without ID",
			input: WarehouseToEngineerFormInput{
				LaptopID:           1,
				SoftwareEngineerID: 0, // Optional when name is provided
				EngineerName:       "Jane Smith",
				EngineerEmail:      "jane@bairesdev.com",
				EngineerAddress:    "456 Tech Ave",
				EngineerCity:       "San Francisco",
				EngineerState:      "CA",
				EngineerZip:        "94102",
				CourierName:        "FedEx",
				TrackingNumber:     "123456789",
				JiraTicketNumber:   "SCOP-12345",
			},
			shouldBeValid: true,
		},
		{
			name: "invalid - missing engineer address",
			input: WarehouseToEngineerFormInput{
				LaptopID:           1,
				SoftwareEngineerID: 5,
				EngineerName:       "Jane Smith",
				EngineerEmail:      "jane@bairesdev.com",
				EngineerAddress:    "", // Missing
				EngineerCity:       "San Francisco",
				EngineerState:      "CA",
				EngineerZip:        "94102",
				CourierName:        "FedEx",
				TrackingNumber:     "123456789",
				JiraTicketNumber:   "SCOP-12345",
			},
			shouldBeValid: false,
			errorContains: "address",
		},
		{
			name: "invalid - missing engineer city",
			input: WarehouseToEngineerFormInput{
				LaptopID:           1,
				SoftwareEngineerID: 5,
				EngineerName:       "Jane Smith",
				EngineerEmail:      "jane@bairesdev.com",
				EngineerAddress:    "456 Tech Ave",
				EngineerCity:       "", // Missing
				EngineerState:      "CA",
				EngineerZip:        "94102",
				CourierName:        "FedEx",
				TrackingNumber:     "123456789",
				JiraTicketNumber:   "SCOP-12345",
			},
			shouldBeValid: false,
			errorContains: "city",
		},
		{
			name: "invalid - missing engineer state",
			input: WarehouseToEngineerFormInput{
				LaptopID:           1,
				SoftwareEngineerID: 5,
				EngineerName:       "Jane Smith",
				EngineerEmail:      "jane@bairesdev.com",
				EngineerAddress:    "456 Tech Ave",
				EngineerCity:       "San Francisco",
				EngineerState:      "", // Missing
				EngineerZip:        "94102",
				CourierName:        "FedEx",
				TrackingNumber:     "123456789",
				JiraTicketNumber:   "SCOP-12345",
			},
			shouldBeValid: false,
			errorContains: "state",
		},
		{
			name: "invalid - missing engineer zip",
			input: WarehouseToEngineerFormInput{
				LaptopID:           1,
				SoftwareEngineerID: 5,
				EngineerName:       "Jane Smith",
				EngineerEmail:      "jane@bairesdev.com",
				EngineerAddress:    "456 Tech Ave",
				EngineerCity:       "San Francisco",
				EngineerState:      "CA",
				EngineerZip:        "", // Missing
				CourierName:        "FedEx",
				TrackingNumber:     "123456789",
				JiraTicketNumber:   "SCOP-12345",
			},
			shouldBeValid: false,
			errorContains: "ZIP",
		},
		{
			name: "invalid - invalid state code",
			input: WarehouseToEngineerFormInput{
				LaptopID:           1,
				SoftwareEngineerID: 5,
				EngineerName:       "Jane Smith",
				EngineerEmail:      "jane@bairesdev.com",
				EngineerAddress:    "456 Tech Ave",
				EngineerCity:       "San Francisco",
				EngineerState:      "XX", // Invalid
				EngineerZip:        "94102",
				CourierName:        "FedEx",
				TrackingNumber:     "123456789",
				JiraTicketNumber:   "SCOP-12345",
			},
			shouldBeValid: false,
			errorContains: "state",
		},
		{
			name: "valid - courier info optional",
			input: WarehouseToEngineerFormInput{
				LaptopID:           1,
				SoftwareEngineerID: 5,
				EngineerName:       "Jane Smith",
				EngineerEmail:      "jane@bairesdev.com",
				EngineerAddress:    "456 Tech Ave",
				EngineerCity:       "San Francisco",
				EngineerState:      "CA",
				EngineerZip:        "94102",
				CourierName:        "", // Optional
				TrackingNumber:     "", // Optional
				JiraTicketNumber:   "SCOP-12345",
			},
			shouldBeValid: true,
		},
		{
			name: "invalid - missing JIRA ticket",
			input: WarehouseToEngineerFormInput{
				LaptopID:           1,
				SoftwareEngineerID: 5,
				EngineerName:       "Jane Smith",
				EngineerEmail:      "jane@bairesdev.com",
				EngineerAddress:    "456 Tech Ave",
				EngineerCity:       "San Francisco",
				EngineerState:      "CA",
				EngineerZip:        "94102",
				CourierName:        "FedEx",
				TrackingNumber:     "123456789",
				JiraTicketNumber:   "", // Missing
			},
			shouldBeValid: false,
			errorContains: "JIRA ticket",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateWarehouseToEngineerForm(tt.input)
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

