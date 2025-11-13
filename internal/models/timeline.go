package models

import "time"

// TimelineItem represents a single item in a shipment's tracking timeline
type TimelineItem struct {
	Label          string     // Display label (e.g., "Pickup Scheduled")
	Status         ShipmentStatus // The actual status value
	Timestamp      *time.Time // When this status was reached (nil if not reached yet)
	IsCompleted    bool       // Whether this status has been completed
	IsCurrent      bool       // Whether this is the current status
	IsPending      bool       // Whether this status is yet to be reached
	IsTransit      bool       // Whether this is a transit status (for special coloring)
	Icon           string     // Icon/emoji for the status
	TrackingNumber string     // Tracking number associated with this status (if applicable)
}

// BuildTimeline creates a complete timeline from the shipment's current state
// The timeline is filtered based on shipment type:
// - single_full_journey: Full timeline from pickup to delivery
// - bulk_to_warehouse: Only pickup to warehouse arrival
// - warehouse_to_engineer: Only warehouse release to delivery
func BuildTimeline(s *Shipment) []TimelineItem {
	// All possible statuses in order
	allStatuses := []struct {
		Status    ShipmentStatus
		Label     string
		Icon      string
		IsTransit bool
		GetTime   func(*Shipment) *time.Time
	}{
		{
			Status:  ShipmentStatusPendingPickup,
			Label:   "Pending Pickup",
			Icon:    "clock",
			GetTime: func(s *Shipment) *time.Time { return nil },
		},
		{
			Status:  ShipmentStatusPickupScheduled,
			Label:   "Pickup Scheduled",
			Icon:    "calendar",
			GetTime: func(s *Shipment) *time.Time { return s.PickupScheduledDate },
		},
		{
			Status:  ShipmentStatusPickedUpFromClient,
			Label:   "Picked Up from Client",
			Icon:    "check",
			GetTime: func(s *Shipment) *time.Time { return s.PickedUpAt },
		},
		{
			Status:    ShipmentStatusInTransitToWarehouse,
			Label:     "In Transit to Warehouse",
			Icon:      "truck",
			IsTransit: true,
			GetTime:   func(s *Shipment) *time.Time { return nil }, // No specific timestamp field
		},
		{
			Status:  ShipmentStatusAtWarehouse,
			Label:   "Arrived at Warehouse",
			Icon:    "home",
			GetTime: func(s *Shipment) *time.Time { return s.ArrivedWarehouseAt },
		},
		{
			Status:  ShipmentStatusReleasedFromWarehouse,
			Label:   "Released from Warehouse",
			Icon:    "truck",
			GetTime: func(s *Shipment) *time.Time { return s.ReleasedWarehouseAt },
		},
		{
			Status:    ShipmentStatusInTransitToEngineer,
			Label:     "In Transit to Engineer",
			Icon:      "truck",
			IsTransit: true,
			GetTime:   func(s *Shipment) *time.Time { return nil }, // No specific timestamp field
		},
		{
			Status:  ShipmentStatusDelivered,
			Label:   "Delivered Successfully",
			Icon:    "badge",
			GetTime: func(s *Shipment) *time.Time { return s.DeliveredAt },
		},
	}

	// Filter statuses based on shipment type
	var filteredStatuses []struct {
		Status    ShipmentStatus
		Label     string
		Icon      string
		IsTransit bool
		GetTime   func(*Shipment) *time.Time
	}

	switch s.ShipmentType {
	case ShipmentTypeBulkToWarehouse:
		// Only show pickup to warehouse arrival (first 5 statuses)
		filteredStatuses = allStatuses[:5]
	case ShipmentTypeWarehouseToEngineer:
		// Only show warehouse release to delivery (last 3 statuses)
		filteredStatuses = allStatuses[5:]
	default:
		// Single full journey or unspecified: show all statuses
		filteredStatuses = allStatuses
	}

	// Find the index of the current status in the filtered list
	currentStatusIndex := -1
	for i, statusInfo := range filteredStatuses {
		if statusInfo.Status == s.Status {
			currentStatusIndex = i
			break
		}
	}

	// Build timeline items
	timeline := make([]TimelineItem, 0, len(filteredStatuses))
	for i, statusInfo := range filteredStatuses {
		timestamp := statusInfo.GetTime(s)
		
		item := TimelineItem{
			Label:       statusInfo.Label,
			Status:      statusInfo.Status,
			Timestamp:   timestamp,
			IsTransit:   statusInfo.IsTransit,
			Icon:        statusInfo.Icon,
			IsCompleted: i < currentStatusIndex || (i == currentStatusIndex && timestamp != nil),
			IsCurrent:   i == currentStatusIndex,
			IsPending:   i > currentStatusIndex,
		}

		// Add tracking number for pickup scheduled status
		if statusInfo.Status == ShipmentStatusPickupScheduled && s.TrackingNumber != "" {
			item.TrackingNumber = s.TrackingNumber
		}

		// For in-transit statuses that are current but have no timestamp,
		// mark them as current (not completed)
		if i == currentStatusIndex && timestamp == nil {
			item.IsCompleted = false
			item.IsCurrent = true
		}

		timeline = append(timeline, item)
	}

	return timeline
}

