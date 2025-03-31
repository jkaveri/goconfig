.PHONY: all test lint clean build

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOCLEAN=$(GOCMD) clean
GOMOD=$(GOCMD) mod
BINARY_NAME=goconfig

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) ./...

test:
	$(GOTEST) ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

deps:
	$(GOGET) ./...

tidy:
	$(GOMOD) tidy

lint:
	revive -config revive.toml --formatter friendly --exclude *_test.go ./...


# Development tools
.PHONY: install-tools
install-tools:
	go install github.com/mgechev/revive@latest