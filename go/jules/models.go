package jules

// SessionState represents the state of a session.
type SessionState string

const (
	StateUnspecified      SessionState = "STATE_UNSPECIFIED"
	StateQueued           SessionState = "QUEUED"
	StatePlanning         SessionState = "PLANNING"
	StateAwaitingApproval SessionState = "AWAITING_PLAN_APPROVAL"
	StateAwaitingFeedback SessionState = "AWAITING_USER_FEEDBACK"
	StateInProgress       SessionState = "IN_PROGRESS"
	StatePaused           SessionState = "PAUSED"
	StateFailed           SessionState = "FAILED"
	StateCompleted        SessionState = "COMPLETED"
)

// GitHubBranch represents a GitHub branch.
type GitHubBranch struct {
	DisplayName string `json:"displayName,omitempty"`
}

// GitHubRepo represents a GitHub repository.
type GitHubRepo struct {
	Owner         string         `json:"owner,omitempty"`
	Repo          string         `json:"repo,omitempty"`
	IsPrivate     bool           `json:"isPrivate,omitempty"`
	DefaultBranch *GitHubBranch  `json:"defaultBranch,omitempty"`
	Branches      []GitHubBranch `json:"branches,omitempty"`
}

// Source represents an input source of data for a session.
type Source struct {
	Name       string      `json:"name,omitempty"`
	ID         string      `json:"id,omitempty"`
	GitHubRepo *GitHubRepo `json:"githubRepo,omitempty"`
}

// GitHubRepoContext represents context to use a GitHubRepo in a session.
type GitHubRepoContext struct {
	StartingBranch string `json:"startingBranch,omitempty"`
}

// SourceContext represents context for how to use a source in a session.
type SourceContext struct {
	Source            string             `json:"source,omitempty"`
	GitHubRepoContext *GitHubRepoContext `json:"githubRepoContext,omitempty"`
}

// PullRequest represents a pull request.
type PullRequest struct {
	URL         string `json:"url,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

// SessionOutput represents an output of a session.
type SessionOutput struct {
	PullRequest *PullRequest `json:"pullRequest,omitempty"`
}

// Session represents a contiguous amount of work within the same context.
type Session struct {
	Name                string          `json:"name,omitempty"`
	ID                  string          `json:"id,omitempty"`
	Prompt              string          `json:"prompt,omitempty"`
	SourceContext       *SourceContext  `json:"sourceContext,omitempty"`
	Title               string          `json:"title,omitempty"`
	RequirePlanApproval bool            `json:"requirePlanApproval,omitempty"`
	CreateTime          string          `json:"createTime,omitempty"`
	UpdateTime          string          `json:"updateTime,omitempty"`
	State               SessionState    `json:"state,omitempty"`
	URL                 string          `json:"url,omitempty"`
	Outputs             []SessionOutput `json:"outputs,omitempty"`
}

// PlanStep represents a step in a plan.
type PlanStep struct {
	ID          string `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Index       int    `json:"index,omitempty"`
}

// Plan represents a sequence of steps that the agent will take to complete the task.
type Plan struct {
	ID         string     `json:"id,omitempty"`
	Steps      []PlanStep `json:"steps,omitempty"`
	CreateTime string     `json:"createTime,omitempty"`
}

// GitPatch represents a patch in Git format.
type GitPatch struct {
	UnidiffPatch           string `json:"unidiffPatch,omitempty"`
	BaseCommitID           string `json:"baseCommitId,omitempty"`
	SuggestedCommitMessage string `json:"suggestedCommitMessage,omitempty"`
}

// ChangeSet represents a change set artifact.
type ChangeSet struct {
	Source   string    `json:"source,omitempty"`
	GitPatch *GitPatch `json:"gitPatch,omitempty"`
}

// Media represents a media artifact.
type Media struct {
	Data     string `json:"data,omitempty"`
	MimeType string `json:"mimeType,omitempty"`
}

// BashOutput represents a bash output artifact.
type BashOutput struct {
	Command  string `json:"command,omitempty"`
	Output   string `json:"output,omitempty"`
	ExitCode int    `json:"exitCode,omitempty"`
}

// Artifact represents a single unit of data produced by an activity step.
type Artifact struct {
	ChangeSet  *ChangeSet  `json:"changeSet,omitempty"`
	Media      *Media      `json:"media,omitempty"`
	BashOutput *BashOutput `json:"bashOutput,omitempty"`
}

// Activity represents a single unit of work within a session.
type Activity struct {
	Name             string            `json:"name,omitempty"`
	ID               string            `json:"id,omitempty"`
	Description      string            `json:"description,omitempty"`
	CreateTime       string            `json:"createTime,omitempty"`
	Originator       string            `json:"originator,omitempty"`
	Artifacts        []Artifact        `json:"artifacts,omitempty"`
	AgentMessaged    map[string]string `json:"agentMessaged,omitempty"`
	UserMessaged     map[string]string `json:"userMessaged,omitempty"`
	PlanGenerated    map[string]any    `json:"planGenerated,omitempty"`
	PlanApproved     map[string]string `json:"planApproved,omitempty"`
	ProgressUpdated  map[string]string `json:"progressUpdated,omitempty"`
	SessionCompleted map[string]any    `json:"sessionCompleted,omitempty"`
	SessionFailed    map[string]string `json:"sessionFailed,omitempty"`
}
