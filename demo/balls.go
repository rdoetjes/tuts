package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Ball struct {
	Pos Vec3
	Vel Vec3
}

type BallSystem struct {
	Data [2]Ball
}

func (bs *BallSystem) Update(s *Scene, floorY float64) {
	time := rl.GetTime()

	// Update Balls
	for i := range bs.Data {
		b := &bs.Data[i]

		zCenter := 15.0
		zRange := 20.0
		zSpeed := 1.2
		offset := float64(i)*math.Pi + float64(100*i)
		b.Pos.Z = zCenter + math.Sin(time*zSpeed+offset)*zRange*2.0

		bounceHeight := 9.0
		bounceSpeed := 3.5
		b.Pos.Y = floorY + 1.0 + math.Abs(math.Sin(time*bounceSpeed+offset))*bounceHeight
	}
}

func (bs *BallSystem) Draw(s *Scene) {
	bs.drawBalls(s)
}

func (bs *BallSystem) drawBalls(s *Scene) {
	floorY := s.floor.Y
	for _, b := range bs.Data {
		ballPos := b.Pos

		shadowPos := Vec3{X: b.Pos.X, Y: floorY, Z: b.Pos.Z}
		shadowProj := shadowPos.Project(s.screenWidth, s.screenHeight, s.fov, s.viewerDist)
		heightFactor := (b.Pos.Y - floorY)
		shadowSize := 40.0 / (1.0 + heightFactor*0.2)
		if shadowPos.Z > -s.viewerDist+1 {
			rl.DrawEllipse(int32(shadowProj.X), int32(shadowProj.Y), float32(shadowSize), float32(shadowSize/2), rl.NewColor(0, 0, 0, 150))
		}

		proj := ballPos.Project(s.screenWidth, s.screenHeight, s.fov, s.viewerDist)
		if ballPos.Z > -s.viewerDist+1 {
			radius := float32(40.0 / (s.viewerDist + ballPos.Z) * (s.fov / 10.0))
			center := rl.NewVector2(float32(proj.X), float32(proj.Y))

			c1 := rl.LightGray
			c2 := rl.NewColor(40, 40, 40, 255)

			rl.DrawCircleGradient(center, radius, c1, c2)
		}
	}
}
