.PHONY: all build test test-verbose clean install

BINARY_NAME=pim
BUILD_DIR=.
GO=go

all: test build

build:
	@echo "Building $(BINARY_NAME)..."
	$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME) .

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
	$(GO) install .
