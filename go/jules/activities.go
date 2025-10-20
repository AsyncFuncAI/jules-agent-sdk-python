package jules

import (
	"context"
	"fmt"
	"net/url"
)

// ActivityListResponse represents a response from listing activities.
type ActivityListResponse struct {
	Activities    []Activity `json:"activities"`
	NextPageToken string     `json:"nextPageToken"`
}

// ListActivities lists all activities for a session.
func (c *Client) ListActivities(ctx context.Context, sessionID string, pageSize int, pageToken string) (*ActivityListResponse, error) {
	query := url.Values{}
	if pageSize > 0 {
		query.Set("pageSize", fmt.Sprintf("%d", pageSize))
	}
	if pageToken != "" {
		query.Set("pageToken", pageToken)
	}

	path := fmt.Sprintf("/sessions/%s/activities", sessionID)
	if len(query) > 0 {
		path += "?" + query.Encode()
	}

	var response ActivityListResponse
	if err := c.doRequest(ctx, "GET", path, nil, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
