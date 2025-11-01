package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/yourusername/laptop-tracking-system/internal/config"
	"github.com/yourusername/laptop-tracking-system/internal/database"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
)

func main() {
	fmt.Println(colorCyan + "=== Database Connection Test ===" + colorReset)
	fmt.Println()

	// Load environment variables
	fmt.Print(colorYellow + "[1/4] Loading environment variables..." + colorReset)
	if err := godotenv.Load(); err != nil {
		fmt.Println(colorYellow + " (using system environment)" + colorReset)
	} else {
		fmt.Println(colorGreen + " ✓" + colorReset)
	}

	// Load configuration
	fmt.Print(colorYellow + "[2/4] Loading configuration..." + colorReset)
	cfg := config.Load()
	fmt.Println(colorGreen + " ✓" + colorReset)

	// Display connection parameters
	fmt.Println()
	fmt.Println(colorCyan + "Connection Parameters:" + colorReset)
	fmt.Printf("  Host:     %s\n", cfg.Database.Host)
	fmt.Printf("  Port:     %s\n", cfg.Database.Port)
	fmt.Printf("  Database: %s\n", cfg.Database.Name)
	fmt.Printf("  User:     %s\n", cfg.Database.User)
	fmt.Printf("  SSL Mode: %s\n", cfg.Database.SSLMode)
	fmt.Println()

	// Attempt database connection
	fmt.Print(colorYellow + "[3/4] Connecting to database..." + colorReset)
	db, err := database.Connect(cfg.Database)
	if err != nil {
		fmt.Println(colorRed + " ✗" + colorReset)
		fmt.Println()
		fmt.Println(colorRed + "Connection Failed!" + colorReset)
		fmt.Printf("Error: %v\n", err)
		fmt.Println()
		printTroubleshootingTips(err)
		os.Exit(1)
	}
	defer db.Close()
	fmt.Println(colorGreen + " ✓" + colorReset)

	// Query database information
	fmt.Print(colorYellow + "[4/4] Querying database information..." + colorReset)
	if err := queryDatabaseInfo(db); err != nil {
		fmt.Println(colorRed + " ✗" + colorReset)
		fmt.Printf("Error querying database: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(colorGreen + " ✓" + colorReset)

	// Success summary
	fmt.Println()
	fmt.Println(colorGreen + "========================================" + colorReset)
	fmt.Println(colorGreen + "  DATABASE CONNECTION SUCCESSFUL! ✓" + colorReset)
	fmt.Println(colorGreen + "========================================" + colorReset)
	fmt.Println()
	fmt.Println("Your database is properly configured and accessible.")
}

func queryDatabaseInfo(db *sql.DB) error {
	fmt.Println()
	fmt.Println(colorCyan + "Database Information:" + colorReset)

	// Get PostgreSQL version
	var version string
	err := db.QueryRow("SELECT version()").Scan(&version)
	if err != nil {
		return fmt.Errorf("failed to get PostgreSQL version: %w", err)
	}
	fmt.Printf("  Version:  %s\n", version[:80]+"...")

	// Get current database
	var currentDB string
	err = db.QueryRow("SELECT current_database()").Scan(&currentDB)
	if err != nil {
		return fmt.Errorf("failed to get current database: %w", err)
	}
	fmt.Printf("  Current:  %s\n", currentDB)

	// Get current user
	var currentUser string
	err = db.QueryRow("SELECT current_user").Scan(&currentUser)
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}
	fmt.Printf("  User:     %s\n", currentUser)

	// Count tables in public schema
	var tableCount int
	err = db.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public'").Scan(&tableCount)
	if err != nil {
		log.Printf("Warning: failed to count tables: %v", err)
	} else {
		fmt.Printf("  Tables:   %d in public schema\n", tableCount)
	}

	return nil
}

func printTroubleshootingTips(err error) {
	fmt.Println(colorYellow + "Troubleshooting Tips:" + colorReset)
	fmt.Println()

	errStr := err.Error()

	if contains(errStr, "connection refused") {
		fmt.Println("  • PostgreSQL server is not running")
		fmt.Println("    → Start PostgreSQL service")
		fmt.Println("    → Check if PostgreSQL is listening on the specified port")
		fmt.Println()
	}

	if contains(errStr, "database") && contains(errStr, "does not exist") {
		fmt.Println("  • The database does not exist")
		fmt.Println("    → Create the database using: createdb laptop_tracking_dev")
		fmt.Println("    → Or run the database setup script")
		fmt.Println()
	}

	if contains(errStr, "authentication failed") || contains(errStr, "password") {
		fmt.Println("  • Authentication failed")
		fmt.Println("    → Check username and password in .env file")
		fmt.Println("    → Verify PostgreSQL user permissions")
		fmt.Println()
	}

	if contains(errStr, "no such host") || contains(errStr, "unknown host") {
		fmt.Println("  • Cannot resolve hostname")
		fmt.Println("    → Check DB_HOST in .env file")
		fmt.Println("    → Verify network connectivity")
		fmt.Println()
	}

	fmt.Println("  General checks:")
	fmt.Println("    1. Ensure PostgreSQL is installed and running")
	fmt.Println("    2. Verify .env file exists and has correct values")
	fmt.Println("    3. Check firewall settings")
	fmt.Println("    4. Review PostgreSQL logs for more details")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
