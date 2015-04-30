package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cryptix/wav"

	"code.google.com/p/portaudio-go/portaudio"
)

var inputDevice, outputDevice *portaudio.DeviceInfo
var writer *wav.Writer

func main() {

	err := portaudio.Initialize()
	if err != nil {
		fmt.Println("Init", err)
		return
	}
	outputDevice, err := portaudio.DefaultOutputDevice()
	outputParams := portaudio.StreamDeviceParameters{
		Device:   outputDevice,
		Channels: 2,
	}

	inputDevice, err := portaudio.DefaultInputDevice()
	inputParams := portaudio.StreamDeviceParameters{
		Device:   inputDevice,
		Channels: 2,
	}

	fmt.Printf("%+v\n", outputDevice)
	fmt.Printf("%+v\n", inputDevice)

	p := portaudio.StreamParameters{
		Output:          outputParams,
		Input:           inputParams,
		SampleRate:      44100,
		FramesPerBuffer: 256,
		Flags:           portaudio.NoFlag,
	}

	f, _ := os.Create("output.wav")

	meta := wav.File{
		Channels:        1,
		SampleRate:      44100,
		SignificantBits: 32,
	}

	writer, err = meta.NewWriter(f)
	defer f.Close()
	defer writer.Close()

	buf := make([]int32, 10000)

	s, err := portaudio.OpenStream(p, recorder, &buf)
	if err != nil {
		fmt.Println("OpenStream", err)
		fmt.Println(p)
		return
	}
	log.Println(s.Info())
	s.Start()

	time.Sleep(5 * time.Second)

	s.Stop()

}

func recorder(in []int32, out []int32, timeInfo portaudio.StreamCallbackTimeInfo, flags portaudio.StreamCallbackFlags) {
	fmt.Printf("Recording... %d %d \n", len(in), len(out))
	for i := 0; i < len(out); i++ {
		_ = writer.WriteInt32(out[i])
	}
}
