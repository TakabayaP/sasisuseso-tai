package components

import (
	"image"
	"math"
	"sort"

	"github.com/takabayap/sasisuseso-tai/server/utils"
	"gocv.io/x/gocv"
)

type Marker struct {
	ID      int
	Corners []gocv.Point2f
}

func (m *Marker) Center() (float32, float32) {
	var x, y float32
	for _, c := range m.Corners {
		x += c.X
		y += c.Y
	}
	return x / 4, y / 4
}

func (m *Marker) Angle() float32 {
	return float32(math.Atan2(float64(m.Corners[1].Y-m.Corners[0].Y), float64(m.Corners[1].X-m.Corners[0].X)))
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

func MarkersFromDetectMarkers(fieldCorners [][]gocv.Point2f, ids []int, rejectedCandidates [][]gocv.Point2f) []Marker {
	var markers []Marker
	for i, c := range fieldCorners {
		markers = append(markers, Marker{ID: ids[i], Corners: c})
	}
	return markers
}
