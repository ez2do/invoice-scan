package storage

import "context"

type FileStorage interface {
	Save(ctx context.Context, filename string, data []byte, contentType string) (string, error)
	Get(ctx context.Context, path string) ([]byte, error)
	Delete(ctx context.Context, path string) error
	GetURL(path string) string
}
