package main

import (
	"net"
	"os"
	"os/signal"

	pb "github.com/takabayap/sasisuseso-tai/protocol"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	pb.RegisterRobotServer(s, NewGRPCServer())
	go func() {
		s.Serve(lis)
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	s.GracefulStop()
}
