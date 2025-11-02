package jira

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

// User represents a JIRA user
type User struct {
	AccountID    string `json:"accountId"`
	EmailAddress string `json:"emailAddress"`
	DisplayName  string `json:"displayName"`
	Active       bool   `json:"active"`
}

// Project represents a JIRA project
type Project struct {
	ID          string `json:"id"`
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ProjectType string `json:"projectTypeKey"`
}

// GetCurrentUser retrieves information about the authenticated user
func (c *Client) GetCurrentUser() (*User, error) {
	url := fmt.Sprintf("%s/rest/api/3/myself", c.config.URL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", c.createAuthHeader())
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var user User
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &user, nil
}

// ListProjects retrieves all projects accessible to the authenticated user
func (c *Client) ListProjects() ([]Project, error) {
	url := fmt.Sprintf("%s/rest/api/3/project/search", c.config.URL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", c.createAuthHeader())
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get projects: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response struct {
		Values []Project `json:"values"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return response.Values, nil
}

// GetProjectDetails retrieves details about a specific project
func (c *Client) GetProjectDetails(key string) (*Project, error) {
	if key == "" {
		return nil, errors.New("project key is required")
	}

	url := fmt.Sprintf("%s/rest/api/3/project/%s", c.config.URL, key)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", c.createAuthHeader())
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("project not found")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var project Project
	if err := json.Unmarshal(body, &project); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &project, nil
}

