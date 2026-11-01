package audio

import (
	"fmt"
	"os"

	"github.com/yalue/onnxruntime_go"
)

// VAD handles Voice Activity Detection using Silero VAD.
type VAD struct {
	session *onnxruntime_go.DynamicAdvancedSession
	state   []float32 // Hidden state
	sr      int64     // Sample rate
}

// NewVAD creates a new VAD instance.
func NewVAD(modelPath string) (*VAD, error) {
	// Initialize ONNX Runtime
	// Note: The user must ensure the shared library is available.
	// We assume onnxruntime_go.SetSharedLibraryPath has been called or is not needed if in standard path.
	// For simplicity, we'll try to load the shared library from a common location or skip if already loaded.
	// In a real app, we might want to make this configurable.

	// Check if the model file exists
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("model file not found: %s", modelPath)
	}

	onnxruntime_go.SetSharedLibraryPath("/mnt/main/ra0_0d/repositories/personal/voice-assetance/third_party/onnxruntime/lib/libonnxruntime.so") // Default assumption for Linux
	err := onnxruntime_go.InitializeEnvironment()
	if err != nil {
		// It might be already initialized
		fmt.Printf("ONNX Environment init warning (might be already init): %v\n", err)
	}

	inputNames := []string{"input", "state", "sr"}
	outputNames := []string{"output", "stateN"}

	session, err := onnxruntime_go.NewDynamicAdvancedSession(modelPath, inputNames, outputNames, nil)

	if err != nil {

		return nil, fmt.Errorf("failed to create session: %v", err)
	}

	// Initialize states (2 layers, batch 1, 128 dim)
	state := make([]float32, 2*1*128)

	return &VAD{
		session: session,
		state:   state,
		sr:      16000,
	}, nil
}

// Process processes a chunk of audio and returns the speech probability.
// Audio chunk should be 512, 1024, or 1536 samples for 16kHz.
func (v *VAD) Process(audio []float32) (float32, error) {
	// if len(audio) != 512 && len(audio) != 1024 && len(audio) != 1536 {
	// 	return 0, fmt.Errorf("invalid audio chunk size: %d", len(audio))
	// }

	// Prepare inputs
	inputTensor, err := onnxruntime_go.NewTensor(onnxruntime_go.NewShape(1, int64(len(audio))), audio)
	if err != nil {
		return 0, err
	}
	defer inputTensor.Destroy()

	// sr is a scalar int64, but NewTensor might need a dimension. Let's try [1].
	srTensor, err := onnxruntime_go.NewTensor(onnxruntime_go.NewShape(1), []int64{v.sr})
	if err != nil {
		fmt.Printf("invalid sr tensor: %v\n", []int64{v.sr})
		return 0, err
	}
	defer srTensor.Destroy()

	stateTensor, err := onnxruntime_go.NewTensor(onnxruntime_go.NewShape(2, 1, 128), v.state)
	if err != nil {
		fmt.Printf("invalid state tensor: %v\n", v.state)
		return 0, err
	}
	defer stateTensor.Destroy()

	// Prepare outputs
	// Output 0: Probability (1, 1)
	// Output 1: stateN (2, 1, 128)

	outProb := make([]float32, 1)                                                          // Batch 1
	outProbTensor, err := onnxruntime_go.NewTensor(onnxruntime_go.NewShape(1, 1), outProb) // Shape might be (1, 1) or just (1)
	if err != nil {
		fmt.Printf("invalid outProb tensor: %v\n", outProb)
		return 0, err
	}
	defer outProbTensor.Destroy()

	outState := make([]float32, 2*1*128)
	outStateTensor, err := onnxruntime_go.NewTensor(onnxruntime_go.NewShape(2, 1, 128), outState)
	if err != nil {
		fmt.Printf("invalid outState tensor: %v\n", outState)
		return 0, err
	}
	defer outStateTensor.Destroy()

	inputs := []onnxruntime_go.Value{inputTensor, stateTensor, srTensor}
	outputs := []onnxruntime_go.Value{outProbTensor, outStateTensor}
	// Run inference
	err = v.session.Run(inputs, outputs)
	if err != nil {
		fmt.Printf("VAD error: %v\n", err)
		return 0, err
	}

	// Get output probability
	probData := outProbTensor.GetData()
	prob := probData[0]

	// Update states
	v.state = outStateTensor.GetData()

	return prob, nil
}

// Int16ToFloat32 converts int16 audio samples to float32 [-1, 1]
func Int16ToFloat32(data []byte) []float32 {

	floats := make([]float32, len(data)/2)
	for i := 0; i < len(floats); i++ {
		sample := int16(data[i*2]) | int16(data[i*2+1])<<8
		floats[i] = float32(sample) / 32768.0
	}
	return floats
}

// Uint8ToFloat32 converts uint8 audio samples to float32 [-1, 1]
func Uint8ToFloat32(data []uint8) []float32 {
	floats := make([]float32, len(data))
	for i := 0; i < len(data); i++ {
		// Convert uint8 (0-255) to float32 (-1.0 to ~1.0)
		// Shift by 128 and scale by 128.0
		floats[i] = float32(data[i]-128) / 128.0
	}
	return floats
}

// Float32ToInt16 converts a single float32 sample to int16
func Float32ToInt16(sample float32) int {
	// Scale float32 [-1, 1] to int16 range [-32768, 32767]
	scaled := sample * 32768.0
	// Clamp values to int16 range
	if scaled > 32767.0 {
		return 32767
	} else if scaled < -32768.0 {
		return -32768
	} else {
		return int(scaled)
	}
}
