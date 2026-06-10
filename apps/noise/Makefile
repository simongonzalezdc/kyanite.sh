# Build automation for noise.sh
SHELL := /usr/bin/env bash
GO ?= go
GOLANGCI_LINT ?= golangci-lint
BUILD_DIR := bin
DIST_DIR := dist
BINARY_NAME := noise
COVERAGE_FILE := coverage.out
EXE_SUFFIX :=
VERSION ?= dev
COMMIT ?= none
DATE ?= unknown
LD_FLAGS := -s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)
RELEASE_TOOL := $(GO) run ./tools/release/main.go

ifeq ($(OS),Windows_NT)
	EXE_SUFFIX := .exe
endif

BINARY := $(BUILD_DIR)/$(BINARY_NAME)$(EXE_SUFFIX)

.PHONY: all build test lint clean run coverage test-themes launch theme-test comprehensive-test

all: build

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

build: $(BUILD_DIR)
	$(GO) build -trimpath -ldflags "$(LD_FLAGS)" -o $(BINARY) ./cmd/noise

test:
	$(GO) test ./... -coverprofile=$(COVERAGE_FILE) -covermode=atomic

lint:
	$(GOLANGCI_LINT) run ./...

coverage: test
	$(GO) tool cover -func=$(COVERAGE_FILE)

clean:
	rm -f $(COVERAGE_FILE)
	rm -rf $(BUILD_DIR) $(DIST_DIR)

release-notes:
ifeq ($(strip $(VERSION_TAG)),)
	$(error VERSION_TAG is required (use VERSION_TAG=vX.Y.Z))
endif
ifeq ($(strip $(OUTPUT)),)
	$(error OUTPUT is required (use OUTPUT=build/release-notes.md))
endif
	$(RELEASE_TOOL) notes --from "$(FROM)" --to "$(TO)" --version "$(VERSION_TAG)" --output "$(OUTPUT)"

release-changelog:
ifeq ($(strip $(VERSION_TAG)),)
	$(error VERSION_TAG is required (use VERSION_TAG=vX.Y.Z))
endif
ifeq ($(strip $(NOTES)),)
	$(error NOTES is required (use NOTES=build/release-notes.md))
endif
	$(RELEASE_TOOL) changelog --version "$(VERSION_TAG)" --notes "$(NOTES)"

release-checksums:
ifeq ($(strip $(CHECKSUM_OUTPUT)),)
	$(error CHECKSUM_OUTPUT is required (use CHECKSUM_OUTPUT=build/checksums.txt))
endif
	$(RELEASE_TOOL) checksums --input-dir "$(DIST_DIR)" --output "$(CHECKSUM_OUTPUT)"

verify-checksums:
ifeq ($(strip $(CHECKSUM_FILE)),)
	$(error CHECKSUM_FILE is required (use CHECKSUM_FILE=build/checksums.txt))
endif
	$(RELEASE_TOOL) verify-checksums --input-dir "$(DIST_DIR)" --checksums "$(CHECKSUM_FILE)"

run: build
	./$(BINARY)

# Theme testing targets
test-themes: build
	@echo "Running theme system tests..."
	go run scripts/test_themes.go

launch: build
	@echo "Launching noise.sh..."
	./$(BINARY)

launch-debug: build
	@echo "Launching noise.sh with debug mode..."
	./$(BINARY) --debug

launch-quick: build
	@echo "Launching noise.sh in quick mode..."
	./$(BINARY) quick

theme-test: test-themes
	@echo "Theme testing completed. Launching application for manual testing..."
	./$(BINARY) --debug

comprehensive-test: build
	@echo "Running comprehensive theme testing..."
	go run tools/theme_test.go

# Windows-specific targets
launch-windows:
	@if exist scripts\build_and_launch.bat (
		scripts\build_and_launch.bat
	) else (
		echo "Windows launch script not found. Using default launch..."
		make run
	)

