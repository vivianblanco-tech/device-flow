package validator

import (
	"errors"
	"strings"
)

// WarehouseToEngineerFormInput represents the form data for warehouse-to-engineer shipments
type WarehouseToEngineerFormInput struct {
	LaptopID            int64
	SoftwareEngineerID  int64
	EngineerName        string
	EngineerEmail       string
	EngineerAddress     string
	EngineerCity        string
	EngineerCountry     string // Required for international addresses
	EngineerState       string // Optional for international addresses
	EngineerZip         string // Optional for international addresses (postal code)
	CourierName         string
	TrackingNumber      string
	JiraTicketNumber    string
	SpecialInstructions string
}

// ValidateWarehouseToEngineerForm validates the warehouse-to-engineer shipment form
func ValidateWarehouseToEngineerForm(input WarehouseToEngineerFormInput) error {
	// Laptop selection validation (REQUIRED)
	if input.LaptopID == 0 {
		return errors.New("laptop selection is required")
	}

	// Software engineer validation (REQUIRED)
	if input.SoftwareEngineerID == 0 && strings.TrimSpace(input.EngineerName) == "" {
		return errors.New("software engineer is required")
	}

	// Engineer address validation (international format for warehouse-to-engineer)
	if err := validateInternationalAddress(input.EngineerAddress, input.EngineerCity, input.EngineerCountry, input.EngineerState, input.EngineerZip); err != nil {
		return err
	}

	// Courier information validation (optional initially, required before shipping)
	// For form submission, we'll make it optional

	// JIRA ticket validation
	if err := validateJiraTicket(input.JiraTicketNumber); err != nil {
		return err
	}

	return nil
}

// validateInternationalAddress validates international address fields
// Required: address, city, country
// Optional: state, postal code
func validateInternationalAddress(address, city, country, state, postalCode string) error {
	// Validate address (REQUIRED)
	if strings.TrimSpace(address) == "" {
		return errors.New("address is required")
	}

	// Validate city (REQUIRED)
	if strings.TrimSpace(city) == "" {
		return errors.New("city is required")
	}

	// Validate country (REQUIRED)
	if strings.TrimSpace(country) == "" {
		return errors.New("country is required")
	}

	// State and postal code are optional for international addresses
	// No validation needed for these fields

	return nil
}

