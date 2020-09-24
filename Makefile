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

build:
	echo "GOPATH=$(GOPATH)"
	$(GO) build -o $(BINARY) $(GOFLAGS)

# initial impl. Feel free to override. Please keep ARTIFACT_ROOT coming from env though. CI/CD pipeline relies on this
build-all:
	echo "GOPATH=$(GOPATH)"
	mkdir -p ${ARTIFACT_ROOT}
	GOOS=linux GOARCH=amd64 $(GO) build -o ${ARTIFACT_ROOT}/$(BINARY)-linux-amd64 $(GOFLAGS)
	GOOS=linux GOARCH=arm64 $(GO) build -o ${ARTIFACT_ROOT}/$(BINARY)-linux-arm64 $(GOFLAGS)
	GOOS=darwin GOARCH=amd64 $(GO) build -o ${ARTIFACT_ROOT}/$(BINARY)-mac $(GOFLAGS)
	GOOS=windows GOARCH=amd64 $(GO) build -o ${ARTIFACT_ROOT}/$(BINARY)-win.exe $(GOFLAGS)

	tar -czvf ${ARTIFACT_ROOT}/$(BINARY)-linux-amd64-$(VERSION).tar.gz -C ${ARTIFACT_ROOT} $(BINARY)-linux-amd64
	tar -czvf ${ARTIFACT_ROOT}/$(BINARY)-linux-arm64-$(VERSION).tar.gz -C ${ARTIFACT_ROOT} $(BINARY)-linux-arm64
	tar -czvf ${ARTIFACT_ROOT}/$(BINARY)-mac-$(VERSION).tar.gz -C ${ARTIFACT_ROOT} $(BINARY)-mac
	zip -j ${ARTIFACT_ROOT}/$(BINARY)-win-$(VERSION).zip ${ARTIFACT_ROOT}/$(BINARY)-win.exe


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
	-rm $(BINARY)
	-rm -rf $(ARTIFACT_ROOT)
