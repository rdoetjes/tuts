package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Ball struct {
	Pos Vec3
	Vel Vec3
}

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
	balls        [2]Ball
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
		scrollerText: "------------------------------ PHONAX DELIVERS AGAIN! --- ANOTHER CRACK FOR US ELITE DUDES! --- WE HEARD THE LAMERS AT 'THE WEAKLINGS' ARE STILL TRYING TO FIGURE OUT THE NOP SLIDE... MAYBE TRY TOUCHING SOME GRASS INSTEAD....OR BETTER YET, START SMOKING SOME GRASS!!! --- PHONAX IS RAIDING THE SEVEN DIGITAL SEAS WHILE YOUR COMMODORE IS STILL LOADING FROM TAPE! --- GREETS TO THE REAL ONES ON THE WHIRLWIND BBS! --- BIG FUCKS TO THE PIRACY PATROL - CATCH US IF YOU CAN, SUCKERS! --- WE'RE NOT JUST CRACKING THE CODE, WE'RE CRACKING YOUR MOM'S FAVORITE HIGH SCORES! --- PHONAX: THE ONLY GROUP THAT CAN COMPILE IN THEIR SLEEP! --- KEEP YOUR EYES PEELED FOR OUR NEXT RELEASE OR YOU'LL BE STUCK PLAYING PONG FOREVER! --- PHONAX OWNS 1988!!! ------------------------",
		fontSize:     100.0,
	}

	s.font = rl.LoadFont("assets/fonts/Impact.ttf")
	defer rl.UnloadFont(s.font)

	s.sf = NewStarField(int(s.screenWidth * s.screenHeight * 0.0005))
	s.word = GetWord()
	s.scrollerX = float32(s.screenWidth)

	// Initialize balls (X is now static)
	s.balls[0] = Ball{Pos: Vec3{X: -12, Y: 0, Z: 10}}
	s.balls[1] = Ball{Pos: Vec3{X: 12, Y: 0, Z: 10}}

	rl.SetTargetFPS(60)
	rl.HideCursor()

	for !rl.WindowShouldClose() {
		s.update()
		s.draw()
	}
}

var MoodyPalette = []rl.Color{
	//{R: 20, G: 0, B: 40, A: 255},   // Midnight Purple (Deepest)
	{R: 40, G: 0, B: 80, A: 255},   // Dark Indigo
	{R: 120, G: 0, B: 180, A: 255}, // Moody Violet
	{R: 180, G: 0, B: 120, A: 255}, // Dim Magenta
	{R: 60, G: 20, B: 150, A: 255}, // Deep Electric Blue
	//{R: 20, G: 0, B: 40, A: 255},   // Back to Midnight
}

func (s *Scene) update() {
	s.angle += 0.015
	s.sf.Process()
	time := rl.GetTime()
	floorY := -8.0

	s.updateBalls(time, floorY)
	s.updateScroller()
}

func (s *Scene) updateBalls(time, floorY float64) {
	for i := range s.balls {
		b := &s.balls[i]

		// Depth movement (Z) - Moving near and far
		zCenter := 15.0
		zRange := 20.0
		zSpeed := 1.2
		offset := float64(i)*math.Pi + float64(100*i) // Out of phase
		b.Pos.Z = zCenter + math.Sin(time*zSpeed+offset)*zRange*2.0

		// Vertical bounce (Y)
		bounceHeight := 9.0
		bounceSpeed := 3.5
		// Use absolute sine for a "bouncing" motion off the floor
		b.Pos.Y = floorY + 1.0 + math.Abs(math.Sin(time*bounceSpeed+offset))*bounceHeight
	}
}

func (s *Scene) updateScroller() {
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

	// 4. Draw balls and shadows
	s.drawBalls()

	// 3. Draw word
	s.drawWord()

	// 5. Draw rainbow scroller
	s.drawScroller()

	rl.EndDrawing()
}

