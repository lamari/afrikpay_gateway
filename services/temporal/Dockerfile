# ----------- Build Stage -----------
FROM golang:1.21-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build API binary
RUN CGO_ENABLED=0 go build -o /out/api ./cmd/api/main.go
# Build Worker binary
RUN CGO_ENABLED=0 go build -o /out/worker ./cmd/worker/main.go

# ----------- Final Stage -----------
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /out/api /app/api
COPY --from=builder /out/worker /app/worker
# Copy config if needed at runtime (optional, can be mounted)
COPY config /config

# Set default command (can be overridden in docker-compose)
CMD ["/app/api"]
