# ============================================================================ #
# BUILD
# ============================================================================ #
## build/cmd: build the cmd/ application
.PHONY: build/server
build/cmd:
	@echo "Building cmd/..."
	go build -o=./bin ./server.go