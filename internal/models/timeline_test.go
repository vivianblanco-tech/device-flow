package models

import (
	"testing"
	"time"
)

func TestBuildTimeline(t *testing.T) {
	t.Run("timeline for shipment in pending pickup status", func(t *testing.T) {
		now := time.Now()
		shipment := Shipment{
			Status:              ShipmentStatusPendingPickup,
			PickupScheduledDate: &now,
		}

		timeline := BuildTimeline(&shipment)

		if len(timeline) != 8 {
			t.Errorf("Expected 8 timeline items, got %d", len(timeline))
		}

		// First item should be current (pending pickup)
		if !timeline[0].IsCurrent {
			t.Error("First item should be current")
		}

		// First item (pending pickup) should not have timestamp
		if timeline[0].Timestamp != nil {
			t.Error("Pending pickup should not have timestamp")
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
}

