package main

import (
	"math"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth   = 1024
	screenHeight  = 768
	starCount     = 400
	heartbeatSize = 60
	fireWidth     = 320 // Lower resolution for the fire effect
	fireHeight    = 100
)

var (
	heartbeatTable [heartbeatSize]float32
	fireBuffer     [fireWidth * fireHeight]uint8
)

func updateFire(isFlatline bool) {
	if !isFlatline {
		// Randomize bottom row (fire source) with "hot spots" for clumping
		for x := 0; x < fireWidth; x++ {
			if rl.GetRandomValue(0, 10) > 1 {
				fireBuffer[(fireHeight-1)*fireWidth+x] = uint8(rl.GetRandomValue(200, 255))
			} else {
				fireBuffer[(fireHeight-1)*fireWidth+x] = uint8(rl.GetRandomValue(0, 100))
			}
		}
	} else {
		// Cool down the source when flatlined
		for x := 0; x < fireWidth; x++ {
			if fireBuffer[(fireHeight-1)*fireWidth+x] > 8 {
				fireBuffer[(fireHeight-1)*fireWidth+x] -= 8
			} else {
				fireBuffer[(fireHeight-1)*fireWidth+x] = 0
			}
		}
	}

	// Propagate fire upwards with random horizontal drift
	for x := 0; x < fireWidth; x++ {
		for y := 1; y < fireHeight; y++ {
			// Get heat from below
			srcIdx := y*fireWidth + x
			pixel := fireBuffer[srcIdx]

			if pixel == 0 {
				fireBuffer[(y-1)*fireWidth+x] = 0
			} else {
				// Random horizontal drift (x - 1, x, or x + 1)
				randOffset := rl.GetRandomValue(0, 2) - 1
				dstX := (x + int(randOffset) + fireWidth) % fireWidth

				// Random cooling (increased if flatlined to fade out faster)
				coolingBase := int32(1)
				coolingMax := int32(4)
				if isFlatline {
					coolingBase = 3
					coolingMax = 6
				}
				cooling := uint8(rl.GetRandomValue(coolingBase, coolingMax))

				if uint32(pixel) > uint32(cooling) {
					fireBuffer[(y-1)*fireWidth+dstX] = pixel - cooling
				} else {
					fireBuffer[(y-1)*fireWidth+dstX] = 0
				}
			}
		}
	}
}

func drawFire() {
	scaleX := float32(screenWidth) / float32(fireWidth)
	scaleY := float32(250) / float32(fireHeight) // Slightly taller
	startY := float32(screenHeight - 220)

	for y := 0; y < fireHeight; y++ {
		for x := 0; x < fireWidth; x++ {
			val := fireBuffer[y*fireWidth+x]
			if val < 10 { // Threshold for visibility
				continue
			}

			// Map heat value to vivid fire palette
			var color rl.Color
			if val < 70 {
				// Deep red for the tips
				color = rl.NewColor(val*3, 0, 0, uint8(val*2))
			} else if val < 150 {
				// Fiery Orange
				g := uint8((float64(val) - 70) * 1.5)
				color = rl.NewColor(255, g, 0, 200)
			} else {
				// Bright Yellow core
				g := uint8(120 + (float64(val)-150)*1.3)
				color = rl.NewColor(255, g, 0, 255)
			}
			// Add some vertical tapering to the alpha based on height
			heightAlpha := uint8(float32(fireHeight-y) / float32(fireHeight) * 255)
			if color.A > heightAlpha {
				color.A = heightAlpha
			}

			rl.DrawRectangle(int32(float32(x)*scaleX), int32(startY+float32(y)*scaleY), int32(scaleX)+1, int32(scaleY)+1, color)
			rl.DrawRectangle(int32(float32(x)*scaleX), int32(rl.GetScreenHeight())-int32(startY+float32(y)*scaleY), int32(scaleX)+1, int32(scaleY)+1, color)

		}
	}
}

