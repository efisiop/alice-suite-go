# Alice Suite Development Makefile

.PHONY: help setup start stop build test clean proxy-setup

help: ## Show this help message
	@echo "Alice Suite Development Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

setup: proxy-setup ## Run initial setup (proxy configuration)
	@echo "✅ Setup complete!"

proxy-setup: ## Configure Safari proxy bypass for localhost
	@./setup_safari_proxy.sh

build: ## Build the server
	@echo "Building server..."
	@go build -o bin/server ./cmd/server
	@echo "✅ Build complete: ./bin/server"

start: build ## Start the development server
	@echo "Starting development server..."
	@./start_dev_server.sh

run: start ## Alias for start

stop: ## Stop the server if running
	@echo "Stopping server..."
	@-pkill -f "bin/server" || true
	@-lsof -ti:8080 | xargs kill -9 2>/dev/null || true
	@echo "✅ Server stopped"

restart: stop start ## Restart the server

test: ## Run tests
	@go test ./...

clean: ## Clean build artifacts
	@rm -rf bin/
	@go clean
	@echo "✅ Cleaned"

check: ## Check if server is running
	@if lsof -Pi :8080 -sTCP:LISTEN -t >/dev/null 2>&1; then \
		echo "✅ Server is running on port 8080"; \
		echo "   Access at: http://127.0.0.1:8080/reader/login"; \
	else \
		echo "❌ Server is NOT running"; \
		echo "   Start with: make start"; \
	fi

