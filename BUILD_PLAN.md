# gh-tui: GitHub CLI Interactive TUI - Build Prompt

**Project**: Interactive GitHub CLI wrapper using Bubbletea
**Template**: ~/projects/TUITemplate
**Reference**: github-api-integration-ideas.md (this directory)
**Timeline**: 3-4 weeks (MVP in 1-2 weeks)

---

## 🚀 CURRENT STATUS (Session 1 - 2025-10-31)

### ✅ COMPLETED (Core Architecture - 100%)
- ✅ Project scaffolded from TUITemplate
- ✅ All GitHub data types defined (PullRequest, Issue, Repository, WorkflowRun, Gist)
- ✅ ViewType enum and View interface implemented
- ✅ GitHub CLI integration (github.go) - all fetch functions complete
- ✅ Model initialization with tabbed view system
- ✅ Tabbed view rendering system in view.go
- ✅ Message handling for all GitHub data types in update.go
- ✅ Keyboard navigation (Tab switching, number keys 1-5, refresh)
- ✅ GitHub-inspired dark theme styling
- ✅ Config updated to default to tabbed layout
- ✅ GitHub auth check in main.go

### 🔄 IN PROGRESS (View Implementations)
Subagents created these files but need integration:
- 🔄 view_pullrequests.go (created by subagent)
- 🔄 view_issues.go (created by subagent)
- 🔄 view_repositories.go (created by subagent)
- 🔄 view_actions.go (created by subagent)
- 🔄 view_gists.go (created by subagent)

### 📋 NEXT SESSION TODO
1. **Review & integrate subagent view files** - Check what they created
2. **Fix compilation errors** - Ensure all views implement View interface correctly
3. **Add helper functions** - truncateString, formatTime, etc.
4. **Build and test** - `go build` and run with real GitHub data
5. **Fix any runtime issues** - Adjust layouts, handle edge cases
6. **Polish UI** - Fine-tune spacing, colors, icons

### 🎯 READY TO TEST
Core architecture is complete! Once view files are integrated, the app should:
- Launch with GitHub auth check
- Show tabbed interface (PRs, Issues, Repos, Actions, Gists)
- Allow navigation with Tab/Shift+Tab or 1-5 keys
- Fetch and display GitHub data via `gh` CLI
- Refresh views with 'r' key

---

## Project Overview

Build **gh-tui**, an interactive terminal interface for the GitHub CLI (`gh`) that provides a menu-driven, keyboard-navigable interface similar to lazygit, but for GitHub operations (PRs, issues, repos, actions, etc.).

**Differentiation**:
- vs **gh-dash**: More comprehensive (covers repos, actions, releases, not just PRs/issues)
- vs **GitHub web**: Terminal-native, keyboard-driven, faster navigation
- vs **raw `gh` CLI**: Visual interface, no need to remember commands

**Target Users**: Developers who live in the terminal and want fast GitHub operations without context switching to browser.

---

## Architecture Plan

### Phase 1: Foundation (Days 1-3)
Use TUITemplate to scaffold the base application with tabbed multi-view layout.

### Phase 2: Core Views (Days 4-10)
Implement 5 primary views: Pull Requests, Issues, Repositories, Actions, Gists.

### Phase 3: Details & Actions (Days 11-18)
Add detail panes, interactive operations (create PR, merge, comment, etc.).

### Phase 4: Polish & Extras (Days 19-28)
Real-time updates, configuration, keyboard shortcuts, help system.

---

## Step-by-Step Implementation Plan

### Step 1: Scaffold Project from TUITemplate

**Location**: Create new directory `~/projects/gh-tui`

**Command**:
```bash
cd ~/projects/TUITemplate
./scripts/new_project.sh

# Interactive prompts:
# App name: gh-tui
# Title: gh-tui - GitHub CLI Interactive Interface
# Description: Interactive TUI for GitHub CLI operations
# Author: [Your Name]
# Layout: tabbed        # Use tabbed layout for multiple views
# Components: list,preview,input,dialog,status
```

**Generated Structure**:
```
gh-tui/
├── main.go              # Entry point
├── types.go             # Struct definitions
├── model.go             # Model initialization
├── update.go            # Message dispatcher
├── update_keyboard.go   # Keyboard handling
├── update_mouse.go      # Mouse handling
├── view.go              # View rendering
├── styles.go            # Lipgloss styles
├── config.go            # Configuration
├── go.mod
└── .config/
    └── gh-tui/
        └── config.yaml
```

