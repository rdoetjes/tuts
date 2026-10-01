package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func drawBall(angle float64, radius float64, color rl.Color, ballRadius float32) {
	centerX := float64(rl.GetScreenWidth()) / 2
	centerY := float64(rl.GetScreenHeight()) / 2

	x := int32(math.Cos(angle)*radius + centerX)
	y := int32(math.Sin(angle)*radius + centerY)

	// --- Chrome Shading Effect ---

	// 1. Draw the base sphere with the primary color
	rl.DrawCircle(x, y, ballRadius, color)

	// 2. Add a darker "bottom" shadow to give it 3D depth
	shadowColor := rl.NewColor(0, 0, 0, 100)
	rl.DrawCircleSector(rl.NewVector2(float32(x), float32(y)), ballRadius, 0, 180, 0, shadowColor)

	// 3. Draw a top-down gradient/highlight
	// We'll simulate this by drawing a lighter, smaller circle offset slightly upwards
	highlightColor := rl.NewColor(255, 255, 255, 80)
	rl.DrawCircle(x-int32(ballRadius/4), y-int32(ballRadius/4), ballRadius*0.7, highlightColor)

	// 4. Add the "Specular Glint" (the sharp white reflection)
	rl.DrawCircle(x-int32(ballRadius/2.5), y-int32(ballRadius/2.5), ballRadius*0.2, rl.White)

	// 5. Add a subtle rim light/outline to make it pop against the black
	rl.DrawCircleLines(x, y, ballRadius, rl.NewColor(200, 200, 200, 150))
}

func main() {

	rl.InitWindow(1024, 768, "demo")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	var i float64 = 0
	for !rl.WindowShouldClose() {
		i += float64(rl.GetFrameTime()) * 2.5 // Slower rotation for aesthetic
		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)
		rl.DrawFPS(10, 10)
		nrBalls := 24
		spreadAngle := (2 * math.Pi) / float64(nrBalls)

		for j := 0; j < nrBalls; j++ {
			// Calculate position angle
			angle := float64(j)*spreadAngle + i

			// t determines the color and pulse phase
			t := i*2.0 + float64(j)*0.3

			// Metallic Synthwave Palette
			// High contrast colors work best for chrome
			r := uint8(150 + 105*math.Cos(t))
			g := uint8(80 + 80*math.Sin(t*0.7))
			b := uint8(200 + 55*math.Sin(t*0.5))
			color := rl.NewColor(r, g, b, 255)

			// Spiral Effect
			baseRadius := 200.0
			amplitude := 10.0
			dynamicRadius := baseRadius + amplitude*math.Sin(i*4.0+float64(j)*0.4)

			drawBall(angle, dynamicRadius, color, 18)
		}

		rl.EndDrawing()
	}
}
