# Jules API Open Source SDK Context: Session, Activity, and Source API Definitions

This document unifies the resource definitions, data structures, and available methods for the **Session**, **Activity**, and **Source** APIs, providing a comprehensive context for building an open-source SDK.

## 0. Authentication
To authenticate your requests, pass the API key in the X-Goog-Api-Key header of your API calls.
e.g: 
curl 'https://jules.googleapis.com/v1alpha/sources' \
    -H 'X-Goog-Api-Key: YOUR_API_KEY'
    
## 1. Resource: Session

A session is a contiguous amount of work within the same context.

### 1.1 Session JSON Representation

```json
{
  "name": string,
  "id": string,
  "prompt": string,
  "sourceContext": {
    object (SourceContext)
  },
  "title": string,
  "requirePlanApproval": boolean,
  "createTime": string,
  "updateTime": string,
  "state": enum (State),
  "url": string,
  "outputs": [
    {
      object (SessionOutput)
    }
  ]
}
```

### 1.2 Session Fields

| Field | Type | Description | Attributes |
| :--- | :--- | :--- | :--- |
| `name` | `string` | The full resource name (e.g., `"sessions/{session}"`). | Output only. Identifier. |
| `id` | `string` | The id of the session. | Output only. |
| `prompt` | `string` | The prompt to start the session with. | Required. |
| `sourceContext` | `object (SourceContext)` | The source to use in this session, with additional context. | Required. |
| `title` | `string` | Optional title. If not provided, the system will generate one. | Optional. |
| `requirePlanApproval` | `boolean` | If `true`, plans require explicit plan approval. If not set, plans are auto-approved. | Optional. Input only. |
| `createTime` | `string` (Timestamp) | The time the session was created (RFC 3339). | Output only. |
| `updateTime` | `string` (Timestamp) | The time the session was last updated (RFC 3339). | Output only. |
| `state` | `enum (State)` | The state of the session. | Output only. |
| `url` | `string` | The URL to view the session in the Jules web app. | Output only. |
| `outputs[]` | `object (SessionOutput)` | The outputs of the session, if any. | Output only. |

### 1.3 Data Structure: SourceContext

Context for how to use a source in a session.

```json
{
  "source": string,
  "githubRepoContext": {
    object (GitHubRepoContext)
  }
}
```

| Field | Type | Description | Attributes |
| :--- | :--- | :--- | :--- |
| `source` | `string` | The name of the source (e.g., `"sources/{source}"`). | Required. |
| `githubRepoContext` | `object (GitHubRepoContext)` | Context to use a GitHubRepo in a session. | Union field `context`. |

### 1.4 Data Structure: GitHubRepoContext

Context to use a GitHubRepo in a session.

```json
{
  "startingBranch": string
}
```

| Field | Type | Description | Attributes |
| :--- | :--- | :--- | :--- |
| `startingBranch` | `string` | The name of the branch to start the session from. | Required. |

### 1.5 Enum: State (Session State)

| Enum Value | Description |
| :--- | :--- |
| `STATE_UNSPECIFIED` | The state is unspecified. |
| `QUEUED` | The session is queued. |
| `PLANNING` | The agent is planning. |
| `AWAITING_PLAN_APPROVAL` | The agent is waiting for plan approval. |
| `AWAITING_USER_FEEDBACK` | The agent is waiting for user feedback. |
| `IN_PROGRESS` | The session is in progress. |
| `PAUSED` | The session is paused. |
| `FAILED` | The session has failed. |
| `COMPLETED` | The session has completed. |

### 1.6 Data Structure: SessionOutput

An output of a session.

```json
{
  "pullRequest": {
    object (PullRequest)
  }
}
```

| Field | Type | Description | Attributes |
| :--- | :--- | :--- | :--- |
| `pullRequest` | `object (PullRequest)` | A pull request created by the session. | Union field `output`. |

### 1.7 Data Structure: PullRequest

A pull request.

```json
{
  "url": string,
  "title": string,
  "description": string
}
```

| Field | Type | Description |
| :--- | :--- | :--- |
| `url` | `string` | The URL of the pull request. |
| `title` | `string` | The title of the pull request. |
| `description` | `string` | The description of the pull request. |

### 1.8 Session Methods

