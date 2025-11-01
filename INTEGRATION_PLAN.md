# Integration Plan: tkan + TFE → gh-tui

## 🎯 Executive Summary

Integrate proven features from your existing projects:
- **tkan**: Kanban board with GitHub Projects API backend + drag & drop
- **TFE**: Table views with sortable headers + git operations context menus

This will transform gh-tui from a read-only viewer into a **powerful GitHub workflow hub**.

---

## 📦 What We're Bringing Over

### From tkan (~/projects/tkan)

**Core Features:**
- ✅ Kanban board view (BACKLOG, TODO, IN PROGRESS, REVIEW, DONE)
- ✅ Solitaire-style card stacking
- ✅ Mouse drag & drop between columns
- ✅ Visual drop indicators (green line)
- ✅ Ghost cards during drag
- ✅ Detail panel (33% width, toggleable)
- ✅ GitHub Projects API integration (backend_github.go)

**Files to Port:**
```
tkan/backend_github.go     → gh-tui/github_projects.go
tkan/update_mouse.go        → Enhance gh-tui/update_mouse.go
tkan/view.go (board render) → gh-tui/view_projects.go
```

### From TFE (~/projects/TFE)

**Core Features:**
- ✅ Sortable table headers (click to sort)
- ✅ Context menus with right-click
- ✅ Git operations (pull, push, sync, fetch)
- ✅ Multiple view modes (list vs table)

**Files to Port:**
```
TFE/context_menu.go           → gh-tui/context_menu.go
TFE/render_file_list.go       → Patterns for gh-tui table views
Git operations logic           → gh-tui/git_operations.go
```

---

## 🏗️ New Architecture

### Tab Structure (Enhanced)

```
┌─────────────────────────────────────────────────────────────┐
│  [1] PRs  [2] Issues  [3] Repos  [4] Actions  [5] Gists    │
│  [6] Projects  [7] Stars  [8] Notifications                 │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  View content (list, table, or kanban based on tab)        │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### View Modes Per Tab

| Tab | Views Available | Default |
|-----|----------------|---------|
| PRs | List, Table, Detail | List |
| Issues | List, Table, Detail | List |
| **Repos** | **List, Table** | **Table** |
| Actions | List, Table | List |
| Gists | List, Table | List |
| **Projects** | **Kanban, Table** | **Kanban** |
| Stars | List, Table | List |
| Notifications | List, Table | List |

**Toggle view mode**: Press `v` to cycle through available views

---

## 🎨 Tab 6: Projects (Kanban Board)

### Integration Strategy

**New Tab Structure:**
```go
// view_projects.go
type ProjectsView struct {
    // Kanban state
    board          *Board
    columns        []Column
    selectedColumn int
    selectedCard   int
    showDetails    bool

    // Table view state
    viewMode       ViewMode  // Kanban or Table
    sortColumn     string
    sortAscending  bool

    // Mouse drag state (from tkan)
    draggingCard   *Card
    dragFromColumn int
    dropTarget     DropTarget

    // Data
    data    []GitHubProject
    loading bool
    err     error
}
```

### Kanban View Layout

```
┌─────────────────────────────────────────────────────┬──────────────────┐
│ 📋 BACKLOG │ ✅ TODO │ 🚧 IN PROGRESS │ 👀 REVIEW │ ✅ DONE │ 📦 ARCHIVE   │
│    (5)     │   (8)   │      (12)       │    (3)    │  (23)   │            │
│                                                                             │
│ ┌─────────┐ ┌──────  ┌─────────────┐  ┌────────  ┌────────              │
│ │ Setup    │ │ Fix     │ Refactor    │  │ Review   │ Deploy               │
│ │ auth     │ │ login   │ API layer   │  │ PR #42   │ v2.0                 │
│ ┌─────────┐ ┌──────  ┌─────────────┐  ┌────────                          │
│ │ Add      │ │ Write   │ Add tests   │  │ Test     │  👈 Drag & drop     │
│ │ docs     │ │ tests   │ for API     │  │ deploy   │     with mouse!     │
│ └─────────┘ └──────  └─────────────┘  └────────                          │
│                       │             │                                      │
│                       │ 🐛 Bug fix  │                                      │
│                       │             │                                      │
│                       └─────────────┘                                      │
├─────────────────────────────────────────────────────┴──────────────────────┤
│ Detail Panel: Refactor API layer                                           │
│                                                                             │
│ Description: Split monolithic API into microservices                       │
│ Assignee: @alice                                                            │
│ Labels: enhancement, backend                                                │
│ Updated: 2h ago                                                             │
│                                                                             │
│ [d]rag [e]dit [m]ove [o]pen [Tab] toggle details                          │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Table View Layout

