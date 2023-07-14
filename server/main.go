package main

import (
	"os"
	"os/signal"

	"github.com/takabayap/sasisuseso-tai/server/constants"
	"github.com/takabayap/sasisuseso-tai/server/controllers"
)

func main() {
	fc := controllers.NewFieldController(constants.WebcamID)
	go fc.WatchField()

	hd := controllers.NewHotwordDetector()
	hd.HandleFunc(controllers.Snowboy, func(s string) {
		fc.CallRobot()
	})
	go hd.ReadAndDetect()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
