FROM caddy:2-alpine

# Create non-root user for security
# The base image already includes a properly configured caddy user
# with necessary capabilities for binding to ports 80 and 443

# Copy Caddyfile
COPY Caddyfile /etc/caddy/Caddyfile

# Use non-root user
USER caddy

# Expose ports 80 and 443
EXPOSE 80 443

# Healthcheck using Caddy's admin API (works regardless of CADDY_DOMAIN)
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:2019/config/ || exit 1