```
┌─────────────────────────────────────────────────────────────────────┐
│ Projects (Table View)                             [v] Switch to Board │
├──────────┬─────────────────────────┬───────┬──────────┬─────────────┤
│ Status ▲ │ Title                   │ Owner │ Updated  │ Items       │
├──────────┼─────────────────────────┼───────┼──────────┼─────────────┤
│ 🚧 Active │ Backend Refactor v2     │ alice │ 2h ago   │ 12 open     │
│ 🚧 Active │ Frontend Migration      │ bob   │ 5h ago   │ 8 open      │
│ 📦 Done   │ Q4 Releases            │ eve   │ 1d ago   │ 0 open      │
│ ⏸️ Hold   │ Mobile App MVP          │ dave  │ 3d ago   │ 5 open      │
└──────────┴─────────────────────────┴───────┴──────────┴─────────────┘

Click headers to sort • Right-click for menu
```

### Key Bindings

**Kanban Mode:**
- `←/→` or `h/l` - Navigate columns
- `↑/↓` or `k/j` - Navigate cards
- `Tab` - Toggle detail panel
- `a` - Toggle archive column
- `v` - Switch to table view
- `n` - New project/card
- `e` - Edit selected card
- `o` - Open project in browser
- `d` - Delete card (with confirmation)
- **Mouse**: Drag & drop cards between columns

**Table Mode:**
- `↑/↓` or `k/j` - Navigate rows
- `←/→` or `h/l` - Navigate columns
- Click header - Sort by column
- `v` - Switch to kanban view
- `o` - Open selected project
- Right-click - Context menu

---

## 📊 Tab 3: Repos (Enhanced with Table View)

### Integration from TFE

**Current (List Only):**
```
Repositories (5)

▶ GGPrompts/gh-tui        ⭐ 0  🍴 0  Go
  GGPrompts/tkan          ⭐ 2  🍴 1  Go
  GGPrompts/TFE           ⭐ 5  🍴 0  Go
```

**Enhanced (Table View):**
```
┌────────────────────────────────────────────────────────────────────────┐
│ Repositories (Table View)                          [v] Switch to List  │
├──────────────────────────┬──────┬──────┬──────────┬──────────┬─────────┤
│ Name ▼                   │ ⭐   │ 🍴   │ Language │ Updated  │ Vis     │
├──────────────────────────┼──────┼──────┼──────────┼──────────┼─────────┤
│ GGPrompts/gh-tui        │ 0    │ 0    │ Go       │ 1h ago   │ 🌐 Pub  │
│ GGPrompts/tkan          │ 2    │ 1    │ Go       │ 2d ago   │ 🌐 Pub  │
│ GGPrompts/TFE           │ 5    │ 0    │ Go       │ 5d ago   │ 🌐 Pub  │
└──────────────────────────┴──────┴──────┴──────────┴──────────┴─────────┘

Click headers to sort • Right-click for menu
```

### Context Menu (Right-Click)

```
┌──────────────────────────┐
│ GGPrompts/gh-tui         │
├──────────────────────────┤
│ 🌐 Open in browser      │
│ 📋 Copy URL             │
│ 🍴 Fork                  │
│ ⭐ Star/Unstar          │
│ ──────────────────────  │
│ Git Operations          │
│   ↓ Pull                │
│   ↑ Push                │
│   🔄 Sync (Pull + Push) │
│   🔍 Fetch              │
│ ──────────────────────  │
│ 🌿 Open in lazygit      │
│ 📁 Open in terminal     │
│ ──────────────────────  │
│ ❌ Cancel               │
└──────────────────────────┘
```

### Key Features

1. **Sortable Headers**: Click any column header to sort
2. **Multi-sort**: Shift+click for secondary sort
3. **Git Operations**: Built-in git commands (no external tools needed)
4. **Context Menus**: Right-click on any repo
5. **View Toggle**: Press `v` to switch between list and table

---

## 🛠️ Implementation Plan

### Phase 1: Projects Tab (Kanban) - Week 1

