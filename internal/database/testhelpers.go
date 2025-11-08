package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sync"
	"testing"

	_ "github.com/lib/pq"
)

// testDBMutex ensures that only one test accesses the database at a time
// This prevents race conditions when tests run in parallel
var testDBMutex sync.Mutex

// SetupTestDB creates a test database connection
// It uses the TEST_DATABASE_URL environment variable if set,
// otherwise falls back to a default test database configuration
func SetupTestDB(t *testing.T) (*sql.DB, func()) {
	t.Helper()

	// Lock the mutex to ensure only one test accesses the database at a time
	testDBMutex.Lock()

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

	// Clean up test tables in reverse order of dependencies BEFORE the test runs
	// This ensures each test starts with a clean slate, preventing race conditions
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
			t.Logf("Pre-cleanup warning: %v", err)
		}
	}

	// Cleanup function to close the connection and clean up test data after the test
	cleanup := func() {
		// Clean up test tables again after test completion
		for _, query := range cleanupQueries {
			_, err := db.Exec(query)
			if err != nil {
				t.Logf("Post-cleanup warning: %v", err)
			}
		}

		db.Close()
		
		// Unlock the mutex to allow other tests to run
		testDBMutex.Unlock()
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

