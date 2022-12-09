
MG_PATH=./pkg/repository/migrations
REPO_PATH=github.com/xyedo/blindate/pkg/repository

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
	mockgen -destination pkg/repository/mock/$(mock_name) -package mockrepo $(REPO_PATH) $(interface) 

test :
	go test ./... 

test-repo:
	go test -timeout 2m github.com/xyedo/blindate/pkg/repository

.PHONY: migrate-up migrate-down migrate-create build-up up down mock-repo test test-repo
