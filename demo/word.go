package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type WordSystem struct {
	Data Letter
}

func (ws *WordSystem) Draw(s *Scene, isReflection bool) {
	rotX := s.angle * 2
	rotY := s.angle * 3
	rotZ := 0.0
	time := rl.GetTime()
	floorY := s.floor.Y

	for _, line := range ws.Data.Lines {
		p1 := transform(line.P1, rotX, rotY, rotZ)
		p2 := transform(line.P2, rotX, rotY, rotZ)

		scale := math.Abs((math.Sin(time) * 1.5))
		p1 = p1.Scale(scale)
		p2 = p2.Scale(scale)

		if isReflection {
			p1.Y = floorY - (p1.Y - floorY)
			p2.Y = floorY - (p2.Y - floorY)
			if p1.Y > floorY || p2.Y > floorY {
				continue
			}
		}

		v1 := p1.Project(s.screenWidth, s.screenHeight, s.fov, s.viewerDist)
		v2 := p2.Project(s.screenWidth, s.screenHeight, s.fov, s.viewerDist)

		for i := 0; i < 1; i++ {
			offset := float32(i) * 2.0
			colorIdx := (i + int(time*5.0)) % (len(MoodyPalette) - 1)
			color := MoodyPalette[colorIdx]

			if isReflection {
				color.A = 100
			}

			p1_2d := rl.NewVector2(float32(v1.X)+offset, float32(v1.Y)+offset)
			p2_2d := rl.NewVector2(float32(v2.X)+offset, float32(v2.Y)+offset)

			rl.DrawLineEx(p1_2d, p2_2d, 3.0, color)
		}
	}
}
