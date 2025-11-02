package jira

import (
	"errors"
	"fmt"
	"net/http"
)

// Config holds the JIRA client configuration
type Config struct {
	URL          string
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

// Client represents a JIRA API client
type Client struct {
	config     Config
	httpClient *http.Client
}

// NewClient creates a new JIRA client with the provided configuration
func NewClient(config Config) (*Client, error) {
	// Validate required configuration fields
	if config.URL == "" {
		return nil, errors.New("JIRA URL is required")
	}
	if config.ClientID == "" {
		return nil, errors.New("JIRA ClientID is required")
	}
	if config.ClientSecret == "" {
		return nil, errors.New("JIRA ClientSecret is required")
	}

	return &Client{
		config:     config,
		httpClient: &http.Client{},
	}, nil
}

// TestConnection validates the connection to JIRA API using the provided access token
func (c *Client) TestConnection(accessToken string) error {
	if accessToken == "" {
		return errors.New("access token is required")
	}

	// Make a request to JIRA API to validate the connection
	url := fmt.Sprintf("%s/rest/api/3/myself", c.config.URL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Add authorization header
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Accept", "application/json")

	// Make the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to JIRA: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode == http.StatusUnauthorized {
		return errors.New("unauthorized: invalid access token")
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

