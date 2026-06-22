.PHONY: build test run clean lint test-cover

BINARY=openlibing

build:
	go build -o bin/$(BINARY) ./cmd/openlibing/

test:
	go test ./... -v -count=1

test-cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out

run: build
	./bin/$(BINARY)

clean:
	rm -rf bin/

lint:
	go vet ./...
