/**
 * Extension test suite
 */

import * as assert from 'assert';
import * as vscode from 'vscode';

suite('Extension Test Suite', () => {
    vscode.window.showInformationMessage('Start all tests.');

    test('Extension should be present', () => {
        assert.ok(vscode.extensions.getExtension('jules-agent.jules-agent'));
    });

    test('Extension should activate', async () => {
        const ext = vscode.extensions.getExtension('jules-agent.jules-agent');
        assert.ok(ext);
        await ext.activate();
        assert.strictEqual(ext.isActive, true);
    });

    test('Commands should be registered', async () => {
        const commands = await vscode.commands.getCommands(true);
        const julesCommands = [
            'jules.startSession',
            'jules.openChat',
            'jules.listSessions',
            'jules.approvePlan',
            'jules.refreshActivities',
            'jules.stopSession',
            'jules.openSettings'
        ];

        for (const cmd of julesCommands) {
            assert.ok(commands.includes(cmd), `Command ${cmd} should be registered`);
        }
    });

    test('Configuration should have default values', () => {
        const config = vscode.workspace.getConfiguration('jules');
        assert.strictEqual(config.get('autoApprovePlans'), false);
        assert.strictEqual(config.get('pythonPath'), 'python3');
        assert.strictEqual(config.get('sessionTimeout'), 600);
        assert.strictEqual(config.get('pollInterval'), 5);
    });
});
