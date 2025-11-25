---
phase: testing
title: Invoice Scan MVP - Testing Strategy
description: Define testing approach, test cases, and quality assurance for the invoice scanning MVP
---

# Invoice Scan MVP - Testing Strategy

## Test Coverage Goals

**What level of testing do we aim for?**

- Unit test coverage target: 80%+ for business logic (services, utils)
- Integration test scope: API endpoints, Gemini integration
- End-to-end test scenarios: Full scan-to-reconcile flow
- Manual testing: Camera functionality, PWA installation, cross-device

## Unit Tests

**What individual components need testing?**

### Server: Gemini Service

- [ ] Test case: Successfully parses valid JSON response from Gemini
- [ ] Test case: Handles malformed JSON response gracefully
- [ ] Test case: Extracts JSON from markdown-wrapped response
- [ ] Test case: Returns error for empty response
- [ ] Test case: Validates extracted data structure

### Server: Extract Route

- [ ] Test case: Returns 400 for missing image parameter
- [ ] Test case: Returns 400 for non-string image parameter
- [ ] Test case: Returns 400 for image exceeding size limit
- [ ] Test case: Returns 200 with valid invoice data on success
- [ ] Test case: Returns 500 with error message on Gemini failure

### Client: useCamera Hook

- [ ] Test case: Requests camera with environment facing mode
- [ ] Test case: Handles permission denied error
- [ ] Test case: Handles no camera available error
- [ ] Test case: Properly releases camera stream on cleanup

### Client: useExtraction Hook

- [ ] Test case: Sets loading state while fetching
- [ ] Test case: Sets data on successful response
- [ ] Test case: Sets error on failed response
- [ ] Test case: Clears error on retry

### Client: Invoice State Management

- [ ] Test case: Updates header field correctly
- [ ] Test case: Updates line item field correctly
- [ ] Test case: Recalculates total on line item change
- [ ] Test case: Adds new line item with default values
- [ ] Test case: Removes line item and updates total

## Integration Tests

**How do we test component interactions?**

### API Integration Tests

- [ ] POST /api/extract with valid image returns structured invoice data
- [ ] POST /api/extract handles Gemini timeout (mock long response)
- [ ] POST /api/extract handles Gemini rate limit response
- [ ] GET /api/health returns status ok

### Client-Server Integration

- [ ] Camera capture → API upload → Display results flow
- [ ] Error display when server is unreachable
- [ ] Retry functionality works after initial failure

## End-to-End Tests

**What user flows need validation?**

### Primary Flow: Scan and Reconcile

- [ ] **E2E 1**: User opens app → Camera view loads → Frame guide visible
- [ ] **E2E 2**: User captures image → Loading indicator shown → Results displayed
- [ ] **E2E 3**: User edits vendor name → Change persists in UI
- [ ] **E2E 4**: User expands line item → Details visible → Edit quantity → Total updates
- [ ] **E2E 5**: User adds new line item → Form appears → Fill and see in list

### Error Flows

- [ ] **E2E 6**: Camera permission denied → Helpful message shown
- [ ] **E2E 7**: Network offline → Error with retry button → Retry works when online
- [ ] **E2E 8**: Extraction fails → Error message → Retry button works

### PWA Flows

- [ ] **E2E 9**: Install prompt appears on supported browsers
- [ ] **E2E 10**: App works after installation (standalone mode)
- [ ] **E2E 11**: App loads from cache when revisited

## Test Data

**What data do we use for testing?**

### Mock Invoice Images

Create test fixtures in `/server/test/fixtures/`:

```
fixtures/
├── printed-vietnamese.jpg       # Clear printed invoice
├── handwritten-vietnamese.jpg   # Handwritten invoice
├── mixed-invoice.jpg            # Both printed and handwritten
├── blurry-invoice.jpg           # Poor quality for error testing
└── non-invoice.jpg              # Random image for edge case
```

### Expected Extraction Results

```typescript
export const mockPrintedInvoiceResult: InvoiceData = {
  header: {
    vendorName: 'Công ty TNHH ABC',
    invoiceNumber: 'HD-2024-001',
    invoiceDate: '2024-01-15',
    dueDate: '2024-02-15',
    vendorAddress: '123 Nguyễn Huệ, Quận 1, TP.HCM'
  },
  lineItems: [
    {
      id: '1',
      description: 'Giấy A4',
      quantity: 5,
      unitPrice: 50000,
      lineTotal: 250000
    }
  ],
  totals: {
    grandTotal: 250000
  },
  metadata: {
    processedAt: '2024-01-15T10:00:00Z',
    modelUsed: 'gemini-1.5-flash'
  }
};
```

