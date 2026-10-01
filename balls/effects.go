package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func DrawGlitches() {
	if rl.GetRandomValue(0, 100) < 2 {
		glitchY := int32(rl.GetRandomValue(0, screenHeight))
		glitchH := int32(rl.GetRandomValue(5, 20))
		rl.DrawRectangle(0, glitchY, screenWidth, glitchH, rl.NewColor(255, 255, 255, 100))
	}
}

func DrawScanlines() {
	for y := 0; y < screenHeight; y += 3 {
		rl.DrawLine(0, int32(y), screenWidth, int32(y), rl.NewColor(0, 0, 0, 80))
	}
}

func DrawBorder() {
	rl.DrawRectangleLinesEx(rl.NewRectangle(0, 0, float32(screenWidth), float32(screenHeight)), 10, rl.NewColor(100, 100, 100, 255))
}
