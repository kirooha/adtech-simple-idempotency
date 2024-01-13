PSQL_DSN:=$(shell grep -A1 'pg_dsn' config/config.json | sed "s/pg_dsn//" | sed "s/[:}]//" | sed "s/[ \"]//g")
ROOT_PSQL_DSN:=$(shell echo $(PSQL_DSN) | sed 's/adtech_simple/postgres/g')

bin-deps:
	go install github.com/pressly/goose/v3/cmd/goose@v3.15.1
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.23.0
	go install golang.org/x/tools/cmd/stringer@v0.14.0

regenerate-db: regenerate-db-drop-create db-migrate db-gen-structure

regenerate-db-drop-create:
	psql "$(ROOT_PSQL_DSN)" --command="DROP DATABASE IF EXISTS adtech_simple;"
	psql "$(ROOT_PSQL_DSN)" --command="CREATE DATABASE adtech_simple WITH OWNER postgres ENCODING = 'UTF8';"

db-migrate:
	goose -dir db/migrations postgres "$(PSQL_DSN)" up

db-gen-structure:
	pg_dump "$(PSQL_DSN)" --schema-only --no-owner --no-privileges --no-tablespaces --no-security-labels --no-comments > db/structure.sql

generate: tidy sqlc stringer

sqlc:
	sqlc generate

stringer:
	cd ./internal/pkg/model;stringer -type=JobType,QueueType

tidy:
	go mod tidy

db-dsn:
	echo $(PSQL_DSN)