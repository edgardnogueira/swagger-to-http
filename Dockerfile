# Build stage
FROM golang:1.21-alpine AS builder

# Set up build environment
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git make

# Copy Go module files first for better caching
COPY go.mod go.sum* ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN make build

# Final stage
FROM alpine:3.18

# Add necessary runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Set up non-root user for security
RUN adduser -D appuser
USER appuser

WORKDIR /home/appuser

# Copy binary from builder
COPY --from=builder /app/bin/swagger-to-http /usr/local/bin/

# Command to run
ENTRYPOINT ["swagger-to-http"]
CMD ["--help"]
