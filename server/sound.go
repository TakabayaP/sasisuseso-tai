package main

import (
	"bytes"
	"encoding/binary"

	"github.com/gordonklaus/portaudio"
)

type Sound struct {
	stream *portaudio.Stream
	data   []int16
}

func (s *Sound) Init() error {
	inputChannels := 1
	outputChannels := 0
	sampleRate := 16000
	s.data = make([]int16, 1024)

	err := portaudio.Initialize()
	if err != nil {
		return err
	}

	stream, err := portaudio.OpenDefaultStream(inputChannels, outputChannels, float64(sampleRate), len(s.data), s.data)
	if err != nil {
		return err
	}

	err = stream.Start()
	if err != nil {
		return err
	}

	s.stream = stream
	return nil
}

func (s *Sound) Close() error {
	err := s.stream.Close()
	if err != nil {
		return err
	}
	portaudio.Terminate()
	if err != nil {
		return err
	}
	return nil
}

func (s *Sound) Read(p []byte) (int, error) {
	s.stream.Read()

	buf := &bytes.Buffer{}
	for _, v := range s.data {
		binary.Write(buf, binary.LittleEndian, v)
	}

	copy(p, buf.Bytes())
	return len(p), nil
}
