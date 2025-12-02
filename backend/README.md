# Invoice Scan Backend API

Backend API server for invoice image extraction using Google Gemini Vision API.

## Prerequisites

- Go 1.21+ (or 1.24+ for Gemini SDK)
- Google Gemini API key

## Setup

1. **Install dependencies**:
   ```bash
   go mod download
   ```

2. **Configure environment variables**:
   ```bash
   cp .env.example .env
   ```
   Edit `.env` and set your `GEMINI_API_KEY`.

3. **Run the server**:
   ```bash
   go run cmd/server/main.go
   ```

The server will start on `http://localhost:3001` by default.

## Environment Variables

- `PORT` - Server port (default: 3001)
- `HOST` - Server host (default: localhost)
- `GEMINI_API_KEY` - Google Gemini API key (required)
- `CORS_ORIGIN` - Allowed CORS origin (default: http://localhost:5173)

## API Endpoints

### GET `/api/health`
Health check endpoint.

**Response**:
```json
{
  "status": "ok",
  "timestamp": "2025-12-02T23:30:00Z"
}
```

### POST `/api/extract`
Extract invoice data from invoice image using Gemini Vision API.

**Request**: Multipart form data
- `image` (file): Invoice image file (JPEG, PNG, etc.)
- Max file size: 10MB

**cURL Example**:
```bash
curl -X POST http://localhost:3001/api/extract \
  -F "image=@invoice.jpg"
```

**Response (Success)**:
```json
{
  "success": true,
  "data": {
    "keyValuePairs": [
      {"key": "Invoice Number", "value": "INV-001", "confidence": 0.95},
      {"key": "Date", "value": "2025-12-02", "confidence": 0.90}
    ],
    "table": {
      "headers": ["Item", "Quantity", "Price"],
      "rows": [["Item 1", "2", "100000"]]
    },
    "summary": [
      {"key": "Total", "value": "200000", "confidence": 0.95}
    ]
  },
  "processingTime": 1234
}
```

**Response (Error)**:
```json
{
  "success": false,
  "error": "Error message here"
}
```

**Status Codes**:
- `200 OK` - Success
- `400 Bad Request` - Invalid request or image
- `502 Bad Gateway` - Gemini API error
- `504 Gateway Timeout` - Request timeout
- `500 Internal Server Error` - Server error

## Project Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── internal/
│   ├── config/              # Configuration management
│   ├── handlers/            # HTTP handlers
│   ├── services/            # Business logic
│   ├── models/              # Data models
│   └── middleware/          # HTTP middleware
├── go.mod
├── go.sum
├── .env.example
└── README.md
```

## Development

Run with hot reload (requires air or similar tool):
```bash
air
```

Build binary:
```bash
go build -o bin/server cmd/server/main.go
```

## Testing

Run tests:
```bash
go test ./...
```
