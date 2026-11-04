# Jules Agent - Quick Start Guide

## ✅ Build Status: SUCCESS

The extension has been successfully built and is ready to use!

## 🚀 Get Started in 3 Steps

### Step 1: Install Dependencies

```bash
cd jules-vscode-extension

# Install Node.js dependencies
npm install

# Install Python dependencies
pip install jules-agent-sdk
```

### Step 2: Get Your API Key

1. Visit [https://jules.google.com](https://jules.google.com)
2. Sign in and navigate to API keys
3. Create a new API key
4. Copy the key (you'll need it in Step 3)

### Step 3: Run the Extension

#### Method A: Development Mode (Recommended for Testing)

1. Open the `jules-vscode-extension` folder in VSCode
2. Press **F5** (or Run → Start Debugging)
3. A new VSCode window opens with the extension loaded
4. Enter your API key when prompted
5. You're ready to use Jules!

#### Method B: Install as VSIX

```bash
# Build the extension package
npm install -g @vscode/vsce
vsce package

# Install in VSCode
code --install-extension jules-agent-*.vsix
```

## 🎯 First Session

Once the extension is running:

1. **Open Command Palette**: `Cmd+Shift+P` (Mac) or `Ctrl+Shift+P` (Windows/Linux)

2. **Run**: `Jules: Start New Session`

3. **Select Your Repository**: Choose from your connected sources

4. **Enter a Task**: For example:
   ```
   Add error handling to the authentication module
   ```

5. **Watch Jules Work**: Monitor progress in the Activities view

## 📱 Key Features

### Chat Interface
- Click the Jules icon in the Activity Bar
- Chat naturally with Jules
- Send follow-up messages to active sessions

### Code Actions
1. Select any code in your editor
2. Right-click → **Jules** menu
3. Choose:
   - Fix this code
   - Refactor this code
   - Write tests for this code

### Session Management
- View all sessions in the Sessions panel
- Click any session to see details
- Monitor real-time progress in Activities panel

### Plan Approval
When a session requires approval:
1. You'll see a notification
2. Review the plan in your browser
3. Run `Jules: Approve Plan` to continue

## ⚙️ Configuration

### Set API Key Permanently

**Option 1: VSCode Settings**
1. `Cmd/Ctrl + ,` to open Settings
2. Search for "jules"
3. Set `Jules: Api Key`

**Option 2: Environment Variable**
```bash
export JULES_API_KEY="your-api-key-here"
```

### Optional Settings

```json
{
  "jules.pythonPath": "python3",           // Custom Python path
  "jules.autoApprovePlans": false,         // Auto-approve plans
  "jules.sessionTimeout": 600,             // Session timeout (seconds)
  "jules.pollInterval": 5                  // Poll interval (seconds)
}
```

## 🐛 Troubleshooting

### Extension Won't Activate

**Problem**: Extension doesn't load or shows errors

**Solutions**:
1. Check Python is installed: `python3 --version`
2. Check Jules SDK is installed: `pip show jules-agent-sdk`
3. View extension logs: Output → Jules Agent
4. Restart VSCode

### Python Bridge Error

**Problem**: "Failed to start Python bridge"

**Solutions**:
1. Set correct Python path in settings
2. Install SDK: `pip install jules-agent-sdk`
3. Check permissions: `chmod +x python/bridge.py`

### API Key Issues

**Problem**: "Invalid API key" or authentication errors

**Solutions**:
1. Verify API key is correct (no extra spaces)
2. Check API key is active in Jules dashboard
3. Try setting via environment variable
4. Clear and re-enter in settings

### Session Not Creating

**Problem**: Session creation fails

**Solutions**:
1. Check you have sources in Jules dashboard
2. Verify network connectivity
3. Check API key permissions
4. Review session timeout settings

## 📚 Learn More

- **Full Documentation**: [README.md](README.md)
- **Development Guide**: [SETUP.md](SETUP.md)
- **Implementation Details**: [IMPLEMENTATION_SUMMARY.md](IMPLEMENTATION_SUMMARY.md)
- **Jules API Docs**: https://developers.google.com/jules/api/

## 💡 Tips

1. **Use Descriptive Prompts**: Be specific about what you want Jules to do
   - ❌ "fix the code"
   - ✅ "Add null checks and error handling to the login function"

2. **Monitor Progress**: Keep the Activities panel open to see what Jules is doing

3. **Approve Plans**: Review plans before approval for better results

4. **Use Chat for Follow-ups**: Send additional instructions via chat after session starts

5. **Check PRs**: Jules creates pull requests - review and merge them

## 🎉 Success!

You're all set! Try creating your first session and watch Jules transform your coding workflow.

**Need help?** Check the [README.md](README.md) or open an issue on GitHub.

---

**Happy Coding with Jules! 🚀**
