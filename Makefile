.PHONY: build launcher test cover tidy vulncheck verify fmt lint clean release snapshot sync

APPS := focus noise syntax prism
BINARIES := focus focus-mcp noise syntax prism kyanite
BUILD_DIR := ./dist

build: build-apps launcher

build-apps:
	@for app in $(APPS); do \
		echo "Building $$app..."; \
		go build -o $(BUILD_DIR)/ ./apps/$$app/cmd/...; \
	done

launcher:
	@echo "Building kyanite launcher..."
	@go build -o $(BUILD_DIR)/kyanite ./cmd/kyanite/

desktop: build
	@echo "Creating desktop launcher..."
	@cp -r desktop/mac/kyanite.app $(BUILD_DIR)/kyanite.app
	@cp desktop/linux/kyanite.desktop $(BUILD_DIR)/kyanite.desktop
	@echo "Done. Run: open $(BUILD_DIR)/kyanite.app"

TEST_DIRS := pkg/design pkg/config pkg/ai pkg/tui pkg/cache pkg/session pkg/testutil pkg/appnames pkg/export apps/focus apps/noise apps/syntax apps/prism cmd/kyanite

test:
	@failed=0; \
	for dir in $(TEST_DIRS); do \
		if [ -d "$$dir" ]; then \
			echo "Testing $$dir..."; \
			(cd $$dir && go test -race -count=1 ./...) || failed=$$((failed + 1)); \
		else \
			echo "Skipping $$dir (not found)"; \
		fi; \
	done; \
	if [ $$failed -ne 0 ]; then echo "FAILED: $$failed module(s) failed"; exit 1; fi

cover:
	@for dir in $(TEST_DIRS); do \
		if [ -d "$$dir" ]; then \
			echo "Coverage $$dir..."; \
			(cd $$dir && go test -race -coverprofile=coverage.out ./... && go tool cover -func=coverage.out && rm -f coverage.out); \
		fi; \
	done

tidy:
	@for dir in $(TEST_DIRS); do \
		if [ -d "$$dir/go.mod" ]; then \
			echo "Tidying $$dir..."; \
			(cd $$dir && go mod tidy); \
		fi; \
	done

vulncheck:
	@for dir in $(TEST_DIRS); do \
		if [ -f "$$dir/go.mod" ]; then \
			echo "Vulncheck $$dir..."; \
			(cd $$dir && govulncheck ./...); \
		fi; \
	done

verify: build test lint
	@echo "All checks passed."

fmt:
	@gofmt -l -s -w .

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
