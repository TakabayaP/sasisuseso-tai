package controllers

import (
	"fmt"

	"github.com/brentnd/go-snowboy"
	"github.com/takabayap/sasisuseso-tai/server/components"
)

type hotword struct {
	Name      string
	modelFile string
}

const resourceFile = "../resources/common.res"

var (
	Snowboy = hotword{
		Name:      "snowboy",
		modelFile: "../resources/snowboy.umdl",
	}
)

type HotwordDetector struct {
	hotwords        []hotword
	snowboyDetector *snowboy.Detector
	mic             *components.Sound
}

func NewHotwordDetector() *HotwordDetector {
	snowboyDetector := snowboy.NewDetector(resourceFile)

	mic := &components.Sound{}
	err := mic.Init()
	if err != nil {
		panic(err)
	}
	return &HotwordDetector{
		snowboyDetector: &snowboyDetector,
		mic:             mic,
	}
}

func (h *HotwordDetector) ReadAndDetect() {
	h.snowboyDetector.ReadAndDetect(h.mic)
}

func (h *HotwordDetector) Close() {
	h.snowboyDetector.Close()
	h.mic.Close()
}

func (h *HotwordDetector) HandleFunc(hw hotword, f func(string)) {
	h.hotwords = append(h.hotwords, hw)
	h.snowboyDetector.HandleFunc(snowboy.NewHotword(hw.modelFile, 0.5), f)
	sr, nc, bd := h.snowboyDetector.AudioFormat()
	fmt.Printf("sample rate=%d, num channels=%d, bit depth=%d\n", sr, nc, bd)
}
