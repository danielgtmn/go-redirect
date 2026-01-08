# Go Redirect

A minimal Docker container that redirects HTTP requests to another URL.

## Quick Start

```bash
docker run -d \
  -p 8080:8080 \
  -e REDIRECT_TARGET=https://example.com \
  ghcr.io/danielgtmn/go-redirect:latest
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `REDIRECT_TARGET` | Target URL (required) | - |
| `REDIRECT_CODE` | HTTP status code (301 or 302) | `301` |
| `PRESERVE_PATH` | Keep path and query string | `true` |
| `PORT` | Server port | `8080` |

## Examples

### Domain Redirect

Redirect `old-domain.com` to `new-domain.com`:

```bash
docker run -d \
  -p 80:8080 \
  -e REDIRECT_TARGET=https://new-domain.com \
  -e PRESERVE_PATH=true \
  ghcr.io/danielgtmn/go-redirect:latest
```

Requests to `old-domain.com/page?q=1` will redirect to `new-domain.com/page?q=1`.

### Temporary Redirect

```bash
docker run -d \
  -p 8080:8080 \
  -e REDIRECT_TARGET=https://example.com \
  -e REDIRECT_CODE=302 \
  ghcr.io/danielgtmn/go-redirect:latest
```

### Docker Compose

```yaml
services:
  redirect:
    image: ghcr.io/danielgtmn/go-redirect:latest
    ports:
      - "80:8080"
    environment:
      - REDIRECT_TARGET=https://new-domain.com
      - REDIRECT_CODE=301
      - PRESERVE_PATH=true
    restart: unless-stopped
```

## Health Check

The container exposes a `/health` endpoint:

```bash
curl http://localhost:8080/health
# OK
```

## Build Locally

```bash
docker build -t go-redirect .
docker run -p 8080:8080 -e REDIRECT_TARGET=https://example.com go-redirect
```

## Image Size

The final image is ~5MB (scratch base with static Go binary).
