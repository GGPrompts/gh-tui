package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// IssueView displays a list of issues
type IssueView struct {
	data     []Issue
	cursor   int
	focused  bool
	err      error
	loading  bool
	width    int
	height   int
}

// NewIssueView creates a new issue view
func NewIssueView() *IssueView {
	return &IssueView{
		data:    []Issue{},
		cursor:  0,
		focused: false,
		loading: true,
	}
}

// Update handles messages for the issue view
func (v *IssueView) Update(msg tea.Msg) (View, tea.Cmd) {
	switch msg := msg.(type) {
	case issuesLoadedMsg:
		v.loading = false
		if msg.err != nil {
			v.err = msg.err
		} else {
			v.data = msg.issues
			if len(v.data) > 0 && v.cursor >= len(v.data) {
				v.cursor = len(v.data) - 1
			}
		}

	case tea.KeyMsg:
		if !v.focused {
			return v, nil
		}

		switch msg.String() {
		case "up", "k":
			if v.cursor > 0 {
				v.cursor--
			}
		case "down", "j":
			if v.cursor < len(v.data)-1 {
				v.cursor++
			}
		case "r":
			v.loading = true
			v.err = nil
			return v, fetchIssues("")
		case "b":
			// Open issue in browser
			if len(v.data) > 0 && v.cursor < len(v.data) {
				issue := v.data[v.cursor]
				return v, openInBrowser("issue", fmt.Sprintf("%d", issue.Number), "")
			}
		case "n":
			// Create new issue
			return v, createNewIssue()
		}

	case tea.MouseMsg:
		if !v.focused {
			return v, nil
		}

		switch msg.Type {
		case tea.MouseWheelUp:
			if v.cursor > 0 {
				v.cursor--
			}
		case tea.MouseWheelDown:
			if v.cursor < len(v.data)-1 {
				v.cursor++
			}
		}
	}

	return v, nil
}

// View renders the issue view
func (v *IssueView) View(width, height int) string {
	v.width = width
	v.height = height

	if v.loading {
		return lipgloss.Place(width, height,
			lipgloss.Center, lipgloss.Center,
			infoStyle.Render("Loading issues..."))
	}

	if v.err != nil {
		return lipgloss.Place(width, height,
			lipgloss.Center, lipgloss.Center,
			errorStyle.Render(fmt.Sprintf("Error: %v", v.err)))
	}

	if len(v.data) == 0 {
		return lipgloss.Place(width, height,
			lipgloss.Center, lipgloss.Center,
			dimmedStyle.Render("No issues found"))
	}

	// Split view: list on left, detail on right
	listWidth := width / 2
	detailWidth := width - listWidth - 1

	listContent := v.renderList(listWidth, height)
	detailContent := v.renderDetail(detailWidth, height)

	// Combine list and detail with a divider
	return lipgloss.JoinHorizontal(lipgloss.Top,
		listContent,
		dividerStyle.Render("│"),
		detailContent,
	)
}

// renderList renders the issue list
func (v *IssueView) renderList(width, height int) string {
	var lines []string

	// Header
	title := listTitleStyle.Render(fmt.Sprintf(" Issues (%d)", len(v.data)))
	lines = append(lines, title)
	lines = append(lines, "")

	// Calculate visible range
	maxVisible := height - 3
	start := v.cursor - maxVisible/2
	if start < 0 {
		start = 0
	}
	end := start + maxVisible
	if end > len(v.data) {
		end = len(v.data)
		start = max(0, end-maxVisible)
	}

	// Render visible issues
	for i := start; i < end; i++ {
		issue := v.data[i]
		cursor := "  "
		style := listItemStyle

		if i == v.cursor {
			cursor = "▶ "
			style = listSelectedStyle
		}

		// Format: "▶ #123 Title - author • 2h ago"
		line := fmt.Sprintf("%s#%d %s",
			cursor,
			issue.Number,
			truncateString(issue.Title, width-20))

		meta := fmt.Sprintf("%s • %s",
			issue.Author.Login,
			formatTimeAgo(issue.UpdatedAt))

		if len(line)+len(meta)+3 < width {
			line = padRight(line, width-len(meta)-3) + dimmedStyle.Render(meta)
		}

		lines = append(lines, style.Render(line))
	}

	content := strings.Join(lines, "\n")
	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Render(content)
}

// renderDetail renders the detail pane for the selected issue
func (v *IssueView) renderDetail(width, height int) string {
	if v.cursor < 0 || v.cursor >= len(v.data) {
		return ""
	}

	issue := v.data[v.cursor]
	var lines []string

	// Title
	title := detailTitleStyle.Render(fmt.Sprintf(" Issue #%d Details", issue.Number))
	lines = append(lines, title)
	lines = append(lines, "")

	// Issue Title (wrapped)
	titleLines := wrapText(issue.Title, width-4)
	for _, line := range titleLines {
		lines = append(lines, highlightStyle.Render(line))
	}
	lines = append(lines, "")

	// Metadata
	lines = append(lines, fmt.Sprintf("Author:    %s", issue.Author.Login))
	lines = append(lines, fmt.Sprintf("State:     %s", formatIssueState(issue.State)))
	lines = append(lines, fmt.Sprintf("Created:   %s", formatTime(issue.CreatedAt)))
	lines = append(lines, fmt.Sprintf("Updated:   %s", formatTimeAgo(issue.UpdatedAt)))

	if len(issue.Labels) > 0 {
		labels := formatLabels(issue.Labels)
		lines = append(lines, fmt.Sprintf("Labels:    %s", truncateString(labels, width-15)))
	}

	if len(issue.Assignees) > 0 {
		assignees := make([]string, len(issue.Assignees))
		for i, a := range issue.Assignees {
			assignees[i] = a.Login
		}
		lines = append(lines, fmt.Sprintf("Assignees: %s", truncateString(strings.Join(assignees, ", "), width-15)))
	}

	if issue.Milestone != nil {
		lines = append(lines, fmt.Sprintf("Milestone: %s", issue.Milestone.Title))
	}

	lines = append(lines, "")
	// Truncate URL to prevent wrapping in narrow detail pane
	displayURL := truncateString(issue.URL, width-10)
	clickableURL := makeHyperlink(issue.URL, displayURL)
	lines = append(lines, dimmedStyle.Render(fmt.Sprintf("URL: %s", clickableURL)))

	// Keyboard hints
	lines = append(lines, "")
	lines = append(lines, helpStyle.Render("↑/↓: Navigate • b: Browser • n: New Issue • r: Refresh • q: Quit"))

	content := strings.Join(lines, "\n")
	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Padding(1, 2).
		Render(content)
}

// Focus sets the view as focused
func (v *IssueView) Focus() {
	v.focused = true
}

// Blur sets the view as unfocused
func (v *IssueView) Blur() {
	v.focused = false
}
