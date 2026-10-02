package effect

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Moon struct {
	radius float32
	timer  float32
}

func NewMoon(radius float32) *Moon {
	return &Moon{
		radius: radius,
		timer:  0,
	}
}

func (m *Moon) Process() {
	m.timer += 0.01
}

func (m *Moon) Draw() {
	centerX := float32(rl.GetScreenWidth()) / 2
	centerY := float32(rl.GetScreenHeight()) / 2

	// Dynamic pulse for the whole moon
	globalPulse := float32(math.Sin(float64(m.timer*8.0))*0.15 + 0.85)

	// Large, soft eerie glow - Deep Blood Red
	for i := 45; i > 0; i-- {
		glowRadius := m.radius + float32(i)*7.0 + (globalPulse * 20.0)
		alpha := uint8(float32(45-i) * 1.8 * globalPulse)
		rl.DrawCircleV(rl.Vector2{X: centerX, Y: centerY}, glowRadius, rl.Color{120, 0, 0, alpha})
	}
}

var _ BaseEffect = (*Moon)(nil)
