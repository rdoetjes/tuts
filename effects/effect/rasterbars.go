package effect

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type RasterBars struct {
	palette []rl.Color
	timer   float32
}

func NewRasterBars(pal []rl.Color) *RasterBars {
	return &RasterBars{
		palette: pal,
		timer:   0,
	}
}

func (r *RasterBars) Process() {
	r.timer += 0.06
}

func (r *RasterBars) Draw() {
	screenWidth := float32(rl.GetScreenWidth())
	screenHeight := float32(rl.GetScreenHeight())

	numBars := 8
	barHeight := float32(28.0) // Slightly taller to accommodate internal shading
	gap := float32(5.0)        // 5 pixel gap as requested

	for i := 0; i < numBars; i++ {
		bounce := float64(r.timer * 0.8)
		offset := float32(math.Sin(bounce)) * (screenHeight * 0.3)

		// Base Y calculation including the gap
		totalBarStep := barHeight + gap
		y := screenHeight/2 + offset - (float32(numBars) * totalBarStep / 2) + (float32(i) * totalBarStep)

		shift := int(r.timer * 5)
		colorIdx := (i + shift) % len(r.palette) / 2
		baseColor := r.palette[colorIdx]

		// Draw rounded shading for each bar
		// We use 8 internal slices to create a convex/round appearance
		slices := 12
		for s := 0; s < slices; s++ {
			t := float32(s) / float32(slices-1)

			// Triangle wave for shading (0.0 at edges, 1.0 at center)
			shade := 1.0 - math.Abs(float64(t-0.5)*2.0)

			// Modulate color: darken edges, brighten center
			r_val := uint8(float32(baseColor.R) * float32(0.3+shade*0.7))
			g_val := uint8(float32(baseColor.G) * float32(0.3+shade*0.7))
			b_val := uint8(float32(baseColor.B) * float32(0.3+shade*0.7))

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
	// The derivative of sin(r.timer * 0.8) is 0.8 * cos(r.timer * 0.8)
	// If cos is negative, it's moving UP (Y value is decreasing)
	return math.Cos(float64(r.timer*0.8)) < 0
}

var _ BaseEffect = (*RasterBars)(nil)
