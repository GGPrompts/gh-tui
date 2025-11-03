# gh-tui Hotkeys & Commands

Quick reference for GitHub TUI keyboard shortcuts and workflows.

## 🎯 Global Navigation

### Tab Navigation
```
1 - Pull Requests view
2 - Issues view
3 - Repositories view
4 - Workflow Runs view
5 - Gists view
```

### Quick Switch
```
Tab - Next view
Shift+Tab - Previous view
```

## 📋 List Navigation

### Move in Lists
```
↑ / k - Move up
↓ / j - Move down
g - Jump to top
G - Jump to bottom
Home - Jump to top (alternative)
End - Jump to bottom (alternative)
```

### Page Scroll
```
Ctrl+U - Scroll up half page
Ctrl+D - Scroll down half page
Page Up - Scroll up one page
Page Down - Scroll down one page
```

## 🔍 Search & Filter

### Search
```
/ - Search/filter current view
n - Next search result
N - Previous search result
Esc - Clear search filter
```

### Filters (Coming Soon)
```
f o - Filter: Open
f c - Filter: Closed
f a - Filter: All
f m - Filter: Mine
```

## 📝 Pull Requests View

### PR Actions
```
Enter - View PR details
o - Open PR in browser
c - Checkout PR locally
r - Review PR
m - Merge PR
```

### PR Details
```
Space - Toggle PR description
d - Show diff
w - View in web browser
```

### PR Status Icons
```
🟢 - Open
🟣 - Draft
🔴 - Closed
✅ - Merged
🔥 - Conflicts
👁️ - Changes requested
✓ - Approved
```

## 🐛 Issues View

### Issue Actions
```
Enter - View issue details
o - Open issue in browser
n - Create new issue
e - Edit issue
c - Close issue
r - Reopen issue
```

### Issue Management
```
l - Add/edit labels
a - Assign to user
m - Add milestone
```

### Issue Status Icons
```
🟢 - Open
🟣 - In progress
🔴 - Closed
🐛 - Bug
✨ - Feature
📚 - Documentation
```

## 📦 Repositories View

### Repository Actions
```
Enter - View repository details
o - Open repo in browser
c - Clone repository
f - Fork repository
s - Star repository
```

### Repository Info
```
i - Show repository info
r - Show recent commits
b - List branches
t - List tags
```

### Repository Icons
```
⭐ - Starred
🍴 - Forked
📦 - Private
🌍 - Public
```

## ⚙️ Workflow Runs View

### Workflow Actions
```
Enter - View workflow details
o - Open workflow in browser
r - Rerun workflow
c - Cancel workflow
```

### Workflow Status Icons
```
✅ - Success
❌ - Failed
🟡 - In progress
⏸️ - Queued
🔵 - Waiting
```

### Job Navigation
```
j/k - Navigate jobs in workflow
Enter - View job logs
l - Show job logs
```

## 📄 Gists View

### Gist Actions
```
Enter - View gist content
o - Open gist in browser
n - Create new gist
e - Edit gist
d - Delete gist
```

### Gist Management
```
s - Star gist
f - Fork gist
c - Clone gist
```

## 🎨 Detail Panel

### Navigation
```
↑ / k - Scroll up
↓ / j - Scroll down
g - Scroll to top
G - Scroll to bottom
```

### Content Actions
```
c - Copy content to clipboard
o - Open in browser
Esc - Close detail panel
```

## 🔧 Global Actions

### Refresh & Reload
```
r - Refresh current view
R - Force refresh all data
F5 - Refresh (alternative)
```

### View Options
```
v - Toggle view mode (list/detail)
w - Toggle wrap text
```

### Application Control
```
q - Quit gh-tui
Ctrl+C - Force quit
? - Show help
h - Show help (alternative)
```

## 📊 Common Workflows

### Review Pull Request
```bash
1                    # Switch to PR view
↓↓                   # Navigate to PR
Enter                # View details
d                    # Show diff
r                    # Start review
# Review changes
c                    # Checkout locally (optional)
# Test changes
m                    # Merge when ready
```

### Work on Issue
```bash
2                    # Switch to Issues
/bug                 # Search for bugs
↓                    # Select issue
Enter                # View details
c                    # Close issue when done
```

