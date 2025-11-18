package models

import (
	"strings"
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
				ShipmentType:        ShipmentTypeSingleFullJourney,
				ClientCompanyID:     1,
				LaptopCount:         1,
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
				ShipmentType:     ShipmentTypeSingleFullJourney,
				ClientCompanyID:  2,
				LaptopCount:      1,
				Status:           ShipmentStatusPendingPickup,
				JiraTicketNumber: "PROJECT-12345",
			},
			wantErr: false,
		},
		{
			name: "invalid - missing JIRA ticket number",
			shipment: Shipment{
				ShipmentType:    ShipmentTypeSingleFullJourney,
				ClientCompanyID: 1,
				LaptopCount:     1,
				Status:          ShipmentStatusPendingPickup,
			},
			wantErr: true,
			errMsg:  "JIRA ticket number is required",
		},
		{
			name: "invalid - empty JIRA ticket number",
			shipment: Shipment{
				ShipmentType:     ShipmentTypeSingleFullJourney,
				ClientCompanyID:  1,
				LaptopCount:      1,
				Status:           ShipmentStatusPendingPickup,
				JiraTicketNumber: "",
			},
			wantErr: true,
			errMsg:  "JIRA ticket number is required",
		},
		{
			name: "invalid - malformed JIRA ticket (missing dash)",
			shipment: Shipment{
				ShipmentType:     ShipmentTypeSingleFullJourney,
				ClientCompanyID:  1,
				LaptopCount:      1,
				Status:           ShipmentStatusPendingPickup,
				JiraTicketNumber: "SCOP67702",
			},
			wantErr: true,
			errMsg:  "JIRA ticket number must be in format PROJECT-NUMBER (e.g., SCOP-67702)",
		},
		{
			name: "invalid - malformed JIRA ticket (no project key)",
			shipment: Shipment{
				ShipmentType:     ShipmentTypeSingleFullJourney,
				ClientCompanyID:  1,
				LaptopCount:      1,
				Status:           ShipmentStatusPendingPickup,
				JiraTicketNumber: "-67702",
			},
			wantErr: true,
			errMsg:  "JIRA ticket number must be in format PROJECT-NUMBER (e.g., SCOP-67702)",
		},
		{
			name: "invalid - malformed JIRA ticket (no number)",
			shipment: Shipment{
				ShipmentType:     ShipmentTypeSingleFullJourney,
				ClientCompanyID:  1,
				LaptopCount:      1,
				Status:           ShipmentStatusPendingPickup,
				JiraTicketNumber: "SCOP-",
			},
			wantErr: true,
			errMsg:  "JIRA ticket number must be in format PROJECT-NUMBER (e.g., SCOP-67702)",
		},
		{
			name: "invalid - malformed JIRA ticket (lowercase)",
			shipment: Shipment{
				ShipmentType:     ShipmentTypeSingleFullJourney,
				ClientCompanyID:  1,
				LaptopCount:      1,
				Status:           ShipmentStatusPendingPickup,
				JiraTicketNumber: "scop-67702",
			},
			wantErr: true,
			errMsg:  "JIRA ticket number must be in format PROJECT-NUMBER (e.g., SCOP-67702)",
		},
		{
			name: "invalid - malformed JIRA ticket (special characters)",
			shipment: Shipment{
				ShipmentType:     ShipmentTypeSingleFullJourney,
				ClientCompanyID:  1,
				LaptopCount:      1,
				Status:           ShipmentStatusPendingPickup,
				JiraTicketNumber: "SC@P-67702",
			},
			wantErr: true,
			errMsg:  "JIRA ticket number must be in format PROJECT-NUMBER (e.g., SCOP-67702)",
		},
		{
			name: "invalid - missing client company ID",
			shipment: Shipment{
				ShipmentType:     ShipmentTypeSingleFullJourney,
				LaptopCount:      1,
				Status:           ShipmentStatusPendingPickup,
				JiraTicketNumber: "SCOP-67702",
			},
			wantErr: true,
			errMsg:  "client company ID is required",
		},
		{
			name: "invalid - missing status",
			shipment: Shipment{
				ShipmentType:     ShipmentTypeSingleFullJourney,
				ClientCompanyID:  1,
				LaptopCount:      1,
				JiraTicketNumber: "SCOP-67702",
			},
			wantErr: true,
			errMsg:  "status is required",
		},
		{
			name: "invalid - invalid status",
			shipment: Shipment{
				ShipmentType:     ShipmentTypeSingleFullJourney,
				ClientCompanyID:  1,
				LaptopCount:      1,
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
				ShipmentType: ShipmentTypeSingleFullJourney, // Full journey through all statuses
				Status:       tt.currentStatus,
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
				ShipmentType: ShipmentTypeSingleFullJourney, // Full journey through all statuses
				Status:       tt.currentStatus,
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

func TestShipment_EngineerAssignmentRules(t *testing.T) {
	engineerID := int64(1)
	
	tests := []struct {
		name          string
		shipmentType  ShipmentType
		status        ShipmentStatus
		engineerID    *int64
		shouldBeValid bool
		errorContains string
	}{
		{
			name:          "single_full_journey can have engineer assigned before release",
			shipmentType:  ShipmentTypeSingleFullJourney,
			status:        ShipmentStatusAtWarehouse,
			engineerID:    nil,
			shouldBeValid: true,
		},
		{
			name:          "single_full_journey can have engineer assigned",
			shipmentType:  ShipmentTypeSingleFullJourney,
			status:        ShipmentStatusAtWarehouse,
			engineerID:    &engineerID,
			shouldBeValid: true,
		},
		{
			name:          "warehouse_to_engineer must have engineer assigned",
			shipmentType:  ShipmentTypeWarehouseToEngineer,
			status:        ShipmentStatusReleasedFromWarehouse,
			engineerID:    nil,
			shouldBeValid: false,
			errorContains: "must have software engineer assigned",
		},
		{
			name:          "warehouse_to_engineer with engineer is valid",
			shipmentType:  ShipmentTypeWarehouseToEngineer,
			status:        ShipmentStatusReleasedFromWarehouse,
			engineerID:    &engineerID,
			shouldBeValid: true,
		},
		{
			name:          "bulk_to_warehouse cannot have engineer assigned",
			shipmentType:  ShipmentTypeBulkToWarehouse,
			status:        ShipmentStatusAtWarehouse,
			engineerID:    &engineerID,
			shouldBeValid: false,
			errorContains: "cannot have software engineer assigned",
		},
		{
			name:          "bulk_to_warehouse without engineer is valid",
			shipmentType:  ShipmentTypeBulkToWarehouse,
			status:        ShipmentStatusAtWarehouse,
			engineerID:    nil,
			shouldBeValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Shipment{
				ShipmentType:       tt.shipmentType,
				Status:             tt.status,
				SoftwareEngineerID: tt.engineerID,
				ClientCompanyID:    1,
				JiraTicketNumber:   "SCOP-12345",
			}

			err := s.ValidateEngineerAssignment()
			if tt.shouldBeValid && err != nil {
				t.Errorf("Expected valid, got error: %v", err)
			}
			if !tt.shouldBeValid {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if tt.errorContains != "" {
					if !strings.Contains(err.Error(), tt.errorContains) {
						t.Errorf("Expected error to contain '%s', got: %v", tt.errorContains, err)
					}
				}
			}
		})
	}
}

func TestShipment_TypeSpecificStatusFlows(t *testing.T) {
	tests := []struct {
		name          string
		shipmentType  ShipmentType
		currentStatus ShipmentStatus
		nextStatus    ShipmentStatus
		shouldBeValid bool
	}{
		// Single full journey - full flow
		{
			name:          "single_full_journey allows full flow to delivered",
			shipmentType:  ShipmentTypeSingleFullJourney,
			currentStatus: ShipmentStatusInTransitToEngineer,
			nextStatus:    ShipmentStatusDelivered,
			shouldBeValid: true,
		},
		{
			name:          "single_full_journey allows warehouse to release transition",
			shipmentType:  ShipmentTypeSingleFullJourney,
			currentStatus: ShipmentStatusAtWarehouse,
			nextStatus:    ShipmentStatusReleasedFromWarehouse,
			shouldBeValid: true,
		},
		// Bulk to warehouse - stops at warehouse
		{
			name:          "bulk_to_warehouse can reach at_warehouse",
			shipmentType:  ShipmentTypeBulkToWarehouse,
			currentStatus: ShipmentStatusInTransitToWarehouse,
			nextStatus:    ShipmentStatusAtWarehouse,
			shouldBeValid: true,
		},
		{
			name:          "bulk_to_warehouse cannot go past at_warehouse to released",
			shipmentType:  ShipmentTypeBulkToWarehouse,
			currentStatus: ShipmentStatusAtWarehouse,
			nextStatus:    ShipmentStatusReleasedFromWarehouse,
			shouldBeValid: false,
		},
		{
			name:          "bulk_to_warehouse cannot reach in_transit_to_engineer",
			shipmentType:  ShipmentTypeBulkToWarehouse,
			currentStatus: ShipmentStatusReleasedFromWarehouse,
			nextStatus:    ShipmentStatusInTransitToEngineer,
			shouldBeValid: false,
		},
		// Warehouse to engineer - starts from released
		{
			name:          "warehouse_to_engineer starts from released_from_warehouse",
			shipmentType:  ShipmentTypeWarehouseToEngineer,
			currentStatus: ShipmentStatusReleasedFromWarehouse,
			nextStatus:    ShipmentStatusInTransitToEngineer,
			shouldBeValid: true,
		},
		{
			name:          "warehouse_to_engineer can reach delivered",
			shipmentType:  ShipmentTypeWarehouseToEngineer,
			currentStatus: ShipmentStatusInTransitToEngineer,
			nextStatus:    ShipmentStatusDelivered,
			shouldBeValid: true,
		},
		{
			name:          "warehouse_to_engineer cannot have pending_pickup status",
			shipmentType:  ShipmentTypeWarehouseToEngineer,
			currentStatus: ShipmentStatusPendingPickup,
			nextStatus:    ShipmentStatusPickupScheduled,
			shouldBeValid: false,
		},
		{
			name:          "warehouse_to_engineer cannot have at_warehouse status",
			shipmentType:  ShipmentTypeWarehouseToEngineer,
			currentStatus: ShipmentStatusAtWarehouse,
			nextStatus:    ShipmentStatusReleasedFromWarehouse,
			shouldBeValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Shipment{
				ShipmentType:     tt.shipmentType,
				Status:           tt.currentStatus,
				ClientCompanyID:  1,
				JiraTicketNumber: "SCOP-12345",
			}

			isValid := s.IsValidStatusTransition(tt.nextStatus)
			if tt.shouldBeValid && !isValid {
				t.Errorf("Expected transition from %s to %s to be valid for %s", tt.currentStatus, tt.nextStatus, tt.shipmentType)
			}
			if !tt.shouldBeValid && isValid {
				t.Errorf("Expected transition from %s to %s to be invalid for %s", tt.currentStatus, tt.nextStatus, tt.shipmentType)
			}
		})
	}
}

