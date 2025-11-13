package validator

import (
	"errors"
)

// EditShipmentDetailsInput represents form input for logistics editing shipment details
type EditShipmentDetailsInput struct {
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
	LaptopModel            string
	LaptopRAMGB            string
	LaptopSSDGB            string
	EngineerName           string
	IncludeAccessories     bool
	AccessoriesDescription string
}

// ValidateEditShipmentDetails validates form input for editing shipment details
// Note: JIRA ticket, Client Company, and Laptop Serial Number are NOT editable
func ValidateEditShipmentDetails(input EditShipmentDetailsInput) error {
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

	// Laptop model validation (optional for editing, validated if provided)
	if input.LaptopModel != "" && len(input.LaptopModel) > 200 {
		return errors.New("laptop model must be less than 200 characters")
	}

	// Laptop RAM validation (optional for editing, validated if provided)
	if input.LaptopRAMGB != "" && len(input.LaptopRAMGB) > 50 {
		return errors.New("laptop RAM must be less than 50 characters")
	}

	// Laptop SSD validation (optional for editing, validated if provided)
	if input.LaptopSSDGB != "" && len(input.LaptopSSDGB) > 50 {
		return errors.New("laptop SSD must be less than 50 characters")
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

