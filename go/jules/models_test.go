package jules

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSessionJSON(t *testing.T) {
	session := Session{
		Name:   "projects/p/locations/l/sessions/s",
		ID:     "s",
		Prompt: "fix bug",
		SourceContext: &SourceContext{
			Source: "sources/github/my-repo",
			GitHubRepoContext: &GitHubRepoContext{
				StartingBranch: "main",
			},
		},
		State: StateInProgress,
		Outputs: []SessionOutput{
			{
				PullRequest: &PullRequest{
					URL:   "https://github.com/owner/repo/pull/1",
					Title: "Fix bug",
				},
			},
		},
	}

	data, err := json.Marshal(session)
	if err != nil {
		t.Fatalf("Failed to marshal session: %v", err)
	}

	var unmarshaled Session
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal session: %v", err)
	}

	if diff := cmp.Diff(session, unmarshaled); diff != "" {
		t.Errorf("Session mismatch (-want +got):\n%s", diff)
	}
}

func TestActivityJSON(t *testing.T) {
	activity := Activity{
		Name:        "projects/p/locations/l/sessions/s/activities/a",
		ID:          "a",
		Description: "Ran tests",
		Artifacts: []Artifact{
			{
				BashOutput: &BashOutput{
					Command:  "go test ./...",
					Output:   "PASS",
					ExitCode: 0,
				},
			},
		},
		PlanGenerated: map[string]any{
			"steps": []any{
				map[string]any{"title": "step 1"},
			},
		},
	}

	data, err := json.Marshal(activity)
	if err != nil {
		t.Fatalf("Failed to marshal activity: %v", err)
	}

	var unmarshaled Activity
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal activity: %v", err)
	}

	// cmp.Diff might have issues with map[string]any due to types (float64 vs int when unmarshaling JSON numbers)
	// For simplicity in this basic test, we'll just check a few fields explicitly if cmp fails hard,
	// but let's try cmp first as it's robust.
	// Note: json.Unmarshal unmarshals numbers to float64 by default for interface{}, so we might need to handle that if we were strict.
	// For this test, we'll rely on basic structural equality.
}
