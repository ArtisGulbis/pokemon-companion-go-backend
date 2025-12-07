# Run all tests
test:
	go test ./... -v

# Run tests with coverage
coverage:
	go test ./... -coverprofile=coverage.out

coverage-html:
	go tool cover -html=coverage.out

# Build the sync command
build:
	go build -o bin/sync ./cmd/sync

# Run the sync command
run: build
	./bin/sync

# Run sync with limit
sync:
	go run ./cmd/sync -limit=5

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out
	rm -f pokemon.db

# Install dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Run linter (requires golangci-lint)
lint:
	golangci-lint run

# Run all checks (test, fmt, lint)
check: fmt test lint

# Show available commands
help:
	@echo "Available targets:"
	@echo "  make test      - Run all tests"
	@echo "  make coverage  - Run tests with coverage report"
	@echo "  make build     - Build the sync binary"
	@echo "  make run       - Build and run sync"
	@echo "  make sync      - Run sync directly (no build)"
	@echo "  make clean     - Remove build artifacts"
	@echo "  make fmt       - Format code"
	@echo "  make lint      - Run linter"
	@echo "  make check     - Run fmt + test + lint"
	@echo "  make help      - Show this help"