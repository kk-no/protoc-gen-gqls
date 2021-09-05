.DEFAULT_GOAL:=generate

BIN := $(CURDIR)/.bin
GEN := $(CURDIR)/gen
PATH := $(abspath $(BIN)):$(PATH)
UNAME_OS := $(shell uname -s)
UNAME_ARCH := $(shell uname -m)

$(BIN):
	mkdir -p $(BIN)

$(GEN):
	mkdir -p $(GEN)

# golangci-lint setting
GOLANGCLI_LINT := $(BIN)/golangci-lint
GOLANGCLI_LINT_VERSION := v1.42.0
$(GOLANGCLI_LINT): | $(BIN)
	@curl -sSfL "https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh" | sh -s -- -b $(BIN) $(GOLANGCLI_LINT_VERSION)
	@chmod +x "$(BIN)/golangci-lint"

PROTOC_GEN_GQLS := $(BIN)/protoc-gen-gqls
$(PROTOC_GEN_GQLS): | $(BIN)
	@go build -o $(BIN)/protoc-gen-gqls main.go

.PHONY: run
run: ## Run applications
	go run main.go

.PHONY: mod
mod: ## Download modules
	go mod tidy

.PHONY: test
test: ## Run unit test
	go test ./...

.PHONY: generate
generate: $(GEN) $(PROTOC_GEN_GQLS) ## Run generate plugin
	protoc -I./example/proto:./example/proto/third_party --plugin=$(PROTOC_GEN_GQLS) --gqls_out=$(GEN) example/proto/*.proto

.PHONY: regen
regen: ## Run generate plugin
	$(MAKE) clean
	$(MAKE)

.PHONY: lint
lint: | $(GOLANGCLI_LINT) ## Run linter
	$(BIN)/golangci-lint -verbose run ./...

.PHONY: clean
clean: ## Remove .bin and gen directory
	rm -rf "$(BIN)"
	rm -rf "$(GEN)"

.PHONY: os
os:  ## Print OS name
	@echo "$(UNAME_OS)"

.PHONY: help
help: ## Print help
	@grep -E '^[/a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'