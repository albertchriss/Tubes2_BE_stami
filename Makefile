APP_NAME=api
APP_DIR=cmd/$(APP_NAME)
BINARY_DIR=bin
DOCKER_DIR=deployments/docker
MIGRATION_DIR=internal/db/migrations

.PHONY: all
all: clean generate-docs build run

.PHONY: build
build:
	@echo "Building the application..."
	@go build -o $(BINARY_DIR)/$(APP_NAME) $(APP_DIR)/main.go
	@echo "Build completed."

.PHONY: run
run: build
	@echo "Running the application..."
	@./$(BINARY_DIR)/$(APP_NAME)

.PHONY: generate-docs
generate-docs:
	@echo "Generating documentation..."
	@mkdir -p docs
	@echo "Generating API documentation..."
	@swag init -g $(APP_DIR)/main.go -o docs
	@echo "API documentation generated in docs/ directory."

.PHONY: clean
clean:
	@echo "Cleaning up..."
	@rm -rf $(BINARY_DIR)/$(APP_NAME)
	@echo "Clean up completed."

.PHONY: up-dev
up-dev:
	@echo "Running docker compose..."
	@docker compose --env-file $(DOCKER_DIR)/.env.dev -f $(DOCKER_DIR)/docker-compose.dev.yaml up --build -d

.PHONY: down-dev
down-dev:
	@echo "Stopping docker compose..."
	@docker compose --env-file $(DOCKER_DIR)/.env.dev -f $(DOCKER_DIR)/docker-compose.dev.yaml down

.PHONY: up-prod
up-prod:
	@echo "Running docker compose..."
	@docker compose --env-file $(DOCKER_DIR)/.env.prod -f $(DOCKER_DIR)/docker-compose.prod.yaml up --build -d

.PHONY: down-prod
down-prod:
	@echo "Stopping docker compose..."
	@docker compose --env-file $(DOCKER_DIR)/.env.prod -f $(DOCKER_DIR)/docker-compose.prod.yaml down