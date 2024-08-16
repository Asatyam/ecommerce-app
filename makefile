## help: print this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## run/api: run the application
run/api:
	@go run ./cmd/api -db-dsn=${ECOMMERCE_DB_DSN}

## db/psql: connect to the database
db/psql:
	@psql ${ECOMMERCE_DB_DSN}

## db/migrations/up: Run up migration files
db/migrations/up:
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${ECOMMERCE_DB_DSN} up

## db/migrations/new: Create a new migration file with the given name
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}