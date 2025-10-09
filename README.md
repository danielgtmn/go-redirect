# Caddy Redirect Docker Container

A flexible Docker container based on Caddy that configures redirects and reverse proxy functions via environment variables.

## Quick Start

### Using Pre-built Images from GHCR

The easiest way is to use pre-built images from GitHub Container Registry:

```bash
# Using latest version
docker run -d \
  --name caddy-redirect \
  -p 80:80 \
  -p 443:443 \
  -e CADDY_DOMAIN=yourdomain.com \
  -e CADDY_EMAIL=admin@yourdomain.com \
  -e CADDY_UPSTREAM=http://your-app:8080 \
  -v caddy_data:/data \
  -v caddy_config:/config \
  ghcr.io/YOUR_USERNAME/caddy-redirect:latest
```

Replace `YOUR_USERNAME` with your GitHub username.

### Option 1: Docker Compose (Recommended)

1. Copy the example configuration:
   ```bash
   cp env.example .env
   ```

2. Edit the `.env` file with your settings

3. Start the container:
   ```bash
   docker-compose up -d
   ```

### Option 2: Direct Docker Run (Local Build)

You can also build and run the container locally:

```bash
# Build the image
docker build -t caddy-redirect:latest .

# Run the container
docker run -d \
  --name caddy-redirect \
  -p 80:80 \
  -p 443:443 \
  -e CADDY_DOMAIN=yourdomain.com \
  -e CADDY_EMAIL=admin@yourdomain.com \
  -e CADDY_UPSTREAM=http://your-app:8080 \
  -v caddy_data:/data \
  -v caddy_config:/config \
  caddy-redirect:latest
```

## Environment Variables

### Basic Configuration

| Variable | Description | Example |
|----------|-------------|---------|
| `CADDY_DOMAIN` | Main domain (required) | `example.com` |
| `CADDY_EMAIL` | Email for Let's Encrypt | `admin@example.com` |
| `CADDY_UPSTREAM` | Target server for redirection | `http://localhost:8080` |

### Advanced Configuration

| Variable | Description | Example |
|----------|-------------|---------|
| `CADDY_HEADERS` | Additional HTTP headers | `header X-Forwarded-Proto {scheme}` |
| `CADDY_RATE_LIMIT` | Rate limiting rules | `rate_limit { zone static { key {remote} window 1m burst 10 } }` |
| `CADDY_BASIC_AUTH` | Basic authentication | `basicauth { user $2a$14$... }` |
| `CADDY_REDIRECTS` | URL redirects | `redir /old /new permanent` |
| `CADDY_PROXY_HEADERS` | Proxy-specific headers | `header_up Host {host}` |
| `CADDY_ADDITIONAL_DOMAINS` | Additional domains | `www.example.com,api.example.com` |

## Examples

### Simple Redirection

**Docker Compose:**
```bash
# .env
CADDY_DOMAIN=my-app.de
CADDY_UPSTREAM=http://localhost:3000
```

**Docker Run:**
```bash
docker run -d --name caddy-redirect \
  -p 80:80 -p 443:443 \
  -e CADDY_DOMAIN=my-app.de \
  -e CADDY_UPSTREAM=http://localhost:3000 \
  -v caddy_data:/data \
  caddy-redirect:latest
```

### With HTTPS and Basic Auth

**Docker Compose:**
```bash
# .env
CADDY_DOMAIN=secure-app.com
CADDY_EMAIL=admin@secure-app.com
CADDY_UPSTREAM=http://localhost:8080
CADDY_BASIC_AUTH=basicauth { admin $2a$14$9xv0MtX8mJ2N/.5jUwJcO8qhQ5QE2ZrGj5JvZK5vT5rXnHkYjx6y }
```

**Docker Run:**
```bash
docker run -d --name caddy-redirect \
  -p 80:80 -p 443:443 \
  -e CADDY_DOMAIN=secure-app.com \
  -e CADDY_EMAIL=admin@secure-app.com \
  -e CADDY_UPSTREAM=http://localhost:8080 \
  -e CADDY_BASIC_AUTH="basicauth { admin \$2a\$14\$9xv0MtX8mJ2N/.5jUwJcO8qhQ5QE2ZrGj5JvZK5vT5rXnHkYjx6y }" \
  -v caddy_data:/data \
  -v caddy_config:/config \
  caddy-redirect:latest
```

### With Rate Limiting and Redirects

