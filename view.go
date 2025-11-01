package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// view.go - View Rendering
// Purpose: Top-level view rendering and layout
// When to extend: Add new view modes or modify layout logic

// View renders the entire application
func (m model) View() string {
	// Check if terminal size is sufficient
	if !m.isValidSize() {
		return m.renderMinimalView()
	}

	// Handle errors
	if m.err != nil {
		return m.renderErrorView()
	}

	// Render based on layout type
	switch m.config.Layout.Type {
	case "single":
		return m.renderSinglePane()

	case "dual_pane":
		return m.renderDualPane()

	case "multi_panel":
		return m.renderMultiPanel()

	case "tabbed":
		return m.renderTabbed()

	default:
		return m.renderSinglePane()
	}
}

// renderSinglePane renders a single-pane layout
func (m model) renderSinglePane() string {
	var sections []string

	// Title bar
	if m.config.UI.ShowTitle {
		sections = append(sections, m.renderTitleBar())
	}

	// Main content
	sections = append(sections, m.renderMainContent())

	// Status bar
	if m.config.UI.ShowStatus {
		sections = append(sections, m.renderStatusBar())
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// renderDualPane renders a dual-pane layout (side-by-side)
func (m model) renderDualPane() string {
	var sections []string

	// Title bar
	if m.config.UI.ShowTitle {
		sections = append(sections, m.renderTitleBar())
	}

	// Calculate pane dimensions
	leftWidth, rightWidth := m.calculateDualPaneLayout()

	// Left pane
	leftPane := m.renderLeftPane(leftWidth)

	// Divider
	divider := ""
	if m.config.Layout.ShowDivider {
		divider = m.renderDivider()
	}

	// Right pane
	rightPane := m.renderRightPane(rightWidth)

	// Join panes horizontally
	panes := lipgloss.JoinHorizontal(lipgloss.Top, leftPane, divider, rightPane)
	sections = append(sections, panes)

	// Status bar
	if m.config.UI.ShowStatus {
		sections = append(sections, m.renderStatusBar())
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// renderMultiPanel renders a multi-panel layout
func (m model) renderMultiPanel() string {
	// Implement multi-panel layout
	// This is a placeholder - customize based on your needs
	return m.renderSinglePane()
}

// renderTabbed renders a tabbed interface
func (m model) renderTabbed() string {
	var sections []string

	// Title bar
	if m.config.UI.ShowTitle {
		sections = append(sections, m.renderTitleBar())
	}

	// Tabs
	sections = append(sections, m.renderTabs())

	// Active view content
	contentWidth, contentHeight := m.calculateLayout()
	contentHeight -= 3 // Subtract space for tabs

	view := m.views[m.activeView]
	if view != nil {
		content := view.View(contentWidth, contentHeight)
		sections = append(sections, content)
	} else {
		sections = append(sections, "View not implemented")
	}

	// Status bar
	if m.config.UI.ShowStatus {
		sections = append(sections, m.renderStatusBar())
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// renderTabs renders the tab bar
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

// Component rendering functions

// renderTitleBar renders the title bar
func (m model) renderTitleBar() string {
	title := titleStyle.Render("gh-tui - GitHub CLI Interactive Interface")
	padding := m.width - lipgloss.Width(title)
	if padding < 0 {
		padding = 0
	}
	return title + strings.Repeat(" ", padding)
}

// renderStatusBar renders the status bar
func (m model) renderStatusBar() string {
	status := m.statusMsg
	width := m.width - lipgloss.Width(status) - 4

	if width < 0 {
		// Truncate status if too long
		maxLen := m.width - 4
		if maxLen > 0 && len(status) > maxLen {
			status = status[:maxLen-3] + "..."
		}
		width = 0
	}

	return statusStyle.Render(status + strings.Repeat(" ", width))
}

// renderMainContent renders the main content area
func (m model) renderMainContent() string {
	contentWidth, contentHeight := m.calculateLayout()

	// Implement your main content rendering here
	// Example:
	// return m.renderItemList(contentWidth, contentHeight)

	placeholder := "Main content area\n\n"
	placeholder += "Implement your content rendering in renderMainContent()\n\n"
	placeholder += "Press ? for help\n"
	placeholder += "Press q to quit"

	return contentStyle.Width(contentWidth).Height(contentHeight).Render(placeholder)
}

// renderLeftPane renders the left pane in dual-pane mode
func (m model) renderLeftPane(width int) string {
	_, contentHeight := m.calculateLayout()

	// Implement left pane content
	content := "Left Pane\n\n"
	content += "Width: " + string(rune(width))

	return leftPaneStyle.Width(width).Height(contentHeight).Render(content)
}

// renderRightPane renders the right pane in dual-pane mode
func (m model) renderRightPane(width int) string {
	_, contentHeight := m.calculateLayout()

	// Implement right pane content
	content := "Right Pane (Preview)\n\n"
	content += "Width: " + string(rune(width))

	return rightPaneStyle.Width(width).Height(contentHeight).Render(content)
}

// renderDivider renders the vertical divider between panes
func (m model) renderDivider() string {
	_, contentHeight := m.calculateLayout()
	divider := strings.Repeat("│\n", contentHeight)
	return dividerStyle.Render(divider)
}

// Error and minimal views

// renderErrorView renders an error message
func (m model) renderErrorView() string {
	content := "Error: " + m.err.Error() + "\n\n"
	content += "Press q to quit"
	return errorStyle.Render(content)
}

// renderMinimalView renders a minimal view for small terminals
func (m model) renderMinimalView() string {
	content := "Terminal too small\n"
	content += "Minimum: 40x10\n"
	content += "Press q to quit"
	return errorStyle.Render(content)
}

// Helper functions

// centerString centers a string within the given width
func centerString(s string, width int) string {
	strWidth := lipgloss.Width(s)
	if strWidth >= width {
		return s
	}
	leftPad := (width - strWidth) / 2
	rightPad := width - strWidth - leftPad
	return strings.Repeat(" ", leftPad) + s + strings.Repeat(" ", rightPad)
}
