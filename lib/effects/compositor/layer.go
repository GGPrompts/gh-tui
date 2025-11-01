package compositor

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Layer represents a renderable layer that can be composited
type Layer interface {
	Render() string
	Width() int
	Height() int
}

// Compositor manages multiple layers and composites them together
type Compositor struct {
	layers []Layer
	width  int
	height int
}

// NewCompositor creates a new layer compositor
func NewCompositor(width, height int) *Compositor {
	return &Compositor{
		layers: make([]Layer, 0),
		width:  width,
		height: height,
	}
}

// AddLayer adds a layer to the compositor (layers are drawn in order added)
func (c *Compositor) AddLayer(layer Layer) {
	c.layers = append(c.layers, layer)
}

// Composite renders all layers from bottom to top
// Uses ANSI-aware overlaying to properly handle styled text
func (c *Compositor) Composite() string {
	if len(c.layers) == 0 {
		return ""
	}

	// Start with the bottom layer
	result := c.layers[0].Render()
	resultLines := strings.Split(result, "\n")

	// Ensure we have enough lines
	for len(resultLines) < c.height {
		resultLines = append(resultLines, strings.Repeat(" ", c.width))
	}

	// Overlay each subsequent layer
	for i := 1; i < len(c.layers); i++ {
		layer := c.layers[i]
		layerContent := layer.Render()
		layerLines := strings.Split(layerContent, "\n")

		// Center the layer if it's smaller than the compositor
		startY := (c.height - len(layerLines)) / 2
		if startY < 0 {
			startY = 0
		}

		// Overlay each line
		for j, layerLine := range layerLines {
			y := startY + j
			if y >= 0 && y < len(resultLines) {
				// Center horizontally
				layerWidth := lipgloss.Width(layerLine)
				startX := (c.width - layerWidth) / 2
				if startX < 0 {
					startX = 0
				}

				resultLines[y] = overlayString(resultLines[y], layerLine, startX, c.width)
			}
		}
	}

	return strings.Join(resultLines, "\n")
}

// Resize updates the compositor dimensions
func (c *Compositor) Resize(width, height int) {
	c.width = width
	c.height = height
}

// Clear removes all layers
func (c *Compositor) Clear() {
	c.layers = make([]Layer, 0)
}

// overlayString overlays src onto dst at position x
// Preserves ANSI escape codes and handles visual width properly
func overlayString(dst, src string, x, maxWidth int) string {
	if x < 0 {
		return dst
	}

	// Get visual widths (ignoring ANSI codes)
	dstWidth := lipgloss.Width(dst)
	srcWidth := lipgloss.Width(src)

	var result strings.Builder

	// Left side: extract x visible characters from dst
	if x > 0 {
		leftPart := extractVisibleChars(dst, 0, x)
		result.WriteString(leftPart)
	}

	// Middle: the overlay
	result.WriteString(src)

	// Right side: extract remaining characters from dst after the overlay
	rightStart := x + srcWidth
	if rightStart < dstWidth {
		rightPart := extractVisibleChars(dst, rightStart, dstWidth-rightStart)
		result.WriteString(rightPart)
	}

	return result.String()
}

// extractVisibleChars extracts count visible characters starting from position start
// Properly handles ANSI escape codes
func extractVisibleChars(s string, start, count int) string {
	if count <= 0 {
		return ""
	}

	runes := []rune(s)
	visibleCount := 0
	inEscape := false
	startIdx := -1
	endIdx := -1

	for i, r := range runes {
		// Track ANSI escape sequences
		if r == '\x1b' {
			inEscape = true
		} else if inEscape && r == 'm' {
			inEscape = false
			continue
		}

		// Only count visible characters
		if !inEscape {
			if visibleCount == start && startIdx == -1 {
				startIdx = i
			}
			if visibleCount >= start {
				visibleCount++
				if visibleCount-start >= count {
					endIdx = i + 1
					break
				}
			} else {
				visibleCount++
			}
		}
	}

	if startIdx == -1 {
		return ""
	}
	if endIdx == -1 {
		endIdx = len(runes)
	}

	return string(runes[startIdx:endIdx])
}
