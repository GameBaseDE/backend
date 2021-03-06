PROJECT_NAME := "gamebase-daemon"
PKG := "gitlab.tandashi.de/GameBase/$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

.PHONY: all dep build clean test coverage coverhtml lint

all: build

lint: ## Lint the files
	@golint -set_exit_status ${PKG_LIST}

test: ## Run unittests
	@go test -short ${PKG_LIST}

race: ## Run data race detector
	@go test -race -short ${PKG_LIST}

msan: ## Run memory sanitizer
	@go test -msan -short ${PKG_LIST}

coverage: ## Generate global code coverage report
	@chmod +x tools/coverage.sh;
	bash tools/coverage.sh;

coverhtml: ## Generate global code coverage report in HTML
	@chmod +x tools/coverage.sh;
	bash tools/coverage.sh html;

dep: ## Get the dependencies
	@go get -v -d ./...

build: ## Build the binary file
	@go build -i -v -o out/server

build-static: ## Build the binary file and link it statically
	@go build -i -v -o out/server --ldflags '-linkmode external -extldflags "-static"'

clean: ## Remove previous build
	@rm -f $(PROJECT_NAME)

generate: ## Generate the server stub from the latest openapi specification
	dos2unix tools/generate_openapi.sh | true
	@chmod +x tools/coverage.sh
	bash tools/generate_openapi.sh

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
