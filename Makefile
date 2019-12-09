.PHONY: build
build:
	go build -v ./cmd/apiserver

.PHONY: test
test:
	go test -v -cover  -race  ./...

.DEFAULT_GOAL := build