package models

import (
	"testing"
)

// TestGetLaptopStatusDisplayName tests the display name conversion for laptop statuses
func TestGetLaptopStatusDisplayName(t *testing.T) {
	tests := []struct {
		name     string
		status   LaptopStatus
		expected string
	}{
		{
			name:     "available shows as Available at Warehouse",
			status:   LaptopStatusAvailable,
			expected: "Available at Warehouse",
		},
		{
			name:     "at_warehouse shows as Received at Warehouse",
			status:   LaptopStatusAtWarehouse,
			expected: "Received at Warehouse",
		},
		{
			name:     "in_transit_to_warehouse remains as In Transit To Warehouse",
			status:   LaptopStatusInTransitToWarehouse,
			expected: "In Transit To Warehouse",
		},
		{
			name:     "in_transit_to_engineer remains as In Transit To Engineer",
			status:   LaptopStatusInTransitToEngineer,
			expected: "In Transit To Engineer",
		},
		{
			name:     "delivered remains as Delivered",
			status:   LaptopStatusDelivered,
			expected: "Delivered",
		},
		{
			name:     "retired remains as Retired",
			status:   LaptopStatusRetired,
			expected: "Retired",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetLaptopStatusDisplayName(tt.status)
			if got != tt.expected {
				t.Errorf("GetLaptopStatusDisplayName(%v) = %v, want %v", tt.status, got, tt.expected)
			}
		})
	}
}

// TestGetLaptopStatusesInOrder tests that statuses are returned in logical order
func TestGetLaptopStatusesInOrder(t *testing.T) {
	expectedOrder := []LaptopStatus{
		LaptopStatusInTransitToWarehouse,
		LaptopStatusAtWarehouse, // "Received at Warehouse"
		LaptopStatusAvailable,   // "Available at Warehouse"
		LaptopStatusInTransitToEngineer,
		LaptopStatusDelivered,
		LaptopStatusRetired,
	}

	got := GetLaptopStatusesInOrder()

	if len(got) != len(expectedOrder) {
		t.Fatalf("GetLaptopStatusesInOrder() returned %d statuses, want %d", len(got), len(expectedOrder))
	}

	for i, status := range expectedOrder {
		if got[i] != status {
			t.Errorf("GetLaptopStatusesInOrder()[%d] = %v, want %v", i, got[i], status)
		}
	}
}

// TestGetLaptopStatusesForNewLaptop tests that only appropriate statuses are shown when adding a new laptop
func TestGetLaptopStatusesForNewLaptop(t *testing.T) {
	// Only "Received at Warehouse" should be available for new laptops
	// This ensures warehouse users must create a reception report for every laptop
	// before it becomes "Available at Warehouse"
	expectedStatuses := []LaptopStatus{
		LaptopStatusAtWarehouse, // "Received at Warehouse" - the only status for new laptops
	}

	got := GetLaptopStatusesForNewLaptop()

	if len(got) != len(expectedStatuses) {
		t.Fatalf("GetLaptopStatusesForNewLaptop() returned %d statuses, want %d", len(got), len(expectedStatuses))
	}

	for i, status := range expectedStatuses {
		if got[i] != status {
			t.Errorf("GetLaptopStatusesForNewLaptop()[%d] = %v, want %v", i, got[i], status)
		}
	}

	// Explicitly verify that "Available at Warehouse" is NOT included
	for _, status := range got {
		if status == LaptopStatusAvailable {
			t.Errorf("GetLaptopStatusesForNewLaptop() should not include 'Available at Warehouse' status, but it does")
		}
	}
}
