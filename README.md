# Caddy Redirect Docker Container

[![Docker Build](https://github.com/danielgtmn/caddy-redirect/actions/workflows/docker-build-push.yml/badge.svg)](https://github.com/danielgtmn/caddy-redirect/actions)
[![Docker Image](https://img.shields.io/badge/docker-ghcr.io%2Fdanielgtmn%2Fcaddy--redirect-blue)](https://ghcr.io/danielgtmn/caddy-redirect)

A flexible Docker container based on Caddy that configures redirects and reverse proxy functions via environment variables.

**🚀 Pre-built images available at: [`ghcr.io/danielgtmn/caddy-redirect`](https://ghcr.io/danielgtmn/caddy-redirect)**

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
  ghcr.io/danielgtmn/caddy-redirect:latest
```

Replace `danielgtmn` with your GitHub username if you're using a fork.

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

### Automatic Builds & Releases

Images are automatically built and pushed when:
- Pushing to the `main` branch (development builds)
- Opening pull requests
- After successful releases (versioned builds)

**Releases are created** by pushing version tags (e.g., `v1.0.0`) to the repository

### Image Tags

Available tags in `ghcr.io/danielgtmn/caddy-redirect`:
- `latest` - Latest build from main branch
- `main` - Latest build from main branch
- `v1.2.3` - Specific version tags (auto-generated)
- `v1.2` - Major.minor version tags
- `v1` - Major version tags

### Release Workflow

```bash
# Development on main branch
git checkout main
# Make your changes...
git add .
git commit -m "feat: add support for custom headers"
git push origin main

# → CI: Tests run, Docker image built and pushed

# When ready for release, create release branch
git checkout -b release
git push origin release

# → Release: Semantic release creates GitHub release
#   1. Analyzes commits since last release
#   2. Bumps version (patch/minor/major)
#   3. Creates GitHub release with changelog
#   4. Triggers build of versioned Docker image

# After release, merge back to main
git checkout main
git merge release
git push origin main
```

### Using Built Images

```bash
# Pull and run the latest version
docker run -d \
  --name caddy-redirect \
  -p 80:80 -p 443:443 \
  -e CADDY_DOMAIN=yourdomain.com \
  -e CADDY_UPSTREAM=http://your-app:8080 \
  ghcr.io/danielgtmn/caddy-redirect:latest
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

## Testing

This repository includes comprehensive tests to ensure everything works correctly.

### Running Tests Locally

```bash
# Run all tests
./test.sh

# Test with Docker Compose
docker-compose -f docker-compose.test.yml up -d
curl http://localhost:8080
docker-compose -f docker-compose.test.yml down
```

### GitHub Actions Tests

Tests run automatically on:
- **Push to main**: Full test suite
- **Pull Requests**: Full test suite
- **Manual trigger**: Via Actions tab

The tests cover:
- ✅ Caddyfile syntax validation
- ✅ Docker build success
- ✅ Container startup
- ✅ HTTP proxy functionality
- ✅ Environment variable processing

### Test Coverage

| Test Type | Description | Status |
|-----------|-------------|---------|
| Syntax Validation | Caddyfile syntax check | ✅ Automated |
| Docker Build | Container builds successfully | ✅ Automated |
| Startup Test | Container starts without errors | ✅ Automated |
| HTTP Proxy | Requests are properly forwarded | ✅ Automated |
| ENV Vars | Configuration via environment | ✅ Automated |
| Integration | Full stack with backend | ✅ Manual/Local |

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Development Setup

1. **Install dependencies:**
   ```bash
   pnpm install  # or npm install
   ```

2. **Initialize Husky hooks:**
   ```bash
   pnpm run prepare  # or npm run prepare
   ```

### Release Process

**For Maintainers:**

1. **Ensure all changes are merged** to main branch
2. **Run the release script** to create and push a version tag
3. **Manually trigger** the GitHub Release workflow

```bash
# 1. Create and push version tag
./release.sh

# 2. Go to GitHub Actions → "Create Release" workflow
# 3. Click "Run workflow"
# 4. Enter the version number (without 'v')
# 5. Choose if it's a pre-release
```

**Alternative manual approach:**
```bash
# Create tag manually
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# Then trigger release workflow manually
```

**What happens:**
- ✅ Git Tag wird erstellt und gepusht
- ✅ Docker Images werden automatisch gebaut (latest + version)
- ✅ GitHub Release wird manuell erstellt mit Release Notes
- ✅ Release Notes werden automatisch aus Commits generiert

### Branch Protection

This repository uses branch protection to ensure code quality:

#### `main` Branch (Development)
- ✅ **Require PR reviews** - All changes need approval
- ✅ **Require status checks** - Tests must pass
- ✅ **Require up-to-date branches** - Must merge latest changes

#### `release` Branch (Production)
- ✅ **Require PR reviews** - Release changes need approval
- ✅ **Require status checks** - All tests must pass
- ✅ **Restrict pushes** - Only maintainers can push directly
- ✅ **Require up-to-date branches** - Must be current

**Setup:** Go to Settings → Branches → Add rule for each branch

```bash
# Create release branch from main
git checkout main
git pull origin main
git checkout -b release
git push origin release

# After successful release, merge back
git checkout main
git merge release
git push origin main
```

### Commit Convention

This project uses descriptive commit messages. While not strictly enforced, following conventional patterns helps with release notes generation.

**Valid commit types:**
```bash
# Features
feat: add new functionality

# Bug fixes
fix: resolve issue with proxy headers

# Documentation
docs: update README with new examples

# Build/CI changes
build: update Dockerfile
ci: modify GitHub Actions workflow

# Breaking changes
feat!: change API that breaks backward compatibility

# Other changes
chore: update dependencies
refactor: restructure code
test: add tests
style: format code
perf: improve performance
```

**Pre-commit hooks will:**
- ✅ Validate Caddyfile syntax
- ✅ Check commit message format

**Invalid commits will be rejected automatically.**

### For Fork Users

If you're using a fork of this repository, update the image references in the examples to point to your own GHCR registry:

```bash
# Replace danielgtmn with your GitHub username
ghcr.io/YOUR_USERNAME/caddy-redirect:latest
```

## Caddy Documentation

For advanced configuration options, see the [official Caddy documentation](https://caddyserver.com/docs/).

---

**Repository**: [github.com/danielgtmn/caddy-redirect](https://github.com/danielgtmn/caddy-redirect) |
**Docker Images**: [ghcr.io/danielgtmn/caddy-redirect](https://ghcr.io/danielgtmn/caddy-redirect)
