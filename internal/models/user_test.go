package models

import (
	"testing"
	"time"
)

func TestUser_Validate(t *testing.T) {
	tests := []struct {
		name    string
		user    User
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid user with all fields",
			user: User{
				Email:        "test@bairesdev.com",
				PasswordHash: "hashedpassword123",
				Role:         RoleLogistics,
			},
			wantErr: false,
		},
		{
			name: "valid user with Google ID",
			user: User{
				Email:    "test@bairesdev.com",
				Role:     RoleClient,
				GoogleID: stringPtr("google_oauth_id_123"),
			},
			wantErr: false,
		},
		{
			name: "invalid - missing email",
			user: User{
				PasswordHash: "hashedpassword123",
				Role:         RoleLogistics,
			},
			wantErr: true,
			errMsg:  "email is required",
		},
		{
			name: "invalid - invalid email format",
			user: User{
				Email:        "invalid-email",
				PasswordHash: "hashedpassword123",
				Role:         RoleLogistics,
			},
			wantErr: true,
			errMsg:  "invalid email format",
		},
		{
			name: "invalid - missing role",
			user: User{
				Email:        "test@bairesdev.com",
				PasswordHash: "hashedpassword123",
			},
			wantErr: true,
			errMsg:  "role is required",
		},
		{
			name: "invalid - invalid role",
			user: User{
				Email:        "test@bairesdev.com",
				PasswordHash: "hashedpassword123",
				Role:         "invalid_role",
			},
			wantErr: true,
			errMsg:  "invalid role",
		},
		{
			name: "invalid - missing both password and Google ID",
			user: User{
				Email: "test@bairesdev.com",
				Role:  RoleLogistics,
			},
			wantErr: true,
			errMsg:  "either password_hash or google_id must be provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("User.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" && err.Error() != tt.errMsg {
				t.Errorf("User.Validate() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestUser_IsValidRole(t *testing.T) {
	tests := []struct {
		name string
		role UserRole
		want bool
	}{
		{"logistics role", RoleLogistics, true},
		{"client role", RoleClient, true},
		{"warehouse role", RoleWarehouse, true},
		{"project_manager role", RoleProjectManager, true},
		{"invalid role", "admin", false},
		{"empty role", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidRole(tt.role); got != tt.want {
				t.Errorf("IsValidRole() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_HasRole(t *testing.T) {
	user := User{
		Email: "test@bairesdev.com",
		Role:  RoleLogistics,
	}

	tests := []struct {
		name     string
		role     UserRole
		expected bool
	}{
		{"matching role", RoleLogistics, true},
		{"non-matching role", RoleClient, false},
		{"warehouse role", RoleWarehouse, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := user.HasRole(tt.role); got != tt.expected {
				t.Errorf("User.HasRole() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestUser_IsGoogleUser(t *testing.T) {
	tests := []struct {
		name     string
		user     User
		expected bool
	}{
		{
			name: "user with Google ID",
			user: User{
				Email:    "test@bairesdev.com",
				GoogleID: stringPtr("google_id_123"),
				Role:     RoleClient,
			},
			expected: true,
		},
		{
			name: "user without Google ID",
			user: User{
				Email:        "test@bairesdev.com",
				PasswordHash: "hashedpassword",
				Role:         RoleLogistics,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.user.IsGoogleUser(); got != tt.expected {
				t.Errorf("User.IsGoogleUser() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestUser_TableName(t *testing.T) {
	user := User{}
	expected := "users"
	if got := user.TableName(); got != expected {
		t.Errorf("User.TableName() = %v, want %v", got, expected)
	}
}

// Helper function for creating string pointers
func stringPtr(s string) *string {
	return &s
}

func TestUser_BeforeCreate(t *testing.T) {
	user := &User{
		Email:        "test@bairesdev.com",
		PasswordHash: "hashedpassword",
		Role:         RoleLogistics,
	}

	user.BeforeCreate()

	// Check that timestamps are set
	if user.CreatedAt.IsZero() {
		t.Error("User.BeforeCreate() did not set CreatedAt")
	}
	if user.UpdatedAt.IsZero() {
		t.Error("User.BeforeCreate() did not set UpdatedAt")
	}

	// Check that CreatedAt and UpdatedAt are approximately equal (within 1 second)
	diff := user.UpdatedAt.Sub(user.CreatedAt)
	if diff < 0 || diff > time.Second {
		t.Errorf("User.BeforeCreate() CreatedAt and UpdatedAt differ by %v, expected them to be nearly equal", diff)
	}
}

func TestUser_BeforeUpdate(t *testing.T) {
	user := &User{
		Email:        "test@bairesdev.com",
		PasswordHash: "hashedpassword",
		Role:         RoleLogistics,
		CreatedAt:    time.Now().Add(-24 * time.Hour), // 1 day ago
		UpdatedAt:    time.Now().Add(-24 * time.Hour),
	}

	oldUpdatedAt := user.UpdatedAt
	time.Sleep(10 * time.Millisecond) // Small delay to ensure time difference

	user.BeforeUpdate()

	// Check that UpdatedAt was updated
	if !user.UpdatedAt.After(oldUpdatedAt) {
		t.Error("User.BeforeUpdate() did not update UpdatedAt")
	}

	// Check that CreatedAt was not modified
	if user.CreatedAt != time.Now().Add(-24*time.Hour).Truncate(time.Second) {
		// Allow small time drift
		diff := time.Now().Add(-24 * time.Hour).Sub(user.CreatedAt)
		if diff < -time.Second || diff > time.Second {
			t.Error("User.BeforeUpdate() should not modify CreatedAt")
		}
	}
}
