# Docker Compose Setup

This document provides instructions for running the Invoice Scan application using Docker Compose.

## Overview

The project uses separate Docker Compose configurations for backend and frontend:

- **Backend dependencies**: `backend/deployment/docker-compose.dev.yml` - MySQL database
- **Frontend**: `frontend/docker-compose.yml` - Frontend development server

## Prerequisites

- Docker Engine 20.10+
- Docker Compose 2.0+
- Google Gemini API key (for backend)

## Quick Start

### Option 1: Separate Services (Recommended for Development)

**Start Backend Dependencies (MySQL)**:
```bash
cd backend/deployment
cp env.example .env  # Edit if needed
docker-compose -f docker-compose.dev.yml up -d
```

**Start Frontend**:
```bash
cd frontend
cp env.example .env  # Edit if needed
docker-compose up -d
```

**Run Backend Locally**:
```bash
cd backend
# Set environment variables or use config.yaml
go run cmd/server/main.go
```

### Option 2: Full Stack (Root docker-compose.yml)

For running everything together, use the root `docker-compose.yml`:

1. **Copy environment file**:
   ```bash
   cp docker-compose.env.example .env
   ```

2. **Configure environment variables**:
   Edit `.env` and set your `GEMINI_API_KEY`:
   ```env
   GEMINI_API_KEY=your_actual_api_key_here
   ```

3. **Start all services**:
   ```bash
   docker-compose up -d
   ```

4. **Access the application**:
   - Frontend: http://localhost:5173
   - Backend API: http://localhost:3001/api
   - Health check: http://localhost:3001/api/health

## Project Structure

```
invoice-scan/
├── docker-compose.yml              # Full stack (all services)
├── docker-compose.env.example      # Root environment template
├── backend/
│   ├── deployment/
│   │   ├── docker-compose.dev.yml  # MySQL only (dependencies)
│   │   ├── .env.example           # Backend env template
│   │   └── README.md              # Backend deployment docs
│   └── Dockerfile                  # Backend container
└── frontend/
    ├── docker-compose.yml          # Frontend service
    ├── .env.example                # Frontend env template
    ├── DOCKER.md                   # Frontend Docker docs
    └── Dockerfile                  # Frontend container
```

## Services

### MySQL Database (Backend Dependencies)
- **Location**: `backend/deployment/docker-compose.dev.yml`
- **Container**: `invoice-scan-mysql`
- **Port**: 3306 (configurable via `DB_PORT`)
- **Database**: `invoice_scan` (configurable via `DB_NAME`)
- **User**: `invoice_user` (configurable via `DB_USER`)
- **Password**: Set via `DB_PASSWORD` in `.env`
- **Data persistence**: Stored in `mysql_data` volume

### Backend API
- **Dockerfile**: `backend/Dockerfile`
- **Port**: 3001 (configurable via `BACKEND_PORT`)
- **Health check**: `/api/health`
- **Uploads**: Stored in `backend_uploads` volume (when using root docker-compose)

### Frontend
- **Location**: `frontend/docker-compose.yml`
- **Container**: `invoice-scan-frontend`
- **Port**: 5173 (configurable via `FRONTEND_PORT`)
- **Hot reload**: Enabled for development

## Environment Variables

Docker Compose reads environment variables from:
1. **`.env` file** in the same directory as the `docker-compose.yml` file (recommended)
2. Shell environment variables (exported in your terminal)

**Where to define variables:**

- **Backend dependencies**: `backend/deployment/.env` (copy from `backend/deployment/env.example`)
- **Frontend**: `frontend/.env` (copy from `frontend/env.example`)
- **Full stack**: Root `.env` (copy from `docker-compose.env.example`)

**Important**: `.env` files are gitignored for security. Never commit passwords or API keys.

