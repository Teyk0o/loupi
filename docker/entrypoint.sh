#!/bin/sh
# =============================================================================
# Loupi — All-in-one entrypoint
# Starts the Go API, Next.js frontend, and Caddy reverse proxy.
# Monitors child processes and exits if any child dies.
# =============================================================================

set -e

# Forward signals to child processes
cleanup() {
    echo "[loupi] Shutting down..."
    kill "$API_PID" "$FRONTEND_PID" "$CADDY_PID" 2>/dev/null || true
    wait "$API_PID" "$FRONTEND_PID" "$CADDY_PID" 2>/dev/null || true
    exit 0
}
trap cleanup SIGTERM SIGINT

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

# Start Caddy in the background
echo "[loupi] Starting Caddy reverse proxy..."
caddy run --config /etc/caddy/Caddyfile --adapter caddyfile &
CADDY_PID=$!

# Monitor all child processes — exit if any dies
while true; do
    for PID in $API_PID $FRONTEND_PID $CADDY_PID; do
        if ! kill -0 "$PID" 2>/dev/null; then
            echo "[loupi] Process $PID died, shutting down..."
            cleanup
        fi
    done
    sleep 5
done
