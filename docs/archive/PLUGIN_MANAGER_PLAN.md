# Plugin Manager Implementation Plan

## Overview

Add a comprehensive Plugin Management tab to gh-tui that integrates with Claude Code's plugin system and uses GitHub Gists as the distribution mechanism.

**IMPORTANT**: This implementation will use TFE's tree view pattern with expand/collapse functionality for both:
1. **Plugins tab** - Browse plugins with expandable categories and multi-file visualization
2. **Gists tab enhancement** - Show multi-file gists as expandable trees

## Context: Why This Matters

### Current Problem
- All slash commands, skills, and agents load into context window every session
- Massive system prompt waste on unused features
- Can't maintain large library without token bloat
- Manual Obsidian vault syncing is cumbersome

### Plugin Solution Benefits
- ✅ On-demand activation (enable only what you need)
- ✅ Smaller context windows (faster responses, less waste)
- ✅ Massive library storage (100+ skills, use 5 at a time)
- ✅ Project-specific configurations
- ✅ Team consistency via `.claude/settings.json`
- ✅ Gist-based sharing (version control, forkable, discoverable)

## UI Design

### Tab 6: Plugins View (Tree View with TFE-style expand/collapse)

```
┌─ Plugin Sources (Tree) ────┬─ Plugin Details ──────────────┐
│ 🔌 Plugins                 │ Name: TUI Development Suite   │
│ ├─▼ My Gists (23)          │ Type: Skill + Commands        │
│ │  ├─▼ [Plugin] TUI Dev    │ Gist ID: abc123...            │
│ │  │  ├─ skills/           │ Files: 3                      │
│ │  │  │  └─ bubbletea.md   │ Updated: 2 days ago           │
│ │  │  └─ commands/         │                               │
│ │  │     ├─ tui-new.md     │ Description:                  │
│ │  │     └─ tui-comp.md    │ Comprehensive TUI development │
│ │  ├─▶ [Plugin] React      │ skills with Bubbletea...      │
│ │  └─▶ [Plugin] DevOps     │                               │
│ ├─▼ Installed (5)          │ Files (expandable):           │
│ │  ├─▼✓ bubbletea-helpers  │ └─▼ skills/                   │
│ │  │  └─ skills/           │    └─ bubbletea.md            │
│ │  │     └─ bubbletea.md   │ └─▼ commands/                 │
│ │  ├─○ testing-suite       │    ├─ tui-new.md              │
│ │  └─▼✓ git-workflow       │    └─ tui-component.md        │
│ │     ├─ commands/         │                               │
│ │     │  └─ git-sync.md    │ Installation Path:            │
│ │     └─ hooks/            │ ~/.config/claude-code/        │
│ │        └─ pre-commit.sh  │                               │
│ ├─▶ Marketplaces (3)       │ Status: Not installed         │
│ └─▶ Local (2)              │                               │
│                            │ [i] Install  [u] Update       │
│                            │ [e] Enable   [x] Disable      │
│                            │ [p] Publish to Gist           │
└────────────────────────────┴───────────────────────────────┘

Keys: →/←: Expand/Collapse • Enter: Toggle • i=Install • e=Enable
      x=Disable • p=Publish • n=New • r=Refresh • f=Fork

Symbols:
  ▶ = Collapsed   ▼ = Expanded   ✓ = Enabled   ○ = Disabled
  ├─ = Tree branch   └─ = Last item
```

### Enhanced Tab 5: Gists View (Tree View for Multi-file Gists)

```
┌─ Gists (Tree View) ────────┬─ Gist Details ────────────────┐
│ 🌐 Public (15)             │ Opus' Pre-Weekly-Limits Plan  │
│ ├─▶ Quick Scripts         │ Gist ID: dbbd08f2...          │
│ ├─▼ Opus Plan (1 file)    │ Public • Updated 2 days ago   │
│ │  └─ Opus' Pre-Weekly... │                               │
│ └─▼ Config Files (3 files)│ Files:                        │
│    ├─ config.yaml         │ └─▼ Opus' Pre-Weekly-Limits   │
│    ├─ script.sh           │    (14.2 KB)                  │
│    └─ README.md           │                               │
│ 🔒 Private (8)             │ Content Preview:              │
│ ├─▼ [Plugin] TUI Dev      │ Operation Maximum Extraction  │
│ │  ├─ skills/             │ The August 2025 Claude Gold   │
│ │  │  └─ bubbletea.md     │ Rush Strategy...              │
│ │  └─ commands/           │                               │
│ │     ├─ tui-new.md       │ [Click URL to open in browser]│
│ │     └─ tui-comp.md      │ https://gist.github.com/...   │
│ └─▶ Personal Notes        │                               │
│                            │                               │
└────────────────────────────┴───────────────────────────────┘

Keys: →/←: Expand/Collapse • Enter: Toggle expansion
      o=View • e=Edit • n=New • i=Install as Plugin
```

