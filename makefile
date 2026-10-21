SHELL := /bin/sh

api:
	echo "Starting server"
	go run cmd/main.go

build:
	go build -o bin/main cmd/main.go

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

migrate-up:
	@if [ -z "$$DB_URL" ]; then . ./.env; fi; migrate -path ./migrations -database "$$DB_URL" -verbose up

migrate-down:
	@if [ -z "$$DB_URL" ]; then . ./.env; fi; migrate -path ./migrations -database "$$DB_URL" down
