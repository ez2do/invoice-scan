# Phase 02: Backend File Storage

**Parent Plan**: [plan.md](./plan.md)  
**Dependencies**: Phase 01  
**Status**: Pending  
**Priority**: High

## Overview

Implement file storage following DDD architecture. Interface defined in domain layer, local filesystem implementation in adapters layer. Design allows easy swap to S3 in future.

## Key Insights

- **DDD Architecture**: Interface in domain, implementation in adapters
- Interface pattern enables future S3 migration without code changes
- Local storage sufficient for MVP/self-hosted deployment
- ULID-based naming prevents filename collisions

## DDD Directory Structure

```
backend/internal/
├── domain/
│   └── storage/
│       └── repository.go      # FileStorage interface
├── adapters/
│   └── storage/
│       ├── local_storage.go   # Local filesystem implementation
│       └── s3_storage.go      # Future S3 implementation
```

## Requirements

1. Define FileStorage interface in domain layer
2. Implement LocalStorage in adapters layer
3. Add upload directory configuration
4. Serve static files via Gin

## Related Code Files

- NEW: `backend/internal/domain/storage/repository.go` - Interface definition
- NEW: `backend/internal/adapters/storage/local_storage.go` - Local implementation
- `backend/internal/config/config.go` - Add upload path config
- `backend/cmd/server/main.go` - Add static file serving

## Implementation Steps

### Step 1: Create Storage Interface (Domain Layer)

**`internal/domain/storage/repository.go`**:
```go
package storage

import "context"

type FileStorage interface {
    Save(ctx context.Context, filename string, data []byte, contentType string) (string, error)
    Get(ctx context.Context, path string) ([]byte, error)
    Delete(ctx context.Context, path string) error
    GetURL(path string) string
}
```

### Step 2: Implement Local Storage (Adapters Layer)

**`internal/adapters/storage/local_storage.go`**:
```go
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
```

### Step 3: Update Config

Add to `config.go`:
```go
type Config struct {
    // ... existing fields
    UploadPath string
    BaseURL    string
}

func Load() (*Config, error) {
    // ...
    config := &Config{
        // ... existing
        UploadPath: getEnv("UPLOAD_PATH", "./uploads"),
        BaseURL:    getEnv("BASE_URL", "http://localhost:3001"),
    }
    return config, nil
}
```

### Step 4: Configure Static Serving

In `main.go`:
```go
router.Static("/uploads", cfg.UploadPath)
```

### Step 5: Initialize Storage

In `main.go`:
```go
import adapterstorage "invoice-scan/backend/internal/adapters/storage"

fileStorage, err := adapterstorage.NewLocalStorage(cfg.UploadPath, cfg.BaseURL)
if err != nil {
    log.Fatalf("Failed to create file storage: %v", err)
}
```

## Todo List

- [ ] Create internal/domain/storage/ directory
- [ ] Create repository.go with FileStorage interface
- [ ] Create internal/adapters/storage/ directory
- [ ] Create local_storage.go with LocalStorage implementation
- [ ] Add UploadPath and BaseURL to config
- [ ] Create upload directory on startup
- [ ] Configure static file serving in main.go
- [ ] Write unit tests for local storage

## Success Criteria

- [ ] FileStorage interface in domain layer (no infrastructure deps)
- [ ] LocalStorage implementation in adapters layer
- [ ] Can save image file via storage interface
- [ ] Can retrieve image via HTTP request
- [ ] File paths stored in database correctly
- [ ] Interface is swappable to S3

## Future S3 Implementation Stub

**`internal/adapters/storage/s3_storage.go`** (future):
```go
package storage

import (
    "context"
    
    "github.com/aws/aws-sdk-go-v2/service/s3"
    domainstorage "invoice-scan/backend/internal/domain/storage"
)

var _ domainstorage.FileStorage = (*S3Storage)(nil)

type S3Storage struct {
    client     *s3.Client
    bucketName string
    region     string
    baseURL    string
}

func NewS3Storage(client *s3.Client, bucketName, region, baseURL string) *S3Storage {
    return &S3Storage{
        client:     client,
        bucketName: bucketName,
        region:     region,
        baseURL:    baseURL,
    }
}

func (s *S3Storage) Save(ctx context.Context, filename string, data []byte, contentType string) (string, error) {
    // S3 PutObject implementation
}

func (s *S3Storage) Get(ctx context.Context, path string) ([]byte, error) {
    // S3 GetObject implementation
}

func (s *S3Storage) Delete(ctx context.Context, path string) error {
    // S3 DeleteObject implementation
}

func (s *S3Storage) GetURL(path string) string {
    // Return S3 URL or CloudFront URL
}
```

## Environment Variables

```env
UPLOAD_PATH=./uploads
BASE_URL=http://localhost:3001
```

## Security Considerations

- Validate file types before saving (image/* only)
- Limit file size (already 10MB limit in handler)
- Generate safe filenames using ULID (no path traversal)

## Next Steps

After completion, proceed to [Phase 03: Backend Async API](./phase-03-backend-async-api.md)