func (s *Scene) drawBalls() {
	floorY := -8.0
	for _, b := range s.balls {
		// 1. Draw shadow on the grid
		shadowPos := Vec3{X: b.Pos.X, Y: floorY, Z: b.Pos.Z}
		shadowProj := shadowPos.Project(s.screenWidth, s.screenHeight, s.fov, s.viewerDist)

		// Shadow size scales with height
		heightFactor := (b.Pos.Y - floorY)
		shadowSize := 40.0 / (1.0 + heightFactor*0.2)
		if shadowPos.Z > -s.viewerDist+1 {
			rl.DrawEllipse(int32(shadowProj.X), int32(shadowProj.Y), float32(shadowSize), float32(shadowSize/2), rl.NewColor(0, 0, 0, 150))
		}

		// 2. Draw "Chrome" Ball
		proj := b.Pos.Project(s.screenWidth, s.screenHeight, s.fov, s.viewerDist)
		if b.Pos.Z > -s.viewerDist+1 {
			radius := float32(40.0 / (s.viewerDist + b.Pos.Z) * (s.fov / 10.0))
			center := rl.NewVector2(float32(proj.X), float32(proj.Y))

			// 1. Chrome Gradient (Main Body)
			// Mid-grey to light-grey for that metallic "tapered" shading
			rl.DrawCircleGradient(center, radius, rl.LightGray, rl.NewColor(40, 40, 40, 255))
		}
	}
}

func (s *Scene) drawScroller() {
	currentX := s.scrollerX
	y := float32(s.screenHeight) - 120
	time := rl.GetTime()

	for i, char := range s.scrollerText {
		charStr := string(char)

		// Wavy effect
		waveY := y + float32(math.Sin(float64(i)*0.15+time*5.0)*25.0)

		// Calculate palette index and interpolation factor
		val := math.Mod(time*2.0+float64(i)*0.08, 1.0)
		if val < 0 {
			val += 1.0
		}

		idx := val * float64(len(MoodyPalette)-1)
		i1 := int(math.Floor(idx))
		i2 := i1 + 1
		frac := float32(idx - float64(i1))

		color := rl.ColorLerp(MoodyPalette[i1], MoodyPalette[i2], frac)

		position := rl.NewVector2(currentX, waveY)

		// 80s "glow" shadow
		rl.DrawTextEx(s.font, charStr, rl.NewVector2(position.X+3, position.Y+3), s.fontSize, 2, rl.NewColor(50, 0, 50, 200))
		rl.DrawTextEx(s.font, charStr, position, s.fontSize, 2, color)

		charWidth := rl.MeasureTextEx(s.font, charStr, s.fontSize, 2).X
		currentX += charWidth
	}
}

func (s *Scene) drawWord() {
	// Entire word rotations
	rotX := s.angle * 2
	rotY := s.angle * 3
	rotZ := 0.0
	time := rl.GetTime()

	for _, line := range s.word.Lines {
		// 1. Rotate
		p1 := transform(line.P1, rotX, rotY, rotZ)
		p2 := transform(line.P2, rotX, rotY, rotZ)

		// 2. Scale in and out
		scale := math.Abs((math.Sin(time) * 1.5))
		p1 = p1.Scale(scale)
		p2 = p2.Scale(scale)

		// 3. Project to 2D
		v1 := p1.Project(s.screenWidth, s.screenHeight, s.fov, s.viewerDist)
		v2 := p2.Project(s.screenWidth, s.screenHeight, s.fov, s.viewerDist)

		// 4. Draw 3 times with offset and colors from MoodyPalette
		for i := 0; i < 1; i++ {
			offset := float32(i) * 2.0

			// Select color from MoodyPalette based on layer index and time
			colorIdx := (i + int(time*5.0)) % (len(MoodyPalette) - 1)
			color := MoodyPalette[colorIdx]

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
				s := len(MoodyPalette)
				color = MoodyPalette[int(totalScroll)%s]
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
