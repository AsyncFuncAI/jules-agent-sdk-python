/**
 * Git Integration - Handle PR and branch operations
 */

import * as vscode from 'vscode';
import { JulesSession } from '../types';

export class GitIntegration {
    constructor(private outputChannel: vscode.OutputChannel) {}

    /**
     * Open PR in browser
     */
    async openPullRequest(session: JulesSession): Promise<void> {
        if (!session.outputs || session.outputs.length === 0) {
            vscode.window.showWarningMessage('No outputs available for this session.');
            return;
        }

        const pr = session.outputs[0]?.pullRequest;
        if (!pr || !pr.url) {
            vscode.window.showWarningMessage('No pull request found for this session.');
            return;
        }

        await vscode.env.openExternal(vscode.Uri.parse(pr.url));
    }

    /**
     * Show PR details in a webview
     */
    async showPullRequestDetails(session: JulesSession): Promise<void> {
        if (!session.outputs || session.outputs.length === 0) {
            vscode.window.showWarningMessage('No outputs available for this session.');
            return;
        }

        const pr = session.outputs[0]?.pullRequest;
        if (!pr) {
            vscode.window.showWarningMessage('No pull request found for this session.');
            return;
        }

        const panel = vscode.window.createWebviewPanel(
            'julesPR',
            `PR: ${pr.title}`,
            vscode.ViewColumn.One,
            { enableScripts: true }
        );

        panel.webview.html = this.getPRWebviewContent(pr.title, pr.description, pr.url);
    }

    private getPRWebviewContent(title: string, description: string, url: string): string {
        return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>${title}</title>
    <style>
        body {
            font-family: var(--vscode-font-family);
            color: var(--vscode-foreground);
            background-color: var(--vscode-editor-background);
            padding: 20px;
            line-height: 1.6;
        }
        h1 {
            color: var(--vscode-editor-foreground);
            border-bottom: 1px solid var(--vscode-panel-border);
            padding-bottom: 10px;
        }
        .description {
            white-space: pre-wrap;
            padding: 15px;
            background-color: var(--vscode-textBlockQuote-background);
            border-left: 3px solid var(--vscode-textBlockQuote-border);
            margin: 20px 0;
        }
        .actions {
            margin-top: 30px;
        }
        a {
            display: inline-block;
            padding: 10px 20px;
            background-color: var(--vscode-button-background);
            color: var(--vscode-button-foreground);
            text-decoration: none;
            border-radius: 3px;
            margin-right: 10px;
        }
        a:hover {
            background-color: var(--vscode-button-hoverBackground);
        }
    </style>
</head>
<body>
    <h1>${title}</h1>
    <div class="description">${description}</div>
    <div class="actions">
        <a href="${url}" target="_blank">View on GitHub</a>
    </div>
</body>
</html>`;
    }

    /**
     * Get current git branch
     */
    async getCurrentBranch(): Promise<string | undefined> {
        const gitExtension = vscode.extensions.getExtension('vscode.git')?.exports;
        if (!gitExtension) {
            this.outputChannel.appendLine('Git extension not found');
            return undefined;
        }

        const api = gitExtension.getAPI(1);
        if (api.repositories.length === 0) {
            return undefined;
        }

        const repo = api.repositories[0];
        return repo.state.HEAD?.name;
    }

    /**
     * Get current repository info
     */
    async getRepositoryInfo(): Promise<{ owner: string; repo: string } | undefined> {
        const gitExtension = vscode.extensions.getExtension('vscode.git')?.exports;
        if (!gitExtension) {
            return undefined;
        }

        const api = gitExtension.getAPI(1);
        if (api.repositories.length === 0) {
            return undefined;
        }

        const repo = api.repositories[0];
        const remotes = repo.state.remotes;

        if (remotes.length === 0) {
            return undefined;
        }

        // Try to parse GitHub URL from remote
        const remote = remotes[0];
        const url = remote.fetchUrl || remote.pushUrl;

        if (!url) {
            return undefined;
        }

        // Match github.com URLs
        const match = url.match(/github\.com[:/]([^/]+)\/([^/.]+)/);
        if (match) {
            return {
                owner: match[1],
                repo: match[2]
            };
        }

        return undefined;
    }
}
