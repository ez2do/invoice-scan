package invoice

import "context"

type ExtractionService interface {
	Extract(ctx context.Context, imageBytes []byte, mimeType string) (ExtractedData, error)
	Close() error
}
