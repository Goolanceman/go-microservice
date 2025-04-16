.PHONY: build run test clean docker-build docker-run swagger

# Build the application
build:
	go build -o bin/api ./cmd/api

# Run the application
run:
	go run ./cmd/api

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Build Docker image
docker-build:
	docker build -t go-microservice:latest .

# Run Docker container
docker-run:
	docker run -p 8080:8080 go-microservice:latest

# Generate Swagger documentation
swagger:
	swag init -g cmd/api/main.go -o docs

# Install dependencies
deps:
	go mod download

# Format code
fmt:
	go fmt ./...

# Vet code
vet:
	go vet ./...

# Lint code
lint:
	golangci-lint run

# All-in-one development setup
dev: deps fmt vet lint test build 