# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/auth-service

# Final stage
FROM alpine:latest
WORKDIR /app

# Install CA certificates
RUN apk --no-cache add ca-certificates

# Copy the binary from builder
COPY --from=builder /app/main /app/main

# Copy any required config files (like .env if you have one)
# COPY --from=builder /app/.env /app/.env

EXPOSE 8080
CMD ["/app/main"]