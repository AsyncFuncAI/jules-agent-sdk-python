package jules

import (
	"fmt"
	"time"
)

const (
	DefaultPollInterval = 5 * time.Second
)

// SessionsService is the service for interacting with the Sessions API.
type SessionsService struct {
	client *Client
}

// NewSessionsService creates a new SessionsService.
func NewSessionsService(client *Client) *SessionsService {
	return &SessionsService{client: client}
}

type CreateSessionRequest struct {
	Prompt              string         `json:"prompt"`
	SourceContext       *SourceContext `json:"sourceContext"`
	Title               string         `json:"title,omitempty"`
	RequirePlanApproval bool           `json:"requirePlanApproval,omitempty"`
}

// Create creates a new session.
func (s *SessionsService) Create(req *CreateSessionRequest) (*Session, error) {
	var session Session
	_, err := s.client.post("sessions", req, &session)
	return &session, err
}

// Get retrieves a session by ID.
func (s *SessionsService) Get(sessionID string) (*Session, error) {
	var session Session
	path := fmt.Sprintf("sessions/%s", sessionID)
	_, err := s.client.get(path, &session)
	return &session, err
}

type ListSessionsResponse struct {
	Sessions      []Session `json:"sessions"`
	NextPageToken string    `json:"nextPageToken"`
}

// List lists all sessions.
func (s *SessionsService) List() (*ListSessionsResponse, error) {
	var resp ListSessionsResponse
	_, err := s.client.get("sessions", &resp)
	return &resp, err
}

// ApprovePlan approves a session's plan.
func (s *SessionsService) ApprovePlan(sessionID string) error {
	path := fmt.Sprintf("sessions/%s:approvePlan", sessionID)
	_, err := s.client.post(path, nil, nil)
	return err
}

type SendMessageRequest struct {
	Prompt string `json:"prompt"`
}

// SendMessage sends a message to a session.
func (s *SessionsService) SendMessage(sessionID string, prompt string) error {
	path := fmt.Sprintf("sessions/%s:sendMessage", sessionID)
	req := &SendMessageRequest{Prompt: prompt}
	_, err := s.client.post(path, req, nil)
	return err
}

// WaitForCompletion polls a session until it completes or fails.
func (s *SessionsService) WaitForCompletion(sessionID string) (*Session, error) {
	timeout := time.After(DefaultTimeout)
	ticker := time.NewTicker(DefaultPollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return nil, fmt.Errorf("timed out waiting for session %s to complete", sessionID)
		case <-ticker.C:
			session, err := s.Get(sessionID)
			if err != nil {
				return nil, err
			}

			switch session.State {
			case Completed:
				return session, nil
			case Failed:
				return session, fmt.Errorf("session %s failed", sessionID)
			}
		}
	}
}