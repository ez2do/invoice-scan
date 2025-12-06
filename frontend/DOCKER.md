# Frontend Docker Setup

This document provides instructions for running the frontend using Docker Compose.

## Prerequisites

- Docker Engine 20.10+
- Docker Compose 2.0+

## Quick Start

1. **Copy environment file**:
   ```bash
   cp env.example .env
   ```

2. **Configure environment variables** (optional):
   Docker Compose reads variables from `.env` file in the same directory.
   Edit `.env` if you need to change defaults:
   ```env
   FRONTEND_PORT=5173
   VITE_API_URL=http://localhost:3001/api
   ```

3. **Start frontend service**:
   ```bash
   docker-compose up -d
   ```

4. **View logs**:
   ```bash
   docker-compose logs -f
   ```

5. **Access the application**:
   - Frontend: http://localhost:5173

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `FRONTEND_PORT` | `5173` | Frontend dev server port |
| `VITE_API_URL` | `http://localhost:3001/api` | Backend API URL |

## Common Commands

### Start service
```bash
docker-compose up -d
```

### Stop service
```bash
docker-compose down
```

### View logs
```bash
# All logs
docker-compose logs -f

# Follow logs
docker-compose logs -f frontend
```

### Rebuild service
```bash
docker-compose build
docker-compose up -d
```

### Execute commands in container
```bash
# Shell access
docker-compose exec frontend sh

# Run npm commands
docker-compose exec frontend npm install
docker-compose exec frontend npm run build
```

### Clean up
```bash
# Stop and remove container
docker-compose down

# Remove volumes
docker-compose down -v
```

## Development

The frontend runs in development mode with:
- Hot module replacement (HMR)
- Source maps
- Live reloading
- Volume mounting for code changes

Code changes are reflected immediately without rebuilding the container.

## Production

For production builds, build the static files and serve with a web server:

```bash
# Build production files
docker-compose exec frontend npm run build

# Files will be in ./dist directory
```

## Troubleshooting

### Port already in use
Change `FRONTEND_PORT` in `.env` file or stop the conflicting service.

### Cannot connect to backend
- Verify `VITE_API_URL` matches your backend URL
- Ensure backend is running and accessible
- Check network connectivity

### Hot reload not working
- Ensure volume mounts are correct
- Check file permissions
- Restart the container: `docker-compose restart`

