package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Header struct {
	font rl.Font
}

func (h *Header) DrawLogo(timer float64, pulse float32) {
	headerText := "CRU JONES & PHONAX '89"
	fontSize := float32(80 + pulse*8)

	textSize := rl.MeasureTextEx(h.font, headerText, fontSize, 2)
	headerX := (float32(screenWidth) - textSize.X) / 2

	rl.DrawTextEx(h.font, headerText, rl.NewVector2(headerX+6, 56), fontSize, 2, rl.Maroon)
	rl.DrawTextEx(h.font, headerText, rl.NewVector2(headerX, 50), fontSize, 2, rl.Gold)
}

func (h *Header) DrawSubHeader(timer float64, isFlatline bool) {
	subText := "<< CRACKED BY D'ELITE >>"
	baseFontSize := float32(35)
	spacing := float32(2)

	totalWidth := float32(0)
	for _, char := range subText {
		charSize := rl.MeasureTextEx(h.font, string(char), baseFontSize, spacing)
		totalWidth += charSize.X
	}

	startX := (float32(screenWidth) - totalWidth) / 2
	baseY := float32(135)

	for i, char := range subText {
		charTimer := timer*5.0 + float64(i)*0.5
		offsetY := float32(math.Sin(charTimer)) * 15.0
		offsetX := float32(math.Cos(charTimer*0.8)) * 5.0

		charPulse := float32(math.Sin(charTimer*0.7))*0.2 + 1.0
		charFontSize := baseFontSize * charPulse

		if isFlatline {
			charPulse = 0.0
			charFontSize = baseFontSize
			offsetX = 0.0
			offsetY = 0.0
		}

		t := float64(math.Sin(timer*3.0+float64(i)*0.3))*0.5 + 0.5
		g := uint8(t * 220)
		textColor := rl.NewColor(255, g, 0, 255)

		charStr := string(char)
		charSize := rl.MeasureTextEx(h.font, charStr, charFontSize, spacing)
		pos := rl.NewVector2(startX+offsetX, baseY+offsetY-charSize.Y/2)

		rl.DrawTextEx(h.font, charStr, rl.NewVector2(pos.X+3, pos.Y+3), charFontSize, spacing, rl.NewColor(0, 0, 0, 150))

		glowColor := rl.NewColor(255, g, 0, 80)
		rl.DrawTextEx(h.font, charStr, rl.NewVector2(pos.X-1, pos.Y-1), charFontSize, spacing, glowColor)
		rl.DrawTextEx(h.font, charStr, rl.NewVector2(pos.X+1, pos.Y+1), charFontSize, spacing, glowColor)

		rl.DrawTextEx(h.font, charStr, pos, charFontSize, spacing, textColor)

		baseCharSize := rl.MeasureTextEx(h.font, charStr, baseFontSize, spacing)
		startX += baseCharSize.X
	}
}
