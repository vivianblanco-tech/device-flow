package email

import (
	"strings"
	"testing"
	"time"
)

func TestNewEmailTemplates(t *testing.T) {
	templates := NewEmailTemplates()

	if templates == nil {
		t.Fatal("NewEmailTemplates() returned nil")
	}

	// Verify that all expected templates are loaded
	expectedTemplates := []string{
		"magic_link",
		"address_confirmation",
		"pickup_confirmation",
		"warehouse_pre_alert",
		"release_notification",
		"delivery_confirmation",
	}

	for _, name := range expectedTemplates {
		if _, exists := templates.templates[name]; !exists {
			t.Errorf("Expected template '%s' not loaded", name)
		}
	}
}

func TestEmailTemplates_RenderTemplate_MagicLink(t *testing.T) {
	templates := NewEmailTemplates()

	data := MagicLinkData{
		RecipientName: "John Doe",
		MagicLink:     "https://example.com/form?token=abc123",
		ExpiresAt:     time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
		FormType:      "pickup",
	}

	html, err := templates.RenderTemplate("magic_link", data)
	if err != nil {
		t.Fatalf("RenderTemplate() error = %v", err)
	}

	// Check that the rendered HTML contains expected content
	expectedContent := []string{
		"John Doe",
		"https://example.com/form?token=abc123",
		"pickup",
		"Access Form",
		"Security Notice",
	}

	for _, expected := range expectedContent {
		if !strings.Contains(html, expected) {
			t.Errorf("Rendered HTML missing expected content: %s", expected)
		}
	}

	// Verify HTML structure
	if !strings.Contains(html, "<!DOCTYPE html>") {
		t.Error("Rendered HTML missing DOCTYPE declaration")
	}

	if !strings.Contains(html, "<html") {
		t.Error("Rendered HTML missing html tag")
	}
}

func TestEmailTemplates_RenderTemplate_AddressConfirmation(t *testing.T) {
	templates := NewEmailTemplates()

	data := AddressConfirmationData{
		EngineerName:    "Jane Smith",
		CompanyName:     "BairesDev",
		ProjectName:     "Project Phoenix",
		ExpectedDate:    "January 15, 2024",
		ConfirmationURL: "https://example.com/confirm-address?id=123",
	}

	html, err := templates.RenderTemplate("address_confirmation", data)
	if err != nil {
		t.Fatalf("RenderTemplate() error = %v", err)
	}

	expectedContent := []string{
		"Jane Smith",
		"Project Phoenix",
		"January 15, 2024",
		"https://example.com/confirm-address?id=123",
		"Confirm Address",
	}

	for _, expected := range expectedContent {
		if !strings.Contains(html, expected) {
			t.Errorf("Rendered HTML missing expected content: %s", expected)
		}
	}
}

func TestEmailTemplates_RenderTemplate_PickupConfirmation(t *testing.T) {
	templates := NewEmailTemplates()

	data := PickupConfirmationData{
		ClientName:       "Alice Johnson",
		ClientCompany:    "TechCorp",
		TrackingNumber:   "UPS123456789",
		PickupDate:       "December 20, 2024",
		PickupTimeSlot:   "Morning (9AM - 12PM)",
		NumberOfDevices:  2,
		ConfirmationCode: "CONF-2024-001",
	}

	html, err := templates.RenderTemplate("pickup_confirmation", data)
	if err != nil {
		t.Fatalf("RenderTemplate() error = %v", err)
	}

	expectedContent := []string{
		"Alice Johnson",
		"UPS123456789",
		"December 20, 2024",
		"Morning (9AM - 12PM)",
		"CONF-2024-001",
		"2",
	}

	for _, expected := range expectedContent {
		if !strings.Contains(html, expected) {
			t.Errorf("Rendered HTML missing expected content: %s", expected)
		}
	}
}

func TestEmailTemplates_RenderTemplate_WarehousePreAlert(t *testing.T) {
	templates := NewEmailTemplates()

	// Test with single shipment
	data := WarehousePreAlertData{
		TrackingNumber:    "UPS987654321",
		ExpectedDate:      "December 22, 2024",
		ShipperName:       "Bob Wilson",
		ShipperCompany:    "ClientCorp",
		DeviceDescription: "Dell Latitude Laptop",
		ProjectName:       "Project Alpha",
		TrackingURL:       "https://www.ups.com/track?tracknum=UPS987654321",
		IsSingleShipment:  true,
		SerialNumber:      "SN123456789",
		Brand:             "Dell",
		Model:             "Latitude",
		CPU:               "Intel i7",
		RAMGB:             "16GB",
		SSDGB:             "512GB",
	}

	html, err := templates.RenderTemplate("warehouse_pre_alert", data)
	if err != nil {
		t.Fatalf("RenderTemplate() error = %v", err)
	}

	expectedContent := []string{
		"UPS987654321",
		"December 22, 2024",
		"Bob Wilson",
		"Project Alpha",
		"Track Shipment",
		"Laptop Details",
		"SN123456789",
		"Dell",
		"Latitude",
		"Intel i7",
		"16GB",
		"512GB",
	}

	for _, expected := range expectedContent {
		if !strings.Contains(html, expected) {
			t.Errorf("Rendered HTML missing expected content: %s", expected)
		}
	}
}