---

### Step 2: Define Core Data Structures

**File**: `types.go`

**Key Structs**:

```go
package main

import (
    "time"
    tea "github.com/charmbracelet/bubbletea"
)

// Main model
type model struct {
    width, height    int
    activeView       ViewType
    views            map[ViewType]View

    // GitHub data
    pullRequests     []PullRequest
    issues           []Issue
    repositories     []Repository
    workflowRuns     []WorkflowRun
    gists            []Gist

    // UI state
    loading          bool
    error            string
    lastSync         time.Time

    // Config
    config           Config
}

// View types (tabs)
type ViewType int

const (
    ViewPullRequests ViewType = iota
    ViewIssues
    ViewRepositories
    ViewActions
    ViewGists
)

// View interface
type View interface {
    Update(tea.Msg) (View, tea.Cmd)
    View(width, height int) string
    Focus()
    Blur()
}

// GitHub data structures
type PullRequest struct {
    Number       int       `json:"number"`
    Title        string    `json:"title"`
    State        string    `json:"state"`
    Author       string    `json:"author"`
    CreatedAt    time.Time `json:"createdAt"`
    UpdatedAt    time.Time `json:"updatedAt"`
    Repository   string    `json:"repository"`
    HeadRef      string    `json:"headRef"`
    BaseRef      string    `json:"baseRef"`
    IsDraft      bool      `json:"isDraft"`
    ReviewStatus string    `json:"reviewStatus"` // APPROVED, CHANGES_REQUESTED, PENDING
    Mergeable    string    `json:"mergeable"`
    URL          string    `json:"url"`
}

type Issue struct {
    Number       int       `json:"number"`
    Title        string    `json:"title"`
    State        string    `json:"state"`
    Author       string    `json:"author"`
    CreatedAt    time.Time `json:"createdAt"`
    UpdatedAt    time.Time `json:"updatedAt"`
    Repository   string    `json:"repository"`
    Labels       []string  `json:"labels"`
    Assignees    []string  `json:"assignees"`
    Milestone    string    `json:"milestone"`
    URL          string    `json:"url"`
}

type Repository struct {
    Name         string `json:"name"`
    FullName     string `json:"fullName"`
    Description  string `json:"description"`
    Stars        int    `json:"stars"`
    Forks        int    `json:"forks"`
    OpenIssues   int    `json:"openIssues"`
    Language     string `json:"language"`
    Visibility   string `json:"visibility"` // public, private
    URL          string `json:"url"`
}

type WorkflowRun struct {
    ID           int64     `json:"id"`
    Name         string    `json:"name"`
    Status       string    `json:"status"`     // queued, in_progress, completed
    Conclusion   string    `json:"conclusion"` // success, failure, cancelled
    HeadBranch   string    `json:"headBranch"`
    HeadCommit   string    `json:"headCommit"`
    RunNumber    int       `json:"runNumber"`
    CreatedAt    time.Time `json:"createdAt"`
    Repository   string    `json:"repository"`
    URL          string    `json:"url"`
}

type Gist struct {
    ID          string    `json:"id"`
    Description string    `json:"description"`
    Public      bool      `json:"public"`
    Files       []string  `json:"files"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
    URL         string    `json:"url"`
}

// Messages for async operations
type prLoadedMsg struct {
    prs []PullRequest
    err error
}

type issuesLoadedMsg struct {
    issues []Issue
    err    error
}

type reposLoadedMsg struct {
    repos []Repository
    err   error
}

type workflowsLoadedMsg struct {
    runs []WorkflowRun
    err  error
}

type gistsLoadedMsg struct {
    gists []Gist
    err   error
}

// Config
type Config struct {
    DefaultRepo    string   `yaml:"default_repo"`
    DefaultOrg     string   `yaml:"default_org"`
    Theme          string   `yaml:"theme"`
    AutoRefresh    bool     `yaml:"auto_refresh"`
    RefreshInterval int     `yaml:"refresh_interval"` // seconds
    FavoriteRepos  []string `yaml:"favorite_repos"`
    KeyBindings    string   `yaml:"keybindings"` // default, vim, emacs
}
```

---

### Step 3: Implement GitHub CLI Integration

**File**: `github.go` (new file)

**Pattern**: Use `gh` CLI as backend (same as TKAN pattern)

```go
package main

