FROM golang:1.16.4

RUN mkdir -p /go/src/github.com/navidshaikh/test-webhook
WORKDIR /go/src/github.com/navidshaikh/test-webhook
COPY . .
EXPOSE 8080
CMD ["go", "run", "main.go"]
