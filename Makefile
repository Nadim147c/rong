GO       ?= go
BUN      ?= bun
NAME     ?= rong
VERSION  ?= $(shell git describe --tags)
PREFIX   ?= /usr/local/
TOOL_MOD ?= -modfile tool.go.mod
TOOL     ?= $(GO) tool $(TOOL_MOD)

BUILD_DIR ?= build
BUILD_BIN ?= $(BUILD_DIR)/$(NAME)
BUILD_COMPLETION_DIR ?= $(BUILD_DIR)/completions

INSTALL_BIN                 = $(shell realpath -m "$(PREFIX)/bin/$(NAME)")
INSTALL_LICENSE             = $(shell realpath -m "$(PREFIX)/share/licenses/$(NAME)/LICENSE")
INSTALL_BASH_COMPLETION_DIR = $(shell realpath -m "$(PREFIX)/share/bash-completion/completions")
INSTALL_ZSH_COMPLETION_DIR  = $(shell realpath -m "$(PREFIX)/share/zsh/site-functions")
INSTALL_FISH_COMPLETION_DIR = $(shell realpath -m "$(PREFIX)/share/fish/vendor_completions.d")

-include Makefile.local

.PHONY: build install test docs-dev generate-completion

build:
	$(GO) build -trimpath -ldflags '-s -w -X main.Version=$(VERSION)' -o $(BUILD_BIN)

install:
	install -Dm755 $(BUILD_BIN) "$(INSTALL_BIN)"
	install -Dm644 LICENSE "$(INSTALL_LICENSE)"
	install -Dm644 "$(BUILD_COMPLETION_DIR)/$(NAME).bash" "$(INSTALL_BASH_COMPLETION_DIR)/$(NAME)"
	install -Dm644 "$(BUILD_COMPLETION_DIR)/$(NAME).zsh"  "$(INSTALL_ZSH_COMPLETION_DIR)/_$(NAME)"
	install -Dm644 "$(BUILD_COMPLETION_DIR)/$(NAME).fish" "$(INSTALL_FISH_COMPLETION_DIR)/$(NAME).fish"

generate-completion: build
	mkdir -p "$(BUILD_COMPLETION_DIR)"
	$(BUILD_BIN) _carapace bash > "$(BUILD_COMPLETION_DIR)/$(NAME).bash"
	$(BUILD_BIN) _carapace zsh  > "$(BUILD_COMPLETION_DIR)/$(NAME).zsh"
	$(BUILD_BIN) _carapace fish > "$(BUILD_COMPLETION_DIR)/$(NAME).fish"

tools-install:
	$(GO) get $(TOOL_MOD) -tool github.com/mgechev/revive@latest
	$(GO) get $(TOOL_MOD) -tool github.com/segmentio/golines@latest
	$(GO) get $(TOOL_MOD) -tool mvdan.cc/gofumpt@latest
	$(GO) mod tidy $(TOOL_MOD)

format:
	find -iname '*.go' -print0 | xargs -0 $(TOOL) golines --max-len 80 -w
	find -iname '*.go' -print0 | xargs -0 $(TOOL) gofumpt -w

test:
	$(GO) test -v ./...
	$(TOOL) revive -config revive.toml -formatter friendly ./...

docs-dev:
	$(BUN) run dev
