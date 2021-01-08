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
	mkdir -p ${ARTIFACT_ROOT}
	mkdir -p ${ARTIFACT_ROOT_AMD64}
	mkdir -p ${ARTIFACT_ROOT_ARM64}
	mkdir -p ${ARTIFACT_ROOT_MAC}
	mkdir -p ${ARTIFACT_ROOT_WIN}
	GOOS=linux GOARCH=amd64 $(GO) build -o ${ARTIFACT_ROOT_AMD64}/$(BINARY)-linux-amd64 $(GOFLAGS); cd ${ARTIFACT_ROOT_AMD64}; ln -s $(BINARY)-linux-amd64  $(BINARY)
	GOOS=linux GOARCH=arm64 $(GO) build -o ${ARTIFACT_ROOT_ARM64}/$(BINARY)-linux-arm64 $(GOFLAGS); cd ${ARTIFACT_ROOT_ARM64}; ln -s $(BINARY)-linux-arm64 $(BINARY)
	GOOS=darwin GOARCH=amd64 $(GO) build -o ${ARTIFACT_ROOT_MAC}/$(BINARY)-mac $(GOFLAGS); cd ${ARTIFACT_ROOT_MAC}; ln -s $(BINARY)-mac $(BINARY)
	GOOS=windows GOARCH=amd64 $(GO) build -o ${ARTIFACT_ROOT_WIN}/$(BINARY)-win.exe $(GOFLAGS); cd ${ARTIFACT_ROOT_WIN}; ln -s  $(BINARY)-win.exe $(BINARY).exe

	tar -czvf ${ARTIFACT_ROOT}/$(BINARY)-linux-amd64-$(VERSION).tar.gz Attribution.txt LICENSE -C ${ARTIFACT_ROOT_AMD64} $(BINARY)-linux-amd64 $(BINARY)
	tar -czvf ${ARTIFACT_ROOT}/$(BINARY)-linux-arm64-$(VERSION).tar.gz Attribution.txt LICENSE -C ${ARTIFACT_ROOT_ARM64} $(BINARY)-linux-arm64 $(BINARY)
	tar -czvf ${ARTIFACT_ROOT}/$(BINARY)-mac-$(VERSION).tar.gz Attribution.txt LICENSE -C ${ARTIFACT_ROOT_MAC} $(BINARY)-mac $(BINARY)
	zip -j --symlinks ${ARTIFACT_ROOT}/$(BINARY)-win-$(VERSION).zip ${ARTIFACT_ROOT_WIN}/$(BINARY)-win.exe ${ARTIFACT_ROOT_WIN}/$(BINARY).exe Attribution.txt LICENSE


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
