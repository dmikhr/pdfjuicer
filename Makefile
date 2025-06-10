LOCAL_BIN:=$(CURDIR)/bin

.PHONY: third_party_licenses go-lic lint install-golangci-lint test build build-run run

run:
	./bin/pdfjuicer

build:
	go build -o ./bin/pdfjuicer

build-run:
	go build -o ./bin/pdfjuicer ./cmd && ./bin/pdfjuicer

test:
	go test ./... -v

install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8

lint:
	$(LOCAL_BIN)/golangci-lint run --config=.golangci.yaml ./...

validate-lint-config:
	$(LOCAL_BIN)/golangci-lint config verify --config=.golangci.yaml

go-lic-install:
	GOBIN=$(LOCAL_BIN) go install github.com/google/go-licenses@v1.6.0

third_party_licenses:
	$(LOCAL_BIN)/go-licenses report ./... > THIRD_PARTY_LICENSES