import (
    "encoding/json"
    "fmt"
    "os/exec"
    "strings"
    tea "github.com/charmbracelet/bubbletea"
)

// Fetch pull requests using gh CLI
func fetchPullRequests(repo string) tea.Cmd {
    return func() tea.Msg {
        // gh pr list --json number,title,state,author,createdAt,headRefName,baseRefName
        cmd := exec.Command("gh", "pr", "list",
            "--repo", repo,
            "--json", "number,title,state,author,createdAt,updatedAt,headRefName,baseRefName,isDraft,reviewDecision,mergeable,url",
            "--limit", "100")

        output, err := cmd.Output()
        if err != nil {
            return prLoadedMsg{err: fmt.Errorf("gh pr list failed: %w", err)}
        }

        var prs []PullRequest
        if err := json.Unmarshal(output, &prs); err != nil {
            return prLoadedMsg{err: fmt.Errorf("parse error: %w", err)}
        }

        return prLoadedMsg{prs: prs}
    }
}

// Fetch issues using gh CLI
func fetchIssues(repo string) tea.Cmd {
    return func() tea.Msg {
        cmd := exec.Command("gh", "issue", "list",
            "--repo", repo,
            "--json", "number,title,state,author,createdAt,updatedAt,labels,assignees,milestone,url",
            "--limit", "100")

        output, err := cmd.Output()
        if err != nil {
            return issuesLoadedMsg{err: fmt.Errorf("gh issue list failed: %w", err)}
        }

        var issues []Issue
        if err := json.Unmarshal(output, &issues); err != nil {
            return issuesLoadedMsg{err: fmt.Errorf("parse error: %w", err)}
        }

        return issuesLoadedMsg{issues: issues}
    }
}

// Fetch repositories (user or org)
func fetchRepositories(owner string) tea.Cmd {
    return func() tea.Msg {
        cmd := exec.Command("gh", "repo", "list", owner,
            "--json", "name,nameWithOwner,description,stargazerCount,forkCount,openIssuesCount,primaryLanguage,visibility,url",
            "--limit", "100")

        output, err := cmd.Output()
        if err != nil {
            return reposLoadedMsg{err: fmt.Errorf("gh repo list failed: %w", err)}
        }

        var repos []Repository
        if err := json.Unmarshal(output, &repos); err != nil {
            return reposLoadedMsg{err: fmt.Errorf("parse error: %w", err)}
        }

        return reposLoadedMsg{repos: repos}
    }
}

// Fetch workflow runs
func fetchWorkflowRuns(repo string) tea.Cmd {
    return func() tea.Msg {
        cmd := exec.Command("gh", "run", "list",
            "--repo", repo,
            "--json", "databaseId,name,status,conclusion,headBranch,headSha,runNumber,createdAt,url",
            "--limit", "50")

        output, err := cmd.Output()
        if err != nil {
            return workflowsLoadedMsg{err: fmt.Errorf("gh run list failed: %w", err)}
        }

        var runs []WorkflowRun
        if err := json.Unmarshal(output, &runs); err != nil {
            return workflowsLoadedMsg{err: fmt.Errorf("parse error: %w", err)}
        }

        return workflowsLoadedMsg{runs: runs}
    }
}

// Fetch gists
func fetchGists() tea.Cmd {
    return func() tea.Msg {
        cmd := exec.Command("gh", "gist", "list",
            "--json", "id,description,public,files,createdAt,updatedAt,url",
            "--limit", "50")

        output, err := cmd.Output()
        if err != nil {
            return gistsLoadedMsg{err: fmt.Errorf("gh gist list failed: %w", err)}
        }

        var gists []Gist
        if err := json.Unmarshal(output, &gists); err != nil {
            return gistsLoadedMsg{err: fmt.Errorf("parse error: %w", err)}
        }

        return gistsLoadedMsg{gists: gists}
    }
}

