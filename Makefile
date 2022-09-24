PATH=./pkg/service/postgre/migrations

migrate-up: 
	@migrate -path ${PATH} -database ${DSN} up	
