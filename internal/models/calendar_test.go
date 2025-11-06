package models

import (
	"database/sql"
	"testing"
	"time"

	"github.com/yourusername/laptop-tracking-system/internal/database"
)

// TestGetCalendarEvents tests retrieving calendar events from shipments
func TestGetCalendarEvents(t *testing.T) {
	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	// Create test client company
	clientCompany := &ClientCompany{
		Name:        "Test Corp",
		ContactInfo: "contact@testcorp.com",
	}
	err := createClientCompany(db, clientCompany)
	if err != nil {
		t.Fatalf("Failed to create client company: %v", err)
	}

	// Create test software engineer
	engineer := &SoftwareEngineer{
		Name:    "John Doe",
		Email:   "john@example.com",
		Address: "123 Main St",
	}
	err = createSoftwareEngineer(db, engineer)
	if err != nil {
		t.Fatalf("Failed to create software engineer: %v", err)
	}

	// Create test shipments with various dates
	now := time.Now()
	tomorrow := now.AddDate(0, 0, 1)
	nextWeek := now.AddDate(0, 0, 7)

	shipment1 := &Shipment{
		ClientCompanyID:     clientCompany.ID,
		SoftwareEngineerID:  &engineer.ID,
		Status:              ShipmentStatusPendingPickup,
		JiraTicketNumber:    "TEST-600",
		PickupScheduledDate: &tomorrow,
	}
	shipment1.BeforeCreate()
	err = createShipment(db, shipment1)
	if err != nil {
		t.Fatalf("Failed to create shipment 1: %v", err)
	}

	shipment2 := &Shipment{
		ClientCompanyID:     clientCompany.ID,
		SoftwareEngineerID:  &engineer.ID,
		Status:              ShipmentStatusDelivered,
		JiraTicketNumber:    "TEST-601",
		PickupScheduledDate: &now,
		PickedUpAt:          &now,
		DeliveredAt:         &nextWeek,
	}
	shipment2.BeforeCreate()
	err = createShipment(db, shipment2)
	if err != nil {
		t.Fatalf("Failed to create shipment 2: %v", err)
	}

	// Test: Get calendar events
	events, err := GetCalendarEvents(db, now.AddDate(0, 0, -1), now.AddDate(0, 0, 14))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(events) == 0 {
		t.Error("Expected to retrieve calendar events, got none")
	}

	// Verify event types exist
	hasPickupEvent := false
	hasDeliveryEvent := false
	for _, event := range events {
		if event.Type == CalendarEventTypePickup {
			hasPickupEvent = true
		}
		if event.Type == CalendarEventTypeDelivery {
			hasDeliveryEvent = true
		}
	}

	if !hasPickupEvent {
		t.Error("Expected to find pickup events")
	}
	if !hasDeliveryEvent {
		t.Error("Expected to find delivery events")
	}
}

// TestGetCalendarEventsWithDateFilter tests calendar events with date filtering
func TestGetCalendarEventsWithDateFilter(t *testing.T) {
	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	// Create test data
	clientCompany := &ClientCompany{
		Name:        "Test Corp",
		ContactInfo: "contact@testcorp.com",
	}
	err := createClientCompany(db, clientCompany)
	if err != nil {
		t.Fatalf("Failed to create client company: %v", err)
	}

	now := time.Now()
	futureDate := now.AddDate(0, 1, 0) // 1 month in the future

	shipment := &Shipment{
		ClientCompanyID:     clientCompany.ID,
		Status:              ShipmentStatusPendingPickup,
		JiraTicketNumber:    "TEST-602",
		PickupScheduledDate: &futureDate,
	}
	shipment.BeforeCreate()
	err = createShipment(db, shipment)
	if err != nil {
		t.Fatalf("Failed to create shipment: %v", err)
	}

	// Test: Query within date range (should include the event)
	events, err := GetCalendarEvents(db, now, now.AddDate(0, 2, 0))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(events) == 0 {
		t.Error("Expected to find events within date range")
	}

	// Test: Query outside date range (should not include the event)
	events, err = GetCalendarEvents(db, now.AddDate(0, -2, 0), now.AddDate(0, -1, 0))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(events) != 0 {
		t.Errorf("Expected 0 events outside date range, got %d", len(events))
	}
}

// TestCalendarEventFormatting tests the formatting of calendar events
func TestCalendarEventFormatting(t *testing.T) {
	now := time.Now()

	event := &CalendarEvent{
		ID:          1,
		Type:        CalendarEventTypePickup,
		Title:       "Pickup from Test Corp",
		Date:        now,
		ShipmentID:  123,
		Description: "Scheduled pickup",
	}

	// Test: Title is not empty
	if event.Title == "" {
		t.Error("Expected non-empty title")
	}

	// Test: Type is valid
	if !IsValidCalendarEventType(event.Type) {
		t.Error("Expected valid event type")
	}

	// Test: Date is set
	if event.Date.IsZero() {
		t.Error("Expected non-zero date")
	}

	// Test: GetColorClass returns appropriate color for event type
	colorClass := event.GetColorClass()
	if colorClass == "" {
		t.Error("Expected non-empty color class")
	}
}

// TestCalendarEventTypeValidation tests validation of calendar event types
func TestCalendarEventTypeValidation(t *testing.T) {
	testCases := []struct {
		name      string
		eventType CalendarEventType
		expected  bool
	}{
		{"Valid pickup", CalendarEventTypePickup, true},
		{"Valid delivery", CalendarEventTypeDelivery, true},
		{"Valid in-transit", CalendarEventTypeInTransit, true},
		{"Valid at warehouse", CalendarEventTypeAtWarehouse, true},
		{"Invalid type", CalendarEventType("invalid"), false},
		{"Empty type", CalendarEventType(""), false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := IsValidCalendarEventType(tc.eventType)
			if result != tc.expected {
				t.Errorf("Expected %v for %s, got %v", tc.expected, tc.eventType, result)
			}
		})
	}
}

// Helper function to create a software engineer
func createSoftwareEngineer(db *sql.DB, engineer *SoftwareEngineer) error {
	query := `
		INSERT INTO software_engineers (name, email, address, phone, address_confirmed, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	return db.QueryRow(
		query,
		engineer.Name,
		engineer.Email,
		engineer.Address,
		engineer.Phone,
		engineer.AddressConfirmed,
		time.Now(),
	).Scan(&engineer.ID)
}
