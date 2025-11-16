package validator

import (
	"errors"
	"strings"
	"time"
)

// SingleFullJourneyFormInput represents the pickup form data for single full journey shipments
type SingleFullJourneyFormInput struct {
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

	// Laptop information (required)
	LaptopSerialNumber string
	LaptopBrand        string
	LaptopModel        string
	LaptopCPU          string
	LaptopRAMGB        string
	LaptopSSDGB        string

	// Engineer assignment (optional - can be assigned later)
	EngineerName string

	// Accessories (optional)
	IncludeAccessories     bool
	AccessoriesDescription string
}

// ValidateSingleFullJourneyForm validates the single full journey pickup form
func ValidateSingleFullJourneyForm(input SingleFullJourneyFormInput) error {
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

	// Laptop serial number validation (REQUIRED)
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

	// Laptop CPU validation (REQUIRED)
	if strings.TrimSpace(input.LaptopCPU) == "" {
		return errors.New("laptop CPU is required")
	}
	if len(input.LaptopCPU) > 200 {
		return errors.New("laptop CPU must be less than 200 characters")
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

	// Engineer name (optional - can be assigned later)
	// No validation - can be empty

	// Accessories validation
	if input.IncludeAccessories && strings.TrimSpace(input.AccessoriesDescription) == "" {
		return errors.New("accessories description is required when including accessories")
	}

	return nil
}

// validateContactInfo validates contact information fields
func validateContactInfo(name, email, phone string) error {
	// Validate contact name
	if strings.TrimSpace(name) == "" {
		return errors.New("contact name is required")
	}

	// Validate contact email
	if strings.TrimSpace(email) == "" {
		return errors.New("contact email is required")
	}
	if !isValidEmail(email) {
		return errors.New("invalid email format")
	}

	// Validate contact phone
	if strings.TrimSpace(phone) == "" {
		return errors.New("contact phone is required")
	}

	return nil
}

// validateAddress validates address fields
func validateAddress(address, city, state, zip string) error {
	// Validate address
	if strings.TrimSpace(address) == "" {
		return errors.New("address is required")
	}

	// Validate city
	if strings.TrimSpace(city) == "" {
		return errors.New("city is required")
	}

	// Validate state
	if strings.TrimSpace(state) == "" {
		return errors.New("state is required")
	}
	if !isValidUSState(state) {
		return errors.New("invalid US state code")
	}

	// Validate ZIP code
	if strings.TrimSpace(zip) == "" {
		return errors.New("ZIP code is required")
	}
	if !isValidZipCode(zip) {
		return errors.New("ZIP code must be 5 digits")
	}

	return nil
}

// validatePickupDateTime validates pickup date and time slot
func validatePickupDateTime(date, timeSlot string) error {
	// Validate pickup date
	if strings.TrimSpace(date) == "" {
		return errors.New("pickup date is required")
	}

	pickupDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return errors.New("invalid date format")
	}

	// Check if date is in the future (allowing same day)
	today := time.Now().Truncate(24 * time.Hour)
	if pickupDate.Before(today) {
		return errors.New("pickup date must be in the future")
	}

	// Validate time slot
	if strings.TrimSpace(timeSlot) == "" {
		return errors.New("pickup time slot is required")
	}
	if !isValidTimeSlot(timeSlot) {
		return errors.New("invalid time slot")
	}

	return nil
}

// validateJiraTicket validates JIRA ticket number
func validateJiraTicket(ticket string) error {
	if strings.TrimSpace(ticket) == "" {
		return errors.New("JIRA ticket number is required")
	}
	if !isValidJiraTicketFormat(ticket) {
		return errors.New("JIRA ticket number must be in format PROJECT-NUMBER (e.g., SCOP-67702)")
	}
	return nil
}