### Check CI Status
```bash
4                    # Switch to Workflows
↓                    # Navigate to recent run
Enter                # View details
j/k                  # Navigate jobs
l                    # View logs
```

### Create and Share Gist
```bash
5                    # Switch to Gists
n                    # New gist
# Edit in $EDITOR
# Save
o                    # Open in browser to share
```

### Monitor Repository
```bash
3                    # Switch to Repos
↓                    # Select repository
r                    # View recent commits
b                    # Check branches
```

## 🚀 Power User Combos

### Quick PR Review
```bash
1                    # PRs
/                    # Search
# Type PR number
Enter Enter          # Select and view
d                    # Diff
o                    # Open in browser (if needed)
```

### Multi-Repo Check
```bash
3                    # Repos
↓                    # First repo
r                    # Recent commits
Esc                  # Back to list
↓                    # Next repo
r                    # Recent commits
# Repeat
```

### CI/CD Monitoring
```bash
4                    # Workflows
r                    # Refresh
↓↓↓                  # Scan for failures
Enter                # View failed workflow
l                    # Check logs
r                    # Rerun if needed
```

## 📱 GitHub CLI Integration

### Direct gh Commands
```bash
# gh-tui uses gh CLI under the hood
# You can also use gh directly:

gh pr list           # List PRs
gh pr view 123       # View PR
gh pr checkout 123   # Checkout PR
gh issue list        # List issues
gh repo view         # View repo
gh workflow list     # List workflows
```

### Authentication
```bash
# gh-tui requires gh auth
gh auth login        # First time setup
gh auth status       # Check auth
gh auth refresh      # Refresh token
```

## ⚙️ Configuration

### Settings (Coming Soon)
```
s - Settings
t - Change theme
k - Customize keybindings
```

### Cache
```
Ctrl+R - Clear cache and refresh
```

## 🛠️ Troubleshooting

### Refresh Issues
```
R - Force full refresh
r - Refresh current view
Ctrl+C, then restart - Hard reset
```

### Authentication Issues
```bash
# Outside gh-tui:
gh auth login        # Re-authenticate
gh auth status       # Verify
```

### Performance
```
v - Toggle to list-only view
# Reduces detail panel updates
```

## ⌨️ Quick Reference Card

```
Views:          1-5        PR/Issue/Repo/Workflow/Gist
Navigation:     hjkl       Vim-style
                g/G        Top/bottom
                /          Search

Actions:        Enter      View details
                o          Open in browser
                r          Refresh/Review/Rerun
                c          Close/Clone/Checkout

Global:         q          Quit
                ?          Help
                R          Force refresh
```

## 📊 Status Indicators

### Pull Request States
```
🟢 Open           Ready for review
🟣 Draft          Work in progress
✅ Merged         Successfully merged
🔴 Closed         Closed without merge
🔥 Conflicts      Merge conflicts present
👁️ Reviewed       Changes requested
✓ Approved        Approved for merge
```

### Workflow States
```
✅ Success        Workflow passed
❌ Failed         Workflow failed
🟡 Running        In progress
⏸️ Queued         Waiting to start
🔵 Pending        Approval required
```

### Repository Info
```
⭐ Stars          Popularity
🍴 Forks          Community interest
📦 Private        Access restricted
🌍 Public         Open source
```

## 🎯 Tips & Tricks

### Keyboard-Only Flow
```bash
# Never touch mouse:
1-5                  # Jump to views
hjkl                 # Navigate
Enter                # Select
o                    # Browser when needed
q                    # Quit
```

### Monitor CI/CD
```bash
# Keep gh-tui open in tmux pane
4                    # Workflows view
r                    # Auto-refresh (planned)
# Watch for failures in real-time
```

### Quick Issue Triage
```bash
2                    # Issues
/label:bug           # Filter bugs
↓                    # Review each
a                    # Assign
l                    # Add labels
```

---

**Version**: gh-tui v1.0
**Last Updated**: 2024-11-02
**Requires**: GitHub CLI (`gh`)
**Platform**: Linux | macOS | WSL | Termux
