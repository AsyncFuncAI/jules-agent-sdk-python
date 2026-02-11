package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/sashimikun/jules-agent-sdk-go/jules"
)

func main() {
	// Create client
	client, err := jules.NewClient(os.Getenv("JULES_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	ctx := context.Background()

	// Create a session
	session, err := client.Sessions.Create(ctx, &jules.CreateSessionRequest{
		Prompt: "Fix the bug in the login function",
		Source: "sources/my-repo",
		Title:  "Fix Login Bug",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Session created: %s\n", session.ID)
	fmt.Printf("State: %s\n", session.State)

	// Wait for completion
	result, err := client.Sessions.WaitForCompletion(ctx, session.ID, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Final state: %s\n", result.State)
	if result.Output != nil && result.Output.PullRequest != nil {
		fmt.Printf("PR URL: %s\n", result.Output.PullRequest.URL)
	}
}
