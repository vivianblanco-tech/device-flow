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
				LaptopID:           1,
				SoftwareEngineerID: 5,
				EngineerName:       "Jane Smith",
				EngineerEmail:      "jane@bairesdev.com",
				EngineerAddress:    "456 Tech Ave",
				EngineerCity:       "San Francisco",
				EngineerCountry:    "United States",
				EngineerState:      "CA",
				EngineerZip:        "94102",
				CourierName:        "FedEx",
				TrackingNumber:     "123456789",
				JiraTicketNumber:   "SCOP-12345",
			},
			shouldBeValid: true,
		},
		{
			name: "invalid - missing laptop ID",
			input: WarehouseToEngineerFormInput{
				LaptopID:           0, // Missing
				SoftwareEngineerID: 5,
				EngineerName:       "Jane Smith",
				EngineerEmail:      "jane@bairesdev.com",
				EngineerAddress:    "456 Tech Ave",
				EngineerCity:       "San Francisco",
				EngineerCountry:    "United States",
				EngineerState:      "CA",
				EngineerZip:        "94102",
				CourierName:        "FedEx",
				TrackingNumber:     "123456789",
				JiraTicketNumber:   "SCOP-12345",
			},
			shouldBeValid: false,
			errorContains: "laptop selection is required",
		},
		{
			name: "invalid - missing engineer ID and name",
			input: WarehouseToEngineerFormInput{
				LaptopID:           1,
				SoftwareEngineerID: 0,  // Missing
				EngineerName:       "", // Missing
				EngineerEmail:      "jane@bairesdev.com",
				EngineerAddress:    "456 Tech Ave",
				EngineerCity:       "San Francisco",
				EngineerCountry:    "United States",
				EngineerState:      "CA",
				EngineerZip:        "94102",
				CourierName:        "FedEx",
				TrackingNumber:     "123456789",
				JiraTicketNumber:   "SCOP-12345",
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
				EngineerCountry:    "United States",
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
				EngineerCountry:    "United States",
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
				EngineerCountry:    "United States",
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
				EngineerCountry:    "United States",
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
			name: "valid - missing engineer state (optional for international)",
			input: WarehouseToEngineerFormInput{
				LaptopID:           1,
				SoftwareEngineerID: 5,
				EngineerName:       "Jane Smith",
				EngineerEmail:      "jane@bairesdev.com",
				EngineerAddress:    "456 Tech Ave",
				EngineerCity:       "Dublin",
				EngineerCountry:    "Ireland",
				EngineerState:      "", // Optional for international
				EngineerZip:        "D02 XY45",
				CourierName:        "FedEx",
				TrackingNumber:     "123456789",
				JiraTicketNumber:   "SCOP-12345",
			},
			shouldBeValid: true,
		},
		{
			name: "valid - missing engineer zip (optional for international)",
			input: WarehouseToEngineerFormInput{
				LaptopID:           1,
				SoftwareEngineerID: 5,
				EngineerName:       "Jane Smith",
				EngineerEmail:      "jane@bairesdev.com",
				EngineerAddress:    "456 Tech Ave",
				EngineerCity:       "Tokyo",
				EngineerCountry:    "Japan",
				EngineerState:      "Tokyo",
				EngineerZip:        "", // Optional for international
				CourierName:        "FedEx",
				TrackingNumber:     "123456789",
				JiraTicketNumber:   "SCOP-12345",
			},
			shouldBeValid: true,
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
				EngineerCountry:    "United States",
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
				EngineerCountry:    "United States",
				EngineerState:      "CA",
				EngineerZip:        "94102",
				CourierName:        "FedEx",
				TrackingNumber:     "123456789",
				JiraTicketNumber:   "", // Missing
			},
			shouldBeValid: false,
			errorContains: "JIRA ticket",
		},
		// International address format tests
		{
			name: "valid - international address with country, city, address (state and postal optional)",
			input: WarehouseToEngineerFormInput{
				LaptopID:           1,
				SoftwareEngineerID: 5,
				EngineerName:       "Jane Smith",
				EngineerEmail:      "jane@bairesdev.com",
				EngineerAddress:    "123 Main Street",
				EngineerCity:       "London",
				EngineerCountry:    "United Kingdom",
				EngineerState:      "", // Optional
				EngineerZip:        "", // Optional
				JiraTicketNumber:   "SCOP-12345",
			},
			shouldBeValid: true,
		},
		{
			name: "valid - international address with all fields including optional state and postal",
			input: WarehouseToEngineerFormInput{
				LaptopID:           1,
				SoftwareEngineerID: 5,
				EngineerName:       "Jane Smith",
				EngineerEmail:      "jane@bairesdev.com",
				EngineerAddress:    "456 Tech Avenue, Apt 12B",
				EngineerCity:       "Buenos Aires",
				EngineerCountry:    "Argentina",
				EngineerState:      "Buenos Aires", // Optional
				EngineerZip:        "C1000ABC",     // Optional
				JiraTicketNumber:   "SCOP-12345",
			},
			shouldBeValid: true,
		},
		{
			name: "invalid - missing country (required for international)",
			input: WarehouseToEngineerFormInput{
				LaptopID:           1,
				SoftwareEngineerID: 5,
				EngineerName:       "Jane Smith",
				EngineerEmail:      "jane@bairesdev.com",
				EngineerAddress:    "123 Main Street",
				EngineerCity:       "London",
				EngineerCountry:    "", // Missing - required
				EngineerState:      "",
				EngineerZip:        "",
				JiraTicketNumber:   "SCOP-12345",
			},
			shouldBeValid: false,
			errorContains: "country",
		},
		{
			name: "valid - international address without state (optional)",
			input: WarehouseToEngineerFormInput{
				LaptopID:           1,
				SoftwareEngineerID: 5,
				EngineerName:       "Jane Smith",
				EngineerEmail:      "jane@bairesdev.com",
				EngineerAddress:    "789 Oak Road",
				EngineerCity:       "Dublin",
				EngineerCountry:    "Ireland",
				EngineerState:      "", // Optional - not required
				EngineerZip:        "D02 XY45",
				JiraTicketNumber:   "SCOP-12345",
			},
			shouldBeValid: true,
		},
		{
			name: "valid - international address without postal code (optional)",
			input: WarehouseToEngineerFormInput{
				LaptopID:           1,
				SoftwareEngineerID: 5,
				EngineerName:       "Jane Smith",
				EngineerEmail:      "jane@bairesdev.com",
				EngineerAddress:    "321 Pine Street",
				EngineerCity:       "Tokyo",
				EngineerCountry:    "Japan",
				EngineerState:      "Tokyo",
				EngineerZip:        "", // Optional - not required
				JiraTicketNumber:   "SCOP-12345",
			},
			shouldBeValid: true,
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
