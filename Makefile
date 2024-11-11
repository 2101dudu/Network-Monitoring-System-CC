# Variables
BINARY_NAME := nms
CMD_PATH := ./cmd/nms
AGENT_RUNNER := ./internal/agent/runner.go
SERVER_RUNNER := ./internal/server/runner.go

# Default target
.PHONY: all
all: build

# Build and run the agent component
.PHONY: build-agent
build-agent:
	@echo "Building and running agent..."
	go run $(AGENT_RUNNER)
	@echo "Agent execution finished."

# Build and run the server component
.PHONY: build-server
build-server:
	@echo "Building and running server..."
	go run $(SERVER_RUNNER)
	@echo "Server execution finished."

# Clean built binaries and Go cache
.PHONY: clean
clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)
	go clean -cache
	@echo "Clean complete."


