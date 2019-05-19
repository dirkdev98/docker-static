# Build app
FROM golang:latest as builder
WORKDIR /go/src/github.com/dirkdev98/docker-static
COPY . .
RUN go get -d -v

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/docker-static

# Build final image
FROM scratch
WORKDIR /
COPY --from=builder /go/bin .
ENTRYPOINT ["/docker-static"]