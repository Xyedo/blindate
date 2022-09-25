cc=gcc

PATH=./pkg/repository/migrations
DSN=postgres://blindate:pa55word@localhost:5432/blindate?sslmode=disable

migrate-up: 
	migrate -path $(PATH) -database $(DSN) up

migrate-down:
	migrate -path $(PATH) -database $(DSN) down

migrate-create:
	migrate create -dir $(PATH) -seq -ext .sql $()

build-up:
	docker compose up -d --build
up:
	docker compose up

down: 
	docker compose down


.PHONY: migrate-up migrate-down migrate-create build-up up down