**Day 1-2: GitHub Projects API Integration**
```bash
# Port tkan's backend_github.go
cp ~/projects/tkan/backend_github.go ~/projects/gh-tui/github_projects.go

# Adapt to gh-tui structure
- Change package name
- Integrate with existing github.go
- Add fetch functions for projects
```

**Day 3-4: Kanban View**
```bash
# Create new view file
touch ~/projects/gh-tui/view_projects_kanban.go

# Port from tkan:
- Board rendering logic
- Card stacking system
- Column navigation
- Detail panel
```

**Day 5-7: Mouse Drag & Drop**
```bash
# Enhance update_mouse.go
- Add drag detection
- Ghost card rendering
- Drop indicators
- Column drop logic
```

**Files Created:**
- `github_projects.go` - API integration
- `view_projects.go` - Main projects view
- `view_projects_kanban.go` - Kanban rendering
- `view_projects_table.go` - Table rendering (Phase 2)

### Phase 2: Table View + Context Menus - Week 2

**Day 1-2: Table View Foundation**
```bash
# Create table view system
touch ~/projects/gh-tui/table_view.go

# Add sortable headers
- Click detection on headers
- Sort state management
- Visual sort indicators (▲▼)
```

**Day 3-4: Context Menu System**
```bash
# Port from TFE
cp ~/projects/TFE/context_menu.go ~/projects/gh-tui/context_menu.go

# Adapt to gh-tui:
- Generic context menu component
- Per-view menu configurations
- Keyboard + mouse navigation
```

**Day 5-7: Git Operations**
```bash
# Create git integration
touch ~/projects/gh-tui/git_operations.go

# Implement:
- Pull, push, sync, fetch
- Status checking
- Error handling
- Progress indicators
```

**Files Created:**
- `table_view.go` - Reusable table component
- `context_menu.go` - Context menu system
- `git_operations.go` - Git command wrappers

### Phase 3: Apply to All Tabs - Week 3

**Convert Each Tab:**
1. **Repos** - Add table view + git ops (highest priority)
2. **PRs** - Add table view + context menu
3. **Issues** - Add table view + context menu
4. **Actions** - Add table view
5. **Gists** - Add table view

**Pattern for Each View:**
```go
// view_repositories.go (enhanced)
type RepositoryView struct {
    // Existing fields
    data   []Repository
    cursor int

    // New fields
    viewMode      ViewMode      // List or Table
    tableState    *TableState   // Sortable table
    contextMenu   *ContextMenu  // Right-click menu
}

// Implement both render methods
func (v *RepositoryView) renderList() string { ... }
func (v *RepositoryView) renderTable() string { ... }
```

---

## 📝 Code Patterns

### 1. View Mode Toggle

```go
// In update_keyboard.go
case "v":
    if view, ok := m.views[m.activeView].(ViewModeToggler); ok {
        view.ToggleViewMode()
    }
    return m, nil
```

### 2. Sortable Table Headers

```go
// table_view.go
type TableState struct {
    headers       []Header
    sortColumn    int
    sortAscending bool
    data          [][]string
}

func (t *TableState) HandleHeaderClick(x, y int) {
    column := t.getColumnFromPosition(x, y)
    if column == t.sortColumn {
        t.sortAscending = !t.sortAscending
    } else {
        t.sortColumn = column
        t.sortAscending = true
    }
    t.sort()
}

func (t *TableState) RenderHeaders() string {
    var headers []string
    for i, h := range t.headers {
        indicator := "  "
        if i == t.sortColumn {
            if t.sortAscending {
                indicator = " ▲"
            } else {
                indicator = " ▼"
            }
        }
        headers = append(headers, h.Name + indicator)
    }
    return strings.Join(headers, " │ ")
}
```

### 3. Context Menu

```go
// context_menu.go
type ContextMenu struct {
    items    []MenuItem
    selected int
    x, y     int
    visible  bool
}

type MenuItem struct {
    Label  string
    Action string
    Icon   string
}

func (c *ContextMenu) Show(x, y int, items []MenuItem) {
    c.x = x
    c.y = y
    c.items = items
    c.visible = true
    c.selected = 0
}

func (c *ContextMenu) HandleClick(x, y int) string {
    if !c.visible return ""

    index := c.getItemAtPosition(x, y)
    if index >= 0 && index < len(c.items) {
        action := c.items[index].Action
        c.Hide()
        return action
    }
    return ""
}
```

### 4. Git Operations

