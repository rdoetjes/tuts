package cvhelper

import (
	"fmt"
	"image"
	"image/color"
	"time"

	// Import gocv package for OpenCV wrappers and bindings
	"gocv.io/x/gocv"
)

// Fps struct to track FPS calculation related data (statefull)
type Fps struct {
	avg_tot    float64
	framecount int64
	last_time  int64
}

func (fps *Fps) Init() {
	fps.avg_tot = 0
	fps.framecount = 0
	fps.last_time = time.Now().UnixNano()
}

func NewFps() *Fps {
	fps := Fps{}
	fps.Init()
	return &fps
}

// Converts Frame from BFR to HSV and applies the given HSV range
// returns the filtered frame in the target Mat
func FilterHue(frame *gocv.Mat, target *gocv.Mat, hue_min gocv.Scalar, hue_max gocv.Scalar) {
	gocv.CvtColor(*frame, target, gocv.ColorBGRToHSV)
	gocv.InRangeWithScalar(*frame, hue_min, hue_max, target)
}

// Solely a wrapper around InRangeWithScalar to filter the frame in the RGB range
func FilterRGB(frame *gocv.Mat, target *gocv.Mat, rgb_min gocv.Scalar, rgb_max gocv.Scalar) {
	gocv.InRangeWithScalar(*frame, rgb_min, rgb_max, target)
}

// AddFpsOnFrame draws the calculated FPS on the frame
func AddFpsOnFrame(frame *gocv.Mat, fps_data *Fps) {
	// store the new data in the struct
	current_time := time.Now().UnixNano()
	fps_data.avg_tot += 1 / ((float64(current_time) - float64(fps_data.last_time)) / 1e9)
	fps_data.framecount += 1

	// calculate fps
	fps := fmt.Sprintf("%.0f", fps_data.avg_tot/float64(fps_data.framecount))

	// draw the fps on the frame
	pt := image.Point{X: 10, Y: 20}
	gocv.PutText(frame, fps, pt, gocv.FontHersheyPlain, 1.2, color.RGBA{255, 0, 255, 1}, 2)

	// update the last_time to the current time, for the next run
	fps_data.last_time = current_time
}
