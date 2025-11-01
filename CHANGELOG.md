# Changelog

All notable changes to gh-tui will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased] - 2025-11-01

### 🎉 Recent Updates

**Interactive Actions - Phase 1 Progress** - Adding user actions across all views
**Help System Complete** - Comprehensive keyboard shortcut reference

### ✨ Added

#### Help Screen (`?` key)
- **Interactive Help Overlay** - Press `?` to toggle comprehensive help screen
  - Beautiful centered modal with GitHub theme
  - Organized by sections: Global Keys, Navigation, Per-Tab Actions
  - Shows all current keybindings and coming soon features
  - Press `?` or `Esc` to close
- **Complete Keyboard Reference**
  - Global keys (quit, refresh, help)
  - Navigation (tab switching, arrow keys, vim keys)
  - Per-tab actions for all 5 tabs
  - Clear visual hierarchy with colors
- **Improved Discoverability** - New users can press `?` to learn all shortcuts

#### Interactive Browser Integration (`b` key)
- **Open in Browser** - Universal `b` key across all tabs
  - Pull Requests: Open PR in browser
  - Issues: Open issue in browser
  - Repositories: Open repo in browser
  - Actions: Open workflow run in browser
  - Gists: Open gist in browser (note: `o` opens in micro for viewing)
- **Helper Function** - New `openInBrowser()` utility in helpers.go
  - Supports all GitHub item types (pr, issue, repo, run, gist)
  - Uses `gh` CLI with `--web` flag
  - Proper error handling with `errMsg` type
- **Updated Help Text** - All views show `b: Browser` in keyboard hints

### 🔧 Technical Details

**Files Modified (Help Screen):**
- `update_keyboard.go` - Fixed `toggleHelp()` to actually toggle `showHelp` boolean, added Esc handler
- `view.go` - Added `renderHelpScreen()` function with comprehensive keybinding list
- `styles.go` - Added `helpSectionStyle` and `helpKeyStyle` for GitHub-themed formatting

**Files Modified (Browser Integration):**
- `helpers.go` - Added `openInBrowser()` function
- `view_pullrequests.go` - Added `b` key handler
- `view_issues.go` - Added `b` key handler
- `view_repositories.go` - Added `b` key handler
- `view_actions.go` - Added `b` key handler, fixed `DatabaseId` casing
- `view_gists.go` - Added `b` key handler

**Implementation:**
- Consistent key binding across all views
- Opens default browser to GitHub web interface
- Works with current context (no repo selection needed)
- Builds successfully with no errors

---

## [0.1.0 - MVP Complete] - 2025-11-01

### 🎉 Major Milestones

**Core MVP Complete** - All foundational components built and working

### ✨ Added

#### Foundation & Architecture
- **Project Scaffolding** - Complete Bubbletea-based TUI application structure
- **Tabbed Interface** - 5-tab navigation system (Pull Requests, Issues, Repositories, Actions, Gists)
- **Configuration System** - YAML-based config with sensible defaults
- **GitHub CLI Integration** - Complete backend using `gh` CLI for all API calls
- **Dual-Pane Layout** - List view + detail view for all tabs
- **Keyboard Navigation** - Full vim-style navigation (hjkl) + arrow keys + Tab switching
- **Mouse Support** - Click to select, scroll, and interact with UI elements

#### UI Components
- **Animated Landing Page** - Stunning wavy grid animation with GitHub dark theme
  - Ocean-themed metaballs flowing across screen
  - Smooth 20fps animation
  - Auto-dismisses after 3 seconds or on any key press
- **Table View Component** - Reusable sortable table with:
  - Sortable headers (click to sort)
  - Visual sort indicators (▲▼)
  - Responsive column widths
  - Proper overflow handling
- **Tree View Component** - Expandable/collapsible tree renderer with:
  - Tree symbols (├─ └─ │ ▶ ▼)
  - Multi-level expansion tracking
  - Keyboard navigation (→ expand, ← collapse)
  - Foundation for file/folder hierarchies

#### View Implementations

