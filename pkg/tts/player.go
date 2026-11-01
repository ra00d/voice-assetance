package tts

import (
	"bytes"
	"time"

	"github.com/hajimehoshi/oto/v2"
)

// Player handles audio playback.
type Player struct {
	ctx *oto.Context
}

// NewPlayer creates a new audio player.
func NewPlayer() (*Player, error) {
	// Piper default is usually 22050Hz or 16000Hz depending on model.
	// We'll assume 22050Hz, 1 channel, 16-bit (2 bytes) for now.
	// This should be configurable.
	const sampleRate = 22050
	const channelCount = 1
	const bitDepthInBytes = 2

	ctx, ready, err := oto.NewContext(sampleRate, channelCount, bitDepthInBytes)
	if err != nil {
		return nil, err
	}

	<-ready // Wait for audio context to be ready
	return &Player{ctx: ctx}, nil
}

// Play plays the given audio data.
func (p *Player) Play(data []byte) error {
	player := p.ctx.NewPlayer(bytes.NewReader(data))
	player.Play()

	// Wait for playback to finish
	// Oto doesn't have a direct "wait" so we sleep or check state.
	// A simple way is to calculate duration.
	// Duration = len(data) / (sampleRate * channels * bytesPerSample)
	
	// For now, we just let it play in background or block?
	// If we want to block until finished:
	for player.IsPlaying() {
		time.Sleep(10 * time.Millisecond)
	}
	
	return nil
}
