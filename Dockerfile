FROM scratch

COPY ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY bin/osba /app/osba

CMD ["/app/osba"]

EXPOSE 8080