// Check if gh CLI is authenticated
func checkGitHubAuth() error {
    cmd := exec.Command("gh", "auth", "status")
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("gh not authenticated. Run: gh auth login")
    }
    return nil
}
```

---

### Step 4: Implement Tabbed View System

**File**: `view.go`

**Layout**:
```
┌──────────────────────────────────────────────────────────────┐
│  gh-tui - GitHub CLI Interactive Interface                   │
│  [Pull Requests] [Issues] [Repos] [Actions] [Gists]          │ ← Tabs
├──────────────────────────────────────────────────────────────┤
│ ┌────────────────────────────────┬─────────────────────────┐ │
│ │  List View                     │  Detail View            │ │
│ │  (PRs/Issues/Repos/etc)        │  (Selected item details)│ │
│ │                                │                         │ │
│ │  [✓] #123 Fix bug              │  PR #123: Fix bug       │ │
│ │  [x] #124 Add feature          │  Author: user123        │ │
│ │  [•] #125 Update docs          │  Created: 2h ago        │ │
│ │                                │  Status: ✓ Approved     │ │
│ │  ↑/↓: Navigate                 │  Base: main             │ │
│ │  Enter: Details                │  Head: feature-branch   │ │
│ │  m: Merge                      │                         │ │
│ │  c: Comment                    │  [m] Merge              │ │
│ │  r: Refresh                    │  [c] Comment            │ │
│ │                                │  [o] Open in browser    │ │
│ └────────────────────────────────┴─────────────────────────┘ │
├──────────────────────────────────────────────────────────────┤
│  Tab: Switch view | r: Refresh | q: Quit | ?: Help          │ ← Status bar
└──────────────────────────────────────────────────────────────┘
```

**Implementation**:

```go
package main

import (
    "fmt"
    "strings"
    "github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
    if m.width == 0 {
        return "Loading..."
    }

    // Build layout sections
    header := m.renderHeader()
    tabs := m.renderTabs()
    content := m.renderContent()
    status := m.renderStatusBar()

    return lipgloss.JoinVertical(lipgloss.Left,
        header,
        tabs,
        content,
        status,
    )
}

func (m model) renderHeader() string {
    title := titleStyle.Render("gh-tui - GitHub CLI Interactive Interface")
    return title
}

func (m model) renderTabs() string {
    tabs := []string{
        "Pull Requests",
        "Issues",
        "Repositories",
        "Actions",
        "Gists",
    }

    var renderedTabs []string
    for i, tab := range tabs {
        if ViewType(i) == m.activeView {
            renderedTabs = append(renderedTabs, activeTabStyle.Render(tab))
        } else {
            renderedTabs = append(renderedTabs, inactiveTabStyle.Render(tab))
        }
    }

    return lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
}

func (m model) renderContent() string {
    // Delegate to active view
    view := m.views[m.activeView]
    if view == nil {
        return "View not implemented"
    }

    // Calculate available space
    contentWidth := m.width
    contentHeight := m.height - 5 // header + tabs + status

    return view.View(contentWidth, contentHeight)
}

func (m model) renderStatusBar() string {
    left := "Tab: Switch view | r: Refresh | q: Quit | ?: Help"
    right := fmt.Sprintf("Last sync: %s", m.lastSync.Format("15:04:05"))

    padding := m.width - lipgloss.Width(left) - lipgloss.Width(right)
    if padding < 0 {
        padding = 0
    }

    return statusBarStyle.Render(
        left + strings.Repeat(" ", padding) + right,
    )
}
```

---

### Step 5: Implement Pull Request View

**File**: `view_pullrequests.go` (new file)

```go
package main

import (
    "fmt"
    "strings"
    "github.com/charmbracelet/lipgloss"
)

type PullRequestView struct {
    prs          []PullRequest
    cursor       int
    selectedPR   *PullRequest
    showDetail   bool
    loading      bool
    error        string
}

func NewPullRequestView() *PullRequestView {
    return &PullRequestView{
        prs:        []PullRequest{},
        showDetail: true,
    }
}

