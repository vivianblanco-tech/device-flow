package models

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

// ShipmentStatus represents the status of a shipment in the system
type ShipmentStatus string

// Shipment status constants matching the process steps
const (
	ShipmentStatusPendingPickup          ShipmentStatus = "pending_pickup_from_client"
	ShipmentStatusPickupScheduled        ShipmentStatus = "pickup_from_client_scheduled"
	ShipmentStatusPickedUpFromClient     ShipmentStatus = "picked_up_from_client"
	ShipmentStatusInTransitToWarehouse   ShipmentStatus = "in_transit_to_warehouse"
	ShipmentStatusAtWarehouse            ShipmentStatus = "at_warehouse"
	ShipmentStatusReleasedFromWarehouse  ShipmentStatus = "released_from_warehouse"
	ShipmentStatusInTransitToEngineer    ShipmentStatus = "in_transit_to_engineer"
	ShipmentStatusDelivered              ShipmentStatus = "delivered"
)

// Courier name constants
const (
	CourierUPS   = "UPS"
	CourierFedEx = "FedEx"
	CourierDHL   = "DHL"
)

// Shipment represents a shipment of laptops through the delivery pipeline
type Shipment struct {
	ID                  int64           `json:"id" db:"id"`
	ClientCompanyID     int64           `json:"client_company_id" db:"client_company_id"`
	SoftwareEngineerID  *int64          `json:"software_engineer_id,omitempty" db:"software_engineer_id"`
	Status              ShipmentStatus  `json:"status" db:"status"`
	JiraTicketNumber    string          `json:"jira_ticket_number" db:"jira_ticket_number"`
	CourierName         string          `json:"courier_name,omitempty" db:"courier_name"`
	TrackingNumber      string          `json:"tracking_number,omitempty" db:"tracking_number"`
	
	// Tracking dates for each step
	PickupScheduledDate *time.Time      `json:"pickup_scheduled_date,omitempty" db:"pickup_scheduled_date"`
	PickedUpAt          *time.Time      `json:"picked_up_at,omitempty" db:"picked_up_at"`
	ArrivedWarehouseAt  *time.Time      `json:"arrived_warehouse_at,omitempty" db:"arrived_warehouse_at"`
	ReleasedWarehouseAt *time.Time      `json:"released_warehouse_at,omitempty" db:"released_warehouse_at"`
	ETAToEngineer       *time.Time      `json:"eta_to_engineer,omitempty" db:"eta_to_engineer"`
	DeliveredAt         *time.Time      `json:"delivered_at,omitempty" db:"delivered_at"`
	
	Notes               string          `json:"notes,omitempty" db:"notes"`
	CreatedAt           time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time       `json:"updated_at" db:"updated_at"`

	// Relations (not stored in DB directly)
	ClientCompany     *ClientCompany     `json:"client_company,omitempty" db:"-"`
	SoftwareEngineer  *SoftwareEngineer  `json:"software_engineer,omitempty" db:"-"`
	Laptops           []Laptop           `json:"laptops,omitempty" db:"-"`
}

// Validate validates the Shipment model
func (s *Shipment) Validate() error {
	// Client company validation
	if s.ClientCompanyID == 0 {
		return errors.New("client company ID is required")
	}

	// Status validation
	if s.Status == "" {
		return errors.New("status is required")
	}
	if !IsValidShipmentStatus(s.Status) {
		return errors.New("invalid status")
	}

	// JIRA ticket validation
	if s.JiraTicketNumber == "" {
		return errors.New("JIRA ticket number is required")
	}
	if !IsValidJiraTicketFormat(s.JiraTicketNumber) {
		return errors.New("JIRA ticket number must be in format PROJECT-NUMBER (e.g., SCOP-67702)")
	}

	return nil
}

// IsValidJiraTicketFormat validates the JIRA ticket number format (PROJECT-NUMBER)
func IsValidJiraTicketFormat(ticket string) bool {
	// Pattern: uppercase letters, dash, digits
	// Example: SCOP-67702, PROJECT-12345
	pattern := `^[A-Z]+\-[0-9]+$`
	matched, _ := regexp.MatchString(pattern, ticket)
	return matched
}

// JiraTicketValidator is a function type that validates if a JIRA ticket exists
// Returns nil if ticket exists, error otherwise
type JiraTicketValidator func(ticketKey string) error

// ValidateJiraTicketExists validates that a JIRA ticket exists using the provided validator
// If validator is nil, validation is skipped (for sample/test data)
func ValidateJiraTicketExists(ticketKey string, validator JiraTicketValidator) error {
	// Skip validation if no validator provided (sample data mode)
	if validator == nil {
		return nil
	}

	// Use the validator to check if ticket exists
	return validator(ticketKey)
}

