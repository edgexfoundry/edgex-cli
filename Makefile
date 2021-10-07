.PHONY: build test clean install

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
 
VERSION=$(shell cat ./VERSION 2>/dev/null || echo 0.0.0)
TIME=$(shell date)
GOFLAGS=-ldflags "-X 'github.com/edgexfoundry/edgex-cli.BuildVersion=$(VERSION)' -X 'github.com/edgexfoundry/edgex-cli.BuildTime=$(TIME)'"
ARTIFACT_ROOT?=bin

tidy:
	go mod tidy

build:
	@echo "GOPATH=$(GOPATH)"
	$(GO) build -o ${ARTIFACT_ROOT}/$(BINARY) $(GOFLAGS) ./cmd/edgex-cli

# initial impl. Feel free to override. Please keep ARTIFACT_ROOT coming from env though. CI/CD pipeline relies on this
build-all:
	@echo "GOPATH=$(GOPATH)"

	GOOS=linux GOARCH=amd64 $(GO) build -o ${ARTIFACT_ROOT}/$(BINARY)-linux-amd64 $(GOFLAGS) ./cmd/edgex-cli
	GOOS=linux GOARCH=arm64 $(GO) build -o ${ARTIFACT_ROOT}/$(BINARY)-linux-arm64 $(GOFLAGS) ./cmd/edgex-cli
	GOOS=darwin GOARCH=amd64 $(GO) build -o ${ARTIFACT_ROOT}/$(BINARY)-mac $(GOFLAGS) ./cmd/edgex-cli
	GOOS=windows GOARCH=amd64 $(GO) build -o ${ARTIFACT_ROOT}/$(BINARY).exe $(GOFLAGS) ./cmd/edgex-cli

	tar -czvf ${ARTIFACT_ROOT}/$(BINARY)-linux-amd64-$(VERSION).tar.gz Attribution.txt LICENSE res/sample-configuration.toml  ${ARTIFACT_ROOT}/$(BINARY)-linux-amd64 --transform 's/-linux-amd64//'
	tar -czvf ${ARTIFACT_ROOT}/$(BINARY)-linux-arm64-$(VERSION).tar.gz Attribution.txt LICENSE res/sample-configuration.toml  ${ARTIFACT_ROOT}/$(BINARY)-linux-arm64 --transform 's/-linux-arm64//'
	tar -czvf ${ARTIFACT_ROOT}/$(BINARY)-mac-$(VERSION).tar.gz Attribution.txt LICENSE res/sample-configuration.toml ${ARTIFACT_ROOT}/$(BINARY)-mac --transform 's/-mac//'
	zip  ${ARTIFACT_ROOT}/$(BINARY)-win-$(VERSION).zip ${ARTIFACT_ROOT}/$(BINARY).exe Attribution.txt LICENSE res/sample-configuration.toml

test:
	$(GO) test ./... -coverprofile coverage.out
	GO111MODULE=on go vet ./...
	gofmt -l $$(find . -type f -name '*.go'| grep -v "/vendor/")
	[ "`gofmt -l $$(find . -type f -name '*.go'| grep -v "/vendor/")`" = "" ]
	./bin/test-attribution-txt.sh

install:
	@echo "GOBIN=$(GOBIN)"
	$(GO) install $(GOFLAGS)

uninstall:
	@echo "GOBIN=$(GOBIN)"
	rm -f $(GOBIN)/$(BINARY)
	rm -rf $(HOME)/.edgex-cli


clean: uninstall
	rm -f $(BINARY)*
	rm -rf $(ARTIFACT_ROOT)/$(BINARY)*

vendor:
	$(GO) mod vendor

