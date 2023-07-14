package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
	pb "github.com/takabayap/sasisuseso-tai/protocol"
)

type grpcServer struct {
	pb.UnimplementedRobotServer
	pins       map[int]rpio.Pin
	isMockMode bool
	mutex      sync.Mutex
}

func NewGRPCServer() *grpcServer {
	if err := rpio.Open(); err != nil {
		fmt.Println("error opening rpio, running server as mock")
		return &grpcServer{
			UnimplementedRobotServer: pb.UnimplementedRobotServer{},
			isMockMode:               true,
		}
	}
	pins := make(map[int]rpio.Pin)
	for _, pin := range []int{4, 17, 9, 11} {
		pins[pin] = rpio.Pin(pin)
		pins[pin].Output()
	}
	return &grpcServer{
		UnimplementedRobotServer: pb.UnimplementedRobotServer{},
		pins:                     pins,
		isMockMode:               false,
		mutex:                    sync.Mutex{},
	}
}

func (g *grpcServer) Turn(ctx context.Context, req *pb.TurnRequest) (*pb.TurnResponse, error) {
	g.mutex.Lock()
	fmt.Println("turn", req)
	if g.isMockMode {
		time.Sleep(time.Second)
		return &pb.TurnResponse{
			Success: true,
		}, nil
	}
	g.pins[4].High()
	g.pins[11].High()
	time.Sleep(time.Duration(time.Second))
	g.pins[4].Low()
	g.pins[11].Low()
	g.mutex.Unlock()
	return &pb.TurnResponse{
		Success: true,
	}, nil
}

func (g *grpcServer) Move(ctx context.Context, req *pb.MoveRequest) (*pb.MoveResponse, error) {
	g.mutex.Lock()
	fmt.Println("move", req)
	if g.isMockMode {
		time.Sleep(time.Second)
		return &pb.MoveResponse{
			Success: true,
		}, nil
	}
	if req.Forward {
		g.pins[4].High()
		g.pins[9].High()
	} else {
		g.pins[11].High()
		g.pins[17].High()
	}
	time.Sleep(time.Duration(time.Millisecond * time.Duration(req.Distance)))
	if req.Forward {
		g.pins[4].Low()
		g.pins[9].Low()
	} else {
		g.pins[11].Low()
		g.pins[17].Low()
	}
	g.mutex.Unlock()
	return &pb.MoveResponse{
		Success: true,
	}, nil
}

func (g *grpcServer) SetPin(ctx context.Context, req *pb.SetPinRequest) (*pb.SetPinResponse, error) {
	fmt.Println("set pin", req)
	if g.isMockMode {
		time.Sleep(time.Second)
		return &pb.SetPinResponse{
			Success: true,
		}, nil
	}
	for key, v := range req.Pins {
		if v {
			g.pins[int(key)].High()
		} else {
			g.pins[int(key)].Low()
		}
	}
	return &pb.SetPinResponse{
		Success: true,
	}, nil
}
