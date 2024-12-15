api:
	echo "Starting server"
	go run cmd/main.go

build:
	go build -o bin/main cmd/main.go

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

migrate-up:
	migrate -path ./migrations -database ${DATABASE_URL} up

migrate-down:
	migrate -path ./migrations -database ${DATABASE_URL} down
