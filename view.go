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

	// Show help screen (takes priority over everything)
	if m.showHelp {
		return m.renderHelpScreen()
	}

	// Show landing page if enabled
	if m.showLandingPage && m.landingPage != nil {
		return m.landingPage.Render()
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

// renderHelpScreen renders the help overlay
func (m model) renderHelpScreen() string {
	// Create help content
	var sections []string

	// Title
	title := titleStyle.Render("gh-tui - Keyboard Shortcuts")
	sections = append(sections, title)
	sections = append(sections, "")

	// Global Keys
	sections = append(sections, helpSectionStyle.Render("Global Keys"))
	sections = append(sections, helpKeyStyle.Render("  ?        ")+"  Toggle this help screen")
	sections = append(sections, helpKeyStyle.Render("  q        ")+"  Quit application")
	sections = append(sections, helpKeyStyle.Render("  r        ")+"  Refresh current view")
	sections = append(sections, helpKeyStyle.Render("  Esc      ")+"  Close help / dialogs")
	sections = append(sections, "")

	// Navigation Keys
	sections = append(sections, helpSectionStyle.Render("Navigation"))
	sections = append(sections, helpKeyStyle.Render("  Tab      ")+"  Next tab")
	sections = append(sections, helpKeyStyle.Render("  Shift+Tab")+"  Previous tab")
	sections = append(sections, helpKeyStyle.Render("  1-5      ")+"  Jump to tab (1=PRs, 2=Issues, 3=Repos, 4=Actions, 5=Gists)")
	sections = append(sections, helpKeyStyle.Render("  ↑/↓, j/k ")+"  Navigate list items")
	sections = append(sections, "")

	// Common Actions
	sections = append(sections, helpSectionStyle.Render("Common Actions (All Tabs)"))
	sections = append(sections, helpKeyStyle.Render("  b        ")+"  Open item in browser")
	sections = append(sections, "")

	// Pull Requests Tab
	sections = append(sections, helpSectionStyle.Render("Pull Requests Tab"))
	sections = append(sections, helpKeyStyle.Render("  b        ")+"  Open PR in browser")
	sections = append(sections, helpKeyStyle.Render("  d        ")+"  View diff in pager")
	sections = append(sections, helpKeyStyle.Render("  m        ")+"  Merge PR (coming soon)")
	sections = append(sections, "")

	// Issues Tab
	sections = append(sections, helpSectionStyle.Render("Issues Tab"))
	sections = append(sections, helpKeyStyle.Render("  b        ")+"  Open issue in browser")
	sections = append(sections, helpKeyStyle.Render("  n        ")+"  Create new issue")
	sections = append(sections, helpKeyStyle.Render("  e        ")+"  Edit issue (coming soon)")
	sections = append(sections, "")

	// Repositories Tab
	sections = append(sections, helpSectionStyle.Render("Repositories Tab"))
	sections = append(sections, helpKeyStyle.Render("  b        ")+"  Open repo in browser")
	sections = append(sections, helpKeyStyle.Render("  v        ")+"  Toggle list/table view")
	sections = append(sections, helpKeyStyle.Render("  s        ")+"  Star/unstar repository")
	sections = append(sections, helpKeyStyle.Render("  c        ")+"  Clone repository")
	sections = append(sections, "")

	// Actions Tab
	sections = append(sections, helpSectionStyle.Render("Actions Tab"))
	sections = append(sections, helpKeyStyle.Render("  b        ")+"  Open workflow run in browser")
	sections = append(sections, helpKeyStyle.Render("  l        ")+"  View logs (coming soon)")
	sections = append(sections, helpKeyStyle.Render("  r        ")+"  Re-run workflow (coming soon)")
	sections = append(sections, "")

	// Gists Tab
	sections = append(sections, helpSectionStyle.Render("Gists Tab"))
	sections = append(sections, helpKeyStyle.Render("  o        ")+"  View gist in micro (read-only)")
	sections = append(sections, helpKeyStyle.Render("  e        ")+"  Edit gist with micro")
	sections = append(sections, helpKeyStyle.Render("  n        ")+"  Create new gist")
	sections = append(sections, helpKeyStyle.Render("  b        ")+"  Open gist in browser")
	sections = append(sections, "")

	// Footer
	sections = append(sections, dimmedStyle.Render("Press ? or Esc to close this help screen"))

	// Join all sections
	content := strings.Join(sections, "\n")

	// Create a box around the content
	helpBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#58a6ff")).
		Padding(1, 2).
		Width(m.width - 4).
		MaxWidth(100).
		Render(content)

	// Center the box on screen
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		helpBox,
	)
}
