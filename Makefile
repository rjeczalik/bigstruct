all: help

PROJECT   := github.com/rjeczalik/bigstruct
GOVERSION := 1.15.5

# build Builds bigstruct binary
.PHONY: build
build: GIT_COMMIT:= $(or ${GIT_COMMIT},$(shell git rev-list -1 HEAD))
build:
	@go build -ldflags '-X "main.gitCommit=$(GIT_COMMIT)"' -o bin/bigstruct ./cmd/bigstruct

# docker-build Builds bigstruct binary with golang container
.PHONY: docker-build
docker-build: GIT_COMMIT:= $(shell git rev-list -1 HEAD)
docker-build:
	@docker run --rm -u $(shell id -u):$(shell id -g) \
		-e GIT_COMMIT=$(GIT_COMMIT) \
		-e GOCACHE=/go/cache \
		-e GOMODCACHE=/go/modcache \
		-v $(PWD):/go/bigstruct \
		-v $(shell go env GOCACHE):/go/cache \
		-v $(shell go env GOMODCACHE):/go/modcache \
		-w /go/bigstruct \
	golang:$(GOVERSION) ./make.sh

	@docker run --rm -u $(shell id -u):$(shell id -g) \
		-v ${PWD}:/app \
		gruebel/upx --best --lzma /app/bin/bigstruct

# help Shows targets list
.PHONY: help
help:
	@awk -F ' ' '/^# / {cmd=$$2; $$2=""; printf "\033[36m%-25s\033[0m %s\n", cmd, tolower($$0)}' $(MAKEFILE_LIST) | sort
