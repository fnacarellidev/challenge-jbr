FROM golang:1.22 as builder

ENV CGO_ENABLED=0

WORKDIR /build

RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.25.0

RUN go env -w GOCACHE=/go-cache

COPY . .
WORKDIR ./backend
RUN sqlc generate
RUN --mount=type=cache,target=/go-cache go build -o app

FROM cgr.dev/chainguard/wolfi-base:latest@sha256:3eff851ab805966c768d2a8107545a96218426cee1e5cc805865505edbe6ce92

WORKDIR /app
RUN apk add curl
COPY --from=builder /build/backend/app .

ENTRYPOINT ["./app"]
