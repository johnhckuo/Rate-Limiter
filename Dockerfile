FROM golang:1.13-alpine
WORKDIR /rate-limiter
ADD . /rate-limiter
RUN cd /rate-limiter && go build -o rate-limiter cmd/server/main.go
CMD ["./rate-limiter"]
