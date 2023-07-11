package main

import (
	"context"
	"fmt"

	pb "github.com/takabayap/sasisuseso-tai/protocol"
)

type grpcServer struct {
	pb.UnimplementedRobotServer
}

func NewGRPCServer() *grpcServer {
	return &grpcServer{}
}

func (g *grpcServer) Turn(ctx context.Context, req *pb.TurnRequest) (*pb.TurnResponse, error) {
	fmt.Println("turn", req)
	return &pb.TurnResponse{
		Success: false,
	}, nil
}

func (g *grpcServer) Move(ctx context.Context, req *pb.MoveRequest) (*pb.MoveResponse, error) {
	fmt.Println("move", req)
	return &pb.MoveResponse{
		Success: false,
	}, nil
}
