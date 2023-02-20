.PHONY: fmt check test

all: fmt check test

fmt:
	gofmt -s -w -l .
	@echo 'goimports' && goimports -w -local gobe $(shell find . -type f -name '*.go' -not -path "./internal/*")
	gci write -s standard -s default -s "Prefix(gobe)" --skip-generated .
	go mod tidy

check:
	revive -exclude pkg/... -formatter friendly -config test/tools/revive.toml  ./...
	golangci-lint run
	go vet -all ./...
	misspell -error */**
	@echo 'staticcheck' && staticcheck $(shell go list ./... | grep -v internal)

test:
	go test ./...