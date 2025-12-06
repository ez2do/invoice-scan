package handlers

import (
	"encoding/json"
	"invoice-scan/backend/internal/domain/invoice"
)

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

type InvoiceData struct {
	ID            string      `json:"id"`
	Status        string      `json:"status"`
	ImagePath     string      `json:"image_path"`
	CreatedAt     string      `json:"created_at"`
	UpdatedAt     string      `json:"updated_at,omitempty"`
	ExtractedData interface{} `json:"extracted_data,omitempty"`
	ErrorMessage  *string     `json:"error_message,omitempty"`
}

func NewInvoiceData(inv *invoice.Invoice, imageURL string) InvoiceData {
	data := InvoiceData{
		ID:        inv.ID.String(),
		Status:    inv.Status.String(),
		ImagePath: imageURL,
		CreatedAt: inv.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if !inv.UpdatedAt.IsZero() {
		data.UpdatedAt = inv.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
	}

	if inv.Status == invoice.StatusCompleted && len(inv.ExtractedData) > 0 {
		var extractedData interface{}
		if err := json.Unmarshal(inv.ExtractedData, &extractedData); err == nil {
			data.ExtractedData = extractedData
		}
	}

	if inv.ErrorMessage != nil {
		data.ErrorMessage = inv.ErrorMessage
	}

	return data
}
