# Build automation for noise.sh
SHELL := /usr/bin/env bash
GO ?= go
GOLANGCI_LINT ?= golangci-lint
BUILD_DIR := bin
BINARY_NAME := noise
COVERAGE_FILE := coverage.out
EXE_SUFFIX :=

ifeq ($(OS),Windows_NT)
	EXE_SUFFIX := .exe
endif

BINARY := $(BUILD_DIR)/$(BINARY_NAME)$(EXE_SUFFIX)

.PHONY: all build test lint clean run coverage

all: build

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

build: $(BUILD_DIR)
	$(GO) build -o $(BINARY) ./cmd/noise

test:
	$(GO) test ./... -coverprofile=$(COVERAGE_FILE) -covermode=atomic

lint:
	$(GOLANGCI_LINT) run ./...

coverage: test
	$(GO) tool cover -func=$(COVERAGE_FILE)

clean:
	rm -f $(COVERAGE_FILE)
	rm -rf $(BUILD_DIR)

run: build
	./$(BINARY)