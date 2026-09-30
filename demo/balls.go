package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Ball struct {
	Pos Vec3
	Vel Vec3
}

type Spark struct {
	Pos      Vec3
	Vel      Vec3
	Life     float32
	ColorIdx int
}

type BallSystem struct {
	Data   [2]Ball
	Sparks []Spark
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

	// Update Sparks
	for i := len(bs.Sparks) - 1; i >= 0; i-- {
		bs.Sparks[i].Pos.X += bs.Sparks[i].Vel.X
		bs.Sparks[i].Pos.Y += bs.Sparks[i].Vel.Y
		bs.Sparks[i].Pos.Z += bs.Sparks[i].Vel.Z
		bs.Sparks[i].Vel.Y -= 0.01
		bs.Sparks[i].Life -= 0.02
		if bs.Sparks[i].Life <= 0 {
			bs.Sparks = append(bs.Sparks[:i], bs.Sparks[i+1:]...)
		}
	}
}

func (bs *BallSystem) Draw(s *Scene) {
	bs.DrawBalls(s, false)
}

func (bs *BallSystem) DrawBalls(s *Scene, isReflection bool) {
	floorY := -8.0
	for _, b := range bs.Data {
		ballPos := b.Pos
		if isReflection {
			ballPos.Y = floorY - (ballPos.Y - floorY)
			if ballPos.Y > floorY {
				continue
			}
		} else {
			shadowPos := Vec3{X: b.Pos.X, Y: floorY, Z: b.Pos.Z}
			shadowProj := shadowPos.Project(s.screenWidth, s.screenHeight, s.fov, s.viewerDist)
			heightFactor := (b.Pos.Y - floorY)
			shadowSize := 40.0 / (1.0 + heightFactor*0.2)
			if shadowPos.Z > -s.viewerDist+1 {
				rl.DrawEllipse(int32(shadowProj.X), int32(shadowProj.Y), float32(shadowSize), float32(shadowSize/2), rl.NewColor(0, 0, 0, 150))
			}
		}

		proj := ballPos.Project(s.screenWidth, s.screenHeight, s.fov, s.viewerDist)
		if ballPos.Z > -s.viewerDist+1 {
			radius := float32(40.0 / (s.viewerDist + ballPos.Z) * (s.fov / 10.0))
			center := rl.NewVector2(float32(proj.X), float32(proj.Y))
			c1 := rl.LightGray
			c2 := rl.NewColor(40, 40, 40, 255)
			if isReflection {
				c1.A = 100
				c2.A = 100
			}
			rl.DrawCircleGradient(center, radius, c1, c2)
		}
	}
}
