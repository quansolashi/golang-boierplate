.PHONY: help setup install gqlgen entdesc entgen protoc fmt vet lint lint-fix test docs build dev 

GOLANGCI_VERSION := 1.57.2
MOCKGEN_VERSION := 0.4.0

LINT_PACKAGES := $(shell go list $(CURDIR)/... | grep -v -e "mock" -v -e "tmp")
TEST_PACKAGES := $(shell go list $(CURDIR)/internal/... $(CURDIR)/pkg/...)

setup: install
	wget -O - -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v${GOLANGCI_VERSION}

install:
	go install go.uber.org/mock/mockgen@v${MOCKGEN_VERSION}

gqlgen:
	go run github.com/99designs/gqlgen generate

entdesc:
	go run -mod=mod entgo.io/ent/cmd/ent describe ./ent/schema

entgen:
	go generate ./ent

fmt:
	! gofmt -d -s ./pkg | grep '^'

vet:
	go vet $(LINT_PACKAGES)

lint:
	./bin/golangci-lint run -c .golangci.yaml ./...

lint-fix:
	./bin/golangci-lint run -c .golangci.yaml ./... --fix

test:
	go test -v -cover -coverprofile=coverage.txt -covermode=atomic $(TEST_PACKAGES)

docs:
	swag init -d ./internal/web/controller -g ./controller.go --pd -o ./docs/web

build:
	go build -o ./app ./cmd/main.go

dev:
	air -c ./.air.toml
