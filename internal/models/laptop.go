package models

import (
	"errors"
	"time"
)

// LaptopStatus represents the status of a laptop in the system
type LaptopStatus string

// Laptop status constants
const (
	LaptopStatusAvailable            LaptopStatus = "available"
	LaptopStatusInTransitToWarehouse LaptopStatus = "in_transit_to_warehouse"
	LaptopStatusAtWarehouse          LaptopStatus = "at_warehouse"
	LaptopStatusInTransitToEngineer  LaptopStatus = "in_transit_to_engineer"
	LaptopStatusDelivered            LaptopStatus = "delivered"
	LaptopStatusRetired              LaptopStatus = "retired"
)

// Laptop represents a laptop device in the inventory
type Laptop struct {
	ID                 int64        `json:"id" db:"id"`
	SerialNumber       string       `json:"serial_number" db:"serial_number"`
	SKU                string       `json:"sku,omitempty" db:"sku"`
	Brand              string       `json:"brand,omitempty" db:"brand"`
	Model              string       `json:"model" db:"model"`
	CPU                string       `json:"cpu" db:"cpu"`
	RAMGB              string       `json:"ram_gb" db:"ram_gb"`
	SSDGB              string       `json:"ssd_gb" db:"ssd_gb"`
	Status             LaptopStatus `json:"status" db:"status"`
	ClientCompanyID    *int64       `json:"client_company_id,omitempty" db:"client_company_id"`
	SoftwareEngineerID *int64       `json:"software_engineer_id,omitempty" db:"software_engineer_id"`
	CreatedAt          time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time    `json:"updated_at" db:"updated_at"`

	// Relations (not stored in DB directly, populated by queries with joins)
	ClientCompanyName    string `json:"client_company_name,omitempty" db:"client_company_name"`
	SoftwareEngineerName string `json:"software_engineer_name,omitempty" db:"software_engineer_name"`
	EmployeeID           string `json:"employee_id,omitempty" db:"employee_id"`
	
	// Reception Report information (populated by queries with joins for warehouse users)
	HasReceptionReport     bool    `json:"has_reception_report" db:"has_reception_report"`
	ReceptionReportID      *int64  `json:"reception_report_id,omitempty" db:"reception_report_id"`
	ReceptionReportStatus  string  `json:"reception_report_status,omitempty" db:"reception_report_status"`
}

// GenerateAndSetSKU generates and sets the SKU for the laptop if it's not already set
func (l *Laptop) GenerateAndSetSKU() {
	// Only generate if SKU is empty
	if l.SKU == "" {
		l.SKU = GenerateSKU(l.Model, l.CPU, l.RAMGB, l.SSDGB)
	}
}

// Validate validates the Laptop model
func (l *Laptop) Validate() error {
	// Serial number validation
	if l.SerialNumber == "" {
		return errors.New("serial number is required")
	}

	// Brand validation (required)
	if l.Brand == "" {
		return errors.New("brand is required")
	}

	// Model validation (required)
	if l.Model == "" {
		return errors.New("laptop model is required")
	}

	// CPU validation (required)
	if l.CPU == "" {
		return errors.New("laptop CPU is required")
	}

	// RAM validation (required)
	if l.RAMGB == "" {
		return errors.New("laptop RAM is required")
	}

	// SSD validation (required)
	if l.SSDGB == "" {
		return errors.New("laptop SSD is required")
	}

	// Client company validation (required)
	if l.ClientCompanyID == nil {
		return errors.New("client company is required")
	}

	// Status validation
	if l.Status == "" {
		return errors.New("status is required")
	}
	if !IsValidLaptopStatus(l.Status) {
		return errors.New("invalid status")
	}

	return nil
}

// ValidateStatusChange validates that a laptop status change is allowed based on business rules
func (l *Laptop) ValidateStatusChange(receptionReport *ReceptionReport) error {
	// Rule 1: Cannot set status to "available" without an approved reception report
	if l.Status == LaptopStatusAvailable {
		if receptionReport == nil || !receptionReport.IsApproved() {
			return errors.New("cannot set status to available without an approved reception report")
		}
	}

	// Rule 2: Cannot set status to "in transit to engineer" without an assigned engineer
	if l.Status == LaptopStatusInTransitToEngineer {
		if l.SoftwareEngineerID == nil {
			return errors.New("cannot set status to in transit to engineer without an assigned engineer")
		}
	}

	return nil
}

// IsValidLaptopStatus checks if a given status is valid
func IsValidLaptopStatus(status LaptopStatus) bool {
	switch status {
	case LaptopStatusAvailable,
		LaptopStatusInTransitToWarehouse,
		LaptopStatusAtWarehouse,
		LaptopStatusInTransitToEngineer,
		LaptopStatusDelivered,
		LaptopStatusRetired:
		return true
	}
	return false
}

// TableName returns the table name for the Laptop model
func (l *Laptop) TableName() string {
	return "laptops"
}

