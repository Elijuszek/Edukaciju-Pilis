build:
	@go build -o bin/Edukaciju-Pilis cmd/main.go

test:
	@go test -v ./...
	
run: build
	@./bin/Edukaciju-Pilis