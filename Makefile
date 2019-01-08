COMMIT := $(shell git rev-parse --short HEAD)

all: deps proto build

deps:
	@ if ! which dep > /dev/null; then \
		echo "warning: dep not installed - doing it now" >&2; \
		go get -u -v github.com/golang/dep/cmd/dep; \
	fi
	dep ensure

build:
	go build -race -ldflags "-s -w -X main.Version=DEV-SNAPSHOT -X main.Commit=$(COMMIT)" -o dist/intercert github.com/evenh/intercert

proto:
	@ if ! which protoc > /dev/null; then \
		echo "error: protoc not installed" >&2; \
		exit 1; \
	fi
	go get -u -v github.com/golang/protobuf/protoc-gen-go
	for file in $$(git ls-files '*.proto'); do \
		protoc -I $$(dirname $$file) --go_out=plugins=grpc:$$(dirname $$file) $$file; \
	done

test:
	go test -v -cpu 1,4 github.com/evenh/intercert/...

clean:
	go clean -i github.com/evenh/intercert/...

coverage:
	./coverage.sh --coveralls

.PHONY: \
	all \
	deps \
	build \
	proto \
	test \
	clean \
	coverage
