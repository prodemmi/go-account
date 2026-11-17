# syntax=docker/dockerfile:1
# Multi-stage build: compile with golang, run on minimal alpine
FROM golang:1.27.0-alpine AS builder

WORKDIR /app/src

# Configure Alpine apk to use Runflare mirror (per docs.runflare.com/mirror/mirror-alpine/)
RUN APK_VER=$(cat /etc/alpine-release | cut -d. -f1-2) && \
    echo "http://mirror-linux.runflare.com/alpine/v${APK_VER}/main" > /etc/apk/repositories && \
    echo "http://mirror-linux.runflare.com/alpine/v${APK_VER}/community" >> /etc/apk/repositories && \
    apk update

# Set Go module proxy (Iran mirror)
ENV GOPROXY=https://go.iranserver.com/repository/go/
ENV GOFLAGS=-mod=mod

# Install build dependencies
RUN --mount=type=cache,target=/var/cache/apk \
    apk add --no-cache git

# Copy go.mod/go.sum first for module download layer caching
COPY go.mod go.sum ./

# Download Go modules first (layer caching with buildkit cache)
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

# Copy source code (context is src/)
COPY . ./

# Build the binary (buildkit cache for incremental builds)
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -v -o /app/server .

# Runtime stage
FROM alpine:latest

# Configure Alpine apk to use Runflare mirror (per docs.runflare.com/mirror/mirror-alpine/)
RUN APK_VER=$(cat /etc/alpine-release | cut -d. -f1-2) && \
    echo "http://mirror-linux.runflare.com/alpine/v${APK_VER}/main" > /etc/apk/repositories && \
    echo "http://mirror-linux.runflare.com/alpine/v${APK_VER}/community" >> /etc/apk/repositories && \
    apk update

WORKDIR /app

# Copy the compiled binary
COPY --from=builder /app/server .

# Copy static CSS files
COPY --from=builder /app/src/internal/adapters/inbound/http/pages/css /app/static/css

EXPOSE 9090

CMD ["./server"]
