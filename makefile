SHELL := /bin/zsh

api:
	echo "Starting server"
	go run cmd/main.go

build:
	go build -o bin/main cmd/main.go

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

migrate-up:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/daily-expense" -verbose up || echo "Migration up failed"

migrate-down:
	source .env && migrate -path ./migrations -database ${DATABASE_URL} down || echo "Migration down failed"
