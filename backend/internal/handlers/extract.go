package handlers

import (
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

	var req models.ExtractRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ExtractResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	if req.Image == "" {
		c.JSON(http.StatusBadRequest, models.ExtractResponse{
			Success: false,
			Error:   "Image field is required",
		})
		return
	}

	invoiceData, err := h.extractionService.Extract(c.Request.Context(), req.Image)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "Gemini API") {
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
