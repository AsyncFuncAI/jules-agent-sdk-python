package jules

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestSessionsAPICreate(t *testing.T) {
	now := time.Now()
	mockServer := NewMockServer(t, []MockResponse{
		{
			StatusCode: 200,
			Body: map[string]interface{}{
				"name":   "sessions/test-123",
				"id":     "test-123",
				"prompt": "Test prompt",
				"state":  "QUEUED",
				"sourceContext": map[string]interface{}{
					"source": "sources/test-repo",
				},
				"createTime": now.Format(time.RFC3339),
			},
		},
	})
	defer mockServer.Close()

	config := NewTestConfig(mockServer.URL, "test-key")
	baseClient, _ := NewBaseClient(config)
	defer baseClient.Close()

	api := NewSessionsAPI(baseClient)

	ctx := context.Background()
	session, err := api.Create(ctx, &CreateSessionRequest{
		Prompt: "Test prompt",
		Source: "sources/test-repo",
		Title:  "Test Session",
	})

	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if session.ID != "test-123" {
		t.Errorf("Session.ID = %v, want test-123", session.ID)
	}
	if session.State != SessionStateQueued {
		t.Errorf("Session.State = %v, want QUEUED", session.State)
	}

	// Verify request
	req := mockServer.GetLastRequest()
	if req.Method != http.MethodPost {
		t.Errorf("Request method = %v, want POST", req.Method)
	}
	if !strings.HasSuffix(req.URL.Path, "/sessions") {
		t.Errorf("Request path = %v, want /sessions", req.URL.Path)
	}
}

func TestSessionsAPICreateWithBranch(t *testing.T) {
	mockServer := NewMockServer(t, []MockResponse{
		{
			StatusCode: 200,
			Body: map[string]interface{}{
				"id":    "test-123",
				"state": "QUEUED",
				"sourceContext": map[string]interface{}{
					"source": "sources/test-repo",
					"githubRepoContext": map[string]interface{}{
						"startingBranch": "develop",
					},
				},
			},
		},
	})
	defer mockServer.Close()

	config := NewTestConfig(mockServer.URL, "test-key")
	baseClient, _ := NewBaseClient(config)
	defer baseClient.Close()

	api := NewSessionsAPI(baseClient)

	ctx := context.Background()
	session, err := api.Create(ctx, &CreateSessionRequest{
		Prompt:         "Test prompt",
		Source:         "sources/test-repo",
		StartingBranch: "develop",
	})

	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if session.SourceContext == nil || session.SourceContext.GitHubRepoContext == nil {
		t.Fatal("Expected GitHub repo context")
	}
	if session.SourceContext.GitHubRepoContext.StartingBranch != "develop" {
		t.Errorf("StartingBranch = %v, want develop", session.SourceContext.GitHubRepoContext.StartingBranch)
	}
}

func TestSessionsAPIGet(t *testing.T) {
	mockServer := NewMockServer(t, []MockResponse{
		{
			StatusCode: 200,
			Body: map[string]interface{}{
				"name":   "sessions/test-456",
				"id":     "test-456",
				"state":  "IN_PROGRESS",
				"prompt": "Get test",
			},
		},
	})
	defer mockServer.Close()

	config := NewTestConfig(mockServer.URL, "test-key")
	baseClient, _ := NewBaseClient(config)
	defer baseClient.Close()

	api := NewSessionsAPI(baseClient)

	ctx := context.Background()
	session, err := api.Get(ctx, "test-456")

	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if session.ID != "test-456" {
		t.Errorf("Session.ID = %v, want test-456", session.ID)
	}
	if session.State != SessionStateInProgress {
		t.Errorf("Session.State = %v, want IN_PROGRESS", session.State)
	}

	// Verify request path
	req := mockServer.GetLastRequest()
	if !strings.Contains(req.URL.Path, "test-456") {
		t.Errorf("Request path = %v, should contain test-456", req.URL.Path)
	}
}

