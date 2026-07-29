.PHONY: build test lint run clean install release

BINARY := antipatterns
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

build:
	go build -o $(BINARY) ./cmd/antipatterns

test:
	go test ./... -v -race -cover

lint:
	golangci-lint run ./...

run: build
	./$(BINARY) scan .

clean:
	rm -f $(BINARY)
	rm -rf dist/

install: build
	cp $(BINARY) $$(go env GOPATH)/bin/

# Cross-compilation para releases
release: clean
	mkdir -p dist/
	GOOS=linux   GOARCH=amd64 go build -o dist/$(BINARY)-linux-amd64   ./cmd/antipatterns
	GOOS=linux   GOARCH=arm64 go build -o dist/$(BINARY)-linux-arm64   ./cmd/antipatterns
	GOOS=darwin  GOARCH=amd64 go build -o dist/$(BINARY)-darwin-amd64  ./cmd/antipatterns
	GOOS=darwin  GOARCH=arm64 go build -o dist/$(BINARY)-darwin-arm64  ./cmd/antipatterns
	GOOS=windows GOARCH=amd64 go build -o dist/$(BINARY)-windows-amd64.exe ./cmd/antipatterns
	ls -lh dist/
