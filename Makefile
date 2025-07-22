GO     ?= go
REVIVE ?= revive
OUTBIN ?= rong
PREFIX = /usr/local/

-include Makefile.local

build:
	$(GO) build -trimpath -o $(OUTBIN)

install:
	install -Dm755 $(OUTBIN) "$(PREFIX)/bin/$(OUTBIN)"
	install -Dm644 LICENSE "$(PREFIX)/share/licenses/$(OUTBIN)/LICENSE"

lint:
	$(REVIVE) -config revive.toml -formatter friendly ./...

docs-dev:
	bun run dev
