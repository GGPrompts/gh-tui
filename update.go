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

	// GitHub data loaded messages
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
