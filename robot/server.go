package main

import (
	"context"
	"fmt"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
	pb "github.com/takabayap/sasisuseso-tai/protocol"
)

type grpcServer struct {
	pb.UnimplementedRobotServer
	pins map[int]rpio.Pin
}

func NewGRPCServer() *grpcServer {
	rpio.Open()
	pins := make(map[int]rpio.Pin)
	for _, pin := range []int{4, 17, 9, 11} {
		pins[pin] = rpio.Pin(pin)
		pins[pin].Output()
	}
	return &grpcServer{
		UnimplementedRobotServer: pb.UnimplementedRobotServer{},
		pins:                     pins,
	}
}

func (g *grpcServer) Turn(ctx context.Context, req *pb.TurnRequest) (*pb.TurnResponse, error) {
	fmt.Println("turn", req)
	g.pins[4].High()
	g.pins[11].High()
	time.Sleep(time.Duration(time.Second))
	g.pins[4].Low()
	g.pins[11].Low()
	return &pb.TurnResponse{
		Success: true,
	}, nil
}

func (g *grpcServer) Move(ctx context.Context, req *pb.MoveRequest) (*pb.MoveResponse, error) {
	fmt.Println("move", req)
	g.pins[4].High()
	g.pins[17].High()
	time.Sleep(time.Duration(time.Second))
	g.pins[4].Low()
	g.pins[17].Low()
	return &pb.MoveResponse{
		Success: true,
	}, nil
}
