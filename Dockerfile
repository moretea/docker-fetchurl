FROM golang:alpine AS builder

# Add certificates to ensure that we can fetch https files.
RUN apk add ca-certificates
COPY ./fetchurl.go .
RUN GOOS=linux GARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo ./fetchurl.go

FROM scratch
ADD ./tmp /tmp
COPY --from=builder /go/fetchurl /bin/fetchurl
COPY --from=builder /etc/ssl/certs /etc/ssl/certs

ENTRYPOINT ["/bin/fetchurl", "-template", "-url"]
