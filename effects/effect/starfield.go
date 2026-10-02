package effect

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Star struct {
	x, y  float32
	speed float32
	color rl.Color
	size  float32
}

type Starfield struct {
	stars []Star
}

func NewStarfield(count int) *Starfield {
	stars := make([]Star, count)
	screenWidth := int32(rl.GetScreenWidth())
	screenHeight := int32(rl.GetScreenHeight())

	for i := 0; i < count; i++ {
		// Parallax speed based on "depth" (1 to 4)
		depth := float32(rl.GetRandomValue(10, 40)) / 10.0

		// Pick a "hot" color from the FirePalette
		colorIdx := int(rl.GetRandomValue(15, int32(len(FirePalette)-1)))

		stars[i] = Star{
			x:     float32(rl.GetRandomValue(0, screenWidth)),
			y:     float32(rl.GetRandomValue(0, screenHeight)),
			speed: depth * 0.5,
			size:  depth * 0.8,
			color: FirePalette[colorIdx],
		}
	}

	return &Starfield{
		stars: stars,
	}
}

func (s *Starfield) Process() {
	screenWidth := float32(rl.GetScreenWidth())
	for i := range s.stars {
		s.stars[i].x -= s.stars[i].speed
		if s.stars[i].x < 0 {
			s.stars[i].x = screenWidth
		}
	}
}

func (s *Starfield) Draw() {
	for _, star := range s.stars {
		// Draw the star with a small glow if it's "hot" (larger)
		if star.size > 2.0 {
			glowColor := star.color
			glowColor.A = 50
			rl.DrawCircleV(rl.Vector2{X: star.x, Y: star.y}, star.size*2.0, glowColor)
		}
		rl.DrawCircleV(rl.Vector2{X: star.x, Y: star.y}, star.size, star.color)
	}
}

var _ BaseEffect = (*Starfield)(nil)
