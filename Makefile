include .env

MG_PATH=./migrations
DB_DSN =postgresql://${PG_USER}:${PG_PASSWORD}@localhost:5432/${PG_DB}?sslmode=disable

migrate-up: 
	migrate -path $(MG_PATH) -database $(DB_DSN) up

migrate-down:
	migrate -path $(MG_PATH) -database $(DB_DSN) down

migrate-create:
	migrate create -dir $(MG_PATH) -seq -ext .sql $(name)

up:
	docker compose --env-file ./.env up -d

migrate-force:
	migrate -path $(MG_PATH) -database $(DB_DSN) force $(n)

build:
	docker build --tag blindate .
	
down: 
	docker compose --env-file ./.env down 

mock-repo:
	mockgen -destination pkg/infra/repository/mock/$(mock_name).go -package mockrepo -mock_names Repository=Mock$(mock_interface) github.com/xyedo/blindate/pkg/domain/$(domain_name) Repository 

test :
	go test ./... 

test-repo:
	go test -timeout 2m github.com/xyedo/blindate/pkg/repository

.PHONY: migrate-up migrate-down migrate-create up down mock-repo test test-repo build
