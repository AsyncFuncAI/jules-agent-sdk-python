package jules

import (
	"context"
	"strings"
	"testing"
)

func TestActivitiesAPIGet(t *testing.T) {
	mockServer := NewMockServer(t, []MockResponse{
		{
			StatusCode: 200,
			Body: map[string]interface{}{
				"name":        "activities/activity-123",
				"id":          "activity-123",
				"description": "Test activity",
				"originator":  "AGENT",
			},
		},
	})
	defer mockServer.Close()

	config := NewTestConfig(mockServer.URL, "test-key")
	baseClient, _ := NewBaseClient(config)
	defer baseClient.Close()

	api := NewActivitiesAPI(baseClient)

	ctx := context.Background()
	activity, err := api.Get(ctx, "session-123", "activity-123")

	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if activity.ID != "activity-123" {
		t.Errorf("Activity.ID = %v, want activity-123", activity.ID)
	}
	if activity.Originator != "AGENT" {
		t.Errorf("Activity.Originator = %v, want AGENT", activity.Originator)
	}

	// Verify request path
	req := mockServer.GetLastRequest()
	if !strings.Contains(req.URL.Path, "session-123") {
		t.Errorf("Request path should contain session-123")
	}
	if !strings.Contains(req.URL.Path, "activity-123") {
		t.Errorf("Request path should contain activity-123")
	}
}

func TestActivitiesAPIList(t *testing.T) {
	mockServer := NewMockServer(t, []MockResponse{
		{
			StatusCode: 200,
			Body: map[string]interface{}{
				"activities": []interface{}{
					map[string]interface{}{
						"id":          "activity-1",
						"description": "First activity",
						"originator":  "AGENT",
					},
					map[string]interface{}{
						"id":          "activity-2",
						"description": "Second activity",
						"originator":  "USER",
					},
				},
				"nextPageToken": "next-token",
			},
		},
	})
	defer mockServer.Close()

	config := NewTestConfig(mockServer.URL, "test-key")
	baseClient, _ := NewBaseClient(config)
	defer baseClient.Close()

	api := NewActivitiesAPI(baseClient)

	ctx := context.Background()
	response, err := api.List(ctx, "session-123", &ListOptions{
		PageSize:  10,
		PageToken: "",
	})

	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(response.Activities) != 2 {
		t.Errorf("Activities count = %v, want 2", len(response.Activities))
	}
	if response.NextPageToken != "next-token" {
		t.Errorf("NextPageToken = %v, want next-token", response.NextPageToken)
	}
	if response.Activities[0].ID != "activity-1" {
		t.Errorf("First activity ID = %v, want activity-1", response.Activities[0].ID)
	}

	// Verify request
	req := mockServer.GetLastRequest()
	if !strings.Contains(req.URL.Path, "session-123") {
		t.Errorf("Request path should contain session-123")
	}
	if !strings.Contains(req.URL.Path, "activities") {
		t.Errorf("Request path should contain activities")
	}
}

func TestActivitiesAPIListAll(t *testing.T) {
	// Mock two pages of results
	mockServer := NewMockServer(t, []MockResponse{
		{
			StatusCode: 200,
			Body: map[string]interface{}{
				"activities": []interface{}{
					map[string]interface{}{"id": "activity-1", "description": "First"},
					map[string]interface{}{"id": "activity-2", "description": "Second"},
				},
				"nextPageToken": "page-2-token",
			},
		},
		{
			StatusCode: 200,
			Body: map[string]interface{}{
				"activities": []interface{}{
					map[string]interface{}{"id": "activity-3", "description": "Third"},
					map[string]interface{}{"id": "activity-4", "description": "Fourth"},
				},
				"nextPageToken": "",
			},
		},
	})
	defer mockServer.Close()

	config := NewTestConfig(mockServer.URL, "test-key")
	baseClient, _ := NewBaseClient(config)
	defer baseClient.Close()

	api := NewActivitiesAPI(baseClient)

	ctx := context.Background()
	activities, err := api.ListAll(ctx, "session-123")

	if err != nil {
		t.Fatalf("ListAll() error = %v", err)
	}

	if len(activities) != 4 {
		t.Errorf("Activities count = %v, want 4", len(activities))
	}

	// Verify all activities are present
	expectedIDs := []string{"activity-1", "activity-2", "activity-3", "activity-4"}
	for i, expected := range expectedIDs {
		if activities[i].ID != expected {
			t.Errorf("Activity[%d].ID = %v, want %v", i, activities[i].ID, expected)
		}
	}

	// Verify pagination worked (should have made 2 requests)
	if mockServer.GetRequestCount() != 2 {
		t.Errorf("Request count = %v, want 2", mockServer.GetRequestCount())
	}
}

