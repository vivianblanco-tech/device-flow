package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/yourusername/laptop-tracking-system/internal/jira"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found, using system environment variables")
	}

	// Get JIRA configuration from environment
	jiraURL := os.Getenv("JIRA_URL")
	jiraUsername := os.Getenv("JIRA_USERNAME")
	jiraAPIToken := os.Getenv("JIRA_API_TOKEN")

	// Validate configuration
	if jiraURL == "" || jiraUsername == "" || jiraAPIToken == "" {
		fmt.Println("❌ JIRA configuration incomplete!")
		fmt.Println("\nRequired environment variables:")
		fmt.Println("  JIRA_URL:", jiraURL)
		fmt.Println("  JIRA_USERNAME:", jiraUsername)
		fmt.Println("  JIRA_API_TOKEN:", maskToken(jiraAPIToken))
		fmt.Println("\nPlease ensure all JIRA configuration is set in your .env file.")
		os.Exit(1)
	}

	// Create JIRA client
	client, err := jira.NewClient(jira.Config{
		URL:      jiraURL,
		Username: jiraUsername,
		APIToken: jiraAPIToken,
	})
	if err != nil {
		fmt.Printf("❌ Failed to create JIRA client: %v\n", err)
		os.Exit(1)
	}

	// Test connection
	fmt.Println("Testing JIRA connection...")
	fmt.Println("URL:", jiraURL)
	fmt.Println("Username:", jiraUsername)
	fmt.Println()

	err = client.TestConnection()
	if err != nil {
		fmt.Printf("❌ JIRA Connection Failed: %v\n", err)
		fmt.Println("\nPossible issues:")
		fmt.Println("  - Check that your API token is valid")
		fmt.Println("  - Verify your username (email) is correct")
		fmt.Println("  - Ensure the JIRA URL is correct")
		os.Exit(1)
	}

	fmt.Println("✅ JIRA Connection Successful!")
	fmt.Println()

	// Get current user information
	fmt.Println("Fetching user information...")
	user, err := client.GetCurrentUser()
	if err != nil {
		fmt.Printf("❌ Failed to get user information: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("User Information:")
	fmt.Printf("  Name: %s\n", user.DisplayName)
	fmt.Printf("  Email: %s\n", user.EmailAddress)
	fmt.Printf("  Account ID: %s\n", user.AccountID)
	fmt.Printf("  Active: %v\n", user.Active)
	fmt.Println()

	// List accessible projects
	fmt.Println("Fetching accessible projects...")
	projects, err := client.ListProjects()
	if err != nil {
		fmt.Printf("❌ Failed to list projects: %v\n", err)
		os.Exit(1)
	}

	if len(projects) == 0 {
		fmt.Println("⚠️  No projects found. You may not have access to any JIRA projects.")
		os.Exit(0)
	}

	fmt.Printf("Accessible Projects (%d total):\n", len(projects))
	fmt.Println()

	for i, project := range projects {
		fmt.Printf("  %d. %s - %s\n", i+1, project.Key, project.Name)
		if project.Description != "" {
			fmt.Printf("     Description: %s\n", truncate(project.Description, 80))
		}
		fmt.Printf("     Type: %s\n", project.ProjectType)
		fmt.Println()
	}

	// Provide recommendation
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("Recommendation:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
	fmt.Println("To use one of these projects for laptop tracking, update your .env file:")
	fmt.Println()
	fmt.Println("  JIRA_DEFAULT_PROJECT=<PROJECT_KEY>")
	fmt.Println()
	fmt.Println("For example, if you want to use the first project:")
	if len(projects) > 0 {
		fmt.Printf("  JIRA_DEFAULT_PROJECT=%s\n", projects[0].Key)
	}
	fmt.Println()
	fmt.Println("Choose the project key that corresponds to where you want to")
	fmt.Println("track hardware deployment tickets.")
}

// maskToken masks the API token for display purposes
func maskToken(token string) string {
	if token == "" {
		return "<not set>"
	}
	if len(token) <= 8 {
		return "****"
	}
	return token[:4] + "..." + token[len(token)-4:]
}

// truncate truncates a string to a maximum length
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

