package invoice

import "context"

type Repository interface {
	NextID() ID
	Create(ctx context.Context, invoice *Invoice) error
	GetByID(ctx context.Context, id ID) (*Invoice, error)
	List(ctx context.Context) (Invoices, error)
	Update(ctx context.Context, invoice *Invoice) error
}
