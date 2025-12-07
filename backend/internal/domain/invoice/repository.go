package invoice

import "context"

// PaginationParams defines pagination parameters for list queries
type PaginationParams struct {
	Page     int
	PageSize int
}

// DefaultPaginationParams returns default pagination values
func DefaultPaginationParams() PaginationParams {
	return PaginationParams{
		Page:     1,
		PageSize: 10,
	}
}

// PaginatedResult contains paginated invoice results with metadata
type PaginatedResult struct {
	Invoices   Invoices
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}

type Repository interface {
	NextID() ID
	Create(ctx context.Context, invoice *Invoice) error
	GetByID(ctx context.Context, id ID) (*Invoice, error)
	List(ctx context.Context, params PaginationParams) (*PaginatedResult, error)
	Update(ctx context.Context, invoice *Invoice, updateFunc func(*Invoice) error) error
	Delete(ctx context.Context, id ID) error
}
