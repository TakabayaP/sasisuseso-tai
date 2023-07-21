package components

import (
	"context"
	"fmt"
	"math"

	"google.golang.org/grpc"

	pb "github.com/takabayap/sasisuseso-tai/protocol"
)

type Robot struct {
	angle      float32
	ids        [2]int
	dist       float32
	getMarkers func() map[int]Marker
	conn       *grpc.ClientConn
	client     *pb.RobotClient
}

// NewRobot returns a new robot struct.
// Markers should have different IDs, and the first marker should be the left one.
func NewRobot(i1, i2 int, addr string, markerDistFromCenter float32, getMarkers func() map[int]Marker) *Robot {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client := pb.NewRobotClient(conn)
	return &Robot{
		angle:      0,
		ids:        [2]int{i1, i2},
		dist:       markerDistFromCenter,
		getMarkers: getMarkers,
		conn:       conn,
		client:     &client,
	}
}

func (r *Robot) Angle() float32 {
	fmt.Println("not implemented")
	return r.angle
}

func (r *Robot) Position() (float32, float32, error) {
	markers := r.getMarkers()
	var left, right Marker
	left, okl := markers[r.ids[0]]
	right, okr := markers[r.ids[1]]

	switch {
	case okl && okr:
		rcx, rcy := left.Center()
		lcx, lcy := right.Center()
		x := (rcx + lcx) / 2
		y := (rcy + lcy) / 2
		return x, y, nil
	case okl:
		cx, cy := left.Center()
		x := cx + r.dist*float32(math.Cos(float64(left.Angle())))
		y := cy + r.dist*float32(math.Sin(float64(left.Angle())))
		return x, y, nil
	case okr:
		cx, cy := right.Center()
		x := cx + r.dist*float32(math.Cos(float64(right.Angle())))
		y := cy + r.dist*float32(math.Sin(float64(right.Angle())))
		return x, y, nil
	default:
		return 0, 0, fmt.Errorf("not enough markers")
	}
}

func (r *Robot) Move(forward bool, dist float32) error {
	message := &pb.MoveRequest{
		Forward:  forward,
		Distance: int32(dist),
	}
	res, err := (*r.client).Move(context.Background(), message)
	if !res.Success {
		return fmt.Errorf("failed to move")
	}
	return err
}

func (r *Robot) Turn(angle float32) error {
	message := &pb.TurnRequest{
		Angle: int32(angle),
	}
	res, err := (*r.client).Turn(context.Background(), message)
	if !res.Success {
		return fmt.Errorf("failed to turn")
	}
	return err
}

func (r *Robot) Close() {
	r.conn.Close()
}
