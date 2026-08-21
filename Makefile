BINARY := goldfish
BUILD_DIR := build

build:
	go build -o $(BUILD_DIR)/$(BINARY) ./cmd/goldfish

run:
	go run ./cmd/goldfish

run-react:
	npm run dev --prefix ./web -- --host

test:
	go test ./...

clean:
	rm -rf $(BUILD_DIR)
