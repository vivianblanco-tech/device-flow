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

// Helper function for creating int64 pointers
func int64Ptr(i int64) *int64 {
	return &i
}
