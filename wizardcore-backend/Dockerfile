# Build stage
FROM golang:1.24-alpine AS builder

# Cache busting argument
ARG CACHE_BUST=1

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application (with cache busting)
RUN echo "Cache bust: ${CACHE_BUST}" && \
    CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .

# Copy migrations (preserve path relative to /app)
COPY --from=builder /app/internal/database/migrations ./internal/database/migrations

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]