# Security Policy

## Reporting a vulnerability

If you discover a security vulnerability in Loupi, please report it responsibly.

**Do not open a public GitHub issue.**

Instead, use GitHub's private vulnerability reporting:

**[Report a vulnerability](https://github.com/Teyk0o/loupi/security/advisories/new)**

Include:
- A description of the vulnerability
- Steps to reproduce
- Potential impact
- Any suggested fix (optional)

We will acknowledge your report within **48 hours** and aim to provide a fix within **7 days** for critical issues.

## Scope

The following are in scope for security reports:

- Authentication and authorization flaws
- Data exposure or leakage
- Injection vulnerabilities (SQL, XSS, CSRF, etc.)
- Encryption or hashing weaknesses
- Access control bypasses
- Server misconfiguration

## Out of scope

- Denial of service attacks
- Social engineering
- Issues in third-party dependencies (report these upstream)
- Issues that require physical access to a user's device

## Security measures

Loupi implements the following security practices:

- HTTPS only (TLS 1.3 via Caddy)
- JWT with short-lived access tokens (15 min) and refresh tokens (7 days)
- bcrypt password hashing (cost >= 12)
- AES-256-GCM encryption for sensitive health data at rest
- Rate limiting on authentication endpoints
- Security headers (CSP, HSTS, X-Frame-Options, X-Content-Type-Options)
- Input validation and sanitization on all endpoints
- CASCADE deletion for complete account removal (GDPR)

## Acknowledgments

We appreciate responsible disclosure and will credit reporters (with permission) in our release notes.