**Pull Requests View** (`view_pullrequests.go`)
- List all PRs with status indicators (✓, ✗, •)
- Show PR number, title, author, date
- Display review status and checks
- Detail pane with full PR metadata
- Navigate with arrow keys

**Issues View** (`view_issues.go`)
- List all issues with state indicators
- Show issue number, title, labels, assignees
- Filter by state (open/closed)
- Detail pane with issue body
- Labels displayed with colors

**Repositories View** (`view_repositories.go`)
- List repos with stats (⭐ stars, 🍴 forks)
- Show language, visibility, description
- Sort by name, stars, update date
- Detail pane with README preview
- Visual indicators for public/private

**Actions View** (`view_actions.go`)
- List workflow runs with status (success/failure/running)
- Show workflow name, branch, conclusion
- Color-coded status indicators
- Time since last run
- Filter by status

**Gists View** (`view_gists.go`)
- List all gists (public + private)
- Show description, file count, visibility
- Display creation/update timestamps
- Detail pane with gist content preview
- Integration with micro editor for editing

#### Styling & Theming
- **GitHub Dark Theme** - Carefully matched to GitHub's official color palette
  - Primary: #58a6ff (GitHub blue for links)
  - Success: #56d364 (GitHub green)
  - Error: #f85149 (GitHub red)
  - Background: #0d1117 (GitHub dark)
  - Borders: #21262d (GitHub border)
- **Status Indicators** - Emoji-based icons for PR/issue states
- **Responsive Layout** - Adapts to terminal size (min 40x10)
- **Lipgloss Styling** - Consistent, beautiful rendering across all views

#### Developer Experience
- **Helper Utilities** (`helpers.go`)
  - `truncateString()` - Smart text truncation with ellipsis
  - `formatTime()` - Human-readable timestamps
  - `formatTimeAgo()` - Relative time (2h ago, 3d ago)
  - `formatDuration()` - Time duration formatting
- **Type System** - Comprehensive Go types for all GitHub entities
- **Error Handling** - Graceful error messages and recovery
- **Config Management** - User config at `~/.config/gh-tui/config.yaml`

#### GitHub API Integration (`github.go`)
- `fetchPullRequests()` - Get PRs from any repo
- `fetchIssues()` - Get issues with full metadata
- `fetchRepositories()` - List user/org repos
- `fetchWorkflowRuns()` - Get GitHub Actions runs
- `fetchGists()` - List public + private gists
- `checkGitHubAuth()` - Verify `gh` CLI authentication
- All functions use async commands with proper error handling

#### Special Features
- **Gist Editor Integration** - Full gist management with micro text editor
  - View gists in read-only mode (`o` key)
  - Edit existing gists (`e` key)
  - Create new gists (`n` key)
  - Launch micro in-place
  - Auto-save and upload to GitHub
  - Syntax highlighting via micro
  - Clickable hyperlinks in gist descriptions
- **Dual-Pane Navigation** - Smooth transitions between list and detail
- **Tab Persistence** - Remembers which tab you were on
- **Loading States** - Shows loading indicators during data fetch
- **Empty States** - Helpful messages when no data available

### 🔧 Technical Details

#### Files Added
```
gh-tui/
├── main.go                 - Application entry point
├── types.go                - All type definitions
├── model.go                - Model initialization
├── update.go               - Message handling
├── update_keyboard.go      - Keyboard input handling
├── update_mouse.go         - Mouse input handling
├── view.go                 - Main view renderer
├── config.go               - Configuration system
├── github.go               - GitHub CLI integration
├── styles.go               - Lipgloss styles
├── helpers.go              - Utility functions
├── landing_page.go         - Animated startup screen
├── table_view.go           - Sortable table component
├── tree_view.go            - Tree view component
├── gist_editor.go          - Gist editing integration
├── view_pullrequests.go    - PR view implementation
├── view_issues.go          - Issues view implementation
├── view_repositories.go    - Repositories view implementation
├── view_actions.go         - Actions view implementation
└── view_gists.go           - Gists view implementation
```

