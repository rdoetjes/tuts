package main

import (
	"fmt"
	"image"
	"image/color"
	"time"

	"gocv.io/x/gocv"
)

type fps struct {
	avg_tot    float64
	framecount int64
	last_time  int64
}

func AddFpsOnFrame(frame *gocv.Mat, fps_data *fps) {
	current_time := time.Now().UnixNano()
	fps_data.avg_tot += 1 / ((float64(current_time) - float64(fps_data.last_time)) / 1e9)
	fps_data.framecount += 1
	fps := fmt.Sprintf("%.0f", fps_data.avg_tot/float64(fps_data.framecount))
	pt := image.Point{X: 100, Y: 100}
	gocv.PutText(frame, fps, pt, gocv.FontHersheyPlain, 1.2, color.RGBA{255, 0, 255, 1}, 2)
	fps_data.last_time = current_time
}

func main() {
	fps_data := fps{avg_tot: 0, framecount: 0, last_time: 0}
	webcam, err := gocv.VideoCaptureDevice(1)
	if err != nil {
		panic(err)
	}
	window := gocv.NewWindow("WebCam")
	img := gocv.NewMat()
	for {
		webcam.Read(&img)
		if window.WaitKey(1) != -1 {
			break
		}
		AddFpsOnFrame(&img, &fps_data)
		window.IMShow(img)
	}
}
