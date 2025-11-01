package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// update.go - Main Update Dispatcher
// Purpose: Message dispatching and non-input event handling
// When to extend: Add new message types or top-level event handlers here

// Update handles all messages and updates the model
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Window resize
	case tea.WindowSizeMsg:
		m.setSize(msg.Width, msg.Height)
		// Update landing page size if it exists
		if m.landingPage != nil {
			m.landingPage.Resize(msg.Width, msg.Height)
		}
		return m, nil

	// Keyboard input
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	// Mouse input
	case tea.MouseMsg:
		return m.handleMouseEvent(msg)

	// Custom messages
	case errMsg:
		m.err = msg.err
		m.statusMsg = "Error: " + msg.err.Error()
		return m, nil

	case statusMsg:
		m.statusMsg = msg.message
		return m, nil

	// Landing page animation tick
	case landingTickMsg:
		if m.showLandingPage && m.landingPage != nil {
			m.landingPage.Update()
			return m, landingTick()
		}
		return m, nil

	// GitHub data loaded messages - forward to views
	case prLoadedMsg:
		m.loading = false
		m.lastSync = time.Now()
		if msg.err != nil {
			m.err = msg.err
			m.statusMsg = "Error loading PRs: " + msg.err.Error()
		} else {
			m.pullRequests = msg.prs
			m.statusMsg = fmt.Sprintf("Loaded %d pull requests", len(msg.prs))
		}
		// Forward to PR view
		if view, ok := m.views[ViewPullRequests]; ok {
			updatedView, _ := view.Update(msg)
			m.views[ViewPullRequests] = updatedView
		}
		return m, nil

	case issuesLoadedMsg:
		m.loading = false
		m.lastSync = time.Now()
		if msg.err != nil {
			m.err = msg.err
			m.statusMsg = "Error loading issues: " + msg.err.Error()
		} else {
			m.issues = msg.issues
			m.statusMsg = fmt.Sprintf("Loaded %d issues", len(msg.issues))
		}
		// Forward to Issues view
		if view, ok := m.views[ViewIssues]; ok {
			updatedView, _ := view.Update(msg)
			m.views[ViewIssues] = updatedView
		}
		return m, nil

	case reposLoadedMsg:
		m.loading = false
		m.lastSync = time.Now()
		if msg.err != nil {
			m.err = msg.err
			m.statusMsg = "Error loading repositories: " + msg.err.Error()
		} else {
			m.repositories = msg.repos
			m.statusMsg = fmt.Sprintf("Loaded %d repositories", len(msg.repos))
		}
		// Forward to Repositories view
		if view, ok := m.views[ViewRepositories]; ok {
			updatedView, _ := view.Update(msg)
			m.views[ViewRepositories] = updatedView
		}
		return m, nil

	case workflowsLoadedMsg:
		m.loading = false
		m.lastSync = time.Now()
		if msg.err != nil {
			m.err = msg.err
			m.statusMsg = "Error loading workflow runs: " + msg.err.Error()
		} else {
			m.workflowRuns = msg.runs
			m.statusMsg = fmt.Sprintf("Loaded %d workflow runs", len(msg.runs))
		}
		// Forward to Actions view
		if view, ok := m.views[ViewActions]; ok {
			updatedView, _ := view.Update(msg)
			m.views[ViewActions] = updatedView
		}
		return m, nil

	case gistsLoadedMsg:
		m.loading = false
		m.lastSync = time.Now()
		if msg.err != nil {
			m.err = msg.err
			m.statusMsg = "Error loading gists: " + msg.err.Error()
		} else {
			m.gists = msg.gists
			m.statusMsg = fmt.Sprintf("Loaded %d gists", len(msg.gists))
		}
		// Forward to Gists view
		if view, ok := m.views[ViewGists]; ok {
			updatedView, _ := view.Update(msg)
			m.views[ViewGists] = updatedView
		}
		return m, nil
	}

	return m, nil
}

// Helper functions for message handling

// sendStatus creates a status message command
func sendStatus(message string) tea.Cmd {
	return func() tea.Msg {
		return statusMsg{message: message}
	}
}

// sendError creates an error message command
func sendError(err error) tea.Cmd {
	return func() tea.Msg {
		return errMsg{err: err}
	}
}

// isSpecialKey checks if a key is a special key (not printable)
func isSpecialKey(key tea.KeyMsg) bool {
	return key.Type != tea.KeyRunes
}