func (v *PullRequestView) Update(msg tea.Msg) (View, tea.Cmd) {
    switch msg := msg.(type) {
    case prLoadedMsg:
        if msg.err != nil {
            v.error = msg.err.Error()
            v.loading = false
            return v, nil
        }
        v.prs = msg.prs
        v.loading = false
        if len(v.prs) > 0 {
            v.selectedPR = &v.prs[0]
        }
        return v, nil

    case tea.KeyMsg:
        switch msg.String() {
        case "up", "k":
            if v.cursor > 0 {
                v.cursor--
                v.selectedPR = &v.prs[v.cursor]
            }
        case "down", "j":
            if v.cursor < len(v.prs)-1 {
                v.cursor++
                v.selectedPR = &v.prs[v.cursor]
            }
        case "enter":
            v.showDetail = !v.showDetail
        case "o":
            // Open in browser
            if v.selectedPR != nil {
                exec.Command("gh", "pr", "view", "--web",
                    fmt.Sprintf("%d", v.selectedPR.Number)).Run()
            }
        case "m":
            // Merge PR
            if v.selectedPR != nil {
                return v, v.mergePR(v.selectedPR.Number)
            }
        }
    }

    return v, nil
}

func (v *PullRequestView) View(width, height int) string {
    if v.loading {
        return "Loading pull requests..."
    }

    if v.error != "" {
        return errorStyle.Render("Error: " + v.error)
    }

    if len(v.prs) == 0 {
        return "No pull requests found"
    }

    // Split view: list on left, detail on right
    listWidth := width / 2
    detailWidth := width - listWidth

    list := v.renderList(listWidth, height)
    detail := ""
    if v.showDetail && v.selectedPR != nil {
        detail = v.renderDetail(detailWidth, height)
    }

    return lipgloss.JoinHorizontal(lipgloss.Top, list, detail)
}

func (v *PullRequestView) renderList(width, height int) string {
    var lines []string

    // Title
    lines = append(lines, listTitleStyle.Render("Pull Requests"))
    lines = append(lines, "")

    // PR list
    visibleStart := v.cursor
    if visibleStart > height-5 {
        visibleStart = height - 5
    }

    for i := visibleStart; i < len(v.prs) && len(lines) < height-2; i++ {
        pr := v.prs[i]

        // Status icon
        icon := "•"
        if pr.ReviewStatus == "APPROVED" {
            icon = "✓"
        } else if pr.ReviewStatus == "CHANGES_REQUESTED" {
            icon = "✗"
        }

        // Line format: [icon] #123 Title
        line := fmt.Sprintf("%s #%-4d %s", icon, pr.Number, pr.Title)

        // Truncate if too long
        if len(line) > width-4 {
            line = line[:width-7] + "..."
        }

        // Highlight selected
        if i == v.cursor {
            line = selectedStyle.Render(line)
        }

        lines = append(lines, line)
    }

    // Help text
    lines = append(lines, "")
    lines = append(lines, helpStyle.Render("↑/↓: Navigate | Enter: Toggle detail | m: Merge | o: Browser"))

    return listPanelStyle.Width(width).Height(height).Render(
        strings.Join(lines, "\n"),
    )
}

func (v *PullRequestView) renderDetail(width, height int) string {
    if v.selectedPR == nil {
        return ""
    }

    pr := v.selectedPR

    var lines []string

    // Title
    lines = append(lines, detailTitleStyle.Render(
        fmt.Sprintf("PR #%d: %s", pr.Number, pr.Title),
    ))
    lines = append(lines, "")

    // Metadata
    lines = append(lines, fmt.Sprintf("Author:     %s", pr.Author))
    lines = append(lines, fmt.Sprintf("Created:    %s", pr.CreatedAt.Format("Jan 2, 15:04")))
    lines = append(lines, fmt.Sprintf("Updated:    %s", pr.UpdatedAt.Format("Jan 2, 15:04")))
    lines = append(lines, fmt.Sprintf("State:      %s", pr.State))
    lines = append(lines, fmt.Sprintf("Review:     %s", pr.ReviewStatus))
    lines = append(lines, fmt.Sprintf("Mergeable:  %s", pr.Mergeable))
    lines = append(lines, "")
    lines = append(lines, fmt.Sprintf("Base:       %s", pr.BaseRef))
    lines = append(lines, fmt.Sprintf("Head:       %s", pr.HeadRef))
    lines = append(lines, "")

    // Actions
    lines = append(lines, actionStyle.Render("[m] Merge PR"))
    lines = append(lines, actionStyle.Render("[c] Add comment"))
    lines = append(lines, actionStyle.Render("[o] Open in browser"))
    lines = append(lines, actionStyle.Render("[r] Request review"))

    return detailPanelStyle.Width(width).Height(height).Render(
        strings.Join(lines, "\n"),
    )
}

