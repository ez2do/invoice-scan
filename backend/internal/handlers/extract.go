package handlers

import (
	"io"
	"net/http"
	"strings"
	"time"

	"invoice-scan/backend/internal/models"
	"invoice-scan/backend/internal/services"

	"github.com/gin-gonic/gin"
)

type ExtractHandler struct {
	extractionService *services.ExtractionService
}

func NewExtractHandler(extractionService *services.ExtractionService) *ExtractHandler {
	return &ExtractHandler{
		extractionService: extractionService,
	}
}

func (h *ExtractHandler) Extract(c *gin.Context) {
	startTime := time.Now()

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ExtractResponse{
			Success: false,
			Error:   "Image file is required: " + err.Error(),
		})
		return
	}

	if file.Size == 0 {
		c.JSON(http.StatusBadRequest, models.ExtractResponse{
			Success: false,
			Error:   "Image file is empty",
		})
		return
	}

	const maxSize = 10 * 1024 * 1024
	if file.Size > maxSize {
		c.JSON(http.StatusBadRequest, models.ExtractResponse{
			Success: false,
			Error:   "Image file too large (max 10MB)",
		})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ExtractResponse{
			Success: false,
			Error:   "Failed to open image file: " + err.Error(),
		})
		return
	}
	defer src.Close()

	imageBytes, err := io.ReadAll(src)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ExtractResponse{
			Success: false,
			Error:   "Failed to read image file: " + err.Error(),
		})
		return
	}

	if len(imageBytes) == 0 {
		c.JSON(http.StatusBadRequest, models.ExtractResponse{
			Success: false,
			Error:   "Image file is empty",
		})
		return
	}

	mimeType := file.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = http.DetectContentType(imageBytes)
	}

	invoiceData, err := h.extractionService.Extract(c.Request.Context(), imageBytes, mimeType)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "gemini API") || strings.Contains(err.Error(), "Gemini API") {
			statusCode = http.StatusBadGateway
		} else if strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "context deadline exceeded") {
			statusCode = http.StatusGatewayTimeout
		} else if strings.Contains(err.Error(), "invalid") {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, models.ExtractResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	processingTime := time.Since(startTime).Milliseconds()

	c.JSON(http.StatusOK, models.ExtractResponse{
		Success:        true,
		Data:           invoiceData,
		ProcessingTime: &processingTime,
	})
}
