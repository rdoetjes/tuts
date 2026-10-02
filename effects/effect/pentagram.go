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
	palette        []rl.Color
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

	numCirclePoints := 80
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

	palette := []rl.Color{
		{40, 10, 0, 255},     // Dark Patina / Oxidation
		{70, 30, 5, 255},     // Weathered Bronze
		{110, 50, 10, 255},   // Antique Copper
		{150, 80, 20, 255},   // Burnished Bronze
		{255, 230, 120, 255}, // Warm Metallic Glint
	}

	return &Pentagram{
		vertices:       vertices,
		circleVertices: circleVertices,
		angle:          0,
		palette:        palette,
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

func (p *Pentagram) drawBeveledLine(p1, p2 rl.Vector2, avgScale float32, depth float32, heatPulse float32) {
	baseThickness := 12.0 * avgScale
	if baseThickness < 4.0 {
		baseThickness = 4.0
	}

	dx := p2.X - p1.X
	dy := p2.Y - p1.Y
	angle := float32(math.Atan2(float64(dy), float64(dx)))

	depthShade := depth * 0.8
	if depthShade > 1.2 {
		depthShade = 1.2
	}

	numSteps := 2
	for step := 0; step < numSteps; step++ {
		t := 1.0 - (float32(step) / float32(numSteps))
		thickness := baseThickness * t

		colorIdx := step * len(p.palette) / numSteps
		if colorIdx >= len(p.palette) {
			colorIdx = len(p.palette) - 1
		}
		baseColor := p.palette[colorIdx]

		// Apply depth shading and pulsing ember effect
		// Fade out as it moves away from the camera
		// avgScale decreases as it moves away
		fade := avgScale * 4.0
		if fade > 1.0 {
			fade = 1.0
		}
		if fade < 0.0 {
			fade = 0.0
		}

		r := float32(baseColor.R) * depthShade * fade * heatPulse
		g := float32(baseColor.G) * depthShade * fade * heatPulse
		b := float32(baseColor.B) * depthShade * fade * heatPulse
		alpha := uint8(255.0 * fade)

		// Add directional lighting (specular) - Warm Metallic
		if step > 16 {
			spec := float32(math.Sin(float64(angle)*2.0+float64(p.angle))) * 0.25
			if spec > 0 {
				r += spec * 255 * fade
				g += spec * 200 * fade
				b += spec * 100 * fade
			}
		}

		rl.DrawLineEx(p1, p2, thickness, rl.Color{
			uint8(math.Min(255, float64(r))),
			uint8(math.Min(255, float64(g))),
			uint8(math.Min(255, float64(b))),
			alpha,
		})
	}

	// Warm metallic glint on the ridge
	glintFactor := float32(math.Cos(float64(angle) - float64(p.angle*0.5)))
	if glintFactor > 0.2 {
		fade := avgScale * 4.0
		if fade > 1.0 {
			fade = 1.0
		}
		alpha := uint8((glintFactor - 0.8) * 5.0 * 200 * fade)
		rl.DrawLineEx(p1, p2, 1.5*avgScale, rl.Color{255, 240, 180, alpha})
	}
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

	// Subtle "heat soak" pulse - looks like it's warm from the fire
	heatPulse := float32(math.Sin(rl.GetTime()*2.0)*0.1 + 0.9)

	// Draw Pentagram
	indices := []int{0, 2, 4, 1, 3, 0}
	for i := 0; i < 5; i++ {
		p1 := projected[indices[i]]
		p2 := projected[indices[i+1]]
		avgScale := (scales[indices[i]] + scales[indices[i+1]]) / 2.0
		// Depth factor for shading
		depth := avgScale * 5.0
		p.drawBeveledLine(p1, p2, avgScale, depth, heatPulse)
	}

	// Draw Circle Frame
	for i := 0; i < len(projectedCircle); i++ {
		p1 := projectedCircle[i]
		p2 := projectedCircle[(i+1)%len(projectedCircle)]
		avgScale := (scalesCircle[i] + scalesCircle[(i+1)%len(projectedCircle)]) / 2.0
		depth := avgScale * 5.0
		p.drawBeveledLine(p1, p2, avgScale, depth, heatPulse)
	}
}

var _ BaseEffect = (*Pentagram)(nil)