| Method | HTTP Request | Description | Parameters/Body | Response |
| :--- | :--- | :--- | :--- | :--- |
| **`sessions.create`** | `POST https://jules.googleapis.com/v1alpha/sessions` | Creates a new session. | Request body: `object (Session)` | Response body: `object (Session)` |
| **`sessions.get`** | `GET https://jules.googleapis.com/v1alpha/{name=sessions/*}` | Gets a single session. | Path param: `name` (e.g., `"sessions/{session}"`) | Response body: `object (Session)` |
| **`sessions.list`** | `GET https://jules.googleapis.com/v1alpha/sessions` | Lists all sessions. | Query params: `pageSize` (int), `pageToken` (string) | Response body: `{"sessions": [Session], "nextPageToken": string}` |
| **`sessions.approvePlan`** | `POST https://jules.googleapis.com/v1alpha/{session=sessions/*}:approvePlan` | Approves a plan in a session. | Path param: `session` (e.g., `"sessions/{session}"`). Request body: empty. | Response body: empty. |
| **`sessions.sendMessage`** | `POST https://jules.googleapis.com/v1alpha/{session=sessions/*}:sendMessage` | Sends a message from the user to a session. | Path param: `session`. Request body: `{"prompt": string}` | Response body: empty. |

***

## 2. REST Resource: sessions.activities

An **Activity** is a single unit of work within a session.

### 2.1 Activity JSON Representation

```json
{
  "name": string,
  "id": string,
  "description": string,
  "createTime": string,
  "originator": string,
  "artifacts": [
    {
      object (Artifact)
    }
  ],
  "agentMessaged": { object (AgentMessaged) },
  "userMessaged": { object (UserMessaged) },
  "planGenerated": { object (PlanGenerated) },
  "planApproved": { object (PlanApproved) },
  "progressUpdated": { object (ProgressUpdated) },
  "sessionCompleted": { object (SessionCompleted) },
  "sessionFailed": { object (SessionFailed) }
}
```

### 2.2 Activity Fields

| Field | Type | Description | Attributes |
| :--- | :--- | :--- | :--- |
| `name` | `string` | The full resource name (e.g., `"sessions/{session}/activities/{activity}"`). | Identifier. |
| `id` | `string` | The id of the activity. | Output only. |
| `description` | `string` | A description of this activity. | Output only. |
| `createTime` | `string` (Timestamp) | The time at which this activity was created (RFC 3339). | Output only. |
| `originator` | `string` | The entity that this activity originated from (e.g. `"user"`, `"agent"`, `"system"`). | |
| `artifacts[]` | `object (Artifact)` | The artifacts produced by this activity. | Output only. |
| **Union Field `activity`** | | **The activity content (one of the following):** | |
| `agentMessaged` | `object (AgentMessaged)` | The agent posted a message. | |
| `userMessaged` | `object (UserMessaged)` | The user posted a message. | |
| `planGenerated` | `object (PlanGenerated)` | A plan was generated. | |
| `planApproved` | `object (PlanApproved)` | A plan was approved. | |
| `progressUpdated` | `object (ProgressUpdated)` | There was a progress update. | |
| `sessionCompleted` | `object (SessionCompleted)` | The session was completed. | |
| `sessionFailed` | `object (SessionFailed)` | The session failed. | |

### 2.3 Activity Content Structures

| Structure | JSON Representation | Fields |
| :--- | :--- | :--- |
| **AgentMessaged** | `{"agentMessage": string}` | `agentMessage`: The message the agent posted. |
| **UserMessaged** | `{"userMessage": string}` | `userMessage`: The message the user posted. |
| **PlanGenerated** | `{"plan": object (Plan)}` | `plan`: The plan that was generated. |
| **PlanApproved** | `{"planId": string}` | `planId`: The ID of the plan that was approved. |
| **ProgressUpdated** | `{"title": string, "description": string}` | `title`: Title of the progress update. `description`: Description of the progress update. |
| **SessionCompleted** | `{}` | No fields. |
| **SessionFailed** | `{"reason": string}` | `reason`: The reason the session failed. |

### 2.4 Data Structure: Plan

A sequence of steps that the agent will take to complete the task.

```json
{
  "id": string,
  "steps": [
    {
      object (PlanStep)
    }
  ],
  "createTime": string
}
```

| Field | Type | Description | Attributes |
| :--- | :--- | :--- | :--- |
| `id` | `string` | ID for this plan; unique within a session. | Output only. |
| `steps[]` | `object (PlanStep)` | The steps in the plan. | Output only. |
| `createTime` | `string` (Timestamp) | Time when the plan was created (RFC 3339). | Output only. |

### 2.5 Data Structure: PlanStep

A step in a plan.

```json
{
  "id": string,
  "title": string,
  "description": string,
  "index": integer
}
```

| Field | Type | Description | Attributes |
| :--- | :--- | :--- | :--- |
| `id` | `string` | ID for this step; unique within a plan. | Output only. |
| `title` | `string` | The title of the step. | Output only. |
| `description` | `string` | The description of the step. | Output only. |
| `index` | `integer` | 0-based index into the plan.steps. | Output only. |