func TestSessionsAPIGetWithFullName(t *testing.T) {
	mockServer := NewMockServer(t, []MockResponse{
		{
			StatusCode: 200,
			Body: map[string]interface{}{
				"name": "sessions/test-789",
				"id":   "test-789",
			},
		},
	})
	defer mockServer.Close()

	config := NewTestConfig(mockServer.URL, "test-key")
	baseClient, _ := NewBaseClient(config)
	defer baseClient.Close()

	api := NewSessionsAPI(baseClient)

	ctx := context.Background()
	session, err := api.Get(ctx, "sessions/test-789")

	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if session.ID != "test-789" {
		t.Errorf("Session.ID = %v, want test-789", session.ID)
	}
}

func TestSessionsAPIList(t *testing.T) {
	mockServer := NewMockServer(t, []MockResponse{
		{
			StatusCode: 200,
			Body: map[string]interface{}{
				"sessions": []interface{}{
					map[string]interface{}{
						"id":    "session-1",
						"state": "COMPLETED",
					},
					map[string]interface{}{
						"id":    "session-2",
						"state": "IN_PROGRESS",
					},
				},
				"nextPageToken": "token-123",
			},
		},
	})
	defer mockServer.Close()

	config := NewTestConfig(mockServer.URL, "test-key")
	baseClient, _ := NewBaseClient(config)
	defer baseClient.Close()

	api := NewSessionsAPI(baseClient)

	ctx := context.Background()
	response, err := api.List(ctx, &ListOptions{
		PageSize:  10,
		PageToken: "",
	})

	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(response.Sessions) != 2 {
		t.Errorf("Sessions count = %v, want 2", len(response.Sessions))
	}
	if response.NextPageToken != "token-123" {
		t.Errorf("NextPageToken = %v, want token-123", response.NextPageToken)
	}

	// Verify request query params
	req := mockServer.GetLastRequest()
	query := req.URL.Query()
	if query.Get("pageSize") != "10" {
		t.Errorf("pageSize query param = %v, want 10", query.Get("pageSize"))
	}
}

func TestSessionsAPIApprovePlan(t *testing.T) {
	mockServer := NewMockServer(t, []MockResponse{
		{StatusCode: 200, Body: map[string]interface{}{}},
	})
	defer mockServer.Close()

	config := NewTestConfig(mockServer.URL, "test-key")
	baseClient, _ := NewBaseClient(config)
	defer baseClient.Close()

	api := NewSessionsAPI(baseClient)

	ctx := context.Background()
	err := api.ApprovePlan(ctx, "session-123")

	if err != nil {
		t.Fatalf("ApprovePlan() error = %v", err)
	}

	// Verify request
	req := mockServer.GetLastRequest()
	if req.Method != http.MethodPost {
		t.Errorf("Request method = %v, want POST", req.Method)
	}
	if !strings.Contains(req.URL.Path, "approvePlan") {
		t.Errorf("Request path = %v, should contain approvePlan", req.URL.Path)
	}
}

func TestSessionsAPISendMessage(t *testing.T) {
	mockServer := NewMockServer(t, []MockResponse{
		{StatusCode: 200, Body: map[string]interface{}{}},
	})
	defer mockServer.Close()

	config := NewTestConfig(mockServer.URL, "test-key")
	baseClient, _ := NewBaseClient(config)
	defer baseClient.Close()

	api := NewSessionsAPI(baseClient)

	ctx := context.Background()
	err := api.SendMessage(ctx, "session-123", "Additional instruction")

	if err != nil {
		t.Fatalf("SendMessage() error = %v", err)
	}

	// Verify request
	req := mockServer.GetLastRequest()
	if req.Method != http.MethodPost {
		t.Errorf("Request method = %v, want POST", req.Method)
	}
	if !strings.Contains(req.URL.Path, "sendMessage") {
		t.Errorf("Request path = %v, should contain sendMessage", req.URL.Path)
	}
}

