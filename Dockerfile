# Builder image to add TLS certificates.
FROM alpine:latest AS builder
RUN apk add -U make ca-certificates

# Final image with binary built by GoReleaser.
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY labrador /
ENTRYPOINT ["/labrador"]
