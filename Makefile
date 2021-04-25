GO = go

tidify-deps:
	$(GO) mod tidy
.PHONY: tidify-deps

install-deps: tidify-deps
	$(GO) get
.PHONY: install-deps

build: repeat
	$(GO) build
.PHONY: build