func TestEmailTemplates_RenderTemplate_ReleaseNotification(t *testing.T) {
	templates := NewEmailTemplates()

	data := ReleaseNotificationData{
		CourierName:        "Courier Service",
		CourierCompany:     "FastShip Inc",
		PickupDate:         "December 25, 2024",
		PickupTimeSlot:     "Afternoon (1PM - 5PM)",
		WarehouseAddress:   "123 Warehouse St, City, State 12345",
		ContactPerson:      "Warehouse Manager",
		ContactPhone:       "+1-555-0123",
		DeviceSerialNumber: "SN123456789",
		EngineerName:       "Chris Davis",
		TrackingNumber:     "UPS111222333",
	}

	html, err := templates.RenderTemplate("release_notification", data)
	if err != nil {
		t.Fatalf("RenderTemplate() error = %v", err)
	}

	expectedContent := []string{
		"Courier Service",
		"December 25, 2024",
		"Afternoon (1PM - 5PM)",
		"123 Warehouse St, City, State 12345",
		"Warehouse Manager",
		"1-555-0123", // Phone number (+ might be HTML escaped)
		"SN123456789",
		"Chris Davis",
		"UPS111222333",
	}

	for _, expected := range expectedContent {
		if !strings.Contains(html, expected) {
			t.Errorf("Rendered HTML missing expected content: %s", expected)
		}
	}
}

func TestEmailTemplates_RenderTemplate_DeliveryConfirmation(t *testing.T) {
	templates := NewEmailTemplates()

	data := DeliveryConfirmationData{
		EngineerName:       "Diana Martinez",
		DeviceSerialNumber: "SN987654321",
		DeviceModel:        "HP EliteBook 840",
		DeliveryDate:       "December 26, 2024",
		TrackingNumber:     "UPS444555666",
		ProjectName:        "Project Beta",
	}

	html, err := templates.RenderTemplate("delivery_confirmation", data)
	if err != nil {
		t.Fatalf("RenderTemplate() error = %v", err)
	}

	expectedContent := []string{
		"Diana Martinez",
		"SN987654321",
		"HP EliteBook 840",
		"December 26, 2024",
		"UPS444555666",
		"Project Beta",
	}

	for _, expected := range expectedContent {
		if !strings.Contains(html, expected) {
			t.Errorf("Rendered HTML missing expected content: %s", expected)
		}
	}
}

func TestEmailTemplates_RenderTemplate_InvalidTemplate(t *testing.T) {
	templates := NewEmailTemplates()

	data := MagicLinkData{
		RecipientName: "Test User",
		MagicLink:     "https://example.com/test",
		ExpiresAt:     time.Now(),
		FormType:      "test",
	}

	_, err := templates.RenderTemplate("non_existent_template", data)
	if err == nil {
		t.Error("RenderTemplate() should return error for non-existent template")
	}

	expectedErr := "template 'non_existent_template' not found"
	if err.Error() != expectedErr {
		t.Errorf("RenderTemplate() error = %v, want %v", err.Error(), expectedErr)
	}
}

func TestEmailTemplates_GetSubject(t *testing.T) {
	templates := NewEmailTemplates()

	tests := []struct {
		name         string
		templateName string
		data         interface{}
		want         string
	}{
		{
			name:         "magic link subject",
			templateName: "magic_link",
			data: MagicLinkData{
				FormType: "pickup",
			},
			want: "Access Your Form - pickup",
		},
		{
			name:         "address confirmation subject",
			templateName: "address_confirmation",
			data:         AddressConfirmationData{},
			want:         "Confirm Your Delivery Address",
		},
		{
			name:         "pickup confirmation subject",
			templateName: "pickup_confirmation",
			data: PickupConfirmationData{
				ConfirmationCode: "CONF-123",
			},
			want: "Pickup Confirmation - CONF-123",
		},
		{
			name:         "warehouse pre-alert subject",
			templateName: "warehouse_pre_alert",
			data: WarehousePreAlertData{
				TrackingNumber: "UPS123",
			},
			want: "Incoming Shipment Alert - UPS123",
		},
		{
			name:         "release notification subject",
			templateName: "release_notification",
			data: ReleaseNotificationData{
				TrackingNumber: "UPS456",
			},
			want: "Hardware Release for Pickup - UPS456",
		},
		{
			name:         "delivery confirmation subject",
			templateName: "delivery_confirmation",
			data:         DeliveryConfirmationData{},
			want:         "Device Delivered Successfully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := templates.GetSubject(tt.templateName, tt.data)
			if got != tt.want {
				t.Errorf("GetSubject() = %v, want %v", got, tt.want)
			}
		})
	}
}

