package jira

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Ticket represents a JIRA issue/ticket
type Ticket struct {
	Key         string `json:"key"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Assignee    string `json:"assignee"`
	Created     string `json:"created"`
	Updated     string `json:"updated"`
}

// SearchResults represents the results of a JIRA search
type SearchResults struct {
	Total  int      `json:"total"`
	Issues []Ticket `json:"issues"`
}

// GetTicket fetches a specific JIRA ticket by key
func (c *Client) GetTicket(issueKey string) (*Ticket, error) {
	if issueKey == "" {
		return nil, errors.New("issue key is required")
	}

	// Build the request URL
	url := fmt.Sprintf("%s/rest/api/3/issue/%s", c.config.URL, issueKey)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authorization header with Basic Auth
	req.Header.Set("Authorization", c.createAuthHeader())
	req.Header.Set("Accept", "application/json")

	// Make the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch ticket: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("ticket not found")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse the response
	var response struct {
		Key    string `json:"key"`
		Fields struct {
			Summary     string `json:"summary"`
			Description string `json:"description"`
			Status      struct {
				Name string `json:"name"`
			} `json:"status"`
			Assignee struct {
				DisplayName  string `json:"displayName"`
				EmailAddress string `json:"emailAddress"`
			} `json:"assignee"`
			Created string `json:"created"`
			Updated string `json:"updated"`
		} `json:"fields"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Map to our Ticket struct
	ticket := &Ticket{
		Key:         response.Key,
		Summary:     response.Fields.Summary,
		Description: response.Fields.Description,
		Status:      response.Fields.Status.Name,
		Assignee:    response.Fields.Assignee.DisplayName,
		Created:     response.Fields.Created,
		Updated:     response.Fields.Updated,
	}

	return ticket, nil
}

// SearchTickets searches for JIRA tickets using JQL (JIRA Query Language)
func (c *Client) SearchTickets(jql string) (*SearchResults, error) {
	if jql == "" {
		return nil, errors.New("JQL query is required")
	}

	// Build the request URL with query parameters
	baseURL := fmt.Sprintf("%s/rest/api/3/search", c.config.URL)
	params := url.Values{}
	params.Add("jql", jql)
	params.Add("maxResults", "50") // Default max results

	requestURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authorization header with Basic Auth
	req.Header.Set("Authorization", c.createAuthHeader())
	req.Header.Set("Accept", "application/json")

	// Make the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to search tickets: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse the response
	var response struct {
		Total  int `json:"total"`
		Issues []struct {
			Key    string `json:"key"`
			Fields struct {
				Summary string `json:"summary"`
				Status  struct {
					Name string `json:"name"`
				} `json:"status"`
			} `json:"fields"`
		} `json:"issues"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Map to our SearchResults struct
	results := &SearchResults{
		Total:  response.Total,
		Issues: make([]Ticket, len(response.Issues)),
	}

	for i, issue := range response.Issues {
		results.Issues[i] = Ticket{
			Key:     issue.Key,
			Summary: issue.Fields.Summary,
			Status:  issue.Fields.Status.Name,
		}
	}

	return results, nil
}

