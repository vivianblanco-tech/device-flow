package jira

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/laptop-tracking-system/internal/models"
)

// 🟥 RED: Test for creating a JIRA ticket from shipment
// This test verifies that we can create a new JIRA ticket from shipment data
func TestClient_CreateTicket(t *testing.T) {
	// Create a test server that simulates JIRA API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST request, got %s", r.Method)
		}
		if r.URL.Path == "/rest/api/3/issue" {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{
				"id": "10000",
				"key": "PROJ-126",
				"self": "https://example.atlassian.net/rest/api/3/issue/10000"
			}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	config := Config{
		URL:          server.URL,
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURL:  "http://localhost:8080/auth/jira/callback",
	}

	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("expected no error creating client, got %v", err)
	}

	// Create JIRA ticket
	ticketRequest := &CreateTicketRequest{
		ProjectKey: "PROJ",
		Summary:    "Hardware deployment for shipment #1",
		IssueType:  "Task",
	}

	response, err := client.CreateTicket(ticketRequest, "mock-access-token")
	if err != nil {
		t.Fatalf("expected no error creating ticket, got %v", err)
	}

	// Verify response
	if response.Key != "PROJ-126" {
		t.Errorf("expected ticket key PROJ-126, got %s", response.Key)
	}
}

// 🟥 RED: Test for creating a JIRA ticket without access token
// This test verifies that ticket creation requires an access token
func TestClient_CreateTicket_NoToken(t *testing.T) {
	config := Config{
		URL:          "https://bairesdev.atlassian.net",
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURL:  "http://localhost:8080/auth/jira/callback",
	}

	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("expected no error creating client, got %v", err)
	}

	ticketRequest := &CreateTicketRequest{
		ProjectKey: "PROJ",
		Summary:    "Test ticket",
		IssueType:  "Task",
	}

	_, err = client.CreateTicket(ticketRequest, "")
	if err == nil {
		t.Error("expected error for missing token, got nil")
	}
}

// 🟥 RED: Test for creating a JIRA ticket with invalid request
// This test verifies proper validation of ticket creation request
func TestClient_CreateTicket_InvalidRequest(t *testing.T) {
	config := Config{
		URL:          "https://bairesdev.atlassian.net",
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURL:  "http://localhost:8080/auth/jira/callback",
	}

	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("expected no error creating client, got %v", err)
	}

	tests := []struct {
		name    string
		request *CreateTicketRequest
	}{
		{
			name: "missing project key",
			request: &CreateTicketRequest{
				Summary:   "Test ticket",
				IssueType: "Task",
			},
		},
		{
			name: "missing summary",
			request: &CreateTicketRequest{
				ProjectKey: "PROJ",
				IssueType:  "Task",
			},
		},
		{
			name: "missing issue type",
			request: &CreateTicketRequest{
				ProjectKey: "PROJ",
				Summary:    "Test ticket",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.CreateTicket(tt.request, "mock-token")
			if err == nil {
				t.Error("expected error for invalid request, got nil")
			}
		})
	}
}

// 🟥 RED: Test for updating a JIRA ticket status
// This test verifies that we can update a ticket's status
func TestClient_UpdateTicketStatus(t *testing.T) {
	// Create a test server that simulates JIRA API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST request, got %s", r.Method)
		}
		if r.URL.Path == "/rest/api/3/issue/PROJ-123/transitions" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	config := Config{
		URL:          server.URL,
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURL:  "http://localhost:8080/auth/jira/callback",
	}

	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("expected no error creating client, got %v", err)
	}

	// Update ticket status
	err = client.UpdateTicketStatus("PROJ-123", "In Progress", "mock-access-token")
	if err != nil {
		t.Fatalf("expected no error updating status, got %v", err)
	}
}

// 🟥 RED: Test for updating ticket with comment
// This test verifies that we can add comments to a ticket
func TestClient_AddComment(t *testing.T) {
	// Create a test server that simulates JIRA API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST request, got %s", r.Method)
		}
		if r.URL.Path == "/rest/api/3/issue/PROJ-123/comment" {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{
				"id": "10001",
				"body": "Shipment has been picked up from client"
			}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	config := Config{
		URL:          server.URL,
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURL:  "http://localhost:8080/auth/jira/callback",
	}

	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("expected no error creating client, got %v", err)
	}

	// Add comment to ticket
	err = client.AddComment("PROJ-123", "Shipment has been picked up from client", "mock-access-token")
	if err != nil {
		t.Fatalf("expected no error adding comment, got %v", err)
	}
}

// 🟥 RED: Test for building ticket request from shipment
// This test verifies that we can build a proper JIRA ticket request from shipment
func TestBuildTicketFromShipment(t *testing.T) {
	shipment := &models.Shipment{
		ID:              1,
		ClientCompanyID: 100,
		Status:          models.ShipmentStatusPendingPickup,
		Notes:           "Deploy laptop to John Doe",
		CreatedAt:       time.Now(),
	}

	clientCompany := &models.ClientCompany{
		Name: "Acme Corp",
	}

	laptops := []models.Laptop{
		{
			SerialNumber: "SN123456789",
			Brand:        "Dell",
			Model:        "Latitude 7420",
		},
	}

	// Build ticket request
	request := BuildTicketFromShipment(shipment, clientCompany, laptops, "PROJ")
	
	// Verify request
	if request.ProjectKey != "PROJ" {
		t.Errorf("expected project key PROJ, got %s", request.ProjectKey)
	}
	if request.Summary == "" {
		t.Error("expected summary to be populated")
	}
	if request.IssueType != "Task" {
		t.Errorf("expected issue type Task, got %s", request.IssueType)
	}
	if request.Description == "" {
		t.Error("expected description to be populated")
	}
}

// 🟥 RED: Test for syncing shipment status to JIRA
// This test verifies that we can sync shipment status changes to JIRA
func TestSyncShipmentStatusToJira(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/rest/api/3/issue/PROJ-123/transitions" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if r.URL.Path == "/rest/api/3/issue/PROJ-123/comment" {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"id": "10001"}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	config := Config{
		URL:          server.URL,
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURL:  "http://localhost:8080/auth/jira/callback",
	}

	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("expected no error creating client, got %v", err)
	}

	shipment := &models.Shipment{
		ID:     1,
		Status: models.ShipmentStatusPickedUpFromClient,
	}

	// Sync shipment status to JIRA
	err = client.SyncShipmentStatusToJira("PROJ-123", shipment, "mock-access-token")
	if err != nil {
		t.Fatalf("expected no error syncing status, got %v", err)
	}
}

