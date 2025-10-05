GO       ?= go
REVIVE   ?= revive
BIN_NAME ?= rong
VERSION  ?= $(shell git describe --tags)
PREFIX   ?= /usr/local/

BIN_FILE            = $(shell realpath -m "$(PREFIX)/bin/$(BIN_NAME)")
LICENSE_FILE        = $(shell realpath -m "$(PREFIX)/share/licenses/$(BIN_NAME)/LICENSE")
BASH_COMPLETION_DIR = $(shell realpath -m "$(PREFIX)/share/bash-completion/completions")
ZSH_COMPLETION_DIR  = $(shell realpath -m "$(PREFIX)/share/zsh/site-functions")
FISH_COMPLETION_DIR = $(shell realpath -m "$(PREFIX)/share/fish/vendor_completions.d")

-include Makefile.local

.PHONY: build install test docs-dev

build:
	$(GO) build -trimpath -ldflags '-X main.Version=$(VERSION)' -o $(BIN_NAME)

install:
	install -Dm755 $(BIN_NAME) "$(BIN_FILE)"
	install -Dm644 LICENSE "$(LICENSE_FILE)"
	$(BIN_NAME) _carapace bash | install -Dm644 /dev/stdin "$(BASH_COMPLETION_DIR)/$(BIN_NAME)"
	$(BIN_NAME) _carapace zsh  | install -Dm644 /dev/stdin "$(ZSH_COMPLETION_DIR)/_$(BIN_NAME)"
	$(BIN_NAME) _carapace fish | install -Dm644 /dev/stdin "$(FISH_COMPLETION_DIR)/$(BIN_NAME).fish"

test:
	$(GO) test -v ./...
	$(REVIVE) -config revive.toml -formatter friendly ./...
docs-dev:
	bun run dev