func initHeartbeat() {
	for i := 0; i < heartbeatSize; i++ {
		t := float64(i) / float64(heartbeatSize)
		var val float64

		// Realistic PQRST Waveform
		if t >= 0.1 && t <= 0.2 { // P wave: small bump
			val = 0.15 * math.Sin((t-0.1)*10*math.Pi)
		} else if t > 0.22 && t <= 0.24 { // Q wave: small dip
			val = -0.2 * math.Sin((t-0.22)*50*math.Pi)
		} else if t > 0.24 && t <= 0.28 { // R wave: massive spike
			val = 1.0 * math.Sin((t-0.24)*25*math.Pi)
		} else if t > 0.28 && t <= 0.30 { // S wave: sharp dip
			val = -0.4 * math.Sin((t-0.28)*50*math.Pi)
		} else if t >= 0.4 && t <= 0.6 { // T wave: medium bump
			val = 0.25 * math.Sin((t-0.4)*5*math.Pi)
		}

		heartbeatTable[i] = float32(val)
	}
}

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

		// Projection check for recycling
		// If z is too small or the star has flown past the screen boundaries, recycle it
		isOffScreen := false
		if stars[i].z > 0 {
			sx := (stars[i].x/stars[i].z)*100 + float32(screenWidth/2)
			sy := (stars[i].y/stars[i].z)*100 + float32(screenHeight/2)
			if sx < 0 || sx > screenWidth || sy < 0 || sy > screenHeight {
				isOffScreen = true
			}
		}

		if stars[i].z <= 0 || isOffScreen {
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
				color = rl.NewColor(255, 0, 0, 180) // Vivid Red
			} else {
				color = rl.NewColor(255, 230, 0, 180) // Bright Yellow
			}
			// Use intensity to modulate alpha for the glow effect
			color.A = uint8(float64(intensity) * 0.7)
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
	//rl.DrawCircleLines(x, y, ballRadius, rl.NewColor(255, 255, 255, 180))
}

func drawOrbitingBalls(timer float64) {
	nrBalls := 12
	spreadAngle := (2 * math.Pi) / float64(nrBalls)
	for j := 0; j < nrBalls; j++ {
		angle := float64(j)*spreadAngle + timer*2.0

		// Vivid Fire palette color for the balls (Red -> Orange -> Yellow)
		t := timer*3.0 + float64(j)*0.5
		g := uint8(115 + 115*math.Sin(t)) // G up to 230 for bright yellow
		color := rl.NewColor(255, g, 0, 255)

		drawBall(angle, 280, color, 30, timer)
	}
}

func drawScroller(font rl.Font, timer float64, scrollText string, scrollPos float32) {
	for i, char := range scrollText {
		charX := scrollPos + float32(i*45)
		if charX < -40 || charX > screenWidth {
			continue
		}
		// Classic sine wave scroller
		charY := float32(screenHeight-120) + float32(math.Sin(timer*4+float64(i)*0.25)*60)

		pos := rl.NewVector2(charX, charY)

		// Vivid Fire Palette: Red to Yellow
		t := timer*3.0 - float64(i)*0.2
		g := uint8(110 + 110*math.Cos(t)) // G up to 220
		// Main Animated Text Color
		textColor := rl.NewColor(255, g, 0, 255)

		// Pulsing Outline Color (Deep Red/Maroon)
		ot := timer*5.0 + float64(i)*0.3
		or := uint8(60 + 40*math.Sin(ot))
		og := uint8(0)
		ob := uint8(0)
		outlineColor := rl.NewColor(or, og, ob, 255)

		thickness := float32(5)
		offsets := []rl.Vector2{
			{-thickness, -thickness}, {0, -thickness}, {thickness, -thickness},
			{-thickness, 0}, {thickness, 0},
			{-thickness, thickness}, {0, thickness}, {thickness, thickness},
		}

		for _, off := range offsets {
			rl.DrawTextEx(font, string(char), rl.NewVector2(pos.X+off.X, pos.Y+off.Y), 60, 2, outlineColor)
		}

		// Main Animated Text
		rl.DrawTextEx(font, string(char), pos, 60, 2, textColor)
	}
}