func TestShipment_LaptopCountTracking(t *testing.T) {
	tests := []struct {
		name          string
		shipmentType  ShipmentType
		laptopCount   int
		shouldBeValid bool
		errorContains string
	}{
		{
			name:          "single_full_journey must have count = 1",
			shipmentType:  ShipmentTypeSingleFullJourney,
			laptopCount:   1,
			shouldBeValid: true,
		},
		{
			name:          "single_full_journey cannot have count > 1",
			shipmentType:  ShipmentTypeSingleFullJourney,
			laptopCount:   2,
			shouldBeValid: false,
			errorContains: "exactly 1 laptop",
		},
		{
			name:          "single_full_journey cannot have count = 0",
			shipmentType:  ShipmentTypeSingleFullJourney,
			laptopCount:   0,
			shouldBeValid: false,
			errorContains: "exactly 1 laptop",
		},
		{
			name:          "bulk_to_warehouse must have count >= 2",
			shipmentType:  ShipmentTypeBulkToWarehouse,
			laptopCount:   2,
			shouldBeValid: true,
		},
		{
			name:          "bulk_to_warehouse can have count > 2",
			shipmentType:  ShipmentTypeBulkToWarehouse,
			laptopCount:   10,
			shouldBeValid: true,
		},
		{
			name:          "bulk_to_warehouse cannot have count = 1",
			shipmentType:  ShipmentTypeBulkToWarehouse,
			laptopCount:   1,
			shouldBeValid: false,
			errorContains: "at least 2 laptops",
		},
		{
			name:          "bulk_to_warehouse cannot have count = 0",
			shipmentType:  ShipmentTypeBulkToWarehouse,
			laptopCount:   0,
			shouldBeValid: false,
			errorContains: "at least 2 laptops",
		},
		{
			name:          "warehouse_to_engineer must have count = 1",
			shipmentType:  ShipmentTypeWarehouseToEngineer,
			laptopCount:   1,
			shouldBeValid: true,
		},
		{
			name:          "warehouse_to_engineer cannot have count > 1",
			shipmentType:  ShipmentTypeWarehouseToEngineer,
			laptopCount:   2,
			shouldBeValid: false,
			errorContains: "exactly 1 laptop",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Shipment{
				ShipmentType:     tt.shipmentType,
				LaptopCount:      tt.laptopCount,
				ClientCompanyID:  1,
				Status:           ShipmentStatusPendingPickup,
				JiraTicketNumber: "SCOP-12345",
			}

			err := s.ValidateLaptopCount()
			if tt.shouldBeValid && err != nil {
				t.Errorf("Expected valid, got error: %v", err)
			}
			if !tt.shouldBeValid {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if tt.errorContains != "" {
					if !strings.Contains(err.Error(), tt.errorContains) {
						t.Errorf("Expected error to contain '%s', got: %v", tt.errorContains, err)
					}
				}
			}
		})
	}
}

