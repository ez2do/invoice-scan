---
phase: implementation
title: Invoice Scan MVP - Implementation Guide
description: Technical implementation notes, patterns, and code guidelines for the invoice scanning MVP
---

# Invoice Scan MVP - Implementation Guide

## Development Setup

**How do we get started?**

### Prerequisites

- **Node.js**: 18+ (for pnpm and frontend)
- **pnpm**: Latest version (`npm install -g pnpm`)
- **Go**: 1.21+ (for backend)
- **Docker**: Optional, for containerized development
- **Google Cloud Account**: With Gemini API access and API key

### Initial Setup

1. **Clone and install dependencies:**
```bash
# Install root dependencies (if any)
pnpm install

# Install frontend dependencies
cd frontend && pnpm install

# Install backend dependencies
cd ../backend && go mod download
```

2. **Configure environment variables:**

**Backend `.env`** (create `backend/.env`):
```env
PORT=3001
GEMINI_API_KEY=your_gemini_api_key_here
ALLOWED_ORIGINS=http://localhost:5173
LOG_LEVEL=info
```

**Frontend `.env`** (create `frontend/.env`):
```env
VITE_API_URL=http://localhost:3001
```

3. **Run development servers:**

**Terminal 1 - Backend:**
```bash
cd backend
go run cmd/server/main.go
```

**Terminal 2 - Frontend:**
```bash
cd frontend
pnpm dev
```

### Getting Gemini API Key

