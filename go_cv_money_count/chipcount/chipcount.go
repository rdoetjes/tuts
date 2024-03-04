package chipcount

import (
	"fmt"
	"go/cv/cvhelper"
	"image"

	"gocv.io/x/gocv"
)

const RED = 0
const BLUE = 1
const ORANGE = 2

type CoinProcessing struct {
	kernel       image.Point
	cannyThresh1 float32
	cannyThresh2 float32
	red_lb       gocv.Scalar
	red_ub       gocv.Scalar
	blue_lb      gocv.Scalar
	blue_ub      gocv.Scalar
	orange_lb    gocv.Scalar
	orange_ub    gocv.Scalar
}

func NewDefaultCoinProcessing() *CoinProcessing {
	return &CoinProcessing{
		kernel:       image.Pt(13, 13),
		cannyThresh1: 50,
		cannyThresh2: 150,
		red_lb:       gocv.NewScalar(190, 90, 90, 0),
		red_ub:       gocv.NewScalar(255, 150, 150, 0),
		blue_lb:      gocv.NewScalar(70, 70, 130, 0),
		blue_ub:      gocv.NewScalar(170, 170, 255, 0),
		orange_lb:    gocv.NewScalar(200, 200, 120, 0),
		orange_ub:    gocv.NewScalar(255, 233, 155, 0),
	}
}

// Example function that turns the image into black and white and flips it
// over all the axes.
func preProcessForChipCount(input *gocv.Mat, process *gocv.Mat, config *CoinProcessing) {
	gocv.CvtColor(*input, process, gocv.ColorBGRToRGB)
	gocv.GaussianBlur(*process, process, config.kernel, 0, 0, gocv.BorderDefault)
}

func MaskColorChips(input *gocv.Mat, output *gocv.Mat, color int, config *CoinProcessing) {
	switch color {
	case RED:
		cvhelper.FilterRGB(input, output, config.red_lb, config.red_ub)
	case BLUE:
		cvhelper.FilterRGB(input, output, config.blue_lb, config.blue_ub)
	case ORANGE:
		cvhelper.FilterRGB(input, output, config.orange_lb, config.orange_ub)
	}

	kernel_erode := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(3, 3))
	gocv.Erode(*output, output, kernel_erode)

	kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(18, 18))
	gocv.Dilate(*output, output, kernel)
}

func CountChips(input *gocv.Mat, process *gocv.Mat, config *CoinProcessing) string {
	preProcessForChipCount(input, process, config)

	red := process.Clone()
	defer red.Close()
	MaskColorChips(&red, process, RED, config)
	vec_r := gocv.FindContours(*process, gocv.RetrievalExternal, gocv.ChainApproxSimple)

	blue := process.Clone()
	defer blue.Close()
	MaskColorChips(&blue, &blue, BLUE, config)
	vec_b := gocv.FindContours(blue, gocv.RetrievalExternal, gocv.ChainApproxSimple)

	orange := process.Clone()
	defer orange.Close()
	MaskColorChips(&orange, &orange, ORANGE, config)
	vec_y := gocv.FindContours(orange, gocv.RetrievalExternal, gocv.ChainApproxSimple)

	return fmt.Sprintf("r: %d y: %d b: %d ", vec_r.Size(), vec_y.Size(), vec_b.Size())
}
