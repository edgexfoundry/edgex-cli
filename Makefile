.PHONY: build clean install

GO = CGO_ENABLED=0 GO111MODULE=off go

BINARY=edgex

build:
	$(GO) build -o $(BINARY)

install:
	$(GO) install 

clean:
	-rm $(BINARY)