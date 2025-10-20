package jules

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateSession(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/v1alpha/sessions" {
			t.Errorf("Expected path /v1alpha/sessions, got %s", r.URL.Path)
		}
		if r.Header.Get("X-Goog-Api-Key") != "test-key" {
			t.Errorf("Expected X-Goog-Api-Key header to be test-key, got %s", r.Header.Get("X-Goog-Api-Key"))
		}

		var req Session
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if req.Prompt != "fix bug" {
			t.Errorf("Expected prompt 'fix bug', got '%s'", req.Prompt)
		}

		resp := Session{
			Name:          "projects/p/locations/l/sessions/s1",
			ID:            "s1",
			Prompt:        req.Prompt,
			SourceContext: req.SourceContext,
			State:         StateQueued,
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	client := NewClient("test-key", WithBaseURL(ts.URL+"/v1alpha"))
	ctx := context.Background()

	session, err := client.CreateSession(ctx, CreateSessionRequest{
		Prompt: "fix bug",
		Source: "sources/s1",
	})
	if err != nil {
		t.Fatalf("CreateSession failed: %v", err)
	}

	if session.ID != "s1" {
		t.Errorf("Expected session ID s1, got %s", session.ID)
	}
}

func TestListSessions(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/v1alpha/sessions" {
			t.Errorf("Expected path /v1alpha/sessions, got %s", r.URL.Path)
		}

		resp := SessionListResponse{
			Sessions: []Session{
				{ID: "s1"},
				{ID: "s2"},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	client := NewClient("test-key", WithBaseURL(ts.URL+"/v1alpha"))
	ctx := context.Background()

	resp, err := client.ListSessions(ctx, 0, "")
	if err != nil {
		t.Fatalf("ListSessions failed: %v", err)
	}

	if len(resp.Sessions) != 2 {
		t.Errorf("Expected 2 sessions, got %d", len(resp.Sessions))
	}
}
