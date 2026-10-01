package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type HeartMonitor struct {
	prevX, prevY float32
	hasPrev      bool
}

func (hm *HeartMonitor) Draw(centerX, centerY int32, timer float64, isFlatline bool) {
	width := int32(180)
	height := int32(100)
	x := centerX - width/2
	y := centerY - height/2

	rl.DrawRectangle(x, y, width, height, rl.NewColor(30, 0, 0, 180))
	rl.DrawRectangleLinesEx(rl.NewRectangle(float32(x), float32(y), float32(width), float32(height)), 2, rl.Gold)

	hm.hasPrev = false
	for i := int32(0); i < width; i++ {
		var h float32 = 0
		if !isFlatline {
			histOffset := (timer * 2.0) - float64(i)*0.01
			t := math.Mod(histOffset, 1.0)
			if t < 0 {
				t += 1.0
			}

			idx := int(t*float64(heartbeatSize)) % heartbeatSize
			h = heartbeatTable[idx] * float32(height/2-10)
		}

		colorT := float64(math.Sin(timer*4.0-float64(i)*0.04))*0.5 + 0.5
		g := uint8(colorT * 230)
		traceColor := rl.NewColor(255, g, 0, 255)

		px := float32(x + width - i)
		py := float32(y) + float32(height)/2 - h

		if hm.hasPrev {
			rl.DrawLineEx(rl.NewVector2(hm.prevX, hm.prevY), rl.NewVector2(px, py), 3, traceColor)
		}

		if !isFlatline && i < 20 {
			alpha := uint8(255 - i*12)
			glowColor := rl.NewColor(255, g, 0, alpha)
			rl.DrawCircle(int32(px), int32(py), 2, glowColor)
		}

		hm.prevX, hm.prevY = px, py
		hm.hasPrev = true
	}
}
