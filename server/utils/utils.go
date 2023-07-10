package utils

import (
	"image"

	"gocv.io/x/gocv"
)

func GoCVPoint2ftoImagePoint(p gocv.Point2f) image.Point {
	return image.Point{X: int(p.X), Y: int(p.Y)}
}

func DistanceSquare(p1, p2 image.Point) float32 {
	return float32((p1.X-p2.X)*(p1.X-p2.X) + (p1.Y-p2.Y)*(p1.Y-p2.Y))
}
