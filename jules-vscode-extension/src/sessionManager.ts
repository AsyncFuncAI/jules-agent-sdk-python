/**
 * Session Manager - Manages Jules session lifecycle
 */

import * as vscode from 'vscode';
import { JulesClient } from './julesClient';
import { JulesSession, SessionState } from './types';

export class SessionManager {
    private activeSessions = new Map<string, JulesSession>();
    private sessionPollers = new Map<string, NodeJS.Timeout>();
    private readonly onSessionUpdateEmitter = new vscode.EventEmitter<JulesSession>();
    public readonly onSessionUpdate = this.onSessionUpdateEmitter.event;

    constructor(
        private client: JulesClient,
        private outputChannel: vscode.OutputChannel
    ) {}

    async createSession(
        prompt: string,
        source: string,
        options: {
            startingBranch?: string;
            title?: string;
            requirePlanApproval?: boolean;
        } = {}
    ): Promise<JulesSession> {
        try {
            this.outputChannel.appendLine(`Creating session: ${prompt.substring(0, 50)}...`);

            const session = await this.client.createSession(prompt, source, options);
            this.activeSessions.set(session.id, session);
            this.onSessionUpdateEmitter.fire(session);

            // Start polling for updates
            this.startPolling(session.id);

            vscode.window.showInformationMessage(
                `Jules session created: ${session.title || session.id}`,
                'View Session'
            ).then(selection => {
                if (selection === 'View Session' && session.url) {
                    vscode.env.openExternal(vscode.Uri.parse(session.url));
                }
            });

            return session;
        } catch (error) {
            const errorMsg = error instanceof Error ? error.message : String(error);
            this.outputChannel.appendLine(`Error creating session: ${errorMsg}`);
            vscode.window.showErrorMessage(`Failed to create Jules session: ${errorMsg}`);
            throw error;
        }
    }

    async getSession(sessionId: string, forceRefresh = false): Promise<JulesSession | undefined> {
        if (!forceRefresh && this.activeSessions.has(sessionId)) {
            return this.activeSessions.get(sessionId);
        }

        try {
            const session = await this.client.getSession(sessionId);
            this.activeSessions.set(session.id, session);
            this.onSessionUpdateEmitter.fire(session);
            return session;
        } catch (error) {
            this.outputChannel.appendLine(`Error getting session ${sessionId}: ${error}`);
            return undefined;
        }
    }

    async listSessions(refresh = false): Promise<JulesSession[]> {
        if (!refresh && this.activeSessions.size > 0) {
            return Array.from(this.activeSessions.values());
        }

        try {
            const sessions = await this.client.listSessions();
            this.activeSessions.clear();
            sessions.forEach(session => {
                this.activeSessions.set(session.id, session);
            });
            return sessions;
        } catch (error) {
            this.outputChannel.appendLine(`Error listing sessions: ${error}`);
            throw error;
        }
    }

    async approvePlan(sessionId: string): Promise<void> {
        try {
            await this.client.approvePlan(sessionId);
            vscode.window.showInformationMessage(`Plan approved for session ${sessionId}`);

            // Refresh session state
            await this.getSession(sessionId, true);
        } catch (error) {
            const errorMsg = error instanceof Error ? error.message : String(error);
            vscode.window.showErrorMessage(`Failed to approve plan: ${errorMsg}`);
            throw error;
        }
    }

    async sendMessage(sessionId: string, message: string): Promise<void> {
        try {
            await this.client.sendMessage(sessionId, message);
            this.outputChannel.appendLine(`Message sent to session ${sessionId}`);

            // Refresh session state
            await this.getSession(sessionId, true);
        } catch (error) {
            const errorMsg = error instanceof Error ? error.message : String(error);
            vscode.window.showErrorMessage(`Failed to send message: ${errorMsg}`);
            throw error;
        }
    }

    async stopSession(sessionId: string): Promise<void> {
        this.stopPolling(sessionId);
        this.activeSessions.delete(sessionId);
        vscode.window.showInformationMessage(`Stopped monitoring session ${sessionId}`);
    }

    private startPolling(sessionId: string): void {
        // Don't poll if already polling
        if (this.sessionPollers.has(sessionId)) {
            return;
        }

        const config = vscode.workspace.getConfiguration('jules');
        const pollInterval = config.get<number>('pollInterval', 5) * 1000;

        const poller = setInterval(async () => {
            const session = await this.getSession(sessionId, true);

            if (!session) {
                this.stopPolling(sessionId);
                return;
            }

            // Stop polling if session is in terminal state
            if (this.isTerminalState(session.state)) {
                this.stopPolling(sessionId);
                this.handleSessionCompletion(session);
            }
        }, pollInterval);

        this.sessionPollers.set(sessionId, poller);
    }

    private stopPolling(sessionId: string): void {
        const poller = this.sessionPollers.get(sessionId);
        if (poller) {
            clearInterval(poller);
            this.sessionPollers.delete(sessionId);
        }
    }

    private isTerminalState(state: SessionState): boolean {
        return state === SessionState.COMPLETED || state === SessionState.FAILED;
    }

    private handleSessionCompletion(session: JulesSession): void {
        if (session.state === SessionState.COMPLETED) {
            vscode.window.showInformationMessage(
                `Jules session completed: ${session.title || session.id}`,
                'View Session',
                'View PR'
            ).then(selection => {
                if (selection === 'View Session' && session.url) {
                    vscode.env.openExternal(vscode.Uri.parse(session.url));
                } else if (selection === 'View PR' && session.outputs?.[0]?.pullRequest?.url) {
                    vscode.env.openExternal(vscode.Uri.parse(session.outputs[0].pullRequest.url));
                }
            });
        } else if (session.state === SessionState.FAILED) {
            vscode.window.showErrorMessage(
                `Jules session failed: ${session.title || session.id}`,
                'View Session'
            ).then(selection => {
                if (selection === 'View Session' && session.url) {
                    vscode.env.openExternal(vscode.Uri.parse(session.url));
                }
            });
        }
    }

    getActiveSession(): JulesSession | undefined {
        // Return the most recently updated session
        const sessions = Array.from(this.activeSessions.values());
        if (sessions.length === 0) {
            return undefined;
        }

        return sessions.sort((a, b) => {
            return new Date(b.updateTime).getTime() - new Date(a.updateTime).getTime();
        })[0];
    }

    dispose(): void {
        // Stop all pollers
        this.sessionPollers.forEach((poller, sessionId) => {
            this.stopPolling(sessionId);
        });
        this.activeSessions.clear();
        this.onSessionUpdateEmitter.dispose();
    }
}
