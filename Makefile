
MG_PATH=./pkg/repository/migrations
DB_DSN=postgres://blindate:pa55word@localhost:5433/blindate?sslmode=disable
REPO_PATH=github.com/xyedo/blindate/pkg/repository



migrate-up: 
	migrate -path $(MG_PATH) -database $(DB_DSN) up

migrate-down:
	migrate -path $(MG_PATH) -database $(DB_DSN) down

migrate-create:
	migrate create -dir $(MG_PATH) -seq -ext .sql $()

build-up:
	docker compose up -d --build
up:
	docker compose up -d

down: 
	docker compose down

mock:
	mockgen -destination pkg/repository/mock/$(mock_name) -package mockrepo $(REPO_PATH) $(interface) 

test :
	go test ./...

.PHONY: migrate-up migrate-down migrate-create build-up up down mock test
