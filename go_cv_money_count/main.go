// main function
package main

import (
	"fmt"
	"go/cv/chipcount"
	"go/cv/cvhelper"
	"image"
	"image/color"
	"os"
	"strconv"

	// Import gocv package for OpenCV wrappers and bindings
	"gocv.io/x/gocv"
)

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
	//fps tracking struct
	fps_data := cvhelper.NewFps()
	//config for the coincount preprocessing
	config := chipcount.NewDefaultCoinProcessing()

	//webcam setup
	var webcam_type string = "v4l2"
	var idx int64 = 0

	if len(os.Args) > 1 {
		idx, _ = strconv.ParseInt(os.Args[1], 10, 0)
	}

	if len(os.Args) > 2 {
		webcam_type = os.Args[2]
	}

	webcam := cvhelper.SetupWebcam(int(idx), webcam_type)
	defer webcam.Close()
	webcam.Set(gocv.VideoCaptureFrameWidth, 480)
	webcam.Set(gocv.VideoCaptureFrameHeight, 640)

	//window setup to display the output
	input_w, process_w := setupWindows()
	defer input_w.Close()
	defer process_w.Close()

	//image setup to store the input and output and avoid dynamic allocation in the loop to speed things up
	img, process := setupImages()
	defer img.Close()
	defer process.Close()
	defer process.Close()
	for {
		var errCount = 0
		success := webcam.Read(img)
		for !success {
			success = webcam.Read(img)
			//// Check if the frame is read correctly
			fmt.Println("Device closed")
			errCount += 1
			if errCount == 15 {
				os.Exit(1)
			}
		}
		//*img = gocv.IMRead("./euros.jpg", gocv.IMReadColor)

		// if any key is pressed then exit
		if input_w.WaitKey(1) != -1 || process_w.WaitKey(1) != -1 {
			break
		}

		// add your image processing functions below
		totalAmount := chipcount.CountChips(img, process, config)
		gocv.PutText(img, totalAmount, image.Pt(10, 40), gocv.FontHersheyPlain, 3.2, color.RGBA{255, 0, 255, 1}, 2)
		fmt.Println(totalAmount)

		//add frame count to the upperleft of the frame, the stateful data is held in fps_data and is updated by the function
		cvhelper.AddFpsOnFrame(process, fps_data)

		//display the frame from the webcam wit the fps on it
		process_w.IMShow(*process)
		input_w.IMShow(*img)
	}
}
