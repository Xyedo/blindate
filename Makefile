
MG_PATH=./pkg/migrations

migrate-up: 
	migrate -path $(MG_PATH) -database $(DB_DSN) up

migrate-down:
	migrate -path $(MG_PATH) -database $(DB_DSN) down

migrate-create:
	migrate create -dir $(MG_PATH) -seq -ext .sql $(name)

build-up:
	docker compose --env-file ./.env.dev up -d --build 

down: 
	docker compose --env-file ./.env.dev down 

mock-repo:
	mockgen -destination pkg/infra/repository/mock/$(mock_name).go -package mockrepo -mock_names Repository=Mock$(mock_interface) github.com/xyedo/blindate/pkg/domain/$(domain_name) Repository 

test :
	go test ./... 

test-repo:
	go test -timeout 2m github.com/xyedo/blindate/pkg/repository

.PHONY: migrate-up migrate-down migrate-create build-up up down mock-repo test test-repo
