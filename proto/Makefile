# Protocol Buffers management with buf

.PHONY: generate
generate: ## Generate Go code from proto files
	buf generate

.PHONY: lint
lint: ## Lint proto files
	buf lint

.PHONY: format
format: ## Format proto files
	buf format -w

.PHONY: breaking
breaking: ## Check for breaking changes
	buf breaking --against '.git#branch=main'

.PHONY: mod-update
mod-update: ## Update buf dependencies
	buf mod update

.PHONY: clean
clean: ## Clean generated files
	rm -rf go/

.PHONY: help
help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
