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
	EngineerState       string
	EngineerZip         string
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

	// Engineer address validation
	if err := validateAddress(input.EngineerAddress, input.EngineerCity, input.EngineerState, input.EngineerZip); err != nil {
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

