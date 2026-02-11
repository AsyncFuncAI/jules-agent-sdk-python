package jules

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// SourcesAPI provides methods for interacting with the Sources API
type SourcesAPI struct {
	client *BaseClient
}

// NewSourcesAPI creates a new SourcesAPI instance
func NewSourcesAPI(client *BaseClient) *SourcesAPI {
	return &SourcesAPI{client: client}
}

// SourcesListOptions represents options for listing sources
type SourcesListOptions struct {
	Filter    string
	PageSize  int
	PageToken string
}

// Get retrieves a source by ID
func (s *SourcesAPI) Get(ctx context.Context, sourceID string) (*Source, error) {
	path := s.buildSourcePath(sourceID)

	resp, err := s.client.Get(ctx, path)
	if err != nil {
		return nil, err
	}

	return s.parseSource(resp)
}

// List retrieves a list of sources
func (s *SourcesAPI) List(ctx context.Context, opts *SourcesListOptions) (*SourcesListResponse, error) {
	path := "/sources"

	if opts != nil {
		query := ""
		if opts.Filter != "" {
			query += fmt.Sprintf("filter=%s", opts.Filter)
		}
		if opts.PageSize > 0 {
			if query != "" {
				query += "&"
			}
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
	var result SourcesListResponse
	respBytes, _ := json.Marshal(resp)
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to parse sources list: %w", err)
	}

	return &result, nil
}

// ListAll retrieves all sources (handles pagination automatically)
func (s *SourcesAPI) ListAll(ctx context.Context, filter string) ([]Source, error) {
	var allSources []Source
	pageToken := ""

	for {
		opts := &SourcesListOptions{
			Filter:    filter,
			PageToken: pageToken,
		}

		resp, err := s.List(ctx, opts)
		if err != nil {
			return nil, err
		}

		allSources = append(allSources, resp.Sources...)

		// Check if there are more pages
		if resp.NextPageToken == "" {
			break
		}
		pageToken = resp.NextPageToken
	}

	return allSources, nil
}

// buildSourcePath builds the API path for a source
func (s *SourcesAPI) buildSourcePath(sourceID string) string {
	if strings.HasPrefix(sourceID, "sources/") {
		return "/" + sourceID
	}
	return "/sources/" + sourceID
}

// parseSource parses a source from a response
func (s *SourcesAPI) parseSource(data map[string]interface{}) (*Source, error) {
	// Convert to JSON and back to properly parse the source
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal source data: %w", err)
	}

	var source Source
	if err := json.Unmarshal(jsonBytes, &source); err != nil {
		return nil, fmt.Errorf("failed to parse source: %w", err)
	}

	return &source, nil
}
