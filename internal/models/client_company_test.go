package models

import (
	"testing"
	"time"
)

func TestClientCompany_Validate(t *testing.T) {
	tests := []struct {
		name    string
		company ClientCompany
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid company with all required fields",
			company: ClientCompany{
				Name:        "Acme Corporation",
				ContactInfo: "contact@acme.com, +1-555-0100",
			},
			wantErr: false,
		},
		{
			name: "valid company with minimal fields",
			company: ClientCompany{
				Name: "Tech Startup Inc",
			},
			wantErr: false,
		},
		{
			name: "invalid - missing name",
			company: ClientCompany{
				ContactInfo: "contact@company.com",
			},
			wantErr: true,
			errMsg:  "company name is required",
		},
		{
			name: "invalid - empty name",
			company: ClientCompany{
				Name:        "",
				ContactInfo: "contact@company.com",
			},
			wantErr: true,
			errMsg:  "company name is required",
		},
		{
			name: "invalid - name too short",
			company: ClientCompany{
				Name: "AB",
			},
			wantErr: true,
			errMsg:  "company name must be at least 3 characters",
		},
		{
			name: "valid - name exactly 3 characters",
			company: ClientCompany{
				Name: "ABC",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.company.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ClientCompany.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" && err.Error() != tt.errMsg {
				t.Errorf("ClientCompany.Validate() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestClientCompany_TableName(t *testing.T) {
	company := ClientCompany{}
	expected := "client_companies"
	if got := company.TableName(); got != expected {
		t.Errorf("ClientCompany.TableName() = %v, want %v", got, expected)
	}
}

func TestClientCompany_BeforeCreate(t *testing.T) {
	company := &ClientCompany{
		Name:        "Test Company",
		ContactInfo: "test@company.com",
	}

	company.BeforeCreate()

	// Check that timestamps are set
	if company.CreatedAt.IsZero() {
		t.Error("ClientCompany.BeforeCreate() did not set CreatedAt")
	}
	if company.UpdatedAt.IsZero() {
		t.Error("ClientCompany.BeforeCreate() did not set UpdatedAt")
	}

	// Check that CreatedAt and UpdatedAt are approximately equal (within 1 second)
	diff := company.UpdatedAt.Sub(company.CreatedAt)
	if diff < 0 || diff > time.Second {
		t.Errorf("ClientCompany.BeforeCreate() CreatedAt and UpdatedAt differ by %v, expected them to be nearly equal", diff)
	}
}

func TestClientCompany_BeforeUpdate(t *testing.T) {
	company := &ClientCompany{
		Name:        "Test Company",
		ContactInfo: "test@company.com",
		CreatedAt:   time.Now().Add(-24 * time.Hour), // 1 day ago
		UpdatedAt:   time.Now().Add(-24 * time.Hour),
	}

	oldUpdatedAt := company.UpdatedAt
	time.Sleep(10 * time.Millisecond) // Small delay to ensure time difference

	company.BeforeUpdate()

	// Check that UpdatedAt was updated
	if !company.UpdatedAt.After(oldUpdatedAt) {
		t.Error("ClientCompany.BeforeUpdate() did not update UpdatedAt")
	}

	// Check that CreatedAt was not modified
	if company.CreatedAt != time.Now().Add(-24*time.Hour).Truncate(time.Second) {
		// Allow small time drift
		diff := time.Now().Add(-24 * time.Hour).Sub(company.CreatedAt)
		if diff < -time.Second || diff > time.Second {
			t.Error("ClientCompany.BeforeUpdate() should not modify CreatedAt")
		}
	}
}

func TestClientCompany_GetActiveUsersCount(t *testing.T) {
	// This test will be implemented when we add user relationship methods
	// For now, we'll test that the function exists and returns 0 for a new company
	company := ClientCompany{
		Name: "Test Company",
	}

	count := company.GetActiveUsersCount()
	if count != 0 {
		t.Errorf("ClientCompany.GetActiveUsersCount() for new company = %v, want 0", count)
	}
}

