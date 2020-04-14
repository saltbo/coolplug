.PHONY: all dep lint vet test test-coverage build clean

# custom define
SERVER := "cmd/server/main.go"
CLIENT := "cmd/client/main.go"

all: build

dep: ## Get the dependencies
	@go mod download

lint: ## Lint Golang files
	@golangci-lint --version
	@golangci-lint run -D errcheck

test: ## Run tests with coverage
	go test -coverprofile .coverprofile ./...
	go tool cover --func=.coverprofile

coverage-html: ## show coverage by the html
	go tool cover -html=.coverprofile

build: dep ## Build the binary file
	@go build -i -o build/bin/server $(SERVER)
	@go build -i -o build/bin/client $(CLIENT)

clean: ## Remove previous build
	@rm -f ./build

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

example: ## Build the example plugin
	@go build -v -buildmode=plugin -o=build/plugins/example.so plugin_example/main.go
