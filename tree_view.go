package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// tree_view.go - Generic Tree View Utilities
// Adapted from TFE's tree rendering system
// Purpose: Reusable tree view for Gists and Plugins tabs

// TreeItem represents an item in a tree structure
type TreeItem struct {
	Type        TreeItemType
	Name        string
	Depth       int
	IsLast      bool
	ParentLasts []bool
	Expanded    bool

	// Type-specific data (use type assertion when needed)
	Data interface{} // Can be *Gist, *Plugin, *GistFile, etc.
}

// TreeItemType identifies what kind of tree item this is
type TreeItemType int

const (
	TreeItemCategory TreeItemType = iota // Top-level category
	TreeItemGist                         // A gist
	TreeItemGistFile                     // File within a gist
	TreeItemPlugin                       // A plugin
	TreeItemPluginFile                   // File within a plugin
	TreeItemDirectory                    // Directory/folder
)

// TreeConfig holds configuration for tree rendering
type TreeConfig struct {
	ShowIcons     bool
	IndentSize    int
	ExpandIcon    string
	CollapseIcon  string
	BranchIcon    string
	LastBranchIcon string
	VerticalIcon  string
}

// DefaultTreeConfig returns the default tree configuration
func DefaultTreeConfig() TreeConfig {
	return TreeConfig{
		ShowIcons:      true,
		IndentSize:     2,
		ExpandIcon:     "▶",
		CollapseIcon:   "▼",
		BranchIcon:     "├─",
		LastBranchIcon: "└─",
		VerticalIcon:   "│ ",
	}
}

// renderTreeBranches generates the tree branch lines for an item
// Adapted from TFE's tree rendering logic
func renderTreeBranches(item TreeItem, config TreeConfig) string {
	if item.Depth == 0 {
		return ""
	}

	var branches strings.Builder

	// Draw vertical lines for parent levels
	for i := 0; i < len(item.ParentLasts); i++ {
		if item.ParentLasts[i] {
			// Parent was last item, don't draw vertical line
			branches.WriteString(strings.Repeat(" ", config.IndentSize))
		} else {
			// Parent has more siblings, draw vertical line
			branches.WriteString(config.VerticalIcon)
		}
	}

	// Draw branch for current item
	if item.IsLast {
		branches.WriteString(config.LastBranchIcon)
	} else {
		branches.WriteString(config.BranchIcon)
	}

	return branches.String()
}

// getExpandCollapseIcon returns the appropriate icon for expandable items
func getExpandCollapseIcon(item TreeItem, config TreeConfig) string {
	// Only show expand/collapse for items that can have children
	switch item.Type {
	case TreeItemCategory, TreeItemGist, TreeItemPlugin, TreeItemDirectory:
		if item.Expanded {
			return config.CollapseIcon
		}
		return config.ExpandIcon
	default:
		// Files don't have expand/collapse icons
		return " "
	}
}

// renderTreeItem renders a single tree item with proper formatting
func renderTreeItem(item TreeItem, config TreeConfig, style lipgloss.Style, selected bool) string {
	branches := renderTreeBranches(item, config)
	icon := getExpandCollapseIcon(item, config)

	// Add spacing between icon and text
	iconPart := ""
	if config.ShowIcons {
		iconPart = icon + " "
	}

	// Build the line
	line := branches + iconPart + item.Name

	// Apply style
	if selected {
		return selectedTreeStyle.Render(line)
	}
	return style.Render(line)
}

// TreeViewState tracks the state of a tree view
type TreeViewState struct {
	ExpandedItems map[string]bool // Key is item identifier
	SelectedIndex int
	TreeItems     []TreeItem
	Config        TreeConfig
}

// NewTreeViewState creates a new tree view state
func NewTreeViewState() *TreeViewState {
	return &TreeViewState{
		ExpandedItems: make(map[string]bool),
		SelectedIndex: 0,
		TreeItems:     []TreeItem{},
		Config:        DefaultTreeConfig(),
	}
}

// IsExpanded checks if an item is expanded
func (s *TreeViewState) IsExpanded(itemKey string) bool {
	return s.ExpandedItems[itemKey]
}

// ToggleExpanded toggles the expansion state of an item
func (s *TreeViewState) ToggleExpanded(itemKey string) {
	s.ExpandedItems[itemKey] = !s.ExpandedItems[itemKey]
}

// SetExpanded sets the expansion state of an item
func (s *TreeViewState) SetExpanded(itemKey string, expanded bool) {
	s.ExpandedItems[itemKey] = expanded
}

// CollapseAll collapses all expanded items
func (s *TreeViewState) CollapseAll() {
	s.ExpandedItems = make(map[string]bool)
}

// GetSelectedItem returns the currently selected tree item
func (s *TreeViewState) GetSelectedItem() *TreeItem {
	if s.SelectedIndex >= 0 && s.SelectedIndex < len(s.TreeItems) {
		return &s.TreeItems[s.SelectedIndex]
	}
	return nil
}

