// main function
package main

import (
	"fmt"
	"go/cv/cvhelper"

	// Import gocv package for OpenCV wrappers and bindings
	"gocv.io/x/gocv"
)

func process_example(img *gocv.Mat) {
	gocv.CvtColor(*img, img, gocv.ColorBGRToGray)
	gocv.Flip(*img, img, -1)
}

func main() {
	fps_data := cvhelper.NewFps()

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
		process_example(&img)

		//add frame count to the upperleft of the frame, the stateful data is held in fps_data and is updated by the function
		cvhelper.AddFpsOnFrame(&img, fps_data)

		//display the frame from the webcam wit the fps on it
		window.IMShow(img)
	}
}
