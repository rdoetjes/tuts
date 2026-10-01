package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type OrbitingBalls struct{}

func (ob *OrbitingBalls) Draw(timer float64) {
	nrBalls := 12
	spreadAngle := (2 * math.Pi) / float64(nrBalls)
	for j := 0; j < nrBalls; j++ {
		angle := float64(j)*spreadAngle + timer*2.0

		t := timer*3.0 + float64(j)*0.5
		g := uint8(115 + 115*math.Sin(t))
		color := rl.NewColor(255, g, 0, 255)

		ob.drawBall(angle, 280, color, 30, timer)
	}
}

func (ob *OrbitingBalls) drawBall(angle float64, radius float64, color rl.Color, ballRadius float32, timer float64) {
	centerX := float64(screenWidth) / 2
	centerY := float64(screenHeight) / 2

	radius += math.Sin(timer*3) * 30

	x := int32(math.Cos(angle)*radius + centerX)
	y := int32(math.Sin(angle)*radius + centerY)

	rl.DrawCircle(x, y, ballRadius, color)
	rl.DrawCircleSector(rl.NewVector2(float32(x), float32(y)), ballRadius, 0, 180, 0, rl.NewColor(0, 0, 0, 140))
	rl.DrawCircleGradient(rl.NewVector2(float32(x-int32(ballRadius/3)), float32(y-int32(ballRadius/3))), ballRadius*0.7, rl.NewColor(255, 255, 255, 120), rl.Blank)
	rl.DrawCircle(x-int32(ballRadius/2), y-int32(ballRadius/2), ballRadius*0.15, rl.White)
}
