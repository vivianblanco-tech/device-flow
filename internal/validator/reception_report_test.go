package validator

import (
	"testing"
)

func TestValidateReceptionReport(t *testing.T) {
	tests := []struct {
		name    string
		input   ReceptionReportInput
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid report with all required fields",
			input: ReceptionReportInput{
				ShipmentID:      1,
				WarehouseUserID: 1,
				Notes:           "All items received in good condition",
				PhotoURLs:       []string{"https://example.com/photo1.jpg"},
			},
			wantErr: false,
		},
		{
			name: "valid report without photos",
			input: ReceptionReportInput{
				ShipmentID:      1,
				WarehouseUserID: 1,
				Notes:           "Items received",
				PhotoURLs:       []string{},
			},
			wantErr: false,
		},
		{
			name: "missing shipment ID",
			input: ReceptionReportInput{
				WarehouseUserID: 1,
				Notes:           "Items received",
			},
			wantErr: true,
			errMsg:  "shipment ID is required",
		},
		{
			name: "missing warehouse user ID",
			input: ReceptionReportInput{
				ShipmentID: 1,
				Notes:      "Items received",
			},
			wantErr: true,
			errMsg:  "warehouse user ID is required",
		},
		{
			name: "notes too long",
			input: ReceptionReportInput{
				ShipmentID:      1,
				WarehouseUserID: 1,
				Notes:           generateLongString(1001),
			},
			wantErr: true,
			errMsg:  "notes must not exceed 1000 characters",
		},
		{
			name: "too many photos",
			input: ReceptionReportInput{
				ShipmentID:      1,
				WarehouseUserID: 1,
				Notes:           "Items received",
				PhotoURLs:       []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11"},
			},
			wantErr: true,
			errMsg:  "cannot upload more than 10 photos",
		},
		{
			name: "invalid photo URL",
			input: ReceptionReportInput{
				ShipmentID:      1,
				WarehouseUserID: 1,
				Notes:           "Items received",
				PhotoURLs:       []string{"not-a-url"},
			},
			wantErr: true,
			errMsg:  "invalid photo URL format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateReceptionReport(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateReceptionReport() expected error but got nil")
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("ValidateReceptionReport() error = %v, want %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateReceptionReport() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestValidatePhotoURL(t *testing.T) {
	tests := []struct {
		name  string
		url   string
		valid bool
	}{
		{"valid http URL", "http://example.com/photo.jpg", true},
		{"valid https URL", "https://example.com/photo.jpg", true},
		{"valid URL with path", "https://storage.example.com/bucket/photo.jpg", true},
		{"invalid - no protocol", "example.com/photo.jpg", false},
		{"invalid - empty", "", false},
		{"invalid - spaces", "https://example.com/photo with spaces.jpg", false},
		{"invalid - not a URL", "just-text", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := isValidURL(tt.url)
			if valid != tt.valid {
				t.Errorf("isValidURL(%q) = %v, want %v", tt.url, valid, tt.valid)
			}
		})
	}
}

// Helper function to generate a long string for testing
func generateLongString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = 'a'
	}
	return string(b)
}

