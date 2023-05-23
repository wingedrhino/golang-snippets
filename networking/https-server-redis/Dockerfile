FROM golang:1.10.3
RUN mkdir -p /go/src/github.com/wingedrhino/golang-snippets/networking/https-server-redis
WORKDIR /go/src/github.com/wingedrhino/golang-snippets/networking/https-server-redis
COPY . .
RUN go install .
EXPOSE 8443
VOLUME ["/crt"]
ENTRYPOINT ["https-server-redis", "-cert", "/crt/cert.pem", "-key", "/crt/key.pem"]