// BeforeCreate sets the timestamps before creating a laptop
func (l *Laptop) BeforeCreate() {
	now := time.Now()
	l.CreatedAt = now
	l.UpdatedAt = now
}

// BeforeUpdate sets the updated_at timestamp before updating a laptop
func (l *Laptop) BeforeUpdate() {
	l.UpdatedAt = time.Now()
}

// IsAvailable returns true if the laptop is available for assignment
func (l *Laptop) IsAvailable() bool {
	return l.Status == LaptopStatusAvailable
}

// IsAvailableForWarehouseShipment checks if laptop can be used for warehouse-to-engineer shipment
// Requirements: must be available/at_warehouse, have reception report, and not in active shipment
func (l *Laptop) IsAvailableForWarehouseShipment(hasReceptionReport bool, inActiveShipment bool) bool {
	// Must be available or at warehouse
	if l.Status != LaptopStatusAvailable && l.Status != LaptopStatusAtWarehouse {
		return false
	}
	
	// Must have completed reception report
	if !hasReceptionReport {
		return false
	}
	
	// Must not be in an active shipment
	if inActiveShipment {
		return false
	}
	
	return true
}

// UpdateStatus updates the laptop status
func (l *Laptop) UpdateStatus(status LaptopStatus) {
	l.Status = status
	l.BeforeUpdate()
}

// CanChangeToAvailable checks if a laptop can be changed to available status
// Requirements: must be at_warehouse and have an approved reception report
func (l *Laptop) CanChangeToAvailable(receptionReport *ReceptionReport) bool {
	// Must currently be at warehouse
	if l.Status != LaptopStatusAtWarehouse {
		return false
	}
	
	// Must have a reception report
	if receptionReport == nil {
		return false
	}
	
	// Reception report must be approved
	if !receptionReport.IsApproved() {
		return false
	}
	
	// Reception report must be for this laptop
	if receptionReport.LaptopID != l.ID {
		return false
	}
	
	return true
}

// GetFullDescription returns a full description of the laptop
func (l *Laptop) GetFullDescription() string {
	if l.Brand == "" && l.Model == "" {
		return "Unknown"
	}

	desc := l.Brand
	if l.Model != "" {
		if desc != "" {
			desc += " "
		}
		desc += l.Model
	}

	// Add CPU, RAM and SSD specs
	if l.CPU != "" || l.RAMGB != "" || l.SSDGB != "" {
		specs := ""
		if l.CPU != "" {
			specs = l.CPU
		}
		if l.RAMGB != "" {
			if specs != "" {
				specs += ", "
			}
			specs += l.RAMGB + " RAM"
		}
		if l.SSDGB != "" {
			if specs != "" {
				specs += ", "
			}
			specs += l.SSDGB + " SSD"
		}
		if specs != "" {
			desc += " (" + specs + ")"
		}
	}

	return desc
}

// GetLaptopStatusDisplayName returns the user-friendly display name for a laptop status
func GetLaptopStatusDisplayName(status LaptopStatus) string {
	switch status {
	case LaptopStatusAvailable:
		return "Available at Warehouse"
	case LaptopStatusAtWarehouse:
		return "Received at Warehouse"
	case LaptopStatusInTransitToWarehouse:
		return "In Transit To Warehouse"
	case LaptopStatusInTransitToEngineer:
		return "In Transit To Engineer"
	case LaptopStatusDelivered:
		return "Delivered"
	case LaptopStatusRetired:
		return "Retired"
	default:
		return string(status)
	}
}

// GetLaptopStatusesInOrder returns all laptop statuses in logical order
func GetLaptopStatusesInOrder() []LaptopStatus {
	return []LaptopStatus{
		LaptopStatusInTransitToWarehouse,
		LaptopStatusAtWarehouse, // "Received at Warehouse"
		LaptopStatusAvailable,   // "Available at Warehouse"
		LaptopStatusInTransitToEngineer,
		LaptopStatusDelivered,
		LaptopStatusRetired,
	}
}

// GetLaptopStatusesForNewLaptop returns statuses appropriate for newly added laptops
// Only "Received at Warehouse" is allowed to ensure all laptops go through
// the reception report process before becoming available
func GetLaptopStatusesForNewLaptop() []LaptopStatus {
	return []LaptopStatus{
		LaptopStatusAtWarehouse, // "Received at Warehouse" - the only status for new laptops
	}
}

// GetAllowedStatusesForRole returns laptop statuses that a user role is allowed to see
func GetAllowedStatusesForRole(role UserRole) []LaptopStatus {
	// Warehouse users only see warehouse-relevant statuses
	if role == RoleWarehouse {
		return []LaptopStatus{
			LaptopStatusInTransitToWarehouse,
			LaptopStatusAtWarehouse,
			LaptopStatusAvailable,
		}
	}

	// All other roles see all statuses
	return GetLaptopStatusesInOrder()
}
