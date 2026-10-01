package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func drawBall(angle float64, radius float64, color rl.Color, diameter uint32) {
	centerX := float64(rl.GetScreenWidth()) / 2
	centerY := float64(rl.GetScreenHeight()) / 2

	x := math.Cos(angle) * radius
	y := math.Sin(angle) * radius

	x += centerX
	y += centerY

	rl.DrawCircle(int32(x), int32(y), float32(diameter), color)
}

func main() {

	rl.InitWindow(1024, 768, "demo")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	var i float64 = 0
	for !rl.WindowShouldClose() {
		i += float64(rl.GetFrameTime()) * 1.5 // Slower rotation for aesthetic
		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)
		rl.DrawFPS(10, 10)
		nrBalls := 20 // More balls for a smoother spiral effect
		spreadAngle := (2 * math.Pi) / float64(nrBalls)

		for j := 0; j < nrBalls; j++ {
			// Calculate position angle
			angle := float64(j)*spreadAngle + i

			// t determines the color and pulse phase
			t := i*2.0 + float64(j)*0.4

			// 80s Synthwave Palette
			r := uint8(127 + 127*math.Cos(t))
			g := uint8(127 + 127*math.Sin(t*0.7))
			b := uint8(200 + 55*math.Sin(t*0.5))
			color := rl.NewColor(r, g, b, 255)

			drawBall(angle, 100+(50*math.Sin(i)), color, 15)
		}

		rl.EndDrawing()
	}
}
