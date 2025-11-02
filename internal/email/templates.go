package email

import (
	"bytes"
	"fmt"
	"html/template"
	"time"
)

// TemplateData holds data for rendering email templates
type TemplateData map[string]interface{}

// Common email template data structures

// MagicLinkData contains data for magic link emails
type MagicLinkData struct {
	RecipientName string
	MagicLink     string
	ExpiresAt     time.Time
	FormType      string // "pickup", "delivery", etc.
}

// AddressConfirmationData contains data for address confirmation emails
type AddressConfirmationData struct {
	EngineerName    string
	CompanyName     string
	ProjectName     string
	ExpectedDate    string
	ConfirmationURL string
}

// PickupConfirmationData contains data for pickup confirmation emails
type PickupConfirmationData struct {
	ClientName       string
	ClientCompany    string
	TrackingNumber   string
	PickupDate       string
	PickupTimeSlot   string
	NumberOfDevices  int
	ConfirmationCode string
}

// WarehousePreAlertData contains data for warehouse pre-alert emails
type WarehousePreAlertData struct {
	TrackingNumber    string
	ExpectedDate      string
	ShipperName       string
	ShipperCompany    string
	DeviceDescription string
	ProjectName       string
	TrackingURL       string
}

// ReleaseNotificationData contains data for release notification emails
type ReleaseNotificationData struct {
	CourierName        string
	CourierCompany     string
	PickupDate         string
	PickupTimeSlot     string
	WarehouseAddress   string
	ContactPerson      string
	ContactPhone       string
	DeviceSerialNumber string
	EngineerName       string
	TrackingNumber     string
}

// DeliveryConfirmationData contains data for delivery confirmation emails
type DeliveryConfirmationData struct {
	EngineerName       string
	DeviceSerialNumber string
	DeviceModel        string
	DeliveryDate       string
	TrackingNumber     string
	ProjectName        string
}

// EmailTemplates holds all compiled email templates
type EmailTemplates struct {
	templates map[string]*template.Template
}

// NewEmailTemplates creates a new email templates instance
func NewEmailTemplates() *EmailTemplates {
	et := &EmailTemplates{
		templates: make(map[string]*template.Template),
	}

	// Load all templates
	et.loadTemplates()

	return et
}

