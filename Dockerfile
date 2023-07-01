# Builder image to add TLS certificates. The scratch image doesn't ship
# with any, so outbound TLS connections will break without these.
FROM alpine:latest AS builder
RUN apk add -U make ca-certificates

# Final image with binary built by GoReleaser.
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY labrador /
ENTRYPOINT ["/labrador"]
