package validator

import (
	"testing"
	"time"
)

func TestValidatePickupForm(t *testing.T) {
	tests := []struct {
		name    string
		input   PickupFormInput
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid form with all required fields",
			input: PickupFormInput{
				ClientCompanyID:      1,
				ContactName:          "John Doe",
				ContactEmail:         "john@company.com",
				ContactPhone:         "+1-555-0123",
				PickupAddress:        "123 Main St",
				PickupCity:           "New York",
				PickupState:          "NY",
				PickupZip:            "10001",
				PickupDate:           time.Now().Add(24 * time.Hour).Format("2006-01-02"),
				PickupTimeSlot:       "morning",
				NumberOfLaptops:      3,
				JiraTicketNumber:     "TEST-12345",
				SpecialInstructions:  "Call before arrival",
				NumberOfBoxes:        2,
				AssignmentType:       "single",
				IncludeAccessories:   false,
			},
			wantErr: false,
		},
		{
			name: "missing pickup city",
			input: PickupFormInput{
				ClientCompanyID:     1,
				ContactName:         "John Doe",
				ContactEmail:        "john@company.com",
				ContactPhone:        "+1-555-0123",
				PickupAddress:       "123 Main St",
				PickupState:         "NY",
				PickupZip:           "10001",
				PickupDate:          time.Now().Add(24 * time.Hour).Format("2006-01-02"),
				PickupTimeSlot:      "morning",
				NumberOfLaptops:     1,
				JiraTicketNumber:    "TEST-600",
			},
			wantErr: true,
			errMsg:  "pickup city is required",
		},
		{
			name: "missing pickup state",
			input: PickupFormInput{
				ClientCompanyID:     1,
				ContactName:         "John Doe",
				ContactEmail:        "john@company.com",
				ContactPhone:        "+1-555-0123",
				PickupAddress:       "123 Main St",
				PickupCity:          "New York",
				PickupZip:           "10001",
				PickupDate:          time.Now().Add(24 * time.Hour).Format("2006-01-02"),
				PickupTimeSlot:      "morning",
				NumberOfLaptops:     1,
				JiraTicketNumber:    "TEST-601",
			},
			wantErr: true,
			errMsg:  "pickup state is required",
		},
		{
			name: "invalid state code",
			input: PickupFormInput{
				ClientCompanyID:     1,
				ContactName:         "John Doe",
				ContactEmail:        "john@company.com",
				ContactPhone:        "+1-555-0123",
				PickupAddress:       "123 Main St",
				PickupCity:          "New York",
				PickupState:         "XX",
				PickupZip:           "10001",
				PickupDate:          time.Now().Add(24 * time.Hour).Format("2006-01-02"),
				PickupTimeSlot:      "morning",
				NumberOfLaptops:     1,
				JiraTicketNumber:    "TEST-602",
			},
			wantErr: true,
			errMsg:  "invalid US state code",
		},
		{
			name: "missing pickup zip",
			input: PickupFormInput{
				ClientCompanyID:     1,
				ContactName:         "John Doe",
				ContactEmail:        "john@company.com",
				ContactPhone:        "+1-555-0123",
				PickupAddress:       "123 Main St",
				PickupCity:          "New York",
				PickupState:         "NY",
				PickupDate:          time.Now().Add(24 * time.Hour).Format("2006-01-02"),
				PickupTimeSlot:      "morning",
				NumberOfLaptops:     1,
				JiraTicketNumber:    "TEST-603",
			},
			wantErr: true,
			errMsg:  "pickup ZIP code is required",
		},
		{
			name: "invalid zip code format - too short",
			input: PickupFormInput{
				ClientCompanyID:     1,
				ContactName:         "John Doe",
				ContactEmail:        "john@company.com",
				ContactPhone:        "+1-555-0123",
				PickupAddress:       "123 Main St",
				PickupCity:          "New York",
				PickupState:         "NY",
				PickupZip:           "1234",
				PickupDate:          time.Now().Add(24 * time.Hour).Format("2006-01-02"),
				PickupTimeSlot:      "morning",
				NumberOfLaptops:     1,
				JiraTicketNumber:    "TEST-604",
			},
			wantErr: true,
			errMsg:  "ZIP code must be 5 digits",
		},
		{
			name: "invalid zip code format - not digits",
			input: PickupFormInput{
				ClientCompanyID:     1,
				ContactName:         "John Doe",
				ContactEmail:        "john@company.com",
				ContactPhone:        "+1-555-0123",
				PickupAddress:       "123 Main St",
				PickupCity:          "New York",
				PickupState:         "NY",
				PickupZip:           "ABCDE",
				PickupDate:          time.Now().Add(24 * time.Hour).Format("2006-01-02"),
				PickupTimeSlot:      "morning",
				NumberOfLaptops:     1,
				JiraTicketNumber:    "TEST-605",
			},
			wantErr: true,
			errMsg:  "ZIP code must be 5 digits",
		},
		{
			name: "missing JIRA ticket number",
			input: PickupFormInput{
				ClientCompanyID:     1,
				ContactName:         "John Doe",
				ContactEmail:        "john@company.com",
				ContactPhone:        "+1-555-0123",
				PickupAddress:       "123 Main St",
				PickupCity:          "New York",
				PickupState:         "NY",
				PickupZip:           "10001",
				PickupDate:          time.Now().Add(24 * time.Hour).Format("2006-01-02"),
				PickupTimeSlot:      "morning",
				NumberOfLaptops:     1,
				JiraTicketNumber:    "",
			},
			wantErr: true,
			errMsg:  "JIRA ticket number is required",
		},
		{
			name: "invalid JIRA ticket format - no dash",
			input: PickupFormInput{
				ClientCompanyID:     1,
				ContactName:         "John Doe",
				ContactEmail:        "john@company.com",
				ContactPhone:        "+1-555-0123",
				PickupAddress:       "123 Main St",
				PickupCity:          "New York",
				PickupState:         "NY",
				PickupZip:           "10001",
				PickupDate:          time.Now().Add(24 * time.Hour).Format("2006-01-02"),
				PickupTimeSlot:      "morning",
				NumberOfLaptops:     1,
				JiraTicketNumber:    "TEST12345",
			},
			wantErr: true,
			errMsg:  "JIRA ticket number must be in format PROJECT-NUMBER (e.g., SCOP-67702)",
		},
		{
			name: "invalid JIRA ticket format - lowercase",
			input: PickupFormInput{
				ClientCompanyID:     1,
				ContactName:         "John Doe",
				ContactEmail:        "john@company.com",
				ContactPhone:        "+1-555-0123",
				PickupAddress:       "123 Main St",
				PickupCity:          "New York",
				PickupState:         "NY",
				PickupZip:           "10001",
				PickupDate:          time.Now().Add(24 * time.Hour).Format("2006-01-02"),
				PickupTimeSlot:      "morning",
				NumberOfLaptops:     1,
				JiraTicketNumber:    "test-12345",
			},
			wantErr: true,
			errMsg:  "JIRA ticket number must be in format PROJECT-NUMBER (e.g., SCOP-67702)",
		},
		{
			name: "invalid JIRA ticket format - no number",
			input: PickupFormInput{
				ClientCompanyID:     1,
				ContactName:         "John Doe",
				ContactEmail:        "john@company.com",
				ContactPhone:        "+1-555-0123",
				PickupAddress:       "123 Main St",
				PickupCity:          "New York",
				PickupState:         "NY",
				PickupZip:           "10001",
				PickupDate:          time.Now().Add(24 * time.Hour).Format("2006-01-02"),
				PickupTimeSlot:      "morning",
				NumberOfLaptops:     1,
				JiraTicketNumber:    "TEST-",
			},
			wantErr: true,
			errMsg:  "JIRA ticket number must be in format PROJECT-NUMBER (e.g., SCOP-67702)",
		},
		{
			name: "missing client company ID",
			input: PickupFormInput{
				ContactName:         "John Doe",
				ContactEmail:        "john@company.com",
				ContactPhone:        "+1-555-0123",
				PickupAddress:       "123 Main St",
				PickupCity:          "New York",
				PickupState:         "NY",
				PickupZip:           "10001",
				PickupDate:          time.Now().Add(24 * time.Hour).Format("2006-01-02"),
				PickupTimeSlot:      "morning",
				NumberOfLaptops:     1,
				JiraTicketNumber:    "TEST-500",
			},
			wantErr: true,
			errMsg:  "client company ID is required",
		},
		{
			name: "missing contact name",
			input: PickupFormInput{
				ClientCompanyID:     1,
				ContactEmail:        "john@company.com",
				ContactPhone:        "+1-555-0123",
				PickupAddress:       "123 Main St",
				PickupCity:          "New York",
				PickupState:         "NY",
				PickupZip:           "10001",
				PickupDate:          time.Now().Add(24 * time.Hour).Format("2006-01-02"),
				PickupTimeSlot:      "morning",
				NumberOfLaptops:     1,
				JiraTicketNumber:    "TEST-501",
			},
			wantErr: true,
			errMsg:  "contact name is required",
		},
		{
			name: "invalid email format",
			input: PickupFormInput{
				ClientCompanyID:     1,
				ContactName:         "John Doe",
				ContactEmail:        "invalid-email",
				ContactPhone:        "+1-555-0123",
				PickupAddress:       "123 Main St",
				PickupCity:          "New York",
				PickupState:         "NY",
				PickupZip:           "10001",
				PickupDate:          time.Now().Add(24 * time.Hour).Format("2006-01-02"),
				PickupTimeSlot:      "morning",
				NumberOfLaptops:     1,
				JiraTicketNumber:    "TEST-502",
			},
			wantErr: true,
			errMsg:  "invalid email format",
		},
		{
			name: "missing contact phone",
			input: PickupFormInput{
				ClientCompanyID:     1,
				ContactName:         "John Doe",
				ContactEmail:        "john@company.com",
				PickupAddress:       "123 Main St",
				PickupCity:          "New York",
				PickupState:         "NY",
				PickupZip:           "10001",
				PickupDate:          time.Now().Add(24 * time.Hour).Format("2006-01-02"),
				PickupTimeSlot:      "morning",
				NumberOfLaptops:     1,
				JiraTicketNumber:    "TEST-503",
			},
			wantErr: true,
			errMsg:  "contact phone is required",
		},
		{
			name: "missing pickup address",
			input: PickupFormInput{
				ClientCompanyID:     1,
				ContactName:         "John Doe",
				ContactEmail:        "john@company.com",
				ContactPhone:        "+1-555-0123",
				PickupDate:          time.Now().Add(24 * time.Hour).Format("2006-01-02"),
				PickupTimeSlot:      "morning",
				NumberOfLaptops:     1,
				JiraTicketNumber:    "TEST-504",
			},
			wantErr: true,
			errMsg:  "pickup address is required",
		},
		{
			name: "missing pickup date",
			input: PickupFormInput{
				ClientCompanyID:     1,
				ContactName:         "John Doe",
				ContactEmail:        "john@company.com",
				ContactPhone:        "+1-555-0123",
				PickupAddress:       "123 Main St",
				PickupCity:          "New York",
				PickupState:         "NY",
				PickupZip:           "10001",
				PickupTimeSlot:      "morning",
				NumberOfLaptops:     1,
				JiraTicketNumber:    "TEST-505",
			},
			wantErr: true,
			errMsg:  "pickup date is required",
		},
		{
			name: "invalid date format",
			input: PickupFormInput{
				ClientCompanyID:     1,
				ContactName:         "John Doe",
				ContactEmail:        "john@company.com",
				ContactPhone:        "+1-555-0123",
				PickupAddress:       "123 Main St",
				PickupCity:          "New York",
				PickupState:         "NY",
				PickupZip:           "10001",
				PickupDate:          "invalid-date",
				PickupTimeSlot:      "morning",
				NumberOfLaptops:     1,
				JiraTicketNumber:    "TEST-506",
			},
			wantErr: true,
			errMsg:  "invalid date format",
		},
		{
			name: "date in the past",
			input: PickupFormInput{
				ClientCompanyID:     1,
				ContactName:         "John Doe",
				ContactEmail:        "john@company.com",
				ContactPhone:        "+1-555-0123",
				PickupAddress:       "123 Main St",
				PickupCity:          "New York",
				PickupState:         "NY",
				PickupZip:           "10001",
				PickupDate:          "2020-01-01",
				PickupTimeSlot:      "morning",
				NumberOfLaptops:     1,
				JiraTicketNumber:    "TEST-507",
			},
			wantErr: true,
			errMsg:  "pickup date must be in the future",
		},
		{
			name: "missing time slot",
			input: PickupFormInput{
				ClientCompanyID:     1,
				ContactName:         "John Doe",
				ContactEmail:        "john@company.com",
				ContactPhone:        "+1-555-0123",
				PickupAddress:       "123 Main St",
				PickupCity:          "New York",
				PickupState:         "NY",
				PickupZip:           "10001",
				PickupDate:          time.Now().Add(24 * time.Hour).Format("2006-01-02"),
				NumberOfLaptops:     1,
				JiraTicketNumber:    "TEST-508",
			},
			wantErr: true,
			errMsg:  "pickup time slot is required",
		},
		{
			name: "invalid time slot",
			input: PickupFormInput{
				ClientCompanyID:     1,
				ContactName:         "John Doe",
				ContactEmail:        "john@company.com",
				ContactPhone:        "+1-555-0123",
				PickupAddress:       "123 Main St",
				PickupCity:          "New York",
				PickupState:         "NY",
				PickupZip:           "10001",
				PickupDate:          time.Now().Add(24 * time.Hour).Format("2006-01-02"),
				PickupTimeSlot:      "invalid-slot",
				NumberOfLaptops:     1,
				JiraTicketNumber:    "TEST-509",
			},
			wantErr: true,
			errMsg:  "invalid time slot",
		},
		{
			name: "number of laptops is zero",
			input: PickupFormInput{
				ClientCompanyID:     1,
				ContactName:         "John Doe",
				ContactEmail:        "john@company.com",
				ContactPhone:        "+1-555-0123",
				PickupAddress:       "123 Main St",
				PickupCity:          "New York",
				PickupState:         "NY",
				PickupZip:           "10001",
				PickupDate:          time.Now().Add(24 * time.Hour).Format("2006-01-02"),
				PickupTimeSlot:      "morning",
				NumberOfLaptops:     0,
				JiraTicketNumber:    "TEST-510",
			},
			wantErr: true,
			errMsg:  "number of laptops must be at least 1",
		},
		{
			name: "number of laptops is negative",
			input: PickupFormInput{
				ClientCompanyID:     1,
				ContactName:         "John Doe",
				ContactEmail:        "john@company.com",
				ContactPhone:        "+1-555-0123",
				PickupAddress:       "123 Main St",
				PickupCity:          "New York",
				PickupState:         "NY",
				PickupZip:           "10001",
				PickupDate:          time.Now().Add(24 * time.Hour).Format("2006-01-02"),
				PickupTimeSlot:      "morning",
				NumberOfLaptops:     -1,
				JiraTicketNumber:    "TEST-511",
			},
			wantErr: true,
			errMsg:  "number of laptops must be at least 1",
		},
		// NEW FIELD TESTS: Number of Boxes
		{
			name: "missing number of boxes",
			input: PickupFormInput{
				ClientCompanyID:     1,
				ContactName:         "John Doe",
				ContactEmail:        "john@company.com",
				ContactPhone:        "+1-555-0123",
				PickupAddress:       "123 Main St",
				PickupCity:          "New York",
				PickupState:         "NY",
				PickupZip:           "10001",
				PickupDate:          time.Now().Add(24 * time.Hour).Format("2006-01-02"),
				PickupTimeSlot:      "morning",
				NumberOfLaptops:     1,
				JiraTicketNumber:    "TEST-600",
				NumberOfBoxes:       0,
			},
			wantErr: true,
			errMsg:  "number of boxes must be at least 1",
		},
		{
			name: "negative number of boxes",
			input: PickupFormInput{
				ClientCompanyID:     1,
				ContactName:         "John Doe",
				ContactEmail:        "john@company.com",
				ContactPhone:        "+1-555-0123",
				PickupAddress:       "123 Main St",
				PickupCity:          "New York",
				PickupState:         "NY",
				PickupZip:           "10001",
				PickupDate:          time.Now().Add(24 * time.Hour).Format("2006-01-02"),
				PickupTimeSlot:      "morning",
				NumberOfLaptops:     1,
				JiraTicketNumber:    "TEST-601",
				NumberOfBoxes:       -1,
			},
			wantErr: true,
			errMsg:  "number of boxes must be at least 1",
		},
		// NEW FIELD TESTS: Assignment Type
		{
			name: "missing assignment type",
			input: PickupFormInput{
				ClientCompanyID:     1,
				ContactName:         "John Doe",
				ContactEmail:        "john@company.com",
				ContactPhone:        "+1-555-0123",
				PickupAddress:       "123 Main St",
				PickupCity:          "New York",
				PickupState:         "NY",
				PickupZip:           "10001",
				PickupDate:          time.Now().Add(24 * time.Hour).Format("2006-01-02"),
				PickupTimeSlot:      "morning",
				NumberOfLaptops:     1,
				JiraTicketNumber:    "TEST-602",
				NumberOfBoxes:       1,
				AssignmentType:      "",
			},
			wantErr: true,
			errMsg:  "assignment type is required",
		},
		{
			name: "invalid assignment type",
			input: PickupFormInput{
				ClientCompanyID:     1,
				ContactName:         "John Doe",
				ContactEmail:        "john@company.com",
				ContactPhone:        "+1-555-0123",
				PickupAddress:       "123 Main St",
				PickupCity:          "New York",
				PickupState:         "NY",
				PickupZip:           "10001",
				PickupDate:          time.Now().Add(24 * time.Hour).Format("2006-01-02"),
				PickupTimeSlot:      "morning",
				NumberOfLaptops:     1,
				JiraTicketNumber:    "TEST-603",
				NumberOfBoxes:       1,
				AssignmentType:      "invalid",
			},
			wantErr: true,
			errMsg:  "assignment type must be 'single' or 'bulk'",
		},
		// NEW FIELD TESTS: Bulk Dimensions and Weight
		{
			name: "bulk assignment missing dimensions",
			input: PickupFormInput{
				ClientCompanyID:     1,
				ContactName:         "John Doe",
				ContactEmail:        "john@company.com",
				ContactPhone:        "+1-555-0123",
				PickupAddress:       "123 Main St",
				PickupCity:          "New York",
				PickupState:         "NY",
				PickupZip:           "10001",
				PickupDate:          time.Now().Add(24 * time.Hour).Format("2006-01-02"),
				PickupTimeSlot:      "morning",
				NumberOfLaptops:     1,
				JiraTicketNumber:    "TEST-604",
				NumberOfBoxes:       2,
				AssignmentType:      "bulk",
			},
			wantErr: true,
			errMsg:  "bulk length is required for bulk shipments",
		},
		{
			name: "bulk assignment with zero length",
			input: PickupFormInput{
				ClientCompanyID:     1,
				ContactName:         "John Doe",
				ContactEmail:        "john@company.com",
				ContactPhone:        "+1-555-0123",
				PickupAddress:       "123 Main St",
				PickupCity:          "New York",
				PickupState:         "NY",
				PickupZip:           "10001",
				PickupDate:          time.Now().Add(24 * time.Hour).Format("2006-01-02"),
				PickupTimeSlot:      "morning",
				NumberOfLaptops:     1,
				JiraTicketNumber:    "TEST-605",
				NumberOfBoxes:       2,
				AssignmentType:      "bulk",
				BulkLength:          0,
			},
			wantErr: true,
			errMsg:  "bulk length is required for bulk shipments",
		},
		{
			name: "bulk assignment missing width",
			input: PickupFormInput{
				ClientCompanyID:     1,
				ContactName:         "John Doe",
				ContactEmail:        "john@company.com",
				ContactPhone:        "+1-555-0123",
				PickupAddress:       "123 Main St",
				PickupCity:          "New York",
				PickupState:         "NY",
				PickupZip:           "10001",
				PickupDate:          time.Now().Add(24 * time.Hour).Format("2006-01-02"),
				PickupTimeSlot:      "morning",
				NumberOfLaptops:     1,
				JiraTicketNumber:    "TEST-606",
				NumberOfBoxes:       2,
				AssignmentType:      "bulk",
				BulkLength:          20.5,
			},
			wantErr: true,
			errMsg:  "bulk width is required for bulk shipments",
		},
		{
			name: "bulk assignment missing height",
			input: PickupFormInput{
				ClientCompanyID:     1,
				ContactName:         "John Doe",
				ContactEmail:        "john@company.com",
				ContactPhone:        "+1-555-0123",
				PickupAddress:       "123 Main St",
				PickupCity:          "New York",
				PickupState:         "NY",
				PickupZip:           "10001",
				PickupDate:          time.Now().Add(24 * time.Hour).Format("2006-01-02"),
				PickupTimeSlot:      "morning",
				NumberOfLaptops:     1,
				JiraTicketNumber:    "TEST-607",
				NumberOfBoxes:       2,
				AssignmentType:      "bulk",
				BulkLength:          20.5,
				BulkWidth:           15.0,
			},
			wantErr: true,
			errMsg:  "bulk height is required for bulk shipments",
		},
		{
			name: "bulk assignment missing weight",
			input: PickupFormInput{
				ClientCompanyID:     1,
				ContactName:         "John Doe",
				ContactEmail:        "john@company.com",
				ContactPhone:        "+1-555-0123",
				PickupAddress:       "123 Main St",
				PickupCity:          "New York",
				PickupState:         "NY",
				PickupZip:           "10001",
				PickupDate:          time.Now().Add(24 * time.Hour).Format("2006-01-02"),
				PickupTimeSlot:      "morning",
				NumberOfLaptops:     1,
				JiraTicketNumber:    "TEST-608",
				NumberOfBoxes:       2,
				AssignmentType:      "bulk",
				BulkLength:          20.5,
				BulkWidth:           15.0,
				BulkHeight:          10.0,
			},
			wantErr: true,
			errMsg:  "bulk weight is required for bulk shipments",
		},
		{
			name: "valid bulk assignment with all dimensions",
			input: PickupFormInput{
				ClientCompanyID:     1,
				ContactName:         "John Doe",
				ContactEmail:        "john@company.com",
				ContactPhone:        "+1-555-0123",
				PickupAddress:       "123 Main St",
				PickupCity:          "New York",
				PickupState:         "NY",
				PickupZip:           "10001",
				PickupDate:          time.Now().Add(24 * time.Hour).Format("2006-01-02"),
				PickupTimeSlot:      "morning",
				NumberOfLaptops:     5,
				JiraTicketNumber:    "TEST-609",
				NumberOfBoxes:       2,
				AssignmentType:      "bulk",
				BulkLength:          20.5,
				BulkWidth:           15.0,
				BulkHeight:          10.0,
				BulkWeight:          25.5,
			},
			wantErr: false,
		},
		{
			name: "valid single assignment without dimensions",
			input: PickupFormInput{
				ClientCompanyID:     1,
				ContactName:         "John Doe",
				ContactEmail:        "john@company.com",
				ContactPhone:        "+1-555-0123",
				PickupAddress:       "123 Main St",
				PickupCity:          "New York",
				PickupState:         "NY",
				PickupZip:           "10001",
				PickupDate:          time.Now().Add(24 * time.Hour).Format("2006-01-02"),
				PickupTimeSlot:      "morning",
				NumberOfLaptops:     1,
				JiraTicketNumber:    "TEST-610",
				NumberOfBoxes:       1,
				AssignmentType:      "single",
			},
			wantErr: false,
		},
		// NEW FIELD TESTS: Accessories
		{
			name: "include accessories without description",
			input: PickupFormInput{
				ClientCompanyID:        1,
				ContactName:            "John Doe",
				ContactEmail:           "john@company.com",
				ContactPhone:           "+1-555-0123",
				PickupAddress:          "123 Main St",
				PickupCity:             "New York",
				PickupState:            "NY",
				PickupZip:              "10001",
				PickupDate:             time.Now().Add(24 * time.Hour).Format("2006-01-02"),
				PickupTimeSlot:         "morning",
				NumberOfLaptops:        1,
				JiraTicketNumber:       "TEST-611",
				NumberOfBoxes:          1,
				AssignmentType:         "single",
				IncludeAccessories:     true,
				AccessoriesDescription: "",
			},
			wantErr: true,
			errMsg:  "accessories description is required when including accessories",
		},
		{
			name: "include accessories with description",
			input: PickupFormInput{
				ClientCompanyID:        1,
				ContactName:            "John Doe",
				ContactEmail:           "john@company.com",
				ContactPhone:           "+1-555-0123",
				PickupAddress:          "123 Main St",
				PickupCity:             "New York",
				PickupState:            "NY",
				PickupZip:              "10001",
				PickupDate:             time.Now().Add(24 * time.Hour).Format("2006-01-02"),
				PickupTimeSlot:         "morning",
				NumberOfLaptops:        1,
				JiraTicketNumber:       "TEST-612",
				NumberOfBoxes:          1,
				AssignmentType:         "single",
				IncludeAccessories:     true,
				AccessoriesDescription: "2x YubiKeys, 3x USB-C cables, 1x dock",
			},
			wantErr: false,
		},
		{
			name: "no accessories with description ignored",
			input: PickupFormInput{
				ClientCompanyID:        1,
				ContactName:            "John Doe",
				ContactEmail:           "john@company.com",
				ContactPhone:           "+1-555-0123",
				PickupAddress:          "123 Main St",
				PickupCity:             "New York",
				PickupState:            "NY",
				PickupZip:              "10001",
				PickupDate:             time.Now().Add(24 * time.Hour).Format("2006-01-02"),
				PickupTimeSlot:         "morning",
				NumberOfLaptops:        1,
				JiraTicketNumber:       "TEST-613",
				NumberOfBoxes:          1,
				AssignmentType:         "single",
				IncludeAccessories:     false,
				AccessoriesDescription: "This should be ignored",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePickupForm(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidatePickupForm() expected error but got nil")
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("ValidatePickupForm() error = %v, want %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidatePickupForm() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		valid bool
	}{
		{"valid email", "user@example.com", true},
		{"valid email with subdomain", "user@mail.example.com", true},
		{"valid email with plus", "user+tag@example.com", true},
		{"invalid - no @", "userexample.com", false},
		{"invalid - no domain", "user@", false},
		{"invalid - no user", "@example.com", false},
		{"invalid - spaces", "user @example.com", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := isValidEmail(tt.email)
			if valid != tt.valid {
				t.Errorf("isValidEmail(%q) = %v, want %v", tt.email, valid, tt.valid)
			}
		})
	}
}

func TestValidateTimeSlot(t *testing.T) {
	tests := []struct {
		name     string
		timeSlot string
		valid    bool
	}{
		{"morning slot", "morning", true},
		{"afternoon slot", "afternoon", true},
		{"evening slot", "evening", true},
		{"invalid slot", "night", false},
		{"empty slot", "", false},
		{"uppercase", "MORNING", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := isValidTimeSlot(tt.timeSlot)
			if valid != tt.valid {
				t.Errorf("isValidTimeSlot(%q) = %v, want %v", tt.timeSlot, valid, tt.valid)
			}
		})
	}
}

func TestValidateUSState(t *testing.T) {
	tests := []struct {
		name  string
		state string
		valid bool
	}{
		{"valid state - NY", "NY", true},
		{"valid state - CA", "CA", true},
		{"valid state - TX", "TX", true},
		{"valid state - FL", "FL", true},
		{"invalid state - XX", "XX", false},
		{"invalid state - lowercase", "ny", false},
		{"invalid state - full name", "New York", false},
		{"empty state", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := isValidUSState(tt.state)
			if valid != tt.valid {
				t.Errorf("isValidUSState(%q) = %v, want %v", tt.state, valid, tt.valid)
			}
		})
	}
}

func TestValidateZipCode(t *testing.T) {
	tests := []struct {
		name  string
		zip   string
		valid bool
	}{
		{"valid 5-digit zip", "12345", true},
		{"valid 5-digit zip - zeros", "00000", true},
		{"invalid - 4 digits", "1234", false},
		{"invalid - 6 digits", "123456", false},
		{"invalid - letters", "ABCDE", false},
		{"invalid - mixed", "12A45", false},
		{"empty zip", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := isValidZipCode(tt.zip)
			if valid != tt.valid {
				t.Errorf("isValidZipCode(%q) = %v, want %v", tt.zip, valid, tt.valid)
			}
		})
	}
}

