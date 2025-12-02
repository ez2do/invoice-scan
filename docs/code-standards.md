# Code Standards & Codebase Structure

**Last Updated**: 2025-01-26
**Version**: 1.0.0
**Applies To**: All code within Invoice Scan MVP project

## Overview

This document defines coding standards, file organization patterns, naming conventions, and best practices for the Invoice Scan MVP project. All code must adhere to these standards to ensure consistency, maintainability, and quality.

## Core Development Principles

### KISS (Keep It Simple, Stupid)
- Prefer simple, straightforward solutions
- Avoid unnecessary complexity
- Write code that's easy to understand and modify
- Choose clarity over cleverness

### DRY (Don't Repeat Yourself)
- Eliminate code duplication
- Extract common logic into reusable functions/components
- Use composition and abstraction appropriately
- Maintain single source of truth

### Mobile-First
- Design for mobile screens first (320px+)
- Progressive enhancement for larger screens
- Touch-friendly interfaces (44x44px minimum touch targets)
- Optimize for mobile performance

### Type Safety
- Use TypeScript strict mode
- Define types for all data structures
- Avoid `any` type usage
- Leverage TypeScript's type inference where appropriate

## File Organization Standards

### Directory Structure

```
frontend/src/
├── components/          # React components
│   └── ui/            # Base UI components (Button, Input, etc.)
├── hooks/             # Custom React hooks
├── lib/               # Utilities and API client
├── pages/             # Page components (routes)
├── stores/            # Zustand state management stores
├── types/             # TypeScript type definitions
├── assets/            # Static assets (images, fonts)
├── App.tsx           # Main app component
└── main.tsx          # Entry point
```

### File Naming Conventions

**React Components**:
- Format: `PascalCase.tsx`
- Examples: `TakePicturePage.tsx`, `Button.tsx`, `CameraView.tsx`

**Custom Hooks**:
- Format: `camelCase.ts` with `use` prefix
- Examples: `useCamera.ts`, `useInvoiceExtraction.ts`

**Utilities and Helpers**:
- Format: `camelCase.ts`
- Examples: `api.ts`, `utils.ts`, `imageUtils.ts`

**Type Definitions**:
- Format: `camelCase.ts` or `index.ts`
- Examples: `types/index.ts`, `invoice.ts`

**Configuration Files**:
- Format: `kebab-case.config.js` or `camelCase.config.ts`
- Examples: `vite.config.ts`, `tailwind.config.js`, `tsconfig.json`

## Naming Conventions

### Variables & Functions

**JavaScript/TypeScript**:
- **Variables**: camelCase
  ```typescript
  const userName = 'John Doe';
  const isAuthenticated = true;
  const invoiceData = null;
  ```

- **Functions**: camelCase
  ```typescript
  function calculateTotal(items: Item[]) { }
  const extractInvoice = (image: string) => { };
  const handleCapture = () => { };
  ```

- **React Components**: PascalCase
  ```typescript
  function TakePicturePage() { }
  const CameraView = () => { };
  ```

- **Constants**: UPPER_SNAKE_CASE
  ```typescript
  const API_BASE_URL = 'https://api.example.com';
  const MAX_IMAGE_SIZE = 10_000_000;
  ```

- **Private/Internal**: Prefix with underscore (sparingly)
  ```typescript
  class APIClient {
    private _baseURL: string;
    private _handleError(error: Error) { }
  }
  ```

### TypeScript Types & Interfaces

**Interfaces**: PascalCase
```typescript
interface InvoiceData {
  keyValuePairs: KeyValuePair[];
  table: TableData | null;
}

interface ExtractRequest {
  image: string;
}
```

**Types**: PascalCase
```typescript
type AppState = {
  currentImage: string | null;
  isLoading: boolean;
};

type ButtonVariant = 'primary' | 'secondary' | 'ghost';
```

**Generic Types**: Single uppercase letter
```typescript
function useQuery<T>(queryKey: string): QueryResult<T> { }
```

## Code Style Guidelines

