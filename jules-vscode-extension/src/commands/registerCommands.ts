/**
 * Register all extension commands
 */

import * as vscode from 'vscode';
import { SessionManager } from '../sessionManager';
import { JulesClient } from '../julesClient';
import { JulesPanelProvider } from '../julesPanel';
import { SessionState } from '../types';

export function registerCommands(
    context: vscode.ExtensionContext,
    sessionManager: SessionManager,
    client: JulesClient,
    julesPanel: JulesPanelProvider,
    outputChannel: vscode.OutputChannel
): void {

    // Start new session
    context.subscriptions.push(
        vscode.commands.registerCommand('jules.startSession', async () => {
            try {
                // Get sources
                const sources = await client.listSources();
                if (sources.length === 0) {
                    vscode.window.showErrorMessage('No sources available. Please configure sources in Jules dashboard.');
                    return;
                }

                // Let user select source
                const sourceItems = sources.map(s => ({
                    label: s.githubRepo ? `${s.githubRepo.owner}/${s.githubRepo.repo}` : s.id,
                    description: s.githubRepo?.isPrivate ? '(Private)' : '(Public)',
                    detail: s.name,
                    source: s
                }));

                const selectedSource = await vscode.window.showQuickPick(sourceItems, {
                    placeHolder: 'Select a repository/source'
                });

                if (!selectedSource) {
                    return;
                }

                // Get starting branch if GitHub repo
                let startingBranch: string | undefined;
                if (selectedSource.source.githubRepo) {
                    startingBranch = selectedSource.source.githubRepo.defaultBranch || 'main';
                    const branchInput = await vscode.window.showInputBox({
                        prompt: 'Enter starting branch',
                        value: startingBranch,
                        placeHolder: 'main'
                    });
                    if (branchInput) {
                        startingBranch = branchInput;
                    }
                }

                // Get task prompt
                const prompt = await vscode.window.showInputBox({
                    prompt: 'What would you like Jules to do?',
                    placeHolder: 'e.g., Add error handling to the authentication module',
                    validateInput: (value) => {
                        return value.trim().length > 0 ? null : 'Prompt cannot be empty';
                    }
                });

                if (!prompt) {
                    return;
                }

                // Get optional title
                const title = await vscode.window.showInputBox({
                    prompt: 'Session title (optional)',
                    placeHolder: 'Leave empty for auto-generated title'
                });

                // Check if plan approval required
                const config = vscode.workspace.getConfiguration('jules');
                const autoApprove = config.get<boolean>('autoApprovePlans', false);

                // Create session
                await sessionManager.createSession(
                    prompt,
                    selectedSource.source.name,
                    {
                        startingBranch,
                        title: title || undefined,
                        requirePlanApproval: !autoApprove
                    }
                );

                vscode.window.showInformationMessage('Session created! Check the Jules panel.');

            } catch (error) {
                const errorMsg = error instanceof Error ? error.message : String(error);
                vscode.window.showErrorMessage(`Failed to start session: ${errorMsg}`);
                outputChannel.appendLine(`Error starting session: ${errorMsg}`);
            }
        })
    );

    // Open Jules panel
    context.subscriptions.push(
        vscode.commands.registerCommand('jules.openChat', async () => {
            vscode.commands.executeCommand('workbench.view.extension.jules-sidebar');
        })
    );

    // List sessions
    context.subscriptions.push(
        vscode.commands.registerCommand('jules.listSessions', async () => {
            try {
                const sessions = await sessionManager.listSessions(true);

                if (sessions.length === 0) {
                    vscode.window.showInformationMessage('No sessions found.');
                    return;
                }

                const sessionItems = sessions.map(s => ({
                    label: s.title || s.id,
                    description: s.state,
                    detail: `Created: ${new Date(s.createTime).toLocaleString()}`,
                    session: s
                }));

                const selected = await vscode.window.showQuickPick(sessionItems, {
                    placeHolder: 'Select a session'
                });

                if (selected && selected.session.url) {
                    vscode.env.openExternal(vscode.Uri.parse(selected.session.url));
                }
            } catch (error) {
                vscode.window.showErrorMessage(`Failed to list sessions: ${error}`);
            }
        })
    );

    // Approve plan
    context.subscriptions.push(
        vscode.commands.registerCommand('jules.approvePlan', async () => {
            const activeSession = sessionManager.getActiveSession();

            if (!activeSession) {
                vscode.window.showWarningMessage('No active session. Please start a session first.');
                return;
            }

            if (activeSession.state !== SessionState.AWAITING_PLAN_APPROVAL) {
                vscode.window.showWarningMessage(`Session is not awaiting plan approval (current state: ${activeSession.state})`);
                return;
            }

            const confirm = await vscode.window.showInformationMessage(
                `Approve plan for session: ${activeSession.title || activeSession.id}?`,
                'Approve',
                'Cancel'
            );

            if (confirm === 'Approve') {
                try {
                    await sessionManager.approvePlan(activeSession.id);
                    vscode.window.showInformationMessage('Plan approved. Jules is now executing...');
                } catch (error) {
                    vscode.window.showErrorMessage(`Failed to approve plan: ${error}`);
                }
            }
        })
    );

    // Refresh activities
    context.subscriptions.push(
        vscode.commands.registerCommand('jules.refreshActivities', async () => {
            vscode.window.showInformationMessage('Refreshing activities...');
        })
    );

    // Stop session
    context.subscriptions.push(
        vscode.commands.registerCommand('jules.stopSession', async () => {
            const activeSession = sessionManager.getActiveSession();

            if (!activeSession) {
                vscode.window.showWarningMessage('No active session.');
                return;
            }

            const confirm = await vscode.window.showWarningMessage(
                `Stop monitoring session: ${activeSession.title || activeSession.id}?`,
                'Stop',
                'Cancel'
            );

            if (confirm === 'Stop') {
                await sessionManager.stopSession(activeSession.id);
            }
        })
    );

    // Open session in browser
    context.subscriptions.push(
        vscode.commands.registerCommand('jules.openSession', async (sessionId: string) => {
            const session = await sessionManager.getSession(sessionId);
            if (session && session.url) {
                vscode.env.openExternal(vscode.Uri.parse(session.url));
            }
        })
    );

    // Open settings
    context.subscriptions.push(
        vscode.commands.registerCommand('jules.openSettings', () => {
            vscode.commands.executeCommand('workbench.action.openSettings', 'jules');
        })
    );

    outputChannel.appendLine('Commands registered successfully');
}
