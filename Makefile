build:
	@go build -o bin/Edukaciju-Pilis main.go

test:
	@go test -v ./...
	
run: build
	@./bin/Edukaciju-Pilis

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down
	
docker-build:
	@echo "Building the Docker image..."
	@docker build -f Dockerfile -t go-containerized:latest .