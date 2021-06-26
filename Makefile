.PHONY: clean build

clean:
	rm -rf bin

build:
	go build -ldflags "-s -w" -o bin/apiserver ./cmd/apiserver

.DEFAULT_GOAL := build
