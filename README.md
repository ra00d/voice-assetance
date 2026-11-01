# Kind Echo

Kind Echo is a fully offline, voice-activated interactive robot assistant for seniors.

## Prerequisites

### Hardware

- Microphone
- Speaker
- Display (for the face)

### Software

- Go 1.21+
- GCC (for CGO)
- PortAudio (for Malgo)
  - Ubuntu: `sudo apt-get install portaudio19-dev`
- ONNX Runtime (for Silero VAD)
  - Download `libonnxruntime.so` and place it in a library path or the project root.

### Models

You need to download the following models and place them in the `models/` directory:

1. **Silero VAD**: `silero_vad.onnx`
   - [Download](https://github.com/snakers4/silero-vad/raw/master/files/silero_vad.onnx)
2. **Whisper**: `ggml-base.en.bin`
   - [Download](https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-base.en.bin)
3. **Piper TTS**: `en_US-lessac-medium.onnx` and `en_US-lessac-medium.onnx.json`
   - [Download](https://github.com/rhasspy/piper/releases)
4. **Ollama**: Ensure Ollama is running with `llama3` (or your chosen model).
   - `ollama run llama3`

## Installation

1. **Setup Whisper.cpp**:

   ```bash
   ./setup_whisper.sh
   ```

2. **Build Kind Echo**:
3. **Build Kind Echo**:

   ```bash
   # Download ONNX Runtime (fixes system dependency issues)
   ./setup_onnx.sh

   go mod tidy
   # Set environment variables for CGO and ONNX
   export C_INCLUDE_PATH=$(pwd)/third_party/include
   export LIBRARY_PATH=$(pwd)/third_party/lib
   export LD_LIBRARY_PATH=$(pwd)/third_party/onnxruntime/lib:$LD_LIBRARY_PATH

   go build -o kind-echo cmd/kind-echo/main.go
   ```

## Usage

Run the application (ensure environment variables are set):

```bash
export C_INCLUDE_PATH=$(pwd)/third_party/include
export LIBRARY_PATH=$(pwd)/third_party/lib
export LD_LIBRARY_PATH=$(pwd)/third_party/onnxruntime/lib:$LD_LIBRARY_PATH

./kind-echo \
  --vad models/silero_vad.onnx \
  --whisper models/ggml-base.en.bin \
  --piper-model models/en_US-lessac-medium.onnx \
  --piper-bin /path/to/piper \
  --ollama llama3
```

## Troubleshooting

- **Audio Input**: If no audio is detected, check your default microphone settings.
- **CGO Errors**: Ensure you have `portaudio19-dev` and `libonnxruntime` installed.
