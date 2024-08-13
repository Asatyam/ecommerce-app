run/api:
	go run ./cmd/api -db-dsn=${GREENLIGHT_DB_DSN}

db/psql:
	@psql ${ECOMMERCE_DB_DSN}

db/migrations/up:
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${ECOMMERCE_DB_DSN}

db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}