func TestShipment_ValidateWithType(t *testing.T) {
	engineerID := int64(1)
	
	tests := []struct {
		name          string
		shipment      Shipment
		shouldBeValid bool
		errorContains string
	}{
		{
			name: "valid single_full_journey",
			shipment: Shipment{
				ShipmentType:     ShipmentTypeSingleFullJourney,
				ClientCompanyID:  1,
				LaptopCount:      1,
				Status:           ShipmentStatusPendingPickup,
				JiraTicketNumber: "SCOP-12345",
			},
			shouldBeValid: true,
		},
		{
			name: "invalid single_full_journey - wrong laptop count",
			shipment: Shipment{
				ShipmentType:     ShipmentTypeSingleFullJourney,
				ClientCompanyID:  1,
				LaptopCount:      2,
				Status:           ShipmentStatusPendingPickup,
				JiraTicketNumber: "SCOP-12345",
			},
			shouldBeValid: false,
			errorContains: "exactly 1 laptop",
		},
		{
			name: "invalid single_full_journey - invalid status for type",
			shipment: Shipment{
				ShipmentType:     ShipmentTypeSingleFullJourney,
				ClientCompanyID:  1,
				LaptopCount:      1,
				Status:           "invalid_status",
				JiraTicketNumber: "SCOP-12345",
			},
			shouldBeValid: false,
			errorContains: "invalid status",
		},
		{
			name: "valid bulk_to_warehouse",
			shipment: Shipment{
				ShipmentType:     ShipmentTypeBulkToWarehouse,
				ClientCompanyID:  1,
				LaptopCount:      5,
				Status:           ShipmentStatusAtWarehouse,
				JiraTicketNumber: "SCOP-12345",
			},
			shouldBeValid: true,
		},
		{
			name: "invalid bulk_to_warehouse - has engineer assigned",
			shipment: Shipment{
				ShipmentType:       ShipmentTypeBulkToWarehouse,
				ClientCompanyID:    1,
				LaptopCount:        5,
				Status:             ShipmentStatusAtWarehouse,
				JiraTicketNumber:   "SCOP-12345",
				SoftwareEngineerID: &engineerID,
			},
			shouldBeValid: false,
			errorContains: "cannot have software engineer assigned",
		},
		{
			name: "invalid bulk_to_warehouse - laptop count too low",
			shipment: Shipment{
				ShipmentType:     ShipmentTypeBulkToWarehouse,
				ClientCompanyID:  1,
				LaptopCount:      1,
				Status:           ShipmentStatusAtWarehouse,
				JiraTicketNumber: "SCOP-12345",
			},
			shouldBeValid: false,
			errorContains: "at least 2 laptops",
		},
		{
			name: "invalid bulk_to_warehouse - status past allowed",
			shipment: Shipment{
				ShipmentType:     ShipmentTypeBulkToWarehouse,
				ClientCompanyID:  1,
				LaptopCount:      5,
				Status:           ShipmentStatusInTransitToEngineer,
				JiraTicketNumber: "SCOP-12345",
			},
			shouldBeValid: false,
			errorContains: "not valid for shipment type",
		},
		{
			name: "valid warehouse_to_engineer",
			shipment: Shipment{
				ShipmentType:       ShipmentTypeWarehouseToEngineer,
				ClientCompanyID:    1,
				LaptopCount:        1,
				Status:             ShipmentStatusReleasedFromWarehouse,
				JiraTicketNumber:   "SCOP-12345",
				SoftwareEngineerID: &engineerID,
			},
			shouldBeValid: true,
		},
		{
			name: "invalid warehouse_to_engineer - missing engineer",
			shipment: Shipment{
				ShipmentType:     ShipmentTypeWarehouseToEngineer,
				ClientCompanyID:  1,
				LaptopCount:      1,
				Status:           ShipmentStatusReleasedFromWarehouse,
				JiraTicketNumber: "SCOP-12345",
			},
			shouldBeValid: false,
			errorContains: "must have software engineer assigned",
		},
		{
			name: "invalid warehouse_to_engineer - wrong laptop count",
			shipment: Shipment{
				ShipmentType:       ShipmentTypeWarehouseToEngineer,
				ClientCompanyID:    1,
				LaptopCount:        2,
				Status:             ShipmentStatusReleasedFromWarehouse,
				JiraTicketNumber:   "SCOP-12345",
				SoftwareEngineerID: &engineerID,
			},
			shouldBeValid: false,
			errorContains: "exactly 1 laptop",
		},
		{
			name: "invalid warehouse_to_engineer - invalid status for type",
			shipment: Shipment{
				ShipmentType:       ShipmentTypeWarehouseToEngineer,
				ClientCompanyID:    1,
				LaptopCount:        1,
				Status:             ShipmentStatusPendingPickup,
				JiraTicketNumber:   "SCOP-12345",
				SoftwareEngineerID: &engineerID,
			},
			shouldBeValid: false,
			errorContains: "not valid for shipment type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.shipment.Validate()
			if tt.shouldBeValid && err != nil {
				t.Errorf("Expected valid, got error: %v", err)
			}
			if !tt.shouldBeValid {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if tt.errorContains != "" {
					if !strings.Contains(err.Error(), tt.errorContains) {
						t.Errorf("Expected error to contain '%s', got: %v", tt.errorContains, err)
					}
				}
			}
		})
	}
}

