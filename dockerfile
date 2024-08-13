# syntax=docker/dockerfile:1-labs

ARG go_version=1.23

# workspace
FROM --platform=$BUILDPLATFORM golang:${go_version} AS workspace

ARG TARGETOS TARGETARCH
WORKDIR /work

RUN --mount=type=bind,target=. \
  --mount=type=cache,target=/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  set -eux; \
  go mod download; \
  CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -buildmode pie -buildvcs=false -ldflags "-s -w -extldflags '-static'" -trimpath -o /usr/bin/entry; \
  chmod +x /usr/bin/entry

# production
FROM --platform=$TARGETPLATFORM gcr.io/distroless/base:debug AS production

RUN ["/busybox/sh", "-c", "ln -s /busybox/sh /bin/sh"]
RUN ["/busybox/sh", "-c", "ln -s /busybox/env /usr/bin/env"]

COPY --chmod=755 --from=workspace /usr/bin/entry /usr/bin/entry

ENTRYPOINT ["/usr/bin/entry"]

# development
FROM --platform=$TARGETPLATFORM golang:${go_version} AS development

ENTRYPOINT set -eux; \
  go mod download; \
  CGO_ENABLED=0 go run -gcflags=all="-N -l" .
