# Loupi — Project Instructions

## Project Overview

Loupi is a free, open-source PWA for food & digestive symptom tracking, hosted at loupi.app.

## Architecture

- **Monorepo** with `frontend/`, `api/`, `docker/`, `docs/`
- **Frontend:** Next.js (React) + Tailwind CSS + pnpm
- **Backend:** Go (REST API)
- **Database:** PostgreSQL 16
- **Cache/Sessions:** Redis 7
- **Reverse proxy:** Caddy (auto HTTPS)
- **Deployment:** Docker Compose on Gandi VPS (loupi.app)

## Code Standards

### Language
- All code, comments, variable names, commits, and PRs must be in **English**.
- The user interface is in **French** (France-only market for v1).

### Frontend (Next.js 16)
- Package manager: **pnpm** (never use npm or yarn)
- TypeScript strict mode
- **Tailwind CSS v4** (CSS-first config via `@theme inline` in `globals.css`, NO `tailwind.config.ts`)
- Path alias: `@/*` → `./src/*`
- File-based routing (App Router), source in `src/`
- Mobile-first, responsive design
- PWA-ready (manifest + service worker)

### Design System
- **Style:** Soft & warm (rounded corners, soft shadows, airy spacing)
- **Fonts:** Poppins (headings, labels, buttons) + Inter (body text, fallback)
- **Primary:** `#48C7B0` (mint green) — use `text-primary`, `bg-primary`
- **Secondary:** `#FFE6B3` (soft peach) — use `text-secondary`, `bg-secondary`
- **Dark mode:** Auto via `prefers-color-scheme` — tokens defined as CSS vars in `globals.css`
- **Radii:** `rounded-sm` (8px), `rounded-md` (12px), `rounded-lg` (16px), `rounded-xl` (24px)
- **Icons:** Lucide React (thin stroke, consistent with soft style)
- See `globals.css` for all design tokens

### Backend (Go)
- Project layout follows [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- `cmd/server/` — main entry point
- `internal/` — private application code (config, database, handlers, middleware, models, routes, services, utils)
- Use `net/http` or Gin for HTTP routing
- All endpoints prefixed with `/v1`
- JSON request/response format
- UUID for all primary keys
- Proper error handling with structured error responses

### Database
- PostgreSQL 16 with migrations (versioned SQL files)
- JSONB for flexible data (symptoms, sport metrics)
- AES-256-GCM encryption for sensitive health data at rest
- Audit logging for all data access/modifications

### Security
- HTTPS only (TLS 1.3 via Caddy)
- JWT + refresh tokens for authentication
- bcrypt (cost ≥ 12) for password hashing
- CSRF protection
- Rate limiting on auth endpoints
- Input validation and sanitization on all endpoints
- Security headers (CSP, HSTS, X-Frame-Options, etc.)
- Photos served only through authenticated API endpoints

## Git Conventions

### Commits
- Follow **Conventional Commits** strictly
- Format: `type(scope): description`
- Types: feat, fix, docs, style, refactor, test, chore, ci, perf
- Scopes: frontend, api, docker, docs, deps
- Never reference Claude or AI tools in commits

### Branches
- Format: `type/short-description`
- Examples: `feat/meal-tracking`, `fix/auth-refresh`, `chore/docker-setup`

## Testing
- Frontend: Jest + React Testing Library
- Backend: Go standard testing (`go test`)
- Always propose a test plan before committing
