package jira

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/yourusername/laptop-tracking-system/internal/models"
)

// CreateTicketRequest represents a request to create a JIRA ticket
type CreateTicketRequest struct {
	ProjectKey  string
	Summary     string
	Description string
	IssueType   string
	Labels      []string
	CustomFields map[string]interface{}
}

// CreateTicketResponse represents the response from creating a JIRA ticket
type CreateTicketResponse struct {
	ID   string `json:"id"`
	Key  string `json:"key"`
	Self string `json:"self"`
}

// CreateTicket creates a new JIRA ticket
func (c *Client) CreateTicket(request *CreateTicketRequest, accessToken string) (*CreateTicketResponse, error) {
	if accessToken == "" {
		return nil, errors.New("access token is required")
	}
	if request == nil {
		return nil, errors.New("request cannot be nil")
	}

	// Validate request
	if request.ProjectKey == "" {
		return nil, errors.New("project key is required")
	}
	if request.Summary == "" {
		return nil, errors.New("summary is required")
	}
	if request.IssueType == "" {
		return nil, errors.New("issue type is required")
	}

	// Build the request payload
	payload := map[string]interface{}{
		"fields": map[string]interface{}{
			"project": map[string]string{
				"key": request.ProjectKey,
			},
			"summary": request.Summary,
			"issuetype": map[string]string{
				"name": request.IssueType,
			},
		},
	}

	// Add optional fields
	fields := payload["fields"].(map[string]interface{})
	
	if request.Description != "" {
		fields["description"] = map[string]interface{}{
			"type": "doc",
			"version": 1,
			"content": []map[string]interface{}{
				{
					"type": "paragraph",
					"content": []map[string]interface{}{
						{
							"type": "text",
							"text": request.Description,
						},
					},
				},
			},
		}
	}

	if len(request.Labels) > 0 {
		fields["labels"] = request.Labels
	}

	// Add custom fields
	for key, value := range request.CustomFields {
		fields[key] = value
	}

	// Marshal payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Create the HTTP request
	url := fmt.Sprintf("%s/rest/api/3/issue", c.config.URL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Make the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create ticket: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check response status
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var response CreateTicketResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// UpdateTicketStatus updates a JIRA ticket's status by transitioning it
func (c *Client) UpdateTicketStatus(issueKey, status, accessToken string) error {
	if accessToken == "" {
		return errors.New("access token is required")
	}
	if issueKey == "" {
		return errors.New("issue key is required")
	}
	if status == "" {
		return errors.New("status is required")
	}

	// In a real implementation, we would need to:
	// 1. Get available transitions for the issue
	// 2. Find the transition ID that matches the desired status
	// 3. Execute the transition
	// For now, we'll use a simplified approach

	// Build the request payload
	payload := map[string]interface{}{
		"transition": map[string]interface{}{
			"name": status,
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Create the HTTP request
	url := fmt.Sprintf("%s/rest/api/3/issue/%s/transitions", c.config.URL, issueKey)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Make the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

// AddComment adds a comment to a JIRA ticket
func (c *Client) AddComment(issueKey, comment, accessToken string) error {
	if accessToken == "" {
		return errors.New("access token is required")
	}
	if issueKey == "" {
		return errors.New("issue key is required")
	}
	if comment == "" {
		return errors.New("comment is required")
	}

	// Build the request payload
	payload := map[string]interface{}{
		"body": map[string]interface{}{
			"type": "doc",
			"version": 1,
			"content": []map[string]interface{}{
				{
					"type": "paragraph",
					"content": []map[string]interface{}{
						{
							"type": "text",
							"text": comment,
						},
					},
				},
			},
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Create the HTTP request
	url := fmt.Sprintf("%s/rest/api/3/issue/%s/comment", c.config.URL, issueKey)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Make the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to add comment: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

// BuildTicketFromShipment builds a JIRA ticket request from shipment data
func BuildTicketFromShipment(shipment *models.Shipment, clientCompany *models.ClientCompany, laptops []models.Laptop, projectKey string) *CreateTicketRequest {
	// Build summary
	summary := fmt.Sprintf("Hardware Deployment - Shipment #%d", shipment.ID)
	if clientCompany != nil {
		summary = fmt.Sprintf("Hardware Deployment for %s - Shipment #%d", clientCompany.Name, shipment.ID)
	}

	// Build description
	var descParts []string
	descParts = append(descParts, fmt.Sprintf("Shipment ID: %d", shipment.ID))
	
	if clientCompany != nil {
		descParts = append(descParts, fmt.Sprintf("Client Company: %s", clientCompany.Name))
	}

	if len(laptops) > 0 {
		descParts = append(descParts, "\nDevices:")
		for _, laptop := range laptops {
			descParts = append(descParts, fmt.Sprintf("- %s %s (SN: %s)", laptop.Brand, laptop.Model, laptop.SerialNumber))
		}
	}

	if shipment.Notes != "" {
		descParts = append(descParts, fmt.Sprintf("\nNotes: %s", shipment.Notes))
	}

	description := strings.Join(descParts, "\n")

	// Build custom fields map
	customFields := make(map[string]interface{})
	if len(laptops) > 0 {
		// Add first laptop serial number to custom field
		customFields["customfield_10001"] = laptops[0].SerialNumber
	}

	return &CreateTicketRequest{
		ProjectKey:   projectKey,
		Summary:      summary,
		Description:  description,
		IssueType:    "Task",
		Labels:       []string{"hardware-deployment", "laptop-tracking"},
		CustomFields: customFields,
	}
}

// SyncShipmentStatusToJira syncs shipment status changes to JIRA
func (c *Client) SyncShipmentStatusToJira(issueKey string, shipment *models.Shipment, accessToken string) error {
	if issueKey == "" {
		return errors.New("issue key is required")
	}
	if shipment == nil {
		return errors.New("shipment cannot be nil")
	}

	// Map shipment status to JIRA status
	var jiraStatus string
	switch shipment.Status {
	case models.ShipmentStatusPendingPickup:
		jiraStatus = "Pending Pickup"
	case models.ShipmentStatusPickedUpFromClient:
		jiraStatus = "Picked Up"
	case models.ShipmentStatusInTransitToWarehouse:
		jiraStatus = "In Transit to Warehouse"
	case models.ShipmentStatusAtWarehouse:
		jiraStatus = "At Warehouse"
	case models.ShipmentStatusReleasedFromWarehouse:
		jiraStatus = "Released from Warehouse"
	case models.ShipmentStatusInTransitToEngineer:
		jiraStatus = "In Transit to Engineer"
	case models.ShipmentStatusDelivered:
		jiraStatus = "Delivered"
	default:
		jiraStatus = "Pending Pickup"
	}

	// Update the ticket status
	if err := c.UpdateTicketStatus(issueKey, jiraStatus, accessToken); err != nil {
		return fmt.Errorf("failed to update ticket status: %w", err)
	}

	// Add a comment with the status update
	comment := fmt.Sprintf("Shipment status updated to: %s", shipment.Status)
	if err := c.AddComment(issueKey, comment, accessToken); err != nil {
		// Don't fail if comment fails, just log it
		// In production, this would use proper logging
		return nil
	}

	return nil
}

