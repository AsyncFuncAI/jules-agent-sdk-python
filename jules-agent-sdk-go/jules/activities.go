package jules

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// ActivitiesAPI provides methods for interacting with the Activities API
type ActivitiesAPI struct {
	client *BaseClient
}

// NewActivitiesAPI creates a new ActivitiesAPI instance
func NewActivitiesAPI(client *BaseClient) *ActivitiesAPI {
	return &ActivitiesAPI{client: client}
}

// Get retrieves an activity by ID
func (a *ActivitiesAPI) Get(ctx context.Context, sessionID, activityID string) (*Activity, error) {
	path := a.buildActivityPath(sessionID, activityID)

	resp, err := a.client.Get(ctx, path)
	if err != nil {
		return nil, err
	}

	return a.parseActivity(resp)
}

// List retrieves a list of activities for a session
func (a *ActivitiesAPI) List(ctx context.Context, sessionID string, opts *ListOptions) (*ActivitiesListResponse, error) {
	sessionPath := a.buildSessionPath(sessionID)
	path := sessionPath + "/activities"

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

	resp, err := a.client.Get(ctx, path)
	if err != nil {
		return nil, err
	}

	// Parse response
	var result ActivitiesListResponse
	respBytes, _ := json.Marshal(resp)
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to parse activities list: %w", err)
	}

	return &result, nil
}

// ListAll retrieves all activities for a session (handles pagination automatically)
func (a *ActivitiesAPI) ListAll(ctx context.Context, sessionID string) ([]Activity, error) {
	var allActivities []Activity
	pageToken := ""

	for {
		opts := &ListOptions{
			PageToken: pageToken,
		}

		resp, err := a.List(ctx, sessionID, opts)
		if err != nil {
			return nil, err
		}

		allActivities = append(allActivities, resp.Activities...)

		// Check if there are more pages
		if resp.NextPageToken == "" {
			break
		}
		pageToken = resp.NextPageToken
	}

	return allActivities, nil
}

// buildSessionPath builds the API path for a session
func (a *ActivitiesAPI) buildSessionPath(sessionID string) string {
	if strings.HasPrefix(sessionID, "sessions/") {
		return "/" + sessionID
	}
	return "/sessions/" + sessionID
}

// buildActivityPath builds the API path for an activity
func (a *ActivitiesAPI) buildActivityPath(sessionID, activityID string) string {
	sessionPath := a.buildSessionPath(sessionID)

	if strings.HasPrefix(activityID, "activities/") {
		return sessionPath + "/" + activityID
	}
	return sessionPath + "/activities/" + activityID
}

// parseActivity parses an activity from a response
func (a *ActivitiesAPI) parseActivity(data map[string]interface{}) (*Activity, error) {
	// Convert to JSON and back to properly parse the activity
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal activity data: %w", err)
	}

	var activity Activity
	if err := json.Unmarshal(jsonBytes, &activity); err != nil {
		return nil, fmt.Errorf("failed to parse activity: %w", err)
	}

	return &activity, nil
}