1. Go to [Google AI Studio](https://aistudio.google.com/)
2. Sign in with Google account
3. Click "Get API Key"
4. Create new API key or use existing
5. Copy to `backend/.env`

## Code Structure

**How is the code organized?**

### Frontend Structure

```
frontend/
├── src/
│   ├── components/          # React components
│   │   ├── camera/         # Camera-related components
│   │   ├── verify/         # Verification screen components
│   │   └── common/         # Shared UI components
│   ├── hooks/              # Custom React hooks
│   │   ├── useCamera.ts
│   │   └── useExtractInvoice.ts
│   ├── stores/              # Zustand stores
│   │   └── invoiceStore.ts
│   ├── services/           # API client and services
│   │   └── api.ts
│   ├── types/              # TypeScript type definitions
│   │   └── invoice.ts
│   ├── utils/              # Utility functions
│   │   └── imageUtils.ts
│   ├── App.tsx             # Main app component with routing
│   └── main.tsx            # Entry point
├── public/                 # Static assets
│   ├── icons/             # PWA icons
│   └── manifest.json      # PWA manifest
├── index.html
├── package.json
├── vite.config.ts
├── tailwind.config.js
└── tsconfig.json
```

### Backend Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go         # Application entry point
├── internal/
│   ├── handlers/          # HTTP request handlers
│   │   ├── extract.go     # POST /api/extract
│   │   └── health.go      # GET /api/health
│   ├── services/          # Business logic
│   │   └── gemini.go      # Gemini API integration
│   ├── models/            # Data models
│   │   └── invoice.go
│   ├── config/            # Configuration
│   │   └── config.go
│   └── middleware/        # HTTP middleware
│       ├── cors.go
│       └── logger.go
├── pkg/                   # Public packages (if needed)
├── go.mod
├── go.sum
└── Dockerfile
```

### Naming Conventions

| Type | Convention | Example |
|------|------------|---------|
| **Frontend Components** | PascalCase | `CameraView.tsx` |
| **Frontend Hooks** | camelCase with `use` prefix | `useCamera.ts` |
| **Frontend Files** | kebab-case or camelCase | `api.ts`, `image-utils.ts` |
| **Go Packages** | lowercase, single word | `handlers`, `services` |
| **Go Files** | snake_case | `extract.go`, `invoice.go` |
| **Go Types** | PascalCase | `InvoiceData`, `KeyValuePair` |
| **Go Functions** | PascalCase (exported), camelCase (private) | `ExtractInvoice`, `parseResponse` |

## Implementation Notes

**Key technical details to remember:**

### Core Features

#### Feature 1: Camera Capture

**Frontend Implementation:**

Use MediaDevices API with proper error handling:

```typescript
const constraints = {
  video: {
    facingMode: 'environment',
    width: { ideal: 1920 },
    height: { ideal: 1080 }
  }
};

const stream = await navigator.mediaDevices.getUserMedia(constraints);
```

**Capture to base64:**
```typescript
const captureImage = (video: HTMLVideoElement): string => {
  const canvas = document.createElement('canvas');
  canvas.width = video.videoWidth;
  canvas.height = video.videoHeight;
  const ctx = canvas.getContext('2d');
  ctx?.drawImage(video, 0, 0);
  
  // Compress to reduce size
  return canvas.toDataURL('image/jpeg', 0.8);
};
```

**Cleanup:**
```typescript
useEffect(() => {
  return () => {
    stream?.getTracks().forEach(track => track.stop());
  };
}, [stream]);
```

#### Feature 2: Gemini Vision Integration

**Backend Implementation (Go):**

Use Google's Generative AI Go SDK:

```go
import "github.com/google/generative-ai-go/genai"

func ExtractInvoice(ctx context.Context, imageBase64 string) (*models.InvoiceData, error) {
    client, err := genai.NewClient(ctx, option.WithAPIKey(cfg.GeminiAPIKey))
    if err != nil {
        return nil, err
    }
    defer client.Close()

    model := client.GenerativeModel("gemini-1.5-flash")
    
    // Decode base64 image
    imageData, err := base64.StdEncoding.DecodeString(imageBase64)
    if err != nil {
        return nil, err
    }

    // Create image part
    imgPart := &genai.Blob{
        MimeType: "image/jpeg",
        Data:     imageData,
    }

    // Generate content
    resp, err := model.GenerateContent(ctx, imgPart, genai.Text(extractionPrompt))
    if err != nil {
        return nil, err
    }

    // Parse response
    return parseGeminiResponse(resp)
}
```

**Prompt Design:**
Store prompt in `internal/services/prompt.go` or as a constant. See design doc for full prompt.

#### Feature 3: Flexible Data Display

**Frontend - Key-Value Section:**
```typescript
const KeyValueSection = ({ pairs, editable }: Props) => {
  return (
    <div>
      {pairs.map((pair, idx) => (
        <EditableField
          key={idx}
          label={pair.key}
          value={pair.value}
          onChange={(value) => updatePair(idx, value)}
          editable={editable}
        />
      ))}
    </div>
  );
};
```

**Frontend - Dynamic Table:**
```typescript
const TableSection = ({ table, editable }: Props) => {
  return (
    <table>
      <thead>
        <tr>
          {table.headers.map((header, idx) => (
            <th key={idx}>{header}</th>
          ))}
        </tr>
      </thead>
      <tbody>
        {table.rows.map((row) => (
          <tr key={row.id}>
            {row.cells.map((cell, idx) => (
              <td key={idx}>
                <EditableField
                  value={cell}
                  onChange={(value) => updateCell(row.id, idx, value)}
                  editable={editable}
                />
              </td>
            ))}
          </tr>
        ))}
      </tbody>
    </table>
  );
};
```

#### Feature 4: State Management with Zustand

**Store Definition:**
```typescript
import { create } from 'zustand';

interface InvoiceStore {
  invoiceData: InvoiceData | null;
  setInvoiceData: (data: InvoiceData) => void;
  updateKeyValuePair: (section: 'keyValuePairs' | 'summary', index: number, value: string) => void;
  updateTableCell: (rowId: string, cellIndex: number, value: string) => void;
}

export const useInvoiceStore = create<InvoiceStore>((set) => ({
  invoiceData: null,
  setInvoiceData: (data) => set({ invoiceData: data }),
  updateKeyValuePair: (section, index, value) =>
    set((state) => ({
      invoiceData: {
        ...state.invoiceData!,
        [section]: state.invoiceData![section].map((pair, i) =>
          i === index ? { ...pair, value } : pair
        ),
      },
    })),
  // ... other update methods
}));
```

#### Feature 5: API Integration with TanStack Query

**Query Setup:**
```typescript
import { useMutation } from '@tanstack/react-query';
import { extractInvoice } from '@/services/api';

export const useExtractInvoice = () => {
  return useMutation({
    mutationFn: (imageBase64: string) => extractInvoice(imageBase64),
    onSuccess: (data) => {
      useInvoiceStore.getState().setInvoiceData(data);
    },
  });
};
```

**API Client:**
```typescript
const API_URL = import.meta.env.VITE_API_URL;

export const extractInvoice = async (imageBase64: string): Promise<InvoiceData> => {
  const response = await fetch(`${API_URL}/api/extract`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ image: imageBase64 }),
  });

  if (!response.ok) {
    throw new Error('Extraction failed');
  }

  const result = await response.json();
  if (!result.success) {
    throw new Error(result.error || 'Extraction failed');
  }

  return result.data;
};
```

### Patterns & Best Practices

#### Frontend Patterns

**Component Structure:**
```typescript
// Component with props interface
interface CameraViewProps {
  onCapture: (image: string) => void;
}

