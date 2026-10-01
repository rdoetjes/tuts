package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Scroller struct {
	text  string
	font  rl.Font
	timer float64
}

func NewScroller(font rl.Font, text string) *Scroller {
	return &Scroller{
		font: font,
		text: text,
	}
}

func (s *Scroller) Draw(timer float64, scrollPos float32) {
	for i, char := range s.text {
		charX := scrollPos + float32(i*50)

		if charX < -50 || charX > screenWidth {
			continue
		}

		charY := float32(screenHeight-120) + float32(math.Sin(timer*4+float64(i)*0.25)*60)
		pos := rl.NewVector2(charX, charY)

		t := timer*3.0 - float64(i)*0.2
		g := uint8(110 + 110*math.Cos(t))
		textColor := rl.NewColor(255, g, 0, 255)

		ot := timer*5.0 + float64(i)*0.3
		or := uint8(60 + 40*math.Sin(ot))
		outlineColor := rl.NewColor(or, 0, 0, 255)

		thickness := float32(5)
		offsets := []rl.Vector2{
			{X: -thickness, Y: -thickness}, {X: 0, Y: -thickness}, {X: thickness, Y: -thickness},
			{X: -thickness, Y: 0}, {X: thickness, Y: 0},
			{X: -thickness, Y: thickness}, {X: 0, Y: thickness}, {X: thickness, Y: thickness},
		}

		for _, off := range offsets {
			rl.DrawTextEx(s.font, string(char), rl.NewVector2(pos.X+off.X, pos.Y+off.Y), 60, 2, outlineColor)
		}

		rl.DrawTextEx(s.font, string(char), pos, 60, 2, textColor)
	}
}
