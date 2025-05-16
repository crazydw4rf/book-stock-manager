ARG BINARY_NAME=book-stock-manager-api

FROM golang:1.24-alpine AS builder

ARG BINARY_NAME
ARG APP_VERSION=0.0.1
ARG APP_ENV=production
ENV MODULE_NAME=github.com/crazydw4rf/book-stock-manager
ENV GOBUILD_LDFLAGS="-X ${MODULE_NAME}/internal/config.APP_VERSION=${APP_VERSION} -X ${MODULE_NAME}/internal/config.APP_ENV=${APP_ENV}"

WORKDIR /build

COPY . .

RUN --mount=type=cache,id=go-build,target=/root/.cache/go-build \
  --mount=type=cache,id=go-modules,target=/go \
  CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ${BINARY_NAME} -trimpath -ldflags="-s -w -extldflags '-static' ${GOBUILD_LDFLAGS}"  ./cmd/app/*

FROM gcr.io/distroless/static-debian12:nonroot

ARG BINARY_NAME

COPY --from=builder /build/${BINARY_NAME} /usr/bin/

EXPOSE 8080

CMD ["book-stock-manager-api"]
