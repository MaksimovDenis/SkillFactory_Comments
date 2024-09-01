include .env
build:
	go mod download && CGD_ENABLED=0 GOOS=linux go build -o ./.bin/app ./cmd/main.go

up:
	docker-compose up --build 

down:
	docker-compose down

run:
	go run ./cmd/main.go