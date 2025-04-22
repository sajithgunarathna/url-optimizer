.PHONY: build run test lint clean

build:
	go build -o bin/web-analyzer ./cmd/web-analyzer

run: build
	./bin/web-analyzer

test:
	go test ./...

lint:
	golangci-lint run

clean:
	rm -rf bin/