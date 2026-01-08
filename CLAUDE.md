# CLAUDE.md

This file provides guidance to Claude Code when working with this repository.

## Project Overview

Go Redirect is a minimal Docker container that redirects HTTP requests to another URL. Written in Go, the final image is ~5MB.

## Common Commands

```bash
# Build Docker image
make build

# Run all tests
make test

# Clean up
make clean

# Start with Docker Compose
docker-compose up -d
```

## Architecture

- **cmd/redirect/main.go** - Go HTTP server with redirect logic
- **go.mod** - Go module definition
- **scanner_paths.json** - Default blocked scanner paths
- **Dockerfile** - Multi-stage build (scratch-based)
- **docker-compose.yml** - Docker Compose configuration

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `REDIRECT_TARGET` | Target URL (required) | - |
| `REDIRECT_CODE` | 301 or 302 | `301` |
| `PRESERVE_PATH` | Keep path and query | `true` |
| `BLOCK_SCANNERS` | Block scanner paths | `false` |
| `BLOCKED_PATHS_FILE` | Custom paths JSON | `/scanner_paths.json` |
| `LOG_LEVEL` | debug/info/warn/error/none | `info` |
| `PORT` | Server port | `8080` |

IPs are fully anonymized in logs (`x.x.x.x`).

## Git Hooks

Pre-commit hooks (via Husky):
- **pre-commit**: Runs `make build` to verify the build
- **commit-msg**: Validates conventional commit format

Valid commit types: `feat`, `fix`, `docs`, `build`, `ci`, `chore`, `refactor`, `revert`, `style`, `test`, `perf`

## CI/CD

GitHub Actions workflows:
- **test.yml** - Runs test suite on push/PR
- **docker-build-push.yml** - Builds and pushes images to GHCR
