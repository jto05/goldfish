BINARY := goldfish
BUILD_DIR := build

build:
	go build -o $(BUILD_DIR)/$(BINARY) ./cmd/goldfish

run:
	go run ./cmd/goldfish

test:
	go test ./...

clean:
	rm -rf $(BUILD_DIR)
