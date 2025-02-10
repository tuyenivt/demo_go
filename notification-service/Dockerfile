# Stage 1: Build
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first for dependency caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the service source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o notification-service ./cmd/main.go

# Stage 2: Production-ready image
FROM alpine:3.21

WORKDIR /app
COPY --from=builder /app/notification-service /app/notification-service

# Non-root user (optional security best practice)
# RUN addgroup -S appgroup && adduser -S appuser -G appgroup
# USER appuser

# Expose port if needed
EXPOSE 8080

ENTRYPOINT ["/app/notification-service"]