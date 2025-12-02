# Phase 2: Extract Endpoint Implementation

**Date**: 2025-12-02  
**Priority**: High  
**Status**: Pending  
**Estimated Time**: 3-4 hours

## Context Links

- **Parent Plan**: `plan.md`
- **Dependencies**: Phase 1 (Project Setup)
- **Research**: 
  - `research/researcher-01-go-gin-gemini.md`
  - `research/researcher-02-api-design-patterns.md`
- **Reference Docs**:
  - `frontend/src/lib/api.ts` - Frontend API client
  - `frontend/src/types/index.ts` - Type definitions
  - `docs/system-architecture.md` - API design

## Overview

Implement POST `/api/extract` endpoint that accepts base64-encoded invoice images, validates input, calls Gemini Vision API, parses response, and returns structured invoice data matching frontend expectations.

## Key Insights

- Frontend sends base64 image in JSON body: `{ image: string }`
- Response must match `ExtractResponse` type exactly
- Gemini API requires proper prompt engineering for Vietnamese invoices
- Need to handle base64 decoding and image validation
- Processing time should be included in response

## Requirements

### Functional Requirements
- Accept POST requests to `/api/extract`
- Validate base64 image format
- Decode base64 to image bytes
- Call Gemini Vision API with proper prompt
- Parse Gemini response to InvoiceData structure
- Return JSON response matching frontend types
- Include processing time in response

### Non-Functional Requirements
- Request timeout: 30 seconds
- Max image size: 10MB
- Error handling for all failure cases
- Logging (without sensitive data)

## Architecture

```
Request Flow:
Frontend → Gin Handler → Request Validation → Extraction Service → Gemini API → Response Parsing → Frontend
```

**Components**:
- Handler: HTTP request/response handling
- Service: Business logic and Gemini integration
- Models: Data structures matching frontend types

## Related Code Files

### Files to Create
- `backend/internal/models/extract.go` - Request/response models
- `backend/internal/models/invoice.go` - Invoice data models
- `backend/internal/handlers/extract.go` - Extract endpoint handler
- `backend/internal/services/extraction.go` - Extraction service

### Files to Modify
- `backend/cmd/server/main.go` - Register extract route

## Implementation Steps

1. **Create data models**
   - `ExtractRequest` struct (Image string)
   - `ExtractResponse` struct (Success, Data, Error, ProcessingTime)
   - `InvoiceData` struct (KeyValuePairs, Table, Summary, Confidence)
   - `KeyValuePair` struct (Key, Value, Confidence)
   - `TableData` struct (Headers, Rows)
   - Match frontend types exactly

2. **Create extraction service**
   - Initialize Gemini client
   - Build extraction prompt (Vietnamese invoice focused)
   - Send image + prompt to Gemini
   - Parse JSON response
   - Map to InvoiceData model
   - Handle errors and timeouts

3. **Create extract handler**
   - Bind JSON request
   - Validate base64 image format
   - Validate image size
   - Call extraction service
   - Measure processing time
   - Return formatted response
   - Handle errors appropriately

4. **Implement base64 validation**
   - Check data URL prefix
   - Validate base64 encoding
   - Decode and verify image format
   - Check size limits

5. **Implement Gemini integration**
   - Initialize client with API key
   - Create prompt for invoice extraction
   - Request structured JSON output
   - Handle Vietnamese language
   - Parse response to InvoiceData

6. **Register route in main.go**
   - Add POST `/api/extract` route
   - Wire handler to service
   - Add CORS middleware

7. **Error handling**
   - Validation errors → 400
   - Gemini API errors → 502
   - Internal errors → 500
   - Timeout errors → 504

## Todo List

- [ ] Create ExtractRequest and ExtractResponse models
- [ ] Create InvoiceData, KeyValuePair, TableData models
- [ ] Implement extraction service with Gemini client
- [ ] Create extraction prompt for Vietnamese invoices
- [ ] Implement base64 image validation
- [ ] Create extract handler
- [ ] Add request validation
- [ ] Implement error handling
- [ ] Add processing time measurement
- [ ] Register route in main.go
- [ ] Test with sample image

## Success Criteria

- Endpoint accepts POST requests to `/api/extract`
- Validates base64 image format correctly
- Calls Gemini API successfully
- Returns structured invoice data
- Response matches frontend ExtractResponse type
- Error handling works for all scenarios
- Processing time included in response
- Works with Vietnamese invoice images

## Risk Assessment

| Risk | Impact | Mitigation |
|------|--------|------------|
| Gemini API failures | High | Proper error handling, retry logic (future) |
| Invalid base64 input | Medium | Validation before processing |
| Response parsing errors | Medium | Robust JSON parsing with error handling |
| Timeout issues | Medium | Set appropriate timeout (30s) |
| Prompt not optimal | Medium | Iterate on prompt based on results |

## Security Considerations

- Validate all inputs
- Sanitize error messages (no sensitive data)
- Limit request size
- No image storage
- API key not exposed in logs

## Next Steps

- Proceed to Phase 3: Health Check & Error Handling
- Requires: Completed Phase 2, working extract endpoint
