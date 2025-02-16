SERVER_OUTPUT   =sourcetool/backend
SERVER_BUILD_DIR=./cmd/server

MIGRATION_FILEPATH=postgres/ddl/`date +%Y%m%d%H%M`_$(FILENAME).sql
MIGRATION_TEMPLATE="-- +migrate Up\n\n-- +migrate Down\n"

.PHONY: dc-build
dc-build:
	@docker compose build

.PHONY: dc-up
dc-up:
	@docker compose up --remove-orphans --force-recreate

.PHONY: go-build
go-build:
	@go build -o $(SERVER_OUTPUT) $(SERVER_BUILD_DIR)


.PHONY: gen-migrate
gen-migrate:
	@touch $(MIGRATION_FILEPATH) && echo $(MIGRATION_TEMPLATE) >> ./$(MIGRATION_FILEPATH)

.PHONY: remove-docker-images
remove-docker-images:
	@bash ./script/remove-untagged-docker-images.sh

.PHONY: remove-docker-builder
remove-docker-builder:
	@docker builder prune

.PHONY: remove-postgres-data
remove-postgres-data:
	@rm -rf ./postgres/data

.PHONY: gen-encryption-key
gen-encryption-key:
	cat /dev/urandom | base64 | head -c 32

.PHONY: gen-jwt-key
gen-jwt-key:
	cat /dev/urandom | base64 | head -c 256

.PHONY: gen-swagger
gen-swagger:
	swag init -g ${SERVER_BUILD_DIR}/main.go

.PHONY: open-swagger
open-swagger:
	open http://localhost:8080/swagger/index.html

.PHONY: run-lint
run-lint:
	gofumpt -l -w . && \
		golangci-lint cache clean && \
		golangci-lint run --print-issued-lines --fix --go=1.22

.PHONY: cache-clean-lint
cache-clean-lint:
	golangci-lint cache clean
