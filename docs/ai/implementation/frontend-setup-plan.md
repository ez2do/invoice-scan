---
phase: implementation
title: Frontend Setup Plan - Invoice Scan MVP
description: Detailed implementation plan for React/TypeScript PWA frontend setup
---

# Frontend Setup Plan - Invoice Scan MVP

## Overview

This document outlines the complete frontend implementation plan for the Invoice Scan MVP. The frontend will be a React-based Progressive Web App (PWA) with TypeScript, using modern tools and libraries for optimal mobile experience.

## Tech Stack

| Technology | Purpose | Reasoning |
|------------|---------|-----------|
| **React 18** | UI Framework | Component-based, great mobile support |
| **TypeScript** | Type Safety | Better DX, fewer runtime errors |
| **Vite** | Build Tool | Fast dev server, optimized builds |
| **Tailwind CSS** | Styling | Utility-first, responsive design |
| **React Router** | Navigation | Client-side routing for SPA |
| **TanStack Query** | API State | Server state management, caching |
| **Zustand** | Local State | Lightweight state management |
| **React Hook Form** | Forms | Form validation, performance |
| **Framer Motion** | Animations | Smooth UI transitions |

## Project Structure

```
frontend/
├── public/
│   ├── manifest.json           # PWA manifest
│   ├── sw.js                   # Service worker
│   ├── icons/                  # PWA icons (various sizes)
│   └── favicon.ico
├── src/
│   ├── components/             # Reusable components
│   │   ├── ui/                 # Base UI components
│   │   ├── camera/             # Camera-specific components
│   │   ├── invoice/            # Invoice data components
│   │   └── layout/             # Layout components
│   ├── hooks/                  # Custom React hooks
│   ├── lib/                    # Utilities and configurations
│   ├── pages/                  # Page components
│   ├── stores/                 # Zustand stores
│   ├── types/                  # TypeScript type definitions
│   ├── utils/                  # Helper functions
│   ├── App.tsx                 # Root component
│   ├── main.tsx                # Entry point
│   └── index.css               # Global styles
├── package.json
├── tsconfig.json
├── tailwind.config.js
├── vite.config.ts
└── README.md
```

## Implementation Phases

### Phase 1: Project Scaffolding (30 mins)

#### 1.1 Initialize Vite React TypeScript Project
```bash
npm create vite@latest frontend -- --template react-ts
cd frontend
npm install
```

#### 1.2 Install Core Dependencies
```bash
# Routing and state management
npm install react-router-dom @tanstack/react-query zustand

# UI and styling
npm install tailwindcss autoprefixer postcss
npm install @headlessui/react @heroicons/react
npm install framer-motion

# Forms and validation
npm install react-hook-form @hookform/resolvers zod

# PWA and utilities
npm install workbox-window
npm install clsx tailwind-merge

# Development dependencies
npm install -D @types/node
```

#### 1.3 Configure Tailwind CSS
```bash
npx tailwindcss init -p
```

#### 1.4 Setup TypeScript Configuration
- Configure path aliases
- Strict type checking
- Modern target settings

### Phase 2: Base Configuration (45 mins)

#### 2.1 PWA Configuration
- **Manifest.json**: App metadata, icons, theme colors
- **Service Worker**: Basic caching strategy
- **Icons**: Generate multiple sizes (72x72 to 512x512)

#### 2.2 Vite Configuration
- **PWA Plugin**: Workbox integration
- **Path Aliases**: `@/` for src directory
- **Build Optimization**: Bundle analysis, chunk splitting

#### 2.3 Router Setup
- **Route Structure**:
  ```
  / → Home (landing page)
  /scan → Camera capture
  /verify → Data verification
  ```

#### 2.4 Global Styles and Theme
- **Design System**: Colors, typography, spacing
- **Mobile-First**: Responsive breakpoints
- **Dark Mode**: Optional theme support

### Phase 3: Type Definitions (30 mins)

#### 3.1 Invoice Data Types
```typescript
// Flexible invoice structure
interface InvoiceData {
  keyValuePairs: KeyValuePair[];
  table: TableData | null;
  summary: KeyValuePair[];
  confidence?: number;
}

interface KeyValuePair {
  key: string;
  value: string;
  confidence?: number;
}

interface TableData {
  headers: string[];
  rows: string[][];
}
```

#### 3.2 API Types
```typescript
interface ExtractRequest {
  image: string; // base64 encoded
}

interface ExtractResponse {
  success: boolean;
  data?: InvoiceData;
  error?: string;
  processingTime?: number;
}
```

#### 3.3 App State Types
```typescript
interface AppState {
  currentImage: string | null;
  extractedData: InvoiceData | null;
  isLoading: boolean;
  error: string | null;
}
```

### Phase 4: Core Components (2 hours)

#### 4.1 Layout Components
- **AppLayout**: Main container with navigation
- **MobileLayout**: Mobile-optimized layout
- **LoadingSpinner**: Global loading component
- **ErrorBoundary**: Error handling wrapper

#### 4.2 UI Components Library
- **Button**: Primary, secondary, icon variants
- **Input**: Text input with validation
- **Modal**: Overlay dialogs
- **Toast**: Success/error notifications
- **Card**: Content containers

#### 4.3 Camera Components
- **CameraView**: Main camera interface
- **CameraControls**: Capture, flash, close buttons
- **CameraFrame**: Invoice positioning guide
- **ImagePreview**: Captured image display

#### 4.4 Invoice Components
- **InvoiceDisplay**: Split-screen layout (image + data)
- **KeyValueSection**: Editable key-value pairs
- **TableSection**: Dynamic table with editable cells
- **SummarySection**: Totals and calculations
- **EditableField**: Individual field editing

