---
phase: implementation  
title: Frontend Design Update - Stitch Invoice Scan Alignment
description: Updated frontend implementation to match provided screen designs
completed: true
---

# Frontend Design Update - Stitch Invoice Scan Alignment

## 🎨 Design System Integration

The frontend has been completely updated to match the provided Stitch Invoice Scan screen designs located in `docs/stitch_invoice_scan/`. The new implementation follows the exact visual specifications and user flow.

## 🔄 Updated Navigation Flow

### New Route Structure
```
/list-invoices          # Main dashboard (default route)
    ↓ (FAB button)
/take-picture          # Camera capture interface
    ↓ (Capture button)
/review-picture        # Image preview and confirmation
    ↓ (Extract Data button)
/extract-invoice-data  # Split-screen data editing
    ↓ (Back button)
/list-invoices         # Return to main dashboard
```

### Flow Description
1. **List Invoices**: Shows invoice history with status indicators
2. **Take Picture**: Full-screen camera with framing guide
3. **Review Picture**: Preview captured image before processing
4. **Extract Invoice Data**: Edit extracted data in split-screen layout
5. **Back Navigation**: All screens return to List Invoices on back button

## 🎨 Design Implementation

### Color Scheme (Updated Tailwind Config)
```javascript
colors: {
  primary: '#137fec',                    // Stitch blue
  'background-light': '#f6f7f8',        // Light theme background
  'background-dark': '#101922',          // Dark theme background  
  'text-light': '#0f172a',              // Light theme text
  'text-dark': '#f8fafc',               // Dark theme text
  'surface-light': '#ffffff',           // Light theme surfaces
  'surface-dark': '#1e293b',            // Dark theme surfaces
  'status-green': '#10b981',            // Success status
  'status-yellow': '#f59e0b',           // Warning status
  'status-blue': '#3b82f6'              // Processing status
}
```

### Typography
- **Font Family**: Inter (Stitch design system)
- **Font Weights**: 400 (normal), 500 (medium), 600 (semibold), 700 (bold)
- **Text Hierarchy**: Follows Stitch specifications

### Component Styling
- **Border Radius**: Consistent with Stitch (0.25rem default, 0.75rem for cards)
- **Spacing**: Following Stitch spacing scale
- **Shadows**: Subtle shadows for elevation
- **Icons**: Using text-based icons matching the designs

## 📱 Updated Page Components

### 1. ListInvoicesPage (`/list-invoices`)
**Design Source**: `docs/stitch_invoice_scan/list_invoices/`

**Features**:
- Invoice cards with thumbnails and status indicators
- Three status types: ✓ Done, ⚠ Review Needed, ↻ Extracting
- Floating Action Button (camera icon) for new scan
- Header with back arrow and search icon
- Status color coding (green/yellow/blue)

**Layout**:
```
┌─────────────────────────┐
│ ← Scanned Invoices   🔍 │
├─────────────────────────┤
│ [📄] Invoice #2024-07   │
│      ✓ Extraction Done  │
├─────────────────────────┤
│ [📄] Costco Wholesale   │
│      ⚠ Review Needed    │
├─────────────────────────┤
│ [📄] Office Depot       │
│      ↻ Extracting...    │
└─────────────────────────┘
                       [📷]
```

### 2. TakePicturePage (`/take-picture`)
**Design Source**: `docs/stitch_invoice_scan/take_picture/`

**Features**:
- Full-screen camera view with live video feed
- Dashed rectangular frame guide for invoice positioning
- Top controls: Close (✕) and Flash (⚡) buttons
- Bottom capture button (large blue circle with camera icon)
- Instructional text: "Position invoice within the frame"
- Black overlay with semi-transparent controls

**Layout**:
```
┌─────────────────────────┐
│ ✕               ⚡      │
│                         │
│   ┌─ ─ ─ ─ ─ ─ ─ ─ ─┐   │
│   ┊                 ┊   │
│   ┊    Frame Guide  ┊   │ 
│   ┊                 ┊   │
│   └─ ─ ─ ─ ─ ─ ─ ─ ─┘   │
│                         │
│ Position invoice within │
│     the frame...        │
│                         │
│        [📷]             │
└─────────────────────────┘
```

### 3. ReviewPicturePage (`/review-picture`)  
**Design Source**: `docs/stitch_invoice_scan/review_picture/`