func (v *PullRequestView) mergePR(number int) tea.Cmd {
    return func() tea.Msg {
        cmd := exec.Command("gh", "pr", "merge", fmt.Sprintf("%d", number), "--merge")
        if err := cmd.Run(); err != nil {
            return errorMsg{err}
        }
        return successMsg{"PR merged successfully"}
    }
}

func (v *PullRequestView) Focus() {}
func (v *PullRequestView) Blur() {}
```

---

### Step 6: Implement Issues, Repos, Actions, Gists Views

**Pattern**: Follow the same structure as PullRequestView

**Files**:
- `view_issues.go` - Issue list with filters (state, labels, assignees)
- `view_repositories.go` - Repo list with stats (stars, forks, language)
- `view_actions.go` - Workflow runs with status indicators
- `view_gists.go` - Gist list with file counts

**Each view implements**:
```go
type View interface {
    Update(tea.Msg) (View, tea.Cmd)
    View(width, height int) string
    Focus()
    Blur()
}
```

---

### Step 7: Implement Keyboard Navigation

**File**: `update_keyboard.go`

```go
package main

import tea "github.com/charmbracelet/bubbletea"

func (m model) handleKeyboard(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    // Global keybindings
    switch msg.String() {
    case "q", "ctrl+c":
        return m, tea.Quit

    case "tab":
        // Switch to next view
        m.activeView = (m.activeView + 1) % 5
        return m, nil

    case "shift+tab":
        // Switch to previous view
        if m.activeView == 0 {
            m.activeView = 4
        } else {
            m.activeView--
        }
        return m, nil

    case "r":
        // Refresh current view data
        return m, m.refreshActiveView()

    case "?":
        // Show help
        m.showHelp = !m.showHelp
        return m, nil

    case "1":
        m.activeView = ViewPullRequests
        return m, nil
    case "2":
        m.activeView = ViewIssues
        return m, nil
    case "3":
        m.activeView = ViewRepositories
        return m, nil
    case "4":
        m.activeView = ViewActions
        return m, nil
    case "5":
        m.activeView = ViewGists
        return m, nil
    }

    // Delegate to active view
    view := m.views[m.activeView]
    if view != nil {
        updatedView, cmd := view.Update(msg)
        m.views[m.activeView] = updatedView
        return m, cmd
    }

    return m, nil
}

func (m model) refreshActiveView() tea.Cmd {
    repo := m.config.DefaultRepo

    switch m.activeView {
    case ViewPullRequests:
        return fetchPullRequests(repo)
    case ViewIssues:
        return fetchIssues(repo)
    case ViewRepositories:
        return fetchRepositories(m.config.DefaultOrg)
    case ViewActions:
        return fetchWorkflowRuns(repo)
    case ViewGists:
        return fetchGists()
    }

    return nil
}
```

---

### Step 8: Add Lipgloss Styles

**File**: `styles.go`

```go
package main

import "github.com/charmbracelet/lipgloss"

var (
    // Colors
    primaryColor   = lipgloss.Color("#61AFEF")
    secondaryColor = lipgloss.Color("#C678DD")
    accentColor    = lipgloss.Color("#98C379")
    errorColor     = lipgloss.Color("#E06C75")

    // Title bar
    titleStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(primaryColor).
        Padding(1, 2)

    // Tabs
    activeTabStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("#000")).
        Background(primaryColor).
        Padding(0, 2)

    inactiveTabStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#666")).
        Padding(0, 2)

    // Panels
    listPanelStyle = lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(primaryColor).
        Padding(1, 2)

    detailPanelStyle = lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(secondaryColor).
        Padding(1, 2)

    // List items
    selectedStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(accentColor).
        Background(lipgloss.Color("#1E1E1E"))

    // Text styles
    listTitleStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(primaryColor)

    detailTitleStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(secondaryColor)

    helpStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#555"))

    actionStyle = lipgloss.NewStyle().
        Foreground(accentColor)

    errorStyle = lipgloss.NewStyle().
        Foreground(errorColor).
        Bold(true)

    // Status bar
    statusBarStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#999")).
        Background(lipgloss.Color("#1A1A1A")).
        Padding(0, 1)
)
```

---

### Step 9: Configuration System

**File**: `config.go`

```go
package main

