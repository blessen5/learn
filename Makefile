BINARY    := learn
PREFIX    := /usr/local
DESTDIR   :=
BINDIR    := $(PREFIX)/bin
COMPLDIR  := $(PREFIX)/share/bash-completion/completions
VERSION   := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS   := -ldflags "-s -w -X main.version=$(VERSION)"
GOFLAGS   :=

.PHONY: all build install uninstall clean test fmt vet completions help

all: build

build:
	go build $(GOFLAGS) $(LDFLAGS) -o $(BINARY) .

install: build completions-bash
	install -Dm755 $(BINARY) $(DESTDIR)$(BINDIR)/$(BINARY)
	install -Dm644 bash-completion $(DESTDIR)$(COMPLDIR)/$(BINARY)

uninstall:
	rm -f $(DESTDIR)$(BINDIR)/$(BINARY)
	rm -f $(DESTDIR)$(COMPLDIR)/$(BINARY)

clean:
	rm -f $(BINARY)
	rm -f bash-completion zsh-completion fish-completion

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

completions: completions-bash completions-zsh completions-fish

completions-bash: build
	./$(BINARY) completion bash > bash-completion

completions-zsh: build
	./$(BINARY) completion zsh > zsh-completion

completions-fish: build
	./$(BINARY) completion fish > fish-completion

help:
	@echo "Targets:"
	@echo "  build           Build the binary"
	@echo "  install         Install binary to $(DESTDIR)$(BINDIR) and bash completion"
	@echo "  uninstall       Remove installed binary and completion"
	@echo "  clean           Remove build artifacts"
	@echo "  test            Run tests"
	@echo "  fmt             Format Go source"
	@echo "  vet             Run go vet"
	@echo "  completions     Generate all shell completion scripts"
	@echo "  help            Show this help"