### Status Indicators
- ✓ = Enabled and active
- ○ = Installed but disabled
- ∅ = Not installed
- ⚠ = Update available
- 🔄 = Syncing

## Tree View Architecture (Adapted from TFE)

### Core Tree Structures

```go
// Adapted from TFE's tree view system
type TreeItem struct {
	Type        TreeItemType
	Name        string
	Depth       int
	IsLast      bool
	ParentLasts []bool

	// Type-specific data
	Plugin      *Plugin      // If Type == ItemTypePlugin
	File        *PluginFile  // If Type == ItemTypeFile
	Category    string       // If Type == ItemTypeCategory
	Gist        *Gist        // For Gists tab tree view
}

type TreeItemType int
const (
	ItemTypeCategory TreeItemType = iota  // "My Gists", "Installed", etc.
	ItemTypePlugin                        // Individual plugin
	ItemTypeFile                          // File within plugin
	ItemTypeGist                          // For Gists tab
	ItemTypeGistFile                      // File within gist
)

// Expansion state tracking (like TFE's expandedDirs)
type PluginView struct {
	expandedCategories map[string]bool  // "My Gists" -> true
	expandedPlugins    map[string]bool  // pluginID -> true
	expandedGists      map[string]bool  // gistID -> true (for Gists tab)
	treeItems          []TreeItem       // Cached flattened tree
}
```

### Tree Building (Adapted from TFE's buildTreeItems)

```go
// Build flattened tree from hierarchical data
func (v *PluginView) buildTreeItems(depth int, parentLasts []bool) []TreeItem {
	items := []TreeItem{}

	// Add "My Gists" category
	if v.expandedCategories["My Gists"] {
		items = append(items, TreeItem{
			Type:  ItemTypeCategory,
			Name:  "My Gists",
			Depth: depth,
		})

		// Add gist plugins under category
		for _, plugin := range v.gistPlugins {
			items = append(items, buildPluginTreeItem(plugin, depth+1))

			// If plugin expanded, add files
			if v.expandedPlugins[plugin.ID] {
				items = append(items, buildFileTreeItems(plugin.Files, depth+2)...)
			}
		}
	}

	return items
}
```

## Data Structures

### Plugin Type (types.go)

```go
type Plugin struct {
	ID          string            // Unique identifier (gist ID or local path hash)
	Name        string            // Display name
	Description string            // Plugin description
	Source      PluginSource      // Where it comes from
	GistID      string            // If from gist
	GistURL     string            // GitHub URL
	Files       []PluginFile      // Files in plugin
	Type        PluginType        // What kind of plugin
	Version     string            // Version string
	Author      string            // Creator
	UpdatedAt   time.Time         // Last update
	Installed   bool              // Is it installed?
	Enabled     bool              // Is it active?
	LocalPath   string            // Where it's installed
}

type PluginSource int
const (
	PluginSourceGist PluginSource = iota
	PluginSourceMarketplace
	PluginSourceLocal
)

type PluginType int
const (
	PluginTypeSkill PluginType = iota
	PluginTypeCommand
	PluginTypeAgent
	PluginTypeMCP
	PluginTypeHook
	PluginTypeMixed  // Contains multiple types
)

type PluginFile struct {
	Filename string
	Path     string
	Type     PluginType
	Size     int64
}

type Marketplace struct {
	Name        string
	URL         string
	Description string
	Plugins     []Plugin
}

type pluginsLoadedMsg struct {
	plugins []Plugin
	err     error
}
```

## Implementation Phases

### Phase 0: Tree View Foundation (FIRST!)
**Files to create/adapt from TFE:**
- `tree_view.go` - Generic tree rendering utilities (adapted from TFE)
  - `buildTreeItems()` - Flatten hierarchical data
  - `renderTreeBranches()` - Draw tree lines (├─ └─ │)
  - `getTreeIcon()` - Get expand/collapse icons (▶ ▼)

