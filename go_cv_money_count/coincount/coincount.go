package coincount

import (
	"fmt"
	"image"

	"gocv.io/x/gocv"
)

type CoinProcessing struct {
	kernel   image.Point
	hsvLower gocv.Scalar
	hsvUpper gocv.Scalar
}

func NewDefaultCoinProcessing() *CoinProcessing {
	return &CoinProcessing{
		kernel:   image.Pt(7, 7),
		hsvLower: gocv.NewScalar(0, 10, 10, 0),
		hsvUpper: gocv.NewScalar(15, 155, 155, 0),
	}
}

// Example function that turns the image into black and white and flips it
// over all the axes.
func preProcessForChipCount(input *gocv.Mat, process *gocv.Mat, config *CoinProcessing) {
	gocv.CvtColor(*input, process, gocv.ColorBGRToHLS)
	gocv.GaussianBlur(*process, process, config.kernel, 0, 0, gocv.BorderDefault)
	gocv.InRangeWithScalar(*process, config.hsvLower, config.hsvUpper, process)
}

func getContours(process *gocv.Mat) int {
	var result int = 0
	circles := gocv.NewMat()
	defer circles.Close()
	// 25 48
	gocv.HoughCirclesWithParams(*process, &circles, gocv.HoughGradient, 1, float64(process.Rows()/8), 15, 30, 10, 0) // maxRadius ste yo 0 for some readon it doesn;t seem to work
	for i := 0; i < circles.Cols(); i++ {
		v := circles.GetVecfAt(0, i)

		// if circles are found
		if len(v) > 2 {
			x := int(v[0])
			y := int(v[1])
			r := int(v[2])

			if r > 80 {
				continue
			}

			result += 1
			fmt.Printf("Circle detected at (%d, %d) with radius %d\n", x, y, r)
			// gocv.Circle(input, image.Pt(x, y), r, color.RGBA{0, 0, 255, 0}, 2)
			// gocv.Circle(input, image.Pt(x, y), 2, color.RGBA{255, 0, 255, 0}, 3)
		}
	}
	return result
}

func CountBrownChips(input *gocv.Mat, process *gocv.Mat, config *CoinProcessing) int {
	result := 0

	// brown chips
	config.hsvLower = gocv.NewScalar(0, 10, 10, 0)
	config.hsvUpper = gocv.NewScalar(15, 155, 155, 0)
	preProcessForChipCount(input, process, config) // copper euros
	result += getContours(process)
	return result
}

func CountBlueChips(input *gocv.Mat, process *gocv.Mat, config *CoinProcessing) int {
	result := 0

	// brown chips
	config.hsvLower = gocv.NewScalar(108, 10, 10, 0)
	config.hsvUpper = gocv.NewScalar(122, 155, 155, 0)
	preProcessForChipCount(input, process, config) // copper euros
	result += getContours(process)
	return result
}
