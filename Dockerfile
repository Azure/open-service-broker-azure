FROM scratch

COPY ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY bin/asb /app/asb

CMD ["/app/asb"]

EXPOSE 8080
