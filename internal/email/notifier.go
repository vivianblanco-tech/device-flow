package email

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
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
	var clientName, clientCompany string
	err = n.db.QueryRowContext(ctx,
		`SELECT name, contact_info FROM client_companies WHERE id = $1`,
		shipment.ClientCompanyID,
	).Scan(&clientName, &clientCompany)
	if err != nil {
		return fmt.Errorf("failed to fetch client company: %w", err)
	}

	// Fetch pickup form data to get contact email (FIX: use form contact_email instead of users table)
	var formDataJSON string
	err = n.db.QueryRowContext(ctx,
		`SELECT form_data FROM pickup_forms WHERE shipment_id = $1 ORDER BY submitted_at DESC LIMIT 1`,
		shipmentID,
	).Scan(&formDataJSON)
	if err != nil {
		// If no pickup form exists, we can't send the notification (no contact email)
		return fmt.Errorf("no pickup form found for shipment %d: cannot send notification without contact email", shipmentID)
	}

	// Parse form data to extract contact information
	var formData map[string]interface{}
	if err := json.Unmarshal([]byte(formDataJSON), &formData); err != nil {
		return fmt.Errorf("failed to parse form data: %w", err)
	}

	clientEmail, ok := formData["contact_email"].(string)
	if !ok || clientEmail == "" {
		return fmt.Errorf("contact email not found in pickup form")
	}

	// Extract contact name from form if available
	contactName, _ := formData["contact_name"].(string)
	if contactName == "" {
		contactName = clientName // Fallback to company name
	}

	// Prepare template data
	pickupDate := "To be scheduled"
	if shipment.PickupScheduledDate.Valid {
		pickupDate = shipment.PickupScheduledDate.Time.Format("Monday, January 2, 2006")
	}

	// Extract pickup time slot from form if available
	pickupTimeSlot := "Morning (9AM - 12PM)" // Default time slot
	if timeSlot, ok := formData["pickup_time_slot"].(string); ok && timeSlot != "" {
		switch timeSlot {
		case "morning":
			pickupTimeSlot = "Morning (8AM - 12PM)"
		case "afternoon":
			pickupTimeSlot = "Afternoon (12PM - 5PM)"
		case "evening":
			pickupTimeSlot = "Evening (5PM - 8PM)"
		}
	}

	data := PickupConfirmationData{
		ClientName:       contactName, // Use contact name from form
		ClientCompany:    clientCompany,
		TrackingNumber:   shipment.TrackingNumber.String,
		PickupDate:       pickupDate,
		PickupTimeSlot:   pickupTimeSlot,
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

// SendPickupScheduledNotification sends notification to contact email when pickup is scheduled
func (n *Notifier) SendPickupScheduledNotification(ctx context.Context, shipmentID int64) error {
	// Fetch shipment details
	shipment, err := n.getShipmentDetails(ctx, shipmentID)
	if err != nil {
		return fmt.Errorf("failed to fetch shipment details: %w", err)
	}

	// Fetch client company
	var clientCompany string
	err = n.db.QueryRowContext(ctx,
		`SELECT name FROM client_companies WHERE id = $1`,
		shipment.ClientCompanyID,
	).Scan(&clientCompany)
	if err != nil {
		return fmt.Errorf("failed to fetch client company: %w", err)
	}

	// Fetch pickup form data to get contact email
	var formDataJSON string
	err = n.db.QueryRowContext(ctx,
		`SELECT form_data FROM pickup_forms WHERE shipment_id = $1 ORDER BY submitted_at DESC LIMIT 1`,
		shipmentID,
	).Scan(&formDataJSON)
	if err != nil {
		// If no pickup form exists, we can't send the notification (no contact email)
		// This is not an error - it just means the form hasn't been submitted yet
		return fmt.Errorf("no pickup form found for shipment %d: cannot send notification without contact email", shipmentID)
	}

	// Parse form data to extract contact information
	var formData map[string]interface{}
	if err := json.Unmarshal([]byte(formDataJSON), &formData); err != nil {
		return fmt.Errorf("failed to parse form data: %w", err)
	}

	contactEmail, ok := formData["contact_email"].(string)
	if !ok || contactEmail == "" {
		return fmt.Errorf("contact email not found in pickup form")
	}

	contactName, _ := formData["contact_name"].(string)
	if contactName == "" {
		contactName = "Client Contact"
	}

	// Prepare template data
	// Priority: Use pickup_date from form data first, then fall back to shipment.PickupScheduledDate
	pickupDate := "To be determined"
	if formPickupDate, ok := formData["pickup_date"].(string); ok && formPickupDate != "" {
		// Parse the date from form (format: "2006-01-02")
		if parsedDate, err := time.Parse("2006-01-02", formPickupDate); err == nil {
			pickupDate = parsedDate.Format("Monday, January 2, 2006")
		} else if shipment.PickupScheduledDate.Valid {
			// Fallback to shipment date if form date parsing fails
			pickupDate = shipment.PickupScheduledDate.Time.Format("Monday, January 2, 2006")
		}
	} else if shipment.PickupScheduledDate.Valid {
		// Fallback to shipment date if form date not available
		pickupDate = shipment.PickupScheduledDate.Time.Format("Monday, January 2, 2006")
	}

	pickupTimeSlot, _ := formData["pickup_time_slot"].(string)
	if pickupTimeSlot == "" {
		pickupTimeSlot = "To be confirmed"
	} else {
		// Format time slot
		switch pickupTimeSlot {
		case "morning":
			pickupTimeSlot = "Morning (8AM - 12PM)"
		case "afternoon":
			pickupTimeSlot = "Afternoon (12PM - 5PM)"
		case "evening":
			pickupTimeSlot = "Evening (5PM - 8PM)"
		}
	}

	pickupAddress := ""
	if addr, ok := formData["pickup_address"].(string); ok {
		pickupAddress = addr
		if city, ok := formData["pickup_city"].(string); ok {
			pickupAddress += ", " + city
		}
		if state, ok := formData["pickup_state"].(string); ok {
			pickupAddress += ", " + state
		}
		if zip, ok := formData["pickup_zip"].(string); ok {
			pickupAddress += " " + zip
		}
	}

	data := PickupScheduledData{
		ContactName:    contactName,
		ClientCompany:  clientCompany,
		TrackingNumber: shipment.TrackingNumber.String,
		PickupDate:     pickupDate,
		PickupTimeSlot: pickupTimeSlot,
		PickupAddress:  pickupAddress,
		ShipmentID:     shipmentID,
	}

	// Render template
	htmlBody, err := n.templates.RenderTemplate("pickup_scheduled", data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	// Send email
	message := Message{
		To:       []string{contactEmail},
		Subject:  n.templates.GetSubject("pickup_scheduled", data),
		Body:     n.generatePlainTextFromHTML(htmlBody),
		HTMLBody: htmlBody,
	}

	if err := n.client.Send(message); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	// Log notification
	if err := n.logNotification(ctx, shipmentID, "pickup_scheduled", contactEmail, "sent"); err != nil {
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
	var shipperName string
	var contactInfoJSON string
	err = n.db.QueryRowContext(ctx,
		`SELECT name, contact_info FROM client_companies WHERE id = $1`,
		shipment.ClientCompanyID,
	).Scan(&shipperName, &contactInfoJSON)
	if err != nil {
		return fmt.Errorf("failed to fetch client company: %w", err)
	}

	// Parse and format contact_info JSON to readable format
	shipperCompany := shipperName // Default to company name
	if contactInfoJSON != "" {
		var contactInfo map[string]interface{}
		if err := json.Unmarshal([]byte(contactInfoJSON), &contactInfo); err == nil {
			// Build readable format: email, phone, address, etc.
			var parts []string
			if email, ok := contactInfo["email"].(string); ok && email != "" {
				parts = append(parts, fmt.Sprintf("Email: %s", email))
			}
			if phone, ok := contactInfo["phone"].(string); ok && phone != "" {
				parts = append(parts, fmt.Sprintf("Phone: %s", phone))
			}
			if address, ok := contactInfo["address"].(string); ok && address != "" {
				parts = append(parts, fmt.Sprintf("Address: %s", address))
			}
			if len(parts) > 0 {
				shipperCompany = strings.Join(parts, " | ")
			}
		}
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

	// Get warehouse user email for contact details
	// Note: users table doesn't have a name field, so we'll use email
	var warehouseUserEmail string
	err = n.db.QueryRowContext(ctx,
		`SELECT email FROM users WHERE role = 'warehouse' LIMIT 1`,
	).Scan(&warehouseUserEmail)
	if err != nil {
		warehouseUserEmail = "warehouse@bairesdev.com" // Default fallback
	}

	// Get shipment release date for pickup date
	pickupDate := time.Now().AddDate(0, 0, 1).Format("Monday, January 2, 2006") // Default: tomorrow
	var releasedAt sql.NullTime
	err = n.db.QueryRowContext(ctx,
		`SELECT released_warehouse_at FROM shipments WHERE id = $1`,
		shipmentID,
	).Scan(&releasedAt)
	if err == nil && releasedAt.Valid {
		// Use release date + 1 day as pickup date (hardware ready for pickup next day)
		pickupDate = releasedAt.Time.AddDate(0, 0, 1).Format("Monday, January 2, 2006")
	}

	// Prepare template data
	data := ReleaseNotificationData{
		CourierName:        "Logistics Team",
		CourierCompany:     "BairesDev",
		PickupDate:         pickupDate,
		PickupTimeSlot:     "Morning (9AM - 12PM)", // Default time slot
		WarehouseAddress:   "Please contact warehouse for pickup address", // TODO: Add to config
		ContactPerson:      "Warehouse Team", // Use generic name since users table doesn't have name field
		ContactPhone:       "Email: " + warehouseUserEmail, // Use email since phone not available in users table
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

