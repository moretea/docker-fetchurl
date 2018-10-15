FROM golang:alpine as builder

# Add certificates to ensure that we can fetch https files.
RUN apk add ca-certificates git
RUN go get -u golang.org/x/vgo
WORKDIR /go/src/github.com/moretea/docker-fetchurl

ENV GO111MODULE=on

# Populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN GOOS=linux GARCH=amd64 CGO_ENABLED=0 go install -a -installsuffix cgo ./cmd/fetchurl

FROM scratch
ADD ./tmp /tmp
COPY --from=builder /go/bin/fetchurl /bin/fetchurl
COPY --from=builder /etc/ssl/certs /etc/ssl/certs

ENTRYPOINT ["/bin/fetchurl"]
