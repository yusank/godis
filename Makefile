help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[0-9A-Za-z_-]+:.*?##/ { printf "  \033[36m%-45s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
.PHONY: server-run
server-run: ## run server as defualt mode
	go run cmd/server/main.go
.PHONY: lint
lint: ## run golangci-lint for project
	golangci-lint run ./...
.PHONY: test
test: ## run all test cases
	go test -v ./...