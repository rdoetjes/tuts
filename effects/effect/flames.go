package effects

import rl "github.com/gen2brain/raylib-go/raylib"

var firePalette = []rl.Color{
	{R: 7, G: 7, B: 7, A: 255}, {R: 31, G: 7, B: 7, A: 255}, {R: 47, G: 15, B: 7, A: 255},
	{R: 71, G: 15, B: 7, A: 255}, {R: 87, G: 23, B: 7, A: 255}, {R: 103, G: 31, B: 7, A: 255},
	{R: 119, G: 31, B: 7, A: 255}, {R: 143, G: 39, B: 7, A: 255}, {R: 159, G: 47, B: 7, A: 255},
	{R: 175, G: 63, B: 7, A: 255}, {R: 191, G: 71, B: 7, A: 255}, {R: 199, G: 71, B: 7, A: 255},
	{R: 223, G: 79, B: 7, A: 255}, {R: 223, G: 87, B: 7, A: 255}, {R: 223, G: 87, B: 7, A: 255},
	{R: 215, G: 95, B: 7, A: 255}, {R: 215, G: 103, B: 15, A: 255}, {R: 207, G: 111, B: 15, A: 255},
	{R: 207, G: 119, B: 15, A: 255}, {R: 207, G: 127, B: 15, A: 255}, {R: 207, G: 135, B: 23, A: 255},
	{R: 199, G: 135, B: 23, A: 255}, {R: 199, G: 143, B: 23, A: 255}, {R: 199, G: 151, B: 31, A: 255},
	{R: 191, G: 159, B: 31, A: 255}, {R: 191, G: 159, B: 31, A: 255}, {R: 191, G: 167, B: 39, A: 255},
	{R: 191, G: 167, B: 39, A: 255}, {R: 191, G: 175, B: 47, A: 255}, {R: 183, G: 175, B: 47, A: 255},
	{R: 183, G: 183, B: 47, A: 255}, {R: 183, G: 183, B: 55, A: 255}, {R: 207, G: 207, B: 111, A: 255},
	{R: 223, G: 223, B: 159, A: 255}, {R: 239, G: 239, B: 199, A: 255}, {R: 255, G: 255, B: 255, A: 255},
}

type Flames struct {
	heatmap [][]uint8
	nRows   int
	nCols   int
	isOn    bool
	image   *rl.Image
	texture rl.Texture2D
}

func NewFlames(nCols int, nRows int) *Flames {
	heatmap := make([][]uint8, nRows)
	for i := range heatmap {
		heatmap[i] = make([]uint8, nCols)
	}

	// Create an image and texture for efficient rendering
	img := rl.GenImageColor(nCols, nRows, rl.Blank)
	tex := rl.LoadTextureFromImage(img)

	flames := Flames{
		heatmap: heatmap,
		nRows:   nRows,
		nCols:   nCols,
		isOn:    true,
		image:   img,
		texture: tex,
	}

	return &flames
}

func (s *Flames) IsBurning() bool {
	return s.isOn
}

func (s *Flames) Toggle() {
	s.isOn = !s.isOn
}

func (s *Flames) Process() {
	if !s.isOn {
		for i := 0; i < s.nCols; i++ {
			// extinquish bottom row
			if s.heatmap[s.nRows-1][i] > 0 {
				s.heatmap[s.nRows-1][i] -= 1
			}
		}
	} else {
		for i := 0; i < s.nCols; i++ {
			// Seed bottom row with max intensity
			s.heatmap[s.nRows-1][i] = uint8(rl.GetRandomValue(20, 35))
		}
	}

	// Propagate fire upwards
	// Start from row 1 to s.nRows-1 (bottom-up propagation)
	for y := 1; y < s.nRows; y++ {
		for x := 0; x < s.nCols; x++ {
			s.spreadFire(x, y)
		}
	}

	//update pixels to image structe
	s.updatePixels()
}

func (s *Flames) spreadFire(x int, y int) {
	pixel := s.heatmap[y][x]
	if pixel == 0 {
		s.heatmap[y-1][x] = 0
		return
	}

	// Random horizontal offset and decay
	rnd := int(rl.GetRandomValue(-1, 2))
	dstX := x - rnd + 1

	// Wrap or clamp horizontal offset
	if dstX < 0 {
		dstX = 0
	} else if dstX >= s.nCols {
		dstX = s.nCols - 1
	}

	decay := uint8(rnd & 1) // 0 or 1
	if pixel > decay {
		s.heatmap[y-1][dstX] = pixel - decay
	} else {
		s.heatmap[y-1][dstX] = 0
	}
}

func (s *Flames) updatePixels() {
	// Update the image pixels based on the heatmap
	for y := 0; y < s.nRows; y++ {
		for x := 0; x < s.nCols; x++ {
			intensity := s.heatmap[y][x]
			color := rl.Blank // Transparent
			if intensity > 0 {
				color = firePalette[intensity]
			}
			rl.ImageDrawPixel(s.image, int32(x), int32(y), color)
		}
	}
}

func (s *Flames) Draw() {
	// Upload the updated image to the GPU texture
	rl.UpdateTexture(s.texture, rl.LoadImageColors(s.image))

	srcRect := rl.Rectangle{X: 0, Y: 0, Width: float32(s.nCols), Height: float32(s.nRows)}
	origin := rl.Vector2{X: 0, Y: 0}
	screenWidth := float32(rl.GetScreenWidth())
	screenHeight := float32(rl.GetScreenHeight())

	// Draw floor flames (bottom 320 pixels)
	destRectFloor := rl.Rectangle{X: 0, Y: screenHeight - 320, Width: screenWidth, Height: 320}
	rl.DrawTexturePro(s.texture, srcRect, destRectFloor, origin, 0, rl.White)

	// Draw ceiling flames (top 320 pixels, flipped vertically and horizontally)
	destRectCeiling := rl.Rectangle{X: 0, Y: 0, Width: -screenWidth, Height: -320}
	srcRectFlipped := rl.Rectangle{X: 0, Y: 0, Width: -float32(s.nCols), Height: -float32(s.nRows)}
	rl.DrawTexturePro(s.texture, srcRectFlipped, destRectCeiling, origin, 0, rl.White)
}

// Ensure Flames implements BaseEffect interface
var _ BaseEffect = (*Flames)(nil)
