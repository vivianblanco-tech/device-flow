package jira

import (
	"testing"
	"time"

	"github.com/yourusername/laptop-tracking-system/internal/models"
)

// 🟥 RED: Test for mapping JIRA ticket to shipment data
// This test verifies that we can extract shipment-relevant data from a JIRA ticket
func TestMapTicketToShipmentData(t *testing.T) {
	ticket := &Ticket{
		Key:         "PROJ-123",
		Summary:     "Laptop deployment for John Doe",
		Description: "Deploy laptop to software engineer",
		Status:      "In Progress",
		Assignee:    "Jane Smith",
		Created:     "2023-10-01T10:00:00.000+0000",
		Updated:     "2023-10-02T15:30:00.000+0000",
	}

	// Map the ticket to shipment data
	shipmentData, err := MapTicketToShipmentData(ticket)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify mapped data
	if shipmentData.JiraTicketKey != "PROJ-123" {
		t.Errorf("expected JiraTicketKey PROJ-123, got %s", shipmentData.JiraTicketKey)
	}
	if shipmentData.Summary != "Laptop deployment for John Doe" {
		t.Errorf("expected correct summary, got %s", shipmentData.Summary)
	}
	if shipmentData.Status != "In Progress" {
		t.Errorf("expected status 'In Progress', got %s", shipmentData.Status)
	}
}

// 🟥 RED: Test for extracting custom fields from JIRA ticket
// This test verifies that we can extract custom field data like serial numbers
func TestExtractCustomFields(t *testing.T) {
	// Simulate a ticket response with custom fields
	ticketData := map[string]interface{}{
		"key": "PROJ-124",
		"fields": map[string]interface{}{
			"summary":     "Hardware deployment",
			"description": "Deploy hardware to engineer",
			"customfield_10001": "SN123456789",     // Serial number custom field
			"customfield_10002": "john@example.com", // Engineer email custom field
			"customfield_10003": "Acme Corp",        // Client company custom field
		},
	}

	// Extract custom fields
	customFields, err := ExtractCustomFields(ticketData)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify extracted fields
	if customFields.SerialNumber != "SN123456789" {
		t.Errorf("expected serial number SN123456789, got %s", customFields.SerialNumber)
	}
	if customFields.EngineerEmail != "john@example.com" {
		t.Errorf("expected engineer email john@example.com, got %s", customFields.EngineerEmail)
	}
	if customFields.ClientCompany != "Acme Corp" {
		t.Errorf("expected client company Acme Corp, got %s", customFields.ClientCompany)
	}
}

// 🟥 RED: Test for mapping JIRA status to shipment status
// This test verifies that we can map JIRA workflow statuses to our internal statuses
func TestMapJiraStatusToShipmentStatus(t *testing.T) {
	tests := []struct {
		jiraStatus     string
		expectedStatus models.ShipmentStatus
	}{
		{"To Do", models.ShipmentStatusPendingPickup},
		{"Pending Pickup", models.ShipmentStatusPendingPickup},
		{"Picked Up", models.ShipmentStatusPickedUpFromClient},
		{"In Transit to Warehouse", models.ShipmentStatusInTransitToWarehouse},
		{"At Warehouse", models.ShipmentStatusAtWarehouse},
		{"Released from Warehouse", models.ShipmentStatusReleasedFromWarehouse},
		{"In Transit to Engineer", models.ShipmentStatusInTransitToEngineer},
		{"Delivered", models.ShipmentStatusDelivered},
		{"Done", models.ShipmentStatusDelivered},
	}

	for _, tt := range tests {
		t.Run(tt.jiraStatus, func(t *testing.T) {
			status := MapJiraStatusToShipmentStatus(tt.jiraStatus)
			if status != tt.expectedStatus {
				t.Errorf("expected status %s for JIRA status %s, got %s",
					tt.expectedStatus, tt.jiraStatus, status)
			}
		})
	}
}

// 🟥 RED: Test for mapping unknown JIRA status
// This test verifies that unknown statuses map to a default status
func TestMapJiraStatusToShipmentStatus_Unknown(t *testing.T) {
	status := MapJiraStatusToShipmentStatus("Unknown Status")
	if status != models.ShipmentStatusPendingPickup {
		t.Errorf("expected default status %s, got %s",
			models.ShipmentStatusPendingPickup, status)
	}
}

// 🟥 RED: Test for creating shipment from JIRA ticket
// This test verifies that we can create a complete shipment from a JIRA ticket
func TestCreateShipmentFromTicket(t *testing.T) {
	ticket := &Ticket{
		Key:         "PROJ-125",
		Summary:     "Deploy laptop to engineer",
		Description: "Deploy laptop with serial SN987654321",
		Status:      "Pending Pickup",
		Created:     "2023-10-01T10:00:00.000+0000",
	}

	customFields := &CustomFields{
		SerialNumber:  "SN987654321",
		EngineerEmail: "engineer@example.com",
		ClientCompany: "Test Corp",
	}

	// Create shipment from ticket
	shipment, err := CreateShipmentFromTicket(ticket, customFields)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify shipment data
	if shipment.Status != models.ShipmentStatusPendingPickup {
		t.Errorf("expected status %s, got %s",
			models.ShipmentStatusPendingPickup, shipment.Status)
	}
	if shipment.Notes == "" {
		t.Error("expected notes to be populated")
	}
	if shipment.CreatedAt.IsZero() {
		t.Error("expected created_at to be set")
	}
}

// 🟥 RED: Test for parsing JIRA timestamp
// This test verifies that we can parse JIRA's timestamp format
func TestParseJiraTimestamp(t *testing.T) {
	jiraTime := "2023-10-01T10:00:00.000+0000"
	
	parsedTime, err := ParseJiraTimestamp(jiraTime)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedYear := 2023
	expectedMonth := time.October
	expectedDay := 1

	if parsedTime.Year() != expectedYear {
		t.Errorf("expected year %d, got %d", expectedYear, parsedTime.Year())
	}
	if parsedTime.Month() != expectedMonth {
		t.Errorf("expected month %s, got %s", expectedMonth, parsedTime.Month())
	}
	if parsedTime.Day() != expectedDay {
		t.Errorf("expected day %d, got %d", expectedDay, parsedTime.Day())
	}
}

