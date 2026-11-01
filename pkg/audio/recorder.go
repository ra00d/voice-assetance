package audio

import (
	"fmt"
	"sync"

	"github.com/gen2brain/malgo"
)

// Recorder handles audio capture from the microphone.
type Recorder struct {
	ctx      *malgo.AllocatedContext
	device   *malgo.Device
	stopChan chan struct{}
	wg       sync.WaitGroup
	Stream   chan []byte // Channel to send raw audio data
}

// NewRecorder creates a new Recorder instance.
func NewRecorder() (*Recorder, error) {
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, func(message string) {
		fmt.Printf("MALGO LOG: %v\n", message)
	})
	if err != nil {
		fmt.Println("Failed to init context:", err)
		return nil, err
	}

	return &Recorder{
		ctx:      ctx,
		stopChan: make(chan struct{}),
		Stream:   make(chan []byte, 1024), // Buffer audio chunks
	}, nil
}

// Start begins capturing audio.
func (r *Recorder) Start() error {
	fmt.Println("Starting recorder...")
	deviceConfig := malgo.DefaultDeviceConfig(malgo.Capture)
	deviceConfig.Capture.Format = malgo.FormatS16
	deviceConfig.Capture.Channels = 1
	deviceConfig.SampleRate = 16000
	deviceConfig.Alsa.NoMMap = 1

	onRecvFrames := func(pOutputSample, pInputSamples []byte, framecount uint32) {
		// Send a copy of the data to avoid race conditions if the buffer is reused
		data := make([]byte, len(pInputSamples))
		copy(data, pInputSamples)

		select {
		case r.Stream <- data:
		default:
			// Drop frame if channel is full to avoid blocking the callback
			fmt.Println("Warning: Audio buffer full, dropping frame")
		}
	}

	deviceCallbacks := malgo.DeviceCallbacks{
		Data: onRecvFrames,
	}

	device, err := malgo.InitDevice(r.ctx.Context, deviceConfig, deviceCallbacks)
	if err != nil {
		fmt.Println("Failed to init device:", err)
		return err
	}
	r.device = device

	if err := device.Start(); err != nil {
		fmt.Println("Failed to start device:", err)
		return err
	}

	return nil
}

// Stop stops the audio capture and cleans up resources.
func (r *Recorder) Stop() {
	fmt.Println("Stopping recorder...")
	if r.device != nil {
		r.device.Uninit()
	}
	if r.ctx != nil {
		r.ctx.Free()
	}
	close(r.Stream)
}
