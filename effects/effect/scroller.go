package effect

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Scroller struct {
	text string
	font rl.Font
	x    float32
	pal  []rl.Color
}

func NewScroller(txt string, pal []rl.Color) *Scroller {
	return &Scroller{text: txt, font: rl.LoadFontEx("assets/fonts/Menlo.ttf", 480, nil, 0), x: float32(rl.GetScreenWidth()), pal: pal}
}

func (s *Scroller) Process() {
	s.x -= 8
	// Use fixed width for reset logic: len(text) * spacing
	if s.x < -float32(len(s.text)*280) {
		s.x = float32(rl.GetScreenWidth())
	}
}

func (s *Scroller) Draw() {
	t := rl.GetTime()
	charStep := float32(175) // Fixed width for every character
	for i, char := range s.text {
		y := float32(rl.GetScreenHeight())/2 - 240 + float32(math.Sin(t*3+float64(i)*0.3))*150

		// Calculate fixed position
		posX := (s.x + float32(i)*charStep)

		// Only draw if visible
		if posX+charStep < 0 || posX > float32(rl.GetScreenWidth()) {
			continue
		}

		pos := rl.Vector2{X: posX, Y: y}

		// Draw 3D extrude
		layers := 10
		for l := layers; l > 0; l-- {
			offset := float32(l) * 4.0
			layerColor := rl.Color{R: 40 + uint8(16*l), G: 0, B: 0, A: 255}
			rl.DrawTextEx(s.font, string(char), rl.Vector2{X: pos.X + offset, Y: pos.Y + offset}, 480, 2, layerColor)
		}

		// Draw main face
		// Cycle colors up and down the palette (ping-pong effect)
		numColors := len(s.pal)
		if numColors > 1 {
			cycleLength := 2 * (numColors - 1)
			idx := (int(t*20) + i) % cycleLength
			c := idx
			if idx >= numColors {
				c = cycleLength - idx
			}
			color := s.pal[c]
			rl.DrawTextEx(s.font, string(char), pos, 480, 2, color)
		}

		// Draw highlight
		rl.DrawTextEx(s.font, string(char), rl.Vector2{X: pos.X - 3, Y: pos.Y - 3}, 480, 2, rl.Color{255, 255, 255, 80})
	}
}
