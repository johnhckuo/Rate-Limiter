docker:
	docker build -f ./Dockerfile -t="johnhckuo/rate-limiter" .

test-api:
	go test -v ./test/ -count 1

run:
	go run cmd/server/main.go

lint:
	golint ./...

fmt:
	go fmt ./...
