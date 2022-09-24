PATH=./pkg/repository/postgre/migrations

migrate-up: 
	@migrate -path ${PATH} -database ${DSN} up

migrate-down:
	@migrate -path ${PATH} -database ${DSN} down

migrate-create:
	@migrate create -dir ${PATH} -seq -ext .sql $()

up-build:
	@docker-compose up -d --build

up:
	@docker-compose up

down: 
	@docker-compose down


.PHONY: migrate-up migrate-down migrate-create up-build up down