**Features:**
1. Generic tree item structure
2. Expansion state tracking
3. Tree line rendering
4. Navigate with arrows (→ expand, ← collapse)

### Phase 1: Enhanced Gists Tab with Tree View
**Files to modify:**
- `view_gists.go` - Add tree view mode
- `types.go` - Add `expandedGists map[string]bool`

**Features:**
1. Toggle tree view mode (press 't')
2. Expand gists to show all files
3. Navigate files with arrows
4. Select individual files for view/edit
5. Visual file hierarchy

### Phase 2: Plugin Tab Core (MVP)
**Files to create:**
- `view_plugins.go` - Plugin tree list and detail view
- `plugin_manager.go` - Plugin operations (install, enable, etc.)

**Features:**
1. Display installed plugins from `.claude/` in tree
2. List user's gists with `[Plugin]` prefix
3. Expandable categories (My Gists, Installed, etc.)
4. Expand plugins to show files
5. Show plugin details (files, description, status)

### Phase 3: Installation & Management
**Features:**
1. Install plugin from gist → `.claude/`
2. Enable/disable plugins (update settings.json)
3. Uninstall plugins
4. Update check (compare gist updated time)
5. Auto-update plugins

### Phase 4: Creation & Publishing
**Features:**
1. Create new plugin from local files
2. Publish plugin to gist with `[Plugin]` prefix
3. Edit plugin metadata
4. Fork existing plugins

### Phase 5: Marketplace Integration
**Features:**
1. Fetch plugins from marketplace repos
2. Browse community plugins
3. Search/filter plugins
4. Install from marketplaces

### Phase 6: Advanced Features
**Features:**
1. Plugin dependencies
2. Bulk operations
3. Plugin templates
4. Team workspace sync
5. Plugin health checks

## Technical Implementation Details

### 1. Detecting Plugin Gists

```go
// In github.go
func fetchPluginGists() tea.Cmd {
	return func() tea.Msg {
		// Use gh CLI to get gists
		cmd := exec.Command("gh", "api", "/gists", "--paginate")
		output, err := cmd.Output()

		// Filter for gists with [Plugin] prefix in description
		var plugins []Plugin
		for _, gist := range gists {
			if strings.HasPrefix(gist.Description, "[Plugin]") {
				plugin := parsePluginFromGist(gist)
				plugins = append(plugins, plugin)
			}
		}

		return pluginsLoadedMsg{plugins: plugins}
	}
}
```

### 2. Reading Installed Plugins

```go
// In plugin_manager.go
func getInstalledPlugins() ([]Plugin, error) {
	claudeDir := getClaudeConfigDir() // ~/.config/claude-code/

	var plugins []Plugin

	// Check each subdirectory
	dirs := []string{"skills", "commands", "agents", "hooks", "mcp"}
	for _, dir := range dirs {
		path := filepath.Join(claudeDir, dir)
		files, _ := os.ReadDir(path)

		for _, file := range files {
			plugin := parseLocalPlugin(filepath.Join(path, file.Name()))
			plugins = append(plugins, plugin)
		}
	}

	// Check enabled status from settings.json
	settings := loadClaudeSettings()
	for i := range plugins {
		plugins[i].Enabled = isPluginEnabled(plugins[i], settings)
	}

	return plugins, nil
}
```

### 3. Installing Plugin from Gist

```go
// In plugin_manager.go
func installPluginFromGist(gistID string) tea.Cmd {
	return func() tea.Msg {
		// 1. Download gist files
		tempDir, err := downloadGistToTemp(gistID)
		if err != nil {
			return errMsg{err: err}
		}
		defer os.RemoveAll(tempDir)

		// 2. Parse plugin structure
		plugin, err := parsePluginStructure(tempDir)
		if err != nil {
			return errMsg{err: err}
		}

		// 3. Copy to appropriate .claude/ directories
		claudeDir := getClaudeConfigDir()
		for _, file := range plugin.Files {
			destDir := filepath.Join(claudeDir, getPluginTypeDir(file.Type))
			os.MkdirAll(destDir, 0755)

			src := filepath.Join(tempDir, file.Path)
			dst := filepath.Join(destDir, file.Filename)
			copyFile(src, dst)
		}

		// 4. Update metadata tracking
		savePluginMetadata(plugin)

		return pluginInstalledMsg{plugin: plugin}
	}
}
```

### 4. Enable/Disable Plugin

