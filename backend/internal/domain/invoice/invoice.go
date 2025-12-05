package invoice

import (
	"encoding/json"
	"time"
)

type ID string

func (id ID) String() string {
	return string(id)
}

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

type (
	Invoice struct {
		ID            ID
		Status        Status
		ImagePath     string
		ExtractedData json.RawMessage
		ErrorMessage  *string
		CreatedAt     time.Time
		UpdatedAt     time.Time
	}

	Invoices []*Invoice
)

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
