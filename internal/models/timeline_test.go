package models

import (
	"testing"
	"time"
)

func TestBuildTimeline(t *testing.T) {
	t.Run("timeline for shipment in pending pickup status", func(t *testing.T) {
		now := time.Now()
		createdAt := now.AddDate(0, 0, -5) // Created 5 days ago
		shipment := Shipment{
			Status:    ShipmentStatusPendingPickup,
			CreatedAt: createdAt,
		}

		timeline := BuildTimeline(&shipment)

		if len(timeline) != 8 {
			t.Errorf("Expected 8 timeline items, got %d", len(timeline))
		}

		// First item should be current (pending pickup)
		if !timeline[0].IsCurrent {
			t.Error("First item should be current")
		}

		// First item (pending pickup) should have CreatedAt timestamp
		if timeline[0].Timestamp == nil {
			t.Error("Pending pickup should have CreatedAt timestamp")
		}
		if timeline[0].Timestamp != nil && !timeline[0].Timestamp.Equal(createdAt) {
			t.Errorf("Pending pickup timestamp should be CreatedAt, expected %v, got %v", createdAt, timeline[0].Timestamp)
		}

		// All other items should be pending
		for i := 1; i < len(timeline); i++ {
			if !timeline[i].IsPending {
				t.Errorf("Item %d should be pending", i)
			}
		}
	})

	t.Run("timeline for shipment in transit to warehouse", func(t *testing.T) {
		now := time.Now()
		pickupDate := now.AddDate(0, 0, -2)
		pickedUpAt := now.AddDate(0, 0, -1)

		shipment := Shipment{
			Status:              ShipmentStatusInTransitToWarehouse,
			PickupScheduledDate: &pickupDate,
			PickedUpAt:          &pickedUpAt,
		}

		timeline := BuildTimeline(&shipment)

		// First three items should be completed (pending pickup, pickup scheduled, picked up)
		if !timeline[0].IsCompleted {
			t.Error("Pending pickup should be completed")
		}
		if !timeline[1].IsCompleted {
			t.Error("Pickup scheduled should be completed")
		}
		if !timeline[2].IsCompleted {
			t.Error("Picked up should be completed")
		}

		// Fourth item (in transit to warehouse) should be current and marked as transit
		if !timeline[3].IsCurrent {
			t.Error("In transit to warehouse should be current")
		}
		if !timeline[3].IsTransit {
			t.Error("In transit to warehouse should be marked as transit")
		}
		// In transit to warehouse should have PickedUpAt timestamp
		if timeline[3].Timestamp == nil {
			t.Error("In transit to warehouse should have PickedUpAt timestamp")
		}
		if timeline[3].Timestamp != nil && !timeline[3].Timestamp.Equal(pickedUpAt) {
			t.Errorf("In transit to warehouse timestamp should be PickedUpAt, expected %v, got %v", pickedUpAt, timeline[3].Timestamp)
		}

		// Remaining items should be pending
		for i := 4; i < len(timeline); i++ {
			if !timeline[i].IsPending {
				t.Errorf("Item %d should be pending", i)
			}
		}
	})

	t.Run("timeline for delivered shipment", func(t *testing.T) {
		now := time.Now()
		pickupDate := now.AddDate(0, 0, -10)
		pickedUpAt := now.AddDate(0, 0, -9)
		arrivedWarehouseAt := now.AddDate(0, 0, -8)
		releasedWarehouseAt := now.AddDate(0, 0, -7)
		deliveredAt := now.AddDate(0, 0, -1)

		shipment := Shipment{
			Status:              ShipmentStatusDelivered,
			PickupScheduledDate: &pickupDate,
			PickedUpAt:          &pickedUpAt,
			ArrivedWarehouseAt:  &arrivedWarehouseAt,
			ReleasedWarehouseAt: &releasedWarehouseAt,
			DeliveredAt:         &deliveredAt,
		}

		timeline := BuildTimeline(&shipment)

		// All items with timestamps should be completed
		completedCount := 0
		for _, item := range timeline {
			if item.IsCompleted {
				completedCount++
			}
		}

		// Should have 5 completed items (all except the two transit statuses)
		if completedCount < 5 {
			t.Errorf("Expected at least 5 completed items, got %d", completedCount)
		}

		// Last item should be current (delivered)
		if !timeline[len(timeline)-1].IsCurrent {
			t.Error("Last item (delivered) should be current")
		}
	})

	t.Run("timeline includes in transit to engineer status", func(t *testing.T) {
		now := time.Now()
		releasedAt := now.AddDate(0, 0, -1)

		shipment := Shipment{
			Status:              ShipmentStatusInTransitToEngineer,
			ReleasedWarehouseAt: &releasedAt,
		}

		timeline := BuildTimeline(&shipment)

		// Find the "In Transit to Engineer" item
		var transitItem *TimelineItem
		for i := range timeline {
			if timeline[i].Status == ShipmentStatusInTransitToEngineer {
				transitItem = &timeline[i]
				break
			}
		}

		if transitItem == nil {
			t.Fatal("Timeline should include 'In Transit to Engineer' status")
		}

		if !transitItem.IsCurrent {
			t.Error("In Transit to Engineer should be current")
		}

		if !transitItem.IsTransit {
			t.Error("In Transit to Engineer should be marked as transit")
		}

		// In transit to engineer should have ReleasedWarehouseAt timestamp
		if transitItem.Timestamp == nil {
			t.Error("In transit to engineer should have ReleasedWarehouseAt timestamp")
		}
		if transitItem.Timestamp != nil && !transitItem.Timestamp.Equal(releasedAt) {
			t.Errorf("In transit to engineer timestamp should be ReleasedWarehouseAt, expected %v, got %v", releasedAt, transitItem.Timestamp)
		}
	})

	t.Run("timeline labels are correct", func(t *testing.T) {
		shipment := Shipment{
			Status: ShipmentStatusPendingPickup,
		}

		timeline := BuildTimeline(&shipment)

		expectedLabels := []string{
			"Pending Pickup",
			"Pickup Scheduled",
			"Picked Up from Client",
			"In Transit to Warehouse",
			"Arrived at Warehouse",
			"Released from Warehouse",
			"In Transit to Engineer",
			"Delivered Successfully",
		}

		for i, expected := range expectedLabels {
			if timeline[i].Label != expected {
				t.Errorf("Item %d: expected label '%s', got '%s'", i, expected, timeline[i].Label)
			}
		}
	})

	t.Run("timeline for single_full_journey shipment shows full timeline", func(t *testing.T) {
		now := time.Now()
		shipment := Shipment{
			ShipmentType:        ShipmentTypeSingleFullJourney,
			Status:              ShipmentStatusPendingPickup,
			PickupScheduledDate: &now,
		}

		timeline := BuildTimeline(&shipment)

		// Should have all 8 statuses from Pending Pickup to Delivered
		if len(timeline) != 8 {
			t.Errorf("Single full journey should have 8 timeline items, got %d", len(timeline))
		}

		// Verify first and last statuses
		if timeline[0].Status != ShipmentStatusPendingPickup {
			t.Errorf("First status should be Pending Pickup, got %s", timeline[0].Status)
		}
		if timeline[len(timeline)-1].Status != ShipmentStatusDelivered {
			t.Errorf("Last status should be Delivered, got %s", timeline[len(timeline)-1].Status)
		}
	})

	t.Run("timeline for bulk_to_warehouse shipment shows only pickup to warehouse", func(t *testing.T) {
		now := time.Now()
		shipment := Shipment{
			ShipmentType:        ShipmentTypeBulkToWarehouse,
			Status:              ShipmentStatusPendingPickup,
			PickupScheduledDate: &now,
		}

		timeline := BuildTimeline(&shipment)

		// Should have only 5 statuses from Pending Pickup to At Warehouse
		if len(timeline) != 5 {
			t.Errorf("Bulk to warehouse should have 5 timeline items, got %d", len(timeline))
		}

		// Verify first and last statuses
		if timeline[0].Status != ShipmentStatusPendingPickup {
			t.Errorf("First status should be Pending Pickup, got %s", timeline[0].Status)
		}
		if timeline[len(timeline)-1].Status != ShipmentStatusAtWarehouse {
			t.Errorf("Last status should be At Warehouse, got %s", timeline[len(timeline)-1].Status)
		}

		// Verify it doesn't include warehouse release or delivery statuses
		for _, item := range timeline {
			if item.Status == ShipmentStatusReleasedFromWarehouse ||
				item.Status == ShipmentStatusInTransitToEngineer ||
				item.Status == ShipmentStatusDelivered {
				t.Errorf("Bulk to warehouse timeline should not include status %s", item.Status)
			}
		}
	})

	t.Run("timeline for warehouse_to_engineer shipment shows only warehouse to delivery", func(t *testing.T) {
		now := time.Now()
		shipment := Shipment{
			ShipmentType:        ShipmentTypeWarehouseToEngineer,
			Status:              ShipmentStatusReleasedFromWarehouse,
			ReleasedWarehouseAt: &now,
		}

		timeline := BuildTimeline(&shipment)

		// Should have only 3 statuses from Released from Warehouse to Delivered
		if len(timeline) != 3 {
			t.Errorf("Warehouse to engineer should have 3 timeline items, got %d", len(timeline))
		}

		// Verify first and last statuses
		if timeline[0].Status != ShipmentStatusReleasedFromWarehouse {
			t.Errorf("First status should be Released from Warehouse, got %s", timeline[0].Status)
		}
		if timeline[len(timeline)-1].Status != ShipmentStatusDelivered {
			t.Errorf("Last status should be Delivered, got %s", timeline[len(timeline)-1].Status)
		}

		// Verify it doesn't include client pickup statuses
		for _, item := range timeline {
			if item.Status == ShipmentStatusPendingPickup ||
				item.Status == ShipmentStatusPickupScheduled ||
				item.Status == ShipmentStatusPickedUpFromClient ||
				item.Status == ShipmentStatusInTransitToWarehouse ||
				item.Status == ShipmentStatusAtWarehouse {
				t.Errorf("Warehouse to engineer timeline should not include status %s", item.Status)
			}
		}
	})

	t.Run("timeline shows CreatedAt timestamp for completed pending pickup status", func(t *testing.T) {
		now := time.Now()
		createdAt := now.AddDate(0, 0, -10)
		pickupDate := now.AddDate(0, 0, -2)
		pickedUpAt := now.AddDate(0, 0, -1)

		shipment := Shipment{
			Status:              ShipmentStatusPickedUpFromClient,
			CreatedAt:           createdAt,
			PickupScheduledDate: &pickupDate,
			PickedUpAt:          &pickedUpAt,
		}

		timeline := BuildTimeline(&shipment)

		// Find the "Pending Pickup" item (should be completed)
		var pendingPickupItem *TimelineItem
		for i := range timeline {
			if timeline[i].Status == ShipmentStatusPendingPickup {
				pendingPickupItem = &timeline[i]
				break
			}
		}

		if pendingPickupItem == nil {
			t.Fatal("Timeline should include 'Pending Pickup' status")
		}

		if !pendingPickupItem.IsCompleted {
			t.Error("Pending Pickup should be completed when shipment is picked up")
		}

		// Completed pending pickup should have CreatedAt timestamp
		if pendingPickupItem.Timestamp == nil {
			t.Error("Completed pending pickup should have CreatedAt timestamp")
		}
		if pendingPickupItem.Timestamp != nil && !pendingPickupItem.Timestamp.Equal(createdAt) {
			t.Errorf("Pending pickup timestamp should be CreatedAt, expected %v, got %v", createdAt, pendingPickupItem.Timestamp)
		}
	})

	t.Run("timeline shows PickedUpAt timestamp for completed in transit to warehouse status", func(t *testing.T) {
		now := time.Now()
		pickupDate := now.AddDate(0, 0, -3)
		pickedUpAt := now.AddDate(0, 0, -2)
		arrivedWarehouseAt := now.AddDate(0, 0, -1)

		shipment := Shipment{
			Status:             ShipmentStatusAtWarehouse,
			PickupScheduledDate: &pickupDate,
			PickedUpAt:         &pickedUpAt,
			ArrivedWarehouseAt: &arrivedWarehouseAt,
		}

		timeline := BuildTimeline(&shipment)

		// Find the "In Transit to Warehouse" item (should be completed)
		var transitItem *TimelineItem
		for i := range timeline {
			if timeline[i].Status == ShipmentStatusInTransitToWarehouse {
				transitItem = &timeline[i]
				break
			}
		}

		if transitItem == nil {
			t.Fatal("Timeline should include 'In Transit to Warehouse' status")
		}

		if !transitItem.IsCompleted {
			t.Error("In Transit to Warehouse should be completed when shipment is at warehouse")
		}

		// Completed in transit to warehouse should have PickedUpAt timestamp
		if transitItem.Timestamp == nil {
			t.Error("Completed in transit to warehouse should have PickedUpAt timestamp")
		}
		if transitItem.Timestamp != nil && !transitItem.Timestamp.Equal(pickedUpAt) {
			t.Errorf("In transit to warehouse timestamp should be PickedUpAt, expected %v, got %v", pickedUpAt, transitItem.Timestamp)
		}
	})
}
