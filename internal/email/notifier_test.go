package email

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/yourusername/laptop-tracking-system/internal/database"
	"github.com/yourusername/laptop-tracking-system/internal/models"
)

func TestNewNotifier(t *testing.T) {
	client, err := NewClient(Config{
		Host: "smtp.example.com",
		Port: 587,
		From: "noreply@example.com",
	})
	if err != nil {
		t.Fatalf("Failed to create email client: %v", err)
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	notifier := NewNotifier(client, db)

	if notifier == nil {
		t.Fatal("NewNotifier() returned nil")
	}

	if notifier.client == nil {
		t.Error("Notifier client is nil")
	}

	if notifier.templates == nil {
		t.Error("Notifier templates is nil")
	}

	if notifier.db == nil {
		t.Error("Notifier db is nil")
	}
}

func TestNotifier_logNotification(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client, err := NewClient(Config{
		Host: "smtp.example.com",
		Port: 587,
		From: "noreply@example.com",
	})
	if err != nil {
		t.Fatalf("Failed to create email client: %v", err)
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	notifier := NewNotifier(client, db)

	// Create test user and shipment
	user := &models.User{
		Email:        "test@example.com",
		PasswordHash: "hash",
		Role:         models.RoleLogistics,
	}
	user.BeforeCreate()

	err = db.QueryRowContext(
		context.Background(),
		`INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		user.Email, user.PasswordHash, user.Role, user.CreatedAt, user.UpdatedAt,
	).Scan(&user.ID)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Create test client company
	company := &models.ClientCompany{
		Name:        "Test Company",
		ContactInfo: "contact@test.com",
	}
	company.BeforeCreate()

	err = db.QueryRowContext(
		context.Background(),
		`INSERT INTO client_companies (name, contact_info, created_at)
		VALUES ($1, $2, $3) RETURNING id`,
		company.Name, company.ContactInfo, company.CreatedAt,
	).Scan(&company.ID)
	if err != nil {
		t.Fatalf("Failed to create test company: %v", err)
	}

	// Create test shipment
	shipment := &models.Shipment{
		ClientCompanyID: company.ID,
		Status:          models.ShipmentStatusPendingPickup,
	}
	shipment.BeforeCreate()

	err = db.QueryRowContext(
		context.Background(),
		`INSERT INTO shipments (client_company_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4) RETURNING id`,
		shipment.ClientCompanyID, shipment.Status, shipment.CreatedAt, shipment.UpdatedAt,
	).Scan(&shipment.ID)
	if err != nil {
		t.Fatalf("Failed to create test shipment: %v", err)
	}

	tests := []struct {
		name             string
		shipmentID       int64
		notificationType string
		recipient        string
		status           string
		wantErr          bool
	}{
		{
			name:             "log notification with shipment",
			shipmentID:       shipment.ID,
			notificationType: "pickup_confirmation",
			recipient:        "test@example.com",
			status:           "sent",
			wantErr:          false,
		},
		{
			name:             "log notification without shipment",
			shipmentID:       0,
			notificationType: "magic_link",
			recipient:        "user@example.com",
			status:           "sent",
			wantErr:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := notifier.logNotification(
				context.Background(),
				tt.shipmentID,
				tt.notificationType,
				tt.recipient,
				tt.status,
			)

			if (err != nil) != tt.wantErr {
				t.Errorf("logNotification() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify the log was created
				var count int
				err = db.QueryRowContext(
					context.Background(),
					`SELECT COUNT(*) FROM notification_logs 
					WHERE type = $1 AND recipient = $2 AND status = $3`,
					tt.notificationType, tt.recipient, tt.status,
				).Scan(&count)

				if err != nil {
					t.Fatalf("Failed to query notification log: %v", err)
				}

				if count == 0 {
					t.Error("Notification was not logged in the database")
				}
			}
		})
	}
}

func TestNotifier_getShipmentDetails(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client, err := NewClient(Config{
		Host: "smtp.example.com",
		Port: 587,
		From: "noreply@example.com",
	})
	if err != nil {
		t.Fatalf("Failed to create email client: %v", err)
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	notifier := NewNotifier(client, db)

	// Create test client company with unique name to avoid duplicate key errors
	company := &models.ClientCompany{
		Name:        fmt.Sprintf("Test Company %d", time.Now().UnixNano()),
		ContactInfo: "contact@test.com",
	}
	company.BeforeCreate()

	err = db.QueryRowContext(
		context.Background(),
		`INSERT INTO client_companies (name, contact_info, created_at)
		VALUES ($1, $2, $3) RETURNING id`,
		company.Name, company.ContactInfo, company.CreatedAt,
	).Scan(&company.ID)
	if err != nil {
		t.Fatalf("Failed to create test company: %v", err)
	}

	// Create test shipment with details
	pickupScheduledDate := time.Now().AddDate(0, 0, 1)

	shipment := &models.Shipment{
		ClientCompanyID:     company.ID,
		Status:              models.ShipmentStatusPendingPickup,
		TrackingNumber:      "UPS123456789",
		PickupScheduledDate: &pickupScheduledDate,
	}
	shipment.BeforeCreate()

	err = db.QueryRowContext(
		context.Background(),
		`INSERT INTO shipments (client_company_id, status, tracking_number, 
		pickup_scheduled_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		shipment.ClientCompanyID, shipment.Status, shipment.TrackingNumber,
		shipment.PickupScheduledDate, shipment.CreatedAt, shipment.UpdatedAt,
	).Scan(&shipment.ID)
	if err != nil {
		t.Fatalf("Failed to create test shipment: %v", err)
	}

	// Test fetching shipment details
	details, err := notifier.getShipmentDetails(context.Background(), shipment.ID)
	if err != nil {
		t.Fatalf("getShipmentDetails() error = %v", err)
	}

	if details.ClientCompanyID != company.ID {
		t.Errorf("ClientCompanyID = %v, want %v", details.ClientCompanyID, company.ID)
	}

	if details.TrackingNumber.String != "UPS123456789" {
		t.Errorf("TrackingNumber = %v, want UPS123456789", details.TrackingNumber.String)
	}

	if !details.PickupScheduledDate.Valid {
		t.Error("PickupScheduledDate should be valid")
	}
}

func TestNotifier_SendMagicLink(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Start a mock SMTP server
	mockServer, err := newMockSMTPServer(2526)
	if err != nil {
		t.Fatalf("Failed to start mock SMTP server: %v", err)
	}
	defer mockServer.close()

	// Give the server time to start
	time.Sleep(100 * time.Millisecond)

	client, err := NewClient(Config{
		Host: "127.0.0.1",
		Port: 2526,
		From: "noreply@example.com",
	})
	if err != nil {
		t.Fatalf("Failed to create email client: %v", err)
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	notifier := NewNotifier(client, db)

	// Test sending magic link
	err = notifier.SendMagicLink(
		context.Background(),
		"recipient@example.com",
		"John Doe",
		"https://example.com/form?token=abc123",
		"pickup",
		time.Now().Add(24*time.Hour),
	)

	if err != nil {
		t.Errorf("SendMagicLink() error = %v", err)
	}

	// Verify notification was logged
	var count int
	err = db.QueryRowContext(
		context.Background(),
		`SELECT COUNT(*) FROM notification_logs WHERE type = 'magic_link'`,
	).Scan(&count)

	if err != nil {
		t.Fatalf("Failed to query notification log: %v", err)
	}

	if count == 0 {
		t.Error("Magic link notification was not logged")
	}
}

func TestNotifier_generatePlainTextFromHTML(t *testing.T) {
	client, err := NewClient(Config{
		Host: "smtp.example.com",
		Port: 587,
		From: "noreply@example.com",
	})
	if err != nil {
		t.Fatalf("Failed to create email client: %v", err)
	}

	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	notifier := NewNotifier(client, db)

	html := "<html><body><h1>Test</h1><p>This is a test</p></body></html>"
	plainText := notifier.generatePlainTextFromHTML(html)

	if plainText == "" {
		t.Error("generatePlainTextFromHTML() returned empty string")
	}

	// For now, we just check it returns something
	// In future, could implement proper HTML to text conversion
}

