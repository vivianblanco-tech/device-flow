package models

import (
	"testing"
	"time"
)

// TestCalendarEvent_GetColorClass tests that each event type returns the correct color class
func TestCalendarEvent_GetColorClass(t *testing.T) {
	testCases := []struct {
		name           string
		eventType      CalendarEventType
		expectedClass  string
	}{
		{
			name:          "Pickup event has blue color",
			eventType:     CalendarEventTypePickup,
			expectedClass: "bg-blue-500",
		},
		{
			name:          "Delivery event has green color",
			eventType:     CalendarEventTypeDelivery,
			expectedClass: "bg-green-500",
		},
		{
			name:          "In transit event has yellow color",
			eventType:     CalendarEventTypeInTransit,
			expectedClass: "bg-yellow-500",
		},
		{
			name:          "At warehouse event has purple color",
			eventType:     CalendarEventTypeAtWarehouse,
			expectedClass: "bg-purple-500",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			event := CalendarEvent{
				ID:         1,
				Type:       tc.eventType,
				Title:      "Test Event",
				Date:       time.Now(),
				ShipmentID: 100,
			}

			colorClass := event.GetColorClass()
			if colorClass != tc.expectedClass {
				t.Errorf("Expected color class %s, got %s", tc.expectedClass, colorClass)
			}
		})
	}
}

// TestCalendarEvent_GetBorderColorClass tests that events return the correct border color class
func TestCalendarEvent_GetBorderColorClass(t *testing.T) {
	testCases := []struct {
		name           string
		eventType      CalendarEventType
		expectedClass  string
	}{
		{
			name:          "Pickup event has blue border",
			eventType:     CalendarEventTypePickup,
			expectedClass: "border-blue-500",
		},
		{
			name:          "Delivery event has green border",
			eventType:     CalendarEventTypeDelivery,
			expectedClass: "border-green-500",
		},
		{
			name:          "In transit event has yellow border",
			eventType:     CalendarEventTypeInTransit,
			expectedClass: "border-yellow-500",
		},
		{
			name:          "At warehouse event has purple border",
			eventType:     CalendarEventTypeAtWarehouse,
			expectedClass: "border-purple-500",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			event := CalendarEvent{
				ID:         1,
				Type:       tc.eventType,
				Title:      "Test Event",
				Date:       time.Now(),
				ShipmentID: 100,
			}

			borderClass := event.GetBorderColorClass()
			if borderClass != tc.expectedClass {
				t.Errorf("Expected border class %s, got %s", tc.expectedClass, borderClass)
			}
		})
	}
}

// TestCalendarEvent_GetTextColorClass tests that events return the correct text color class
func TestCalendarEvent_GetTextColorClass(t *testing.T) {
	testCases := []struct {
		name           string
		eventType      CalendarEventType
		expectedClass  string
	}{
		{
			name:          "Pickup event has blue text",
			eventType:     CalendarEventTypePickup,
			expectedClass: "text-blue-700",
		},
		{
			name:          "Delivery event has green text",
			eventType:     CalendarEventTypeDelivery,
			expectedClass: "text-green-700",
		},
		{
			name:          "In transit event has yellow text",
			eventType:     CalendarEventTypeInTransit,
			expectedClass: "text-yellow-700",
		},
		{
			name:          "At warehouse event has purple text",
			eventType:     CalendarEventTypeAtWarehouse,
			expectedClass: "text-purple-700",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			event := CalendarEvent{
				ID:         1,
				Type:       tc.eventType,
				Title:      "Test Event",
				Date:       time.Now(),
				ShipmentID: 100,
			}

			textClass := event.GetTextColorClass()
			if textClass != tc.expectedClass {
				t.Errorf("Expected text class %s, got %s", tc.expectedClass, textClass)
			}
		})
	}
}

// TestCalendarEvent_GetShipmentLink tests that events generate the correct shipment link
func TestCalendarEvent_GetShipmentLink(t *testing.T) {
	event := CalendarEvent{
		ID:         1,
		Type:       CalendarEventTypePickup,
		Title:      "Pickup from Test Corp",
		Date:       time.Now(),
		ShipmentID: 12345,
	}

	expectedLink := "/shipments/12345"
	link := event.GetShipmentLink()

	if link != expectedLink {
		t.Errorf("Expected link %s, got %s", expectedLink, link)
	}
}

// TestCalendarEvent_GetShipmentLinkWithZeroID tests that events with zero ID return empty link
func TestCalendarEvent_GetShipmentLinkWithZeroID(t *testing.T) {
	event := CalendarEvent{
		ID:         1,
		Type:       CalendarEventTypePickup,
		Title:      "Pickup from Test Corp",
		Date:       time.Now(),
		ShipmentID: 0,
	}

	link := event.GetShipmentLink()

	if link != "" {
		t.Errorf("Expected empty link for zero shipment ID, got %s", link)
	}
}

