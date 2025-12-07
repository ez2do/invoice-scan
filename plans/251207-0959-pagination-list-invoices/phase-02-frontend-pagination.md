# Phase 2: Frontend Pagination UI

**Priority**: High | **Status**: ⏳ Pending

## Context

Frontend currently fetches all invoices. Need to add pagination controls and update API client.

## Overview

- Update API client to pass pagination params
- Add pagination controls UI
- Update query key to include page number

## Related Files

- [api.ts](file:///Users/tuananh/projects/invoice-scan/frontend/src/lib/api.ts)
- [types/index.ts](file:///Users/tuananh/projects/invoice-scan/frontend/src/types/index.ts)
- [ListInvoicesPage.tsx](file:///Users/tuananh/projects/invoice-scan/frontend/src/pages/ListInvoicesPage.tsx)
- [index.css](file:///Users/tuananh/projects/invoice-scan/frontend/src/index.css)

---

## Implementation Steps

### 1. Update TypeScript types

#### [MODIFY] [index.ts](file:///Users/tuananh/projects/invoice-scan/frontend/src/types/index.ts)

Add pagination types:
```typescript
export interface PaginationParams {
  page: number;
  pageSize: number;
}

export interface PaginatedResponse<T> {
  success: boolean;
  data: T[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export type PaginatedInvoicesResponse = PaginatedResponse<InvoiceListItem>;
```

---

### 2. Update API client

#### [MODIFY] [api.ts](file:///Users/tuananh/projects/invoice-scan/frontend/src/lib/api.ts)

Update `getInvoices()` to accept pagination params:
```typescript
async getInvoices(page = 1, pageSize = 10): Promise<PaginatedInvoicesResponse>
```

---

### 3. Add pagination UI component styles

#### [MODIFY] [index.css](file:///Users/tuananh/projects/invoice-scan/frontend/src/index.css)

Add pagination button styles using existing design system.

---

### 4. Update List Invoices Page

#### [MODIFY] [ListInvoicesPage.tsx](file:///Users/tuananh/projects/invoice-scan/frontend/src/pages/ListInvoicesPage.tsx)

- Add `page` state with `useState`
- Pass page to API call
- Update query key to include page
- Add pagination controls at bottom
- Show current page / total pages info

---

## Todo

- [ ] Add pagination TypeScript types
- [ ] Update API client method
- [ ] Add pagination CSS styles
- [ ] Implement pagination UI in list page

## Success Criteria

- Page loads with default page 1
- Can navigate between pages
- Shows "Page X of Y" indicator
- Previous/Next buttons work correctly
- Disabled states for first/last page
