build:
	@go build -o bin/ng

run: build
	@./bin/ng

test:
	@go test ./...
