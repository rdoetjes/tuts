// main function
package main

import (
	"fmt"
	"go/cv/coincount"
	"go/cv/cvhelper"

	// Import gocv package for OpenCV wrappers and bindings
	"gocv.io/x/gocv"
)

func setupWebcam() *gocv.VideoCapture {
	// Open webcam
	webcam, err := gocv.VideoCaptureDevice(1)
	if err != nil {
		panic(err)
	}
	webcam.Set(gocv.VideoCaptureFrameWidth, 480)
	webcam.Set(gocv.VideoCaptureFrameHeight, 640)
	return webcam
}

func setupWindows() (*gocv.Window, *gocv.Window) {
	input_w := gocv.NewWindow("WebCam")
	process_w := gocv.NewWindow("Process")
	return input_w, process_w
}

func setupImages() (*gocv.Mat, *gocv.Mat) {
	// Create windows to display the output
	//prealloc image sizes
	img := gocv.NewMat()
	process := gocv.NewMat()
	return &img, &process
}

func main() {
	fps_data := cvhelper.NewFps()
	webcam := setupWebcam()
	defer webcam.Close()

	input_w, process_w := setupWindows()
	defer input_w.Close()
	defer process_w.Close()

	img, process := setupImages()
	defer img.Close()
	defer process.Close()

	config := coincount.NewDefaultCoinProcessing()
	for {
		success := webcam.Read(img)
		// Check if the frame is read correctly
		if !success {
			fmt.Println("Device closed")
			break
		}

		// if any key is pressed then exit
		if input_w.WaitKey(1) != -1 || process_w.WaitKey(1) != -1 {
			break
		}

		// add your image processing functions below
		coincount.CountEuros(img, process, config)

		//add frame count to the upperleft of the frame, the stateful data is held in fps_data and is updated by the function
		cvhelper.AddFpsOnFrame(process, fps_data)

		//display the frame from the webcam wit the fps on it
		process_w.IMShow(*process)
		input_w.IMShow(*img)
	}
}
