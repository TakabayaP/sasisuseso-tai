// Components package represents components used in sasisuseso-tai.
package components

import (
	"fmt"
	"image"

	"gocv.io/x/gocv"
)

// Field is a struct that represents the main field.
// Field should have four markers with the given IDs.
type Field struct {
	Width         float64
	Height        float64
	CornerIDs     [4]int
	arucoDetector *gocv.ArucoDetector
}

// NewField returns a new Field struct.
func NewField(width, height float64, cornerIDs [4]int) *Field {
	a := gocv.NewArucoDetectorWithParams(gocv.GetPredefinedDictionary(gocv.ArucoDict4x4_100), gocv.NewArucoDetectorParameters())
	return &Field{
		Width: width, Height: height,
		CornerIDs:     cornerIDs,
		arucoDetector: &a,
	}
}

// GetFieldImg returns the field image from the given image.
// The given image should contain four markers with the given IDs.
// If the markers are not enough, it returns an error.
// Even if the markers are duplicated, it will not return an error.
func (f *Field) GetFieldImg(img gocv.Mat) (gocv.Mat, error) {
	markers := MarkersFromDetectMarkers(f.arucoDetector.DetectMarkers(img))
	var fieldMarkers [4]Marker
	if len(markers) < 4 {
		return gocv.Mat{}, fmt.Errorf("not enough markers")
	}
	for i, id := range f.CornerIDs {
		for j, m := range markers {
			if m.ID == id {
				fieldMarkers[i] = m
				break
			}
			if j == len(markers)-1 {
				return gocv.Mat{}, fmt.Errorf("not enough markers")
			}
		}
	}
	points := getInnerCorners(fieldMarkers)
	fmt.Println(points)
	markerPoints := gocv.NewPointVectorFromPoints(points[:])
	imagePoints := gocv.NewPointVectorFromPoints([]image.Point{
		{X: 0, Y: 0},
		{X: int(f.Width), Y: 0},
		{X: int(f.Width), Y: int(f.Height)},
		{X: 0, Y: int(f.Height)},
	})
	fmt.Println(points[:])
	trans_mat := gocv.GetPerspectiveTransform(
		markerPoints,
		imagePoints,
	)
	resultImg := gocv.NewMat()
	gocv.WarpPerspective(img, &resultImg, trans_mat, image.Point{X: int(f.Width), Y: int(f.Height)})
	return resultImg, nil
}
