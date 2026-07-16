# Multi-stage Dockerfile for building and running the AMJ HUB server
# Builds a static Go binary in a lightweight builder and copies it into
# a minimal runtime image.

FROM golang:1.21-alpine AS builder
RUN apk add --no-cache git ca-certificates
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -o /out/amjhub-server ./

FROM alpine:3.18
RUN apk add --no-cache ca-certificates wget
WORKDIR /app
COPY --from=builder /out/amjhub-server /app/amjhub-server
RUN addgroup -S app && adduser -S app -G app
USER app
EXPOSE 8000
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s CMD wget -qO- --timeout=2 http://127.0.0.1:8000/health || exit 1
ENTRYPOINT ["/app/amjhub-server"]
