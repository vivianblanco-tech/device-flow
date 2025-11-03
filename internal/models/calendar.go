package models

import (
	"database/sql"
	"fmt"
	"time"
)

// CalendarEventType represents the type of calendar event
type CalendarEventType string

// Calendar event type constants
const (
	CalendarEventTypePickup      CalendarEventType = "pickup"
	CalendarEventTypeDelivery    CalendarEventType = "delivery"
	CalendarEventTypeInTransit   CalendarEventType = "in_transit"
	CalendarEventTypeAtWarehouse CalendarEventType = "at_warehouse"
)

// CalendarEvent represents an event on the calendar
type CalendarEvent struct {
	ID          int64             `json:"id"`
	Type        CalendarEventType `json:"type"`
	Title       string            `json:"title"`
	Date        time.Time         `json:"date"`
	ShipmentID  int64             `json:"shipment_id"`
	Description string            `json:"description,omitempty"`
}

// IsValidCalendarEventType checks if a given event type is valid
func IsValidCalendarEventType(eventType CalendarEventType) bool {
	switch eventType {
	case CalendarEventTypePickup,
		CalendarEventTypeDelivery,
		CalendarEventTypeInTransit,
		CalendarEventTypeAtWarehouse:
		return true
	}
	return false
}

// GetColorClass returns the CSS color class for the event type
func (e *CalendarEvent) GetColorClass() string {
	switch e.Type {
	case CalendarEventTypePickup:
		return "bg-blue-500"
	case CalendarEventTypeDelivery:
		return "bg-green-500"
	case CalendarEventTypeInTransit:
		return "bg-yellow-500"
	case CalendarEventTypeAtWarehouse:
		return "bg-purple-500"
	default:
		return "bg-gray-500"
	}
}

// GetCalendarEvents retrieves calendar events within a date range
func GetCalendarEvents(db *sql.DB, startDate, endDate time.Time) ([]CalendarEvent, error) {
	query := `
		SELECT 
			s.id,
			s.client_company_id,
			s.pickup_scheduled_date,
			s.picked_up_at,
			s.arrived_warehouse_at,
			s.released_warehouse_at,
			s.delivered_at,
			cc.name as client_name,
			se.name as engineer_name
		FROM shipments s
		LEFT JOIN client_companies cc ON s.client_company_id = cc.id
		LEFT JOIN software_engineers se ON s.software_engineer_id = se.id
		WHERE 
			(s.pickup_scheduled_date BETWEEN $1 AND $2)
			OR (s.picked_up_at BETWEEN $1 AND $2)
			OR (s.arrived_warehouse_at BETWEEN $1 AND $2)
			OR (s.released_warehouse_at BETWEEN $1 AND $2)
			OR (s.delivered_at BETWEEN $1 AND $2)
		ORDER BY 
			COALESCE(s.pickup_scheduled_date, s.picked_up_at, s.arrived_warehouse_at, s.released_warehouse_at, s.delivered_at)
	`

	rows, err := db.Query(query, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query calendar events: %w", err)
	}
	defer rows.Close()

	var events []CalendarEvent
	eventID := int64(1) // Auto-increment ID for events

	for rows.Next() {
		var (
			shipmentID          int64
			clientCompanyID     int64
			pickupScheduledDate sql.NullTime
			pickedUpAt          sql.NullTime
			arrivedWarehouseAt  sql.NullTime
			releasedWarehouseAt sql.NullTime
			deliveredAt         sql.NullTime
			clientName          sql.NullString
			engineerName        sql.NullString
		)

		err := rows.Scan(
			&shipmentID,
			&clientCompanyID,
			&pickupScheduledDate,
			&pickedUpAt,
			&arrivedWarehouseAt,
			&releasedWarehouseAt,
			&deliveredAt,
			&clientName,
			&engineerName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan calendar event: %w", err)
		}

		// Create event for pickup scheduled date
		if pickupScheduledDate.Valid {
			event := CalendarEvent{
				ID:          eventID,
				Type:        CalendarEventTypePickup,
				Title:       fmt.Sprintf("Pickup from %s", getStringOrDefault(clientName, "Unknown Client")),
				Date:        pickupScheduledDate.Time,
				ShipmentID:  shipmentID,
				Description: "Scheduled pickup",
			}
			events = append(events, event)
			eventID++
		}

		// Create event for picked up
		if pickedUpAt.Valid {
			event := CalendarEvent{
				ID:          eventID,
				Type:        CalendarEventTypeInTransit,
				Title:       fmt.Sprintf("Picked up from %s", getStringOrDefault(clientName, "Unknown Client")),
				Date:        pickedUpAt.Time,
				ShipmentID:  shipmentID,
				Description: "In transit to warehouse",
			}
			events = append(events, event)
			eventID++
		}

		// Create event for arrived at warehouse
		if arrivedWarehouseAt.Valid {
			event := CalendarEvent{
				ID:          eventID,
				Type:        CalendarEventTypeAtWarehouse,
				Title:       fmt.Sprintf("Arrived at warehouse"),
				Date:        arrivedWarehouseAt.Time,
				ShipmentID:  shipmentID,
				Description: "Shipment at warehouse",
			}
			events = append(events, event)
			eventID++
		}

		// Create event for released from warehouse
		if releasedWarehouseAt.Valid {
			event := CalendarEvent{
				ID:          eventID,
				Type:        CalendarEventTypeInTransit,
				Title:       fmt.Sprintf("Released to %s", getStringOrDefault(engineerName, "Engineer")),
				Date:        releasedWarehouseAt.Time,
				ShipmentID:  shipmentID,
				Description: "In transit to engineer",
			}
			events = append(events, event)
			eventID++
		}

		// Create event for delivered
		if deliveredAt.Valid {
			event := CalendarEvent{
				ID:          eventID,
				Type:        CalendarEventTypeDelivery,
				Title:       fmt.Sprintf("Delivered to %s", getStringOrDefault(engineerName, "Engineer")),
				Date:        deliveredAt.Time,
				ShipmentID:  shipmentID,
				Description: "Delivery completed",
			}
			events = append(events, event)
			eventID++
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating calendar events: %w", err)
	}

	return events, nil
}

// Helper function to get string value or default
func getStringOrDefault(ns sql.NullString, defaultValue string) string {
	if ns.Valid {
		return ns.String
	}
	return defaultValue
}

