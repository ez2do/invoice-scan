package services

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"invoice-scan/backend/internal/models"

	"google.golang.org/genai"
)

type ExtractionService struct {
	client  *genai.Client
	apiKey  string
	timeout time.Duration
}

func NewExtractionService(apiKey string) (*ExtractionService, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	return &ExtractionService{
		client:  client,
		apiKey:  apiKey,
		timeout: 30 * time.Second,
	}, nil
}

func (s *ExtractionService) Close() {
}

func (s *ExtractionService) Extract(ctx context.Context, base64Image string) (*models.InvoiceData, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	imageBytes, mimeType, err := decodeBase64Image(base64Image)
	if err != nil {
		return nil, fmt.Errorf("invalid image: %w", err)
	}

	prompt := buildExtractionPrompt()

	parts := []*genai.Part{
		{Text: prompt},
		{InlineData: &genai.Blob{Data: imageBytes, MIMEType: mimeType}},
	}

	contents := []*genai.Content{
		{Parts: parts},
	}

	result, err := s.client.Models.GenerateContent(ctx, "gemini-2.0-flash-exp", contents, nil)
	if err != nil {
		return nil, fmt.Errorf("Gemini API error: %w", err)
	}

	if result == nil || len(result.Candidates) == 0 {
		return nil, fmt.Errorf("empty response from Gemini API")
	}

	text := result.Text()
	if text == "" {
		return nil, fmt.Errorf("empty text in Gemini response")
	}

	invoiceData, err := parseGeminiResponse(text)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return invoiceData, nil
}

func decodeBase64Image(base64Str string) ([]byte, string, error) {
	const maxSize = 10 * 1024 * 1024

	if len(base64Str) > maxSize {
		return nil, "", fmt.Errorf("image too large")
	}

	var imageData string
	var mimeType string

	if strings.HasPrefix(base64Str, "data:image/") {
		parts := strings.SplitN(base64Str, ",", 2)
		if len(parts) != 2 {
			return nil, "", fmt.Errorf("invalid data URL format")
		}
		mimePart := strings.TrimSuffix(parts[0], ";base64")
		mimeType = strings.TrimPrefix(mimePart, "data:")
		imageData = parts[1]
	} else {
		imageData = base64Str
		mimeType = "image/jpeg"
	}

	decoded, err := base64.StdEncoding.DecodeString(imageData)
	if err != nil {
		return nil, "", fmt.Errorf("invalid base64 encoding: %w", err)
	}

	if len(decoded) == 0 {
		return nil, "", fmt.Errorf("empty image data")
	}

	return decoded, mimeType, nil
}

func buildExtractionPrompt() string {
	return `Extract invoice data from this image and return ONLY valid JSON in the following format:

{
  "keyValuePairs": [
    {"key": "Invoice Number", "value": "...", "confidence": 0.95},
    {"key": "Date", "value": "...", "confidence": 0.90},
    {"key": "Vendor", "value": "...", "confidence": 0.85}
  ],
  "table": {
    "headers": ["Item", "Quantity", "Price", "Total"],
    "rows": [
      ["Item 1", "2", "100000", "200000"],
      ["Item 2", "1", "50000", "50000"]
    ]
  },
  "summary": [
    {"key": "Subtotal", "value": "...", "confidence": 0.95},
    {"key": "Tax", "value": "...", "confidence": 0.90},
    {"key": "Total", "value": "...", "confidence": 0.95}
  ],
  "confidence": 0.90
}

Rules:
- Extract all visible invoice information
- Support Vietnamese language (both printed and handwritten)
- If table exists, extract it with headers and rows
- If no table, set "table" to null
- Include confidence scores (0.0 to 1.0) for each field
- Return ONLY the JSON object, no additional text or markdown
- Use Vietnamese field names if the invoice is in Vietnamese
- Extract dates, amounts, and numbers accurately`
}

func parseGeminiResponse(text string) (*models.InvoiceData, error) {
	text = strings.TrimSpace(text)

	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimPrefix(text, "```")
	text = strings.TrimSuffix(text, "```")
	text = strings.TrimSpace(text)

	var invoiceData models.InvoiceData
	if err := json.Unmarshal([]byte(text), &invoiceData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &invoiceData, nil
}
