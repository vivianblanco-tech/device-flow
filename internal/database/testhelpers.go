package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

// SetupTestDB creates a test database connection
// It uses the TEST_DATABASE_URL environment variable if set,
// otherwise falls back to a default test database configuration
func SetupTestDB(t *testing.T) (*sql.DB, func()) {
	t.Helper()

	// Get test database URL from environment or use default
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:password@localhost:5432/laptop_tracking_test?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping test database: %v", err)
	}

	// Cleanup function to close the connection and clean up test data
	cleanup := func() {
		// Clean up test tables in reverse order of dependencies
		cleanupQueries := []string{
			"DELETE FROM sessions",
			"DELETE FROM magic_links",
			"DELETE FROM notification_logs",
			"DELETE FROM audit_logs",
			"DELETE FROM delivery_forms",
			"DELETE FROM reception_reports",
			"DELETE FROM pickup_forms",
			"DELETE FROM shipment_laptops",
			"DELETE FROM shipments",
			"DELETE FROM laptops",
			"DELETE FROM software_engineers",
			"DELETE FROM users",
			"DELETE FROM client_companies",
		}

		for _, query := range cleanupQueries {
			_, err := db.Exec(query)
			if err != nil {
				t.Logf("Cleanup warning: %v", err)
			}
		}

		db.Close()
	}

	return db, cleanup
}

// ExecTestSQL executes SQL statements for test setup
func ExecTestSQL(ctx context.Context, db *sql.DB, query string, args ...interface{}) error {
	_, err := db.ExecContext(ctx, query, args...)
	return err
}

// QueryRowTestSQL executes a query that returns a single row for testing
func QueryRowTestSQL(ctx context.Context, db *sql.DB, query string, args ...interface{}) *sql.Row {
	return db.QueryRowContext(ctx, query, args...)
}

// CreateTestDatabase creates a test database (for CI/CD setup)
func CreateTestDatabase(dbName string) error {
	// Connect to postgres database first
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}
	defer db.Close()

	// Drop existing test database if it exists
	_, err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName))
	if err != nil {
		return fmt.Errorf("failed to drop existing test database: %w", err)
	}

	// Create new test database
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
	if err != nil {
		return fmt.Errorf("failed to create test database: %w", err)
	}

	return nil
}

