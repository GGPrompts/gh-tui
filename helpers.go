package main

import (
	"fmt"
	"strings"
	"time"
)

// helpers.go - Utility Functions
// Purpose: Common helper functions used across views

// truncateString truncates text to fit within maxWidth, adding "..." if needed
func truncateString(s string, maxWidth int) string {
	if len(s) <= maxWidth {
		return s
	}
	if maxWidth <= 3 {
		return s[:maxWidth]
	}
	return s[:maxWidth-3] + "..."
}

// padRight pads a string with spaces to reach the target width
func padRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}

// formatTime formats timestamps in a user-friendly way
func formatTime(t time.Time) string {
	return t.Format("Jan 2, 15:04")
}

// formatTimeAgo shows relative time (e.g., "2h ago", "3d ago")
func formatTimeAgo(t time.Time) string {
	duration := time.Since(t)

	if duration.Seconds() < 60 {
		return "just now"
	}

	if duration.Minutes() < 60 {
		mins := int(duration.Minutes())
		return fmt.Sprintf("%dm ago", mins)
	}

	if duration.Hours() < 24 {
		hours := int(duration.Hours())
		return fmt.Sprintf("%dh ago", hours)
	}

	days := int(duration.Hours() / 24)
	if days < 30 {
		return fmt.Sprintf("%dd ago", days)
	}

	months := days / 30
	if months < 12 {
		return fmt.Sprintf("%dmo ago", months)
	}

	years := months / 12
	return fmt.Sprintf("%dy ago", years)
}

// formatLabels formats a list of labels as a comma-separated string
func formatLabels(labels []Label) string {
	if len(labels) == 0 {
		return "none"
	}

	names := make([]string, len(labels))
	for i, label := range labels {
		names[i] = label.Name
	}
	return strings.Join(names, ", ")
}

// formatLanguage formats a language name, handling nil cases
func formatLanguage(lang *Language) string {
	if lang == nil {
		return "none"
	}
	return lang.Name
}

// formatVisibility formats visibility string with an icon
func formatVisibility(vis string) string {
	switch vis {
	case "PUBLIC":
		return "🌐 Public"
	case "PRIVATE":
		return "🔒 Private"
	default:
		return vis
	}
}

// formatStatus formats workflow run status with color/icon
func formatStatus(status, conclusion string) string {
	if status == "completed" {
		switch conclusion {
		case "success":
			return "✓ Success"
		case "failure":
			return "✗ Failure"
		case "cancelled":
			return "⊘ Cancelled"
		case "skipped":
			return "⊙ Skipped"
		default:
			return conclusion
		}
	}
	return status
}

// formatPRState formats PR state with an icon
func formatPRState(state string, isDraft bool) string {
	if isDraft {
		return "📝 Draft"
	}
	switch state {
	case "OPEN":
		return "🟢 Open"
	case "CLOSED":
		return "🔴 Closed"
	case "MERGED":
		return "🟣 Merged"
	default:
		return state
	}
}

// formatIssueState formats issue state with an icon
func formatIssueState(state string) string {
	switch state {
	case "OPEN":
		return "🟢 Open"
	case "CLOSED":
		return "🟣 Closed"
	default:
		return state
	}
}

// formatNumber formats numbers with k/M suffixes for large values
func formatNumber(n int) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}
	if n < 1000000 {
		return fmt.Sprintf("%.1fk", float64(n)/1000)
	}
	return fmt.Sprintf("%.1fM", float64(n)/1000000)
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// makeHyperlink creates a clickable terminal hyperlink using OSC 8
// The displayed text can be truncated while the actual URL remains full
func makeHyperlink(url string, displayText string) string {
	// OSC 8 format: \x1b]8;;URL\x1b\\DISPLAY_TEXT\x1b]8;;\x1b\\
	return fmt.Sprintf("\x1b]8;;%s\x1b\\%s\x1b]8;;\x1b\\", url, displayText)
}

// wrapText wraps text to fit within a given width
func wrapText(text string, width int) []string {
	if width <= 0 {
		return []string{}
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}

	var lines []string
	var currentLine strings.Builder

	for _, word := range words {
		if currentLine.Len() == 0 {
			currentLine.WriteString(word)
		} else if currentLine.Len()+1+len(word) <= width {
			currentLine.WriteString(" " + word)
		} else {
			lines = append(lines, currentLine.String())
			currentLine.Reset()
			currentLine.WriteString(word)
		}
	}

	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}

	return lines
}
