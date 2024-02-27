package main

import (
	"fmt"
	"image"
	"image/color"
	"time"

	"gocv.io/x/gocv"
)

type Fps struct {
	avg_tot    float64
	framecount int64
	last_time  int64
}

func AddFpsOnFrame(frame *gocv.Mat, fps_data *Fps) {
	current_time := time.Now().UnixNano()
	fps_data.avg_tot += 1 / ((float64(current_time) - float64(fps_data.last_time)) / 1e9)
	fps_data.framecount += 1
	fps := fmt.Sprintf("%.0f", fps_data.avg_tot/float64(fps_data.framecount))
	pt := image.Point{X: 10, Y: 20}
	gocv.PutText(frame, fps, pt, gocv.FontHersheyPlain, 1.2, color.RGBA{255, 0, 255, 1}, 2)
	fps_data.last_time = current_time
}

func main() {
	fps_data := Fps{avg_tot: 0, framecount: 0, last_time: 0}
	webcam, err := gocv.VideoCaptureDevice(1)
	if err != nil {
		panic(err)
	}
	window := gocv.NewWindow("WebCam")
	img := gocv.NewMat()
	for {
		success := webcam.Read(&img)
		if !success {
			fmt.Println("Device closed")
			break
		}

		if window.WaitKey(1) != -1 {
			break
		}
		AddFpsOnFrame(&img, &fps_data)
		window.IMShow(img)
	}
	webcam.Close()
}
