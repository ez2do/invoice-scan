# Phase 1: Backend API Pagination

**Priority**: High | **Status**: ⏳ Pending

## Context

Current `List()` returns all invoices without limits, which will not scale as data grows.

## Overview

Add offset-based pagination to the invoice list API:
- Accept `page` (default=1) and `page_size` (default=10) query params
- Return paginated results with `total`, `page`, `page_size`, `total_pages` metadata

## Related Files

- [repository.go](file:///Users/tuananh/projects/invoice-scan/backend/internal/domain/invoice/repository.go)
- [invoice_gorm_repo.go](file:///Users/tuananh/projects/invoice-scan/backend/internal/adapters/repo/invoice_gorm_repo.go)
- [invoice.go handler](file:///Users/tuananh/projects/invoice-scan/backend/internal/handlers/invoice.go)
- [dto.go](file:///Users/tuananh/projects/invoice-scan/backend/internal/handlers/dto.go)

---

## Implementation Steps

### 1. Add pagination types to domain layer

#### [MODIFY] [repository.go](file:///Users/tuananh/projects/invoice-scan/backend/internal/domain/invoice/repository.go)

Add `PaginationParams` and `PaginatedResult` types. Update `List()` signature to accept pagination params and return paginated result.

```go
type PaginationParams struct {
    Page     int
    PageSize int
}

type PaginatedResult struct {
    Invoices   Invoices
    Total      int64
    Page       int
    PageSize   int
    TotalPages int
}
```

---

### 2. Update repository implementation

#### [MODIFY] [invoice_gorm_repo.go](file:///Users/tuananh/projects/invoice-scan/backend/internal/adapters/repo/invoice_gorm_repo.go)

Implement paginated `List()`:
- Count total records
- Apply LIMIT/OFFSET
- Calculate total pages

---

### 3. Add paginated response DTO

#### [MODIFY] [dto.go](file:///Users/tuananh/projects/invoice-scan/backend/internal/handlers/dto.go)

Add `PaginatedInvoicesResponse`:
```go
type PaginatedInvoicesResponse struct {
    Data       []InvoiceData `json:"data"`
    Total      int64         `json:"total"`
    Page       int           `json:"page"`
    PageSize   int           `json:"page_size"`
    TotalPages int           `json:"total_pages"`
}
```

---

### 4. Update handler to support pagination

#### [MODIFY] [invoice.go](file:///Users/tuananh/projects/invoice-scan/backend/internal/handlers/invoice.go)

Parse query params `page` and `page_size`, validate, pass to repo, return paginated response.

---

## Todo

- [ ] Add pagination types to domain
- [ ] Update repository interface
- [ ] Implement paginated List in GORM repo
- [ ] Add PaginatedInvoicesResponse DTO
- [ ] Update List handler

## Success Criteria

- API accepts `?page=1&page_size=10` parameters
- Returns correct pagination metadata
- Default to page=1, page_size=10 if not specified
- Handles invalid values gracefully (use defaults)
