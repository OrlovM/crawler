
all: lint test build

lint:
	golangci-lint run

test:
	go test ./... -v

build:
	go build -o bin/crawler ./cmd/crawler/main.go

run:
	go run ./cmd/crawler/main.go
