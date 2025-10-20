# Jules Agent SDK for Go

The official Go SDK for the Jules Agent API. Jules is an AI-powered agent that helps automate software development tasks.

## Features

- **Sessions API**: Create and manage AI agent sessions
- **Activities API**: Track and retrieve session activities
- **Sources API**: Manage source repositories
- **Automatic Retries**: Built-in exponential backoff for failed requests
- **Connection Pooling**: Efficient HTTP connection management
- **Comprehensive Error Handling**: Specific error types for different failure scenarios
- **Context Support**: Full support for Go's context package for cancellation and timeouts

## Installation

```bash
go get github.com/sashimikun/jules-agent-sdk-go
```

## Quick Start

```go
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

    // Wait for completion
    result, err := client.Sessions.WaitForCompletion(ctx, session.ID, nil)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Final state: %s\n", result.State)
}
```

## API Reference

### Client Initialization

#### Basic Client

```go
client, err := jules.NewClient(apiKey)
if err != nil {
    log.Fatal(err)
}
defer client.Close()
```

#### Custom Configuration

```go
config := &jules.Config{
    APIKey:             apiKey,
    BaseURL:            "https://julius.googleapis.com/v1alpha",
    Timeout:            30 * time.Second,
    MaxRetries:         3,
    RetryBackoffFactor: 1.0,
    MaxBackoff:         10 * time.Second,
    VerifySSL:          true,
}

client, err := jules.NewClientWithConfig(config)
if err != nil {
    log.Fatal(err)
}
defer client.Close()
```

### Sessions API

#### Create a Session

```go
session, err := client.Sessions.Create(ctx, &jules.CreateSessionRequest{
    Prompt:              "Add authentication to the API",
    Source:              "sources/my-repo",
    StartingBranch:      "main",
    Title:               "Add Authentication",
    RequirePlanApproval: false,
})
```

#### Get a Session

```go
session, err := client.Sessions.Get(ctx, "session-id")
```

#### List Sessions

```go
response, err := client.Sessions.List(ctx, &jules.ListOptions{
    PageSize:  10,
    PageToken: "",
})
```

#### Approve a Plan

```go
err := client.Sessions.ApprovePlan(ctx, "session-id")
```

#### Send a Message

```go
err := client.Sessions.SendMessage(ctx, "session-id", "Please also add rate limiting")
```

#### Wait for Completion

```go
session, err := client.Sessions.WaitForCompletion(ctx, "session-id", &jules.WaitForCompletionOptions{
    PollInterval: 5 * time.Second,
    Timeout:      600 * time.Second,
})
```

### Activities API

#### Get an Activity

```go
activity, err := client.Activities.Get(ctx, "session-id", "activity-id")
```

#### List Activities

```go
response, err := client.Activities.List(ctx, "session-id", &jules.ListOptions{
    PageSize:  10,
    PageToken: "",
})
```

#### List All Activities (with automatic pagination)

```go
activities, err := client.Activities.ListAll(ctx, "session-id")
```

### Sources API

#### Get a Source

```go
source, err := client.Sources.Get(ctx, "source-id")
```

#### List Sources

```go
response, err := client.Sources.List(ctx, &jules.SourcesListOptions{
    Filter:    "owner:myorg",
    PageSize:  10,
    PageToken: "",
})
```

#### List All Sources (with automatic pagination)

```go
sources, err := client.Sources.ListAll(ctx, "owner:myorg")
```

## Data Models

### Session States

The SDK defines the following session states:

- `SessionStateUnspecified`: Default unspecified state
- `SessionStateQueued`: Session is queued
- `SessionStatePlanning`: Session is in planning phase
- `SessionStateAwaitingPlanApproval`: Session is waiting for plan approval
- `SessionStateAwaitingUserFeedback`: Session is waiting for user feedback
- `SessionStateInProgress`: Session is in progress
- `SessionStatePaused`: Session is paused
- `SessionStateFailed`: Session has failed
- `SessionStateCompleted`: Session has completed

### Error Types

The SDK provides specific error types for different failure scenarios:

- `APIError`: Base error type for all API errors
- `AuthenticationError`: 401 authentication errors
- `NotFoundError`: 404 not found errors
- `ValidationError`: 400 validation errors
- `RateLimitError`: 429 rate limit errors (includes RetryAfter value)
- `ServerError`: 5xx server errors
- `TimeoutError`: Timeout errors

### Error Handling Example

```go
session, err := client.Sessions.Get(ctx, "session-id")
if err != nil {
    switch e := err.(type) {
    case *jules.AuthenticationError:
        log.Printf("Authentication failed: %s", e.Message)
    case *jules.NotFoundError:
        log.Printf("Session not found: %s", e.Message)
    case *jules.RateLimitError:
        log.Printf("Rate limited. Retry after %d seconds", e.RetryAfter)
    default:
        log.Printf("Error: %v", err)
    }
    return
}
```

## Configuration Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| APIKey | string | (required) | Your Jules API key |
| BaseURL | string | `https://julius.googleapis.com/v1alpha` | API base URL |
| Timeout | time.Duration | 30s | HTTP request timeout |
| MaxRetries | int | 3 | Maximum number of retry attempts |
| RetryBackoffFactor | float64 | 1.0 | Exponential backoff factor |
| MaxBackoff | time.Duration | 10s | Maximum backoff duration |
| VerifySSL | bool | true | Enable SSL certificate verification |

## Advanced Features

### Connection Pooling

The SDK automatically manages HTTP connection pooling with:
- 10 max idle connections
- 10 max idle connections per host
- 90 second idle connection timeout

### Automatic Retries

The SDK automatically retries:
- Network errors (connection failures, timeouts)
- 5xx server errors

The SDK does NOT retry:
- 4xx client errors (these indicate a problem with your request)
- 429 rate limit errors (returns immediately with retry-after information)

Retry behavior uses exponential backoff:
```
backoff = RetryBackoffFactor * 2^(attempt-1)
capped at MaxBackoff
```

### Statistics

You can get request statistics from the client:

```go
stats := client.Stats()
fmt.Printf("Total Requests: %d\n", stats["request_count"])
fmt.Printf("Total Errors: %d\n", stats["error_count"])
```

## Examples

See the [examples](./examples) directory for more usage examples:

- `simple_example.go`: Basic session creation and completion
- `basic_usage.go`: Comprehensive example covering all major features

## Development

### Running Examples

```bash
# Set your API key
export JULES_API_KEY="your-api-key-here"

# Run the simple example
go run examples/simple_example.go

# Run the comprehensive example
go run examples/basic_usage.go
```

### Building

```bash
cd jules-agent-sdk-go
go build ./jules
```

### Testing

```bash
go test ./jules -v
```

## License

MIT License - See LICENSE file for details

## Support

For issues and questions:
- GitHub Issues: [https://github.com/sashimikun/jules-agent-sdk-go/issues](https://github.com/sashimikun/jules-agent-sdk-go/issues)
- Documentation: [https://docs.julius.googleapis.com](https://docs.julius.googleapis.com)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
