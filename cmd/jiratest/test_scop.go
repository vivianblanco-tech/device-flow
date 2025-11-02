package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/yourusername/laptop-tracking-system/internal/jira"
)

func main() {
	// Load environment variables
	godotenv.Load()

	// Create JIRA client
	client, err := jira.NewClient(jira.Config{
		URL:      os.Getenv("JIRA_URL"),
		Username: os.Getenv("JIRA_USERNAME"),
		APIToken: os.Getenv("JIRA_API_TOKEN"),
	})
	if err != nil {
		fmt.Printf("❌ Failed to create client: %v\n", err)
		os.Exit(1)
	}

	// Test access to SCOP project
	fmt.Println("Testing access to SCOP project...")
	fmt.Println()

	project, err := client.GetProjectDetails("SCOP")
	if err != nil {
		fmt.Printf("❌ Cannot access SCOP project: %v\n", err)
		fmt.Println()
		fmt.Println("This could mean:")
		fmt.Println("  - The project doesn't exist")
		fmt.Println("  - You don't have permission to access it")
		fmt.Println("  - The project key is different (case-sensitive)")
		os.Exit(1)
	}

	fmt.Println("✅ You have access to SCOP!")
	fmt.Println()
	fmt.Println("Project Details:")
	fmt.Printf("  Key: %s\n", project.Key)
	fmt.Printf("  Name: %s\n", project.Name)
	fmt.Printf("  ID: %s\n", project.ID)
	fmt.Printf("  Type: %s\n", project.ProjectType)
	if project.Description != "" {
		fmt.Printf("  Description: %s\n", project.Description)
	}
	fmt.Println()
	fmt.Println("You can use this project for laptop tracking by setting:")
	fmt.Println("  JIRA_DEFAULT_PROJECT=SCOP")
}

