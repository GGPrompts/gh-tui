# TUI Effects Library

**Beautiful, physics-based animations for terminal user interfaces.**

This library provides reusable animation effects extracted from the Balatro TUI project's landing page. Create stunning visual effects that rival modern GUIs - all in your terminal!

## 🎨 Features

- **🔮 Metaballs** - Lava lamp-style floating blobs with physics simulation
- **🌊 Wave Effects** - Sine wave distortions for grids and content
- **🌈 Rainbow Cycling** - Animated color gradients for text
- **🎭 Layer Compositor** - ANSI-aware multi-layer rendering

## 📦 Installation

```go
import (
    "github.com/GGPrompts/TUITemplate/lib/effects/metaballs"
    "github.com/GGPrompts/TUITemplate/lib/effects/waves"
    "github.com/GGPrompts/TUITemplate/lib/effects/rainbow"
    "github.com/GGPrompts/TUITemplate/lib/effects/compositor"
)
```

## 🚀 Quick Start

### Metaballs - Loading Spinner

```go
// Create engine
engine := metaballs.NewEngine(width, height)

// Add colorful blobs
engine.AddBlob(metaballs.NewBlob(
    x, y,           // Position
    vx, vy,         // Velocity
    radius,         // Size
    color,          // Lipgloss color
))

// In your update loop
engine.Update()

// In your view
return engine.Render()
```

### Wave Distortion - Animated Grid

```go
// Create wavy grid
grid := waves.NewGrid(width, height)

// Customize appearance
grid.SetGridSize(10)
grid.SetColors(waves.GridColors{
    Intersection: lipgloss.Color("129"), // Purple
    Vertical:     lipgloss.Color("61"),
    Horizontal:   lipgloss.Color("61"),
})

// Animate
grid.Update()

// Render
return grid.Render()
```

### Rainbow Text - Animated Colors

```go
// Create cycler
cycler := rainbow.NewCycler()

// Animate
cycler.Update()

// Apply to text
rainbowText := cycler.Render("HELLO WORLD")

// Or apply to multiple lines with wave effect
rainbowLines := cycler.RenderLines([]string{
    "Line 1",
    "Line 2",
    "Line 3",
})
```

### Layer Compositor - Combine Effects

```go
// Create compositor
comp := compositor.NewCompositor(width, height)

// Add layers (bottom to top)
comp.AddLayer(backgroundLayer)
comp.AddLayer(metaballLayer)
comp.AddLayer(textLayer)

// Render composite
return comp.Composite()
```

## 📚 Examples

Check out the `examples/effects/` directory for complete working examples:

- **metaball-spinner** - Loading screen with floating blobs
- **wavy-menu** - Menu with animated wave background
- **rainbow-text** - Color-cycling text animations
- **landing-page** - Full demo combining all effects

### Running Examples

```bash
cd examples/effects/metaball-spinner
go run main.go
```

## 🎭 Effect Details

### Metaballs

Physics-based blob simulation with:
- Real metaball algorithm (field strength = radius² / distance²)
- Organic motion with wobble
- Gradient rendering using Unicode block characters
- Customizable colors and sizes

**Perfect for:**
- Loading screens
- Background effects
- Animated logos
- Visual feedback

### Waves

Sine wave distortion effects:
- Animated grid with flowing lines
- Customizable amplitude, frequency, and speed
- Apply distortion to any coordinates
- Smooth, organic motion

**Perfect for:**
- Background grids
- Menu animations
- Text effects
- Flowing layouts

### Rainbow

Color cycling effects:
- Smooth rainbow gradients
- Character-by-character coloring
- Vertical wave patterns for multi-line text
- Customizable color palettes

**Perfect for:**
- Title screens
- Highlighted text
- Success messages
- Brand effects

### Compositor

Layer management system:
- ANSI-aware overlaying
- Properly handles styled text
- Multiple layer support
- Automatic centering

**Perfect for:**
- Combining multiple effects
- Overlay text on backgrounds
- Complex UI layouts
- Professional polish

## 🔧 Advanced Usage

### Custom Gradient Characters

```go
engine := metaballs.NewEngine(width, height)

// Use different characters for gradient
engine.SetGradient(
    []string{" ", ".", ":", ";", "#", "@"},
    []float64{0.2, 0.5, 1.0, 1.5, 2.0},
)
```

### Custom Wave Parameters

```go
distortion := waves.NewDistortion()
distortion.SetAmplitude(3.0)   // Bigger waves
distortion.SetFrequency(10.0)  // Tighter waves
distortion.SetSpeed(30.0)      // Faster animation
```

### Custom Rainbow Colors

```go
cycler := rainbow.NewCycler()

// Use your brand colors
cycler.SetColors([]lipgloss.Color{
    lipgloss.Color("#FF0000"),
    lipgloss.Color("#00FF00"),
    lipgloss.Color("#0000FF"),
})

cycler.SetSpeed(3) // Faster color changes
```

## 💡 Pro Tips

1. **Performance**: Keep blob count < 10 for smooth 60fps on most terminals
2. **Colors**: Use ANSI 256 colors for best compatibility
3. **Grid Size**: Smaller grid sizes (5-10) look better than larger ones
4. **Compositing**: Order matters - add background layers first
5. **Frame Rate**: 20-30 fps is plenty for smooth animations

## 🎨 Color Palettes

### Neon (Default)
```go
Gold:    lipgloss.Color("226")
Cyan:    lipgloss.Color("51")
Magenta: lipgloss.Color("201")
Green:   lipgloss.Color("46")
Purple:  lipgloss.Color("129")
```

### Pastel
```go
Pink:   lipgloss.Color("213")
Blue:   lipgloss.Color("117")
Green:  lipgloss.Color("156")
Yellow: lipgloss.Color("228")
```

### Fire
```go
Red:    lipgloss.Color("196")
Orange: lipgloss.Color("208")
Yellow: lipgloss.Color("226")
White:  lipgloss.Color("255")
```

## 🏗️ Architecture

All effects follow the same pattern:

```go
type Effect struct {
    Frame int
    // ... effect-specific fields
}

func NewEffect(params) *Effect
func (e *Effect) Update()
func (e *Effect) Render() string
func (e *Effect) Resize(width, height int)
```

This makes them easy to integrate into any Bubbletea application!

## 🤝 Contributing

Want to add more effects? PRs welcome!

Potential ideas:
- Particle systems
- Matrix rain effect
- Fire/smoke simulation
- Audio visualizer bars
- Starfield backgrounds

## 📖 Credits

Effects extracted from **TUI Classics - Balatro** landing page.

Original implementation by the TUI Classics team, inspired by the mesmerizing Balatro game by LocalThunk.

## 📄 License

MIT License - Use these effects in your TUI apps, commercial or otherwise!

---

**Made with ❤️ for the terminal community**

*Because terminals can be beautiful too.*
