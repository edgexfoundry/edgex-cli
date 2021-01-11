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
GOFLAGS=-ldflags "-X github.com/edgexfoundry/edgex-cli/cmd/version.Version=$(VERSION)"
ARTIFACT_ROOT?=bin
ARTIFACT_ROOT_AMD64=${ARTIFACT_ROOT}/AMD64
ARTIFACT_ROOT_ARM64=${ARTIFACT_ROOT}/ARM64
ARTIFACT_ROOT_MAC=${ARTIFACT_ROOT}/MAC
ARTIFACT_ROOT_WIN=${ARTIFACT_ROOT}/WIN
build:
	echo "GOPATH=$(GOPATH)"
	$(GO) build -o $(BINARY) $(GOFLAGS)

# initial impl. Feel free to override. Please keep ARTIFACT_ROOT coming from env though. CI/CD pipeline relies on this
build-all:
	echo "GOPATH=$(GOPATH)"

	GOOS=linux GOARCH=amd64 $(GO) build -o ${ARTIFACT_ROOT}/$(BINARY)-linux-amd64 $(GOFLAGS)
	GOOS=linux GOARCH=arm64 $(GO) build -o ${ARTIFACT_ROOT}/$(BINARY)-linux-arm64 $(GOFLAGS)
	GOOS=darwin GOARCH=amd64 $(GO) build -o ${ARTIFACT_ROOT}/$(BINARY)-mac $(GOFLAGS)
	GOOS=windows GOARCH=amd64 $(GO) build -o ${ARTIFACT_ROOT}/$(BINARY).exe $(GOFLAGS)

	tar -czvf ${ARTIFACT_ROOT}/$(BINARY)-linux-amd64-$(VERSION).tar.gz Attribution.txt LICENSE res/.sample-configuration.toml  ${ARTIFACT_ROOT}/$(BINARY)-linux-amd64 --transform 's/-linux-amd64//'
	tar -czvf ${ARTIFACT_ROOT}/$(BINARY)-linux-arm64-$(VERSION).tar.gz Attribution.txt LICENSE res/.sample-configuration.toml  ${ARTIFACT_ROOT}/$(BINARY)-linux-arm64 --transform 's/-linux-arm64//'
	tar -czvf ${ARTIFACT_ROOT}/$(BINARY)-mac-$(VERSION).tar.gz Attribution.txt LICENSE res/.sample-configuration.toml ${ARTIFACT_ROOT}/$(BINARY)-mac --transform 's/-mac//'
	zip  ${ARTIFACT_ROOT}/$(BINARY)-win-$(VERSION).zip ${ARTIFACT_ROOT}/$(BINARY).exe Attribution.txt LICENSE res/.sample-configuration.toml


test:
	$(GO) test ./... -coverprofile coverage.out
	GO111MODULE=on go vet ./...
	gofmt -l .
	[ "`gofmt -l .`" = "" ]
	./bin/test-go-mod-tidy.sh
	./bin/test-attribution-txt.sh

install:
	echo "GOBIN=$(GOBIN)"
	$(GO) install $(GOFLAGS)
	mkdir -p $(HOME)/.edgex-cli
	cp ./res/configuration.toml $(HOME)/.edgex-cli/configuration.toml

uninstall:
	echo "GOBIN=$(GOBIN)"
	rm $(GOBIN)/$(BINARY)
	rm -rf $(HOME)/.edgex-cli

clean:
	-rm $(BINARY)*
	-rm -rf $(ARTIFACT_ROOT)/$(BINARY)
	-rm -rf ${ARTIFACT_ROOT_AMD64}
	-rm -rf ${ARTIFACT_ROOT_ARM64}
	-rm -rf ${ARTIFACT_ROOT_MAC}
	-rm -rf ${ARTIFACT_ROOT_WIN}
