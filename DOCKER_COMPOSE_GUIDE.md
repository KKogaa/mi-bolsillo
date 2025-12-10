# Mi Bolsillo - Docker Compose Guide

This guide explains how to run the Mi Bolsillo application stack using Docker Compose.

## Project Structure

The Mi Bolsillo project consists of three main components:

1. **Echo Server** (`mi-bolsillo-api/Dockerfile.echo`) - Main REST API backend
2. **Telegram Bot** (`mi-bolsillo-api/Dockerfile.telegram`) - Telegram bot service
3. **Frontend** (`mi-bolsillo-front/Dockerfile`) - React frontend application

## Docker Compose Files

### Root Level (`docker-compose.yml`)
The root `docker-compose.yml` orchestrates all services together. Use this to run the entire application stack.

**Location:** `/mnt/Laplace/projects/mi-bolsillo/docker-compose.yml`

### Individual Service Files
Each project also has its own `docker-compose.yml` for running services independently:
- `mi-bolsillo-api/docker-compose.yml` - Runs backend services only
- `mi-bolsillo-front/docker-compose.yml` - Runs frontend only

## Quick Start

### 1. Environment Setup

```bash
# Navigate to the root directory
cd /mnt/Laplace/projects/mi-bolsillo

# Copy the example environment file
cp .env.example .env

# Edit .env and fill in your actual values
nano .env  # or use your preferred editor
```

### 2. Start All Services

```bash
# Start all services (echo-server, telegram-bot, frontend)
docker-compose up -d

# View logs
docker-compose logs -f

# View logs for specific service
docker-compose logs -f echo-server
```

### 3. Access the Application

- **Frontend:** http://localhost:3000
- **API:** http://localhost:8080
- **Telegram Bot:** No web interface (uses Telegram long polling)

## Common Commands

### Start Services

```bash
# Start all services
docker-compose up -d

# Start specific services
docker-compose up -d echo-server frontend

# Start without detached mode (see logs directly)
docker-compose up
```

### Stop Services

```bash
# Stop all services
docker-compose down

# Stop and remove volumes
docker-compose down -v
```

### Rebuild Services

```bash
# Rebuild all services
docker-compose up -d --build

# Rebuild specific service
docker-compose up -d --build frontend

# Force rebuild without cache
docker-compose build --no-cache
```

### View Logs

```bash
# View all logs
docker-compose logs

# Follow logs (live)
docker-compose logs -f

# View logs for specific service
docker-compose logs -f echo-server

# View last 100 lines
docker-compose logs --tail=100
```

### Service Status

```bash
# Check running services
docker-compose ps

# Check resource usage
docker stats
```

### Execute Commands in Containers

```bash
# Open shell in echo-server
docker-compose exec echo-server sh

# Open shell in frontend
docker-compose exec frontend sh
```

## Service Dependencies

The services have the following dependency chain:

```
echo-server (starts first)
    ↓
telegram-bot (waits for echo-server health check)
    ↓
frontend (waits for echo-server health check)
```

## Health Checks

### Echo Server
- **Endpoint:** `http://localhost:8080/health`
- **Interval:** 30s
- **Timeout:** 10s
- **Retries:** 3

### Frontend
- **Endpoint:** `http://localhost:80/`
- **Interval:** 30s
- **Timeout:** 10s
- **Retries:** 3

## Environment Variables

### Required for All Backend Services
- `DATABASE_URL` - Turso/LibSQL database URL
- `DATABASE_TOKEN` - Database authentication token
- `CLERK_JWKS_URL` - Clerk JWKS URL for authentication
- `GROK_API_KEY` - Grok API key for AI features

