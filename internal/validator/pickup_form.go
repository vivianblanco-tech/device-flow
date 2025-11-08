package validator

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

// PickupFormInput represents the input data for a pickup form
type PickupFormInput struct {
	ClientCompanyID     int64  `json:"client_company_id"`
	ContactName         string `json:"contact_name"`
	ContactEmail        string `json:"contact_email"`
	ContactPhone        string `json:"contact_phone"`
	PickupAddress       string `json:"pickup_address"`
	PickupCity          string `json:"pickup_city"`
	PickupState         string `json:"pickup_state"`
	PickupZip           string `json:"pickup_zip"`
	PickupDate          string `json:"pickup_date"`
	PickupTimeSlot      string `json:"pickup_time_slot"`
	NumberOfLaptops     int    `json:"number_of_laptops"`
	JiraTicketNumber    string `json:"jira_ticket_number"`
	SpecialInstructions string `json:"special_instructions"`
}

// ValidatePickupForm validates the pickup form input
func ValidatePickupForm(input PickupFormInput) error {
	// Validate client company ID
	if input.ClientCompanyID == 0 {
		return errors.New("client company ID is required")
	}

	// Validate contact name
	if strings.TrimSpace(input.ContactName) == "" {
		return errors.New("contact name is required")
	}

	// Validate contact email
	if strings.TrimSpace(input.ContactEmail) == "" {
		return errors.New("contact email is required")
	}
	if !isValidEmail(input.ContactEmail) {
		return errors.New("invalid email format")
	}

	// Validate contact phone
	if strings.TrimSpace(input.ContactPhone) == "" {
		return errors.New("contact phone is required")
	}

	// Validate pickup address
	if strings.TrimSpace(input.PickupAddress) == "" {
		return errors.New("pickup address is required")
	}

	// Validate pickup city
	if strings.TrimSpace(input.PickupCity) == "" {
		return errors.New("pickup city is required")
	}

	// Validate pickup state
	if strings.TrimSpace(input.PickupState) == "" {
		return errors.New("pickup state is required")
	}
	if !isValidUSState(input.PickupState) {
		return errors.New("invalid US state code")
	}

	// Validate pickup ZIP code
	if strings.TrimSpace(input.PickupZip) == "" {
		return errors.New("pickup ZIP code is required")
	}
	if !isValidZipCode(input.PickupZip) {
		return errors.New("ZIP code must be 5 digits")
	}

	// Validate pickup date
	if strings.TrimSpace(input.PickupDate) == "" {
		return errors.New("pickup date is required")
	}

	pickupDate, err := time.Parse("2006-01-02", input.PickupDate)
	if err != nil {
		return errors.New("invalid date format")
	}

	// Check if date is in the future (allowing same day)
	today := time.Now().Truncate(24 * time.Hour)
	if pickupDate.Before(today) {
		return errors.New("pickup date must be in the future")
	}

	// Validate time slot
	if strings.TrimSpace(input.PickupTimeSlot) == "" {
		return errors.New("pickup time slot is required")
	}
	if !isValidTimeSlot(input.PickupTimeSlot) {
		return errors.New("invalid time slot")
	}

	// Validate number of laptops
	if input.NumberOfLaptops < 1 {
		return errors.New("number of laptops must be at least 1")
	}

	// Validate JIRA ticket number
	if strings.TrimSpace(input.JiraTicketNumber) == "" {
		return errors.New("JIRA ticket number is required")
	}
	if !isValidJiraTicketFormat(input.JiraTicketNumber) {
		return errors.New("JIRA ticket number must be in format PROJECT-NUMBER (e.g., SCOP-67702)")
	}

	return nil
}

// isValidJiraTicketFormat validates the JIRA ticket number format (PROJECT-NUMBER)
func isValidJiraTicketFormat(ticket string) bool {
	// Pattern: uppercase letters, dash, digits
	// Example: SCOP-67702, PROJECT-12345, TEST-100
	pattern := `^[A-Z]+\-[0-9]+$`
	matched, _ := regexp.MatchString(pattern, ticket)
	return matched
}

// isValidEmail validates email format using a simple regex
func isValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	if email == "" {
		return false
	}

	// Simple email validation regex
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// isValidTimeSlot checks if the time slot is valid
func isValidTimeSlot(slot string) bool {
	validSlots := []string{"morning", "afternoon", "evening"}
	for _, valid := range validSlots {
		if slot == valid {
			return true
		}
	}
	return false
}

// isValidZipCode validates US ZIP code format (5 digits)
func isValidZipCode(zip string) bool {
	zip = strings.TrimSpace(zip)
	if len(zip) != 5 {
		return false
	}
	// Check if all characters are digits
	pattern := `^[0-9]{5}$`
	matched, _ := regexp.MatchString(pattern, zip)
	return matched
}

// isValidUSState validates US state code (2-letter abbreviation)
func isValidUSState(state string) bool {
	validStates := []string{
		"AL", "AK", "AZ", "AR", "CA", "CO", "CT", "DE", "FL", "GA",
		"HI", "ID", "IL", "IN", "IA", "KS", "KY", "LA", "ME", "MD",
		"MA", "MI", "MN", "MS", "MO", "MT", "NE", "NV", "NH", "NJ",
		"NM", "NY", "NC", "ND", "OH", "OK", "OR", "PA", "RI", "SC",
		"SD", "TN", "TX", "UT", "VT", "VA", "WA", "WV", "WI", "WY",
		"DC", // District of Columbia
	}
	state = strings.TrimSpace(state)
	for _, valid := range validStates {
		if state == valid {
			return true
		}
	}
	return false
}

