package stt

import (
	"fmt"
	"os"

	"github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"
)

// Whisper handles speech-to-text using Whisper.cpp.
type Whisper struct {
	model   whisper.Model
	context whisper.Context
}

// NewWhisper creates a new Whisper instance.
func NewWhisper(modelPath string) (*Whisper, error) {
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("model file not found: %s", modelPath)
	}

	model, err := whisper.New(modelPath)
	if err != nil {
		return nil, err
	}

	context, err := model.NewContext()
	if err != nil {
		return nil, err
	}

	return &Whisper{
		model:   model,
		context: context,
	}, nil
}

// Transcribe converts audio samples to text.
func (w *Whisper) Transcribe(samples []float32) (string, error) {
	if err := w.context.Process(samples, nil, nil,nil); err != nil {
		return "", err
	}

	// Iterate through segments
	var text string
	for {
		segment, err := w.context.NextSegment()
		if err != nil {
			break
		}
		text += segment.Text
	}

	return text, nil
}

// Close releases resources.
func (w *Whisper) Close() {
	// w.context.
	w.model.Close()
}
