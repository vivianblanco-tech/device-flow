# JIRA Integration Guide

## Overview

This guide explains how to use the JIRA integration to sync shipments with JIRA tickets.

## Configuration

### 1. Environment Variables

Add these to your `.env` file:

```bash
# JIRA Configuration (API Token Authentication)
JIRA_URL=https://bairesdev.atlassian.net
JIRA_USERNAME=your-email@bairesdev.com
JIRA_API_TOKEN=your-api-token-here
JIRA_DEFAULT_PROJECT=PROJ
```

### 2. API Token Setup

Generate a JIRA API token:

1. Go to https://id.atlassian.com/manage-profile/security/api-tokens
2. Click "Create API token"
3. Give it a descriptive name (e.g., "Laptop Tracking System")
4. Copy the generated token
5. Add it to your `.env` file as `JIRA_API_TOKEN`

**Important:** Keep your API token secure! Never commit it to version control.

## Usage Examples

### Initialize the JIRA Client

```go
package main

import (
    "os"
    "github.com/yourusername/laptop-tracking-system/internal/jira"
)

func initJiraClient() (*jira.Client, error) {
    config := jira.Config{
        URL:      os.Getenv("JIRA_URL"),
        Username: os.Getenv("JIRA_USERNAME"),
        APIToken: os.Getenv("JIRA_API_TOKEN"),
    }

    client, err := jira.NewClient(config)
    if err != nil {
        return nil, err
    }

    return client, nil
}
```

### Test the Connection

```go
func testJiraConnection(client *jira.Client) error {
    err := client.TestConnection()
    if err != nil {
        return fmt.Errorf("JIRA connection failed: %w", err)
    }
    fmt.Println("✅ JIRA connection successful!")
    return nil
}
```

## Sync Scenarios

### Scenario 1: Create JIRA Ticket When Shipment is Created

**Use Case:** Automatically create a JIRA ticket when a new shipment is created.

```go
func CreateShipmentWithJiraTicket(
    db *database.DB,
    jiraClient *jira.Client,
    shipment *models.Shipment,
    clientCompany *models.ClientCompany,
    laptops []models.Laptop,
) error {
    // 1. Create the shipment in the database
    err := db.CreateShipment(shipment)
    if err != nil {
        return fmt.Errorf("failed to create shipment: %w", err)
    }

    // 2. Build JIRA ticket from shipment data
    ticketRequest := jira.BuildTicketFromShipment(
        shipment,
        clientCompany,
        laptops,
        "PROJ", // Your JIRA project key
    )

    // 3. Create the JIRA ticket
    response, err := jiraClient.CreateTicket(ticketRequest)
    if err != nil {
        // Log error but don't fail the shipment creation
        log.Printf("Warning: Failed to create JIRA ticket: %v", err)
        return nil
    }

    // 4. Store the JIRA ticket key in the shipment
    // (You'll need to add a jira_ticket_key field to shipments table)
    err = db.UpdateShipmentJiraKey(shipment.ID, response.Key)
    if err != nil {
        log.Printf("Warning: Failed to save JIRA key: %v", err)
    }

    log.Printf("✅ Created JIRA ticket: %s for shipment #%d", response.Key, shipment.ID)
    return nil
}
```

### Scenario 2: Sync Shipment Status to JIRA

**Use Case:** Update JIRA ticket status when shipment status changes.

```go
func UpdateShipmentStatus(
    db *database.DB,
    jiraClient *jira.Client,
    shipmentID int64,
    newStatus models.ShipmentStatus,
    jiraTicketKey string,
) error {
    // 1. Update shipment status in database
    shipment, err := db.GetShipment(shipmentID)
    if err != nil {
        return err
    }

    shipment.UpdateStatus(newStatus)
    err = db.UpdateShipment(shipment)
    if err != nil {
        return err
    }

    // 2. Sync to JIRA if ticket exists
    if jiraTicketKey != "" && jiraClient != nil {
        err = jiraClient.SyncShipmentStatusToJira(jiraTicketKey, shipment)
        if err != nil {
            // Log but don't fail the status update
            log.Printf("Warning: Failed to sync to JIRA: %v", err)
        } else {
            log.Printf("✅ Synced status to JIRA ticket: %s", jiraTicketKey)
        }
    }

    return nil
}
```

### Scenario 3: Import Existing JIRA Ticket

**Use Case:** Create a shipment from an existing JIRA ticket.

