# Justfile for btui project

# Default recipe - show available commands
default:
    @just --list

# Build the application
build:
    mkdir -p dist
    go build -o dist/btui .

# Run the application with list-devices command
run:
    go run . list-devices

# Run tests
test:
    go test ./...

# Format code
fmt:
    go fmt ./...

# Vet code for issues
vet:
    go vet ./...

# Tidy dependencies
tidy:
    go mod tidy

# Clean build artifacts
clean:
    rm -rf dist/

# Build and run
build-run: build
    ./dist/btui list-devices

# Run all checks (fmt, vet, test)
check: fmt vet test

# Full pipeline: check, build
all: check build