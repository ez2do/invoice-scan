---
phase: requirements
title: Invoice Scan MVP - Requirements & Problem Understanding
description: Clarify the problem space, gather requirements, and define success criteria for the invoice scanning MVP
---

# Invoice Scan MVP - Requirements & Problem Understanding

## Problem Statement

**What problem are we solving?**

- **Core Problem**: Extracting data from physical invoices, especially handwritten ones in Vietnamese, is tedious, time-consuming, and error-prone when done manually.
- **Who is affected**: Internal finance/accounting staff who need to process invoices regularly.
- **Current situation**: Staff manually read invoices and type data into systems, leading to:
  - High time investment per invoice
  - Transcription errors, especially with handwritten content
  - Difficulty reading Vietnamese diacritics in poor handwriting
  - No easy way to verify extracted data against the original document

## Goals & Objectives

**What do we want to achieve?**

### Primary Goals
- Enable users to scan physical invoices using their mobile device camera
- Automatically extract invoice data using AI (Gemini Vision)
- Display extracted data alongside the original image for user verification
- Support Vietnamese language invoices (printed and handwritten)
- Support various invoice formats without predefined field constraints

### Secondary Goals
- Provide a mobile-first, installable experience via PWA
- Keep infrastructure simple and maintainable (self-hosted)
- Minimize per-invoice processing costs

### Non-Goals (Out of Scope for MVP)
- Saving/persisting invoice data to a database
- Multi-user authentication and authorization
- Invoice history and search
- Batch processing of multiple invoices
- Offline processing capability
- Export to accounting software formats
- Integration with external systems

## User Stories & Use Cases

**How will users interact with the solution?**

### Primary User Story
> As a **finance staff member**, I want to **scan a physical invoice with my phone and see the extracted data displayed alongside the original image**, so that I can **quickly verify the extraction accuracy**.

### Key Workflows

#### Workflow 1: Scan and Verify Invoice
1. User opens the PWA on their mobile device
2. User positions the invoice within the camera frame
3. User taps the capture button
4. System uploads image and extracts data via Gemini API
5. User sees extracted data displayed alongside the invoice image
6. User visually compares extracted data with the original image
7. User can edit any incorrectly extracted fields
8. User confirms the data is correct

#### Edge Cases
- Poor lighting conditions → Show guidance message
- Blurry image → Allow retake
- Partially visible invoice → Show warning
- Gemini API unavailable → Show error with retry option
- Handwriting too unclear → Extract what's possible, user corrects manually

## Invoice Data Model

**Flexible structure to support various invoice formats**

The system should NOT enforce a rigid schema. Instead, it extracts data based on the common anatomy of invoices:

### Generic Invoice Anatomy

```
┌─────────────────────────────────────────┐
│           KEY-VALUE SECTION             │
│  (Header information, metadata)         │
│  - Vendor name, address                 │
│  - Invoice number, date                 │
│  - Customer info                        │
│  - Any other key-value pairs            │
└─────────────────────────────────────────┘
┌─────────────────────────────────────────┐
│           TABLE SECTION                 │
│  ┌─────────────────────────────────┐    │
│  │  Table Header (column names)   │    │
│  ├─────────────────────────────────┤    │
│  │  Item row 1                    │    │
│  │  Item row 2                    │    │
│  │  Item row 3                    │    │
│  │  ...                           │    │
│  └─────────────────────────────────┘    │
└─────────────────────────────────────────┘
┌─────────────────────────────────────────┐
│           SUMMARY SECTION               │
│  (Totals, taxes, discounts, etc.)       │
│  - Subtotal                             │
│  - Tax                                  │
│  - Grand total                          │
│  - Any other summary fields             │
└─────────────────────────────────────────┘
```

### Data Structure Principles

1. **Key-Value Pairs**: Any metadata fields found at the top/header of the invoice
2. **Table Data**: Column headers + rows of items (columns vary by invoice type)
3. **Summary**: Any totals/calculations at the bottom
4. **No Fixed Schema**: The AI extracts whatever fields exist in each invoice

### Example: Different Invoice Types

| Invoice Type | Key-Value Fields | Table Columns | Summary Fields |
|--------------|-----------------|---------------|----------------|
| Office Supplies | Vendor, Date, Invoice # | Item, Qty, Price | Subtotal, Tax, Total |
| Restaurant | Restaurant name, Date, Table # | Dish, Qty, Price | Subtotal, Service, Total |
| Medical | Clinic, Patient, Date | Service, Code, Amount | Insurance, Co-pay, Due |
| Utility Bill | Provider, Account #, Period | Usage, Rate | Consumption, Tax, Amount Due |

The system adapts to whatever format is presented.

## Success Criteria

**How will we know when we're done?**

### Functional Criteria
- [ ] Camera capture works on iOS and Android mobile browsers
- [ ] Invoice image successfully sent to server for processing
- [ ] Gemini Vision API extracts invoice data dynamically (no fixed schema)
- [ ] Extracted data displayed in a flexible format alongside invoice image
- [ ] User can edit extracted fields
- [ ] PWA is installable on mobile devices
- [ ] Vietnamese text (including diacritics) correctly extracted

### Performance Benchmarks
- [ ] Image upload + extraction completes in < 5 seconds (average)
- [ ] PWA loads in < 3 seconds on 4G connection
- [ ] Cost per invoice extraction < $0.005

### Quality Criteria
- [ ] Extraction accuracy > 90% for printed invoices
- [ ] Extraction accuracy > 75% for handwritten invoices
- [ ] UI displays extracted data in a readable, organized format
- [ ] UI is responsive and works on screens 320px and wider

## Constraints & Assumptions

**What limitations do we need to work within?**

### Technical Constraints
- Must work as a Progressive Web App (no native app stores)
- Self-hosted server infrastructure
- Gemini API requires Google Cloud account with billing enabled
- Camera access requires HTTPS (SSL certificate needed)

### Business Constraints
- MVP scope only - no data persistence
- Internal tool - no public access required
- Single language focus: Vietnamese (with English fallback)

### Assumptions
- Users have modern smartphones with camera access
- Users have stable internet connection during scanning
- Invoice formats vary but follow general invoice anatomy (key-values, table, summary)
- Server has outbound internet access to call Gemini API

## Questions & Open Items

**What do we still need to clarify?**

### Resolved
- [x] LLM Choice → Gemini 1.5 Flash (cost-effective, good Vietnamese support)
- [x] Deployment → Self-hosted server
- [x] Invoice types → Various formats, no fixed schema
- [x] Invoice language → Mixed (printed + handwritten), mostly Vietnamese
- [x] Offline support → Not required (keep it simple)

### Open Questions
- [ ] What server technology stack? (Node.js recommended for simplicity)
- [ ] Hosting environment? (VPS, Docker, bare metal?)
- [ ] Domain name for the PWA?
- [ ] SSL certificate approach? (Let's Encrypt recommended)
- [ ] UI layout for flexible data display? (key-value list + dynamic table + summary)
