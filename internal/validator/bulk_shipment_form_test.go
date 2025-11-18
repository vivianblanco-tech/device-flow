package validator

import (
	"strings"
	"testing"
	"time"
)

func TestValidateBulkToWarehouseForm(t *testing.T) {
	tests := []struct {
		name          string
		input         BulkToWarehouseFormInput
		shouldBeValid bool
		errorContains string
	}{
		{
			name: "valid bulk form",
			input: BulkToWarehouseFormInput{
				ClientCompanyID:  1,
				ContactName:      "John Doe",
				ContactEmail:     "john@company.com",
				ContactPhone:     "+1-555-0123",
				PickupAddress:    "123 Main St",
				PickupCity:       "New York",
				PickupState:      "NY",
				PickupZip:        "10001",
				PickupDate:       time.Now().AddDate(0, 0, 2).Format("2006-01-02"),
				PickupTimeSlot:   "morning",
				JiraTicketNumber: "SCOP-12345",
				NumberOfLaptops:  5,
				BulkLength:       30.0,
				BulkWidth:        20.0,
				BulkHeight:       15.0,
				BulkWeight:       50.0,
			},
			shouldBeValid: true,
		},
		{
			name: "invalid - laptop count too low",
			input: BulkToWarehouseFormInput{
				ClientCompanyID:  1,
				ContactName:      "John Doe",
				ContactEmail:     "john@company.com",
				ContactPhone:     "+1-555-0123",
				PickupAddress:    "123 Main St",
				PickupCity:       "New York",
				PickupState:      "NY",
				PickupZip:        "10001",
				PickupDate:       time.Now().AddDate(0, 0, 2).Format("2006-01-02"),
				PickupTimeSlot:   "morning",
				JiraTicketNumber: "SCOP-12345",
				NumberOfLaptops:  1, // Too low for bulk
				BulkLength:       30.0,
				BulkWidth:        20.0,
				BulkHeight:       15.0,
				BulkWeight:       50.0,
			},
			shouldBeValid: false,
			errorContains: "at least 2 laptops",
		},
		{
			name: "invalid - missing bulk length",
			input: BulkToWarehouseFormInput{
				ClientCompanyID:  1,
				ContactName:      "John Doe",
				ContactEmail:     "john@company.com",
				ContactPhone:     "+1-555-0123",
				PickupAddress:    "123 Main St",
				PickupCity:       "New York",
				PickupState:      "NY",
				PickupZip:        "10001",
				PickupDate:       time.Now().AddDate(0, 0, 2).Format("2006-01-02"),
				PickupTimeSlot:   "morning",
				JiraTicketNumber: "SCOP-12345",
				NumberOfLaptops:  5,
				BulkLength:       0, // Missing
				BulkWidth:        20.0,
				BulkHeight:       15.0,
				BulkWeight:       50.0,
			},
			shouldBeValid: false,
			errorContains: "bulk dimensions",
		},
		{
			name: "invalid - missing bulk width",
			input: BulkToWarehouseFormInput{
				ClientCompanyID:  1,
				ContactName:      "John Doe",
				ContactEmail:     "john@company.com",
				ContactPhone:     "+1-555-0123",
				PickupAddress:    "123 Main St",
				PickupCity:       "New York",
				PickupState:      "NY",
				PickupZip:        "10001",
				PickupDate:       time.Now().AddDate(0, 0, 2).Format("2006-01-02"),
				PickupTimeSlot:   "morning",
				JiraTicketNumber: "SCOP-12345",
				NumberOfLaptops:  5,
				BulkLength:       30.0,
				BulkWidth:        0, // Missing
				BulkHeight:       15.0,
				BulkWeight:       50.0,
			},
			shouldBeValid: false,
			errorContains: "bulk dimensions",
		},
		{
			name: "invalid - missing bulk height",
			input: BulkToWarehouseFormInput{
				ClientCompanyID:  1,
				ContactName:      "John Doe",
				ContactEmail:     "john@company.com",
				ContactPhone:     "+1-555-0123",
				PickupAddress:    "123 Main St",
				PickupCity:       "New York",
				PickupState:      "NY",
				PickupZip:        "10001",
				PickupDate:       time.Now().AddDate(0, 0, 2).Format("2006-01-02"),
				PickupTimeSlot:   "morning",
				JiraTicketNumber: "SCOP-12345",
				NumberOfLaptops:  5,
				BulkLength:       30.0,
				BulkWidth:        20.0,
				BulkHeight:       0, // Missing
				BulkWeight:       50.0,
			},
			shouldBeValid: false,
			errorContains: "bulk dimensions",
		},
		{
			name: "invalid - missing bulk weight",
			input: BulkToWarehouseFormInput{
				ClientCompanyID:  1,
				ContactName:      "John Doe",
				ContactEmail:     "john@company.com",
				ContactPhone:     "+1-555-0123",
				PickupAddress:    "123 Main St",
				PickupCity:       "New York",
				PickupState:      "NY",
				PickupZip:        "10001",
				PickupDate:       time.Now().AddDate(0, 0, 2).Format("2006-01-02"),
				PickupTimeSlot:   "morning",
				JiraTicketNumber: "SCOP-12345",
				NumberOfLaptops:  5,
				BulkLength:       30.0,
				BulkWidth:        20.0,
				BulkHeight:       15.0,
				BulkWeight:       0, // Missing
			},
			shouldBeValid: false,
			errorContains: "bulk dimensions",
		},
		{
			name: "invalid - negative bulk dimensions",
			input: BulkToWarehouseFormInput{
				ClientCompanyID:  1,
				ContactName:      "John Doe",
				ContactEmail:     "john@company.com",
				ContactPhone:     "+1-555-0123",
				PickupAddress:    "123 Main St",
				PickupCity:       "New York",
				PickupState:      "NY",
				PickupZip:        "10001",
				PickupDate:       time.Now().AddDate(0, 0, 2).Format("2006-01-02"),
				PickupTimeSlot:   "morning",
				JiraTicketNumber: "SCOP-12345",
				NumberOfLaptops:  5,
				BulkLength:       -30.0, // Negative
				BulkWidth:        20.0,
				BulkHeight:       15.0,
				BulkWeight:       50.0,
			},
			shouldBeValid: false,
			errorContains: "positive",
		},
		{
			name: "valid with accessories",
			input: BulkToWarehouseFormInput{
				ClientCompanyID:        1,
				ContactName:            "John Doe",
				ContactEmail:           "john@company.com",
				ContactPhone:           "+1-555-0123",
				PickupAddress:          "123 Main St",
				PickupCity:             "New York",
				PickupState:            "NY",
				PickupZip:              "10001",
				PickupDate:             time.Now().AddDate(0, 0, 2).Format("2006-01-02"),
				PickupTimeSlot:         "morning",
				JiraTicketNumber:       "SCOP-12345",
				NumberOfLaptops:        10,
				BulkLength:             40.0,
				BulkWidth:              30.0,
				BulkHeight:             20.0,
				BulkWeight:             100.0,
				IncludeAccessories:     true,
				AccessoriesDescription: "Chargers and mice for all laptops",
			},
			shouldBeValid: true,
		},
		{
			name: "invalid - accessories description missing when including accessories",
			input: BulkToWarehouseFormInput{
				ClientCompanyID:        1,
				ContactName:            "John Doe",
				ContactEmail:           "john@company.com",
				ContactPhone:           "+1-555-0123",
				PickupAddress:          "123 Main St",
				PickupCity:             "New York",
				PickupState:            "NY",
				PickupZip:              "10001",
				PickupDate:             time.Now().AddDate(0, 0, 2).Format("2006-01-02"),
				PickupTimeSlot:         "morning",
				JiraTicketNumber:       "SCOP-12345",
				NumberOfLaptops:        5,
				BulkLength:             30.0,
				BulkWidth:              20.0,
				BulkHeight:             15.0,
				BulkWeight:             50.0,
				IncludeAccessories:     true,
				AccessoriesDescription: "", // Missing
			},
			shouldBeValid: false,
			errorContains: "accessories description is required",
		},
		{
			name: "invalid - missing contact email",
			input: BulkToWarehouseFormInput{
				ClientCompanyID:  1,
				ContactName:      "John Doe",
				ContactEmail:     "", // Missing
				ContactPhone:     "+1-555-0123",
				PickupAddress:    "123 Main St",
				PickupCity:       "New York",
				PickupState:      "NY",
				PickupZip:        "10001",
				PickupDate:       time.Now().AddDate(0, 0, 2).Format("2006-01-02"),
				PickupTimeSlot:   "morning",
				JiraTicketNumber: "SCOP-12345",
				NumberOfLaptops:  5,
				BulkLength:       30.0,
				BulkWidth:        20.0,
				BulkHeight:       15.0,
				BulkWeight:       50.0,
			},
			shouldBeValid: false,
			errorContains: "email",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateBulkToWarehouseForm(tt.input)
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
