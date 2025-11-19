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
	events, err := GetCalendarEvents(db, now.AddDate(0, 0, -1), now.AddDate(0, 0, 14), nil, nil)
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
	events, err := GetCalendarEvents(db, now, now.AddDate(0, 2, 0), nil, nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(events) == 0 {
		t.Error("Expected to find events within date range")
	}

	// Test: Query outside date range (should not include the event)
	events, err = GetCalendarEvents(db, now.AddDate(0, -2, 0), now.AddDate(0, -1, 0), nil, nil)
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

// TestGetCalendarEventsWithClientCompanyFilter tests filtering events by client company
func TestGetCalendarEventsWithClientCompanyFilter(t *testing.T) {
	db, cleanup := database.SetupTestDB(t)
	defer cleanup()

	// Create two client companies
	clientCompany1 := &ClientCompany{
		Name:        "Company A",
		ContactInfo: "contact@companya.com",
	}
	err := createClientCompany(db, clientCompany1)
	if err != nil {
		t.Fatalf("Failed to create client company 1: %v", err)
	}

	clientCompany2 := &ClientCompany{
		Name:        "Company B",
		ContactInfo: "contact@companyb.com",
	}
	err = createClientCompany(db, clientCompany2)
	if err != nil {
		t.Fatalf("Failed to create client company 2: %v", err)
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

	// Create shipments for both companies
	now := time.Now()
	tomorrow := now.AddDate(0, 0, 1)

	// Shipment for Company A
	shipmentA := &Shipment{
		ClientCompanyID:     clientCompany1.ID,
		SoftwareEngineerID:  &engineer.ID,
		Status:              ShipmentStatusPendingPickup,
		JiraTicketNumber:    "TEST-700",
		PickupScheduledDate: &tomorrow,
	}
	shipmentA.BeforeCreate()
	err = createShipment(db, shipmentA)
	if err != nil {
		t.Fatalf("Failed to create shipment A: %v", err)
	}

	// Shipment for Company B
	shipmentB := &Shipment{
		ClientCompanyID:     clientCompany2.ID,
		SoftwareEngineerID:  &engineer.ID,
		Status:              ShipmentStatusPendingPickup,
		JiraTicketNumber:    "TEST-701",
		PickupScheduledDate: &tomorrow,
	}
	shipmentB.BeforeCreate()
	err = createShipment(db, shipmentB)
	if err != nil {
		t.Fatalf("Failed to create shipment B: %v", err)
	}

	// Test: Get all events without filter (should return both)
	allEvents, err := GetCalendarEvents(db, now.AddDate(0, 0, -1), now.AddDate(0, 0, 14), nil, nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(allEvents) < 2 {
		t.Errorf("Expected at least 2 events without filter, got %d", len(allEvents))
	}

	// Test: Filter by Company A (should only return Company A's events)
	companyAID := clientCompany1.ID
	eventsA, err := GetCalendarEvents(db, now.AddDate(0, 0, -1), now.AddDate(0, 0, 14), &companyAID, nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(eventsA) == 0 {
		t.Error("Expected to find events for Company A")
	}

	// Verify all events belong to Company A
	for _, event := range eventsA {
		if event.ShipmentID != shipmentA.ID {
			t.Errorf("Expected event to belong to Company A shipment, got shipment ID %d", event.ShipmentID)
		}
	}

	// Test: Filter by Company B (should only return Company B's events)
	companyBID := clientCompany2.ID
	eventsB, err := GetCalendarEvents(db, now.AddDate(0, 0, -1), now.AddDate(0, 0, 14), &companyBID, nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(eventsB) == 0 {
		t.Error("Expected to find events for Company B")
	}

	// Verify all events belong to Company B
	for _, event := range eventsB {
		if event.ShipmentID != shipmentB.ID {
			t.Errorf("Expected event to belong to Company B shipment, got shipment ID %d", event.ShipmentID)
		}
	}

	// Test: Filter by non-existent company (should return no events)
	nonExistentID := int64(99999)
	eventsEmpty, err := GetCalendarEvents(db, now.AddDate(0, 0, -1), now.AddDate(0, 0, 14), &nonExistentID, nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(eventsEmpty) != 0 {
		t.Errorf("Expected 0 events for non-existent company, got %d", len(eventsEmpty))
	}
}

// TestGetCalendarEventsWithWarehouseRoleFilter tests filtering events for warehouse users
func TestGetCalendarEventsWithWarehouseRoleFilter(t *testing.T) {
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

	// Create shipments with different statuses
	now := time.Now()
	tomorrow := now.AddDate(0, 0, 1)
	nextWeek := now.AddDate(0, 0, 7)

	// Shipment 1: In transit to warehouse (warehouse users should see)
	shipment1 := &Shipment{
		ClientCompanyID:     clientCompany.ID,
		SoftwareEngineerID:  &engineer.ID,
		Status:              ShipmentStatusInTransitToWarehouse,
		JiraTicketNumber:    "TEST-800",
		PickupScheduledDate: &tomorrow,
		PickedUpAt:          &now,
	}
	shipment1.BeforeCreate()
	err = createShipment(db, shipment1)
	if err != nil {
		t.Fatalf("Failed to create shipment 1: %v", err)
	}

	// Shipment 2: At warehouse (warehouse users should see)
	shipment2 := &Shipment{
		ClientCompanyID:     clientCompany.ID,
		SoftwareEngineerID:  &engineer.ID,
		Status:              ShipmentStatusAtWarehouse,
		JiraTicketNumber:    "TEST-801",
		PickupScheduledDate: &tomorrow,
		ArrivedWarehouseAt:  &now,
	}
	shipment2.BeforeCreate()
	err = createShipment(db, shipment2)
	if err != nil {
		t.Fatalf("Failed to create shipment 2: %v", err)
	}

	// Shipment 3: Released from warehouse (warehouse users should see)
	shipment3 := &Shipment{
		ClientCompanyID:      clientCompany.ID,
		SoftwareEngineerID:   &engineer.ID,
		Status:               ShipmentStatusReleasedFromWarehouse,
		JiraTicketNumber:    "TEST-802",
		PickupScheduledDate:  &tomorrow,
		ReleasedWarehouseAt: &now,
	}
	shipment3.BeforeCreate()
	err = createShipment(db, shipment3)
	if err != nil {
		t.Fatalf("Failed to create shipment 3: %v", err)
	}

	// Shipment 4: Pending pickup (warehouse users should NOT see)
	shipment4 := &Shipment{
		ClientCompanyID:     clientCompany.ID,
		SoftwareEngineerID:  &engineer.ID,
		Status:              ShipmentStatusPendingPickup,
		JiraTicketNumber:    "TEST-803",
		PickupScheduledDate: &tomorrow,
	}
	shipment4.BeforeCreate()
	err = createShipment(db, shipment4)
	if err != nil {
		t.Fatalf("Failed to create shipment 4: %v", err)
	}

	// Shipment 5: Delivered (warehouse users should NOT see)
	shipment5 := &Shipment{
		ClientCompanyID:     clientCompany.ID,
		SoftwareEngineerID:  &engineer.ID,
		Status:              ShipmentStatusDelivered,
		JiraTicketNumber:    "TEST-804",
		PickupScheduledDate: &tomorrow,
		DeliveredAt:         &nextWeek,
	}
	shipment5.BeforeCreate()
	err = createShipment(db, shipment5)
	if err != nil {
		t.Fatalf("Failed to create shipment 5: %v", err)
	}

	// Test: Get all events without role filter (should return all)
	allEvents, err := GetCalendarEvents(db, now.AddDate(0, 0, -1), now.AddDate(0, 0, 14), nil, nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(allEvents) == 0 {
		t.Error("Expected to find events without filter")
	}

	// Test: Filter by warehouse role (should only return warehouse-related shipments)
	warehouseRole := RoleWarehouse
	warehouseEvents, err := GetCalendarEvents(db, now.AddDate(0, 0, -1), now.AddDate(0, 0, 14), nil, &warehouseRole)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(warehouseEvents) == 0 {
		t.Error("Expected to find events for warehouse users")
	}

	// Verify all events belong to warehouse-related shipments
	allowedShipmentIDs := map[int64]bool{
		shipment1.ID: true,
		shipment2.ID: true,
		shipment3.ID: true,
	}
	excludedShipmentIDs := map[int64]bool{
		shipment4.ID: true,
		shipment5.ID: true,
	}

	for _, event := range warehouseEvents {
		if excludedShipmentIDs[event.ShipmentID] {
			t.Errorf("Warehouse user should not see event for shipment %d (status: %s)", event.ShipmentID, shipment4.Status)
		}
		if !allowedShipmentIDs[event.ShipmentID] && !excludedShipmentIDs[event.ShipmentID] {
			// This is fine - might be from other test data
		}
	}

	// Verify at least one event from each allowed shipment is present
	foundShipment1 := false
	foundShipment2 := false
	foundShipment3 := false
	for _, event := range warehouseEvents {
		if event.ShipmentID == shipment1.ID {
			foundShipment1 = true
		}
		if event.ShipmentID == shipment2.ID {
			foundShipment2 = true
		}
		if event.ShipmentID == shipment3.ID {
			foundShipment3 = true
		}
	}

	if !foundShipment1 {
		t.Error("Expected to find events for shipment in transit to warehouse")
	}
	if !foundShipment2 {
		t.Error("Expected to find events for shipment at warehouse")
	}
	if !foundShipment3 {
		t.Error("Expected to find events for shipment released from warehouse")
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
