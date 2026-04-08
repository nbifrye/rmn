VERSION ?= dev
COMMIT  := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE    := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

.PHONY: build test vet lint clean install

build:
	go build $(LDFLAGS) -o rmn ./cmd/rmn/

test:
	go test ./...

vet:
	go vet ./...

lint: vet
	@which golangci-lint > /dev/null 2>&1 && golangci-lint run || echo "golangci-lint not installed, skipped"

clean:
	rm -f rmn

install:
	go install $(LDFLAGS) ./cmd/rmn/
