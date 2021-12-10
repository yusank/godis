MSG=$(msg)
help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[0-9A-Za-z_-]+:.*?##/ { printf "  \033[36m%-45s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
.PHONY: server-run
server-run: ## run server as default mode
	CGO_ENABLED=0 go run cmd/server/main.go
.PHONY: build-linux
build-linux: ## build server binary for linux
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o go_build_server cmd/server/main.go
.PHONY: lint
lint: ## run golangci-lint for project
	golangci-lint run ./... -v
.PHONY: test
test: ## run all test cases
	CGO_ENABLED=0 go test -v ./...
.PHONEY: cmt
cmt:## git commit with message
ifeq ($(strip $(MSG)),)
	@echo "must input commit msg"
	exit 1
endif
	git add .
	git commit -m '$(MSG)'
	@echo "msg:$(MSG)"
.PHONEY: gen_cmd
gen_cmd: ## gen redis cmd code
	go build -o go_build_gen_redis_cmd cmd/gen_redis_cmd/main.go
	./go_build_gen_redis_cmd -d ./redis
.PHONEY: clean
clean: ## clean all generated code
	rm -rf redis/*.cmd.go