.PHONY: build run test

build:
	go build -o app ./cmd/

run:
	docker-compose up -d

tests:
	go test ./...
