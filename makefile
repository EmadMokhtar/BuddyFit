# Makefile for building Go binaries and running the API server and Vue.js app

# Define the commands
CMDS = fetcher importer buddyfit api

# Variables
GO_CMD=go
NPM_CMD=npm
API_DIR=cmd/api
FRONTEND_DIR=buddyfit-bot-chat
DOCKER_DIR=docker

# Default target
all: $(CMDS)

# Build each command
$(CMDS):
	$(GO_CMD) build -o ./bin/$@ ./cmd/$@

# Build the Vue.js app
frontend:
	cd $(FRONTEND_DIR) && $(NPM_CMD) install && $(NPM_CMD) run build && $(NPM_CMD) install -g serve && serve -s dist

# Run the Go API server
run-api:
	@echo "Running Go API server..."
	cd $(API_DIR) && $(GO_CMD) run main.go

# Run the Vue.js app
run-frontend:
	@echo "Running Vue.js app..."
	cd $(FRONTEND_DIR) && $(NPM_CMD) install && $(NPM_CMD) run dev

# Clean up binaries
clean:
	rm -f bin/*

# Database commands
migrate-up:
	@migrate -database '$(BF_DB_URL)' -path ./migrations up

database-start:
	cd $(DOCKER_DIR) && docker compose up -d timescaledb

database-stop:
	cd $(DOCKER_DIR) && docker compose down timescaledb

vectorizer-worker-start:
	cd $(DOCKER_DIR) && docker compose up -d vectorizer-worker


.PHONY: all clean $(CMDS) api frontend