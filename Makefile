.PHONY: build
build:
	GOOS=linux go build -v ./cmd/apiserver

.PHONY: test
test:
	go test -v -cover  -race  ./...

clean:
	rm -rf ./apiserver

.DEFAULT_GOAL := build