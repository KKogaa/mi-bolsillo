# Render Deployment Guide

This guide explains how to deploy the Mi Bolsillo API services to Render.

## Prerequisites

1. Docker images are automatically built and pushed to Docker Hub via GitHub Actions
2. You need a Render account
3. Environment variables values ready (see `.env.example` in `mi-bolsillo-api/`)

## Architecture

The project consists of three services:
- **Echo Server** (`mi-bolsillo-echo`): Main HTTP API server
- **Telegram Bot** (`mi-bolsillo-telegram`): Telegram bot service
- **Frontend** (`mi-bolsillo-front`): React/Vite frontend application

## Deployment Steps

### 1. Create Web Services on Render

For each service (Echo Server, Telegram Bot, and Frontend):

1. Go to Render Dashboard → New → Web Service
2. Choose "Deploy an existing image from a registry"
3. Enter the Docker image URL:
   - Echo Server: `your-dockerhub-username/mi-bolsillo-echo:latest`
   - Telegram Bot: `your-dockerhub-username/mi-bolsillo-telegram:latest`
   - Frontend: `your-dockerhub-username/mi-bolsillo-front:latest`

### 2. Configure Service Settings

#### Echo Server Configuration
- **Name**: `mi-bolsillo-echo` (or your preferred name)
- **Region**: Choose closest to your users
- **Instance Type**: Choose based on your needs (Free tier available)
- **Port**: Will be set via `PORT` env variable (Render sets this automatically)

#### Telegram Bot Configuration
- **Name**: `mi-bolsillo-telegram`
- **Region**: Same as Echo Server for consistency
- **Instance Type**: Choose based on your needs
- **Note**: No port configuration needed (uses long polling)

#### Frontend Configuration
- **Name**: `mi-bolsillo-front`
- **Region**: Same as backend services
- **Instance Type**: Choose based on your needs
- **Port**: Depends on your frontend server configuration

### 3. Set Environment Variables

In each service's Environment tab, add the required variables:

#### Echo Server Environment Variables
```
DATABASE_URL=libsql://your-database.turso.io
DATABASE_TOKEN=your-database-token
CLERK_JWKS_URL=https://your-app.clerk.accounts.dev/.well-known/jwks.json
GROK_API_KEY=your-grok-api-key
OTP_EXPIRATION_MINUTES=5
EMAIL_PROVIDER_URL=https://your-email-provider.com/api (optional)
EMAIL_PROVIDER_TOKEN=your-email-token (optional)
```

**Note**: Render automatically sets the `PORT` variable, so you don't need to add it.

#### Telegram Bot Environment Variables
```
DATABASE_URL=libsql://your-database.turso.io
DATABASE_TOKEN=your-database-token
CLERK_JWKS_URL=https://your-app.clerk.accounts.dev/.well-known/jwks.json
GROK_API_KEY=your-grok-api-key
TELEGRAM_BOT_TOKEN=your-telegram-bot-token
```

#### Frontend Environment Variables
Add any `VITE_*` or `REACT_APP_*` variables your frontend needs:
```
VITE_API_URL=https://your-echo-server.onrender.com
VITE_CLERK_PUBLISHABLE_KEY=your-clerk-key
```

### 4. Deploy

1. Click "Create Web Service"
2. Render will pull the Docker image and start the service
3. Monitor the logs for any startup issues
4. Once deployed, you'll get a public URL like `https://your-service.onrender.com`

## Auto-Deploy on Push

To enable automatic deployments when you push to GitHub:

1. In each Render service, go to Settings → Build & Deploy
2. Enable "Auto-Deploy"
3. Choose "Deploy on image update"
4. Render will automatically pull the latest `:latest` tag when GitHub Actions pushes a new image

Alternatively, you can set up webhooks:
1. Get the Deploy Hook URL from Render service settings
2. Add it as a repository secret in GitHub
3. Update the GitHub Actions workflow to trigger the deploy hook after building images

## Image Update Workflow

The current setup:
1. You push code to `master` branch
2. GitHub Actions builds Docker images
3. Images are pushed to Docker Hub with `:latest` tag
4. Render can auto-deploy or you manually trigger deployment

## Monitoring and Logs

- **Logs**: View real-time logs in Render Dashboard → Your Service → Logs
- **Metrics**: Monitor CPU, memory, and request metrics in the Metrics tab
- **Health Checks**: Configure health check endpoints in Settings → Health & Alerts

## Troubleshooting

### Service won't start
- Check logs for error messages
- Verify all required environment variables are set
- Ensure Docker image was built successfully in GitHub Actions

### Database connection issues
- Verify `DATABASE_URL` and `DATABASE_TOKEN` are correct
- Check if Turso database is accessible
- Ensure database region is accessible from Render

### Port binding errors
- For Echo Server: Ensure your Go app reads `PORT` from environment
- Render sets `PORT` automatically - don't hardcode it
- Check `config/config.go:35` reads from `os.Getenv("PORT")`

### Environment variables not loading
- Verify variables are set in Render Dashboard → Environment
- Restart the service after adding/updating variables
- Check logs to see what values are being read

## Cost Optimization

- Free tier available for basic usage
- Consider using shared databases across services
- Monitor usage in Render Dashboard
- Scale down during low-traffic periods

## Security Best Practices

1. Never commit `.env` files to version control
2. Use Render's environment variable encryption
3. Rotate secrets regularly
4. Use least-privilege database credentials
5. Enable HTTPS (Render provides this automatically)
6. Set up proper CORS configuration in your API

## Next Steps

After deployment:
1. Test all API endpoints
2. Verify Telegram bot is responding
3. Check database connections
4. Set up monitoring and alerts
5. Configure custom domain (optional)
6. Set up SSL certificates if using custom domain

## Support

- Render Documentation: https://render.com/docs
- GitHub Actions logs for build issues
- Check application logs in Render Dashboard
