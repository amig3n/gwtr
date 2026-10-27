PHONY: test clean

test:
	@echo "Running tests..."
	@go test ./...

build:
	@echo "Building the project..."
	@mkdir -p bin
	@go build -o bin/gwtr .

clean:
	@echo "Cleaning up..."
	@rm -rf bin

