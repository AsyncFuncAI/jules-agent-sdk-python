"""Data models for Jules API resources."""

from dataclasses import dataclass, field
from typing import Optional, List, Dict, Any
from enum import Enum


class SessionState(str, Enum):
    """Session state enumeration."""

    STATE_UNSPECIFIED = "STATE_UNSPECIFIED"
    QUEUED = "QUEUED"
    PLANNING = "PLANNING"
    AWAITING_PLAN_APPROVAL = "AWAITING_PLAN_APPROVAL"
    AWAITING_USER_FEEDBACK = "AWAITING_USER_FEEDBACK"
    IN_PROGRESS = "IN_PROGRESS"
    PAUSED = "PAUSED"
    FAILED = "FAILED"
    COMPLETED = "COMPLETED"


@dataclass
class GitHubBranch:
    """A GitHub branch."""

    display_name: str

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "GitHubBranch":
        """Create from API response dictionary."""
        return cls(display_name=data.get("displayName", ""))

    def to_dict(self) -> Dict[str, Any]:
        """Convert to API request dictionary."""
        return {"displayName": self.display_name}


@dataclass
class GitHubRepo:
    """A GitHub repository."""

    owner: str
    repo: str
    is_private: bool = False
    default_branch: Optional[GitHubBranch] = None
    branches: List[GitHubBranch] = field(default_factory=list)

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "GitHubRepo":
        """Create from API response dictionary."""
        default_branch = None
        if data.get("defaultBranch"):
            default_branch = GitHubBranch.from_dict(data["defaultBranch"])

        branches = []
        if data.get("branches"):
            branches = [GitHubBranch.from_dict(b) for b in data["branches"]]

        return cls(
            owner=data.get("owner", ""),
            repo=data.get("repo", ""),
            is_private=data.get("isPrivate", False),
            default_branch=default_branch,
            branches=branches,
        )

    def to_dict(self) -> Dict[str, Any]:
        """Convert to API request dictionary."""
        result: Dict[str, Any] = {
            "owner": self.owner,
            "repo": self.repo,
            "isPrivate": self.is_private,
        }
        if self.default_branch:
            result["defaultBranch"] = self.default_branch.to_dict()
        if self.branches:
            result["branches"] = [b.to_dict() for b in self.branches]
        return result


@dataclass
class Source:
    """An input source of data for a session."""

    name: str
    id: str
    github_repo: Optional[GitHubRepo] = None

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "Source":
        """Create from API response dictionary."""
        github_repo = None
        if data.get("githubRepo"):
            github_repo = GitHubRepo.from_dict(data["githubRepo"])

        return cls(
            name=data.get("name", ""),
            id=data.get("id", ""),
            github_repo=github_repo,
        )

    def to_dict(self) -> Dict[str, Any]:
        """Convert to API request dictionary."""
        result: Dict[str, Any] = {"name": self.name, "id": self.id}
        if self.github_repo:
            result["githubRepo"] = self.github_repo.to_dict()
        return result


@dataclass
class GitHubRepoContext:
    """Context to use a GitHubRepo in a session."""

    starting_branch: str

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "GitHubRepoContext":
        """Create from API response dictionary."""
        return cls(starting_branch=data.get("startingBranch", ""))

    def to_dict(self) -> Dict[str, Any]:
        """Convert to API request dictionary."""
        return {"startingBranch": self.starting_branch}


@dataclass
class SourceContext:
    """Context for how to use a source in a session."""

    source: str
    github_repo_context: Optional[GitHubRepoContext] = None

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "SourceContext":
        """Create from API response dictionary."""
        github_context = None
        if data.get("githubRepoContext"):
            github_context = GitHubRepoContext.from_dict(data["githubRepoContext"])

        return cls(source=data.get("source", ""), github_repo_context=github_context)

    def to_dict(self) -> Dict[str, Any]:
        """Convert to API request dictionary."""
        result: Dict[str, Any] = {"source": self.source}
        if self.github_repo_context:
            result["githubRepoContext"] = self.github_repo_context.to_dict()
        return result


@dataclass
class PullRequest:
    """A pull request."""

    url: str
    title: str
    description: str

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "PullRequest":
        """Create from API response dictionary."""
        return cls(
            url=data.get("url", ""),
            title=data.get("title", ""),
            description=data.get("description", ""),
        )

    def to_dict(self) -> Dict[str, Any]:
        """Convert to API request dictionary."""
        return {"url": self.url, "title": self.title, "description": self.description}


