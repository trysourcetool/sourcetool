.PHONY: help up up-ee down down-ee build build-ee clean clean-ee logs logs-ee ps ps-ee \
	gen-keys gen-encryption-key gen-jwt-key \
	swagger swagger-open \
	backend-lint frontend-lint go-sdk-lint remove-docker-images remove-docker-builder \
	db-migrate \
	proto-generate proto-lint proto-format proto-breaking proto-mod-update proto-clean \
	go-sdk-test backend-test go-mod-tidy

# Default target
help:
	@echo "Available commands:"
	@echo ""
	@echo "Docker Compose Commands:"
	@echo "  make up              - Start the Community Edition (CE) services"
	@echo "  make up-ee           - Start the Enterprise Edition (EE) services"
	@echo "  make down            - Stop the CE services"
	@echo "  make down-ee         - Stop the EE services"
	@echo "  make build           - Build the CE services"
	@echo "  make build-ee        - Build the EE services"
	@echo "  make clean           - Stop and remove CE containers, networks, volumes"
	@echo "  make clean-ee        - Stop and remove EE containers, networks, volumes"
	@echo "  make logs            - View logs for CE services"
	@echo "  make logs-ee         - View logs for EE services"
	@echo "  make ps              - List running CE services"
	@echo "  make ps-ee           - List running EE services"
	@echo ""
	@echo "Development Commands:"
	@echo "  make gen-keys        - Generate both encryption and JWT keys"
	@echo "  make gen-encryption-key - Generate a random encryption key"
	@echo "  make gen-jwt-key     - Generate a random JWT key"
	@echo "  make swagger         - Generate Swagger documentation"
	@echo "  make swagger-open    - Open Swagger UI in browser"
	@echo "  make backend-lint    - Run linters on both CE and EE codebases (includes cache clean)"
	@echo "  make frontend-lint   - Run linters on frontend codebase"
	@echo "  make go-sdk-lint     - Run linters on Go SDK"
	@echo "  make go-sdk-test     - Run tests on Go SDK"
	@echo "  make backend-test    - Run tests on backend codebase"
	@echo "  make go-mod-tidy     - Run go mod tidy on both backend and Go SDK"
	@echo ""
	@echo "Database Commands:"
	@echo "  make db-migrate      - Run database migrations"
	@echo ""
	@echo "Protocol Buffer Commands:"
	@echo "  make proto-generate  - Generate Go code from proto files"
	@echo "  make proto-lint      - Lint proto files"
	@echo "  make proto-format    - Format proto files"
	@echo "  make proto-breaking  - Check for breaking changes"
	@echo "  make proto-mod-update - Update buf dependencies"
	@echo "  make proto-clean     - Clean generated proto files"
	@echo ""
	@echo "Maintenance Commands:"
	@echo "  make remove-docker-images - Remove untagged Docker images"
	@echo "  make remove-docker-builder - Prune Docker builder cache"

# Community Edition (CE) commands
up:
	docker compose up -d

down:
	docker compose down

build:
	docker compose build

clean:
	docker compose down -v

logs:
	docker compose logs -f

ps:
	docker compose ps

# Enterprise Edition (EE) commands
up-ee:
	docker compose -f compose.ee.yaml up -d

down-ee:
	docker compose -f compose.ee.yaml down

build-ee:
	docker compose -f compose.ee.yaml build

clean-ee:
	docker compose -f compose.ee.yaml down -v

logs-ee:
	docker compose -f compose.ee.yaml logs -f

ps-ee:
	docker compose -f compose.ee.yaml ps

# Key generation commands
gen-keys: gen-encryption-key gen-jwt-key

gen-encryption-key:
	@echo "Generating encryption key..."
	@cat /dev/urandom | base64 | head -c 32
	@echo ""

gen-jwt-key:
	@echo "Generating JWT key..."
	@cat /dev/urandom | base64 | head -c 256
	@echo ""

# Swagger commands
swagger:
	@echo "Generating Swagger documentation..."
	@cd backend && swag init -g cmd/server/main.go

swagger-open:
	@echo "Opening Swagger UI in browser..."
	@open http://localhost:8080/swagger/index.html

# Linting commands
backend-lint:
	@echo "Cleaning linter cache..."
	@cd backend && golangci-lint cache clean
	@echo "Running linters on codebase..."
	@cd backend && gofumpt -l -w . && \
		golangci-lint run --print-issued-lines --fix --go=1.22

frontend-lint:
	@echo "Running frontend linters..."
	@cd frontend && yarn lint

go-sdk-lint:
	@echo "Running Go SDK linters..."
	@cd sdk/go && gofumpt -l -w . && \
		golangci-lint run --print-issued-lines --fix --go=1.22

# Maintenance commands
remove-docker-images:
	@echo "Removing untagged Docker images..."
	@bash ./devtools/remove_untagged_docker_images.sh

remove-docker-builder:
	@echo "Pruning Docker builder cache..."
	@docker builder prune

# Database commands
db-migrate:
	@echo "Running database migrations..."
	@cd backend && go run ./devtools/cmd/db/main.go migrate

# Protocol Buffer commands
proto-generate:
	@echo "Generating Go code from proto files..."
	@cd proto && buf generate

proto-lint:
	@echo "Linting proto files..."
	@cd proto && buf lint

proto-format:
	@echo "Formatting proto files..."
	@cd proto && buf format -w

proto-breaking:
	@echo "Checking for breaking changes in proto files..."
	@cd proto && buf breaking --against '.git#branch=main'

proto-mod-update:
	@echo "Updating buf dependencies..."
	@cd proto && buf mod update

proto-clean:
	@echo "Cleaning generated proto files..."
	@cd proto && rm -rf go/

# Go SDK commands
go-sdk-test:
	@echo "Running Go SDK tests..."
	@cd sdk/go && go test -v ./...

# Backend test commands
backend-test:
	@echo "Running backend tests..."
	@cd backend && go test -v ./...

# Go module commands
go-mod-tidy:
	@echo "Running go mod tidy on backend..."
	@cd backend && go mod tidy
	@echo "Running go mod tidy on Go SDK..."
	@cd sdk/go && go mod tidy
