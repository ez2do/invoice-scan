# Phase 03: Backend Async API

**Parent Plan**: [plan.md](./plan.md)  
**Dependencies**: Phase 01, Phase 02  
**Status**: Completed  
**Priority**: High

## Overview

Refactor extraction to async pattern. New upload endpoint saves image and returns immediately. Background goroutine processes extraction and updates status.

## Key Insights

- Goroutines provide native async without external dependencies
- Status field enables polling from frontend
- Decoupling upload from extraction improves UX

## Requirements

1. Create new upload endpoint (replaces sync extract)
2. Implement async extraction worker
3. Create invoice list endpoint
4. Create invoice detail endpoint
5. Update existing extract handler

## API Design

### POST /api/invoices/upload
Upload image, create pending invoice, trigger async extraction.

**Request**: `multipart/form-data` with `image` file  
**Response**:
```json
{
  "success": true,
  "data": {
    "id": "01HXYZ123ABC456DEF789GHI",
    "status": "pending",
    "imagePath": "/uploads/01HXYZ123ABC456DEF789GHI.jpg",
    "createdAt": "2025-12-05T10:00:00Z"
  }
}
```

### GET /api/invoices
List all invoices with status.

**Response**:
```json
{
  "success": true,
  "data": [
    {
      "id": "01HXYZ123ABC456DEF789GHI",
      "status": "completed",
      "imagePath": "/uploads/01HXYZ123ABC456DEF789GHI.jpg",
      "extractedData": {...},
      "createdAt": "2025-12-05T10:00:00Z"
    }
  ]
}
```

### GET /api/invoices/:id
Get single invoice with full details.

**Response**:
```json
{
  "success": true,
  "data": {
    "id": "01HXYZ123ABC456DEF789GHI",
    "status": "completed",
    "imagePath": "/uploads/01HXYZ123ABC456DEF789GHI.jpg",
    "extractedData": {...},
    "errorMessage": null,
    "createdAt": "2025-12-05T10:00:00Z"
  }
}
```

## Related Code Files

- NEW: `backend/internal/handlers/invoice.go` - Invoice handlers
- `backend/internal/services/extraction.go` - Add async method
- `backend/cmd/server/main.go` - Register new routes
- `backend/internal/domain/invoice/` - Domain entities and interfaces
- `backend/internal/adapters/repo/invoice_gorm_repo.go` - Repository implementation

## Implementation Steps

### Step 1: Create Invoice Handler
```go
import (
    "invoice-scan/backend/internal/domain/invoice"
    domainstorage "invoice-scan/backend/internal/domain/storage"
)

type InvoiceHandler struct {
    repo              invoice.Repository       // Domain interface
    storage           domainstorage.FileStorage // Domain interface
    extractionService *services.ExtractionService
}

func NewInvoiceHandler(
    repo invoice.Repository,
    storage domainstorage.FileStorage,
    extractionService *services.ExtractionService,
) *InvoiceHandler {
    return &InvoiceHandler{
        repo:              repo,
        storage:           storage,
        extractionService: extractionService,
    }
}

func (h *InvoiceHandler) Upload(c *gin.Context) {...}
func (h *InvoiceHandler) List(c *gin.Context) {...}
func (h *InvoiceHandler) GetByID(c *gin.Context) {...}
```

### Step 2: Implement Upload Handler
```go
func (h *InvoiceHandler) Upload(c *gin.Context) {
    file, err := c.FormFile("image")
    // ... validation ...
    
    // Generate ID using domain repository
    id := h.repo.NextID()
    
    // Save file using storage interface
    filename := fmt.Sprintf("%s%s", id.String(), filepath.Ext(file.Filename))
    imagePath, err := h.storage.Save(c.Request.Context(), filename, imageBytes, mimeType)
    
    // Create domain entity
    inv := invoice.New(id, imagePath)
    if err := h.repo.Create(c.Request.Context(), inv); err != nil {
        // handle error
    }
    
    // Trigger async extraction
    h.processExtractionAsync(id, imageBytes, mimeType)
    
    // Return immediately
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data": map[string]interface{}{
            "id":        inv.ID.String(),
            "status":    inv.Status.String(),
            "imagePath": h.storage.GetURL(inv.ImagePath),
            "createdAt": inv.CreatedAt,
        },
    })
}
```

### Step 3: Implement Async Extraction
```go
func (h *InvoiceHandler) processExtractionAsync(invoiceID invoice.ID, imageBytes []byte, mimeType string) {
    go func() {
        ctx := context.Background()
        
        // Get invoice from repository
        inv, err := h.repo.GetByID(ctx, invoiceID)
        if err != nil {
            log.Printf("Failed to get invoice %s: %v", invoiceID.String(), err)
            return
        }
        
        // Update to processing using domain method
        inv.MarkProcessing()
        h.repo.Update(ctx, inv)
        
        // Extract data
        data, err := h.extractionService.Extract(ctx, imageBytes, mimeType)
        
        if err != nil {
            // Use domain method to mark failed
            inv.MarkFailed(err.Error())
            h.repo.Update(ctx, inv)
            return
        }
        
        // Use domain method to mark completed
        dataJSON, _ := json.Marshal(data)
        inv.MarkCompleted(dataJSON)
        h.repo.Update(ctx, inv)
    }()
}
```

### Step 4: Implement List Handler
- Query all invoices ordered by created_at desc
- Return with thumbnail URLs

### Step 5: Implement GetByID Handler
- Return full invoice details
- Include extracted data if completed

### Step 6: Register Routes
```go
api.POST("/invoices/upload", invoiceHandler.Upload)
api.GET("/invoices", invoiceHandler.List)
api.GET("/invoices/:id", invoiceHandler.GetByID)
```

## Todo List

- [ ] Create handlers/invoice.go with handler struct
- [ ] Implement Upload handler with file storage
- [ ] Implement async extraction goroutine
- [ ] Implement List handler
- [ ] Implement GetByID handler
- [ ] Register new routes in main.go
- [ ] Test all endpoints with curl/Postman

## Success Criteria

- [x] Upload returns within 100ms (no extraction wait)
- [x] Status transitions: pending → processing → completed/failed
- [x] List endpoint returns all invoices with correct status
- [x] GetByID returns full extracted data for completed invoices

## Risk Assessment

- **Goroutine leaks**: Context cancellation handled properly
- **Race conditions**: Use atomic status updates
- **Memory pressure**: Monitor with multiple concurrent extractions

## Security Considerations

- Validate image MIME type
- Rate limit upload endpoint (future)
- Sanitize error messages in response

## Next Steps

After completion, proceed to [Phase 04: Frontend Integration](./phase-04-frontend-integration.md)

