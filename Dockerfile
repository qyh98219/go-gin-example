FROM golang:latest

ENV GOPROXY https://goproxy.cn, direct
WORKDIR $GOPATH/src/github.com/go-gin-example
COPY . $GOPAHT/src/github.com/go-gin-example
RUN go build .

EXPOSE 8000
ENTRYPOINT ["./go-gin-example"]