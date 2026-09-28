package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 800
	screenHeight = 600
	fov          = 400.0
	viewerDist   = 5.0
)

func main() {
	rl.InitWindow(screenWidth, screenHeight, "3D Vector Art - R A Y")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	lettersMap := GetLetters()

	// Define positions for each letter
	type LetterState struct {
		char   string
		pos    Vec3
		rotate Vec3
	}

	states := []LetterState{
		{char: "R", pos: Vec3{X: -3, Y: 0, Z: 0}},
		{char: "A", pos: Vec3{X: 0, Y: 0, Z: 0}},
		{char: "Y", pos: Vec3{X: 3, Y: 0, Z: 0}},
	}

	angle := 0.0

	for !rl.WindowShouldClose() {
		// Update
		angle += 0.02

		// Draw
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		rl.DrawText("3D Vector Art using 2D Lines", 10, 10, 20, rl.Gray)

		for i, state := range states {
			letter := lettersMap[state.char]

			// Individual rotations for a "twist and tumble" effect
			rotX := angle * (0.5 + float64(i)*0.1)
			rotY := angle * (0.8 + float64(i)*0.1)
			rotZ := angle * (0.3 + float64(i)*0.1)

			// Floating effect
			offsetX := math.Sin(angle+float64(i)) * 0.5
			offsetY := math.Cos(angle*0.7+float64(i)) * 0.5

			for _, line := range letter.Lines {
				// 1. Rotate
				p1 := transform(line.P1, rotX, rotY, rotZ)
				p2 := transform(line.P2, rotX, rotY, rotZ)

				// 2. Translate to its position + offset
				p1.X += state.pos.X + offsetX
				p1.Y += state.pos.Y + offsetY
				p1.Z += state.pos.Z

				p2.X += state.pos.X + offsetX
				p2.Y += state.pos.Y + offsetY
				p2.Z += state.pos.Z

				// 3. Project to 2D
				v1 := p1.Project(float64(screenWidth), float64(screenHeight), fov, viewerDist)
				v2 := p2.Project(float64(screenWidth), float64(screenHeight), fov, viewerDist)

				// 4. Draw
				rl.DrawLine(int32(v1.X), int32(v1.Y), int32(v2.X), int32(v2.Y), rl.RayWhite)
			}
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
