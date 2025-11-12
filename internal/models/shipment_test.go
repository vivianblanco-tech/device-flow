package models

import (
	"testing"
	"time"
)

func TestShipment_Validate(t *testing.T) {
	tests := []struct {
		name     string
		shipment Shipment
		wantErr  bool
		errMsg   string
	}{
		{
			name: "valid shipment with all fields",
			shipment: Shipment{
				ClientCompanyID:     1,
				SoftwareEngineerID:  int64Ptr(10),
				Status:              ShipmentStatusPendingPickup,
				JiraTicketNumber:    "SCOP-67702",
				CourierName:         "FedEx",
				TrackingNumber:      "TRACK123456",
				PickupScheduledDate: timePtr(time.Now().Add(24 * time.Hour)),
			},
			wantErr: false,
		},
		{
			name: "valid shipment with minimal fields and JIRA ticket",
			shipment: Shipment{
				ClientCompanyID:  2,
				Status:           ShipmentStatusPendingPickup,
				JiraTicketNumber: "PROJECT-12345",
			},
			wantErr: false,
		},
		{
			name: "invalid - missing JIRA ticket number",
			shipment: Shipment{
				ClientCompanyID: 1,
				Status:          ShipmentStatusPendingPickup,
			},
			wantErr: true,
			errMsg:  "JIRA ticket number is required",
		},
		{
			name: "invalid - empty JIRA ticket number",
			shipment: Shipment{
				ClientCompanyID:  1,
				Status:           ShipmentStatusPendingPickup,
				JiraTicketNumber: "",
			},
			wantErr: true,
			errMsg:  "JIRA ticket number is required",
		},
		{
			name: "invalid - malformed JIRA ticket (missing dash)",
			shipment: Shipment{
				ClientCompanyID:  1,
				Status:           ShipmentStatusPendingPickup,
				JiraTicketNumber: "SCOP67702",
			},
			wantErr: true,
			errMsg:  "JIRA ticket number must be in format PROJECT-NUMBER (e.g., SCOP-67702)",
		},
		{
			name: "invalid - malformed JIRA ticket (no project key)",
			shipment: Shipment{
				ClientCompanyID:  1,
				Status:           ShipmentStatusPendingPickup,
				JiraTicketNumber: "-67702",
			},
			wantErr: true,
			errMsg:  "JIRA ticket number must be in format PROJECT-NUMBER (e.g., SCOP-67702)",
		},
		{
			name: "invalid - malformed JIRA ticket (no number)",
			shipment: Shipment{
				ClientCompanyID:  1,
				Status:           ShipmentStatusPendingPickup,
				JiraTicketNumber: "SCOP-",
			},
			wantErr: true,
			errMsg:  "JIRA ticket number must be in format PROJECT-NUMBER (e.g., SCOP-67702)",
		},
		{
			name: "invalid - malformed JIRA ticket (lowercase)",
			shipment: Shipment{
				ClientCompanyID:  1,
				Status:           ShipmentStatusPendingPickup,
				JiraTicketNumber: "scop-67702",
			},
			wantErr: true,
			errMsg:  "JIRA ticket number must be in format PROJECT-NUMBER (e.g., SCOP-67702)",
		},
		{
			name: "invalid - malformed JIRA ticket (special characters)",
			shipment: Shipment{
				ClientCompanyID:  1,
				Status:           ShipmentStatusPendingPickup,
				JiraTicketNumber: "SC@P-67702",
			},
			wantErr: true,
			errMsg:  "JIRA ticket number must be in format PROJECT-NUMBER (e.g., SCOP-67702)",
		},
		{
			name: "invalid - missing client company ID",
			shipment: Shipment{
				Status:           ShipmentStatusPendingPickup,
				JiraTicketNumber: "SCOP-67702",
			},
			wantErr: true,
			errMsg:  "client company ID is required",
		},
		{
			name: "invalid - missing status",
			shipment: Shipment{
				ClientCompanyID:  1,
				JiraTicketNumber: "SCOP-67702",
			},
			wantErr: true,
			errMsg:  "status is required",
		},
		{
			name: "invalid - invalid status",
			shipment: Shipment{
				ClientCompanyID:  1,
				Status:           "invalid_status",
				JiraTicketNumber: "SCOP-67702",
			},
			wantErr: true,
			errMsg:  "invalid status",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.shipment.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Shipment.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" && err.Error() != tt.errMsg {
				t.Errorf("Shipment.Validate() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestShipment_IsValidStatus(t *testing.T) {
	tests := []struct {
		name   string
		status ShipmentStatus
		want   bool
	}{
		{"pending_pickup_from_client", ShipmentStatusPendingPickup, true},
		{"pickup_from_client_scheduled", ShipmentStatusPickupScheduled, true},
		{"picked_up_from_client", ShipmentStatusPickedUpFromClient, true},
		{"in_transit_to_warehouse", ShipmentStatusInTransitToWarehouse, true},
		{"at_warehouse", ShipmentStatusAtWarehouse, true},
		{"released_from_warehouse", ShipmentStatusReleasedFromWarehouse, true},
		{"in_transit_to_engineer", ShipmentStatusInTransitToEngineer, true},
		{"delivered", ShipmentStatusDelivered, true},
		{"invalid status", "unknown", false},
		{"empty status", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidShipmentStatus(tt.status); got != tt.want {
				t.Errorf("IsValidShipmentStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShipment_TableName(t *testing.T) {
	shipment := Shipment{}
	expected := "shipments"
	if got := shipment.TableName(); got != expected {
		t.Errorf("Shipment.TableName() = %v, want %v", got, expected)
	}
}

func TestShipment_BeforeCreate(t *testing.T) {
	shipment := &Shipment{
		ClientCompanyID: 1,
		Status:          ShipmentStatusPendingPickup,
	}

	shipment.BeforeCreate()

	// Check that timestamps are set
	if shipment.CreatedAt.IsZero() {
		t.Error("Shipment.BeforeCreate() did not set CreatedAt")
	}
	if shipment.UpdatedAt.IsZero() {
		t.Error("Shipment.BeforeCreate() did not set UpdatedAt")
	}
}

func TestShipment_BeforeUpdate(t *testing.T) {
	shipment := &Shipment{
		ClientCompanyID: 1,
		Status:          ShipmentStatusPendingPickup,
		CreatedAt:       time.Now().Add(-24 * time.Hour),
		UpdatedAt:       time.Now().Add(-24 * time.Hour),
	}

	oldUpdatedAt := shipment.UpdatedAt
	time.Sleep(10 * time.Millisecond)

	shipment.BeforeUpdate()

	// Check that UpdatedAt was updated
	if !shipment.UpdatedAt.After(oldUpdatedAt) {
		t.Error("Shipment.BeforeUpdate() did not update UpdatedAt")
	}
}

func TestShipment_UpdateStatus(t *testing.T) {
	shipment := &Shipment{
		ClientCompanyID: 1,
		Status:          ShipmentStatusPendingPickup,
	}

	if shipment.Status != ShipmentStatusPendingPickup {
		t.Error("Expected initial status to be pending_pickup_from_client")
	}

	shipment.UpdateStatus(ShipmentStatusPickedUpFromClient)

	if shipment.Status != ShipmentStatusPickedUpFromClient {
		t.Errorf("UpdateStatus() did not update status, got %v, want %v", shipment.Status, ShipmentStatusPickedUpFromClient)
	}
}

func TestShipment_IsDelivered(t *testing.T) {
	tests := []struct {
		name     string
		shipment Shipment
		expected bool
	}{
		{
			name: "delivered shipment",
			shipment: Shipment{
				ClientCompanyID: 1,
				Status:          ShipmentStatusDelivered,
			},
			expected: true,
		},
		{
			name: "pending shipment",
			shipment: Shipment{
				ClientCompanyID: 1,
				Status:          ShipmentStatusPendingPickup,
			},
			expected: false,
		},
		{
			name: "in transit shipment",
			shipment: Shipment{
				ClientCompanyID: 1,
				Status:          ShipmentStatusInTransitToEngineer,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.shipment.IsDelivered(); got != tt.expected {
				t.Errorf("Shipment.IsDelivered() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestShipment_IsAtWarehouse(t *testing.T) {
	tests := []struct {
		name     string
		shipment Shipment
		expected bool
	}{
		{
			name: "at warehouse",
			shipment: Shipment{
				ClientCompanyID: 1,
				Status:          ShipmentStatusAtWarehouse,
			},
			expected: true,
		},
		{
			name: "pending pickup",
			shipment: Shipment{
				ClientCompanyID: 1,
				Status:          ShipmentStatusPendingPickup,
			},
			expected: false,
		},
		{
			name: "delivered",
			shipment: Shipment{
				ClientCompanyID: 1,
				Status:          ShipmentStatusDelivered,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.shipment.IsAtWarehouse(); got != tt.expected {
				t.Errorf("Shipment.IsAtWarehouse() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestShipment_GetLaptopCount(t *testing.T) {
	// Test for shipment with no laptops
	shipment := Shipment{
		ClientCompanyID: 1,
		Status:          ShipmentStatusPendingPickup,
	}

	count := shipment.GetLaptopCount()
	if count != 0 {
		t.Errorf("GetLaptopCount() for empty shipment = %v, want 0", count)
	}

	// Test for shipment with laptops
	shipment.Laptops = []Laptop{
		{SerialNumber: "SN1", Status: LaptopStatusAvailable},
		{SerialNumber: "SN2", Status: LaptopStatusAvailable},
		{SerialNumber: "SN3", Status: LaptopStatusAvailable},
	}

	count = shipment.GetLaptopCount()
	if count != 3 {
		t.Errorf("GetLaptopCount() = %v, want 3", count)
	}
}

// TestShipment_GetTrackingURL tests the generation of tracking URLs for different couriers
func TestShipment_GetTrackingURL(t *testing.T) {
	tests := []struct {
		name           string
		courierName    string
		trackingNumber string
		expectedURL    string
	}{
		{
			name:           "UPS tracking URL",
			courierName:    "UPS",
			trackingNumber: "1Z9999999999999999",
			expectedURL:    "https://www.ups.com/track?tracknum=1Z9999999999999999",
		},
		{
			name:           "DHL tracking URL",
			courierName:    "DHL",
			trackingNumber: "1234567890",
			expectedURL:    "http://www.dhl.com/en/express/tracking.html?AWB=1234567890",
		},
		{
			name:           "FedEx tracking URL",
			courierName:    "FedEx",
			trackingNumber: "999999999999",
			expectedURL:    "https://www.fedex.com/fedextrack/?tracknumbers=999999999999",
		},
		{
			name:           "Case insensitive - ups (lowercase)",
			courierName:    "ups",
			trackingNumber: "1Z9999999999999999",
			expectedURL:    "https://www.ups.com/track?tracknum=1Z9999999999999999",
		},
		{
			name:           "Case insensitive - fedex (lowercase)",
			courierName:    "fedex",
			trackingNumber: "999999999999",
			expectedURL:    "https://www.fedex.com/fedextrack/?tracknumbers=999999999999",
		},
		{
			name:           "Case insensitive - dhl (lowercase)",
			courierName:    "dhl",
			trackingNumber: "1234567890",
			expectedURL:    "http://www.dhl.com/en/express/tracking.html?AWB=1234567890",
		},
		{
			name:           "FedEx with service type - FedEx Express",
			courierName:    "FedEx Express",
			trackingNumber: "999999999999",
			expectedURL:    "https://www.fedex.com/fedextrack/?tracknumbers=999999999999",
		},
		{
			name:           "UPS with service type - UPS Next Day Air",
			courierName:    "UPS Next Day Air",
			trackingNumber: "1Z9999999999999999",
			expectedURL:    "https://www.ups.com/track?tracknum=1Z9999999999999999",
		},
		{
			name:           "DHL with service type - DHL Express",
			courierName:    "DHL Express",
			trackingNumber: "1234567890",
			expectedURL:    "http://www.dhl.com/en/express/tracking.html?AWB=1234567890",
		},
		{
			name:           "Unknown courier returns empty string",
			courierName:    "Unknown Courier",
			trackingNumber: "TRACK123",
			expectedURL:    "",
		},
		{
			name:           "Empty courier name returns empty string",
			courierName:    "",
			trackingNumber: "TRACK123",
			expectedURL:    "",
		},
		{
			name:           "Empty tracking number still generates URL",
			courierName:    "UPS",
			trackingNumber: "",
			expectedURL:    "https://www.ups.com/track?tracknum=",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shipment := Shipment{
				CourierName:    tt.courierName,
				TrackingNumber: tt.trackingNumber,
			}
			got := shipment.GetTrackingURL()
			if got != tt.expectedURL {
				t.Errorf("Shipment.GetTrackingURL() = %v, want %v", got, tt.expectedURL)
			}
		})
	}
}

// TestShipment_UpdateStatus_WithETA tests that UpdateStatus properly handles ETA for in_transit_to_engineer status
func TestShipment_UpdateStatus_WithETA(t *testing.T) {
	shipment := &Shipment{
		ClientCompanyID: 1,
		Status:          ShipmentStatusAtWarehouse,
	}

	// Update status to in_transit_to_engineer with ETA
	eta := time.Now().Add(48 * time.Hour) // ETA in 48 hours
	shipment.UpdateStatusWithETA(ShipmentStatusInTransitToEngineer, &eta)

	// Verify status is updated
	if shipment.Status != ShipmentStatusInTransitToEngineer {
		t.Errorf("UpdateStatusWithETA() did not update status, got %v, want %v", shipment.Status, ShipmentStatusInTransitToEngineer)
	}

	// Verify ETA is set
	if shipment.ETAToEngineer == nil {
		t.Error("UpdateStatusWithETA() did not set ETAToEngineer")
	}
	if shipment.ETAToEngineer != nil && !shipment.ETAToEngineer.Equal(eta) {
		t.Errorf("UpdateStatusWithETA() ETAToEngineer = %v, want %v", shipment.ETAToEngineer, eta)
	}
}

// TestShipment_UpdateStatus_WithoutETA tests backward compatibility when no ETA is provided
func TestShipment_UpdateStatus_WithoutETA(t *testing.T) {
	shipment := &Shipment{
		ClientCompanyID: 1,
		Status:          ShipmentStatusAtWarehouse,
	}

	// Update status without providing ETA
	shipment.UpdateStatusWithETA(ShipmentStatusInTransitToEngineer, nil)

	// Verify status is updated
	if shipment.Status != ShipmentStatusInTransitToEngineer {
		t.Errorf("UpdateStatusWithETA() did not update status, got %v, want %v", shipment.Status, ShipmentStatusInTransitToEngineer)
	}

	// Verify ETA remains nil (optional field)
	if shipment.ETAToEngineer != nil {
		t.Error("UpdateStatusWithETA() should allow nil ETA")
	}
}

// TestShipment_GetNextAllowedStatuses tests getting the next valid statuses for sequential transitions
func TestShipment_GetNextAllowedStatuses(t *testing.T) {
	tests := []struct {
		name          string
		currentStatus ShipmentStatus
		expectedNext  []ShipmentStatus
	}{
		{
			name:          "from pending_pickup_from_client",
			currentStatus: ShipmentStatusPendingPickup,
			expectedNext:  []ShipmentStatus{ShipmentStatusPickupScheduled},
		},
		{
			name:          "from pickup_from_client_scheduled",
			currentStatus: ShipmentStatusPickupScheduled,
			expectedNext:  []ShipmentStatus{ShipmentStatusPickedUpFromClient},
		},
		{
			name:          "from picked_up_from_client",
			currentStatus: ShipmentStatusPickedUpFromClient,
			expectedNext:  []ShipmentStatus{ShipmentStatusInTransitToWarehouse},
		},
		{
			name:          "from in_transit_to_warehouse",
			currentStatus: ShipmentStatusInTransitToWarehouse,
			expectedNext:  []ShipmentStatus{ShipmentStatusAtWarehouse},
		},
		{
			name:          "from at_warehouse",
			currentStatus: ShipmentStatusAtWarehouse,
			expectedNext:  []ShipmentStatus{ShipmentStatusReleasedFromWarehouse},
		},
		{
			name:          "from released_from_warehouse",
			currentStatus: ShipmentStatusReleasedFromWarehouse,
			expectedNext:  []ShipmentStatus{ShipmentStatusInTransitToEngineer},
		},
		{
			name:          "from in_transit_to_engineer",
			currentStatus: ShipmentStatusInTransitToEngineer,
			expectedNext:  []ShipmentStatus{ShipmentStatusDelivered},
		},
		{
			name:          "from delivered (final status)",
			currentStatus: ShipmentStatusDelivered,
			expectedNext:  []ShipmentStatus{}, // No next statuses
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shipment := &Shipment{
				Status: tt.currentStatus,
			}
			got := shipment.GetNextAllowedStatuses()

			// Check length matches
			if len(got) != len(tt.expectedNext) {
				t.Errorf("GetNextAllowedStatuses() returned %d statuses, expected %d", len(got), len(tt.expectedNext))
				return
			}

			// Check each expected status is present
			for _, expected := range tt.expectedNext {
				found := false
				for _, actual := range got {
					if actual == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("GetNextAllowedStatuses() missing expected status %v", expected)
				}
			}
		})
	}
}

// TestShipment_IsValidStatusTransition tests sequential status transition validation
func TestShipment_IsValidStatusTransition(t *testing.T) {
	tests := []struct {
		name          string
		currentStatus ShipmentStatus
		newStatus     ShipmentStatus
		expected      bool
	}{
		// Valid transitions (sequential, forward only)
		{
			name:          "valid: pending_pickup_from_client -> pickup_from_client_scheduled",
			currentStatus: ShipmentStatusPendingPickup,
			newStatus:     ShipmentStatusPickupScheduled,
			expected:      true,
		},
		{
			name:          "valid: pickup_from_client_scheduled -> picked_up_from_client",
			currentStatus: ShipmentStatusPickupScheduled,
			newStatus:     ShipmentStatusPickedUpFromClient,
			expected:      true,
		},
		{
			name:          "valid: picked_up_from_client -> in_transit_to_warehouse",
			currentStatus: ShipmentStatusPickedUpFromClient,
			newStatus:     ShipmentStatusInTransitToWarehouse,
			expected:      true,
		},
		{
			name:          "valid: in_transit_to_warehouse -> at_warehouse",
			currentStatus: ShipmentStatusInTransitToWarehouse,
			newStatus:     ShipmentStatusAtWarehouse,
			expected:      true,
		},
		{
			name:          "valid: at_warehouse -> released_from_warehouse",
			currentStatus: ShipmentStatusAtWarehouse,
			newStatus:     ShipmentStatusReleasedFromWarehouse,
			expected:      true,
		},
		{
			name:          "valid: released_from_warehouse -> in_transit_to_engineer",
			currentStatus: ShipmentStatusReleasedFromWarehouse,
			newStatus:     ShipmentStatusInTransitToEngineer,
			expected:      true,
		},
		{
			name:          "valid: in_transit_to_engineer -> delivered",
			currentStatus: ShipmentStatusInTransitToEngineer,
			newStatus:     ShipmentStatusDelivered,
			expected:      true,
		},

		// Invalid transitions - skipping statuses
		{
			name:          "invalid: pending_pickup_from_client -> picked_up_from_client (skipping pickup_from_client_scheduled)",
			currentStatus: ShipmentStatusPendingPickup,
			newStatus:     ShipmentStatusPickedUpFromClient,
			expected:      false,
		},
		{
			name:          "invalid: pending_pickup_from_client -> at_warehouse (skipping multiple)",
			currentStatus: ShipmentStatusPendingPickup,
			newStatus:     ShipmentStatusAtWarehouse,
			expected:      false,
		},
		{
			name:          "invalid: pickup_from_client_scheduled -> at_warehouse (skipping multiple)",
			currentStatus: ShipmentStatusPickupScheduled,
			newStatus:     ShipmentStatusAtWarehouse,
			expected:      false,
		},

		// Invalid transitions - going backwards
		{
			name:          "invalid: at_warehouse -> pending_pickup_from_client (backwards)",
			currentStatus: ShipmentStatusAtWarehouse,
			newStatus:     ShipmentStatusPendingPickup,
			expected:      false,
		},
		{
			name:          "invalid: delivered -> in_transit_to_engineer (backwards)",
			currentStatus: ShipmentStatusDelivered,
			newStatus:     ShipmentStatusInTransitToEngineer,
			expected:      false,
		},
		{
			name:          "invalid: in_transit_to_warehouse -> picked_up_from_client (backwards)",
			currentStatus: ShipmentStatusInTransitToWarehouse,
			newStatus:     ShipmentStatusPickedUpFromClient,
			expected:      false,
		},

		// Invalid transitions - same status
		{
			name:          "invalid: pending_pickup_from_client -> pending_pickup_from_client (same status)",
			currentStatus: ShipmentStatusPendingPickup,
			newStatus:     ShipmentStatusPendingPickup,
			expected:      false,
		},
		{
			name:          "invalid: at_warehouse -> at_warehouse (same status)",
			currentStatus: ShipmentStatusAtWarehouse,
			newStatus:     ShipmentStatusAtWarehouse,
			expected:      false,
		},

		// Invalid transitions - from delivered (final status)
		{
			name:          "invalid: delivered -> any status (final status)",
			currentStatus: ShipmentStatusDelivered,
			newStatus:     ShipmentStatusDelivered,
			expected:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shipment := &Shipment{
				Status: tt.currentStatus,
			}
			got := shipment.IsValidStatusTransition(tt.newStatus)
			if got != tt.expected {
				t.Errorf("IsValidStatusTransition() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestShipmentType_Validation(t *testing.T) {
	validTypes := []ShipmentType{
		ShipmentTypeSingleFullJourney,
		ShipmentTypeBulkToWarehouse,
		ShipmentTypeWarehouseToEngineer,
	}

	for _, shipmentType := range validTypes {
		t.Run(string(shipmentType), func(t *testing.T) {
			if !IsValidShipmentType(shipmentType) {
				t.Errorf("Expected %s to be valid", shipmentType)
			}
		})
	}

	// Test invalid type
	t.Run("invalid type", func(t *testing.T) {
		if IsValidShipmentType("invalid_type") {
			t.Error("Expected invalid_type to be invalid")
		}
	})
}

// Helper function for creating int64 pointers
func int64Ptr(i int64) *int64 {
	return &i
}
