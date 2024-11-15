GO ?= go

.PHONY: build
build:
	$(GO) build

.PHONY: fmt
fmt:
	$(GO) fmt ./...

.PHONY: vet
vet: build
	$(GO) vet ./...
	$(GO) vet -vettool=khulnasoft-vet ./...

.PHONY: lint
lint:
	@hash golangci-lint > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		export BINARY="golangci-lint"; \
		curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell $(GO) env GOPATH)/bin v1.52.2; \
	fi
	golangci-lint run --timeout 10m
