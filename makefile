# Makefile for building Go binaries and running the API server and Vue.js app

# Define the commands
CMDS = fetcher importer buddyfit

# Variables
GO_CMD=go
NPM_CMD=npm
API_DIR=cmd/api
FRONTEND_DIR=buddyfit-bot-chat

# Default target
all: $(CMDS) api frontend

# Build each command
$(CMDS):
	$(GO_CMD) build -o ./bin/$@ ./cmd/$@

# Run the Go API server
api:
	@echo "Running Go API server..."
	cd $(API_DIR) && $(GO_CMD) run main.go

# Run the Vue.js app
frontend:
	@echo "Running Vue.js app..."
	cd $(FRONTEND_DIR) && $(NPM_CMD) install && $(NPM_CMD) run dev

# Clean up binaries
clean:
	rm -f bin/*

.PHONY: all clean $(CMDS) api frontend