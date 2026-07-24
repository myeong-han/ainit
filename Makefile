.PHONY: all build run test clean install fmt tidy help

BINARY_NAME=bin/ainit
MAIN_PACKAGE=./cmd/ainit

all: build

## build: Build the ainit TUI binary
build:
	@echo "🔨 Building $(BINARY_NAME)..."
	@mkdir -p bin
	go build -o $(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "✅ Build complete: $(BINARY_NAME)"

## run: Build and run the TUI tool
run: build
	./$(BINARY_NAME)

## test: Run unit tests
test:
	@echo "🧪 Running unit tests..."
	go test -v ./...

## fmt: Format Go source code
fmt:
	@echo "🎨 Formatting code..."
	go fmt ./...

## tidy: Tidy Go module dependencies
tidy:
	@echo "🧹 Tidying Go dependencies..."
	go mod tidy

## install: Install binary to $GOPATH/bin
install:
	@echo "📦 Installing ainit binary..."
	go install $(MAIN_PACKAGE)

## clean: Remove built binaries
clean:
	@echo "🗑️ Cleaning bin/..."
	rm -rf bin

## help: Show Makefile targets
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':'
