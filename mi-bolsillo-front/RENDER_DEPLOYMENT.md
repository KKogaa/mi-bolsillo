# Render Deployment Guide for Mi Bolsillo Frontend

## Environment Variables Configuration

**IMPORTANT**: For this app to work properly on Render, you need to configure the following environment variables in your Render service settings.

### Required Environment Variables

Go to your Render service → **Environment** tab and add these variables:

1. **VITE_CLERK_PUBLISHABLE_KEY**
   - Your Clerk publishable key
   - Example: `pk_test_...` or `pk_live_...`
   - Get this from: https://dashboard.clerk.com

2. **VITE_API_URL**
   - Your backend API URL
   - Example: `https://mi-bolsillo-backend.onrender.com`
   - This should point to your deployed backend service

3. **VITE_TELEGRAM_BOT_USERNAME**
   - Your Telegram bot username (without @)
   - Example: `mi_bolsillo_bot`

4. **PORT** (Optional - usually auto-set by Render)
   - The port nginx will listen on
   - Render typically sets this automatically to 10000 or uses 443 for HTTPS

### How to Add Environment Variables on Render

1. Go to your Render Dashboard
2. Select your frontend service
3. Click on **Environment** in the left sidebar
4. Click **Add Environment Variable**
5. Add each variable with its key and value
6. Click **Save Changes**

### Important Notes

- **These are runtime environment variables**: They are injected when the container starts, not during build time
- **After adding/changing variables**: You must manually trigger a new deploy or the changes won't take effect
- **Check the logs**: After deployment, check the container logs. You should see:
  ```
  === Starting container initialization ===
  PORT: 443
  VITE_CLERK_PUBLISHABLE_KEY is set: YES
  VITE_API_URL is set: YES
  VITE_TELEGRAM_BOT_USERNAME is set: YES
  === Injecting runtime environment variables ===
  Runtime config created successfully
  ...
  ```

### Troubleshooting

#### Error: "Missing Clerk Publishable Key"

This means `VITE_CLERK_PUBLISHABLE_KEY` is not being read. Check:

1. Is the variable set in Render's Environment tab?
2. Did you redeploy after adding the variable?
3. Check the container logs - do you see "VITE_CLERK_PUBLISHABLE_KEY is set: YES"?
4. Open browser console and check if `window.RUNTIME_CONFIG` exists and has the key

#### Only Favicon Shows

This could mean:

1. Environment variables are missing (causing React to error before rendering)
2. Port mismatch - check that nginx is listening on the correct PORT
3. Check browser console for JavaScript errors

#### How to Verify Variables are Loaded

1. After deployment, open your app in a browser
2. Open browser DevTools (F12)
3. Go to Console tab
4. Type: `window.RUNTIME_CONFIG`
5. You should see your environment variables (publishable key will be visible - this is normal for Clerk public keys)

### Docker Build Arguments (Alternative Approach - Not Recommended)

If you want to inject variables at build time instead, you can use Docker build arguments in Render:

1. Go to your service settings
2. Under **Docker Command**, you can pass build args like:
   ```
   --build-arg VITE_CLERK_PUBLISHABLE_KEY=$VITE_CLERK_PUBLISHABLE_KEY
   ```

However, the current setup (runtime injection) is better because:
- You can change variables without rebuilding
- Faster deployments
- More flexible for different environments

## Deployment Checklist

- [ ] All environment variables are set in Render
- [ ] Repository is connected to Render
- [ ] Dockerfile is in the repository root
- [ ] Manual deploy or auto-deploy is triggered
- [ ] Check logs for successful variable injection
- [ ] Test the deployed app
- [ ] Verify browser console shows `window.RUNTIME_CONFIG` with your values
