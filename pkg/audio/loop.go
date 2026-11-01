package audio

import (
	"fmt"
)

// Loop manages the audio capture and VAD loop.
type Loop struct {
	recorder      *Recorder
	vad           *VAD
	Output        chan []float32 // Channel to send detected speech segments
	SpeechStarted chan bool      // Channel to signal start of speech
}

// NewLoop creates a new audio loop.
func NewLoop(recorder *Recorder, vad *VAD) *Loop {
	return &Loop{
		recorder:      recorder,
		vad:           vad,
		Output:        make(chan []float32, 10),
		SpeechStarted: make(chan bool, 10),
	}
}

// Start begins the audio loop.
func (l *Loop) Start() {
	go l.run()
}

func (l *Loop) run() {
	var speechSegment []float32
	isSpeaking := false
	silenceCounter := 0
	const silenceThreshold = 10 // Number of frames of silence to consider speech ended (approx 0.5s)
	const speechThreshold = 0.5 // VAD probability threshold

	fmt.Println("Audio loop started. Listening...", l.recorder.Stream)

	for data := range l.recorder.Stream {
		// Convert to float32 for VAD
		floats := Uint8ToFloat32(data)

		// Process in chunks of 512 (32ms at 16kHz)
		// Note: malgo might return different chunk sizes, so we need to buffer/chunk properly.
		// For simplicity, we assume malgo is configured or we handle it.
		// Here we just process what we get if it matches VAD requirements, or skip/buffer.
		// Silero VAD supports 512, 1024, 1536.

		// Simple buffering logic (omitted for brevity, assuming 512 or similar from recorder)
		// In a real app, we'd use a ring buffer.
		// Let's assume we get 512 samples for now.

		// If we get more, we split.
		chunkSize := 512
		for i := 0; i < len(floats); i += chunkSize {
			end := i + chunkSize
			if end > len(floats) {
				break // Drop incomplete chunk
			}
			chunk := floats[i:end]

			prob, err := l.vad.Process(chunk)
			chunkInt16 := Float32ToInt16(prob)
			fmt.Println("VAD chunkInt16:", chunkInt16)

			if err != nil {
				fmt.Printf("loop VAD error: %v\n", err)
				continue
			}
			if prob > speechThreshold {
				if !isSpeaking {
					fmt.Println("Speech detected...")
					isSpeaking = true
					// Signal speech start non-blocking
					select {
					case l.SpeechStarted <- true:
					default:
					}
				}
				silenceCounter = 0
				speechSegment = append(speechSegment, chunk...)
			} else {
				if isSpeaking {
					silenceCounter++
					speechSegment = append(speechSegment, chunk...) // Keep trailing silence for naturalness

					if silenceCounter > silenceThreshold {
						fmt.Println("Speech ended.")
						isSpeaking = false
						// Send segment
						l.Output <- speechSegment
						speechSegment = nil // Reset
						silenceCounter = 0
					}
				}
			}
		}
	}
}
