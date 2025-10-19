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

.PHONY: all build test lint clean run coverage

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