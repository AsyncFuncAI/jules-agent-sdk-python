package jules

import "time"

type SessionState string

const (
	StateUnspecified      SessionState = "STATE_UNSPECIFIED"
	Queued                SessionState = "QUEUED"
	Planning              SessionState = "PLANNING"
	AwaitingPlanApproval  SessionState = "AWAITING_PLAN_APPROVAL"
	AwaitingUserFeedback  SessionState = "AWAITING_USER_FEEDBACK"
	InProgress            SessionState = "IN_PROGRESS"
	Paused                SessionState = "PAUSED"
	Failed                SessionState = "FAILED"
	Completed             SessionState = "COMPLETED"
)

type GitHubBranch struct {
	DisplayName string `json:"displayName"`
}

type GitHubRepo struct {
	Owner         string         `json:"owner"`
	Repo          string         `json:"repo"`
	IsPrivate     bool           `json:"isPrivate"`
	DefaultBranch *GitHubBranch  `json:"defaultBranch,omitempty"`
	Branches      []GitHubBranch `json:"branches,omitempty"`
}

type Source struct {
	Name       string      `json:"name"`
	ID         string      `json:"id"`
	GithubRepo *GitHubRepo `json:"githubRepo,omitempty"`
}

type GitHubRepoContext struct {
	StartingBranch string `json:"startingBranch"`
}

type SourceContext struct {
	Source            string             `json:"source"`
	GithubRepoContext *GitHubRepoContext `json:"githubRepoContext,omitempty"`
}

type PullRequest struct {
	URL         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type SessionOutput struct {
	PullRequest *PullRequest `json:"pullRequest,omitempty"`
}

type Session struct {
	Prompt               string          `json:"prompt"`
	SourceContext        SourceContext   `json:"sourceContext"`
	Name                 string          `json:"name,omitempty"`
	ID                   string          `json:"id,omitempty"`
	Title                string          `json:"title,omitempty"`
	RequirePlanApproval  bool            `json:"requirePlanApproval,omitempty"`
	CreateTime           time.Time       `json:"createTime,omitempty"`
	UpdateTime           time.Time       `json:"updateTime,omitempty"`
	State                SessionState    `json:"state,omitempty"`
	URL                  string          `json:"url,omitempty"`
	Outputs              []SessionOutput `json:"outputs,omitempty"`
}

type PlanStep struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Index       int    `json:"index"`
}

type Plan struct {
	ID        string     `json:"id"`
	Steps     []PlanStep `json:"steps"`
	CreateTime time.Time  `json:"createTime,omitempty"`
}

type GitPatch struct {
	UnidiffPatch           string `json:"unidiffPatch"`
	BaseCommitID           string `json:"baseCommitId"`
	SuggestedCommitMessage string `json:"suggestedCommitMessage"`
}

type ChangeSet struct {
	Source   string    `json:"source"`
	GitPatch *GitPatch `json:"gitPatch,omitempty"`
}

type Media struct {
	Data     string `json:"data"`
	MimeType string `json:"mimeType"`
}

type BashOutput struct {
	Command  string `json:"command"`
	Output   string `json:"output"`
	ExitCode int    `json:"exitCode"`
}

type Artifact struct {
	ChangeSet  *ChangeSet  `json:"changeSet,omitempty"`
	Media      *Media      `json:"media,omitempty"`
	BashOutput *BashOutput `json:"bashOutput,omitempty"`
}

type Activity struct {
	Name              string                 `json:"name"`
	ID                string                 `json:"id,omitempty"`
	Description       string                 `json:"description,omitempty"`
	CreateTime        time.Time              `json:"createTime,omitempty"`
	Originator        string                 `json:"originator,omitempty"`
	Artifacts         []Artifact             `json:"artifacts,omitempty"`
	AgentMessaged     map[string]string      `json:"agentMessaged,omitempty"`
	UserMessaged      map[string]string      `json:"userMessaged,omitempty"`
	PlanGenerated     map[string]interface{} `json:"planGenerated,omitempty"`
	PlanApproved      map[string]string      `json:"planApproved,omitempty"`
	ProgressUpdated   map[string]string      `json:"progressUpdated,omitempty"`
	SessionCompleted  map[string]interface{} `json:"sessionCompleted,omitempty"`
	SessionFailed     map[string]string      `json:"sessionFailed,omitempty"`
}