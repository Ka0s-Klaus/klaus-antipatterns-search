.PHONY: build test lint run clean install

BINARY := antipatterns
MODULE  := github.com/Ka0s-Klaus/Klaus-antipatterns-search
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X $(MODULE)/cmd/antipatterns.version=$(VERSION)"

build:
	go build $(LDFLAGS) -o $(BINARY) ./cmd/antipatterns

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
	mv $(BINARY) $(GOPATH)/bin/

# Cross-compilation para releases
release:
	GOOS=linux   GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY)-linux-amd64   ./cmd/antipatterns
	GOOS=linux   GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY)-linux-arm64   ./cmd/antipatterns
	GOOS=darwin  GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY)-darwin-amd64  ./cmd/antipatterns
	GOOS=darwin  GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY)-darwin-arm64  ./cmd/antipatterns
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY)-windows-amd64.exe ./cmd/antipatterns
