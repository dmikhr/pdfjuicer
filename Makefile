LOCAL_BIN:=$(CURDIR)/bin

.PHONY: lint install-golangci-lint test build build-run run

run:
	go run ./bin/pdfjuicer

build:
	go build -o ./bin/pdfjuicer

build-run:
	go build -o ./bin/pdfjuicer ./cmd && ./bin/pdfjuicer

test:
	go test ./... -v

install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8

lint:
	@$(LOCAL_BIN)/golangci-lint run ./... --config .golangci.yaml
