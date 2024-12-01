.PHONY: generate test lint clean buf-generate

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Binary names
BINARY_NAME=protoc-gen-mcp

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) cmd/protoc-gen-mcp/main.go

test:
	$(GOTEST) -v ./...

lint:
	golangci-lint run
	buf lint

generate: build
	rm -rf examples/**/gen
	cd examples/basic/proto && buf generate

clean:
	rm -f bin/$(BINARY_NAME)
	rm -rf examples/**/gen
	rm -f protoc-gen-mcp

install:
	$(GOBUILD) -o $(GOPATH)/bin/$(BINARY_NAME) cmd/protoc-gen-mcp/main.go

# Buf-specific commands
buf-lint:
	buf lint

buf-generate:
	buf generate

# Development helpers
dev-deps:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/bufbuild/buf/cmd/buf@latest
