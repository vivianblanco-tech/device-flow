package models

import (
	"testing"
)

// TestGetAllowedStatusesForRole tests filtering statuses by user role
func TestGetAllowedStatusesForRole(t *testing.T) {
	tests := []struct {
		name           string
		role           UserRole
		expectedCount  int
		expectedStatus []LaptopStatus
	}{
		{
			name:          "Warehouse user sees only 3 statuses",
			role:          RoleWarehouse,
			expectedCount: 3,
			expectedStatus: []LaptopStatus{
				LaptopStatusInTransitToWarehouse,
				LaptopStatusAtWarehouse,
				LaptopStatusAvailable,
			},
		},
		{
			name:          "Logistics user sees all statuses",
			role:          RoleLogistics,
			expectedCount: 6, // all statuses
			expectedStatus: []LaptopStatus{
				LaptopStatusAvailable,
				LaptopStatusInTransitToWarehouse,
				LaptopStatusAtWarehouse,
				LaptopStatusInTransitToEngineer,
				LaptopStatusDelivered,
				LaptopStatusRetired,
			},
		},
		{
			name:          "Client user sees all statuses",
			role:          RoleClient,
			expectedCount: 6,
			expectedStatus: []LaptopStatus{
				LaptopStatusAvailable,
				LaptopStatusInTransitToWarehouse,
				LaptopStatusAtWarehouse,
				LaptopStatusInTransitToEngineer,
				LaptopStatusDelivered,
				LaptopStatusRetired,
			},
		},
		{
			name:          "Project Manager user sees all statuses",
			role:          RoleProjectManager,
			expectedCount: 6,
			expectedStatus: []LaptopStatus{
				LaptopStatusAvailable,
				LaptopStatusInTransitToWarehouse,
				LaptopStatusAtWarehouse,
				LaptopStatusInTransitToEngineer,
				LaptopStatusDelivered,
				LaptopStatusRetired,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statuses := GetAllowedStatusesForRole(tt.role)

			// Check count
			if len(statuses) != tt.expectedCount {
				t.Errorf("Expected %d statuses for role %s, got %d", tt.expectedCount, tt.role, len(statuses))
			}

			// Check if all expected statuses are present
			statusMap := make(map[LaptopStatus]bool)
			for _, status := range statuses {
				statusMap[status] = true
			}

			for _, expectedStatus := range tt.expectedStatus {
				if !statusMap[expectedStatus] {
					t.Errorf("Expected status %s to be present for role %s, but it was not found", expectedStatus, tt.role)
				}
			}
		})
	}
}

// TestGetAllowedStatusesForRole_WarehouseExcludesOthers tests warehouse users don't see restricted statuses
func TestGetAllowedStatusesForRole_WarehouseExcludesOthers(t *testing.T) {
	statuses := GetAllowedStatusesForRole(RoleWarehouse)

	// Convert to map for easy lookup
	statusMap := make(map[LaptopStatus]bool)
	for _, status := range statuses {
		statusMap[status] = true
	}

	// Statuses that warehouse users should NOT see
	restrictedStatuses := []LaptopStatus{
		LaptopStatusInTransitToEngineer,
		LaptopStatusDelivered,
		LaptopStatusRetired,
	}

	for _, restrictedStatus := range restrictedStatuses {
		if statusMap[restrictedStatus] {
			t.Errorf("Warehouse user should not see status %s, but it was present", restrictedStatus)
		}
	}
}

// TestGetAllowedStatusesForRole_WarehouseIncludesRequired tests warehouse users see all required statuses
func TestGetAllowedStatusesForRole_WarehouseIncludesRequired(t *testing.T) {
	statuses := GetAllowedStatusesForRole(RoleWarehouse)

	// Convert to map for easy lookup
	statusMap := make(map[LaptopStatus]bool)
	for _, status := range statuses {
		statusMap[status] = true
	}

	// Statuses that warehouse users MUST see
	requiredStatuses := []LaptopStatus{
		LaptopStatusInTransitToWarehouse,
		LaptopStatusAtWarehouse,
		LaptopStatusAvailable,
	}

	for _, requiredStatus := range requiredStatuses {
		if !statusMap[requiredStatus] {
			t.Errorf("Warehouse user must see status %s, but it was not present", requiredStatus)
		}
	}
}

