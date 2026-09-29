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

	sf := NewStarField(int(float64(rl.GetScreenWidth()*rl.GetScreenHeight()) * 0.0005))

	rl.SetTargetFPS(60)
	rl.HideCursor()

	word := GetWord()
	screenWidth := rl.GetScreenWidth()
	screenHeight := rl.GetScreenHeight()
	angle := 0.0

	for !rl.WindowShouldClose() {
		update(sf, &angle)
		draw(sf, word, angle, float64(screenWidth), float64(screenHeight), fov, viewerDist)
	}
}

func update(sf *StarField, angle *float64) {
	*angle += 0.015
	sf.Process()
}

func draw(sf *StarField, word Letter, angle float64, screenWidth, screenHeight, fov, viewerDist float64) {
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)

	// 1. Draw stars
	sf.Draw()

	// 2. Draw scrolling floor
	drawFloor(angle, screenWidth, screenHeight, fov, 30.0)

	// 3. Draw word
	drawWord(word, angle, screenWidth, screenHeight, fov, 30.0)

	rl.EndDrawing()
}

func drawWord(word Letter, angle float64, screenWidth, screenHeight, fov, viewerDist float64) {
	// Entire word rotations
	rotX := angle * 0.4
	rotY := angle * 0.5
	rotZ := angle * 0.2

	for _, line := range word.Lines {
		// 1. Rotate
		p1 := transform(line.P1, rotX, rotY, rotZ)
		p2 := transform(line.P2, rotX, rotY, rotZ)

		// 2. Scale in and out
		scale := (math.Sin(rl.GetTime()) * 0.8) + 1.0
		p1 = p1.Scale(scale)
		p2 = p2.Scale(scale)

		// 3. Project to 2D
		v1 := p1.Project(screenWidth, screenHeight, fov, viewerDist)
		v2 := p2.Project(screenWidth, screenHeight, fov, viewerDist)

		// 4. Draw
		color := rl.RayWhite
		rl.DrawLineEx(rl.NewVector2(float32(v1.X), float32(v1.Y)), rl.NewVector2(float32(v2.X), float32(v2.Y)), 3.0, color)
	}
}

func transform(v Vec3, rx, ry, rz float64) Vec3 {
	v = v.RotateX(rx)
	v = v.RotateY(ry)
	v = v.RotateZ(rz)

	return v
}

func drawFloor(time, screenWidth, screenHeight, fov, viewerDist float64) {
	gridSize := 5.0
	numCols := 16
	numRows := 20
	floorY := -8.0

	// Scrolling speed and wrap-around
	speed := 15.0
	totalScroll := time * speed
	zOffset := math.Mod(totalScroll, gridSize)
	scrollIndex := int(math.Floor(totalScroll / gridSize))

	startX := -float64(numCols) * gridSize / 2.0
	startZ := -viewerDist + 5.0 // Start a bit in front of the viewer

	for i := 0; i < numRows; i++ {
		for j := 0; j < numCols; j++ {
			x0 := startX + float64(j)*gridSize
			x1 := x0 + gridSize
			z0 := startZ + float64(i)*gridSize - zOffset
			z1 := z0 + gridSize

			// Checkered pattern - pinned to virtual grid coordinates to prevent snapping
			color := rl.White
			if (i+j+scrollIndex)%2 == 0 {
				color = rl.Red
			}

			// Smooth distance fade (fog) based on actual Z position to prevent flickering
			currentZ := float64(i)*gridSize - zOffset
			maxDistance := float64(numRows) * gridSize
			distFade := 1.0 - (currentZ / maxDistance)
			if distFade < 0 {
				distFade = 0
			}
			r := uint8(float64(color.R) * distFade)
			g := uint8(float64(color.G) * distFade)
			b := uint8(float64(color.B) * distFade)
			fadeColor := rl.NewColor(r, g, b, 255)

			drawProjectedQuad(
				Vec3{X: x0, Y: floorY, Z: z0},
				Vec3{X: x1, Y: floorY, Z: z0},
				Vec3{X: x1, Y: floorY, Z: z1},
				Vec3{X: x0, Y: floorY, Z: z1},
				fadeColor,
				screenWidth, screenHeight, fov, viewerDist,
			)
		}
	}
}

func drawProjectedQuad(c1, c2, c3, c4 Vec3, color rl.Color, screenWidth, screenHeight, fov, viewerDist float64) {
	// Near plane clipping
	if c1.Z <= -viewerDist+1 || c2.Z <= -viewerDist+1 || c3.Z <= -viewerDist+1 || c4.Z <= -viewerDist+1 {
		return
	}

	p1 := c1.Project(screenWidth, screenHeight, fov, viewerDist)
	p2 := c2.Project(screenWidth, screenHeight, fov, viewerDist)
	p3 := c3.Project(screenWidth, screenHeight, fov, viewerDist)
	p4 := c4.Project(screenWidth, screenHeight, fov, viewerDist)

	v1 := rl.Vector2{X: float32(p1.X), Y: float32(p1.Y)}
	v2 := rl.Vector2{X: float32(p2.X), Y: float32(p2.Y)}
	v3 := rl.Vector2{X: float32(p3.X), Y: float32(p3.Y)}
	v4 := rl.Vector2{X: float32(p4.X), Y: float32(p4.Y)}

	rl.DrawTriangle(v1, v2, v3, color)
	rl.DrawTriangle(v1, v3, v4, color)
}
