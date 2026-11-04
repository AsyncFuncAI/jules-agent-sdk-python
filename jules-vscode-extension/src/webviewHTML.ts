export function getWebviewHTML(): string {
    return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Jules Agent</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: var(--vscode-font-family);
            color: var(--vscode-foreground);
            background: var(--vscode-editor-background);
            height: 100vh;
            display: flex;
            flex-direction: column;
            overflow: hidden;
        }

        /* Tabs */
        #tabs {
            display: flex;
            background: var(--vscode-sideBar-background);
            border-bottom: 1px solid var(--vscode-panel-border);
        }
        .tab {
            flex: 1;
            padding: 8px 12px;
            cursor: pointer;
            text-align: center;
            font-size: 12px;
            border-bottom: 2px solid transparent;
            opacity: 0.7;
        }
        .tab:hover {
            background: var(--vscode-list-hoverBackground);
            opacity: 1;
        }
        .tab.active {
            opacity: 1;
            border-bottom-color: var(--vscode-button-background);
            background: var(--vscode-editor-background);
        }

        /* Tab Content */
        .tab-content {
            display: none;
            flex: 1;
            flex-direction: column;
            overflow: hidden;
        }
        .tab-content.active {
            display: flex;
        }

        /* Sessions Tab */
        #sessions-list {
            flex: 1;
            overflow-y: auto;
        }
        .session-item {
            padding: 8px 12px;
            border-bottom: 1px solid var(--vscode-panel-border);
            cursor: pointer;
        }
        .session-item:hover {
            background: var(--vscode-list-hoverBackground);
        }
        .session-item.selected {
            background: var(--vscode-list-activeSelectionBackground);
        }
        .session-header {
            display: flex;
            align-items: center;
            gap: 8px;
            margin-bottom: 4px;
        }
        .session-icon {
            font-size: 14px;
        }
        .session-title {
            font-weight: 500;
            font-size: 13px;
            flex: 1;
        }
        .session-state {
            font-size: 10px;
            padding: 2px 6px;
            border-radius: 3px;
            background: var(--vscode-badge-background);
            color: var(--vscode-badge-foreground);
        }
        .session-prompt {
            font-size: 11px;
            opacity: 0.7;
            white-space: nowrap;
            overflow: hidden;
            text-overflow: ellipsis;
            padding-left: 22px;
        }
        .session-meta {
            font-size: 10px;
            opacity: 0.5;
            margin-top: 4px;
            padding-left: 22px;
        }

        /* Activities Tab - Vercel-like Chat UI */
        #activities-list {
            flex: 1;
            overflow-y: auto;
            padding: 16px 12px;
            display: flex;
            flex-direction: column;
            gap: 16px;
        }
        .activity-message {
            display: flex;
            flex-direction: column;
            gap: 8px;
        }
        .message-header {
            display: flex;
            align-items: center;
            gap: 8px;
            font-size: 11px;
            opacity: 0.6;
        }
        .message-avatar {
            width: 20px;
            height: 20px;
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 10px;
        }
        .avatar-user {
            background: var(--vscode-button-background);
        }
        .avatar-agent {
            background: var(--vscode-editorInfo-foreground);
        }
        .message-content {
            padding: 12px 16px;
            border-radius: 8px;
            font-size: 13px;
            line-height: 1.6;
            max-width: 100%;
        }
        .user-message .message-content {
            background: var(--vscode-input-background);
            border: 1px solid var(--vscode-input-border);
        }
        .agent-message .message-content {
            background: var(--vscode-editor-background);
            border: 1px solid var(--vscode-panel-border);
        }
        .message-text {
            margin-bottom: 12px;
            white-space: pre-wrap;
            word-wrap: break-word;
        }
        .message-text:last-child {
            margin-bottom: 0;
        }

        /* Code blocks - Vercel style */
        .code-section {
            margin-top: 12px;
        }
        .code-header {
            font-size: 11px;
            font-weight: 500;
            margin-bottom: 6px;
            opacity: 0.8;
            display: flex;
            align-items: center;
            gap: 6px;
        }
        .code-block {
            background: var(--vscode-textCodeBlock-background);
            border: 1px solid var(--vscode-panel-border);
            border-radius: 6px;
            padding: 12px;
            font-family: var(--vscode-editor-font-family), 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', Consolas, 'Courier New', monospace;
            font-size: 12px;
            line-height: 1.5;
            overflow-x: auto;
            white-space: pre;
        }
        .code-block::-webkit-scrollbar {
            height: 8px;
        }
        .code-block::-webkit-scrollbar-thumb {
            background: var(--vscode-scrollbarSlider-background);
            border-radius: 4px;
        }

        /* Terminal output - Vercel style */
        .terminal-output {
            background: #0a0a0a;
            border: 1px solid rgba(255,255,255,0.1);
            border-radius: 6px;
            padding: 12px;
            font-family: 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', Consolas, 'Courier New', monospace;
            font-size: 12px;
            line-height: 1.5;
            color: #e5e5e5;
        }
        .terminal-command {
            color: #00d9ff;
            margin-bottom: 8px;
            display: flex;
            align-items: center;
            gap: 6px;
        }
        .terminal-command::before {
            content: '❯';
            color: #888;
        }
        .terminal-output-text {
            color: #a8a8a8;
            white-space: pre-wrap;
            word-wrap: break-word;
        }
        .terminal-exit {
            margin-top: 8px;
            font-size: 11px;
            opacity: 0.5;
        }

        /* Progress indicator */
        .progress-message {
            display: flex;
            align-items: center;
            gap: 8px;
            padding: 10px 14px;
            background: var(--vscode-editor-inactiveSelectionBackground);
            border-radius: 6px;
            font-size: 12px;
            border-left: 3px solid var(--vscode-editorInfo-foreground);
        }
        .progress-spinner {
            animation: spin 1s linear infinite;
        }
        @keyframes spin {
            from { transform: rotate(0deg); }
            to { transform: rotate(360deg); }
        }

        /* Commit message styling */
        .commit-message {
            font-size: 11px;
            padding: 8px 12px;
            background: var(--vscode-editor-inactiveSelectionBackground);
            border-radius: 4px;
            margin-bottom: 8px;
            font-style: italic;
            opacity: 0.9;
        }

        /* Chat Tab */
        #chat-history {
            flex: 1;
            overflow-y: auto;
            padding: 8px;
        }
        #chat-input-area {
            border-top: 1px solid var(--vscode-panel-border);
            padding: 8px;
            background: var(--vscode-sideBar-background);
        }
        #message-input {
            width: 100%;
            min-height: 50px;
            padding: 6px;
            background: var(--vscode-input-background);
            color: var(--vscode-input-foreground);
            border: 1px solid var(--vscode-input-border);
            border-radius: 3px;
            resize: vertical;
            font-family: var(--vscode-font-family);
            font-size: 12px;
            margin-bottom: 6px;
        }
        button {
            width: 100%;
            padding: 6px 12px;
            background: var(--vscode-button-background);
            color: var(--vscode-button-foreground);
            border: none;
            border-radius: 3px;
            cursor: pointer;
            font-size: 12px;
        }
        button:hover {
            background: var(--vscode-button-hoverBackground);
        }
        .btn-secondary {
            background: var(--vscode-button-secondaryBackground);
            color: var(--vscode-button-secondaryForeground);
        }

        /* Toolbar */
        #toolbar {
            padding: 6px 8px;
            background: var(--vscode-sideBar-background);
            border-bottom: 1px solid var(--vscode-panel-border);
            display: flex;
            gap: 6px;
        }
        #toolbar button {
            width: auto;
            font-size: 11px;
            padding: 4px 8px;
        }

        /* Session banner */
        #session-banner {
            padding: 8px;
            background: var(--vscode-editor-inactiveSelectionBackground);
            border-bottom: 1px solid var(--vscode-panel-border);
            font-size: 11px;
            display: none;
        }
        #session-banner.visible {
            display: block;
        }
        .banner-title {
            font-weight: 500;
            margin-bottom: 2px;
        }

        .empty-state {
            text-align: center;
            padding: 40px 20px;
            opacity: 0.5;
            font-size: 12px;
        }
    </style>
