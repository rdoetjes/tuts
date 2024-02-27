# GO-CV 
Is the scaffolding to quickly build a cross-platform Go application with OpenCV, reading from the web cam by default

## MacOS
The opencv is installed from home brew, incase you also have gcc/g++ installed with home brew you need to force the longer to clang
```
export CXX=/usr/bin/clang++ export CC=/usr/bin/clang
```
Also on an Mac you need to set device for the camera to 1 (is now the default), for other systems this is most likely 0

## Linux Windows
Be sure to change the camera id in main.go to 0
```
	webcam, err := gocv.VideoCaptureDevice(1)
```

Should become

```
	webcam, err := gocv.VideoCaptureDevice(0)
```