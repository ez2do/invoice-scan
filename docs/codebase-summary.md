# Codebase Summary

**Last Updated**: 2025-12-06
**Version**: 1.0.0
**Repository**: invoice-scan

## Overview

Invoice Scan MVP is a React-based Progressive Web Application for scanning and extracting data from invoices using AI-powered vision recognition. The application is built with a modern frontend stack and designed for mobile-first usage, supporting Vietnamese language invoices (both printed and handwritten).

## Project Structure

```
invoice-scan/
├── frontend/                    # React PWA application
│   ├── src/
│   │   ├── components/         # React components
│   │   │   └── ui/             # Base UI components (Button, Input)
│   │   ├── hooks/              # Custom React hooks
│   │   │   ├── useCamera.ts
│   │   │   ├── useInvoiceExtraction.ts
│   │   │   ├── usePWANavigation.ts
│   │   │   └── usePWARouter.ts
│   │   ├── lib/                # Utilities and API client
│   │   │   ├── api.ts
│   │   │   └── utils.ts
│   │   ├── pages/              # Page components
│   │   │   ├── ListInvoicesPage.tsx
│   │   │   ├── TakePicturePage.tsx
│   │   │   ├── ReviewPicturePage.tsx
│   │   │   └── ExtractInvoiceDataPage.tsx
│   │   ├── stores/             # Zustand state management
│   │   │   └── app-store.ts
│   │   ├── types/              # TypeScript definitions
│   │   │   └── index.ts
│   │   ├── App.tsx            # Main app component
│   │   ├── main.tsx           # Entry point
│   │   └── index.css           # Global styles
│   ├── public/                # Static assets
│   │   └── icons/             # PWA icons (SVG)
│   ├── index.html
│   ├── package.json
│   ├── vite.config.ts         # Vite configuration with PWA plugin
│   ├── tailwind.config.js
│   └── tsconfig.json
├── docs/                      # Project documentation
│   ├── assets/               # Documentation assets
│   ├── project-overview-pdr.md
│   ├── codebase-summary.md
│   ├── code-standards.md
│   ├── system-architecture.md
│   └── project-roadmap.md
├── plans/                     # Implementation plans
│   └── templates/            # Plan templates
├── .cursor/                   # Cursor IDE configuration
│   ├── commands/
│   └── rules/
├── .gitignore
├── .repomixignore
├── AGENTS.md
├── CLAUDE.md
└── README.md
```

## Core Technologies

### Frontend Stack
- **React**: 19.2.0 - UI library
- **TypeScript**: 5.9.3 - Type safety
- **Vite**: 7.2.4 - Build tool and dev server
- **React Router**: 7.9.6 - Client-side routing
- **TanStack Query**: 5.90.11 - Data fetching and caching
- **Zustand**: 5.0.8 - State management
- **Tailwind CSS**: 3.4.18 - Utility-first CSS
- **vite-plugin-pwa**: 1.2.0 - PWA support

### Development Tools
- **ESLint**: 9.39.1 - Code linting
- **TypeScript ESLint**: 8.46.4 - TypeScript linting
- **PostCSS**: 8.5.6 - CSS processing
- **Autoprefixer**: 10.4.22 - CSS vendor prefixes

### PWA Features
- Service Worker for offline support
- Web App Manifest for installation
- App icons (72x72 to 512x512)
- SSL support (basic-ssl plugin)

## Key Components

### 1. Pages

**ListInvoicesPage** (`src/pages/ListInvoicesPage.tsx`):
- Landing page with invoice list
- Navigation to scan functionality
- Entry point for the application

**TakePicturePage** (`src/pages/TakePicturePage.tsx`):
- Camera interface with live preview
- Frame guide overlay for invoice positioning
- Capture button and controls
- Camera permission handling
- Error state management

**ReviewPicturePage** (`src/pages/ReviewPicturePage.tsx`):
- Preview captured invoice image
- Option to retake or proceed
- Image validation

**ExtractInvoiceDataPage** (`src/pages/ExtractInvoiceDataPage.tsx`):
- Side-by-side layout (image + extracted data)
- Loading state during extraction
- Error handling and display
- Editable extracted fields
- Key-value pairs section
- Dynamic table display
- Summary section
- Save and complete actions

