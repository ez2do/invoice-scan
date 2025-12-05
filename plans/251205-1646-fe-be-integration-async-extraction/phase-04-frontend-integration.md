# Phase 04: Frontend Integration

**Parent Plan**: [plan.md](./plan.md)  
**Dependencies**: Phase 03  
**Status**: Pending  
**Priority**: High

## Overview

Remove frontend mocks, integrate with backend APIs, implement polling for status updates, update UI to reflect async extraction flow.

## Key Insights

- TanStack Query's refetchInterval enables efficient polling
- Zustand store needs update for invoice list management
- Flow changes from sync extraction to async upload + list view

## Requirements

1. Remove mock data from ListInvoicesPage
2. Update API client with new endpoints
3. Update types for new response formats
4. Implement polling for status updates
5. Update pages for new async flow

## Related Code Files

- `frontend/src/lib/api.ts` - Add new endpoints
- `frontend/src/types/index.ts` - Update type definitions
- `frontend/src/stores/app-store.ts` - Update for invoice list
- `frontend/src/pages/ListInvoicesPage.tsx` - Remove mocks, add API
- `frontend/src/pages/ReviewPicturePage.tsx` - Change to upload flow
- `frontend/src/pages/ExtractInvoiceDataPage.tsx` - Fetch from API

## New Types

```typescript
export interface InvoiceListItem {
  id: string;  // ULID string (26 chars)
  status: 'pending' | 'processing' | 'completed' | 'failed';
  imagePath: string;
  extractedData?: InvoiceData;
  errorMessage?: string;
  createdAt: string;
  updatedAt: string;
}

export interface UploadResponse {
  success: boolean;
  data?: InvoiceListItem;
  error?: string;
}

export interface InvoicesResponse {
  success: boolean;
  data?: InvoiceListItem[];
  error?: string;
}
```

## Implementation Steps

### Step 1: Update API Client
```typescript
class APIClient {
  async uploadInvoice(imageDataUrl: string): Promise<UploadResponse> {...}
  async getInvoices(): Promise<InvoicesResponse> {...}
  async getInvoice(id: string): Promise<InvoiceResponse> {...}  // ULID string
}
```

### Step 2: Update Types
Add InvoiceListItem, UploadResponse, InvoicesResponse to types/index.ts.

### Step 3: Update Store
```typescript
interface AppStore {
  invoices: InvoiceListItem[];
  selectedInvoiceId: string | null;  // ULID string
  setInvoices: (invoices: InvoiceListItem[]) => void;
  setSelectedInvoiceId: (id: string | null) => void;
  // ... existing state
}
```

### Step 4: Update ListInvoicesPage
- Remove hardcoded mock data
- Use TanStack Query to fetch invoices
- Implement polling with refetchInterval
- Map API data to list items
- Handle loading/error states

```typescript
const { data, isLoading, error } = useQuery({
  queryKey: ['invoices'],
  queryFn: () => apiClient.getInvoices(),
  refetchInterval: 3000, // Poll every 3 seconds
});
```

### Step 5: Update ReviewPicturePage
- On "Continue" button: call uploadInvoice API
- Store returned invoice ID
- Navigate to list page (not extract page)

### Step 6: Update ExtractInvoiceDataPage
- Accept invoice ID from route params or store
- Fetch invoice data via API
- Poll until status is completed/failed
- Display extracted data or error

```typescript
const { data } = useQuery({
  queryKey: ['invoice', invoiceId],
  queryFn: () => apiClient.getInvoice(invoiceId),
  refetchInterval: (data) => {
    if (data?.data?.status === 'completed' || data?.data?.status === 'failed') {
      return false;
    }
    return 2000;
  },
});
```

### Step 7: Update Status Display
- pending: "Waiting..." with spinner
- processing: "Extracting..." with spinner
- completed: green checkmark
- failed: red X with error message

## Todo List

- [ ] Update types/index.ts with new interfaces
- [ ] Add new API methods to api.ts
- [ ] Update app-store.ts for invoice list
- [ ] Refactor ListInvoicesPage to use API
- [ ] Update ReviewPicturePage to use upload API
- [ ] Refactor ExtractInvoiceDataPage for API fetch
- [ ] Add polling logic with TanStack Query
- [ ] Update status icons and colors
- [ ] Test full flow end-to-end

## Success Criteria

- [ ] List page shows real invoices from database
- [ ] New scan uploads and shows in list immediately
- [ ] Status updates automatically via polling
- [ ] Clicking completed invoice shows extracted data
- [ ] Error handling displays user-friendly messages

## UI/UX Considerations

- Show loading skeleton while fetching list
- Optimistic update: add pending item to list immediately
- Stop polling when all items are completed/failed
- Pull-to-refresh on mobile (future enhancement)

## Environment Configuration

Ensure `VITE_API_URL` points to backend:
```env
VITE_API_URL=http://localhost:3001/api
```

## Next Steps

After completion, proceed to [Phase 05: Testing & Polish](./phase-05-testing-polish.md)

