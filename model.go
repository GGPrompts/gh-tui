package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

// model.go - Model Management
// Purpose: Model initialization and layout calculations
// When to extend: Add new initialization logic or layout calculation functions here

// initialModel creates the initial application state
func initialModel(cfg Config) model {
	m := model{
		config:           cfg,
		width:            0,
		height:           0,
		focusedComponent: "main",
		statusMsg:        "gh-tui - GitHub CLI Interactive Interface - Press ? for help",
		activeView:       ViewPullRequests,
		views:            make(map[ViewType]View),
		loading:          false,
		showHelp:         false,
	}

	// Initialize all views
	m.views[ViewPullRequests] = NewPullRequestView()
	m.views[ViewIssues] = NewIssueView()
	m.views[ViewRepositories] = NewRepositoryView()
	m.views[ViewActions] = NewActionsView()
	m.views[ViewGists] = NewGistView()

	// Focus the initial view (Pull Requests)
	if view, ok := m.views[ViewPullRequests]; ok {
		view.Focus()
	}

	return m
}

// setSize updates the model dimensions and recalculates layouts
func (m *model) setSize(width, height int) {
	m.width = width
	m.height = height

	// Recalculate any layout-dependent values here
	// Example:
	// m.viewportHeight = height - 4 // account for title and status bars
	// m.maxVisible = m.viewportHeight - 2
}

// calculateLayout computes layout dimensions based on config
func (m model) calculateLayout() (int, int) {
	contentWidth := m.width
	contentHeight := m.height

	// Adjust for UI elements
	if m.config.UI.ShowTitle {
		contentHeight -= 2 // title bar height
	}
	if m.config.UI.ShowStatus {
		contentHeight -= 1 // status bar height
	}

	return contentWidth, contentHeight
}

// calculateDualPaneLayout computes left and right pane widths
func (m model) calculateDualPaneLayout() (int, int) {
	contentWidth, _ := m.calculateLayout()

	dividerWidth := 0
	if m.config.Layout.ShowDivider {
		dividerWidth = 1
	}

	leftWidth := int(float64(contentWidth-dividerWidth) * m.config.Layout.SplitRatio)
	rightWidth := contentWidth - leftWidth - dividerWidth

	return leftWidth, rightWidth
}

// Helper functions for common operations

// getContentArea returns the available content area dimensions
func (m model) getContentArea() (width, height int) {
	return m.calculateLayout()
}

// isValidSize checks if the terminal size is sufficient
func (m model) isValidSize() bool {
	return m.width >= 40 && m.height >= 10
}

// Init initializes the model and fetches initial data
func (m model) Init() tea.Cmd {
	return tea.Batch(
		fetchPullRequests(""),
		fetchIssues(""),
		fetchRepositories(""),
		fetchWorkflowRuns(""),
		fetchGists(),
	)
}
