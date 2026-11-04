/**
 * Activities Tree Provider - Display session activities
 */

import * as vscode from 'vscode';
import { JulesClient } from '../julesClient';
import { SessionManager } from '../sessionManager';
import { Activity } from '../types';

type TreeItemType = ActivityTreeItem | DetailTreeItem | InfoTreeItem;

export class ActivitiesTreeProvider implements vscode.TreeDataProvider<TreeItemType> {
    private _onDidChangeTreeData = new vscode.EventEmitter<TreeItemType | undefined | void>();
    readonly onDidChangeTreeData = this._onDidChangeTreeData.event;

    constructor(
        private client: JulesClient,
        private sessionManager: SessionManager
    ) {}

    refresh(): void {
        this._onDidChangeTreeData.fire();
    }

    getTreeItem(element: TreeItemType): vscode.TreeItem {
        return element;
    }

    async getChildren(element?: TreeItemType): Promise<TreeItemType[]> {
        if (!element) {
            // Show activities for active session
            const activeSession = this.sessionManager.getActiveSession();
            if (!activeSession) {
                return [new InfoTreeItem('No active session')];
            }

            try {
                const activities = await this.client.listActivities(activeSession.id);
                if (activities.length === 0) {
                    return [new InfoTreeItem('No activities yet')];
                }
                return activities.map(activity => new ActivityTreeItem(activity));
            } catch (error) {
                return [new InfoTreeItem(`Error: ${error}`)];
            }
        }

        // Show activity details
        if (element instanceof ActivityTreeItem) {
            return this.getActivityDetails(element.activity);
        }

        return [];
    }

    private getActivityDetails(activity: Activity): TreeItemType[] {
        const items: TreeItemType[] = [];

        // Add description
        if (activity.description) {
            items.push(new DetailTreeItem('Description', activity.description));
        }

        // Add originator
        if (activity.originator) {
            items.push(new DetailTreeItem('Originator', activity.originator));
        }

        // Add artifacts
        if (activity.artifacts && activity.artifacts.length > 0) {
            for (const artifact of activity.artifacts) {
                if (artifact.changeSet?.gitPatch) {
                    items.push(new DetailTreeItem(
                        'Git Patch',
                        artifact.changeSet.gitPatch.suggestedCommitMessage || 'View patch'
                    ));
                }
                if (artifact.bashOutput) {
                    items.push(new DetailTreeItem(
                        'Bash',
                        `${artifact.bashOutput.command} (exit ${artifact.bashOutput.exitCode})`
                    ));
                }
            }
        }

        return items;
    }
}

class ActivityTreeItem extends vscode.TreeItem {
    constructor(public readonly activity: Activity) {
        super(activity.description || activity.name, vscode.TreeItemCollapsibleState.Collapsed);

        this.description = new Date(activity.createTime).toLocaleTimeString();
        this.tooltip = `ID: ${activity.id}\nCreated: ${new Date(activity.createTime).toLocaleString()}\n\n${activity.description}`;
        this.iconPath = new vscode.ThemeIcon('pulse');
        this.contextValue = 'activity';
    }
}

class DetailTreeItem extends vscode.TreeItem {
    constructor(label: string, description: string) {
        super(label, vscode.TreeItemCollapsibleState.None);
        this.description = description;
        this.iconPath = new vscode.ThemeIcon('info');
    }
}

class InfoTreeItem extends vscode.TreeItem {
    constructor(message: string) {
        super(message, vscode.TreeItemCollapsibleState.None);
        this.iconPath = new vscode.ThemeIcon('info');
    }
}
