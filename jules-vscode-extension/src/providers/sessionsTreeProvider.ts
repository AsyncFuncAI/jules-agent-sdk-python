/**
 * Sessions Tree Provider - Display sessions in tree view
 */

import * as vscode from 'vscode';
import { SessionManager } from '../sessionManager';
import { JulesSession, SessionState } from '../types';

export class SessionsTreeProvider implements vscode.TreeDataProvider<SessionTreeItem> {
    private _onDidChangeTreeData = new vscode.EventEmitter<SessionTreeItem | undefined | void>();
    readonly onDidChangeTreeData = this._onDidChangeTreeData.event;

    constructor(private sessionManager: SessionManager) {}

    refresh(): void {
        this._onDidChangeTreeData.fire();
    }

    getTreeItem(element: SessionTreeItem): vscode.TreeItem {
        return element;
    }

    async getChildren(element?: SessionTreeItem): Promise<SessionTreeItem[]> {
        if (!element) {
            // Root level - show all sessions
            try {
                const sessions = await this.sessionManager.listSessions();
                return sessions.map(session => new SessionTreeItem(session));
            } catch (error) {
                vscode.window.showErrorMessage(`Failed to load sessions: ${error}`);
                return [];
            }
        }
        return [];
    }
}

class SessionTreeItem extends vscode.TreeItem {
    constructor(public readonly session: JulesSession) {
        super(session.title || session.id, vscode.TreeItemCollapsibleState.None);

        this.description = this.getStateDescription(session.state);
        this.tooltip = this.getTooltip();
        this.iconPath = this.getIcon();
        this.contextValue = 'session';

        // Add commands - clicking session will activate it for chat
        this.command = {
            command: 'jules.selectSession',
            title: 'Select Session for Chat',
            arguments: [session.id]
        };
    }

    private getStateDescription(state: SessionState): string {
        const stateMap: { [key in SessionState]: string } = {
            [SessionState.STATE_UNSPECIFIED]: 'Unknown',
            [SessionState.QUEUED]: 'Queued',
            [SessionState.PLANNING]: 'Planning',
            [SessionState.AWAITING_PLAN_APPROVAL]: 'Needs Approval',
            [SessionState.AWAITING_USER_FEEDBACK]: 'Needs Feedback',
            [SessionState.IN_PROGRESS]: 'In Progress',
            [SessionState.PAUSED]: 'Paused',
            [SessionState.FAILED]: 'Failed',
            [SessionState.COMPLETED]: 'Completed'
        };
        return stateMap[state] || 'Unknown';
    }

    private getTooltip(): string {
        const { session } = this;
        return `ID: ${session.id}\nState: ${session.state}\nCreated: ${new Date(session.createTime).toLocaleString()}\nUpdated: ${new Date(session.updateTime).toLocaleString()}\n\nPrompt: ${session.prompt}`;
    }

    private getIcon(): vscode.ThemeIcon {
        switch (this.session.state) {
            case SessionState.COMPLETED:
                return new vscode.ThemeIcon('check', new vscode.ThemeColor('charts.green'));
            case SessionState.FAILED:
                return new vscode.ThemeIcon('error', new vscode.ThemeColor('charts.red'));
            case SessionState.IN_PROGRESS:
                return new vscode.ThemeIcon('sync~spin', new vscode.ThemeColor('charts.blue'));
            case SessionState.AWAITING_PLAN_APPROVAL:
            case SessionState.AWAITING_USER_FEEDBACK:
                return new vscode.ThemeIcon('warning', new vscode.ThemeColor('charts.yellow'));
            default:
                return new vscode.ThemeIcon('circle-outline');
        }
    }
}
