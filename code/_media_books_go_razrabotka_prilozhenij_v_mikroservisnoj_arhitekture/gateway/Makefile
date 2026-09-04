.PHONY: build run test lint docker-build docker-push

build:
	go build -o bin/app ./cmd

run:
	go run ./cmd

test:
	go test -v -race ./...

lint:
	golangci-lint run

docker-build:
	docker build -t gateway:latest .

docker-push:
	docker tag gateway:latest yuliapopova/gateway:latest
	docker push yuliapopova/gateway:latest