```go
// git_operations.go
func GitPull(repoPath string) tea.Cmd {
    return func() tea.Msg {
        cmd := exec.Command("git", "-C", repoPath, "pull")
        output, err := cmd.CombinedOutput()

        if err != nil {
            return gitOperationMsg{
                operation: "pull",
                success:   false,
                message:   string(output),
                err:       err,
            }
        }

        return gitOperationMsg{
            operation: "pull",
            success:   true,
            message:   "Successfully pulled changes",
        }
    }
}

func GitSync(repoPath string) tea.Cmd {
    return tea.Sequence(
        GitPull(repoPath),
        GitPush(repoPath),
    )
}
```

---

## 🎯 Benefits

### For Users

1. **Unified GitHub Workflow**
   - No need to switch between tools
   - One TUI for everything GitHub

2. **Visual Project Management**
   - Drag & drop kanban boards
   - Real-time GitHub Projects sync

3. **Powerful Data Views**
   - Sort by any column
   - Quick context actions
   - Fast keyboard navigation

4. **Git Integration**
   - Pull/push from any repo
   - No terminal command needed
   - Visual feedback

### For Development

1. **Proven Code**
   - Already works in tkan & TFE
   - Battle-tested patterns
   - Known performance

2. **Modular Design**
   - Components work independently
   - Easy to add to new tabs
   - Reusable across views

3. **Incremental Integration**
   - One feature at a time
   - Test as you go
   - Low risk

---

## 🚀 Quick Wins First

### Week 1 Priorities (Highest Impact, Lowest Effort)

1. **Table View for Repos** (2 days)
   - Most requested feature
   - Patterns exist in TFE
   - Immediate value

2. **Context Menu System** (1 day)
   - Works for all views
   - Generic component
   - Easy to add actions

3. **Git Operations** (2 days)
   - Pull/push/sync
   - Works with local repos
   - Power user feature

### Demo After Week 1

```bash
cd ~/projects/gh-tui
./gh-tui

# Try it:
1. Press 3 (Repos tab)
2. Press 'v' (Switch to table view)
3. Click any header (Sort)
4. Right-click repo (Context menu)
5. Select "Pull" (Git operation runs)
```

---

## 💾 File Organization

```
gh-tui/
├── main.go
├── types.go
├── model.go
├── update.go
├── view.go
│
├── GitHub API Integration
│   ├── github.go              # Existing
│   ├── github_projects.go     # NEW (from tkan)
│   └── git_operations.go      # NEW (from TFE)
│
├── Reusable Components
│   ├── table_view.go          # NEW
│   ├── context_menu.go        # NEW (from TFE)
│   └── drag_drop.go           # NEW (from tkan)
│
├── View Implementations
│   ├── view_pullrequests.go
│   ├── view_issues.go
│   ├── view_repositories.go   # Enhanced with table
│   ├── view_actions.go
│   ├── view_gists.go
│   └── view_projects.go       # NEW (kanban from tkan)
│
└── View Helpers
    ├── view_projects_kanban.go  # NEW
    ├── view_projects_table.go   # NEW
    └── view_table_*.go          # Table variants per view
```

---

## 🎓 Learning Resources

### GitHub Projects API (for kanban backend)

tkan already has this working in `backend_github.go`:
- Uses GraphQL API
- Handles authentication
- Manages card creation/movement
- Syncs with GitHub Projects

**Note**: The automation features in tkan were just brainstorming and not implemented yet.

### Table Sorting (from TFE)

TFE demonstrates:
- Click detection on headers
- Multi-column sorting
- Visual indicators
- Responsive layout

### Mouse Handling (from tkan)

tkan shows production drag & drop:
- Drag detection with threshold
- Ghost rendering
- Drop target calculation
- Smooth animations

---

## ✅ Next Steps

1. **Review this plan** - Any changes needed?
2. **Pick Phase 1 start** - Projects tab or table view first?
3. **Set up integration branch** - `feature/tkan-integration`?
4. **Start with quick win** - Table view for repos (easiest)?

Would you like me to start implementing any of these? I recommend:

### 🎯 Recommended: Start with Table View (2-3 hours)

**Why:**
- Immediate value for current repos tab
- Foundation for all other tabs
- Easier than kanban (no drag & drop)
- Users will see instant improvement

**Steps:**
1. Create `table_view.go` component
2. Add table mode to repos view
3. Add `v` key toggle
4. Add click-to-sort

Should I start with this? 🚀
