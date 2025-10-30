package models

import (
	"errors"
	"time"
)

// LaptopStatus represents the status of a laptop in the system
type LaptopStatus string

// Laptop status constants
const (
	LaptopStatusAvailable              LaptopStatus = "available"
	LaptopStatusInTransitToWarehouse   LaptopStatus = "in_transit_to_warehouse"
	LaptopStatusAtWarehouse            LaptopStatus = "at_warehouse"
	LaptopStatusInTransitToEngineer    LaptopStatus = "in_transit_to_engineer"
	LaptopStatusDelivered              LaptopStatus = "delivered"
	LaptopStatusRetired                LaptopStatus = "retired"
)

// Laptop represents a laptop device in the inventory
type Laptop struct {
	ID           int64        `json:"id" db:"id"`
	SerialNumber string       `json:"serial_number" db:"serial_number"`
	Brand        string       `json:"brand,omitempty" db:"brand"`
	Model        string       `json:"model,omitempty" db:"model"`
	Specs        string       `json:"specs,omitempty" db:"specs"`
	Status       LaptopStatus `json:"status" db:"status"`
	CreatedAt    time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at" db:"updated_at"`
}

// Validate validates the Laptop model
func (l *Laptop) Validate() error {
	// Serial number validation
	if l.SerialNumber == "" {
		return errors.New("serial number is required")
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

// UpdateStatus updates the laptop status
func (l *Laptop) UpdateStatus(status LaptopStatus) {
	l.Status = status
	l.BeforeUpdate()
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

	if l.Specs != "" {
		desc += " (" + l.Specs + ")"
	}

	return desc
}

