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

	lyrics := "HEY, YOU ... YEAH, YOU ... SOMETIMES I THINK THIS WOMAN IS KINDA HOT ... SOMETIMES I THINK THIS WOMAN IS SOMETIMES NOT ... PUTS ME DOWN, FOOL ME AROUND ... WHY SHE DO IT TO ME? ... OUT FOR SATISFACTION, ANY PIECE OF ACTION ... THAT AIN'T THE WAY IT SHOULD BE ... SHE NEEDS LOVIN', KNOWS I'M THE MAN ... SHE'S GOTTA SEE ... POURS MY BEER, LICKS MY EAR ... BRINGS OUT THE DEVIL IN ME ... HELL AIN'T A BAD PLACE TO BE ... SPEND MY MONEY, DRINKS MY BOOZE ... STAYS OUT EVERY NIGHT ... BUT I GOT TO THINKIN', 'HEY, JUST A MINUTE, SOMETHIN' AIN'T RIGHT' ... OH, DISILLUSIONS AND CONFUSION ... MAKE ME WANNA CRY ... ALL THE SHAME, YOU PLAYIN' YOUR GAMES ... TELLIN' ME THOSE LIES ... DON'T MIND HER PLAYIN' DEMON ... AS LONG AS IT'S WITH ME ... IF THIS IS HELL, THEN YOU COULD SAY IT'S HEAVENLY ... HELL AIN'T A BAD PLACE TO BE ... END OF THE NIGHT, TURNS DOWN THE LIGHT ... CLOSES UP ON ME ... OPENS MY HEART, TEARS IT APART ... BRINGS OUT THE DEVIL IN ME ... HELL AIN'T A BAD PLACE TO BE ... I SAID, HELL AIN'T A BAD PLACE TO BE ... HELL AIN'T A BAD PLACE TO BE ... YOU KNOW THAT HELL AIN'T A BAD PLACE TO BE"

	scroll := effect.NewScroller(lyrics, effect.FirePalette)

	for !rl.WindowShouldClose() {
		if rl.GetKeyPressed() > 0 {
			fire.Toggle()
		}

		fire.Process()
		scroll.Process()

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		scroll.Draw()
		fire.Draw()

		rl.DrawFPS(10, 10)
		rl.EndDrawing()
	}
}
