# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY main.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o redirect .

# Runtime stage
FROM scratch

COPY --from=builder /app/redirect /redirect

EXPOSE 8080

ENTRYPOINT ["/redirect"]