All environment variables can be configured in `.env` files:

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_HOST` | `mysql` | MySQL hostname |
| `DB_PORT` | `3306` | MySQL port |
| `DB_USER` | `invoice_user` | MySQL username |
| `DB_PASSWORD` | `rootpassword` | MySQL password |
| `DB_NAME` | `invoice_scan` | Database name |
| `BACKEND_PORT` | `3001` | Backend API port |
| `GEMINI_API_KEY` | - | **Required** - Google Gemini API key |
| `CORS_ORIGIN` | `http://localhost:5173` | Allowed CORS origin |
| `STORAGE_BASE_URL` | `http://localhost:3001` | Base URL for file storage |
| `FRONTEND_PORT` | `5173` | Frontend dev server port |
| `VITE_API_URL` | `http://localhost:3001/api` | Backend API URL for frontend |

## Common Commands

### Backend Dependencies (MySQL)

```bash
cd backend/deployment

# Start MySQL
docker-compose -f docker-compose.dev.yml up -d

# Stop MySQL
docker-compose -f docker-compose.dev.yml down

# View logs
docker-compose -f docker-compose.dev.yml logs -f mysql

# Connect to MySQL
docker-compose -f docker-compose.dev.yml exec mysql mysql -u invoice_user -p invoice_scan
```

### Frontend

```bash
cd frontend

# Start frontend
docker-compose up -d

# Stop frontend
docker-compose down

# View logs
docker-compose logs -f

# Rebuild
docker-compose build
```

### Full Stack (Root)

```bash
# Start all services
docker-compose up -d

# Stop all services
docker-compose down

# View logs
docker-compose logs -f

# Rebuild specific service
docker-compose build backend
docker-compose build frontend
```

## Database Migrations

Migrations are located in `backend/db/migrations/`. The backend does not automatically run migrations on startup. You need to run them manually:

**Option 1: Using sql-migrate CLI (if installed locally)**
```bash
cd backend
export DB_USER=invoice_user
export DB_PASSWORD=rootpassword
export DB_HOST=localhost
export DB_PORT=3306
export DB_NAME=invoice_scan
sql-migrate up -config=db/dbconfig.yml
```

**Option 2: Using Docker exec**
```bash
# First, install sql-migrate in the container
docker-compose exec backend sh -c "go install github.com/rubenv/sql-migrate/...@latest"

# Then run migrations
docker-compose exec backend sh -c "cd /app && sql-migrate up -config=db/dbconfig.yml"
```

**Option 3: Create database manually**
If migrations fail, you can create the database and table manually:
```bash
docker-compose exec mysql mysql -u invoice_user -prootpassword invoice_scan < backend/db/migrations/1764948381_create_invoices.sql
```

**Note**: Migrations only need to be run once when setting up the database for the first time.

## Troubleshooting

### Backend fails to start
- Check MySQL is healthy: `docker-compose ps`
- Verify `GEMINI_API_KEY` is set in `.env`
- Check backend logs: `docker-compose logs backend`

### Frontend can't connect to backend
- Verify `VITE_API_URL` matches backend URL
- Check CORS configuration in backend
- Ensure both services are on the same network

### Database connection errors
- Wait for MySQL health check to pass (may take 30-60 seconds)
- Verify database credentials in `.env`
- Check MySQL logs: `docker-compose logs mysql`

### Port conflicts
- Change ports in `.env` file
- Ensure ports are not in use: `lsof -i :3001` or `lsof -i :5173`

## Production Considerations

For production deployment:

1. **Use production builds**:
   - Frontend: Build static files and serve with nginx
   - Backend: Use optimized Go binary

2. **Security**:
   - Use strong database passwords
   - Set proper CORS origins
   - Use HTTPS for all services
   - Store secrets securely (Docker secrets, vault, etc.)

3. **Performance**:
   - Use production MySQL configuration
   - Configure connection pooling
   - Set appropriate resource limits

4. **Monitoring**:
   - Add health check endpoints
   - Set up logging aggregation
   - Monitor resource usage

## Development vs Production

The current setup is optimized for **development**:
- Frontend runs in dev mode with hot reload
- Backend runs with debug logging
- Volumes are mounted for live code changes

For production, consider:
- Multi-stage builds for smaller images
- Production-optimized frontend build
- Separate compose file for production (`docker-compose.prod.yml`)
- Reverse proxy (nginx/traefik)
- SSL/TLS termination

