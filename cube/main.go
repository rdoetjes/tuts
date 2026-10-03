package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Vector3 represents a point in 3D space.
type Vector3 struct {
	X, Y, Z float64
}

// Vector2 represents a point in 2D screen space.
type Vector2 struct {
	X, Y float32
}

func main() {
	// Initialization
	const screenWidth = 800
	const screenHeight = 600

	rl.InitWindow(screenWidth, screenHeight, "Raylib Go - CPU Plasma & Wireframes")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	// --- 1. Plasma Setup (CPU-rendered) ---
	const plasmaWidth = 200
	const plasmaHeight = 150

	// Create a texture to display the CPU-rendered plasma
	// We'll update this texture every frame with our pixel data
	plasmaTexture := rl.LoadRenderTexture(int32(plasmaWidth), int32(plasmaHeight)).Texture
	defer rl.UnloadTexture(plasmaTexture)

	// Buffer to hold pixel data (RGBA, 4 bytes per pixel)
	plasmaPixels := make([]byte, plasmaWidth*plasmaHeight*4)

	// --- 2. Define Shapes ---
	// Cube
	cubeVertices := []Vector3{
		{-0.7, -0.7, 0.7}, {0.7, -0.7, 0.7}, {0.7, 0.7, 0.7}, {-0.7, 0.7, 0.7},
		{-0.7, -0.7, -0.7}, {0.7, -0.7, -0.7}, {0.7, 0.7, -0.7}, {-0.7, 0.7, -0.7},
	}
	cubeEdges := [][2]int{
		{0, 1}, {1, 2}, {2, 3}, {3, 0},
		{4, 5}, {5, 6}, {6, 7}, {7, 4},
		{0, 4}, {1, 5}, {2, 6}, {3, 7},
	}

	var angleX, angleY, angleZ float64
	var time float64

	for !rl.WindowShouldClose() {
		// Update
		time += 0.03
		angleX += 0.01
		angleY += 0.02
		angleZ += 0.015

		// --- 3. Update Plasma (CPU) ---
		updatePlasma(plasmaPixels, plasmaWidth, plasmaHeight, time)
		rl.UpdateTexture(plasmaTexture, plasmaPixels)

		// Draw
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		// Draw Plasma Background (scaled to screen)
		rl.DrawTexturePro(
			plasmaTexture,
			rl.NewRectangle(0, 0, float32(plasmaWidth), float32(plasmaHeight)),
			rl.NewRectangle(0, 0, float32(screenWidth), float32(screenHeight)),
			rl.NewVector2(0, 0),
			0,
			rl.White,
		)

		// Draw Cube (shifted left)
		drawShape(cubeVertices, cubeEdges, Vector3{0, 0, 0}, angleX, angleY, angleZ, rl.RayWhite)

		rl.DrawFPS(10, 10)
		rl.EndDrawing()
	}
}

// updatePlasma generates an old-school plasma effect on the CPU into a byte buffer.
func updatePlasma(pixels []byte, width, height int, t float64) {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			fx := float64(x)
			fy := float64(y)

			// Traditional plasma formula: sum of several sine waves
			v1 := math.Sin(fx/16.0 + t)
			v2 := math.Sin(fy/8.0 + t*1.5)
			v3 := math.Sin((fx + fy + t*10.0) / 16.0)
			v4 := math.Sin(math.Sqrt(fx*fx+fy*fy)/8.0 + t)

			val := v1 + v2 + v3 + v4

			// Map val [-4, 4] to color components
			r := uint8(128 + 127*math.Sin(val*math.Pi))
			g := uint8(128 + 127*math.Sin(val*math.Pi+2.0*math.Pi/3.0))
			b := uint8(128 + 127*math.Sin(val*math.Pi+4.0*math.Pi/3.0))

			idx := (y*width + x) * 4
			pixels[idx] = r
			pixels[idx+1] = g
			pixels[idx+2] = b
			pixels[idx+3] = 255
		}
	}
}

// drawShape handles rotation, projection, and rendering of a wireframe shape.
func drawShape(vertices []Vector3, edges [][2]int, offset Vector3, ax, ay, az float64, color rl.Color) {
	const screenWidth = 800
	const screenHeight = 600

	projectedPoints := make([]Vector2, len(vertices))

	for i, v := range vertices {
		// 1. Rotate around local origin
		rotated := rotate(v, ax, ay, az)

		// 2. Translate to world position
		worldPos := Vector3{
			X: rotated.X + offset.X,
			Y: rotated.Y + offset.Y,
			Z: rotated.Z + offset.Z,
		}

		// 3. Perspective Projection
		distance := 5.0
		fov := 400.0
		z := 1 / (distance - worldPos.Z)

		px := float32(worldPos.X*z*fov) + screenWidth/2
		py := float32(worldPos.Y*z*fov) + screenHeight/2

		projectedPoints[i] = Vector2{px, py}
	}

	for _, edge := range edges {
		p1 := projectedPoints[edge[0]]
		p2 := projectedPoints[edge[1]]
		rl.DrawLine(int32(p1.X), int32(p1.Y), int32(p2.X), int32(p2.Y), color)
	}
}

func rotate(v Vector3, ax, ay, az float64) Vector3 {
	res := v
	// Rotate around X
	y := res.Y*math.Cos(ax) - res.Z*math.Sin(ax)
	z := res.Y*math.Sin(ax) + res.Z*math.Cos(ax)
	res.Y, res.Z = y, z

	// Rotate around Y
	x := res.X*math.Cos(ay) + res.Z*math.Sin(ay)
	z = -res.X*math.Sin(ay) + res.Z*math.Cos(ay)
	res.X, res.Z = x, z

	// Rotate around Z
	x = res.X*math.Cos(az) - res.Y*math.Sin(az)
	y = res.X*math.Sin(az) + res.Y*math.Cos(az)
	res.X, res.Y = x, y

	return res
}
