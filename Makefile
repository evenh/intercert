COMMIT := $(shell git rev-parse --short HEAD)

all: clean deps proto test build

deps:
	go mod download

build:
	go build -race -ldflags "-s -w -X main.Version=DEV-SNAPSHOT -X main.Commit=$(COMMIT)" -o dist/intercert github.com/evenh/intercert

proto:
	@ if ! which buf > /dev/null; then \
		echo "error: buf not installed" >&2; \
		exit 1; \
	fi
	buf generate

test:
	go test -v github.com/evenh/intercert/...

clean:
	rm -rf ./dist
	go clean -i github.com/evenh/intercert/...

.PHONY: \
	all \
	deps \
	build \
	proto \
	test \
	clean
