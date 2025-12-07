package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "Image file is required: " + err.Error(),
		})
		return
	}

	if file.Size == 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "Image file is empty",
		})
		return
	}

	const maxSize = 10 * 1024 * 1024
	if file.Size > maxSize {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "Image file too large (max 10MB)",
		})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "Failed to open image file: " + err.Error(),
		})
		return
	}
	defer src.Close()

	imageBytes, err := io.ReadAll(src)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "Failed to read image file: " + err.Error(),
		})
		return
	}

	if len(imageBytes) == 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "Image file is empty",
		})
		return
	}

	mimeType := file.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = http.DetectContentType(imageBytes)
	}

	if !strings.HasPrefix(mimeType, "image/") {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "Invalid file type. Only images are allowed",
		})
		return
	}

	id := h.repo.NextID()
	filename := fmt.Sprintf("%s%s", id.String(), filepath.Ext(file.Filename))

	imagePath, err := h.storage.Save(c.Request.Context(), filename, imageBytes, mimeType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error:   "Failed to save image: " + err.Error(),
		})
		return
	}

	inv := invoice.New(id, imagePath)
	if err := h.repo.Create(c.Request.Context(), inv); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error:   "Failed to create invoice: " + err.Error(),
		})
		return
	}

	h.processExtractionAsync(id, imageBytes, mimeType)

	data := NewInvoiceData(inv, getImagePath(inv.ImagePath))

	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Data:    data,
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
	// Parse pagination params with defaults
	params := invoice.DefaultPaginationParams()

	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			params.Page = page
		}
	}

	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 && pageSize <= 100 {
			params.PageSize = pageSize
		}
	}

	result, err := h.repo.List(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error:   "Failed to list invoices: " + err.Error(),
		})
		return
	}

	data := make([]InvoiceData, len(result.Invoices))
	for i, inv := range result.Invoices {
		data[i] = NewInvoiceData(inv, getImagePath(inv.ImagePath))
	}

	c.JSON(http.StatusOK, PaginatedInvoicesResponse{
		Success:    true,
		Data:       data,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	})
}

func (h *InvoiceHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id := invoice.ID(idStr)

	inv, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Success: false,
			Error:   "Invoice not found",
		})
		return
	}

	data := NewInvoiceData(inv, getImagePath(inv.ImagePath))

	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Data:    data,
	})
}

func (h *InvoiceHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id := invoice.ID(idStr)

	inv, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Success: false,
			Error:   "Invoice not found",
		})
		return
	}

	// Delete image file
	if inv.ImagePath != "" {
		if err := h.storage.Delete(c.Request.Context(), inv.ImagePath); err != nil {
			log.Printf("Failed to delete image file %s: %v", inv.ImagePath, err)
		}
	}

	// Delete database record
	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error:   "Failed to delete invoice: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Data:    nil,
	})
}

func (h *InvoiceHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id := invoice.ID(idStr)

	var req UpdateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
		return
	}

	inv, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Success: false,
			Error:   "Invoice not found",
		})
		return
	}

	if err := h.repo.Update(c.Request.Context(), inv, func(i *invoice.Invoice) error {
		i.ExtractedData = req.ExtractedData
		i.UpdatedAt = time.Now()
		return nil
	}); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error:   "Failed to update invoice: " + err.Error(),
		})
		return
	}

	data := NewInvoiceData(inv, getImagePath(inv.ImagePath))

	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Data:    data,
	})
}


func getImagePath(fullPath string) string {
	filename := filepath.Base(fullPath)
	return fmt.Sprintf("/uploads/%s", filename)
}