```go
func ImportFromJiraTicket(
    db *database.DB,
    jiraClient *jira.Client,
    ticketKey string,
) (*models.Shipment, error) {
    // 1. Fetch the JIRA ticket
    ticket, err := jiraClient.GetTicket(ticketKey)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch ticket: %w", err)
    }

    // 2. Extract custom fields (if configured)
    // Note: You'll need the full ticket data for this
    // This is a simplified example
    customFields := &jira.CustomFields{
        SerialNumber:  "extracted-serial",
        EngineerEmail: "extracted-email",
        ClientCompany: "extracted-company",
    }

    // 3. Create shipment from ticket
    shipment, err := jira.CreateShipmentFromTicket(ticket, customFields)
    if err != nil {
        return nil, err
    }

    // 4. Save to database
    err = db.CreateShipment(shipment)
    if err != nil {
        return nil, err
    }

    log.Printf("✅ Imported shipment from JIRA ticket: %s", ticketKey)
    return shipment, nil
}
```

### Scenario 4: Search JIRA Tickets

**Use Case:** Find JIRA tickets related to hardware deployments.

```go
func SearchHardwareTickets(
    jiraClient *jira.Client,
    projectKey string,
) (*jira.SearchResults, error) {
    // Build JQL query
    jql := fmt.Sprintf(
        "project = %s AND labels = hardware-deployment AND status != Done",
        projectKey,
    )

    // Search tickets
    results, err := jiraClient.SearchTickets(jql)
    if err != nil {
        return nil, err
    }

    log.Printf("Found %d hardware deployment tickets", results.Total)
    return results, nil
}
```

### Scenario 5: Add Comment to JIRA Ticket

**Use Case:** Add notes when important events occur.

```go
func NotifyJiraOnDelivery(
    jiraClient *jira.Client,
    jiraTicketKey string,
    shipment *models.Shipment,
    engineer *models.SoftwareEngineer,
) error {
    comment := fmt.Sprintf(
        "✅ Laptop delivered successfully!\n\n"+
            "Engineer: %s\n"+
            "Email: %s\n"+
            "Delivery Date: %s\n"+
            "Shipment ID: #%d",
        engineer.Name,
        engineer.Email,
        shipment.DeliveredAt.Format("2006-01-02"),
        shipment.ID,
    )

    err := jiraClient.AddComment(jiraTicketKey, comment)
    if err != nil {
        return err
    }

    log.Printf("✅ Added delivery comment to JIRA ticket: %s", jiraTicketKey)
    return nil
}
```

## Integration with Handlers

### Add JIRA Client to Application Context

```go
// In cmd/web/main.go or your app initialization

type App struct {
    DB         *database.DB
    JiraClient *jira.Client
    // ... other fields
}

func main() {
    // Initialize JIRA client
    jiraClient, err := initJiraClient()
    if err != nil {
        log.Printf("Warning: JIRA integration disabled: %v", err)
        jiraClient = nil // Graceful degradation
    }

    app := &App{
        DB:         db,
        JiraClient: jiraClient,
    }

    // ... rest of your app setup
}
```

### Update Pickup Form Handler

```go
// In internal/handlers/pickup_form.go

func (h *Handler) HandlePickupFormSubmit(w http.ResponseWriter, r *http.Request) {
    // ... existing form validation ...

    // Create shipment
    err = h.DB.CreateShipment(shipment)
    if err != nil {
        // handle error
    }

    // Create JIRA ticket if client is available
    if h.JiraClient != nil {
        go func() {
            // Async JIRA ticket creation
            ticketRequest := jira.BuildTicketFromShipment(
                shipment,
                clientCompany,
                laptops,
                "PROJ",
            )
            _, err := h.JiraClient.CreateTicket(ticketRequest)
            if err != nil {
                log.Printf("Failed to create JIRA ticket: %v", err)
            }
        }()
    }

    // ... rest of handler ...
}
```

### Update Status Change Handler

```go
// In internal/handlers/shipments.go

func (h *Handler) UpdateShipmentStatus(w http.ResponseWriter, r *http.Request) {
    // ... parse request ...

    // Update shipment status
    shipment.UpdateStatus(newStatus)
    err = h.DB.UpdateShipment(shipment)
    if err != nil {
        // handle error
    }

    // Sync to JIRA
    if h.JiraClient != nil && shipment.JiraTicketKey != "" {
        go func() {
            // Async JIRA sync
            err := h.JiraClient.SyncShipmentStatusToJira(
                shipment.JiraTicketKey,
                shipment,
            )
            if err != nil {
                log.Printf("Failed to sync to JIRA: %v", err)
            }
        }()
    }

    // ... rest of handler ...
}
```

