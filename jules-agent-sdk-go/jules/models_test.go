package jules

import (
	"encoding/json"
	"testing"
	"time"
)

func TestSessionStateIsTerminal(t *testing.T) {
	tests := []struct {
		name     string
		state    SessionState
		terminal bool
	}{
		{"unspecified", SessionStateUnspecified, false},
		{"queued", SessionStateQueued, false},
		{"planning", SessionStatePlanning, false},
		{"awaiting plan approval", SessionStateAwaitingPlanApproval, false},
		{"awaiting user feedback", SessionStateAwaitingUserFeedback, false},
		{"in progress", SessionStateInProgress, false},
		{"paused", SessionStatePaused, false},
		{"failed", SessionStateFailed, true},
		{"completed", SessionStateCompleted, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.state.IsTerminal(); got != tt.terminal {
				t.Errorf("SessionState.IsTerminal() = %v, want %v", got, tt.terminal)
			}
		})
	}
}

func TestSessionJSONMarshaling(t *testing.T) {
	now := time.Now()
	session := &Session{
		Name:   "sessions/test-123",
		ID:     "test-123",
		Prompt: "test prompt",
		SourceContext: &SourceContext{
			Source: "sources/test-repo",
			GitHubRepoContext: &GitHubRepoContext{
				StartingBranch: "main",
			},
		},
		Title:      "Test Session",
		State:      SessionStateInProgress,
		URL:        "https://example.com/session/test-123",
		CreateTime: &now,
		UpdateTime: &now,
		Output: &SessionOutput{
			PullRequest: &PullRequest{
				URL:         "https://github.com/owner/repo/pull/1",
				Title:       "Test PR",
				Description: "Test description",
			},
		},
	}

	// Marshal to JSON
	data, err := json.Marshal(session)
	if err != nil {
		t.Fatalf("Failed to marshal session: %v", err)
	}

	// Unmarshal back
	var unmarshaled Session
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal session: %v", err)
	}

	// Verify fields
	if unmarshaled.Name != session.Name {
		t.Errorf("Session.Name = %v, want %v", unmarshaled.Name, session.Name)
	}
	if unmarshaled.ID != session.ID {
		t.Errorf("Session.ID = %v, want %v", unmarshaled.ID, session.ID)
	}
	if unmarshaled.Prompt != session.Prompt {
		t.Errorf("Session.Prompt = %v, want %v", unmarshaled.Prompt, session.Prompt)
	}
	if unmarshaled.State != session.State {
		t.Errorf("Session.State = %v, want %v", unmarshaled.State, session.State)
	}
	if unmarshaled.SourceContext.Source != session.SourceContext.Source {
		t.Errorf("Session.SourceContext.Source = %v, want %v", unmarshaled.SourceContext.Source, session.SourceContext.Source)
	}
}

func TestActivityJSONMarshaling(t *testing.T) {
	now := time.Now()
	activity := &Activity{
		Name:        "activities/test-456",
		ID:          "test-456",
		Description: "Test activity",
		CreateTime:  &now,
		Originator:  "AGENT",
		Artifacts: []Artifact{
			{
				BashOutput: &BashOutput{
					Command:  "ls -la",
					Output:   "total 0",
					ExitCode: 0,
				},
			},
		},
		AgentMessaged: map[string]interface{}{
			"message": "test message",
		},
	}

	// Marshal to JSON
	data, err := json.Marshal(activity)
	if err != nil {
		t.Fatalf("Failed to marshal activity: %v", err)
	}

	// Unmarshal back
	var unmarshaled Activity
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal activity: %v", err)
	}

	// Verify fields
	if unmarshaled.Name != activity.Name {
		t.Errorf("Activity.Name = %v, want %v", unmarshaled.Name, activity.Name)
	}
	if unmarshaled.ID != activity.ID {
		t.Errorf("Activity.ID = %v, want %v", unmarshaled.ID, activity.ID)
	}
	if unmarshaled.Description != activity.Description {
		t.Errorf("Activity.Description = %v, want %v", unmarshaled.Description, activity.Description)
	}
	if unmarshaled.Originator != activity.Originator {
		t.Errorf("Activity.Originator = %v, want %v", unmarshaled.Originator, activity.Originator)
	}
}

