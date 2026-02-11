package jules

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestSourcesAPIGet(t *testing.T) {
	mockServer := NewMockServer(t, []MockResponse{
		{
			StatusCode: 200,
			Body: map[string]interface{}{
				"id":   "source-123",
				"name": "sources/source-123",
				"githubRepo": map[string]interface{}{
					"owner":         "testowner",
					"repo":          "testrepo",
					"isPrivate":     true,
					"defaultBranch": "main",
				},
			},
		},
	})
	defer mockServer.Close()

	config := NewTestConfig(mockServer.URL, "test-key")
	baseClient, _ := NewBaseClient(config)
	defer baseClient.Close()

	api := NewSourcesAPI(baseClient)

	ctx := context.Background()
	source, err := api.Get(ctx, "source-123")

	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if source.ID != "source-123" {
		t.Errorf("Source.ID = %v, want source-123", source.ID)
	}
	if source.GitHubRepo == nil {
		t.Fatal("Expected GitHubRepo to be non-nil")
	}
	if source.GitHubRepo.Owner != "testowner" {
		t.Errorf("GitHubRepo.Owner = %v, want testowner", source.GitHubRepo.Owner)
	}
	if source.GitHubRepo.Repo != "testrepo" {
		t.Errorf("GitHubRepo.Repo = %v, want testrepo", source.GitHubRepo.Repo)
	}
	if !source.GitHubRepo.IsPrivate {
		t.Error("Expected GitHubRepo.IsPrivate to be true")
	}

	// Verify request path
	req := mockServer.GetLastRequest()
	if !strings.Contains(req.URL.Path, "source-123") {
		t.Errorf("Request path should contain source-123")
	}
}

func TestSourcesAPIList(t *testing.T) {
	mockServer := NewMockServer(t, []MockResponse{
		{
			StatusCode: 200,
			Body: map[string]interface{}{
				"sources": []interface{}{
					map[string]interface{}{
						"id":   "source-1",
						"name": "sources/repo-1",
					},
					map[string]interface{}{
						"id":   "source-2",
						"name": "sources/repo-2",
					},
				},
				"nextPageToken": "next-token-abc",
			},
		},
	})
	defer mockServer.Close()

	config := NewTestConfig(mockServer.URL, "test-key")
	baseClient, _ := NewBaseClient(config)
	defer baseClient.Close()

	api := NewSourcesAPI(baseClient)

	ctx := context.Background()
	response, err := api.List(ctx, &SourcesListOptions{
		PageSize:  10,
		PageToken: "",
	})

	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(response.Sources) != 2 {
		t.Errorf("Sources count = %v, want 2", len(response.Sources))
	}
	if response.NextPageToken != "next-token-abc" {
		t.Errorf("NextPageToken = %v, want next-token-abc", response.NextPageToken)
	}
	if response.Sources[0].ID != "source-1" {
		t.Errorf("First source ID = %v, want source-1", response.Sources[0].ID)
	}

	// Verify request
	req := mockServer.GetLastRequest()
	if req.Method != http.MethodGet {
		t.Errorf("Request method = %v, want GET", req.Method)
	}
}

func TestSourcesAPIListWithFilter(t *testing.T) {
	mockServer := NewMockServer(t, []MockResponse{
		{
			StatusCode: 200,
			Body: map[string]interface{}{
				"sources":       []interface{}{},
				"nextPageToken": "",
			},
		},
	})
	defer mockServer.Close()

	config := NewTestConfig(mockServer.URL, "test-key")
	baseClient, _ := NewBaseClient(config)
	defer baseClient.Close()

	api := NewSourcesAPI(baseClient)

	ctx := context.Background()
	_, err := api.List(ctx, &SourcesListOptions{
		Filter:    "owner:testorg",
		PageSize:  5,
		PageToken: "",
	})

	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	// Verify query params
	req := mockServer.GetLastRequest()
	query := req.URL.Query()
	if query.Get("filter") != "owner:testorg" {
		t.Errorf("filter = %v, want owner:testorg", query.Get("filter"))
	}
	if query.Get("pageSize") != "5" {
		t.Errorf("pageSize = %v, want 5", query.Get("pageSize"))
	}
}

func TestSourcesAPIListAll(t *testing.T) {
	// Mock two pages of results
	mockServer := NewMockServer(t, []MockResponse{
		{
			StatusCode: 200,
			Body: map[string]interface{}{
				"sources": []interface{}{
					map[string]interface{}{"id": "source-1", "name": "sources/repo-1"},
					map[string]interface{}{"id": "source-2", "name": "sources/repo-2"},
				},
				"nextPageToken": "page-2-token",
			},
		},
		{
			StatusCode: 200,
			Body: map[string]interface{}{
				"sources": []interface{}{
					map[string]interface{}{"id": "source-3", "name": "sources/repo-3"},
				},
				"nextPageToken": "",
			},
		},
	})
	defer mockServer.Close()

	config := NewTestConfig(mockServer.URL, "test-key")
	baseClient, _ := NewBaseClient(config)
	defer baseClient.Close()

	api := NewSourcesAPI(baseClient)

	ctx := context.Background()
	sources, err := api.ListAll(ctx, "")

	if err != nil {
		t.Fatalf("ListAll() error = %v", err)
	}

	if len(sources) != 3 {
		t.Errorf("Sources count = %v, want 3", len(sources))
	}

	// Verify all sources are present
	expectedIDs := []string{"source-1", "source-2", "source-3"}
	for i, expected := range expectedIDs {
		if sources[i].ID != expected {
			t.Errorf("Source[%d].ID = %v, want %v", i, sources[i].ID, expected)
		}
	}

	// Verify pagination worked (should have made 2 requests)
	if mockServer.GetRequestCount() != 2 {
		t.Errorf("Request count = %v, want 2", mockServer.GetRequestCount())
	}
}