### Phase 5: Custom Hooks (1 hour)

#### 5.1 Camera Hook
```typescript
const useCamera = () => {
  // Camera access, capture, error handling
  const { stream, captureImage, error, isSupported } = ...
  return { stream, captureImage, error, isSupported };
}
```

#### 5.2 API Integration Hook
```typescript
const useInvoiceExtraction = () => {
  // TanStack Query integration for API calls
  const mutation = useMutation({
    mutationFn: extractInvoiceData,
    onSuccess: (data) => ...,
    onError: (error) => ...,
  });
  return mutation;
}
```

#### 5.3 Local Storage Hook
```typescript
const useLocalStorage = <T>(key: string, initialValue: T) => {
  // Persistent local storage with type safety
}
```

### Phase 6: State Management (30 mins)

#### 6.1 Zustand Store
```typescript
interface AppStore {
  // Current app state
  currentImage: string | null;
  extractedData: InvoiceData | null;
  
  // Actions
  setCurrentImage: (image: string) => void;
  setExtractedData: (data: InvoiceData) => void;
  clearData: () => void;
}
```

#### 6.2 Query Client Setup
- **TanStack Query**: Global configuration
- **Caching Strategy**: 5-minute cache for API calls
- **Error Handling**: Global error boundaries

### Page Components (1.5 hours)

#### 7.1 List Invoices Page  
- **Invoice History**: List of previously scanned invoices
- **Status Indicators**: Done, Review Needed, Extracting
- **Floating Action Button**: Camera icon for new scan
- **Search Functionality**: Find specific invoices

#### 7.2 Take Picture Page
- **Full-Screen Camera**: Environment camera with overlay
- **Framing Guide**: Dashed border for invoice positioning
- **Camera Controls**: Capture button, flash toggle, close
- **Instruction Text**: "Position invoice within frame"

#### 7.3 Review Picture Page
- **Image Preview**: Full-width captured invoice image
- **Action Buttons**: Retake or Extract Data
- **Dark Theme**: Black background with white text
- **Quality Check**: "Is the invoice clear and readable?"

#### 7.4 Extract Invoice Data Page
- **Split Layout**: Image on left, extracted data on right
- **Editable Fields**: Inline editing for all extracted values
- **Structured Display**: Key-values, table data, summary sections
- **Save & Complete**: Final action to finish processing

### Phase 8: PWA Features (45 mins)

#### 8.1 Service Worker
- **Cache Strategy**: Cache-first for static assets
- **Network Strategy**: Network-first for API calls
- **Offline Handling**: Show offline status

#### 8.2 Install Prompt
- **A2HS**: Add to Home Screen prompt
- **Install Button**: Manual install trigger
- **Installation Events**: Track install success

#### 8.3 Responsive Design
- **Mobile-First**: Optimized for 320px+ screens
- **Touch Targets**: 44px minimum for buttons
- **Viewport Meta**: Proper mobile scaling

## Implementation Steps

### Step 1: Bootstrap Project (Execute Now)
1. Create Vite project with React TypeScript template
2. Install all dependencies
3. Configure Tailwind CSS
4. Setup basic project structure
5. Configure Vite for PWA

### Step 2: Type System Setup
1. Create comprehensive TypeScript interfaces
2. Setup API client types
3. Configure path aliases and imports

### Step 3: UI Foundation
1. Build core UI component library
2. Implement design system
3. Setup responsive layout components

### Step 4: Camera Feature
1. Implement camera access hook
2. Build camera UI components
3. Add image capture and compression

### Step 5: Data Display
1. Create flexible invoice data components
2. Implement editing functionality
3. Add validation and error handling

### Step 6: Integration
1. Connect frontend to backend API
2. Add loading and error states
3. Implement data persistence

### Step 7: PWA Setup
1. Configure service worker
2. Add PWA manifest
3. Test installation flow

### Step 8: Testing & Optimization
1. Test on target mobile devices
2. Optimize bundle size
3. Performance testing

## Success Criteria

- [ ] PWA installs successfully on iOS and Android
- [ ] Camera captures high-quality images
- [ ] UI is responsive on screens 320px and wider
- [ ] Data extraction displays in organized, editable format
- [ ] App loads in < 3 seconds on 4G
- [ ] Offline handling works gracefully
- [ ] TypeScript compilation with zero errors
- [ ] Bundle size < 1MB compressed

## File Ownership

This implementation will create/modify these files:

```
frontend/                       # New directory
├── package.json               # Dependencies and scripts
├── tsconfig.json              # TypeScript configuration
├── vite.config.ts            # Build configuration
├── tailwind.config.js        # Styling configuration
├── public/manifest.json      # PWA manifest
├── public/sw.js             # Service worker
├── src/main.tsx             # App entry point
├── src/App.tsx              # Root component
├── src/types/index.ts       # Type definitions
├── src/lib/api.ts           # API client
├── src/hooks/useCamera.ts   # Camera functionality
├── src/pages/HomePage.tsx   # Landing page
├── src/pages/ScanPage.tsx   # Camera capture
├── src/pages/VerifyPage.tsx # Data verification
└── src/components/          # All UI components
```

## Next Steps After Implementation

1. **Backend Integration**: Connect to Node.js API
2. **Device Testing**: Test on actual iOS/Android devices  
3. **Performance Optimization**: Bundle analysis and optimization
4. **Deployment**: Build and deploy to production server
5. **SSL Setup**: Configure HTTPS for camera access

---

*This plan provides a comprehensive roadmap for implementing the frontend. Each phase builds upon the previous one, ensuring a solid foundation for the Invoice Scan MVP.*