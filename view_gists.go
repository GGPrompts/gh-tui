package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// GistView displays a list of gists
type GistView struct {
	data              []Gist
	cursor            int
	focused           bool
	err               error
	loading           bool
	width             int
	height            int
	awaitingGistInput bool   // waiting for description/visibility input for new gist
	newGistFilePath   string // temp file path for new gist being created
}

// NewGistView creates a new gist view
func NewGistView() *GistView {
	return &GistView{
		data:    []Gist{},
		cursor:  0,
		focused: false,
		loading: true,
	}
}

// Update handles messages for the gist view
func (v *GistView) Update(msg tea.Msg) (View, tea.Cmd) {
	switch msg := msg.(type) {
	case gistsLoadedMsg:
		v.loading = false
		if msg.err != nil {
			v.err = msg.err
		} else {
			v.data = msg.gists
			if len(v.data) > 0 && v.cursor >= len(v.data) {
				v.cursor = len(v.data) - 1
			}
		}

	case gistEditorFinishedMsg:
		// Handle editor exit
		if msg.err != nil {
			v.err = fmt.Errorf("editor error: %w", msg.err)
			return v, nil
		}

		if msg.isNewGist {
			// Handle new gist creation
			if msg.wasModified {
				// User wrote content, now create the gist
				// For simplicity, create as private by default
				// TODO: Add dialog to ask for description and public/private
				return v, createGistFromFile(msg.tempFilePath, "", false)
			} else {
				// User didn't write anything, just clean up
				os.Remove(msg.tempFilePath)
			}
		} else if msg.wasModified {
			// Handle gist edit - upload changes
			return v, uploadGistChanges(msg.gistID, msg.tempFilePath)
		} else {
			// No changes made, just clean up temp file
			os.Remove(msg.tempFilePath)
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
			return v, fetchGists()
		case "o":
			// Open/view gist in read-only mode
			if len(v.data) > 0 && v.cursor < len(v.data) {
				gist := v.data[v.cursor]
				if !checkMicroAvailable() {
					v.err = fmt.Errorf("micro editor not found - please install micro")
					return v, nil
				}
				if len(gist.Files) == 0 {
					v.err = fmt.Errorf("gist has no files")
					return v, nil
				}
				// Open first file in read-only mode
				return v, openGistInMicro(gist.ID, gist.Files[0].Filename, true)
			}
		case "e":
			// Edit gist
			if len(v.data) > 0 && v.cursor < len(v.data) {
				gist := v.data[v.cursor]
				if !checkMicroAvailable() {
					v.err = fmt.Errorf("micro editor not found - please install micro")
					return v, nil
				}
				if len(gist.Files) == 0 {
					v.err = fmt.Errorf("gist has no files")
					return v, nil
				}
				// Open first file for editing
				return v, openGistInMicro(gist.ID, gist.Files[0].Filename, false)
			}
		case "n":
			// Create new gist
			if !checkMicroAvailable() {
				v.err = fmt.Errorf("micro editor not found - please install micro")
				return v, nil
			}
			return v, createNewGistInMicro()
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

// View renders the gist view
func (v *GistView) View(width, height int) string {
	v.width = width
	v.height = height

	if v.loading {
		return lipgloss.Place(width, height,
			lipgloss.Center, lipgloss.Center,
			infoStyle.Render("Loading gists..."))
	}

	if v.err != nil {
		return lipgloss.Place(width, height,
			lipgloss.Center, lipgloss.Center,
			errorStyle.Render(fmt.Sprintf("Error: %v", v.err)))
	}

	if len(v.data) == 0 {
		return lipgloss.Place(width, height,
			lipgloss.Center, lipgloss.Center,
			dimmedStyle.Render("No gists found"))
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

// renderList renders the gist list
func (v *GistView) renderList(width, height int) string {
	var lines []string

	// Header
	title := listTitleStyle.Render(fmt.Sprintf(" Gists (%d)", len(v.data)))
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

	// Render visible gists
	for i := start; i < end; i++ {
		gist := v.data[i]
		cursor := "  "
		style := listItemStyle

		if i == v.cursor {
			cursor = "▶ "
			style = listSelectedStyle
		}

		// Get description or first file name
		displayName := gist.Description
		if displayName == "" && len(gist.Files) > 0 {
			displayName = gist.Files[0].Filename
		}
		if displayName == "" {
			displayName = gist.ID
		}

		// Format: "▶ Description - visibility • age"
		visibility := "🔒"
		if gist.Public {
			visibility = "🌐"
		}

		line := fmt.Sprintf("%s%s %s",
			cursor,
			visibility,
			truncateString(displayName, width-25))

		meta := fmt.Sprintf("%d files • %s",
			len(gist.Files),
			formatTimeAgo(gist.UpdatedAt))

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

// renderDetail renders the detail pane for the selected gist
func (v *GistView) renderDetail(width, height int) string {
	if v.cursor < 0 || v.cursor >= len(v.data) {
		return ""
	}

	gist := v.data[v.cursor]
	var lines []string

	// Title
	title := detailTitleStyle.Render(" Gist Details")
	lines = append(lines, title)
	lines = append(lines, "")

	// Description or ID
	if gist.Description != "" {
		descLines := wrapText(gist.Description, width-4)
		for _, line := range descLines {
			lines = append(lines, highlightStyle.Render(line))
		}
	} else {
		lines = append(lines, highlightStyle.Render(gist.ID))
	}
	lines = append(lines, "")

	// Metadata
	visibility := "Private"
	if gist.Public {
		visibility = "Public"
	}
	lines = append(lines, fmt.Sprintf("ID:         %s", gist.ID))
	lines = append(lines, fmt.Sprintf("Visibility: %s", visibility))
	lines = append(lines, fmt.Sprintf("Created:    %s", formatTime(gist.CreatedAt)))
	lines = append(lines, fmt.Sprintf("Updated:    %s", formatTimeAgo(gist.UpdatedAt)))
	lines = append(lines, fmt.Sprintf("Files:      %d", len(gist.Files)))

	// List files
	if len(gist.Files) > 0 {
		lines = append(lines, "")
		lines = append(lines, dimmedStyle.Render("Files:"))
		maxFiles := min(5, len(gist.Files))
		for i := 0; i < maxFiles; i++ {
			lines = append(lines, fmt.Sprintf("  • %s", gist.Files[i].Filename))
		}
		if len(gist.Files) > 5 {
			lines = append(lines, dimmedStyle.Render(fmt.Sprintf("  ... and %d more", len(gist.Files)-5)))
		}
	}

	lines = append(lines, "")
	// Create clickable hyperlink with truncated display text
	displayURL := truncateString(gist.URL, width-10)
	clickableURL := makeHyperlink(gist.URL, displayURL)
	lines = append(lines, dimmedStyle.Render(fmt.Sprintf("URL: %s", clickableURL)))

	// Keyboard hints
	lines = append(lines, "")
	lines = append(lines, helpStyle.Render("↑/↓: Navigate • o: View • e: Edit • n: New"))
	lines = append(lines, helpStyle.Render("r: Refresh • q: Quit"))

	content := strings.Join(lines, "\n")
	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Padding(1, 2).
		Render(content)
}

// Focus sets the view as focused
func (v *GistView) Focus() {
	v.focused = true
}

// Blur sets the view as unfocused
func (v *GistView) Blur() {
	v.focused = false
}
