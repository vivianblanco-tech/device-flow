package jira

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// 🟥 RED: Test for JIRA client initialization
// This test verifies that we can create a new JIRA client with proper configuration
func TestNewClient(t *testing.T) {
	config := Config{
		URL:          "https://bairesdev.atlassian.net",
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURL:  "http://localhost:8080/auth/jira/callback",
	}

	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if client == nil {
		t.Fatal("expected client to be non-nil")
	}

	if client.config.URL != config.URL {
		t.Errorf("expected URL %s, got %s", config.URL, client.config.URL)
	}
}

// 🟥 RED: Test for JIRA client initialization with invalid config
// This test verifies that client creation fails with missing required fields
func TestNewClient_InvalidConfig(t *testing.T) {
	tests := []struct {
		name   string
		config Config
	}{
		{
			name: "missing URL",
			config: Config{
				ClientID:     "test-client-id",
				ClientSecret: "test-client-secret",
				RedirectURL:  "http://localhost:8080/auth/jira/callback",
			},
		},
		{
			name: "missing ClientID",
			config: Config{
				URL:          "https://bairesdev.atlassian.net",
				ClientSecret: "test-client-secret",
				RedirectURL:  "http://localhost:8080/auth/jira/callback",
			},
		},
		{
			name: "missing ClientSecret",
			config: Config{
				URL:         "https://bairesdev.atlassian.net",
				ClientID:    "test-client-id",
				RedirectURL: "http://localhost:8080/auth/jira/callback",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewClient(tt.config)
			if err == nil {
				t.Error("expected error for invalid config, got nil")
			}
		})
	}
}

// 🟥 RED: Test for JIRA connection validation
// This test verifies that we can test the connection to JIRA API
func TestClient_TestConnection(t *testing.T) {
	// Create a test server that simulates JIRA API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate a successful response from JIRA API
		if r.URL.Path == "/rest/api/3/myself" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"accountId":"test-account-id","emailAddress":"test@example.com"}`))
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

	// Test connection with a mock access token
	err = client.TestConnection("mock-access-token")
	if err != nil {
		t.Errorf("expected no error testing connection, got %v", err)
	}
}

// 🟥 RED: Test for JIRA connection validation with invalid token
// This test verifies that connection test fails with unauthorized access
func TestClient_TestConnection_Unauthorized(t *testing.T) {
	// Create a test server that returns 401 Unauthorized
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"errorMessages":["Unauthorized"]}`))
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

	// Test connection with an invalid token
	err = client.TestConnection("invalid-token")
	if err == nil {
		t.Error("expected error for unauthorized connection, got nil")
	}
}

// 🟥 RED: Test for JIRA connection validation without token
// This test verifies that connection test fails when no token is provided
func TestClient_TestConnection_NoToken(t *testing.T) {
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

	// Test connection without a token
	err = client.TestConnection("")
	if err == nil {
		t.Error("expected error for missing token, got nil")
	}
}