func TestActivitiesAPIListWithFullSessionName(t *testing.T) {
	mockServer := NewMockServer(t, []MockResponse{
		{
			StatusCode: 200,
			Body: map[string]interface{}{
				"activities":     []interface{}{},
				"nextPageToken": "",
			},
		},
	})
	defer mockServer.Close()

	config := NewTestConfig(mockServer.URL, "test-key")
	baseClient, _ := NewBaseClient(config)
	defer baseClient.Close()

	api := NewActivitiesAPI(baseClient)

	ctx := context.Background()
	_, err := api.List(ctx, "sessions/session-456", nil)

	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	// Verify request path handles full session name
	req := mockServer.GetLastRequest()
	if !strings.Contains(req.URL.Path, "sessions/session-456") {
		t.Errorf("Request path = %v, should contain sessions/session-456", req.URL.Path)
	}
}

func TestActivitiesAPIGetWithFullNames(t *testing.T) {
	mockServer := NewMockServer(t, []MockResponse{
		{
			StatusCode: 200,
			Body: map[string]interface{}{
				"id": "activity-789",
			},
		},
	})
	defer mockServer.Close()

	config := NewTestConfig(mockServer.URL, "test-key")
	baseClient, _ := NewBaseClient(config)
	defer baseClient.Close()

	api := NewActivitiesAPI(baseClient)

	ctx := context.Background()
	activity, err := api.Get(ctx, "sessions/session-789", "activities/activity-789")

	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if activity.ID != "activity-789" {
		t.Errorf("Activity.ID = %v, want activity-789", activity.ID)
	}
}

func TestBuildActivityPath(t *testing.T) {
	api := &ActivitiesAPI{}

	tests := []struct {
		sessionID  string
		activityID string
		want       string
	}{
		{
			"session-1",
			"activity-1",
			"/sessions/session-1/activities/activity-1",
		},
		{
			"sessions/session-2",
			"activity-2",
			"/sessions/session-2/activities/activity-2",
		},
		{
			"session-3",
			"activities/activity-3",
			"/sessions/session-3/activities/activity-3",
		},
		{
			"sessions/session-4",
			"activities/activity-4",
			"/sessions/session-4/activities/activity-4",
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := api.buildActivityPath(tt.sessionID, tt.activityID)
			if got != tt.want {
				t.Errorf("buildActivityPath(%v, %v) = %v, want %v", tt.sessionID, tt.activityID, got, tt.want)
			}
		})
	}
}

func TestActivitiesAPIWithArtifacts(t *testing.T) {
	mockServer := NewMockServer(t, []MockResponse{
		{
			StatusCode: 200,
			Body: map[string]interface{}{
				"id": "activity-123",
				"artifacts": []interface{}{
					map[string]interface{}{
						"bashOutput": map[string]interface{}{
							"command":  "ls -la",
							"output":   "total 0",
							"exitCode": float64(0),
						},
					},
					map[string]interface{}{
						"media": map[string]interface{}{
							"data":     "base64data",
							"mimeType": "image/png",
						},
					},
				},
			},
		},
	})
	defer mockServer.Close()

	config := NewTestConfig(mockServer.URL, "test-key")
	baseClient, _ := NewBaseClient(config)
	defer baseClient.Close()

	api := NewActivitiesAPI(baseClient)

	ctx := context.Background()
	activity, err := api.Get(ctx, "session-123", "activity-123")

	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if len(activity.Artifacts) != 2 {
		t.Errorf("Artifacts count = %v, want 2", len(activity.Artifacts))
	}

	// Verify bash output artifact
	if activity.Artifacts[0].BashOutput == nil {
		t.Error("First artifact should have BashOutput")
	} else {
		if activity.Artifacts[0].BashOutput.Command != "ls -la" {
			t.Errorf("BashOutput.Command = %v, want ls -la", activity.Artifacts[0].BashOutput.Command)
		}
	}

	// Verify media artifact
	if activity.Artifacts[1].Media == nil {
		t.Error("Second artifact should have Media")
	} else {
		if activity.Artifacts[1].Media.MimeType != "image/png" {
			t.Errorf("Media.MimeType = %v, want image/png", activity.Artifacts[1].Media.MimeType)
		}
	}
}

func TestActivitiesAPIListWithPagination(t *testing.T) {
	mockServer := NewMockServer(t, []MockResponse{
		{
			StatusCode: 200,
			Body: map[string]interface{}{
				"activities": []interface{}{
					map[string]interface{}{"id": "a1"},
				},
				"nextPageToken": "",
			},
		},
	})
	defer mockServer.Close()

	config := NewTestConfig(mockServer.URL, "test-key")
	baseClient, _ := NewBaseClient(config)
	defer baseClient.Close()

	api := NewActivitiesAPI(baseClient)

	ctx := context.Background()
	response, err := api.List(ctx, "session-123", &ListOptions{
		PageSize:  5,
		PageToken: "existing-token",
	})

	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(response.Activities) != 1 {
		t.Errorf("Activities count = %v, want 1", len(response.Activities))
	}

	// Verify query params
	req := mockServer.GetLastRequest()
	query := req.URL.Query()
	if query.Get("pageSize") != "5" {
		t.Errorf("pageSize = %v, want 5", query.Get("pageSize"))
	}
	if query.Get("pageToken") != "existing-token" {
		t.Errorf("pageToken = %v, want existing-token", query.Get("pageToken"))
	}
}
