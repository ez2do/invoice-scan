# Implementation Plan: FE-BE Integration with Async Extraction

**Created**: 2025-12-05  
**Status**: Planning  
**Priority**: High

## Overview

Integrate frontend with backend, implement MySQL/GORM persistence, file storage with S3-ready interface, and async extraction processing. Users can take pictures continuously while extractions process in background. List screen shows real-time extraction status.

## Objectives

1. Remove all frontend mocks, integrate with real backend APIs
2. Implement MySQL database with GORM for invoice persistence
3. Create file storage interface (local now, S3-ready)
4. Implement async extraction processing
5. Display extraction status on list screen (pending/processing/completed/failed)

## Research References

- [GORM/MySQL & File Storage](./research/researcher-01-gorm-mysql-file-storage.md)
- [Async Extraction & Polling](./research/researcher-02-async-extraction-polling.md)

## Implementation Phases

| Phase | Name | Status | Progress |
|-------|------|--------|----------|
| 01 | [Backend Database Setup](./phase-01-backend-database-setup.md) | Completed | 100% |
| 02 | [Backend File Storage](./phase-02-backend-file-storage.md) | Completed | 100% |
| 03 | [Backend Async API](./phase-03-backend-async-api.md) | Completed | 100% |
| 04 | [Frontend Integration](./phase-04-frontend-integration.md) | Pending | 0% |
| 05 | [Testing & Polish](./phase-05-testing-polish.md) | Pending | 0% |

## High-Level Architecture

```
┌─────────────┐    POST /upload    ┌─────────────────┐
│   Frontend  │ ─────────────────► │   Backend API   │
│    (PWA)    │                    │   (Gin + DDD)   │
└─────────────┘                    └────────┬────────┘
       │                                    │
       │  GET /invoices                     │ goroutine
       │  GET /invoices/:id                 ▼
       │                           ┌─────────────────┐
       └──────────────────────────►│  MySQL Database │
                                   └─────────────────┘
                                           │
                                   ┌───────┴───────┐
                                   │               │
                              FileStorage    Gemini API
                              (Local/S3)
```

## Backend DDD Structure

```
backend/internal/
├── domain/
│   └── invoice/
│       ├── invoice.go       # Entity + types (ID, Status)
│       └── repository.go    # Repository interface
├── adapters/
│   └── repo/
│       └── invoice_gorm_repo.go  # GORM impl (gormInvoice)
├── handlers/
├── services/
└── config/
```

## Key Decisions

1. **Polling over WebSocket**: Simpler to implement, adequate for use case
2. **Local storage first**: S3 interface ready but local implementation for MVP
3. **JSON column for extracted data**: Flexible schema for varying invoice formats
4. **Goroutines for async**: Native Go concurrency, no external job queue needed
5. **[sql-migrate](https://github.com/rubenv/sql-migrate) for migrations**: Version-controlled SQL files with up/down support, production-ready
6. **ULID for IDs**: Sortable, URL-safe string IDs (26 chars) via `oklog/ulid`
7. **DDD Architecture**: Domain models/interfaces in `domain/`, GORM implementations in `adapters/repo/`

## Success Criteria

- [ ] Frontend displays real invoice list from database
- [ ] Invoice upload returns immediately with pending status
- [ ] Extraction runs in background, status updates visible
- [ ] Completed invoices show extracted data on detail page
- [ ] Failed extractions show error message
- [ ] Image files persist and serve correctly

## Dependencies

- MySQL server (local or Docker)
- Existing Gemini API integration
- Existing frontend PWA structure

## Risks

1. **Database connection issues**: Mitigate with proper error handling and health checks
2. **File storage permissions**: Test on target deployment environment
3. **Memory usage with goroutines**: Monitor and add limits if needed

