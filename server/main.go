package main

import (
	"context"
	"fmt"
	"image"
	"os"
	"os/signal"

	"github.com/brentnd/go-snowboy"
	"github.com/takabayap/sasisuseso-tai/server/components"
	"gocv.io/x/gocv"
	"google.golang.org/grpc"

	pb "github.com/takabayap/sasisuseso-tai/protocol"
)

func main() {
	webcam, err := gocv.OpenVideoCapture(0)
	if err != nil {
		panic(err)
	}

	webcamImg := gocv.NewMat()
	windowImg := gocv.NewMat()
	window := gocv.NewWindow("Detection")
	field := components.NewField(1000, 1000, [4]int{0, 1, 2, 3}, gocv.ArucoDict4x4_50)
	// should run this as a separate goroutine in the future
	var webcamOK bool

	go func() {
		for {
			webcamOK = webcam.Read(&webcamImg)
			window.WaitKey(1)
			// some webcam has a problem reading the frame sometimes
			if !webcamOK || webcamImg.Empty() {
				continue
			}
			fieldImg, err := field.GetFieldImg(webcamImg)
			if err != nil {
				// if the field marker is not detected
				gocv.Resize(webcamImg, &windowImg, image.Point{1920, 1080}, 0, 0, gocv.InterpolationLinear)
				window.IMShow(windowImg)
			} else {
				window.IMShow(fieldImg)
			}
		}
	}()

	mic := &Sound{}
	err = mic.Init()
	if err != nil {
		panic(err)
	}
	defer mic.Close()

	snowboyDetector := snowboy.NewDetector("common.res")
	defer snowboyDetector.Close()

	snowboyDetector.HandleFunc(snowboy.NewHotword("snowboy.umdl", 0.5), func(string) {
		fmt.Println("detected!")
	})
	sr, nc, bd := snowboyDetector.AudioFormat()
	fmt.Printf("sample rate=%d, num channels=%d, bit depth=%d\n", sr, nc, bd)
	go snowboyDetector.ReadAndDetect(mic)

	conn, err := grpc.Dial("localhost:9000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := pb.NewRobotClient(conn)
	message := &pb.MoveRequest{
		Forward:  true,
		Distance: 100,
	}
	res, err := client.Move(context.Background(), message)
	if err != nil {
		panic(err)
	}
	fmt.Println("res", res.Success)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
