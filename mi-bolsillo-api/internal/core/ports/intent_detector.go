package ports

import "github.com/KKogaa/mi-bolsillo-api/internal/core/domain/entities"

// IntentDetector defines the outbound port for detecting user intent from text
// This is an interface that external adapters (like GrokClient) will implement
type IntentDetector interface {
	DetectIntent(userText string) (*entities.Intent, error)
}
