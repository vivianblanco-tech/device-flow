package database

import (
	"testing"

	"github.com/yourusername/laptop-tracking-system/internal/config"
)

func TestConnect(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("successful connection with valid config", func(t *testing.T) {
		cfg := config.DatabaseConfig{
			Host:     "localhost",
			Port:     "5432",
			User:     "postgres",
			Password: "password",
			Name:     "laptop_tracking_test",
			SSLMode:  "disable",
		}

		db, err := Connect(cfg)
		if err != nil {
			t.Fatalf("Failed to connect to database: %v", err)
		}
		defer db.Close()

		// Test ping
		err = db.Ping()
		if err != nil {
			t.Errorf("Failed to ping database: %v", err)
		}

		// Verify connection pool settings
		stats := db.Stats()
		if stats.MaxOpenConnections != 25 {
			t.Errorf("Expected MaxOpenConnections to be 25, got %d", stats.MaxOpenConnections)
		}
	})

	t.Run("connection fails with invalid host", func(t *testing.T) {
		cfg := config.DatabaseConfig{
			Host:     "invalid-host-that-does-not-exist",
			Port:     "5432",
			User:     "postgres",
			Password: "password",
			Name:     "laptop_tracking_test",
			SSLMode:  "disable",
		}

		db, err := Connect(cfg)
		if err == nil {
			db.Close()
			t.Error("Expected error with invalid host, got nil")
		}
	})

	t.Run("connection fails with invalid database name", func(t *testing.T) {
		cfg := config.DatabaseConfig{
			Host:     "localhost",
			Port:     "5432",
			User:     "postgres",
			Password: "password",
			Name:     "nonexistent_database_12345",
			SSLMode:  "disable",
		}

		db, err := Connect(cfg)
		if err == nil {
			db.Close()
			t.Error("Expected error with invalid database name, got nil")
		}
	})

	t.Run("connection fails with invalid credentials", func(t *testing.T) {
		cfg := config.DatabaseConfig{
			Host:     "localhost",
			Port:     "5432",
			User:     "invalid_user",
			Password: "invalid_password",
			Name:     "laptop_tracking_test",
			SSLMode:  "disable",
		}

		db, err := Connect(cfg)
		if err == nil {
			db.Close()
			t.Error("Expected error with invalid credentials, got nil")
		}
	})
}

func TestDatabaseConnectionPool(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	cfg := config.DatabaseConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "password",
		Name:     "laptop_tracking_test",
		SSLMode:  "disable",
	}

	db, err := Connect(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	t.Run("connection pool has correct settings", func(t *testing.T) {
		stats := db.Stats()

		// Check MaxOpenConnections
		if stats.MaxOpenConnections != 25 {
			t.Errorf("Expected MaxOpenConnections to be 25, got %d", stats.MaxOpenConnections)
		}
	})

	t.Run("can execute queries", func(t *testing.T) {
		var result int
		err := db.QueryRow("SELECT 1").Scan(&result)
		if err != nil {
			t.Errorf("Failed to execute query: %v", err)
		}
		if result != 1 {
			t.Errorf("Expected result to be 1, got %d", result)
		}
	})
}
