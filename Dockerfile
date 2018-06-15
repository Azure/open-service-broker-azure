FROM quay.io/deis/lightweight-docker-go:v0.2.0
ARG BASE_PACKAGE_NAME
ARG LDFLAGS
ENV CGO_ENABLED=0
WORKDIR /go/src/$BASE_PACKAGE_NAME/
COPY src/open-service-broker-azure/cmd/broker cmd/broker
COPY src/open-service-broker-azure/pkg/ pkg/
COPY src/open-service-broker-azure/vendor/ vendor/
RUN go build -o bin/broker -ldflags "$LDFLAGS" ./cmd/broker

FROM scratch
ARG BASE_PACKAGE_NAME
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=0 /go/src/$BASE_PACKAGE_NAME/bin/broker /app/broker
CMD ["/app/broker"]
EXPOSE 8080
