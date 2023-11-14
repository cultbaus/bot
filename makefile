BINDIR    := $(CURDIR)/bin
BIN       := bot
PKG       := ./...
SRC       := $(shell find . -type f -name '*.go' -print) go.mod go.sum
TESTS     := .
TESTFLAGS :=
VERSION   ?= latest

LDFLAGS     := -w -s
CGO_ENABLED := 0

all: build

build: $(BINDIR)/$(BIN)

$(BINDIR)/$(BIN): $(SRC)
	CGO_ENABLED=$(CGO_ENABLED) go build $(GOFLAGS) -trimpath -ldflags '$(LDFLAGS)' -o '$(BINDIR)'/$(BIN) ./cmd/bot

test: build
test: TESTFLAGS += -race -v -count=1
test: unit

unit:
	go test $(GOFLAGS) -run $(TESTS) $(PKG) $(TESTFLAGS)

clean:
	@rm -rf '$(BINDIR)'

image:
	docker build -t cultbaus/bot:$(VERSION) build/package/bot

.PHONY: build test clean
