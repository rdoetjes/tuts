package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Scene struct {
	starfield    *StarField
	word         Letter
	font         rl.Font
	fov          float64
	viewerDist   float64
	screenWidth  float64
	screenHeight float64
	angle        float64
	scroller     Scroller
	balls        BallSystem
}

var MoodyPalette = []rl.Color{
	{R: 20, G: 0, B: 40, A: 255},   // Midnight Purple (Deepest)
	{R: 40, G: 0, B: 80, A: 255},   // Dark Indigo
	{R: 120, G: 0, B: 180, A: 255}, // Moody Violet
	{R: 180, G: 0, B: 120, A: 255}, // Dim Magenta
	{R: 60, G: 20, B: 150, A: 255}, // Deep Electric Blue
}

func main() {
	rl.SetConfigFlags(rl.FlagFullscreenMode)
	rl.InitWindow(0, 0, "3D Vector Art - R A Y")
	defer rl.CloseWindow()

	s := &Scene{
		fov:          500.0,
		viewerDist:   30.0,
		screenWidth:  float64(rl.GetScreenWidth()),
		screenHeight: float64(rl.GetScreenHeight()),
	}

	s.font = rl.LoadFont("assets/fonts/Impact.ttf")
	defer rl.UnloadFont(s.font)

	s.starfield = NewStarField(int(s.screenWidth * s.screenHeight * 0.0005))
	s.word = GetWord()
	s.scroller = Scroller{
		Text:     "--- PHONAX DELIVERS AGAIN! --- ANOTHER 0-DAY CRACK FOR THE ELITE DUDES! --- WE HEARD THE LAMERS AT 'THE WEAKLINGS' ARE STILL TRYING TO FIGURE OUT THE NOP SLIDE... MAYBE TRY POKING SOME GRASS INSTEAD! --- PHONAX IS RAIDING THE SEVEN DIGITAL SEAS WHILE YOUR COMMODORE IS STILL LOADING FROM TAPE! --- GREETS TO THE REAL ONES ON THE WHIRLWIND BBS! --- BIG FUCKS TO THE PIRACY PATROL - CATCH US IF YOU CAN, SUCKERS! --- WE'RE NOT JUST CRACKING THE CODE, WE'RE CRACKING YOUR MOM'S FAVORITE HIGH SCORES! --- PHONAX: THE ONLY GROUP THAT CAN COMPILE IN THEIR SLEEP! --- KEEP YOUR EYES PEELED FOR OUR NEXT RELEASE OR YOU'LL BE STUCK PLAYING PONG FOREVER! --- PHONAX OWNS 1988!!! ---",
		FontSize: 100.0,
		X:        float32(s.screenWidth),
	}

	// Initialize balls
	s.balls.Data[0] = Ball{Pos: Vec3{X: -12, Y: 0, Z: 10}}
	s.balls.Data[1] = Ball{Pos: Vec3{X: 12, Y: 0, Z: 10}}

	rl.SetTargetFPS(60)
	rl.HideCursor()

	for !rl.WindowShouldClose() {
		s.Update()
		s.Draw()
	}
}

func (s *Scene) Update() {
	s.angle += 0.015
	s.starfield.Update()
	floorY := -8.0

	s.balls.Update(s, floorY)
	s.scroller.Update(s)
}

func (s *Scene) Draw() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)

	s.drawRasterBars()
	s.starfield.Draw()
	s.drawReflections()
	s.drawFloor()
	s.drawWord(false)
	s.balls.Draw(s)
	s.scroller.Draw(s)

	rl.EndDrawing()
}
