
test:
	go test ./...

build:
	go build -o bin/crawler ./cmd/crawler/main.go

run:
	go run ./cmd/crawler/main.go