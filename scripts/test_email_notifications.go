package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/yourusername/laptop-tracking-system/internal/email"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
)

// MailhogMessage represents a message from Mailhog API
type MailhogMessage struct {
	ID   string `json:"ID"`
	From struct {
		Mailbox string `json:"Mailbox"`
		Domain  string `json:"Domain"`
	} `json:"From"`
	To []struct {
		Mailbox string `json:"Mailbox"`
		Domain  string `json:"Domain"`
	} `json:"To"`
	Content struct {
		Headers map[string][]string `json:"Headers"`
		Body    string              `json:"Body"`
	} `json:"Content"`
	Created time.Time `json:"Created"`
}

// MailhogResponse represents the response from Mailhog API
type MailhogResponse struct {
	Total int              `json:"total"`
	Count int              `json:"count"`
	Start int              `json:"start"`
	Items []MailhogMessage `json:"items"`
}

// TestResult holds the result of a notification test
type TestResult struct {
	NotificationType string
	Success          bool
	Error            string
	EmailID          string
	Subject          string
	Recipient        string
	TimeSent         time.Time
}

func main() {
	fmt.Println(colorCyan + "==================================================" + colorReset)
	fmt.Println(colorCyan + "  Email Notifications Test - Mailhog Verification" + colorReset)
	fmt.Println(colorCyan + "==================================================" + colorReset)
	fmt.Println()

	// Load configuration from environment
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/laptop_tracking_test?sslmode=disable"
		fmt.Println(colorYellow + "⚠️  Using default DATABASE_URL (set DATABASE_URL env var to override)" + colorReset)
	}

	mailhogURL := os.Getenv("MAILHOG_URL")
	if mailhogURL == "" {
		mailhogURL = "http://localhost:8025"
		fmt.Println(colorYellow + "⚠️  Using default MAILHOG_URL (set MAILHOG_URL env var to override)" + colorReset)
	}

	smtpHost := os.Getenv("SMTP_HOST")
	if smtpHost == "" {
		smtpHost = "localhost"
	}

	smtpPort := os.Getenv("SMTP_PORT")
	if smtpPort == "" {
		smtpPort = "1025"
	}

	fmt.Println()
	fmt.Println(colorBlue + "Configuration:" + colorReset)
	fmt.Printf("  Database: %s\n", dbURL)
	fmt.Printf("  Mailhog:  %s\n", mailhogURL)
	fmt.Printf("  SMTP:     %s:%s\n", smtpHost, smtpPort)
	fmt.Println()

	// Connect to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf(colorRed+"❌ Failed to connect to database: %v\n"+colorReset, err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		fmt.Printf(colorRed+"❌ Failed to ping database: %v\n"+colorReset, err)
		os.Exit(1)
	}

	fmt.Println(colorGreen + "✅ Connected to database" + colorReset)

	// Check Mailhog is accessible
	if err := checkMailhog(mailhogURL); err != nil {
		fmt.Printf(colorRed+"❌ Failed to connect to Mailhog: %v\n"+colorReset, err)
		fmt.Println(colorYellow + "💡 Make sure Mailhog is running: mailhog" + colorReset)
		os.Exit(1)
	}

	fmt.Println(colorGreen + "✅ Connected to Mailhog" + colorReset)
	fmt.Println()

	// Clear Mailhog messages
	if err := clearMailhog(mailhogURL); err != nil {
		fmt.Printf(colorYellow+"⚠️  Could not clear Mailhog: %v\n"+colorReset, err)
	} else {
		fmt.Println(colorGreen + "✅ Cleared Mailhog messages" + colorReset)
	}
	fmt.Println()

	// Create email client
	emailClient, err := email.NewClient(email.Config{
		Host:     smtpHost,
		Port:     1025, // Mailhog default port
		Username: "",
		Password: "",
		From:     "test@localhost",
	})
	if err != nil {
		fmt.Printf(colorRed+"❌ Failed to create email client: %v\n"+colorReset, err)
		os.Exit(1)
	}

	// Create notifier
	notifier := email.NewNotifier(emailClient, db)

	// Set up test data
	ctx := context.Background()
	testData, err := setupTestData(ctx, db)
	if err != nil {
		fmt.Printf(colorRed+"❌ Failed to set up test data: %v\n"+colorReset, err)
		os.Exit(1)
	}
	fmt.Println(colorGreen + "✅ Test data created" + colorReset)
	fmt.Println()

	// Run tests
	results := []TestResult{}

	fmt.Println(colorPurple + "===========================================" + colorReset)
	fmt.Println(colorPurple + "  Running Email Notification Tests" + colorReset)
	fmt.Println(colorPurple + "===========================================" + colorReset)
	fmt.Println()

	// Test 1: Magic Link
	fmt.Println(colorBlue + "📧 Test 1: Magic Link Email" + colorReset)
	result := testMagicLink(ctx, notifier, mailhogURL)
	results = append(results, result)
	printTestResult(result)
	time.Sleep(500 * time.Millisecond)

	// Test 2: Pickup Confirmation
	fmt.Println(colorBlue + "📧 Test 2: Pickup Confirmation" + colorReset)
	result = testPickupConfirmation(ctx, notifier, mailhogURL, testData.ShipmentID)
	results = append(results, result)
	printTestResult(result)
	time.Sleep(500 * time.Millisecond)

	// Test 3: Pickup Scheduled
	fmt.Println(colorBlue + "📧 Test 3: Pickup Scheduled Notification" + colorReset)
	result = testPickupScheduled(ctx, notifier, mailhogURL, testData.ShipmentID)
	results = append(results, result)
	printTestResult(result)
	time.Sleep(500 * time.Millisecond)

	// Test 4: Warehouse Pre-Alert
	fmt.Println(colorBlue + "📧 Test 4: Warehouse Pre-Alert" + colorReset)
	result = testWarehousePreAlert(ctx, notifier, mailhogURL, testData.ShipmentID)
	results = append(results, result)
	printTestResult(result)
	time.Sleep(500 * time.Millisecond)

	// Test 5: Release Notification
	fmt.Println(colorBlue + "📧 Test 5: Release Notification" + colorReset)
	result = testReleaseNotification(ctx, notifier, mailhogURL, testData.ShipmentID)
	results = append(results, result)
	printTestResult(result)
	time.Sleep(500 * time.Millisecond)

	// Test 6: Delivery Confirmation
	fmt.Println(colorBlue + "📧 Test 6: Delivery Confirmation" + colorReset)
	result = testDeliveryConfirmation(ctx, notifier, mailhogURL, testData.ShipmentID)
	results = append(results, result)
	printTestResult(result)
	time.Sleep(500 * time.Millisecond)

	// Print summary
	fmt.Println()
	printSummary(results)

	// Clean up test data
	if err := cleanupTestData(ctx, db, testData); err != nil {
		fmt.Printf(colorYellow+"⚠️  Warning: Failed to clean up test data: %v\n"+colorReset, err)
	} else {
		fmt.Println(colorGreen + "✅ Test data cleaned up" + colorReset)
	}

	// Exit with appropriate code
	allPassed := true
	for _, r := range results {
		if !r.Success {
			allPassed = false
			break
		}
	}

	if allPassed {
		fmt.Println()
		fmt.Println(colorGreen + "🎉 All email notifications passed!" + colorReset)
		os.Exit(0)
	} else {
		fmt.Println()
		fmt.Println(colorRed + "❌ Some email notifications failed" + colorReset)
		os.Exit(1)
	}
}

