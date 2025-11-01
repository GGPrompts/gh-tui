package rainbow

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Cycler handles rainbow color cycling for text
type Cycler struct {
	Frame  int
	Colors []lipgloss.Color
	Speed  int // How many frames before color shifts (default: 5)
}

// NewCycler creates a new rainbow cycler with default colors
func NewCycler() *Cycler {
	return &Cycler{
		Frame: 0,
		Colors: []lipgloss.Color{
			lipgloss.Color("196"), // Red
			lipgloss.Color("214"), // Orange
			lipgloss.Color("226"), // Yellow
			lipgloss.Color("46"),  // Green
			lipgloss.Color("51"),  // Cyan
			lipgloss.Color("39"),  // Blue
			lipgloss.Color("201"), // Magenta
		},
		Speed: 5,
	}
}

// Update advances the animation frame
func (c *Cycler) Update() {
	c.Frame++
}

// Render applies rainbow colors to text, with each character getting a different color
// The colors cycle through the rainbow and shift with the animation frame
func (c *Cycler) Render(text string) string {
	var b strings.Builder

	for charIdx, char := range text {
		// Skip styling for whitespace and special characters
		if char == ' ' || char == '\n' {
			b.WriteRune(char)
			continue
		}

		// Calculate color index based on character position and frame
		colorIdx := (charIdx + c.Frame/c.Speed) % len(c.Colors)
		color := c.Colors[colorIdx]

		colored := lipgloss.NewStyle().
			Foreground(color).
			Bold(true).
			Render(string(char))

		b.WriteString(colored)
	}

	return b.String()
}

// RenderLines applies rainbow colors to multi-line text
// Each line gets different base color offset for a wave effect
func (c *Cycler) RenderLines(lines []string) string {
	var result strings.Builder

	for lineIdx, line := range lines {
		for charIdx, char := range line {
			if char == ' ' {
				result.WriteRune(char)
				continue
			}

			// Rainbow effect with line offset for vertical wave
			colorIdx := (charIdx + lineIdx + c.Frame/c.Speed) % len(c.Colors)
			color := c.Colors[colorIdx]

			colored := lipgloss.NewStyle().
				Foreground(color).
				Bold(true).
				Render(string(char))

			result.WriteString(colored)
		}
		if lineIdx < len(lines)-1 {
			result.WriteString("\n")
		}
	}

	return result.String()
}

// SetColors allows customizing the rainbow color palette
func (c *Cycler) SetColors(colors []lipgloss.Color) {
	if len(colors) > 0 {
		c.Colors = colors
	}
}

// SetSpeed updates the animation speed (higher = slower color changes)
func (c *Cycler) SetSpeed(speed int) {
	if speed > 0 {
		c.Speed = speed
	}
}

// GetColor returns the current color for a given index
// Useful for applying rainbow colors to UI elements
func (c *Cycler) GetColor(index int) lipgloss.Color {
	colorIdx := (index + c.Frame/c.Speed) % len(c.Colors)
	return c.Colors[colorIdx]
}
