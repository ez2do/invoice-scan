package repo

import (
	"context"
	"database/sql"
	"invoice-scan/backend/pkg/ulid"
	"time"

	"invoice-scan/backend/internal/domain/invoice"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type gormInvoice struct {
	ID            string         `gorm:"column:id;primaryKey"`
	Status        string         `gorm:"column:status"`
	ImagePath     string         `gorm:"column:image_path"`
	ExtractedData datatypes.JSON `gorm:"column:extracted_data"`
	ErrorMessage  sql.NullString `gorm:"column:error_message"`
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
	return invoice.ID(ulid.GenerateULID())
}

func (r *InvoiceGormRepo) Create(ctx context.Context, inv *invoice.Invoice) error {
	var (
		db          = getDBFromContext(ctx, r.db)
		gormInvoice = r.toGorm(inv)
	)

	return db.WithContext(ctx).Create(&gormInvoice).Error
}

func (r *InvoiceGormRepo) GetByID(ctx context.Context, id invoice.ID) (*invoice.Invoice, error) {
	var gormInv gormInvoice
	err := r.db.WithContext(ctx).
		First(&gormInv, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return r.toDomain(&gormInv), nil
}

func (r *InvoiceGormRepo) List(ctx context.Context, params invoice.PaginationParams) (*invoice.PaginatedResult, error) {
	var (
		gormInvoices []gormInvoice
		total        int64
	)

	// Count total records
	if err := r.db.WithContext(ctx).Model(&gormInvoice{}).Count(&total).Error; err != nil {
		return nil, err
	}

	// Calculate offset
	offset := (params.Page - 1) * params.PageSize

	// Fetch paginated records
	if err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(params.PageSize).
		Offset(offset).
		Find(&gormInvoices).Error; err != nil {
		return nil, err
	}

	// Convert to domain models
	invoices := make([]*invoice.Invoice, len(gormInvoices))
	for i, m := range gormInvoices {
		invoices[i] = r.toDomain(&m)
	}

	// Calculate total pages
	totalPages := int(total) / params.PageSize
	if int(total)%params.PageSize > 0 {
		totalPages++
	}

	return &invoice.PaginatedResult{
		Invoices:   invoices,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (r *InvoiceGormRepo) Update(ctx context.Context, inv *invoice.Invoice, updateFunc func(invoice2 *invoice.Invoice) error) error {
	var (
		db = getDBFromContext(ctx, r.db)
	)

	if err := updateFunc(inv); err != nil {
		return err
	}

	gormInv := r.toGorm(inv)
	return db.Updates(&gormInv).Error
}

func (r *InvoiceGormRepo) Delete(ctx context.Context, id invoice.ID) error {
	db := getDBFromContext(ctx, r.db)
	err := db.WithContext(ctx).Delete(&gormInvoice{}, "id = ?", id.String()).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *InvoiceGormRepo) toGorm(inv *invoice.Invoice) *gormInvoice {
	var errorMsg sql.NullString
	if inv.ErrorMessage != nil {
		errorMsg = sql.NullString{
			String: *inv.ErrorMessage,
			Valid:  true,
		}
	}
	return &gormInvoice{
		ID:            inv.ID.String(),
		Status:        inv.Status.String(),
		ImagePath:     inv.ImagePath,
		ExtractedData: datatypes.JSON(inv.ExtractedData),
		ErrorMessage:  errorMsg,
		CreatedAt:     inv.CreatedAt,
		UpdatedAt:     inv.UpdatedAt,
	}
}

func (r *InvoiceGormRepo) toDomain(m *gormInvoice) *invoice.Invoice {
	var errorMsg *string
	if m.ErrorMessage.Valid {
		errorMsg = &m.ErrorMessage.String
	}
	return &invoice.Invoice{
		ID:            invoice.ID(m.ID),
		Status:        invoice.Status(m.Status),
		ImagePath:     m.ImagePath,
		ExtractedData: []byte(m.ExtractedData),
		ErrorMessage:  errorMsg,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}
