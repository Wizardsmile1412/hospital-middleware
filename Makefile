.PHONY: run build test docker-up docker-down migrate-up migrate-down

run:
	go run ./cmd/main.go

build:
	go build -o bin/server ./cmd/main.go

test:
	go test -v ./...

docker-up:
	docker compose up --build

docker-down:
	docker compose down

migrate-up:
	migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" down
