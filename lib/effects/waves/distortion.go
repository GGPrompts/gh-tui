package waves

import "math"

// Distortion applies sine wave distortion to coordinates
type Distortion struct {
	Frame      int
	Amplitude  float64 // Wave height (default: 2.0)
	Frequency  float64 // Wave frequency (default: 5.0)
	Speed      float64 // Animation speed (default: 20.0)
}

// NewDistortion creates a new wave distortion effect
func NewDistortion() *Distortion {
	return &Distortion{
		Frame:     0,
		Amplitude: 2.0,
		Frequency: 5.0,
		Speed:     20.0,
	}
}

// Update advances the animation frame
func (d *Distortion) Update() {
	d.Frame++
}

// ApplyX calculates the X distortion for a given Y position
func (d *Distortion) ApplyX(y int) float64 {
	return math.Sin(float64(y)/d.Frequency+float64(d.Frame)/d.Speed) * d.Amplitude
}

// ApplyY calculates the Y distortion for a given X position
func (d *Distortion) ApplyY(x int) float64 {
	return math.Sin(float64(x)/d.Frequency+float64(d.Frame)/d.Speed) * d.Amplitude
}

// Apply calculates both X and Y distortion for a given position
func (d *Distortion) Apply(x, y int) (float64, float64) {
	return d.ApplyX(y), d.ApplyY(x)
}

// SetAmplitude updates the wave height
func (d *Distortion) SetAmplitude(amplitude float64) {
	d.Amplitude = amplitude
}

// SetFrequency updates the wave frequency (higher = tighter waves)
func (d *Distortion) SetFrequency(frequency float64) {
	if frequency > 0 {
		d.Frequency = frequency
	}
}

// SetSpeed updates the animation speed (higher = faster)
func (d *Distortion) SetSpeed(speed float64) {
	if speed > 0 {
		d.Speed = speed
	}
}
