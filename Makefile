.PHONY: all build test clean lint vet generate install uninstall release docker

VERSION := $(shell cat .version 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS := -ldflags "-X neocut/internal/config.Commit=$(COMMIT)"
BINARY := neocut
OUTDIR := bin

all: build

generate:
	go generate ./...

build: generate
	go build $(LDFLAGS) -o $(BINARY) ./cmd/neocut/

build-all: generate
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(OUTDIR)/neocut-linux-amd64 ./cmd/neocut/
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(OUTDIR)/neocut-linux-arm64 ./cmd/neocut/
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(OUTDIR)/neocut-darwin-amd64 ./cmd/neocut/
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(OUTDIR)/neocut-darwin-arm64 ./cmd/neocut/
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(OUTDIR)/neocut-windows-amd64.exe ./cmd/neocut/
	GOOS=windows GOARCH=arm64 go build $(LDFLAGS) -o $(OUTDIR)/neocut-windows-arm64.exe ./cmd/neocut/

test:
	go test -count=1 ./...

vet:
	go vet ./...

lint:
	golangci-lint run ./... 2>/dev/null || go vet ./...

clean:
	go clean
	rm -f $(BINARY) $(BINARY).exe
	rm -rf $(OUTDIR)

install: build
	cp $(BINARY) $(GOPATH)/bin/$(BINARY) 2>/dev/null || true

uninstall:
	rm -f $(GOPATH)/bin/$(BINARY) 2>/dev/null || true

release: build-all
	@echo "Release $(VERSION) built in $(OUTDIR)/"

docker:
	docker build -t neocut:$(VERSION) .

docker-run:
	docker run --rm -v "$(PWD):/workspace" neocut:$(VERSION)

.DEFAULT_GOAL := all
