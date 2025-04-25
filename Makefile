.PHONY: test build build-run run

run:
	go run ./bin/pdfjuicer

build:
	go build -o ./bin/pdfjuicer

build-run:
	go build -o ./bin/pdfjuicer ./cmd && ./bin/pdfjuicer

test:
	go test ./... -v
