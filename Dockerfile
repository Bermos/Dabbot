FROM scratch

COPY ca-certificates.crt /etc/ssl/certs/
COPY main /
COPY dabs /dabs

CMD ["/main"]