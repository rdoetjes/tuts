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
	yAmplitude   float32
	pal          []rl.Color
	flashSpeed   int
	perCharColor bool
	frames       int
	charOffsets  []float32
	totalWidth   float32
}

func NewScroller(txt string, pal []rl.Color, fontPath string, fontSize float32, yPos float32, speed float32, yAmp float32, flashSpeed int, perCharColor bool) *Scroller {
	font := rl.LoadFontEx(fontPath, int32(fontSize), nil, 0)

	// Pre-calculate character offsets for proportional spacing
	charOffsets := make([]float32, len(txt))
	var currentX float32 = 0
	for i, char := range txt {
		charOffsets[i] = currentX
		// Measure each character to find its width
		charSize := rl.MeasureTextEx(font, string(char), fontSize, 2)
		currentX += charSize.X
	}

	return &Scroller{
		text:         txt,
		font:         font,
		fontSize:     fontSize,
		x:            float32(rl.GetScreenWidth()),
		yPos:         yPos,
		speed:        speed,
		yAmplitude:   yAmp,
		pal:          pal,
		flashSpeed:   flashSpeed,
		perCharColor: perCharColor,
		charOffsets:  charOffsets,
		totalWidth:   currentX,
	}
}

func (s *Scroller) Process() bool {
	s.x -= s.speed
	s.frames++
	// Reset when the entire string has scrolled off
	if s.x < -s.totalWidth {
		s.x = float32(rl.GetScreenWidth())
		return true
	}
	return false
}

func (s *Scroller) GetCharX(index int) float32 {
	if index < 0 || index >= len(s.charOffsets) {
		return -1000
	}
	return s.x + s.charOffsets[index]
}

func (s *Scroller) Draw() {
	t := rl.GetTime()
	for i, char := range s.text {
		y := s.yPos + float32(math.Sin(t*3+float64(i)*0.3))*s.yAmplitude

		// Calculate position using pre-calculated proportional offsets
		posX := s.x + s.charOffsets[i]

		// Get character width for visibility check
		charWidth := rl.MeasureTextEx(s.font, string(char), s.fontSize, 2).X

		// Only draw if visible
		if posX+charWidth < 0 || posX > float32(rl.GetScreenWidth()) {
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
