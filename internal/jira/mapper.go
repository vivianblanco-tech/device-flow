package jira

import (
	"errors"
	"strings"
	"time"

	"github.com/yourusername/laptop-tracking-system/internal/models"
)

// ShipmentData represents shipment-relevant data extracted from a JIRA ticket
type ShipmentData struct {
	JiraTicketKey string
	Summary       string
	Description   string
	Status        string
	Assignee      string
	Created       time.Time
	Updated       time.Time
}

// CustomFields represents custom field data from JIRA
type CustomFields struct {
	SerialNumber  string
	EngineerEmail string
	ClientCompany string
}

// MapTicketToShipmentData extracts shipment-relevant data from a JIRA ticket
func MapTicketToShipmentData(ticket *Ticket) (*ShipmentData, error) {
	if ticket == nil {
		return nil, errors.New("ticket cannot be nil")
	}

	// Parse timestamps
	created, _ := ParseJiraTimestamp(ticket.Created)
	updated, _ := ParseJiraTimestamp(ticket.Updated)

	return &ShipmentData{
		JiraTicketKey: ticket.Key,
		Summary:       ticket.Summary,
		Description:   ticket.Description,
		Status:        ticket.Status,
		Assignee:      ticket.Assignee,
		Created:       created,
		Updated:       updated,
	}, nil
}

// ExtractCustomFields extracts custom field data from JIRA ticket response
func ExtractCustomFields(ticketData map[string]interface{}) (*CustomFields, error) {
	if ticketData == nil {
		return nil, errors.New("ticket data cannot be nil")
	}

	fields, ok := ticketData["fields"].(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid fields structure")
	}

	customFields := &CustomFields{}

	// Extract serial number (customfield_10001)
	if val, ok := fields["customfield_10001"].(string); ok {
		customFields.SerialNumber = val
	}

	// Extract engineer email (customfield_10002)
	if val, ok := fields["customfield_10002"].(string); ok {
		customFields.EngineerEmail = val
	}

	// Extract client company (customfield_10003)
	if val, ok := fields["customfield_10003"].(string); ok {
		customFields.ClientCompany = val
	}

	return customFields, nil
}

// MapJiraStatusToShipmentStatus maps JIRA status to internal shipment status
func MapJiraStatusToShipmentStatus(jiraStatus string) models.ShipmentStatus {
	// Normalize status string (lowercase, trim spaces)
	status := strings.ToLower(strings.TrimSpace(jiraStatus))

	// Map JIRA statuses to shipment statuses
	switch status {
	case "to do", "pending pickup":
		return models.ShipmentStatusPendingPickup
	case "picked up":
		return models.ShipmentStatusPickedUpFromClient
	case "in transit to warehouse":
		return models.ShipmentStatusInTransitToWarehouse
	case "at warehouse":
		return models.ShipmentStatusAtWarehouse
	case "released from warehouse":
		return models.ShipmentStatusReleasedFromWarehouse
	case "in transit to engineer":
		return models.ShipmentStatusInTransitToEngineer
	case "delivered", "done":
		return models.ShipmentStatusDelivered
	default:
		// Default to pending pickup for unknown statuses
		return models.ShipmentStatusPendingPickup
	}
}

// CreateShipmentFromTicket creates a shipment model from JIRA ticket and custom fields
func CreateShipmentFromTicket(ticket *Ticket, customFields *CustomFields) (*models.Shipment, error) {
	if ticket == nil {
		return nil, errors.New("ticket cannot be nil")
	}

	// Map JIRA status to shipment status
	status := MapJiraStatusToShipmentStatus(ticket.Status)

	// Parse created timestamp
	createdAt, err := ParseJiraTimestamp(ticket.Created)
	if err != nil {
		createdAt = time.Now()
	}

	// Build notes from ticket information
	notes := ticket.Summary
	if ticket.Description != "" {
		notes += "\n\nDescription: " + ticket.Description
	}
	if customFields != nil && customFields.SerialNumber != "" {
		notes += "\n\nSerial Number: " + customFields.SerialNumber
	}
	notes += "\n\nJIRA Ticket: " + ticket.Key

	// Create the shipment
	shipment := &models.Shipment{
		Status:    status,
		Notes:     notes,
		CreatedAt: createdAt,
		UpdatedAt: time.Now(),
	}

	return shipment, nil
}

// ParseJiraTimestamp parses JIRA's timestamp format
func ParseJiraTimestamp(timestamp string) (time.Time, error) {
	if timestamp == "" {
		return time.Time{}, errors.New("timestamp is empty")
	}

	// JIRA uses ISO 8601 format: 2023-10-01T10:00:00.000+0000
	layout := "2006-01-02T15:04:05.000-0700"
	
	parsedTime, err := time.Parse(layout, timestamp)
	if err != nil {
		return time.Time{}, err
	}

	return parsedTime, nil
}

