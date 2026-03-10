<p align="center">
  <a href="https://loupi.app">
    <img src="docs/logo.png" alt="Loupi" width="220" />
  </a>
</p>

<p align="center">
  <strong>Track your meals, symptoms & wellness — all in one place.</strong><br />
  Free, open-source, and privacy-first.
</p>

<p align="center">
  <a href="https://loupi.app"><img src="https://img.shields.io/badge/loupi.app-Visit-48C7B0?style=for-the-badge&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSIyNCIgaGVpZ2h0PSIyNCIgdmlld0JveD0iMCAwIDI0IDI0IiBmaWxsPSJub25lIiBzdHJva2U9IndoaXRlIiBzdHJva2Utd2lkdGg9IjIiIHN0cm9rZS1saW5lY2FwPSJyb3VuZCIgc3Ryb2tlLWxpbmVqb2luPSJyb3VuZCI+PGNpcmNsZSBjeD0iMTIiIGN5PSIxMiIgcj0iMTAiLz48cGF0aCBkPSJNMiAxMmgyMCIvPjxwYXRoIGQ9Ik0xMiAyYTE1LjMgMTUuMyAwIDAgMSA0IDEwIDE1LjMgMTUuMyAwIDAgMS00IDEwIDE1LjMgMTUuMyAwIDAgMS00LTEwIDE1LjMgMTUuMyAwIDAgMSA0LTEwIi8+PC9zdmc+" alt="Website" /></a>
  <a href="https://github.com/teyk0o/loupi"><img src="https://img.shields.io/github/stars/teyk0o/loupi?style=for-the-badge&color=48C7B0&logo=github&logoColor=white" alt="Stars" /></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-Apache_2.0-48C7B0?style=for-the-badge" alt="License" /></a>
  <a href="https://github.com/teyk0o/loupi/issues"><img src="https://img.shields.io/github/issues/teyk0o/loupi?style=for-the-badge&color=48C7B0&logo=github&logoColor=white" alt="Issues" /></a>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Next.js_16-213547?style=flat-square&logo=next.js&logoColor=white" alt="Next.js" />
  <img src="https://img.shields.io/badge/Tailwind_CSS_4-213547?style=flat-square&logo=tailwindcss&logoColor=white" alt="Tailwind" />
  <img src="https://img.shields.io/badge/Go-213547?style=flat-square&logo=go&logoColor=white" alt="Go" />
  <img src="https://img.shields.io/badge/PostgreSQL_16-213547?style=flat-square&logo=postgresql&logoColor=white" alt="PostgreSQL" />
  <img src="https://img.shields.io/badge/Redis_7-213547?style=flat-square&logo=redis&logoColor=white" alt="Redis" />
  <img src="https://img.shields.io/badge/Docker-213547?style=flat-square&logo=docker&logoColor=white" alt="Docker" />
</p>

---

## About

