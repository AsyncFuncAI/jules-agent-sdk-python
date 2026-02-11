package jules

import "time"

// SessionState represents the state of a session
type SessionState string

const (
	// SessionStateUnspecified is the default unspecified state
	SessionStateUnspecified SessionState = "STATE_UNSPECIFIED"
	// SessionStateQueued indicates the session is queued
	SessionStateQueued SessionState = "QUEUED"
	// SessionStatePlanning indicates the session is in planning phase
	SessionStatePlanning SessionState = "PLANNING"
	// SessionStateAwaitingPlanApproval indicates the session is waiting for plan approval
	SessionStateAwaitingPlanApproval SessionState = "AWAITING_PLAN_APPROVAL"
	// SessionStateAwaitingUserFeedback indicates the session is waiting for user feedback
	SessionStateAwaitingUserFeedback SessionState = "AWAITING_USER_FEEDBACK"
	// SessionStateInProgress indicates the session is in progress
	SessionStateInProgress SessionState = "IN_PROGRESS"
	// SessionStatePaused indicates the session is paused
	SessionStatePaused SessionState = "PAUSED"
	// SessionStateFailed indicates the session has failed
	SessionStateFailed SessionState = "FAILED"
	// SessionStateCompleted indicates the session has completed
	SessionStateCompleted SessionState = "COMPLETED"
)

// IsTerminal returns true if the session state is terminal (completed or failed)
func (s SessionState) IsTerminal() bool {
	return s == SessionStateCompleted || s == SessionStateFailed
}

// Session represents a Jules session
type Session struct {
	Name          string         `json:"name,omitempty"`
	ID            string         `json:"id,omitempty"`
	Prompt        string         `json:"prompt,omitempty"`
	SourceContext *SourceContext `json:"sourceContext,omitempty"`
	Title         string         `json:"title,omitempty"`
	State         SessionState   `json:"state,omitempty"`
	URL           string         `json:"url,omitempty"`
	CreateTime    *time.Time     `json:"createTime,omitempty"`
	UpdateTime    *time.Time     `json:"updateTime,omitempty"`
	Output        *SessionOutput `json:"output,omitempty"`
}

// SessionOutput represents the output of a session
type SessionOutput struct {
	PullRequest *PullRequest `json:"pullRequest,omitempty"`
}

// SourceContext represents the source context for a session
type SourceContext struct {
	Source           string            `json:"source,omitempty"`
	GitHubRepoContext *GitHubRepoContext `json:"githubRepoContext,omitempty"`
}

// GitHubRepoContext represents GitHub repository context
type GitHubRepoContext struct {
	StartingBranch string `json:"startingBranch,omitempty"`
}

// PullRequest represents a GitHub pull request
type PullRequest struct {
	URL         string `json:"url,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

// Source represents a source repository
type Source struct {
	ID         string      `json:"id,omitempty"`
	Name       string      `json:"name,omitempty"`
	GitHubRepo *GitHubRepo `json:"githubRepo,omitempty"`
}

// GitHubRepo represents a GitHub repository
type GitHubRepo struct {
	Owner         string         `json:"owner,omitempty"`
	Repo          string         `json:"repo,omitempty"`
	IsPrivate     bool           `json:"isPrivate,omitempty"`
	DefaultBranch string         `json:"defaultBranch,omitempty"`
	Branches      []GitHubBranch `json:"branches,omitempty"`
}

// GitHubBranch represents a GitHub branch
type GitHubBranch struct {
	DisplayName string `json:"displayName,omitempty"`
}

// Activity represents an activity in a session
type Activity struct {
	Name        string                 `json:"name,omitempty"`
	ID          string                 `json:"id,omitempty"`
	Description string                 `json:"description,omitempty"`
	CreateTime  *time.Time             `json:"createTime,omitempty"`
	Originator  string                 `json:"originator,omitempty"`
	Artifacts   []Artifact             `json:"artifacts,omitempty"`

	// Activity event types (only one should be set)
	AgentMessaged      map[string]interface{} `json:"agentMessaged,omitempty"`
	UserMessaged       map[string]interface{} `json:"userMessaged,omitempty"`
	PlanGenerated      map[string]interface{} `json:"planGenerated,omitempty"`
	PlanApproved       map[string]interface{} `json:"planApproved,omitempty"`
	ProgressUpdated    map[string]interface{} `json:"progressUpdated,omitempty"`
	SessionCompleted   map[string]interface{} `json:"sessionCompleted,omitempty"`
	SessionFailed      map[string]interface{} `json:"sessionFailed,omitempty"`
}

// Artifact represents an artifact in an activity
type Artifact struct {
	ChangeSet  *ChangeSet  `json:"changeSet,omitempty"`
	Media      *Media      `json:"media,omitempty"`
	BashOutput *BashOutput `json:"bashOutput,omitempty"`
}

// ChangeSet represents a set of changes
type ChangeSet struct {
	Source   string    `json:"source,omitempty"`
	GitPatch *GitPatch `json:"gitPatch,omitempty"`
}

// GitPatch represents a git patch
type GitPatch struct {
	UnidiffPatch           string `json:"unidiffPatch,omitempty"`
	BaseCommitID           string `json:"baseCommitId,omitempty"`
	SuggestedCommitMessage string `json:"suggestedCommitMessage,omitempty"`
}

// Media represents media content
type Media struct {
	Data     string `json:"data,omitempty"`
	MimeType string `json:"mimeType,omitempty"`
}

// BashOutput represents bash command output
type BashOutput struct {
	Command  string `json:"command,omitempty"`
	Output   string `json:"output,omitempty"`
	ExitCode int    `json:"exitCode,omitempty"`
}

// Plan represents a plan for a session
type Plan struct {
	ID         string      `json:"id,omitempty"`
	Steps      []PlanStep  `json:"steps,omitempty"`
	CreateTime *time.Time  `json:"createTime,omitempty"`
}

// PlanStep represents a step in a plan
type PlanStep struct {
	ID          string `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Index       int    `json:"index,omitempty"`
}

// ListResponse represents a paginated list response
type ListResponse struct {
	NextPageToken string `json:"nextPageToken,omitempty"`
}

// SessionsListResponse represents a list of sessions
type SessionsListResponse struct {
	Sessions      []Session `json:"sessions,omitempty"`
	NextPageToken string    `json:"nextPageToken,omitempty"`
}

// ActivitiesListResponse represents a list of activities
type ActivitiesListResponse struct {
	Activities    []Activity `json:"activities,omitempty"`
	NextPageToken string     `json:"nextPageToken,omitempty"`
}

// SourcesListResponse represents a list of sources
type SourcesListResponse struct {
	Sources       []Source `json:"sources,omitempty"`
	NextPageToken string   `json:"nextPageToken,omitempty"`
}