func TestSourcesAPIListAllWithFilter(t *testing.T) {
	mockServer := NewMockServer(t, []MockResponse{
		{
			StatusCode: 200,
			Body: map[string]interface{}{
				"sources": []interface{}{
					map[string]interface{}{"id": "source-1"},
				},
				"nextPageToken": "",
			},
		},
	})
	defer mockServer.Close()

	config := NewTestConfig(mockServer.URL, "test-key")
	baseClient, _ := NewBaseClient(config)
	defer baseClient.Close()

	api := NewSourcesAPI(baseClient)

	ctx := context.Background()
	sources, err := api.ListAll(ctx, "owner:myorg")

	if err != nil {
		t.Fatalf("ListAll() error = %v", err)
	}

	if len(sources) != 1 {
		t.Errorf("Sources count = %v, want 1", len(sources))
	}

	// Verify filter was passed
	req := mockServer.GetLastRequest()
	query := req.URL.Query()
	if query.Get("filter") != "owner:myorg" {
		t.Errorf("filter = %v, want owner:myorg", query.Get("filter"))
	}
}

func TestSourcesAPIGetWithFullName(t *testing.T) {
	mockServer := NewMockServer(t, []MockResponse{
		{
			StatusCode: 200,
			Body: map[string]interface{}{
				"id":   "source-456",
				"name": "sources/source-456",
			},
		},
	})
	defer mockServer.Close()

	config := NewTestConfig(mockServer.URL, "test-key")
	baseClient, _ := NewBaseClient(config)
	defer baseClient.Close()

	api := NewSourcesAPI(baseClient)

	ctx := context.Background()
	source, err := api.Get(ctx, "sources/source-456")

	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if source.ID != "source-456" {
		t.Errorf("Source.ID = %v, want source-456", source.ID)
	}
}

func TestBuildSourcePath(t *testing.T) {
	api := &SourcesAPI{}

	tests := []struct {
		input string
		want  string
	}{
		{"source-123", "/sources/source-123"},
		{"sources/source-456", "/sources/source-456"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := api.buildSourcePath(tt.input)
			if got != tt.want {
				t.Errorf("buildSourcePath(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestSourcesAPIWithBranches(t *testing.T) {
	mockServer := NewMockServer(t, []MockResponse{
		{
			StatusCode: 200,
			Body: map[string]interface{}{
				"id":   "source-789",
				"name": "sources/source-789",
				"githubRepo": map[string]interface{}{
					"owner": "testowner",
					"repo":  "testrepo",
					"branches": []interface{}{
						map[string]interface{}{"displayName": "main"},
						map[string]interface{}{"displayName": "develop"},
						map[string]interface{}{"displayName": "feature/test"},
					},
				},
			},
		},
	})
	defer mockServer.Close()

	config := NewTestConfig(mockServer.URL, "test-key")
	baseClient, _ := NewBaseClient(config)
	defer baseClient.Close()

	api := NewSourcesAPI(baseClient)

	ctx := context.Background()
	source, err := api.Get(ctx, "source-789")

	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if source.GitHubRepo == nil {
		t.Fatal("Expected GitHubRepo to be non-nil")
	}

	if len(source.GitHubRepo.Branches) != 3 {
		t.Errorf("Branches count = %v, want 3", len(source.GitHubRepo.Branches))
	}

	expectedBranches := []string{"main", "develop", "feature/test"}
	for i, expected := range expectedBranches {
		if source.GitHubRepo.Branches[i].DisplayName != expected {
			t.Errorf("Branch[%d].DisplayName = %v, want %v", i, source.GitHubRepo.Branches[i].DisplayName, expected)
		}
	}
}

func TestSourcesAPIListEmpty(t *testing.T) {
	mockServer := NewMockServer(t, []MockResponse{
		{
			StatusCode: 200,
			Body: map[string]interface{}{
				"sources":       []interface{}{},
				"nextPageToken": "",
			},
		},
	})
	defer mockServer.Close()

	config := NewTestConfig(mockServer.URL, "test-key")
	baseClient, _ := NewBaseClient(config)
	defer baseClient.Close()

	api := NewSourcesAPI(baseClient)

	ctx := context.Background()
	response, err := api.List(ctx, nil)

	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(response.Sources) != 0 {
		t.Errorf("Sources count = %v, want 0", len(response.Sources))
	}
	if response.NextPageToken != "" {
		t.Errorf("NextPageToken = %v, want empty string", response.NextPageToken)
	}
}
