FROM golang:1.24-alpine AS builder
WORKDIR /build
COPY . .

RUN --mount=type=cache,id=go-build,target=/root/.cache/go-build \
--mount=type=cache,id=go-modules,target=/go \
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o migrate -trimpath -ldflags="-s -w -extldflags '-static'" ./db/migrate.go

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /build/db /db
COPY --from=builder /build/migrate /usr/bin
ENTRYPOINT ["/usr/bin/migrate"]
