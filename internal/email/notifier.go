package email

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/yourusername/laptop-tracking-system/internal/config"
	"github.com/yourusername/laptop-tracking-system/internal/models"
)

// Notifier handles sending email notifications for various events
type Notifier struct {
	client    *Client
	templates *EmailTemplates
	db        *sql.DB
	config    *config.SMTPConfig // Optional config for default emails
}

// NewNotifier creates a new email notifier instance
func NewNotifier(client *Client, db *sql.DB) *Notifier {
	return &Notifier{
		client:    client,
		templates: NewEmailTemplates(),
		db:        db,
		config:    nil, // Config is optional for backward compatibility
	}
}

// NewNotifierWithConfig creates a new email notifier instance with config
func NewNotifierWithConfig(client *Client, db *sql.DB, cfg *config.SMTPConfig) *Notifier {
	return &Notifier{
		client:    client,
		templates: NewEmailTemplates(),
		db:        db,
		config:    cfg,
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
	// Fetch shipment details including shipment type
	var shipmentType string
	shipment, err := n.getShipmentDetails(ctx, shipmentID)
	if err != nil {
		return fmt.Errorf("failed to fetch shipment details: %w", err)
	}
	
	// Get shipment type
	err = n.db.QueryRowContext(ctx,
		`SELECT shipment_type FROM shipments WHERE id = $1`,
		shipmentID,
	).Scan(&shipmentType)
	if err != nil {
		return fmt.Errorf("failed to fetch shipment type: %w", err)
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
	warehouseEmail, err := n.getWarehouseEmail(ctx)
	if err != nil {
		// Log warning but continue with default email
		fmt.Printf("Warning: %v\n", err)
	}

	// Get pickup date from form data (preferred) or shipment
	var pickupDate time.Time
	pickupDateFound := false
	
	// Try to get pickup date from pickup form first
	var formDataJSON string
	err = n.db.QueryRowContext(ctx,
		`SELECT form_data FROM pickup_forms WHERE shipment_id = $1 ORDER BY submitted_at DESC LIMIT 1`,
		shipmentID,
	).Scan(&formDataJSON)
	if err == nil {
		var formData map[string]interface{}
		if err := json.Unmarshal([]byte(formDataJSON), &formData); err == nil {
			if formPickupDate, ok := formData["pickup_date"].(string); ok && formPickupDate != "" {
				// Parse the date from form (format: "2006-01-02")
				if parsedDate, err := time.Parse("2006-01-02", formPickupDate); err == nil {
					pickupDate = parsedDate
					pickupDateFound = true
				}
			}
		}
	}
	
	// Fallback to shipment.PickupScheduledDate if form date not available
	if !pickupDateFound && shipment.PickupScheduledDate.Valid {
		pickupDate = shipment.PickupScheduledDate.Time
		pickupDateFound = true
	}

	// Prepare template data
	expectedDate := "To be determined"
	if pickupDateFound {
		// Expected delivery is pickup date + 1 day
		expectedDate = pickupDate.AddDate(0, 0, 1).Format("Monday, January 2, 2006")
	}

	// Determine if single or bulk shipment
	isSingleShipment := shipmentType == string(models.ShipmentTypeSingleFullJourney) || 
		shipmentType == string(models.ShipmentTypeWarehouseToEngineer)
	isBulkShipment := shipmentType == string(models.ShipmentTypeBulkToWarehouse)

	data := WarehousePreAlertData{
		TrackingNumber:    shipment.TrackingNumber.String,
		ExpectedDate:      expectedDate,
		ShipperName:       shipperName,
		ShipperCompany:    shipperCompany,
		DeviceDescription: fmt.Sprintf("%d device(s)", shipment.LaptopCount),
		ProjectName:       "", // Can be added if needed
		TrackingURL:       fmt.Sprintf("https://www.ups.com/track?tracknum=%s", shipment.TrackingNumber.String),
		IsSingleShipment:  isSingleShipment,
		IsBulkShipment:    isBulkShipment,
		LaptopCount:       shipment.LaptopCount,
		NumberOfBoxes:     0, // Will be set for bulk shipments below
	}

	// Fetch laptop details for single shipments
	if isSingleShipment {
		var serialNumber, brand, model, cpu, ramGB, ssdGB, sku sql.NullString
		err = n.db.QueryRowContext(ctx,
			`SELECT l.serial_number, l.brand, l.model, l.cpu, l.ram_gb, l.ssd_gb, l.sku
			FROM laptops l
			JOIN shipment_laptops sl ON sl.laptop_id = l.id
			WHERE sl.shipment_id = $1
			LIMIT 1`,
			shipmentID,
		).Scan(&serialNumber, &brand, &model, &cpu, &ramGB, &ssdGB, &sku)
		
		if err == nil {
			if serialNumber.Valid {
				data.SerialNumber = serialNumber.String
			}
			if brand.Valid {
				data.Brand = brand.String
			}
			if model.Valid {
				data.Model = model.String
			}
			if cpu.Valid {
				data.CPU = cpu.String
			}
			if ramGB.Valid {
				data.RAMGB = ramGB.String
			}
			if ssdGB.Valid {
				data.SSDGB = ssdGB.String
			}
			if sku.Valid {
				data.SKU = sku.String
			}
		}
	}

	// For bulk shipments, fetch actual count from pickup form data
	if isBulkShipment {
		// Try to get actual laptop count and number of boxes from pickup form
		var numberOfLaptops, numberOfBoxes int
		if formDataJSON != "" {
			var formData map[string]interface{}
			if err := json.Unmarshal([]byte(formDataJSON), &formData); err == nil {
				// Get number_of_laptops from form data
				if laptops, ok := formData["number_of_laptops"].(float64); ok {
					numberOfLaptops = int(laptops)
				} else if laptops, ok := formData["number_of_laptops"].(int); ok {
					numberOfLaptops = laptops
				}
				
				// Get number_of_boxes from form data
				if boxes, ok := formData["number_of_boxes"].(float64); ok {
					numberOfBoxes = int(boxes)
				} else if boxes, ok := formData["number_of_boxes"].(int); ok {
					numberOfBoxes = boxes
				}
			}
		}
		
		// Use form data if available, otherwise fall back to shipment.LaptopCount
		if numberOfLaptops > 0 {
			data.LaptopCount = numberOfLaptops
		} else if shipment.LaptopCount > 0 {
			data.LaptopCount = shipment.LaptopCount
		}
		
		// Set number of boxes
		if numberOfBoxes > 0 {
			data.NumberOfBoxes = numberOfBoxes
		}
		
		// Build bulk description with laptop count and boxes if available
		if numberOfBoxes > 0 && data.LaptopCount > 0 {
			data.BulkDescription = fmt.Sprintf("Bulk shipment containing %d device(s) in %d box(es)", data.LaptopCount, numberOfBoxes)
		} else if data.LaptopCount > 0 {
			data.BulkDescription = fmt.Sprintf("Bulk shipment containing %d device(s)", data.LaptopCount)
		} else {
			data.BulkDescription = "Bulk shipment"
		}
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

// SendShipmentPickedUpNotification sends a notification to the client when shipment is picked up
func (n *Notifier) SendShipmentPickedUpNotification(ctx context.Context, shipmentID int64) error {
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

	// Extract contact name from form if available
	contactName, _ := formData["contact_name"].(string)
	if contactName == "" {
		contactName = clientCompany // Fallback to company name
	}

	// Get picked up date and courier from shipment
	var pickedUpAt sql.NullTime
	var courierName sql.NullString
	var shipmentType string
	err = n.db.QueryRowContext(ctx,
		`SELECT picked_up_at, courier_name, shipment_type FROM shipments WHERE id = $1`,
		shipmentID,
	).Scan(&pickedUpAt, &courierName, &shipmentType)
	if err != nil {
		return fmt.Errorf("failed to fetch shipment pickup details: %w", err)
	}

	// Get picked up date
	pickedUpDate := time.Now().Format("Monday, January 2, 2006")
	if pickedUpAt.Valid {
		pickedUpDate = pickedUpAt.Time.Format("Monday, January 2, 2006")
	}

	// Calculate expected arrival (typically 3 days after pickup)
	expectedArrivalDate := time.Now().AddDate(0, 0, 3)
	if pickedUpAt.Valid {
		expectedArrivalDate = pickedUpAt.Time.AddDate(0, 0, 3)
	}
	expectedArrival := expectedArrivalDate.Format("Monday, January 2, 2006")

	// Build tracking URL based on courier
	trackingURL := ""
	courierNameStr := "Courier"
	if courierName.Valid && courierName.String != "" {
		courierNameStr = courierName.String
	}
	if shipment.TrackingNumber.String != "" {
		switch strings.ToUpper(courierNameStr) {
		case "UPS":
			trackingURL = fmt.Sprintf("https://www.ups.com/track?tracknum=%s", shipment.TrackingNumber.String)
		case "FEDEX", "FEDEX EXPRESS":
			trackingURL = fmt.Sprintf("https://www.fedex.com/fedextrack/?trknbr=%s", shipment.TrackingNumber.String)
		case "DHL":
			trackingURL = fmt.Sprintf("https://www.dhl.com/en/express/tracking.html?AWB=%s", shipment.TrackingNumber.String)
		default:
			trackingURL = fmt.Sprintf("https://www.google.com/search?q=track+%s", shipment.TrackingNumber.String)
		}
	}

	// Prepare template data
	data := ShipmentPickedUpData{
		ContactName:     contactName,
		ClientCompany:   clientCompany,
		TrackingNumber:  shipment.TrackingNumber.String,
		CourierName:     courierNameStr,
		PickedUpDate:    pickedUpDate,
		ExpectedArrival: expectedArrival,
		TrackingURL:     trackingURL,
		ShipmentType:    shipmentType,
	}

	// Render template
	htmlBody, err := n.templates.RenderTemplate("shipment_picked_up", data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	// Send email
	message := Message{
		To:       []string{contactEmail},
		Subject:  n.templates.GetSubject("shipment_picked_up", data),
		Body:     n.generatePlainTextFromHTML(htmlBody),
		HTMLBody: htmlBody,
	}

	if err := n.client.Send(message); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	// Log notification
	if err := n.logNotification(ctx, shipmentID, "shipment_picked_up", contactEmail, "sent"); err != nil {
		fmt.Printf("Warning: failed to log notification: %v\n", err)
	}

	return nil
}

// SendPickupFormSubmittedNotification sends notification to logistics when pickup form is submitted
func (n *Notifier) SendPickupFormSubmittedNotification(ctx context.Context, shipmentID int64) error {
	// Fetch shipment details
	shipment, err := n.getShipmentDetails(ctx, shipmentID)
	if err != nil {
		return fmt.Errorf("failed to fetch shipment details: %w", err)
	}

	// Get shipment type and other details
	var shipmentType string
	var jiraTicket string
	var laptopCount int
	err = n.db.QueryRowContext(ctx,
		`SELECT shipment_type, jira_ticket_number, laptop_count FROM shipments WHERE id = $1`,
		shipmentID,
	).Scan(&shipmentType, &jiraTicket, &laptopCount)
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

	// Fetch pickup form data
	var formDataJSON string
	err = n.db.QueryRowContext(ctx,
		`SELECT form_data FROM pickup_forms WHERE shipment_id = $1 ORDER BY submitted_at DESC LIMIT 1`,
		shipmentID,
	).Scan(&formDataJSON)
	if err != nil {
		return fmt.Errorf("no pickup form found for shipment %d", shipmentID)
	}

	// Parse form data
	var formData map[string]interface{}
	if err := json.Unmarshal([]byte(formDataJSON), &formData); err != nil {
		return fmt.Errorf("failed to parse form data: %w", err)
	}

	// Extract contact information
	contactName, _ := formData["contact_name"].(string)
	contactEmail, _ := formData["contact_email"].(string)
	contactPhone, _ := formData["contact_phone"].(string)

	// Build pickup address
	pickupAddress := ""
	if addr, ok := formData["pickup_address"].(string); ok && addr != "" {
		pickupAddress = addr
		if city, ok := formData["pickup_city"].(string); ok && city != "" {
			pickupAddress += ", " + city
		}
		if state, ok := formData["pickup_state"].(string); ok && state != "" {
			pickupAddress += ", " + state
		}
		if zip, ok := formData["pickup_zip"].(string); ok && zip != "" {
			pickupAddress += " " + zip
		}
	}

	// Get pickup date
	pickupDate := "To be scheduled"
	if formPickupDate, ok := formData["pickup_date"].(string); ok && formPickupDate != "" {
		if parsedDate, err := time.Parse("2006-01-02", formPickupDate); err == nil {
			pickupDate = parsedDate.Format("Monday, January 2, 2006")
		}
	}

	// Get logistics email
	var logisticsEmail string
	err = n.db.QueryRowContext(ctx,
		`SELECT email FROM users WHERE role = 'logistics' LIMIT 1`,
	).Scan(&logisticsEmail)
	if err != nil {
		// Fallback to default logistics email
		logisticsEmail = "international@bairesdev.com"
	}

	// Build shipment URL (assuming base URL from config or environment)
	shipmentURL := fmt.Sprintf("/shipments/%d", shipmentID)

	// Prepare template data
	data := PickupFormSubmittedData{
		ShipmentID:      shipmentID,
		ShipmentType:    shipmentType,
		ClientCompany:   clientCompany,
		ContactName:     contactName,
		ContactEmail:    contactEmail,
		ContactPhone:    contactPhone,
		PickupAddress:   pickupAddress,
		PickupDate:      pickupDate,
		NumberOfDevices: laptopCount,
		JiraTicket:      jiraTicket,
		ShipmentURL:     shipmentURL,
	}

	// Render template
	htmlBody, err := n.templates.RenderTemplate("pickup_form_submitted_logistics", data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	// Send email
	message := Message{
		To:       []string{logisticsEmail},
		Subject:  n.templates.GetSubject("pickup_form_submitted_logistics", data),
		Body:     n.generatePlainTextFromHTML(htmlBody),
		HTMLBody: htmlBody,
	}

	if err := n.client.Send(message); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	// Log notification
	if err := n.logNotification(ctx, shipmentID, "pickup_form_submitted_logistics", logisticsEmail, "sent"); err != nil {
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
	courierEmail, err := n.getLogisticsEmail(ctx)
	if err != nil {
		// Log warning but continue with default email
		fmt.Printf("Warning: %v\n", err)
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
	warehouseUserEmail, err := n.getWarehouseEmail(ctx)
	if err != nil {
		// Log warning but continue with default email
		fmt.Printf("Warning: %v\n", err)
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

	// Get actual delivery date from shipment
	deliveryDate := time.Now().Format("Monday, January 2, 2006") // Default to current date
	var deliveredAt sql.NullTime
	err = n.db.QueryRowContext(ctx,
		`SELECT delivered_at FROM shipments WHERE id = $1`,
		shipmentID,
	).Scan(&deliveredAt)
	if err == nil && deliveredAt.Valid {
		deliveryDate = deliveredAt.Time.Format("Monday, January 2, 2006")
	}

	// Get laptop model from the first laptop in the shipment
	deviceModel := "Configured Laptop" // Default fallback
	var laptopBrand, laptopModel string
	err = n.db.QueryRowContext(ctx,
		`SELECT l.brand, l.model 
		FROM laptops l
		JOIN shipment_laptops sl ON sl.laptop_id = l.id
		WHERE sl.shipment_id = $1
		LIMIT 1`,
		shipmentID,
	).Scan(&laptopBrand, &laptopModel)
	if err == nil {
		// Format as "Brand Model" if both available, otherwise just model
		if laptopBrand != "" && laptopModel != "" {
			deviceModel = fmt.Sprintf("%s %s", laptopBrand, laptopModel)
		} else if laptopModel != "" {
			deviceModel = laptopModel
		} else if laptopBrand != "" {
			deviceModel = laptopBrand
		}
	}

	// Prepare template data
	data := DeliveryConfirmationData{
		EngineerName:       engineerName,
		DeviceSerialNumber: "SN-" + shipment.TrackingNumber.String,
		DeviceModel:        deviceModel,
		DeliveryDate:       deliveryDate,
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

// SendEngineerDeliveryNotificationToClient sends notification to client when device is delivered to engineer
func (n *Notifier) SendEngineerDeliveryNotificationToClient(ctx context.Context, shipmentID int64) error {
	// Fetch shipment details
	shipment, err := n.getShipmentDetails(ctx, shipmentID)
	if err != nil {
		return fmt.Errorf("failed to fetch shipment details: %w", err)
	}

	// Get shipment type - only send for applicable types
	var shipmentType string
	err = n.db.QueryRowContext(ctx,
		`SELECT shipment_type FROM shipments WHERE id = $1`,
		shipmentID,
	).Scan(&shipmentType)
	if err != nil {
		return fmt.Errorf("failed to fetch shipment type: %w", err)
	}

	// Only send for single_full_journey and warehouse_to_engineer
	if shipmentType != string(models.ShipmentTypeSingleFullJourney) && 
		shipmentType != string(models.ShipmentTypeWarehouseToEngineer) {
		return fmt.Errorf("engineer delivery notification not applicable for shipment type: %s", shipmentType)
	}

	// Get engineer details
	if !shipment.SoftwareEngineerID.Valid {
		return fmt.Errorf("shipment has no assigned engineer")
	}

	var engineerName string
	err = n.db.QueryRowContext(ctx,
		`SELECT name FROM software_engineers WHERE id = $1`,
		shipment.SoftwareEngineerID.Int64,
	).Scan(&engineerName)
	if err != nil {
		return fmt.Errorf("failed to fetch engineer details: %w", err)
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

	// Get JIRA ticket and delivery date
	var jiraTicket string
	var deliveredAt sql.NullTime
	err = n.db.QueryRowContext(ctx,
		`SELECT jira_ticket_number, delivered_at FROM shipments WHERE id = $1`,
		shipmentID,
	).Scan(&jiraTicket, &deliveredAt)
	if err != nil {
		return fmt.Errorf("failed to fetch shipment details: %w", err)
	}

	// Get delivery date
	deliveryDate := time.Now().Format("Monday, January 2, 2006")
	if deliveredAt.Valid {
		deliveryDate = deliveredAt.Time.Format("Monday, January 2, 2006")
	}

	// Try to get contact email from pickup form first
	contactEmail := ""
	contactName := clientCompany
	
	var formDataJSON string
	err = n.db.QueryRowContext(ctx,
		`SELECT form_data FROM pickup_forms WHERE shipment_id = $1 ORDER BY submitted_at DESC LIMIT 1`,
		shipmentID,
	).Scan(&formDataJSON)
	
	if err == nil {
		// Parse form data
		var formData map[string]interface{}
		if err := json.Unmarshal([]byte(formDataJSON), &formData); err == nil {
			if email, ok := formData["contact_email"].(string); ok && email != "" {
				contactEmail = email
			}
			if name, ok := formData["contact_name"].(string); ok && name != "" {
				contactName = name
			}
		}
	}
	
	// Fallback: Try to get contact email from client company
	if contactEmail == "" {
		var clientContactInfo sql.NullString
		err = n.db.QueryRowContext(ctx,
			`SELECT contact_info FROM client_companies WHERE id = $1`,
			shipment.ClientCompanyID,
		).Scan(&clientContactInfo)
		
		if err == nil && clientContactInfo.Valid && clientContactInfo.String != "" {
			// Try to extract email from contact_info (could be JSON or plain text)
			var contactInfoMap map[string]interface{}
			if err := json.Unmarshal([]byte(clientContactInfo.String), &contactInfoMap); err == nil {
				// If it's JSON, try to get email field
				if email, ok := contactInfoMap["email"].(string); ok && email != "" {
					contactEmail = email
				} else if email, ok := contactInfoMap["contact_email"].(string); ok && email != "" {
					contactEmail = email
				}
			} else {
				// If it's plain text, check if it looks like an email
				contactStr := clientContactInfo.String
				if strings.Contains(contactStr, "@") {
					// Simple email extraction - take first email-like string
					parts := strings.Fields(contactStr)
					for _, part := range parts {
						if strings.Contains(part, "@") && strings.Contains(part, ".") {
							contactEmail = strings.Trim(part, ",;")
							break
						}
					}
				}
			}
		}
	}
	
	// Final fallback: Try to get email from client company users
	if contactEmail == "" {
		var clientUserEmail sql.NullString
		err = n.db.QueryRowContext(ctx,
			`SELECT email FROM users WHERE client_company_id = $1 LIMIT 1`,
			shipment.ClientCompanyID,
		).Scan(&clientUserEmail)
		
		if err == nil && clientUserEmail.Valid && clientUserEmail.String != "" {
			contactEmail = clientUserEmail.String
		}
	}
	
	// If still no email found, return error
	if contactEmail == "" {
		return fmt.Errorf("no contact email found for shipment %d: tried pickup form, client company contact info, and client company users", shipmentID)
	}

	// Prepare template data
	data := EngineerDeliveryClientData{
		ContactName:    contactName,
		ClientCompany: clientCompany,
		EngineerName:  engineerName,
		DeliveryDate:  deliveryDate,
		TrackingNumber: shipment.TrackingNumber.String,
		JiraTicket:    jiraTicket,
		ProjectName:   "", // Can be added if available
	}

	// Render template
	htmlBody, err := n.templates.RenderTemplate("engineer_delivery_notification_to_client", data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	// Send email
	message := Message{
		To:       []string{contactEmail},
		Subject:  n.templates.GetSubject("engineer_delivery_notification_to_client", data),
		Body:     n.generatePlainTextFromHTML(htmlBody),
		HTMLBody: htmlBody,
	}

	if err := n.client.Send(message); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	// Log notification
	if err := n.logNotification(ctx, shipmentID, "engineer_delivery_notification_to_client", contactEmail, "sent"); err != nil {
		fmt.Printf("Warning: failed to log notification: %v\n", err)
	}

	return nil
}

// SendInTransitToEngineerNotification sends notification to engineer when device is in transit
func (n *Notifier) SendInTransitToEngineerNotification(ctx context.Context, shipmentID int64) error {
	// Fetch shipment details
	shipment, err := n.getShipmentDetails(ctx, shipmentID)
	if err != nil {
		return fmt.Errorf("failed to fetch shipment details: %w", err)
	}

	// Get shipment type - only send for applicable types
	var shipmentType string
	var courierName sql.NullString
	var etaToEngineer sql.NullTime
	err = n.db.QueryRowContext(ctx,
		`SELECT shipment_type, courier_name, eta_to_engineer FROM shipments WHERE id = $1`,
		shipmentID,
	).Scan(&shipmentType, &courierName, &etaToEngineer)
	if err != nil {
		return fmt.Errorf("failed to fetch shipment details: %w", err)
	}

	// Only send for single_full_journey and warehouse_to_engineer
	if shipmentType != string(models.ShipmentTypeSingleFullJourney) && 
		shipmentType != string(models.ShipmentTypeWarehouseToEngineer) {
		return fmt.Errorf("in transit to engineer notification not applicable for shipment type: %s", shipmentType)
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

	// Get ETA - this is required
	eta := "To be determined"
	if etaToEngineer.Valid {
		eta = etaToEngineer.Time.Format("Monday, January 2, 2006 at 3:04 PM")
	}

	// Get courier name
	courierNameStr := "Courier"
	if courierName.Valid && courierName.String != "" {
		courierNameStr = courierName.String
	}

	// Get complete laptop details from the first laptop in the shipment
	var serialNumber, brand, model, cpu, ramGB, ssdGB, sku sql.NullString
	err = n.db.QueryRowContext(ctx,
		`SELECT l.serial_number, l.brand, l.model, l.cpu, l.ram_gb, l.ssd_gb, l.sku
		FROM laptops l
		JOIN shipment_laptops sl ON sl.laptop_id = l.id
		WHERE sl.shipment_id = $1
		LIMIT 1`,
		shipmentID,
	).Scan(&serialNumber, &brand, &model, &cpu, &ramGB, &ssdGB, &sku)
	
	// Set defaults if laptop details not found
	serialNumberStr := ""
	brandStr := ""
	deviceModel := ""
	cpuStr := ""
	ramGBStr := ""
	ssdGBStr := ""
	skuStr := ""
	
	if err == nil {
		if serialNumber.Valid {
			serialNumberStr = serialNumber.String
		}
		if brand.Valid {
			brandStr = brand.String
		}
		if model.Valid {
			deviceModel = model.String
			// Format as "Brand Model" if both available
			if brandStr != "" {
				deviceModel = fmt.Sprintf("%s %s", brandStr, model.String)
			}
		}
		if cpu.Valid {
			cpuStr = cpu.String
		}
		if ramGB.Valid {
			ramGBStr = ramGB.String
		}
		if ssdGB.Valid {
			ssdGBStr = ssdGB.String
		}
		if sku.Valid {
			skuStr = sku.String
		}
	}

	// Build tracking URL based on courier
	trackingURL := ""
	if shipment.TrackingNumber.String != "" {
		switch strings.ToUpper(courierNameStr) {
		case "UPS":
			trackingURL = fmt.Sprintf("https://www.ups.com/track?tracknum=%s", shipment.TrackingNumber.String)
		case "FEDEX", "FEDEX EXPRESS":
			trackingURL = fmt.Sprintf("https://www.fedex.com/fedextrack/?trknbr=%s", shipment.TrackingNumber.String)
		case "DHL":
			trackingURL = fmt.Sprintf("https://www.dhl.com/en/express/tracking.html?AWB=%s", shipment.TrackingNumber.String)
		default:
			trackingURL = fmt.Sprintf("https://www.google.com/search?q=track+%s", shipment.TrackingNumber.String)
		}
	}

	// Get logistics contact info for support
	logisticsEmail, err := n.getLogisticsEmail(ctx)
	if err != nil {
		// Log warning but continue with default email
		fmt.Printf("Warning: %v\n", err)
	}
	contactInfo := fmt.Sprintf("If you have any questions or concerns, please contact logistics at %s", logisticsEmail)

	// Prepare template data
	data := InTransitToEngineerData{
		EngineerName:   engineerName,
		SerialNumber:   serialNumberStr,
		Brand:          brandStr,
		DeviceModel:    deviceModel,
		CPU:            cpuStr,
		RAMGB:          ramGBStr,
		SSDGB:          ssdGBStr,
		SKU:            skuStr,
		TrackingNumber: shipment.TrackingNumber.String,
		CourierName:    courierNameStr,
		ETA:            eta,
		ShipmentURL:    trackingURL, // Use tracking URL for shipment tracking
		ContactInfo:    contactInfo,
	}

	// Render template
	htmlBody, err := n.templates.RenderTemplate("in_transit_to_engineer", data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	// Send email
	message := Message{
		To:       []string{engineerEmail},
		Subject:  n.templates.GetSubject("in_transit_to_engineer", data),
		Body:     n.generatePlainTextFromHTML(htmlBody),
		HTMLBody: htmlBody,
	}

	if err := n.client.Send(message); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	// Log notification
	if err := n.logNotification(ctx, shipmentID, "in_transit_to_engineer", engineerEmail, "sent"); err != nil {
		fmt.Printf("Warning: failed to log notification: %v\n", err)
	}

	return nil
}

// SendReceptionReportApprovalRequest sends notification to logistics when reception report is created
func (n *Notifier) SendReceptionReportApprovalRequest(ctx context.Context, reportID int64) error {
	// Fetch reception report details
	var report models.ReceptionReport
	var warehouseUserEmail string
	err := n.db.QueryRowContext(ctx,
		`SELECT rr.id, rr.laptop_id, rr.shipment_id, rr.client_company_id, rr.tracking_number,
		rr.warehouse_user_id, rr.received_at, rr.notes, rr.photo_serial_number, 
		rr.photo_external_condition, rr.photo_working_condition, rr.status,
		u.email as warehouse_user_email
		FROM reception_reports rr
		LEFT JOIN users u ON u.id = rr.warehouse_user_id
		WHERE rr.id = $1`,
		reportID,
	).Scan(
		&report.ID, &report.LaptopID, &report.ShipmentID, &report.ClientCompanyID, &report.TrackingNumber,
		&report.WarehouseUserID, &report.ReceivedAt, &report.Notes, &report.PhotoSerialNumber,
		&report.PhotoExternalCondition, &report.PhotoWorkingCondition, &report.Status,
		&warehouseUserEmail,
	)
	if err != nil {
		return fmt.Errorf("failed to fetch reception report: %w", err)
	}

	// Get laptop serial number
	var serialNumber string
	err = n.db.QueryRowContext(ctx,
		`SELECT serial_number FROM laptops WHERE id = $1`,
		report.LaptopID,
	).Scan(&serialNumber)
	if err != nil {
		return fmt.Errorf("failed to fetch laptop serial number: %w", err)
	}

	// Get client company name
	// First try from reception report, then fall back to laptop's client_company_id
	clientCompany := "Unknown Company"
	var clientCompanyID *int64 = report.ClientCompanyID
	
	// If reception report doesn't have client_company_id, get it from laptop
	if clientCompanyID == nil {
		var laptopClientCompanyID sql.NullInt64
		err = n.db.QueryRowContext(ctx,
			`SELECT client_company_id FROM laptops WHERE id = $1`,
			report.LaptopID,
		).Scan(&laptopClientCompanyID)
		if err == nil && laptopClientCompanyID.Valid {
			clientCompanyID = &laptopClientCompanyID.Int64
		}
	}
	
	// Fetch company name if we have an ID
	if clientCompanyID != nil {
		err = n.db.QueryRowContext(ctx,
			`SELECT name FROM client_companies WHERE id = $1`,
			*clientCompanyID,
		).Scan(&clientCompany)
		if err != nil {
			// Non-critical, use default
			clientCompany = "Unknown Company"
		}
	}

	// Get warehouse user name/email
	warehouseUser := warehouseUserEmail
	if warehouseUser == "" {
		warehouseUser = "Warehouse Team"
	}

	// Build photo URLs array
	photoURLs := []string{}
	if report.PhotoSerialNumber != "" {
		photoURLs = append(photoURLs, report.PhotoSerialNumber)
	}
	if report.PhotoExternalCondition != "" {
		photoURLs = append(photoURLs, report.PhotoExternalCondition)
	}
	if report.PhotoWorkingCondition != "" {
		photoURLs = append(photoURLs, report.PhotoWorkingCondition)
	}

	// Format received date
	receivedDate := report.ReceivedAt.Format("Monday, January 2, 2006 at 3:04 PM")

	// Build URLs
	reportURL := fmt.Sprintf("/reception-reports/%d", reportID)
	approvalURL := fmt.Sprintf("/reception-reports/%d/approve", reportID)

	// Get logistics email
	logisticsEmail, err := n.getLogisticsEmail(ctx)
	if err != nil {
		// Log warning but continue with default email
		fmt.Printf("Warning: %v\n", err)
	}

	// Prepare template data
	data := ReceptionReportApprovalData{
		ShipmentID:     0, // Will be set if shipment exists
		TrackingNumber: report.TrackingNumber,
		ClientCompany:  clientCompany,
		ReceivedDate:   receivedDate,
		WarehouseUser:  warehouseUser,
		Notes:          report.Notes,
		PhotoURLs:      photoURLs,
		SerialNumber:   serialNumber,
		ReportURL:      reportURL,
		ApprovalURL:    approvalURL,
	}

	// Set shipment ID if available
	if report.ShipmentID != nil {
		data.ShipmentID = *report.ShipmentID
	}

	// Render template
	htmlBody, err := n.templates.RenderTemplate("reception_report_approval_request", data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	// Send email
	message := Message{
		To:       []string{logisticsEmail},
		Subject:  n.templates.GetSubject("reception_report_approval_request", data),
		Body:     n.generatePlainTextFromHTML(htmlBody),
		HTMLBody: htmlBody,
	}

	if err := n.client.Send(message); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	// Log notification
	if err := n.logNotification(ctx, data.ShipmentID, "reception_report_approval_request", logisticsEmail, "sent"); err != nil {
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

// getLogisticsEmail retrieves the logistics team email
// First tries to get from users table, falls back to config if available
func (n *Notifier) getLogisticsEmail(ctx context.Context) (string, error) {
	var logisticsEmail string
	err := n.db.QueryRowContext(ctx,
		`SELECT email FROM users WHERE role = 'logistics' LIMIT 1`,
	).Scan(&logisticsEmail)
	
	if err == nil && logisticsEmail != "" {
		return logisticsEmail, nil
	}
	
	// Fallback to config if available
	if n.config != nil && n.config.LogisticsEmail != "" {
		return n.config.LogisticsEmail, nil
	}
	
	// Final fallback
	if err != nil {
		return "international@bairesdev.com", fmt.Errorf("no logistics user found, using default: %w", err)
	}
	
	return logisticsEmail, nil
}

// getWarehouseEmail retrieves a warehouse user email
// First tries to get from users table, falls back to config if available
func (n *Notifier) getWarehouseEmail(ctx context.Context) (string, error) {
	var warehouseEmail string
	err := n.db.QueryRowContext(ctx,
		`SELECT email FROM users WHERE role = 'warehouse' LIMIT 1`,
	).Scan(&warehouseEmail)
	
	if err == nil && warehouseEmail != "" {
		return warehouseEmail, nil
	}
	
	// Fallback to config if available
	if n.config != nil && n.config.WarehouseEmail != "" {
		return n.config.WarehouseEmail, nil
	}
	
	// Final fallback
	if err != nil {
		return "warehouse@bairesdev.com", fmt.Errorf("no warehouse user found, using default: %w", err)
	}
	
	return warehouseEmail, nil
}

// getContactEmailFromForm retrieves contact email from pickup form
func (n *Notifier) getContactEmailFromForm(ctx context.Context, shipmentID int64) (string, error) {
	var formDataJSON string
	err := n.db.QueryRowContext(ctx,
		`SELECT form_data FROM pickup_forms WHERE shipment_id = $1 ORDER BY submitted_at DESC LIMIT 1`,
		shipmentID,
	).Scan(&formDataJSON)
	if err != nil {
		return "", fmt.Errorf("no pickup form found for shipment %d: %w", shipmentID, err)
	}

	// Parse form data to extract contact email
	var formData map[string]interface{}
	if err := json.Unmarshal([]byte(formDataJSON), &formData); err != nil {
		return "", fmt.Errorf("failed to parse form data: %w", err)
	}

	contactEmail, ok := formData["contact_email"].(string)
	if !ok || contactEmail == "" {
		return "", fmt.Errorf("contact email not found in pickup form")
	}

	return contactEmail, nil
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

