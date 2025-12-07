# Deployment Guide

This guide explains how to deploy Invoice Scan to a VPS using GitHub Actions.

## Prerequisites

- VPS with Docker installed (Ubuntu 20.04+ recommended)
- Existing MySQL instance
- GitHub repository with Actions enabled

## Architecture

```
                    ┌─────────────────────────────────────────────────────┐
                    │                        VPS                          │
   HTTPS :443       │  ┌─────────────────────────────────────────────┐   │
  ─────────────────▶│  │              Frontend (nginx)               │   │
                    │  │  - Serves static files (React app)          │   │
   HTTP :80         │  │  - SSL termination                          │   │
  ──── redirect ───▶│  │  - Proxies /api/* to backend                │   │
                    │  └──────────────────┬──────────────────────────┘   │
                    │                     │ internal network             │
                    │                     ▼ http://backend:3001          │
                    │  ┌─────────────────────────────────────────────┐   │
                    │  │              Backend (Go)                    │   │
                    │  │  - REST API                                  │   │
                    │  │  - Not exposed externally                    │   │
                    │  └──────────────────┬──────────────────────────┘   │
                    │                     │                              │
                    └─────────────────────┼──────────────────────────────┘
                                          │
                                 ┌────────▼────────┐
                                 │  MySQL (yours)  │
                                 └─────────────────┘
```

**Key points:**
- Frontend runs on HTTPS (required for camera access)
- Backend is only accessible via internal Docker network
- Nginx proxies `/api/*` requests to backend

## Setup Steps

### 1. Configure GitHub Secrets

Go to your GitHub repository → Settings → Secrets and variables → Actions → New repository secret

| Secret Name | Description | Example |
|-------------|-------------|---------|
| `VPS_HOST` | Your VPS IP or hostname | `123.45.67.89` |
| `VPS_USER` | SSH username | `root` or `deploy` |
| `VPS_SSH_KEY` | Private SSH key for VPS access | `-----BEGIN OPENSSH PRIVATE KEY-----...` |
| `DEPLOY_PATH` | Deployment directory on VPS | `/opt/invoice-scan` |
| `SITE_URL` | Your site URL (HTTPS) | `https://invoice.yourdomain.com` or `https://YOUR_VPS_IP` |
| `GEMINI_API_KEY` | Google Gemini API key | `AIza...` |
| `DATABASE_HOST` | MySQL host address | `your-mysql-host` |
| `DATABASE_PORT` | MySQL port | `3306` |
| `DATABASE_USER` | MySQL username | `invoice_user` |
| `DATABASE_PASSWORD` | MySQL password | `your-password` |
| `DATABASE_NAME` | Database name | `invoice_scan` |

### 2. Prepare VPS (One-time setup)

SSH into your VPS and install Docker:

```bash
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker $USER
exit
```

That's it! The GitHub Actions workflow will automatically:
- Create the deployment directory
- Generate self-signed SSL certificate (if not exists)
- Copy `docker-compose.prod.yml` to VPS
- Create `.env` file from secrets
- Pull and start containers

### 3. Run Database Migration (One-time)

Ensure your MySQL has the required table:

```sql
CREATE TABLE IF NOT EXISTS invoices (
    id VARCHAR(26) PRIMARY KEY,
    file_path VARCHAR(500) NOT NULL,
    original_filename VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    extracted_data JSON,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### 4. Deploy

Just push to main branch:

```bash
git add .
git commit -m "feat: add deployment configuration"
git push origin main
```

GitHub Actions will automatically:
1. Build Docker images for backend and frontend
2. Push images to GitHub Container Registry (ghcr.io)
3. Copy `docker-compose.prod.yml` to VPS
4. Generate SSL cert (if needed)
5. Create/update `.env` file from GitHub secrets
6. Pull latest images and restart containers

**Every push to `main` triggers automatic deployment!**

## SSL Certificates

### Default: Self-Signed (Demo)

The workflow automatically generates a self-signed certificate on first deployment. This works for demos but browsers will show a security warning.

### Production: Let's Encrypt

For production, replace self-signed certs with Let's Encrypt:

```bash
# On VPS
sudo apt install certbot -y

# Stop frontend temporarily
cd /opt/invoice-scan
docker compose -f docker-compose.prod.yml stop frontend

# Get certificate (replace with your domain)
sudo certbot certonly --standalone -d invoice.yourdomain.com

# Copy certs to ssl directory
sudo cp /etc/letsencrypt/live/invoice.yourdomain.com/fullchain.pem ssl/cert.pem
sudo cp /etc/letsencrypt/live/invoice.yourdomain.com/privkey.pem ssl/key.pem
sudo chown $USER:$USER ssl/*.pem

# Start frontend
docker compose -f docker-compose.prod.yml up -d frontend
```

### Cloudflare (Easiest)

1. Add your domain to Cloudflare
2. Point DNS to your VPS IP (proxied - orange cloud)
3. Set SSL/TLS mode to "Full"
4. The self-signed cert on your server will work with Cloudflare's edge SSL

## Access Points

After deployment:
- **App**: `https://YOUR_VPS_IP` or `https://your-domain.com`
- **Health Check**: `https://YOUR_VPS_IP/api/v1/health`

## Troubleshooting

### Check container status
```bash
cd /opt/invoice-scan
docker compose -f docker-compose.prod.yml ps
```

### View logs
```bash
docker compose -f docker-compose.prod.yml logs -f
docker compose -f docker-compose.prod.yml logs -f backend
docker compose -f docker-compose.prod.yml logs -f frontend
```

### Restart services
```bash
docker compose -f docker-compose.prod.yml restart
```

### SSL certificate issues
```bash
# Check if certs exist
ls -la /opt/invoice-scan/ssl/

# Regenerate self-signed cert
cd /opt/invoice-scan
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout ssl/key.pem \
  -out ssl/cert.pem \
  -subj "/CN=localhost"

docker compose -f docker-compose.prod.yml restart frontend
```

### Cannot connect to database
- Verify DATABASE_HOST is accessible from VPS
- Check firewall allows connection to MySQL port
- Verify credentials are correct

### GitHub Actions failed
- Check Actions tab in GitHub for error details
- Verify all secrets are set correctly
- Ensure VPS SSH key has correct permissions

## Environment Variables Reference

| Variable | Required | Description |
|----------|----------|-------------|
| `GITHUB_REPOSITORY` | Auto | GitHub repo in format `user/repo` |
| `GEMINI_API_KEY` | Yes | Google Gemini API key |
| `SITE_URL` | Yes | Full site URL with https |
| `DATABASE_HOST` | Yes | MySQL host address |
| `DATABASE_PORT` | No | MySQL port (default: 3306) |
| `DATABASE_USER` | Yes | MySQL username |
| `DATABASE_PASSWORD` | Yes | MySQL password |
| `DATABASE_NAME` | Yes | Database name |
| `SSL_CERT_PATH` | No | Path to SSL certs (default: ./ssl) |
