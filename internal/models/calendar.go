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

// GetColorClass returns the CSS background color class for the event type
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

// GetBorderColorClass returns the CSS border color class for the event type
func (e *CalendarEvent) GetBorderColorClass() string {
	switch e.Type {
	case CalendarEventTypePickup:
		return "border-blue-500"
	case CalendarEventTypeDelivery:
		return "border-green-500"
	case CalendarEventTypeInTransit:
		return "border-yellow-500"
	case CalendarEventTypeAtWarehouse:
		return "border-purple-500"
	default:
		return "border-gray-500"
	}
}

// GetTextColorClass returns the CSS text color class for the event type
func (e *CalendarEvent) GetTextColorClass() string {
	switch e.Type {
	case CalendarEventTypePickup:
		return "text-blue-700"
	case CalendarEventTypeDelivery:
		return "text-green-700"
	case CalendarEventTypeInTransit:
		return "text-yellow-700"
	case CalendarEventTypeAtWarehouse:
		return "text-purple-700"
	default:
		return "text-gray-700"
	}
}

// GetShipmentLink returns the URL to the shipment detail page
func (e *CalendarEvent) GetShipmentLink() string {
	if e.ShipmentID == 0 {
		return ""
	}
	return fmt.Sprintf("/shipments/%d", e.ShipmentID)
}

// GetCalendarEvents retrieves calendar events within a date range
// If clientCompanyID is provided (non-nil), only events for that company are returned
// If clientCompanyID is nil, all events are returned
func GetCalendarEvents(db *sql.DB, startDate, endDate time.Time, clientCompanyID *int64) ([]CalendarEvent, error) {
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
			((s.pickup_scheduled_date BETWEEN $1 AND $2)
			OR (s.picked_up_at BETWEEN $1 AND $2)
			OR (s.arrived_warehouse_at BETWEEN $1 AND $2)
			OR (s.released_warehouse_at BETWEEN $1 AND $2)
			OR (s.delivered_at BETWEEN $1 AND $2))
	`

	var rows *sql.Rows
	var err error

	// Add client company filter if provided
	if clientCompanyID != nil {
		query += ` AND s.client_company_id = $3`
		query += `
		ORDER BY 
			COALESCE(s.pickup_scheduled_date, s.picked_up_at, s.arrived_warehouse_at, s.released_warehouse_at, s.delivered_at)
		`
		rows, err = db.Query(query, startDate, endDate, *clientCompanyID)
	} else {
		query += `
		ORDER BY 
			COALESCE(s.pickup_scheduled_date, s.picked_up_at, s.arrived_warehouse_at, s.released_warehouse_at, s.delivered_at)
		`
		rows, err = db.Query(query, startDate, endDate)
	}
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
				Title:       "Arrived at warehouse",
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

// CalendarDay represents a single day in the calendar grid
type CalendarDay struct {
	Date           time.Time       `json:"date"`
	IsCurrentMonth bool            `json:"is_current_month"`
	Events         []CalendarEvent `json:"events"`
}

// GenerateCalendarGrid creates a calendar grid for the given month
// Returns a 2D array where each inner array represents a week (7 days)
func GenerateCalendarGrid(year int, month time.Month) [][]CalendarDay {
	// Get the first day of the month
	firstOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	
	// Get the last day of the month
	lastOfMonth := firstOfMonth.AddDate(0, 1, 0).Add(-24 * time.Hour)
	
	// Calculate the starting Sunday (may be in previous month)
	// If first day is Sunday (0), start there; otherwise go back to previous Sunday
	startDate := firstOfMonth
	for startDate.Weekday() != time.Sunday {
		startDate = startDate.Add(-24 * time.Hour)
	}
	
	// Calculate the ending Saturday (may be in next month)
	endDate := lastOfMonth
	for endDate.Weekday() != time.Saturday {
		endDate = endDate.Add(24 * time.Hour)
	}
	
	// Build the grid
	var grid [][]CalendarDay
	var currentWeek []CalendarDay
	
	currentDate := startDate
	for !currentDate.After(endDate) {
		day := CalendarDay{
			Date:           currentDate,
			IsCurrentMonth: currentDate.Month() == month,
			Events:         []CalendarEvent{},
		}
		currentWeek = append(currentWeek, day)
		
		// If we've completed a week (7 days), add it to the grid
		if len(currentWeek) == 7 {
			grid = append(grid, currentWeek)
			currentWeek = []CalendarDay{}
		}
		
		currentDate = currentDate.Add(24 * time.Hour)
	}
	
	// Add any remaining days (shouldn't happen with proper calculation)
	if len(currentWeek) > 0 {
		grid = append(grid, currentWeek)
	}
	
	return grid
}

// GenerateCalendarGridWithEvents creates a calendar grid and populates it with events
func GenerateCalendarGridWithEvents(year int, month time.Month, events []CalendarEvent) [][]CalendarDay {
	grid := GenerateCalendarGrid(year, month)
	
	// Create a map of date strings to events for faster lookup
	eventsByDate := make(map[string][]CalendarEvent)
	for _, event := range events {
		dateKey := event.Date.Format("2006-01-02")
		eventsByDate[dateKey] = append(eventsByDate[dateKey], event)
	}
	
	// Add events to the appropriate days
	for i := range grid {
		for j := range grid[i] {
			dateKey := grid[i][j].Date.Format("2006-01-02")
			if dayEvents, exists := eventsByDate[dateKey]; exists {
				grid[i][j].Events = dayEvents
			}
		}
	}
	
	return grid
}

