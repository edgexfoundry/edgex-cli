.PHONY: build clean install

GO = CGO_ENABLED=0 GO111MODULE=on go

# If GOPATH undefined, base within module to avoid collisions"
ifndef GOPATH
  	GOMOD=$(shell go env GOMOD)
  GOPATH=$(dir ${GOMOD})go
  export GOPATH
endif

ifndef GOBIN
  	GOBIN=$(GOPATH)/bin
  export GOBIN
endif



BINARY=edgex-cli

VERSION=$(shell cat ./VERSION)
GOFLAGS=-ldflags "-X github.com/edgexfoundry-holding/edgex-cli/cmd/version.version=$(VERSION)"

build:
	echo "GOPATH=$(GOPATH)"
	$(GO) build -o $(BINARY) $(GOFLAGS)

test:
	$(GO) test ./... -coverprofile coverage.out
install:
	echo "GOBIN=$(GOBIN)"
	$(GO) install $(GOFLAGS)

clean:
	-rm $(BINARY)
