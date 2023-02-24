ARG GOLANG_VERSION=1.20-alpine

FROM golang:${GOLANG_VERSION} AS build
WORKDIR /build
COPY . .

ARG VERSION="0.0.1"

RUN go mod vendor

RUN go build -o /bin/main -mod=vendor -ldflags "-X main.version=$VERSION"

FROM alpine:latest AS dev

COPY --from=build /bin/main /bin/main

EXPOSE 8080
ENTRYPOINT ["/bin/main"]
CMD []
