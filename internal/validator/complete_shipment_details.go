package validator

import (
	"errors"
	"strings"
)

// CompleteShipmentDetailsInput represents form input for client completing shipment details via magic link
type CompleteShipmentDetailsInput struct {
	ShipmentID             int64
	ContactName            string
	ContactEmail           string
	ContactPhone           string
	PickupAddress          string
	PickupCity             string
	PickupState            string
	PickupZip              string
	PickupDate             string
	PickupTimeSlot         string
	SpecialInstructions    string
	LaptopSerialNumber     string
	LaptopSpecs            string
	EngineerName           string
	IncludeAccessories     bool
	AccessoriesDescription string
}

// ValidateCompleteShipmentDetails validates form input for completing shipment details
// Note: ClientCompanyID and JiraTicketNumber are NOT validated here (already set by logistics)
func ValidateCompleteShipmentDetails(input CompleteShipmentDetailsInput) error {
	// Shipment ID validation
	if input.ShipmentID == 0 {
		return errors.New("shipment ID is required")
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

	// Laptop serial number validation (REQUIRED for completing shipment details)
	if strings.TrimSpace(input.LaptopSerialNumber) == "" {
		return errors.New("laptop serial number is required")
	}

	// Laptop specs validation (optional but validated if provided)
	// No strict validation - just check length
	if len(input.LaptopSpecs) > 500 {
		return errors.New("laptop specs must be less than 500 characters")
	}

	// Engineer name validation (optional)
	if input.EngineerName != "" && len(input.EngineerName) > 100 {
		return errors.New("engineer name must be less than 100 characters")
	}

	// Special instructions validation (optional)
	if len(input.SpecialInstructions) > 1000 {
		return errors.New("special instructions must be less than 1000 characters")
	}

	// Accessories validation
	if input.IncludeAccessories && len(input.AccessoriesDescription) > 500 {
		return errors.New("accessories description must be less than 500 characters")
	}

	return nil
}