func TestSourceJSONMarshaling(t *testing.T) {
	source := &Source{
		ID:   "test-source",
		Name: "sources/test-source",
		GitHubRepo: &GitHubRepo{
			Owner:         "testowner",
			Repo:          "testrepo",
			IsPrivate:     true,
			DefaultBranch: "main",
			Branches: []GitHubBranch{
				{DisplayName: "main"},
				{DisplayName: "develop"},
			},
		},
	}

	// Marshal to JSON
	data, err := json.Marshal(source)
	if err != nil {
		t.Fatalf("Failed to marshal source: %v", err)
	}

	// Unmarshal back
	var unmarshaled Source
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal source: %v", err)
	}

	// Verify fields
	if unmarshaled.ID != source.ID {
		t.Errorf("Source.ID = %v, want %v", unmarshaled.ID, source.ID)
	}
	if unmarshaled.Name != source.Name {
		t.Errorf("Source.Name = %v, want %v", unmarshaled.Name, source.Name)
	}
	if unmarshaled.GitHubRepo.Owner != source.GitHubRepo.Owner {
		t.Errorf("Source.GitHubRepo.Owner = %v, want %v", unmarshaled.GitHubRepo.Owner, source.GitHubRepo.Owner)
	}
	if unmarshaled.GitHubRepo.Repo != source.GitHubRepo.Repo {
		t.Errorf("Source.GitHubRepo.Repo = %v, want %v", unmarshaled.GitHubRepo.Repo, source.GitHubRepo.Repo)
	}
	if unmarshaled.GitHubRepo.IsPrivate != source.GitHubRepo.IsPrivate {
		t.Errorf("Source.GitHubRepo.IsPrivate = %v, want %v", unmarshaled.GitHubRepo.IsPrivate, source.GitHubRepo.IsPrivate)
	}
	if len(unmarshaled.GitHubRepo.Branches) != len(source.GitHubRepo.Branches) {
		t.Errorf("Source.GitHubRepo.Branches length = %v, want %v", len(unmarshaled.GitHubRepo.Branches), len(source.GitHubRepo.Branches))
	}
}

func TestArtifactJSONMarshaling(t *testing.T) {
	artifact := &Artifact{
		ChangeSet: &ChangeSet{
			Source: "sources/test",
			GitPatch: &GitPatch{
				UnidiffPatch:           "diff --git a/file.txt b/file.txt",
				BaseCommitID:           "abc123",
				SuggestedCommitMessage: "Update file",
			},
		},
		Media: &Media{
			Data:     "base64data",
			MimeType: "image/png",
		},
		BashOutput: &BashOutput{
			Command:  "echo test",
			Output:   "test",
			ExitCode: 0,
		},
	}

	// Marshal to JSON
	data, err := json.Marshal(artifact)
	if err != nil {
		t.Fatalf("Failed to marshal artifact: %v", err)
	}

	// Unmarshal back
	var unmarshaled Artifact
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal artifact: %v", err)
	}

	// Verify fields
	if unmarshaled.ChangeSet.Source != artifact.ChangeSet.Source {
		t.Errorf("Artifact.ChangeSet.Source = %v, want %v", unmarshaled.ChangeSet.Source, artifact.ChangeSet.Source)
	}
	if unmarshaled.ChangeSet.GitPatch.BaseCommitID != artifact.ChangeSet.GitPatch.BaseCommitID {
		t.Errorf("Artifact.ChangeSet.GitPatch.BaseCommitID = %v, want %v", unmarshaled.ChangeSet.GitPatch.BaseCommitID, artifact.ChangeSet.GitPatch.BaseCommitID)
	}
	if unmarshaled.Media.MimeType != artifact.Media.MimeType {
		t.Errorf("Artifact.Media.MimeType = %v, want %v", unmarshaled.Media.MimeType, artifact.Media.MimeType)
	}
	if unmarshaled.BashOutput.ExitCode != artifact.BashOutput.ExitCode {
		t.Errorf("Artifact.BashOutput.ExitCode = %v, want %v", unmarshaled.BashOutput.ExitCode, artifact.BashOutput.ExitCode)
	}
}

func TestPlanJSONMarshaling(t *testing.T) {
	now := time.Now()
	plan := &Plan{
		ID: "plan-123",
		Steps: []PlanStep{
			{
				ID:          "step-1",
				Title:       "Step 1",
				Description: "First step",
				Index:       0,
			},
			{
				ID:          "step-2",
				Title:       "Step 2",
				Description: "Second step",
				Index:       1,
			},
		},
		CreateTime: &now,
	}

	// Marshal to JSON
	data, err := json.Marshal(plan)
	if err != nil {
		t.Fatalf("Failed to marshal plan: %v", err)
	}

	// Unmarshal back
	var unmarshaled Plan
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal plan: %v", err)
	}

	// Verify fields
	if unmarshaled.ID != plan.ID {
		t.Errorf("Plan.ID = %v, want %v", unmarshaled.ID, plan.ID)
	}
	if len(unmarshaled.Steps) != len(plan.Steps) {
		t.Errorf("Plan.Steps length = %v, want %v", len(unmarshaled.Steps), len(plan.Steps))
	}
	if len(unmarshaled.Steps) > 0 {
		if unmarshaled.Steps[0].Title != plan.Steps[0].Title {
			t.Errorf("Plan.Steps[0].Title = %v, want %v", unmarshaled.Steps[0].Title, plan.Steps[0].Title)
		}
		if unmarshaled.Steps[0].Index != plan.Steps[0].Index {
			t.Errorf("Plan.Steps[0].Index = %v, want %v", unmarshaled.Steps[0].Index, plan.Steps[0].Index)
		}
	}
}
