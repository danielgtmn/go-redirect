.PHONY: build test clean help

help:
	@echo "Available commands:"
	@echo "  make build    - Build Docker image"
	@echo "  make test     - Run all tests"
	@echo "  make clean    - Clean up containers and images"

build:
	@echo "Building Docker image..."
	docker build -t go-redirect:test .

test: build
	@echo "Running tests..."
	@./test.sh

clean:
	@echo "Cleaning up..."
	docker stop redirect-test 2>/dev/null || true
	docker rm redirect-test 2>/dev/null || true
	docker rmi go-redirect:test 2>/dev/null || true
