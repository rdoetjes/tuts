package effect

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const BounceSpeed = 0.8

type RasterBars struct {
	palette           []rl.Color
	barColors         []rl.Color
	timer             float32
	currentScale      float32
	currentBrightness float32
}

func NewRasterBars(pal []rl.Color) *RasterBars {
	numBars := 8
	barColors := make([]rl.Color, numBars)

	// Specific 80s color sequence: Brown, Red, Orange, Yellow
	sequence := []rl.Color{
		{110, 50, 20, 255}, // 80s Brown
		{220, 0, 0, 255},   // Red
		{255, 120, 0, 255}, // Orange
		{255, 220, 0, 255}, // Yellow
		{255, 120, 0, 255}, // Orange
		{220, 0, 0, 255},   // Red
		{110, 50, 20, 255}, // 80s brown
	}

	for i := 0; i < numBars; i++ {
		// Assign colors in the requested repeating order
		barColors[i] = sequence[i%len(sequence)]
	}

	return &RasterBars{
		palette:           pal,
		barColors:         barColors,
		timer:             0,
		currentScale:      1.0,
		currentBrightness: 1.0,
	}
}

func (r *RasterBars) Process() {
	r.timer += 0.06

	// Target scale: 1.0 when moving up, 0.8 when moving down
	targetScale := float32(1.0)
	targetBrightness := float32(1.0)
	if !r.IsMovingUp() {
		targetScale = 0.8
		targetBrightness = 0.8
	}

	// Smoothly transition over 5 frames (step = 0.2 / 5 = 0.04)
	step := float32(0.06)

	// Scale transition
	if r.currentScale < targetScale {
		r.currentScale += step
		if r.currentScale > targetScale {
			r.currentScale = targetScale
		}
	} else if r.currentScale > targetScale {
		r.currentScale -= step
		if r.currentScale < targetScale {
			r.currentScale = targetScale
		}
	}

	// Brightness transition
	if r.currentBrightness < targetBrightness {
		r.currentBrightness += step
		if r.currentBrightness > targetBrightness {
			r.currentBrightness = targetBrightness
		}
	} else if r.currentBrightness > targetBrightness {
		r.currentBrightness -= step
		if r.currentBrightness < targetBrightness {
			r.currentBrightness = targetBrightness
		}
	}
}

func (r *RasterBars) Draw() {
	screenWidth := float32(rl.GetScreenWidth())
	screenHeight := float32(rl.GetScreenHeight())

	numBars := 7
	barHeight := float32(28.0) * r.currentScale
	gap := float32(3.0)

	for i := 0; i < numBars; i++ {
		bounce := float64(r.timer * BounceSpeed)
		offset := float32(math.Sin(bounce)) * (screenHeight * 0.3)

		// Base Y calculation including the gap
		totalBarStep := barHeight + gap
		y := screenHeight/2 + offset - (float32(numBars) * totalBarStep / 2) + (float32(i) * totalBarStep)

		// Use the fixed random color assigned to this bar
		baseColor := r.barColors[i]

		// Draw rounded shading for each bar
		// We use 8 internal slices to create a convex/round appearance
		slices := 12
		for s := 0; s < slices; s++ {
			t := float32(s) / float32(slices-1)

			// Triangle wave for shading (0.0 at edges, 1.0 at center)
			shade := 1.0 - math.Abs(float64(t-0.5)*2.0)

			// Modulate color: darken edges, brighten center, apply currentBrightness
			r_val := uint8(float32(baseColor.R) * float32(0.3+shade*0.7) * r.currentBrightness)
			g_val := uint8(float32(baseColor.G) * float32(0.3+shade*0.7) * r.currentBrightness)
			b_val := uint8(float32(baseColor.B) * float32(0.3+shade*0.7) * r.currentBrightness)

			sliceHeight := barHeight / float32(slices)
			sliceY := y + float32(s)*sliceHeight

			sliceColor := rl.Color{r_val, g_val, b_val, 255}

			// Highlight on top slice
			if s == 2 {
				sliceColor.R = uint8(math.Min(255, float64(sliceColor.R)+40))
				sliceColor.G = uint8(math.Min(255, float64(sliceColor.G)+40))
				sliceColor.B = uint8(math.Min(255, float64(sliceColor.B)+40))
			}

			rl.DrawRectangle(0, int32(sliceY), int32(screenWidth), int32(sliceHeight)+1, sliceColor)
		}
	}
}

func (r *RasterBars) IsMovingUp() bool {
	// The derivative of sin(r.timer * BounceSpeed) is BounceSpeed * cos(r.timer * BounceSpeed)
	// If cos is negative, it's moving UP (Y value is decreasing)
	return math.Cos(float64(r.timer*BounceSpeed)) < 0
}

var _ BaseEffect = (*RasterBars)(nil)
