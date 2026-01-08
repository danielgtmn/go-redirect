# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY cmd/ ./cmd/

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o redirect ./cmd/redirect

# Runtime stage
FROM scratch

COPY --from=builder /app/redirect /redirect
COPY scanner_paths.json /scanner_paths.json

EXPOSE 8080

ENTRYPOINT ["/redirect"]
