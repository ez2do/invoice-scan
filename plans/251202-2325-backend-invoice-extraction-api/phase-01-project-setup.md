# Phase 1: Project Setup & Configuration

**Date**: 2025-12-02  
**Priority**: High  
**Status**: Pending  
**Estimated Time**: 1-2 hours

## Context Links

- **Parent Plan**: `plan.md`
- **Research**: `research/researcher-01-go-gin-gemini.md`
- **Dependencies**: None
- **Reference Docs**: 
  - `docs/system-architecture.md`
  - `docs/code-standards.md`
  - `docs/project-overview-pdr.md`

## Overview

Set up Go backend project structure with proper directory organization, initialize Go modules, configure dependencies, and set up environment variable management. Establish foundation for API implementation.

## Key Insights

- Use standard Go project layout (`cmd/`, `internal/`, `pkg/`)
- Gin framework for HTTP server
- Google Gemini Go SDK for vision API
- Environment variables for configuration
- No database needed for MVP

## Requirements

### Functional Requirements
- Go project initialized with proper structure
- Dependencies configured (Gin, Gemini SDK, godotenv)
- Environment variable loading
- Configuration struct defined
- Basic server setup

### Non-Functional Requirements
- Follow Go best practices
- Clean project structure
- Environment-based configuration
- No hardcoded values

## Architecture

```
backend/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── handlers/
│   ├── services/
│   ├── models/
│   └── middleware/
├── go.mod
├── go.sum
├── .env.example
└── README.md
```

## Related Code Files

### Files to Create
- `backend/cmd/server/main.go` - Application entry point
- `backend/internal/config/config.go` - Configuration management
- `backend/go.mod` - Go module definition
- `backend/.env.example` - Environment variables template
- `backend/README.md` - Backend documentation
- `backend/.gitignore` - Git ignore rules

## Implementation Steps

1. **Create backend directory structure**
   - Create `backend/cmd/server/` directory
   - Create `backend/internal/config/` directory
   - Create `backend/internal/handlers/` directory
   - Create `backend/internal/services/` directory
   - Create `backend/internal/models/` directory
   - Create `backend/internal/middleware/` directory

2. **Initialize Go module**
   - Run `go mod init` in backend directory
   - Set module name (e.g., `invoice-scan/backend`)

3. **Install dependencies**
   - `go get github.com/gin-gonic/gin`
   - `go get google.generativeai/go`
   - `go get github.com/joho/godotenv`
   - `go get github.com/gin-contrib/cors`

4. **Create configuration package**
   - Define Config struct with fields:
     - Port (string)
     - Host (string)
     - GeminiAPIKey (string)
     - CORSOrigin (string)
   - Load from environment variables
   - Provide default values

5. **Create main.go**
   - Initialize Gin router
   - Load configuration
   - Set up basic server
   - Add placeholder routes

6. **Create .env.example**
   - PORT=3001
   - HOST=localhost
   - GEMINI_API_KEY=your_api_key_here
   - CORS_ORIGIN=http://localhost:5173

7. **Create .gitignore**
   - Ignore `.env` file
   - Ignore Go build artifacts
   - Ignore IDE files

8. **Create README.md**
   - Setup instructions
   - Environment variables
   - Running instructions

## Todo List

- [ ] Create backend directory structure
- [ ] Initialize Go module
- [ ] Install dependencies
- [ ] Create config package
- [ ] Create main.go with basic setup
- [ ] Create .env.example
- [ ] Create .gitignore
- [ ] Create README.md
- [ ] Verify project builds successfully

## Success Criteria

- Go project initializes without errors
- All dependencies install successfully
- Configuration loads from environment variables
- Server starts and responds to requests
- Project structure follows Go best practices
- No hardcoded configuration values

## Risk Assessment

| Risk | Impact | Mitigation |
|------|--------|------------|
| Wrong Go version | Medium | Specify Go 1.21+ requirement |
| Dependency conflicts | Low | Use latest stable versions |
| Missing API key | High | Document in README and .env.example |

## Security Considerations

- API key stored in environment variables only
- .env file excluded from git
- No secrets in code or config files
- Environment variable validation

## Next Steps

- Proceed to Phase 2: Extract Endpoint Implementation
- Requires: Completed Phase 1, Gemini API key configured
