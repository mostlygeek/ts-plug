BUILD_DIR = build

# Get the current Git hash
GIT_HASH := $(shell git rev-parse --short HEAD)
ifneq ($(shell git status --porcelain),)
    # There are untracked changes
    GIT_HASH := $(GIT_HASH)+
endif

# Capture the current build date in RFC3339 format
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")


all: examples darwin linux

darwin:
	GOOS=darwin GOARCH=arm64 go build -o build/ts-plug-darwin-arm64 ./cmd/ts-plug-web

linux:
	GOOS=linux GOARCH=arm64 go build -o build/ts-plug-linux-arm64 ./cmd/ts-plug-web
	GOOS=linux GOARCH=amd64 go build -o build/ts-plug-linux-amd64 ./cmd/ts-plug-web

clean:
	rm -rf $(BUILD_DIR)/*

examples:
	go build -o $(BUILD_DIR)/hello ./cmd/examples/hello-go

# use cached test results while developing
test: examples
#	go test -race -timeout 30s -short ./internal/...
	staticcheck ./... || true

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

.PHONY: all test examples clean