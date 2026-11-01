package brain

import (
	"strings"
	"time"
)

// Intent represents a detected user intent.
type Intent struct {
	Type    string
	Payload interface{}
}

// IntentEngine analyzes text for commands.
type IntentEngine struct {
	// In a real app, we might use the LLM to classify intents or a regex-based approach.
}

// NewIntentEngine creates a new IntentEngine.
func NewIntentEngine() *IntentEngine {
	return &IntentEngine{}
}

// Analyze determines the intent from the user's text.
func (ie *IntentEngine) Analyze(text string) *Intent {
	lowerText := strings.ToLower(text)

	if strings.Contains(lowerText, "remind me") || strings.Contains(lowerText, "set a reminder") {
		// Simple rule-based extraction (placeholder)
		return &Intent{
			Type: "set_reminder",
			Payload: map[string]string{
				"raw_text": text,
				"time":     time.Now().Add(1 * time.Hour).Format(time.RFC3339), // Default to 1 hour later
			},
		}
	}

	return &Intent{
		Type: "chat",
	}
}
