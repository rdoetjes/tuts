// main function
package main

import (
	"fmt"
	"go/cv/coincount"
	"go/cv/cvhelper"
	"image"
	"image/color"
	"os"

	// Import gocv package for OpenCV wrappers and bindings
	"gocv.io/x/gocv"
)

func setupWebcam(os string) *gocv.VideoCapture {
	// Open webcam this is fo linux (V4L2)
	if os == "linux" {
		webcam, err := gocv.VideoCaptureDeviceWithAPI(1, gocv.VideoCaptureV4L2)
		if err != nil {
			panic(err)
		}
		return webcam
	} else if os == "macos" {
		//which uses device 1 for some reason
		webcam, err := gocv.VideoCaptureDevice(1)
		if err != nil {
			panic(err)
		}
		return webcam
	} else {
		//the standard webcam
		webcam, err := gocv.VideoCaptureDevice(0)
		if err != nil {
			panic(err)
		}
		return webcam
	}
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
	//fps tracking struct
	fps_data := cvhelper.NewFps()
	//config for the coincount preprocessing
	config := coincount.NewDefaultCoinProcessing()

	//webcam setup
	var webcam_type string = "linux"
	if len(os.Args) > 1 {
		webcam_type = os.Args[1]
	}
	webcam := setupWebcam(webcam_type)
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

	for {
		var errCount = 0
		success := webcam.Read(img)
		for !success {
			success = webcam.Read(img)
			//// Check if the frame is read correctly
			fmt.Println("Device closed")
			errCount += 1
			if errCount == 15 {
				break
			}
		}
		//*img = gocv.IMRead("./euros.jpg", gocv.IMReadColor)

		//// if any key is pressed then exit
		if input_w.WaitKey(1) != -1 || process_w.WaitKey(1) != -1 {
			break
		}

		// add your image processing functions below
		totalAmount := coincount.CountEuros(img, process, config)
		fmt.Printf("%d\n", totalAmount)
		//add frame count to the upperleft of the frame, the stateful data is held in fps_data and is updated by the function
		cvhelper.AddFpsOnFrame(process, fps_data)
		gocv.PutText(img, fmt.Sprintf("Total Amount: %d", totalAmount), image.Pt(10, 20), gocv.FontHersheyPlain, 1.2, color.RGBA{255, 0, 255, 1}, 2)

		//display the frame from the webcam wit the fps on it
		process_w.IMShow(*process)
		input_w.IMShow(*img)
	}
}
