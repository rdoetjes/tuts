package main

import (
	"effects/effect"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Engine struct {
	fire       *effect.Flames
	pentagram  *effect.Pentagram
	rasterBars *effect.RasterBars
	embers     *effect.EmberField
	scroll1    *effect.Scroller
	scroll2    *effect.Scroller
	shader     rl.Shader
	target     rl.RenderTexture2D
	phraseIdx  int
}

func NewEngine() *Engine {
	e := &Engine{}

	e.fire = effect.NewFlames(400, 150)
	e.pentagram = effect.NewPentagram()
	e.rasterBars = effect.NewRasterBars(effect.FirePalette)
	e.embers = effect.NewStarfield(400)

	lyrics1 := "                  I GUESS, D'ELITE TEAM IS GOING TO HELL... AND LET us TELL YOU... IN THE IMMORTAL WORDS OF BON SCOTT..... HELL AIN'T A BAD PLACE TO BE!!!...          "
	e.scroll1 = effect.NewScroller(lyrics1, effect.FirePalette, "assets/fonts/Menlo.ttf", 400, float32(rl.GetScreenHeight())/2-250, 8, 100, 12, true)

	candyPalette := []rl.Color{
		{255, 0, 0, 255},   // Bright Red
		{255, 100, 0, 255}, // Vibrant Orange
		{255, 200, 0, 255}, // Hot Yellow
		{120, 40, 0, 255},  // Deep Burnt Sienna
	}

	lyrics2 := "HERE WE ARE, PHONAX AND CRU JONES, FROM D'ELITE.....              YEAH, YOU FUCKING LAME ASS BITCHES, YOU REALLY THOUGHT THIS ONE WAS SAFE.....                              HAHAHAHAHA!!!                      WRONG!!!                      WE'VE DONE IT AGAIN, BITCHES!!!                         WE TOOK YOUR PRECIOUS LITTLE PROTECTION, KICKED THE FUCKING DOOR IN, AND LEFT THE WHOLE THING SPRAWLED ACROSS THE BED LIKE A CHEAP, DISCARDED, HARBOR CRACKWHORE, AFTER A SAILOR ON HIS PAY DAY, HAD HIS WAY WITH HER!!!!......                         YOUR CODE?....WAS OUR, \"MATURE\", PLAYTHING.....                         YOUR SECURITY? FUCKING GONE, LIKE YOH MOMMA'S VIRGINITY WHEN SHE WAS 11.....                         YOUR PROTECTION? DEAD AND BURIED.                         AND YOUR BEAUTIFUL LITTLE PROGRAM?                                   CRACKED, PACKED, AND HANDED TO THE WORLD FOR FUCKING FREE!!!.....                         NO KEYS.                         NO PASSWORDS.                         NO REGISTRATION.                         NO MERCY.                         JUST ANOTHER LAME PROTECTION SCHEME ADDED TO OUR TROPHY ROOM AKA \"THE NAUGHT DUNGEON\", IF YOU KNOW WHAT WE MEAN.....                         SO TO THE CODERS WHO SPENT ALL NIGHT BUILDING THIS SHIT:                         NICE TRY, BITCHES!!!                         MAYBE NEXT TIME DON'T MAKE THE PROTECTION SO FUCKING EASY.....                         AND NOW FOR THE IMPORTANT PART.....                         GREETZ FLY OUT TO ALL THE DUDES STILL HAMMERING AWAY AT THEIR KEYBOARDS AT THREE IN THE FUCKING MORNING.....                         RESPECT TO THE REAL CRACKERS, THE CODERS, THE SWAPPERS, THE BBS FREAKS, AND EVERYONE KEEPING THIS SCENE ALIVE.....                         AND TO THE LAMERS, THE SUIT WEARING BETA MEN, THE SO CALLED COPY-PROTECTION \"GURUS\", AND EVERYONE WHO THINKS THEY CAN STOP US.....                                   FUCK YOU!!!.....                         YOU CAN SEND YOUR COMPLAINTS STRAIGHT TO HELL.....                         THAT'S WHERE WE'LL BE ANYWAY.....                         EATING PUSSY, DRINKING VODKA, AND HANGING AROUND SATAN'S TABLE LIKE WE OWN THE FUCKING PLACE.....                         ADOLF'S THERE.                         MAO'S THERE.                         SATAN'S THERE.                         AND SOMEHOW WE'RE STILL THE FUCKING PARTY ANIMALS AT THE TABLE!!!.....                         HAHAHAHAHA!!!                         I GUESS WE'RE GOING TO HELL.....                         ...BUT HEY.....                         HELL AIN'T A BAD PLACE TO BE!!!.....                         PHONAX.....                         CRU JONES.....                         D'ELITE.....                         WE CAME.                         WE SAW.                         WE CRACKED.                         AND THEN WE GAVE IT AWAY!!!.....                         FUCK YOUR COPY PROTECTION!!!                         FUCK YOUR REGISTRATION SCHEME!!!                         FUCK YOUR SECURITY!!!                         SEE YOU ON THE NEXT RELEASE, YOU LAMERS.....                                   ...IF YOU CAN KEEP UP!!!.....                         D'ELITE '90.....                         SIGNING OFF.....                         KEEP CRACKING.....                         KEEP SWAPPING.....                         AND KEEP THE FUCKING SCENE ALIVE!!!"
	e.scroll2 = effect.NewScroller(lyrics2, candyPalette, "assets/fonts/Impact.ttf", 60, float32(rl.GetScreenHeight())-80, 5, 15, 3, false)

	phrase := "KEEP THE FUCKING SCENE ALIVE!!!"
	e.phraseIdx = strings.Index(lyrics2, phrase)

	e.shader = rl.LoadShader("", "assets/shaders/crt.fs")
	resolutionLoc := rl.GetShaderLocation(e.shader, "resolution")
	rl.SetShaderValue(e.shader, resolutionLoc, []float32{float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight())}, rl.ShaderUniformVec2)

	e.target = rl.LoadRenderTexture(1024, 768)

	return e
}

