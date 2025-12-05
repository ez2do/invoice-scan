package httputil

import (
	"invoice-scan/backend/pkg/errors"
	"invoice-scan/backend/pkg/locale"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GinResponse provides common response structures for Gin handlers
// All API responses follow this consistent format:
//
// Success Response:
//
//	{
//	  "success": true,
//	  "message": "optional success message",
//	  "data": { ... actual response data ... }
//	}
//
// Error Response:
//
//	{
//	  "success": false,
//	  "message": "error description",
//	  "detail": { ... optional error details ... }
//	}
type GinResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Detail  interface{} `json:"detail,omitempty"`
}

// SuccessResponse sends a successful response with data
func SuccessResponse(c *gin.Context, data interface{}, message ...string) {
	msg := ""
	if len(message) > 0 {
		msg = message[0]
	}

	c.JSON(http.StatusOK, GinResponse{
		Success: true,
		Message: msg,
		Data:    data,
	})
}

// CreatedResponse sends a 201 Created response with data
func CreatedResponse(c *gin.Context, data interface{}, message ...string) {
	msg := ""
	if len(message) > 0 {
		msg = message[0]
	}

	c.JSON(http.StatusCreated, GinResponse{
		Success: true,
		Message: msg,
		Data:    data,
	})
}

// ErrorResponse sends an error response with a custom message
func ErrorResponse(c *gin.Context, statusCode int, message string, detail ...interface{}) {
	response := GinResponse{
		Success: false,
		Message: locale.TL(c.GetHeader(HeaderContentLanguage), message),
	}

	if len(detail) > 0 {
		response.Detail = detail[0]
	}

	c.JSON(statusCode, response)
}

// BadRequestResponse sends a 400 Bad Request response
func BadRequestResponse(c *gin.Context, message string, detail ...interface{}) {
	ErrorResponse(c, http.StatusBadRequest, message, detail...)
}

// UnauthorizedResponse sends a 401 Unauthorized response
func UnauthorizedResponse(c *gin.Context, message string, detail ...interface{}) {
	ErrorResponse(c, http.StatusUnauthorized, message, detail...)
}

// ForbiddenResponse sends a 403 Forbidden response
func ForbiddenResponse(c *gin.Context, message string, detail ...interface{}) {
	ErrorResponse(c, http.StatusForbidden, message, detail...)
}

// NotFoundResponse sends a 404 Not Found response
func NotFoundResponse(c *gin.Context, message string, detail ...interface{}) {
	ErrorResponse(c, http.StatusNotFound, message, detail...)
}

// InternalServerErrorResponse sends a 500 Internal Server Error response
func InternalServerErrorResponse(c *gin.Context, message string, detail ...interface{}) {
	ErrorResponse(c, http.StatusInternalServerError, message, detail...)
}

// ValidationErrorResponse sends a 422 Unprocessable Entity response for validation errors
func ValidationErrorResponse(c *gin.Context, message string, detail ...interface{}) {
	ErrorResponse(c, http.StatusUnprocessableEntity, message, detail...)
}

// IErrorResponse sends an error response using IError instance
func IErrorResponse(c *gin.Context, ierr *errors.IError) {
	// Set error headers
	c.Header(XErrorIDHeader, ierr.ID())
	c.Header(XErrorCodeHeader, ierr.Code())

	errDetail := translateIError(ierr, c.GetHeader(HeaderContentLanguage))

	c.JSON(ierr.HTTPCode(), GinResponse{
		Success: false,
		Message: errDetail["error_summary"].(string),
		Detail:  errDetail,
	})
}

// BindingErrorResponse handles JSON binding errors with detailed validation messages
func BindingErrorResponse(c *gin.Context, err error) {
	BadRequestResponse(c, "Invalid request format", gin.H{
		"binding_error": err.Error(),
	})
}
