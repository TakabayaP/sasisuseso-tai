package controllers

import (
	"image"
	"sync"
	"time"

	"github.com/takabayap/sasisuseso-tai/server/components"
	"github.com/takabayap/sasisuseso-tai/server/constants"
	"gocv.io/x/gocv"
)

type FieldController struct {
	field         *components.Field
	robot         *components.Robot
	webcam        *gocv.VideoCapture
	window        *gocv.Window
	arucoDetector *gocv.ArucoDetector
	fieldImg      gocv.Mat
	fieldImgMutex sync.RWMutex
	markers       map[int]components.Marker
}

func NewFieldController(webcamID int) *FieldController {
	webcam, err := gocv.OpenVideoCapture(webcamID)
	if err != nil {
		panic(err)
	}
	window := gocv.NewWindow("Detection")
	field := components.NewField(1000, 1000, [4]int{0, 1, 2, 3})
	arucoDetector := gocv.NewArucoDetectorWithParams(gocv.GetPredefinedDictionary(gocv.ArucoDict4x4_100), gocv.NewArucoDetectorParameters())

	f := &FieldController{
		field:         field,
		webcam:        webcam,
		window:        window,
		arucoDetector: &arucoDetector,
		fieldImgMutex: sync.RWMutex{},
	}
	f.robot = components.NewRobot(10, 11, constants.RobotAddr, constants.MarkerDistFromCenter, f.getMarkers)
	return f
}

func (f *FieldController) getMarkers() map[int]components.Marker {
	f.fieldImgMutex.RLock()
	markers := components.MarkersFromDetectMarkers(f.arucoDetector.DetectMarkers(f.fieldImg))
	f.fieldImgMutex.RUnlock()
	for _, m := range markers {
		f.markers[m.ID] = m
	}
	return f.markers
}

func (f *FieldController) WatchField() {
	var webcamOK bool
	webcamImg := gocv.NewMat()
	windowImg := gocv.NewMat()
	for {
		webcamOK = f.webcam.Read(&webcamImg)
		f.window.WaitKey(1)
		// some webcam have problems reading the frame sometimes
		if !webcamOK || webcamImg.Empty() {
			continue
		}
		fieldImg, err := f.field.GetFieldImg(webcamImg)
		if err != nil {
			// if the field marker is not detected
			gocv.Resize(webcamImg, &windowImg, image.Point{1920, 1080}, 0, 0, gocv.InterpolationLinear)
			f.window.IMShow(windowImg)
			continue
		}
		f.window.IMShow(fieldImg)
		f.fieldImgMutex.Lock()
		f.fieldImg = fieldImg
		f.fieldImgMutex.Unlock()
	}
}

func (f *FieldController) CallRobot() {
	f.robot.Move(true, 300)
	f.robot.Turn(90)
	f.robot.Move(true, 50)
	time.Sleep(3 * time.Second)
	f.robot.Move(false, 50)
	f.robot.Turn(90)
	f.robot.Move(false, 300)
}
