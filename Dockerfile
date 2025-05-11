FROM golang:1.24-alpine AS builder

WORKDIR /build

COPY . .

RUN --mount=type=cache,id=go-build,target=/root/.cache/go-build \
    --mount=type=cache,id=go-modules,target=/go \
    CGO_ENABLED=0 go build -tags viper_bind_struct -o app.bin -trimpath -ldflags='-s -w -extldflags "-static"'  cmd/app/main.go


FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=builder /build/app.bin .

EXPOSE 8080

CMD ["/app/app.bin"]
