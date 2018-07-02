FROM quay.io/deis/lightweight-docker-go:v0.2.0
ARG BASE_PACKAGE_NAME
ARG LDFLAGS
ENV CGO_ENABLED=0
WORKDIR /go/src/$BASE_PACKAGE_NAME/
COPY cmd/broker cmd/broker
COPY pkg/ pkg/
COPY vendor/ vendor/
RUN go build -o bin/broker -ldflags "$LDFLAGS" ./cmd/broker

RUN mkdir /app && \
    cp bin/broker /app/broker
CMD ["/app/broker"]
EXPOSE 8080
