package validator

import (
	"errors"
	"strings"
)

// BulkToWarehouseFormInput represents the pickup form data for bulk shipments
type BulkToWarehouseFormInput struct {
	ClientCompanyID     int64
	ContactName         string
	ContactEmail        string
	ContactPhone        string
	PickupAddress       string
	PickupCity          string
	PickupState         string
	PickupZip           string
	PickupDate          string
	PickupTimeSlot      string
	JiraTicketNumber    string
	SpecialInstructions string

	// Laptop count (must be >= 2)
	NumberOfLaptops int

	// Bulk dimensions (REQUIRED)
	BulkLength float64
	BulkWidth  float64
	BulkHeight float64
	BulkWeight float64

	// Accessories (optional)
	IncludeAccessories     bool
	AccessoriesDescription string
}

// ValidateBulkToWarehouseForm validates the bulk to warehouse pickup form
func ValidateBulkToWarehouseForm(input BulkToWarehouseFormInput) error {
	// Client company validation
	if input.ClientCompanyID == 0 {
		return errors.New("client company is required")
	}

	// Contact information validation
	if err := validateContactInfo(input.ContactName, input.ContactEmail, input.ContactPhone); err != nil {
		return err
	}

	// Pickup address validation
	if err := validateAddress(input.PickupAddress, input.PickupCity, input.PickupState, input.PickupZip); err != nil {
		return err
	}

	// Pickup date and time validation
	if err := validatePickupDateTime(input.PickupDate, input.PickupTimeSlot); err != nil {
		return err
	}

	// JIRA ticket validation
	if err := validateJiraTicket(input.JiraTicketNumber); err != nil {
		return err
	}

	// Laptop count validation (must be >= 2 for bulk)
	if input.NumberOfLaptops < 2 {
		return errors.New("bulk shipments must have at least 2 laptops")
	}

	// Bulk dimensions validation (REQUIRED for bulk shipments)
	if input.BulkLength <= 0 || input.BulkWidth <= 0 || input.BulkHeight <= 0 || input.BulkWeight <= 0 {
		return errors.New("bulk dimensions (length, width, height, weight) are required and must be positive")
	}

	// Accessories validation
	if input.IncludeAccessories && strings.TrimSpace(input.AccessoriesDescription) == "" {
		return errors.New("accessories description is required when including accessories")
	}

	return nil
}