func TestShipment_SyncLaptopStatusOnUpdate(t *testing.T) {
	tests := []struct {
		name                 string
		shipmentType         ShipmentType
		shipmentStatus       ShipmentStatus
		expectedLaptopStatus LaptopStatus
		shouldSync           bool
	}{
		// Single full journey - should sync
		{
			name:                 "single_full_journey syncs laptop status - in transit to warehouse",
			shipmentType:         ShipmentTypeSingleFullJourney,
			shipmentStatus:       ShipmentStatusInTransitToWarehouse,
			expectedLaptopStatus: LaptopStatusInTransitToWarehouse,
			shouldSync:           true,
		},
		{
			name:                 "single_full_journey syncs laptop status - at warehouse",
			shipmentType:         ShipmentTypeSingleFullJourney,
			shipmentStatus:       ShipmentStatusAtWarehouse,
			expectedLaptopStatus: LaptopStatusAtWarehouse,
			shouldSync:           true,
		},
		{
			name:                 "single_full_journey syncs laptop status - in transit to engineer",
			shipmentType:         ShipmentTypeSingleFullJourney,
			shipmentStatus:       ShipmentStatusInTransitToEngineer,
			expectedLaptopStatus: LaptopStatusInTransitToEngineer,
			shouldSync:           true,
		},
		{
			name:                 "single_full_journey syncs laptop status - delivered",
			shipmentType:         ShipmentTypeSingleFullJourney,
			shipmentStatus:       ShipmentStatusDelivered,
			expectedLaptopStatus: LaptopStatusDelivered,
			shouldSync:           true,
		},
		// Warehouse to engineer - should sync
		{
			name:                 "warehouse_to_engineer syncs laptop status - in transit",
			shipmentType:         ShipmentTypeWarehouseToEngineer,
			shipmentStatus:       ShipmentStatusInTransitToEngineer,
			expectedLaptopStatus: LaptopStatusInTransitToEngineer,
			shouldSync:           true,
		},
		{
			name:                 "warehouse_to_engineer syncs laptop status - delivered",
			shipmentType:         ShipmentTypeWarehouseToEngineer,
			shipmentStatus:       ShipmentStatusDelivered,
			expectedLaptopStatus: LaptopStatusDelivered,
			shouldSync:           true,
		},
		// Bulk to warehouse - should NOT sync
		{
			name:           "bulk_to_warehouse does not sync laptop status",
			shipmentType:   ShipmentTypeBulkToWarehouse,
			shipmentStatus: ShipmentStatusInTransitToWarehouse,
			shouldSync:     false,
		},
		{
			name:           "bulk_to_warehouse at warehouse does not sync",
			shipmentType:   ShipmentTypeBulkToWarehouse,
			shipmentStatus: ShipmentStatusAtWarehouse,
			shouldSync:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Shipment{
				ShipmentType:     tt.shipmentType,
				Status:           tt.shipmentStatus,
				ClientCompanyID:  1,
				JiraTicketNumber: "SCOP-12345",
			}

			shouldSync := s.ShouldSyncLaptopStatus()
			if shouldSync != tt.shouldSync {
				t.Errorf("Expected shouldSync=%v, got %v", tt.shouldSync, shouldSync)
			}

			if tt.shouldSync {
				laptopStatus := s.GetLaptopStatusForShipmentStatus()
				if laptopStatus != tt.expectedLaptopStatus {
					t.Errorf("Expected laptop status %s, got %s", tt.expectedLaptopStatus, laptopStatus)
				}
			}
		})
	}
}

