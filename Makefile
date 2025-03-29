.PHONY: all build clean deps

GO := go
GOFLAGS := -v
BIN_DIR := bin
PLUGIN_DIR := $(BIN_DIR)/plugins
SRC_DIR := src
HOST_SRC := $(wildcard $(SRC_DIR)/host/*.go)

all: build

build: deps host plugins

host: $(HOST_SRC)
	@mkdir -p $(BIN_DIR)
	$(GO) build $(GOFLAGS) -o $(BIN_DIR)/host $(SRC_DIR)/host/*.go

plugins: calculator string_utils date_utils

calculator:
	@mkdir -p $(PLUGIN_DIR)
	$(GO) build $(GOFLAGS) -o $(PLUGIN_DIR)/calculator $(SRC_DIR)/plugins/calculator/*.go

string_utils:
	@mkdir -p $(PLUGIN_DIR)
	$(GO) build $(GOFLAGS) -o $(PLUGIN_DIR)/string_utils $(SRC_DIR)/plugins/string_utils/*.go

date_utils:
	@mkdir -p $(PLUGIN_DIR)
	$(GO) build $(GOFLAGS) -o $(PLUGIN_DIR)/date_utils $(SRC_DIR)/plugins/date_utils/*.go

deps:
	$(GO) mod download
	$(GO) mod verify

clean:
	rm -rf $(BIN_DIR)
