build:
	@go build -o bin/Edukaciju-Pilis main.go

test:
	@go test -v ./...
	
run: build
	@./bin/Edukaciju-Pilis