## API Token Security

Best practices for API token management:

```go
// Load token from environment at startup
func loadJiraConfig() jira.Config {
    return jira.Config{
        URL:      os.Getenv("JIRA_URL"),
        Username: os.Getenv("JIRA_USERNAME"),
        APIToken: os.Getenv("JIRA_API_TOKEN"),
    }
}

// Store config securely, not in plain text
// Use environment variables or secure secret management
// Never log or expose the API token

// Token validation
func validateJiraConfig(config jira.Config) error {
    if config.URL == "" || config.Username == "" || config.APIToken == "" {
        return errors.New("incomplete JIRA configuration")
    }
    return nil
}
```

## Database Schema Update

You'll want to add a field to track JIRA tickets:

```sql
-- Add to shipments table
ALTER TABLE shipments 
ADD COLUMN jira_ticket_key VARCHAR(50),
ADD COLUMN jira_ticket_id VARCHAR(50);

-- Add index for faster lookups
CREATE INDEX idx_shipments_jira_ticket_key ON shipments(jira_ticket_key);
```

## Error Handling Best Practices

```go
func syncToJiraSafely(
    jiraClient *jira.Client,
    ticketKey string,
    shipment *models.Shipment,
) {
    // Don't let JIRA failures break your main workflow
    defer func() {
        if r := recover(); r != nil {
            log.Printf("JIRA sync panic recovered: %v", r)
        }
    }()

    if jiraClient == nil {
        return // JIRA integration disabled
    }

    if ticketKey == "" {
        return // No ticket linked
    }

    err := jiraClient.SyncShipmentStatusToJira(ticketKey, shipment)
    if err != nil {
        log.Printf("JIRA sync failed: %v", err)
        // Could store failed syncs in a queue for retry
    }
}
```

## Testing with Mock Data

```go
func TestJiraIntegration(t *testing.T) {
    // The JIRA tests use httptest for mocking
    // You can use the same approach in your handlers

    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Mock JIRA API responses
        w.WriteHeader(http.StatusCreated)
        w.Write([]byte(`{"key":"PROJ-123"}`))
    }))
    defer server.Close()

    config := jira.Config{
        URL:          server.URL,
        ClientID:     "test",
        ClientSecret: "test",
    }

    client, _ := jira.NewClient(config)
    // Test your integration...
}
```

## JQL Query Examples

```go
// Common JQL queries for hardware tracking

// Find all open hardware deployments
"project = PROJ AND labels = hardware-deployment AND status != Done"

// Find tickets for specific client
"project = PROJ AND customfield_10003 = 'Acme Corp'"

// Find tickets with specific serial number
"project = PROJ AND customfield_10001 = 'SN123456789'"

// Find overdue deliveries
"project = PROJ AND duedate < now() AND status != Delivered"

// Find tickets assigned to specific user
"project = PROJ AND assignee = 'john.doe@example.com'"
```

## Monitoring and Logging

```go
// Add metrics for JIRA sync operations

type JiraMetrics struct {
    TicketsCreated  int64
    TicketsFailed   int64
    SyncsSucceeded  int64
    SyncsFailed     int64
}

func logJiraOperation(operation string, success bool, duration time.Duration) {
    log.Printf(
        "JIRA %s: success=%v duration=%v",
        operation,
        success,
        duration,
    )
}
```

## Summary

**Key Points:**
1. Initialize JIRA client at app startup with API token from environment
2. Use async operations (goroutines) for JIRA calls to avoid blocking
3. Implement graceful degradation if JIRA is unavailable
4. Store JIRA ticket keys in your database
5. Keep API token secure - never commit to version control
6. Don't let JIRA failures break your main workflow
7. Log all JIRA operations for debugging
8. API token authentication is simpler than OAuth for server-to-server

**Common Integration Points:**
- ✅ Pickup form submission → Create ticket
- ✅ Status changes → Update ticket
- ✅ Delivery completion → Add comment
- ✅ Dashboard → Search tickets
- ✅ Import form → Fetch ticket data

The integration is designed to be **optional** and **non-blocking** - your app should work fine even if JIRA is unavailable!

