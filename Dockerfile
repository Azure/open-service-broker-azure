FROM quay.io/deis/lightweight-docker-go:v0.3.0
ARG BASE_PACKAGE_NAME
ARG LDFLAGS
ENV CGO_ENABLED=0
ENV GO111MODULE=on
WORKDIR /go/src/$BASE_PACKAGE_NAME/
COPY cmd/broker cmd/broker
COPY pkg/ pkg/
COPY go.mod go.mod
COPY go.sum go.sum
COPY vendor/ vendor/
RUN go build -o bin/broker -ldflags "$LDFLAGS" -mod vendor ./cmd/broker

FROM scratch
ARG BASE_PACKAGE_NAME
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=0 /go/src/$BASE_PACKAGE_NAME/bin/broker /app/broker
CMD ["/app/broker"]
