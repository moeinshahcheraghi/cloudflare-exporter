FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /build/cloudflare-exporter \
    ./cmd/exporter

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata wget

RUN addgroup -g 1000 exporter && \
    adduser -D -u 1000 -G exporter exporter

WORKDIR /app

COPY --from=builder /build/cloudflare-exporter .

RUN chown -R exporter:exporter /app

USER exporter

EXPOSE 9199

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:9199/health || exit 1

ENTRYPOINT ["./cloudflare-exporter"]
