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
| `REDIRECT_TARGET` | Target URL | - |
| `REDIRECT_DOMAIN_MAP` | Domain mapping (`source:target`) | - |
| `REDIRECT_SCHEME` | Scheme for domain mapping redirects | `https` |
| `REDIRECT_CODE` | HTTP status code (301 or 302) | `301` |
| `PRESERVE_PATH` | Keep path and query string | `true` |
| `BLOCK_SCANNERS` | Block common scanner paths with 404 | `false` |
| `BLOCKED_PATHS_FILE` | Custom JSON file with blocked paths | `/scanner_paths.json` |
| `LOG_LEVEL` | debug, info, warn, error, none | `info` |
| `PORT` | Server port | `8080` |

Either `REDIRECT_TARGET` or `REDIRECT_DOMAIN_MAP` is required. Both can be set together — domain mapping takes priority, with `REDIRECT_TARGET` as fallback for unmatched hosts.

**Privacy:** IPs are fully anonymized in logs (`x.x.x.x`).

## Scanner Blocking

Enable `BLOCK_SCANNERS=true` to return 404 for common vulnerability scanner paths like `.env`, `swagger`, `actuator`, `wp-admin`, etc.

The default blocked paths are defined in `scanner_paths.json`. You can mount your own file:

```bash
docker run -d \
  -p 8080:8080 \
  -e REDIRECT_TARGET=https://example.com \
  -e BLOCK_SCANNERS=true \
  -v ./my_paths.json:/scanner_paths.json \
  ghcr.io/danielgtmn/go-redirect:latest
```

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

### Domain Mapping (Wildcard Subdomains)

Redirect `*.domain.com` to `*.domain.de`, preserving subdomains and paths:

```bash
docker run -d \
  -p 80:8080 \
  -e REDIRECT_DOMAIN_MAP=domain.com:domain.de \
  ghcr.io/danielgtmn/go-redirect:latest
```

- `domain.com/page` → `https://domain.de/page`
- `blog.domain.com/post?id=1` → `https://blog.domain.de/post?id=1`
- `a.b.domain.com/` → `https://a.b.domain.de/`

Requests with unmatched hosts return 404 (or redirect to `REDIRECT_TARGET` if set).

### With Scanner Blocking

```bash
docker run -d \
  -p 80:8080 \
  -e REDIRECT_TARGET=https://new-domain.com \
  -e PRESERVE_PATH=true \
  -e BLOCK_SCANNERS=true \
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
      - BLOCK_SCANNERS=true
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
