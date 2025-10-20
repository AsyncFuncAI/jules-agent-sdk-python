package jules

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

// CreateSessionRequest represents a request to create a session.
type CreateSessionRequest struct {
	Prompt         string `json:"prompt"`
	Source         string `json:"source"`
	StartingBranch string `json:"startingBranch,omitempty"`
}

// SessionListResponse represents a response from listing sessions.
type SessionListResponse struct {
	Sessions      []Session `json:"sessions"`
	NextPageToken string    `json:"nextPageToken"`
}

// CreateSession creates a new session.
func (c *Client) CreateSession(ctx context.Context, req CreateSessionRequest) (*Session, error) {
	session := &Session{
		Prompt: req.Prompt,
		SourceContext: &SourceContext{
			Source: req.Source,
		},
	}

	if req.StartingBranch != "" {
		session.SourceContext.GitHubRepoContext = &GitHubRepoContext{
			StartingBranch: req.StartingBranch,
		}
	}

	var createdSession Session
	if err := c.doRequest(ctx, "POST", "/sessions", session, &createdSession); err != nil {
		return nil, err
	}

	return &createdSession, nil
}

// GetSession gets a session by ID.
func (c *Client) GetSession(ctx context.Context, sessionID string) (*Session, error) {
	var session Session
	if err := c.doRequest(ctx, "GET", fmt.Sprintf("/sessions/%s", sessionID), nil, &session); err != nil {
		return nil, err
	}

	return &session, nil
}

// ListSessions lists all sessions.
func (c *Client) ListSessions(ctx context.Context, pageSize int, pageToken string) (*SessionListResponse, error) {
	query := url.Values{}
	if pageSize > 0 {
		query.Set("pageSize", fmt.Sprintf("%d", pageSize))
	}
	if pageToken != "" {
		query.Set("pageToken", pageToken)
	}

	path := "/sessions"
	if len(query) > 0 {
		path += "?" + query.Encode()
	}

	var response SessionListResponse
	if err := c.doRequest(ctx, "GET", path, nil, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// DeleteSession deletes a session by ID.
func (c *Client) DeleteSession(ctx context.Context, sessionID string) error {
	return c.doRequest(ctx, "DELETE", fmt.Sprintf("/sessions/%s", sessionID), nil, nil)
}

// ContinueSession sends a user message to a session.
func (c *Client) ContinueSession(ctx context.Context, sessionID, message string) error {
	body := map[string]string{"message": message}
	return c.doRequest(ctx, "POST", fmt.Sprintf("/sessions/%s:continue", sessionID), body, nil)
}

// WaitForSessionCompletion waits for a session to complete.
func (c *Client) WaitForSessionCompletion(ctx context.Context, sessionID string, pollInterval time.Duration) (*Session, error) {
	if pollInterval == 0 {
		pollInterval = 5 * time.Second
	}

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			session, err := c.GetSession(ctx, sessionID)
			if err != nil {
				return nil, err
			}

			if session.State == StateCompleted || session.State == StateFailed {
				return session, nil
			}
		}
	}
}
