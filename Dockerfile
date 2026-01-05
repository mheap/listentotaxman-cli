# Build stage
FROM golang:1.25.5-alpine AS builder

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -X main.Version=dev -X main.GitCommit=local -X main.BuildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    -o listentotaxman .

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /build/listentotaxman /usr/local/bin/listentotaxman

RUN adduser -D -u 1000 appuser
USER appuser

ENTRYPOINT ["/usr/local/bin/listentotaxman"]
CMD ["--help"]
