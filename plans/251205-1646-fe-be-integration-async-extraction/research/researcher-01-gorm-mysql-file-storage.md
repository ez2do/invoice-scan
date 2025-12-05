# Research Report: GORM/MySQL, sql-migrate, ULID & File Storage Interface Pattern

**Date**: 2025-12-05  
**Researcher**: Agent 01  
**Focus**: Backend database integration, schema migrations, ULID IDs, file storage abstraction

## 1. GORM with MySQL Setup

### Installation
```bash
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql
go get -u gorm.io/datatypes
go get -u github.com/oklog/ulid/v2
```

### Model Definition Pattern (with ULID string ID)
```go
type Invoice struct {
    ID            string         `gorm:"type:varchar(26);primaryKey"`
    Status        string         `gorm:"column:status"`
    ImagePath     string         `gorm:"column:image_path"`
    ExtractedData datatypes.JSON `gorm:"column:extracted_data"`
    ErrorMessage  *string        `gorm:"column:error_message"`
    CreatedAt     time.Time      `gorm:"column:created_at"`
    UpdatedAt     time.Time      `gorm:"column:updated_at"`
}
```

### ULID Generation
```go
import (
    "crypto/rand"
    "time"
    "github.com/oklog/ulid/v2"
)

func NextID() string {
    entropy := ulid.Monotonic(rand.Reader, 0)
    return ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
}
```

### ULID Benefits
- **Sortable**: Lexicographically sortable by creation time
- **Compact**: 26 characters vs UUID's 36
- **URL-safe**: No special characters
- **Unique**: Cryptographically random component

### Status Enum Values
- `pending`: Invoice uploaded, extraction queued
- `processing`: Extraction in progress
- `completed`: Extraction successful
- `failed`: Extraction failed

### Connection Configuration
```go
dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
```

## 2. sql-migrate for Schema Migrations

### Why sql-migrate over GORM AutoMigrate
- Version-controlled SQL migration files
- Up/down migration support for rollbacks
- Production-ready migration tracking
- Clear audit trail of schema changes

### Installation
```bash
go install github.com/rubenv/sql-migrate/...@latest
go get -u github.com/rubenv/sql-migrate
```

### Configuration File (db/dbconfig.yml)
```yaml
development:
  dialect: mysql
  datasource: root:password@tcp(localhost:3306)/invoice_scan?parseTime=true
  dir: db/migrations
  table: migrations
```

### Creating New Migrations
```bash
cd backend
sql-migrate new -config=db/dbconfig.yml create_invoices
```

### Migration File Format
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
);

-- +migrate Down
DROP TABLE IF EXISTS invoices;
```

### Running Migrations Programmatically
```go
import migrate "github.com/rubenv/sql-migrate"

func RunMigrations(db *sql.DB, migrationsDir string) (int, error) {
    migrations := &migrate.FileMigrationSource{
        Dir: migrationsDir,
    }
    return migrate.Exec(db, "mysql", migrations, migrate.Up)
}
```

### MySQL Caveat
Must append `?parseTime=true` to datasource for proper time handling.

## 3. File Storage Interface Pattern

### Interface Definition
```go
type FileStorage interface {
    Save(ctx context.Context, filename string, data []byte) (string, error)
    Get(ctx context.Context, path string) ([]byte, error)
    Delete(ctx context.Context, path string) error
    GetURL(path string) string
}
```

### Local Implementation
```go
type LocalStorage struct {
    basePath string
    baseURL  string
}

func (s *LocalStorage) Save(ctx context.Context, filename string, data []byte) (string, error) {
    path := filepath.Join(s.basePath, filename)
    return path, os.WriteFile(path, data, 0644)
}
```

### Future S3 Implementation Signature
```go
type S3Storage struct {
    client     *s3.Client
    bucketName string
    region     string
}
```

## 4. DDD Architecture for Backend

### Directory Structure
```
backend/internal/
├── domain/                     # Domain layer (business logic)
│   ├── invoice/
│   │   ├── invoice.go         # Entity + types (ID, Status)
│   │   └── repository.go      # Repository interface
│   └── storage/
│       └── repository.go      # FileStorage interface
├── adapters/                   # Infrastructure layer
│   ├── repo/
│   │   └── invoice_gorm_repo.go  # gormInvoice model
│   └── storage/
│       └── local_storage.go
```

### Domain Types (invoice.go)
```go
// internal/domain/invoice/invoice.go
type ID string
type Status string

const (
    StatusPending    Status = "pending"
    StatusProcessing Status = "processing"
    StatusCompleted  Status = "completed"
    StatusFailed     Status = "failed"
)
```

### Repository Interface (Domain Layer)
```go
// internal/domain/invoice/repository.go
type Repository interface {
    NextID() ID
    Create(ctx context.Context, invoice *Invoice) error
    GetByID(ctx context.Context, id ID) (*Invoice, error)
    List(ctx context.Context) ([]*Invoice, error)
    Update(ctx context.Context, invoice *Invoice) error
}
```

### GORM Implementation (Adapters Layer)
```go
// internal/adapters/repo/invoice_gorm_repo.go
type gormInvoice struct {
    ID            string         `gorm:"type:varchar(26);primaryKey"`
    Status        string         `gorm:"column:status"`
    // ...
}

func (gormInvoice) TableName() string {
    return "invoices"
}

type InvoiceGormRepo struct {
    db *gorm.DB
}

func (r *InvoiceGormRepo) NextID() invoice.ID {
    entropy := ulid.Monotonic(rand.Reader, 0)
    return invoice.ID(ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String())
}
```

## 5. Best Practices

### Database Migrations
- Use sql-migrate for all schema changes
- Create migrations with `sql-migrate new -config=db/dbconfig.yml {name}`
- Always write both Up and Down migrations
- Test rollbacks in development

### JSON Storage
- Use `datatypes.JSON` from gorm.io/datatypes
- Allows flexible schema for varying invoice formats

### File Naming
- Use ULID for unique filenames
- Format: `{ulid}.{ext}`

## 6. Key Considerations

1. **Transaction Support**: GORM supports database transactions
2. **Indexes**: Add indexes on `status` and `created_at` via migrations
3. **Connection Pool**: Configure GORM connection pool for production
4. **Migration Safety**: sql-migrate tracks applied migrations in `migrations` table

## References
- [GORM Documentation](https://gorm.io/docs/)
- [GORM MySQL Driver](https://gorm.io/docs/connecting_to_the_database.html#MySQL)
- [sql-migrate GitHub](https://github.com/rubenv/sql-migrate)
- [ULID Spec](https://github.com/ulid/spec)
- [oklog/ulid Go Library](https://github.com/oklog/ulid)
