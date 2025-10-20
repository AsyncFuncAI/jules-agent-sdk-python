package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jules-labs/jules-go-sdk"
)

func main() {
	apiKey := os.Getenv("JULES_API_KEY")
	if apiKey == "" {
		log.Fatal("JULES_API_KEY environment variable not set")
	}

	client := jules.NewJulesClient(apiKey)

	// List sources
	sourcesResp, err := client.Sources.List()
	if err != nil {
		log.Fatalf("Error listing sources: %v", err)
	}

	if len(sourcesResp.Sources) == 0 {
		log.Fatal("No sources found")
	}

	fmt.Printf("Found %d sources\n", len(sourcesResp.Sources))

	// Create a session
	createReq := &jules.CreateSessionRequest{
		Prompt: "Add error handling to the authentication module",
		SourceContext: &jules.SourceContext{
			Source: sourcesResp.Sources[0].Name,
		},
	}

	session, err := client.Sessions.Create(createReq)
	if err != nil {
		log.Fatalf("Error creating session: %v", err)
	}

	fmt.Printf("Session created: %s\n", session.ID)
	fmt.Printf("View at: %s\n", session.URL)
}