</head>
<body>
    <div id="toolbar">
        <button onclick="createSession()">+ New</button>
        <button class="btn-secondary" onclick="refreshData()">↻ Refresh</button>
        <button class="btn-secondary" onclick="openInBrowser()" id="browser-btn" style="display:none">🌐 Browser</button>
    </div>

    <div id="session-banner">
        <div class="banner-title" id="banner-title"></div>
        <div id="banner-state"></div>
    </div>

    <div id="tabs">
        <div class="tab active" data-tab="sessions" onclick="switchTab('sessions')">📋 Sessions</div>
        <div class="tab" data-tab="activities" onclick="switchTab('activities')">⚡ Activities</div>
        <div class="tab" data-tab="chat" onclick="switchTab('chat')">💬 Chat</div>
    </div>

    <div id="sessions-tab" class="tab-content active">
        <div id="sessions-list">
            <div class="empty-state">Loading...</div>
        </div>
    </div>

    <div id="activities-tab" class="tab-content">
        <div id="activities-list">
            <div class="empty-state">Select a session</div>
        </div>
    </div>

    <div id="chat-tab" class="tab-content">
        <div id="chat-history">
            <div class="empty-state">Select a session to chat</div>
        </div>
        <div id="chat-input-area">
            <textarea id="message-input" placeholder="Send a message to Jules..."></textarea>
            <button onclick="sendMessage()">Send Message</button>
        </div>
    </div>

    <script>
        const vscode = acquireVsCodeApi();
        let sessions = [];
        let activities = [];
        let selectedSessionId = null;

        window.addEventListener('DOMContentLoaded', () => {
            vscode.postMessage({ type: 'ready' });
        });

        window.addEventListener('message', event => {
            const msg = event.data;
            switch (msg.type) {
                case 'sessions': renderSessions(msg.sessions); break;
                case 'sessionSelected': selectSessionUI(msg.session); break;
                case 'activities': renderActivities(msg.activities); break;
                case 'messageSent': addChatMessage(msg.message); break;
            }
        });

        function switchTab(tab) {
            document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
            document.querySelectorAll('.tab-content').forEach(t => t.classList.remove('active'));
            document.querySelector(\`.tab[data-tab="\${tab}"]\`).classList.add('active');
            document.getElementById(\`\${tab}-tab\`).classList.add('active');
        }

        function renderSessions(list) {
            sessions = list;
            const container = document.getElementById('sessions-list');
            if (!list.length) {
                container.innerHTML = '<div class="empty-state">No sessions</div>';
                return;
            }
            container.innerHTML = list.map(s => \`
                <div class="session-item \${s.id === selectedSessionId ? 'selected' : ''}" onclick="selectSession('\${s.id}')">
                    <div class="session-header">
                        <span class="session-icon">\${getIcon(s.state)}</span>
                        <span class="session-title">\${escapeHtml(s.title)}</span>
                        <span class="session-state">\${s.state}</span>
                    </div>
                    <div class="session-prompt">\${escapeHtml(s.prompt)}</div>
                    <div class="session-meta">\${formatTime(s.updateTime)}</div>
                </div>
            \`).join('');
        }

        function selectSession(id) {
            selectedSessionId = id;
            vscode.postMessage({ type: 'selectSession', sessionId: id });
            renderSessions(sessions);
        }

        function selectSessionUI(session) {
            const banner = document.getElementById('session-banner');
            document.getElementById('banner-title').textContent = session.title;
            document.getElementById('banner-state').textContent = \`State: \${session.state} • \${session.prompt}\`;
            banner.className = 'visible';
            document.getElementById('browser-btn').style.display = 'inline-block';
        }

        function renderActivities(list) {
            activities = list;
            const container = document.getElementById('activities-list');
            if (!list.length) {
                container.innerHTML = '<div class="empty-state">No activities</div>';
                return;
            }
            container.innerHTML = list.map((a, i) => renderActivity(a, i)).join('');
            container.scrollTop = container.scrollHeight;
        }

        function renderActivity(a, i) {
            // User message
            if (a.userMessage) {
                return \`
                    <div class="activity-message user-message">
                        <div class="message-header">
                            <div class="message-avatar avatar-user">👤</div>
                            <span>You</span>
                            <span>•</span>
                            <span>\${formatTime(a.createTime)}</span>
                        </div>
                        <div class="message-content">
                            <div class="message-text">\${escapeHtml(a.userMessage)}</div>
                        </div>
                    </div>
                \`;
            }

            // Agent message
            if (a.agentMessage) {
                let content = \`<div class="message-text">\${escapeHtml(a.agentMessage)}</div>\`;

                // Add code changes
                a.artifacts?.forEach(art => {
                    if (art.gitPatch?.patch) {
                        content += \`
                            <div class="code-section">
                                <div class="code-header">
                                    <span>📝</span>
                                    <span>Code Changes</span>
                                </div>
                                \${art.gitPatch.commitMessage ? \`<div class="commit-message">\${escapeHtml(art.gitPatch.commitMessage)}</div>\` : ''}
                                <div class="code-block">\${escapeHtml(art.gitPatch.patch)}</div>
                            </div>
                        \`;
                    }
                    if (art.bashOutput) {
                        content += \`
                            <div class="code-section">
                                <div class="code-header">
                                    <span>💻</span>
                                    <span>Terminal</span>
                                </div>
                                <div class="terminal-output">
                                    <div class="terminal-command">\${escapeHtml(art.bashOutput.command)}</div>
                                    <div class="terminal-output-text">\${escapeHtml(art.bashOutput.output)}</div>
                                    <div class="terminal-exit">Exit code: \${art.bashOutput.exitCode}</div>
                                </div>
                            </div>
                        \`;
                    }
                });

                return \`
                    <div class="activity-message agent-message">
                        <div class="message-header">
                            <div class="message-avatar avatar-agent">🤖</div>
                            <span>Jules</span>
                            <span>•</span>
                            <span>\${formatTime(a.createTime)}</span>
                        </div>
                        <div class="message-content">
                            \${content}
                        </div>
                    </div>
                \`;
            }

            // Progress update
            if (a.progressUpdate) {
                return \`
                    <div class="progress-message">
                        <span class="progress-spinner">⏳</span>
                        <div>
                            <strong>\${escapeHtml(a.progressUpdate.title)}</strong><br>
                            <span style="opacity:0.7">\${escapeHtml(a.progressUpdate.description)}</span>
                        </div>
                    </div>
                \`;
            }

            // Code-only activities
            let hasContent = false;
            let content = '';

            a.artifacts?.forEach(art => {
                if (art.gitPatch?.patch) {
                    hasContent = true;
                    content += \`
                        <div class="code-section">
                            <div class="code-header">
                                <span>📝</span>
                                <span>Code Changes</span>
                            </div>
                            \${art.gitPatch.commitMessage ? \`<div class="commit-message">\${escapeHtml(art.gitPatch.commitMessage)}</div>\` : ''}
                            <div class="code-block">\${escapeHtml(art.gitPatch.patch)}</div>
                        </div>
                    \`;
                }
                if (art.bashOutput) {
                    hasContent = true;
                    content += \`
                        <div class="code-section">
                            <div class="code-header">
                                <span>💻</span>
                                <span>Terminal</span>
                            </div>
                            <div class="terminal-output">
                                <div class="terminal-command">\${escapeHtml(art.bashOutput.command)}</div>
                                <div class="terminal-output-text">\${escapeHtml(art.bashOutput.output)}</div>
                                <div class="terminal-exit">Exit code: \${art.bashOutput.exitCode}</div>
                            </div>
                        </div>
                    \`;
                }
            });

            if (hasContent) {
                return \`
                    <div class="activity-message agent-message">
                        <div class="message-header">
                            <div class="message-avatar avatar-agent">🤖</div>
                            <span>Jules</span>
                            <span>•</span>
                            <span>\${formatTime(a.createTime)}</span>
                        </div>
                        <div class="message-content">
                            \${content}
                        </div>
                    </div>
                \`;
            }

            return '';
        }

        function addChatMessage(msg) {
            const chat = document.getElementById('chat-history');
            const div = document.createElement('div');
            div.className = \`message-bubble \${msg.role}-message\`;
            div.textContent = msg.content;
            chat.appendChild(div);
            chat.scrollTop = chat.scrollHeight;
        }

        function sendMessage() {
            const input = document.getElementById('message-input');
            const msg = input.value.trim();
            if (!msg) return;
            vscode.postMessage({ type: 'sendMessage', message: msg });
            input.value = '';
        }

        function createSession() {
            vscode.postMessage({ type: 'createSession' });
        }

        function refreshData() {
            vscode.postMessage({ type: 'refreshSessions' });
        }

        function openInBrowser() {
            vscode.postMessage({ type: 'openInBrowser' });
        }

        function getIcon(state) {
            const icons = {
                'COMPLETED': '✅', 'FAILED': '❌', 'IN_PROGRESS': '🔄',
                'AWAITING_PLAN_APPROVAL': '⏸️', 'AWAITING_USER_FEEDBACK': '💬',
                'PLANNING': '🤔', 'QUEUED': '⏳', 'PAUSED': '⏸️'
            };
            return icons[state] || '⚪';
        }

        function formatTime(t) {
            const d = new Date(t);
            const diff = Date.now() - d;
            if (diff < 60000) return 'Just now';
            if (diff < 3600000) return Math.floor(diff / 60000) + 'm ago';
            if (diff < 86400000) return Math.floor(diff / 3600000) + 'h ago';
            return Math.floor(diff / 86400000) + 'd ago';
        }

        function escapeHtml(text) {
            const div = document.createElement('div');
            div.textContent = text;
            return div.innerHTML;
        }

        document.getElementById('message-input').addEventListener('keydown', (e) => {
            if ((e.ctrlKey || e.metaKey) && e.key === 'Enter') sendMessage();
        });
    </script>
</body>
</html>`;
}