#### Dependencies
- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/lipgloss` - Terminal styling
- `github.com/charmbracelet/bubbles` - UI components
- `gopkg.in/yaml.v3` - Config file parsing

#### Build System
- Go 1.21+ required
- Single binary compilation (`go build`)
- Cross-platform support (Linux, macOS, Windows)
- ~5MB binary size

### 📚 Documentation Added
- **README.md** - Complete installation and usage guide
  - Prerequisites and setup
  - Feature overview
  - Keyboard shortcuts
  - Installation instructions
  - Quick start guide
  - GitHub authentication setup

### 🐛 Bug Fixes

#### Data Handling
- Fixed API field name mismatches for workflow runs
- Corrected gist file structure to handle mixed types
- Fixed arrow key navigation with proper focus management
- Enabled data forwarding to views for proper display

#### UI/UX
- Fixed tab switching focus management
- Corrected window resize handling
- Fixed status bar overflow for long messages
- Improved error message display

### 🎨 UI Improvements
- GitHub-inspired color scheme throughout
- Consistent icon usage (⭐🍴🟢🔴✓✗•)
- Smooth animations in landing page
- Better contrast for readability
- Responsive layouts for different terminal sizes

---

## Development Session History

### Session 1 - 2025-10-31 (Foundation)
- Scaffolded project from TUITemplate
- Defined all GitHub data types
- Implemented tabbed view system
- Created GitHub CLI integration
- Added keyboard navigation
- Built all 5 view files
- Configured GitHub dark theme

### Session 2 - 2025-11-01 (Enhancement)
- Added stunning animated landing page
- Created tree view foundation (Phase 0)
- Integrated micro editor for Gists
- Added clickable hyperlinks
- Refined all view implementations
- Fixed navigation bugs
- Added helper utilities

---

## Git Commit History (Recent)

```
15c6bbe Add: Stunning animated landing page with GitHub dark theme
836006e Add: Tree view foundation (Phase 0)
e5dd626 Update plugin plan: Add tree view for Gists + Plugins tabs
88b8bb8 Add: Micro editor integration for Gists + clickable hyperlinks
39c96c2 Add integration plan for tkan and TFE features
0786bd9 Add comprehensive enhancement roadmap
f9f8479 Fix: Enable arrow key navigation with focus management
4eb06df Fix: Forward data loaded messages to views
1aa71c0 Fix: Correct gist files struct to handle mixed types
fb77fa7 Fix: Correct API field names for workflow runs and gists
```

---

## What's Next?

See [PLAN.md](PLAN.md) for the complete roadmap of upcoming features:

**Immediate Priorities:**
1. Interactive actions (open in browser, merge PR, etc.)
2. Help screen (`?` key)
3. Enhanced Gists with tree view for multi-file gists
4. Search and filter functionality

**Medium-Term:**
- Plugin manager tab
- Context menus and git operations
- Additional tabs (Notifications, Stars, Releases)

**Long-Term:**
- Kanban board integration (from tkan)
- Multi-repo workspace support
- Real-time webhook updates

---

## Acknowledgments

**Built with:**
- [Bubbletea](https://github.com/charmbracelet/bubbletea) - The amazing TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [GitHub CLI](https://cli.github.com/) - GitHub API access

**Inspired by:**
- [gh-dash](https://github.com/dlvhdr/gh-dash) - GitHub dashboard
- [lazygit](https://github.com/jesseduffield/lazygit) - Git TUI excellence
- GitHub's own dark theme color palette

**Related Projects:**
- [tkan](~/projects/tkan) - Kanban board with drag & drop
- [TFE](~/projects/TFE) - Tree view and context menus
- [TUITemplate](~/projects/TUITemplate) - Project scaffolding

---

## Statistics

**Timeline:** 2 development sessions (~8-10 hours total)
**Lines of Code:** ~2,200 Go source lines
**Files:** 18 Go source files + config/docs
**Features:** 5 views, 3 reusable components, 1 editor integration
**Status:** ✅ MVP Complete - Ready for Interactive Actions Phase
