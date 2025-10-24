BUILD_DIR = build

all: examples
	go build -o build/ts-plug ./cmd/ts-plug-web

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