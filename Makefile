.PHONY: all build test test-verbose clean install

BINARY_NAME=pim
BUILD_DIR=.
GO=go

# Version information
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS=-ldflags "-s -w -X github.com/hubblew/pim/cmd.Version=$(VERSION) -X github.com/hubblew/pim/cmd.Commit=$(COMMIT) -X github.com/hubblew/pim/cmd.Date=$(DATE)"

all: test build

build:
	@echo "Building $(BINARY_NAME)..."
	$(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .

test:
	@echo "Running tests..."
	$(GO) test ./...

test-verbose:
	@echo "Running tests (verbose)..."
	$(GO) test -v ./...

clean:
	@echo "Cleaning build artifacts..."
	rm -f $(BUILD_DIR)/$(BINARY_NAME)
	rm -f coverage.out coverage.html

install:
	@echo "Installing $(BINARY_NAME)..."
	$(GO) install $(LDFLAGS) .
