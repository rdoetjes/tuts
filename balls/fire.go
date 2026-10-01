package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type FireEffect struct {
	buffer [fireWidth * fireHeight]uint8
}

func NewFireEffect() *FireEffect {
	return &FireEffect{}
}

func (f *FireEffect) Update(isFlatline bool) {
	if !isFlatline {
		// Randomize bottom row (fire source) with "hot spots" for clumping
		for x := 0; x < fireWidth; x++ {
			if rl.GetRandomValue(0, 10) > 1 {
				f.buffer[(fireHeight-1)*fireWidth+x] = uint8(rl.GetRandomValue(200, 255))
			} else {
				f.buffer[(fireHeight-1)*fireWidth+x] = uint8(rl.GetRandomValue(0, 100))
			}
		}
	} else {
		// Cool down the source when flatlined
		for x := 0; x < fireWidth; x++ {
			if f.buffer[(fireHeight-1)*fireWidth+x] > 8 {
				f.buffer[(fireHeight-1)*fireWidth+x] -= 8
			} else {
				f.buffer[(fireHeight-1)*fireWidth+x] = 0
			}
		}
	}

	// Propagate fire upwards with random horizontal drift
	for x := 0; x < fireWidth; x++ {
		for y := 1; y < fireHeight; y++ {
			// Get heat from below
			srcIdx := y*fireWidth + x
			pixel := f.buffer[srcIdx]

			if pixel == 0 {
				f.buffer[(y-1)*fireWidth+x] = 0
			} else {
				// Random horizontal drift (x - 1, x, or x + 1)
				randOffset := rl.GetRandomValue(0, 2) - 1
				dstX := (x + int(randOffset) + fireWidth) % fireWidth

				// Random cooling
				coolingBase := int32(1)
				coolingMax := int32(4)
				if isFlatline {
					coolingBase = 3
					coolingMax = 6
				}
				cooling := uint8(rl.GetRandomValue(coolingBase, coolingMax))

				if uint32(pixel) > uint32(cooling) {
					f.buffer[(y-1)*fireWidth+dstX] = pixel - cooling
				} else {
					f.buffer[(y-1)*fireWidth+dstX] = 0
				}
			}
		}
	}
}

func (f *FireEffect) Draw(isTop bool) {
	scaleX := float32(screenWidth) / float32(fireWidth)
	scaleY := float32(250) / float32(fireHeight)
	startY := float32(screenHeight - 220)

	if isTop {
		startY = -30 // Mirror position for top
	}

	for y := 0; y < fireHeight; y++ {
		drawY := y
		if isTop {
			drawY = fireHeight - 1 - y
		}

		for x := 0; x < fireWidth; x++ {
			val := f.buffer[y*fireWidth+x]
			if val < 10 {
				continue
			}

			var color rl.Color
			if val < 70 {
				color = rl.NewColor(val*3, 0, 0, uint8(val*2))
			} else if val < 150 {
				g := uint8((float64(val) - 70) * 1.5)
				color = rl.NewColor(255, g, 0, 200)
			} else {
				g := uint8(120 + (float64(val)-150)*1.3)
				color = rl.NewColor(255, g, 0, 255)
			}

			// Height-based alpha tapering
			hFactor := float32(y) / float32(fireHeight)
			if isTop {
				hFactor = float32(fireHeight-1-y) / float32(fireHeight)
			}

			heightAlpha := uint8(hFactor * 255)
			if color.A > heightAlpha {
				color.A = heightAlpha
			}

			posY := startY + float32(drawY)*scaleY
			rl.DrawRectangle(int32(float32(x)*scaleX), int32(posY), int32(scaleX)+1, int32(scaleY)+1, color)
		}
	}
}
