GO ?= go
GOLANGCI_LINT ?= golangci-lint
GOLANGCI_LINT_VERSION ?= v1.30.0

BINARY="resolver"

LDFLAGS ?= -s -w
ifdef COMMIT
LDFLAGS += -X github.com/ethersphere/resolver.commit="$(COMMIT)"
endif

.PHONY: all
all: build lint vet test binary

.PHONY: binary
binary: export CGO_ENABLED=0
binary: dist
	$(GO) version
	$(GO) build -trimpath -ldflags "$(LDFLAGS)" -o dist/$(BINARY) ./cmd/resolver

dist:
	mkdir $@

.PHONY: lint
lint: linter
	$(GOLANGCI_LINT) run

.PHONY: linter
linter:
	which $(GOLANGCI_LINT) || ( cd /tmp && GO111MODULE=on $(GO) get -u github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION) )

.PHONY: vet
vet:
	$(GO) vet ./...

.PHONY: test
test:
	$(GO) test -v ./...

.PHONY: test-integration
test-integration:
	$(GO) test -tags=integration -v ./...

.PHONY: build
build: export CGO_ENABLED=0
build:
	$(GO) build -trimpath -ldflags "$(LDFLAGS)" ./...

.PHONY: clean
clean:
	$(GO) clean
	rm -rf dist/

