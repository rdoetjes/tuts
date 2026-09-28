package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	const fov = 500.0
	viewerDist := 10.0

	// Set fullscreen flag before initialization
	rl.SetConfigFlags(rl.FlagFullscreenMode)

	rl.InitWindow(0, 0, "3D Vector Art - R A Y")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)
	rl.HideCursor()

	word := GetWord()
	screenWidth := rl.GetScreenWidth()
	screenHeight := rl.GetScreenHeight()
	angle := 0.0

	for !rl.WindowShouldClose() {
		// Update
		angle += 0.015
		// Draw
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		// Entire word rotations
		rotX := angle * 0.4
		rotY := angle * 0.6
		rotZ := angle * 0.2

		viewerDist = 30 + (20 * (math.Sin(rl.GetTime()))) // pulse the distance in and out
		for _, line := range word.Lines {
			// 1. Rotate
			p1 := transform(line.P1, rotX, rotY, rotZ)
			p2 := transform(line.P2, rotX, rotY, rotZ)

			// 2. Project to 2D
			v1 := p1.Project(float64(screenWidth), float64(screenHeight), fov, viewerDist)
			v2 := p2.Project(float64(screenWidth), float64(screenHeight), fov, viewerDist)

			// 3. Draw
			color := rl.RayWhite
			rl.DrawLine(int32(v1.X), int32(v1.Y), int32(v2.X), int32(v2.Y), color)
		}

		rl.EndDrawing()
	}
}

func transform(v Vec3, rx, ry, rz float64) Vec3 {
	v = v.RotateX(rx)
	v = v.RotateY(ry)
	v = v.RotateZ(rz)
	return v
}
