FROM debian:stretch

RUN adduser --system \
  --shell /bin/bash \
  --disabled-password \
  --home /var/lib/azure-service-broker \
  --group \
  asb

COPY bin/asb /opt/azure-service-broker/bin/asb

# Fix some permissions since we'll be running as a non-root user
RUN chown -R asb:asb /opt/azure-service-broker

USER asb

WORKDIR /var/lib/azure-service-broker

CMD ["/opt/azure-service-broker/bin/asb"]