### General Formatting

**Indentation**:
- Use 2 spaces (not tabs)
- Consistent indentation throughout file
- No trailing whitespace

**Line Length**:
- Preferred: 80-100 characters
- Hard limit: 120 characters
- Break long lines logically

**Whitespace**:
- One blank line between functions/methods
- Two blank lines between classes/components
- Space after keywords: `if (`, `for (`, `while (`
- No space before function parentheses: `function name(`

### Comments & Documentation

**File Headers** (Optional):
```typescript
/**
 * TakePicturePage Component
 *
 * Handles camera access, image capture, and navigation to review page.
 *
 * @module pages/TakePicturePage
 */
```

**Function Documentation**:
```typescript
/**
 * Captures an image from the video stream and converts it to base64.
 *
 * @param video - HTML video element with active camera stream
 * @returns Base64-encoded image string or null if capture fails
 */
function captureImage(video: HTMLVideoElement): string | null {
  // Implementation
}
```

**Inline Comments**:
- Explain WHY, not WHAT
- Complex logic requires explanation
- TODO comments include context
```typescript
// TODO: Add image quality validation before capture
// Compress image to reduce upload size (target: < 2MB)
const compressed = compressImage(canvas);
```

### React Component Structure

**Component Organization**:
```typescript
// 1. Imports
import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

// 2. Types/Interfaces
interface ComponentProps {
  onCapture: (image: string) => void;
}

// 3. Component
export default function ComponentName({ onCapture }: ComponentProps) {
  // 4. Hooks
  const navigate = useNavigate();
  const [state, setState] = useState(null);

  // 5. Effects
  useEffect(() => {
    // Effect logic
  }, [dependencies]);

  // 6. Event Handlers
  const handleClick = () => {
    // Handler logic
  };

  // 7. Render
  return (
    <div>
      {/* JSX */}
    </div>
  );
}
```

**Component Props**:
- Always define prop types/interfaces
- Use destructuring for props
- Provide default values where appropriate
```typescript
interface ButtonProps {
  variant?: 'primary' | 'secondary';
  onClick?: () => void;
  children: React.ReactNode;
}

function Button({ variant = 'primary', onClick, children }: ButtonProps) {
  // Component implementation
}
```

### Hooks Patterns

**Custom Hook Structure**:
```typescript
export function useCamera() {
  // 1. State
  const [stream, setStream] = useState<MediaStream | null>(null);
  const [error, setError] = useState<string | null>(null);

  // 2. Effects
  useEffect(() => {
    // Setup logic
    return () => {
      // Cleanup logic
    };
  }, []);

  // 3. Functions
  const startCamera = async () => {
    // Implementation
  };

  // 4. Return
  return {
    stream,
    error,
    startCamera,
    stopCamera,
  };
}
```

**Hook Naming**:
- Always start with `use`
- Descriptive names indicating purpose
- Examples: `useCamera`, `useInvoiceExtraction`, `usePWANavigation`

## Error Handling

### Try-Catch Blocks

**Always Use Try-Catch for Async Operations**:
```typescript
async function extractInvoice(image: string) {
  try {
    const response = await apiClient.extractInvoice({ image });
    return response.data;
  } catch (error) {
    console.error('Extraction failed:', error);
    throw new Error('Failed to extract invoice data');
  }
}
```

### Error Types

**Create Custom Error Classes** (when needed):
```typescript
class CameraError extends Error {
  constructor(message: string, public code: string) {
    super(message);
    this.name = 'CameraError';
  }
}

// Usage
throw new CameraError('Camera access denied', 'PERMISSION_DENIED');
```

### Error Logging

**Structured Error Logging**:
```typescript
try {
  // Operation
} catch (error) {
  console.error('Operation failed:', {
    error: error instanceof Error ? error.message : 'Unknown error',
    context: { userId, invoiceId },
    timestamp: new Date().toISOString(),
  });
  throw error;
}
```

## State Management

### Zustand Store Pattern

