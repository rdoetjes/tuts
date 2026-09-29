package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Scene struct {
	sf           *StarField
	word         Letter
	font         rl.Font
	fov          float64
	viewerDist   float64
	screenWidth  float64
	screenHeight float64
	angle        float64
	scrollerX    float32
	scrollerText string
	fontSize     float32
}

func main() {
	// Set fullscreen flag before initialization
	rl.SetConfigFlags(rl.FlagFullscreenMode)

	rl.InitWindow(0, 0, "3D Vector Art - R A Y")
	defer rl.CloseWindow()

	s := &Scene{
		fov:          500.0,
		viewerDist:   30.0,
		screenWidth:  float64(rl.GetScreenWidth()),
		screenHeight: float64(rl.GetScreenHeight()),
		scrollerText: "------------------------------ PHONAX DELIVERS AGAIN! --- ANOTHER 0-DAY CRACK FOR THE ELITE DUDES! --- WE HEARD THE LAMERS AT 'THE WEAKLINGS' ARE STILL TRYING TO FIGURE OUT THE NOP SLIDE... MAYBE TRY POKING SOME GRASS INSTEAD! --- PHONAX IS RAIDING THE SEVEN DIGITAL SEAS WHILE YOUR COMMODORE IS STILL LOADING FROM TAPE! --- GREETS TO THE REAL ONES ON THE WHIRLWIND BBS! --- BIG FUCKS TO THE PIRACY PATROL - CATCH US IF YOU CAN, SUCKERS! --- WE'RE NOT JUST CRACKING THE CODE, WE'RE CRACKING YOUR MOM'S FAVORITE HIGH SCORES! --- PHONAX: THE ONLY GROUP THAT CAN COMPILE IN THEIR SLEEP! --- KEEP YOUR EYES PEELED FOR OUR NEXT RELEASE OR YOU'LL BE STUCK PLAYING PONG FOREVER! --- PHONAX OWNS 1988!!! ----------------",
		fontSize:     100.0,
	}

	s.font = rl.LoadFont("assets/fonts/Impact.ttf")
	defer rl.UnloadFont(s.font)

	s.sf = NewStarField(int(s.screenWidth * s.screenHeight * 0.0005))
	s.word = GetWord()
	s.scrollerX = float32(s.screenWidth)

	rl.SetTargetFPS(60)
	rl.HideCursor()

	for !rl.WindowShouldClose() {
		s.update()
		s.draw()
	}
}

func (s *Scene) update() {
	s.angle += 0.015
	s.sf.Process()

	// Update scroller
	s.scrollerX -= 4.5
	textSize := rl.MeasureTextEx(s.font, s.scrollerText, s.fontSize, 2)
	if s.scrollerX < -textSize.X {
		s.scrollerX = float32(s.screenWidth)
	}
}

func (s *Scene) draw() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)

	// 1. Draw stars
	s.sf.Draw()

	// 2. Draw scrolling floor
	drawFloor(s.angle, s.screenWidth, s.screenHeight, s.fov, s.viewerDist)

	// 3. Draw word
	s.drawWord()

	// 4. Draw rainbow scroller
	s.drawScroller()

	rl.EndDrawing()
}

func (s *Scene) drawScroller() {
	currentX := s.scrollerX
	y := float32(s.screenHeight) - 150
	time := rl.GetTime()

	for i, char := range s.scrollerText {
		charStr := string(char)
		hue := uint8(int(time*700+float64(i)*15) % 200)
		color := rl.NewColor(hue, hue, hue, 255)

		// Wavy effect
		waveY := y + float32(math.Sin(float64(i)*0.15+time*5.0)*25.0)

		position := rl.NewVector2(currentX, waveY)
		rl.DrawTextEx(s.font, charStr, position, s.fontSize, 2, color)

		charWidth := rl.MeasureTextEx(s.font, charStr, s.fontSize, 2).X
		currentX += charWidth
	}
}