```go
// In plugin_manager.go
func togglePlugin(pluginID string, enable bool) tea.Cmd {
	return func() tea.Msg {
		// Read .claude/settings.json
		settingsPath := filepath.Join(getClaudeConfigDir(), "settings.json")
		settings := loadSettings(settingsPath)

		// Update enabledPlugins array
		if enable {
			settings.EnabledPlugins = append(settings.EnabledPlugins, pluginID)
		} else {
			settings.EnabledPlugins = removeFromSlice(settings.EnabledPlugins, pluginID)
		}

		// Save back
		saveSettings(settingsPath, settings)

		return pluginToggledMsg{pluginID: pluginID, enabled: enable}
	}
}
```

### 5. Publishing Plugin to Gist

```go
// In plugin_manager.go
func publishPluginToGist(plugin Plugin, description string, public bool) tea.Cmd {
	return func() tea.Msg {
		// 1. Collect all plugin files
		files := collectPluginFiles(plugin)

		// 2. Create temp directory with files
		tempDir, _ := os.MkdirTemp("", "plugin-publish-")
		defer os.RemoveAll(tempDir)

		for _, file := range files {
			dest := filepath.Join(tempDir, file.Filename)
			copyFile(file.Path, dest)
		}

		// 3. Create gist with [Plugin] prefix
		desc := fmt.Sprintf("[Plugin] %s", description)

		args := []string{"gist", "create", "-d", desc}
		if public {
			args = append(args, "-p")
		}

		// Add all files
		for _, file := range files {
			args = append(args, filepath.Join(tempDir, file.Filename))
		}

		cmd := exec.Command("gh", args...)
		output, err := cmd.Output()

		if err != nil {
			return errMsg{err: err}
		}

		gistURL := strings.TrimSpace(string(output))
		return pluginPublishedMsg{url: gistURL}
	}
}
```

### 6. Plugin Metadata Tracking

Store metadata in `~/.config/gh-tui/plugins.json`:

```json
{
  "plugins": [
    {
      "id": "abc123",
      "name": "TUI Development Suite",
      "gist_id": "abc123def456",
      "installed_at": "2025-01-15T10:30:00Z",
      "updated_at": "2025-01-20T14:22:00Z",
      "files": [
        {
          "filename": "bubbletea.md",
          "path": "~/.config/claude-code/skills/bubbletea.md",
          "type": "skill"
        }
      ]
    }
  ]
}
```

## File Structure

```
gh-tui/
├── view_plugins.go           # NEW - Plugin list/detail view
├── plugin_manager.go          # NEW - Plugin operations
├── types.go                   # UPDATE - Add Plugin types
├── github.go                  # UPDATE - Add fetchPluginGists()
├── model.go                   # UPDATE - Add Plugins view
└── config.go                  # NEW - Claude config management
```

## Key Bindings

### Navigation
- **Tab/1-6**: Switch between tabs
- **↑/↓**: Navigate plugin list
- **Enter**: View plugin details
- **Esc**: Back to list

### Actions
- **i**: Install plugin from gist
- **u**: Update plugin (re-download from gist)
- **d**: Delete/uninstall plugin
- **e**: Enable plugin
- **x**: Disable plugin
- **p**: Publish current selection as plugin to gist
- **n**: Create new plugin
- **f**: Fork plugin (copy gist and customize)
- **r**: Refresh plugin list
- **s**: Search/filter plugins

## Workflow Examples

### Example 1: Install Plugin from Gist

```
1. User navigates to Plugins tab
2. Sees "My Gists" section with [Plugin] TUI Dev
3. Presses 'i' to install
4. gh-tui downloads gist files
5. Copies to ~/.config/claude-code/skills/
6. Updates plugins.json metadata
7. Shows success message
8. Plugin appears in "Installed" section
```

### Example 2: Create & Publish Plugin

```
1. User has local skill: ~/.config/claude-code/skills/my-skill.md
2. Navigates to Plugins tab → "Local Plugins"
3. Selects "my-skill.md"
4. Presses 'p' to publish
5. Prompted for description: "My Custom Skill"
6. Prompted for visibility: Public/Private
7. gh-tui creates gist with [Plugin] prefix
8. Returns gist URL
9. Plugin now appears in "My Gists" section
```

### Example 3: Enable/Disable for Context Management

