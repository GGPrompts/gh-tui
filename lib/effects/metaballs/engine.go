package metaballs

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Engine manages multiple blobs and renders the metaball effect
type Engine struct {
	Blobs  []*Blob
	Width  int
	Height int
	Frame  int

	// Rendering options
	GradientChars []string       // Characters to use for gradient (lightest to darkest)
	Thresholds    []float64      // Field strength thresholds for each gradient level
	DefaultColor  lipgloss.Color // Color for empty space
}

// NewEngine creates a new metaball engine with default settings
func NewEngine(width, height int) *Engine {
	return &Engine{
		Blobs:  make([]*Blob, 0),
		Width:  width,
		Height: height,
		Frame:  0,
		// Default gradient: light to dark using Unicode block characters
		GradientChars: []string{" ", "░", "▒", "▓", "█"},
		Thresholds:    []float64{0.3, 0.8, 1.5, 2.5},
		DefaultColor:  lipgloss.Color("0"),
	}
}

// AddBlob adds a blob to the engine
func (e *Engine) AddBlob(blob *Blob) {
	blob.colorIndex = len(e.Blobs)
	e.Blobs = append(e.Blobs, blob)
}

// Update advances the animation frame and updates all blobs
func (e *Engine) Update() {
	e.Frame++
	for i, blob := range e.Blobs {
		blob.Update(e.Frame, i, e.Width, e.Height)
	}
}

// Render generates the metaball effect as a string
func (e *Engine) Render() string {
	// Calculate field strength and closest blob for each position
	field := make([][]float64, e.Height)
	colorMap := make([][]int, e.Height)

	for y := 0; y < e.Height; y++ {
		field[y] = make([]float64, e.Width)
		colorMap[y] = make([]int, e.Width)

		for x := 0; x < e.Width; x++ {
			totalField := 0.0
			closestBlob := 0
			maxField := 0.0

			// Calculate combined field from all blobs
			for i, blob := range e.Blobs {
				blobField := blob.Field(float64(x), float64(y))
				totalField += blobField

				// Track which blob contributes most (for coloring)
				if blobField > maxField {
					maxField = blobField
					closestBlob = i
				}
			}

			field[y][x] = totalField
			colorMap[y][x] = closestBlob
		}
	}

	// Render with gradient characters
	var b strings.Builder

	for y := 0; y < e.Height; y++ {
		for x := 0; x < e.Width; x++ {
			strength := field[y][x]
			blobIndex := colorMap[y][x]

			// Choose character based on field strength
			char := e.getGradientChar(strength)

			// Color from closest blob
			var color lipgloss.Color
			if char == " " {
				color = e.DefaultColor
			} else if blobIndex < len(e.Blobs) {
				color = e.Blobs[blobIndex].Color
			} else {
				color = e.DefaultColor
			}

			styled := lipgloss.NewStyle().
				Foreground(color).
				Render(char)

			b.WriteString(styled)
		}
		if y < e.Height-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

// getGradientChar returns the appropriate character for a given field strength
func (e *Engine) getGradientChar(strength float64) string {
	if strength < e.Thresholds[0] {
		return e.GradientChars[0] // Empty space
	}

	for i, threshold := range e.Thresholds {
		if strength < threshold {
			return e.GradientChars[i]
		}
	}

	// Strongest field - use darkest character
	return e.GradientChars[len(e.GradientChars)-1]
}

// Resize updates the engine dimensions
func (e *Engine) Resize(width, height int) {
	e.Width = width
	e.Height = height
}

// SetGradient allows customizing the gradient characters and thresholds
func (e *Engine) SetGradient(chars []string, thresholds []float64) {
	if len(chars) > 0 {
		e.GradientChars = chars
	}
	if len(thresholds) > 0 {
		e.Thresholds = thresholds
	}
}
