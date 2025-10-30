# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm ci

# Copy source code
COPY . .

# Accept build-time environment variables (optional, for build-time injection)
ARG VITE_CLERK_PUBLISHABLE_KEY=""
ARG VITE_API_URL=""
ARG VITE_TELEGRAM_BOT_USERNAME=""

# Set environment variables for build (if provided)
ENV VITE_CLERK_PUBLISHABLE_KEY=$VITE_CLERK_PUBLISHABLE_KEY
ENV VITE_API_URL=$VITE_API_URL
ENV VITE_TELEGRAM_BOT_USERNAME=$VITE_TELEGRAM_BOT_USERNAME

# Build the application
RUN npm run build

# Production stage
FROM nginx:alpine

# Copy built assets from builder
COPY --from=builder /app/dist /usr/share/nginx/html

# Create nginx configuration that supports PORT environment variable
RUN echo 'server { \
    listen ${PORT}; \
    server_name _; \
    root /usr/share/nginx/html; \
    index index.html; \
    location / { \
        try_files $uri $uri/ /index.html; \
    } \
    # Security headers \
    add_header X-Frame-Options "SAMEORIGIN" always; \
    add_header X-Content-Type-Options "nosniff" always; \
    add_header X-XSS-Protection "1; mode=block" always; \
}' > /etc/nginx/conf.d/default.conf.template

# Copy entrypoint script that injects runtime env vars
COPY <<'EOF' /docker-entrypoint.sh
#!/bin/sh
set -e

# Use PORT environment variable or default to 80
PORT=${PORT:-80}

# Replace PORT placeholder in nginx config
envsubst '${PORT}' < /etc/nginx/conf.d/default.conf.template > /etc/nginx/conf.d/default.conf

# Inject runtime environment variables into the built JavaScript
if [ -n "$VITE_CLERK_PUBLISHABLE_KEY" ] || [ -n "$VITE_API_URL" ] || [ -n "$VITE_TELEGRAM_BOT_USERNAME" ]; then
    echo "Injecting runtime environment variables..."

    # Create runtime config that will be loaded before the app
    cat > /usr/share/nginx/html/runtime-config.js <<RUNTIME_EOF
window.RUNTIME_CONFIG = {
    VITE_CLERK_PUBLISHABLE_KEY: "${VITE_CLERK_PUBLISHABLE_KEY}",
    VITE_API_URL: "${VITE_API_URL}",
    VITE_TELEGRAM_BOT_USERNAME: "${VITE_TELEGRAM_BOT_USERNAME}"
};
RUNTIME_EOF

    # Inject script tag into index.html before other scripts
    sed -i 's|<head>|<head>\n    <script src="/runtime-config.js"></script>|' /usr/share/nginx/html/index.html
fi

# Start nginx
exec nginx -g 'daemon off;'
EOF

RUN chmod +x /docker-entrypoint.sh

# Expose PORT (will be set by Render)
EXPOSE ${PORT}

CMD ["/docker-entrypoint.sh"]
