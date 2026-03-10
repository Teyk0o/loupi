# 🍽️ Loupi

**Loupi** is a free, open-source progressive web app (PWA) that helps you track your meals, digestive symptoms, and daily wellness — all in one place.

Designed for people with digestive issues (IBS, food intolerances, etc.), Loupi helps you identify patterns between what you eat and how you feel.

🌐 **[loupi.app](https://loupi.app)**

---

## Features

- **Meal tracking** — Log meals with photos, descriptions, and categories
- **Symptom check-ins** — Record digestive symptoms linked to meals or standalone
- **Wellness tracking** — Monitor stress, mood, energy, sleep, exercise, and hydration daily
- **Journal & history** — Review your daily timeline and browse past entries
- **Secure & private** — Your health data is encrypted and GDPR-compliant

## Tech Stack

| Component     | Technology        |
|--------------|-------------------|
| Frontend     | Next.js + Tailwind CSS |
| Backend      | Go (REST API)     |
| Database     | PostgreSQL 16     |
| Cache        | Redis 7           |
| Proxy        | Caddy             |
| Deployment   | Docker Compose    |

## Getting Started

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/) & Docker Compose
- [Node.js](https://nodejs.org/) 20+ & [pnpm](https://pnpm.io/)
- [Go](https://go.dev/) 1.22+

### Development

```bash
# Clone the repository
git clone https://github.com/teyk0o/loupi.git
cd loupi

# Start infrastructure (PostgreSQL, Redis)
docker compose -f docker/docker-compose.dev.yml up -d

# Start the API
cd api && go run cmd/server/main.go

# Start the frontend
cd frontend && pnpm install && pnpm dev
```

## Project Structure

```
loupi/
├── frontend/          # Next.js PWA
├── api/               # Go REST API
│   ├── cmd/server/    # Entry point
│   └── internal/      # Application code
├── docker/            # Docker Compose & configs
├── docs/              # Documentation
├── CAHIER_DES_CHARGES.md
└── README.md
```

## License

This project is open source and available under the [MIT License](LICENSE).

## Contributing

Contributions are welcome! Please read the contributing guidelines before submitting a pull request.

---

Made with care for people who want to understand their gut better. 💚
