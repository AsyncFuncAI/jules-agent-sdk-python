# Change Log

All notable changes to the "jules-agent" extension will be documented in this file.

## [0.1.0] - 2025-01-06

### Added
- Initial release of Jules Agent VSCode Extension
- Interactive chat interface for Jules sessions
- Session management (create, monitor, stop)
- Activity stream viewer with real-time updates
- Code action providers (fix, refactor, test generation)
- Git/PR integration
- Plan approval workflow
- Status bar integration
- Configuration management with secure API key storage
- Python bridge for Jules SDK communication
- Comprehensive command palette commands
- Tree view for sessions and activities
- Session state monitoring with automatic polling
- Support for multiple sources/repositories

### Features
- **Chat View**: Side panel for conversing with Jules
- **Sessions View**: Browse all sessions with state indicators
- **Activities View**: Monitor session activities in real-time
- **Code Actions**: Right-click menu integration for quick tasks
- **Status Bar**: Quick access indicator
- **Commands**: Full command palette support

### Configuration Options
- API key management (settings, env var, or prompt)
- Custom Python path support
- Configurable session timeout and poll intervals
- Auto-approve plans option
- Custom base URL support

### Developer Tools
- TypeScript + Webpack build system
- ESLint configuration
- Test suite with Mocha
- Comprehensive error handling
- Debug logging to output channel

## [Unreleased]

### Planned Features
- Inline diff preview for code changes
- Session templates for common tasks
- Bulk session operations
- Export session history
- Custom prompt snippets
- Team collaboration features
- Session sharing and replay
- Advanced Git operations (cherry-pick, rebase support)
- Multi-workspace support
- Keybinding customization
