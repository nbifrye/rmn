VERSION ?= dev
COMMIT  := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE    := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

.PHONY: build test vet lint clean install cover

build:
	go build $(LDFLAGS) -o rmn ./cmd/rmn/

test:
	go test ./...

vet:
	go vet ./...

lint: vet
	@which golangci-lint > /dev/null 2>&1 && golangci-lint run || echo "golangci-lint not installed, skipped"

cover:
	go test -coverprofile=coverage.out -covermode=atomic ./...
	@TOTAL=$$(go tool cover -func=coverage.out | grep '^total:' | awk '{print $$NF}' | tr -d '%'); \
	echo "Coverage: $${TOTAL}%"; \
	if [ $$(echo "$${TOTAL} < 100.0" | bc -l) -eq 1 ]; then \
		echo "FAIL: Coverage below 100%"; \
		go tool cover -func=coverage.out | grep -v '100.0%' | grep -v 'total:'; \
		exit 1; \
	fi

clean:
	rm -f rmn coverage.out

install:
	go install $(LDFLAGS) ./cmd/rmn/