// TestData holds IDs of created test data
type TestData struct {
	ShipmentID    int64
	ClientID      int64
	EngineerID    int64
	WarehouseUser int64
	LogisticsUser int64
	LaptopID      int64
}

func setupTestData(ctx context.Context, db *sql.DB) (*TestData, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	data := &TestData{}

	// Create client company
	err = tx.QueryRowContext(ctx,
		`INSERT INTO client_companies (name, contact_info, created_at, updated_at)
		VALUES ($1, $2, $3, $4) RETURNING id`,
		"Test Client Company", "Test Contact Info", time.Now(), time.Now(),
	).Scan(&data.ClientID)
	if err != nil {
		return nil, fmt.Errorf("create client company: %w", err)
	}

	// Create software engineer
	err = tx.QueryRowContext(ctx,
		`INSERT INTO software_engineers (name, email, address, city, state, zip_code, country, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
		"Test Engineer", "engineer@test.com", "123 Test St", "Test City", "TS", "12345", "Test Country",
		time.Now(), time.Now(),
	).Scan(&data.EngineerID)
	if err != nil {
		return nil, fmt.Errorf("create engineer: %w", err)
	}

	// Create warehouse user
	err = tx.QueryRowContext(ctx,
		`INSERT INTO users (email, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4) RETURNING id`,
		"warehouse@test.com", "warehouse", time.Now(), time.Now(),
	).Scan(&data.WarehouseUser)
	if err != nil {
		return nil, fmt.Errorf("create warehouse user: %w", err)
	}

	// Create logistics user
	err = tx.QueryRowContext(ctx,
		`INSERT INTO users (email, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4) RETURNING id`,
		"logistics@test.com", "logistics", time.Now(), time.Now(),
	).Scan(&data.LogisticsUser)
	if err != nil {
		return nil, fmt.Errorf("create logistics user: %w", err)
	}

	// Create laptop
	err = tx.QueryRowContext(ctx,
		`INSERT INTO laptops (serial_number, model, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"TEST-SERIAL-123", "Test Model", "available", time.Now(), time.Now(),
	).Scan(&data.LaptopID)
	if err != nil {
		return nil, fmt.Errorf("create laptop: %w", err)
	}

	// Create shipment
	err = tx.QueryRowContext(ctx,
		`INSERT INTO shipments (client_company_id, software_engineer_id, status, tracking_number, 
		pickup_scheduled_date, created_at, updated_at, shipment_type)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
		data.ClientID, data.EngineerID, "pending_pickup", "TEST-TRACK-123",
		time.Now().AddDate(0, 0, 1), time.Now(), time.Now(), "full_journey",
	).Scan(&data.ShipmentID)
	if err != nil {
		return nil, fmt.Errorf("create shipment: %w", err)
	}

	// Link laptop to shipment
	_, err = tx.ExecContext(ctx,
		`INSERT INTO shipment_laptops (shipment_id, laptop_id, created_at)
		VALUES ($1, $2, $3)`,
		data.ShipmentID, data.LaptopID, time.Now(),
	)
	if err != nil {
		return nil, fmt.Errorf("link laptop to shipment: %w", err)
	}

	// Create pickup form with contact email
	formData := map[string]interface{}{
		"contact_name":     "Test Contact",
		"contact_email":    "contact@test.com",
		"pickup_address":   "123 Pickup St",
		"pickup_city":      "Pickup City",
		"pickup_state":     "PS",
		"pickup_zip":       "54321",
		"pickup_time_slot": "morning",
	}
	formDataJSON, _ := json.Marshal(formData)

	_, err = tx.ExecContext(ctx,
		`INSERT INTO pickup_forms (shipment_id, submitted_by_user_id, submitted_at, form_data)
		VALUES ($1, $2, $3, $4)`,
		data.ShipmentID, data.WarehouseUser, time.Now(), formDataJSON,
	)
	if err != nil {
		return nil, fmt.Errorf("create pickup form: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return data, nil
}

func cleanupTestData(ctx context.Context, db *sql.DB, data *TestData) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete in reverse order of creation
	_, _ = tx.ExecContext(ctx, "DELETE FROM notification_logs WHERE shipment_id = $1", data.ShipmentID)
	_, _ = tx.ExecContext(ctx, "DELETE FROM pickup_forms WHERE shipment_id = $1", data.ShipmentID)
	_, _ = tx.ExecContext(ctx, "DELETE FROM shipment_laptops WHERE shipment_id = $1", data.ShipmentID)
	_, _ = tx.ExecContext(ctx, "DELETE FROM shipments WHERE id = $1", data.ShipmentID)
	_, _ = tx.ExecContext(ctx, "DELETE FROM laptops WHERE id = $1", data.LaptopID)
	_, _ = tx.ExecContext(ctx, "DELETE FROM users WHERE id = $1", data.WarehouseUser)
	_, _ = tx.ExecContext(ctx, "DELETE FROM users WHERE id = $1", data.LogisticsUser)
	_, _ = tx.ExecContext(ctx, "DELETE FROM software_engineers WHERE id = $1", data.EngineerID)
	_, _ = tx.ExecContext(ctx, "DELETE FROM client_companies WHERE id = $1", data.ClientID)

	return tx.Commit()
}

func checkMailhog(mailhogURL string) error {
	resp, err := http.Get(mailhogURL + "/api/v2/messages")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func clearMailhog(mailhogURL string) error {
	req, err := http.NewRequest("DELETE", mailhogURL+"/api/v1/messages", nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func getLatestMailhogMessage(mailhogURL string) (*MailhogMessage, error) {
	resp, err := http.Get(mailhogURL + "/api/v2/messages?limit=1")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var mailhogResp MailhogResponse
	if err := json.Unmarshal(body, &mailhogResp); err != nil {
		return nil, err
	}

	if mailhogResp.Count == 0 {
		return nil, fmt.Errorf("no messages found")
	}

	return &mailhogResp.Items[0], nil
}

func testMagicLink(ctx context.Context, notifier *email.Notifier, mailhogURL string) TestResult {
	result := TestResult{NotificationType: "Magic Link"}

	err := notifier.SendMagicLink(ctx, "test@example.com", "Test User",
		"https://example.com/magic-link", "pickup", time.Now().Add(24*time.Hour))

	if err != nil {
		result.Error = err.Error()
		return result
	}

	// Wait a moment for email to arrive
	time.Sleep(200 * time.Millisecond)

	msg, err := getLatestMailhogMessage(mailhogURL)
	if err != nil {
		result.Error = fmt.Sprintf("Email not found in Mailhog: %v", err)
		return result
	}

	result.EmailID = msg.ID
	result.TimeSent = msg.Created
	result.Recipient = fmt.Sprintf("%s@%s", msg.To[0].Mailbox, msg.To[0].Domain)

	if headers := msg.Content.Headers["Subject"]; len(headers) > 0 {
		result.Subject = headers[0]
	}

	// Verify subject contains expected text
	if !strings.Contains(result.Subject, "Access Your Form") {
		result.Error = fmt.Sprintf("Unexpected subject: %s", result.Subject)
		return result
	}

	// Verify recipient
	if result.Recipient != "test@example.com" {
		result.Error = fmt.Sprintf("Wrong recipient: expected test@example.com, got %s", result.Recipient)
		return result
	}

	result.Success = true
	return result
}

func testPickupConfirmation(ctx context.Context, notifier *email.Notifier, mailhogURL string, shipmentID int64) TestResult {
	result := TestResult{NotificationType: "Pickup Confirmation"}

	err := notifier.SendPickupConfirmation(ctx, shipmentID)
	if err != nil {
		result.Error = err.Error()
		return result
	}

	time.Sleep(200 * time.Millisecond)

	msg, err := getLatestMailhogMessage(mailhogURL)
	if err != nil {
		result.Error = fmt.Sprintf("Email not found in Mailhog: %v", err)
		return result
	}

	result.EmailID = msg.ID
	result.TimeSent = msg.Created
	result.Recipient = fmt.Sprintf("%s@%s", msg.To[0].Mailbox, msg.To[0].Domain)

	if headers := msg.Content.Headers["Subject"]; len(headers) > 0 {
		result.Subject = headers[0]
	}

	if !strings.Contains(result.Subject, "Pickup Confirmation") {
		result.Error = fmt.Sprintf("Unexpected subject: %s", result.Subject)
		return result
	}

	result.Success = true
	return result
}

func testPickupScheduled(ctx context.Context, notifier *email.Notifier, mailhogURL string, shipmentID int64) TestResult {
	result := TestResult{NotificationType: "Pickup Scheduled"}

	err := notifier.SendPickupScheduledNotification(ctx, shipmentID)
	if err != nil {
		result.Error = err.Error()
		return result
	}

	time.Sleep(200 * time.Millisecond)

	msg, err := getLatestMailhogMessage(mailhogURL)
	if err != nil {
		result.Error = fmt.Sprintf("Email not found in Mailhog: %v", err)
		return result
	}

	result.EmailID = msg.ID
	result.TimeSent = msg.Created
	result.Recipient = fmt.Sprintf("%s@%s", msg.To[0].Mailbox, msg.To[0].Domain)

	if headers := msg.Content.Headers["Subject"]; len(headers) > 0 {
		result.Subject = headers[0]
	}

	if !strings.Contains(result.Subject, "Pickup Scheduled") {
		result.Error = fmt.Sprintf("Unexpected subject: %s", result.Subject)
		return result
	}

	if result.Recipient != "contact@test.com" {
		result.Error = fmt.Sprintf("Wrong recipient: expected contact@test.com, got %s", result.Recipient)
		return result
	}

	result.Success = true
	return result
}

func testWarehousePreAlert(ctx context.Context, notifier *email.Notifier, mailhogURL string, shipmentID int64) TestResult {
	result := TestResult{NotificationType: "Warehouse Pre-Alert"}

	err := notifier.SendWarehousePreAlert(ctx, shipmentID)
	if err != nil {
		result.Error = err.Error()
		return result
	}

	time.Sleep(200 * time.Millisecond)

	msg, err := getLatestMailhogMessage(mailhogURL)
	if err != nil {
		result.Error = fmt.Sprintf("Email not found in Mailhog: %v", err)
		return result
	}

	result.EmailID = msg.ID
	result.TimeSent = msg.Created
	result.Recipient = fmt.Sprintf("%s@%s", msg.To[0].Mailbox, msg.To[0].Domain)

	if headers := msg.Content.Headers["Subject"]; len(headers) > 0 {
		result.Subject = headers[0]
	}

	if !strings.Contains(result.Subject, "Incoming Shipment Alert") {
		result.Error = fmt.Sprintf("Unexpected subject: %s", result.Subject)
		return result
	}

	if result.Recipient != "warehouse@test.com" {
		result.Error = fmt.Sprintf("Wrong recipient: expected warehouse@test.com, got %s", result.Recipient)
		return result
	}

	result.Success = true
	return result
}

func testReleaseNotification(ctx context.Context, notifier *email.Notifier, mailhogURL string, shipmentID int64) TestResult {
	result := TestResult{NotificationType: "Release Notification"}

	err := notifier.SendReleaseNotification(ctx, shipmentID)
	if err != nil {
		result.Error = err.Error()
		return result
	}

	time.Sleep(200 * time.Millisecond)

	msg, err := getLatestMailhogMessage(mailhogURL)
	if err != nil {
		result.Error = fmt.Sprintf("Email not found in Mailhog: %v", err)
		return result
	}

	result.EmailID = msg.ID
	result.TimeSent = msg.Created
	result.Recipient = fmt.Sprintf("%s@%s", msg.To[0].Mailbox, msg.To[0].Domain)

	if headers := msg.Content.Headers["Subject"]; len(headers) > 0 {
		result.Subject = headers[0]
	}

	if !strings.Contains(result.Subject, "Hardware Release for Pickup") {
		result.Error = fmt.Sprintf("Unexpected subject: %s", result.Subject)
		return result
	}

	if result.Recipient != "logistics@test.com" {
		result.Error = fmt.Sprintf("Wrong recipient: expected logistics@test.com, got %s", result.Recipient)
		return result
	}

	result.Success = true
	return result
}

func testDeliveryConfirmation(ctx context.Context, notifier *email.Notifier, mailhogURL string, shipmentID int64) TestResult {
	result := TestResult{NotificationType: "Delivery Confirmation"}

	err := notifier.SendDeliveryConfirmation(ctx, shipmentID)
	if err != nil {
		result.Error = err.Error()
		return result
	}

	time.Sleep(200 * time.Millisecond)

	msg, err := getLatestMailhogMessage(mailhogURL)
	if err != nil {
		result.Error = fmt.Sprintf("Email not found in Mailhog: %v", err)
		return result
	}

	result.EmailID = msg.ID
	result.TimeSent = msg.Created
	result.Recipient = fmt.Sprintf("%s@%s", msg.To[0].Mailbox, msg.To[0].Domain)

	if headers := msg.Content.Headers["Subject"]; len(headers) > 0 {
		result.Subject = headers[0]
	}

	if !strings.Contains(result.Subject, "Device Delivered Successfully") {
		result.Error = fmt.Sprintf("Unexpected subject: %s", result.Subject)
		return result
	}

	if result.Recipient != "engineer@test.com" {
		result.Error = fmt.Sprintf("Wrong recipient: expected engineer@test.com, got %s", result.Recipient)
		return result
	}

	result.Success = true
	return result
}

func printTestResult(result TestResult) {
	if result.Success {
		fmt.Printf("  %s✅ SUCCESS%s\n", colorGreen, colorReset)
		fmt.Printf("     Subject:   %s\n", result.Subject)
		fmt.Printf("     Recipient: %s\n", result.Recipient)
		fmt.Printf("     Email ID:  %s\n", result.EmailID)
	} else {
		fmt.Printf("  %s❌ FAILED%s\n", colorRed, colorReset)
		fmt.Printf("     Error: %s\n", result.Error)
	}
	fmt.Println()
}

func printSummary(results []TestResult) {
	fmt.Println(colorPurple + "===========================================" + colorReset)
	fmt.Println(colorPurple + "  Test Summary" + colorReset)
	fmt.Println(colorPurple + "===========================================" + colorReset)
	fmt.Println()

	passed := 0
	failed := 0

	for _, r := range results {
		if r.Success {
			passed++
			fmt.Printf("%s✅%s %s\n", colorGreen, colorReset, r.NotificationType)
		} else {
			failed++
			fmt.Printf("%s❌%s %s - %s\n", colorRed, colorReset, r.NotificationType, r.Error)
		}
	}

	fmt.Println()
	fmt.Printf("Total Tests: %d\n", len(results))
	fmt.Printf("%sPassed: %d%s\n", colorGreen, passed, colorReset)
	if failed > 0 {
		fmt.Printf("%sFailed: %d%s\n", colorRed, failed, colorReset)
	} else {
		fmt.Printf("Failed: %d\n", failed)
	}
	fmt.Println()
}
