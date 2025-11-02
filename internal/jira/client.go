package jira

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
)

// Config holds the JIRA client configuration
type Config struct {
	URL      string
	Username string // JIRA account email (e.g., user@example.com)
	APIToken string // API token from JIRA user settings
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
	if config.Username == "" {
		return nil, errors.New("JIRA Username is required")
	}
	if config.APIToken == "" {
		return nil, errors.New("JIRA APIToken is required")
	}

	return &Client{
		config:     config,
		httpClient: &http.Client{},
	}, nil
}

// createAuthHeader generates Basic Auth header from username and API token
func (c *Client) createAuthHeader() string {
	auth := c.config.Username + ":" + c.config.APIToken
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

// TestConnection validates the connection to JIRA API using Basic Auth
func (c *Client) TestConnection() error {
	// Make a request to JIRA API to validate the connection
	url := fmt.Sprintf("%s/rest/api/3/myself", c.config.URL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Add authorization header with Basic Auth
	req.Header.Set("Authorization", c.createAuthHeader())
	req.Header.Set("Accept", "application/json")

	// Make the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to JIRA: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode == http.StatusUnauthorized {
		return errors.New("unauthorized: invalid credentials")
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

