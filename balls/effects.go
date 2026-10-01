package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func DrawBorder() {
	rl.DrawRectangleLinesEx(rl.NewRectangle(0, 0, float32(screenWidth), float32(screenHeight)), 10, rl.NewColor(100, 100, 100, 255))
}