export const CameraView = ({ onCapture }: CameraViewProps) => {
  // Component logic
};
```

**Custom Hooks Pattern:**
```typescript
export const useCamera = () => {
  const [stream, setStream] = useState<MediaStream | null>(null);
  const [error, setError] = useState<string | null>(null);

  const startCamera = async () => {
    try {
      const mediaStream = await navigator.mediaDevices.getUserMedia({
        video: { facingMode: 'environment' }
      });
      setStream(mediaStream);
    } catch (err) {
      setError('Camera access denied');
    }
  };

  useEffect(() => {
    return () => {
      stream?.getTracks().forEach(track => track.stop());
    };
  }, [stream]);

  return { stream, error, startCamera };
};
```

#### Backend Patterns

**Handler Pattern (Gin):**
```go
func ExtractInvoice(c *gin.Context) {
    var req ExtractRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid request"})
        return
    }

    // Validate image
    if len(req.Image) > 10_000_000 {
        c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Image too large"})
        return
    }

    // Extract invoice
    invoiceData, err := geminiService.ExtractInvoice(c.Request.Context(), req.Image)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Extraction failed"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"success": true, "data": invoiceData})
}
```

**Service Pattern:**
```go
type GeminiService struct {
    client *genai.Client
    apiKey string
}

func NewGeminiService(apiKey string) (*GeminiService, error) {
    ctx := context.Background()
    client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
    if err != nil {
        return nil, err
    }
    return &GeminiService{client: client, apiKey: apiKey}, nil
}

func (s *GeminiService) ExtractInvoice(ctx context.Context, imageBase64 string) (*models.InvoiceData, error) {
    // Implementation
}
```

**Error Handling Pattern:**
```go
if err != nil {
    log.Printf("Error extracting invoice: %v", err)
    return nil, fmt.Errorf("failed to extract invoice: %w", err)
}
```

## Integration Points

**How do pieces connect?**

### Frontend → Backend API

**API Client Setup:**
```typescript
// frontend/src/services/api.ts
const API_URL = import.meta.env.VITE_API_URL;

export const api = {
  extract: (image: string) => 
    fetch(`${API_URL}/api/extract`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ image }),
    }),
  
  health: () => 
    fetch(`${API_URL}/api/health`),
};
```

### Backend → Gemini API

**Gemini Client Setup:**
```go
// backend/internal/services/gemini.go
import (
    "github.com/google/generative-ai-go/genai"
    "google.golang.org/api/option"
)

func NewGeminiClient(ctx context.Context, apiKey string) (*genai.Client, error) {
    return genai.NewClient(ctx, option.WithAPIKey(apiKey))
}
```

### CORS Configuration

**Backend Middleware:**
```go
// backend/internal/middleware/cors.go
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        origin := c.GetHeader("Origin")
        allowedOrigins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
        
        for _, allowed := range allowedOrigins {
            if origin == allowed {
                c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
                break
            }
        }
        
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        
        c.Next()
    }
}
```

## Error Handling

**How do we handle failures?**

### Frontend Error Handling

**TanStack Query Error Handling:**
```typescript
const { mutate, error, isError } = useExtractInvoice();

useEffect(() => {
  if (isError) {
    // Show error toast/notification
    toast.error(error?.message || 'Failed to extract invoice');
  }
}, [isError, error]);
```

**Camera Error Handling:**
```typescript
const handleCameraError = (error: Error) => {
  if (error.name === 'NotAllowedError') {
    setError('Camera permission denied. Please enable camera access.');
  } else if (error.name === 'NotFoundError') {
    setError('No camera found on this device.');
  } else {
    setError('Failed to access camera.');
  }
};
```

### Backend Error Handling

**Structured Error Responses:**
```go
type ErrorResponse struct {
    Success bool   `json:"success"`
    Error   string `json:"error"`
}

func handleError(c *gin.Context, statusCode int, message string) {
    c.JSON(statusCode, ErrorResponse{
        Success: false,
        Error:   message,
    })
}
```

**Logging:**
```go
import "log"

log.Printf("Error: %v", err)
// Or use structured logging with log/slog (Go 1.21+)
```

**Timeout Handling:**
```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

