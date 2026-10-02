package main

import (
	"effects/effect"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(1024, 768, "EFFECTS")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	fire := effect.NewFlames(400, 150)
	pentagram := effect.NewPentagram()
	rasterBars := effect.NewRasterBars(effect.FirePalette)
	stars := effect.NewStarfield(400)

	lyrics := "HEY, YOU ... YEAH, YOU ... SOMETIMES I THINK THIS WOMAN IS KINDA HOT ... SOMETIMES I THINK THIS WOMAN IS SOMETIMES NOT ... PUTS ME DOWN, FOOL ME AROUND ... WHY SHE DO IT TO ME? ... OUT FOR SATISFACTION, ANY PIECE OF ACTION ... THAT AIN'T THE WAY IT SHOULD BE ... SHE NEEDS LOVIN', KNOWS I'M THE MAN ... SHE'S GOTTA SEE ... POURS MY BEER, LICKS MY EAR ... BRINGS OUT THE DEVIL IN ME ... HELL AIN'T A BAD PLACE TO BE ... SPEND MY MONEY, DRINKS MY BOOZE ... STAYS OUT EVERY NIGHT ... BUT I GOT TO THINKIN', 'HEY, JUST A MINUTE, SOMETHIN' AIN'T RIGHT' ... OH, DISILLUSIONS AND CONFUSION ... MAKE ME WANNA CRY ... ALL THE SHAME, YOU PLAYIN' YOUR GAMES ... TELLIN' ME THOSE LIES ... DON'T MIND HER PLAYIN' DEMON ... AS LONG AS IT'S WITH ME ... IF THIS IS HELL, THEN YOU COULD SAY IT'S HEAVENLY ... HELL AIN'T A BAD PLACE TO BE ... END OF THE NIGHT, TURNS DOWN THE LIGHT ... CLOSES UP ON ME ... OPENS MY HEART, TEARS IT APART ... BRINGS OUT THE DEVIL IN ME ... HELL AIN'T A BAD PLACE TO BE ... I SAID, HELL AIN'T A BAD PLACE TO BE ... HELL AIN'T A BAD PLACE TO BE ... YOU KNOW THAT HELL AIN'T A BAD PLACE TO BE"

	scroll := effect.NewScroller(lyrics, effect.FirePalette, "assets/fonts/Menlo.ttf", 400, float32(rl.GetScreenHeight())/2-250, 8, 150, 100, 12, true)

	candyPalette := []rl.Color{
		{255, 0, 50, 255},  // Bright Candy Red
		{255, 120, 0, 255}, // Neon Orange
		{255, 255, 0, 255}, // Electric Yellow
		{110, 50, 20, 255}, // Rich Chocolate Brown
	}

	lyrics2 := "-------------------------PHONAX AND CRU JONES, FUCKING DID IT AGAIN! WE CRACKED THIS BITCH, WE SURE AS HELL ARE GOING TO HELL... BUT HEY... HELL AIN'T A BAD PLACE TO BE! YOU'LL FIND US EATING PUSSY AND DRINKING VODKA NEAR SATAN, ADOLF, MAO AND STALIN'S TABLE!.... HELL YEAH!!!!----------------------------------------"
	scroll2 := effect.NewScroller(lyrics2, candyPalette, "assets/fonts/Impact.ttf", 60, float32(rl.GetScreenHeight())-80, 5, 45, 15, 3, false)

	for !rl.WindowShouldClose() {
		if rl.GetKeyPressed() > 0 {
			fire.Toggle()
		}

		fire.Process()
		pentagram.Process()
		rasterBars.Process()
		stars.Process()
		scroll.Process()
		scroll2.Process()

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

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

		rl.DrawFPS(10, 10)
		rl.EndDrawing()
	}
}
