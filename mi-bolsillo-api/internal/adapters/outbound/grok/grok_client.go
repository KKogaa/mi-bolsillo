package grok

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/KKogaa/mi-bolsillo-api/internal/core/domain/entities"
)

type GrokClient struct {
	apiKey     string
	httpClient *http.Client
}

func NewGrokClient(apiKey string) *GrokClient {
	return &GrokClient{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type BillItem struct {
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
}

type ParsedBillData struct {
	Items        []BillItem `json:"items"`
	TotalAmount  float64    `json:"total_amount"`
	Currency     string     `json:"currency"`
	Date         string     `json:"date"`
	MerchantName string     `json:"merchant_name"`
}

type grokRequest struct {
	Messages []grokMessage `json:"messages"`
	Model    string        `json:"model"`
	Stream   bool          `json:"stream"`
}

type grokMessage struct {
	Role    string        `json:"role"`
	Content []grokContent `json:"content"`
}

type grokContent struct {
	Type     string          `json:"type"`
	Text     string          `json:"text,omitempty"`
	ImageURL *grokImageURL   `json:"image_url,omitempty"`
}

type grokImageURL struct {
	URL string `json:"url"`
}

type grokResponse struct {
	ID      string       `json:"id"`
	Choices []grokChoice `json:"choices"`
}

type grokChoice struct {
	Message grokResponseMessage `json:"message"`
}

type grokResponseMessage struct {
	Content string `json:"content"`
}

func (c *GrokClient) ParseBillImage(imageData []byte) (*ParsedBillData, error) {
	// Encode image to base64
	base64Image := base64.StdEncoding.EncodeToString(imageData)
	dataURL := fmt.Sprintf("data:image/jpeg;base64,%s", base64Image)

	// Create the request payload
	reqBody := grokRequest{
		Model:  "grok-2-vision-1212",
		Stream: false,
		Messages: []grokMessage{
			{
				Role: "user",
				Content: []grokContent{
					{
						Type: "image_url",
						ImageURL: &grokImageURL{
							URL: dataURL,
						},
					},
					{
						Type: "text",
						Text: `Analyze this bill/receipt image and extract the following information in JSON format:
{
  "items": [
    {
      "description": "item name",
      "amount": numeric_amount,
      "category": "Food|Transportation|Entertainment|Shopping|Utilities|Healthcare|Other"
    }
  ],
  "total_amount": numeric_total,
  "currency": "USD|PEN|EUR|etc",
  "date": "YYYY-MM-DD",
  "merchant_name": "store/restaurant name"
}

Rules:
- Extract ALL line items from the receipt
- Categorize each item appropriately
- Use the currency symbol or text to determine the currency (default to USD if unclear)
- Extract the date in YYYY-MM-DD format (use today's date if not visible)
- Return ONLY valid JSON, no additional text or explanation`,
					},
				},
			},
		},
	}

	// Marshal request to JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", "https://api.x.ai/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("grok API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse Grok response
	var grokResp grokResponse
	if err := json.Unmarshal(body, &grokResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal grok response: %w", err)
	}

	if len(grokResp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in grok response")
	}

	// Extract JSON from response content
	content := grokResp.Choices[0].Message.Content

	// Parse the bill data from the content
	var parsedData ParsedBillData
	if err := json.Unmarshal([]byte(content), &parsedData); err != nil {
		return nil, fmt.Errorf("failed to parse bill data from response: %w", err)
	}

	return &parsedData, nil
}

// Intent detection types
type grokTextRequest struct {
	Messages []grokTextMessage `json:"messages"`
	Model    string            `json:"model"`
	Stream   bool              `json:"stream"`
}

type grokTextMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// DetectIntent analyzes user text and determines their intent
// Implements the ports.IntentDetector interface
func (c *GrokClient) DetectIntent(userText string) (*entities.Intent, error) {
	systemPrompt := `You are an intent classifier for a bill/expense management application.
Analyze the user's message and determine their intent. Return ONLY a valid JSON object with this structure:
{
  "type": "list_bills|summary_bills|upload_bill|unknown",
  "confidence": 0.0-1.0,
  "parameters": {
    "period": "last_month|this_month|last_week|all_time" (for summary_bills),
    "limit": number (for list_bills, how many bills to show)
  }
}

Intent types:
- list_bills: User wants to see a list of recent bills (e.g., "show my bills", "list my last bills", "what are my recent expenses")
- summary_bills: User wants a summary/total of bills for a period (e.g., "how much did I spend last month", "summary of this month", "total expenses")
- upload_bill: User wants to upload/add a bill (this is detected when they send an image, not text)
- unknown: Cannot determine intent or asking something else

Return ONLY valid JSON, no additional text.`

	reqBody := grokTextRequest{
		Model:  "grok-2-1212",
		Stream: false,
		Messages: []grokTextMessage{
			{
				Role:    "system",
				Content: systemPrompt,
			},
			{
				Role:    "user",
				Content: userText,
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.x.ai/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("grok API error (status %d): %s", resp.StatusCode, string(body))
	}

	var grokResp grokResponse
	if err := json.Unmarshal(body, &grokResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal grok response: %w", err)
	}

	if len(grokResp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in grok response")
	}

	content := grokResp.Choices[0].Message.Content

	var intent entities.Intent
	if err := json.Unmarshal([]byte(content), &intent); err != nil {
		return nil, fmt.Errorf("failed to parse intent from response: %w", err)
	}

	return &intent, nil
}
