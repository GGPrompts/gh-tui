package main

import (
	"github.com/charmbracelet/lipgloss"
)

// styles.go - Visual Styling
// Purpose: All Lipgloss style definitions
// When to extend: Add new styles when introducing new visual components

// Color palette - GitHub-inspired dark theme
var (
	colorPrimary    = lipgloss.Color("#58A6FF") // GitHub blue
	colorSecondary  = lipgloss.Color("#BC8CFF") // Purple
	colorBackground = lipgloss.Color("#0D1117") // GitHub dark
	colorForeground = lipgloss.Color("#C9D1D9") // Light gray
	colorAccent     = lipgloss.Color("#3FB950") // GitHub green
	colorError      = lipgloss.Color("#F85149") // GitHub red
	colorWarning    = lipgloss.Color("#D29922") // GitHub yellow
	colorInfo       = lipgloss.Color("#79C0FF") // Light blue

	// Semantic colors
	colorSelected = lipgloss.Color("#58A6FF")
	colorFocused  = lipgloss.Color("#3FB950")
	colorDimmed   = lipgloss.Color("#8B949E")
	colorBorder   = lipgloss.Color("#30363D")
)

// Base styles

var baseStyle = lipgloss.NewStyle().
	Foreground(colorForeground)

// Layout styles

var titleStyle = lipgloss.NewStyle().
	Foreground(colorPrimary).
	Bold(true).
	Padding(0, 1)

var statusStyle = lipgloss.NewStyle().
	Foreground(colorDimmed).
	Padding(0, 1)

var contentStyle = lipgloss.NewStyle().
	Padding(0, 1)

var leftPaneStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(colorBorder).
	BorderLeft(false).
	BorderTop(false).
	BorderBottom(false).
	Padding(0, 1)

var rightPaneStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(colorBorder).
	BorderRight(false).
	BorderTop(false).
	BorderBottom(false).
	Padding(0, 1)

var dividerStyle = lipgloss.NewStyle().
	Foreground(colorBorder)

// Component styles

var selectedStyle = lipgloss.NewStyle().
	Foreground(colorSelected).
	Bold(true).
	Background(lipgloss.Color("#3E4451"))

var focusedStyle = lipgloss.NewStyle().
	Foreground(colorFocused).
	Bold(true)

var dimmedStyle = lipgloss.NewStyle().
	Foreground(colorDimmed)

var highlightStyle = lipgloss.NewStyle().
	Foreground(colorAccent).
	Bold(true)

// List styles

var listItemStyle = lipgloss.NewStyle().
	Foreground(colorForeground)

var listSelectedStyle = lipgloss.NewStyle().
	Foreground(colorSelected).
	Bold(true).
	Background(lipgloss.Color("#3E4451")).
	Padding(0, 1)

var listCursorStyle = lipgloss.NewStyle().
	Foreground(colorAccent).
	Bold(true)

// Button styles

var buttonStyle = lipgloss.NewStyle().
	Foreground(colorForeground).
	Background(colorBorder).
	Padding(0, 2).
	MarginRight(1)

var buttonActiveStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#282C34")).
	Background(colorPrimary).
	Padding(0, 2).
	MarginRight(1).
	Bold(true)

// Dialog styles

var dialogBoxStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(colorPrimary).
	Padding(1, 2).
	Width(50)

var dialogTitleStyle = lipgloss.NewStyle().
	Foreground(colorPrimary).
	Bold(true).
	Align(lipgloss.Center)

var dialogContentStyle = lipgloss.NewStyle().
	Foreground(colorForeground).
	MarginTop(1).
	MarginBottom(1)

// Input styles

var inputStyle = lipgloss.NewStyle().
	Foreground(colorForeground).
	Background(colorBorder).
	Padding(0, 1)

var inputFocusedStyle = lipgloss.NewStyle().
	Foreground(colorForeground).
	Background(colorBorder).
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(colorPrimary).
	Padding(0, 1)

// Message styles

var errorStyle = lipgloss.NewStyle().
	Foreground(colorError).
	Bold(true).
	Padding(1, 2)

var warningStyle = lipgloss.NewStyle().
	Foreground(colorWarning).
	Bold(true).
	Padding(1, 2)

var infoStyle = lipgloss.NewStyle().
	Foreground(colorInfo).
	Padding(1, 2)

var successStyle = lipgloss.NewStyle().
	Foreground(colorAccent).
	Bold(true).
	Padding(1, 2)

// Table styles

var tableHeaderStyle = lipgloss.NewStyle().
	Foreground(colorPrimary).
	Bold(true).
	BorderStyle(lipgloss.NormalBorder()).
	BorderBottom(true).
	BorderForeground(colorBorder)

var tableCellStyle = lipgloss.NewStyle().
	Foreground(colorForeground).
	Padding(0, 1)

var tableSelectedStyle = lipgloss.NewStyle().
	Foreground(colorSelected).
	Bold(true).
	Background(lipgloss.Color("#3E4451"))

// Menu styles

var menuItemStyle = lipgloss.NewStyle().
	Foreground(colorForeground).
	Padding(0, 2)

var menuSelectedStyle = lipgloss.NewStyle().
	Foreground(colorBackground).
	Background(colorPrimary).
	Bold(true).
	Padding(0, 2)

var menuSeparatorStyle = lipgloss.NewStyle().
	Foreground(colorBorder)

// Helper functions for dynamic styling

// applyTheme applies a theme to all styles
func applyTheme(theme ThemeColors) {
	colorPrimary = lipgloss.Color(theme.Primary)
	colorSecondary = lipgloss.Color(theme.Secondary)
	colorBackground = lipgloss.Color(theme.Background)
	colorForeground = lipgloss.Color(theme.Foreground)
	colorAccent = lipgloss.Color(theme.Accent)
	colorError = lipgloss.Color(theme.Error)

	// Update all styles
	titleStyle = titleStyle.Foreground(colorPrimary)
	statusStyle = statusStyle.Foreground(colorDimmed)
	selectedStyle = selectedStyle.Foreground(colorPrimary)
	// ... update other styles as needed
}

// getTheme returns the current theme colors
func getTheme() ThemeColors {
	return ThemeColors{
		Primary:    string(colorPrimary),
		Secondary:  string(colorSecondary),
		Background: string(colorBackground),
		Foreground: string(colorForeground),
		Accent:     string(colorAccent),
		Error:      string(colorError),
	}
}

// Tab styles for tabbed interface
var activeTabStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#000")).
	Background(colorPrimary).
	Padding(0, 2)

var inactiveTabStyle = lipgloss.NewStyle().
	Foreground(colorDimmed).
	Padding(0, 2)

// List and panel styles
var listPanelStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(colorBorder).
	Padding(1, 2)

var detailPanelStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(colorSecondary).
	Padding(1, 2)

var listTitleStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(colorPrimary)

var detailTitleStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(colorSecondary)

var helpStyle = lipgloss.NewStyle().
	Foreground(colorDimmed)

var actionStyle = lipgloss.NewStyle().
	Foreground(colorAccent)

// Help screen styles
var helpSectionStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#58a6ff")).
	Underline(true)

var helpKeyStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#3fb950"))
