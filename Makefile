
MG_PATH=./pkg/repository/migrations
DB_DSN=postgres://blindate:pa55word@localhost:5433/blindate?sslmode=disable
REPO_PATH=github.com/xyedo/blindate/pkg/repository

migrate-up: 
	migrate -path $(MG_PATH) -database $(DB_DSN) up

migrate-down:
	migrate -path $(MG_PATH) -database $(DB_DSN) down

migrate-create:
	migrate create -dir $(MG_PATH) -seq -ext .sql $(name)

build-up:
	docker compose up -d --build
up:
	docker compose up -d

down: 
	docker compose down

mock-repo:
	mockgen -destination pkg/repository/mock/$(mock_name) -package mockrepo $(REPO_PATH) $(interface) 

test :
	go test ./...

test-repo:
	go test -timeout 2m -coverprofile=C:\Users\ACER\AppData\Local\Temp\vscode-googaDrR\go-code-cover github.com/xyedo/blindate/pkg/repository

.PHONY: migrate-up migrate-down migrate-create build-up up down mock-repo test test-repository
