# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Caddy Redirect is a Docker container based on Caddy that configures reverse proxy and redirects via environment variables. Pre-built images are published to `ghcr.io/danielgtmn/caddy-redirect`.

## Common Commands

```bash
# Build Docker image
make build                    # Builds caddy-redirect:test image

# Run all tests (includes build)
make test                     # or ./test.sh

# Validate Caddyfile syntax only
make validate

# Clean up test containers and images
make clean

# Start with Docker Compose
docker-compose up -d          # Uses .env file for configuration

# Start test environment with mock backend
docker-compose -f docker-compose.test.yml up -d
```

## Architecture

The project is a minimal Docker container:

- **Caddyfile** - Caddy server configuration using environment variable substitution (`{$VAR_NAME:default}` syntax)
- **Dockerfile** - Simple Alpine-based image that copies Caddyfile into the container
- **docker-compose.yml** - Production compose file reading from `.env`
- **docker-compose.test.yml** - Test compose with nginx backend on port 8080

Key environment variables that configure the Caddyfile at runtime:
- `CADDY_DOMAIN` (required) - Main domain to serve
- `CADDY_UPSTREAM` - Backend server URL for reverse proxy
- `CADDY_TLS` - TLS configuration (e.g., `tls admin@example.com` for Let's Encrypt)
- `CADDY_SECURITY_HEADERS` - Additional security headers (default headers: X-Content-Type-Options, X-Frame-Options, Referrer-Policy)
- `CADDY_HEADERS`, `CADDY_RATE_LIMIT`, `CADDY_BASIC_AUTH`, `CADDY_REDIRECTS`, `CADDY_PROXY_HEADERS`, `CADDY_ADDITIONAL_DOMAINS` - Optional configuration blocks

## Git Hooks and Commit Convention

Pre-commit hooks (via Husky) enforce:
- **pre-commit**: Runs `make validate` to check Caddyfile syntax
- **commit-msg**: Validates conventional commit format via commitlint

Valid commit types: `feat`, `fix`, `docs`, `build`, `ci`, `chore`, `refactor`, `revert`, `style`, `test`, `perf`

Setup hooks after cloning: `pnpm install` (or `npm install`)

## CI/CD

GitHub Actions workflows:
- **test.yml** - Runs test suite on push/PR to main
- **docker-build-push.yml** - Builds and pushes images to GHCR
- **release.yml** - Semantic release triggered on release branch

Images are tagged: `latest`, `main`, `vX.Y.Z`, `vX.Y`, `vX`
