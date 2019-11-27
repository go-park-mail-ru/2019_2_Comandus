.PHONY: build
build:
	go build -v ./cmd/apiserver

.PHONY: test
test:
	go test -v -cover -coverpkg=./... -race  ./...

.DEFAULT_GOAL := build