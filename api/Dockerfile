# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /build

# Install dependencies first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /build/server ./cmd/server

# Runtime stage
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy binary and migrations
COPY --from=builder /build/server .
COPY --from=builder /build/migrations ./migrations

EXPOSE 8080

CMD ["./server"]
