FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum for both project and module
COPY go.work go.work.sum ./
COPY services/auth/go.mod services/auth/go.sum ./services/auth/

# Copy shared directory if it exists and is used
COPY shared/ ./shared/

# Download dependencies
RUN cd services/auth && go mod download

# Copy the source code
COPY services/auth/ ./services/auth/

# Build the application
RUN cd services/auth && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o auth ./cmd/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Install necessary tools for healthcheck
RUN apk --no-cache add curl ca-certificates

# Copy config and create logs directory
COPY --from=builder /app/services/auth/api /app/api
RUN mkdir -p /app/logs /app/config/keys

# Copy the binary from the builder stage
COPY --from=builder /app/services/auth/auth /app/auth

# Set executable permission
RUN chmod +x /app/auth

# Expose port
EXPOSE 8001

# Run the application
CMD ["/app/auth"]
