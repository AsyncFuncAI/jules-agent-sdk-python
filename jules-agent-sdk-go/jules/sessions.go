package jules

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// SessionsAPI provides methods for interacting with the Sessions API
type SessionsAPI struct {
	client *BaseClient
}

// NewSessionsAPI creates a new SessionsAPI instance
func NewSessionsAPI(client *BaseClient) *SessionsAPI {
	return &SessionsAPI{client: client}
}

// CreateSessionRequest represents a request to create a new session
type CreateSessionRequest struct {
	Prompt              string         `json:"prompt"`
	Source              string         `json:"source"`
	StartingBranch      string         `json:"startingBranch,omitempty"`
	Title               string         `json:"title,omitempty"`
	RequirePlanApproval bool           `json:"requirePlanApproval,omitempty"`
}

// Create creates a new session
func (s *SessionsAPI) Create(ctx context.Context, req *CreateSessionRequest) (*Session, error) {
	// Build request body
	body := map[string]interface{}{
		"prompt": req.Prompt,
		"sourceContext": map[string]interface{}{
			"source": req.Source,
		},
	}

	if req.StartingBranch != "" {
		sourceContext := body["sourceContext"].(map[string]interface{})
		sourceContext["githubRepoContext"] = map[string]interface{}{
			"startingBranch": req.StartingBranch,
		}
	}

	if req.Title != "" {
		body["title"] = req.Title
	}

	if req.RequirePlanApproval {
		body["requirePlanApproval"] = true
	}

	// Make request
	resp, err := s.client.Post(ctx, "/sessions", body)
	if err != nil {
		return nil, err
	}

	// Parse response
	return s.parseSession(resp)
}

// Get retrieves a session by ID
func (s *SessionsAPI) Get(ctx context.Context, sessionID string) (*Session, error) {
	// Handle both short IDs and full names
	path := s.buildSessionPath(sessionID)

	resp, err := s.client.Get(ctx, path)
	if err != nil {
		return nil, err
	}

	return s.parseSession(resp)
}

// ListOptions represents options for listing sessions
type ListOptions struct {
	PageSize  int
	PageToken string
}

// List retrieves a list of sessions
func (s *SessionsAPI) List(ctx context.Context, opts *ListOptions) (*SessionsListResponse, error) {
	path := "/sessions"

	if opts != nil {
		query := ""
		if opts.PageSize > 0 {
			query += fmt.Sprintf("pageSize=%d", opts.PageSize)
		}
		if opts.PageToken != "" {
			if query != "" {
				query += "&"
			}
			query += fmt.Sprintf("pageToken=%s", opts.PageToken)
		}
		if query != "" {
			path += "?" + query
		}
	}

	resp, err := s.client.Get(ctx, path)
	if err != nil {
		return nil, err
	}

	// Parse response
	var result SessionsListResponse
	respBytes, _ := json.Marshal(resp)
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to parse sessions list: %w", err)
	}

	return &result, nil
}

// ApprovePlan approves a session plan
func (s *SessionsAPI) ApprovePlan(ctx context.Context, sessionID string) error {
	path := s.buildSessionPath(sessionID) + ":approvePlan"

	_, err := s.client.Post(ctx, path, nil)
	return err
}

// SendMessage sends a message to a session
func (s *SessionsAPI) SendMessage(ctx context.Context, sessionID string, prompt string) error {
	path := s.buildSessionPath(sessionID) + ":sendMessage"

	body := map[string]interface{}{
		"prompt": prompt,
	}

	_, err := s.client.Post(ctx, path, body)
	return err
}

// WaitForCompletionOptions represents options for waiting for session completion
type WaitForCompletionOptions struct {
	PollInterval time.Duration
	Timeout      time.Duration
}

// WaitForCompletion polls a session until it reaches a terminal state
func (s *SessionsAPI) WaitForCompletion(ctx context.Context, sessionID string, opts *WaitForCompletionOptions) (*Session, error) {
	pollInterval := DefaultPollInterval
	timeout := DefaultSessionTimeout

	if opts != nil {
		if opts.PollInterval > 0 {
			pollInterval = opts.PollInterval
		}
		if opts.Timeout > 0 {
			timeout = opts.Timeout
		}
	}

	// Create timeout context
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		// Get session status
		session, err := s.Get(ctx, sessionID)
		if err != nil {
			return nil, err
		}

		// Check if terminal state reached
		if session.State.IsTerminal() {
			if session.State == SessionStateFailed {
				return session, fmt.Errorf("session failed")
			}
			return session, nil
		}

		// Wait for next poll or timeout
		select {
		case <-ticker.C:
			continue
		case <-ctx.Done():
			return nil, NewTimeoutError(fmt.Sprintf("session did not complete within %v", timeout))
		}
	}
}

// buildSessionPath builds the API path for a session
func (s *SessionsAPI) buildSessionPath(sessionID string) string {
	if strings.HasPrefix(sessionID, "sessions/") {
		return "/" + sessionID
	}
	return "/sessions/" + sessionID
}

// parseSession parses a session from a response
func (s *SessionsAPI) parseSession(data map[string]interface{}) (*Session, error) {
	// Convert to JSON and back to properly parse the session
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal session data: %w", err)
	}

	var session Session
	if err := json.Unmarshal(jsonBytes, &session); err != nil {
		return nil, fmt.Errorf("failed to parse session: %w", err)
	}

	return &session, nil
}
