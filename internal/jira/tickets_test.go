package jira

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// 🟥 RED: Test for fetching JIRA ticket by key
// This test verifies that we can fetch a specific ticket from JIRA
func TestClient_GetTicket(t *testing.T) {
	// Create a test server that simulates JIRA API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/rest/api/3/issue/PROJ-123" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"key": "PROJ-123",
				"fields": {
					"summary": "Test Issue",
					"description": "Test Description",
					"status": {
						"name": "In Progress"
					},
					"assignee": {
						"displayName": "John Doe",
						"emailAddress": "john@example.com"
					},
					"created": "2023-10-01T10:00:00.000+0000",
					"updated": "2023-10-02T15:30:00.000+0000"
				}
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

	// Fetch the ticket
	ticket, err := client.GetTicket("PROJ-123", "mock-access-token")
	if err != nil {
		t.Fatalf("expected no error fetching ticket, got %v", err)
	}

	// Verify ticket data
	if ticket.Key != "PROJ-123" {
		t.Errorf("expected key PROJ-123, got %s", ticket.Key)
	}
	if ticket.Summary != "Test Issue" {
		t.Errorf("expected summary 'Test Issue', got %s", ticket.Summary)
	}
	if ticket.Status != "In Progress" {
		t.Errorf("expected status 'In Progress', got %s", ticket.Status)
	}
}

// 🟥 RED: Test for fetching non-existent ticket
// This test verifies proper error handling when ticket doesn't exist
func TestClient_GetTicket_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"errorMessages":["Issue does not exist or you do not have permission to see it."]}`))
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

	_, err = client.GetTicket("NONEXISTENT-123", "mock-access-token")
	if err == nil {
		t.Error("expected error for non-existent ticket, got nil")
	}
}

// 🟥 RED: Test for fetching ticket without access token
// This test verifies that fetching requires an access token
func TestClient_GetTicket_NoToken(t *testing.T) {
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

	_, err = client.GetTicket("PROJ-123", "")
	if err == nil {
		t.Error("expected error for missing token, got nil")
	}
}

// 🟥 RED: Test for searching JIRA tickets
// This test verifies that we can search for tickets using JQL
func TestClient_SearchTickets(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/rest/api/3/search" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"total": 2,
				"issues": [
					{
						"key": "PROJ-123",
						"fields": {
							"summary": "First Issue",
							"status": {"name": "To Do"}
						}
					},
					{
						"key": "PROJ-124",
						"fields": {
							"summary": "Second Issue",
							"status": {"name": "In Progress"}
						}
					}
				]
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

	// Search tickets
	results, err := client.SearchTickets("project = PROJ", "mock-access-token")
	if err != nil {
		t.Fatalf("expected no error searching tickets, got %v", err)
	}

	// Verify results
	if results.Total != 2 {
		t.Errorf("expected 2 total results, got %d", results.Total)
	}
	if len(results.Issues) != 2 {
		t.Errorf("expected 2 issues, got %d", len(results.Issues))
	}
	if results.Issues[0].Key != "PROJ-123" {
		t.Errorf("expected first issue key PROJ-123, got %s", results.Issues[0].Key)
	}
}

