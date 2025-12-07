# Pagination for List Invoices Page

Implement cursor-based pagination for the invoice listing, both backend API and frontend UI.

## Phases

| Phase | Description | Status |
|-------|-------------|--------|
| [Phase 1](./phase-01-backend-pagination.md) | Backend API pagination | ⏳ Pending |
| [Phase 2](./phase-02-frontend-pagination.md) | Frontend pagination UI | ⏳ Pending |

## Summary

- **Approach**: Offset-based pagination (simpler, suitable for small-medium datasets)
- **Default page size**: 10 items
- **Backend changes**: Modify repository `List()` and handler
- **Frontend changes**: Add pagination controls with page numbers

## Quick Links

- [Backend Invoice Handler](file:///Users/tuananh/projects/invoice-scan/backend/internal/handlers/invoice.go)
- [Frontend List Page](file:///Users/tuananh/projects/invoice-scan/frontend/src/pages/ListInvoicesPage.tsx)
- [API Client](file:///Users/tuananh/projects/invoice-scan/frontend/src/lib/api.ts)
