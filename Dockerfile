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
RUN bun install
COPY frontend ./
RUN bun run build

FROM gcr.io/distroless/static-debian12:latest
WORKDIR /app
COPY --from=backend /out/weather ./
COPY --from=frontend /app/build ./public

EXPOSE 3000

ENTRYPOINT ["./weather"]