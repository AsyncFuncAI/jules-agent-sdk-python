package jules

import "fmt"

// ActivitiesService is the service for interacting with the Activities API.
type ActivitiesService struct {
	client *Client
}

// NewActivitiesService creates a new ActivitiesService.
func NewActivitiesService(client *Client) *ActivitiesService {
	return &ActivitiesService{client: client}
}

// Get retrieves an activity by session and activity ID.
func (s *ActivitiesService) Get(sessionID, activityID string) (*Activity, error) {
	var activity Activity
	path := fmt.Sprintf("sessions/%s/activities/%s", sessionID, activityID)
	_, err := s.client.get(path, &activity)
	return &activity, err
}

type ListActivitiesResponse struct {
	Activities    []Activity `json:"activities"`
	NextPageToken string     `json:"nextPageToken"`
}

// List lists all activities for a session.
func (s *ActivitiesService) List(sessionID string) (*ListActivitiesResponse, error) {
	var resp ListActivitiesResponse
	path := fmt.Sprintf("sessions/%s/activities", sessionID)
	_, err := s.client.get(path, &resp)
	return &resp, err
}