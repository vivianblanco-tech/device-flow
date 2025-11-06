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
				PickupAddress:        "123 Main St, City, State 12345",
				PickupDate:           time.Now().Add(24 * time.Hour).Format("2006-01-02"),
				PickupTimeSlot:       "morning",
				NumberOfLaptops:      3,
				JiraTicketNumber:     "TEST-12345",
				SpecialInstructions:  "Call before arrival",
			},
			wantErr: false,
		},
		{
			name: "missing JIRA ticket number",
			input: PickupFormInput{
				ClientCompanyID:     1,
				ContactName:         "John Doe",
				ContactEmail:        "john@company.com",
				ContactPhone:        "+1-555-0123",
				PickupAddress:       "123 Main St",
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
				PickupDate:          time.Now().Add(24 * time.Hour).Format("2006-01-02"),
				PickupTimeSlot:      "morning",
				NumberOfLaptops:     -1,
				JiraTicketNumber:    "TEST-511",
			},
			wantErr: true,
			errMsg:  "number of laptops must be at least 1",
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

