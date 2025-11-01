package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// table_view.go - Reusable Table View Component
// Purpose: Provides sortable table headers and table rendering for any data
// Inspired by TFE's detail view implementation

// ViewMode represents different display modes for data
type ViewMode int

const (
	ViewModeList ViewMode = iota
	ViewModeTable
)

func (v ViewMode) String() string {
	switch v {
	case ViewModeList:
		return "List"
	case ViewModeTable:
		return "Table"
	default:
		return "Unknown"
	}
}

// TableColumn defines a column in the table
type TableColumn struct {
	Header    string // Column header text
	Width     int    // Column width in characters
	Sortable  bool   // Whether this column can be sorted
	SortKey   string // Key used for sorting (e.g., "name", "stars")
	Alignment string // "left", "right", "center"
}

// TableState holds the state of a table view
type TableState struct {
	Columns       []TableColumn
	SortColumn    int  // Index of currently sorted column
	SortAscending bool // Sort direction
}

// NewTableState creates a new table state with given columns
func NewTableState(columns []TableColumn) *TableState {
	return &TableState{
		Columns:       columns,
		SortColumn:    0,
		SortAscending: true,
	}
}

// RenderHeader renders the table header with sort indicators
func (t *TableState) RenderHeader(width int) string {
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("87")) // Bright blue

	var headerParts []string
	for i, col := range t.Columns {
		headerText := col.Header

		// Reserve space for sort indicator on sortable columns
		// This ensures consistent width whether sorted or not
		sortIndicator := ""
		if col.Sortable {
			if i == t.SortColumn {
				if t.SortAscending {
					sortIndicator = " ▲"
				} else {
					sortIndicator = " ▼"
				}
			} else {
				// Reserve space even when not sorted to maintain alignment
				sortIndicator = "  " // Two spaces same width as " ▲"
			}
		}

		headerTextWithIndicator := headerText + sortIndicator

		// Pad or truncate to column width
		if len(headerTextWithIndicator) > col.Width {
			headerTextWithIndicator = headerTextWithIndicator[:col.Width-2] + ".."
		} else {
			switch col.Alignment {
			case "right":
				headerTextWithIndicator = fmt.Sprintf("%*s", col.Width, headerTextWithIndicator)
			case "center":
				padding := col.Width - len(headerTextWithIndicator)
				leftPad := padding / 2
				rightPad := padding - leftPad
				headerTextWithIndicator = strings.Repeat(" ", leftPad) + headerTextWithIndicator + strings.Repeat(" ", rightPad)
			default: // left
				headerTextWithIndicator = padRight(headerTextWithIndicator, col.Width)
			}
		}

		headerParts = append(headerParts, headerTextWithIndicator)
	}

	// Add consistent left padding (2 spaces) to match row padding
	header := "  " + strings.Join(headerParts, " │ ")
	return headerStyle.Render(header)
}

// RenderRow renders a single table row
func (t *TableState) RenderRow(cells []string, selected bool, width int) string {
	var rowParts []string
	for i, cell := range cells {
		if i >= len(t.Columns) {
			break
		}
		col := t.Columns[i]

		// Truncate or pad cell to column width
		if len(cell) > col.Width {
			cell = cell[:col.Width-2] + ".."
		} else {
			switch col.Alignment {
			case "right":
				cell = fmt.Sprintf("%*s", col.Width, cell)
			case "center":
				padding := col.Width - len(cell)
				leftPad := padding / 2
				rightPad := padding - leftPad
				cell = strings.Repeat(" ", leftPad) + cell + strings.Repeat(" ", rightPad)
			default: // left
				cell = padRight(cell, col.Width)
			}
		}

		rowParts = append(rowParts, cell)
	}

	row := "  " + strings.Join(rowParts, " │ ")

	// Apply selection style if selected
	if selected {
		return listSelectedStyle.Render(row)
	}
	return listItemStyle.Render(row)
}

// HandleHeaderClick handles a click on a table header
func (t *TableState) HandleHeaderClick(x, y int) bool {
	// Calculate which column was clicked based on x position
	// This is a simplified implementation - assumes fixed column positions
	col := t.GetColumnAtX(x)
	if col >= 0 && col < len(t.Columns) && t.Columns[col].Sortable {
		if col == t.SortColumn {
			// Toggle sort direction if clicking the same column
			t.SortAscending = !t.SortAscending
		} else {
			// Switch to new column, default to ascending
			t.SortColumn = col
			t.SortAscending = true
		}
		return true
	}
	return false
}

// GetColumnAtX returns the column index for a given x coordinate
func (t *TableState) GetColumnAtX(x int) int {
	// Account for left padding (2 spaces)
	x -= 2
	if x < 0 {
		return -1
	}

	// Calculate column based on accumulated widths
	pos := 0
	for i, col := range t.Columns {
		// Add column width + separator (3 chars: " │ ")
		nextPos := pos + col.Width
		if x >= pos && x < nextPos {
			return i
		}
		pos = nextPos + 3 // Add separator width
	}
	return -1
}

// SortRepositories sorts repositories based on current sort state
func (t *TableState) SortRepositories(repos []Repository) {
	if t.SortColumn < 0 || t.SortColumn >= len(t.Columns) {
		return
	}

	sortKey := t.Columns[t.SortColumn].SortKey

	sort.Slice(repos, func(i, j int) bool {
		var less bool
		switch sortKey {
		case "name":
			less = repos[i].NameWithOwner < repos[j].NameWithOwner
		case "stars":
			less = repos[i].StargazerCount < repos[j].StargazerCount
		case "forks":
			less = repos[i].ForkCount < repos[j].ForkCount
		case "language":
			langI := formatLanguage(repos[i].PrimaryLanguage)
			langJ := formatLanguage(repos[j].PrimaryLanguage)
			less = langI < langJ
		case "issues":
			less = repos[i].OpenIssuesCount < repos[j].OpenIssuesCount
		case "visibility":
			// Sort by visibility: private < public (alphabetically)
			less = repos[i].Visibility < repos[j].Visibility
		default:
			less = repos[i].NameWithOwner < repos[j].NameWithOwner
		}

		if !t.SortAscending {
			less = !less
		}
		return less
	})
}

// GetTotalWidth returns the total width needed for all columns
func (t *TableState) GetTotalWidth() int {
	total := 2 // Left padding
	for i, col := range t.Columns {
		total += col.Width
		if i < len(t.Columns)-1 {
			total += 3 // Separator " │ "
		}
	}
	return total
}