// MoveUp moves the selection up
func (s *TreeViewState) MoveUp() {
	if s.SelectedIndex > 0 {
		s.SelectedIndex--
	}
}

// MoveDown moves the selection down
func (s *TreeViewState) MoveDown() {
	if s.SelectedIndex < len(s.TreeItems)-1 {
		s.SelectedIndex++
	}
}

// renderTreeView renders the entire tree view
func renderTreeView(state *TreeViewState, width, height int) string {
	if len(state.TreeItems) == 0 {
		return lipgloss.Place(width, height,
			lipgloss.Center, lipgloss.Center,
			dimmedStyle.Render("No items"))
	}

	var lines []string

	// Calculate visible range (scrolling support)
	maxVisible := height - 2 // Reserve space for padding
	start := 0
	end := len(state.TreeItems)

	if len(state.TreeItems) > maxVisible {
		start = state.SelectedIndex - maxVisible/2
		if start < 0 {
			start = 0
		}
		end = start + maxVisible
		if end > len(state.TreeItems) {
			end = len(state.TreeItems)
			start = max(0, end-maxVisible)
		}
	}

	// Render visible items
	for i := start; i < end; i++ {
		item := state.TreeItems[i]
		selected := (i == state.SelectedIndex)

		// Choose style based on item type
		style := listItemStyle
		switch item.Type {
		case TreeItemCategory:
			style = listTitleStyle
		case TreeItemGist, TreeItemPlugin:
			style = highlightStyle
		case TreeItemGistFile, TreeItemPluginFile:
			style = dimmedStyle
		}

		line := renderTreeItem(item, state.Config, style, selected)
		lines = append(lines, line)
	}

	content := strings.Join(lines, "\n")
	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Render(content)
}

// Styles for tree view (reuse existing styles from gh-tui)
var selectedTreeStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("62")).
	Foreground(lipgloss.Color("230")).
	Bold(true)

// buildTreeItemsRecursive is a helper for building tree structures recursively
// This is a generic version - specific implementations will use this pattern
func buildTreeItemsRecursive(
	items []TreeItem,
	depth int,
	parentLasts []bool,
	expandedItems map[string]bool,
	getChildren func(TreeItem) []TreeItem,
	getItemKey func(TreeItem) string,
) []TreeItem {
	result := []TreeItem{}

	for i, item := range items {
		isLast := i == len(items)-1

		// Add current item
		item.Depth = depth
		item.IsLast = isLast
		item.ParentLasts = append([]bool{}, parentLasts...)
		item.Expanded = expandedItems[getItemKey(item)]

		result = append(result, item)

		// If expanded, recursively add children
		if item.Expanded {
			children := getChildren(item)
			if len(children) > 0 {
				newParentLasts := append(parentLasts, isLast)
				childItems := buildTreeItemsRecursive(
					children,
					depth+1,
					newParentLasts,
					expandedItems,
					getChildren,
					getItemKey,
				)
				result = append(result, childItems...)
			}
		}
	}

	return result
}

// Example usage functions (to be used by views):

// BuildGistTree builds a tree of gists with files
func BuildGistTree(gists []Gist, expandedGists map[string]bool) []TreeItem {
	rootItems := make([]TreeItem, len(gists))

	for i, gist := range gists {
		rootItems[i] = TreeItem{
			Type: TreeItemGist,
			Name: formatGistName(gist),
			Data: &gist,
		}
	}

	getChildren := func(item TreeItem) []TreeItem {
		if item.Type == TreeItemGist {
			gist := item.Data.(*Gist)
			children := make([]TreeItem, len(gist.Files))
			for i, file := range gist.Files {
				children[i] = TreeItem{
					Type: TreeItemGistFile,
					Name: file.Filename,
					Data: &file,
				}
			}
			return children
		}
		return nil
	}

	getItemKey := func(item TreeItem) string {
		if item.Type == TreeItemGist {
			return item.Data.(*Gist).ID
		}
		return ""
	}

	return buildTreeItemsRecursive(rootItems, 0, []bool{}, expandedGists, getChildren, getItemKey)
}

// formatGistName formats a gist for tree display
func formatGistName(gist Gist) string {
	name := gist.Description
	if name == "" && len(gist.Files) > 0 {
		name = gist.Files[0].Filename
	}
	if name == "" {
		name = gist.ID
	}

	// Add file count
	fileCount := fmt.Sprintf("(%d file", len(gist.Files))
	if len(gist.Files) != 1 {
		fileCount += "s"
	}
	fileCount += ")"

	// Add visibility indicator
	visibility := "🔒"
	if gist.Public {
		visibility = "🌐"
	}

	return fmt.Sprintf("%s %s %s", visibility, name, fileCount)
}
