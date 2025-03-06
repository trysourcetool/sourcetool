.PHONY: dc-up
dc-up:
	@docker compose up --remove-orphans --force-recreate

.PHONY: remove-docker-images
remove-docker-images:
	@bash ./script/remove-untagged-docker-images.sh

.PHONY: remove-docker-builder
remove-docker-builder:
	@docker builder prune

.PHONY: gen-encryption-key
gen-encryption-key:
	cat /dev/urandom | base64 | head -c 32

.PHONY: gen-jwt-key
gen-jwt-key:
	cat /dev/urandom | base64 | head -c 256

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
