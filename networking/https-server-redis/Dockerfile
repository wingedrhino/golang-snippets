FROM golang:1.9
RUN mkdir -p /go/src/bitbucket.org/mahtuag/learn-go/networking/https-server-redis
WORKDIR /go/src/bitbucket.org/mahtuag/learn-go/networking/https-server-redis
COPY . .
RUN go install .
EXPOSE 8443
VOLUME ["/crt"]
ENTRYPOINT ["https-server-redis", "-cert", "/crt/cert.pem", "-key", "/crt/key.pem"]

