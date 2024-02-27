package main

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

// main function
func main() {
	fps_data := NewFps()

	webcam, err := gocv.VideoCaptureDevice(1)
	if err != nil {
		panic(err)
	}
	defer webcam.Close()

	window := gocv.NewWindow("WebCam")
	img := gocv.NewMat()
	for {
		success := webcam.Read(&img)
		// Check if the frame is read correctly
		if !success {
			fmt.Println("Device closed")
			break
		}

		// if any key is pressed then exit
		if window.WaitKey(1) != -1 {
			break
		}

		// add your image processing functions below
		//...

		//add frame count to the upperleft of the frame, the stateful data is held in fps_data and is updated by the function
		AddFpsOnFrame(&img, fps_data)

		//display the frame from the webcam wit the fps on it
		window.IMShow(img)
	}
}
