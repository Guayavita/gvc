# Makefile for Guayavita (gvc)

# Detect module path from go.mod so ldflags package paths are correct
MODULE := $(shell awk '/^module /{print $$2}' go.mod)
COMMONS_PKG := $(MODULE)/internal/commons
BINARY ?= guayavita
BIN_DIR ?= bin
OUT := $(BIN_DIR)/$(BINARY)

# Git metadata (fall back to sensible defaults when not in a git repo)
DIRTY := $(shell test -n "$$(git status --porcelain 2>/dev/null)" && echo "-dirty" || true)
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)$(DIRTY)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo unknown)
# Version must be SemVer: use latest tag if present; otherwise fallback to 0.0.0
VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || echo 0.0.0)
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS := \
  -X '$(COMMONS_PKG).Version=$(VERSION)' \
  -X '$(COMMONS_PKG).Build=$(BRANCH)' \
  -X '$(COMMONS_PKG).GitCommit=$(COMMIT)'

.PHONY: all build clean run print-vars install

all: build

$(BIN_DIR):
	@mkdir -p $(BIN_DIR)

build: $(BIN_DIR)
	@echo "Building $(OUT)"
	GOFLAGS= CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(OUT) .
	@echo "Built $(OUT)"

run:
	@echo "Running with injected build metadata"
	GOFLAGS= CGO_ENABLED=0 go run -ldflags "$(LDFLAGS)" . -- version || true

install:
	@echo "Installing $(BINARY) to GOPATH/bin"
	GOFLAGS= CGO_ENABLED=0 go install -ldflags "$(LDFLAGS)" .

print-vars:
	@echo MODULE: $(MODULE)
	@echo COMMONS_PKG: $(COMMONS_PKG)
	@echo VERSION: $(VERSION)
	@echo BRANCH: $(BRANCH)
	@echo COMMIT: $(COMMIT)
	@echo DATE: $(DATE)
	@echo LDFLAGS: $(LDFLAGS)

clean:
	rm -rf $(BIN_DIR)
