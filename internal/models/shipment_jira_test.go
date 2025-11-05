package models

import (
	"errors"
	"testing"
)

// MockJiraClient is a mock implementation of the JIRA client for testing
type MockJiraClient struct {
	GetTicketFunc func(issueKey string) (*MockJiraTicket, error)
}

type MockJiraTicket struct {
	Key     string
	Summary string
	Status  string
}

func (m *MockJiraClient) GetTicket(issueKey string) (*MockJiraTicket, error) {
	if m.GetTicketFunc != nil {
		return m.GetTicketFunc(issueKey)
	}
	return nil, errors.New("not implemented")
}

func TestValidateJiraTicketExists(t *testing.T) {
	tests := []struct {
		name        string
		ticketKey   string
		mockClient  *MockJiraClient
		wantErr     bool
		errContains string
	}{
		{
			name:      "valid - ticket exists in JIRA",
			ticketKey: "SCOP-67702",
			mockClient: &MockJiraClient{
				GetTicketFunc: func(issueKey string) (*MockJiraTicket, error) {
					return &MockJiraTicket{
						Key:     "SCOP-67702",
						Summary: "Test ticket",
						Status:  "Open",
					}, nil
				},
			},
			wantErr: false,
		},
		{
			name:      "invalid - ticket does not exist in JIRA",
			ticketKey: "SCOP-99999",
			mockClient: &MockJiraClient{
				GetTicketFunc: func(issueKey string) (*MockJiraTicket, error) {
					return nil, errors.New("ticket not found")
				},
			},
			wantErr:     true,
			errContains: "JIRA ticket SCOP-99999 does not exist",
		},
		{
			name:      "invalid - JIRA API error",
			ticketKey: "SCOP-67702",
			mockClient: &MockJiraClient{
				GetTicketFunc: func(issueKey string) (*MockJiraTicket, error) {
					return nil, errors.New("API connection failed")
				},
			},
			wantErr:     true,
			errContains: "failed to validate JIRA ticket",
		},
		{
			name:       "valid - nil client (skip validation for sample data)",
			ticketKey:  "SCOP-67702",
			mockClient: nil,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a validator function that uses the mock client
			var validator JiraTicketValidator
			if tt.mockClient != nil {
				validator = func(ticketKey string) error {
					ticket, err := tt.mockClient.GetTicket(ticketKey)
					if err != nil {
						if err.Error() == "ticket not found" {
							return errors.New("JIRA ticket " + ticketKey + " does not exist")
						}
						return errors.New("failed to validate JIRA ticket: " + err.Error())
					}
					if ticket == nil {
						return errors.New("JIRA ticket " + ticketKey + " does not exist")
					}
					return nil
				}
			}

			err := ValidateJiraTicketExists(tt.ticketKey, validator)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJiraTicketExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.errContains != "" {
				if !contains(err.Error(), tt.errContains) {
					t.Errorf("ValidateJiraTicketExists() error = %v, should contain %v", err.Error(), tt.errContains)
				}
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > 0 && len(substr) > 0 && stringContains(s, substr)))
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

