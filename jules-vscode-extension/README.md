# Jules Agent for VSCode

> AI-powered coding assistant using the Jules API - similar to Claude Code and GitHub Copilot

![Jules Agent](https://img.shields.io/badge/Jules-Agent-blue)
![VSCode](https://img.shields.io/badge/VSCode-Extension-green)
![Python](https://img.shields.io/badge/Python-3.8+-blue)

## Features

🚀 **AI-Powered Coding Assistant**
- Create coding sessions with natural language prompts
- Real-time session monitoring and activity tracking
- Automatic code generation, refactoring, and bug fixes

💬 **Interactive Chat Interface**
- Side panel chat for conversing with Jules
- Session history and context management
- Quick access to common tasks

🔍 **Code Actions**
- Right-click on code to ask Jules to fix, refactor, or explain
- Generate tests for selected code
- Quick inline suggestions

📊 **Session Management**
- View all active and past sessions
- Monitor session progress in real-time
- Approve/reject execution plans
- Automatic PR creation

🔗 **Git Integration**
- Automatic branch creation
- Pull request generation and management
- View PR details directly in VSCode

## Installation

### Prerequisites

1. **Python 3.8+** with `jules-agent-sdk` installed:
   ```bash
   pip install jules-agent-sdk
   ```

2. **Jules API Key**: Get yours from [Jules Dashboard](https://jules.google.com)

### From VSIX

1. Download the latest `.vsix` file from releases
2. In VSCode: `Extensions` → `...` → `Install from VSIX...`
3. Select the downloaded file

### From Source

```bash
cd jules-vscode-extension
npm install
npm run compile
```

Press `F5` in VSCode to launch extension development host.

## Setup

### 1. Configure API Key

**Option A: Through Extension**
- The extension will prompt for your API key on first activation
- Or use Command Palette: `Jules: Open Settings`

**Option B: Environment Variable**
```bash
export JULES_API_KEY="your-api-key-here"
```

**Option C: VSCode Settings**
```json
{
  "jules.apiKey": "your-api-key-here"
}
```

### 2. Configure Python Path

If Python is not in your PATH or you use a virtual environment:

```json
{
  "jules.pythonPath": "/path/to/python3"
}
```

## Usage

### Starting a Session

1. **Command Palette** (`Cmd+Shift+P` / `Ctrl+Shift+P`):
   - `Jules: Start New Session`
   - Select your repository/source
   - Enter your task prompt
   - Jules will create a session and start working!

2. **Via Chat**:
   - Click the Jules icon in the Activity Bar
   - Type your task in the chat
   - Follow prompts to create a session

### Code Actions

1. Select code in your editor
2. Right-click → **Jules** menu
3. Choose:
   - **Ask Jules to fix this** - Fix bugs or issues
   - **Ask Jules to refactor this** - Improve code quality
   - **Ask Jules to write tests** - Generate test cases

### Monitoring Sessions

- **Sessions View**: See all your sessions and their states
- **Activities View**: Monitor what Jules is doing in real-time
- **Status Bar**: Quick access to Jules features

### Approving Plans

When `requirePlanApproval` is enabled:
1. Jules will create a plan and wait
2. Review the plan in the web interface
3. Use `Jules: Approve Plan` command
4. Jules will execute the approved plan

## Configuration

| Setting | Default | Description |
|---------|---------|-------------|
| `jules.apiKey` | `""` | Your Jules API key |
| `jules.baseUrl` | `https://jules.googleapis.com/v1alpha` | API base URL |
| `jules.pythonPath` | `python3` | Path to Python executable |
| `jules.autoApprovePlans` | `false` | Auto-approve plans without confirmation |
| `jules.sessionTimeout` | `600` | Session timeout in seconds |
| `jules.pollInterval` | `5` | Session polling interval in seconds |

### Example Configuration

```json
{
  "jules.apiKey": "your-api-key",
  "jules.pythonPath": "/usr/local/bin/python3",
  "jules.autoApprovePlans": false,
  "jules.sessionTimeout": 900,
  "jules.pollInterval": 3
}
```

## Commands

| Command | Description |
|---------|-------------|
| `Jules: Start New Session` | Create a new coding session |
| `Jules: Open Chat` | Open the chat interface |
| `Jules: List Sessions` | View all sessions |
| `Jules: Approve Plan` | Approve a session plan |
| `Jules: Stop Session` | Stop monitoring a session |
| `Jules: Refresh Activities` | Refresh activity view |
| `Jules: Open Settings` | Open extension settings |

## Architecture

```
┌─────────────────────────────────────┐
│   VSCode Extension (TypeScript)     │
│  ┌─────────────────────────────┐   │
│  │  Chat UI (Webview)          │   │
│  │  Session Manager            │   │
│  │  Activity Viewer            │   │
│  │  Code Action Provider       │   │
│  └─────────────────────────────┘   │
│            ↕ JSON-RPC               │
│  ┌─────────────────────────────┐   │
│  │  Python Bridge              │   │
│  │  (bridge.py)                │   │
│  └─────────────────────────────┘   │
│            ↕ HTTP                   │
│  ┌─────────────────────────────┐   │
│  │  Jules Agent SDK            │   │
│  │  (jules-agent-sdk)          │   │
│  └─────────────────────────────┘   │
└─────────────────────────────────────┘
               ↕ HTTPS
     ┌──────────────────────┐
     │   Jules API          │
     │  (jules.googleapis)  │
     └──────────────────────┘
```

### Components

- **TypeScript Extension**: VSCode UI and interaction layer
- **Python Bridge**: JSON-RPC server using Jules Python SDK
- **Jules SDK**: Official Python SDK for Jules API
- **Webview UI**: Chat interface and session details

## Development

### Setup Development Environment

```bash
# Clone repository
git clone <repo-url>
cd jules-vscode-extension

# Install dependencies
npm install

# Install Python dependencies
cd python
pip install -r requirements.txt
cd ..

# Compile TypeScript
npm run compile

# Watch mode
npm run watch
```

### Running Tests

```bash
# Run all tests
npm test

# Lint code
npm run lint
```

### Building VSIX Package

```bash
npm run package
vsce package
```

## Troubleshooting

### Extension won't activate

1. Check Python is installed: `python3 --version`
2. Check Jules SDK is installed: `pip show jules-agent-sdk`
3. Check extension logs: `Output` → `Jules Agent`

### API Key Issues

1. Verify API key is correct
2. Check environment variable: `echo $JULES_API_KEY`
3. Try setting in VSCode settings directly

### Python Bridge Errors

1. Check Python path in settings
2. Install Jules SDK: `pip install jules-agent-sdk`
3. Check bridge.py permissions: `chmod +x python/bridge.py`

### Session Not Starting

1. Verify you have sources configured in Jules dashboard
2. Check network connectivity
3. Review session timeout settings

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## Related Projects

- [Jules API Python SDK](https://github.com/your-repo/jules-api-python-sdk) - Official Python SDK
- [Jules API](https://developers.google.com/jules/api/) - Official API documentation

## License

MIT License - See [LICENSE](LICENSE) file for details

## Support

- **Issues**: [GitHub Issues](https://github.com/your-repo/issues)
- **Documentation**: [Full Docs](https://docs.jules.dev)
- **Jules Dashboard**: [https://jules.google.com](https://jules.google.com)

## Acknowledgments

Built with:
- [VSCode Extension API](https://code.visualstudio.com/api)
- [Jules Agent SDK](https://pypi.org/project/jules-agent-sdk/)
- Inspired by Claude Code and GitHub Copilot

---

**Made with ❤️ for developers who want AI-powered coding assistance**
