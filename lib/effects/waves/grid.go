package waves

import (
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Grid represents an animated grid with sine wave distortion
type Grid struct {
	Width    int
	Height   int
	Frame    int
	GridSize int            // Distance between grid lines
	Colors   GridColors     // Colors for different grid elements
}

// GridColors defines the color scheme for the wavy grid
type GridColors struct {
	Intersection lipgloss.Color // Color for grid intersections (+)
	Vertical     lipgloss.Color // Color for vertical lines (│)
	Horizontal   lipgloss.Color // Color for horizontal lines (─)
	Background   lipgloss.Color // Background color
}

// DefaultGridColors returns a purple-themed color scheme
func DefaultGridColors() GridColors {
	return GridColors{
		Intersection: lipgloss.Color("129"), // Purple
		Vertical:     lipgloss.Color("61"),  // Dark purple
		Horizontal:   lipgloss.Color("61"),  // Dark purple
		Background:   lipgloss.Color("0"),   // Black
	}
}

// NewGrid creates a new wavy grid with default settings
func NewGrid(width, height int) *Grid {
	return &Grid{
		Width:    width,
		Height:   height,
		Frame:    0,
		GridSize: 10,
		Colors:   DefaultGridColors(),
	}
}

// Update advances the animation frame
func (g *Grid) Update() {
	g.Frame++
}

// Render generates the wavy grid as a string
func (g *Grid) Render() string {
	var b strings.Builder

	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			// Calculate wave offset using sine waves
			waveX := math.Sin(float64(y)/5.0+float64(g.Frame)/20.0) * 2
			waveY := math.Sin(float64(x)/5.0+float64(g.Frame)/20.0) * 2

			// Apply wave distortion to grid coordinates
			gridX := int(float64(x) + waveX)
			gridY := int(float64(y) + waveY)

			// Determine if this position should be a grid line
			isGridLine := (gridX%g.GridSize == 0) || (gridY%g.GridSize == 0)

			var char string
			var color lipgloss.Color

			if isGridLine {
				if gridX%g.GridSize == 0 && gridY%g.GridSize == 0 {
					char = "+"
					color = g.Colors.Intersection
				} else if gridX%g.GridSize == 0 {
					char = "│"
					color = g.Colors.Vertical
				} else {
					char = "─"
					color = g.Colors.Horizontal
				}
			} else {
				char = " "
				color = g.Colors.Background
			}

			styled := lipgloss.NewStyle().
				Foreground(color).
				Render(char)

			b.WriteString(styled)
		}
		if y < g.Height-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

// Resize updates the grid dimensions
func (g *Grid) Resize(width, height int) {
	g.Width = width
	g.Height = height
}

// SetColors updates the color scheme
func (g *Grid) SetColors(colors GridColors) {
	g.Colors = colors
}

// SetGridSize updates the distance between grid lines
func (g *Grid) SetGridSize(size int) {
	if size > 0 {
		g.GridSize = size
	}
}
