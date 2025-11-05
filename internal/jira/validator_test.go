package jira

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClient_CreateTicketValidator(t *testing.T) {
	tests := []struct {
		name           string
		ticketKey      string
		serverResponse string
		statusCode     int
		wantErr        bool
		errContains    string
	}{
		{
			name:      "valid ticket exists",
			ticketKey: "SCOP-67702",
			serverResponse: `{
				"key": "SCOP-67702",
				"fields": {
					"summary": "Test ticket",
					"description": "Test description",
					"status": {"name": "Open"},
					"assignee": {"displayName": "John Doe", "emailAddress": "john@example.com"},
					"created": "2024-01-01T00:00:00Z",
					"updated": "2024-01-02T00:00:00Z"
				}
			}`,
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:           "ticket not found",
			ticketKey:      "SCOP-99999",
			serverResponse: `{"errorMessages":["Issue does not exist"],"errors":{}}`,
			statusCode:     http.StatusNotFound,
			wantErr:        true,
			errContains:    "JIRA ticket SCOP-99999 does not exist",
		},
		{
			name:           "API error",
			ticketKey:      "SCOP-67702",
			serverResponse: `{"errorMessages":["Internal server error"]}`,
			statusCode:     http.StatusInternalServerError,
			wantErr:        true,
			errContains:    "failed to validate JIRA ticket",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock HTTP server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify the request path
				expectedPath := "/rest/api/3/issue/" + tt.ticketKey
				if r.URL.Path != expectedPath {
					t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
				}

				// Return the mock response
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.serverResponse))
			}))
			defer server.Close()

			// Create a JIRA client with the mock server URL
			client := &Client{
				config: Config{
					URL:      server.URL,
					Username: "test@example.com",
					APIToken: "test-token",
				},
				httpClient: &http.Client{},
			}

			// Create the validator
			validator := client.CreateTicketValidator()

			// Test the validator
			err := validator(tt.ticketKey)

			// Check error expectation
			if (err != nil) != tt.wantErr {
				t.Errorf("validator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check error message contains expected text
			if err != nil && tt.errContains != "" {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("validator() error = %v, should contain %v", err.Error(), tt.errContains)
				}
			}
		})
	}
}

