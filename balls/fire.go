package main

import (
	"balls/parallel"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type FireEffect struct {
	buffer     [fireWidth * fireHeight]uint8
	pixels     []rl.Color
	texture    rl.Texture2D
	texUpdated bool
}

func NewFireEffect() *FireEffect {
	return &FireEffect{
		pixels: make([]rl.Color, fireWidth*fireHeight),
	}
}

func (f *FireEffect) initTexture() {
	if f.texture.ID == 0 {
		img := rl.GenImageColor(fireWidth, fireHeight, rl.Blank)
		f.texture = rl.LoadTextureFromImage(img)
		rl.UnloadImage(img)
	}
}

func (f *FireEffect) Update(isFlatline bool) {
	// ... (source update stays the same)
	if !isFlatline {
		for x := 0; x < fireWidth; x++ {
			if rl.GetRandomValue(0, 10) > 1 {
				f.buffer[(fireHeight-1)*fireWidth+x] = uint8(rl.GetRandomValue(200, 255))
			} else {
				f.buffer[(fireHeight-1)*fireWidth+x] = uint8(rl.GetRandomValue(2, 100))
			}
		}
	} else {
		for x := 0; x < fireWidth; x++ {
			if f.buffer[(fireHeight-1)*fireWidth+x] > 8 {
				f.buffer[(fireHeight-1)*fireWidth+x] -= 1
			} else {
				f.buffer[(fireHeight-1)*fireWidth+x] = 0
			}
		}
	}

	parallel.ForEach(fireWidth, func(x int) {
		r := rand.New(rand.NewSource(int64(x + int(rl.GetRandomValue(0, 1000)))))
		for y := 1; y < fireHeight; y++ {
			srcIdx := y*fireWidth + x
			pixel := f.buffer[srcIdx]
			if pixel == 0 {
				f.buffer[(y-1)*fireWidth+x] = 0
			} else {
				randOffset := r.Intn(3) - 1
				dstX := (x + randOffset + fireWidth) % fireWidth
				cooling := uint8(1 + r.Intn(4))
				if isFlatline {
					cooling += 4
				}
				if uint32(pixel) > uint32(cooling) {
					f.buffer[(y-1)*fireWidth+dstX] = pixel - cooling
				} else {
					f.buffer[(y-1)*fireWidth+dstX] = 0
				}
			}
		}
	})

	// Prepare pixels in parallel
	parallel.ForEach(fireHeight, func(y int) {
		for x := 0; x < fireWidth; x++ {
			val := f.buffer[y*fireWidth+x]
			idx := y*fireWidth + x
			if val < 15 {
				f.pixels[idx] = rl.Blank
				continue
			}

			var color rl.Color
			if val < 70 {
				color = rl.NewColor(255, 0, 0, 120-val)
			} else if val < 150 {
				g := uint8((float64(val) - 70) * 1.5)
				color = rl.NewColor(255, val+g, 0, 200)
			} else {
				g := uint8(120 + (float64(val)-150)*1.3)
				color = rl.NewColor(255, g, 0, 255)
			}

			// Height-based alpha tapering: Tips (y=0) are transparent, Source (y=fireHeight) is opaque
			hFactor := float32(y) / float32(fireHeight)
			heightAlpha := uint8(hFactor * 255)
			if color.A > heightAlpha {
				color.A = heightAlpha
			}
			f.pixels[idx] = color
		}
	})
	f.texUpdated = false
}

func (f *FireEffect) Draw() {
	f.initTexture()
	if !f.texUpdated {
		rl.UpdateTexture(f.texture, f.pixels)
		f.texUpdated = true
	}

	destRect := rl.NewRectangle(0, screenHeight-220, screenWidth, 250)
	sourceRect := rl.NewRectangle(0, 0, fireWidth, fireHeight)

	rl.DrawTexturePro(f.texture, sourceRect, destRect, rl.NewVector2(0, 0), 0, rl.White)

	destRect.Y = -30
	sourceRect.Height = -fireHeight // Flip texture
	rl.DrawTexturePro(f.texture, sourceRect, destRect, rl.NewVector2(0, 0), 0, rl.White)
}
