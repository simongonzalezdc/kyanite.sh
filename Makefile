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
	@for app in $(APPS); do \
		echo "Linting $$app..."; \
		cd apps/$$app && golangci-lint run ./... && cd ../..; \
	done

clean:
	rm -rf $(BUILD_DIR)

release:
	goreleaser release --clean

snapshot:
	goreleaser build --snapshot --clean

sync:
	go work sync
