# Research Report: Async Extraction & Frontend Polling Pattern

**Date**: 2025-12-05  
**Researcher**: Agent 02  
**Focus**: Async processing, frontend-backend integration patterns

## 1. Async Processing Architecture

### Pattern: Background Worker with Status Polling
1. Client uploads image → receives invoice ID immediately
2. Backend queues extraction job
3. Worker processes extraction asynchronously
4. Client polls for status updates

### Go Goroutine Worker Pattern
```go
func (s *ExtractionService) ProcessAsync(invoiceID uint) {
    go func() {
        // Update status to processing
        s.repo.UpdateStatus(ctx, invoiceID, "processing", nil, nil)
        
        // Perform extraction
        data, err := s.Extract(ctx, imageBytes, mimeType)
        
        if err != nil {
            errStr := err.Error()
            s.repo.UpdateStatus(ctx, invoiceID, "failed", nil, &errStr)
            return
        }
        
        s.repo.UpdateStatus(ctx, invoiceID, "completed", data, nil)
    }()
}
```

## 2. API Design for Async Flow

### Upload Endpoint
```
POST /api/invoices/upload
Request: multipart/form-data (image file)
Response: { "id": 123, "status": "pending" }
```

### Status Endpoint
```
GET /api/invoices/:id
Response: {
  "id": 123,
  "status": "completed|pending|processing|failed",
  "data": {...} | null,
  "error": "..." | null,
  "imagePath": "/uploads/..."
}
```

### List Endpoint
```
GET /api/invoices
Response: [
  { "id": 123, "status": "completed", ... },
  { "id": 124, "status": "processing", ... }
]
```

## 3. Frontend Polling Implementation

### TanStack Query Polling
```typescript
const { data, isLoading } = useQuery({
  queryKey: ['invoice', id],
  queryFn: () => fetchInvoice(id),
  refetchInterval: (data) => {
    if (data?.status === 'completed' || data?.status === 'failed') {
      return false; // stop polling
    }
    return 2000; // poll every 2 seconds
  },
});
```

### List Page Auto-Refresh
```typescript
const { data: invoices } = useQuery({
  queryKey: ['invoices'],
  queryFn: fetchInvoices,
  refetchInterval: 3000, // refresh list every 3 seconds
});
```

## 4. Zustand Store Updates

### New Store Structure
```typescript
interface InvoiceListItem {
  id: number;
  status: 'pending' | 'processing' | 'completed' | 'failed';
  thumbnailUrl: string;
  createdAt: string;
  extractedData?: InvoiceData;
  error?: string;
}

interface AppStore {
  invoices: InvoiceListItem[];
  setInvoices: (invoices: InvoiceListItem[]) => void;
  updateInvoiceStatus: (id: number, status: string, data?: InvoiceData) => void;
}
```

## 5. Frontend Flow Changes

### Current Flow (Sync)
1. Take picture → Review → Extract (wait) → Show data

### New Flow (Async)
1. Take picture → Review → Upload (get ID) → Navigate to list
2. List shows all invoices with status
3. Click on completed invoice → View extracted data

## 6. Image Serving Options

### Option A: Static File Serving
```go
router.Static("/uploads", "./uploads")
```

### Option B: API Endpoint
```go
router.GET("/api/invoices/:id/image", handler.GetImage)
```

Recommend Option A for simplicity, with path stored in DB.

## 7. Key Benefits of Async Pattern

1. **Better UX**: User can continue scanning while extraction processes
2. **Scalability**: Multiple extractions can run concurrently
3. **Resilience**: Failed extractions don't block the user
4. **Progress Visibility**: User sees all invoice states in list

## References
- [TanStack Query Polling](https://tanstack.com/query/latest/docs/react/guides/polling)
- [Go Concurrency Patterns](https://go.dev/blog/pipelines)

