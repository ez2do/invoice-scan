package storage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	domainstorage "invoice-scan/backend/internal/domain/storage"
)

var _ domainstorage.FileStorage = (*LocalStorage)(nil)

type LocalStorage struct {
	basePath string
	baseURL  string
}

func NewLocalStorage(basePath, baseURL string) (*LocalStorage, error) {
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	return &LocalStorage{
		basePath: basePath,
		baseURL:  strings.TrimSuffix(baseURL, "/"),
	}, nil
}

func (s *LocalStorage) Save(ctx context.Context, filename string, data []byte, contentType string) (string, error) {
	if !strings.HasPrefix(contentType, "image/") {
		return "", fmt.Errorf("invalid content type: %s", contentType)
	}

	filePath := filepath.Join(s.basePath, filename)

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return filePath, nil
}

func (s *LocalStorage) Get(ctx context.Context, path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return data, nil
}

func (s *LocalStorage) Delete(ctx context.Context, path string) error {
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

func (s *LocalStorage) GetURL(path string) string {
	filename := filepath.Base(path)
	return fmt.Sprintf("%s/uploads/%s", s.baseURL, filename)
}
