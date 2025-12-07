# Edit and Save Invoice Data Feature - Implementation Plan

**Date**: 2025-12-07  
**Type**: Feature Implementation  
**Status**: Planning  
**Priority**: High

## Executive Summary

Enable users to edit extracted invoice data and persist changes back to the database. The UI already supports editing via `AutoExpandTextarea` components, but clicking "Save & Complete" currently only navigates back without persisting changes. This feature adds an UPDATE API endpoint and integrates it with the frontend, including unsaved changes detection and confirmation dialogs.

## User Requirements

Based on user clarification:
1. **Save edited data to `extracted_data` field** - Overwrite the AI-extracted data with user edits
2. **No validation for now** - Accept any user input without field-level validation  
3. **Confirmation dialog before leaving** - Warn users when leaving page with unsaved changes

## Context Links

- **Related Plans**: 
  - [251205-1646-fe-be-integration-async-extraction](file:///Users/tuananh/projects/invoice-scan/plans/251205-1646-fe-be-integration-async-extraction/plan.md) - Initial async extraction implementation
- **Dependencies**: 
  - Backend: Go, Gin framework, GORM
  - Frontend: React, TanStack Query, Zustand
  - Database: MySQL with `invoices` table
- **Reference Docs**: 
  - [System Architecture](file:///Users/tuananh/projects/invoice-scan/docs/system-architecture.md)
  - [Code Standards](file:///Users/tuananh/projects/invoice-scan/docs/code-standards.md)

## Implementation Phases

### Phase 1: Backend API - Update Endpoint

**Files to modify:**
- [backend/internal/handlers/invoice.go](file:///Users/tuananh/projects/invoice-scan/backend/internal/handlers/invoice.go) - Add `Update` handler
- [backend/internal/handlers/dto.go](file:///Users/tuananh/projects/invoice-scan/backend/internal/handlers/dto.go) - Add update request DTO
- [backend/cmd/server/main.go](file:///Users/tuananh/projects/invoice-scan/backend/cmd/server/main.go) - Register PUT route

**Tasks:**
1. Add `UpdateInvoiceRequest` DTO struct with `extracted_data` field
2. Implement `Update(c *gin.Context)` handler method that:
   - Validates invoice ID from URL parameter
   - Parses request body into DTO
   - Fetches existing invoice from repository
   - Updates `ExtractedData` field via repository's `Update` method
   - Returns updated invoice response
3. Register `PUT /api/v1/invoices/:id` route

**Acceptance Criteria:**
- PUT endpoint accepts JSON with `extracted_data` field
- Successfully updates invoice in database
- Returns 200 with updated invoice data
- Returns 404 if invoice not found
- Returns 400 for invalid request body

---

### Phase 2: Frontend API Client

**Files to modify:**
- [frontend/src/lib/api.ts](file:///Users/tuananh/projects/invoice-scan/frontend/src/lib/api.ts) - Add `updateInvoice` method

**Tasks:**
1. Add `updateInvoice(id: string, data: ExtractedData)` method to `APIClient` class
2. Method should:
   - Make PUT request to `/api/v1/invoices/:id`
   - Send `extracted_data` in request body
   - Return `InvoiceResponse` on success
   - Handle errors appropriately

**Acceptance Criteria:**
- Method correctly formats request payload
- Handles successful response (200)
- Properly propagates errors from backend

---

### Phase 3: Frontend State & UI Integration

**Files to modify:**
- [frontend/src/stores/app-store.ts](file:///Users/tuananh/projects/invoice-scan/frontend/src/stores/app-store.ts) - Add dirty state tracking
- [frontend/src/pages/ExtractInvoiceDataPage.tsx](file:///Users/tuananh/projects/invoice-scan/frontend/src/pages/ExtractInvoiceDataPage.tsx) - Integrate save logic

**Tasks:**
1. **Update `app-store.ts`:**
   - Add `isDirty` boolean state (tracks if data was modified)
   - Add `setDirty(value: boolean)` action
   - Modify `updateKeyValue`, `updateTableCell`, `updateSummary` to set `isDirty = true`
   - Add `resetDirty()` action

2. **Update `ExtractInvoiceDataPage.tsx`:**
   - Import `useMutation` from TanStack Query
   - Create mutation for `apiClient.updateInvoice`
   - Update `handleComplete` to:
     - Collect current `extractedData` from store
     - Call mutation with invoice ID and data
     - Show loading state during save
     - On success: reset dirty state, show success message, navigate back
     - On error: show error message, stay on page
   - Add `useEffect` with `beforeunload` event listener when `isDirty` is true
   - Add custom confirmation dialog when navigating back with unsaved changes

**Acceptance Criteria:**
- Save button triggers API call with current data
- Loading indicator shown during save
- Success: updates persist in database, user navigated back
- Error: error message displayed, user stays on page
- Browser warns before closing/refreshing with unsaved changes
- Back button warns before navigating with unsaved changes

---

## Verification Plan

### Backend Tests
```bash
# Navigate to backend directory
cd backend

# Run all tests including new update endpoint tests
go test ./... -v
```

**Manual Backend Testing:**
```bash
# Test update endpoint with curl
curl -X PUT http://localhost:3001/api/v1/invoices/{invoice_id} \
  -H "Content-Type: application/json" \
  -d '{
    "extracted_data": {
      "key_value_pairs": [{"key": "Test", "value": "Updated"}],
      "table": {"headers": [], "rows": []},
      "summary": [],
      "confidence": 0.95
    }
  }'
```

### Frontend Tests

**Manual Testing Flow:**
1. Start backend: `cd backend && go run cmd/server/main.go`
2. Start frontend: `cd frontend && npm run dev`
3. Upload an invoice and wait for extraction to complete
4. Navigate to invoice detail page
5. Edit some fields (key-value pairs, table cells, summary)
6. Verify browser shows warning when trying to close tab
7. Click back button - verify confirmation dialog appears
8. Cancel dialog - verify staying on page
9. Click "Save & Complete" button
10. Verify loading state appears
11. Verify success message/navigation after save
12. Navigate back to invoice detail page
13. Verify edited data persisted correctly

### Integration Testing
- Verify data round-trip: Edit → Save → Reload → Data matches edits
- Verify dirty state resets after successful save
- Verify no warning when navigating after save

---

## Security Considerations

- [x] No authentication required yet (MVP phase)
- [x] Input sanitization handled by JSON marshaling
- [ ] Future: Add user ownership validation when auth is implemented
- [ ] Future: Add field-level validation when requirements are defined

---

## Risk Assessment

| Risk | Impact | Mitigation |
|------|--------|------------|
| Data loss on failed save | High | Show clear error message, keep user on page with data intact |
| User accidentally loses edits | Medium | Implement beforeunload and navigation warnings |
| Invalid JSON structure breaks backend | Medium | Backend validates JSON structure, returns 400 on invalid input |
| Race condition if extraction still processing | Low | UI only allows edit when status is 'completed' |

---

## Success Criteria

- [ ] Backend UPDATE endpoint implemented and tested
- [ ] Frontend API client method added
- [ ] Save functionality integrated with UI
- [ ] Unsaved changes warning implemented (browser close)
- [ ] Unsaved changes warning implemented (back navigation)
- [ ] Manual testing verified successful
- [ ] Data persists correctly after save
- [ ] Error handling works as expected
