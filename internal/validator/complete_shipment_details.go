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
	LaptopBrand            string
	LaptopModel            string
	LaptopRAMGB            string
	LaptopSSDGB            string
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

	// Laptop brand validation (REQUIRED)
	if strings.TrimSpace(input.LaptopBrand) == "" {
		return errors.New("laptop brand is required")
	}
	if len(input.LaptopBrand) > 100 {
		return errors.New("laptop brand must be less than 100 characters")
	}

	// Laptop model validation (REQUIRED)
	if strings.TrimSpace(input.LaptopModel) == "" {
		return errors.New("laptop model is required")
	}
	if len(input.LaptopModel) > 200 {
		return errors.New("laptop model must be less than 200 characters")
	}

	// Laptop RAM validation (REQUIRED)
	if strings.TrimSpace(input.LaptopRAMGB) == "" {
		return errors.New("laptop RAM is required")
	}
	if len(input.LaptopRAMGB) > 50 {
		return errors.New("laptop RAM must be less than 50 characters")
	}

	// Laptop SSD validation (REQUIRED)
	if strings.TrimSpace(input.LaptopSSDGB) == "" {
		return errors.New("laptop SSD is required")
	}
	if len(input.LaptopSSDGB) > 50 {
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

