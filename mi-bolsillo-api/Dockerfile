# Build stage
FROM golang:1.25.1-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code (excluding .env files)
COPY . .

# Remove any .env files that might have been copied
RUN rm -f .env .env.* || true

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .

# Copy config files (messages.json, etc.) if needed
COPY --from=builder /app/config ./config

# Expose port 8080
EXPOSE 8080

# Environment variables will be injected at runtime by Render
# No need to copy .env files - all config comes from environment

# Run the binary
CMD ["./main"]
