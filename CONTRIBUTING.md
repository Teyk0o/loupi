# Contributing to Loupi

Thanks for your interest in contributing to Loupi! Here's everything you need to get started.

## Getting started

1. **Fork** the repository and clone your fork
2. Follow the [Getting started](README.md#getting-started) guide to set up the dev environment
3. Create a feature branch from `main`

## Branch naming

Use the format: `type/short-description`

```
feat/meal-photos
fix/auth-refresh-loop
chore/update-dependencies
docs/api-endpoints
```

## Commit messages

We follow [Conventional Commits](https://www.conventionalcommits.org/) strictly.

```
type(scope): short description

Optional longer description.
```

**Types:** `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`, `ci`, `perf`

**Scopes:** `frontend`, `api`, `docker`, `docs`, `deps`

**Examples:**
```
feat(frontend): add dark mode toggle to settings
fix(api): handle duplicate email registration gracefully
docs: update API endpoint reference in README
```

## Code guidelines

- **Language**: All code, comments, variable names, and documentation must be in **English**
- **UI text**: The user interface is in **French** (v1 targets France only)
- **Frontend**: TypeScript strict mode, Tailwind CSS v4, pnpm only (never npm/yarn)
- **Backend**: Standard Go project layout, Gin framework, pgx for PostgreSQL
- **Testing**: Always include a test plan when submitting a PR

## Pull requests

1. Keep PRs focused — one feature or fix per PR
2. Fill out the PR template completely
3. Make sure CI checks pass
4. Request a review from a maintainer

## Reporting issues

Use the [issue templates](https://github.com/teyk0o/loupi/issues/new/choose) to report bugs or request features. Please search existing issues first to avoid duplicates.

## Security

If you find a security vulnerability, please do **not** open a public issue. See [SECURITY.md](SECURITY.md) for responsible disclosure instructions.

## License

By contributing, you agree that your contributions will be licensed under the [Apache License 2.0](LICENSE).
