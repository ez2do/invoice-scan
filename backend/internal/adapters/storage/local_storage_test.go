package storage

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	domainstorage "invoice-scan/backend/internal/domain/storage"
)

func setupTestStorage(t *testing.T) (*LocalStorage, string) {
	tmpDir := t.TempDir()
	baseURL := "http://localhost:3001"

	storage, err := NewLocalStorage(tmpDir, baseURL)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	return storage, tmpDir
}

func TestNewLocalStorage(t *testing.T) {
	tmpDir := t.TempDir()
	baseURL := "http://localhost:3001"

	storage, err := NewLocalStorage(tmpDir, baseURL)
	if err != nil {
		t.Fatalf("NewLocalStorage() error = %v", err)
	}

	if storage.basePath != tmpDir {
		t.Errorf("Expected basePath %v, got %v", tmpDir, storage.basePath)
	}
	if storage.baseURL != baseURL {
		t.Errorf("Expected baseURL %v, got %v", baseURL, storage.baseURL)
	}

	if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
		t.Error("Upload directory should be created")
	}
}

func TestLocalStorage_Save(t *testing.T) {
	storage, tmpDir := setupTestStorage(t)
	ctx := context.Background()

	filename := "test.jpg"
	data := []byte("test image data")
	contentType := "image/jpeg"

	path, err := storage.Save(ctx, filename, data, contentType)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	expectedPath := filepath.Join(tmpDir, filename)
	if path != expectedPath {
		t.Errorf("Expected path %v, got %v", expectedPath, path)
	}

	savedData, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read saved file: %v", err)
	}

	if string(savedData) != string(data) {
		t.Errorf("Expected data %v, got %v", data, savedData)
	}
}

func TestLocalStorage_Save_InvalidContentType(t *testing.T) {
	storage, _ := setupTestStorage(t)
	ctx := context.Background()

	_, err := storage.Save(ctx, "test.txt", []byte("data"), "text/plain")
	if err == nil {
		t.Error("Expected error for invalid content type")
	}
}

func TestLocalStorage_Get(t *testing.T) {
	storage, tmpDir := setupTestStorage(t)
	ctx := context.Background()

	filename := "test.jpg"
	data := []byte("test image data")
	filePath := filepath.Join(tmpDir, filename)

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	got, err := storage.Get(ctx, filePath)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if string(got) != string(data) {
		t.Errorf("Expected data %v, got %v", data, got)
	}
}

func TestLocalStorage_Get_NotFound(t *testing.T) {
	storage, _ := setupTestStorage(t)
	ctx := context.Background()

	_, err := storage.Get(ctx, "/nonexistent/file.jpg")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestLocalStorage_Delete(t *testing.T) {
	storage, tmpDir := setupTestStorage(t)
	ctx := context.Background()

	filename := "test.jpg"
	filePath := filepath.Join(tmpDir, filename)

	if err := os.WriteFile(filePath, []byte("data"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	if err := storage.Delete(ctx, filePath); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Error("File should be deleted")
	}
}

func TestLocalStorage_Delete_NotFound(t *testing.T) {
	storage, _ := setupTestStorage(t)
	ctx := context.Background()

	err := storage.Delete(ctx, "/nonexistent/file.jpg")
	if err != nil {
		t.Errorf("Delete() should not error for non-existent file, got %v", err)
	}
}

func TestLocalStorage_GetURL(t *testing.T) {
	storage, tmpDir := setupTestStorage(t)

	filePath := filepath.Join(tmpDir, "test.jpg")
	url := storage.GetURL(filePath)

	expectedURL := "http://localhost:3001/uploads/test.jpg"
	if url != expectedURL {
		t.Errorf("Expected URL %v, got %v", expectedURL, url)
	}
}

func TestLocalStorage_ImplementsInterface(t *testing.T) {
	var _ domainstorage.FileStorage = (*LocalStorage)(nil)
}

