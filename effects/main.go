package main

import (
	"effects/effect"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(1024, 768, "EFFECTS")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)
	rl.DisableCursor()

	fire := effect.NewFlames(400, 150)
	pentagram := effect.NewPentagram()
	rasterBars := effect.NewRasterBars(effect.FirePalette)
	stars := effect.NewStarfield(400)
	moon := effect.NewMoon(200)

	lyrics := "                  I GUESS, D'ELITE TEAM IS GOING TO HELL... AND LET US TELL YOU... IN THE IMMORTAL WORDS OF BON SCOTT..... HELL AIN'T A BAD PLACE TO BE!!!...          "

	scroll := effect.NewScroller(lyrics, effect.FirePalette, "assets/fonts/Menlo.ttf", 400, float32(rl.GetScreenHeight())/2-250, 8, 100, 12, true)

	candyPalette := []rl.Color{
		{255, 0, 0, 255},   // Bright Red
		{255, 100, 0, 255}, // Vibrant Orange
		{255, 200, 0, 255}, // Hot Yellow
		{120, 40, 0, 255},  // Deep Burnt Sienna (Brown replacement)
	}

	lyrics2 := "HERE WE ARE, PHONAX AND CRU JONES, FROM D'ELITE.....              YEAH, YOU FUCKING LAME ASS BITCHES, YOU REALLY THOUGHT THIS ONE WAS SAFE.....                              HAHAHAHAHA!!!                      WRONG!!!                      WE'VE DONE IT AGAIN, BITCHES!!!                         WE TOOK YOUR PRECIOUS LITTLE PROTECTION, KICKED THE FUCKING DOOR IN, AND LEFT THE WHOLE THING SPRAWLED ACROSS THE BED LIKE A CHEAP, DISCARDED, HARBOR CRACKWHORE WITH A LOYALTY CARD FROM WALLGREENS.........                         YOUR CODE?....WAS OUR, \"MATURE\", PLAYTHING.....                         YOUR SECURITY? FUCKING GONE, LIKE YOH MOMMA'S VIRGINITY WHEN SHE WAS 11.....                         YOUR PROTECTION? DEAD AND BURIED.                         AND YOUR BEAUTIFUL LITTLE PROGRAM?                                   CRACKED, PACKED, AND HANDED TO THE WORLD FOR FUCKING FREE!!!.....                         NO KEYS.                         NO PASSWORDS.                         NO REGISTRATION.                         NO MERCY.                         JUST ANOTHER LAME PROTECTION SCHEME ADDED TO OUR TROPHY ROOM AKA \"THE NAUGHT DUNGEON\", IF YOU KNOW WHAT WE MEAN.....                         SO TO THE CODERS WHO SPENT ALL NIGHT BUILDING THIS SHIT:                         NICE TRY, BITCHES!!!                         MAYBE NEXT TIME DON'T MAKE THE PROTECTION SO FUCKING EASY.....                         AND NOW FOR THE IMPORTANT PART.....                         GREETZ FLY OUT TO ALL THE DUDES STILL HAMMERING AWAY AT THEIR KEYBOARDS AT THREE IN THE FUCKING MORNING.....                         RESPECT TO THE REAL CRACKERS, THE CODERS, THE SWAPPERS, THE BBS FREAKS, AND EVERYONE KEEPING THIS SCENE ALIVE.....                         AND TO THE LAMERS, THE SUIT WEARING BETA MEN, THE SO CALLED COPY-PROTECTION \"GURUS\", AND EVERYONE WHO THINKS THEY CAN STOP US.....                                   FUCK YOU!!!.....                         YOU CAN SEND YOUR COMPLAINTS STRAIGHT TO HELL.....                         THAT'S WHERE WE'LL BE ANYWAY.....                         EATING PUSSY, DRINKING VODKA, AND HANGING AROUND SATAN'S TABLE LIKE WE OWN THE FUCKING PLACE.....                         ADOLF'S THERE.                         MAO'S THERE.                         SATAN'S THERE.                         AND SOMEHOW WE'RE STILL THE FUCKING PARTY ANIMALS AT THE TABLE!!!.....                         HAHAHAHAHA!!!                         I GUESS WE'RE GOING TO HELL.....                         ...BUT HEY.....                         HELL AIN'T A BAD PLACE TO BE!!!.....                         PHONAX.....                         CRU JONES.....                         D'ELITE.....                         WE CAME.                         WE SAW.                         WE CRACKED.                         AND THEN WE GAVE IT AWAY!!!.....                         FUCK YOUR COPY PROTECTION!!!                         FUCK YOUR REGISTRATION SCHEME!!!                         FUCK YOUR SECURITY!!!                         SEE YOU ON THE NEXT RELEASE, YOU LAMERS.....                                   ...IF YOU CAN KEEP UP!!!.....                         D'ELITE '90.....                         SIGNING OFF.....                         KEEP CRACKING.....                         KEEP SWAPPING.....                         AND KEEP THE FUCKING SCENE ALIVE!!!"

	scroll2 := effect.NewScroller(lyrics2, candyPalette, "assets/fonts/Impact.ttf", 60, float32(rl.GetScreenHeight())-80, 5, 15, 3, false)

	// Load CRT shader
	shader := rl.LoadShader("", "assets/shaders/crt.fs")
	resolutionLoc := rl.GetShaderLocation(shader, "resolution")
	rl.SetShaderValue(shader, resolutionLoc, []float32{float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight())}, rl.ShaderUniformVec2)

	// Create render texture for post-processing
	target := rl.LoadRenderTexture(1024, 768)
	defer rl.UnloadRenderTexture(target)
	defer rl.UnloadShader(shader)

	for !rl.WindowShouldClose() {
		if rl.GetKeyPressed() > 0 {
			fire.Toggle()
		}

		fire.Process()
		pentagram.Process()
		rasterBars.Process()
		stars.Process()
		moon.Process()
		scroll.Process()
		scroll2.Process()

		// Draw to render texture
		rl.BeginTextureMode(target)
		rl.ClearBackground(rl.Black)

		moon.Draw()
		stars.Draw()

		if rasterBars.IsMovingUp() {
			pentagram.Draw()
			rasterBars.Draw()
		} else {
			rasterBars.Draw()
			pentagram.Draw()
		}

		scroll.Draw()
		fire.Draw()
		scroll2.Draw()
		rl.EndTextureMode()

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		// Draw render texture with CRT shader
		rl.BeginShaderMode(shader)
		// NOTE: Render texture must be Y-flipped because of OpenGL coordinates
		rl.DrawTextureRec(target.Texture, rl.Rectangle{X: 0, Y: 0, Width: float32(target.Texture.Width), Height: -float32(target.Texture.Height)}, rl.Vector2{X: 0, Y: 0}, rl.White)
		rl.EndShaderMode()

		rl.DrawFPS(10, 10)
		rl.EndDrawing()
	}
}
