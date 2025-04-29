BINARY_NAME=hotel-booking-service
APP_PATH=./cmd/app
MIGRATE_PATH=./cmd/migrate

build:
	go build -o $(BINARY_NAME) $(APP_PATH)

run:
	go run $(APP_PATH)/main.go

migrate:
	go run $(MIGRATE_PATH)/main.go

fmt:
	go fmt ./...

lint:
	golangci-lint run ./...

deps:
	go mod tidy

clean:
	rm -f $(BINARY_NAME)

build-all:
	go build ./...

docker-migrate-up:
	docker-compose run --rm app migrate -path=/app/migrations -database="postgres://postgres:postgres@postgres:5432/hotel_booking?sslmode=disable" up

docker-migrate-down:
	docker-compose run --rm app migrate -path=/app/migrations -database="postgres://postgres:postgres@postgres:5432/hotel_booking?sslmode=disable" down

up:
	docker-compose up --build

down:
	docker-compose down --volumes --remove-orphans

.PHONY: build run migrate fmt lint deps clean build-all docker-migrate-up docker-migrate-down up down