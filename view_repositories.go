package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// RepositoryView displays a list of repositories
type RepositoryView struct {
	data       []Repository
	cursor     int
	focused    bool
	err        error
	loading    bool
	width      int
	height     int
	viewMode   ViewMode    // List or Table view
	tableState *TableState // Table state for sorting
}

// NewRepositoryView creates a new repository view
func NewRepositoryView() *RepositoryView {
	// Define table columns
	// Note: Unicode arrows (▲▼) are 3 bytes each in UTF-8, so columns need extra width
	columns := []TableColumn{
		{Header: "Name", Width: 30, Sortable: true, SortKey: "name", Alignment: "left"},
		{Header: "⭐", Width: 10, Sortable: true, SortKey: "stars", Alignment: "right"},
		{Header: "🍴", Width: 10, Sortable: true, SortKey: "forks", Alignment: "right"},
		{Header: "Language", Width: 14, Sortable: true, SortKey: "language", Alignment: "left"},
		{Header: "Issues", Width: 10, Sortable: true, SortKey: "issues", Alignment: "right"},
		{Header: "Visibility", Width: 14, Sortable: true, SortKey: "visibility", Alignment: "left"},
	}

	return &RepositoryView{
		data:       []Repository{},
		cursor:     0,
		focused:    false,
		loading:    true,
		viewMode:   ViewModeList, // Default to list view
		tableState: NewTableState(columns),
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
		case "b":
			// Open repo in browser
			if len(v.data) > 0 && v.cursor < len(v.data) {
				repo := v.data[v.cursor]
				return v, openInBrowser("repo", repo.NameWithOwner, "")
			}
		case "s":
			// Star/unstar repository
			if len(v.data) > 0 && v.cursor < len(v.data) {
				repo := v.data[v.cursor]
				return v, toggleRepoStar(repo.NameWithOwner)
			}
		case "c":
			// Clone repository
			if len(v.data) > 0 && v.cursor < len(v.data) {
				repo := v.data[v.cursor]
				return v, cloneRepository(repo.NameWithOwner)
			}
		case "v":
			// Toggle view mode between list and table
			if v.viewMode == ViewModeList {
				v.viewMode = ViewModeTable
			} else {
				v.viewMode = ViewModeList
			}
		}

	case tea.MouseMsg:
		if !v.focused {
			return v, nil
		}

		switch msg.Type {
		case tea.MouseLeft:
			// Handle mouse clicks on table headers
			if v.viewMode == ViewModeTable {
				// Header is at a specific line within the table view:
				// Line 0: Title ("Repositories (N) - Table View")
				// Line 1: Empty line
				// Line 2: Header row with column names
				// Line 3: Empty line
				// Line 4+: Data rows
				//
				// Account for offset from title bar + tabs (typically 2 lines)
				// So absolute Y position of header is around 4 (2 for title/tabs + 2 for title+empty in view)
				headerY := 4 // Approximate position of header in full screen

				if msg.Y >= headerY-1 && msg.Y <= headerY+1 {
					if v.tableState.HandleHeaderClick(msg.X, msg.Y) {
						// Resort data after column change
						v.tableState.SortRepositories(v.data)
					}
				}
			}

		case tea.MouseWheelUp:
			// Scroll up - move cursor up
			if v.cursor > 0 {
				v.cursor--
			}

		case tea.MouseWheelDown:
			// Scroll down - move cursor down
			if v.cursor < len(v.data)-1 {
				v.cursor++
			}
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

	// Render based on view mode
	if v.viewMode == ViewModeTable {
		// Table view - full width
		return v.renderTable(width, height)
	}

	// List view - split view: list on left, detail on right
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

		meta := fmt.Sprintf("⭐ %s • %s",
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
	lines = append(lines, fmt.Sprintf("🔤 Language:   %s", formatLanguage(repo.PrimaryLanguage)))
	lines = append(lines, fmt.Sprintf("🔓 Visibility: %s", formatVisibility(repo.Visibility)))

	lines = append(lines, "")
	// Truncate URL to prevent wrapping in narrow detail pane
	displayURL := truncateString(repo.URL, width-10)
	clickableURL := makeHyperlink(repo.URL, displayURL)
	lines = append(lines, dimmedStyle.Render(fmt.Sprintf("URL: %s", clickableURL)))

	// Keyboard hints
	lines = append(lines, "")
	lines = append(lines, helpStyle.Render("↑/↓: Navigate • b: Browser • s: Star • c: Clone • v: View • r: Refresh"))

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

// renderTable renders the repository table view
func (v *RepositoryView) renderTable(width, height int) string {
	var lines []string

	// Title with view mode indicator
	title := listTitleStyle.Render(fmt.Sprintf(" Repositories (%d) - Table View", len(v.data)))
	viewToggle := dimmedStyle.Render(" [v] Switch to List")
	titleLine := title + "  " + viewToggle
	lines = append(lines, titleLine)
	lines = append(lines, "")

	// Sort data if needed
	sortedData := make([]Repository, len(v.data))
	copy(sortedData, v.data)
	v.tableState.SortRepositories(sortedData)

	// Render table header
	header := v.tableState.RenderHeader(width)
	lines = append(lines, header)
	lines = append(lines, "")

	// Calculate visible range
	maxVisible := height - 5 // Account for title, header, and padding
	start := v.cursor - maxVisible/2
	if start < 0 {
		start = 0
	}
	end := start + maxVisible
	if end > len(sortedData) {
		end = len(sortedData)
		start = max(0, end-maxVisible)
	}

	// Render visible rows
	for i := start; i < end; i++ {
		repo := sortedData[i]

		// Build row cells
		cells := []string{
			truncateString(repo.NameWithOwner, 30),
			formatNumber(repo.StargazerCount),
			formatNumber(repo.ForkCount),
			formatLanguage(repo.PrimaryLanguage),
			formatNumber(repo.OpenIssuesCount),
			formatVisibility(repo.Visibility),
		}

		selected := i == v.cursor
		row := v.tableState.RenderRow(cells, selected, width)
		lines = append(lines, row)
	}

	// Add keyboard hints
	lines = append(lines, "")
	hints := helpStyle.Render("↑/↓: Navigate • v: Toggle View • Click headers to sort • r: Refresh • q: Quit")
	lines = append(lines, hints)

	content := strings.Join(lines, "\n")
	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Render(content)
}
