
AGENT_RUNNER := ./cmd/nms/agent/runner.go
SERVER_RUNNER := ./cmd/nms/server/runner.go

.PHONY: all build build-agent build-server agent server clean
all: build-agent build-server

build-agent:
	@echo "Building agent..."
	@mkdir -p out/bin
	go build -o out/bin/agent $(AGENT_RUNNER)

build-server:
	@echo "Building server..."
	@mkdir -p out/bin
	go build -o out/bin/server $(SERVER_RUNNER)

agent:
	@echo "Running agent..."
	@./out/bin/agent

server:
	@echo "Running server..."	
	@./out/bin/server

clean:
	@echo "Cleaning up..."
	rm -rf out
	@echo "Clean complete."


