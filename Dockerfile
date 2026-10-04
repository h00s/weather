FROM golang:alpine AS backend
WORKDIR /app
ENV CGO_ENABLED=0
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend ./
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /out/weather

FROM oven/bun:latest AS frontend
WORKDIR /app
COPY frontend/package.json frontend/bun.lock ./
RUN bun install --frozen-lockfile
COPY frontend ./
RUN bun run build

FROM gcr.io/distroless/static-debian13:latest
WORKDIR /app
COPY --from=backend /out/weather ./
COPY --from=frontend /app/build ./public
# Raptor binds 127.0.0.1 by default, which is unreachable from outside the container.
ENV SERVER_ADDRESS=0.0.0.0
# Behind the HTTPS reverse proxy: rate-limit and log the client's address, not the proxy's.
ENV SERVER_IP_EXTRACTOR=x-forwarded-for
# HSTS starts small; raise it to 31536000 once HTTPS is settled.
ENV APP_SECURE_HSTS_MAX_AGE=300

EXPOSE 3000

ENTRYPOINT ["./weather"]