**Docker Compose:**
```bash
# .env
CADDY_DOMAIN=api.example.com
CADDY_EMAIL=admin@example.com
CADDY_UPSTREAM=http://localhost:8080
CADDY_RATE_LIMIT=rate_limit { zone static { key {remote} window 1m burst 100 } }
CADDY_REDIRECTS=redir /v1 /v2 permanent
CADDY_HEADERS=header X-Forwarded-Proto {scheme}
```

**Docker Run:**
```bash
docker run -d --name caddy-redirect \
  -p 80:80 -p 443:443 \
  -e CADDY_DOMAIN=api.example.com \
  -e CADDY_EMAIL=admin@example.com \
  -e CADDY_UPSTREAM=http://localhost:8080 \
  -e CADDY_RATE_LIMIT="rate_limit { zone static { key {remote} window 1m burst 100 } }" \
  -e CADDY_REDIRECTS="redir /v1 /v2 permanent" \
  -e CADDY_HEADERS="header X-Forwarded-Proto {scheme}" \
  -v caddy_data:/data \
  -v caddy_config:/config \
  caddy-redirect:latest
```

### Generate Password for Basic Auth

```bash
# Install caddy locally or use Docker
docker run --rm -it caddy:2-alpine caddy hash-password
# Enter your desired password
```

## Docker Commands

### Docker Compose

```bash
# Build and start container
docker-compose up -d

# View logs
docker-compose logs -f caddy

# Stop container
docker-compose down

# Restart container with new build
docker-compose up -d --build
```

### Direct Docker Run

**Important:** For HTTPS certificates to persist, use named Docker volumes.

```bash
# Create persistent volumes for certificates
docker volume create caddy_data
docker volume create caddy_config

# Build the image first
docker build -t caddy-redirect:latest .

# Start container
docker run -d --name caddy-redirect \
  -p 80:80 -p 443:443 \
  -e CADDY_DOMAIN=yourdomain.com \
  -e CADDY_UPSTREAM=http://your-app:8080 \
  -v caddy_data:/data \
  -v caddy_config:/config \
  caddy-redirect:latest

# View logs
docker logs -f caddy-redirect

# Stop container
docker stop caddy-redirect

# Remove container
docker rm caddy-redirect

# Update and restart (keeps certificates)
docker build -t caddy-redirect:latest .
docker run -d --name caddy-redirect \
  -p 80:80 -p 443:443 \
  -e CADDY_DOMAIN=yourdomain.com \
  -e CADDY_UPSTREAM=http://your-app:8080 \
  -v caddy_data:/data \
  -v caddy_config:/config \
  caddy-redirect:latest
```

## CI/CD with GitHub Actions

This repository includes GitHub Actions workflows that automatically build and push Docker images to GitHub Container Registry (GHCR).

### Automatic Builds

Images are automatically built and pushed when:
- Pushing to the `main` branch
- Creating tags (e.g., `v1.0.0`)
- Opening pull requests

### Image Tags

Available tags in `ghcr.io/YOUR_USERNAME/caddy-redirect`:
- `latest` - Latest build from main branch
- `main` - Latest build from main branch
- `v1.2.3` - Specific version tags
- `v1.2` - Major.minor version tags
- `v1` - Major version tags

### Using Built Images

```bash
# Pull and run the latest version
docker run -d \
  --name caddy-redirect \
  -p 80:80 -p 443:443 \
  -e CADDY_DOMAIN=yourdomain.com \
  -e CADDY_UPSTREAM=http://your-app:8080 \
  ghcr.io/YOUR_USERNAME/caddy-redirect:latest
```

### Manual Workflow Triggers

You can also manually trigger the build workflow from the GitHub Actions tab.

## HTTPS Certificates

When `CADDY_EMAIL` is set, HTTPS certificates are automatically obtained from Let's Encrypt. Certificates are stored in Docker volumes and reused on restarts.

## Multiple Domains

Use `CADDY_ADDITIONAL_DOMAINS` to serve multiple domains:

```bash
CADDY_ADDITIONAL_DOMAINS=www.example.com,api.example.com,staging.example.com
```

## Troubleshooting

### Container won't start
- Check the `.env` file for syntax errors
- Ensure required environment variables are set

### HTTPS doesn't work
- Ensure `CADDY_EMAIL` is set
- Check that port 443 is not blocked
- Wait up to 5 minutes for certificate issuance

### Upstream server unreachable
- Ensure the upstream server is running
- Check the `CADDY_UPSTREAM` URL
- Use `docker-compose logs` for detailed errors

## Caddy Documentation

For advanced configuration options, see the [official Caddy documentation](https://caddyserver.com/docs/).
