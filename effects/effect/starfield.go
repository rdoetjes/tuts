package effect

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Ember struct {
	x, y  float32
	speed float32
	drift float32 // Horizontal wiggling
	color rl.Color
	size  float32
}

type EmberField struct {
	embers []Ember
}

func NewStarfield(count int) *EmberField {
	return NewEmberField(count)
}

func NewEmberField(count int) *EmberField {
	embers := make([]Ember, count)
	screenWidth := int32(rl.GetScreenWidth())
	screenHeight := int32(rl.GetScreenHeight())

	for i := 0; i < count; i++ {
		depth := float32(rl.GetRandomValue(10, 40)) / 10.0
		// Pick hot colors
		colorIdx := int(rl.GetRandomValue(int32(len(FirePalette)/2), int32(len(FirePalette)-1)))

		embers[i] = Ember{
			x:     float32(rl.GetRandomValue(0, screenWidth)),
			y:     float32(rl.GetRandomValue(0, screenHeight)),
			speed: depth * 0.4,
			drift: float32(rl.GetRandomValue(-100, 100)) / 100.0,
			size:  depth * 0.7,
			color: FirePalette[colorIdx],
		}
	}

	return &EmberField{
		embers: embers,
	}
}

func (s *EmberField) Process() {
	screenWidth := float32(rl.GetScreenWidth())
	screenHeight := float32(rl.GetScreenHeight())
	t := float64(rl.GetTime())

	for i := range s.embers {
		// Flow up
		s.embers[i].y -= s.embers[i].speed

		// Wiggle/Drift using a sine wave
		s.embers[i].x += float32(math.Sin(t*1.5+float64(i))) * 0.5
		s.embers[i].x += s.embers[i].drift * 0.2

		// Cull and recycle when off screen
		if s.embers[i].y < -20 {
			s.embers[i].y = screenHeight + 20
			s.embers[i].x = float32(rl.GetRandomValue(0, int32(screenWidth)))
		}
		// Wrap horizontally too
		if s.embers[i].x < -20 {
			s.embers[i].x = screenWidth + 20
		} else if s.embers[i].x > screenWidth+20 {
			s.embers[i].x = -20
		}
	}
}

func (s *EmberField) Draw() {
	for _, ember := range s.embers {
		// Subtle glow for hot embers
		if ember.size > 1.5 {
			glowColor := ember.color
			glowColor.A = 40
			rl.DrawCircleV(rl.Vector2{X: ember.x, Y: ember.y}, ember.size*3.0, glowColor)
		}
		rl.DrawCircleV(rl.Vector2{X: ember.x, Y: ember.y}, ember.size, ember.color)
	}
}

var _ BaseEffect = (*EmberField)(nil)
