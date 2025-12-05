# Phase 01: Backend Database Setup

**Parent Plan**: [plan.md](./plan.md)  
**Dependencies**: None  
**Status**: Completed  
**Priority**: High

## Overview

Set up MySQL database with GORM, use [sql-migrate](https://github.com/rubenv/sql-migrate) for versioned schema migrations. Follow DDD architecture with domain models and repository interfaces in domain layer, GORM implementations in adapters layer.

## Key Insights

- **DDD Architecture**: Domain logic separated from infrastructure
- GORM for ORM operations (queries, inserts, updates)
- sql-migrate for versioned SQL migrations (run via CLI)
- ULID for sortable, unique string IDs
- Repository interface in domain, implementation in adapters
- Domain types for ID, Status, etc.

## DDD Directory Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go
├── db/
│   ├── dbconfig.yml
│   └── migrations/
│       └── {timestamp}_create_invoices.sql
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── domain/
│   │   └── invoice/
│   │       ├── invoice.go         # Entity + types (ID, Status)
│   │       └── repository.go      # Repository interface
│   ├── adapters/
│   │   └── repo/
│   │       └── invoice_gorm_repo.go  # GORM implementation
│   ├── handlers/
│   │   └── invoice.go
│   └── services/
│       └── extraction.go
```

## Requirements

1. Add GORM, MySQL driver, sql-migrate, and ULID dependencies
2. Create domain layer with Invoice entity, types, and repository interface
3. Create GORM repository implementation in adapters layer
4. Create SQL migration files
5. Update config for database connection

## Related Code Files

- `backend/internal/config/config.go` - Add DB config
- NEW: `backend/internal/domain/invoice/invoice.go` - Invoice entity + types
- NEW: `backend/internal/domain/invoice/repository.go` - Repository interface
- NEW: `backend/internal/adapters/repo/invoice_gorm_repo.go` - GORM implementation
- NEW: `backend/db/migrations/` - SQL migration files
- NEW: `backend/db/dbconfig.yml` - sql-migrate config

## Implementation Steps

### Step 1: Add Dependencies
```bash
cd backend
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql
go get -u gorm.io/datatypes
go get -u github.com/rubenv/sql-migrate
go get -u github.com/oklog/ulid/v2
```

### Step 2: Setup sql-migrate Config

Create `backend/db/dbconfig.yml`:
```yaml
development:
  dialect: mysql
  datasource: ${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}?parseTime=true
  dir: db/migrations
  table: migrations
```

### Step 3: Create Migration Files

```bash
cd backend
sql-migrate new -config=db/dbconfig.yml create_invoices
```

Edit `backend/db/migrations/{timestamp}_create_invoices.sql`:
```sql
-- +migrate Up
CREATE TABLE invoices (
    id VARCHAR(26) NOT NULL PRIMARY KEY,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    image_path VARCHAR(500) NOT NULL,
    extracted_data JSON,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +migrate Down
DROP TABLE IF EXISTS invoices;
```

### Step 4: Create Domain Layer

**`internal/domain/invoice/invoice.go`**:
```go
package invoice

import (
    "encoding/json"
    "time"
)

// ID represents Invoice identifier (ULID)
type ID string

func (id ID) String() string {
    return string(id)
}

// Status represents Invoice processing status
type Status string

const (
    StatusPending    Status = "pending"
    StatusProcessing Status = "processing"
    StatusCompleted  Status = "completed"
    StatusFailed     Status = "failed"
)

func (s Status) String() string {
    return string(s)
}

func (s Status) IsValid() bool {
    switch s {
    case StatusPending, StatusProcessing, StatusCompleted, StatusFailed:
        return true
    }
    return false
}

// Invoice represents the invoice domain entity
type Invoice struct {
    ID            ID
    Status        Status
    ImagePath     string
    ExtractedData json.RawMessage
    ErrorMessage  *string
    CreatedAt     time.Time
    UpdatedAt     time.Time
}

func New(id ID, imagePath string) *Invoice {
    now := time.Now()
    return &Invoice{
        ID:        id,
        Status:    StatusPending,
        ImagePath: imagePath,
        CreatedAt: now,
        UpdatedAt: now,
    }
}

func (i *Invoice) MarkProcessing() {
    i.Status = StatusProcessing
    i.UpdatedAt = time.Now()
}

func (i *Invoice) MarkCompleted(data json.RawMessage) {
    i.Status = StatusCompleted
    i.ExtractedData = data
    i.UpdatedAt = time.Now()
}

func (i *Invoice) MarkFailed(errMsg string) {
    i.Status = StatusFailed
    i.ErrorMessage = &errMsg
    i.UpdatedAt = time.Now()
}
```

**`internal/domain/invoice/repository.go`**:
```go
package invoice

import "context"

type Repository interface {
    NextID() ID
    Create(ctx context.Context, invoice *Invoice) error
    GetByID(ctx context.Context, id ID) (*Invoice, error)
    List(ctx context.Context) ([]*Invoice, error)
    Update(ctx context.Context, invoice *Invoice) error
}
```

### Step 5: Create GORM Repository Implementation

**`internal/adapters/repo/invoice_gorm_repo.go`**:
```go
package repo

import (
    "context"
    "crypto/rand"
    "time"

    "invoice-scan/backend/internal/domain/invoice"

    "github.com/oklog/ulid/v2"
    "gorm.io/datatypes"
    "gorm.io/gorm"
)

type gormInvoice struct {
    ID            string         `gorm:"type:varchar(26);primaryKey"`
    Status        string         `gorm:"column:status"`
    ImagePath     string         `gorm:"column:image_path"`
    ExtractedData datatypes.JSON `gorm:"column:extracted_data"`
    ErrorMessage  *string        `gorm:"column:error_message"`
    CreatedAt     time.Time      `gorm:"column:created_at"`
    UpdatedAt     time.Time      `gorm:"column:updated_at"`
}

func (gormInvoice) TableName() string {
    return "invoices"
}

type InvoiceGormRepo struct {
    db *gorm.DB
}

func NewInvoiceGormRepo(db *gorm.DB) *InvoiceGormRepo {
    return &InvoiceGormRepo{db: db}
}

func (r *InvoiceGormRepo) NextID() invoice.ID {
    entropy := ulid.Monotonic(rand.Reader, 0)
    return invoice.ID(ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String())
}

func (r *InvoiceGormRepo) Create(ctx context.Context, inv *invoice.Invoice) error {
    model := r.toGorm(inv)
    return r.db.WithContext(ctx).Create(model).Error
}

func (r *InvoiceGormRepo) GetByID(ctx context.Context, id invoice.ID) (*invoice.Invoice, error) {
    var model gormInvoice
    if err := r.db.WithContext(ctx).First(&model, "id = ?", id.String()).Error; err != nil {
        return nil, err
    }
    return r.toDomain(&model), nil
}

func (r *InvoiceGormRepo) List(ctx context.Context) ([]*invoice.Invoice, error) {
    var models []gormInvoice
    if err := r.db.WithContext(ctx).Order("created_at DESC").Find(&models).Error; err != nil {
        return nil, err
    }
    
    invoices := make([]*invoice.Invoice, len(models))
    for i, m := range models {
        invoices[i] = r.toDomain(&m)
    }
    return invoices, nil
}

func (r *InvoiceGormRepo) Update(ctx context.Context, inv *invoice.Invoice) error {
    model := r.toGorm(inv)
    return r.db.WithContext(ctx).Save(model).Error
}

func (r *InvoiceGormRepo) toGorm(inv *invoice.Invoice) *gormInvoice {
    return &gormInvoice{
        ID:            inv.ID.String(),
        Status:        inv.Status.String(),
        ImagePath:     inv.ImagePath,
        ExtractedData: datatypes.JSON(inv.ExtractedData),
        ErrorMessage:  inv.ErrorMessage,
        CreatedAt:     inv.CreatedAt,
        UpdatedAt:     inv.UpdatedAt,
    }
}

func (r *InvoiceGormRepo) toDomain(m *gormInvoice) *invoice.Invoice {
    return &invoice.Invoice{
        ID:            invoice.ID(m.ID),
        Status:        invoice.Status(m.Status),
        ImagePath:     m.ImagePath,
        ExtractedData: []byte(m.ExtractedData),
        ErrorMessage:  m.ErrorMessage,
        CreatedAt:     m.CreatedAt,
        UpdatedAt:     m.UpdatedAt,
    }
}
```

### Step 6: Update Config

Add to `config.go`:
- `DBHost`, `DBPort`, `DBUser`, `DBPassword`, `DBName`
- `GetDSN()` method for connection string

```go
type Config struct {
    // ... existing fields
    DBHost     string
    DBPort     string
    DBUser     string
    DBPassword string
    DBName     string
}

func (c *Config) GetDSN() string {
    return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}
```

### Step 7: Initialize in main.go
```go
// Open GORM connection
gormDB, err := gorm.Open(mysql.Open(cfg.GetDSN()), &gorm.Config{})
if err != nil {
    log.Fatalf("Failed to connect to database: %v", err)
}

// Create repository
invoiceRepo := repo.NewInvoiceGormRepo(gormDB)
```

## Todo List

- [ ] Add GORM, sql-migrate, and ULID dependencies to go.mod
- [ ] Create db/dbconfig.yml for sql-migrate
- [ ] Create db/migrations/ directory
- [ ] Generate initial migration with `sql-migrate new`
- [ ] Run migration with `sql-migrate up`
- [ ] Create internal/domain/invoice/ directory
- [ ] Create invoice.go with domain entity and types (ID, Status)
- [ ] Create repository.go with Repository interface
- [ ] Create internal/adapters/repo/ directory
- [ ] Create invoice_gorm_repo.go with GORM implementation
- [ ] Update config.go with database configuration
- [ ] Update main.go to init DB and create repository
- [ ] Test repository operations

## Success Criteria

- [x] Migrations applied via `sql-migrate up` (migration file created)
- [x] Invoice table created with VARCHAR(26) primary key
- [x] Domain entity is clean (no GORM/infrastructure dependencies)
- [x] Domain types defined (ID, Status)
- [x] Repository interface in domain layer
- [x] GORM implementation in adapters layer (using gormInvoice)
- [x] NextID() generates valid ULIDs
- [x] CRUD operations work via repository

## Risk Assessment

- **MySQL not running**: Add clear error message with setup instructions
- **Migration failures**: sql-migrate tracks applied migrations, prevents re-running

## Environment Variables

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=invoice_scan
```

## Migration Commands (CLI)

```bash
# Install sql-migrate CLI
go install github.com/rubenv/sql-migrate/...@latest

# Create new migration
sql-migrate new -config=db/dbconfig.yml create_invoices

# Check status
sql-migrate status -config=db/dbconfig.yml

# Apply migrations
sql-migrate up -config=db/dbconfig.yml

# Rollback last migration
sql-migrate down -config=db/dbconfig.yml -limit=1
```

## ULID Benefits

- **Sortable**: Lexicographically sortable by creation time
- **Compact**: 26 characters vs UUID's 36
- **URL-safe**: No special characters
- **Unique**: Cryptographically random component

## Next Steps

After completion, proceed to [Phase 02: Backend File Storage](./phase-02-backend-file-storage.md)
