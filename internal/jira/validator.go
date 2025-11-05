package jira

import (
	"errors"
	"fmt"
)

// CreateTicketValidator creates a ticket validator function that uses the JIRA client
// This validator can be used with models.ValidateJiraTicketExists
func (c *Client) CreateTicketValidator() func(ticketKey string) error {
	return func(ticketKey string) error {
		ticket, err := c.GetTicket(ticketKey)
		if err != nil {
			// Check if it's a "not found" error
			if err.Error() == "ticket not found" {
				return fmt.Errorf("JIRA ticket %s does not exist", ticketKey)
			}
			// Other errors (network, auth, etc.)
			return fmt.Errorf("failed to validate JIRA ticket: %w", err)
		}
		
		// Additional validation: ensure ticket is not nil
		if ticket == nil {
			return errors.New("JIRA ticket " + ticketKey + " does not exist")
		}
		
		return nil
	}
}

