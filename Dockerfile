# Ephemeral runner for `cli`. Intended to be used with `--rm` so nothing
# (binaries, the local clid socket, the pass key you type) survives past
# the single `docker run`.
#
# Build stage runs natively on the builder's own platform ($BUILDPLATFORM)
# and cross-compiles the Go binaries for the requested $TARGETOS/$TARGETARCH
# (CGO disabled) — no QEMU emulation needed even for multi-platform builds,
# since the final stage only COPYs, it never executes a foreign-arch binary.
FROM --platform=$BUILDPLATFORM golang:1-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG TARGETOS TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /out/cli ./cmd/cli \
 && CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /out/clid ./cmd/clid

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=build /out/cli /out/clid /usr/local/bin/
COPY docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
RUN chmod +x /usr/local/bin/docker-entrypoint.sh
ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]