# Linux/Mac-specific targets
launch-unix:
	@if [ -f scripts/build_and_launch.sh ]; then \
		chmod +x scripts/build_and_launch.sh && \
		./scripts/build_and_launch.sh; \
	else \
		echo "Unix launch script not found. Using default launch..."; \
		make run; \
	fi

# =============================================================================
# Voice-to-text dependencies (OPTIONAL - for advanced users)
# =============================================================================
# NOTE: Voice models are downloaded AUTOMATICALLY on first use.
# These commands are only needed if you want to:
# - Pre-download models before first run
# - Use a different model than the default (base.en)
# - Build with native whisper.cpp support (requires CGO)

WHISPER_DIR := deps/whisper.cpp
MODELS_DIR := data/models

.PHONY: whisper-deps whisper-clean download-model init-voice

# Clone and build whisper.cpp library (ADVANCED: only needed for CGO build)
whisper-deps:
	@echo "Setting up whisper.cpp dependencies..."
	@if [ ! -d "$(WHISPER_DIR)" ]; then \
		mkdir -p deps && \
		git clone https://github.com/ggerganov/whisper.cpp $(WHISPER_DIR); \
	fi
	@echo "Building whisper.cpp library..."
	cd $(WHISPER_DIR) && make libwhisper.a
	@echo "whisper.cpp library built successfully"

# Clean whisper.cpp build artifacts
whisper-clean:
	@if [ -d "$(WHISPER_DIR)" ]; then \
		cd $(WHISPER_DIR) && make clean; \
	fi

# Pre-download whisper model (OPTIONAL - downloads automatically on first use)
download-model:
	@echo "Pre-downloading whisper model..."
	@mkdir -p $(MODELS_DIR)
	@if [ ! -f "$(MODELS_DIR)/ggml-base.en.bin" ]; then \
		curl -L -o $(MODELS_DIR)/ggml-base.en.bin \
			https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-base.en.bin; \
		echo "Model downloaded to $(MODELS_DIR)/ggml-base.en.bin"; \
	else \
		echo "Model already exists at $(MODELS_DIR)/ggml-base.en.bin"; \
	fi

# Download tiny model (faster, less accurate)
download-model-tiny:
	@mkdir -p $(MODELS_DIR)
	@if [ ! -f "$(MODELS_DIR)/ggml-tiny.en.bin" ]; then \
		curl -L -o $(MODELS_DIR)/ggml-tiny.en.bin \
			https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-tiny.en.bin; \
	fi

# Download small model (better accuracy, slower)
download-model-small:
	@mkdir -p $(MODELS_DIR)
	@if [ ! -f "$(MODELS_DIR)/ggml-small.en.bin" ]; then \
		curl -L -o $(MODELS_DIR)/ggml-small.en.bin \
			https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-small.en.bin; \
	fi

# Initialize voice-to-text (ADVANCED: for CGO builds)
init-voice: whisper-deps download-model
	@echo "Voice-to-text initialization complete"

# =============================================================================
# PWA Sync directories (OPTIONAL - created automatically)
# =============================================================================
# NOTE: Sync directories are created AUTOMATICALLY when the app starts.
# This command is only for manual setup if needed.

.PHONY: init-sync

init-sync:
	@echo "Creating sync directories..."
	@mkdir -p data/sync/media/voice
	@mkdir -p data/sync/media/photos
	@echo "Sync directories created"

# =============================================================================
# Build variants
# =============================================================================

# Standard build (voice support via Go bindings, no CGO required)
# Models download automatically on first use
build: $(BUILD_DIR)
	$(GO) build -trimpath -ldflags "$(LD_FLAGS)" -o $(BINARY) ./cmd/noise

# Build with native whisper.cpp (ADVANCED: requires whisper-deps first)
build-native-voice: $(BUILD_DIR)
	CGO_ENABLED=1 $(GO) build -tags whisper -trimpath -ldflags "$(LD_FLAGS)" -o $(BINARY) ./cmd/noise