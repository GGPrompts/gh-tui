package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// PullRequestView displays a list of pull requests
type PullRequestView struct {
	data     []PullRequest
	cursor   int
	focused  bool
	err      error
	loading  bool
	width    int
	height   int
}

// NewPullRequestView creates a new pull request view
func NewPullRequestView() *PullRequestView {
	return &PullRequestView{
		data:    []PullRequest{},
		cursor:  0,
		focused: true,
		loading: true,
	}
}

// Update handles messages for the pull request view
func (v *PullRequestView) Update(msg tea.Msg) (View, tea.Cmd) {
	switch msg := msg.(type) {
	case prLoadedMsg:
		v.loading = false
		if msg.err != nil {
			v.err = msg.err
		} else {
			v.data = msg.prs
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
			return v, fetchPullRequests("")
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

// View renders the pull request view
func (v *PullRequestView) View(width, height int) string {
	v.width = width
	v.height = height

	if v.loading {
		return lipgloss.Place(width, height,
			lipgloss.Center, lipgloss.Center,
			infoStyle.Render("Loading pull requests..."))
	}

	if v.err != nil {
		return lipgloss.Place(width, height,
			lipgloss.Center, lipgloss.Center,
			errorStyle.Render(fmt.Sprintf("Error: %v", v.err)))
	}

	if len(v.data) == 0 {
		return lipgloss.Place(width, height,
			lipgloss.Center, lipgloss.Center,
			dimmedStyle.Render("No pull requests found"))
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

// renderList renders the PR list
func (v *PullRequestView) renderList(width, height int) string {
	var lines []string

	// Header
	title := listTitleStyle.Render(fmt.Sprintf(" Pull Requests (%d)", len(v.data)))
	lines = append(lines, title)
	lines = append(lines, "")

	// Calculate visible range
	maxVisible := height - 3 // Account for title and padding
	start := v.cursor - maxVisible/2
	if start < 0 {
		start = 0
	}
	end := start + maxVisible
	if end > len(v.data) {
		end = len(v.data)
		start = max(0, end-maxVisible)
	}

	// Render visible PRs
	for i := start; i < end; i++ {
		pr := v.data[i]
		cursor := "  "
		style := listItemStyle

		if i == v.cursor {
			cursor = "▶ "
			style = listSelectedStyle
		}

		// Format: "▶ #123 Title - author • 2h ago"
		line := fmt.Sprintf("%s#%d %s",
			cursor,
			pr.Number,
			truncateString(pr.Title, width-20))

		meta := fmt.Sprintf("%s • %s",
			pr.Author.Login,
			formatTimeAgo(pr.UpdatedAt))

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

// renderDetail renders the detail pane for the selected PR
func (v *PullRequestView) renderDetail(width, height int) string {
	if v.cursor < 0 || v.cursor >= len(v.data) {
		return ""
	}

	pr := v.data[v.cursor]
	var lines []string

	// Title
	title := detailTitleStyle.Render(fmt.Sprintf(" PR #%d Details", pr.Number))
	lines = append(lines, title)
	lines = append(lines, "")

	// PR Title (wrapped)
	titleLines := wrapText(pr.Title, width-4)
	for _, line := range titleLines {
		lines = append(lines, highlightStyle.Render(line))
	}
	lines = append(lines, "")

	// Metadata
	lines = append(lines, fmt.Sprintf("Author:    %s", pr.Author.Login))
	lines = append(lines, fmt.Sprintf("State:     %s", formatPRState(pr.State, pr.IsDraft)))
	lines = append(lines, fmt.Sprintf("Branch:    %s → %s", pr.HeadRefName, pr.BaseRefName))
	lines = append(lines, fmt.Sprintf("Created:   %s", formatTime(pr.CreatedAt)))
	lines = append(lines, fmt.Sprintf("Updated:   %s", formatTimeAgo(pr.UpdatedAt)))

	if pr.ReviewDecision != "" {
		lines = append(lines, fmt.Sprintf("Reviews:   %s", pr.ReviewDecision))
	}

	if pr.Mergeable != "" {
		lines = append(lines, fmt.Sprintf("Mergeable: %s", pr.Mergeable))
	}

	lines = append(lines, "")
	// Truncate URL to prevent wrapping in narrow detail pane
	displayURL := truncateString(pr.URL, width-10)
	clickableURL := makeHyperlink(pr.URL, displayURL)
	lines = append(lines, dimmedStyle.Render(fmt.Sprintf("URL: %s", clickableURL)))

	// Keyboard hints
	lines = append(lines, "")
	lines = append(lines, helpStyle.Render("↑/↓: Navigate • r: Refresh • q: Quit"))

	content := strings.Join(lines, "\n")
	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Padding(1, 2).
		Render(content)
}

// Focus sets the view as focused
func (v *PullRequestView) Focus() {
	v.focused = true
}

// Blur sets the view as unfocused
func (v *PullRequestView) Blur() {
	v.focused = false
}
