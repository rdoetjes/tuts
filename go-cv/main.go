package main

import (
	"fmt"
	"image"
	"image/color"
	"time"

	"gocv.io/x/gocv"
)

func add_fps(frame *gocv.Mat, avg_tot float64, framecount int64, last_time int64) {
	current_time := time.Now().UnixNano()
	avg_tot += 1 / ((float64(current_time) - float64(last_time)) / 1e9)
	framecount += 1

	fps := fmt.Sprintf("%.0f", avg_tot/float64(framecount))
	pt := image.Point{X: 100, Y: 100}
	gocv.PutText(frame, fps, pt, gocv.FontHersheyPlain, 1.2, color.RGBA{255, 0, 255, 1}, 2)
}

func main() {
	webcam, _ := gocv.VideoCaptureDevice(1)
	window := gocv.NewWindow("Wbecam")
	img := gocv.NewMat()
	var framecount int64 = 0
	var avg_tot float64 = 0

	for {
		last_time := time.Now().UnixNano()
		webcam.Read(&img)
		if window.WaitKey(1) != -1 {
			break
		}
		add_fps(&img, avg_tot, framecount, last_time)
		window.IMShow(img)
	}
}
