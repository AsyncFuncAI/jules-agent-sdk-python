# Jules Agent VSCode Extension - Setup Guide

## Quick Start

### 1. Install Dependencies

```bash
# Install Node.js dependencies
cd jules-vscode-extension
npm install

# Install Python dependencies
cd python
pip install -r requirements.txt
# or use the SDK from parent directory
pip install -e ../../
cd ..
```

### 2. Configure Development Environment

Create a `.env` file (optional):
```env
JULES_API_KEY=your-api-key-here
```

### 3. Build the Extension

```bash
# Compile TypeScript
npm run compile

# Or watch for changes
npm run watch
```

### 4. Run Extension in Development Mode

1. Open VSCode in the `jules-vscode-extension` directory
2. Press `F5` to launch Extension Development Host
3. The extension will be loaded in a new VSCode window

### 5. Test the Extension

```bash
# Run tests
npm test

# Run linter
npm run lint
```

## Development Workflow

### Project Structure

```
jules-vscode-extension/
├── src/
│   ├── extension.ts              # Main entry point
│   ├── julesClient.ts            # Python bridge client
│   ├── chatView.ts               # Chat webview
│   ├── sessionManager.ts         # Session lifecycle
│   ├── types.ts                  # Type definitions
│   ├── commands/                 # Command implementations
│   │   └── registerCommands.ts
│   ├── providers/                # VSCode providers
│   │   ├── sessionsTreeProvider.ts
│   │   ├── activitiesTreeProvider.ts
│   │   ├── codeActionProvider.ts
│   │   └── gitIntegration.ts
│   └── test/                     # Test suite
├── python/
│   ├── bridge.py                 # JSON-RPC bridge
│   └── requirements.txt
├── media/                        # Icons and assets
├── package.json                  # Extension manifest
├── tsconfig.json                 # TypeScript config
└── webpack.config.js             # Build config
```

### Key Files

- **[extension.ts](src/extension.ts)**: Extension activation and initialization
- **[julesClient.ts](src/julesClient.ts)**: Communication with Python bridge
- **[bridge.py](python/bridge.py)**: Python server using Jules SDK
- **[package.json](package.json)**: Extension metadata and commands

### Testing Python Bridge Standalone

```bash
cd python
python bridge.py
```

Then send JSON-RPC requests via stdin:
```json
{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"api_key": "your-key"}}
```

### Debugging

1. **TypeScript Code**: Set breakpoints in VSCode, press F5
2. **Python Bridge**: Add `print()` statements in `bridge.py`
3. **Extension Logs**: View in Output → Jules Agent

### Common Issues

**Extension won't load:**
- Run `npm run compile` first
- Check for TypeScript errors
- Verify all dependencies installed

**Python bridge errors:**
- Check Python path in settings
- Ensure Jules SDK installed: `pip show jules-agent-sdk`
- Check bridge.py has execute permissions

**API errors:**
- Verify API key is correct
- Check network connectivity
- Review Jules API documentation

## Building for Production

### Create VSIX Package

```bash
# Install vsce if not already installed
npm install -g @vscode/vsce

# Package extension
npm run package
vsce package
```

This creates a `.vsix` file you can:
- Install locally: Extensions → Install from VSIX
- Share with others
- Publish to marketplace

### Publishing to Marketplace

1. Create publisher account at https://marketplace.visualstudio.com/
2. Generate Personal Access Token
3. Login: `vsce login <publisher-name>`
4. Publish: `vsce publish`

## Configuration for Users

After installation, users should:

1. **Set API Key**:
   - Command Palette: `Jules: Open Settings`
   - Set `jules.apiKey`
   - Or set environment variable: `JULES_API_KEY`

2. **Configure Python Path** (if needed):
   ```json
   {
     "jules.pythonPath": "/path/to/python3"
   }
   ```

3. **Optional Settings**:
   ```json
   {
     "jules.autoApprovePlans": false,
     "jules.sessionTimeout": 600,
     "jules.pollInterval": 5
   }
   ```

## Next Steps

1. Review [README.md](README.md) for usage instructions
2. Check [CHANGELOG.md](CHANGELOG.md) for version history
3. Read [VSCode Extension API docs](https://code.visualstudio.com/api)
4. Explore [Jules API documentation](https://developers.google.com/jules/api/)

## Contributing

See main project CONTRIBUTING.md for guidelines.

## Support

- GitHub Issues: Report bugs and request features
- Documentation: Full API reference in docs/
- Jules Dashboard: https://jules.google.com
