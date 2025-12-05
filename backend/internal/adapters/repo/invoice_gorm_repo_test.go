package repo

import (
	"context"
	"testing"
	"time"

	"invoice-scan/backend/internal/domain/invoice"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	err = db.Exec(`
		CREATE TABLE invoices (
			id VARCHAR(26) NOT NULL PRIMARY KEY,
			status VARCHAR(20) NOT NULL DEFAULT 'pending',
			image_path VARCHAR(500) NOT NULL,
			extracted_data TEXT,
			error_message TEXT,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)
	`).Error
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	return db
}

func TestInvoiceGormRepo_NextID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewInvoiceGormRepo(db)

	id := repo.NextID()
	if id.String() == "" {
		t.Error("NextID() should return a non-empty ID")
	}
	if len(id.String()) != 26 {
		t.Errorf("Expected ULID length 26, got %d", len(id.String()))
	}
}

func TestInvoiceGormRepo_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewInvoiceGormRepo(db)
	ctx := context.Background()

	id := repo.NextID()
	inv := invoice.New(id, "/uploads/test.jpg")

	err := repo.Create(ctx, inv)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	var count int64
	db.Model(&gormInvoice{}).Count(&count)
	if count != 1 {
		t.Errorf("Expected 1 invoice, got %d", count)
	}
}

func TestInvoiceGormRepo_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewInvoiceGormRepo(db)
	ctx := context.Background()

	id := repo.NextID()
	inv := invoice.New(id, "/uploads/test.jpg")
	repo.Create(ctx, inv)

	got, err := repo.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}

	if got.ID != id {
		t.Errorf("Expected ID %v, got %v", id, got.ID)
	}
	if got.ImagePath != "/uploads/test.jpg" {
		t.Errorf("Expected image path '/uploads/test.jpg', got '%s'", got.ImagePath)
	}
}

func TestInvoiceGormRepo_GetByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewInvoiceGormRepo(db)
	ctx := context.Background()

	id := invoice.ID("01HXYZ123ABC456DEF789GHI")
	_, err := repo.GetByID(ctx, id)
	if err == nil {
		t.Error("Expected error for non-existent invoice")
	}
}

func TestInvoiceGormRepo_List(t *testing.T) {
	db := setupTestDB(t)
	repo := NewInvoiceGormRepo(db)
	ctx := context.Background()

	id1 := repo.NextID()
	inv1 := invoice.New(id1, "/uploads/test1.jpg")
	repo.Create(ctx, inv1)

	time.Sleep(10 * time.Millisecond)

	id2 := repo.NextID()
	inv2 := invoice.New(id2, "/uploads/test2.jpg")
	repo.Create(ctx, inv2)

	invoices, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(invoices) != 2 {
		t.Errorf("Expected 2 invoices, got %d", len(invoices))
	}

	if invoices[0].ID != id2 {
		t.Error("List should return invoices in descending order by created_at")
	}
}

func TestInvoiceGormRepo_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewInvoiceGormRepo(db)
	ctx := context.Background()

	id := repo.NextID()
	inv := invoice.New(id, "/uploads/test.jpg")
	repo.Create(ctx, inv)

	inv.MarkProcessing()
	err := repo.Update(ctx, inv)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	got, _ := repo.GetByID(ctx, id)
	if got.Status != invoice.StatusProcessing {
		t.Errorf("Expected status %v, got %v", invoice.StatusProcessing, got.Status)
	}
}

func TestInvoiceGormRepo_toGorm(t *testing.T) {
	db := setupTestDB(t)
	repo := NewInvoiceGormRepo(db)

	id := invoice.ID("01HXYZ123ABC456DEF789GHI")
	inv := invoice.New(id, "/uploads/test.jpg")
	inv.MarkCompleted([]byte(`{"key": "value"}`))

	model := repo.toGorm(inv)
	if model.ID != id.String() {
		t.Errorf("Expected ID %v, got %v", id.String(), model.ID)
	}
	if model.Status != "completed" {
		t.Errorf("Expected status 'completed', got '%s'", model.Status)
	}
}

func TestInvoiceGormRepo_toDomain(t *testing.T) {
	db := setupTestDB(t)
	repo := NewInvoiceGormRepo(db)

	model := &gormInvoice{
		ID:            "01HXYZ123ABC456DEF789GHI",
		Status:        "completed",
		ImagePath:     "/uploads/test.jpg",
		ExtractedData: []byte(`{"key": "value"}`),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	inv := repo.toDomain(model)
	if inv.ID.String() != model.ID {
		t.Errorf("Expected ID %v, got %v", model.ID, inv.ID.String())
	}
	if inv.Status.String() != model.Status {
		t.Errorf("Expected status %v, got %v", model.Status, inv.Status.String())
	}
}

