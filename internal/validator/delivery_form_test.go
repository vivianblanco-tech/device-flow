package validator

import (
	"testing"
)

func TestValidateDeliveryForm(t *testing.T) {
	tests := []struct {
		name    string
		input   DeliveryFormInput
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid form with all required fields",
			input: DeliveryFormInput{
				ShipmentID:  1,
				EngineerID:  1,
				Notes:       "Device delivered in good condition",
				PhotoURLs:   []string{"https://example.com/photo1.jpg"},
			},
			wantErr: false,
		},
		{
			name: "valid form without photos",
			input: DeliveryFormInput{
				ShipmentID: 1,
				EngineerID: 1,
				Notes:      "Delivered successfully",
				PhotoURLs:  []string{},
			},
			wantErr: false,
		},
		{
			name: "missing shipment ID",
			input: DeliveryFormInput{
				EngineerID: 1,
				Notes:      "Delivered",
			},
			wantErr: true,
			errMsg:  "shipment ID is required",
		},
		{
			name: "missing engineer ID",
			input: DeliveryFormInput{
				ShipmentID: 1,
				Notes:      "Delivered",
			},
			wantErr: true,
			errMsg:  "engineer ID is required",
		},
		{
			name: "notes too long",
			input: DeliveryFormInput{
				ShipmentID: 1,
				EngineerID: 1,
				Notes:      generateLongString(1001),
			},
			wantErr: true,
			errMsg:  "notes must not exceed 1000 characters",
		},
		{
			name: "too many photos",
			input: DeliveryFormInput{
				ShipmentID: 1,
				EngineerID: 1,
				Notes:      "Delivered",
				PhotoURLs:  []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11"},
			},
			wantErr: true,
			errMsg:  "cannot upload more than 10 photos",
		},
		{
			name: "invalid photo URL",
			input: DeliveryFormInput{
				ShipmentID: 1,
				EngineerID: 1,
				Notes:      "Delivered",
				PhotoURLs:  []string{"not-a-valid-url"},
			},
			wantErr: true,
			errMsg:  "invalid photo URL format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDeliveryForm(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateDeliveryForm() expected error but got nil")
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("ValidateDeliveryForm() error = %v, want %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateDeliveryForm() unexpected error = %v", err)
				}
			}
		})
	}
}