### Required for Frontend
- `VITE_CLERK_PUBLISHABLE_KEY` - Clerk publishable key
- `VITE_API_URL` - Backend API URL (default: http://localhost:8080)
- `VITE_TELEGRAM_BOT_USERNAME` - Telegram bot username

### Required for Telegram Bot
- `TELEGRAM_BOT_TOKEN` - Bot token from @BotFather

### Optional
- `EMAIL_PROVIDER_URL` - Email service URL
- `EMAIL_PROVIDER_TOKEN` - Email service token
- `OTP_EXPIRATION_MINUTES` - OTP expiration time (default: 5)

See `.env.example` for complete documentation.

## Network Configuration

All services run on a shared Docker network called `mi-bolsillo-network`. This allows services to communicate with each other using their service names.

### Service Communication
- Frontend can reach API at: `http://echo-server:8080`
- Telegram bot can reach API at: `http://echo-server:8080`

## Port Mapping

| Service | Internal Port | External Port | URL |
|---------|--------------|---------------|-----|
| Echo Server | 8080 | 8080 | http://localhost:8080 |
| Frontend | 80 | 3000 | http://localhost:3000 |
| Telegram Bot | - | - | No web interface |

## Running Individual Projects

### Backend Only

```bash
cd mi-bolsillo-api
docker-compose up -d
```

This starts:
- Echo server on port 8080
- Telegram bot

### Frontend Only

```bash
cd mi-bolsillo-front
docker-compose up -d
```

This starts:
- Frontend on port 3000

**Note:** Make sure `VITE_API_URL` points to a running backend instance.

## Troubleshooting

### Services Won't Start

```bash
# Check service status
docker-compose ps

# Check logs for errors
docker-compose logs

# Restart services
docker-compose restart
```

### Port Already in Use

```bash
# Check what's using the port
sudo lsof -i :8080
sudo lsof -i :3000

# Stop the conflicting service or change port in docker-compose.yml
```

### Database Connection Issues

1. Verify `DATABASE_URL` and `DATABASE_TOKEN` in `.env`
2. Check if Turso database is accessible
3. Review echo-server logs: `docker-compose logs echo-server`

### Build Failures

```bash
# Clear build cache and rebuild
docker-compose build --no-cache

# Remove all containers and rebuild
docker-compose down
docker-compose up -d --build
```

### Frontend Environment Variables Not Working

The frontend uses runtime environment variable injection. If variables aren't being picked up:

1. Check `.env` file has correct values
2. Rebuild frontend: `docker-compose up -d --build frontend`
3. Check `runtime-config.js` is being loaded in browser console

## Production Deployment

For production deployment, consider:

1. **Use production Dockerfiles** - Ensure multi-stage builds are optimized
2. **Set secure environment variables** - Never commit `.env` to version control
3. **Use secrets management** - Consider Docker secrets or external secret managers
4. **Configure reverse proxy** - Use nginx or Traefik for SSL/TLS termination
5. **Monitor resources** - Set up resource limits in docker-compose.yml
6. **Enable logging** - Configure proper log aggregation
7. **Backup database** - Regular backups of Turso database

### Adding Resource Limits

Add to each service in `docker-compose.yml`:

```yaml
deploy:
  resources:
    limits:
      cpus: '1'
      memory: 512M
    reservations:
      cpus: '0.5'
      memory: 256M
```

## Development Tips

### Hot Reload for Frontend

For development with hot reload:

```bash
# In mi-bolsillo-front directory
npm run dev
```

This runs Vite dev server with hot module replacement instead of Docker.

### Live Backend Changes

For Go backend development, you can use `air` for hot reload:

```bash
# Install air
go install github.com/cosmtrek/air@latest

# Run with air in mi-bolsillo-api
air
```

### Debugging

To debug a service:

```bash
# Attach to service logs
docker-compose logs -f echo-server

# Execute shell in running container
docker-compose exec echo-server sh

# Check environment variables
docker-compose exec echo-server env
```

## Maintenance

### Cleanup

```bash
# Remove stopped containers
docker-compose down

# Remove containers, networks, and volumes
docker-compose down -v

# Clean up Docker system
docker system prune -a
```

### Updates

```bash
# Pull latest code
git pull

# Rebuild services
docker-compose up -d --build
```

## Support

For issues and questions:
- Check logs: `docker-compose logs -f`
- Review environment variables in `.env`
- Verify database connectivity
- Check network configuration

## Additional Resources

- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [Turso Database Docs](https://docs.turso.tech/)
- [Clerk Authentication](https://clerk.com/docs)
- [Telegram Bot API](https://core.telegram.org/bots/api)