@dataclass
class SessionOutput:
    """An output of a session."""

    pull_request: Optional[PullRequest] = None

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "SessionOutput":
        """Create from API response dictionary."""
        pr = None
        if data.get("pullRequest"):
            pr = PullRequest.from_dict(data["pullRequest"])
        return cls(pull_request=pr)

    def to_dict(self) -> Dict[str, Any]:
        """Convert to API request dictionary."""
        result: Dict[str, Any] = {}
        if self.pull_request:
            result["pullRequest"] = self.pull_request.to_dict()
        return result


@dataclass
class Session:
    """A session is a contiguous amount of work within the same context."""

    prompt: str
    source_context: SourceContext
    name: str = ""
    id: str = ""
    title: str = ""
    automation_mode: Optional[str] = None
    require_plan_approval: bool = False
    create_time: str = ""
    update_time: str = ""
    state: SessionState = SessionState.STATE_UNSPECIFIED
    url: str = ""
    outputs: List[SessionOutput] = field(default_factory=list)

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "Session":
        """Create from API response dictionary."""
        source_context = SourceContext.from_dict(data.get("sourceContext", {}))
        outputs = []
        if data.get("outputs"):
            outputs = [SessionOutput.from_dict(o) for o in data["outputs"]]

        state = SessionState.STATE_UNSPECIFIED
        if data.get("state"):
            try:
                state = SessionState(data["state"])
            except ValueError:
                state = SessionState.STATE_UNSPECIFIED

        return cls(
            name=data.get("name", ""),
            id=data.get("id", ""),
            prompt=data.get("prompt", ""),
            source_context=source_context,
            title=data.get("title", ""),
            automation_mode=data.get("automationMode"),
            require_plan_approval=data.get("requirePlanApproval", False),
            create_time=data.get("createTime", ""),
            update_time=data.get("updateTime", ""),
            state=state,
            url=data.get("url", ""),
            outputs=outputs,
        )

    def to_dict(self) -> Dict[str, Any]:
        """Convert to API request dictionary."""
        result: Dict[str, Any] = {
            "prompt": self.prompt,
            "sourceContext": self.source_context.to_dict(),
        }
        if self.name:
            result["name"] = self.name
        if self.title:
            result["title"] = self.title
        if self.automation_mode:
            result["automationMode"] = self.automation_mode
        if self.require_plan_approval:
            result["requirePlanApproval"] = self.require_plan_approval
        if self.outputs:
            result["outputs"] = [o.to_dict() for o in self.outputs]
        return result


@dataclass
class PlanStep:
    """A step in a plan."""

    id: str
    title: str
    description: str
    index: int

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "PlanStep":
        """Create from API response dictionary."""
        return cls(
            id=data.get("id", ""),
            title=data.get("title", ""),
            description=data.get("description", ""),
            index=data.get("index", 0),
        )

    def to_dict(self) -> Dict[str, Any]:
        """Convert to API request dictionary."""
        return {
            "id": self.id,
            "title": self.title,
            "description": self.description,
            "index": self.index,
        }


@dataclass
class Plan:
    """A sequence of steps that the agent will take to complete the task."""

    id: str
    steps: List[PlanStep]
    create_time: str = ""

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "Plan":
        """Create from API response dictionary."""
        steps = []
        if data.get("steps"):
            steps = [PlanStep.from_dict(s) for s in data["steps"]]

        return cls(
            id=data.get("id", ""), steps=steps, create_time=data.get("createTime", "")
        )

    def to_dict(self) -> Dict[str, Any]:
        """Convert to API request dictionary."""
        return {
            "id": self.id,
            "steps": [s.to_dict() for s in self.steps],
            "createTime": self.create_time,
        }


@dataclass
class GitPatch:
    """A patch in Git format."""

    unidiff_patch: str
    base_commit_id: str
    suggested_commit_message: str

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "GitPatch":
        """Create from API response dictionary."""
        return cls(
            unidiff_patch=data.get("unidiffPatch", ""),
            base_commit_id=data.get("baseCommitId", ""),
            suggested_commit_message=data.get("suggestedCommitMessage", ""),
        )

    def to_dict(self) -> Dict[str, Any]:
        """Convert to API request dictionary."""
        return {
            "unidiffPatch": self.unidiff_patch,
            "baseCommitId": self.base_commit_id,
            "suggestedCommitMessage": self.suggested_commit_message,
        }


