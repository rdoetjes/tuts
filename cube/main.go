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

	rl.InitWindow(screenWidth, screenHeight, "Raylib Go - Spinning Cube from Scratch")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	// 1. Define the Cube
	// A cube has 8 vertices (corners). We define them in a 3D coordinate system
	// where (0,0,0) is the center of the cube.
	vertices := []Vector3{
		{-1, -1, 1}, {1, -1, 1}, {1, 1, 1}, {-1, 1, 1}, // Front face (Z=1)
		{-1, -1, -1}, {1, -1, -1}, {1, 1, -1}, {-1, 1, -1}, // Back face (Z=-1)
	}

	// 2. Define Edges
	// We connect the vertices to form the wireframe. Each pair is an index into the vertices array.
	edges := [][2]int{
		{0, 1}, {1, 2}, {2, 3}, {3, 0}, // Front face loop
		{4, 5}, {5, 6}, {6, 7}, {7, 4}, // Back face loop
		{0, 4}, {1, 5}, {2, 6}, {3, 7}, // Lines connecting front and back
	}

	var angleX, angleY, angleZ float64

	for !rl.WindowShouldClose() {
		// Update
		angleX += 0.01
		angleY += 0.02
		angleZ += 0.015

		// Draw
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		projectedPoints := make([]Vector2, len(vertices))

		for i, v := range vertices {
			// 2. Rotation Logic
			// We apply rotation matrices for X, Y, and Z axes.
			rotated := rotate(v, angleX, angleY, angleZ)

			// 3. Projection Logic
			// To convert 3D coordinates (x, y, z) to 2D (x, y), we use perspective projection.
			// Formula: projected_coord = coord / (z + offset) * scale
			// 'offset' moves the object away from the camera so it doesn't clip.
			// 'scale' acts as a field of view/zoom factor.

			distance := 4.0 // Distance from "camera"
			fov := 400.0    // Scaling factor to map to screen pixels

			// We add the distance to Z to prevent division by zero and move the cube forward.
			z := 1 / (distance - rotated.Z)

			// Project X and Y
			px := float32(rotated.X*z*fov) + screenWidth/2
			py := float32(rotated.Y*z*fov) + screenHeight/2

			projectedPoints[i] = Vector2{px, py}
		}

		// 4. Rendering Edges
		// Draw lines between the projected 2D points.
		for _, edge := range edges {
			p1 := projectedPoints[edge[0]]
			p2 := projectedPoints[edge[1]]
			rl.DrawLine(int32(p1.X), int32(p1.Y), int32(p2.X), int32(p2.Y), rl.RayWhite)
		}

		rl.DrawText("Spinning Cube - Custom 3D Logic", 10, 10, 20, rl.Gray)
		rl.EndDrawing()
	}
}

// rotate applies 3D rotation matrices to a vector.
// Math Background:
// Rotation around X:
// x' = x
// y' = y*cos(a) - z*sin(a)
// z' = y*sin(a) + z*cos(a)
//
// Rotation around Y:
// x' = x*cos(a) + z*sin(a)
// y' = y
// z' = -x*sin(a) + z*cos(a)
//
// Rotation around Z:
// x' = x*cos(a) - y*sin(a)
// y' = x*sin(a) + y*cos(a)
// z' = z
func rotate(v Vector3, ax, ay, az float64) Vector3 {
	// Rotate around X
	res := v
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
