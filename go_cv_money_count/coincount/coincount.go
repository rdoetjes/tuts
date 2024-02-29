package coincount

import (
	"fmt"
	"image"
	"image/color"

	"gocv.io/x/gocv"
)

type CoinProcessing struct {
	kernel       image.Point
	cannyThresh1 float32
	cannyThresh2 float32
}

func NewDefaultCoinProcessing() *CoinProcessing {
	return &CoinProcessing{
		kernel:       image.Pt(7, 7),
		cannyThresh1: 50,
		cannyThresh2: 150,
	}
}

// Example function that turns the image into black and white and flips it
// over all the axes.
func preProcessForChipCount(input *gocv.Mat, process *gocv.Mat, config *CoinProcessing) {
	gocv.CvtColor(*input, process, gocv.ColorBGRToHLS)
	gocv.GaussianBlur(*process, process, config.kernel, 0, 0, gocv.BorderDefault)
	gocv.InRangeWithScalar(*process, gocv.NewScalar(0, 10, 10, 0), gocv.NewScalar(15, 155, 155, 0), process)
}

func getContours(input *gocv.Mat, process *gocv.Mat) int {
	var result int = 0
	circles := gocv.NewMat()
	defer circles.Close()
	//25 48
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
			gocv.Circle(input, image.Pt(x, y), r, color.RGBA{0, 0, 255, 0}, 2)
			gocv.Circle(input, image.Pt(x, y), 2, color.RGBA{255, 0, 255, 0}, 3)
		}
	}
	return result
}

func CountChips(input *gocv.Mat, process *gocv.Mat, config *CoinProcessing) int {
	preProcessForChipCount(input, process, config)
	return getContours(input, process)
}
