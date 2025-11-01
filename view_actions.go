package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ActionsView displays a list of workflow runs
type ActionsView struct {
	data     []WorkflowRun
	cursor   int
	focused  bool
	err      error
	loading  bool
	width    int
	height   int
}

// NewActionsView creates a new actions view
func NewActionsView() *ActionsView {
	return &ActionsView{
		data:    []WorkflowRun{},
		cursor:  0,
		focused: false,
		loading: true,
	}
}

// Update handles messages for the actions view
func (v *ActionsView) Update(msg tea.Msg) (View, tea.Cmd) {
	switch msg := msg.(type) {
	case workflowsLoadedMsg:
		v.loading = false
		if msg.err != nil {
			v.err = msg.err
		} else {
			v.data = msg.runs
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
			return v, fetchWorkflowRuns("")
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

// View renders the actions view
func (v *ActionsView) View(width, height int) string {
	v.width = width
	v.height = height

	if v.loading {
		return lipgloss.Place(width, height,
			lipgloss.Center, lipgloss.Center,
			infoStyle.Render("Loading workflow runs..."))
	}

	if v.err != nil {
		return lipgloss.Place(width, height,
			lipgloss.Center, lipgloss.Center,
			errorStyle.Render(fmt.Sprintf("Error: %v", v.err)))
	}

	if len(v.data) == 0 {
		return lipgloss.Place(width, height,
			lipgloss.Center, lipgloss.Center,
			dimmedStyle.Render("No workflow runs found"))
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

// renderList renders the workflow runs list
func (v *ActionsView) renderList(width, height int) string {
	var lines []string

	// Header
	title := listTitleStyle.Render(fmt.Sprintf(" Workflow Runs (%d)", len(v.data)))
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

	// Render visible workflow runs
	for i := start; i < end; i++ {
		run := v.data[i]
		cursor := "  "
		style := listItemStyle

		if i == v.cursor {
			cursor = "▶ "
			style = listSelectedStyle
		}

		// Format: "▶ Workflow Name - status • branch"
		status := formatStatus(run.Status, run.Conclusion)
		line := fmt.Sprintf("%s%s %s",
			cursor,
			status,
			truncateString(run.Name, width-35))

		meta := fmt.Sprintf("%s • %s",
			run.HeadBranch,
			formatTimeAgo(run.CreatedAt))

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

// renderDetail renders the detail pane for the selected workflow run
func (v *ActionsView) renderDetail(width, height int) string {
	if v.cursor < 0 || v.cursor >= len(v.data) {
		return ""
	}

	run := v.data[v.cursor]
	var lines []string

	// Title
	title := detailTitleStyle.Render(fmt.Sprintf(" Workflow Run #%d", run.RunNumber))
	lines = append(lines, title)
	lines = append(lines, "")

	// Workflow name (wrapped)
	nameLines := wrapText(run.Name, width-4)
	for _, line := range nameLines {
		lines = append(lines, highlightStyle.Render(line))
	}
	lines = append(lines, "")

	// Metadata
	lines = append(lines, fmt.Sprintf("Status:     %s", formatStatus(run.Status, run.Conclusion)))
	lines = append(lines, fmt.Sprintf("Branch:     %s", run.HeadBranch))
	lines = append(lines, fmt.Sprintf("Commit:     %s", run.HeadSha[:8]))
	lines = append(lines, fmt.Sprintf("Run Number: #%d", run.RunNumber))
	lines = append(lines, fmt.Sprintf("Created:    %s", formatTime(run.CreatedAt)))
	lines = append(lines, fmt.Sprintf("Age:        %s", formatTimeAgo(run.CreatedAt)))

	lines = append(lines, "")
	// Truncate URL to prevent wrapping in narrow detail pane
	displayURL := truncateString(run.URL, width-10)
	clickableURL := makeHyperlink(run.URL, displayURL)
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
func (v *ActionsView) Focus() {
	v.focused = true
}

// Blur sets the view as unfocused
func (v *ActionsView) Blur() {
	v.focused = false
}
