.PHONY: build test lint clean release

APPS := focus noise syntax prism
BINARIES := focus focus-mcp noise syntax prism
BUILD_DIR := ./dist

build:
	@for app in $(APPS); do \
		echo "Building $$app..."; \
		go build -o $(BUILD_DIR)/ ./apps/$$app/cmd/...; \
	done

test:
	@for app in $(APPS); do \
		echo "Testing $$app..."; \
		cd apps/$$app && go test ./... && cd ../..; \
	done

lint:
	golangci-lint run ./apps/...

clean:
	rm -rf $(BUILD_DIR)

release:
	goreleaser release --rm-dist

snapshot:
	goreleaser build --snapshot --rm-dist

sync:
	go work sync
