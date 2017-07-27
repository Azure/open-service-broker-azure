FROM scratch

COPY bin/asb /app/asb

CMD ["/app/asb"]

EXPOSE 8080