// TestShipment_SecondTrackingNumber tests that second tracking number can be set and retrieved
func TestShipment_SecondTrackingNumber(t *testing.T) {
	tests := []struct {
		name                  string
		shipment              Shipment
		expectedSecondTracking string
	}{
		{
			name: "shipment with second tracking number",
			shipment: Shipment{
				ShipmentType:         ShipmentTypeSingleFullJourney,
				ClientCompanyID:      1,
				LaptopCount:          1,
				Status:               ShipmentStatusInTransitToEngineer,
				JiraTicketNumber:     "SCOP-12345",
				TrackingNumber:       "TRACK123456",
				SecondTrackingNumber:  "TRACK789012",
			},
			expectedSecondTracking: "TRACK789012",
		},
		{
			name: "shipment without second tracking number",
			shipment: Shipment{
				ShipmentType:     ShipmentTypeSingleFullJourney,
				ClientCompanyID:  1,
				LaptopCount:      1,
				Status:           ShipmentStatusInTransitToEngineer,
				JiraTicketNumber: "SCOP-12345",
				TrackingNumber:   "TRACK123456",
			},
			expectedSecondTracking: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shipment.SecondTrackingNumber != tt.expectedSecondTracking {
				t.Errorf("Expected second tracking number %s, got %s", tt.expectedSecondTracking, tt.shipment.SecondTrackingNumber)
			}
		})
	}
}

// Helper function for creating int64 pointers
func int64Ptr(i int64) *int64 {
	return &i
}
