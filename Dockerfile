FROM golang:1.22-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG VERSION=dev
RUN go build \
    -tags native \
    -ldflags "-X github.com/Ka0s-Klaus/Klaus-antipatterns-search/cmd/antipatterns.version=${VERSION}" \
    -o /antipatterns \
    ./cmd/antipatterns

# ── Runtime image ──────────────────────────────────────────────────────────
FROM alpine:3.20

RUN addgroup -S app && adduser -S app -G app
COPY --from=builder /antipatterns /usr/local/bin/antipatterns

USER app
ENTRYPOINT ["antipatterns"]
CMD ["--help"]