func drawLogo(font rl.Font, timer float64, pulse float32) {
	headerText := "CRU JONES & PHONAX '89"
	fontSize := float32(80 + pulse*8) // Pulse the font size like a heartbeat

	textSize := rl.MeasureTextEx(font, headerText, fontSize, 2)
	headerX := (float32(screenWidth) - textSize.X) / 2

	// Logo shadow
	rl.DrawTextEx(font, headerText, rl.NewVector2(headerX+6, 56), fontSize, 2, rl.Maroon)
	// Logo main
	rl.DrawTextEx(font, headerText, rl.NewVector2(headerX, 50), fontSize, 2, rl.Gold)
}

func drawHeartRateMonitor(centerX, centerY int32, timer float64, isFlatline bool) {
	width := int32(180)
	height := int32(100)
	x := centerX - width/2
	y := centerY - height/2

	// Draw "Postage Stamp" Box with Maroon/Dark theme
	rl.DrawRectangle(x, y, width, height, rl.NewColor(30, 0, 0, 180))
	rl.DrawRectangleLinesEx(rl.NewRectangle(float32(x), float32(y), float32(width), float32(height)), 2, rl.Gold)

	var prevX, prevY float32
	hasPrev := false

	// Draw the trace inside
	for i := int32(0); i < width; i++ {
		var h float32 = 0
		if !isFlatline {
			// Calculate historical pulse based on horizontal position
			histOffset := (timer * 2.0) - float64(i)*0.01
			t := math.Mod(histOffset, 1.0)
			if t < 0 {
				t += 1.0
			}

			// Map t to our heartbeat table
			idx := int(t*float64(heartbeatSize)) % heartbeatSize
			h = heartbeatTable[idx] * float32(height/2-10)
		}

		// Vivid Fire palette color for the trace (Red to Yellow)
		colorT := float64(math.Sin(timer*4.0-float64(i)*0.04))*0.5 + 0.5
		g := uint8(colorT * 230) // G up to 230
		traceColor := rl.NewColor(255, g, 0, 255)

		px := float32(x + width - i)
		py := float32(y) + float32(height)/2 - h

		if hasPrev {
			rl.DrawLineEx(rl.NewVector2(prevX, prevY), rl.NewVector2(px, py), 3, traceColor)
		}

		// Add a little glow/tail (only if active)
		if !isFlatline && i < 20 {
			alpha := uint8(255 - i*12)
			glowColor := rl.NewColor(255, g, 0, alpha)
			rl.DrawCircle(int32(px), int32(py), 2, glowColor)
		}

		prevX, prevY = px, py
		hasPrev = true
	}
}

func drawSubHeader(font rl.Font, timer float64, isFlatline bool) {
	subText := "<< CRACKED BY D'ELITE >>"
	baseFontSize := float32(35)
	spacing := float32(2)

	// Calculate total width to center the whole string
	totalWidth := float32(0)
	for _, char := range subText {
		charSize := rl.MeasureTextEx(font, string(char), baseFontSize, spacing)
		totalWidth += charSize.X
	}

	startX := (float32(screenWidth) - totalWidth) / 2
	baseY := float32(135)

	for i, char := range subText {
		// Individual character wobbling
		charTimer := timer*5.0 + float64(i)*0.5
		offsetY := float32(math.Sin(charTimer)) * 15.0
		offsetX := float32(math.Cos(charTimer*0.8)) * 5.0

		// Individual scaling
		charPulse := float32(math.Sin(charTimer*0.7))*0.2 + 1.0
		charFontSize := baseFontSize * charPulse

		if isFlatline {
			charPulse = 0.0
			charFontSize = baseFontSize
			offsetX = 0.0
			offsetY = 0.0
		}

		// Vivid Fire: Transition from Red (255,0,0) to Yellow (255,220,0)
		t := float64(math.Sin(timer*3.0+float64(i)*0.3))*0.5 + 0.5
		g := uint8(t * 220)
		textColor := rl.NewColor(255, g, 0, 255)

		charStr := string(char)
		charSize := rl.MeasureTextEx(font, charStr, charFontSize, spacing)

		pos := rl.NewVector2(startX+offsetX, baseY+offsetY-charSize.Y/2)

		// Drop shadow
		rl.DrawTextEx(font, charStr, rl.NewVector2(pos.X+3, pos.Y+3), charFontSize, spacing, rl.NewColor(0, 0, 0, 150))

		// Character glow (using the same red/yellow shade with alpha)
		glowColor := rl.NewColor(255, g, 0, 80)
		rl.DrawTextEx(font, charStr, rl.NewVector2(pos.X-1, pos.Y-1), charFontSize, spacing, glowColor)
		rl.DrawTextEx(font, charStr, rl.NewVector2(pos.X+1, pos.Y+1), charFontSize, spacing, glowColor)

		// Main character
		rl.DrawTextEx(font, charStr, pos, charFontSize, spacing, textColor)

		// Advance startX for the next character (using base size for consistent spacing)
		baseCharSize := rl.MeasureTextEx(font, charStr, baseFontSize, spacing)
		startX += baseCharSize.X
	}
}

