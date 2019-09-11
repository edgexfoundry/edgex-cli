.PHONY: build clean install

GO = CGO_ENABLED=0 GO111MODULE=on go

BINARY=edgex

VERSION=$(shell cat ./VERSION)
GOFLAGS=-ldflags "-X github.com/edgexfoundry-holding/edgex-cli/cmd/version.version=$(VERSION)"

build:
	$(GO) build -o $(BINARY) $(GOFLAGS)

install:
	$(GO) install $(GOFLAGS)

clean:
	-rm $(BINARY)