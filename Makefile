.PHONY: up up-ee down down-ee build build-ee clean clean-ee logs logs-ee ps ps-ee help

# Default target
help:
	@echo "Available commands:"
	@echo "  make up        - Start the Community Edition (CE) services"
	@echo "  make up-ee     - Start the Enterprise Edition (EE) services"
	@echo "  make down      - Stop the CE services"
	@echo "  make down-ee   - Stop the EE services"
	@echo "  make build     - Build the CE services"
	@echo "  make build-ee  - Build the EE services"
	@echo "  make clean     - Stop and remove CE containers, networks, volumes"
	@echo "  make clean-ee  - Stop and remove EE containers, networks, volumes"
	@echo "  make logs      - View logs for CE services"
	@echo "  make logs-ee   - View logs for EE services"
	@echo "  make ps        - List running CE services"
	@echo "  make ps-ee     - List running EE services"

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
