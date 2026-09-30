package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Scroller struct {
	X        float32
	Text     string
	FontSize float32
}

func (sc *Scroller) Update(s *Scene) {
	sc.X -= 4.5
	textSize := rl.MeasureTextEx(s.font, sc.Text, sc.FontSize, 2)
	if sc.X < -textSize.X {
		sc.X = float32(s.screenWidth)
	}
}

func (sc *Scroller) Draw(s *Scene) {
	currentX := sc.X
	y := float32(s.screenHeight) - 120
	time := rl.GetTime()

	for i, char := range sc.Text {
		charStr := string(char)
		waveY := y + float32(math.Sin(float64(i)*0.15+time*5.0)*25.0)

		val := math.Mod(time*2.0+float64(i)*0.08, 1.0)
		if val < 0 {
			val += 1.0
		}

		idx := val * float64(len(MoodyPalette)-1)
		i1 := int(math.Floor(idx))
		i2 := i1 + 1
		frac := float32(idx - float64(i1))

		color := rl.ColorLerp(MoodyPalette[i1], MoodyPalette[i2], frac)
		position := rl.NewVector2(currentX, waveY)

		rl.DrawTextEx(s.font, charStr, rl.NewVector2(position.X+3, position.Y+3), sc.FontSize, 2, rl.NewColor(50, 0, 50, 200))
		rl.DrawTextEx(s.font, charStr, position, sc.FontSize, 2, color)

		charWidth := rl.MeasureTextEx(s.font, charStr, sc.FontSize, 2).X
		currentX += charWidth
	}
}

func (s *Scene) drawRasterBars() {
	time := rl.GetTime()
	barHeight := 60
	for i := 0; i < 3; i++ {
		offset := float64(i) * 0.8
		yCenter := float32(s.screenHeight/2) + float32(math.Sin(time*1.5+offset)*(s.screenHeight*0.35))

		for j := 0; j < barHeight; j++ {
			factor := 1.0 - math.Abs(float64(j-barHeight/2))/float64(barHeight/2)
			palVal := math.Mod(float64(i)*0.33+factor*0.5, 1.0)
			idx := palVal * float64(len(MoodyPalette)-1)
			i1 := int(math.Floor(idx))
			i2 := i1 + 1
			frac := float32(idx - float64(i1))

			baseColor := rl.ColorLerp(MoodyPalette[i1], MoodyPalette[i2], frac)
			color := baseColor
			color.R = uint8(float32(color.R) * float32(factor))
			color.G = uint8(float32(color.G) * float32(factor))
			color.B = uint8(float32(color.B) * float32(factor))

			yPos := int32(yCenter) + int32(j-barHeight/2)
			rl.DrawRectangle(0, yPos, int32(s.screenWidth), 1, color)
		}
	}
}