**Features**:
- Dark theme background (#101922)
- Header with back arrow and "Review Invoice" title  
- Centered captured invoice image with rounded corners
- Question text: "Is the invoice clear and all details readable?"
- Two action buttons: "Retake" (outlined) and "Extract Data" (filled)
- White text on dark background

**Layout**:
```
┌─────────────────────────┐
│ ← Review Invoice        │
├─────────────────────────┤
│                         │
│ Is the invoice clear... │
│                         │
│ ┌─────────────────────┐ │
│ │                     │ │
│ │   Invoice Image     │ │
│ │                     │ │
│ └─────────────────────┘ │
│                         │
├─────────────────────────┤
│ [↻ Retake] [📄Extract]  │
└─────────────────────────┘
```

### 4. ExtractInvoiceDataPage (`/extract-invoice-data`)
**Design Source**: `docs/stitch_invoice_scan/extract_invoice_data/`

**Features**:
- Split-screen layout: Image left, data right
- Left panel: "Viewing Invoice Image" with captured photo
- Right panel: Structured extracted data with inline editing
- Three data sections: Invoice Information, Line Items, Summary
- Editable input fields for all extracted values
- Table structure for line items with dynamic columns
- "Save & Complete" button at bottom

**Layout**:
```
┌─────────────────────────┐
│ ← Invoice Data          │
├─────────────┬───────────┤
│ Viewing     │ Invoice   │
│ Invoice     │ Info:     │
│ Image       │ [Field 1] │
│             │ [Field 2] │
│ ┌─────────┐ │           │
│ │ Photo   │ │ Line Items│
│ │         │ │ ┌───┬───┐ │
│ └─────────┘ │ │   │   │ │
│             │ └───┴───┘ │
│             │           │
│             │ Summary:  │
│             │ [Total]   │
│             │ [Save &   │
│             │ Complete] │
└─────────────┴───────────┘
```

## 🔧 Technical Implementation

### Updated Dependencies
All existing dependencies maintained, with updated Tailwind configuration to support the new design system.

### Component Updates
- **Updated Tailwind Config**: New color scheme and design tokens
- **Icon System**: Text-based icons matching Stitch designs  
- **Dark/Light Theme**: Proper dark mode support for Review page
- **Responsive Layout**: Mobile-first approach with proper spacing

### Navigation Logic
```typescript
// Navigation flow implementation
ListInvoices → TakePicture → ReviewPicture → ExtractInvoiceData
     ↑                                                ↓
     ←─────────────── Back button ──────────────────
```

### State Management
- **Current Image**: Stored in Zustand for cross-page access
- **Extracted Data**: Flexible structure for dynamic invoice formats
- **Loading States**: Proper loading indicators during extraction
- **Error Handling**: User-friendly error messages

## ✅ Design Compliance Checklist

- [x] **Color Scheme**: Matches Stitch blue (#137fec) and backgrounds
- [x] **Typography**: Inter font family with proper weights  
- [x] **Layout Structure**: Exact layout matching screen designs
- [x] **Navigation Flow**: Correct 4-step user journey
- [x] **Interactive Elements**: Buttons, inputs, and controls match designs
- [x] **Status Indicators**: Proper status colors and icons
- [x] **Responsive Design**: Works on mobile devices (320px+)
- [x] **Dark Theme**: Review page uses dark background
- [x] **Loading States**: Extraction progress indicators
- [x] **Error Handling**: User-friendly error messages

## 🎯 User Experience Improvements

### Enhanced Flow
1. **Clear Intent**: Each page has a single, clear purpose
2. **Visual Feedback**: Status indicators show processing state
3. **Easy Navigation**: Back button always returns to main list
4. **Quality Check**: Review step ensures good image quality
5. **Inline Editing**: Direct editing of extracted data

### Mobile Optimization
- **Touch Targets**: All buttons are appropriately sized for touch
- **Visual Hierarchy**: Clear information hierarchy
- **Gesture Support**: Proper touch interactions
- **Performance**: Optimized for mobile device capabilities

## 🚀 Build Status

```bash
✓ TypeScript compilation successful
✓ Build completed in 758ms
✓ PWA assets generated
✓ Total bundle: ~309KB (99KB gzipped)
```

## 📱 Ready for Testing

The updated frontend is now ready for:

1. **Device Testing**: Camera functionality on iOS/Android
2. **Backend Integration**: API connection for data extraction  
3. **User Testing**: Flow validation with real users
4. **Performance Testing**: Real-world usage scenarios

---

**🎉 Design Update Complete!**

The frontend now perfectly matches the provided Stitch Invoice Scan designs and implements the specified navigation flow. All visual elements, colors, typography, and layouts align with the design specifications.