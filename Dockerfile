# =============================================================================
# Loupi — All-in-one production image
# Contains: Go API + Next.js frontend + Caddy reverse proxy
# Requires: External PostgreSQL and Redis
# =============================================================================

# ---------------------------------------------------------------------------
# Stage 1: Build Go API
# ---------------------------------------------------------------------------
FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS api-builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /build

COPY api/go.mod api/go.sum ./
RUN go mod download

COPY api/ .
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-s -w" -o /build/server ./cmd/server

# ---------------------------------------------------------------------------
# Stage 2: Install frontend dependencies
# ---------------------------------------------------------------------------
FROM --platform=$BUILDPLATFORM node:22-alpine AS frontend-deps

RUN corepack enable && corepack prepare pnpm@9.15.4 --activate

WORKDIR /app

COPY frontend/package.json frontend/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

# ---------------------------------------------------------------------------
# Stage 3: Build Next.js frontend
# ---------------------------------------------------------------------------
FROM --platform=$BUILDPLATFORM node:22-alpine AS frontend-builder

RUN corepack enable && corepack prepare pnpm@9.15.4 --activate

WORKDIR /app

COPY --from=frontend-deps /app/node_modules ./node_modules
COPY frontend/ .

ENV NEXT_TELEMETRY_DISABLED=1
ENV NODE_ENV=production

RUN pnpm build

# ---------------------------------------------------------------------------
# Stage 4: Runtime image
# ---------------------------------------------------------------------------
FROM alpine:3.21

RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    nodejs \
    caddy \
    tini \
    && addgroup -S loupi \
    && adduser -S -G loupi loupi

WORKDIR /app

# Copy Go API binary and migrations
COPY --from=api-builder /build/server ./api-server
COPY --from=api-builder /build/migrations ./migrations

# Copy Next.js standalone build
COPY --from=frontend-builder /app/public ./frontend/public
COPY --from=frontend-builder /app/.next/standalone ./frontend/
COPY --from=frontend-builder /app/.next/static ./frontend/.next/static

# Copy Caddy configuration
COPY docker/Caddyfile /etc/caddy/Caddyfile

# Copy entrypoint
COPY docker/entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

# Create data directories
RUN mkdir -p /app/photos /app/data /config/caddy /data/caddy \
    && chown -R loupi:loupi /app /config/caddy /data/caddy

# Labels (OCI standard)
LABEL org.opencontainers.image.title="Loupi"
LABEL org.opencontainers.image.description="Free, open-source food & digestive symptom tracker"
LABEL org.opencontainers.image.url="https://loupi.app"
LABEL org.opencontainers.image.source="https://github.com/Teyk0o/loupi"
LABEL org.opencontainers.image.licenses="Apache-2.0"

# Next.js standalone must bind to all interfaces inside the container
ENV HOSTNAME=0.0.0.0
ENV PORT=3000

# Caddy serves on 80 (HTTP) and 443 (HTTPS)
EXPOSE 80 443

USER loupi

# Use tini as PID 1 for proper signal handling
ENTRYPOINT ["/sbin/tini", "--"]
CMD ["/app/entrypoint.sh"]
