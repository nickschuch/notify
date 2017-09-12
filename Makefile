#!/usr/bin/make -f

# This allows us to have a GOPATH per peroject.
# Ensuring that we are all building with the same dependencies.
GOPATH := $(shell pwd):$(shell pwd)/vendor
export GOPATH
export CGO_ENABLED=0

NAME=notify
PACKAGE=github.com/nickschuch/$(NAME)

# Build binaries for linux/amd64 and darwin/amd64
build:
	gox -os='linux darwin' -arch='amd64' -output='bin/$(NAME)_{{.OS}}_{{.Arch}}' -ldflags='-extldflags "-static"' $(PACKAGE)

# Run all lint checking with exit codes for CI
lint:
	golint -set_exit_status $(PACKAGE)/...

# Run tests with coverage reporting
test:
	go test -cover $(PACKAGE)/...