// loadTemplates compiles all email templates
func (et *EmailTemplates) loadTemplates() {
	// Base template with common styling
	baseTemplate := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Subject}}</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f5f5f5;
        }
        .email-container {
            background-color: #ffffff;
            padding: 40px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .header {
            border-bottom: 3px solid #0052CC;
            padding-bottom: 20px;
            margin-bottom: 30px;
        }
        .header h1 {
            color: #0052CC;
            margin: 0;
            font-size: 24px;
        }
        .content {
            margin-bottom: 30px;
        }
        .button {
            display: inline-block;
            padding: 12px 24px;
            background-color: #0052CC;
            color: #ffffff !important;
            text-decoration: none;
            border-radius: 4px;
            margin: 20px 0;
            font-weight: 600;
        }
        .button:hover {
            background-color: #0747A6;
        }
        .info-box {
            background-color: #F4F5F7;
            padding: 20px;
            border-radius: 4px;
            margin: 20px 0;
        }
        .info-box h3 {
            margin-top: 0;
            color: #0052CC;
        }
        .info-row {
            margin: 10px 0;
        }
        .info-label {
            font-weight: 600;
            color: #172B4D;
        }
        .footer {
            margin-top: 40px;
            padding-top: 20px;
            border-top: 1px solid #DFE1E6;
            font-size: 12px;
            color: #6B778C;
            text-align: center;
        }
        .warning {
            background-color: #FFF3CD;
            border-left: 4px solid #FFC107;
            padding: 15px;
            margin: 20px 0;
        }
        .success {
            background-color: #D4EDDA;
            border-left: 4px solid #28A745;
            padding: 15px;
            margin: 20px 0;
        }
    </style>
</head>
<body>
    <div class="email-container">
        {{template "content" .}}
        <div class="footer">
            <p>This is an automated message from the Laptop Tracking System.</p>
            <p>© {{.Year}} BairesDev. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`

	// Magic Link Template
	et.templates["magic_link"] = template.Must(template.New("base").Parse(baseTemplate))
	template.Must(et.templates["magic_link"].New("content").Parse(`
        <div class="header">
            <h1>🔗 Access Your Form</h1>
        </div>
        <div class="content">
            <p>Hello {{.RecipientName}},</p>
            <p>You've been granted access to complete the {{.FormType}} form. Click the button below to get started:</p>
            <p style="text-align: center;">
                <a href="{{.MagicLink}}" class="button">Access Form</a>
            </p>
            <div class="warning">
                <strong>⚠️ Security Notice:</strong> This link is valid for one use only and expires on <strong>{{.ExpiresAtFormatted}}</strong>. Do not share this link with others.
            </div>
            <p>If you didn't request this form, please ignore this email.</p>
        </div>
    `))

	// Address Confirmation Template
	et.templates["address_confirmation"] = template.Must(template.New("base").Parse(baseTemplate))
	template.Must(et.templates["address_confirmation"].New("content").Parse(`
        <div class="header">
            <h1>📍 Confirm Your Delivery Address</h1>
        </div>
        <div class="content">
            <p>Hello {{.EngineerName}},</p>
            <p>We're preparing to ship configured hardware for the <strong>{{.ProjectName}}</strong> project. The expected delivery date is <strong>{{.ExpectedDate}}</strong>.</p>
            <p>Please confirm or update your delivery address to ensure successful delivery:</p>
            <p style="text-align: center;">
                <a href="{{.ConfirmationURL}}" class="button">Confirm Address</a>
            </p>
            <div class="info-box">
                <h3>What You Need to Do:</h3>
                <ol>
                    <li>Click the button above</li>
                    <li>Verify your current address is correct</li>
                    <li>Update if necessary</li>
                    <li>Submit the confirmation</li>
                </ol>
            </div>
            <p>If you have any questions, please contact your project manager.</p>
        </div>
    `))

	// Pickup Confirmation Template
	et.templates["pickup_confirmation"] = template.Must(template.New("base").Parse(baseTemplate))
	template.Must(et.templates["pickup_confirmation"].New("content").Parse(`
        <div class="header">
            <h1>✅ Pickup Request Confirmed</h1>
        </div>
        <div class="content">
            <p>Hello {{.ClientName}},</p>
            <div class="success">
                Thank you for completing the hardware shipping form. Your pickup has been scheduled!
            </div>
            <div class="info-box">
                <h3>📦 Pickup Details</h3>
                <div class="info-row">
                    <span class="info-label">Confirmation Code:</span> {{.ConfirmationCode}}
                </div>
                <div class="info-row">
                    <span class="info-label">Pickup Date:</span> {{.PickupDate}}
                </div>
                <div class="info-row">
                    <span class="info-label">Time Slot:</span> {{.PickupTimeSlot}}
                </div>
                <div class="info-row">
                    <span class="info-label">Number of Devices:</span> {{.NumberOfDevices}}
                </div>
                {{if .TrackingNumber}}
                <div class="info-row">
                    <span class="info-label">Tracking Number:</span> {{.TrackingNumber}}
                </div>
                {{end}}
            </div>
            <div class="info-box">
                <h3>📋 What Happens Next:</h3>
                <ol>
                    <li>You'll receive UPS shipping labels via email</li>
                    <li>Print and attach the labels to your package</li>
                    <li>Have the device(s) ready for pickup</li>
                    <li>The courier will collect the package during the scheduled time slot</li>
                </ol>
            </div>
            <p>If you need to make any changes, please contact us immediately.</p>
        </div>
    `))

	// Warehouse Pre-Alert Template
	et.templates["warehouse_pre_alert"] = template.Must(template.New("base").Parse(baseTemplate))
	template.Must(et.templates["warehouse_pre_alert"].New("content").Parse(`
        <div class="header">
            <h1>📬 Incoming Shipment Alert</h1>
        </div>
        <div class="content">
            <p>Hello Warehouse Team,</p>
            <p>Please be advised that a hardware shipment is scheduled for delivery to our facility.</p>
            <div class="info-box">
                <h3>📦 Shipment Details</h3>
                <div class="info-row">
                    <span class="info-label">Tracking Number:</span> {{.TrackingNumber}}
                </div>
                <div class="info-row">
                    <span class="info-label">Expected Delivery:</span> {{.ExpectedDate}}
                </div>
                <div class="info-row">
                    <span class="info-label">Shipper:</span> {{.ShipperName}} ({{.ShipperCompany}})
                </div>
                <div class="info-row">
                    <span class="info-label">Contents:</span> {{.DeviceDescription}}
                </div>
                {{if .ProjectName}}
                <div class="info-row">
                    <span class="info-label">Project:</span> {{.ProjectName}}
                </div>
                {{end}}
            </div>
            {{if .TrackingURL}}
            <p style="text-align: center;">
                <a href="{{.TrackingURL}}" class="button">Track Shipment</a>
            </p>
            {{end}}
            <div class="info-box">
                <h3>✅ Action Required</h3>
                <p>Upon receipt of this package, please:</p>
                <ol>
                    <li>Verify the package condition and contents</li>
                    <li>Complete the Hardware Reception Report</li>
                    <li>Upload photos of the device</li>
                    <li>Submit the report immediately</li>
                </ol>
            </div>
            <p>Please confirm receipt of this notification and contact logistics immediately if there are any issues.</p>
        </div>
    `))

	// Release Notification Template
	et.templates["release_notification"] = template.Must(template.New("base").Parse(baseTemplate))
	template.Must(et.templates["release_notification"].New("content").Parse(`
        <div class="header">
            <h1>🚚 Hardware Release for Pickup</h1>
        </div>
        <div class="content">
            <p>Hello {{.CourierName}},</p>
            <p>Hardware has been released from our warehouse and is ready for pickup and delivery to the engineer.</p>
            <div class="info-box">
                <h3>📦 Pickup Details</h3>
                <div class="info-row">
                    <span class="info-label">Pickup Date:</span> {{.PickupDate}}
                </div>
                <div class="info-row">
                    <span class="info-label">Time Slot:</span> {{.PickupTimeSlot}}
                </div>
                <div class="info-row">
                    <span class="info-label">Tracking Number:</span> {{.TrackingNumber}}
                </div>
            </div>
            <div class="info-box">
                <h3>📍 Pickup Location</h3>
                <div class="info-row">
                    <span class="info-label">Address:</span> {{.WarehouseAddress}}
                </div>
                <div class="info-row">
                    <span class="info-label">Contact Person:</span> {{.ContactPerson}}
                </div>
                <div class="info-row">
                    <span class="info-label">Contact Phone:</span> {{.ContactPhone}}
                </div>
            </div>
            <div class="info-box">
                <h3>📋 Device Information</h3>
                <div class="info-row">
                    <span class="info-label">Serial Number:</span> {{.DeviceSerialNumber}}
                </div>
                <div class="info-row">
                    <span class="info-label">Deliver To:</span> {{.EngineerName}}
                </div>
            </div>
            <p>Please confirm pickup and update the tracking status once the device is collected.</p>
        </div>
    `))

	// Delivery Confirmation Template
	et.templates["delivery_confirmation"] = template.Must(template.New("base").Parse(baseTemplate))
	template.Must(et.templates["delivery_confirmation"].New("content").Parse(`
        <div class="header">
            <h1>✅ Device Delivered Successfully</h1>
        </div>
        <div class="content">
            <p>Hello {{.EngineerName}},</p>
            <div class="success">
                Your device has been successfully delivered! Welcome to the team!
            </div>
            <div class="info-box">
                <h3>📦 Delivery Details</h3>
                <div class="info-row">
                    <span class="info-label">Delivery Date:</span> {{.DeliveryDate}}
                </div>
                <div class="info-row">
                    <span class="info-label">Device Model:</span> {{.DeviceModel}}
                </div>
                <div class="info-row">
                    <span class="info-label">Serial Number:</span> {{.DeviceSerialNumber}}
                </div>
                {{if .TrackingNumber}}
                <div class="info-row">
                    <span class="info-label">Tracking Number:</span> {{.TrackingNumber}}
                </div>
                {{end}}
                {{if .ProjectName}}
                <div class="info-row">
                    <span class="info-label">Project:</span> {{.ProjectName}}
                </div>
                {{end}}
            </div>
            <div class="info-box">
                <h3>📋 Next Steps</h3>
                <ol>
                    <li>Inspect the device for any shipping damage</li>
                    <li>Set up your device following the included instructions</li>
                    <li>Install required software</li>
                    <li>Contact IT support if you encounter any issues</li>
                </ol>
            </div>
            <p>If you have any questions or concerns about your device, please contact your project manager.</p>
        </div>
    `))
}

// RenderTemplate renders an email template with the given data
func (et *EmailTemplates) RenderTemplate(templateName string, data interface{}) (string, error) {
	tmpl, exists := et.templates[templateName]
	if !exists {
		return "", fmt.Errorf("template '%s' not found", templateName)
	}

	// Add common data
	dataMap := make(map[string]interface{})
	dataMap["Year"] = time.Now().Year()

	// Merge with provided data
	switch v := data.(type) {
	case map[string]interface{}:
		for k, val := range v {
			dataMap[k] = val
		}
	case MagicLinkData:
		dataMap["RecipientName"] = v.RecipientName
		dataMap["MagicLink"] = v.MagicLink
		dataMap["ExpiresAtFormatted"] = v.ExpiresAt.Format("Monday, January 2, 2006 at 3:04 PM")
		dataMap["FormType"] = v.FormType
		dataMap["Subject"] = "Access Your Form - " + v.FormType
	case AddressConfirmationData:
		dataMap["EngineerName"] = v.EngineerName
		dataMap["CompanyName"] = v.CompanyName
		dataMap["ProjectName"] = v.ProjectName
		dataMap["ExpectedDate"] = v.ExpectedDate
		dataMap["ConfirmationURL"] = v.ConfirmationURL
		dataMap["Subject"] = "Confirm Your Delivery Address"
	case PickupConfirmationData:
		dataMap["ClientName"] = v.ClientName
		dataMap["ClientCompany"] = v.ClientCompany
		dataMap["TrackingNumber"] = v.TrackingNumber
		dataMap["PickupDate"] = v.PickupDate
		dataMap["PickupTimeSlot"] = v.PickupTimeSlot
		dataMap["NumberOfDevices"] = v.NumberOfDevices
		dataMap["ConfirmationCode"] = v.ConfirmationCode
		dataMap["Subject"] = "Pickup Confirmation - " + v.ConfirmationCode
	case WarehousePreAlertData:
		dataMap["TrackingNumber"] = v.TrackingNumber
		dataMap["ExpectedDate"] = v.ExpectedDate
		dataMap["ShipperName"] = v.ShipperName
		dataMap["ShipperCompany"] = v.ShipperCompany
		dataMap["DeviceDescription"] = v.DeviceDescription
		dataMap["ProjectName"] = v.ProjectName
		dataMap["TrackingURL"] = v.TrackingURL
		dataMap["Subject"] = "Incoming Shipment Alert - " + v.TrackingNumber
	case ReleaseNotificationData:
		dataMap["CourierName"] = v.CourierName
		dataMap["CourierCompany"] = v.CourierCompany
		dataMap["PickupDate"] = v.PickupDate
		dataMap["PickupTimeSlot"] = v.PickupTimeSlot
		dataMap["WarehouseAddress"] = v.WarehouseAddress
		dataMap["ContactPerson"] = v.ContactPerson
		dataMap["ContactPhone"] = v.ContactPhone
		dataMap["DeviceSerialNumber"] = v.DeviceSerialNumber
		dataMap["EngineerName"] = v.EngineerName
		dataMap["TrackingNumber"] = v.TrackingNumber
		dataMap["Subject"] = "Hardware Release for Pickup - " + v.TrackingNumber
	case DeliveryConfirmationData:
		dataMap["EngineerName"] = v.EngineerName
		dataMap["DeviceSerialNumber"] = v.DeviceSerialNumber
		dataMap["DeviceModel"] = v.DeviceModel
		dataMap["DeliveryDate"] = v.DeliveryDate
		dataMap["TrackingNumber"] = v.TrackingNumber
		dataMap["ProjectName"] = v.ProjectName
		dataMap["Subject"] = "Device Delivered Successfully"
	default:
		return "", fmt.Errorf("unsupported data type for template")
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, dataMap); err != nil {
		return "", fmt.Errorf("failed to render template: %w", err)
	}

	return buf.String(), nil
}

// GetSubject extracts the subject from template data
func (et *EmailTemplates) GetSubject(templateName string, data interface{}) string {
	switch v := data.(type) {
	case MagicLinkData:
		return "Access Your Form - " + v.FormType
	case AddressConfirmationData:
		return "Confirm Your Delivery Address"
	case PickupConfirmationData:
		return "Pickup Confirmation - " + v.ConfirmationCode
	case WarehousePreAlertData:
		return "Incoming Shipment Alert - " + v.TrackingNumber
	case ReleaseNotificationData:
		return "Hardware Release for Pickup - " + v.TrackingNumber
	case DeliveryConfirmationData:
		return "Device Delivered Successfully"
	default:
		return "Notification from Laptop Tracking System"
	}
}