### 2. Custom Hooks

**useCamera** (`src/hooks/useCamera.ts`):
- Camera stream management
- MediaDevices API integration
- Image capture functionality
- Error handling
- Camera permission handling

**useInvoiceExtraction** (`src/hooks/useInvoiceExtraction.ts`):
- API integration for extraction
- TanStack Query mutation
- Loading and error states
- Success handling

**usePWANavigation** (`src/hooks/usePWANavigation.ts`):
- PWA-specific navigation handling
- Install prompt management
- Service worker updates

**usePWARouter** (`src/hooks/usePWARouter.ts`):
- Router configuration for PWA
- Route handling

### 3. State Management

**app-store** (`src/stores/app-store.ts`):
- Zustand store for application state
- Current image management
- Extracted data storage
- Loading and error states
- Data update methods (key-value, table, summary)

### 4. API Client

**api.ts** (`src/lib/api.ts`):
- APIClient class
- extractInvoice method
- healthCheck method
- Error handling
- Base URL configuration

### 5. Type Definitions

**types/index.ts**:
- `KeyValuePair` - Key-value data structure
- `TableData` - Table structure with headers and rows
- `InvoiceData` - Complete invoice data structure
- `ExtractRequest` - API request format
- `ExtractResponse` - API response format
- `AppState` - Application state interface
- `CameraConfig` - Camera configuration
- UI component prop types

## Data Flow

### Invoice Scanning Flow

```
1. User opens app (ListInvoicesPage)
   ↓
2. User taps "Scan Invoice" → Navigate to TakePicturePage
   ↓
3. Camera permission requested → Camera stream starts
   ↓
4. User positions invoice → Frame guide visible
   ↓
5. User captures image → Image stored in app-store
   ↓
6. Navigate to ReviewPicturePage → Preview image
   ↓
7. User confirms → Navigate to ExtractInvoiceDataPage
   ↓
8. Image sent to API → Loading state shown
   ↓
9. Gemini API processes → Extracted data received
   ↓
10. Data displayed → User verifies and edits
   ↓
11. User completes → Data cleared, return to ListInvoicesPage
```

### State Management Flow

```
App Store (Zustand)
├── currentImage: string | null
├── extractedData: InvoiceData | null
├── isLoading: boolean
├── error: string | null
└── Methods:
    ├── setCurrentImage()
    ├── setExtractedData()
    ├── setLoading()
    ├── setError()
    ├── updateKeyValue()
    ├── updateTableCell()
    ├── updateSummary()
    └── clearData()
```

## API Integration

### Endpoints

**POST `/api/extract`**:
- Request: `{ image: string }` (base64 encoded)
- Response: `{ success: boolean, data?: InvoiceData, error?: string }`
- Purpose: Extract invoice data from image

**GET `/api/health`**:
- Response: `{ status: string }`
- Purpose: Health check endpoint

### API Client Pattern

```typescript
const apiClient = new APIClient(API_BASE_URL);
const response = await apiClient.extractInvoice({ image: base64Image });
```

## PWA Configuration

### Service Worker
- Auto-update registration
- Asset caching (JS, CSS, HTML, images)
- Runtime caching for fonts
- Offline fallback to index.html

### Web App Manifest
- Name: "Invoice Scanner"
- Short name: "InvoiceScan"
- Theme color: #3b82f6
- Display: standalone
- Orientation: portrait
- Icons: 72x72 to 512x512 (SVG)

### SSL Support
- Basic SSL plugin for development
- Custom SSL certificate support (ssl/key.pem, ssl/cert.pem)
- HTTPS required for camera access

## Styling Architecture

### Tailwind CSS
- Utility-first approach
- Custom color scheme (primary, background, surface, border)
- Dark mode support
- Responsive breakpoints
- Mobile-first design

### Component Styling
- Base components (Button, Input) with variants
- Consistent spacing and typography
- Accessible color contrasts
- Touch-friendly targets (44x44px minimum)

## Build Configuration

