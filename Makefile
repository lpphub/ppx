VERSION ?= dev
GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

.PHONY: build clean install test

build:
	go build -ldflags "-X 'ppx/cmd.Version=$(VERSION)' \
					 -X 'ppx/cmd.GitCommit=$(GIT_COMMIT)' \
					 -X 'ppx/cmd.BuildTime=$(BUILD_TIME)'" -o ppx .

clean:
	rm -f ppx

install: build
	mv ppx $(GOPATH)/bin/

test:
	go test ./...