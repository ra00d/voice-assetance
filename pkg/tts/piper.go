package tts

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

// Piper handles text-to-speech using the Piper binary.
// Note: We use the binary via os/exec because pure Go bindings are complex/rare.
type Piper struct {
	modelPath string
	binaryPath string
}

// NewPiper creates a new Piper instance.
func NewPiper(modelPath string, binaryPath string) *Piper {
	return &Piper{
		modelPath:  modelPath,
		binaryPath: binaryPath,
	}
}

// Synthesize generates audio from text and returns the WAV data.
func (p *Piper) Synthesize(text string) ([]byte, error) {
	// echo "text" | piper --model model.onnx --output_raw
	
	cmd := exec.Command(p.binaryPath, "--model", p.modelPath, "--output_raw")
	cmd.Stdin = bytes.NewBufferString(text)
	
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr // Redirect stderr for debugging

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("piper execution failed: %w", err)
	}

	return out.Bytes(), nil
}
