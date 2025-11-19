package models

import (
	"testing"
	"time"
)

func TestSoftwareEngineer_Validate(t *testing.T) {
	tests := []struct {
		name     string
		engineer SoftwareEngineer
		wantErr  bool
		errMsg   string
	}{
		{
			name: "valid engineer with all fields",
			engineer: SoftwareEngineer{
				Name:                  "John Doe",
				Email:                 "john.doe@bairesdev.com",
				Address:               "123 Main St, City, State 12345",
				Phone:                 "+1-555-0100",
				AddressConfirmed:      true,
				AddressConfirmationAt: timePtr(time.Now()),
			},
			wantErr: false,
		},
		{
			name: "valid engineer with minimal fields",
			engineer: SoftwareEngineer{
				Name:  "Jane Smith",
				Email: "jane@bairesdev.com",
			},
			wantErr: false,
		},
		{
			name: "invalid - missing name",
			engineer: SoftwareEngineer{
				Email: "test@bairesdev.com",
			},
			wantErr: true,
			errMsg:  "engineer name is required",
		},
		{
			name: "invalid - empty name",
			engineer: SoftwareEngineer{
				Name:  "",
				Email: "test@bairesdev.com",
			},
			wantErr: true,
			errMsg:  "engineer name is required",
		},
		{
			name: "invalid - missing email",
			engineer: SoftwareEngineer{
				Name: "John Doe",
			},
			wantErr: true,
			errMsg:  "engineer email is required",
		},
		{
			name: "invalid - invalid email format",
			engineer: SoftwareEngineer{
				Name:  "John Doe",
				Email: "invalid-email",
			},
			wantErr: true,
			errMsg:  "invalid email format",
		},
		{
			name: "valid - name exactly 2 characters",
			engineer: SoftwareEngineer{
				Name:  "Jo",
				Email: "jo@bairesdev.com",
			},
			wantErr: false,
		},
		{
			name: "valid - engineer with employee number",
			engineer: SoftwareEngineer{
				Name:           "John Doe",
				Email:          "john@bairesdev.com",
				EmployeeNumber: "EMP-12345",
			},
			wantErr: false,
		},
		{
			name: "valid - engineer without employee number (optional field)",
			engineer: SoftwareEngineer{
				Name:  "Jane Smith",
				Email: "jane@bairesdev.com",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.engineer.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("SoftwareEngineer.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" && err.Error() != tt.errMsg {
				t.Errorf("SoftwareEngineer.Validate() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestSoftwareEngineer_TableName(t *testing.T) {
	engineer := SoftwareEngineer{}
	expected := "software_engineers"
	if got := engineer.TableName(); got != expected {
		t.Errorf("SoftwareEngineer.TableName() = %v, want %v", got, expected)
	}
}

func TestSoftwareEngineer_BeforeCreate(t *testing.T) {
	engineer := &SoftwareEngineer{
		Name:  "John Doe",
		Email: "john@bairesdev.com",
	}

	engineer.BeforeCreate()

	// Check that timestamps are set
	if engineer.CreatedAt.IsZero() {
		t.Error("SoftwareEngineer.BeforeCreate() did not set CreatedAt")
	}
	if engineer.UpdatedAt.IsZero() {
		t.Error("SoftwareEngineer.BeforeCreate() did not set UpdatedAt")
	}

	// Check that CreatedAt and UpdatedAt are approximately equal (within 1 second)
	diff := engineer.UpdatedAt.Sub(engineer.CreatedAt)
	if diff < 0 || diff > time.Second {
		t.Errorf("SoftwareEngineer.BeforeCreate() CreatedAt and UpdatedAt differ by %v, expected them to be nearly equal", diff)
	}
}

func TestSoftwareEngineer_BeforeUpdate(t *testing.T) {
	engineer := &SoftwareEngineer{
		Name:      "John Doe",
		Email:     "john@bairesdev.com",
		CreatedAt: time.Now().Add(-24 * time.Hour),
		UpdatedAt: time.Now().Add(-24 * time.Hour),
	}

	oldUpdatedAt := engineer.UpdatedAt
	time.Sleep(10 * time.Millisecond)

	engineer.BeforeUpdate()

	// Check that UpdatedAt was updated
	if !engineer.UpdatedAt.After(oldUpdatedAt) {
		t.Error("SoftwareEngineer.BeforeUpdate() did not update UpdatedAt")
	}
}

func TestSoftwareEngineer_HasConfirmedAddress(t *testing.T) {
	tests := []struct {
		name     string
		engineer SoftwareEngineer
		expected bool
	}{
		{
			name: "confirmed address",
			engineer: SoftwareEngineer{
				AddressConfirmed: true,
			},
			expected: true,
		},
		{
			name: "unconfirmed address",
			engineer: SoftwareEngineer{
				AddressConfirmed: false,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.engineer.HasConfirmedAddress(); got != tt.expected {
				t.Errorf("SoftwareEngineer.HasConfirmedAddress() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSoftwareEngineer_ConfirmAddress(t *testing.T) {
	engineer := &SoftwareEngineer{
		Name:             "John Doe",
		Email:            "john@bairesdev.com",
		AddressConfirmed: false,
	}

	if engineer.AddressConfirmed {
		t.Error("Expected AddressConfirmed to be false initially")
	}
	if engineer.AddressConfirmationAt != nil {
		t.Error("Expected AddressConfirmationAt to be nil initially")
	}

	engineer.ConfirmAddress()

	if !engineer.AddressConfirmed {
		t.Error("ConfirmAddress() did not set AddressConfirmed to true")
	}
	if engineer.AddressConfirmationAt == nil {
		t.Error("ConfirmAddress() did not set AddressConfirmationAt")
	}
	if engineer.AddressConfirmationAt.IsZero() {
		t.Error("ConfirmAddress() set AddressConfirmationAt to zero time")
	}
}

// Helper function for creating time pointers
func timePtr(t time.Time) *time.Time {
	return &t
}
