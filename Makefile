.PHONY: test clean checkenv deps

checkenv:
	@echo "Checking environment..."
	@which go > /dev/null || (echo "Go is not installed. Please install Go." && exit 1)
	@echo "Go version: $$(go version)"

deps: checkenv
	@echo "Installing dependencies..."
	@go mod tidy && echo "go.mod file prepared"
	@go mod download && echo "Dependencies downloaded"
	@echo "Dependencies updated"


test: 
	@echo "Running all tests..."
	@go test ./...

build: deps
	@echo "Building the project..."
	@mkdir -p bin
	@go build -o bin/gwtr .

clean:
	@echo "Cleaning up..."
	@rm -rf bin
	@echo "Done."