**Store Structure**:
```typescript
import { create } from 'zustand';

interface AppStore {
  // State
  currentImage: string | null;
  extractedData: InvoiceData | null;
  isLoading: boolean;

  // Actions
  setCurrentImage: (image: string | null) => void;
  setExtractedData: (data: InvoiceData | null) => void;
  setLoading: (loading: boolean) => void;
}

export const useAppStore = create<AppStore>((set) => ({
  // Initial state
  currentImage: null,
  extractedData: null,
  isLoading: false,

  // Actions
  setCurrentImage: (image) => set({ currentImage: image }),
  setExtractedData: (data) => set({ extractedData: data }),
  setLoading: (loading) => set({ isLoading: loading }),
}));
```

### TanStack Query Pattern

**Query/Mutation Setup**:
```typescript
import { useMutation } from '@tanstack/react-query';
import { apiClient } from '@/lib/api';

export function useExtractInvoice() {
  return useMutation({
    mutationFn: (image: string) => apiClient.extractInvoice({ image }),
    onSuccess: (data) => {
      // Handle success
    },
    onError: (error) => {
      // Handle error
    },
  });
}
```

## API Integration

### API Client Pattern

**Structured API Client**:
```typescript
class APIClient {
  private baseURL: string;

  constructor(baseURL: string) {
    this.baseURL = baseURL;
  }

  async extractInvoice(request: ExtractRequest): Promise<ExtractResponse> {
    const response = await fetch(`${this.baseURL}/extract`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(request),
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    return await response.json();
  }
}
```

### Error Handling in API Calls

**Consistent Error Handling**:
```typescript
try {
  const response = await apiClient.extractInvoice({ image });
  if (!response.success) {
    throw new Error(response.error || 'Extraction failed');
  }
  return response.data;
} catch (error) {
  if (error instanceof Error) {
    throw error;
  }
  throw new Error('Unknown error occurred');
}
```

## TypeScript Best Practices

### Type Definitions

**Comprehensive Type Coverage**:
```typescript
// Define types for all data structures
interface InvoiceData {
  keyValuePairs: KeyValuePair[];
  table: TableData | null;
  summary: KeyValuePair[];
}

// Use union types for variants
type ButtonVariant = 'primary' | 'secondary' | 'ghost';

// Use generics for reusable types
interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
}
```

### Avoid `any` Type

**Use Proper Types**:
```typescript
// BAD
function processData(data: any) { }

// GOOD
function processData(data: InvoiceData) { }

// If truly unknown, use unknown
function processData(data: unknown) {
  if (typeof data === 'object' && data !== null) {
    // Type guard
  }
}
```

### Type Guards

**Use Type Guards for Runtime Safety**:
```typescript
function isInvoiceData(data: unknown): data is InvoiceData {
  return (
    typeof data === 'object' &&
    data !== null &&
    'keyValuePairs' in data &&
    Array.isArray(data.keyValuePairs)
  );
}
```

## Styling Standards

### Tailwind CSS Usage

**Utility-First Approach**:
```typescript
// Use Tailwind utilities directly
<div className="flex items-center justify-between p-4 bg-white dark:bg-gray-800">

// Extract repeated patterns to components
function Card({ children }: { children: React.ReactNode }) {
  return (
    <div className="rounded-lg border border-gray-200 p-4 shadow-sm">
      {children}
    </div>
  );
}
```

### Dark Mode Support

**Always Consider Dark Mode**:
```typescript
<div className="bg-white dark:bg-gray-800 text-gray-900 dark:text-white">
```

### Responsive Design

**Mobile-First Breakpoints**:
```typescript
// Mobile first, then enhance for larger screens
<div className="w-full md:w-1/2 lg:w-1/3">
```

## Testing Standards (Planned)

### Test File Organization

```
frontend/src/
├── components/
│   └── Button.test.tsx
├── hooks/
│   └── useCamera.test.ts
└── lib/
    └── api.test.ts
```

### Test Naming

