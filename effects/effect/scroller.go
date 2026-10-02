package effect

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Scroller struct {
	text         string
	font         rl.Font
	fontSize     float32
	x            float32
	yPos         float32
	speed        float32
	spacing      float32
	yAmplitude   float32
	pal          []rl.Color
	flashSpeed   int
	perCharColor bool
	frames       int
}

func NewScroller(txt string, pal []rl.Color, fontPath string, fontSize float32, yPos float32, speed float32, spacing float32, yAmp float32, flashSpeed int, perCharColor bool) *Scroller {
	return &Scroller{
		text:         txt,
		font:         rl.LoadFontEx(fontPath, int32(fontSize), nil, 0),
		fontSize:     fontSize,
		x:            float32(rl.GetScreenWidth()),
		yPos:         yPos,
		speed:        speed,
		spacing:      spacing,
		yAmplitude:   yAmp,
		pal:          pal,
		flashSpeed:   flashSpeed,
		perCharColor: perCharColor,
	}
}

func (s *Scroller) Process() {
	s.x -= s.speed
	s.frames++
	// Use spacing for reset logic
	if s.x < -float32(len(s.text))*s.spacing*1.5 {
		s.x = float32(rl.GetScreenWidth())
	}
}

func (s *Scroller) Draw() {
	t := rl.GetTime()
	for i, char := range s.text {
		y := s.yPos + float32(math.Sin(t*3+float64(i)*0.3))*s.yAmplitude

		// Calculate position
		posX := (s.x + float32(i)*s.spacing)

		// Only draw if visible
		if posX+s.spacing < 0 || posX > float32(rl.GetScreenWidth()) {
			continue
		}

		pos := rl.Vector2{X: posX, Y: y}

		// Draw 3D extrude
		layers := 5
		if s.fontSize > 300 {
			layers = 10
		}
		for l := layers; l > 0; l-- {
			offset := float32(l) * (s.fontSize / 120.0)
			layerColor := rl.Color{R: 40 + uint8(16*l), G: 0, B: 0, A: 255}
			rl.DrawTextEx(s.font, string(char), rl.Vector2{X: pos.X + offset, Y: pos.Y + offset}, s.fontSize, 2, layerColor)
		}

		// Draw main face
		numColors := len(s.pal)
		if numColors > 1 {
			// Update the color index based on flashSpeed
			divisor := s.flashSpeed
			if divisor < 1 {
				divisor = 1
			}
			step := float64(s.frames / divisor)
			phase := step * 0.8
			if s.perCharColor {
				phase += float64(i) * 0.2
			}
			norm := math.Abs(math.Mod(phase, 2.0) - 1.0)
			c := int(norm * float64(numColors-1))

			color := s.pal[c]

			// Pulse the brightness/alpha slightly
			pulse := float32(math.Sin(t*4.0)*0.2 + 0.8)
			color.R = uint8(float32(color.R) * pulse)
			color.G = uint8(float32(color.G) * pulse)
			color.B = uint8(float32(color.B) * pulse)

			rl.DrawTextEx(s.font, string(char), pos, s.fontSize, 2, color)
		} else if numColors == 1 {
			rl.DrawTextEx(s.font, string(char), pos, s.fontSize, 2, s.pal[0])
		}

		// Draw highlight
		rl.DrawTextEx(s.font, string(char), rl.Vector2{X: pos.X - 2, Y: pos.Y - 2}, s.fontSize, 2, rl.Color{255, 255, 255, 80})
	}
}
