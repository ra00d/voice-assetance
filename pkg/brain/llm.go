package brain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// LLM handles interaction with the Ollama API.
type LLM struct {
	baseURL string
	model   string
	history []Message
}

// Message represents a chat message.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// NewLLM creates a new LLM client.
func NewLLM(model string) *LLM {
	return &LLM{
		baseURL: "http://localhost:11434",
		model:   model,
		history: []Message{
			{Role: "system", Content: "You are Kind Echo, a friendly robot assistant for seniors. Keep your answers short, clear, and warm."},
		},
	}
}

// GenerateResponse sends a prompt to the LLM and returns the response.
func (l *LLM) GenerateResponse(prompt string) (string, error) {
	l.history = append(l.history, Message{Role: "user", Content: prompt})

	requestBody := map[string]interface{}{
		"model":    l.model,
		"messages": l.history,
		"stream":   false,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(l.baseURL+"/api/chat", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama api error: %s", string(body))
	}

	var response struct {
		Message Message `json:"message"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	l.history = append(l.history, response.Message)
	return response.Message.Content, nil
}

// ClearHistory resets the conversation history.
func (l *LLM) ClearHistory() {
	l.history = []Message{
		{Role: "system", Content: "You are Kind Echo, a friendly robot assistant for seniors. Keep your answers short, clear, and warm."},
	}
}
