package main

import (
	effects "effects/effect"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(1024, 768, "EFFECTS")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	var fire *effects.Flames = effects.NewFlames(400, 150)

	for !rl.WindowShouldClose() {
		if rl.GetKeyPressed() > 0 {
			fire.Toggle()
		}

		fire.Process()

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		fire.Draw()

		rl.DrawFPS(10, 10)
		rl.EndDrawing()
	}
}