func TestSessionsAPIWaitForCompletion(t *testing.T) {
	// First call returns IN_PROGRESS, second returns COMPLETED
	mockServer := NewMockServer(t, []MockResponse{
		{
			StatusCode: 200,
			Body: map[string]interface{}{
				"id":    "session-123",
				"state": "IN_PROGRESS",
			},
		},
		{
			StatusCode: 200,
			Body: map[string]interface{}{
				"id":    "session-123",
				"state": "COMPLETED",
			},
		},
	})
	defer mockServer.Close()

	config := NewTestConfig(mockServer.URL, "test-key")
	baseClient, _ := NewBaseClient(config)
	defer baseClient.Close()

	api := NewSessionsAPI(baseClient)

	ctx := context.Background()
	session, err := api.WaitForCompletion(ctx, "session-123", &WaitForCompletionOptions{
		PollInterval: 100 * time.Millisecond,
		Timeout:      5 * time.Second,
	})

	if err != nil {
		t.Fatalf("WaitForCompletion() error = %v", err)
	}

	if session.State != SessionStateCompleted {
		t.Errorf("Session.State = %v, want COMPLETED", session.State)
	}

	// Verify it polled multiple times
	if mockServer.GetRequestCount() < 2 {
		t.Errorf("Request count = %v, want at least 2", mockServer.GetRequestCount())
	}
}

func TestSessionsAPIWaitForCompletionFailed(t *testing.T) {
	mockServer := NewMockServer(t, []MockResponse{
		{
			StatusCode: 200,
			Body: map[string]interface{}{
				"id":    "session-123",
				"state": "FAILED",
			},
		},
	})
	defer mockServer.Close()

	config := NewTestConfig(mockServer.URL, "test-key")
	baseClient, _ := NewBaseClient(config)
	defer baseClient.Close()

	api := NewSessionsAPI(baseClient)

	ctx := context.Background()
	session, err := api.WaitForCompletion(ctx, "session-123", nil)

	if err == nil {
		t.Fatal("WaitForCompletion() expected error for failed session, got nil")
	}

	if session.State != SessionStateFailed {
		t.Errorf("Session.State = %v, want FAILED", session.State)
	}
}

func TestSessionsAPIWaitForCompletionTimeout(t *testing.T) {
	// Always return IN_PROGRESS
	mockServer := NewMockServer(t, []MockResponse{
		{StatusCode: 200, Body: map[string]interface{}{"id": "s1", "state": "IN_PROGRESS"}},
		{StatusCode: 200, Body: map[string]interface{}{"id": "s1", "state": "IN_PROGRESS"}},
		{StatusCode: 200, Body: map[string]interface{}{"id": "s1", "state": "IN_PROGRESS"}},
		{StatusCode: 200, Body: map[string]interface{}{"id": "s1", "state": "IN_PROGRESS"}},
		{StatusCode: 200, Body: map[string]interface{}{"id": "s1", "state": "IN_PROGRESS"}},
	})
	defer mockServer.Close()

	config := NewTestConfig(mockServer.URL, "test-key")
	baseClient, _ := NewBaseClient(config)
	defer baseClient.Close()

	api := NewSessionsAPI(baseClient)

	ctx := context.Background()
	_, err := api.WaitForCompletion(ctx, "s1", &WaitForCompletionOptions{
		PollInterval: 50 * time.Millisecond,
		Timeout:      200 * time.Millisecond,
	})

	if err == nil {
		t.Fatal("WaitForCompletion() expected timeout error, got nil")
	}

	if _, ok := err.(*TimeoutError); !ok {
		t.Errorf("Expected TimeoutError, got %T", err)
	}
}

func TestBuildSessionPath(t *testing.T) {
	api := &SessionsAPI{}

	tests := []struct {
		input string
		want  string
	}{
		{"test-123", "/sessions/test-123"},
		{"sessions/test-456", "/sessions/test-456"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := api.buildSessionPath(tt.input)
			if got != tt.want {
				t.Errorf("buildSessionPath(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