```
# Starting TUI project
1. Navigate to Plugins tab
2. Find "bubbletea-helpers" (currently disabled ○)
3. Press 'e' to enable
4. Status changes to ✓
5. ~/.claude/settings.json updated
6. Next Claude Code session loads bubbletea plugin

# Switching to React project
1. Press 'x' on "bubbletea-helpers" (disable)
2. Press 'e' on "react-patterns" (enable)
3. Context window now optimized for React work
```

## Integration with Claude Code

### Settings.json Format

```json
{
  "enabledPlugins": [
    "bubbletea-helpers",
    "git-workflow"
  ],
  "extraKnownMarketplaces": [
    "danavila/claude-plugins",
    "sethhobson/agents"
  ]
}
```

### Plugin Directory Structure

```
~/.config/claude-code/
├── settings.json              # Plugin enable/disable state
├── skills/
│   ├── bubbletea.md          # From plugin
│   └── react-patterns.md     # From plugin
├── commands/
│   ├── tui-new.md            # From plugin
│   └── react-component.md    # From plugin
├── agents/
│   └── code-reviewer.md      # From plugin
└── hooks/
    └── pre-commit.sh         # From plugin
```

### Metadata Tracking

```
~/.config/gh-tui/
└── plugins.json              # Track installed plugins, sources, versions
```

## Success Metrics

1. **Install plugin in < 5 seconds** from gist
2. **Zero manual file copying** required
3. **One-key enable/disable** for context optimization
4. **Searchable plugin library** (100+ plugins manageable)
5. **Team sync via gist URLs** (share link, instant install)

## Future Enhancements

### Phase 6: Advanced Features
- Plugin dependencies (install A requires B)
- Bulk operations (install all from marketplace)
- Plugin templates (scaffolding)
- Team workspace sync (auto-sync team plugins)
- Plugin health checks (validate structure)
- Version pinning (install specific gist version)
- Conflict resolution (handle duplicate names)
- Plugin metrics (usage tracking)

### Phase 7: Community Features
- Star/bookmark plugins
- Plugin ratings/reviews
- Usage statistics
- Trending plugins
- Plugin categories/tags
- Curated collections

## Related Documentation

- Claude Code Plugins: https://docs.claude.com/en/docs/claude-code/plugins.md
- Plugin Marketplaces: https://docs.claude.com/en/docs/claude-code/plugin-marketplaces.md
- MCP Integration: https://modelcontextprotocol.io/

## Gists Tab Enhancements

### New Features for Gists Tab

**Multi-file Gist Visualization:**
```
Before (current):
┌─ Gists ────────────────────┐
│ ▶ Config Files (3 files)   │  <- Can't see what files
│ ▶ [Plugin] TUI Dev         │
└────────────────────────────┘

After (tree view):
┌─ Gists ────────────────────┐
│ ▼ Config Files (3 files)   │
│   ├─ config.yaml           │  <- Expand to see files!
│   ├─ script.sh             │
│   └─ README.md             │
│ ▶ [Plugin] TUI Dev         │
└────────────────────────────┘
```

**New Key Bindings:**
- **t**: Toggle tree view mode
- **→**: Expand gist to show files
- **←**: Collapse gist
- **Enter**: Toggle expansion
- **o/e**: View/edit now works on specific file (when expanded)
- **i**: Install multi-file gist as plugin (new!)

**Benefits:**
1. See all files in multi-file gists without opening
2. Select specific file to view/edit (not just first file)
3. Visual organization of gist contents
4. Quick "install as plugin" for gists with [Plugin] prefix

## Implementation Timeline

- **Phase 0 (Tree Foundation)**: 1-2 hours - Generic tree utilities from TFE
- **Phase 1 (Gists Enhancement)**: 1-2 hours - Add tree view to Gists tab
- **Phase 2 (Plugin MVP)**: 2-3 hours - Basic plugin tab with tree
- **Phase 3**: 2-3 hours - Install/enable/disable functionality
- **Phase 4**: 1-2 hours - Create and publish features
- **Phase 5**: 2-3 hours - Marketplace integration
- **Phase 6+**: Future iterations based on usage

Total MVP (Phases 0-2): ~4-7 hours for tree foundation + enhanced gists + basic plugins
Full System (Phases 0-5): ~9-13 hours for complete plugin management

## Notes

- All operations use `gh` CLI (already authenticated)
- Gist-based plugins are just convention (description prefix)
- Compatible with official Claude Code plugin system
- Can install both gist plugins AND marketplace plugins
- Metadata tracking allows smart updates and management
