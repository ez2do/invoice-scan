# Researcher Report: API Design Patterns and Error Handling

**Researcher**: Researcher-02  
**Date**: 2025-12-02  
**Topic**: RESTful API design patterns, request validation, error handling, and response formatting

## Overview

Research on API design patterns, request/response validation, error handling strategies, and best practices for invoice extraction API.

## Key Findings

### API Design Patterns

**RESTful Principles**:
- Use HTTP methods correctly (POST for extraction)
- Consistent URL structure (`/api/extract`, `/api/health`)
- JSON request/response format
- Proper HTTP status codes
- Clear error messages

**Request Validation**:
- Validate JSON structure
- Check required fields
- Validate base64 image format
- Enforce image size limits
- Sanitize input data

**Response Format**:
- Consistent structure across endpoints
- Success/error flag
- Data payload or error message
- Optional metadata (processingTime)

### Error Handling Strategy

**Error Types**:

1. **Validation Errors (400)**:
   - Missing required fields
   - Invalid image format
   - Image too large
   - Malformed base64

2. **External API Errors (502)**:
   - Gemini API failures
   - Network timeouts
   - API quota exceeded

3. **Internal Errors (500)**:
   - Server configuration issues
   - Unexpected errors
   - Processing failures

**Error Response Format**:
```json
{
  "success": false,
  "error": "Human-readable error message",
  "code": "ERROR_CODE" // Optional
}
```

### Request/Response Models

**ExtractRequest**:
```go
type ExtractRequest struct {
    Image string `json:"image" binding:"required"`
}
```

**ExtractResponse**:
```go
type ExtractResponse struct {
    Success        bool        `json:"success"`
    Data           *InvoiceData `json:"data,omitempty"`
    Error          string      `json:"error,omitempty"`
    ProcessingTime *int64      `json:"processingTime,omitempty"`
}
```

**InvoiceData** (matches frontend):
```go
type KeyValuePair struct {
    Key        string  `json:"key"`
    Value      string  `json:"value"`
    Confidence *float64 `json:"confidence,omitempty"`
}

type TableData struct {
    Headers []string   `json:"headers"`
    Rows    [][]string `json:"rows"`
}

type InvoiceData struct {
    KeyValuePairs []KeyValuePair `json:"keyValuePairs"`
    Table         *TableData     `json:"table,omitempty"`
    Summary       []KeyValuePair `json:"summary"`
    Confidence    *float64       `json:"confidence,omitempty"`
}
```

### Validation Patterns

**Base64 Image Validation**:
- Check prefix (`data:image/...`)
- Validate base64 encoding
- Decode and verify image format
- Check image dimensions (optional)
- Enforce size limits (e.g., 10MB)

**Gin Binding Validation**:
- Use struct tags (`binding:"required"`)
- Custom validators for base64
- Return clear validation errors

### Middleware Patterns

**CORS Middleware**:
- Allow frontend origin only
- Support credentials if needed
- Configure allowed methods (POST, GET)
- Set appropriate headers

**Logging Middleware**:
- Log request method and path
- Log response status
- Log processing time
- Exclude sensitive data (images)

**Recovery Middleware**:
- Catch panics
- Return 500 error
- Log error details
- Prevent server crash

### Health Check Endpoint

**GET `/api/health`**:
- Simple status check
- Verify Gemini API connectivity (optional)
- Return server status
- Used for monitoring/load balancing

**Response**:
```json
{
  "status": "ok",
  "timestamp": "2025-12-02T23:25:00Z"
}
```

### Security Patterns

**Input Sanitization**:
- Validate all inputs
- Reject malformed data
- Limit request size
- Prevent injection attacks

**CORS Configuration**:
- Whitelist frontend origin
- Restrict methods
- Set appropriate headers
- No wildcard origins

**API Key Protection**:
- Store in environment variables
- Never expose in logs
- Rotate periodically
- Use different keys for dev/prod

### Performance Patterns

**Request Timeout**:
- Set timeout for Gemini API calls
- Default: 30 seconds
- Configurable via env

**Concurrent Processing**:
- Go goroutines for parallel requests
- Connection pooling
- Rate limiting (future)

**Response Caching**:
- Not applicable (stateless extraction)
- No caching needed for MVP

## Implementation Recommendations

1. **Structured Error Handling**:
   - Define error types
   - Centralized error handler
   - Consistent error format

2. **Request Validation**:
   - Use Gin binding
   - Custom validators
   - Clear error messages

3. **Response Formatting**:
   - Match frontend types exactly
   - Include processing time
   - Handle edge cases

4. **Logging**:
   - Structured logging
   - Log levels (info, error)
   - Exclude sensitive data

## Best Practices

- Validate early, fail fast
- Return clear error messages
- Log errors for debugging
- Handle timeouts gracefully
- Test error scenarios
- Document API contracts

## References

- REST API Best Practices: https://restfulapi.net/
- Gin Validation: https://gin-gonic.com/docs/examples/binding-and-validation/
- Go Error Handling: https://go.dev/blog/error-handling-and-go

## Unresolved Questions

1. Should health check verify Gemini API connectivity?
2. Request timeout duration (30s default)?
3. Maximum image size limit (10MB)?
4. Rate limiting strategy for MVP?
