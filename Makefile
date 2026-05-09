.PHONY: help install test lint build run clean docker-build docker-up docker-down

help:
	@echo "PPOB Microservices Makefile"
	@echo ""
	@echo "Available commands:"
	@echo "  make install         - Install dependencies for all services"
	@echo "  make test            - Run tests for all services"
	@echo "  make lint            - Run linters for all services"
	@echo "  make build           - Build all services"
	@echo "  make run             - Run all services locally"
	@echo "  make clean           - Clean build artifacts"
	@echo "  make docker-build    - Build Docker images"
	@echo "  make docker-up       - Start services with Docker Compose"
	@echo "  make docker-down     - Stop Docker Compose services"

SERVICES = auth-service user-service wallet-service transaction-service product-service integration-service

install:
	@echo "Installing dependencies..."
	@for service in $(SERVICES); do \
		echo "Installing dependencies for $$service..."; \
		cd $$service && go mod tidy && cd ..; \
	done

test:
	@echo "Running tests..."
	@for service in $(SERVICES); do \
		echo "Testing $$service..."; \
		cd $$service && go test -v -coverprofile=coverage.out ./... && cd ..; \
	done

test-coverage:
	@echo "Running tests with coverage..."
	@for service in $(SERVICES); do \
		echo "Testing $$service with coverage..."; \
		cd $$service && go test -coverprofile=coverage.out -covermode=atomic ./... && cd ..; \
	done

lint:
	@echo "Running linters..."
	@for service in $(SERVICES); do \
		echo "Linting $$service..."; \
		cd $$service && go vet ./... && cd ..; \
	done

build:
	@echo "Building all services..."
	@for service in $(SERVICES); do \
		echo "Building $$service..."; \
		cd $$service && CGO_ENABLED=0 GOOS=linux go build -o bin/$$service ./cmd && cd ..; \
	done

run:
	@echo "Starting all services..."
	@for service in $(SERVICES); do \
		echo "Starting $$service on port $$(shell echo $$service | cut -d'-' -f1 | tr '[:lower:]' '[:upper:]' | sed 's/-/ /')"; \
		cd $$service && go run cmd/main.go & cd ..; \
	done

clean:
	@echo "Cleaning build artifacts..."
	@for service in $(SERVICES); do \
		rm -rf $$service/bin; \
		rm -f $$service/coverage.out; \
	done
	rm -rf v1/*/bin

docker-build:
	@echo "Building Docker images..."
	@for service in $(SERVICES); do \
		echo "Building Docker image for $$service..."; \
		docker build -t ppob-$$service:latest -f $$service/Dockerfile .; \
	done

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

migrate-up:
	@echo "Running database migrations..."
	@docker-compose exec -T postgres psql -U postgres -d ppob -f /docker-entrypoint-initdb.d/001_initial_schema.sql

migrate-create:
	@echo "Usage: make migrate-create name=migration_name"
	@if [ -z "$(name)" ]; then echo "Please provide name"; exit 1; fi
	@echo "Creating migration $(name)..."
	@echo "-- Up Migration" > migrations/$$(date +%Y%m%d%H%M%S)_$(name).sql
	@echo "-- Down Migration" >> migrations/$$(date +%Y%m%d%H%M%S)_$(name).sql

proto-generate:
	@echo "Generating protobuf code..."
	@for service in $(SERVICES); do \
		if [ -d "proto" ]; then \
			cd proto && protoc --go_out=../$$service internal/ --go-grpc_out=../$$service internal/ && cd ..; \
		fi \
	done

tidy:
	@for service in $(SERVICES); do \
		cd $$service && go mod tidy && cd ..; \
	done

fmt:
	@for service in $(SERVICES); do \
		cd $$service && gofmt -w internal/ cmd/ && cd ..; \
	done

dev-auth:
	cd auth-service && go run cmd/main.go

dev-user:
	cd user-service && go run cmd/main.go

dev-wallet:
	cd wallet-service && go run cmd/main.go

dev-transaction:
	cd transaction-service && go run cmd/main.go

dev-product:
	cd product-service && go run cmd/main.go

dev-integration:
	cd integration-service && go run cmd/main.go