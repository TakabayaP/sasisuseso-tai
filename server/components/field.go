// Components package represents components used in sasisuseso-tai.
package components

import (
	"fmt"
	"image"
	"math"
	"sort"

	"github.com/takabayap/sasisuseso-tai/server/utils"
	"gocv.io/x/gocv"
)

// Field is a struct that represents the main field.
// Field should have four markers with the given IDs.
type Field struct {
	Width         float64
	Height        float64
	CornerIDs     [4]int
	Dict          gocv.ArucoDictionaryCode
	arucoDetector *gocv.ArucoDetector
}

type Marker struct {
	ID      int
	Corners []gocv.Point2f
}

func getInnerCorners(markers [4]Marker) [4]image.Point {
	var x, y float32
	for _, p := range markers {
		x += (p.Corners[0].X + p.Corners[1].X + p.Corners[2].X + p.Corners[3].X) / 4
		y += (p.Corners[0].Y + p.Corners[1].Y + p.Corners[2].Y + p.Corners[3].Y) / 4
	}
	center := image.Point{X: int(x / 4), Y: int(y / 4)}
	var corners []image.Point

	for _, p := range markers {
		var minDist float32 = math.MaxFloat32
		var minDistCorner image.Point
		for _, c := range p.Corners {
			dist := utils.DistanceSquare(center, utils.GoCVPoint2ftoImagePoint(c))
			if dist < minDist {
				minDist = dist
				minDistCorner = utils.GoCVPoint2ftoImagePoint(c)
			}
		}
		corners = append(corners, minDistCorner)
	}
	sort.Slice(corners, func(i, j int) bool {
		return (corners[j].X-center.X)*(corners[i].Y-center.Y) < (corners[i].X-center.X)*(corners[j].Y-center.Y)
	})
	result := [4]image.Point{corners[0], corners[1], corners[2], corners[3]}
	return result
}

func newMarkersFromDetectMarkers(corners [][]gocv.Point2f, ids []int, rejectedCandidates [][]gocv.Point2f) []Marker {
	var markers []Marker
	for i, c := range corners {
		markers = append(markers, Marker{ID: ids[i], Corners: c})
	}
	return markers
}

// NewField returns a new Field struct.
func NewField(width, height float64, cornerIDs [4]int, dict gocv.ArucoDictionaryCode) *Field {
	a := gocv.NewArucoDetectorWithParams(gocv.GetPredefinedDictionary(dict), gocv.NewArucoDetectorParameters())
	return &Field{
		Width: width, Height: height,
		CornerIDs:     cornerIDs,
		Dict:          dict,
		arucoDetector: &a,
	}
}

// GetFieldImg returns the field image from the given image.
// The given image should contain four markers with the given IDs.
// If the markers are not enough, it returns an error.
// Even if the markers are duplicated, it will not return an error.
func (f *Field) GetFieldImg(img gocv.Mat) (gocv.Mat, error) {
	markers := newMarkersFromDetectMarkers(f.arucoDetector.DetectMarkers(img))
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
