package jules

import "fmt"

// SourcesService is the service for interacting with the Sources API.
type SourcesService struct {
	client *Client
}

// NewSourcesService creates a new SourcesService.
func NewSourcesService(client *Client) *SourcesService {
	return &SourcesService{client: client}
}

// Get retrieves a source by ID.
func (s *SourcesService) Get(sourceID string) (*Source, error) {
	var source Source
	path := fmt.Sprintf("sources/%s", sourceID)
	_, err := s.client.get(path, &source)
	return &source, err
}

type ListSourcesResponse struct {
	Sources       []Source `json:"sources"`
	NextPageToken string   `json:"nextPageToken"`
}

// List lists all sources.
func (s *SourcesService) List() (*ListSourcesResponse, error) {
	var resp ListSourcesResponse
	_, err := s.client.get("sources", &resp)
	return &resp, err
}