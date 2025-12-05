package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"invoice-scan/backend/internal/domain/invoice"
	domainstorage "invoice-scan/backend/internal/domain/storage"

	"github.com/gin-gonic/gin"
)

type InvoiceHandler struct {
	repo              invoice.Repository
	storage           domainstorage.FileStorage
	extractionService invoice.ExtractionService
}

func NewInvoiceHandler(
	repo invoice.Repository,
	storage domainstorage.FileStorage,
	extractionService invoice.ExtractionService,
) *InvoiceHandler {
	return &InvoiceHandler{
		repo:              repo,
		storage:           storage,
		extractionService: extractionService,
	}
}

func (h *InvoiceHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Image file is required: " + err.Error(),
		})
		return
	}

	if file.Size == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Image file is empty",
		})
		return
	}

	const maxSize = 10 * 1024 * 1024
	if file.Size > maxSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Image file too large (max 10MB)",
		})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Failed to open image file: " + err.Error(),
		})
		return
	}
	defer src.Close()

	imageBytes, err := io.ReadAll(src)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Failed to read image file: " + err.Error(),
		})
		return
	}

	if len(imageBytes) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Image file is empty",
		})
		return
	}

	mimeType := file.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = http.DetectContentType(imageBytes)
	}

	if !strings.HasPrefix(mimeType, "image/") {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid file type. Only images are allowed",
		})
		return
	}

	id := h.repo.NextID()
	filename := fmt.Sprintf("%s%s", id.String(), filepath.Ext(file.Filename))

	imagePath, err := h.storage.Save(c.Request.Context(), filename, imageBytes, mimeType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to save image: " + err.Error(),
		})
		return
	}

	inv := invoice.New(id, imagePath)
	if err := h.repo.Create(c.Request.Context(), inv); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to create invoice: " + err.Error(),
		})
		return
	}

	h.processExtractionAsync(id, imageBytes, mimeType)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"id":        inv.ID.String(),
			"status":    inv.Status.String(),
			"imagePath": h.storage.GetURL(inv.ImagePath),
			"createdAt": inv.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	})
}

func (h *InvoiceHandler) processExtractionAsync(invoiceID invoice.ID, imageBytes []byte, mimeType string) {
	go func() {
		ctx := context.Background()

		inv, err := h.repo.GetByID(ctx, invoiceID)
		if err != nil {
			log.Printf("Failed to get invoice %s: %v", invoiceID.String(), err)
			return
		}

		if err := h.repo.Update(ctx, inv, func(i *invoice.Invoice) error {
			i.MarkProcessing()
			return nil
		}); err != nil {
			log.Printf("Failed to update invoice %s to processing: %v", invoiceID.String(), err)
			return
		}

		data, err := h.extractionService.Extract(ctx, imageBytes, mimeType)
		if err != nil {
			if updateErr := h.repo.Update(ctx, inv, func(i *invoice.Invoice) error {
				i.MarkFailed(err.Error())
				return nil
			}); updateErr != nil {
				log.Printf("Failed to update invoice %s to failed: %v", invoiceID.String(), updateErr)
			}
			return
		}

		dataJSON, err := json.Marshal(data)
		if err != nil {
			if updateErr := h.repo.Update(ctx, inv, func(i *invoice.Invoice) error {
				i.MarkFailed("Failed to marshal extracted data: " + err.Error())
				return nil
			}); updateErr != nil {
				log.Printf("Failed to update invoice %s to failed: %v", invoiceID.String(), updateErr)
			}
			return
		}

		if err := h.repo.Update(ctx, inv, func(i *invoice.Invoice) error {
			i.MarkCompleted(dataJSON)
			return nil
		}); err != nil {
			log.Printf("Failed to update invoice %s to completed: %v", invoiceID.String(), err)
		}
	}()
}

func (h *InvoiceHandler) List(c *gin.Context) {
	invoices, err := h.repo.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to list invoices: " + err.Error(),
		})
		return
	}

	data := make([]gin.H, len(invoices))
	for i, inv := range []*invoice.Invoice(invoices) {
		item := gin.H{
			"id":        inv.ID.String(),
			"status":    inv.Status.String(),
			"imagePath": h.storage.GetURL(inv.ImagePath),
			"createdAt": inv.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			"updatedAt": inv.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}

		if inv.Status == invoice.StatusCompleted && len(inv.ExtractedData) > 0 {
			var extractedData interface{}
			if err := json.Unmarshal(inv.ExtractedData, &extractedData); err == nil {
				item["extractedData"] = extractedData
			}
		}

		if inv.ErrorMessage != nil {
			item["errorMessage"] = *inv.ErrorMessage
		}

		data[i] = item
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}

func (h *InvoiceHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id := invoice.ID(idStr)

	inv, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Invoice not found",
		})
		return
	}

	data := gin.H{
		"id":        inv.ID.String(),
		"status":    inv.Status.String(),
		"imagePath": h.storage.GetURL(inv.ImagePath),
		"createdAt": inv.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		"updatedAt": inv.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if len(inv.ExtractedData) > 0 {
		var extractedData interface{}
		if err := json.Unmarshal(inv.ExtractedData, &extractedData); err == nil {
			data["extractedData"] = extractedData
		}
	}

	if inv.ErrorMessage != nil {
		data["errorMessage"] = *inv.ErrorMessage
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}
