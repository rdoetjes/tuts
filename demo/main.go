package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 800
	screenHeight = 600
	fov          = 400.0
)

func main() {
	viewerDist := 10.0
	rl.InitWindow(screenWidth, screenHeight, "3D Vector Art - R A Y")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	word := GetWord()

	angle := 0.0

	for !rl.WindowShouldClose() {
		// Update
		angle += 0.015

		// Draw
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		rl.DrawText("3D Vector Word - R A Y", 10, 10, 20, rl.Gray)

		// Entire word rotations
		rotX := angle * 0.4
		rotY := angle * 0.6
		rotZ := angle * 0.2

		// Floating effect for the whole word
		offsetY := math.Sin(angle) * 0.01

		for _, line := range word.Lines {
			//rotZ = math.Sin(rl.GetTime()) * 0.1

			// 1. Rotate
			p1 := transform(line.P1, rotX, rotY, rotZ)
			p2 := transform(line.P2, rotX, rotY, rotZ)

			// 2. Center the word and add float
			p1.Y += offsetY
			p2.Y += offsetY

			// 3. Project to 2D
			viewerDist = 30 + (20 * (math.Sin(rl.GetTime())))
			v1 := p1.Project(float64(screenWidth), float64(screenHeight), fov, viewerDist)
			v2 := p2.Project(float64(screenWidth), float64(screenHeight), fov, viewerDist)

			// 4. Draw
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