import (
    "os"
    "path/filepath"
    "gopkg.in/yaml.v3"
)

func loadConfig() (Config, error) {
    // Default config
    cfg := Config{
        DefaultRepo:     "",  // Will prompt user
        DefaultOrg:      "",  // Will prompt user
        Theme:           "dark",
        AutoRefresh:     false,
        RefreshInterval: 60,
        FavoriteRepos:   []string{},
        KeyBindings:     "default",
    }

    // Config file path
    configDir := filepath.Join(os.Getenv("HOME"), ".config", "gh-tui")
    configPath := filepath.Join(configDir, "config.yaml")

    // Create config dir if not exists
    if err := os.MkdirAll(configDir, 0755); err != nil {
        return cfg, err
    }

    // Read config file
    data, err := os.ReadFile(configPath)
    if err != nil {
        // Config doesn't exist, create default
        if os.IsNotExist(err) {
            data, _ := yaml.Marshal(cfg)
            os.WriteFile(configPath, data, 0644)
        }
        return cfg, nil
    }

    // Parse config
    if err := yaml.Unmarshal(data, &cfg); err != nil {
        return cfg, err
    }

    return cfg, nil
}

func saveConfig(cfg Config) error {
    configDir := filepath.Join(os.Getenv("HOME"), ".config", "gh-tui")
    configPath := filepath.Join(configDir, "config.yaml")

    data, err := yaml.Marshal(cfg)
    if err != nil {
        return err
    }

    return os.WriteFile(configPath, data, 0644)
}
```

**Default config file** (`~/.config/gh-tui/config.yaml`):
```yaml
default_repo: "owner/repo"
default_org: "your-org"
theme: "dark"
auto_refresh: false
refresh_interval: 60
favorite_repos:
  - "owner/repo1"
  - "owner/repo2"
keybindings: "default"
```

---

### Step 10: Main Entry Point

**File**: `main.go`

```go
package main

import (
    "fmt"
    "os"
    tea "github.com/charmbracelet/bubbletea"
)

func main() {
    // Check gh authentication
    if err := checkGitHubAuth(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        fmt.Fprintf(os.Stderr, "\nPlease run: gh auth login\n")
        os.Exit(1)
    }

    // Load config
    cfg, err := loadConfig()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
        os.Exit(1)
    }

    // Initialize model
    m := newModel(cfg)

    // Create program
    p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())

    // Run
    if _, err := p.Run(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}

func newModel(cfg Config) model {
    m := model{
        config:     cfg,
        activeView: ViewPullRequests,
        views:      make(map[ViewType]View),
    }

    // Initialize views
    m.views[ViewPullRequests] = NewPullRequestView()
    m.views[ViewIssues] = NewIssueView()
    m.views[ViewRepositories] = NewRepositoryView()
    m.views[ViewActions] = NewActionsView()
    m.views[ViewGists] = NewGistView()

    return m
}

