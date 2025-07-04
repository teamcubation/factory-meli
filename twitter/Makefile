.DEFAULT_GOAL := run

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# Handle PGHOST for development vs production
ifeq ($(ENV),production)
    PGHOST_OVERRIDE=$(PGHOST)
else
    PGHOST_OVERRIDE=localhost
endif

PG_CONN_STRING=postgresql://$(PGUSER):$(PGPASSWORD)@$(PGHOST_OVERRIDE):$(PGPORT)/$(PGDATABASE)?sslmode=disable

# General tooling
vet: 
	go vet ./... 
	go fmt ./...
.PHONY: vet fmt

vendor:
	go mod vendor
.PHONY: vendor

# Adjust build path to reflect new main.go location
build: vendor
	go build -mod=vendor -o ./cmd/http/twitter ./cmd/http
.PHONY: build

# Run using Air
run:
	air
.PHONY: run

test:
	go test ./... -v -count=1
.PHONY: test

# Goose migrations
m-create:
ifndef name
	$(error name is required, e.g., `make m-create name=article_alter_table_add_column_content`)
endif
	goose create -s -dir core/ports/output/postgres/migrations $(name) sql
.PHONY: m-create

m-up:
	goose -dir core/ports/output/postgres/migrations postgres "$(PG_CONN_STRING)" up
.PHONY: m-up

m-down:
	goose -dir core/ports/output/postgres/migrations postgres "$(PG_CONN_STRING)" down
.PHONY: m-down

m-status:
	goose -dir core/ports/output/postgres/migrations postgres "$(PG_CONN_STRING)" status
.PHONY: m-status
