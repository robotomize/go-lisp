APP_NAME := go-lisp-cli
CMD := ./cmd/go-lisp-cli
BIN_DIR := bin

export GOBIN := $(PWD)/$(BIN_DIR)
export PATH := $(GOBIN):$(PATH)

LDFLAGS := -s -w

GOWRAP_VERSION := 1.4.0
GOLANG_CI_LINT_VERSION ?= v2.12.2
GO_VERSION := $(shell grep '^go ' go.mod | awk '{print $$2}')

UNAME_S := $(shell uname -s)
UNAME_P := $(shell uname -p)

ifeq ($(UNAME_S),Linux)
	OSFLAG = linux
endif

ifeq ($(UNAME_S),Darwin)
	OSFLAG = osx
  ifeq ($(UNAME_P),arm)
    # protobuf team doesn't create releases
    # for Apple M1 (arm) chipset
	  UNAME_P = x86_64
  endif
  ifeq ($(UNAME_P),i386)
    # for Rosetta 2 emulator
	  UNAME_P = x86_64
  endif
endif

.PHONY: all
all: test build

./bin:
	mkdir -p ./bin

./bin/gomock: | ./bin
	go install go.uber.org/mock/mockgen@latest

.PHONY: ./bin/golangci-lint
./bin/golangci-lint: | ./bin
	GOTOOLCHAIN=go$(GO_VERSION) GOBIN=$(PWD)/bin go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANG_CI_LINT_VERSION)

.PHONY: lint
lint: ./bin/golangci-lint
	./bin/golangci-lint run --timeout 5m ./...

.PHONY: fix-lint
fix-lint: ./bin/golangci-lint
	./bin/golangci-lint run --fix --timeout=5m ./...

.PHONY: test
test:
	GOGC=off go test -short ./..

.PHONY: build-cli
build: build-cli

.PHONY: build-cli
build-cli: | $(BIN_DIR)
	go build -trimpath \
		-ldflags "$(LDFLAGS)" \
		-o $(BIN_DIR)/$(APP_NAME) \
		$(CMD)

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: clean
clean:
	rm -rf $(BIN_DIR)

$(BIN_DIR):
	mkdir -p $(BIN_DIR)
