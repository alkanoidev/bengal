all: build

# Build the web application
build:
	@echo "Building web..."
	@go build -o main.exe cmd/api/main.go

# Build the cli application
build-cli:
	@echo "Building cli..."
	@go build -o email-checker.exe cmd/cli/main.go

# Run the web application
run:
	@go run cmd/api/main.go

# Run the cli application
run-cli:
	@go run cmd/cli/main.go

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload

watch:
	@air

.PHONY: all build run test clean watch
