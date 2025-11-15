# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o blitz .

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1000 blitz && \
    adduser -D -u 1000 -G blitz blitz

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/blitz .

# Copy wordlists
COPY usernames.txt passwords.txt ./

# Change ownership
RUN chown -R blitz:blitz /app

# Switch to non-root user
USER blitz

# Set entrypoint
ENTRYPOINT ["./blitz"]

# Default command shows help
CMD ["--help"]
