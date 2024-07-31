FROM golang:1.22 as builder

ENV CGO_ENABLED=0

WORKDIR /build

RUN go env -w GOCACHE=/go-cache

COPY . .
WORKDIR ./graphql-api
RUN --mount=type=cache,target=/go-cache go build -o app

FROM cgr.dev/chainguard/wolfi-base:latest@sha256:3eff851ab805966c768d2a8107545a96218426cee1e5cc805865505edbe6ce92

WORKDIR /app
COPY --from=builder /build/graphql-api/app .

ENTRYPOINT ["./app"]