@dataclass
class ChangeSet:
    """A change set artifact."""

    source: str
    git_patch: Optional[GitPatch] = None

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "ChangeSet":
        """Create from API response dictionary."""
        git_patch = None
        if data.get("gitPatch"):
            git_patch = GitPatch.from_dict(data["gitPatch"])
        return cls(source=data.get("source", ""), git_patch=git_patch)

    def to_dict(self) -> Dict[str, Any]:
        """Convert to API request dictionary."""
        result: Dict[str, Any] = {"source": self.source}
        if self.git_patch:
            result["gitPatch"] = self.git_patch.to_dict()
        return result


@dataclass
class Media:
    """A media artifact."""

    data: str
    mime_type: str

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "Media":
        """Create from API response dictionary."""
        return cls(data=data.get("data", ""), mime_type=data.get("mimeType", ""))

    def to_dict(self) -> Dict[str, Any]:
        """Convert to API request dictionary."""
        return {"data": self.data, "mimeType": self.mime_type}


@dataclass
class BashOutput:
    """A bash output artifact."""

    command: str
    output: str
    exit_code: int

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "BashOutput":
        """Create from API response dictionary."""
        return cls(
            command=data.get("command", ""),
            output=data.get("output", ""),
            exit_code=data.get("exitCode", 0),
        )

    def to_dict(self) -> Dict[str, Any]:
        """Convert to API request dictionary."""
        return {"command": self.command, "output": self.output, "exitCode": self.exit_code}


@dataclass
class Artifact:
    """An artifact is a single unit of data produced by an activity step."""

    change_set: Optional[ChangeSet] = None
    media: Optional[Media] = None
    bash_output: Optional[BashOutput] = None

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "Artifact":
        """Create from API response dictionary."""
        change_set = None
        media = None
        bash_output = None

        if data.get("changeSet"):
            change_set = ChangeSet.from_dict(data["changeSet"])
        if data.get("media"):
            media = Media.from_dict(data["media"])
        if data.get("bashOutput"):
            bash_output = BashOutput.from_dict(data["bashOutput"])

        return cls(change_set=change_set, media=media, bash_output=bash_output)

    def to_dict(self) -> Dict[str, Any]:
        """Convert to API request dictionary."""
        result: Dict[str, Any] = {}
        if self.change_set:
            result["changeSet"] = self.change_set.to_dict()
        if self.media:
            result["media"] = self.media.to_dict()
        if self.bash_output:
            result["bashOutput"] = self.bash_output.to_dict()
        return result


@dataclass
class Activity:
    """An Activity is a single unit of work within a session."""

    name: str
    id: str = ""
    description: str = ""
    create_time: str = ""
    originator: str = ""
    artifacts: List[Artifact] = field(default_factory=list)
    agent_messaged: Optional[Dict[str, str]] = None
    user_messaged: Optional[Dict[str, str]] = None
    plan_generated: Optional[Dict[str, Any]] = None
    plan_approved: Optional[Dict[str, str]] = None
    progress_updated: Optional[Dict[str, str]] = None
    session_completed: Optional[Dict[str, Any]] = None
    session_failed: Optional[Dict[str, str]] = None

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "Activity":
        """Create from API response dictionary."""
        artifacts = []
        if data.get("artifacts"):
            artifacts = [Artifact.from_dict(a) for a in data["artifacts"]]

        return cls(
            name=data.get("name", ""),
            id=data.get("id", ""),
            description=data.get("description", ""),
            create_time=data.get("createTime", ""),
            originator=data.get("originator", ""),
            artifacts=artifacts,
            agent_messaged=data.get("agentMessaged"),
            user_messaged=data.get("userMessaged"),
            plan_generated=data.get("planGenerated"),
            plan_approved=data.get("planApproved"),
            progress_updated=data.get("progressUpdated"),
            session_completed=data.get("sessionCompleted"),
            session_failed=data.get("sessionFailed"),
        )

    def to_dict(self) -> Dict[str, Any]:
        """Convert to API request dictionary."""
        result: Dict[str, Any] = {"name": self.name}
        if self.id:
            result["id"] = self.id
        if self.description:
            result["description"] = self.description
        if self.artifacts:
            result["artifacts"] = [a.to_dict() for a in self.artifacts]
        return result
