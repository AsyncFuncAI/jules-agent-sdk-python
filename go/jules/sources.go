package jules

import (
	"context"
	"fmt"
	"net/url"
)

// SourceListResponse represents a response from listing sources.
type SourceListResponse struct {
	Sources       []Source `json:"sources"`
	NextPageToken string   `json:"nextPageToken"`
}

// ListSources lists all sources.
func (c *Client) ListSources(ctx context.Context, pageSize int, pageToken string) (*SourceListResponse, error) {
	query := url.Values{}
	if pageSize > 0 {
		query.Set("pageSize", fmt.Sprintf("%d", pageSize))
	}
	if pageToken != "" {
		query.Set("pageToken", pageToken)
	}

	path := "/sources"
	if len(query) > 0 {
		path += "?" + query.Encode()
	}

	var response SourceListResponse
	if err := c.doRequest(ctx, "GET", path, nil, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
