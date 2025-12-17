# Vault-inator Makefile
# Run the backend service first, then the frontend

.PHONY: help run run-backend run-frontend build-backend build-frontend clean stop

# Load environment variables from .env file
ifneq (,$(wildcard .env))
    include .env
    export
endif

# Default target
help:
	@echo "Available targets:"
	@echo "  run          - Start backend and frontend (backend first)"
	@echo "  run-backend  - Start only the backend service"
	@echo "  run-frontend - Start only the frontend (requires backend running)"
	@echo "  build-backend - Build the Go backend"
	@echo "  build-frontend - Build the React frontend"
	@echo "  clean        - Stop all services and clean up"
	@echo "  stop         - Stop all running services"

# Main target: run backend first, then frontend
run: run-backend
	@echo "Backend started. Starting frontend in 3 seconds..."
	@sleep 3
	@$(MAKE) run-frontend

# Run backend service
run-backend:
	@echo "Starting backend service..."
	@cd cmd/vault-inator && HOST=localhost go run main.go &
	@echo "Backend service started on http://localhost:8080"

# Run frontend (requires backend to be running)
run-frontend:
	@echo "Starting frontend..."
	@cd web && HOST=localhost npm start
	@echo "Frontend started on http://localhost:3000"

# Build backend
build-backend:
	@echo "Building backend..."
	@cd cmd/vault-inator && go build -o main main.go
	@echo "Backend built successfully"

# Build frontend
build-frontend:
	@echo "Building frontend..."
	@cd web && npm run build
	@echo "Frontend built successfully"

# Build both
build: build-backend build-frontend
	@echo "Both backend and frontend built successfully"

# Stop all services
stop:
	@echo "Stopping all services..."
	@pkill -f "go run main.go" || true
	@pkill -f "npm start" || true
	@echo "All services stopped"

# Clean up
clean: stop
	@echo "Cleaning up..."
	@rm -f cmd/vault-inator/main
	@rm -rf web/build
	@echo "Cleanup completed"

# Development target with hot reload
dev: run-backend
	@echo "Backend started. Starting frontend in development mode..."
	@sleep 3
	@cd web && HOST=localhost npm start

# Production target
prod: build
	@echo "Starting production services..."
	@cd cmd/vault-inator && HOST=localhost ./main &
	@echo "Backend started. Frontend build available in web/build/"
