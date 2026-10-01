package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type CopperBars struct{}

func (cb *CopperBars) Draw(timer float64) {
	numBars := 6
	barHeight := float32(30)
	for i := 0; i < numBars; i++ {
		yPos := float32(screenHeight/2-100) + float32(math.Sin(timer*1.5+float64(i)*0.4)*250)

		for j := 0; j < int(barHeight); j++ {
			intensity := uint8(255 - math.Abs(float64(j)-float64(barHeight/2))*15)

			var color rl.Color
			if i%2 == 0 {
				color = rl.NewColor(255, 0, 0, 180) // Vivid Red
			} else {
				color = rl.NewColor(255, 230, 0, 180) // Bright Yellow
			}
			color.A = uint8(float64(intensity) * 0.7)
			rl.DrawLine(0, int32(yPos)+int32(j), screenWidth, int32(yPos)+int32(j), color)
		}
	}
}
