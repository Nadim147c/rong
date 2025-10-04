GO       ?= go
REVIVE   ?= revive
BIN_NAME ?= rong
VERSION  ?= $(shell git describe --tags)
PREFIX   ?= /usr/local/

BIN_FILE     = $(shell realpath "$(PREFIX)/bin/$(BIN_NAME)")
LICENSE_FILE = $(shell realpath "$(PREFIX)/share/licenses/$(BIN_NAME)/LICENSE")

-include Makefile.local

build:
	$(GO) build -trimpath -ldflags '-s -w -X main.Version=$(VERSION)' -o $(BIN_NAME)

install:
	install -Dm755 $(BIN_NAME) "$(BIN_FILE)"
	install -Dm644 LICENSE "$(LICENSE_FILE)"

test:
	$(GO) test -v ./...
	$(GO) run github.com/mgechev/revive@latest -config revive.toml -formatter friendly ./...

docs-dev:
	bun run dev