### Vite Configuration
- React plugin
- PWA plugin with manifest and service worker
- Basic SSL plugin (development)
- Path aliases (@/ for src/)
- Code splitting (vendor, router, query chunks)
- Custom SSL certificate support

### TypeScript Configuration
- Strict mode enabled
- Path aliases configured
- React JSX support
- ES2020 target
- Module resolution: bundler

## Development Workflow

### Local Development
```bash
# Install dependencies
cd frontend && npm install

# Start dev server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

### Environment Variables
- `VITE_API_URL`: Backend API URL (default: http://localhost:3000/api)

### Mobile Development
- HTTPS required for camera access
- Use `mobile-dev-http.sh` script for local HTTPS
- Or configure custom SSL certificates in `ssl/` directory

## Testing Strategy

### Current State
- No test suite implemented yet
- Manual testing approach

### Planned Testing
- Unit tests for hooks
- Component tests for pages
- Integration tests for API client
- E2E tests for user flows

## Performance Considerations

### Image Handling
- Base64 encoding for API transmission
- Image compression before upload
- Canvas-based capture and processing

### Code Splitting
- Vendor chunk (React, React DOM)
- Router chunk (React Router)
- Query chunk (TanStack Query)
- Lazy loading for routes (planned)

### Caching Strategy
- Service worker caches static assets
- Runtime caching for fonts
- Network-first strategy for API calls

## Security Considerations

### Current Implementation
- HTTPS required for camera access
- API URL configuration via environment variables
- Input validation in components
- Error handling and sanitization

### Planned Security
- CORS configuration on backend
- API key protection (backend)
- Input sanitization
- XSS prevention

## Backend Configuration

### Configuration Management
**Location**: `backend/pkg/config/`
**Technology**: Viper with YAML and environment variable support

**Configuration Structure** (`backend/pkg/config/config.yaml`):
- Server settings (host, port)
- Gemini API key
- CORS origin
- Database connection (host, port, user, password, name)

**Environment Variable Support**:
- Automatic mapping from dot notation to env vars (e.g., `server.port` → `SERVER_PORT`)
- Supports both YAML config file and environment variables
- Environment variables take precedence over YAML values

**Usage Pattern**:
```go
host := config.GetStringWithDefaultValue("server.host", "localhost")
port := config.GetStringWithDefaultValue("server.port", "3001")
geminiAPIKey := config.GetStringWithDefaultValue("gemini.api_key", "")
```

## Known Limitations

### Current Limitations
- Backend API partially implemented (extraction endpoint working)
- Database persistence implemented but not fully integrated
- No user authentication
- No error tracking service
- No analytics integration

### Browser Support
- Modern browsers with MediaDevices API
- iOS Safari 14+
- Android Chrome 8+
- Desktop browsers (for testing)

## Dependencies Overview

### Production Dependencies
- React ecosystem (react, react-dom, react-router-dom)
- Data fetching (@tanstack/react-query)
- State management (zustand)
- Forms (react-hook-form, @hookform/resolvers, zod)
- UI libraries (@headlessui/react, @heroicons/react)
- Styling (tailwindcss, tailwind-merge, clsx)
- Animations (framer-motion)
- PWA (workbox-* packages)

### Development Dependencies
- Build tools (vite, @vitejs/plugin-react, @vitejs/plugin-basic-ssl)
- PWA plugin (vite-plugin-pwa)
- TypeScript and types
- Linting (eslint, typescript-eslint)
- CSS processing (postcss, autoprefixer, tailwindcss)

## File Statistics

**Frontend Source Files**: ~15 TypeScript/TSX files
**Components**: 4 pages + 2 base UI components
**Hooks**: 4 custom hooks
**Stores**: 1 Zustand store
**Types**: Comprehensive TypeScript definitions

## Related Documentation

- [Project Overview PDR](./project-overview-pdr.md)
- [Code Standards](./code-standards.md)
- [System Architecture](./system-architecture.md)

## Unresolved Questions

1. **Backend Implementation**: When will backend API be implemented?
2. **Testing**: What testing framework to use (Vitest, Jest)?
3. **Error Tracking**: Integration with Sentry or similar?
4. **Analytics**: Google Analytics or custom solution?
5. **Deployment**: CI/CD pipeline setup?