// IsValidShipmentStatus checks if a given status is valid
func IsValidShipmentStatus(status ShipmentStatus) bool {
	switch status {
	case ShipmentStatusPendingPickup,
		ShipmentStatusPickupScheduled,
		ShipmentStatusPickedUpFromClient,
		ShipmentStatusInTransitToWarehouse,
		ShipmentStatusAtWarehouse,
		ShipmentStatusReleasedFromWarehouse,
		ShipmentStatusInTransitToEngineer,
		ShipmentStatusDelivered:
		return true
	}
	return false
}

// IsValidCourier checks if a given courier name is valid
func IsValidCourier(courier string) bool {
	switch courier {
	case CourierUPS, CourierFedEx, CourierDHL:
		return true
	}
	return false
}

// TableName returns the table name for the Shipment model
func (s *Shipment) TableName() string {
	return "shipments"
}

// BeforeCreate sets the timestamps before creating a shipment
func (s *Shipment) BeforeCreate() {
	now := time.Now()
	s.CreatedAt = now
	s.UpdatedAt = now
}

// BeforeUpdate sets the updated_at timestamp before updating a shipment
func (s *Shipment) BeforeUpdate() {
	s.UpdatedAt = time.Now()
}

// UpdateStatus updates the shipment status and the corresponding timestamp
func (s *Shipment) UpdateStatus(status ShipmentStatus) {
	s.Status = status
	now := time.Now()

	// Update the appropriate timestamp based on status
	switch status {
	case ShipmentStatusPickupScheduled:
		// If pickup scheduled date is not already set, set it to now
		if s.PickupScheduledDate == nil {
			s.PickupScheduledDate = &now
		}
	case ShipmentStatusPickedUpFromClient:
		s.PickedUpAt = &now
	case ShipmentStatusAtWarehouse:
		s.ArrivedWarehouseAt = &now
	case ShipmentStatusReleasedFromWarehouse:
		s.ReleasedWarehouseAt = &now
	case ShipmentStatusDelivered:
		s.DeliveredAt = &now
	}

	s.BeforeUpdate()
}

// UpdateStatusWithETA updates the shipment status and optionally sets an ETA for in_transit_to_engineer status
func (s *Shipment) UpdateStatusWithETA(status ShipmentStatus, eta *time.Time) {
	s.Status = status
	now := time.Now()

	// Update the appropriate timestamp based on status
	switch status {
	case ShipmentStatusPickupScheduled:
		// If pickup scheduled date is not already set, set it to now
		if s.PickupScheduledDate == nil {
			s.PickupScheduledDate = &now
		}
	case ShipmentStatusPickedUpFromClient:
		s.PickedUpAt = &now
	case ShipmentStatusAtWarehouse:
		s.ArrivedWarehouseAt = &now
	case ShipmentStatusReleasedFromWarehouse:
		s.ReleasedWarehouseAt = &now
	case ShipmentStatusInTransitToEngineer:
		// Set ETA if provided
		if eta != nil {
			s.ETAToEngineer = eta
		}
	case ShipmentStatusDelivered:
		s.DeliveredAt = &now
	}

	s.BeforeUpdate()
}

// IsDelivered returns true if the shipment has been delivered
func (s *Shipment) IsDelivered() bool {
	return s.Status == ShipmentStatusDelivered
}

// IsAtWarehouse returns true if the shipment is currently at the warehouse
func (s *Shipment) IsAtWarehouse() bool {
	return s.Status == ShipmentStatusAtWarehouse
}

// GetLaptopCount returns the number of laptops in this shipment
func (s *Shipment) GetLaptopCount() int {
	return len(s.Laptops)
}

// GetTrackingURL returns the courier's tracking URL for this shipment's tracking number
// Returns an empty string if the courier is not recognized or if courier name is empty
// Supports courier names with service types (e.g., "FedEx Express", "UPS Next Day Air")
func (s *Shipment) GetTrackingURL() string {
	if s.CourierName == "" || s.TrackingNumber == "" {
		if s.CourierName == "" {
			return ""
		}
	}

	// Normalize courier name to lowercase for comparison
	courierLower := strings.ToLower(strings.TrimSpace(s.CourierName))

	// Check for courier name using substring matching to support service types
	// e.g., "FedEx Express", "UPS Next Day Air", "DHL Express"
	var baseURL string
	if strings.Contains(courierLower, "ups") {
		baseURL = "https://www.ups.com/track?tracknum="
	} else if strings.Contains(courierLower, "dhl") {
		baseURL = "http://www.dhl.com/en/express/tracking.html?AWB="
	} else if strings.Contains(courierLower, "fedex") {
		baseURL = "https://www.fedex.com/fedextrack/?tracknumbers="
	} else {
		return ""
	}

	return baseURL + s.TrackingNumber
}

