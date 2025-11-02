package email

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// Notifier handles sending email notifications for various events
type Notifier struct {
	client    *Client
	templates *EmailTemplates
	db        *sql.DB
}

// NewNotifier creates a new email notifier instance
func NewNotifier(client *Client, db *sql.DB) *Notifier {
	return &Notifier{
		client:    client,
		templates: NewEmailTemplates(),
		db:        db,
	}
}

// SendPickupConfirmation sends a pickup confirmation email to the client
func (n *Notifier) SendPickupConfirmation(ctx context.Context, shipmentID int64) error {
	// Fetch shipment details
	shipment, err := n.getShipmentDetails(ctx, shipmentID)
	if err != nil {
		return fmt.Errorf("failed to fetch shipment details: %w", err)
	}

	// Fetch client company
	var clientName, clientCompany, clientEmail string
	err = n.db.QueryRowContext(ctx,
		`SELECT name, contact_info FROM client_companies WHERE id = $1`,
		shipment.ClientCompanyID,
	).Scan(&clientName, &clientCompany)
	if err != nil {
		return fmt.Errorf("failed to fetch client company: %w", err)
	}

	// Get client user email
	err = n.db.QueryRowContext(ctx,
		`SELECT email FROM users WHERE id IN (
			SELECT id FROM users WHERE role = 'client' LIMIT 1
		)`,
	).Scan(&clientEmail)
	if err != nil {
		clientEmail = "noreply@example.com" // Fallback
	}

	// Prepare template data
	pickupDate := "To be scheduled"
	if shipment.PickupScheduledDate.Valid {
		pickupDate = shipment.PickupScheduledDate.Time.Format("Monday, January 2, 2006")
	}

	data := PickupConfirmationData{
		ClientName:       clientName,
		ClientCompany:    clientCompany,
		TrackingNumber:   shipment.TrackingNumber.String,
		PickupDate:       pickupDate,
		PickupTimeSlot:   "Morning (9AM - 12PM)", // Default time slot
		NumberOfDevices:  shipment.LaptopCount,
		ConfirmationCode: fmt.Sprintf("CONF-%d", shipmentID),
	}

	// Render template
	htmlBody, err := n.templates.RenderTemplate("pickup_confirmation", data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	// Send email
	message := Message{
		To:       []string{clientEmail},
		Subject:  n.templates.GetSubject("pickup_confirmation", data),
		Body:     n.generatePlainTextFromHTML(htmlBody),
		HTMLBody: htmlBody,
	}

	if err := n.client.Send(message); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	// Log notification
	if err := n.logNotification(ctx, shipmentID, "pickup_confirmation", clientEmail, "sent"); err != nil {
		// Log error but don't fail
		fmt.Printf("Warning: failed to log notification: %v\n", err)
	}

	return nil
}

// SendWarehousePreAlert sends a pre-alert email to warehouse about incoming shipment
func (n *Notifier) SendWarehousePreAlert(ctx context.Context, shipmentID int64) error {
	// Fetch shipment details
	shipment, err := n.getShipmentDetails(ctx, shipmentID)
	if err != nil {
		return fmt.Errorf("failed to fetch shipment details: %w", err)
	}

	// Fetch client company
	var shipperName, shipperCompany string
	err = n.db.QueryRowContext(ctx,
		`SELECT name, contact_info FROM client_companies WHERE id = $1`,
		shipment.ClientCompanyID,
	).Scan(&shipperName, &shipperCompany)
	if err != nil {
		return fmt.Errorf("failed to fetch client company: %w", err)
	}

	// Get warehouse email
	var warehouseEmail string
	err = n.db.QueryRowContext(ctx,
		`SELECT email FROM users WHERE role = 'warehouse' LIMIT 1`,
	).Scan(&warehouseEmail)
	if err != nil {
		return fmt.Errorf("no warehouse user found: %w", err)
	}

	// Prepare template data
	expectedDate := "To be determined"
	if shipment.PickupScheduledDate.Valid {
		// Estimate arrival as 3 days after pickup
		expectedDate = shipment.PickupScheduledDate.Time.AddDate(0, 0, 3).Format("Monday, January 2, 2006")
	}

	data := WarehousePreAlertData{
		TrackingNumber:    shipment.TrackingNumber.String,
		ExpectedDate:      expectedDate,
		ShipperName:       shipperName,
		ShipperCompany:    shipperCompany,
		DeviceDescription: fmt.Sprintf("%d device(s)", shipment.LaptopCount),
		ProjectName:       "", // Can be added if needed
		TrackingURL:       fmt.Sprintf("https://www.ups.com/track?tracknum=%s", shipment.TrackingNumber.String),
	}

	// Render template
	htmlBody, err := n.templates.RenderTemplate("warehouse_pre_alert", data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	// Send email
	message := Message{
		To:       []string{warehouseEmail},
		Subject:  n.templates.GetSubject("warehouse_pre_alert", data),
		Body:     n.generatePlainTextFromHTML(htmlBody),
		HTMLBody: htmlBody,
	}

	if err := n.client.Send(message); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	// Log notification
	if err := n.logNotification(ctx, shipmentID, "warehouse_pre_alert", warehouseEmail, "sent"); err != nil {
		fmt.Printf("Warning: failed to log notification: %v\n", err)
	}

	return nil
}

// SendReleaseNotification sends notification when hardware is released from warehouse
func (n *Notifier) SendReleaseNotification(ctx context.Context, shipmentID int64) error {
	// Fetch shipment details
	shipment, err := n.getShipmentDetails(ctx, shipmentID)
	if err != nil {
		return fmt.Errorf("failed to fetch shipment details: %w", err)
	}

	// Get courier/logistics email
	var courierEmail string
	err = n.db.QueryRowContext(ctx,
		`SELECT email FROM users WHERE role = 'logistics' LIMIT 1`,
	).Scan(&courierEmail)
	if err != nil {
		return fmt.Errorf("no logistics user found: %w", err)
	}

	// Get engineer name if assigned
	var engineerName string
	if shipment.SoftwareEngineerID.Valid {
		err = n.db.QueryRowContext(ctx,
			`SELECT name FROM software_engineers WHERE id = $1`,
			shipment.SoftwareEngineerID.Int64,
		).Scan(&engineerName)
		if err != nil {
			engineerName = "Engineer"
		}
	} else {
		engineerName = "To be assigned"
	}

	// Prepare template data
	data := ReleaseNotificationData{
		CourierName:        "Logistics Team",
		CourierCompany:     "BairesDev",
		PickupDate:         time.Now().AddDate(0, 0, 1).Format("Monday, January 2, 2006"),
		PickupTimeSlot:     "Morning (9AM - 12PM)",
		WarehouseAddress:   "Warehouse Address", // Should be from config
		ContactPerson:      "Warehouse Manager",
		ContactPhone:       "Contact Phone",     // Should be from config
		DeviceSerialNumber: "SN-" + shipment.TrackingNumber.String,
		EngineerName:       engineerName,
		TrackingNumber:     shipment.TrackingNumber.String,
	}

	// Render template
	htmlBody, err := n.templates.RenderTemplate("release_notification", data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	// Send email
	message := Message{
		To:       []string{courierEmail},
		Subject:  n.templates.GetSubject("release_notification", data),
		Body:     n.generatePlainTextFromHTML(htmlBody),
		HTMLBody: htmlBody,
	}

	if err := n.client.Send(message); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	// Log notification
	if err := n.logNotification(ctx, shipmentID, "release_notification", courierEmail, "sent"); err != nil {
		fmt.Printf("Warning: failed to log notification: %v\n", err)
	}

	return nil
}

// SendDeliveryConfirmation sends confirmation when device is delivered to engineer
func (n *Notifier) SendDeliveryConfirmation(ctx context.Context, shipmentID int64) error {
	// Fetch shipment details
	shipment, err := n.getShipmentDetails(ctx, shipmentID)
	if err != nil {
		return fmt.Errorf("failed to fetch shipment details: %w", err)
	}

	// Get engineer details
	if !shipment.SoftwareEngineerID.Valid {
		return fmt.Errorf("shipment has no assigned engineer")
	}

	var engineerName, engineerEmail string
	err = n.db.QueryRowContext(ctx,
		`SELECT name, email FROM software_engineers WHERE id = $1`,
		shipment.SoftwareEngineerID.Int64,
	).Scan(&engineerName, &engineerEmail)
	if err != nil {
		return fmt.Errorf("failed to fetch engineer details: %w", err)
	}

	// Prepare template data
	data := DeliveryConfirmationData{
		EngineerName:       engineerName,
		DeviceSerialNumber: "SN-" + shipment.TrackingNumber.String,
		DeviceModel:        "Configured Laptop",
		DeliveryDate:       time.Now().Format("Monday, January 2, 2006"),
		TrackingNumber:     shipment.TrackingNumber.String,
		ProjectName:        "",
	}

	// Render template
	htmlBody, err := n.templates.RenderTemplate("delivery_confirmation", data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	// Send email
	message := Message{
		To:       []string{engineerEmail},
		Subject:  n.templates.GetSubject("delivery_confirmation", data),
		Body:     n.generatePlainTextFromHTML(htmlBody),
		HTMLBody: htmlBody,
	}

	if err := n.client.Send(message); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	// Log notification
	if err := n.logNotification(ctx, shipmentID, "delivery_confirmation", engineerEmail, "sent"); err != nil {
		fmt.Printf("Warning: failed to log notification: %v\n", err)
	}

	return nil
}

// SendMagicLink sends a magic link email for form access
func (n *Notifier) SendMagicLink(ctx context.Context, recipientEmail, recipientName, magicLink, formType string, expiresAt time.Time) error {
	// Prepare template data
	data := MagicLinkData{
		RecipientName: recipientName,
		MagicLink:     magicLink,
		ExpiresAt:     expiresAt,
		FormType:      formType,
	}

	// Render template
	htmlBody, err := n.templates.RenderTemplate("magic_link", data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	// Send email
	message := Message{
		To:       []string{recipientEmail},
		Subject:  n.templates.GetSubject("magic_link", data),
		Body:     n.generatePlainTextFromHTML(htmlBody),
		HTMLBody: htmlBody,
	}

	if err := n.client.Send(message); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	// Log notification (no shipment ID for magic links)
	if err := n.logNotification(ctx, 0, "magic_link", recipientEmail, "sent"); err != nil {
		fmt.Printf("Warning: failed to log notification: %v\n", err)
	}

	return nil
}

// Helper methods

type shipmentDetails struct {
	ClientCompanyID      int64
	SoftwareEngineerID   sql.NullInt64
	TrackingNumber       sql.NullString
	PickupScheduledDate  sql.NullTime
	ArrivedWarehouseAt   sql.NullTime
	LaptopCount          int
}

func (n *Notifier) getShipmentDetails(ctx context.Context, shipmentID int64) (*shipmentDetails, error) {
	var details shipmentDetails

	// Get basic shipment info
	err := n.db.QueryRowContext(ctx,
		`SELECT client_company_id, software_engineer_id, tracking_number, 
		pickup_scheduled_date, arrived_warehouse_at
		FROM shipments WHERE id = $1`,
		shipmentID,
	).Scan(
		&details.ClientCompanyID,
		&details.SoftwareEngineerID,
		&details.TrackingNumber,
		&details.PickupScheduledDate,
		&details.ArrivedWarehouseAt,
	)
	if err != nil {
		return nil, err
	}

	// Count laptops in this shipment
	err = n.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM shipment_laptops WHERE shipment_id = $1`,
		shipmentID,
	).Scan(&details.LaptopCount)
	if err != nil {
		details.LaptopCount = 0
	}

	return &details, nil
}

func (n *Notifier) logNotification(ctx context.Context, shipmentID int64, notificationType, recipient, status string) error {
	var shipmentIDPtr *int64
	if shipmentID > 0 {
		shipmentIDPtr = &shipmentID
	}

	_, err := n.db.ExecContext(ctx,
		`INSERT INTO notification_logs (shipment_id, type, recipient, sent_at, status)
		VALUES ($1, $2, $3, $4, $5)`,
		shipmentIDPtr,
		notificationType,
		recipient,
		time.Now(),
		status,
	)
	return err
}

// generatePlainTextFromHTML creates a simple plain text version from HTML
// This is a basic implementation - could be improved with proper HTML parsing
func (n *Notifier) generatePlainTextFromHTML(html string) string {
	// For now, just strip HTML tags - could be improved
	// In production, consider using a library like bluemonday or goquery
	return "Please view this email in an HTML-capable email client."
}

