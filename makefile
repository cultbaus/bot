BINDIR    := $(CURDIR)/bin
BIN       := bot
PKG       := ./...
SRC       := $(shell find . -type f -name '*.go' -print) go.mod go.sum
VERSION   ?= latest

LDFLAGS     := -w -s
CGO_ENABLED := 0

all: build

build: $(BINDIR)/$(BIN)

$(BINDIR)/$(BIN): $(SRC)
	CGO_ENABLED=$(CGO_ENABLED) go build $(GOFLAGS) -trimpath -ldflags '$(LDFLAGS)' -o '$(BINDIR)'/$(BIN) ./cmd/$(BIN)

clean:
	@rm -rf '$(BINDIR)'

image:
	docker build -t cultbaus/$(BIN):$(VERSION) build/package/$(BIN)

.PHONY: build clean image
