PROJECT_NAME := $(shell basename "$(PWD)")
PKG := "github.com/kott/$(PROJECT_NAME)"
PKG_LIST := $(shell go list ./...)
GO_FILES := $(shell find . -name '*.go' | grep -v _test.go)
BUILD_DIR := "bin"

.PHONY: all dep build clean test lint fmt

all: build

lint: ## Lint the files
	@golint -set_exit_status ./...

fmt: ## Format the files
	@go fmt ./...

vet: ## Vet the files
	@go vet ./...

test: ## Run unittests
	@go test -short ${PKG_LIST}

dep: ## Get the dependencies
	@go get -v -d ./...

build: dep ## Build the binary file
	@go build -o $(BUILD_DIR)/$(PROJECT_NAME) -v ./cmd/${PROJECT_NAME}

clean: ## Remove previous build
	@rm -f $(BUILD_DIR)/*

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
