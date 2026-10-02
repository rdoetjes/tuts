package effect

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Vec4 struct {
	X, Y, Z, W float32
}

type Pentagram struct {
	vertices       []Vec4
	circleVertices []Vec4
	angle          float32
}

func NewPentagram() *Pentagram {
	vertices := make([]Vec4, 5)
	for i := 0; i < 5; i++ {
		theta := float64(i) * 2.0 * math.Pi / 5.0
		vertices[i] = Vec4{
			X: float32(math.Cos(theta)),
			Y: float32(math.Sin(theta)),
			Z: 0,
			W: 0,
		}
	}

	numCirclePoints := 64
	circleVertices := make([]Vec4, numCirclePoints)
	radius := float32(1.15) // Slightly larger than the pentagram
	for i := 0; i < numCirclePoints; i++ {
		theta := float64(i) * 2.0 * math.Pi / float64(numCirclePoints)
		circleVertices[i] = Vec4{
			X: float32(math.Cos(theta)) * radius,
			Y: float32(math.Sin(theta)) * radius,
			Z: 0,
			W: 0,
		}
	}

	return &Pentagram{
		vertices:       vertices,
		circleVertices: circleVertices,
		angle:          0,
	}
}

func (p *Pentagram) Process() {
	p.angle += 0.02
}

func (p *Pentagram) rotate(v Vec4) Vec4 {
	// 4D rotations in multiple planes
	angle := float64(p.angle)

	// XY Plane (around Z)
	cosXY := float32(math.Cos(angle))
	sinXY := float32(math.Sin(angle))

	// XW Plane
	cosXW := float32(math.Cos(angle * 0.5))
	sinXW := float32(math.Sin(angle * 0.5))

	// YZ Plane
	cosYZ := float32(math.Cos(angle * 0.7))
	sinYZ := float32(math.Sin(angle * 0.7))

	// ZW Plane
	cosZW := float32(math.Cos(angle * 0.3))
	sinZW := float32(math.Sin(angle * 0.3))

	// Apply rotations
	x, y, z, w := v.X, v.Y, v.Z, v.W

	// XY
	x, y = x*cosXY-y*sinXY, x*sinXY+y*cosXY
	// XW
	x, w = x*cosXW-w*sinXW, x*sinXW+w*cosXW
	// YZ
	y, z = y*cosYZ-z*sinYZ, y*sinYZ+z*cosYZ
	// ZW
	z, w = z*cosZW-w*sinZW, z*sinZW+w*cosZW

	return Vec4{X: x, Y: y, Z: z, W: w}
}

func (p *Pentagram) project(v Vec4) (rl.Vector2, float32) {
	// Increased distances to reduce wide-angle skewing/distortion
	d4 := float32(5.0)
	d3 := float32(5.0)

	// Project 4D to 3D
	scale4 := 1.0 / (d4 - v.W)
	x3d := v.X * scale4
	y3d := v.Y * scale4
	z3d := v.Z * scale4

	// Project 3D to 2D
	scale3 := 1.0 / (d3 - z3d)

	screenWidth := float32(rl.GetScreenWidth())
	screenHeight := float32(rl.GetScreenHeight())

	// Use a fixed base size relative to screen height to maintain aspect ratio
	// Multiplier increased to account for larger projection distances
	size := screenHeight * 3.5

	return rl.Vector2{
		X: screenWidth/2 + x3d*scale3*size,
		Y: screenHeight/2 + y3d*scale3*size,
	}, scale4 * scale3 * 25.0 // Normalised scale for thickness
}

func (p *Pentagram) Draw() {
	// Project pentagram vertices
	projected := make([]rl.Vector2, len(p.vertices))
	scales := make([]float32, len(p.vertices))
	for i, v := range p.vertices {
		rotated := p.rotate(v)
		projected[i], scales[i] = p.project(rotated)
	}

	// Project circle vertices
	projectedCircle := make([]rl.Vector2, len(p.circleVertices))
	scalesCircle := make([]float32, len(p.circleVertices))
	for i, v := range p.circleVertices {
		rotated := p.rotate(v)
		projectedCircle[i], scalesCircle[i] = p.project(rotated)
	}

	palette := []rl.Color{
		{0, 50, 150, 255},    // Deep Blue
		{50, 0, 150, 255},    // Indigo
		{100, 0, 100, 255},   // Purple
		{200, 0, 50, 255},    // Dark Red
		{255, 50, 0, 255},    // Bright Red
		{255, 200, 100, 255}, // Hot highlight
	}

	// Helper to draw beveled lines
	drawBeveledLine := func(p1, p2 rl.Vector2, avgScale float32) {
		// Base thickness adjusted for the new normalised scale
		baseThickness := 10.0 * avgScale
		if baseThickness < 3.0 {
			baseThickness = 3.0
		}

		for step := 0; step < 10; step++ {
			// Thickness tapers from baseThickness down to a highlight
			t := 1.0 - (float32(step) / 10.0)
			thickness := baseThickness * t
			if thickness < 1.0 {
				thickness = 1.0
			}
			colorIdx := step * len(palette) / 10
			if colorIdx >= len(palette) {
				colorIdx = len(palette) - 1
			}
			rl.DrawLineEx(p1, p2, thickness, palette[colorIdx])
		}
	}

	// Draw Pentagram
	indices := []int{0, 2, 4, 1, 3, 0}
	for i := 0; i < 5; i++ {
		p1 := projected[indices[i]]
		p2 := projected[indices[i+1]]
		avgScale := (scales[indices[i]] + scales[indices[i+1]]) / 2.0
		drawBeveledLine(p1, p2, avgScale)
	}

	// Draw Circle Frame
	for i := 0; i < len(projectedCircle); i++ {
		p1 := projectedCircle[i]
		p2 := projectedCircle[(i+1)%len(projectedCircle)]
		avgScale := (scalesCircle[i] + scalesCircle[(i+1)%len(projectedCircle)]) / 2.0
		drawBeveledLine(p1, p2, avgScale)
	}

	// Draw dots at pentagram vertices (rivets)
	for i := 0; i < 5; i++ {
		radius := 12.0 * scales[i]
		rl.DrawCircleV(projected[i], radius, rl.Color{0, 0, 100, 255})
		rl.DrawCircleV(projected[i], radius*0.8, rl.Color{255, 0, 0, 255})
		rl.DrawCircleV(projected[i], radius*0.3, rl.Color{255, 255, 255, 200})
	}
}

var _ BaseEffect = (*Pentagram)(nil)
