package main

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// update_keyboard.go - Keyboard Event Handling
// Purpose: All keyboard input processing
// When to extend: Add new keyboard shortcuts or key bindings here

// handleKeyPress handles keyboard input
func (m model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Landing page has its own controls
	if m.showLandingPage {
		return m.handleLandingPageKeys(msg)
	}

	// Global keybindings (work in all modes)
	switch {
	case key.Matches(msg, keys.Quit):
		return m, tea.Quit

	case key.Matches(msg, keys.Help):
		return m.toggleHelp()

	case key.Matches(msg, keys.Refresh):
		return m.refresh()
	}

	// Mode-specific keybindings
	switch m.focusedComponent {
	case "main":
		return m.handleMainKeys(msg)

	// Add handlers for other components/modes
	// case "dialog":
	//     return m.handleDialogKeys(msg)
	//
	// case "menu":
	//     return m.handleMenuKeys(msg)
	}

	return m, nil
}

// handleLandingPageKeys handles keyboard input on the landing page
func (m model) handleLandingPageKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Check for Enter key using Type instead of String()
	if msg.Type == tea.KeyEnter || msg.String() == " " {
		// Exit landing page and jump to selected view
		if m.landingPage != nil {
			selectedIdx := m.landingPage.GetSelectedItem()
			m.showLandingPage = false
			m.switchToView(ViewType(selectedIdx))

			// Start loading GitHub data now that we're entering the app
			return m, tea.Batch(
				fetchPullRequests(""),
				fetchIssues(""),
				fetchRepositories(""),
				fetchWorkflowRuns(""),
				fetchGists(),
			)
		}
		return m, nil
	}

	switch msg.String() {
	case "q", "ctrl+c", "esc":
		return m, tea.Quit

	case "up", "k":
		if m.landingPage != nil {
			m.landingPage.SelectPrev()
		}
		return m, nil

	case "down", "j":
		if m.landingPage != nil {
			m.landingPage.SelectNext()
		}
		return m, nil
	}

	return m, nil
}

// handleMainKeys handles keys in main view
func (m model) handleMainKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {

	// Tab switching
	case "tab":
		newView := (m.activeView + 1) % 5
		m.switchToView(newView)
		return m, nil

	case "shift+tab":
		var newView ViewType
		if m.activeView == 0 {
			newView = 4
		} else {
			newView = m.activeView - 1
		}
		m.switchToView(newView)
		return m, nil

	// Direct tab access
	case "1":
		m.switchToView(ViewPullRequests)
		return m, nil
	case "2":
		m.switchToView(ViewIssues)
		return m, nil
	case "3":
		m.switchToView(ViewRepositories)
		return m, nil
	case "4":
		m.switchToView(ViewActions)
		return m, nil
	case "5":
		m.switchToView(ViewGists)
		return m, nil

	// Refresh current view
	case "r":
		return m, m.refreshActiveView()
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

// refreshActiveView refreshes the current view's data
func (m model) refreshActiveView() tea.Cmd {
	switch m.activeView {
	case ViewPullRequests:
		m.loading = true
		return fetchPullRequests("")
	case ViewIssues:
		m.loading = true
		return fetchIssues("")
	case ViewRepositories:
		m.loading = true
		return fetchRepositories("")
	case ViewActions:
		m.loading = true
		return fetchWorkflowRuns("")
	case ViewGists:
		m.loading = true
		return fetchGists()
	}
	return nil
}

// Navigation helper functions

func (m model) moveUp() (tea.Model, tea.Cmd) {
	// Implement up navigation
	// Example: m.cursor = max(0, m.cursor-1)
	return m, nil
}

func (m model) moveDown() (tea.Model, tea.Cmd) {
	// Implement down navigation
	// Example: m.cursor = min(len(m.items)-1, m.cursor+1)
	return m, nil
}

func (m model) moveLeft() (tea.Model, tea.Cmd) {
	// Implement left navigation
	return m, nil
}

func (m model) moveRight() (tea.Model, tea.Cmd) {
	// Implement right navigation
	return m, nil
}

func (m model) pageUp() (tea.Model, tea.Cmd) {
	// Implement page up
	// Example: m.cursor = max(0, m.cursor-m.viewportHeight)
	return m, nil
}

func (m model) pageDown() (tea.Model, tea.Cmd) {
	// Implement page down
	// Example: m.cursor = min(len(m.items)-1, m.cursor+m.viewportHeight)
	return m, nil
}

func (m model) moveToTop() (tea.Model, tea.Cmd) {
	// Example: m.cursor = 0
	return m, nil
}

func (m model) moveToBottom() (tea.Model, tea.Cmd) {
	// Example: m.cursor = len(m.items) - 1
	return m, nil
}

// Action helper functions

func (m model) selectItem() (tea.Model, tea.Cmd) {
	// Implement item selection
	return m, nil
}

func (m model) toggleSelection() (tea.Model, tea.Cmd) {
	// Implement toggle selection
	return m, nil
}

func (m model) switchFocus() (tea.Model, tea.Cmd) {
	// Implement focus switching between components
	return m, nil
}

func (m model) toggleHelp() (tea.Model, tea.Cmd) {
	// Toggle help dialog
	m.statusMsg = "Help: q=quit, ?=help, ↑↓=navigate, enter=select"
	return m, nil
}

func (m model) refresh() (tea.Model, tea.Cmd) {
	// Refresh the current view
	m.statusMsg = "Refreshed"
	return m, nil
}

// switchToView changes the active view and manages focus
func (m *model) switchToView(newView ViewType) {
	// Blur the old view
	if oldView, ok := m.views[m.activeView]; ok {
		oldView.Blur()
	}

	// Update active view
	m.activeView = newView

	// Focus the new view
	if view, ok := m.views[m.activeView]; ok {
		view.Focus()
	}
}

// Key bindings definition
type keyMap struct {
	Quit    key.Binding
	Help    key.Binding
	Refresh key.Binding
	Up      key.Binding
	Down    key.Binding
	Left    key.Binding
	Right   key.Binding
	Select  key.Binding
	Toggle  key.Binding
}

var keys = keyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("ctrl+r"),
		key.WithHelp("ctrl+r", "refresh"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "right"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Toggle: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "toggle"),
	),
}
