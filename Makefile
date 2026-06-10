.PHONY: build launcher test lint clean release snapshot sync

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

test:
	@for dir in pkg/design pkg/config pkg/ai apps/focus apps/noise apps/syntax apps/prism cmd/kyanite; do \
		echo "Testing $$dir..."; \
		cd $$dir && go test ./... && cd ../..; \
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
