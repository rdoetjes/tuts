package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 1024
	screenHeight = 768
	starCount    = 200
)

type Star struct {
	x, y, z float32
}

func initStars() []Star {
	stars := make([]Star, starCount)
	for i := range stars {
		stars[i] = Star{
			x: float32(rl.GetRandomValue(-screenWidth, screenWidth)),
			y: float32(rl.GetRandomValue(-screenHeight, screenHeight)),
			z: float32(rl.GetRandomValue(1, screenWidth)),
		}
	}
	return stars
}

func updateStars(stars []Star) {
	speed := float32(8.0)
	for i := range stars {
		stars[i].z -= speed
		if stars[i].z <= 0 {
			stars[i].x = float32(rl.GetRandomValue(-screenWidth, screenWidth))
			stars[i].y = float32(rl.GetRandomValue(-screenHeight, screenHeight))
			stars[i].z = float32(screenWidth)
		}
	}
}

func drawStars(stars []Star) {
	for _, s := range stars {
		// Project 3D to 2D
		sx := (s.x/s.z)*100 + float32(screenWidth/2)
		sy := (s.y/s.z)*100 + float32(screenHeight/2)

		if sx < 0 || sx > screenWidth || sy < 0 || sy > screenHeight {
			continue
		}

		size := (1.0 - s.z/screenWidth) * 4
		alpha := uint8((1.0 - s.z/screenWidth) * 255)
		rl.DrawCircle(int32(sx), int32(sy), size, rl.NewColor(255, 255, 255, alpha))
	}
}

func drawCopperBars(timer float64) {
	numBars := 6
	barHeight := float32(30)
	for i := 0; i < numBars; i++ {
		yPos := float32(screenHeight/2-100) + float32(math.Sin(timer*1.5+float64(i)*0.4)*250)

		// Draw a gradient bar
		for j := 0; j < int(barHeight); j++ {
			intensity := uint8(255 - math.Abs(float64(j)-float64(barHeight/2))*15)

			var color rl.Color
			if i%2 == 0 {
				color = rl.NewColor(intensity, 0, intensity, 180) // Purple/Magenta
			} else {
				color = rl.NewColor(0, intensity, intensity, 180) // Cyan
			}
			rl.DrawLine(0, int32(yPos)+int32(j), screenWidth, int32(yPos)+int32(j), color)
		}
	}
}

func drawBall(angle float64, radius float64, color rl.Color, ballRadius float32, timer float64) {
	centerX := float64(screenWidth) / 2
	centerY := float64(screenHeight) / 2

	// Wobble the radius for more dynamic movement
	radius += math.Sin(timer*3) * 30

	x := int32(math.Cos(angle)*radius + centerX)
	y := int32(math.Sin(angle)*radius + centerY)

	// --- Enhanced Chrome Shading Effect ---
	rl.DrawCircle(x, y, ballRadius, color)

	// Dark bottom half for depth
	rl.DrawCircleSector(rl.NewVector2(float32(x), float32(y)), ballRadius, 0, 180, 0, rl.NewColor(0, 0, 0, 140))

	// Metallic highlight (top-left)
	rl.DrawCircleGradient(rl.NewVector2(float32(x-int32(ballRadius/3)), float32(y-int32(ballRadius/3))), ballRadius*0.7, rl.NewColor(255, 255, 255, 120), rl.Blank)

	// Specular glint (sharp reflection)
	rl.DrawCircle(x-int32(ballRadius/2), y-int32(ballRadius/2), ballRadius*0.15, rl.White)

	// Glowing rim
	rl.DrawCircleLines(x, y, ballRadius, rl.NewColor(255, 255, 255, 180))
}

func main() {
	rl.InitWindow(screenWidth, screenHeight, "CRU JONES - 1989 CRACKTRO")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	stars := initStars()

	scrollText := "    CRU JONES PRESENTS... THE 1989 ULTIMATE CRACKTRO DEMO!    CODED IN GO USING RAYLIB-GO...    GREETINGS TO: FAIRLIGHT - RAZOR 1911 - SKID ROW - GENESIS - TRSI - THE SILENTS - PHENOMENA - ANTHROX - TITAN...    WE BRING YOU THE BEST RELEASES, CRACKED AND PACKED FOR YOUR PLEASURE!    REMEMBER: ONLY THE GOOD DIE YOUNG...    STAY RAD!    "
	scrollPos := float32(screenWidth)

	var timer float64 = 0
	for !rl.WindowShouldClose() {
		dt := float64(rl.GetFrameTime())
		timer += dt

		updateStars(stars)
		scrollPos -= 5.0 // Scroll speed
		if scrollPos < -float32(len(scrollText)*35) {
			scrollPos = float32(screenWidth)
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		// Background layer
		drawStars(stars)
		drawCopperBars(timer)

		// Middle layer: Orbiting balls
		nrBalls := 12
		spreadAngle := (2 * math.Pi) / float64(nrBalls)
		for j := 0; j < nrBalls; j++ {
			angle := float64(j)*spreadAngle + timer*2.0

			t := timer*3.0 + float64(j)*0.5
			r := uint8(150 + 105*math.Cos(t))
			g := uint8(50 + 50*math.Sin(t*0.8))
			b := uint8(220 + 35*math.Sin(t*0.6))
			color := rl.NewColor(r, g, b, 255)

			drawBall(angle, 280, color, 30, timer)
		}

		// Foreground layer: Scrolling text
		for i, char := range scrollText {
			charX := scrollPos + float32(i*35)
			if charX < -40 || charX > screenWidth {
				continue
			}
			// Classic sine wave scroller
			charY := float32(screenHeight-120) + float32(math.Sin(timer*4+float64(i)*0.25)*60)

			// Text shadow
			rl.DrawText(string(char), int32(charX)+4, int32(charY)+4, 60, rl.DarkPurple)
			// Animated text color
			color := rl.NewColor(uint8(127+127*math.Sin(timer*2+float64(i)*0.1)), 255, 255, 255)
			rl.DrawText(string(char), int32(charX), int32(charY), 60, color)
		}

		// Logo/Header with a pulsing effect
		headerText := "CRU JONES '89"
		logoScale := 80 + int32(math.Sin(timer*2)*5)
		headerX := int32(screenWidth/2 - rl.MeasureText(headerText, logoScale)/2)

		// Logo shadow
		rl.DrawText(headerText, headerX+6, 56, logoScale, rl.Maroon)
		// Logo main
		rl.DrawText(headerText, headerX, 50, logoScale, rl.Gold)

		// Sub-header
		subText := "<< CRACKED BY CRU >>"
		subX := int32(screenWidth/2 - rl.MeasureText(subText, 30)/2)
		rl.DrawText(subText, subX, 130, 30, rl.Lime)

		// Occasional Glitch Effect
		if rl.GetRandomValue(0, 100) < 2 {
			glitchY := int32(rl.GetRandomValue(0, screenHeight))
			glitchH := int32(rl.GetRandomValue(5, 20))
			rl.DrawRectangle(0, glitchY, screenWidth, glitchH, rl.NewColor(255, 255, 255, 100))
		}

		// Scanlines / CRT Effect
		for y := 0; y < screenHeight; y += 3 {
			rl.DrawLine(0, int32(y), screenWidth, int32(y), rl.NewColor(0, 0, 0, 80))
		}

		// Border frame
		rl.DrawRectangleLinesEx(rl.NewRectangle(0, 0, float32(screenWidth), float32(screenHeight)), 10, rl.NewColor(100, 100, 100, 255))

		// rl.DrawFPS(10, 10)
		rl.EndDrawing()
	}
}
