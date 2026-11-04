# Jules Agent VSCode Extension - Usage Guide

## ✅ Setup Complete!

Your Jules Agent extension is now working! Here's how to use it effectively.

## 📱 Interface Overview

The Jules Agent sidebar has 3 panels:

### 1. **CHAT Panel**
- Send messages to your active session
- View conversation history
- Get system notifications

### 2. **SESSIONS Panel**
- View all your Jules sessions
- See session state (Completed, In Progress, Failed, etc.)
- Click any session to activate it for chat

### 3. **ACTIVITIES Panel**
- Shows activities for the currently active session
- View what Jules is doing in real-time
- See git patches, bash commands, and progress updates

## 🎯 How to Use

### Step 1: Select a Session

1. Look at the **SESSIONS** panel
2. **Click on any session** - this activates it for chat
3. You'll see a notification showing the selected session
4. The **CHAT** panel will show the active session details

### Step 2: Send Messages

1. Go to the **CHAT** panel
2. Type your message in the text box
3. Click **Send** or press `Ctrl/Cmd + Enter`
4. Your message is sent to the active Jules session

### Step 3: Monitor Progress

1. The **ACTIVITIES** panel shows what Jules is doing
2. You'll see:
   - Agent messages
   - Code changes (git patches)
   - Bash commands executed
   - Progress updates

## 🚀 Common Workflows

### Starting a New Session

1. Open Command Palette (`Cmd/Ctrl + Shift + P`)
2. Run: `Jules: Start New Session`
3. Select your repository
4. Enter your task prompt
5. The session appears in SESSIONS panel

### Continuing an Existing Session

1. Find your session in **SESSIONS** panel
2. Click it to activate
3. Go to **CHAT** panel
4. Send follow-up messages

### Approving a Plan

When a session shows "Needs Approval" state:
1. Click the session to view it
2. Open Command Palette
3. Run: `Jules: Approve Plan`
4. Or click "View in Browser" to review the plan first

### Viewing Session in Browser

1. Click any session in **SESSIONS** panel
2. Click **"View in Browser"** in the notification
3. Or right-click session → context menu options

## 💡 Tips

### Session States

| State | What It Means | What You Can Do |
|-------|---------------|-----------------|
| **Queued** | Waiting to start | Wait for it to begin |
| **Planning** | Jules is making a plan | Wait for plan completion |
| **Needs Approval** | Plan ready for review | Approve or reject the plan |
| **In Progress** | Jules is working | Monitor activities, send messages |
| **Needs Feedback** | Jules has questions | Send a message with your response |
| **Completed** | Session finished | View results, check PR |
| **Failed** | Something went wrong | Check activities for errors |

### Chat Best Practices

1. **Be specific**: "Add error handling to the login function in src/auth.ts"
2. **Provide context**: "The tests are failing because..."
3. **Ask questions**: "What did you change in the last commit?"
4. **Give feedback**: "That looks good, but also add logging"

### When Activities Show "No active session"

This means you haven't selected a session yet. Click any session in the **SESSIONS** panel to activate it.

## 🔧 Keyboard Shortcuts

| Action | Shortcut |
|--------|----------|
| Open Command Palette | `Cmd/Ctrl + Shift + P` |
| Open Jules Sidebar | Click Jules icon in Activity Bar |
| Send Message | `Ctrl/Cmd + Enter` (in chat) |

## 📊 Understanding Sessions

Each session you see represents:
- A specific coding task
- Work on a particular branch
- A conversation history with Jules
- Generated code changes and PRs

Sessions are **persistent** - you can come back to them later and continue the conversation.

## 🎨 Visual Indicators

**Session Icons:**
- ✅ Green checkmark = Completed
- ⚠️ Yellow warning = Needs attention (approval/feedback)
- 🔄 Blue spinning = In progress
- ❌ Red X = Failed
- ⚪ Circle outline = Other states

## 🆘 Troubleshooting

### Chat not working?
- Make sure you've clicked a session to activate it
- Check the session state - it might be completed
- Try selecting a different session

### Activities panel empty?
- Activities only show for the active session
- Click a session in the SESSIONS panel first
- Check if the session has any activities yet

### Can't approve plan?
- Only sessions in "Needs Approval" state can be approved
- Check the session state in the SESSIONS panel
- Make sure the session is active

## 🌐 Web Interface

Each session has a web URL where you can:
- See full conversation history
- View code changes with syntax highlighting
- Review and approve plans
- Access additional features

Click any session and choose "View in Browser" to access it.

## 📝 Example Session

Here's a typical workflow:

1. **Start**: `Jules: Start New Session` → "Add unit tests for the auth module"
2. **Planning**: Jules analyzes code and creates a plan
3. **Approve**: Review plan and approve
4. **Monitor**: Watch activities as Jules writes tests
5. **Feedback**: "Also add integration tests" (via chat)
6. **Complete**: Jules creates a PR with all tests
7. **Review**: Check the PR and merge

## 🎉 You're Ready!

Your Jules Agent extension is fully configured and ready to use. Click any session in the SESSIONS panel to get started!

**Remember**: The key to using Jules effectively is clear communication. The more context you provide, the better the results!
