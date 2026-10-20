SHELL := /bin/sh

api:
	echo "Starting server"
	go run cmd/main.go

build:
	go build -o bin/main cmd/main.go

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

migrate-up:
	. ./.env && migrate -path ./migrations -database "$$DB_URL" -verbose up || echo "Migration up failed"

migrate-down:
	. ./.env && migrate -path ./migrations -database "$$DB_URL" down || echo "Migration down failed"