```typescript
describe('useCamera', () => {
  it('should start camera stream when called', async () => {
    // Test implementation
  });

  it('should handle camera permission denial', async () => {
    // Test implementation
  });
});
```

## Git Standards

### Commit Messages

**Format**: Conventional Commits
```
type(scope): description

[optional body]

[optional footer]
```

**Types**:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Test additions/changes
- `chore`: Maintenance tasks
- `perf`: Performance improvements

**Examples**:
```
feat(camera): add image compression before upload

fix(extraction): handle API timeout errors

docs: update README with setup instructions

refactor(store): simplify app store structure
```

### Branch Naming

**Format**: `type/description`

**Types**:
- `feature/` - New features
- `fix/` - Bug fixes
- `refactor/` - Code refactoring
- `docs/` - Documentation updates

**Examples**:
```
feature/camera-capture
fix/extraction-error-handling
refactor/state-management
```

## Security Standards

### Input Validation

**Validate All Inputs**:
```typescript
function validateImage(base64: string): boolean {
  if (!base64.startsWith('data:image/')) {
    return false;
  }
  if (base64.length > 10_000_000) {
    return false;
  }
  return true;
}
```

### Sensitive Data Handling

**Never Log Sensitive Data**:
```typescript
// BAD
console.log('Image data:', imageBase64);

// GOOD
console.log('Image captured, size:', imageBase64.length);
```

### Environment Variables

**Use Environment Variables for Configuration**:
```typescript
const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:3000/api';
```

## Performance Standards

### Code Splitting

**Lazy Load Routes** (when implemented):
```typescript
import { lazy } from 'react';

const ExtractInvoiceDataPage = lazy(() => import('./pages/ExtractInvoiceDataPage'));
```

### Image Optimization

**Compress Images Before Upload**:
```typescript
function compressImage(canvas: HTMLCanvasElement, quality = 0.8): string {
  return canvas.toDataURL('image/jpeg', quality);
}
```

### Memoization

**Use React.memo and useMemo Appropriately**:
```typescript
const MemoizedComponent = React.memo(ExpensiveComponent);

const expensiveValue = useMemo(() => {
  return computeExpensiveValue(data);
}, [data]);
```

## Documentation Standards

### Code Documentation

**Self-Documenting Code**:
- Clear variable and function names
- Logical code organization
- Minimal comments needed

**When to Comment**:
- Complex algorithms or business logic
- Non-obvious optimizations
- Workarounds for browser limitations
- Public API functions

### README Files

**Component README Structure** (when needed):
```markdown
# ComponentName

Brief description of component purpose.

## Props

| Prop | Type | Required | Description |
|------|------|----------|-------------|
| onCapture | function | Yes | Callback when image is captured |

## Usage

\`\`\`tsx
<ComponentName onCapture={handleCapture} />
\`\`\`
```

## Quality Assurance

### Pre-Commit Checklist

- ✅ No console.log statements (use proper logging)
- ✅ No TypeScript errors
- ✅ No ESLint errors
- ✅ Code follows style guidelines
- ✅ Types are properly defined
- ✅ Error handling is implemented
- ✅ Mobile-responsive design verified

### Code Review Focus

- Functionality correctness
- Type safety
- Error handling
- Performance considerations
- Accessibility
- Mobile responsiveness
- Code clarity and maintainability

## References

### Internal Documentation
- [Project Overview PDR](./project-overview-pdr.md)
- [Codebase Summary](./codebase-summary.md)
- [System Architecture](./system-architecture.md)

### External Standards
- [TypeScript Handbook](https://www.typescriptlang.org/docs/)
- [React Documentation](https://react.dev/)
- [Tailwind CSS Documentation](https://tailwindcss.com/docs)
- [Conventional Commits](https://www.conventionalcommits.org/)

## Unresolved Questions

1. **Testing Framework**: Vitest vs Jest?
2. **E2E Testing**: Playwright vs Cypress?
3. **Code Coverage**: Target percentage?
4. **Linting Rules**: Additional ESLint rules needed?
