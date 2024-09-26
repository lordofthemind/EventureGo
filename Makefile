# Environment variables for MongoDB
MONGODB_CONTAINER_NAME ?= MegaMongoContainer
MONGODB_IMAGE_TAG ?= latest
MONGODB_DB_NAME ?= EventureGo
MONGODB_PORT ?= 27017

# Environment variables for PostgreSQL
PG_CONTAINER_NAME ?= MegaPostgresContainer
PG_IMAGE_TAG ?= latest
PG_DB_NAME ?= EventureGo
PG_DB_USERNAME ?= postgres
PG_DB_PASSWORD ?= EventureGoSecret
PG_PORT ?= 5432

# Colors for help command
CYAN := \033[36m
RESET := \033[0m

# Go-related commands
build: ## Build the Go project
	@echo "Building Go project..."
	go build -o bin/EventureGo main.go

run: ## Run the Go project
	@echo "Running Go project..."
	go run main.go

test: ## Run all tests
	@echo "Running tests..."
	go test ./...

lint: ## Run golangci-lint to check for linting issues
	@echo "Running linter..."
	golangci-lint run ./...

fmt: ## Format the Go code
	@echo "Formatting Go code..."
	go fmt ./...

clean: ## Remove binaries and other build artifacts
	@echo "Cleaning up..."
	rm -rf bin/

# Docker MongoDB commands
crtmgcnt: ## Create and start the MongoDB container
	@echo "Creating and starting MongoDB container..."
	docker run --name $(MONGODB_CONTAINER_NAME) -p $(MONGODB_PORT):27017 -d mongo:$(MONGODB_IMAGE_TAG)

strmgcnt: ## Start the MongoDB container
	@echo "Starting MongoDB container..."
	docker start $(MONGODB_CONTAINER_NAME)

stpmgcnt: ## Stop the MongoDB container
	@echo "Stopping MongoDB container..."
	docker stop $(MONGODB_CONTAINER_NAME)

rmvmgcnt: ## Remove the MongoDB container
	@echo "Removing MongoDB container..."
	docker rm $(MONGODB_CONTAINER_NAME)

# MongoDB database commands
crtmgdb: strmgcnt ## Create MongoDB database
	@echo "Creating MongoDB database..."
	docker exec -it $(MONGODB_CONTAINER_NAME) mongosh --eval "use $(MONGODB_DB_NAME)"

drpmgdb: strmgcnt ## Drop MongoDB database
	@echo "Dropping MongoDB database..."
	docker exec -it $(MONGODB_CONTAINER_NAME) mongosh --eval "db.getSiblingDB('$(MONGODB_DB_NAME)').dropDatabase()"

# Docker PostgreSQL commands
crtpgcnt: ## Create and start the PostgreSQL container
	@echo "Creating and starting PostgreSQL container..."
	docker run --name $(PG_CONTAINER_NAME) -p $(PG_PORT):5432 -e POSTGRES_DB=$(PG_DB_NAME) -e POSTGRES_USER=$(PG_DB_USERNAME) -e POSTGRES_PASSWORD=$(PG_DB_PASSWORD) -d postgres:$(PG_IMAGE_TAG)

strpgcnt: ## Start the PostgreSQL container
	@echo "Starting PostgreSQL container..."
	docker start $(PG_CONTAINER_NAME)

stppgcnt: ## Stop the PostgreSQL container
	@echo "Stopping PostgreSQL container..."
	docker stop $(PG_CONTAINER_NAME)

rmvpgcnt: ## Remove the PostgreSQL container
	@echo "Removing PostgreSQL container..."
	docker rm $(PG_CONTAINER_NAME)

# PostgreSQL database commands
crtpgdb: strpgcnt ## Create PostgreSQL database
	@echo "Creating PostgreSQL database..."
	docker exec -it $(PG_CONTAINER_NAME) createdb -U $(PG_DB_USERNAME) $(PG_DB_NAME)

drppgdb: strpgcnt ## Drop PostgreSQL database
	@echo "Dropping PostgreSQL database..."
	docker exec -it $(PG_CONTAINER_NAME) dropdb -U $(PG_DB_USERNAME) $(PG_DB_NAME)

# Docker utility commands
stopall: ## Stop all running containers
	@echo "Stopping all running containers..."
	docker stop $(docker ps -q)

rmvall: ## Remove all containers
	@echo "Removing all containers..."
	docker rm $(docker ps -a -q)

# Go mod commands
modtidy: ## Tidy the Go modules
	@echo "Tidying Go modules..."
	go mod tidy

modvendor: ## Vendor Go modules
	@echo "Vendoring Go modules..."
	go mod vendor

# Help command
help: ## Show this help message
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*##"; printf "\n\033[1m%-12s\033[0m %s\n\n", "Command", "Description"} /^[a-zA-Z_-]+:.*?##/ { printf "\033[36m%-12s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

.PHONY: build run test lint fmt clean crtmgcnt strmgcnt stpmgcnt rmvmgcnt crtmgdb drpmgdb crtpgcnt strpgcnt stppgcnt rmvpgcnt crtpgdb drppgdb stopall rmvall modtidy modvendor help
