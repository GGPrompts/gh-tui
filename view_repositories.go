package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// RepositoryView displays a list of repositories
type RepositoryView struct {
	data     []Repository
	cursor   int
	focused  bool
	err      error
	loading  bool
	width    int
	height   int
}

// NewRepositoryView creates a new repository view
func NewRepositoryView() *RepositoryView {
	return &RepositoryView{
		data:    []Repository{},
		cursor:  0,
		focused: false,
		loading: true,
	}
}

// Update handles messages for the repository view
func (v *RepositoryView) Update(msg tea.Msg) (View, tea.Cmd) {
	switch msg := msg.(type) {
	case reposLoadedMsg:
		v.loading = false
		if msg.err != nil {
			v.err = msg.err
		} else {
			v.data = msg.repos
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
			return v, fetchRepositories("")
		}
	}

	return v, nil
}

// View renders the repository view
func (v *RepositoryView) View(width, height int) string {
	v.width = width
	v.height = height

	if v.loading {
		return lipgloss.Place(width, height,
			lipgloss.Center, lipgloss.Center,
			infoStyle.Render("Loading repositories..."))
	}

	if v.err != nil {
		return lipgloss.Place(width, height,
			lipgloss.Center, lipgloss.Center,
			errorStyle.Render(fmt.Sprintf("Error: %v", v.err)))
	}

	if len(v.data) == 0 {
		return lipgloss.Place(width, height,
			lipgloss.Center, lipgloss.Center,
			dimmedStyle.Render("No repositories found"))
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

// renderList renders the repository list
func (v *RepositoryView) renderList(width, height int) string {
	var lines []string

	// Header
	title := listTitleStyle.Render(fmt.Sprintf(" Repositories (%d)", len(v.data)))
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

	// Render visible repos
	for i := start; i < end; i++ {
		repo := v.data[i]
		cursor := "  "
		style := listItemStyle

		if i == v.cursor {
			cursor = "▶ "
			style = listSelectedStyle
		}

		// Format: "▶ name - ⭐123 • language"
		line := fmt.Sprintf("%s%s",
			cursor,
			truncateString(repo.NameWithOwner, width-25))

		meta := fmt.Sprintf("⭐%s • %s",
			formatNumber(repo.StargazerCount),
			formatLanguage(repo.PrimaryLanguage))

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

// renderDetail renders the detail pane for the selected repository
func (v *RepositoryView) renderDetail(width, height int) string {
	if v.cursor < 0 || v.cursor >= len(v.data) {
		return ""
	}

	repo := v.data[v.cursor]
	var lines []string

	// Title
	title := detailTitleStyle.Render(fmt.Sprintf(" Repository Details"))
	lines = append(lines, title)
	lines = append(lines, "")

	// Repo name
	lines = append(lines, highlightStyle.Render(repo.NameWithOwner))
	lines = append(lines, "")

	// Description (wrapped)
	if repo.Description != "" {
		descLines := wrapText(repo.Description, width-4)
		for _, line := range descLines {
			lines = append(lines, line)
		}
		lines = append(lines, "")
	}

	// Stats
	lines = append(lines, fmt.Sprintf("⭐ Stars:      %s", formatNumber(repo.StargazerCount)))
	lines = append(lines, fmt.Sprintf("🍴 Forks:      %s", formatNumber(repo.ForkCount)))
	lines = append(lines, fmt.Sprintf("📝 Issues:     %d", repo.OpenIssuesCount))
	lines = append(lines, fmt.Sprintf("🔤 Language:   %s", formatLanguage(repo.PrimaryLanguage)))
	lines = append(lines, fmt.Sprintf("🔓 Visibility: %s", formatVisibility(repo.Visibility)))

	lines = append(lines, "")
	lines = append(lines, dimmedStyle.Render(fmt.Sprintf("URL: %s", repo.URL)))

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
func (v *RepositoryView) Focus() {
	v.focused = true
}

// Blur sets the view as unfocused
func (v *RepositoryView) Blur() {
	v.focused = false
}