func (e *Engine) Unload() {
	rl.UnloadRenderTexture(e.target)
	rl.UnloadShader(e.shader)
}

func (e *Engine) Process() {
	if rl.GetKeyPressed() > 0 {
		e.fire.Toggle()
	}

	e.fire.Process()
	e.pentagram.Process()
	e.rasterBars.Process()
	e.embers.Process()
	e.scroll1.Process()

	if e.scroll2.Process() {
		if !e.fire.IsBurning() {
			e.fire.Toggle()
		}
	}

	if e.phraseIdx != -1 && e.fire.IsBurning() && e.scroll2.GetCharX(e.phraseIdx) < 512 {
		e.fire.Toggle()
	}
}

func (e *Engine) Draw() {
	rl.BeginTextureMode(e.target)
	rl.ClearBackground(rl.Black)

	e.embers.Draw()

	if e.rasterBars.IsMovingUp() {
		e.pentagram.Draw()
		e.rasterBars.Draw()
	} else {
		e.rasterBars.Draw()
		e.pentagram.Draw()
	}

	e.scroll1.Draw()
	e.fire.Draw()
	e.scroll2.Draw()
	rl.EndTextureMode()

	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)

	rl.BeginShaderMode(e.shader)
	rl.DrawTextureRec(e.target.Texture, rl.Rectangle{X: 0, Y: 0, Width: float32(e.target.Texture.Width), Height: -float32(e.target.Texture.Height)}, rl.Vector2{X: 0, Y: 0}, rl.White)
	rl.EndShaderMode()

	rl.DrawFPS(10, 10)
	rl.EndDrawing()
}

func main() {
	rl.InitWindow(1024, 768, "EFFECTS")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)
	rl.DisableCursor()

	engine := NewEngine()
	defer engine.Unload()

	for !rl.WindowShouldClose() {
		engine.Process()
		engine.Draw()
	}
}
