# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy dependency manifests
COPY go.mod go.sum ./

# Cache dependencies
RUN go mod download

# Copy application source
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/health-checker ./cmd/server

# Runtime stage
FROM alpine:latest

# Security: install certificates and upgrade packages
RUN apk --no-cache add ca-certificates && \
    apk --no-cache upgrade

# Create non-root system user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/health-checker /app/health-checker

# Change ownership of working directory to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root execution
USER appuser

# Expose API port
EXPOSE 8000

# Execute server
ENTRYPOINT ["/app/health-checker"]
