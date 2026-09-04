.PHONY: build run test lint migrate-up migrate-down docker-build docker-push

build:
	go build -o bin/app ./cmd

run:
	go run ./cmd

test:
	go test -v -race -coverprofile=coverage.out ./...

lint:
	golangci-lint run

migrate-up:
	goose -dir internal/migrations postgres "$(DB_DSN)" up

migrate-down:
	goose -dir internal/migrations postgres "$(DB_DSN)" down

docker-build:
	docker build -t account:latest .

docker-push:
	docker tag account:latest yuliapopova/account:latest
	docker push yuliapopova/account:latest
