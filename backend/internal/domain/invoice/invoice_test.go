package invoice

import (
	"encoding/json"
	"testing"
	"time"
)

func TestID_String(t *testing.T) {
	id := ID("01HXYZ123ABC456DEF789GHI")
	if id.String() != "01HXYZ123ABC456DEF789GHI" {
		t.Errorf("Expected '01HXYZ123ABC456DEF789GHI', got '%s'", id.String())
	}
}

func TestStatus_String(t *testing.T) {
	tests := []struct {
		name     string
		status   Status
		expected string
	}{
		{"Pending", StatusPending, "pending"},
		{"Processing", StatusProcessing, "processing"},
		{"Completed", StatusCompleted, "completed"},
		{"Failed", StatusFailed, "failed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.String(); got != tt.expected {
				t.Errorf("Status.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestStatus_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		status   Status
		expected bool
	}{
		{"Pending", StatusPending, true},
		{"Processing", StatusProcessing, true},
		{"Completed", StatusCompleted, true},
		{"Failed", StatusFailed, true},
		{"Invalid", Status("invalid"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.IsValid(); got != tt.expected {
				t.Errorf("Status.IsValid() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestNew(t *testing.T) {
	id := ID("01HXYZ123ABC456DEF789GHI")
	imagePath := "/uploads/test.jpg"
	inv := New(id, imagePath)

	if inv.ID != id {
		t.Errorf("Expected ID %v, got %v", id, inv.ID)
	}
	if inv.Status != StatusPending {
		t.Errorf("Expected status %v, got %v", StatusPending, inv.Status)
	}
	if inv.ImagePath != imagePath {
		t.Errorf("Expected image path %v, got %v", imagePath, inv.ImagePath)
	}
	if inv.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
	if inv.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should not be zero")
	}
}

func TestInvoice_MarkProcessing(t *testing.T) {
	id := ID("01HXYZ123ABC456DEF789GHI")
	inv := New(id, "/uploads/test.jpg")
	beforeUpdate := inv.UpdatedAt

	time.Sleep(10 * time.Millisecond)
	inv.MarkProcessing()

	if inv.Status != StatusProcessing {
		t.Errorf("Expected status %v, got %v", StatusProcessing, inv.Status)
	}
	if !inv.UpdatedAt.After(beforeUpdate) {
		t.Error("UpdatedAt should be updated")
	}
}

func TestInvoice_MarkCompleted(t *testing.T) {
	id := ID("01HXYZ123ABC456DEF789GHI")
	inv := New(id, "/uploads/test.jpg")
	data := json.RawMessage(`{"key": "value"}`)

	inv.MarkCompleted(data)

	if inv.Status != StatusCompleted {
		t.Errorf("Expected status %v, got %v", StatusCompleted, inv.Status)
	}
	if string(inv.ExtractedData) != string(data) {
		t.Errorf("Expected extracted data %v, got %v", data, inv.ExtractedData)
	}
}

func TestInvoice_MarkFailed(t *testing.T) {
	id := ID("01HXYZ123ABC456DEF789GHI")
	inv := New(id, "/uploads/test.jpg")
	errMsg := "test error"

	inv.MarkFailed(errMsg)

	if inv.Status != StatusFailed {
		t.Errorf("Expected status %v, got %v", StatusFailed, inv.Status)
	}
	if inv.ErrorMessage == nil || *inv.ErrorMessage != errMsg {
		t.Errorf("Expected error message %v, got %v", errMsg, inv.ErrorMessage)
	}
}