func (m model) Init() tea.Cmd {
    return tea.Batch(
        fetchPullRequests(m.config.DefaultRepo),
        fetchIssues(m.config.DefaultRepo),
    )
}
```

---

## Testing Plan

### Manual Testing Checklist

**Phase 1: Basic Navigation**
- [ ] App launches without errors
- [ ] Tabs switch correctly (Tab/Shift+Tab/1-5)
- [ ] Arrow keys navigate lists
- [ ] Selected item highlights properly

**Phase 2: GitHub Integration**
- [ ] PRs load correctly
- [ ] Issues load correctly
- [ ] Repos load correctly
- [ ] Actions load correctly
- [ ] Gists load correctly

**Phase 3: Actions**
- [ ] Open PR in browser works (o key)
- [ ] Merge PR works (m key)
- [ ] Create issue works
- [ ] Trigger workflow works

**Phase 4: UI**
- [ ] Layout renders correctly at 80x24
- [ ] Layout renders correctly at 120x40
- [ ] Detail pane toggles correctly
- [ ] Status bar shows correct info
- [ ] Error messages display properly

---

## MVP Feature Checklist

**Must Have (Week 1)** - ✅ CORE DONE (Need view integration):
- [x] Tabbed interface (5 views) ✅ DONE
- [x] Pull Requests view (list + detail) 🔄 CREATED (needs integration)
- [x] Issues view (list + detail) 🔄 CREATED (needs integration)
- [x] Keyboard navigation (arrows, Tab, number keys) ✅ DONE
- [x] GitHub CLI integration (fetch data) ✅ DONE
- [x] Basic styling (Lipgloss) ✅ DONE

**Should Have (Week 2)**:
- [x] Repositories view 🔄 CREATED (needs integration)
- [x] Actions/Workflows view 🔄 CREATED (needs integration)
- [x] Gists view 🔄 CREATED (needs integration)
- [ ] Interactive operations (merge PR, create issue) - NEXT
- [x] Configuration system ✅ DONE
- [ ] Help screen - NEXT

**Nice to Have (Week 3-4)**:
- [ ] Auto-refresh
- [ ] Favorites system
- [ ] Filters (state, labels, assignees)
- [ ] Search functionality
- [ ] Mouse support
- [ ] Color themes

---

## File Structure Summary

```
gh-tui/
├── main.go                  # Entry point (30 lines)
├── types.go                 # Structs and types (150 lines)
├── model.go                 # Model init and layout (100 lines)
├── update.go                # Message dispatcher (80 lines)
├── update_keyboard.go       # Keyboard handling (150 lines)
├── view.go                  # Main view rendering (100 lines)
├── view_pullrequests.go     # PR view (250 lines)
├── view_issues.go           # Issues view (250 lines)
├── view_repositories.go     # Repos view (200 lines)
├── view_actions.go          # Actions view (200 lines)
├── view_gists.go            # Gists view (150 lines)
├── github.go                # GitHub CLI integration (200 lines)
├── styles.go                # Lipgloss styles (100 lines)
├── config.go                # Config management (100 lines)
├── go.mod
├── go.sum
└── README.md
```

**Total**: ~1,900 lines for MVP

---

## Dependencies

**Required**:
```go
github.com/charmbracelet/bubbletea
github.com/charmbracelet/lipgloss
github.com/charmbracelet/bubbles
gopkg.in/yaml.v3
```

**Optional** (add as needed):
```go
github.com/charmbracelet/glamour       // Markdown rendering for PR/issue bodies
github.com/alecthomas/chroma/v2        // Syntax highlighting for code diffs
```

---

## Usage Examples

**Launch**:
```bash
gh-tui
```

**With specific repo**:
```bash
gh-tui --repo owner/repo
```

**Keyboard shortcuts**:
- `Tab` / `Shift+Tab` - Switch views
- `1-5` - Jump to specific view
- `↑/↓` or `j/k` - Navigate list
- `Enter` - Toggle detail pane
- `r` - Refresh current view
- `m` - Merge PR (in PR view)
- `o` - Open in browser
- `q` - Quit
- `?` - Help

---

## Integration with Opustrator

**As a TUI Tool** in `spawn-options.json`:

```json
{
  "label": "GitHub TUI",
  "command": "gh-tui",
  "terminalType": "tui-tool",
  "icon": "🐙",
  "description": "Interactive GitHub CLI interface"
}
```

Right-click folder in Opustrator → "GitHub TUI" → Opens gh-tui in canvas terminal.

---

## Next Steps After MVP

1. **Add to Opustrator spawn options** - Test integration
2. **Build distributable binary** - `go build -o gh-tui`
3. **Add to PATH** - Install system-wide
4. **GitHub release** - Share with community
5. **User feedback** - Iterate on features

---

## Success Criteria

✅ **Launch**: App opens without errors
✅ **Navigate**: Tab/arrow keys work smoothly
✅ **GitHub**: Pulls PRs, issues, repos correctly via `gh` CLI
✅ **Display**: Data renders cleanly in dual-pane layout
✅ **Actions**: Can merge PR, open in browser, create issue
✅ **Config**: Remembers default repo/org
✅ **Polish**: Looks professional (Lipgloss styling)

---

**Build this now using the Bubbletea skill and TUITemplate!**
**Estimated build time: 1-2 weeks for MVP, 3-4 weeks for polished v1.0**

---

Ready to start? Run:
```bash
cd ~/projects/TUITemplate
./scripts/new_project.sh
# App name: gh-tui
# Layout: tabbed
# Components: list,preview,input,dialog,status
```

Then follow this prompt step-by-step!
