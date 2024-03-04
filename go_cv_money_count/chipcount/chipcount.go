package chipcount

import (
	"go/cv/cvhelper"
	"image"

	"gocv.io/x/gocv"
)

type CoinProcessing struct {
	kernel       image.Point
	cannyThresh1 float32
	cannyThresh2 float32
	r_red_lb     gocv.Scalar
	r_red_ub     gocv.Scalar
}

func NewDefaultCoinProcessing() *CoinProcessing {
	return &CoinProcessing{
		kernel:       image.Pt(13, 13),
		cannyThresh1: 50,
		cannyThresh2: 150,
		r_red_lb:     gocv.NewScalar(190, 90, 90, 0),
		r_red_ub:     gocv.NewScalar(255, 150, 150, 0),
	}
}

// Example function that turns the image into black and white and flips it
// over all the axes.
func preProcessForChipCount(input *gocv.Mat, process *gocv.Mat, config *CoinProcessing) {
	gocv.CvtColor(*input, process, gocv.ColorBGRToRGB)
	gocv.GaussianBlur(*process, process, config.kernel, 0, 0, gocv.BorderDefault)
}

func CountRedChips(input *gocv.Mat, output *gocv.Mat, config *CoinProcessing) {
	cvhelper.FilterRGB(input, output, config.r_red_lb, config.r_red_ub)
	kernel_erode := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(3, 3))
	gocv.Erode(*output, output, kernel_erode)

	kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(15, 15))
	gocv.Dilate(*output, output, kernel)
}

func CountChips(input *gocv.Mat, process *gocv.Mat, config *CoinProcessing) int {
	preProcessForChipCount(input, process, config)
	CountRedChips(process, process, config)
	vec := gocv.FindContours(*process, gocv.RetrievalExternal, gocv.ChainApproxSimple)
	return vec.Size()
}
