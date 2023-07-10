package main

import (
	"image"

	"github.com/takabayap/sasisuseso-tai/server/components"
	"gocv.io/x/gocv"
)

func main() {
	webcam, err := gocv.OpenVideoCapture(0)
	if err != nil {
		panic(err)
	}

	webcamImg := gocv.NewMat()
	windowImg := gocv.NewMat()
	window := gocv.NewWindow("Detection")
	field := components.NewField(1000, 1000, [4]int{0, 1, 2, 3}, gocv.ArucoDict4x4_50)
	// should run this as a separate goroutine in the future
	var webcamOK bool
	for {
		webcamOK = webcam.Read(&webcamImg)
		window.WaitKey(1)
		// some webcam has a problem reading the frame sometimes
		if !webcamOK || webcamImg.Empty() {
			continue
		}
		fieldImg, err := field.GetFieldImg(webcamImg)
		if err != nil {
			// if the field marker is not detected
			gocv.Resize(webcamImg, &windowImg, image.Point{1920, 1080}, 0, 0, gocv.InterpolationLinear)
			window.IMShow(windowImg)
		} else {
			window.IMShow(fieldImg)
		}
	}
}
