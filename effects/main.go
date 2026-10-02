package main

import (
	"effects/effect"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(1024, 768, "EFFECTS")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	cols, rows := 400, 150
	fire := effect.NewFlames(cols, rows)
	lyrics := "HELL AIN'T A BAD PLACE TO BE ... ALL THE GRIEF YOU GIVE ME ... ALL THE PAIN YOU PUT ME THROUGH ... WELL, I'M COMING HOME TO YOU ... HELL AIN'T A BAD PLACE TO BE"
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
