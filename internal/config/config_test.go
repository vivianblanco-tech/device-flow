package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Save original environment
	originalEnv := make(map[string]string)
	envVars := []string{"APP_ENV", "APP_PORT", "DB_HOST", "DB_PORT"}
	for _, key := range envVars {
		originalEnv[key] = os.Getenv(key)
	}

	// Restore environment after test
	defer func() {
		for key, value := range originalEnv {
			if value == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, value)
			}
		}
	}()

	// Test with default values
	t.Run("DefaultValues", func(t *testing.T) {
		// Clear environment variables
		for _, key := range envVars {
			os.Unsetenv(key)
		}

		cfg := Load()

		if cfg.App.Environment != "development" {
			t.Errorf("Expected environment 'development', got '%s'", cfg.App.Environment)
		}

		if cfg.Server.Port != "8080" {
			t.Errorf("Expected port '8080', got '%s'", cfg.Server.Port)
		}

		if cfg.Database.Host != "localhost" {
			t.Errorf("Expected DB host 'localhost', got '%s'", cfg.Database.Host)
		}
	})

	// Test with custom values
	t.Run("CustomValues", func(t *testing.T) {
		os.Setenv("APP_ENV", "production")
		os.Setenv("APP_PORT", "3000")
		os.Setenv("DB_HOST", "db.example.com")
		os.Setenv("DB_PORT", "5433")

		cfg := Load()

		if cfg.App.Environment != "production" {
			t.Errorf("Expected environment 'production', got '%s'", cfg.App.Environment)
		}

		if cfg.Server.Port != "3000" {
			t.Errorf("Expected port '3000', got '%s'", cfg.Server.Port)
		}

		if cfg.Database.Host != "db.example.com" {
			t.Errorf("Expected DB host 'db.example.com', got '%s'", cfg.Database.Host)
		}

		if cfg.Database.Port != "5433" {
			t.Errorf("Expected DB port '5433', got '%s'", cfg.Database.Port)
		}
	})
}

func TestGetEnvAsInt(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		value        string
		defaultValue int
		expected     int
	}{
		{"ValidInteger", "TEST_INT", "42", 10, 42},
		{"InvalidInteger", "TEST_INT", "invalid", 10, 10},
		{"EmptyString", "TEST_INT", "", 10, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != "" {
				os.Setenv(tt.key, tt.value)
			} else {
				os.Unsetenv(tt.key)
			}
			defer os.Unsetenv(tt.key)

			result := getEnvAsInt(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestGetEnvAsInt64(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		value        string
		defaultValue int64
		expected     int64
	}{
		{"ValidInt64", "TEST_INT64", "9999999999", 100, 9999999999},
		{"InvalidInt64", "TEST_INT64", "invalid", 100, 100},
		{"EmptyString", "TEST_INT64", "", 100, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != "" {
				os.Setenv(tt.key, tt.value)
			} else {
				os.Unsetenv(tt.key)
			}
			defer os.Unsetenv(tt.key)

			result := getEnvAsInt64(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}

