include .env

MG_PATH=./migrations
DB_DSN =postgresql://${PG_USER}:${PG_PASSWORD}@${PG_HOST}:${PG_PORT}/${PG_DB}?sslmode=disable

.PHONY: migrate-up 
migrate-up: 
	migrate -path $(MG_PATH) -database $(DB_DSN) up

.PHONY: migrate-down 
migrate-down:
	migrate -path $(MG_PATH) -database $(DB_DSN) down

.PHONY: migrate-create 
migrate-create:
	migrate create -dir $(MG_PATH) -seq -ext .sql $(name)

.PHONY: up
up:
	docker compose --env-file ./.env up -d
	migrate -path $(MG_PATH) -database $(DB_DSN) up 

.PHONY: migrate-force
migrate-force:
	migrate -path $(MG_PATH) -database $(DB_DSN) force $(n)

.PHONY: build
build:
	docker build --tag blindate .

.PHONY: down
down: 
	docker compose --env-file ./.env down 

.PHONY: mock-repo
mock-repo:
	mockgen -destination pkg/infra/repository/mock/$(mock_name).go -package mockrepo -mock_names Repository=Mock$(mock_interface) github.com/xyedo/blindate/pkg/domain/$(domain_name) Repository 

.PHONY: test
test :
	go test ./... 

.PHONY: test-repo
test-repo:
	go test -timeout 2m github.com/xyedo/blindate/pkg/repository

.PHONY: tunnel
tunnel:
	ngrok http http://${APP_HOST}:${APP_PORT}
	