func (s *Scene) drawWord() {
	// Entire word rotations
	rotX := 0.0
	rotY := s.angle * 0.45
	rotZ := 0.0
	time := rl.GetTime()

	for _, line := range s.word.Lines {
		// 1. Rotate
		p1 := transform(line.P1, rotX, rotY, rotZ)
		p2 := transform(line.P2, rotX, rotY, rotZ)

		// 2. Scale in and out
		scale := (math.Sin(time) * 0.8) + 1.0
		p1 = p1.Scale(scale)
		p2 = p2.Scale(scale)

		// 3. Project to 2D
		v1 := p1.Project(s.screenWidth, s.screenHeight, s.fov, s.viewerDist)
		v2 := p2.Project(s.screenWidth, s.screenHeight, s.fov, s.viewerDist)

		// 4. Draw 3 times with offset and hue shift for extruded look
		for i := 0; i < 2; i++ {
			offset := float32(i) * 2.0
			hue := float32(int(time*120+float64(i)*40) % 360)
			color := rl.ColorFromHSV(hue, 0.8, 1.0)

			p1_2d := rl.NewVector2(float32(v1.X)+offset, float32(v1.Y)+offset)
			p2_2d := rl.NewVector2(float32(v2.X)+offset, float32(v2.Y)+offset)

			rl.DrawLineEx(p1_2d, p2_2d, 3.0, color)
		}
	}
}

func transform(v Vec3, rx, ry, rz float64) Vec3 {
	v = v.RotateX(rx)
	v = v.RotateY(ry)
	v = v.RotateZ(rz)

	return v
}

func drawFloor(time, screenWidth, screenHeight, fov, viewerDist float64) {
	gridSize := 5.0
	numCols := 16
	numRows := 20
	floorY := -8.0

	// Scrolling speed and wrap-around
	speed := 15.0
	totalScroll := time * speed
	zOffset := math.Mod(totalScroll, gridSize)
	scrollIndex := int(math.Floor(totalScroll / gridSize))

	startX := -float64(numCols) * gridSize / 2.0
	startZ := -viewerDist + 5.0 // Start a bit in front of the viewer

	for i := 0; i < numRows; i++ {
		for j := 0; j < numCols; j++ {
			x0 := startX + float64(j)*gridSize
			x1 := x0 + gridSize
			z0 := startZ + float64(i)*gridSize - zOffset
			z1 := z0 + gridSize

			// Checkered pattern - pinned to virtual grid coordinates to prevent snapping
			color := rl.White
			if (i+j+scrollIndex)%2 == 0 {
				color = rl.Red
			}

			// Smooth distance fade (fog) based on actual Z position to prevent flickering
			currentZ := float64(i)*gridSize - zOffset
			maxDistance := float64(numRows) * gridSize
			distFade := 1.0 - (currentZ / maxDistance)
			if distFade < 0 {
				distFade = 0
			}
			r := uint8(float64(color.R) * distFade)
			g := uint8(float64(color.G) * distFade)
			b := uint8(float64(color.B) * distFade)
			fadeColor := rl.NewColor(r, g, b, 255)

			drawProjectedQuad(
				Vec3{X: x0, Y: floorY, Z: z0},
				Vec3{X: x1, Y: floorY, Z: z0},
				Vec3{X: x1, Y: floorY, Z: z1},
				Vec3{X: x0, Y: floorY, Z: z1},
				fadeColor,
				screenWidth, screenHeight, fov, viewerDist,
			)
		}
	}
}

func drawProjectedQuad(c1, c2, c3, c4 Vec3, color rl.Color, screenWidth, screenHeight, fov, viewerDist float64) {
	// Near plane clipping
	if c1.Z <= -viewerDist+1 || c2.Z <= -viewerDist+1 || c3.Z <= -viewerDist+1 || c4.Z <= -viewerDist+1 {
		return
	}

	p1 := c1.Project(screenWidth, screenHeight, fov, viewerDist)
	p2 := c2.Project(screenWidth, screenHeight, fov, viewerDist)
	p3 := c3.Project(screenWidth, screenHeight, fov, viewerDist)
	p4 := c4.Project(screenWidth, screenHeight, fov, viewerDist)

	v1 := rl.Vector2{X: float32(p1.X), Y: float32(p1.Y)}
	v2 := rl.Vector2{X: float32(p2.X), Y: float32(p2.Y)}
	v3 := rl.Vector2{X: float32(p3.X), Y: float32(p3.Y)}
	v4 := rl.Vector2{X: float32(p4.X), Y: float32(p4.Y)}

	rl.DrawTriangle(v1, v2, v3, color)
	rl.DrawTriangle(v1, v3, v4, color)
}