func drawGlitches() {
	if rl.GetRandomValue(0, 100) < 2 {
		glitchY := int32(rl.GetRandomValue(0, screenHeight))
		glitchH := int32(rl.GetRandomValue(5, 20))
		rl.DrawRectangle(0, glitchY, screenWidth, glitchH, rl.NewColor(255, 255, 255, 100))
	}
}

func drawScanlines() {
	for y := 0; y < screenHeight; y += 3 {
		rl.DrawLine(0, int32(y), screenWidth, int32(y), rl.NewColor(0, 0, 0, 80))
	}
}

func drawBorder() {
	rl.DrawRectangleLinesEx(rl.NewRectangle(0, 0, float32(screenWidth), float32(screenHeight)), 10, rl.NewColor(100, 100, 100, 255))
}

func main() {
	rl.InitWindow(screenWidth, screenHeight, "D'ELITE - 1989 CRACKTRO")
	defer rl.CloseWindow()

	font := rl.LoadFont("assets/fonts/Impact.ttf")
	defer rl.UnloadFont(font)

	rl.SetTargetFPS(60)

	stars := initStars()
	initHeartbeat()

	scrollText := "-----------------CRU JONES & PHONAX PRESENTS... THE 1989 ULTIMATE CRACKTRO DEMO!    CODED IN GO USING RAYLIB-GO...    GREETINGS TO: FAIRLIGHT - RAZOR 1911 - SKID ROW - GENESIS - TRSI - THE SILENTS - PHENOMENA - ANTHROX - TITAN...    WE BRING YOU THE BEST RELEASES, CRACKED AND PACKED FOR YOUR PLEASURE! ...... AND  REMEMBER: LIVE FAST, DIE YOUNG, LEAVE A GOOD LOOKING BODY!!! ..... AND ALWAYS, STAY RAD!!! --------------------------------------------"
	phrase := "BODY!!!"
	phraseIndex := strings.Index(scrollText, phrase)
	isFlatline := false

	scrollPos := float32(screenWidth)

	var timer float64 = 0
	for !rl.WindowShouldClose() {
		dt := float64(rl.GetFrameTime())
		timer += dt

		updateStars(stars)
		updateFire(isFlatline)
		scrollPos -= 5.0 // Scroll speed

		// Check if the specific phrase is centered on screen
		phraseX := scrollPos + float32(phraseIndex*45)
		if phraseX < float32(screenWidth/2+200) && phraseX > float32(screenWidth/2-400) {
			isFlatline = true
		}

		if scrollPos < -float32(len(scrollText)*45) {
			scrollPos = float32(screenWidth)
			isFlatline = false // Restart after wrap
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		// Get current pulse from pre-calculated table
		var pulse float32 = 0
		if !isFlatline {
			pulseIndex := int(timer*60) % heartbeatSize
			pulse = heartbeatTable[pulseIndex]
		}

		drawStars(stars)
		drawFire()
		drawCopperBars(timer)
		drawOrbitingBalls(timer)

		// Pass isFlatline to stop the EKG scrolling
		drawHeartRateMonitor(screenWidth/2, screenHeight/2, timer, isFlatline)

		drawScroller(font, timer, scrollText, scrollPos)
		drawLogo(font, timer, pulse)
		drawSubHeader(font, timer, isFlatline)
		//drawGlitches()
		drawScanlines()
		drawBorder()

		rl.EndDrawing()
	}
}