### 2.6 Data Structure: Artifact

An artifact is a single unit of data produced by an activity step.

```json
{
  "changeSet": { object (ChangeSet) },
  "media": { object (Media) },
  "bashOutput": { object (BashOutput) }
}
```

| Field | Type | Description | Attributes |
| :--- | :--- | :--- | :--- |
| **Union Field `content`** | | **The artifact content (one of the following):** | |
| `changeSet` | `object (ChangeSet)` | A change set was produced (e.g. code changes). | |
| `media` | `object (Media)` | A media file was produced (e.g. image, video). | |
| `bashOutput` | `object (BashOutput)` | A bash output was produced. | |

### 2.7 Artifact Content Structures

| Structure | JSON Representation | Fields |
| :--- | :--- | :--- |
| **ChangeSet** | `{"source": string, "gitPatch": object (GitPatch)}` | `source`: Name of the source (e.g., `"sources/{source}"`). `gitPatch`: A patch in Git format (Union field `changes`). |
| **Media** | `{"data": string, "mimeType": string}` | `data`: The media data (base64-encoded string). `mimeType`: The media mime type. |
| **BashOutput** | `{"command": string, "output": string, "exitCode": integer}` | `command`: The bash command. `output`: The bash output (stdout and stderr). `exitCode`: The bash exit code. |

### 2.8 Data Structure: GitPatch

A patch in Git format.

```json
{
  "unidiffPatch": string,
  "baseCommitId": string,
  "suggestedCommitMessage": string
}
```

| Field | Type | Description |
| :--- | :--- | :--- |
| `unidiffPatch` | `string` | The patch in unidiff format. |
| `baseCommitId` | `string` | The base commit id of the patch. |
| `suggestedCommitMessage` | `string` | A suggested commit message for the patch. |

### 2.9 Activity Methods

| Method | HTTP Request | Description | Parameters/Body | Response |
| :--- | :--- | :--- | :--- | :--- |
| **`sessions.activities.get`** | `GET https://jules.googleapis.com/v1alpha/{name=sessions/*/activities/*}` | Gets a single activity. | Path param: `name` (e.g., `"sessions/{session}/activities/{activity}"`) | Response body: `object (Activity)` |
| **`sessions.activities.list`** | `GET https://jules.googleapis.com/v1alpha/{parent=sessions/*}/activities` | Lists activities for a session. | Path param: `parent` (e.g., `"sessions/{session}"`). Query params: `pageSize` (int), `pageToken` (string) | Response body: `{"activities": [Activity], "nextPageToken": string}` |

***

## 3. REST Resource: sources

An input source of data for a session.

### 3.1 Source JSON Representation

```json
{
  "name": string,
  "id": string,
  "githubRepo": {
    object (GitHubRepo)
  }
}
```

### 3.2 Source Fields

| Field | Type | Description | Attributes |
| :--- | :--- | :--- | :--- |
| `name` | `string` | The full resource name (e.g., `"sources/{source}"`). | Identifier. |
| `id` | `string` | The id of the source. | Output only. |
| `githubRepo` | `object (GitHubRepo)` | A GitHub repo. | Union field `source`. |

### 3.3 Data Structure: GitHubRepo

A GitHub repo.

```json
{
  "owner": string,
  "repo": string,
  "isPrivate": boolean,
  "defaultBranch": {
    object (GitHubBranch)
  },
  "branches": [
    {
      object (GitHubBranch)
    }
  ]
}
```

| Field | Type | Description |
| :--- | :--- | :--- |
| `owner` | `string` | The owner of the repo. |
| `repo` | `string` | The name of the repo. |
| `isPrivate` | `boolean` | Whether this repo is private. |
| `defaultBranch` | `object (GitHubBranch)` | The default branch for this repo. |
| `branches[]` | `object (GitHubBranch)` | The list of active branches for this repo. |

### 3.4 Data Structure: GitHubBranch

A GitHub branch.

```json
{
  "displayName": string
}
```

| Field | Type | Description |
| :--- | :--- | :--- |
| `displayName` | `string` | The name of the GitHub branch. |

### 3.5 Source Methods

| Method | HTTP Request | Description | Parameters/Body | Response |
| :--- | :--- | :--- | :--- | :--- |
| **`sources.get`** | `GET https://jules.googleapis.com/v1alpha/{name=sources/**}` | Gets a single source. | Path param: `name` (e.g., `"sources/{source}"`) | Response body: `object (Source)` |
| **`sources.list`** | `GET https://jules.googleapis.com/v1alpha/sources` | Lists sources. | Query params: `filter` (string), `pageSize` (int), `pageToken` (string) | Response body: `{"sources": [Source], "nextPageToken": string}` |