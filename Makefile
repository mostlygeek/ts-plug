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
	GOOS=darwin GOARCH=arm64 go build -o build/ts-multi-plug-darwin-arm64 ./cmd/ts-multi-plug

linux:
	GOOS=linux GOARCH=arm64 go build -o build/ts-plug-linux-arm64 ./cmd/ts-plug-web
	GOOS=linux GOARCH=amd64 go build -o build/ts-plug-linux-amd64 ./cmd/ts-plug-web
	GOOS=linux GOARCH=arm64 go build -o build/ts-multi-plug-darwin-amd64 ./cmd/ts-multi-plug
	GOOS=linux GOARCH=arm64 go build -o build/ts-multi-plug-darwin-arm64 ./cmd/ts-multi-plug

clean:
	rm -rf $(BUILD_DIR)/*

examples:
	go build -o $(BUILD_DIR)/hello ./cmd/examples/hello/hello.go
	go build -o $(BUILD_DIR)/resolver ./cmd/examples/resolver/resolver.go

# use cached test results while developing
test: examples
#	go test -race -timeout 30s -short ./internal/...
	staticcheck ./... || true

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

.PHONY: all test examples clean