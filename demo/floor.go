package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Floor struct {
	Y              float64
	GridSize       float64
	NumCols        int
	NumRows        int
	RoadWidth      float64
	MountainBuffer float64
	Speed          float64
}

func (f *Floor) getHeight(x, z float64) float64 {
	absX := math.Abs(x)
	roadWidth := f.RoadWidth
	mountainBuffer := f.MountainBuffer

	if absX < roadWidth+mountainBuffer {
		return f.Y
	}

	dist := absX - (roadWidth + mountainBuffer)
	h := (math.Sin(x*0.12)*math.Cos(z*0.1) +
		math.Sin(z*0.15)*math.Sin(x*0.08) +
		math.Sin(x*0.05+z*0.05)) * 12.0

	return f.Y + math.Abs(h) + (dist * 0.4)
}

func (f *Floor) Draw(s *Scene) {
	totalScroll := s.angle * f.Speed
	zOffset := math.Mod(totalScroll, f.GridSize)
	scrollIndex := int(math.Floor(totalScroll / f.GridSize))

	startX := -float64(f.NumCols) * f.GridSize / 2.0
	startZ := -s.viewerDist + 5.0

	for i := 0; i < f.NumRows; i++ {
		for j := 0; j < f.NumCols; j++ {
			x0 := startX + float64(j)*f.GridSize
			x1 := x0 + f.GridSize
			z0 := startZ + float64(i)*f.GridSize - zOffset
			z1 := z0 + f.GridSize

			nearPlane := -s.viewerDist + 2.0
			if z0 < nearPlane || z1 < nearPlane {
				continue
			}

			currentZ := float64(i)*f.GridSize - zOffset
			maxDistance := float64(f.NumRows) * f.GridSize
			distFade := 1.0 - (currentZ / maxDistance)
			if distFade < 0 {
				distFade = 0
			}

			absX := math.Abs(x0)
			if absX < f.RoadWidth {
				// 1. Draw Road (Checkerboard)
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
					Vec3{X: x0, Y: f.Y, Z: z0},
					Vec3{X: x1, Y: f.Y, Z: z0},
					Vec3{X: x1, Y: f.Y, Z: z1},
					Vec3{X: x0, Y: f.Y, Z: z1},
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
