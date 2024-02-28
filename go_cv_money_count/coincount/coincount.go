package coincount

import (
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
		kernel:       image.Pt(5, 5),
		cannyThresh1: 50,
		cannyThresh2: 150,
	}
}

// Example function that turns the image into black and white and flips it
// over all the axes.
func preProcessForCoinCount(input *gocv.Mat, process *gocv.Mat, config *CoinProcessing) {
	gocv.CvtColor(*input, process, gocv.ColorBGRToGray)
	gocv.GaussianBlur(*process, process, config.kernel, 0, 0, gocv.BorderDefault)
	gocv.Canny(*process, process, config.cannyThresh1, config.cannyThresh2)
}

func getContours(process *gocv.Mat) {
	circles := gocv.NewMat()
	defer circles.Close()
	gocv.HoughCirclesWithParams(*process, &circles, gocv.HoughGradient, 1, float64(process.Rows()/8), 100, 100, 10, 300)
	for i := 0; i < circles.Cols(); i++ {
		v := circles.GetVecfAt(0, i)
		// if circles are found
		if len(v) > 2 {
			x := int(v[0])
			y := int(v[1])
			r := int(v[2])

			gocv.Circle(process, image.Pt(x, y), r, color.RGBA{0, 0, 255, 0}, 2)
			gocv.Circle(process, image.Pt(x, y), 2, color.RGBA{255, 0, 255, 0}, 3)
		}
	}

}

func CountEuros(input *gocv.Mat, process *gocv.Mat, config *CoinProcessing) float64 {
	preProcessForCoinCount(input, process, config)
	getContours(process)
	return 0.0
}
