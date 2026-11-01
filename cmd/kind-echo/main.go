package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"kind-echo/pkg/audio"
	"kind-echo/pkg/brain"
	"kind-echo/pkg/stt"
	"kind-echo/pkg/tts"
	"kind-echo/pkg/ui"
)

func main() {
	// Flags for model paths
	vadModel := flag.String("vad", "models/silero_vad.onnx", "Path to Silero VAD ONNX model")
	whisperModel := flag.String("whisper", "models/ggml-base.en.bin", "Path to Whisper model")
	piperModel := flag.String("piper-model", "models/en_US-lessac-medium.onnx", "Path to Piper TTS model")
	piperBinary := flag.String("piper-bin", "piper-tts", "Path to Piper binary")
	ollamaModel := flag.String("ollama", "llama3", "Ollama model name")
	flag.Parse()

	fmt.Println("Starting Kind Echo...")

	// 1. Initialize UI (must run on main thread usually, but WebView.Run blocks)
	// We'll run the logic in a goroutine and UI in main.
	face := ui.NewFace()

	go func() {
		// Give UI time to start
		time.Sleep(1 * time.Second)

		// 2. Initialize Components
		fmt.Println("Initializing components...")

		// Audio Recorder
		recorder, err := audio.NewRecorder()
		if err != nil {
			log.Fatalf("Failed to init recorder: %v", err)
		}
		if err := recorder.Start(); err != nil {
			log.Fatalf("Failed to start recorder: %v", err)
		}
		defer recorder.Stop()

		// VAD
		vad, err := audio.NewVAD(*vadModel)
		if err != nil {
			log.Fatalf("Failed to init VAD: %v", err)
		}

		// STT
		whisper, err := stt.NewWhisper(*whisperModel)
		if err != nil {
			log.Printf("Warning: Failed to init Whisper: %v", err)
		}
		defer whisper.Close()

		// LLM
		llm := brain.NewLLM(*ollamaModel)
		intentEngine := brain.NewIntentEngine()

		// TTS
		piper := tts.NewPiper(*piperModel, *piperBinary)
		player, err := tts.NewPlayer()
		if err != nil {
			log.Printf("Warning: Failed to init Player: %v", err)
		}

		// Audio Loop
		loop := audio.NewLoop(recorder, vad)
		loop.Start()

		// Greet User
		greetText := "Hi, I am Kind Echo. I am here to help you."
		fmt.Println("Robot:", greetText)
		face.SetExpression("talking")
		if audioData, err := piper.Synthesize(greetText); err == nil {
			player.Play(audioData)
		} else {
			log.Printf("Greeting TTS Error: %v", err)
		}
		face.SetExpression("listening")

		// Main Event Loop
		for {
			select {
			case <-loop.SpeechStarted:
				face.SetExpression("recording")

			case segment := <-loop.Output:
				fmt.Println("Processing speech segment...")
				face.SetExpression("processing")

				// STT
				text, err := whisper.Transcribe(segment)
				if err != nil {
					log.Printf("STT Error: %v", err)
					face.SetExpression("sadness")
					time.Sleep(2 * time.Second)
					face.SetExpression("listening")
					continue
				}
				if text == "" {
					face.SetExpression("listening")
					continue
				}
				fmt.Printf("User: %s\n", text)

				// Intent
				intent := intentEngine.Analyze(text)
				if intent.Type == "set_reminder" {
					fmt.Println("Intent: Set Reminder")
					// Handle reminder...
				}

				// LLM
				response, err := llm.GenerateResponse(text)
				if err != nil {
					log.Printf("LLM Error: %v", err)
					response = "I'm sorry, I'm having trouble thinking right now."
				}
				fmt.Printf("Robot: %s\n", response)

				// TTS
				face.SetExpression("talking")
				audioData, err := piper.Synthesize(response)
				if err != nil {
					log.Printf("TTS Error: %v", err)
				} else {
					player.Play(audioData)
				}

				// Back to listening
				face.SetExpression("listening")
			}
		}
	}()

	face.Run()
}
