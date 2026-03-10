#!/bin/sh
# =============================================================================
# Loupi — All-in-one entrypoint
# Starts the Go API, Next.js frontend, and Caddy reverse proxy.
# =============================================================================

set -e

echo "============================================"
echo "  Loupi — Starting all services"
echo "============================================"

# Start the Go API in the background
echo "[loupi] Starting API server on :8080..."
cd /app
./api-server &
API_PID=$!

# Start the Next.js frontend in the background
echo "[loupi] Starting frontend on :3000..."
cd /app/frontend
NODE_ENV=production node server.js &
FRONTEND_PID=$!

# Wait for services to be ready
sleep 2

# Start Caddy in the foreground
echo "[loupi] Starting Caddy reverse proxy..."
exec caddy run --config /etc/caddy/Caddyfile --adapter caddyfile