Loupi is a progressive web app designed for people with digestive issues (IBS, food intolerances, Crohn's, etc.). It helps you identify patterns between what you eat and how you feel by keeping everything in one daily journal.

- **100% free** — No premium tier, no ads, no data selling
- **Privacy-first** — Health data encrypted at rest (AES-256-GCM), GDPR-compliant
- **Open source** — Apache 2.0 licensed, contributions welcome
- **French UI** — Designed for the French market (v1)

## Features

| Feature | Description |
|---|---|
| **Meal journal** | Log meals with descriptions, categories, date & time. Navigate day by day. |
| **Symptom tracking** | Record digestive symptoms linked to meals (with delay) or standalone entries. |
| **Wellness dashboard** | Track stress, mood, energy, sleep quality, hydration, and sport sessions daily. |
| **Custom options** | Personalize your symptom types, meal categories, and sport types from settings. |
| **Responsive design** | Mobile-first with bottom navigation; desktop sidebar layout on larger screens. |
| **Dark mode** | Automatic light/dark theme based on system preference. |
| **PWA-ready** | Installable as a native-like app on mobile and desktop. |

## Architecture

```
loupi/
├── frontend/              # Next.js 16 PWA (TypeScript + Tailwind CSS 4)
│   ├── src/app/           # App Router pages & layouts
│   ├── src/components/    # UI components, layout, settings
│   ├── src/hooks/         # Custom React hooks (auth, custom options)
│   └── src/lib/           # API client, utilities
├── api/                   # Go REST API
│   ├── cmd/server/        # Entry point
│   ├── internal/          # Handlers, services, models, middleware
│   └── migrations/        # Versioned PostgreSQL migrations
├── docker/                # Docker Compose configs
└── docs/                  # Documentation & assets
```

### Tech stack

| Layer | Technology | Purpose |
|---|---|---|
| **Frontend** | Next.js 16 + Tailwind CSS 4 | Server-side rendering, PWA, responsive UI |
| **Backend** | Go + Gin | REST API, JWT auth, business logic |
| **Database** | PostgreSQL 16 | Persistent storage with JSONB for flexible data |
| **Cache** | Redis 7 | Session management, rate limiting |
| **Proxy** | Caddy | Automatic HTTPS (TLS 1.3), reverse proxy |
| **Deploy** | Docker Compose | Single-command deployment on VPS |

### API endpoints

```
Auth        POST /v1/auth/register, /login, /refresh
            GET  /v1/auth/me       DELETE /v1/auth/account

Meals       GET  /v1/meals?date=   POST /v1/meals
            GET  /v1/meals/:id     PUT  /v1/meals/:id     DELETE /v1/meals/:id
            GET  /v1/meals/:id/check-ins   POST /v1/meals/:id/check-ins

Symptoms    GET  /v1/symptoms?date=   POST /v1/symptoms
            PUT  /v1/symptoms/:id     DELETE /v1/symptoms/:id

Wellness    GET  /v1/wellness?date=   POST /v1/wellness   PUT /v1/wellness/:id

Options     GET  /v1/options/:category        POST /v1/options/:category
            PUT  /v1/options/:category/reorder
            PUT  /v1/options/item/:id         DELETE /v1/options/item/:id
```

## Deploy with Docker

Loupi ships as a **single Docker image** containing the API, frontend, and Caddy reverse proxy. You only need to provide PostgreSQL and Redis.

```bash
# 1. Clone and configure
git clone https://github.com/teyk0o/loupi.git
cd loupi
cp .env.example .env   # edit with your secrets

# 2. Start everything
docker compose up -d
```

Or pull the image directly:

```bash
docker pull ghcr.io/teyk0o/loupi:latest
```

> **Supported architectures:** `linux/amd64`, `linux/arm64`

See [`.env.example`](.env.example) for all configuration options. Set `LOUPI_DOMAIN` to your domain (e.g. `loupi.example.com`) for automatic HTTPS via Caddy.

## Development

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/) & Docker Compose
- [Node.js](https://nodejs.org/) 22+ & [pnpm](https://pnpm.io/)
- [Go](https://go.dev/) 1.25+

### Quick start

```bash
# Start everything (API, frontend, PostgreSQL, Redis) with hot-reload
docker compose -f docker/docker-compose.dev.yml up -d
```

The frontend is available at `http://localhost:3000` and the API at `http://localhost:8080`.

### Manual setup

```bash
# Start infrastructure only
docker compose -f docker/docker-compose.dev.yml up -d db redis

# Start the API (in a separate terminal)
cd api
cp .env.example .env   # adjust if needed
go run cmd/server/main.go

# Start the frontend (in a separate terminal)
cd frontend
pnpm install
pnpm dev
```

## Security

- **HTTPS only** with TLS 1.3 via Caddy (production)
- **JWT authentication** with short-lived access tokens (15 min) and refresh tokens (7 days)
- **bcrypt** password hashing (cost >= 12)
- **AES-256-GCM** encryption for sensitive health data at rest
- **Rate limiting** on authentication endpoints
- **Security headers**: CSP, HSTS, X-Frame-Options, X-Content-Type-Options
- **CSRF protection** and input validation on all endpoints
- **GDPR-compliant**: users can delete their account and all associated data

## Contributing

Contributions are welcome! Here's how to get started:

1. Fork the repository
2. Create a feature branch (`feat/my-feature`)
3. Follow [Conventional Commits](https://www.conventionalcommits.org/) for all commit messages
4. Write code and comments in **English**
5. Submit a pull request

Please read the [CLAUDE.md](CLAUDE.md) for detailed coding conventions.

## License

This project is open source under the [Apache License 2.0](LICENSE).

---

<p align="center">
  Made with care for people who want to understand their gut better. 💚
</p>
