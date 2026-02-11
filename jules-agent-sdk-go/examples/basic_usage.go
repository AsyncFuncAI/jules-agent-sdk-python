package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/sashimikun/jules-agent-sdk-go/jules"
)

func main() {
	// Get API key from environment variable
	apiKey := os.Getenv("JULES_API_KEY")
	if apiKey == "" {
		log.Fatal("JULES_API_KEY environment variable is required")
	}

	// Create a new Jules client
	client, err := jules.NewClient(apiKey)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// Example 1: List all sources
	fmt.Println("=== Listing Sources ===")
	sources, err := client.Sources.ListAll(ctx, "")
	if err != nil {
		log.Fatalf("Failed to list sources: %v", err)
	}

	for _, source := range sources {
		fmt.Printf("Source: %s (ID: %s)\n", source.Name, source.ID)
		if source.GitHubRepo != nil {
			fmt.Printf("  GitHub: %s/%s\n", source.GitHubRepo.Owner, source.GitHubRepo.Repo)
		}
	}

	// Example 2: Create a new session
	fmt.Println("\n=== Creating Session ===")
	session, err := client.Sessions.Create(ctx, &jules.CreateSessionRequest{
		Prompt: "Add a new feature to improve error handling",
		Source: sources[0].Name, // Use the first source
		Title:  "Improve Error Handling",
	})
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}

	fmt.Printf("Session created: %s\n", session.ID)
	fmt.Printf("State: %s\n", session.State)
	fmt.Printf("URL: %s\n", session.URL)

	// Example 3: Wait for session completion
	fmt.Println("\n=== Waiting for Completion ===")
	completedSession, err := client.Sessions.WaitForCompletion(ctx, session.ID, nil)
	if err != nil {
		log.Fatalf("Failed to wait for completion: %v", err)
	}

	fmt.Printf("Session completed: %s\n", completedSession.State)
	if completedSession.Output != nil && completedSession.Output.PullRequest != nil {
		fmt.Printf("Pull Request: %s\n", completedSession.Output.PullRequest.URL)
	}

	// Example 4: List activities for the session
	fmt.Println("\n=== Listing Activities ===")
	activities, err := client.Activities.ListAll(ctx, session.ID)
	if err != nil {
		log.Fatalf("Failed to list activities: %v", err)
	}

	for _, activity := range activities {
		fmt.Printf("Activity: %s\n", activity.Description)
		fmt.Printf("  Originator: %s\n", activity.Originator)
		if activity.CreateTime != nil {
			fmt.Printf("  Created: %s\n", activity.CreateTime)
		}
	}

	// Display client statistics
	fmt.Println("\n=== Client Statistics ===")
	stats := client.Stats()
	fmt.Printf("Total Requests: %d\n", stats["request_count"])
	fmt.Printf("Total Errors: %d\n", stats["error_count"])
}
