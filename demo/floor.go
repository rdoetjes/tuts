package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (s *Scene) getHeight(x, z float64) float64 {
	absX := math.Abs(x)
	roadWidth := 30.0
	mountainBuffer := 15.0

	if absX < roadWidth+mountainBuffer {
		return -8.0
	}

	dist := absX - (roadWidth + mountainBuffer)
	h := (math.Sin(x*0.12)*math.Cos(z*0.1) +
		math.Sin(z*0.15)*math.Sin(x*0.08) +
		math.Sin(x*0.05+z*0.05)) * 12.0

	return -8.0 + math.Abs(h) + (dist * 0.4)
}

func (s *Scene) drawFloor() {
	gridSize := 5.0
	numCols := 40
	numRows := 25
	floorY := -8.0
	speed := 15.0
	totalScroll := s.angle * speed
	zOffset := math.Mod(totalScroll, gridSize)
	scrollIndex := int(math.Floor(totalScroll / gridSize))

	startX := -float64(numCols) * gridSize / 2.0
	startZ := -s.viewerDist + 5.0

	for i := 0; i < numRows; i++ {
		for j := 0; j < numCols; j++ {
			x0 := startX + float64(j)*gridSize
			x1 := x0 + gridSize
			z0 := startZ + float64(i)*gridSize - zOffset
			z1 := z0 + gridSize

			nearPlane := -s.viewerDist + 2.0
			if z0 < nearPlane || z1 < nearPlane {
				continue
			}

			currentZ := float64(i)*gridSize - zOffset
			maxDistance := float64(numRows) * gridSize
			distFade := 1.0 - (currentZ / maxDistance)
			if distFade < 0 {
				distFade = 0
			}

			absX := math.Abs(x0)
			if absX < 50 {
				color := rl.White
				if (i+j+scrollIndex)%2 == 0 {
					palSize := len(MoodyPalette)
					color = MoodyPalette[int(totalScroll)%palSize]
				}

				fadeColor := rl.NewColor(
					uint8(float64(color.R)*distFade),
					uint8(float64(color.G)*distFade),
					uint8(float64(color.B)*distFade),
					180,
				)

				drawProjectedQuad(
					Vec3{X: x0, Y: floorY, Z: z0},
					Vec3{X: x1, Y: floorY, Z: z0},
					Vec3{X: x1, Y: floorY, Z: z1},
					Vec3{X: x0, Y: floorY, Z: z1},
					fadeColor,
					s.screenWidth, s.screenHeight, s.fov, s.viewerDist,
				)
			}
		}
	}
}

func drawProjectedQuad(c1, c2, c3, c4 Vec3, color rl.Color, screenWidth, screenHeight, fov, viewerDist float64) {
	if c1.Z <= -viewerDist+1 || c2.Z <= -viewerDist+1 || c3.Z <= -viewerDist+1 || c4.Z <= -viewerDist+1 {
		return
	}

	p1 := c1.Project(screenWidth, screenHeight, fov, viewerDist)
	p2 := c2.Project(screenWidth, screenHeight, fov, viewerDist)
	p3 := c3.Project(screenWidth, screenHeight, fov, viewerDist)
	p4 := c4.Project(screenWidth, screenHeight, fov, viewerDist)

	rl.DrawTriangle(rl.NewVector2(float32(p1.X), float32(p1.Y)), rl.NewVector2(float32(p2.X), float32(p2.Y)), rl.NewVector2(float32(p3.X), float32(p3.Y)), color)
	rl.DrawTriangle(rl.NewVector2(float32(p1.X), float32(p1.Y)), rl.NewVector2(float32(p3.X), float32(p3.Y)), rl.NewVector2(float32(p4.X), float32(p4.Y)), color)
}
