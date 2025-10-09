FROM caddy:2-alpine

# Copy Caddyfile
COPY Caddyfile /etc/caddy/Caddyfile

# Expose ports 80 and 443
EXPOSE 80 443
