package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password",
			password: "securePassword123",
			wantErr:  false,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  true,
		},
		{
			name:     "short password",
			password: "123",
			wantErr:  false,
		},
		{
			name:     "long password",
			password: "thisIsAVeryLongPasswordThatShouldStillWork1234567890!@#$%^&*()",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if hash == "" {
					t.Error("HashPassword() returned empty hash")
				}
				if hash == tt.password {
					t.Error("HashPassword() returned plain password instead of hash")
				}
				// bcrypt hashes should start with $2a$, $2b$, or $2y$
				if len(hash) < 60 {
					t.Errorf("HashPassword() hash length = %d, expected at least 60", len(hash))
				}
			}
		})
	}
}

func TestCheckPasswordHash(t *testing.T) {
	// Generate a known hash for testing
	password := "mySecurePassword123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to generate hash for testing: %v", err)
	}

	tests := []struct {
		name     string
		password string
		hash     string
		want     bool
	}{
		{
			name:     "correct password",
			password: password,
			hash:     hash,
			want:     true,
		},
		{
			name:     "incorrect password",
			password: "wrongPassword",
			hash:     hash,
			want:     false,
		},
		{
			name:     "empty password",
			password: "",
			hash:     hash,
			want:     false,
		},
		{
			name:     "empty hash",
			password: password,
			hash:     "",
			want:     false,
		},
		{
			name:     "case sensitive password",
			password: "MYSECUREPASSWORD123",
			hash:     hash,
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckPasswordHash(tt.password, tt.hash); got != tt.want {
				t.Errorf("CheckPasswordHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHashPasswordConsistency(t *testing.T) {
	password := "testPassword123"

	// Hash the same password multiple times
	hash1, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	hash2, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Hashes should be different (bcrypt uses random salt)
	if hash1 == hash2 {
		t.Error("HashPassword() generated identical hashes, expected different salts")
	}

	// But both should validate against the same password
	if !CheckPasswordHash(password, hash1) {
		t.Error("First hash failed to validate")
	}
	if !CheckPasswordHash(password, hash2) {
		t.Error("Second hash failed to validate")
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "valid password - minimum length",
			password: "Pass123!",
			wantErr:  false,
		},
		{
			name:     "valid password - with special chars",
			password: "MyP@ssw0rd!",
			wantErr:  false,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  true,
			errMsg:   "password is required",
		},
		{
			name:     "too short password",
			password: "Pass12",
			wantErr:  true,
			errMsg:   "password must be at least 8 characters",
		},
		{
			name:     "password without uppercase",
			password: "password123!",
			wantErr:  true,
			errMsg:   "password must contain at least one uppercase letter",
		},
		{
			name:     "password without lowercase",
			password: "PASSWORD123!",
			wantErr:  true,
			errMsg:   "password must contain at least one lowercase letter",
		},
		{
			name:     "password without digit",
			password: "Password!",
			wantErr:  true,
			errMsg:   "password must contain at least one digit",
		},
		{
			name:     "password without special char",
			password: "Password123",
			wantErr:  true,
			errMsg:   "password must contain at least one special character",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("ValidatePassword() error = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}