invoiceData, err := geminiService.ExtractInvoice(ctx, imageBase64)
if err != nil {
    if err == context.DeadlineExceeded {
        return nil, fmt.Errorf("extraction timeout")
    }
    return nil, err
}
```

## Performance Considerations

**How do we keep it fast?**

### Frontend Optimizations

**Image Compression:**
```typescript
const compressImage = (canvas: HTMLCanvasElement, maxSize = 2048): string => {
  const scale = Math.min(1, maxSize / Math.max(canvas.width, canvas.height));
  
  const resized = document.createElement('canvas');
  resized.width = canvas.width * scale;
  resized.height = canvas.height * scale;
  
  const ctx = resized.getContext('2d');
  ctx?.drawImage(canvas, 0, 0, resized.width, resized.height);
  
  return resized.toDataURL('image/jpeg', 0.8);
};
```

**Code Splitting:**
```typescript
// Lazy load routes
const VerifyView = lazy(() => import('./components/verify/VerifyView'));

<Suspense fallback={<LoadingSpinner />}>
  <Routes>
    <Route path="/verify" element={<VerifyView />} />
  </Routes>
</Suspense>
```

**PWA Caching:**
```typescript
// vite.config.ts
import { VitePWA } from 'vite-plugin-pwa';

export default defineConfig({
  plugins: [
    VitePWA({
      registerType: 'autoUpdate',
      workbox: {
        globPatterns: ['**/*.{js,css,html,ico,png,svg,woff2}'],
        runtimeCaching: [
          {
            urlPattern: /^https:\/\/api\./,
            handler: 'NetworkFirst',
          },
        ],
      },
    }),
  ],
});
```

### Backend Optimizations

**Connection Pooling:**
Gemini client handles connection pooling internally. Reuse client instance:

```go
var geminiClient *genai.Client

func initGeminiClient() error {
    ctx := context.Background()
    client, err := genai.NewClient(ctx, option.WithAPIKey(cfg.GeminiAPIKey))
    if err != nil {
        return err
    }
    geminiClient = client
    return nil
}
```

**Concurrent Processing:**
Go's goroutines allow concurrent request handling automatically with Gin.

**Response Compression:**
```go
import "github.com/gin-gonic/gin"

// Enable gzip compression
r.Use(gin.Recovery())
```

## Security Notes

**What security measures are in place?**

### API Key Protection

**Backend:**
- Store API key in environment variable
- Never expose to client
- Use `.env` file (not committed to git)

```go
// backend/internal/config/config.go
type Config struct {
    GeminiAPIKey string
    Port         string
}

func LoadConfig() (*Config, error) {
    apiKey := os.Getenv("GEMINI_API_KEY")
    if apiKey == "" {
        return nil, fmt.Errorf("GEMINI_API_KEY not set")
    }
    return &Config{
        GeminiAPIKey: apiKey,
        Port:         getEnv("PORT", "3001"),
    }, nil
}
```

### Input Validation

**Backend:**
```go
func validateImage(imageBase64 string) error {
    // Check size (remove data URL prefix first)
    imageData := strings.TrimPrefix(imageBase64, "data:image/jpeg;base64,")
    if len(imageData) > 10_000_000 {
        return fmt.Errorf("image too large")
    }
    
    // Validate base64 format
    if _, err := base64.StdEncoding.DecodeString(imageData); err != nil {
        return fmt.Errorf("invalid base64 format")
    }
    
    return nil
}
```

**Frontend:**
```typescript
const validateImage = (base64: string): boolean => {
  if (!base64.startsWith('data:image/')) {
    return false;
  }
  if (base64.length > 10_000_000) {
    return false;
  }
  return true;
};
```

### CORS Configuration

Restrict to known origins only:

```go
allowedOrigins := []string{
    "http://localhost:5173",
    "https://yourdomain.com",
}
```

### HTTPS Requirement

- Required for camera access in production
- Use Let's Encrypt for SSL certificates
- Configure reverse proxy (nginx) for HTTPS termination

### No Data Persistence

- Images processed in-memory only
- No database storage
- No file system writes
- Privacy benefit: data not stored

## Development Workflow

**Daily development commands:**

```bash
# Start backend
cd backend && go run cmd/server/main.go

# Start frontend
cd frontend && pnpm dev

# Run tests (when implemented)
cd backend && go test ./...
cd frontend && pnpm test

# Build for production
cd backend && go build -o bin/server cmd/server/main.go
cd frontend && pnpm build
```

## Docker Development (Optional)

**docker-compose.yml:**
```yaml
version: '3.8'

services:
  backend:
    build: ./backend
    ports:
      - "3001:3001"
    environment:
      - GEMINI_API_KEY=${GEMINI_API_KEY}
      - PORT=3001
    volumes:
      - ./backend:/app

  frontend:
    build: ./frontend
    ports:
      - "5173:5173"
    volumes:
      - ./frontend:/app
      - /app/node_modules
```

