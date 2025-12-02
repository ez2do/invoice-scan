package models

type ExtractResponse struct {
	Success        bool         `json:"success"`
	Data           *InvoiceData `json:"data,omitempty"`
	Error          string       `json:"error,omitempty"`
	ProcessingTime *int64       `json:"processingTime,omitempty"`
}
