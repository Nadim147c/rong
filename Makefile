GO     ?= go
REVIVE ?= revive

-include Makefile.local

lint:
	$(REVIVE) -config revive.toml -formatter friendly ./...
