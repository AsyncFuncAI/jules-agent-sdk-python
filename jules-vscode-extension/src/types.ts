/**
 * Type definitions for Jules Agent VSCode Extension
 */

export interface JulesSession {
    id: string;
    name: string;
    title: string;
    state: SessionState;
    url: string;
    prompt: string;
    createTime: string;
    updateTime: string;
    outputs?: SessionOutput[];
}

export enum SessionState {
    STATE_UNSPECIFIED = 'STATE_UNSPECIFIED',
    QUEUED = 'QUEUED',
    PLANNING = 'PLANNING',
    AWAITING_PLAN_APPROVAL = 'AWAITING_PLAN_APPROVAL',
    AWAITING_USER_FEEDBACK = 'AWAITING_USER_FEEDBACK',
    IN_PROGRESS = 'IN_PROGRESS',
    PAUSED = 'PAUSED',
    FAILED = 'FAILED',
    COMPLETED = 'COMPLETED'
}

export interface SessionOutput {
    pullRequest?: PullRequest;
}

export interface PullRequest {
    url: string;
    title: string;
    description: string;
}

export interface Activity {
    id: string;
    name: string;
    description: string;
    createTime: string;
    originator: string;
    artifacts: Artifact[];
    agentMessaged?: { agentMessage: string };
    userMessaged?: { userMessage: string };
    planGenerated?: any;
    planApproved?: any;
    progressUpdated?: { title: string; description: string };
    sessionCompleted?: any;
    sessionFailed?: any;
}

export interface Artifact {
    changeSet?: ChangeSet;
    bashOutput?: BashOutput;
}

export interface ChangeSet {
    gitPatch?: GitPatch;
}

export interface GitPatch {
    unidiffPatch?: string;
    suggestedCommitMessage?: string;
}

export interface BashOutput {
    command: string;
    output: string;
    exitCode: number;
}

export interface Source {
    id: string;
    name: string;
    githubRepo?: GitHubRepo;
}

export interface GitHubRepo {
    owner: string;
    repo: string;
    isPrivate: boolean;
    defaultBranch?: string;
}

export interface JulesConfig {
    apiKey: string;
    baseUrl?: string;
    pythonPath: string;
    timeout: number;
    autoApprovePlans: boolean;
    sessionTimeout: number;
    pollInterval: number;
}

export interface BridgeResponse<T = any> {
    success: boolean;
    error?: string;
    errorType?: string;
    message?: string;
    [key: string]: any;
}

export interface ChatMessage {
    role: 'user' | 'assistant' | 'system';
    content: string;
    timestamp: Date;
    sessionId?: string;
}
