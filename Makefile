.DEFAULT_GOAL := help

.PHONY: build
build: ## Build the binary
	go build -o promptctl main.go

.PHONY: vet
vet: ## Build the binary
	go vet ./...

.PHONY: test
test: ## Test all the test files recursively
	go -v ./tests/... -coverpkg=./...

.PHONY: --help
--help: ##
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}'

.PHONY: help
help: --help
