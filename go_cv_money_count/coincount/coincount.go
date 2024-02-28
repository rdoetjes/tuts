package coincount

import (
	"image"

	"gocv.io/x/gocv"
)

type CoinProcessing struct {
	kernel       image.Point
	cannyThresh1 float32
	cannyThresh2 float32
}

func NewDefaultCoinProcessing() *CoinProcessing {
	return &CoinProcessing{
		kernel:       image.Pt(5, 5),
		cannyThresh1: 50,
		cannyThresh2: 150,
	}
}

// Example function that turns the image into black and white and flips it
// over all the axes.
func preProcessForCounCount(input *gocv.Mat, process *gocv.Mat, config *CoinProcessing) {
	gocv.CvtColor(*input, process, gocv.ColorBGRToGray)
	gocv.GaussianBlur(*process, process, config.kernel, 0, 0, gocv.BorderDefault)
	gocv.Canny(*process, process, config.cannyThresh1, config.cannyThresh2)
}

func CointEuros(input *gocv.Mat, process *gocv.Mat, config *CoinProcessing) float64 {
	preProcessForCounCount(input, process, config)
	return 0.0
}
