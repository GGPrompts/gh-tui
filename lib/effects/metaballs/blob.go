package metaballs

import (
	"math"

	"github.com/charmbracelet/lipgloss"
)

// Blob represents a single floating blob for metaball effects
type Blob struct {
	X, Y       float64        // Position
	VX, VY     float64        // Velocity
	Radius     float64        // Size
	Color      lipgloss.Color // Color
	colorIndex int            // Internal tracking
}

// NewBlob creates a new blob with the given parameters
func NewBlob(x, y, vx, vy, radius float64, color lipgloss.Color) *Blob {
	return &Blob{
		X:      x,
		Y:      y,
		VX:     vx,
		VY:     vy,
		Radius: radius,
		Color:  color,
	}
}

// Update moves the blob based on its velocity and applies wobble
func (b *Blob) Update(frame int, index int, width, height int) {
	// Update position
	b.X += b.VX
	b.Y += b.VY

	// Bounce off edges
	if b.X < 0 || b.X > float64(width) {
		b.VX = -b.VX
	}
	if b.Y < 0 || b.Y > float64(height) {
		b.VY = -b.VY
	}

	// Add organic wobble using sine/cosine
	wobbleX := math.Sin(float64(frame+index*37)/30.0) * 0.05
	wobbleY := math.Cos(float64(frame+index*41)/25.0) * 0.05
	b.VX += wobbleX
	b.VY += wobbleY

	// Damping to prevent excessive speed
	b.VX *= 0.99
	b.VY *= 0.99
}

// Field calculates the metaball field strength at a given point
// Uses the formula: (radius^2) / (distance^2)
func (b *Blob) Field(x, y float64) float64 {
	dx := x - b.X
	dy := y - b.Y
	distance := math.Sqrt(dx*dx + dy*dy)

	if distance == 0 {
		return b.Radius * b.Radius
	}

	return (b.Radius * b.Radius) / (distance * distance)
}