### Mock Gemini Responses

```typescript
export const mockGeminiSuccess = {
  response: {
    text: () => JSON.stringify(mockPrintedInvoiceResult)
  }
};

export const mockGeminiMalformed = {
  response: {
    text: () => 'This is not valid JSON'
  }
};

export const mockGeminiEmpty = {
  response: {
    text: () => ''
  }
};
```

## Test Reporting & Coverage

**How do we verify and communicate test results?**

### Coverage Commands

```bash
# Server tests with coverage
cd server && npm run test:coverage

# Client tests with coverage
cd client && npm run test:coverage
```

### Coverage Thresholds

```json
{
  "jest": {
    "coverageThreshold": {
      "global": {
        "branches": 70,
        "functions": 80,
        "lines": 80,
        "statements": 80
      }
    }
  }
}
```

### Reporting

- Coverage reports generated in `coverage/` directory
- CI pipeline fails if coverage drops below threshold
- HTML report for detailed line-by-line view

## Manual Testing

**What requires human validation?**

### Device Testing Matrix

| Device | Browser | Tests |
|--------|---------|-------|
| iPhone 12+ | Safari | Camera, PWA install, gestures |
| iPhone SE | Safari | Small screen layout |
| Android Phone | Chrome | Camera, PWA install |
| Android Tablet | Chrome | Larger screen layout |
| Desktop | Chrome | Fallback experience |

### Manual Test Checklist

#### Camera Functionality
- [ ] Camera feed displays correctly
- [ ] Frame guide is visible and properly positioned
- [ ] Flash toggle works (on devices with flash)
- [ ] Capture button responsive to tap
- [ ] Captured image is clear and properly oriented

#### Reconciliation UI
- [ ] Invoice image displays correctly
- [ ] Zoom in/out works smoothly
- [ ] All extracted fields are displayed
- [ ] Fields are editable
- [ ] Collapsible sections work
- [ ] Line item expansion works
- [ ] Add item button works
- [ ] Total updates when editing values

#### PWA
- [ ] Install prompt appears (Chrome Android)
- [ ] App icon appears on home screen
- [ ] Opens in standalone mode
- [ ] Works after closing and reopening

#### Accessibility
- [ ] Touch targets are at least 44x44px
- [ ] Color contrast is sufficient
- [ ] Labels are associated with inputs
- [ ] App works with increased font size

### Vietnamese Text Testing

Test with real Vietnamese invoices containing:
- [ ] Common diacritics: à, á, ả, ã, ạ
- [ ] Special characters: đ, ơ, ư
- [ ] Compound vowels: ươ, iê, uy
- [ ] Currency format: 1.000.000 VND

## Performance Testing

**How do we validate performance?**

### Benchmarks to Measure

| Metric | Target | Measurement |
|--------|--------|-------------|
| PWA initial load | < 3s | Lighthouse |
| Cached load | < 1s | Lighthouse |
| Camera to preview | < 500ms | Manual stopwatch |
| Image capture | < 200ms | Console timing |
| API response (extraction) | < 5s | Network tab |

### Lighthouse Targets

```
Performance: > 90
Accessibility: > 90
Best Practices: > 90
PWA: Installable
```

### Load Testing (Optional for MVP)

If needed later:
- Use Artillery or k6
- Target: 10 concurrent extractions
- Response time < 10s under load

## Bug Tracking

**How do we manage issues?**

### Issue Labels

| Label | Description |
|-------|-------------|
| `bug` | Something isn't working |
| `enhancement` | New feature request |
| `camera` | Camera-related issues |
| `extraction` | AI extraction issues |
| `vietnamese` | Vietnamese language issues |
| `pwa` | PWA/installation issues |

### Severity Levels

| Level | Description | Response |
|-------|-------------|----------|
| Critical | App unusable | Fix immediately |
| High | Major feature broken | Fix in 24h |
| Medium | Feature partially broken | Fix in sprint |
| Low | Minor issue | Backlog |

### Regression Testing

Before each release:
1. Run all unit tests
2. Run integration tests
3. Complete manual test checklist
4. Test on at least 1 iOS and 1 Android device
5. Verify PWA installation works

